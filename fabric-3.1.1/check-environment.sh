#!/bin/bash

# Hyperledger Fabric Environment Check Script
# Author: Phong Ngo
# Date: June 15, 2025

set -e

echo "=== Hyperledger Fabric Environment Check ==="
echo "Date: $(date)"
echo "Directory: $(pwd)"
echo

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
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
        "DEBUG")
            echo -e "${PURPLE}🔍 DEBUG${NC}: $message"
            ;;
        "SUCCESS")
            echo -e "${CYAN}🎉 SUCCESS${NC}: $message"
            ;;
    esac
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check file exists
file_exists() {
    [ -f "$1" ]
}

# Function to check directory exists
dir_exists() {
    [ -d "$1" ]
}

# Function to get command version
get_version() {
    local cmd=$1
    if command_exists "$cmd"; then
        $cmd --version 2>/dev/null | head -n1 || echo "Version not available"
    else
        echo "Command not found"
    fi
}

# Function to check Go environment
check_go() {
    echo "=== Go Environment Check ==="
    
    if command_exists go; then
        GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        print_status "PASS" "Go found: $GO_VERSION"
        
        # Check Go version compatibility
        GO_MAJOR=$(echo $GO_VERSION | cut -d. -f1)
        GO_MINOR=$(echo $GO_VERSION | cut -d. -f2)
        
        if [ "$GO_MAJOR" -gt 1 ] || ([ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -ge 21 ]); then
            print_status "PASS" "Go version is compatible (1.21+)"
        else
            print_status "WARN" "Go version may be too old. Recommended: 1.21+"
        fi
        
        # Check Go environment variables
        print_status "INFO" "GOPATH: ${GOPATH:-'Not set'}"
        print_status "INFO" "GOROOT: ${GOROOT:-'Not set'}"
        
        # Check Go modules
        if [ -f "go.mod" ]; then
            print_status "PASS" "Go modules enabled (go.mod found)"
        else
            print_status "WARN" "No go.mod found in current directory"
        fi
        
    else
        print_status "FAIL" "Go not found in PATH"
        print_status "INFO" "Install Go from: https://golang.org/dl/"
    fi
    
    echo
}

# Function to check CGO
check_cgo() {
    echo "=== CGO Check ==="
    
    if [ "$CGO_ENABLED" = "1" ]; then
        print_status "PASS" "CGO is enabled"
    else
        print_status "FAIL" "CGO is not enabled"
        print_status "INFO" "Set CGO_ENABLED=1 to enable CGO"
    fi
    
    # Check CGO environment
    print_status "INFO" "CGO_ENABLED: ${CGO_ENABLED:-'Not set'}"
    print_status "INFO" "CC: ${CC:-'Default (gcc)'}"
    print_status "INFO" "CXX: ${CXX:-'Default (g++)'}"
    
    echo
}

# Function to check C compiler
check_compiler() {
    echo "=== C Compiler Check ==="
    
    if command_exists gcc; then
        GCC_VERSION=$(gcc --version | head -n1)
        print_status "PASS" "GCC found: $GCC_VERSION"
        
        # Check if gcc can compile
        echo 'int main() { return 0; }' > /tmp/test_compile.c
        if gcc -o /tmp/test_compile /tmp/test_compile.c 2>/dev/null; then
            print_status "PASS" "GCC can compile C code"
        else
            print_status "FAIL" "GCC cannot compile C code"
        fi
        rm -f /tmp/test_compile.c /tmp/test_compile
        
    else
        print_status "FAIL" "GCC not found"
        print_status "INFO" "Install build-essential package"
    fi
    
    echo
}

# Function to check OpenSSL
check_openssl() {
    echo "=== OpenSSL Check ==="
    
    # Check OpenSSL command line tool
    if command_exists openssl; then
        OPENSSL_VERSION=$(openssl version)
        print_status "PASS" "OpenSSL CLI: $OPENSSL_VERSION"
    else
        print_status "FAIL" "OpenSSL CLI not found"
    fi
    
    # Check OpenSSL development libraries
    if pkg-config --modversion openssl >/dev/null 2>&1; then
        OPENSSL_DEV_VERSION=$(pkg-config --modversion openssl)
        print_status "PASS" "OpenSSL dev libraries: $OPENSSL_DEV_VERSION"
        
        # Check if we can link against OpenSSL
        echo '#include <openssl/evp.h>' > /tmp/test_openssl.c
        echo 'int main() { return 0; }' >> /tmp/test_openssl.c
        if gcc -o /tmp/test_openssl /tmp/test_openssl.c -lssl -lcrypto 2>/dev/null; then
            print_status "PASS" "Can link against OpenSSL libraries"
        else
            print_status "FAIL" "Cannot link against OpenSSL libraries"
        fi
        rm -f /tmp/test_openssl.c /tmp/test_openssl
        
    else
        print_status "FAIL" "OpenSSL development libraries not found"
        print_status "INFO" "Install libssl-dev package"
    fi
    
    echo
}

