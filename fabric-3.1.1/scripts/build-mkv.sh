#!/bin/bash

# Build MKV Library Script
# This script builds the MKV256 encryption library

set -e

echo "🔧 Building MKV256 Encryption Library..."

# Check if we're in the right directory
if [ ! -d "core/ledger/kvledger/txmgmt/statedb/mkv" ]; then
    echo "❌ Error: MKV directory not found"
    echo "Please run this script from the fabric-3.1.1 root directory"
        exit 1
    fi
    
# Go to MKV directory
cd core/ledger/kvledger/txmgmt/statedb/mkv

# Clean previous builds
echo "🧹 Cleaning previous builds..."
make clean

# Build MKV library
echo "⚙️  Building MKV library..."
make

                    if [ $? -eq 0 ]; then
    echo "✅ MKV library built successfully!"
    echo "📦 Generated files:"
    echo "   - libmkv.so (MKV encryption library)"
    echo "   - mkv.o (Object file)"
    echo "   - MKV256.o (Algorithm object file)"
else
    echo "❌ Failed to build MKV library"
                exit 1
            fi
            
# Go back to root
cd ../../../../..

echo "🎉 MKV library build completed!"