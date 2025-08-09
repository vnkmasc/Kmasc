#!/bin/bash

# MKV API Management Script
# Usage: ./mkv-api.sh [COMMAND] [OPTIONS]

API_BASE_URL="http://localhost:9876/api/v1"
API_KEY="mkv_api_secret_2025"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
API_SERVER_DIR="$ROOT_DIR/core/ledger/kvledger/txmgmt/statedb/mkv-api-server"

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
        print_error "MKV API Server is not running on port 9876"
        return 1
    fi
    return 0
}

# Start server
start_server() {
    print_info "Starting MKV API Server..."
    
    if check_server; then
        print_warning "API Server is already running"
        return 0
    fi
    
    # Build if needed
    if [ ! -f "$API_SERVER_DIR/mkv-api-server" ]; then
        print_info "Building API server..."
        cd "$API_SERVER_DIR"
        go build .
        cd "$ROOT_DIR"
    fi
    
    # Start server
    cd "$API_SERVER_DIR"
    export LD_LIBRARY_PATH=../mkv
    nohup ./mkv-api-server > mkv_api.log 2>&1 &
    API_SERVER_PID=$!
    echo $API_SERVER_PID > mkv_api.pid
    cd "$ROOT_DIR"
    
    sleep 3
    
    if check_server; then
        print_success "MKV API Server started successfully (PID: $API_SERVER_PID)"
        print_info "Log file: $API_SERVER_DIR/mkv_api.log"
    else
        print_error "Failed to start MKV API Server"
        return 1
    fi
}

# Stop server
stop_server() {
    print_info "Stopping MKV API Server..."
    
    if [ -f "$API_SERVER_DIR/mkv_api.pid" ]; then
        local pid=$(cat "$API_SERVER_DIR/mkv_api.pid")
        if kill -0 "$pid" 2>/dev/null; then
            kill "$pid"
            rm -f "$API_SERVER_DIR/mkv_api.pid"
            print_success "MKV API Server stopped (PID: $pid)"
        else
            print_warning "Process $pid is not running"
            rm -f "$API_SERVER_DIR/mkv_api.pid"
        fi
    else
        print_warning "No PID file found. Trying to find and kill process..."
        pkill -f "mkv-api-server" && print_success "MKV API Server stopped" || print_warning "No running server found"
    fi
}

# Health check
health_check() {
    print_info "Checking server health..."
    if check_server; then
        local response=$(curl -s "$API_BASE_URL/health")
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
        print_success "Server is healthy"
    else
        print_error "Server is not responding"
        return 1
    fi
}

# Get status
get_status() {
    print_info "Getting system status..."
    if check_server; then
        local response=$(curl -s -H "X-API-Key: $API_KEY" "$API_BASE_URL/status")
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        return 1
    fi
}

# Change password
change_password() {
    local old_password="$1"
    local new_password="$2"
    
    if [ -z "$old_password" ] || [ -z "$new_password" ]; then
        print_error "Usage: $0 change <old_password> <new_password>"
        return 1
    fi
    
    if ! check_server; then
        return 1
    fi
    
    print_info "Changing password..."
    local response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -H "X-API-Key: $API_KEY" \
        -d "{\"old_password\": \"$old_password\", \"new_password\": \"$new_password\"}" \
        "$API_BASE_URL/change-password")
    
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
    
    if echo "$response" | grep -q '"status":"success"'; then
        print_success "Password changed successfully"
    else
        print_error "Failed to change password"
        return 1
    fi
}

# Test password
test_password() {
    local password="$1"
    
    if [ -z "$password" ]; then
        print_error "Usage: $0 test <password>"
        return 1
    fi
    
    if ! check_server; then
        return 1
    fi
    
    print_info "Testing password..."
    local response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -H "X-API-Key: $API_KEY" \
        -d "{\"password\": \"$password\"}" \
        "$API_BASE_URL/test-password")
    
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
    
    if echo "$response" | grep -q '"valid":true'; then
        print_success "Password is valid"
    else
        print_error "Password is invalid"
        return 1
    fi
}

# Show help
show_help() {
    echo "MKV API Management Script"
    echo "========================="
    echo ""
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  start                    - Start the MKV API server"
    echo "  stop                     - Stop the MKV API server"
    echo "  restart                  - Restart the MKV API server"
    echo "  health                   - Check server health"
    echo "  status                   - Get system status"
    echo "  change <old> <new>       - Change password"
    echo "  test <password>          - Test if password is valid"
    echo "  help                     - Show this help"
    echo ""
    echo "Examples:"
    echo "  $0 start                           # Start API server"
    echo "  $0 change fabric_mkv_password_2025 mynewpass  # Change password"
    echo "  $0 test mynewpass                  # Test password"
    echo "  $0 stop                            # Stop server"
    echo ""
    echo "API Information:"
    echo "  URL: $API_BASE_URL"
    echo "  API Key: $API_KEY"
    echo "  Server Directory: $API_SERVER_DIR"
}

# Main script logic
main() {
    case "$1" in
        "start")
            start_server
            ;;
        "stop")
            stop_server
            ;;
        "restart")
            stop_server
            sleep 2
            start_server
            ;;
        "health")
            health_check
            ;;
        "status")
            get_status
            ;;
        "change")
            change_password "$2" "$3"
            ;;
        "test")
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
