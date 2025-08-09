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

# Step 5: Build MKV library
step5_build_mkv() {
    echo "Step 5: Building MKV library..."
    run_script "scripts/build-mkv.sh" "MKV library build"
}

# Step 6: Build Fabric with MKV
step6_build_fabric() {
    echo "Step 6: Building Fabric with MKV integration..."
    run_script "scripts/build-fabric.sh" "Fabric build with MKV"
}

# Step 7: Initialize MKV system
step7_initialize_mkv() {
    echo "Step 7: Initializing MKV system..."
    run_script "scripts/init-mkv-keys.sh" "MKV key initialization"
}

# Step 8: Start network
step8_start_network() {
    echo "Step 8: Starting test network..."
    run_script "scripts/start-network.sh" "test network startup"
}

# Step 9: Test MKV integration
step9_test_mkv() {
    echo "Step 9: Testing MKV integration..."
    run_script "scripts/test-mkv.sh" "MKV integration test"
}

# Step 10: Show next steps
step10_show_next_steps() {
    echo "Step 10: Showing next steps and usage information..."
    print_status "INFO" "Quick start completed successfully!"
    echo
    echo "🎉 === FABRIC WITH MKV IS READY! ==="
    echo
    echo "📊 What was set up:"
    echo "   ✅ MKV256 encryption library"
    echo "   ✅ Fabric peer with MKV integration"
    echo "   ✅ Test network with encryption"
    echo "   ✅ MKV key management system"
    echo
    echo "🚀 Next steps:"
    echo "   1. Test encryption: docker exec peer0.org1.example.com cat /tmp/state_mkv.log"
    echo "   2. Run chaincode transactions to see encryption in action"
    echo "   3. Change MKV password: cd core/ledger/kvledger/txmgmt/statedb/mkv && ./mkv_client.sh change"
    echo
    echo "📚 For production deployment, run:"
    echo "   ./scripts/quick-start-production.sh"
    echo
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
    step5_build_mkv          # Build MKV library
    step6_build_fabric       # Build Fabric with MKV
    step7_initialize_mkv     # Initialize MKV system
    step8_start_network      # Start test network
    step9_test_mkv          # Test MKV integration
    step10_show_next_steps  # Show next steps and usage info
}

# Check if user wants to continue
echo "This script will:"
echo "1. Fix any repository issues"
echo "2. Set up the environment (Go, OpenSSL, Docker)"
echo "3. Download fabric-samples"
echo "4. Test the environment"
echo "5. Build MKV256 encryption library"
echo "6. Build Fabric with MKV integration"
echo "7. Initialize MKV key management system"
echo "8. Start the test network"
echo "9. Test MKV encryption/decryption"
echo "10. Show next steps and usage examples"
echo
read -p "Do you want to continue? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    main "$@"
else
    echo "Quick start cancelled."
    exit 0
fi 