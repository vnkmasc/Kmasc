#!/bin/bash

# Script để restore encryption gốc
# Chạy: ./restore-encryption.sh

set -e

echo "=== Restoring Original Encryption ==="
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

# Find latest backup directory
find_latest_backup() {
    ENCRYPTION_DIR="core/ledger/kvledger/txmgmt/statedb"
    
    if dir_exists "$ENCRYPTION_DIR"; then
        # Find the most recent backup directory
        LATEST_BACKUP=$(ls -td "${ENCRYPTION_DIR}/backup_"* 2>/dev/null | head -n1)
        if [ -n "$LATEST_BACKUP" ] && dir_exists "$LATEST_BACKUP"; then
            echo "$LATEST_BACKUP"
        else
            echo ""
        fi
    else
        echo ""
    fi
}

# Restore original encryption files
restore_encryption_files() {
    echo "Step 1: Finding backup directory..."
    
    BACKUP_DIR=$(find_latest_backup)
    
    if [ -z "$BACKUP_DIR" ]; then
        print_status "FAIL" "No backup directory found"
        print_status "INFO" "You may need to restore manually from git"
        exit 1
    fi
    
    print_status "INFO" "Found backup directory: $BACKUP_DIR"
    
    ENCRYPTION_DIR="core/ledger/kvledger/txmgmt/statedb"
    
    # Restore original files
    if file_exists "${BACKUP_DIR}/encrypt.go"; then
        cp "${BACKUP_DIR}/encrypt.go" "${ENCRYPTION_DIR}/"
        print_status "PASS" "Restored encrypt.go"
    fi
    
    if file_exists "${BACKUP_DIR}/encrypt.c"; then
        cp "${BACKUP_DIR}/encrypt.c" "${ENCRYPTION_DIR}/"
        print_status "PASS" "Restored encrypt.c"
    fi
    
    if file_exists "${BACKUP_DIR}/encrypt.h"; then
        cp "${BACKUP_DIR}/encrypt.h" "${ENCRYPTION_DIR}/"
        print_status "PASS" "Restored encrypt.h"
    fi
    
    # Restore disabled C files
    if file_exists "${ENCRYPTION_DIR}/encrypt.c.disabled"; then
        mv "${ENCRYPTION_DIR}/encrypt.c.disabled" "${ENCRYPTION_DIR}/encrypt.c"
        print_status "PASS" "Restored encrypt.c from disabled"
    fi
    
    if file_exists "${ENCRYPTION_DIR}/encrypt.h.disabled"; then
        mv "${ENCRYPTION_DIR}/encrypt.h.disabled" "${ENCRYPTION_DIR}/encrypt.h"
        print_status "PASS" "Restored encrypt.h from disabled"
    fi
    
    if file_exists "${ENCRYPTION_DIR}/libencryption.so.disabled"; then
        mv "${ENCRYPTION_DIR}/libencryption.so.disabled" "${ENCRYPTION_DIR}/libencryption.so"
        print_status "PASS" "Restored libencryption.so from disabled"
    fi
    
    print_status "PASS" "Encryption files restored"
}

# Build Fabric with original encryption
build_fabric_with_encryption() {
    echo "Step 2: Building Fabric with original encryption..."
    
    print_status "INFO" "Building Hyperledger Fabric..."
    
    # Set environment variables
    export CGO_ENABLED=1
    
    # Build Fabric
    if make clean && make; then
        print_status "PASS" "Fabric built successfully with encryption enabled"
    else
        print_status "FAIL" "Failed to build Fabric"
        exit 1
    fi
}

# Test the build
test_build() {
    echo "Step 3: Testing the build..."
    
    print_status "INFO" "Running basic tests..."
    
    # Test if binaries were created
    if file_exists "build/bin/peer"; then
        print_status "PASS" "peer binary created"
    else
        print_status "FAIL" "peer binary not found"
    fi
    
    if file_exists "build/bin/orderer"; then
        print_status "PASS" "orderer binary created"
    else
        print_status "FAIL" "orderer binary not found"
    fi
    
    print_status "PASS" "Build test completed"
}

# Clean up backup directory (optional)
cleanup_backup() {
    echo "Step 4: Cleaning up backup directory..."
    
    BACKUP_DIR=$(find_latest_backup)
    
    if [ -n "$BACKUP_DIR" ]; then
        read -p "Do you want to remove the backup directory? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            rm -rf "$BACKUP_DIR"
            print_status "PASS" "Backup directory removed"
        else
            print_status "INFO" "Backup directory kept: $BACKUP_DIR"
        fi
    fi
}

# Main execution
main() {
    echo "Starting restore process..."
    
    # Check if we're in the right directory
    if [ ! -f "go.mod" ]; then
        print_status "FAIL" "Not in Fabric root directory. Please run from fabric-3.1.1/"
        exit 1
    fi
    
    # Check required tools
    if ! command_exists go; then
        print_status "FAIL" "Go not found"
        exit 1
    fi
    
    if ! command_exists make; then
        print_status "FAIL" "Make not found"
        exit 1
    fi
    
    # Execute steps
    restore_encryption_files
    build_fabric_with_encryption
    test_build
    cleanup_backup
    
    echo
    print_status "INFO" "Restore completed successfully at $(date)"
    print_status "INFO" "Fabric now has encryption ENABLED"
}

# Check if script is being run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main
fi 