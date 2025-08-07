#!/bin/bash

# Standalone MKV Keys Initialization Script
# Author: Phong Ngo
# Date: June 15, 2025

set -e

echo "=== Standalone MKV Keys Initialization ==="
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

# Function to check if Docker is running
check_docker() {
    if ! docker info >/dev/null 2>&1; then
        print_status "FAIL" "Docker is not running"
        exit 1
    fi
    print_status "PASS" "Docker is running"
}

# Function to check if containers are running
check_containers() {
    print_status "INFO" "Checking if Fabric containers are running..."
    
    local containers=("peer0.org1.example.com" "peer0.org2.example.com" "orderer.example.com")
    local running_containers=()
    
    for container in "${containers[@]}"; do
        if docker ps --format "table {{.Names}}" | grep -q "^${container}$"; then
            running_containers+=("$container")
            print_status "PASS" "Container $container is running"
        else
            print_status "WARN" "Container $container is not running"
        fi
    done
    
    if [ ${#running_containers[@]} -eq 0 ]; then
        print_status "FAIL" "No Fabric containers are running"
        print_status "INFO" "Please start the Fabric network first:"
        print_status "INFO" "  cd fabric-samples/test-network"
        print_status "INFO" "  ./network.sh up"
        exit 1
    fi
    
    print_status "PASS" "Found ${#running_containers[@]} running containers"
}

# Function to initialize MKV keys in all containers
init_mkv_keys() {
    local password=${1:-"fabric_mkv_password_2025"}
    
    print_status "INFO" "Initializing MKV keys in all containers..."
    print_status "INFO" "Using password: $password"
    
    # Run the main initialization script
    if [ -f "scripts/init-mkv-keys.sh" ]; then
        chmod +x "scripts/init-mkv-keys.sh"
        ./scripts/init-mkv-keys.sh auto -p "$password"
        
        if [ $? -eq 0 ]; then
            print_status "PASS" "MKV keys initialized successfully in all containers"
        else
            print_status "FAIL" "Failed to initialize MKV keys"
            return 1
        fi
    else
        print_status "FAIL" "scripts/init-mkv-keys.sh not found"
        return 1
    fi
}

# Function to test MKV in all containers
test_mkv() {
    print_status "INFO" "Testing MKV in all containers..."
    
    # Run the test script
    if [ -f "scripts/test-mkv-docker.sh" ]; then
        chmod +x "scripts/test-mkv-docker.sh"
        ./scripts/test-mkv-docker.sh all
        
        if [ $? -eq 0 ]; then
            print_status "PASS" "MKV test completed successfully in all containers"
        else
            print_status "FAIL" "MKV test failed"
            return 1
        fi
    else
        print_status "FAIL" "scripts/test-mkv-docker.sh not found"
        return 1
    fi
}

# Function to show help
show_help() {
    echo "Standalone MKV Keys Initialization"
    echo "=================================="
    echo
    echo "Usage: $0 [OPTIONS]"
    echo
    echo "Options:"
    echo "  -p, --password PASSWORD  - Use custom password (default: fabric_mkv_password_2025)"
    echo "  --no-test                - Skip MKV testing after initialization"
    echo "  -h, --help               - Show this help message"
    echo
    echo "Examples:"
    echo "  $0                       # Initialize with default password and test"
    echo "  $0 -p mypassword         # Initialize with custom password"
    echo "  $0 --no-test             # Initialize only, skip testing"
    echo
}

# Parse command line arguments
PASSWORD="fabric_mkv_password_2025"
SKIP_TEST=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -p|--password)
            PASSWORD="$2"
            shift 2
            ;;
        --no-test)
            SKIP_TEST=true
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Main execution
main() {
    echo "Starting standalone MKV keys initialization..."
    echo
    
    # Check if we're in the correct directory
    if [ ! -f "go.mod" ] || [ ! -d "core/ledger/kvledger/txmgmt/statedb/mkv" ]; then
        print_status "FAIL" "Not in Fabric root directory with MKV integration"
        print_status "INFO" "Please run this script from fabric-3.1.1/ directory"
        exit 1
    fi
    
    # Check Docker
    check_docker
    
    # Check containers
    check_containers
    
    # Initialize MKV keys
    init_mkv_keys "$PASSWORD"
    
    # Test MKV (unless skipped)
    if [ "$SKIP_TEST" = false ]; then
        test_mkv
    else
        print_status "INFO" "Skipping MKV test as requested"
    fi
    
    echo
    print_status "INFO" "Standalone MKV keys initialization completed at $(date)"
    echo
    print_status "INFO" "Next steps:"
    print_status "INFO" "  - Check MKV logs: docker exec peer0.org1.example.com cat /root/state_mkv.log"
    print_status "INFO" "  - Test chaincode: cd fabric-samples/test-network && peer chaincode query -C mychannel -n basic -c '{\"function\":\"ReadAsset\",\"Args\":[\"asset1\"]}'"
}

# Run main function
main "$@"
