package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hyperledger/fabric/core/ledger/internal/version"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb/stateleveldb"
)

func main() {
	const (
		ns    = "testns"
		key   = "testkey"
		loops = 10
	)
	testValue := []byte("Sensitive data for encryption test")
	putTimes := make([]float64, 0, loops)
	getTimes := make([]float64, 0, loops)

	// Tạo thư mục DB tạm để test
	dbPath := "./testdb"
	os.RemoveAll(dbPath)
	defer os.RemoveAll(dbPath)

	db, err := stateleveldb.NewVersionedDBProvider(dbPath)
	if err != nil {
		panic(err)
	}
	vdb, err := db.GetDBHandle("testdb", nil)
	if err != nil {
		panic(err)
	}

	var ver = &version.Height{BlockNum: 1, TxNum: 1}

	fmt.Println("Unit test cho 2 hàm put và get (ĐÃ TÍCH HỢP MẬT MÃ)")
	fmt.Println("Lần\tHàm put (s)\tHàm get (s)")
	for i := 1; i <= loops; i++ {
		batch := statedb.NewUpdateBatch()
		startPut := time.Now()
		batch.Put(ns, key, testValue, ver)
		if err := vdb.ApplyUpdates(batch, nil); err != nil {
			panic(err)
		}
		putDur := time.Since(startPut).Seconds()
		putTimes = append(putTimes, putDur)

		startGet := time.Now()
		_, err := vdb.GetState(ns, key)
		if err != nil {
			panic(err)
		}
		getDur := time.Since(startGet).Seconds()
		getTimes = append(getTimes, getDur)

		fmt.Printf("%d\t%.6f\t%.6f\n", i, putDur, getDur)
	}

	// Tính trung bình
	var sumPut, sumGet float64
	for i := 0; i < loops; i++ {
		sumPut += putTimes[i]
		sumGet += getTimes[i]
	}
	fmt.Printf("Trung bình\t%.6f\t%.6f\n", sumPut/float64(loops), sumGet/float64(loops))
}
