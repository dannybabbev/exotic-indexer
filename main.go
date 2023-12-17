package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bitgemtech/exotic-indexer/conf"
	"github.com/bitgemtech/exotic-indexer/esplora"
	"github.com/bitgemtech/exotic-indexer/indexer"
	"github.com/bitgemtech/exotic-indexer/server"
	badger "github.com/dgraph-io/badger/v4"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/toorop/go-bitcoind"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	err := godotenv.Load(".env")
	if err != nil {
		log.WithField("err", err).Trace("Error loading .env file")
	}

	conf := conf.NewConf()

	log.SetLevel(conf.LogLevel)
	log.Info("Starting...")

	log.WithFields(log.Fields{
		"BitcoinRPCHost": conf.BitcoinRPCHost,
		"BitcoinRPCPort": conf.BitcoinRPCPort,
		"EsploraURL":     conf.EsploraURL,
	}).Info("using configuration")

	bc, err := bitcoind.New(
		conf.BitcoinRPCHost,
		conf.BitcoinRPCPort,
		conf.BitcoinRPCUser,
		conf.BitcoinRPCPass,
		false,
	)
	if err != nil {
		log.Fatalln(err)
	}

	dbDir := conf.DataDir
	if len(os.Args) > 1 {
		dbDir = os.Args[1]
	}

	log.WithField("dbDir", dbDir).Info("using database directory")

	opts := badger.DefaultOptions(dbDir).
		WithLoggingLevel(badger.WARNING).
		WithBlockCacheSize(2000 << 20)

	db, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	// Listen for SIGINT (CTRL+C) signal to gracefully shut down.
	sigInt := make(chan os.Signal, 1)
	signal.Notify(sigInt, os.Interrupt, syscall.SIGTERM)

	// Create a channel to communicate when the indexer should stop
	stopChan := make(chan struct{}, 1)

	go func() {
		<-sigInt
		log.Info("Received SIGINT (CTRL+C).")
		stopChan <- struct{}{}
	}()

	indexer := indexer.NewIndexer(bc, db)
	if conf.PeriodFlushToDB != 0 {
		log.WithField("periodFlushToDB", conf.PeriodFlushToDB).Info("using periodFlushToDB from conf")
		indexer = indexer.WithPeriodFlushToDB(conf.PeriodFlushToDB)
	}

	esplora := esplora.NewEsploraAPI(conf.EsploraURL)
	serverModel := server.NewServerModel(indexer, esplora)
	server := server.NewServer(serverModel)

	go server.Start()

	indexer.StartDaemon(stopChan)

	log.Info("Shutting down...")
}
