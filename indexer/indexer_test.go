package indexer

import (
	"fmt"
	"log"
	"testing"

	"github.com/dannybabbev/ord-api/ordinals"
	badger "github.com/dgraph-io/badger/v4"
	"github.com/stretchr/testify/assert"
)

// TestSyncToPizzaRanges requires a local bitcoind instance to be running.
// func TestSyncToPizzaRanges(t *testing.T) {
// 	conf := conf.NewConf()

// 	bc, err := bitcoind.New(
// 		conf.BitcoinRPCHost,
// 		conf.BitcoinRPCPort,
// 		conf.BitcoinRPCUser,
// 		conf.BitcoinRPCPass,
// 		false,
// 	)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	wantRanges := exotic.ReadRangesFromOrdResponse("../data/pizza.json")

// 	opts := badger.DefaultOptions("").WithInMemory(true)
// 	db, err := badger.Open(opts)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	defer db.Close()

// 	blockchain := NewIndexer(bc, db)

// 	stopChan := make(chan struct{})
// 	doneChan := make(chan struct{}, 1)
// 	blockchain.SyncToBlock(57043, stopChan, doneChan)

// 	gotRanges, err := blockchain.GetOrdinalsForUTXO("a1075db55d416d3ca199f55b6084e2115b9345e16c5cf302fc80e9d5fbf5d48d:0")
// 	if err != nil {
// 		t.Errorf("error getting ordinals for utxo: %v", err)
// 	}

// 	assert.Equal(t, wantRanges, gotRanges)

// 	for i, r := range wantRanges {
// 		if r.Start != gotRanges[i].Start {
// 			t.Errorf("want start %d, got start %d", r.Start, gotRanges[i].Start)
// 		}
// 		if r.Size != gotRanges[i].Size {
// 			t.Errorf("want size %d, got size %d", r.Size, gotRanges[i].Size)
// 		}
// 	}

// 	res, err := blockchain.GetOrdinalsForUTXO("1e133f7de73ac7d074e2746a3d6717dfc99ecaa8e9f9fade2cb8b0b20a5e0441:0")
// 	if err == nil {
// 		t.Errorf("expected error getting ordinals for utxo")
// 	}

// 	fmt.Println(res)
// }

