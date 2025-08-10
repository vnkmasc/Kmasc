#!/bin/bash

# Step 6: Initialize MKV Keys
# This script initializes the MKV encryption keys

# Source common functions
if [ -f "scripts/functions.sh" ]; then
    source scripts/functions.sh
else
    # Fallback function if functions.sh is not available
    print_status() {
        local level=$1
        local message=$2
        case $level in
            "INFO") echo "ℹ️  $message" ;;
            "WARN") echo "⚠️  $message" ;;
            "ERROR") echo "❌ $message" ;;
            "PASS") echo "✅ $message" ;;
            *) echo "$message" ;;
        esac
    }
fi

step6_init_mkv_keys() {
    print_status "INFO" "Step 6: Initializing MKV encryption keys..."
    
    echo "Initializing MKV encryption keys..."
    cd core/ledger/kvledger/txmgmt/statedb/mkv
    
    # Check if key_manager.sh exists
    if [ ! -f "key_manager.sh" ]; then
        print_status "ERROR" "key_manager.sh not found"
        exit 1
    fi
    
    # Make it executable
    chmod +x key_manager.sh
    
    # Initialize keys with default password
    echo "Initializing keys with default password 'fabric123'..."
    ./key_manager.sh init fabric123
    
    # Check if keys were created
    if [ -f "k1.key" ] && [ -f "k0.key" ] && [ -f "encrypted_k1.key" ]; then
        print_status "PASS" "MKV keys initialized successfully"
        echo "   - K1 (Data Key): k1.key"
        echo "   - K0 (Derived Key): k0.key"
        echo "   - Encrypted K1: encrypted_k1.key"
    else
        print_status "ERROR" "Failed to initialize MKV keys"
        exit 1
    fi
    
    cd /home/phongnh/go-src/Kmasc/fabric-3.1.1
    
    echo "✅ PASS: MKV encryption keys are ready"
    echo "   - Keys initialized and encrypted"
    echo "   - Ready for Fabric encryption"
}

# Main execution
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    step6_init_mkv_keys
fi
