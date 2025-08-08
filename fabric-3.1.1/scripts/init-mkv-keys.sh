#!/bin/bash

# Initialize MKV Key Management System in Docker
# Author: Phong Ngo
# Date: June 15, 2025

set -e

echo "=== Initialize MKV Key Management System ==="
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

# Function to debug container status
debug_container_status() {
    print_status "INFO" "Debugging container status..."
    print_status "INFO" "All running containers:"
    docker ps --format "table {{.Names}}\t{{.Status}}"
    echo
    
    # Test container ID extraction for each container
    local containers=("peer0.org1.example.com" "peer0.org2.example.com" "orderer.example.com")
    for container_name in "${containers[@]}"; do
        print_status "INFO" "Testing container ID extraction for: ${container_name}"
        local container_id=$(get_container_id "${container_name}")
        if [ $? -eq 0 ]; then
            print_status "PASS" "Container ID for ${container_name}: ${container_id}"
        else
            print_status "FAIL" "Failed to get container ID for ${container_name}"
        fi
    done
    echo
}

# Function to test key_manager.sh locally
test_key_manager_locally() {
    print_status "INFO" "Testing key_manager.sh locally..."
    
    local mkv_dir="core/ledger/kvledger/txmgmt/statedb/mkv"
    
    if [ ! -f "${mkv_dir}/key_manager.sh" ]; then
        print_status "FAIL" "key_manager.sh not found in ${mkv_dir}"
        return 1
    fi
    
    if [ ! -f "${mkv_dir}/libmkv.so" ]; then
        print_status "FAIL" "libmkv.so not found in ${mkv_dir}"
        return 1
    fi
    
    print_status "PASS" "key_manager.sh and libmkv.so found"
    
    # Test if key_manager.sh is executable
    if [ ! -x "${mkv_dir}/key_manager.sh" ]; then
        print_status "INFO" "Making key_manager.sh executable..."
        chmod +x "${mkv_dir}/key_manager.sh"
    fi
    
    print_status "PASS" "key_manager.sh is ready for use"
}

# Function to get container ID
get_container_id() {
    local container_name=$1
    
    # Try multiple formats to get container ID
    local container_id=""
    
    # Method 1: Try table format
    container_id=$(docker ps --format "table {{.ID}}\t{{.Names}}" | grep "^.*\t${container_name}$" | awk '{print $1}' | head -1)
    
    # Method 2: If empty, try simple format
    if [ -z "$container_id" ]; then
        container_id=$(docker ps --format "{{.ID}}\t{{.Names}}" | grep "^.*\t${container_name}$" | awk '{print $1}' | head -1)
    fi
    
    # Method 3: If still empty, try with grep
    if [ -z "$container_id" ]; then
        container_id=$(docker ps -q --filter "name=${container_name}" | head -1)
    fi
    
    # Method 4: Last resort - try exact match
    if [ -z "$container_id" ]; then
        container_id=$(docker ps --format "{{.ID}}\t{{.Names}}" | grep "${container_name}" | awk '{print $1}' | head -1)
    fi
    
    if [ -z "$container_id" ]; then
        print_status "FAIL" "Container ID is empty for ${container_name}"
        print_status "INFO" "Available containers:"
        docker ps --format "table {{.Names}}\t{{.Status}}"
        return 1
    fi
    
    echo "$container_id"
}

# Function to copy files to container
copy_files_to_container() {
    local container_id=$1
    local source_dir=$2
    local dest_dir=$3
    
    print_status "INFO" "Copying MKV files to container..."
    
    # Copy library files
    docker cp "${source_dir}/libmkv.so" "${container_id}:${dest_dir}/"
    docker cp "${source_dir}/mkv.h" "${container_id}:${dest_dir}/"
    docker cp "${source_dir}/mkv.go" "${container_id}:${dest_dir}/"
    docker cp "${source_dir}/key_manager.sh" "${container_id}:${dest_dir}/"
    
    # Make key_manager.sh executable (ignore errors if not permitted)
    docker exec "${container_id}" chmod +x "${dest_dir}/key_manager.sh" 2>/dev/null || true
    
    print_status "PASS" "Files copied successfully"
}

