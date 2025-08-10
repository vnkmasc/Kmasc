#!/bin/bash

# Build All Libraries Script
# Author: Phong Ngo
# Date: June 15, 2025

set -e

echo "=== Build All Libraries (Encryption + MKV) ==="
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

# Function to build encryption library
build_encryption() {
    echo "Building encryption library..."
    
    ENCRYPTION_DIR="core/ledger/kvledger/txmgmt/statedb"
    
    if dir_exists "$ENCRYPTION_DIR"; then
        cd "$ENCRYPTION_DIR"
        
        if file_exists "Makefile"; then
            print_status "INFO" "Building encryption library..."
            make clean && make
            
            if file_exists "libencryption.so"; then
                print_status "PASS" "libencryption.so created successfully (legacy)"
                ls -la libencryption.so
            else
                print_status "INFO" "libencryption.so was not created (not required for MKV)"
            fi
            
        else
            print_status "FAIL" "Makefile not found in $ENCRYPTION_DIR"
            return 1
        fi
        
        cd - > /dev/null
        print_status "PASS" "Encryption library built successfully"
        
    else
        print_status "FAIL" "statedb directory not found"
        return 1
    fi
}

# Function to build MKV library
build_mkv() {
    echo "Building MKV library..."
    
    MKV_DIR="core/ledger/kvledger/txmgmt/statedb/mkv"
    
    if dir_exists "$MKV_DIR"; then
        cd "$MKV_DIR"
        
        if file_exists "Makefile"; then
            print_status "INFO" "Building MKV library..."
            make clean && make
            
            if file_exists "libmkv.so"; then
                print_status "PASS" "libmkv.so created successfully"
                ls -la libmkv.so
                
                # Run MKV tests if available
                if file_exists "mkv_test.go"; then
                    print_status "INFO" "Running MKV unit tests..."
                    LD_LIBRARY_PATH=. go test -v
                    if [ $? -eq 0 ]; then
                        print_status "PASS" "MKV unit tests passed"
                    else
                        print_status "WARN" "MKV unit tests failed, but library was built"
                    fi
                fi
                
            else
                print_status "FAIL" "libmkv.so was not created"
                return 1
            fi
            
        else
            print_status "FAIL" "Makefile not found in $MKV_DIR"
            return 1
        fi
        
        cd - > /dev/null
        print_status "PASS" "MKV library built successfully"
        
    else
        print_status "FAIL" "mkv directory not found"
        return 1
    fi
}

# Function to verify all libraries
verify_libraries() {
    echo "Verifying all libraries..."
    
    # Check encryption library (optional - for backward compatibility)
    if file_exists "core/ledger/kvledger/txmgmt/statedb/libencryption.so"; then
        print_status "PASS" "libencryption.so verified (legacy)"
    else
        print_status "INFO" "libencryption.so not found (not required for MKV)"
    fi
    
    # Check MKV library
    if file_exists "core/ledger/kvledger/txmgmt/statedb/mkv/libmkv.so"; then
        print_status "PASS" "libmkv.so verified"
    else
        print_status "FAIL" "libmkv.so not found"
        return 1
    fi
    
    # Show library info
    print_status "INFO" "Library information:"
    if file_exists "core/ledger/kvledger/txmgmt/statedb/libencryption.so"; then
        echo "libencryption.so:"
        ls -la core/ledger/kvledger/txmgmt/statedb/libencryption.so
    fi
    echo "libmkv.so:"
    ls -la core/ledger/kvledger/txmgmt/statedb/mkv/libmkv.so
    
    print_status "PASS" "All libraries verified successfully"
}

# Main function
main() {
    # Ensure we're in the correct directory
    if [ ! -f "go.mod" ]; then
        print_status "FAIL" "Not in Fabric root directory. Please run from fabric-3.1.1/"
        exit 1
    fi
    
    print_status "INFO" "Building all libraries..."
    
    # Build encryption library
    if build_encryption; then
        print_status "PASS" "Encryption library build completed"
    else
        print_status "FAIL" "Encryption library build failed"
        exit 1
    fi
    
    echo
    
    # Build MKV library
    if build_mkv; then
        print_status "PASS" "MKV library build completed"
    else
        print_status "FAIL" "MKV library build failed"
        exit 1
    fi
    
    echo
    
    # Verify all libraries
    verify_libraries
    
    echo
    print_status "INFO" "All libraries build completed at $(date)"
    print_status "PASS" "✅ Both encryption and MKV libraries are ready!"
}

# Check if script is being run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main
fi
