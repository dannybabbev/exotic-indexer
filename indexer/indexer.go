package indexer

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/bitgemtech/ord-api/ordinals"
	badger "github.com/dgraph-io/badger/v4"
	log "github.com/sirupsen/logrus"
)

type SyncStats struct {
	ChainTip       int    `json:"chainTip"`
	SyncHeight     int    `json:"syncHeight"`
	SyncBlockHash  string `json:"syncBlockHash"`
	ReorgsDetected []int  `json:"reorgsDetected"`
}

type Bitcoind interface {
	GetBlockHash(height uint64) (string, error)
	GetRawBlock(hash string) (string, error)
	GetBlockCount() (uint64, error)
}

type Indexer struct {
	bitcoind         Bitcoind
	db               *badger.DB
	stats            *SyncStats
	utxoIndex        *ordinals.UTXOIndex
	delUTXOs         [][]string
	utxoHistory      [][]string
	blocksChan       chan *ordinals.Block
	periodFlushToDB  int
	keepBlockHistory int
}

const BLOCK_PREFETCH = 12

func NewIndexer(bitcoind Bitcoind, db *badger.DB) *Indexer {
	return &Indexer{
		bitcoind:    bitcoind,
		db:          db,
		stats:       &SyncStats{},
		utxoIndex:   ordinals.NewUTXOIndex(),
		delUTXOs:    make([][]string, 0),
		utxoHistory: make([][]string, 0),
		// buffered channel to allow for some prefetching
		blocksChan:       make(chan *ordinals.Block, BLOCK_PREFETCH),
		periodFlushToDB:  500,
		keepBlockHistory: 6,
	}
}

func (b *Indexer) WithKeepBlockHistory(value int) *Indexer {
	b.keepBlockHistory = value
	return b
}

func (b *Indexer) WithPeriodFlushToDB(value int) *Indexer {
	b.periodFlushToDB = value
	return b
}

func (b *Indexer) StartDaemon(stopChan chan struct{}) {
	// Create a ticker that ticks every n seconds
	n := 10
	ticker := time.NewTicker(time.Duration(n) * time.Second)

	stopIndexerChan := make(chan struct{})
	indexerStoppedChan := make(chan struct{})

	isRunning := false
	tick := func() {
		if !isRunning {
			isRunning = true
			go func() {
				b.SyncToChainTip(stopIndexerChan, indexerStoppedChan)
				isRunning = false
			}()
		}
	}

	tick()
	for {
		select {
		case <-ticker.C:
			tick()
		case <-stopChan:
			if isRunning {
				stopIndexerChan <- struct{}{}
				<-indexerStoppedChan
			}
			return
		}
	}
}

// setDB encodes the data using gob and sets it using the provided badger txn
func (b *Indexer) setDB(key string, data interface{}, wb *badger.WriteBatch) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return err
	}

	bytes := buf.Bytes()

	err := wb.Set([]byte(key), []byte(bytes))
	if err != nil {
		return err
	}

	return nil
}

func (b *Indexer) getUTXODBKey(utxo string) string {
	return fmt.Sprintf("u-%s", utxo)
}

// updateDB updates the database with the current state of the UTXO set in a single transaction.
// It removes the UTXOs from the index as they are processed
func (b *Indexer) updateDB() {
	wb := b.db.NewWriteBatch()
	defer wb.Cancel()

	err := b.setDB("syncStats", b.stats, wb)
	if err != nil {
		log.Panicf("Error setting in db %v", err)
	}

	err = b.setDB("utxoHistory", b.utxoHistory, wb)
	if err != nil {
		log.Panicf("Error setting in db %v", err)
	}

	deleted := 0
	for _, block := range b.delUTXOs {
		for _, v := range block {
			deleted += 1
			key := b.getUTXODBKey(v)
			err := wb.Delete([]byte(key))
			if err != nil {
				log.Infof("Error deleting db: %v\n", err)
			}
		}
	}

	log.Infof("flushing utxos added %d and deleted %d\n", len(b.utxoIndex.Index), deleted)

	// Add the new utxos
	for k, v := range b.utxoIndex.Index {
		key := b.getUTXODBKey(k)

		// save just the ordinals
		saveUTXO := &ordinals.Output{
			Ordinals: v.Ordinals,
		}

		err := b.setDB(key, saveUTXO, wb)
		if err != nil {
			log.Panicf("Error setting in db %v", err)
		}

		// empty the memory as it is being flushed
		delete(b.utxoIndex.Index, k)
	}

	err = wb.Flush()
	if err != nil {
		log.Panicf("Error flushing writes to db %v", err)
	}

	// reset the in memory indexes
	b.utxoIndex = ordinals.NewUTXOIndex()
	b.delUTXOs = make([][]string, 0)
}

