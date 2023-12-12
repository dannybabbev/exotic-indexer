package ordinals

import (
	"testing"
)

func TestSubsidy(t *testing.T) {
	var tests = []struct {
		height int
		want   int64
	}{
		{0, 5000000000},
		{1, 5000000000},
		{210000, 2500000000},
		{420000, 1250000000},
		{630000, 625000000},
	}

	for _, test := range tests {
		if got := Subsidy(test.height); got != test.want {
			t.Errorf("subsidy(%d) = %d; want: %d", test.height, got, test.want)
		}
	}
}

func TestFirstOrdinal(t *testing.T) {
	var tests = []struct {
		height int
		want   int64
	}{
		{0, 0},
		{1, 5000000000},
		{2, 10000000000},
		{3, 15000000000},
		{210000, 1050000000000000},
		{807680, 1948550000000000},
	}

	for _, test := range tests {
		if got := FirstOrdinal(test.height); got != test.want {
			t.Errorf("firstOrdinal(%d) = %d; want: %d", test.height, got, test.want)
		}
	}
}

func TestTransferRanges(t *testing.T) {
	var tests = []struct {
		ordinals        []*Range
		value           int64
		wantTransferred []*Range
		wantRemaining   []*Range
	}{
		{[]*Range{{0, 100}, {100, 100}, {200, 100}}, 150, []*Range{{0, 100}, {100, 50}}, []*Range{{150, 50}, {200, 100}}},
		{[]*Range{{0, 100}}, 50, []*Range{{0, 50}}, []*Range{{50, 50}}},
	}

	for _, test := range tests {
		gotTransferred, gotRemaining := TransferRanges(test.ordinals, test.value)

		for i, r := range gotTransferred {
			e := test.wantTransferred[i]
			if (r.Size != e.Size) || (r.Start != e.Start) {
				t.Errorf("got transferred start: %d size %d; want start %d, size %d", r.Start, r.Size, e.Start, e.Size)
			}
		}

		if len(gotRemaining) != len(test.wantRemaining) {
			t.Errorf("got remaining length: %d; expected length: %d", len(gotRemaining), len(test.wantRemaining))
		}

		for i, r := range gotRemaining {
			e := test.wantRemaining[i]
			if (r.Size != e.Size) || (r.Start != e.Start) {
				t.Errorf("got remaining start: %d, size: %d; expected start: %d, size: %d", r.Start, r.Size, e.Start, e.Size)
			}
		}

		if len(gotTransferred) != len(test.wantTransferred) {
			t.Errorf("got transferred length: %d; expected length: %d", len(gotTransferred), len(test.wantTransferred))
		}
	}
}

func TestAssignOrdinals(t *testing.T) {
	var coinbaseTx = Transaction{
		Txid:   "txcoinbase",
		Inputs: []*Input{},
		Outputs: []*Output{
			{
				Value: 5000000150,
				Address: ScriptPubKey{
					Addresses: []string{"addr_receive_cb"},
				},
			},
		},
	}

	var inputs = []*Input{
		{
			Txid:     "tx1",
			Vout:     0,
			Ordinals: []*Range{{100, 50}, {500, 200}},
		},
		{
			Txid:     "tx2",
			Vout:     0,
			Ordinals: []*Range{{800, 50}, {1000, 200}},
		},
	}

	var outputs = []*Output{
		{
			Value: 100,
			Address: ScriptPubKey{
				Addresses: []string{"addr_john"},
			},
			N: 0,
		},
		{
			Value: 200,
			Address: ScriptPubKey{
				Addresses: []string{"addr_danny"},
			},
			N: 1,
		},
		{
			Value: 50,
			Address: ScriptPubKey{
				Addresses: []string{"addr_john"},
			},
			N: 2,
		},
	}

	utxoIndex := NewUTXOIndex()
	utxoIndex.Index["tx1:0"] = &Output{
		Ordinals: []*Range{{100, 50}, {500, 200}},
		Address: ScriptPubKey{
			Addresses: []string{"addr_alice"},
		},
	}

	utxoIndex.Index["tx2:0"] = &Output{
		Ordinals: []*Range{{800, 50}, {1000, 200}},
		Address: ScriptPubKey{
			Addresses: []string{"addr_bob"},
		},
	}

	var tx1 = Transaction{Txid: "tx3", Inputs: inputs, Outputs: outputs}
	var wantCoinbaseOrdinals = []Range{{500000000000, 5000000000}, {1050, 150}}
	var wantTxOrdinals = []*Output{
		{
			Value:    100,
			Address:  ScriptPubKey{},
			Ordinals: []*Range{{100, 50}, {500, 50}},
		},
		{
			Value:    200,
			Address:  ScriptPubKey{},
			Ordinals: []*Range{{550, 150}, {800, 50}},
		},
		{
			Value:    50,
			Address:  ScriptPubKey{},
			Ordinals: []*Range{{1000, 50}},
		},
	}
	wantDeleted := []string{"tx1:0", "tx2:0"}
	wantAppended := map[string]bool{"tx3:0": true, "tx3:1": true, "tx3:2": true, "txcoinbase:0": true}

	var block = Block{
		Height: 100,
		Transactions: []*Transaction{
			&coinbaseTx,
			&tx1,
		},
	}

	deleted := AssignOrdinals(&block, utxoIndex)

	for i, r := range block.Transactions[0].Outputs[0].Ordinals {
		e := wantCoinbaseOrdinals[i]
		if (r.Size != e.Size) || (r.Start != e.Start) {
			t.Errorf("got coinbase ordinal start: %d, size: %d; expected start: %d, size: %d", r.Start, r.Size, e.Start, e.Size)
		}
	}

	for i, o := range block.Transactions[1].Outputs {
		for j, r := range o.Ordinals {
			e := wantTxOrdinals[i].Ordinals[j]
			if (r.Size != e.Size) || (r.Start != e.Start) {
				t.Errorf("got tx ordinal start: %d, size: %d; expected start: %d, size: %d", r.Start, r.Size, e.Start, e.Size)
			}
		}
	}

	if len(utxoIndex.Index) != len(wantAppended) {
		t.Error("utxoIndex has incorrect number of elements")
	}

	for k, _ := range wantAppended {
		if _, ok := utxoIndex.Index[k]; !ok {
			t.Errorf("%s is not is the utxo index", k)
		}
	}

	for i, w := range wantDeleted {
		if v := deleted[i]; v != w {
			t.Errorf("%s does not exist", w)
		}
	}
}
