# 🌴 Exotic Sat Indexer 

A fast ordinals and exotic sats indexer written in Go. The indexer exposes a REST API which you can use to build sat hunting tools, ordinals wallets and other ordinals services.

## API Docs
- Docs: https://docs.bitgem.tech/
- Public api root: https://api.bitgem.tech
- [Official Website](https://bitgem.tech/)

## Requirements
- Bitcoind node
- [Esplora](https://github.com/Blockstream/esplora) backend (optional, only if you need address index)

### Minimum System Requirements:
- 32 GB RAM
- 8 Core CPU
- ~250 GB SSD for the indexer db (Dec 2023)

## Run with docker
1. Make sure that bitcoind is running and synced. Bitcoind should be available at `http://bitcoind:8332` or `http://localhost:8332`
2. Set up directory for the indexer db
3. Run
```bash
docker run -d --name exotic-indexer \
    -v ./index-dir:/db \
    -v ./bitcoind-root-dir:/bitcoin \
    --stop-timeout 900 \
    -e BITCOIN_RPC_HOST=bitcoin \
    -e BITCOIN_RPC_PORT=8332 \
    -e DATA_DIR=/db \
    -e BITCOIND_DIR=/bitcoin \
    -e ESPLORA_URL=https://blockstream.info/api/ \
    lebonchasseur/exotic-indexer:master
```

If you want to run for testnet there is no difference, just point to the testnet bitcoind directory
```bash
docker run -d --name exotic-indexer \
    -v ./index-dir:/db \
    -v ./bitcoind-root-dir:/bitcoin \
    --stop-timeout 900 \
    -e BITCOIN_RPC_HOST=bitcoin \
    -e BITCOIN_RPC_PORT=8332 \
    -e DATA_DIR=/db \
    -e BITCOIND_DIR=/bitcoin/testnet3 \
    -e ESPLORA_URL=https://blockstream.info/testnet/api/ \
    lebonchasseur/exotic-indexer:master
```
### Build your own docker image
You can also build your own docker image `docker build -t exotic-indexer .` and run it with the same command as above.

## Develop
1. Rename `env` to `.env` and substitute the values for your system
2. Run `docker compose up` in the root directory to bring up the testnet instance of bitcoind
   - the `bitcoind` container is configured to use `arm64v8` by default for Apple sillicon support. Adjust this for your system.
3. `go get -v`
4. Run and debug with VS Code or the editor of your choice

## Community and support
- [Discord](https://discord.gg/STgzjMnkhT)
- [Official Website](https://bitgem.tech/)
