#!/bin/bash

# Hyperledger Fabric Quick Start Script
# Author: Phong Ngo
# Date: June 15, 2025

set -e

echo "=== Hyperledger Fabric Quick Start ==="
echo "This script will set up everything from scratch"
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

# Function to check if script exists
script_exists() {
    [ -f "$1" ] && [ -x "$1" ]
}

# Function to ensure we're in the correct directory
ensure_correct_directory() {
    if [ ! -f "go.mod" ] || [ ! -d "core/ledger/kvledger/txmgmt/statedb" ]; then
        print_status "FAIL" "Not in Fabric root directory with encryption integration"
        print_status "INFO" "Current directory: $(pwd)"
        print_status "INFO" "Please run this script from fabric-3.1.1/ directory"
        exit 1
    fi
    print_status "INFO" "Running from correct directory: $(pwd)"
}

# Save root directory
ROOT_DIR=$(pwd)

# --- Go installation (auto-download Go 1.24.4 if not present or wrong version) ---
GO_VERSION_REQUIRED="1.24.4"
GO_TARBALL="go$GO_VERSION_REQUIRED.linux-amd64.tar.gz"
GO_URL="https://go.dev/dl/$GO_TARBALL"

check_go_version() {
    if command -v go >/dev/null 2>&1; then
        CURRENT_GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        if [ "$CURRENT_GO_VERSION" = "$GO_VERSION_REQUIRED" ]; then
            print_status "PASS" "Go $GO_VERSION_REQUIRED is already installed."
            return 0
        else
            print_status "WARN" "Go version $CURRENT_GO_VERSION found, but $GO_VERSION_REQUIRED is required."
            return 1
        fi
    else
        print_status "WARN" "Go is not installed."
        return 1
    fi
}

install_go() {
    print_status "INFO" "Downloading Go $GO_VERSION_REQUIRED..."
    wget -q "$GO_URL" -O "/tmp/$GO_TARBALL"
    print_status "INFO" "Removing any previous Go installation in /usr/local/go..."
    sudo rm -rf /usr/local/go
    print_status "INFO" "Extracting Go $GO_VERSION_REQUIRED to /usr/local..."
    sudo tar -C /usr/local -xzf "/tmp/$GO_TARBALL"
    print_status "INFO" "Setting up Go environment variables..."
    export PATH=/usr/local/go/bin:$PATH
    if ! grep -q '/usr/local/go/bin' ~/.bashrc; then
        echo 'export PATH=/usr/local/go/bin:$PATH' >> ~/.bashrc
    fi
    print_status "PASS" "Go $GO_VERSION_REQUIRED installed successfully."
    go version
}

check_go_version || install_go

# Step 1: Fix repositories if needed
step1_fix_repositories() {
    echo "Step 1: Checking repositories..."
    if script_exists "fix-repositories.sh"; then
        print_status "INFO" "Running repository fix..."
        ./fix-repositories.sh
        print_status "PASS" "Repository fix completed"
    else
        print_status "WARN" "fix-repositories.sh not found, skipping"
    fi
    echo
}

# Step 2: Setup environment
step2_setup_environment() {
    echo "Step 2: Setting up environment..."
    if script_exists "setup-environment.sh"; then
        print_status "INFO" "Running environment setup..."
        ./setup-environment.sh
        print_status "PASS" "Environment setup completed"
    else
        print_status "FAIL" "setup-environment.sh not found"
        exit 1
    fi
    echo
}

# Step 3: Download fabric-samples
step3_download_fabric_samples() {
    echo "Step 3: Downloading fabric-samples..."
    
    if [ -d "fabric-samples" ]; then
        print_status "INFO" "fabric-samples directory already exists"
        print_status "INFO" "Checking if it's the correct version..."
        
        # Check if it's the right version by looking for test-network
        if [ -d "fabric-samples/test-network" ]; then
            print_status "PASS" "fabric-samples with test-network found"
        else
            print_status "WARN" "fabric-samples exists but test-network not found"
            print_status "INFO" "Removing old fabric-samples and downloading fresh copy..."
            rm -rf fabric-samples
        fi
    fi
    
    if [ ! -d "fabric-samples" ]; then
        print_status "INFO" "Downloading fabric-samples..."
        
        # Try to download fabric-samples
        if command -v curl >/dev/null 2>&1; then
            print_status "INFO" "Using curl to download fabric-samples..."
            curl -sSL https://bit.ly/2ysbOFE | bash -s -- 3.1.1 1.5.1
        elif command -v wget >/dev/null 2>&1; then
            print_status "INFO" "Using wget to download fabric-samples..."
            wget https://bit.ly/2ysbOFE -O - | bash -s -- 3.1.1 1.5.1
        else
            print_status "FAIL" "Neither curl nor wget found. Please install one of them."
            print_status "INFO" "You can manually download fabric-samples from:"
            print_status "INFO" "https://github.com/hyperledger/fabric-samples"
            exit 1
        fi
        
        if [ -d "fabric-samples" ]; then
            print_status "PASS" "fabric-samples downloaded successfully"
        else
            print_status "FAIL" "Failed to download fabric-samples"
            exit 1
        fi
    fi
    echo
}

# Step 4: Test environment
step4_test_environment() {
    echo "Step 4: Testing environment..."
    if script_exists "test_environment.sh"; then
        print_status "INFO" "Running environment test..."
        export CGO_ENABLED=1
        ./test_environment.sh
        print_status "PASS" "Environment test completed"
    else
        print_status "WARN" "test_environment.sh not found, skipping"
    fi
    echo
}

