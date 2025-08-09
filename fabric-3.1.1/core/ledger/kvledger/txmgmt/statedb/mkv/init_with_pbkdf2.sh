#!/bin/bash

# Initialize MKV system with PBKDF2
# Usage: ./init_with_pbkdf2.sh [password]

if [ -z "$1" ]; then
    read -s -p "Enter password for initialization: " password
    echo
else
    password="$1"
fi

echo "Initializing MKV system with PBKDF2..."

# Create a temporary Go file to call the function
cat > temp_init.go << EOF
package main

import (
    "fmt"
    "os"
    "github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb/mkv"
)

func main() {
    err := mkv.InitializeKeyManagement("$password")
    if err != nil {
        fmt.Printf("❌ Failed to initialize: %%v\n", err)
        os.Exit(1)
    }
    fmt.Println("✅ System initialized successfully with PBKDF2!")
    fmt.Println("📁 Files created:")
    fmt.Println("   - k1.key (32 bytes random K1)")
    fmt.Println("   - k0_salt.key (32 bytes random salt)")
    fmt.Println("   - k0.key (32 bytes K0 from PBKDF2)")
    fmt.Println("   - encrypted_k1.key (K1 encrypted with K0)")
}
EOF

# Run the initialization
LD_LIBRARY_PATH=. go run temp_init.go

# Clean up
rm -f temp_init.go 