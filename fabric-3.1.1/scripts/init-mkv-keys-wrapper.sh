#!/bin/bash

# MKV Keys Initialization Wrapper Script
# Author: Phong Ngo
# Date: June 15, 2025

set -e

echo "=== MKV Keys Initialization Wrapper ==="
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

# Function to run init-mkv-keys.sh with auto command
run_auto_init() {
    print_status "INFO" "Running MKV keys auto-initialization..."
    
    if [ -f "scripts/init-mkv-keys.sh" ]; then
        chmod +x "scripts/init-mkv-keys.sh"
        ./scripts/init-mkv-keys.sh auto
        if [ $? -eq 0 ]; then
            print_status "PASS" "MKV keys auto-initialization completed"
        else
            print_status "FAIL" "MKV keys auto-initialization failed"
            return 1
        fi
    else
        print_status "FAIL" "scripts/init-mkv-keys.sh not found"
        return 1
    fi
}

# Function to run init-mkv-keys.sh with all command
run_all_init() {
    print_status "INFO" "Running MKV keys initialization for all containers..."
    
    if [ -f "scripts/init-mkv-keys.sh" ]; then
        chmod +x "scripts/init-mkv-keys.sh"
        ./scripts/init-mkv-keys.sh all
        if [ $? -eq 0 ]; then
            print_status "PASS" "MKV keys initialization for all containers completed"
        else
            print_status "FAIL" "MKV keys initialization for all containers failed"
            return 1
        fi
    else
        print_status "FAIL" "scripts/init-mkv-keys.sh not found"
        return 1
    fi
}

# Main execution
main() {
    local command=${1:-"auto"}
    
    case "$command" in
        auto)
            run_auto_init
            ;;
        all)
            run_all_init
            ;;
        *)
            print_status "FAIL" "Unknown command: $command"
            print_status "INFO" "Available commands: auto, all"
            return 1
            ;;
    esac
    
    echo
    print_status "INFO" "MKV keys initialization wrapper completed at $(date)"
}

# Check if script is being run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi

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

# Run main function with first argument
main "$@"