func TestIndexerIsManagingTheUTXOIndex(t *testing.T) {
	blockHashes := make(map[uint64]string)
	blockHashes[52332] = "000000001b87f065f8a5d4a5240fc8ef7fe03884cf074f568814f7e8d21dc632"

	rawBlocks := make(map[string]string)
	rawBlocks["000000001b87f065f8a5d4a5240fc8ef7fe03884cf074f568814f7e8d21dc632"] = "01000000a0eb1aaf03581db0fed450ab6f57f9d4c139d70666b9530cb10c430600000000d17373f6eab5996fb9d856cccc87cec31e032b157a6c335707b5a43708ae2756e20acf4ba7bc201c5b6320170201000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0704a7bc201c0157ffffffff0100f2052a0100000043410421be99f3512047a98399269cf9a9e494fad367e1c44d5c08675df0f4ac4a6bd7fb3fdee852785c84705de943c94214ad434ba8838c75e0580237abe62570146aac0000000001000000033cc6b815284642a6ccb461a5da3a2d85c3a8a5cdf32c6f75ddece74980f67d400000000048473044022002ec490cad4a97f3d6d1968215e7788f493cd614f4ad34c2f603628440ddd075022008726ab7276fec551803075a917bb52d508db01a138a678b59b1e87edf546f6c01ffffffff449edf2816a08f3aee2dc50a2779182a13222359c4a4a940edd91e9d5e728fd800000000494830450220424746fbe0a7d5ddeb30ebbea23faed597d0e2656c3730230ba257e0b6c50468022100b46bd890981d0dfe7e649da135c0031fa62d6f8b5eeeb3db9a90804d78ba8c5301ffffffff9b91a5c8dc83b8d31ec6a133591b3a7f00fa7eca2292ec83466e295bf02d18da000000004847304402205a18b1ab648e27fc2622f007ac07626461d9c080f4ef03394c07fed303853055022036a592a258c6f7760c3a2706f58e5e3a425e3fbaf846364801bc227acecae08501ffffffff0100d6117e030000001976a91431f19a7d0379f56cb3be0761c21f1f0c9553a47f88ac00000000"

	bitcoind := &MockBitcoind{
		BlockCountReturn: 52332,
		BlockHashReturn:  blockHashes,
		RawBlockReturn:   rawBlocks,
	}

	opts := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	blockchain := NewIndexer(bitcoind, db).WithKeepBlockHistory(0)

	stats := &SyncStats{
		SyncHeight:    52331,
		SyncBlockHash: "0000000006430cb10c53b96606d739c1d4f9576fab50d4feb01d5803af1aeba0",
	}

	blockchain.stats = stats

	// Seed the DB with some UTXOs
	utxoIndexStart := ordinals.NewUTXOIndex()
	ordinalsIn1 := []*ordinals.Range{
		{
			Start: ordinals.FirstOrdinal(52147),
			Size:  ordinals.Subsidy(100),
		},
	}
	utxoIndexStart.Index["407df68049e7ecdd756f2cf3cda5a8c3852d3adaa561b4cca642462815b8c63c:0"] = &ordinals.Output{
		Value:    5e9,
		Ordinals: ordinalsIn1,
	}

	ordinalsIn2 := []*ordinals.Range{
		{
			Start: ordinals.FirstOrdinal(52119),
			Size:  ordinals.Subsidy(100),
		},
	}
	utxoIndexStart.Index["d88f725e9d1ed9ed40a9a4c4592322132a1879270ac52dee3a8fa01628df9e44:0"] = &ordinals.Output{
		Value:    5e9,
		Ordinals: ordinalsIn2,
	}

	ordinalsIn3 := []*ordinals.Range{
		{
			Start: ordinals.FirstOrdinal(52205),
			Size:  ordinals.Subsidy(100),
		},
	}
	utxoIndexStart.Index["da182df05b296e4683ec9222ca7efa007f3a1b5933a1c61ed3b883dcc8a5919b:0"] = &ordinals.Output{
		Value:    5e9,
		Ordinals: ordinalsIn3,
	}

	blockchain.utxoIndex = utxoIndexStart

	// Flush the seed updates to the DB
	blockchain.updateDB()

	// The function we are testing
	stopChan := make(chan struct{})
	doneChan := make(chan struct{}, 1)
	blockchain.SyncToChainTip(stopChan, doneChan)

	// Assertions
	wantUTXO := map[string]*ordinals.Output{
		"a9fc6cd2f8517b9f6790c586f8c4e1dfb9a50eea49fb86116ebeeed2eb0be8e9:0": {
			Value: 5e9,
			Ordinals: []*ordinals.Range{
				{
					Start: ordinals.FirstOrdinal(52332),
					Size:  ordinals.Subsidy(100),
				},
			},
		},
		"ba93e193c270031da6e97294d6e0cc59faab8b448880bf2f9c57cb4eea901523:0": {
			Value: 15e9,
			Ordinals: []*ordinals.Range{
				ordinalsIn1[0],
				ordinalsIn2[0],
				ordinalsIn3[0],
			},
		},
	}

	hasSyncStats := false
	hasUtxoHistory := false
	items := 0
	err = blockchain.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			fmt.Printf("key=%s\n", k)
			items++
			if string(k) == "syncStats" {
				hasSyncStats = true
			} else if string(k) == "utxoHistory" {
				hasUtxoHistory = true
			} else {
				item.Value(func(v []byte) error {
					output := &ordinals.Output{}
					blockchain.decodeBytes(v, output)
					mapKey := string(k)[2:]
					gotOutput, ok := wantUTXO[mapKey]
					assert.True(t, ok)

					assert.Equal(t, len(gotOutput.Ordinals), len(output.Ordinals))
					for i, r := range gotOutput.Ordinals {
						assert.Equal(t, r.Start, output.Ordinals[i].Start)
						assert.Equal(t, r.Size, output.Ordinals[i].Size)
					}
					return nil
				})
			}
		}
		return nil
	})

	assert.True(t, hasSyncStats)
	assert.True(t, hasUtxoHistory)
	assert.Nil(t, err)
	assert.Equal(t, 4, items)

	// TODO: This should be a separate test
	testGetUTXO := "ba93e193c270031da6e97294d6e0cc59faab8b448880bf2f9c57cb4eea901523:0"
	res, err := blockchain.GetOrdinalsForUTXO(testGetUTXO)
	assert.Nil(t, err)

	assert.Equal(t, len(res), 3)

	for i, r := range res {
		if r.Start != wantUTXO[testGetUTXO].Ordinals[i].Start {
			t.Errorf("want start %d, got start %d", r.Start, wantUTXO[testGetUTXO].Ordinals[i].Start)
		}
		if r.Size != wantUTXO[testGetUTXO].Ordinals[i].Size {
			t.Errorf("want size %d, got size %d", r.Size, wantUTXO[testGetUTXO].Ordinals[i].Size)
		}
	}
}

