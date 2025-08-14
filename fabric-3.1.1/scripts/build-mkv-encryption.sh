#!/bin/bash

# MKV Encryption Library Build Script
# Author: Phong Ngo
# Date: June 15, 2025

set -e

echo "=== Building MKV Encryption Library ==="
echo "This script will build the MKV encryption library"
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
    if [ ! -f "go.mod" ] || [ ! -d "core/ledger/kvledger/txmgmt/statedb/mkv" ]; then
        print_status "FAIL" "Not in Fabric root directory with MKV integration"
        print_status "INFO" "Current directory: $(pwd)"
        print_status "INFO" "Please run this script from fabric-3.1.1/ directory"
        exit 1
    fi
    print_status "INFO" "Running from correct directory: $(pwd)"
}

# Function to check dependencies
check_dependencies() {
    print_status "INFO" "Checking dependencies..."
    
    # Check if gcc is available
    if ! command -v gcc >/dev/null 2>&1; then
        print_status "FAIL" "gcc compiler not found. Please install build-essential."
        exit 1
    fi
    
    # Check if make is available
    if ! command -v make >/dev/null 2>&1; then
        print_status "FAIL" "make not found. Please install build-essential."
        exit 1
    fi
    
    print_status "PASS" "All dependencies are available"
}

# Function to build MKV library
build_mkv_library() {
    print_status "INFO" "Building MKV encryption library..."
    
    cd core/ledger/kvledger/txmgmt/statedb/mkv
    
    # Clean previous builds
    if [ -f "Makefile" ]; then
        print_status "INFO" "Cleaning previous build..."
        make clean 2>/dev/null || true
    fi
    
    # Build the library
    print_status "INFO" "Compiling MKV library..."
    if make; then
        print_status "PASS" "MKV library built successfully"
    else
        print_status "FAIL" "Failed to build MKV library"
        exit 1
    fi
    
    # Verify the library was created
    if [ -f "libmkv.so" ]; then
        print_status "PASS" "libmkv.so created successfully"
        ls -la libmkv.so
    else
        print_status "FAIL" "libmkv.so not found after build"
        exit 1
    fi
    
    # Return to root directory
    cd /home/phongnh/go-src/Kmasc/fabric-3.1.1
}

# Function to test the library
test_mkv_library() {
    print_status "INFO" "Testing MKV library..."
    
    cd core/ledger/kvledger/txmgmt/statedb/mkv
    
    # Run Go tests to verify the library works
    if go test -v -run TestEncryptDecryptValueMKV; then
        print_status "PASS" "MKV library test passed"
    else
        print_status "FAIL" "MKV library test failed"
        exit 1
    fi
    
    # Return to root directory
    cd /home/phongnh/go-src/Kmasc/fabric-3.1.1
}

# Main execution
main() {
    echo "Starting MKV library build..."
    echo
    
    ensure_correct_directory
    check_dependencies
    build_mkv_library
    test_mkv_library
    
    echo
    print_status "PASS" "MKV encryption library build completed successfully!"
    echo
    echo "Library files created:"
    echo "  - libmkv.so (shared library)"
    echo "  - *.o (object files)"
    echo
    echo "You can now proceed with initializing MKV keys."
}

# Execute main function directly
echo "This script will:"
echo "1. Check dependencies (gcc, make)"
echo "2. Clean previous builds"
echo "3. Build MKV encryption library"
echo "4. Test the library with Go tests"
echo "5. Verify all files are created"
echo
echo "Starting execution..."
main "$@"
