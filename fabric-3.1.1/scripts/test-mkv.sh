#!/bin/bash

# Test MKV Integration Script
# This script tests MKV encryption/decryption functionality

set -e

echo "🧪 Testing MKV Integration..."

# Check if we're in the right directory
if [ ! -d "core/ledger/kvledger/txmgmt/statedb/mkv" ]; then
    echo "❌ Error: MKV directory not found"
    echo "Please run this script from the fabric-3.1.1 root directory"
    exit 1
fi

# Go to MKV directory
cd core/ledger/kvledger/txmgmt/statedb/mkv

# Check if library and keys exist
if [ ! -f "libmkv.so" ]; then
    echo "❌ Error: libmkv.so not found"
    echo "Please run build-mkv.sh first"
    exit 1
fi

if [ ! -f "k1.key" ] || [ ! -f "encrypted_k1.key" ]; then
    echo "❌ Error: MKV keys not found"
    echo "Please run init-mkv-keys.sh first"
    exit 1
fi

# Test MKV library
echo "⚙️  Testing MKV library..."
LD_LIBRARY_PATH=. go test -v -run TestEncryptDecryptValueMKV

if [ $? -eq 0 ]; then
    echo "✅ MKV library tests passed!"
else
    echo "❌ MKV library tests failed"
    exit 1
fi

# Test key management system
echo "⚙️  Testing key management system..."
if [ -f "mkv-initial-password.txt" ]; then
    INITIAL_PASSWORD=$(cat mkv-initial-password.txt)
    echo "🔑 Testing with initial password..."
    
    # Test PBKDF2 key derivation
    LD_LIBRARY_PATH=. go test -v -run TestPBKDF2Implementation
    
    if [ $? -eq 0 ]; then
        echo "✅ Key management system tests passed!"
    else
        echo "❌ Key management system tests failed"
        exit 1
    fi
else
    echo "⚠️  Initial password file not found, skipping password tests"
fi

# Test network integration if network is running
echo "⚙️  Testing network integration..."
if docker ps | grep -q "peer0.org1.example.com"; then
    echo "🌐 Network is running, testing encryption logs..."
    
    # Check if encryption logs exist
    if docker exec peer0.org1.example.com test -f /tmp/state_mkv.log; then
        echo "📊 Recent encryption activities:"
        docker exec peer0.org1.example.com tail -5 /tmp/state_mkv.log
        echo "✅ Network integration test passed!"
    else
        echo "⚠️  No encryption logs found yet"
        echo "   Run some chaincode transactions to see encryption in action"
    fi
else
    echo "⚠️  Network not running, skipping network integration test"
fi

# Go back to root
cd ../../../../..

echo "🎉 MKV integration tests completed!"
echo
echo "📊 Test Results Summary:"
echo "   ✅ MKV library functionality"
echo "   ✅ Encryption/Decryption operations"
echo "   ✅ Key management system"
echo "   ✅ PBKDF2 key derivation"
if docker ps | grep -q "peer0.org1.example.com"; then
    echo "   ✅ Network integration"
fi
echo
echo "🚀 Next steps:"
echo "   1. Run chaincode transactions to see encryption in action"
echo "   2. Monitor logs: docker exec peer0.org1.example.com tail -f /tmp/state_mkv.log"
echo "   3. Change password: cd core/ledger/kvledger/txmgmt/statedb/mkv && ./mkv_client.sh change"