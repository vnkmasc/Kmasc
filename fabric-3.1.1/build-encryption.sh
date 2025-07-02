#!/bin/bash

# Build Encryption Library Script
# Author: Phong Ngo
# Date: June 15, 2025

set -e

echo "=== Build Encryption Library ==="
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

# Main function to build encryption library
build_encryption() {
    echo "Step 5: Building encryption library..."
    
    # Save current directory
    ROOT_DIR=$(pwd)
    
    # Ensure we're in the correct directory
    if [ ! -f "go.mod" ]; then
        print_status "FAIL" "Not in Fabric root directory. Please run from fabric-3.1.1/"
        exit 1
    fi
    
    ENCRYPTION_DIR="core/ledger/kvledger/txmgmt/statedb"
    
    if dir_exists "$ENCRYPTION_DIR"; then
        print_status "INFO" "Building encryption library..."
        cd "$ROOT_DIR/$ENCRYPTION_DIR"
        
        # Check if Makefile exists
        if file_exists "Makefile"; then
            print_status "INFO" "Found Makefile, running make clean && make..."
            make clean && make
            
            # Check if library was created
            if file_exists "libencryption.so"; then
                print_status "PASS" "libencryption.so created successfully"
                ls -la libencryption.so
                
                # Check library dependencies
                if command -v ldd >/dev/null 2>&1; then
                    print_status "INFO" "Checking library dependencies..."
                    ldd libencryption.so
                fi
                
            else
                print_status "FAIL" "libencryption.so was not created"
                exit 1
            fi
            
        else
            print_status "FAIL" "Makefile not found in $ENCRYPTION_DIR"
            print_status "INFO" "Available files:"
            ls -la
            exit 1
        fi
        
        cd "$ROOT_DIR"
        print_status "PASS" "Encryption library built successfully"
        
    else
        print_status "FAIL" "statedb directory not found at $ENCRYPTION_DIR"
        print_status "INFO" "Current directory: $(pwd)"
        print_status "INFO" "Available directories in core/ledger/kvledger/txmgmt/:"
        ls -la core/ledger/kvledger/txmgmt/ 2>/dev/null || echo "Directory not accessible"
        exit 1
    fi
    
    echo
    print_status "INFO" "Encryption library build completed at $(date)"
}

# Check if script is being run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    build_encryption
fi 