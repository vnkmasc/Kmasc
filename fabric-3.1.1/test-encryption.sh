#!/bin/bash

# Test Encryption Integration Script
# Author: Phong Ngo
# Date: June 15, 2025

set -e

echo "=== Test Encryption Integration ==="
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

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Main function to test encryption integration
test_encryption() {
    echo "Step 6: Testing encryption integration..."
    
    # Save current directory
    ROOT_DIR=$(pwd)
    
    # Ensure we're in the correct directory
    if [ ! -f "go.mod" ]; then
        print_status "FAIL" "Not in Fabric root directory. Please run from fabric-3.1.1/"
        exit 1
    fi
    
    ENCRYPTION_DIR="core/ledger/kvledger/txmgmt/statedb"
    
    if dir_exists "$ENCRYPTION_DIR"; then
        print_status "INFO" "Running encryption tests..."
        cd "$ROOT_DIR/$ENCRYPTION_DIR"
        
        # Check if run_tests.sh exists
        if file_exists "run_tests.sh"; then
            print_status "INFO" "Found run_tests.sh, running it..."
            chmod +x run_tests.sh
            bash run_tests.sh
            print_status "PASS" "Encryption tests completed"
        else
            print_status "WARN" "run_tests.sh not found, running basic Go tests..."
            
            # Check if Go tests exist
            if file_exists "*.go" ]; then
                print_status "INFO" "Running Go tests..."
                go test ./...
                print_status "PASS" "Basic Go tests completed"
            else
                print_status "WARN" "No Go test files found"
            fi
            
            # Check if library exists
            if file_exists "libencryption.so"; then
                print_status "PASS" "libencryption.so exists"
                
                # Test library loading
                if command -v ldd >/dev/null 2>&1; then
                    print_status "INFO" "Testing library dependencies..."
                    ldd libencryption.so
                fi
                
            else
                print_status "FAIL" "libencryption.so not found"
                print_status "INFO" "Run ./build-encryption.sh first"
                exit 1
            fi
        fi
        
        cd "$ROOT_DIR"
        
    else
        print_status "FAIL" "statedb directory not found at $ENCRYPTION_DIR"
        print_status "INFO" "Current directory: $(pwd)"
        print_status "INFO" "Available directories in core/ledger/kvledger/txmgmt/:"
        ls -la core/ledger/kvledger/txmgmt/ 2>/dev/null || echo "Directory not accessible"
        exit 1
    fi
    
    echo
    print_status "INFO" "Encryption integration test completed at $(date)"
}

# Check if script is being run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    test_encryption
fi 