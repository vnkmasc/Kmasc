#!/bin/bash

# Build MKV Encryption and Create Keys Script
# Author: Phong Ngo
# Date: June 15, 2025

set -e

echo "=== Step 5: Building MKV Encryption and Creating Keys ==="
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
    if [ ! -f "go.mod" ] || [ ! -d "core/ledger/kvledger/txmgmt/statedb" ]; then
        print_status "FAIL" "Not in Fabric root directory with encryption integration"
        print_status "INFO" "Current directory: $(pwd)"
        print_status "INFO" "Please run this script from fabric-3.1.1/ directory"
        exit 1
    fi
    print_status "INFO" "Running from correct directory: $(pwd)"
}

# Save root directory
ROOT_DIR=$(pwd)

# Main execution
main() {
    print_status "INFO" "Step 5: Building encryption and creating keys..."
    
    # Ensure we're in the correct directory
    ensure_correct_directory
    
    # Build MKV encryption library
    print_status "INFO" "Building MKV encryption library..."
    cd core/ledger/kvledger/txmgmt/statedb/mkv
    
    if [ -f "Makefile" ]; then
        make clean && make
        print_status "PASS" "MKV library built successfully"
    else
        print_status "FAIL" "Makefile not found in MKV directory"
        exit 1
    fi
    
    # Return to root directory
    cd "$ROOT_DIR"
    
    # Create MKV keys
    print_status "INFO" "Creating MKV keys..."
    
    # Copy necessary files to root directory
    cp core/ledger/kvledger/txmgmt/statedb/mkv/libmkv.so . 2>/dev/null || {
        print_status "WARN" "Could not copy libmkv.so"
    }
    
    cp core/ledger/kvledger/txmgmt/statedb/mkv/mkv.go . 2>/dev/null || {
        print_status "WARN" "Could not copy mkv.go"
    }
    
    # Initialize keys using key_manager.sh
    if [ -f "core/ledger/kvledger/txmgmt/statedb/mkv/key_manager.sh" ]; then
        echo "fabric_mkv_password_2025" | bash core/ledger/kvledger/txmgmt/statedb/mkv/key_manager.sh init
        print_status "PASS" "MKV keys created successfully"
    else
        print_status "FAIL" "key_manager.sh not found"
        exit 1
    fi
    
    # Verify keys were created
    if [ -f "k1.key" ] && [ -f "k0.key" ] && [ -f "encrypted_k1.key" ]; then
        print_status "PASS" "All MKV keys verified"
        echo "   - Keys location: $(pwd)/"
        echo "   - Files: k1.key, k0.key, encrypted_k1.key"
        echo "   - Password: fabric_mkv_password_2025"
    else
        print_status "FAIL" "Some MKV keys are missing"
        ls -la k*.key encrypted_k*.key 2>/dev/null || true
        exit 1
    fi
    
    print_status "PASS" "Step 5 completed successfully"
    echo
    echo "🔑 MKV Encryption Summary:"
    echo "   - Library: libmkv.so (built)"
    echo "   - Keys: k1.key, k0.key, encrypted_k1.key (created)"
    echo "   - Password: fabric_mkv_password_2025"
    echo "   - Location: $(pwd)/"
}

# Run main function
main "$@" 