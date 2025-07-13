#!/bin/bash

# Script để build Hyperledger Fabric với encryption disabled
# Chạy: ./build-fabric-no-encryption.sh

set -e

echo "=== Building Hyperledger Fabric with Encryption DISABLED ==="
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

# Step 1: Backup original encryption files
backup_encryption_files() {
    echo "Step 1: Backing up original encryption files..."
    
    ENCRYPTION_DIR="core/ledger/kvledger/txmgmt/statedb"
    
    if dir_exists "$ENCRYPTION_DIR"; then
        print_status "INFO" "Creating backup of encryption files..."
        
        # Create backup directory
        BACKUP_DIR="${ENCRYPTION_DIR}/backup_$(date +%Y%m%d_%H%M%S)"
        mkdir -p "$BACKUP_DIR"
        
        # Backup original files
        if file_exists "${ENCRYPTION_DIR}/encrypt.go"; then
            cp "${ENCRYPTION_DIR}/encrypt.go" "$BACKUP_DIR/"
            print_status "PASS" "Backed up encrypt.go"
        fi
        
        if file_exists "${ENCRYPTION_DIR}/encrypt.c"; then
            cp "${ENCRYPTION_DIR}/encrypt.c" "$BACKUP_DIR/"
            print_status "PASS" "Backed up encrypt.c"
        fi
        
        if file_exists "${ENCRYPTION_DIR}/encrypt.h"; then
            cp "${ENCRYPTION_DIR}/encrypt.h" "$BACKUP_DIR/"
            print_status "PASS" "Backed up encrypt.h"
        fi
        
        print_status "PASS" "Backup completed: $BACKUP_DIR"
        
    else
        print_status "FAIL" "Encryption directory not found: $ENCRYPTION_DIR"
        exit 1
    fi
}

# Step 2: Create disabled encryption files
create_disabled_encryption() {
    echo "Step 2: Creating disabled encryption files..."
    
    ENCRYPTION_DIR="core/ledger/kvledger/txmgmt/statedb"
    
    # Create disabled encrypt.go
    cat > "${ENCRYPTION_DIR}/encrypt.go" << 'EOF'
//go:build !test
// +build !test

package statedb

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/hyperledger/fabric-lib-go/common/flogging"
)

var encryptLogger = flogging.MustGetLogger("encrypt")

var (
	logFileOnce sync.Once
	logFile     *os.File
	logFileErr  error
	logFileMu   sync.Mutex
)

func logToFile(op, ns, key, status, errMsg string) {
	logFileOnce.Do(func() {
		logFile, logFileErr = os.OpenFile("/root/state_encryption_disabled.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	})
	if logFileErr != nil || logFile == nil {
		return
	}
	logFileMu.Lock()
	defer logFileMu.Unlock()
	now := time.Now().UTC()
	timestamp := fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%06dZ",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
		now.Nanosecond()/1000)
	msg := timestamp + " " + op + " ns=" + ns + " key=" + key + " " + status
	if errMsg != "" {
		msg += " ERROR: " + errMsg
	}
	logFile.WriteString(msg + "\n")
}

// EncryptValue - DISABLED: trả về nguyên input (không mã hóa)
func EncryptValue(value []byte, ns, key string) []byte {
	if value == nil || len(value) == 0 {
		logToFile("ENCRYPT_DISABLED", ns, key, "SKIP_EMPTY", "")
		return value
	}
	logToFile("ENCRYPT_DISABLED", ns, key, "SUCCESS", "Encryption disabled - returning original data")
	return value
}

// DecryptValue - DISABLED: trả về nguyên input (không giải mã)
func DecryptValue(value []byte, ns, key string) []byte {
	if value == nil || len(value) == 0 {
		logToFile("DECRYPT_DISABLED", ns, key, "SKIP_EMPTY", "")
		return value
	}
	logToFile("DECRYPT_DISABLED", ns, key, "SUCCESS", "Decryption disabled - returning original data")
	return value
}
EOF

    # Remove C files to avoid build errors
    if file_exists "${ENCRYPTION_DIR}/encrypt.c"; then
        mv "${ENCRYPTION_DIR}/encrypt.c" "${ENCRYPTION_DIR}/encrypt.c.disabled"
        print_status "PASS" "Disabled encrypt.c"
    fi
    
    if file_exists "${ENCRYPTION_DIR}/encrypt.h"; then
        mv "${ENCRYPTION_DIR}/encrypt.h" "${ENCRYPTION_DIR}/encrypt.h.disabled"
        print_status "PASS" "Disabled encrypt.h"
    fi
    
    if file_exists "${ENCRYPTION_DIR}/libencryption.so"; then
        mv "${ENCRYPTION_DIR}/libencryption.so" "${ENCRYPTION_DIR}/libencryption.so.disabled"
        print_status "PASS" "Disabled libencryption.so"
    fi

    print_status "PASS" "Created disabled encrypt.go"
}

# Step 3: Build Fabric with disabled encryption
build_fabric_no_encryption() {
    echo "Step 3: Building Fabric with encryption disabled..."
    
    print_status "INFO" "Building Hyperledger Fabric..."
    
    # Set environment variables
    export CGO_ENABLED=1
    
    # Build Fabric
    if make clean && make; then
        print_status "PASS" "Fabric built successfully with encryption disabled"
    else
        print_status "FAIL" "Failed to build Fabric"
        exit 1
    fi
}

# Step 4: Test the build
test_build() {
    echo "Step 4: Testing the build..."
    
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

# Main execution
main() {
    echo "Starting build process..."
    
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
    backup_encryption_files
    create_disabled_encryption
    build_fabric_no_encryption
    test_build
    
    echo
    print_status "INFO" "Build completed successfully at $(date)"
    print_status "INFO" "Fabric now has encryption DISABLED"
    print_status "INFO" "You can now run performance tests to compare with encryption enabled"
}

# Check if script is being run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main
fi 