# Function to initialize keys in container
init_keys_in_container() {
    local container_id=$1
    local working_dir=$2
    local password=$3
    
    print_status "INFO" "Initializing MKV keys in container..."
    
    # Set environment variables and run key initialization
    docker exec "${container_id}" bash -c "export LD_LIBRARY_PATH=${working_dir} && cd ${working_dir} && echo '${password}' | bash key_manager.sh init"
    
    if [ $? -eq 0 ]; then
        print_status "PASS" "MKV keys initialized successfully in container"
    else
        print_status "FAIL" "Failed to initialize MKV keys in container"
        return 1
    fi
}

# Function to verify keys in container
verify_keys_in_container() {
    local container_id=$1
    local working_dir=$2
    
    print_status "INFO" "Verifying MKV keys in container..."
    
    # Check if key files exist
    docker exec "${container_id}" bash -c "cd ${working_dir} && ls -la *.key"
    
    # Show key status
    docker exec "${container_id}" bash -c "cd ${working_dir} && ./key_manager.sh status"
    
    print_status "PASS" "MKV keys verified successfully"
}

# Function to setup MKV in all peer containers
setup_mkv_in_peers() {
    local password=${1:-"fabric_mkv_password_2025"}
    local working_dir=${2:-"/home/chaincode/mkv"}
    
    print_status "INFO" "Setting up MKV in all peer containers..."
    
    # List of peer containers
    local peer_containers=("peer0.org1.example.com" "peer0.org2.example.com")
    
    for container_name in "${peer_containers[@]}"; do
        print_status "INFO" "Processing container: ${container_name}"
        
        if check_container_running "${container_name}"; then
            local container_id=$(get_container_id "${container_name}")
            if [ $? -ne 0 ]; then
                print_status "FAIL" "Failed to get container ID for ${container_name}"
                continue
            fi
            
            print_status "INFO" "Container ID: ${container_id}"
            
            # Create working directory in container
            docker exec "${container_id}" mkdir -p "${working_dir}"
            
            # Copy MKV files to container
            copy_files_to_container "${container_id}" "core/ledger/kvledger/txmgmt/statedb/mkv" "${working_dir}"
            
            # Initialize keys in container
            init_keys_in_container "${container_id}" "${working_dir}" "${password}"
            
            # Verify keys in container
            verify_keys_in_container "${container_id}" "${working_dir}"
            
            # Copy keys to current working directory for MKV access
            copy_keys_to_working_dir "${container_id}" "${working_dir}"
            
            print_status "PASS" "MKV setup completed for ${container_name}"
            
        else
            print_status "WARN" "Container ${container_name} is not running, skipping..."
        fi
        
        echo
    done
}

# Function to setup MKV in orderer container
setup_mkv_in_orderer() {
    local password=${1:-"fabric_mkv_password_2025"}
    local working_dir=${2:-"/home/chaincode/mkv"}
    local container_name="orderer.example.com"
    
    print_status "INFO" "Setting up MKV in orderer container..."
    
    if check_container_running "${container_name}"; then
        local container_id=$(get_container_id "${container_name}")
        if [ $? -ne 0 ]; then
            print_status "FAIL" "Failed to get container ID for ${container_name}"
            return 1
        fi
        
        print_status "INFO" "Container ID: ${container_id}"
        
        # Create working directory in container
        docker exec "${container_id}" mkdir -p "${working_dir}"
        
        # Copy MKV files to container
        copy_files_to_container "${container_id}" "core/ledger/kvledger/txmgmt/statedb/mkv" "${working_dir}"
        
        # Initialize keys in container
        init_keys_in_container "${container_id}" "${working_dir}" "${password}"
        
        # Verify keys in container
        verify_keys_in_container "${container_id}" "${working_dir}"
        
        # Copy keys to current working directory for MKV access
        copy_keys_to_working_dir "${container_id}" "${working_dir}"
        
        print_status "PASS" "MKV setup completed for ${container_name}"
        
    else
        print_status "WARN" "Container ${container_name} is not running, skipping..."
    fi
}