func IsInBadgerDB(db *badger.DB, key string) bool {
	var res bool
	err := db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		if err == nil {
			res = true
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return res
}

func TestIndexerIsManagingTheUTXOHistoryAtTheChainTip(t *testing.T) {
	blockHashes := make(map[uint64]string)
	blockHashes[52329] = "0000000011688ef051068b95f35b1f76f691a70f2b951a85dca073f9aa1d8f96"
	blockHashes[52330] = "000000000154a79f8ed2564e7869de5d9243b3387e09dc01912dedaf7aa52c94"
	blockHashes[52331] = "0000000006430cb10c53b96606d739c1d4f9576fab50d4feb01d5803af1aeba0"
	blockHashes[52332] = "000000001b87f065f8a5d4a5240fc8ef7fe03884cf074f568814f7e8d21dc632"
	blockHashes[52333] = "000000000c80e6c43482013903a10f8fd35d3ffab79cfac8148ae31a5f9a41b6"
	blockHashes[52334] = "00000000080236123d70e011faeccf172cb806a5847084e74bc3172d3a9ff453"

	rawBlocks := make(map[string]string)
	rawBlocks["0000000011688ef051068b95f35b1f76f691a70f2b951a85dca073f9aa1d8f96"] = "01000000c783ca947f8a04549fb44c724a63ff1e849e1c694f57dc019918f40100000000941005e5c2807aeb702408d38b79fe9171c36c23ec340eda2e021eef9dbff492e5fbce4ba7bc201cc0c5a2050101000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0804a7bc201c021d06ffffffff0100f2052a01000000434104d53640e6743f1f8d53888fa0ad60d365fdc9f4c31587e62d79b895186cdd2c20013de7fa57937fe4b9c899365c497dfd9e9d09ce0a223e2b8ce21ec5d1dc7103ac00000000"
	rawBlocks["000000000154a79f8ed2564e7869de5d9243b3387e09dc01912dedaf7aa52c94"] = "01000000968f1daaf973a0dc851a952b0fa791f6761f5bf3958b0651f08e681100000000aeeb22ca6b2be23833f75d48244519f31dc942e22d17e6b0168c0b5b333ae872eefdce4ba7bc201cb4c03e020101000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0704a7bc201c0104ffffffff0100f2052a010000004341049099b8f925f3a61d64cb3117076cfb64f0559255236e9185afdaed7fc75b12fdb71364d7f0e88f5ee34993e8d65cce24b92e246294038a7594002c84a109c3d3ac00000000"
	rawBlocks["0000000006430cb10c53b96606d739c1d4f9576fab50d4feb01d5803af1aeba0"] = "01000000942ca57aafed2d9101dc097e38b343925dde69784e56d28e9fa7540100000000de5a762ba5bcb6d18ed6d24cef12a1f431949891f758ceac92a65281f946bd28de04cf4ba7bc201ce302961f0301000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0704a7bc201c0101ffffffff0100f2052a010000004341047f89d6a759b9a644bb606c1ec95d722e9d4f5b1ee4684f7f547932290d57a2889c9d9700b1988a81ba04c17530062633ff347d568f17bcc80aaa495226b2414cac00000000010000000719a01a6aa000c4b9d3ebd7cdf2cff38a15a160ed9c94d05a842d91a0c43a3685000000004a493046022100ceba61e429e59ab09b720923b4f7ab27208e73ccdbd755e39eb153dc24bc80b00221008b6759ab89037c234d85cf4e9b4a36c2a4e9d64585c8f1e1ad58ba3d712b9bda01ffffffff4b0a69c721932e269687e521ab505825f774c8a9402414260b4fe6f4be1b0541000000004a493046022100a5253ebb43cb8293fcbf01524cd99d40bdb8664ff63724913b693d104dca73ad022100e9fe8ef45b30b3c61b350e0e78641ee32ca8f3631eac27f93e9062c5ac0fcc1101ffffffff70fb2e52e7ecdcf2ca0a400bf940f339cc14006102c0e9964284acb47e7e95fb0000000049483045022100d574d87e09e71202c81db6ae46c96c466e5ab5e989f3d88a2a31017a155ba5650220606cb2a29a0305d27d8e43ddcad02fb38d19178f75908d9b882bf7c01f93a1f401ffffffff7f9f29356c6b6a9e973b0a38241b51c923c1acbe920726501c5d7059e7f65a6d000000004a493046022100de4cac98a9be8f3b279669018974db48ada2243011aeb3330b02fba5fa41a748022100a32d828b95debc82a9f10fe457869339d11deb147bdabfbc4a94df96597c4b9d01ffffffffe61d4bffaddb9a9ea3eb4204406bfa4f31fb06124d8f2b70be7029a8fccf801100000000494830450220255008acd3af6386c99bb2e2e0a208b89f103cf7e61944de2507cfb2e671871e02210085f0582923248340e39301c79fa998264e89ebfeca845b86b44a70cd02c0c7a601ffffffffe6c28a7bb5b7491eddb7dd60cd0b65bf857bab4c907ede045d4be07dab0662060000000048473044022012ab455e288bfcf89bcc1f857dfaa10fec3a72534c22c0607c44bb46e8d05265022056c36e748a249d988258d3adf69ab2f6ba4b0c56eaff5003c4ace3fbf159d1be01ffffffffe98bee04da5a8f900f1935c24a6691e506bc89f14bd28a882f9d652a623d06db0000000048473044022069252b55af9dc19f1fba9e1027358b92ce62508167aec65544e5feda1d26ff41022065efec2f9a2f347761fce6394db88defd107cb5ea9d3a39556619af7f33c30dc01ffffffff01009e2926080000001976a914ae04ed787555df42f6ccd36028062dd4129a12d588ac0000000001000000017273dffb0fa52f7676fd14953cf58afd146f37686239b9dd996409d63fbb5717000000004a493046022100d2843f5c5f1caafd1d8a1a0317514a307c0d75f25665100b3c9888914a319687022100cfa5bae41057bebb1aed2899bda1229070860f79313d677e21c506ad359608ac01ffffffff0200286bee00000000434104ce856e3208e5c24944f2932f3cd6e690bdef60061b0486dc6c46663e4ee82c23772cd3ee38b44f153db63253aad4fd7bee29cb689a1a9d7b790ac0045b5a7e03ac00ca9a3b000000001976a914b3fe6e2138c2341af58206151bf3173fa1af90aa88ac00000000"
	rawBlocks["000000001b87f065f8a5d4a5240fc8ef7fe03884cf074f568814f7e8d21dc632"] = "01000000a0eb1aaf03581db0fed450ab6f57f9d4c139d70666b9530cb10c430600000000d17373f6eab5996fb9d856cccc87cec31e032b157a6c335707b5a43708ae2756e20acf4ba7bc201c5b6320170201000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0704a7bc201c0157ffffffff0100f2052a0100000043410421be99f3512047a98399269cf9a9e494fad367e1c44d5c08675df0f4ac4a6bd7fb3fdee852785c84705de943c94214ad434ba8838c75e0580237abe62570146aac0000000001000000033cc6b815284642a6ccb461a5da3a2d85c3a8a5cdf32c6f75ddece74980f67d400000000048473044022002ec490cad4a97f3d6d1968215e7788f493cd614f4ad34c2f603628440ddd075022008726ab7276fec551803075a917bb52d508db01a138a678b59b1e87edf546f6c01ffffffff449edf2816a08f3aee2dc50a2779182a13222359c4a4a940edd91e9d5e728fd800000000494830450220424746fbe0a7d5ddeb30ebbea23faed597d0e2656c3730230ba257e0b6c50468022100b46bd890981d0dfe7e649da135c0031fa62d6f8b5eeeb3db9a90804d78ba8c5301ffffffff9b91a5c8dc83b8d31ec6a133591b3a7f00fa7eca2292ec83466e295bf02d18da000000004847304402205a18b1ab648e27fc2622f007ac07626461d9c080f4ef03394c07fed303853055022036a592a258c6f7760c3a2706f58e5e3a425e3fbaf846364801bc227acecae08501ffffffff0100d6117e030000001976a91431f19a7d0379f56cb3be0761c21f1f0c9553a47f88ac00000000"
	rawBlocks["000000000c80e6c43482013903a10f8fd35d3ffab79cfac8148ae31a5f9a41b6"] = "0100000032c61dd2e8f71488564f07cf8438e07fefc80f24a5d4a5f865f0871b000000001d5f8e0f1824974c8c456393f7ff2df574a6bc5bba1805353f12f17f15fe2b61840bcf4ba7bc201ceb3e65030101000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0804a7bc201c028a00ffffffff0100f2052a0100000043410438fa5488916e61472d8f7cd68f65b0a5f6437cafc965255885fe47c4728faba13e428e1573b4782168ea872949175ab3b42400f5f075df03493bd7793c88b7f9ac00000000"
	rawBlocks["00000000080236123d70e011faeccf172cb806a5847084e74bc3172d3a9ff453"] = "01000000b6419a5f1ae38a14c8fa9cb7fa3f5dd38f0fa10339018234c4e6800c00000000869d2225127e5cb34ba325395f88c5a7711ddae3871188d85ddf0d8de842b8fae40ccf4ba7bc201c7c39540a0101000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0704a7bc201c0105ffffffff0100f2052a01000000434104ab3d5a3a0084e7e80c9544ee3f3510efd3b83bfd52fb2b539dfeb6cab9f2924b6ec2b421051c2cb4c132352285a1a2f6f8460ef8c25669feae4f44f9dc52cc2cac00000000"

	utxoIndex := ordinals.NewUTXOIndex()

	utxoIndex.Index["a"] = simpleOutput(0)
	utxoIndex.Index["b"] = simpleOutput(1)
	utxoIndex.Index["c"] = simpleOutput(2)
	utxoIndex.Index["85363ac4a0912d845ad0949ced60a1158af3cff2cdd7ebd3b9c400a06a1aa019:0"] = simpleOutput(10)
	utxoIndex.Index["41051bbef4e64f0b26142440a9c874f7255850ab21e58796262e9321c7690a4b:0"] = simpleOutput(11)
	utxoIndex.Index["fb957e7eb4ac844296e9c002610014cc39f340f90b400acaf2dcece7522efb70:0"] = simpleOutput(12)
	utxoIndex.Index["6d5af6e759705d1c50260792beacc123c9511b24380a3b979e6a6b6c35299f7f:0"] = simpleOutput(13)
	utxoIndex.Index["1180cffca82970be702b8f4d1206fb314ffa6b400442eba39e9adbadff4b1de6:0"] = simpleOutput(14)
	utxoIndex.Index["066206ab7de04b5d04de7e904cab7b85bf650bcd60ddb7dd1e49b7b57b8ac2e6:0"] = simpleOutput(15)
	utxoIndex.Index["db063d622a659d2f888ad24bf189bc06e591664ac235190f908f5ada04ee8be9:0"] = simpleOutput(16)
	utxoIndex.Index["1757bb3fd6096499ddb9396268376f14fd8af53c9514fd76762fa50ffbdf7372:0"] = simpleOutput(17)
	utxoIndex.Index["407df68049e7ecdd756f2cf3cda5a8c3852d3adaa561b4cca642462815b8c63c:0"] = simpleOutput(18)
	utxoIndex.Index["d88f725e9d1ed9ed40a9a4c4592322132a1879270ac52dee3a8fa01628df9e44:0"] = simpleOutput(19)
	utxoIndex.Index["da182df05b296e4683ec9222ca7efa007f3a1b5933a1c61ed3b883dcc8a5919b:0"] = simpleOutput(20)

	bitcoind := &MockBitcoind{
		BlockCountReturn: 52329, // starting heignt
		BlockHashReturn:  blockHashes,
		RawBlockReturn:   rawBlocks,
	}

	opts := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	blockchain := NewIndexer(bitcoind, db).WithKeepBlockHistory(3)

	stats := &SyncStats{
		SyncHeight:     52329,
		SyncBlockHash:  "0000000011688ef051068b95f35b1f76f691a70f2b951a85dca073f9aa1d8f96",
		ReorgsDetected: make([]int, 0),
	}

	mockStartingUTXOHistory := [][]string{
		{"a"},
		{"b"},
		{"c"},
	}

	blockchain.stats = stats
	blockchain.utxoHistory = mockStartingUTXOHistory
	blockchain.utxoIndex = utxoIndex
	blockchain.updateDB()

	assert.True(t, IsInBadgerDB(db, "u-a"))
	assert.True(t, IsInBadgerDB(db, "u-b"))
	assert.True(t, IsInBadgerDB(db, "u-c"))

	run := func() {
		stopChan := make(chan struct{})
		doneChan := make(chan struct{}, 1)
		bitcoind.BlockCountReturn += 1
		blockchain.SyncToChainTip(stopChan, doneChan)
	}

	// run 1
	run()

	// Assert database
	wantUtxoHistory := [][]string{
		{"b"},
		{"c"},
		{}, // here it is an empty array because we did not load it from DB
	}

	assert.Equal(t, wantUtxoHistory, blockchain.utxoHistory)
	assert.False(t, IsInBadgerDB(db, "u-a"))
	assert.True(t, IsInBadgerDB(db, "u-b"))
	assert.True(t, IsInBadgerDB(db, "u-c"))

	// run 2
	run()

	// Assert database
	wantUtxoHistory = [][]string{
		{"c"},
		nil, // Here loading it from DB makes it nil, ffs
		{"85363ac4a0912d845ad0949ced60a1158af3cff2cdd7ebd3b9c400a06a1aa019:0", "41051bbef4e64f0b26142440a9c874f7255850ab21e58796262e9321c7690a4b:0", "fb957e7eb4ac844296e9c002610014cc39f340f90b400acaf2dcece7522efb70:0", "6d5af6e759705d1c50260792beacc123c9511b24380a3b979e6a6b6c35299f7f:0", "1180cffca82970be702b8f4d1206fb314ffa6b400442eba39e9adbadff4b1de6:0", "066206ab7de04b5d04de7e904cab7b85bf650bcd60ddb7dd1e49b7b57b8ac2e6:0", "db063d622a659d2f888ad24bf189bc06e591664ac235190f908f5ada04ee8be9:0", "1757bb3fd6096499ddb9396268376f14fd8af53c9514fd76762fa50ffbdf7372:0"},
	}

	assert.Equal(t, wantUtxoHistory, blockchain.utxoHistory)
	assert.False(t, IsInBadgerDB(db, "u-a"))
	assert.False(t, IsInBadgerDB(db, "u-b"))
	assert.True(t, IsInBadgerDB(db, "u-c"))
	assert.True(t, IsInBadgerDB(db, "u-85363ac4a0912d845ad0949ced60a1158af3cff2cdd7ebd3b9c400a06a1aa019:0"))
	assert.True(t, IsInBadgerDB(db, "u-41051bbef4e64f0b26142440a9c874f7255850ab21e58796262e9321c7690a4b:0"))

	// run 3
	run()

	// Assert database
	wantUtxoHistory = [][]string{
		nil,
		{"85363ac4a0912d845ad0949ced60a1158af3cff2cdd7ebd3b9c400a06a1aa019:0", "41051bbef4e64f0b26142440a9c874f7255850ab21e58796262e9321c7690a4b:0", "fb957e7eb4ac844296e9c002610014cc39f340f90b400acaf2dcece7522efb70:0", "6d5af6e759705d1c50260792beacc123c9511b24380a3b979e6a6b6c35299f7f:0", "1180cffca82970be702b8f4d1206fb314ffa6b400442eba39e9adbadff4b1de6:0", "066206ab7de04b5d04de7e904cab7b85bf650bcd60ddb7dd1e49b7b57b8ac2e6:0", "db063d622a659d2f888ad24bf189bc06e591664ac235190f908f5ada04ee8be9:0", "1757bb3fd6096499ddb9396268376f14fd8af53c9514fd76762fa50ffbdf7372:0"},
		{"407df68049e7ecdd756f2cf3cda5a8c3852d3adaa561b4cca642462815b8c63c:0", "d88f725e9d1ed9ed40a9a4c4592322132a1879270ac52dee3a8fa01628df9e44:0", "da182df05b296e4683ec9222ca7efa007f3a1b5933a1c61ed3b883dcc8a5919b:0"},
	}

	assert.Equal(t, wantUtxoHistory, blockchain.utxoHistory)

	assert.False(t, IsInBadgerDB(db, "u-a"))
	assert.False(t, IsInBadgerDB(db, "u-b"))
	assert.False(t, IsInBadgerDB(db, "u-c"))
	assert.True(t, IsInBadgerDB(db, "u-85363ac4a0912d845ad0949ced60a1158af3cff2cdd7ebd3b9c400a06a1aa019:0"))
	assert.True(t, IsInBadgerDB(db, "u-41051bbef4e64f0b26142440a9c874f7255850ab21e58796262e9321c7690a4b:0"))

	// run 4
	run()
	assert.True(t, IsInBadgerDB(db, "u-85363ac4a0912d845ad0949ced60a1158af3cff2cdd7ebd3b9c400a06a1aa019:0"))
	assert.True(t, IsInBadgerDB(db, "u-41051bbef4e64f0b26142440a9c874f7255850ab21e58796262e9321c7690a4b:0"))

	run()
	// Assert database
	assert.False(t, IsInBadgerDB(db, "u-85363ac4a0912d845ad0949ced60a1158af3cff2cdd7ebd3b9c400a06a1aa019:0"))
	assert.False(t, IsInBadgerDB(db, "u-41051bbef4e64f0b26142440a9c874f7255850ab21e58796262e9321c7690a4b:0"))
}

func simpleOrdinal(start int) []*ordinals.Range {
	return []*ordinals.Range{
		{
			Start: ordinals.FirstOrdinal(start),
			Size:  ordinals.Subsidy(100),
		},
	}
}

func simpleOutput(start int) *ordinals.Output {
	return &ordinals.Output{
		Value:    5e9,
		Ordinals: simpleOrdinal(start),
	}
}

func TestIndexerHandlesAReorg(t *testing.T) {
	blockHashes := make(map[uint64]string)
	blockHashes[52329] = "0000000011688ef051068b95f35b1f76f691a70f2b951a85dca073f9aa1d8f96"
	blockHashes[52330] = "000000000154a79f8ed2564e7869de5d9243b3387e09dc01912dedaf7aa52c94"
	blockHashes[52331] = "0000000006430cb10c53b96606d739c1d4f9576fab50d4feb01d5803af1aeba0"
	blockHashes[52332] = "000000001b87f065f8a5d4a5240fc8ef7fe03884cf074f568814f7e8d21dc632"

	rawBlocks := make(map[string]string)
	rawBlocks["0000000011688ef051068b95f35b1f76f691a70f2b951a85dca073f9aa1d8f96"] = "01000000c783ca947f8a04549fb44c724a63ff1e849e1c694f57dc019918f40100000000941005e5c2807aeb702408d38b79fe9171c36c23ec340eda2e021eef9dbff492e5fbce4ba7bc201cc0c5a2050101000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0804a7bc201c021d06ffffffff0100f2052a01000000434104d53640e6743f1f8d53888fa0ad60d365fdc9f4c31587e62d79b895186cdd2c20013de7fa57937fe4b9c899365c497dfd9e9d09ce0a223e2b8ce21ec5d1dc7103ac00000000"
	rawBlocks["000000000154a79f8ed2564e7869de5d9243b3387e09dc01912dedaf7aa52c94"] = "01000000968f1daaf973a0dc851a952b0fa791f6761f5bf3958b0651f08e681100000000aeeb22ca6b2be23833f75d48244519f31dc942e22d17e6b0168c0b5b333ae872eefdce4ba7bc201cb4c03e020101000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0704a7bc201c0104ffffffff0100f2052a010000004341049099b8f925f3a61d64cb3117076cfb64f0559255236e9185afdaed7fc75b12fdb71364d7f0e88f5ee34993e8d65cce24b92e246294038a7594002c84a109c3d3ac00000000"
	rawBlocks["0000000006430cb10c53b96606d739c1d4f9576fab50d4feb01d5803af1aeba0"] = "01000000942ca57aafed2d9101dc097e38b343925dde69784e56d28e9fa7540100000000de5a762ba5bcb6d18ed6d24cef12a1f431949891f758ceac92a65281f946bd28de04cf4ba7bc201ce302961f0301000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0704a7bc201c0101ffffffff0100f2052a010000004341047f89d6a759b9a644bb606c1ec95d722e9d4f5b1ee4684f7f547932290d57a2889c9d9700b1988a81ba04c17530062633ff347d568f17bcc80aaa495226b2414cac00000000010000000719a01a6aa000c4b9d3ebd7cdf2cff38a15a160ed9c94d05a842d91a0c43a3685000000004a493046022100ceba61e429e59ab09b720923b4f7ab27208e73ccdbd755e39eb153dc24bc80b00221008b6759ab89037c234d85cf4e9b4a36c2a4e9d64585c8f1e1ad58ba3d712b9bda01ffffffff4b0a69c721932e269687e521ab505825f774c8a9402414260b4fe6f4be1b0541000000004a493046022100a5253ebb43cb8293fcbf01524cd99d40bdb8664ff63724913b693d104dca73ad022100e9fe8ef45b30b3c61b350e0e78641ee32ca8f3631eac27f93e9062c5ac0fcc1101ffffffff70fb2e52e7ecdcf2ca0a400bf940f339cc14006102c0e9964284acb47e7e95fb0000000049483045022100d574d87e09e71202c81db6ae46c96c466e5ab5e989f3d88a2a31017a155ba5650220606cb2a29a0305d27d8e43ddcad02fb38d19178f75908d9b882bf7c01f93a1f401ffffffff7f9f29356c6b6a9e973b0a38241b51c923c1acbe920726501c5d7059e7f65a6d000000004a493046022100de4cac98a9be8f3b279669018974db48ada2243011aeb3330b02fba5fa41a748022100a32d828b95debc82a9f10fe457869339d11deb147bdabfbc4a94df96597c4b9d01ffffffffe61d4bffaddb9a9ea3eb4204406bfa4f31fb06124d8f2b70be7029a8fccf801100000000494830450220255008acd3af6386c99bb2e2e0a208b89f103cf7e61944de2507cfb2e671871e02210085f0582923248340e39301c79fa998264e89ebfeca845b86b44a70cd02c0c7a601ffffffffe6c28a7bb5b7491eddb7dd60cd0b65bf857bab4c907ede045d4be07dab0662060000000048473044022012ab455e288bfcf89bcc1f857dfaa10fec3a72534c22c0607c44bb46e8d05265022056c36e748a249d988258d3adf69ab2f6ba4b0c56eaff5003c4ace3fbf159d1be01ffffffffe98bee04da5a8f900f1935c24a6691e506bc89f14bd28a882f9d652a623d06db0000000048473044022069252b55af9dc19f1fba9e1027358b92ce62508167aec65544e5feda1d26ff41022065efec2f9a2f347761fce6394db88defd107cb5ea9d3a39556619af7f33c30dc01ffffffff01009e2926080000001976a914ae04ed787555df42f6ccd36028062dd4129a12d588ac0000000001000000017273dffb0fa52f7676fd14953cf58afd146f37686239b9dd996409d63fbb5717000000004a493046022100d2843f5c5f1caafd1d8a1a0317514a307c0d75f25665100b3c9888914a319687022100cfa5bae41057bebb1aed2899bda1229070860f79313d677e21c506ad359608ac01ffffffff0200286bee00000000434104ce856e3208e5c24944f2932f3cd6e690bdef60061b0486dc6c46663e4ee82c23772cd3ee38b44f153db63253aad4fd7bee29cb689a1a9d7b790ac0045b5a7e03ac00ca9a3b000000001976a914b3fe6e2138c2341af58206151bf3173fa1af90aa88ac00000000"
	rawBlocks["000000001b87f065f8a5d4a5240fc8ef7fe03884cf074f568814f7e8d21dc632"] = "01000000a0eb1aaf03581db0fed450ab6f57f9d4c139d70666b9530cb10c430600000000d17373f6eab5996fb9d856cccc87cec31e032b157a6c335707b5a43708ae2756e20acf4ba7bc201c5b6320170201000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0704a7bc201c0157ffffffff0100f2052a0100000043410421be99f3512047a98399269cf9a9e494fad367e1c44d5c08675df0f4ac4a6bd7fb3fdee852785c84705de943c94214ad434ba8838c75e0580237abe62570146aac0000000001000000033cc6b815284642a6ccb461a5da3a2d85c3a8a5cdf32c6f75ddece74980f67d400000000048473044022002ec490cad4a97f3d6d1968215e7788f493cd614f4ad34c2f603628440ddd075022008726ab7276fec551803075a917bb52d508db01a138a678b59b1e87edf546f6c01ffffffff449edf2816a08f3aee2dc50a2779182a13222359c4a4a940edd91e9d5e728fd800000000494830450220424746fbe0a7d5ddeb30ebbea23faed597d0e2656c3730230ba257e0b6c50468022100b46bd890981d0dfe7e649da135c0031fa62d6f8b5eeeb3db9a90804d78ba8c5301ffffffff9b91a5c8dc83b8d31ec6a133591b3a7f00fa7eca2292ec83466e295bf02d18da000000004847304402205a18b1ab648e27fc2622f007ac07626461d9c080f4ef03394c07fed303853055022036a592a258c6f7760c3a2706f58e5e3a425e3fbaf846364801bc227acecae08501ffffffff0100d6117e030000001976a91431f19a7d0379f56cb3be0761c21f1f0c9553a47f88ac00000000"

	bitcoind := &MockBitcoind{
		BlockCountReturn: 52332,
		BlockHashReturn:  blockHashes,
		RawBlockReturn:   rawBlocks,
	}

	utxoIndex := ordinals.NewUTXOIndex()

	utxoIndex.Index["85363ac4a0912d845ad0949ced60a1158af3cff2cdd7ebd3b9c400a06a1aa019:0"] = simpleOutput(10)
	utxoIndex.Index["41051bbef4e64f0b26142440a9c874f7255850ab21e58796262e9321c7690a4b:0"] = simpleOutput(11)
	utxoIndex.Index["fb957e7eb4ac844296e9c002610014cc39f340f90b400acaf2dcece7522efb70:0"] = simpleOutput(12)
	utxoIndex.Index["6d5af6e759705d1c50260792beacc123c9511b24380a3b979e6a6b6c35299f7f:0"] = simpleOutput(13)
	utxoIndex.Index["1180cffca82970be702b8f4d1206fb314ffa6b400442eba39e9adbadff4b1de6:0"] = simpleOutput(14)
	utxoIndex.Index["066206ab7de04b5d04de7e904cab7b85bf650bcd60ddb7dd1e49b7b57b8ac2e6:0"] = simpleOutput(15)
	utxoIndex.Index["db063d622a659d2f888ad24bf189bc06e591664ac235190f908f5ada04ee8be9:0"] = simpleOutput(16)
	utxoIndex.Index["1757bb3fd6096499ddb9396268376f14fd8af53c9514fd76762fa50ffbdf7372:0"] = simpleOutput(17)
	utxoIndex.Index["407df68049e7ecdd756f2cf3cda5a8c3852d3adaa561b4cca642462815b8c63c:0"] = simpleOutput(18)
	utxoIndex.Index["d88f725e9d1ed9ed40a9a4c4592322132a1879270ac52dee3a8fa01628df9e44:0"] = simpleOutput(19)
	utxoIndex.Index["da182df05b296e4683ec9222ca7efa007f3a1b5933a1c61ed3b883dcc8a5919b:0"] = simpleOutput(20)

	opts := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	blockchain := NewIndexer(bitcoind, db).WithKeepBlockHistory(3)
	blockchain.utxoIndex = utxoIndex

	stats := &SyncStats{
		SyncHeight:     52331,
		SyncBlockHash:  "wrongblockhash",
		ReorgsDetected: make([]int, 0),
	}

	utxoHistory := [][]string{
		// I am puting an existing UTXO here from the index
		{"85363ac4a0912d845ad0949ced60a1158af3cff2cdd7ebd3b9c400a06a1aa019:0"},
		{"1e133f7de73ac7d074e2746a3d6717dfc99ecaa8e9f9fade2cb8b0b20a5e0441:0"},
		{"1e133f7de73ac7d074e2746a3d6717dfc99ecaa8e9f9fade2cb8b0b20a5e0441:0"},
	}

	blockchain.stats = stats
	blockchain.utxoHistory = utxoHistory

	// Flush the seed updates to the DB
	blockchain.updateDB()

	// The function we are testing
	stopChan := make(chan struct{})
	doneChan := make(chan struct{}, 1)
	blockchain.SyncToChainTip(stopChan, doneChan)

	assert.Equal(t, &SyncStats{
		ChainTip:       52332,
		SyncHeight:     52332,
		SyncBlockHash:  "000000001b87f065f8a5d4a5240fc8ef7fe03884cf074f568814f7e8d21dc632",
		ReorgsDetected: []int{52332},
	}, blockchain.stats)

	// Make sure that the existing UTXO in hostory does not get deleted from the DB after a reorg
	assert.True(t, IsInBadgerDB(db, "u-85363ac4a0912d845ad0949ced60a1158af3cff2cdd7ebd3b9c400a06a1aa019:0"))

	// The UTXO history has to be rebuilt
	wantUtxoHistory := [][]string{
		{},
		{"85363ac4a0912d845ad0949ced60a1158af3cff2cdd7ebd3b9c400a06a1aa019:0", "41051bbef4e64f0b26142440a9c874f7255850ab21e58796262e9321c7690a4b:0", "fb957e7eb4ac844296e9c002610014cc39f340f90b400acaf2dcece7522efb70:0", "6d5af6e759705d1c50260792beacc123c9511b24380a3b979e6a6b6c35299f7f:0", "1180cffca82970be702b8f4d1206fb314ffa6b400442eba39e9adbadff4b1de6:0", "066206ab7de04b5d04de7e904cab7b85bf650bcd60ddb7dd1e49b7b57b8ac2e6:0", "db063d622a659d2f888ad24bf189bc06e591664ac235190f908f5ada04ee8be9:0", "1757bb3fd6096499ddb9396268376f14fd8af53c9514fd76762fa50ffbdf7372:0"},
		{"407df68049e7ecdd756f2cf3cda5a8c3852d3adaa561b4cca642462815b8c63c:0", "d88f725e9d1ed9ed40a9a4c4592322132a1879270ac52dee3a8fa01628df9e44:0", "da182df05b296e4683ec9222ca7efa007f3a1b5933a1c61ed3b883dcc8a5919b:0"},
	}
	assert.Equal(t, wantUtxoHistory, blockchain.utxoHistory)
}
