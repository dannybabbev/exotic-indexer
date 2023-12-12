package ordinals

import "time"

type Range struct {
	Start int64 `json:"start"`
	Size  int64 `json:"size"`
}

type Input struct {
	Txid     string   `json:"txid"`
	Vout     int64    `json:"vout"`
	Ordinals []*Range `json:"ordinals"`
}

type ScriptPubKey struct {
	Addresses []string `json:"addresses"`
	Type      string   `json:"type"`
}

type Output struct {
	Value    int64        `json:"value"`
	Address  ScriptPubKey `json:"scriptPubKey"`
	N        int64        `json:"n"`
	Ordinals []*Range     `json:"ordinals"`
}

type Transaction struct {
	Txid    string    `json:"txid"`
	Inputs  []*Input  `json:"inputs"`
	Outputs []*Output `json:"outputs"`
}

type Block struct {
	Timestamp     time.Time      `json:"timestamp"`
	Height        int            `json:"height"`
	Hash          string         `json:"hash"`
	PrevBlockHash string         `json:"prevBlockHash"`
	Transactions  []*Transaction `json:"transactions"`
}

type UTXOIndex struct {
	Index map[string]*Output
}

func NewUTXOIndex() *UTXOIndex {
	return &UTXOIndex{Index: make(map[string]*Output)}
}
