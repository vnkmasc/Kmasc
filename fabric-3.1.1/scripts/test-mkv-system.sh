#!/bin/bash

# MKV System Test Script
# Author: Phong Ngo
# Date: June 15, 2025

set -e

echo "=== Testing MKV System ==="
echo "This script will test the MKV encryption system"
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

# Function to ensure we're in the correct directory
ensure_correct_directory() {
    if [ ! -f "go.mod" ] || [ ! -d "core/ledger/kvledger/txmgmt/statedb/mkv" ]; then
        print_status "FAIL" "Not in Fabric root directory with MKV integration"
        print_status "INFO" "Current directory: $(pwd)"
        print_status "INFO" "Please run this script from fabric-3.1.1/ directory"
        exit 1
    fi
    print_status "INFO" "Running from correct directory: $(pwd)"
}

# Save root directory
ROOT_DIR=$(pwd)

# Function to check key files
check_key_files() {
    print_status "INFO" "Checking MKV key files..."
    
    local required_files=(
        "core/ledger/kvledger/txmgmt/statedb/mkv/k1.key"
        "core/ledger/kvledger/txmgmt/statedb/mkv/encrypted_k1.key"
        "core/ledger/kvledger/txmgmt/statedb/mkv/k0_salt.key"
        "core/ledger/kvledger/txmgmt/statedb/mkv/password.txt"
    )
    
    local missing_files=()
    
    for file in "${required_files[@]}"; do
        if [ ! -f "$file" ]; then
            missing_files+=("$file")
        else
            print_status "PASS" "Found: $file"
        fi
    done
    
    if [ ${#missing_files[@]} -eq 0 ]; then
        print_status "PASS" "All required key files are present"
    else
        print_status "FAIL" "Missing required key files:"
        for file in "${missing_files[@]}"; do
            echo "  - $file"
        done
        exit 1
    fi
}

# Function to check library files
check_library_files() {
    print_status "INFO" "Checking MKV library files..."
    
    local required_files=(
        "core/ledger/kvledger/txmgmt/statedb/mkv/libmkv.so"
        "core/ledger/kvledger/txmgmt/statedb/mkv/mkv.go"
        "core/ledger/kvledger/txmgmt/statedb/mkv/mkv.h"
        "core/ledger/kvledger/txmgmt/statedb/mkv/MKV256.c"
        "core/ledger/kvledger/txmgmt/statedb/mkv/MKV256.h"
    )
    
    local missing_files=()
    
    for file in "${required_files[@]}"; do
        if [ ! -f "$file" ]; then
            missing_files+=("$file")
        else
            print_status "PASS" "Found: $file"
        fi
    done
    
    if [ ${#missing_files[@]} -eq 0 ]; then
        print_status "PASS" "All required library files are present"
    else
        print_status "FAIL" "Missing required library files:"
        for file in "${missing_files[@]}"; do
            echo "  - $file"
        done
        exit 1
    fi
}

# Function to run Go tests
run_go_tests() {
    print_status "INFO" "Running Go tests..."
    
    cd core/ledger/kvledger/txmgmt/statedb/mkv
    
    # Run basic tests
    if go test -v -run TestEncryptDecryptValueMKV; then
        print_status "PASS" "Basic encryption/decryption test passed"
    else
        print_status "FAIL" "Basic encryption/decryption test failed"
        exit 1
    fi
    
    # Run PBKDF2 tests
    if go test -v -run TestPBKDF2Implementation; then
        print_status "PASS" "PBKDF2 implementation test passed"
    else
        print_status "FAIL" "PBKDF2 implementation test failed"
        exit 1
    fi
    
    # Run key management tests
    if go test -v -run TestKeyManagementSystem; then
        print_status "PASS" "Key management system test passed"
    else
        print_status "FAIL" "Key management system test failed"
        exit 1
    fi
    
    # Return to root directory
    cd "$ROOT_DIR"
}

# Function to test integration with LevelDB
test_leveldb_integration() {
    print_status "INFO" "Testing LevelDB integration..."
    
    # Check if value_encoding.go has been updated
    if grep -q "EncryptValueMKV(v.Value)" core/ledger/kvledger/txmgmt/statedb/stateleveldb/value_encoding.go; then
        print_status "PASS" "LevelDB integration code updated"
    else
        print_status "FAIL" "LevelDB integration code not updated"
        exit 1
    fi
}



# Main execution
main() {
    echo "Starting MKV system test..."
    echo
    
    ensure_correct_directory
    check_key_files
    check_library_files
    run_go_tests
    test_leveldb_integration
    
    echo
    print_status "PASS" "MKV system test completed successfully!"
    echo
    echo "All tests passed:"
    echo "  ✅ Key files are present"
    echo "  ✅ Library files are present"
    echo "  ✅ Go tests passed"
    echo "  ✅ LevelDB integration updated"
    echo
    echo "The MKV system is ready for use!"
}

# Execute main function directly
echo "This script will:"
echo "1. Check MKV key files"
echo "2. Check MKV library files"
echo "3. Run Go tests"
echo "4. Test LevelDB integration"
echo
echo "Starting execution..."
main "$@"
