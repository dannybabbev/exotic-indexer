package indexer

import (
	"encoding/hex"

	"github.com/bitgemtech/ord-api/ordinals"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	log "github.com/sirupsen/logrus"
)

func (b *Indexer) fetchBlock(height int) *ordinals.Block {
	hash, err := b.bitcoind.GetBlockHash(uint64(height))
	if err != nil {
		log.Fatalln(err)
	}

	rawBlock, err := b.bitcoind.GetRawBlock(hash)
	if err != nil {
		log.Fatalln(err)
	}

	blockData, err := hex.DecodeString(rawBlock)
	if err != nil {
		log.Panicf("Failed to decode block: %v", err)
	}

	// Deserialize the bytes into a btcutil.Block.
	block, err := btcutil.NewBlockFromBytes(blockData)
	if err != nil {
		log.Panicf("Failed to parse block: %v", err)
	}

	transactions := block.Transactions()
	txs := make([]*ordinals.Transaction, len(transactions))
	for i, tx := range transactions {
		inputs := []*ordinals.Input{}
		outputs := []*ordinals.Output{}

		for _, v := range tx.MsgTx().TxIn {
			txid := v.PreviousOutPoint.Hash.String()
			vout := v.PreviousOutPoint.Index
			input := &ordinals.Input{Txid: txid, Vout: int64(vout)}
			inputs = append(inputs, input)
		}

		// parse the raw tx values
		for j, v := range tx.MsgTx().TxOut {
			// Determine the type of the script and extract the address
			scyptClass, addrs, _, err := txscript.ExtractPkScriptAddrs(v.PkScript, &chaincfg.MainNetParams)
			if err != nil {
				log.Panicf("Failed to extract address: %v", err)
			}

			addrsString := make([]string, len(addrs))
			for i, x := range addrs {
				addrsString[i] = x.String()
			}

			var receiver ordinals.ScriptPubKey

			if len(addrs) == 0 {
				receiver = ordinals.ScriptPubKey{
					Addresses: []string{"UNKNOWN"},
					Type:      "UNKNOWN",
				}
			} else {
				receiver = ordinals.ScriptPubKey{
					Addresses: addrsString,
					Type:      scyptClass.String(),
				}
			}

			output := &ordinals.Output{Value: v.Value, Address: receiver, N: int64(j)}
			outputs = append(outputs, output)
		}

		txs[i] = &ordinals.Transaction{
			Txid:    tx.Hash().String(),
			Inputs:  inputs,
			Outputs: outputs,
		}
	}

	t := block.MsgBlock().Header.Timestamp
	bl := &ordinals.Block{
		Timestamp:     t,
		Height:        height,
		Hash:          block.Hash().String(),
		PrevBlockHash: block.MsgBlock().Header.PrevBlock.String(),
		Transactions:  txs,
	}

	return bl
}

// Prefetches blocks from bitcoind and sends them to the blocksChan
func (b *Indexer) spawnBlockFetcher(startHeigh int, endHeight int, stopChan chan struct{}) {
	currentHeight := startHeigh
	for currentHeight <= endHeight {
		select {
		case <-stopChan:
			return
		default:
			block := b.fetchBlock(currentHeight)
			b.blocksChan <- block
			currentHeight += 1
		}
	}

	<-stopChan
}

func (b *Indexer) drainBlocksChan() {
	for {
		select {
		case <-b.blocksChan:
		default:
			return
		}
	}
}
