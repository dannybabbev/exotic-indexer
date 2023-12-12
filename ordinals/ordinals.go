package ordinals

import (
	"fmt"
	"log"
	"math"
)

func Subsidy(height int) int64 {
	epoch := int64(math.Floor(float64(height) / 210000))
	return 50 * 100000000 >> epoch
}

func FirstOrdinal(height int) int64 {
	start := int64(0)
	for i := 0; i < height; i++ {
		start += Subsidy(i)
	}
	return start
}

func TransferRanges(ordinals []*Range, value int64) ([]*Range, []*Range) {
	remainingValue := value
	remaining := ordinals
	transferred := make([]*Range, 0)
	for remainingValue > 0 {
		currentRange := remaining[0]
		start := currentRange.Start
		transferSize := currentRange.Size
		if transferSize > remainingValue {
			transferSize = remainingValue
		}

		transferred = append(transferred, &Range{start, transferSize})
		remainingSize := currentRange.Size - transferSize

		if remainingSize == 0 {
			remaining = remaining[1:]
		} else {
			remaining[0] = &Range{start + transferSize, remainingSize}
		}

		remainingValue = remainingValue - transferSize
	}

	return transferred, remaining
}

func AssignOrdinals(block *Block, utxoIndex *UTXOIndex) []string {
	first := FirstOrdinal(block.Height)
	size := Subsidy(block.Height)
	coinbaseOrdinals := []*Range{{first, size}}
	deletedUTXO := make([]string, 0)

	for _, tx := range block.Transactions[1:] {
		ordinals := make([]*Range, 0)
		for _, input := range tx.Inputs {
			// the utxo to be spent in the format txid:vout
			utxoKey := fmt.Sprintf("%s:%d", input.Txid, input.Vout)

			// delete the utxo from the utxo index
			inputUtxo, ok := utxoIndex.Index[utxoKey]
			if !ok {
				log.Panicf("%s does not exist in the utxo index", utxoKey)
			}
			delete(utxoIndex.Index, utxoKey)
			deletedUTXO = append(deletedUTXO, utxoKey)

			// add the utxo's ordinals to the list of ordinals to be transferred
			ordinals = append(ordinals, inputUtxo.Ordinals...)
		}

		for _, output := range tx.Outputs {
			// transfer the ordinals to the output
			transferred, remaining := TransferRanges(ordinals, output.Value)
			output.Ordinals = transferred
			ordinals = remaining

			// add the output to the utxo index
			u := fmt.Sprintf("%s:%d", tx.Txid, output.N)
			utxoIndex.Index[u] = output
		}

		// add the remaining ordinals to the coinbase ordinals
		// those are the ordinals spent on fees
		coinbaseOrdinals = append(coinbaseOrdinals, ordinals...)
	}

	for _, output := range block.Transactions[0].Outputs {
		// transfer the coinbase ordinals to the output
		transferred, remaining := TransferRanges(coinbaseOrdinals, output.Value)
		output.Ordinals = transferred
		coinbaseOrdinals = remaining

		// add the output to the utxo index
		u := fmt.Sprintf("%s:%d", block.Transactions[0].Txid, output.N)
		utxoIndex.Index[u] = output
	}

	// return the modified entries in the UTXO index
	return deletedUTXO
}
