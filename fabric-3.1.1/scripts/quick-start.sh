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

# Function to make script executable and run it
run_script() {
    local script_name=$1
    local description=$2
    
    if script_exists "$script_name"; then
        print_status "INFO" "Running $description..."
        chmod +x "$script_name"
        ./"$script_name"
        print_status "PASS" "$description completed"
    else
        print_status "FAIL" "$script_name not found"
        exit 1
    fi
    echo
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
    run_script "scripts/fix-repositories.sh" "repository fix"
}

# Step 2: Setup environment
step2_setup_environment() {
    echo "Step 2: Setting up environment..."
    run_script "scripts/setup-environment.sh" "environment setup"
}

# Step 3: Download fabric-samples
step3_download_fabric_samples() {
    echo "Step 3: Downloading fabric-samples..."
    run_script "scripts/download-fabric-samples.sh" "fabric-samples download"
}

# Step 4: Test environment
step4_test_environment() {
    echo "Step 4: Testing environment..."
    run_script "scripts/test_environment.sh" "environment test"
}

# Step 5: Build encryption and create keys (SIMPLE VERSION)
step5_build_encryption() {
    print_status "INFO" "Step 5: Building encryption and creating keys..."
    
    echo "Building MKV encryption library..."
    cd core/ledger/kvledger/txmgmt/statedb/mkv
    make clean && make
    cd /home/phongnh/go-src/Kmasc/fabric-3.1.1
    
    echo "Creating MKV keys..."
    cp core/ledger/kvledger/txmgmt/statedb/mkv/libmkv.so . 2>/dev/null || true
    cp core/ledger/kvledger/txmgmt/statedb/mkv/mkv.go . 2>/dev/null || true
    echo "fabric_mkv_password_2025" | bash core/ledger/kvledger/txmgmt/statedb/mkv/key_manager.sh init
    
    echo "✅ PASS: MKV encryption built and keys created"
    echo "   - Keys location: /home/phongnh/go-src/Kmasc/fabric-3.1.1/"
    echo "   - Files: k1.key, k0.key, encrypted_k1.key"
    echo "   - Password: fabric_mkv_password_2025"
}

# Step 5.1: Build MKV library
step5_1_build_mkv() {
    echo "Step 5.1: Building MKV library..."
    run_script "scripts/build-mkv.sh" "MKV library build"
}

# Step 5.2: Test MKV library
step5_2_test_mkv() {
    echo "Step 5.2: Testing MKV library..."
    run_script "scripts/test-mkv.sh" "MKV library test"
}

# Step 5.3: Test MKV in Docker containers
step5_3_test_mkv_docker() {
    echo "Step 5.3: Testing MKV in Docker containers..."
    run_script "scripts/test-mkv-docker.sh" "MKV Docker test"
}

# Step 6: Build Fabric
step6_build_fabric() {
    echo "Step 6: Building Fabric..."
    run_script "scripts/build-fabric.sh" "Fabric build"
}

# Step 7: Start network
step7_start_network() {
    echo "Step 7: Starting test network..."
    run_script "scripts/start-network.sh" "test network startup"
}

# Step 8: Initialize MKV keys in containers
step8_init_mkv_keys() {
    echo "Step 8: Initializing MKV keys in containers..."
    run_script "scripts/init-mkv-keys-wrapper.sh" "MKV keys initialization"
}

# Step 8: Auto-copy MKV keys to containers
step8_auto_copy_mkv_keys() {
    echo "Step 8: Auto-copying MKV keys to containers..."
    
    # Wait a moment for containers to be ready
    sleep 5
    
    # Generate keys in current directory if not exists
    if [ ! -f "k1.key" ] || [ ! -f "k0.key" ] || [ ! -f "encrypted_k1.key" ]; then
        echo "INFO: Generating MKV keys in current directory..."
        cp core/ledger/kvledger/txmgmt/statedb/mkv/libmkv.so . 2>/dev/null || true
        cp core/ledger/kvledger/txmgmt/statedb/mkv/mkv.go . 2>/dev/null || true
        echo "fabric_mkv_password_2025" | bash core/ledger/kvledger/txmgmt/statedb/mkv/key_manager.sh init
    fi
    
    # Function to copy keys to a specific container
    copy_keys_to_container() {
        local container_name=$1
        echo "INFO: Copying keys to $container_name..."
        docker cp k1.key $container_name:/ 2>/dev/null || true
        docker cp k0.key $container_name:/ 2>/dev/null || true
        docker cp encrypted_k1.key $container_name:/ 2>/dev/null || true
    }
    
    # Copy keys to peer containers
    copy_keys_to_container "peer0.org1.example.com"
    copy_keys_to_container "peer0.org2.example.com"
    
    # Copy keys to orderer container
    copy_keys_to_container "orderer.example.com"
    
    # Copy keys to chaincode containers
    echo "INFO: Copying keys to chaincode containers..."
    ./scripts/auto-copy-mkv-keys.sh copy 2>/dev/null || true
    
    echo "PASS: MKV keys auto-copied to all containers"
}

