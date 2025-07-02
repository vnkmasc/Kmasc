#!/bin/bash

# Download Fabric Samples Script
# Author: Phong Ngo
# Date: June 15, 2025

set -e

echo "=== Download Fabric Samples ==="
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

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Main function to download fabric-samples
download_fabric_samples() {
    echo "Step 3: Downloading fabric-samples..."
    
    if [ -d "fabric-samples" ]; then
        print_status "INFO" "fabric-samples directory already exists"
        print_status "INFO" "Checking if it's the correct version..."
        
        # Check if it's the right version by looking for test-network
        if [ -d "fabric-samples/test-network" ]; then
            print_status "PASS" "fabric-samples with test-network found"
            return 0
        else
            print_status "WARN" "fabric-samples exists but test-network not found"
            print_status "INFO" "Removing old fabric-samples and downloading fresh copy..."
            rm -rf fabric-samples
        fi
    fi
    
    if [ ! -d "fabric-samples" ]; then
        print_status "INFO" "Downloading fabric-samples repository only..."
        
        # Clone fabric-samples repository directly
        if command_exists git; then
            print_status "INFO" "Using git to clone fabric-samples..."
            
            # Try to clone with timeout
            if timeout 120s git clone https://github.com/hyperledger/fabric-samples.git; then
                print_status "PASS" "Git clone successful"
                
                # Checkout specific version if needed
                cd fabric-samples
                if git tag | grep -q "v3.1.1"; then
                    print_status "INFO" "Checking out v3.1.1 tag..."
                    git checkout v3.1.1
                else
                    print_status "WARN" "v3.1.1 tag not found, using main branch"
                fi
                cd ..
                
            else
                print_status "FAIL" "Git clone failed (timeout or network issue)"
                print_status "INFO" "Please check your internet connection and try again"
                print_status "INFO" "You can manually download fabric-samples from:"
                print_status "INFO" "https://github.com/hyperledger/fabric-samples"
                exit 1
            fi
            
        else
            print_status "FAIL" "Git not found. Please install git."
            print_status "INFO" "You can manually download fabric-samples from:"
            print_status "INFO" "https://github.com/hyperledger/fabric-samples"
            exit 1
        fi
        
        if [ -d "fabric-samples" ]; then
            print_status "PASS" "fabric-samples repository downloaded successfully"
            
            # Verify the download
            if [ -d "fabric-samples/test-network" ]; then
                print_status "PASS" "test-network directory found"
            else
                print_status "WARN" "test-network directory not found in downloaded samples"
            fi
            
            # Note: bin directory will be created when we build Fabric
            print_status "INFO" "Note: bin directory will be created when building Fabric"
            
        else
            print_status "FAIL" "Failed to download fabric-samples"
            exit 1
        fi
    fi
    
    echo
    print_status "INFO" "Fabric samples download completed at $(date)"
    print_status "INFO" "Docker images and binaries will be downloaded separately when needed"
}

# Check if script is being run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    download_fabric_samples
fi 