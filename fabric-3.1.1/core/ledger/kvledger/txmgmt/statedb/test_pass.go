package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb/mkv"
)

func main() {
	k1, err := mkv.GetCurrentK1("fabric_mkv_password_2025")
	if err != nil {
		fmt.Printf("❌ Password test FAILED: %v\n", err)
		return
	}
	fmt.Printf("✅ Password test SUCCESS!\n")
	fmt.Printf("K1 (hex): %x\n", k1)
}
