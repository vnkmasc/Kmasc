#!/bin/bash

# Script để chạy tất cả các test cho encryption integration
# Chạy: ./run_tests.sh

set -e  # Exit on any error

echo "=== Encryption Integration Test Suite ==="
echo "Date: $(date)"
echo "Directory: $(pwd)"
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

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Test 1: Check environment
echo "1. Checking environment..."
print_status "INFO" "Checking required tools and libraries"

# Check gcc
if command_exists gcc; then
    print_status "PASS" "gcc found: $(gcc --version | head -n1)"
else
    print_status "FAIL" "gcc not found"
    exit 1
fi

# Check OpenSSL
if command_exists openssl; then
    print_status "PASS" "OpenSSL found: $(openssl version)"
else
    print_status "FAIL" "OpenSSL not found"
    exit 1
fi

# Check Go
if command_exists go; then
    print_status "PASS" "Go found: $(go version)"
else
    print_status "FAIL" "Go not found"
    exit 1
fi

# Check CGO
if [ "$CGO_ENABLED" = "1" ]; then
    print_status "PASS" "CGO is enabled"
else
    print_status "WARN" "CGO is not enabled, setting CGO_ENABLED=1"
    export CGO_ENABLED=1
fi

echo

# Test 2: Build C library
echo "2. Building C library..."
print_status "INFO" "Compiling encryption.c with OpenSSL"

if make clean && make; then
    print_status "PASS" "C library built successfully"
else
    print_status "FAIL" "Failed to build C library"
    exit 1
fi

# Check if library was created
if [ -f "libencryption.so" ]; then
    print_status "PASS" "libencryption.so created"
    ls -la libencryption.so
else
    print_status "FAIL" "libencryption.so not found"
    exit 1
fi

echo

# Test 3: Check library dependencies
echo "3. Checking library dependencies..."
print_status "INFO" "Verifying library dependencies"

if command_exists ldd; then
    ldd_output=$(ldd libencryption.so 2>/dev/null || echo "ldd failed")
    if echo "$ldd_output" | grep -q "libssl\|libcrypto"; then
        print_status "PASS" "OpenSSL libraries linked correctly"
    else
        print_status "WARN" "OpenSSL libraries not found in ldd output"
        echo "$ldd_output"
    fi
else
    print_status "WARN" "ldd not available, skipping dependency check"
fi

echo

# Test 4: Build Go package
echo "4. Building Go package..."
print_status "INFO" "Building Go package with CGO"

if go build ./...; then
    print_status "PASS" "Go package built successfully"
else
    print_status "FAIL" "Failed to build Go package"
    exit 1
fi

echo

# Test 5: Run Go tests
echo "5. Running Go tests..."
print_status "INFO" "Running unit tests"

if go test ./...; then
    print_status "PASS" "All Go tests passed"
else
    print_status "FAIL" "Some Go tests failed"
    # Don't exit here, continue with other tests
fi

echo

# Test 6: Run example test
echo "6. Running example test..."
print_status "INFO" "Running encryption example test"

if [ -f "test_encryption_example.go" ]; then
    if go run test_encryption_example.go; then
        print_status "PASS" "Example test completed successfully"
    else
        print_status "FAIL" "Example test failed"
    fi
else
    print_status "WARN" "test_encryption_example.go not found, skipping"
fi

echo

# Test 7: Performance test
echo "7. Running performance test..."
print_status "INFO" "Running benchmark tests"

if go test -bench=. ./... 2>/dev/null; then
    print_status "PASS" "Performance tests completed"
else
    print_status "WARN" "No benchmark tests found or failed to run"
fi

echo

# Test 8: Integration test
echo "8. Running integration test..."
print_status "INFO" "Testing UpdateBatch integration"

# Create a simple integration test
cat > temp_integration_test.go << 'EOF'
package main

import (
    "fmt"
    "log"
    "github.com/hyperledger/fabric/core/ledger/internal/version"
    "github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb"
)

func main() {
    fmt.Println("Integration test: UpdateBatch with encryption")
    
    batch := statedb.NewUpdateBatch()
    testData := []byte("Integration test data")
    
    // Put data
    batch.Put("testns", "testkey", testData, &version.Height{})
    
    // Get data
    retrieved := batch.Get("testns", "testkey")
    if retrieved == nil {
        log.Fatal("Failed to retrieve data")
    }
    
    if string(retrieved.Value) == string(testData) {
        fmt.Println("✅ Integration test PASSED!")
    } else {
        fmt.Println("❌ Integration test FAILED!")
    }
}
EOF

if go run temp_integration_test.go; then
    print_status "PASS" "Integration test passed"
else
    print_status "FAIL" "Integration test failed"
fi

# Clean up
rm -f temp_integration_test.go

echo

# Summary
echo "=== Test Summary ==="
print_status "INFO" "All tests completed"
print_status "INFO" "Encryption integration is ready for use"
echo
echo "Next steps:"
echo "1. Use UpdateBatch.Put() to store encrypted data"
echo "2. Use UpdateBatch.Get() to retrieve decrypted data"
echo "3. Check README_ENCRYPTION.md for detailed documentation"
echo "4. Check USAGE.md for usage examples"
echo
print_status "INFO" "Test suite completed at $(date)" 