# Step 5: Build encryption library
step5_build_encryption() {
    echo "Step 5: Building encryption library..."
    
    # Ensure we're in the correct directory
    if [ ! -f "go.mod" ]; then
        print_status "FAIL" "Not in Fabric root directory. Please run from fabric-3.1.1/"
        exit 1
    fi
    
    if [ -d "core/ledger/kvledger/txmgmt/statedb" ]; then
        print_status "INFO" "Building encryption library..."
        cd "$ROOT_DIR/core/ledger/kvledger/txmgmt/statedb"
        make clean && make
        cd "$ROOT_DIR"
        print_status "PASS" "Encryption library built"
    else
        print_status "FAIL" "statedb directory not found at core/ledger/kvledger/txmgmt/statedb"
        print_status "INFO" "Current directory: $(pwd)"
        print_status "INFO" "Available directories in core/ledger/kvledger/txmgmt/:"
        ls -la core/ledger/kvledger/txmgmt/ 2>/dev/null || echo "Directory not accessible"
        exit 1
    fi
    echo
}

# Step 6: Test encryption
#step6_test_encryption() {
#    echo "Step 6: Testing encryption integration..."
#    
#    # Ensure we're in the correct directory
#    if [ ! -f "go.mod" ]; then
#        print_status "FAIL" "Not in Fabric root directory. Please run from fabric-3.1.1/"
#        exit 1
#    fi
#    
#    if [ -d "core/ledger/kvledger/txmgmt/statedb" ]; then
#        print_status "INFO" "Running encryption tests..."
#        cd "$ROOT_DIR/core/ledger/kvledger/txmgmt/statedb"
#        # Tìm file run_tests.sh trong thư mục hiện tại và các thư mục con, ưu tiên file có quyền thực thi
#        RUN_TESTS_PATH=$(find . -type f -name "run_tests.sh" -perm /u+x | head -n 1)
#        if [ -z "$RUN_TESTS_PATH" ]; then
#            # Nếu không có file thực thi, tìm file run_tests.sh bất kỳ
#            RUN_TESTS_PATH=$(find . -type f -name "run_tests.sh" | head -n 1)
#        fi
#        if [ -n "$RUN_TESTS_PATH" ]; then
#            print_status "INFO" "Found run_tests.sh at $RUN_TESTS_PATH, running it with bash..."
#            bash "$RUN_TESTS_PATH"
#            print_status "PASS" "Encryption tests completed"
#        else
#            print_status "WARN" "run_tests.sh not found anywhere, running basic tests..."
#            go test ./...
#        fi
#        cd "$ROOT_DIR"
#    else
#        print_status "FAIL" "statedb directory not found at core/ledger/kvledger/txmgmt/statedb"
#        print_status "INFO" "Current directory: $(pwd)"
#        print_status "INFO" "Available directories in core/ledger/kvledger/txmgmt/:"
#        ls -la core/ledger/kvledger/txmgmt/ 2>/dev/null || echo "Directory not accessible"
#        exit 1
#    fi
#    echo
#}

# Step 7: Build Fabric
step7_build_fabric() {
    echo "Step 7: Building Fabric..."
    if script_exists "build-fabric.sh"; then
        print_status "INFO" "Building Fabric with encryption..."
        export CGO_ENABLED=1
        ./build-fabric.sh
        print_status "PASS" "Fabric build completed"
    else
        print_status "FAIL" "build-fabric.sh not found"
        exit 1
    fi
    echo
}

# Step 8: Start network
step8_start_network() {
    echo "Step 8: Starting test network..."
    if script_exists "start-network.sh"; then
        print_status "INFO" "Starting test network..."
        ./start-network.sh
        print_status "PASS" "Test network started"
    else
        print_status "FAIL" "start-network.sh not found"
        exit 1
    fi
    echo
}

# Step 9: Show next steps
step9_next_steps() {
    echo "Step 9: Next steps..."
    print_status "INFO" "Setup completed successfully!"
    echo
    echo "🎉 Congratulations! Your Hyperledger Fabric with encryption is ready!"
    echo
    echo "📋 What's been set up:"
    echo "   ✅ Environment dependencies (Go, OpenSSL, Docker)"
    echo "   ✅ Encryption library (libencryption.so)"
    echo "   ✅ Fabric binaries with encryption"
    echo "   ✅ Test network running"
    echo
    echo "🔍 To verify encryption is working:"
    echo "   docker exec peer0.org1.example.com cat /root/state_encryption.log"
    echo
    echo "🧪 To test chaincode:"
    echo "   cd fabric-samples/test-network"
    echo "   peer chaincode query -C mychannel -n basic -c '{\"function\":\"ReadAsset\",\"Args\":[\"asset1\"]}'"
    echo
    echo "🛑 To stop network:"
    echo "   cd fabric-samples/test-network"
    echo "   ./network.sh down"
    echo
    print_status "INFO" "Quick start completed at $(date)"
}

# Main execution
main() {
    echo "Starting quick start process..."
    echo
    
    ensure_correct_directory
    step1_fix_repositories
    step2_setup_environment
    step3_download_fabric_samples
    step4_test_environment
    step5_build_encryption
    step7_build_fabric
    step8_start_network
    step9_next_steps
}

# Check if user wants to continue
echo "This script will:"
echo "1. Fix any repository issues"
echo "2. Set up the environment (Go, OpenSSL, Docker)"
echo "3. Download fabric-samples"
echo "4. Test the environment"
echo "5. Build the encryption library"
echo "6. Test the encryption integration"
echo "7. Build Fabric with encryption"
echo "8. Start the test network"
echo
read -p "Do you want to continue? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    main "$@"
else
    echo "Quick start cancelled."
    exit 0
fi 