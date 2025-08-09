#!/bin/bash

# MKV API Client Script
# Usage: ./mkv_client.sh [COMMAND] [OPTIONS]

API_BASE_URL="http://localhost:9876/api/v1"
API_KEY="mkv_api_secret_2025"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print functions
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if server is running
check_server() {
    local response=$(curl -s -o /dev/null -w "%{http_code}" "$API_BASE_URL/health")
    if [ "$response" != "200" ]; then
        print_error "MKV API Server is not running. Please start it first:"
        print_info "cd ../mkv-api-server && LD_LIBRARY_PATH=../mkv ./mkv-api-server"
        exit 1
    fi
}

# Health check
health_check() {
    print_info "Checking server health..."
    local response=$(curl -s "$API_BASE_URL/health")
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
}

# Get status
get_status() {
    print_info "Getting system status..."
    local response=$(curl -s -H "X-API-Key: $API_KEY" "$API_BASE_URL/status")
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
}

# Initialize system
initialize_system() {
    local password="$1"
    if [ -z "$password" ]; then
        read -s -p "Enter password for initialization: " password
        echo
    fi
    
    print_info "Initializing MKV system..."
    local response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -H "X-API-Key: $API_KEY" \
        -d "{\"password\": \"$password\"}" \
        "$API_BASE_URL/initialize")
    
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
}

# Change password
change_password() {
    local old_password="$1"
    local new_password="$2"
    
    if [ -z "$old_password" ]; then
        read -s -p "Enter old password: " old_password
        echo
    fi
    
    if [ -z "$new_password" ]; then
        read -s -p "Enter new password: " new_password
        echo
    fi
    
    print_info "Changing password..."
    local response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -H "X-API-Key: $API_KEY" \
        -d "{\"old_password\": \"$old_password\", \"new_password\": \"$new_password\"}" \
        "$API_BASE_URL/change-password")
    
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
}

# Test password
test_password() {
    local password="$1"
    
    if [ -z "$password" ]; then
        read -s -p "Enter password to test: " password
        echo
    fi
    
    print_info "Testing password..."
    local response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -H "X-API-Key: $API_KEY" \
        -d "{\"password\": \"$password\"}" \
        "$API_BASE_URL/test-password")
    
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
}

# Start server
start_server() {
    print_info "Starting MKV API Server..."
    print_info "Server will run on port 9876"
    print_info "Press Ctrl+C to stop"
    echo
    
    cd ../mkv-api-server && LD_LIBRARY_PATH=../mkv ./mkv-api-server
}

# Show help
show_help() {
    echo "MKV API Client"
    echo "=============="
    echo ""
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  start                    - Start the MKV API server"
    echo "  health                   - Check server health"
    echo "  status                   - Get system status"
    echo "  init [password]          - Initialize system with password"
    echo "  change [old] [new]       - Change password"
    echo "  test [password]          - Test if password is valid"
    echo "  help                     - Show this help"
    echo ""
    echo "Examples:"
    echo "  $0 start                 # Start API server"
    echo "  $0 init mypassword       # Initialize with password"
    echo "  $0 change oldpass newpass # Change password"
    echo "  $0 test mypassword       # Test password"
    echo ""
    echo "Environment variables:"
    echo "  MKV_API_PORT             - Server port (default: 9876)"
    echo "  MKV_API_KEY              - API key (default: mkv_api_secret_2025)"
}

# Main script logic
main() {
    # Check if jq is installed for pretty JSON output
    if ! command -v jq &> /dev/null; then
        print_warning "jq not found. Install it for better JSON formatting:"
        print_info "sudo apt-get install jq"
    fi
    
    case "$1" in
        "start")
            start_server
            ;;
        "health")
            check_server
            health_check
            ;;
        "status")
            check_server
            get_status
            ;;
        "init")
            check_server
            initialize_system "$2"
            ;;
        "change")
            check_server
            change_password "$2" "$3"
            ;;
        "test")
            check_server
            test_password "$2"
            ;;
        "help"|"--help"|"-h"|"")
            show_help
            ;;
        *)
            print_error "Unknown command: $1"
            show_help
            exit 1
            ;;
    esac
}

# Run main function
main "$@" 