# Function to copy keys to current working directory
copy_keys_to_working_dir() {
    local container_id=$1
    local working_dir=$2
    
    print_status "INFO" "Copying keys to current working directory..."
    
    # Copy keys from working directory to current directory
    docker exec "${container_id}" bash -c "
        if [ -f '${working_dir}/k1.key' ] && [ -f '${working_dir}/k0.key' ] && [ -f '${working_dir}/encrypted_k1.key' ]; then
            cp '${working_dir}/k1.key' . 2>/dev/null || true
            cp '${working_dir}/k0.key' . 2>/dev/null || true
            cp '${working_dir}/encrypted_k1.key' . 2>/dev/null || true
            echo 'Keys copied to current directory successfully'
        else
            echo 'Warning: Some key files not found in ${working_dir}'
        fi
    "
    
    print_status "PASS" "Keys copied to working directory"
}

# Function to show help
show_help() {
    echo "MKV Key Initialization Script"
    echo "============================="
    echo
    echo "Usage: $0 [OPTIONS] [COMMAND]"
    echo
    echo "Commands:"
    echo "  peers     - Setup MKV in all peer containers"
    echo "  orderer   - Setup MKV in orderer container"
    echo "  all       - Setup MKV in all containers (peers + orderer)"
    echo "  auto      - Auto-initialize MKV keys after network start"
    echo "  help      - Show this help message"
    echo
    echo "Options:"
    echo "  -p, --password PASSWORD  - Use custom password (default: fabric_mkv_password_2025)"
    echo "  -d, --dir DIRECTORY      - Working directory in container (default: /home/chaincode/mkv)"
    echo "  -h, --help               - Show this help message"
    echo
    echo "Examples:"
    echo "  $0 all                           # Setup MKV in all containers"
    echo "  $0 peers -p mypassword          # Setup MKV in peers with custom password"
    echo "  $0 orderer -d /opt/mkv          # Setup MKV in orderer with custom directory"
    echo "  $0 auto                          # Auto-initialize MKV keys after network start"
    echo
}

# Parse command line arguments
PASSWORD="fabric_mkv_password_2025"
WORKING_DIR="/home/chaincode/mkv"
COMMAND=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -p|--password)
            PASSWORD="$2"
            shift 2
            ;;
        -d|--dir)
            WORKING_DIR="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        peers|orderer|all|auto|help)
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

# Function to auto-init keys after network start
auto_init_keys_after_network() {
    local password=${1:-"fabric_mkv_password_2025"}
    local working_dir=${2:-"/home/chaincode/mkv"}
    
    print_status "INFO" "Auto-initializing MKV keys after network start..."
    
    # Wait for containers to be ready
    print_status "INFO" "Waiting for containers to be ready..."
    sleep 10
    
    # Test key_manager.sh locally first
    test_key_manager_locally
    
    # Debug container status
    debug_container_status
    
    # Setup MKV in all containers
    setup_mkv_in_peers "$password" "$working_dir"
    echo
    setup_mkv_in_orderer "$password" "$working_dir"
    
    print_status "PASS" "Auto-initialization completed"
}

# Main execution
main() {
    case "${COMMAND:-help}" in
        peers)
            setup_mkv_in_peers "$PASSWORD" "$WORKING_DIR"
            ;;
        orderer)
            setup_mkv_in_orderer "$PASSWORD" "$WORKING_DIR"
            ;;
        all)
            setup_mkv_in_peers "$PASSWORD" "$WORKING_DIR"
            echo
            setup_mkv_in_orderer "$PASSWORD" "$WORKING_DIR"
            ;;
        auto)
            auto_init_keys_after_network "$PASSWORD" "$WORKING_DIR"
            ;;
        help|*)
            show_help
            ;;
    esac
    
    echo
    print_status "INFO" "MKV key initialization completed at $(date)"
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