func (b *Indexer) ForceMajeure(indexerStoppedChan chan struct{}) {
	log.Info("Graceful shutdown received, flushing db...")
	b.updateDB()

	indexerStoppedChan <- struct{}{}
}

func (b *Indexer) handleReorg(currentBlock *ordinals.Block, heightTarget int, stopChan chan struct{}, indexerStoppedChan chan struct{}) {
	// rewind the sync height by the available history
	rewind := len(b.utxoHistory)
	if rewind == 0 {
		log.Panicf("reorg detected at heigh %d but no utxo history was cached", currentBlock.Height)
	}

	b.stats.ReorgsDetected = append(b.stats.ReorgsDetected, currentBlock.Height)

	// rewind the sync height by the available history
	newBlockHeigh := currentBlock.Height - len(b.utxoHistory)
	newBlock := b.fetchBlock(newBlockHeigh)
	b.stats.SyncHeight = newBlockHeigh
	b.stats.SyncBlockHash = newBlock.Hash

	// It is important that we delete the UTXO history to keep the state consitent
	b.utxoHistory = make([][]string, 0)

	b.drainBlocksChan()
	b.updateDB()
	b.SyncToBlock(heightTarget, stopChan, indexerStoppedChan)
}

// SyncToBlock continues from the sync height to the current height
func (b *Indexer) SyncToBlock(height int, stopChan chan struct{}, indexerStoppedChan chan struct{}) {
	b.LoadSyncStatsFromDB()

	if b.stats.SyncHeight == height {
		log.Tracef("already synced to block %d\n", height)
		return
	}

	log.WithFields(log.Fields{
		"currentHeigh": b.stats.SyncHeight,
		"height":       height,
	}).Info("starting sync")

	// if we don't start from precisely this heigh the UTXO index is worthless
	// we need to start from exactly where we left off
	start := b.stats.SyncHeight + 1

	periodProcessedTxs := 0
	startTime := time.Now() // Record the start time

	logProgressPeriod := 1

	stopBlockFetcherChan := make(chan struct{})
	go b.spawnBlockFetcher(start, height, stopBlockFetcherChan)

	for i := start; i <= height; i++ {
		select {
		case <-stopChan:
			b.ForceMajeure(indexerStoppedChan)
			return
		default:
			block := <-b.blocksChan

			// make sure that we are at the correct block height
			if block.Height != i {
				log.Panicf("Expected block height %d, got %d", i, block.Height)
			}

			// detect reorgs
			if i > 0 && block.PrevBlockHash != b.stats.SyncBlockHash {
				log.WithField("height", i).Warn("reorg detected")
				stopBlockFetcherChan <- struct{}{}
				b.handleReorg(block, height, stopChan, indexerStoppedChan)
				return
			}

			periodProcessedTxs += len(block.Transactions)

			b.prefetchIndexesFromDB(block, b.utxoIndex)

			// assign the ordinals and get an array of UTXOs that were spent
			newDelUTXOs := ordinals.AssignOrdinals(block, b.utxoIndex)

			// if we are in the last n blocks, we do not want to directly delete the
			// spent UTXOs, instead we will keep them in the database for n blocks
			// this will allow us to rollback in case of a reorg
			if height-i < b.keepBlockHistory {
				// the utxo history is what we keep
				b.utxoHistory = append(b.utxoHistory, newDelUTXOs)
			} else {
				// the delUTXOs is what we delete
				b.delUTXOs = append(b.delUTXOs, newDelUTXOs)
			}

			// if we have more than n blocks in the history, we need to remove the oldest
			if len(b.utxoHistory) > b.keepBlockHistory {
				b.delUTXOs = append(b.delUTXOs, b.utxoHistory[0])
				b.utxoHistory = b.utxoHistory[1:]
			}

			// Update the sync stats
			b.stats.ChainTip = height
			b.stats.SyncHeight = i
			b.stats.SyncBlockHash = block.Hash

			if block.Height%b.periodFlushToDB == 0 {
				b.updateDB()
			}

			if i%logProgressPeriod == 0 {
				elapsedTime := time.Since(startTime)
				timePerTx := elapsedTime / time.Duration(periodProcessedTxs)
				readableTime := block.Timestamp.Format("2006-01-02 15:04:05")
				log.Infof("processed block %d (%s) with %d transactions took %v (%v per tx)\n", block.Height, readableTime, periodProcessedTxs, elapsedTime, timePerTx)
				startTime = time.Now()
				periodProcessedTxs = 0
			}
		}
	}

	b.updateDB()
	log.Infof("synced to block %d\n", height)
}

