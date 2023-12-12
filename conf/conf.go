package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Conf struct {
	BitcoinRPCHost  string
	BitcoinRPCUser  string
	BitcoinRPCPass  string
	BitcoinRPCPort  int
	DataDir         string
	EsploraURL      string
	LogLevel        log.Level
	PeriodFlushToDB int
}

// GetEnvOrPanic fetches the value of the environment variable named by the key.
// It panics if the key is unset or empty.
func getEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Environment variable %s not set", key))
	}
	return value
}

func NewConf() *Conf {
	bitcoinRPCPort, err := strconv.Atoi(getEnvOrPanic("BITCOIN_RPC_PORT"))
	if err != nil {
		log.Fatalln("Error converting BITCOIN_RPC_PORT to int")
	}

	bitcoindDir := getEnvOrPanic("BITCOIND_DIR")
	cookieFilePath := filepath.Join(bitcoindDir, ".cookie")
	data, err := os.ReadFile(cookieFilePath)
	if err != nil {
		log.Fatalln("Error reading .cookie file:", err)
	}

	// Convert data to a string and split it to get username and password
	contents := string(data)
	parts := strings.SplitN(contents, ":", 2)
	if len(parts) != 2 {
		log.Fatalln("Error parsing .cookie file")
	}

	username := parts[0]
	password := strings.TrimSpace(parts[1])

	ll := os.Getenv("LOG_LEVEL")
	logLevel, err := log.ParseLevel(ll)
	if err != nil {
		logLevel = log.InfoLevel
	}

	periodFlushToDB := 0
	periodFlush := os.Getenv("PERIOD_FLUSH_TO_DB")
	if periodFlush != "" {
		periodFlushToDB, err = strconv.Atoi(periodFlush)
		if err != nil {
			log.Fatalln("Error converting PERIOD_FLUSH_TO_DB to int")
		}
	}

	return &Conf{
		BitcoinRPCUser:  username,
		BitcoinRPCPass:  password,
		BitcoinRPCPort:  bitcoinRPCPort,
		BitcoinRPCHost:  getEnvOrPanic("BITCOIN_RPC_HOST"),
		DataDir:         getEnvOrPanic("DATA_DIR"),
		EsploraURL:      getEnvOrPanic("ESPLORA_URL"),
		LogLevel:        logLevel,
		PeriodFlushToDB: periodFlushToDB,
	}
}
