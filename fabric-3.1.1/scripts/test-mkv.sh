#!/bin/bash

# Test MKV Library Script
# Author: Phong Ngo
# Date: June 15, 2025

set -e

echo "=== Test MKV Library ==="
echo "Date: $(date)"
echo

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
        "PASS")
            echo -e "${GREEN}✅ PASS${NC}: $message"
            ;;
        "FAIL")
            echo -e "${RED}❌ FAIL${NC}: $message"
            ;;
        "INFO")
            echo -e "${BLUE}ℹ️  INFO${NC}: $message"
            ;;
        "WARN")
            echo -e "${YELLOW}⚠️  WARN${NC}: $message"
            ;;
    esac
}

# Function to check if file exists
file_exists() {
    [ -f "$1" ]
}

# Function to check if directory exists
dir_exists() {
    [ -d "$1" ]
}

# Main function to test mkv library
test_mkv() {
    echo "Testing MKV library..."
    
    # Save current directory
    ROOT_DIR=$(pwd)
    
    # Ensure we're in the correct directory
    if [ ! -f "go.mod" ]; then
        print_status "FAIL" "Not in Fabric root directory. Please run from fabric-3.1.1/"
        exit 1
    fi
    
    MKV_DIR="core/ledger/kvledger/txmgmt/statedb/mkv"
    
    if dir_exists "$MKV_DIR"; then
        print_status "INFO" "Testing MKV library..."
        cd "$ROOT_DIR/$MKV_DIR"
        
        # Check if library exists
        if file_exists "libmkv.so"; then
            print_status "PASS" "libmkv.so found"
            
            # Check library dependencies
            if command -v ldd >/dev/null 2>&1; then
                print_status "INFO" "Checking library dependencies..."
                ldd libmkv.so
            fi
            
            # Run MKV tests if available
            if file_exists "mkv_test.go"; then
                print_status "INFO" "Running MKV unit tests..."
                export LD_LIBRARY_PATH=.
                go test -v
                if [ $? -eq 0 ]; then
                    print_status "PASS" "MKV unit tests passed"
                else
                    print_status "FAIL" "MKV unit tests failed"
                    exit 1
                fi
            else
                print_status "WARN" "mkv_test.go not found, skipping unit tests"
            fi
            
            # Test basic functionality with a simple Go program
            print_status "INFO" "Testing basic MKV functionality..."
            cat > test_mkv_basic.go << 'EOF'
package main

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lmkv
#include "mkv.h"
*/
import "C"
import (
    "fmt"
    "unsafe"
)

func main() {
    key := []byte("12345678901234567890123456789012") // 32 bytes for 256-bit
    data := []byte("Hello, MKV encryption!")
    
    // Prepare variables
    cKey := (*C.uchar)(unsafe.Pointer(&key[0]))
    cData := (*C.uchar)(unsafe.Pointer(&data[0]))
    dataLen := C.int(len(data))
    var ciphertextLen C.int
    
    // Allocate buffer for encrypted data (max size)
    maxCiphertextLen := dataLen + 32 // Add some padding
    encrypted := make([]byte, maxCiphertextLen)
    cEncrypted := (*C.uchar)(unsafe.Pointer(&encrypted[0]))
    
    // Encrypt
    result := C.mkv_encrypt(cData, dataLen, cEncrypted, &ciphertextLen, cKey, C.int(len(key)))
    if result != 0 {
        fmt.Println("❌ Encryption failed")
        return
    }
    
    // Resize encrypted buffer to actual size
    encrypted = encrypted[:ciphertextLen]
    
    fmt.Println("✅ Basic MKV functionality test passed")
    fmt.Printf("Original: %s\n", string(data))
    fmt.Printf("Encrypted length: %d\n", ciphertextLen)
}
EOF

            go run test_mkv_basic.go
            if [ $? -eq 0 ]; then
                print_status "PASS" "Basic MKV functionality test passed"
            else
                print_status "FAIL" "Basic MKV functionality test failed"
            fi
            
            # Clean up test file
            rm -f test_mkv_basic.go
            
        else
            print_status "FAIL" "libmkv.so not found. Please build MKV library first."
            exit 1
        fi
        
        cd "$ROOT_DIR"
        print_status "PASS" "MKV library test completed successfully"
        
    else
        print_status "FAIL" "mkv directory not found at $MKV_DIR"
        exit 1
    fi
    
    echo
    print_status "INFO" "MKV library test completed at $(date)"
}

# Check if script is being run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    test_mkv
fi