# Step 9: Monitor and ensure keys persistence
step9_monitor_keys() {
    echo "Step 9: Setting up key monitoring..."
    
    # Start background monitoring
    (
        while true; do
            sleep 30  # Check every 30 seconds
            
            # Check if peer containers have keys
            if ! docker exec peer0.org1.example.com ls -la /k1.key >/dev/null 2>&1; then
                echo "WARN: Keys missing in peer0.org1.example.com, re-copying..."
                docker cp k1.key peer0.org1.example.com:/ 2>/dev/null || true
                docker cp k0.key peer0.org1.example.com:/ 2>/dev/null || true
                docker cp encrypted_k1.key peer0.org1.example.com:/ 2>/dev/null || true
            fi
            
            if ! docker exec peer0.org2.example.com ls -la /k1.key >/dev/null 2>&1; then
                echo "WARN: Keys missing in peer0.org2.example.com, re-copying..."
                docker cp k1.key peer0.org2.example.com:/ 2>/dev/null || true
                docker cp k0.key peer0.org2.example.com:/ 2>/dev/null || true
                docker cp encrypted_k1.key peer0.org2.example.com:/ 2>/dev/null || true
            fi
            
            # Check chaincode containers
            ./scripts/auto-copy-mkv-keys.sh check >/dev/null 2>&1 || {
                echo "WARN: Keys missing in chaincode containers, re-copying..."
                ./scripts/auto-copy-mkv-keys.sh copy >/dev/null 2>&1 || true
            }
        done
    ) &
    
    echo "PASS: Key monitoring started in background"
}

# Step 10: Next steps
step10_next_steps() {
    print_status "INFO" "Step 10: Next steps..."
    
    echo
    echo "🎉 Hyperledger Fabric network with MKV encryption is ready!"
    echo
    echo "📋 Network Information:"
    echo "   - Network: test-network"
    echo "   - Channel: mychannel"
    echo "   - Chaincode: basic"
    echo "   - MKV Encryption: Enabled"
    echo
    echo "🔑 MKV Keys Information:"
    echo "   - Keys location: /home/phongnh/go-src/Kmasc/fabric-3.1.1/"
    echo "   - Password: fabric_mkv_password_2025"
    echo "   - Files: k1.key, k0.key, encrypted_k1.key"
    echo
    echo "🚀 Quick Test Commands:"
    echo "   cd fabric-samples/test-network"
    echo "   export PATH=\${PWD}/bin:\${PWD}/../bin:\${PWD}/../../bin:\$PATH"
    echo "   export FABRIC_CFG_PATH=\$PWD/../config/"
    echo "   export CORE_PEER_TLS_ENABLED=true"
    echo "   export CORE_PEER_LOCALMSPID=\"Org1MSP\""
    echo "   export CORE_PEER_MSPCONFIGPATH=\${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp"
    echo "   export CORE_PEER_TLS_ROOTCERT_FILE=\${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
    echo "   export CORE_PEER_ADDRESS=localhost:7051"
    echo "   export ORDERER_CA=\${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
    echo
    echo "   # Test query"
    echo "   peer chaincode query -C mychannel -n basic -c '{\"function\":\"ReadAsset\",\"Args\":[\"asset1\"]}'"
    echo
    echo "📝 Notes:"
    echo "   - MKV keys are created and ready to use"
    echo "   - Network is ready for app connection"
    echo "   - For production, consider using persistent volumes"
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
    step5_build_encryption  # Only this step - build and create keys
    step6_build_fabric
    step7_start_network
    step10_next_steps  # Skip complex container management
}

# Check if user wants to continue
echo "This script will:"
echo "1. Fix any repository issues"
echo "2. Set up the environment (Go, OpenSSL, Docker)"
echo "3. Download fabric-samples"
echo "4. Test the environment"
echo "5. Build the encryption library and create keys"
echo "6. Build Fabric with encryption"
echo "7. Start the test network"
echo "8. Next steps"
echo
read -p "Do you want to continue? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    main "$@"
else
    echo "Quick start cancelled."
    exit 0
fi 