func (b *Indexer) SyncToChainTip(stopChan chan struct{}, indexerStoppedChan chan struct{}) {
	count, err := b.bitcoind.GetBlockCount()
	if err != nil {
		log.Panicf("failed to get block count %v", err)
	}

	b.SyncToBlock(int(count), stopChan, indexerStoppedChan)
}

func (b *Indexer) decodeBytes(data []byte, target interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(target)
}

func (b *Indexer) getValueFromDB(key string, txn *badger.Txn, target interface{}) error {
	item, err := txn.Get([]byte(key))
	if err != nil {
		return err
	}

	err = item.Value(func(v []byte) error {
		err := b.decodeBytes(v, target)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (b *Indexer) prefetchIndexesFromDB(block *ordinals.Block, utxoIndex *ordinals.UTXOIndex) {
	utxosToFetch := make([]string, 0)

	for _, tx := range block.Transactions {
		for _, input := range tx.Inputs {
			utxo := fmt.Sprintf("%s:%d", input.Txid, input.Vout)

			// check if we already have this utxo in the index
			if _, ok := utxoIndex.Index[utxo]; !ok {
				utxosToFetch = append(utxosToFetch, utxo)
			}
		}
	}

	startTime := time.Now()
	err := b.db.View(func(txn *badger.Txn) error {
		for _, u := range utxosToFetch {
			utxo := &ordinals.Output{}
			dbKey := b.getUTXODBKey(u)
			err := b.getValueFromDB(dbKey, txn, utxo)
			if err == badger.ErrKeyNotFound {
				continue
			} else if err != nil {
				return err
			}

			utxoIndex.Index[u] = utxo
		}

		return nil
	})

	if err != nil {
		log.Panicf("Error prefetching utxos from db: %v", err)
	}

	elapsed := time.Since(startTime)
	log.Tracef("prefetched %d utxos in %v\n", len(utxosToFetch), elapsed)
}

func (b *Indexer) GetOrdinalsForUTXO(utxo string) ([]*ordinals.Range, error) {
	ranges := []*ordinals.Range{}
	err := b.db.View(func(txn *badger.Txn) error {
		output := &ordinals.Output{}
		key := b.getUTXODBKey(utxo)
		err := b.getValueFromDB(key, txn, output)
		if err != nil {
			return err
		}

		ranges = output.Ordinals

		return nil
	})

	return ranges, err
}

func (b *Indexer) LoadSyncStatsFromDB() {
	err := b.db.View(func(txn *badger.Txn) error {
		syncStats := &SyncStats{}
		err := b.getValueFromDB("syncStats", txn, syncStats)
		if err == badger.ErrKeyNotFound {
			log.Info("No sync stats found in db, setting height to -1")
			syncStats.SyncHeight = -1
		} else if err != nil {
			return err
		}

		if syncStats.ReorgsDetected == nil {
			syncStats.ReorgsDetected = make([]int, 0)
		}

		b.stats = syncStats

		utxoHistory := [][]string{}
		err = b.getValueFromDB("utxoHistory", txn, &utxoHistory)
		if err == badger.ErrKeyNotFound {
			log.Info("No utxo history found in db")
		} else if err != nil {
			return err
		}

		b.utxoHistory = utxoHistory

		return nil
	})

	if err != nil {
		log.Panicf("Error loading sync stats from db: %v", err)
	}
}

// TriggerReorg is meant to be used for debugging and tests only
// I used it to simulate a reorg
func (b *Indexer) TriggerReorg() {
	b.LoadSyncStatsFromDB()
	b.stats.SyncBlockHash = "wrong"
	b.updateDB()
}
