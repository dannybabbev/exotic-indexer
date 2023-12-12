package esplora

type UTXOStatus struct {
	Confirmed   bool   `json:"confirmed"`
	BlockHeight int    `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	BlockTime   int    `json:"block_time"`
}

type UTXO struct {
	Txid   string     `json:"txid"`
	Vout   int        `json:"vout"`
	Status UTXOStatus `json:"status"`
	Value  int64      `json:"value"`
}
