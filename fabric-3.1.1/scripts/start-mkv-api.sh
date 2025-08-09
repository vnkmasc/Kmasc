#!/bin/bash

# Start MKV API Server Script
# This script starts the MKV API server and initializes it with default password

set -e

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

# Get root directory
ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
API_SERVER_DIR="$ROOT_DIR/core/ledger/kvledger/txmgmt/statedb/mkv-api-server"

echo "=== Starting MKV API Server ==="
echo "Date: $(date)"
echo

# Build API server if not exists
if [ ! -f "$API_SERVER_DIR/mkv-api-server" ]; then
    print_status "INFO" "Building MKV API Server..."
    cd "$API_SERVER_DIR"
    go build .
    if [ $? -eq 0 ]; then
        print_status "PASS" "MKV API Server built successfully"
    else
        print_status "FAIL" "Failed to build MKV API Server"
        exit 1
    fi
    cd "$ROOT_DIR"
fi

# Check if server is already running
if curl -s http://localhost:9876/api/v1/health > /dev/null 2>&1; then
    print_status "WARN" "MKV API Server is already running on port 9876"
    print_status "INFO" "Skipping server startup"
    exit 0
fi

# Start API server in background
print_status "INFO" "Starting MKV API Server on port 9876..."
cd "$API_SERVER_DIR"
export LD_LIBRARY_PATH=../mkv
nohup ./mkv-api-server > mkv_api.log 2>&1 &
API_SERVER_PID=$!
echo $API_SERVER_PID > mkv_api.pid
cd "$ROOT_DIR"

# Wait for server to start
print_status "INFO" "Waiting for server to start..."
sleep 3

# Test if server is running
MAX_RETRIES=10
RETRY_COUNT=0
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if curl -s http://localhost:9876/api/v1/health > /dev/null 2>&1; then
        print_status "PASS" "MKV API Server started successfully (PID: $API_SERVER_PID)"
        break
    else
        RETRY_COUNT=$((RETRY_COUNT + 1))
        print_status "INFO" "Retry $RETRY_COUNT/$MAX_RETRIES - waiting for server..."
        sleep 2
    fi
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    print_status "FAIL" "Failed to start MKV API Server after $MAX_RETRIES retries"
    print_status "INFO" "Check log file: $API_SERVER_DIR/mkv_api.log"
    exit 1
fi

# Initialize system with default password
print_status "INFO" "Initializing MKV system with default password..."
INIT_RESPONSE=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -H "X-API-Key: mkv_api_secret_2025" \
    -d '{"password": "fabric_mkv_password_2025"}' \
    http://localhost:9876/api/v1/initialize 2>/dev/null)

if echo "$INIT_RESPONSE" | grep -q '"status":"success"' 2>/dev/null; then
    print_status "PASS" "MKV system initialized with default password"
elif echo "$INIT_RESPONSE" | grep -q "already initialized" 2>/dev/null; then
    print_status "INFO" "MKV system was already initialized"
else
    print_status "WARN" "Failed to initialize MKV system, but server is running"
    print_status "INFO" "You may need to initialize manually"
fi

# Show server information
echo
print_status "INFO" "MKV API Server Information:"
echo "   - URL: http://localhost:9876"
echo "   - API Key: mkv_api_secret_2025"
echo "   - PID: $API_SERVER_PID"
echo "   - Log: $API_SERVER_DIR/mkv_api.log"
echo "   - Default Password: fabric_mkv_password_2025"
echo

print_status "PASS" "MKV API Server setup completed successfully!"
echo