# Function to check Docker
check_docker() {
    echo "=== Docker Check ==="
    
    if command_exists docker; then
        DOCKER_VERSION=$(docker --version)
        print_status "PASS" "Docker: $DOCKER_VERSION"
        
        # Check if Docker daemon is running
        if docker info >/dev/null 2>&1; then
            print_status "PASS" "Docker daemon is running"
        else
            print_status "FAIL" "Docker daemon is not running"
            print_status "INFO" "Start Docker: sudo systemctl start docker"
        fi
        
        # Check if user is in docker group
        if groups $USER | grep -q docker; then
            print_status "PASS" "User is in docker group"
        else
            print_status "WARN" "User is not in docker group"
            print_status "INFO" "Add user to docker group: sudo usermod -aG docker $USER"
        fi
        
    else
        print_status "FAIL" "Docker not found"
        print_status "INFO" "Install Docker from: https://docs.docker.com/engine/install/"
    fi
    
    # Check Docker Compose
    if command_exists docker-compose; then
        COMPOSE_VERSION=$(docker-compose --version)
        print_status "PASS" "Docker Compose: $COMPOSE_VERSION"
    else
        print_status "FAIL" "Docker Compose not found"
        print_status "INFO" "Install Docker Compose from: https://docs.docker.com/compose/install/"
    fi
    
    echo
}

# Function to check Fabric files
check_fabric_files() {
    echo "=== Fabric Files Check ==="
    
    # Check if we're in the right directory
    if [ -f "go.mod" ] && grep -q "hyperledger/fabric" go.mod; then
        print_status "PASS" "In Hyperledger Fabric directory"
    else
        print_status "WARN" "Not in Hyperledger Fabric directory"
    fi
    
    # Check encryption files
    ENCRYPTION_DIR="core/ledger/kvledger/txmgmt/statedb"
    if dir_exists "$ENCRYPTION_DIR"; then
        print_status "PASS" "Encryption directory exists: $ENCRYPTION_DIR"
        
        if file_exists "$ENCRYPTION_DIR/encrypt.c"; then
            print_status "PASS" "encrypt.c found"
        else
            print_status "FAIL" "encrypt.c not found"
        fi
        
        if file_exists "$ENCRYPTION_DIR/encrypt.h"; then
            print_status "PASS" "encrypt.h found"
        else
            print_status "FAIL" "encrypt.h not found"
        fi
        
        if file_exists "$ENCRYPTION_DIR/Makefile"; then
            print_status "PASS" "Makefile found"
        else
            print_status "FAIL" "Makefile not found"
        fi
        
        if file_exists "$ENCRYPTION_DIR/libencryption.so"; then
            print_status "PASS" "libencryption.so found"
        else
            print_status "WARN" "libencryption.so not found (run make in $ENCRYPTION_DIR)"
        fi
        
    else
        print_status "FAIL" "Encryption directory not found: $ENCRYPTION_DIR"
    fi
    
    # Check build scripts
    if file_exists "build-fabric.sh"; then
        print_status "PASS" "build-fabric.sh found"
    else
        print_status "WARN" "build-fabric.sh not found"
    fi
    
    if file_exists "start-network.sh"; then
        print_status "PASS" "start-network.sh found"
    else
        print_status "WARN" "start-network.sh not found"
    fi
    
    echo
}

