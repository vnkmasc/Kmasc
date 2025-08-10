#!/bin/bash

# MKV Key Initialization Script with PBKDF2
# This script initializes the MKV encryption system with a password

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "INFO")
            echo -e "${BLUE}[INFO]${NC} $message"
            ;;
        "SUCCESS")
            echo -e "${GREEN}[SUCCESS]${NC} $message"
            ;;
        "WARN")
            echo -e "${YELLOW}[WARN]${NC} $message"
            ;;
        "ERROR")
            echo -e "${RED}[ERROR]${NC} $message"
            ;;
    esac
}

# Check if password is provided
if [ $# -eq 0 ]; then
    echo -e "${YELLOW}Usage: $0 <password>${NC}"
    echo "Example: $0 fabric123"
    exit 1
fi

password=$1

print_status "INFO" "Initializing MKV Key Management System with PBKDF2..."

# Check if required files exist
if [ ! -f "mkv.go" ]; then
    print_status "ERROR" "mkv.go not found in current directory"
    exit 1
fi

if [ ! -f "libmkv.so" ]; then
    print_status "ERROR" "libmkv.so not found. Please build the MKV library first."
    exit 1
fi

# Create a temporary Go file to call the function in /tmp to avoid package conflicts
cat > /tmp/mkv_init.go << EOF
package main
import (
    "fmt"
    "os"
    "github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb/mkv"
)
func main() {
    err := mkv.InitializeKeyManagement("$password")
    if err != nil {
        fmt.Printf("❌ Failed to initialize: %v\n", err)
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

# Run the initialization from current directory but with temp file in /tmp
LD_LIBRARY_PATH=. go run /tmp/mkv_init.go

# Clean up
rm -f /tmp/mkv_init.go

print_status "SUCCESS" "MKV Key Management System initialized successfully!"
print_status "INFO" "Password used: $password"
print_status "INFO" "Check the generated key files in the current directory."
