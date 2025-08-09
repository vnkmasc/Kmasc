#!/bin/bash

# Initialize MKV Keys Script
# This script initializes the MKV key management system

set -e

echo "🔐 Initializing MKV Key Management System..."

# Check if we're in the right directory
if [ ! -d "core/ledger/kvledger/txmgmt/statedb/mkv" ]; then
    echo "❌ Error: MKV directory not found"
    echo "Please run this script from the fabric-3.1.1 root directory"
    exit 1
fi

# Go to MKV directory
cd core/ledger/kvledger/txmgmt/statedb/mkv

# Check if library exists
if [ ! -f "libmkv.so" ]; then
    echo "❌ Error: libmkv.so not found"
    echo "Please run build-mkv.sh first"
    exit 1
fi

# Generate default password if not provided
if [ -z "$1" ]; then
    DEFAULT_PASSWORD="fabric_mkv_$(date +%Y%m%d)_$(openssl rand -hex 8)"
    echo "🔑 Generated default password: $DEFAULT_PASSWORD"
else
    DEFAULT_PASSWORD="$1"
    echo "🔑 Using provided password"
fi

# Initialize MKV system
echo "⚙️  Initializing MKV system with PBKDF2..."
echo "$DEFAULT_PASSWORD" | bash init_with_pbkdf2.sh

if [ $? -eq 0 ]; then
    echo "✅ MKV system initialized successfully!"
    echo "📁 Generated key files:"
    echo "   - k1.key (32 bytes random K1)"
    echo "   - k0_salt.key (32 bytes random salt)"
    echo "   - k0.key (32 bytes K0 from PBKDF2)"
    echo "   - encrypted_k1.key (K1 encrypted with K0)"
    
    # Save password for reference
    echo "$DEFAULT_PASSWORD" > mkv-initial-password.txt
    chmod 600 mkv-initial-password.txt
    echo "🔐 Initial password saved to mkv-initial-password.txt"
    
else
    echo "❌ Failed to initialize MKV system"
    exit 1
fi

# Go back to root
cd ../../../../..

echo "🎉 MKV key initialization completed!"
echo
echo "⚠️  IMPORTANT SECURITY NOTES:"
echo "   1. Change the default password immediately in production"
echo "   2. Store the password securely"
echo "   3. Regular password rotation is recommended"
echo "   4. Backup key files securely"