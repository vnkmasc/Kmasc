#!/bin/bash

# Test MKV in Docker Container Script
# Author: Phong Ngo
# Date: June 15, 2025

set -e

echo "=== Test MKV in Docker Container ==="
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

# Function to check if container is running
check_container_running() {
    local container_name=$1
    if docker ps --format "table {{.Names}}" | grep -q "^${container_name}$"; then
        return 0
    else
        return 1
    fi
}

# Function to get container ID
get_container_id() {
    local container_name=$1
    docker ps --format "table {{.ID}}\t{{.Names}}" | grep "^.*\t${container_name}$" | awk '{print $1}'
}

# Function to test MKV in a specific container
test_mkv_in_container() {
    local container_name=$1
    local working_dir=${2:-"/root/mkv"}
    
    print_status "INFO" "Testing MKV in container: ${container_name}"
    
    if check_container_running "${container_name}"; then
        local container_id=$(get_container_id "${container_name}")
        print_status "INFO" "Container ID: ${container_id}"
        
        # Test MKV library
        print_status "INFO" "Testing MKV library functions..."
        docker exec "${container_id}" bash -c "cd ${working_dir} && export LD_LIBRARY_PATH=${working_dir} && ./key_manager.sh status"
        
        # Test encryption/decryption
        print_status "INFO" "Testing MKV encryption/decryption..."
        docker exec "${container_id}" bash -c "cd ${working_dir} && export LD_LIBRARY_PATH=${working_dir} && echo 'test data' > test_input.txt && echo 'fabric_mkv_password_2025' | ./key_manager.sh test_encryption test_input.txt"
        
        # Check if test files exist
        docker exec "${container_id}" bash -c "cd ${working_dir} && ls -la test_*"
        
        print_status "PASS" "MKV test completed for ${container_name}"
        
    else
        print_status "WARN" "Container ${container_name} is not running, skipping..."
    fi
}

# Function to test MKV in all peer containers
test_mkv_in_peers() {
    local working_dir=${1:-"/root/mkv"}
    
    print_status "INFO" "Testing MKV in all peer containers..."
    
    # List of peer containers
    local peer_containers=("peer0.org1.example.com" "peer0.org2.example.com")
    
    for container_name in "${peer_containers[@]}"; do
        test_mkv_in_container "${container_name}" "${working_dir}"
        echo
    done
}

# Function to test MKV in orderer container
test_mkv_in_orderer() {
    local working_dir=${1:-"/root/mkv"}
    local container_name="orderer.example.com"
    
    test_mkv_in_container "${container_name}" "${working_dir}"
}

# Function to test MKV in all containers
test_mkv_in_all() {
    local working_dir=${1:-"/root/mkv"}
    
    print_status "INFO" "Testing MKV in all containers..."
    
    test_mkv_in_peers "${working_dir}"
    echo
    test_mkv_in_orderer "${working_dir}"
}

# Function to show help
show_help() {
    echo "MKV Docker Test Script"
    echo "======================"
    echo
    echo "Usage: $0 [OPTIONS] [COMMAND]"
    echo
    echo "Commands:"
    echo "  peers     - Test MKV in all peer containers"
    echo "  orderer   - Test MKV in orderer container"
    echo "  all       - Test MKV in all containers (peers + orderer)"
    echo "  help      - Show this help message"
    echo
    echo "Options:"
    echo "  -d, --dir DIRECTORY  - Working directory in container (default: /root/mkv)"
    echo "  -h, --help           - Show this help message"
    echo
    echo "Examples:"
    echo "  $0 all                    # Test MKV in all containers"
    echo "  $0 peers -d /opt/mkv     # Test MKV in peers with custom directory"
    echo
}

# Parse command line arguments
WORKING_DIR="/root/mkv"
COMMAND=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--dir)
            WORKING_DIR="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        peers|orderer|all|help)
            COMMAND="$1"
            shift
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
    case "${COMMAND:-help}" in
        peers)
            test_mkv_in_peers "$WORKING_DIR"
            ;;
        orderer)
            test_mkv_in_orderer "$WORKING_DIR"
            ;;
        all)
            test_mkv_in_all "$WORKING_DIR"
            ;;
        help|*)
            show_help
            ;;
    esac
    
    echo
    print_status "INFO" "MKV Docker test completed at $(date)"
}

# Check if we're in the correct directory
if [ ! -f "go.mod" ] || [ ! -d "core/ledger/kvledger/txmgmt/statedb/mkv" ]; then
    print_status "FAIL" "Not in Fabric root directory with MKV integration"
    print_status "INFO" "Please run this script from fabric-3.1.1/ directory"
    exit 1
fi

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    print_status "FAIL" "Docker is not running"
    exit 1
fi

# Run main function
main "$@"