# Function to check system resources
check_system() {
    echo "=== System Resources Check ==="
    
    # Check available memory
    MEMORY_GB=$(free -g | awk '/^Mem:/{print $2}')
    if [ "$MEMORY_GB" -ge 4 ]; then
        print_status "PASS" "Memory: ${MEMORY_GB}GB (sufficient)"
    else
        print_status "WARN" "Memory: ${MEMORY_GB}GB (recommended: 4GB+)"
    fi
    
    # Check available disk space
    DISK_GB=$(df -BG . | awk 'NR==2{print $4}' | sed 's/G//')
    if [ "$DISK_GB" -ge 10 ]; then
        print_status "PASS" "Disk space: ${DISK_GB}GB available (sufficient)"
    else
        print_status "WARN" "Disk space: ${DISK_GB}GB available (recommended: 10GB+)"
    fi
    
    # Check CPU cores
    CPU_CORES=$(nproc)
    if [ "$CPU_CORES" -ge 2 ]; then
        print_status "PASS" "CPU cores: $CPU_CORES (sufficient)"
    else
        print_status "WARN" "CPU cores: $CPU_CORES (recommended: 2+)"
    fi
    
    echo
}

# Function to run quick tests
run_quick_tests() {
    echo "=== Quick Tests ==="
    
    # Test Go build
    print_status "INFO" "Testing Go build..."
    if go build ./... 2>/dev/null; then
        print_status "PASS" "Go build successful"
    else
        print_status "FAIL" "Go build failed"
    fi
    
    # Test encryption library build
    if dir_exists "core/ledger/kvledger/txmgmt/statedb"; then
        print_status "INFO" "Testing encryption library build..."
        cd core/ledger/kvledger/txmgmt/statedb
        if make clean && make 2>/dev/null; then
            print_status "PASS" "Encryption library build successful"
        else
            print_status "FAIL" "Encryption library build failed"
        fi
        cd ../../../../../..
    fi
    
    echo
}

# Function to provide recommendations
provide_recommendations() {
    echo "=== Recommendations ==="
    
    local issues=0
    
    # Check for common issues
    if ! command_exists go; then
        print_status "WARN" "Install Go: https://golang.org/dl/"
        ((issues++))
    fi
    
    if [ "$CGO_ENABLED" != "1" ]; then
        print_status "WARN" "Set CGO_ENABLED=1 in your environment"
        ((issues++))
    fi
    
    if ! command_exists gcc; then
        print_status "WARN" "Install build-essential: sudo apt-get install build-essential"
        ((issues++))
    fi
    
    if ! pkg-config --modversion openssl >/dev/null 2>&1; then
        print_status "WARN" "Install OpenSSL dev: sudo apt-get install libssl-dev"
        ((issues++))
    fi
    
    if ! command_exists docker; then
        print_status "WARN" "Install Docker: https://docs.docker.com/engine/install/"
        ((issues++))
    fi
    
    if [ $issues -eq 0 ]; then
        print_status "SUCCESS" "Environment looks good! Ready to build Fabric."
    else
        print_status "INFO" "Found $issues potential issues to address"
    fi
    
    echo
}

# Function to generate summary
generate_summary() {
    echo "=== Environment Summary ==="
    
    local total_checks=0
    local passed_checks=0
    
    # Count checks (this is a simplified version)
    if command_exists go; then ((passed_checks++)); fi; ((total_checks++))
    if [ "$CGO_ENABLED" = "1" ]; then ((passed_checks++)); fi; ((total_checks++))
    if command_exists gcc; then ((passed_checks++)); fi; ((total_checks++))
    if pkg-config --modversion openssl >/dev/null 2>&1; then ((passed_checks++)); fi; ((total_checks++))
    if command_exists docker; then ((passed_checks++)); fi; ((total_checks++))
    if command_exists docker-compose; then ((passed_checks++)); fi; ((total_checks++))
    
    local percentage=$((passed_checks * 100 / total_checks))
    
    print_status "INFO" "Environment check completed: $passed_checks/$total_checks checks passed ($percentage%)"
    
    if [ $percentage -ge 80 ]; then
        print_status "SUCCESS" "Environment is ready for Fabric development!"
    elif [ $percentage -ge 60 ]; then
        print_status "WARN" "Environment has some issues but may work"
    else
        print_status "FAIL" "Environment needs significant setup"
    fi
    
    echo
}

# Main execution
main() {
    echo "Starting environment check..."
    echo
    
    check_go
    check_cgo
    check_compiler
    check_openssl
    check_docker
    check_fabric_files
    check_system
    run_quick_tests
    provide_recommendations
    generate_summary
    
    echo "=== Check Complete ==="
    print_status "INFO" "Environment check completed at $(date)"
    echo
    echo "Next steps:"
    echo "1. Fix any issues identified above"
    echo "2. Run: ./build-fabric.sh"
    echo "3. Run: ./start-network.sh"
    echo "4. Check logs for encryption activity"
}

# Run main function
main "$@" 