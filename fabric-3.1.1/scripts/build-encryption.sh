#!/bin/bash

# Build Encryption Libraries and Create MKV Keys
# This script builds the encryption libraries and initializes the MKV key management system

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
        "PASS")
            echo -e "${GREEN}✅ PASS:${NC} $message"
            ;;
        "FAIL")
            echo -e "${RED}❌ FAIL:${NC} $message"
            ;;
        "INFO")
            echo -e "${BLUE}ℹ️  INFO:${NC} $message"
            ;;
        "WARN")
            echo -e "${YELLOW}⚠️  WARN:${NC} $message"
            ;;
    esac
}

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "🔐 Building MKV Encryption System..."
echo "   - Script directory: $SCRIPT_DIR"
echo "   - Fabric root: $ROOT_DIR"

# Step 1: Build encryption library
print_status "INFO" "Step 1: Building encryption library..."
cd "$ROOT_DIR/core/ledger/kvledger/txmgmt/statedb"

if [ -f "Makefile" ]; then
    make clean 2>/dev/null || true
    make
    print_status "PASS" "Encryption library built successfully"
else
    print_status "FAIL" "Makefile not found in encryption directory"
    exit 1
fi

# Step 2: Build MKV library
print_status "INFO" "Step 2: Building MKV library..."
cd "$ROOT_DIR/core/ledger/kvledger/txmgmt/statedb/mkv"

if [ -f "Makefile" ]; then
    make clean 2>/dev/null || true
    make
    print_status "PASS" "MKV library built successfully"
else
    print_status "FAIL" "Makefile not found in MKV directory"
    exit 1
fi

# Step 3: Verify libraries were built
print_status "INFO" "Step 3: Verifying built libraries..."

cd "$ROOT_DIR"

# Check encryption library
if [ -f "core/ledger/kvledger/txmgmt/statedb/libencryption.so" ]; then
    print_status "PASS" "libencryption.so verified"
    ls -la core/ledger/kvledger/txmgmt/statedb/libencryption.so
else
    print_status "FAIL" "libencryption.so not found"
    exit 1
fi

# Check MKV library
if [ -f "core/ledger/kvledger/txmgmt/statedb/mkv/libmkv.so" ]; then
    print_status "PASS" "libmkv.so verified"
    ls -la core/ledger/kvledger/txmgmt/statedb/mkv/libmkv.so
else
    print_status "FAIL" "libmkv.so not found"
    exit 1
fi

# Step 4: Create MKV keys
print_status "INFO" "Step 4: Creating MKV keys..."

# Files remain in their proper location (no copying needed)
print_status "INFO" "Library files kept in proper location: $MKV_DIR/"

# mkv.go remains in proper location within mkv directory

    # Initialize keys using key_manager.sh in proper directory
    MKV_DIR="core/ledger/kvledger/txmgmt/statedb/mkv"
    if [ -f "$MKV_DIR/key_manager.sh" ]; then
        # Change to MKV directory to create keys in the right location
        cd "$MKV_DIR"
        ./key_manager.sh init fabric123
        cd - > /dev/null # Go back to previous directory
        print_status "PASS" "MKV keys created successfully"
    else
        print_status "FAIL" "key_manager.sh not found"
        exit 1
    fi

# Step 5: Verify keys were created
print_status "INFO" "Step 5: Verifying MKV keys..."

    # Verify keys were created in MKV directory
    if [ -f "$MKV_DIR/k1.key" ] && [ -f "$MKV_DIR/k0.key" ] && [ -f "$MKV_DIR/encrypted_k1.key" ]; then
        print_status "PASS" "All MKV keys verified"
        echo "   - Keys location: $(pwd)/$MKV_DIR/"
        echo "   - Files: k1.key, k0.key, encrypted_k1.key"
        echo "   - Password: fabric123"
    else
        print_status "FAIL" "Some MKV keys are missing"
        ls -la "$MKV_DIR"/k*.key "$MKV_DIR"/encrypted_k*.key 2>/dev/null || true
        exit 1
    fi

print_status "PASS" "Step 5 completed successfully"

echo ""
echo "🔑 MKV Encryption Summary:"
echo "   - Library: libmkv.so (built)"
echo "   - Keys: k1.key, k0.key, encrypted_k1.key (created)"
echo "   - Password: fabric123"
echo "   - Location: $(pwd)/$MKV_DIR/"
echo ""
echo "✅ MKV Encryption System built and initialized successfully!"
echo "   - Libraries are ready for use"
echo "   - Keys are generated and secured"
echo "   - System is ready for Fabric integration"
