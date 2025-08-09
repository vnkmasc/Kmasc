#!/bin/bash

# Hyperledger Fabric with MKV - Production Quick Start
# Updated for Production Architecture with Container Integration
# Date: $(date)

set -e

echo "🚀 === Hyperledger Fabric with MKV - Production Quick Start ==="
echo "This script will set up production-ready MKV deployment"
echo "Date: $(date)"
echo

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
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
        "SUCCESS")
            echo -e "${GREEN}🎉 SUCCESS${NC}: $message"
            ;;
        "DEPLOY")
            echo -e "${PURPLE}🚀 DEPLOY${NC}: $message"
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
        print_status "FAIL" "Not in Fabric root directory with MKV integration"
        print_status "INFO" "Current directory: $(pwd)"
        print_status "INFO" "Please run this script from fabric-3.1.1/ directory"
        exit 1
    fi
    print_status "INFO" "Running from correct directory: $(pwd)"
}

# Function to check if deployment already exists
check_existing_deployment() {
    if [ -d "deployment" ]; then
        print_status "WARN" "Existing deployment directory found"
        read -p "Do you want to rebuild? This will overwrite existing deployment (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            print_status "INFO" "Removing existing deployment..."
            rm -rf deployment/
        else
            print_status "INFO" "Using existing deployment"
            return 1
        fi
    fi
    return 0
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

# Step 1: Check prerequisites
step1_check_prerequisites() {
    print_status "DEPLOY" "Checking prerequisites..."
    
    # Check Go
    check_go_version || install_go
    
    # Check Docker
    if ! command -v docker >/dev/null 2>&1; then
        print_status "FAIL" "Docker is not installed. Please install Docker first."
        exit 1
    fi
    print_status "PASS" "Docker is installed"
    
    # Check Docker Compose
    if ! command -v docker-compose >/dev/null 2>&1; then
        print_status "FAIL" "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
    print_status "PASS" "Docker Compose is installed"
    
    # Check OpenSSL
    if ! command -v openssl >/dev/null 2>&1; then
        print_status "FAIL" "OpenSSL is not installed. Please install OpenSSL first."
        exit 1
    fi
    print_status "PASS" "OpenSSL is installed"
    
    echo
}

# Step 2: Build production deployment
step2_build_production() {
    print_status "DEPLOY" "Building production deployment package..."
    
    if script_exists "scripts/build-production-peer.sh"; then
        chmod +x scripts/build-production-peer.sh
        ./scripts/build-production-peer.sh
        print_status "PASS" "Production deployment package built successfully"
    else
        print_status "FAIL" "build-production-peer.sh not found"
        exit 1
    fi
    echo
}

# Step 3: Initialize MKV system
step3_initialize_mkv() {
    print_status "DEPLOY" "Initializing MKV system..."
    
    cd deployment
    
    # Generate secure password
    DEFAULT_PASSWORD="fabric_mkv_$(date +%Y%m%d)_$(openssl rand -hex 8)"
    
    print_status "INFO" "Generated secure default password"
    
    # Start MKV API server
    print_status "INFO" "Starting MKV API server..."
    LD_LIBRARY_PATH=./lib ./bin/mkv-api-server > mkv-api.log 2>&1 &
    API_PID=$!
    
    # Wait for API server to start
    sleep 3
    
    # Check if API server is running
    if curl -f http://localhost:9876/api/v1/health >/dev/null 2>&1; then
        print_status "PASS" "MKV API server started successfully"
    else
        print_status "FAIL" "MKV API server failed to start"
        cat mkv-api.log
        exit 1
    fi
    
    # Initialize MKV system
    print_status "INFO" "Initializing MKV system with secure password..."
    cd mkv-keys
    cp ../lib/libmkv.so .
    echo "$DEFAULT_PASSWORD" | LD_LIBRARY_PATH=. ../bin/mkv_client.sh init
    
    if [ $? -eq 0 ]; then
        print_status "PASS" "MKV system initialized successfully"
        echo "$DEFAULT_PASSWORD" > ../mkv-default-password.txt
        chmod 600 ../mkv-default-password.txt
        print_status "INFO" "Default password saved to mkv-default-password.txt"
    else
        print_status "FAIL" "MKV system initialization failed"
        exit 1
    fi
    
    cd ..
    echo
}

# Step 4: Setup fabric samples and network
step4_setup_fabric_network() {
    print_status "DEPLOY" "Setting up Fabric network..."
    
    cd ..
    
    # Download fabric-samples if not exists
    if [ ! -d "fabric-samples" ]; then
        print_status "INFO" "Downloading fabric-samples..."
        curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh
        chmod +x install-fabric.sh
        ./install-fabric.sh d s b
        print_status "PASS" "fabric-samples downloaded"
    else
        print_status "INFO" "fabric-samples already exists"
    fi
    
    echo
}

# Step 5: Test deployment
step5_test_deployment() {
    print_status "DEPLOY" "Testing deployment..."
    
    cd deployment
    
    # Test MKV system
    print_status "INFO" "Testing MKV system..."
    
    # Test system status
    ./bin/mkv_client.sh status
    if [ $? -eq 0 ]; then
        print_status "PASS" "MKV system status check passed"
    else
        print_status "FAIL" "MKV system status check failed"
        exit 1
    fi
    
    # Test password validation
    DEFAULT_PASSWORD=$(cat mkv-default-password.txt)
    echo "$DEFAULT_PASSWORD" | ./bin/mkv_client.sh test
    if [ $? -eq 0 ]; then
        print_status "PASS" "Password validation test passed"
    else
        print_status "FAIL" "Password validation test failed"
        exit 1
    fi
    
    # Test password change
    NEW_PASSWORD="test_password_$(date +%H%M%S)"
    echo -e "$DEFAULT_PASSWORD\n$NEW_PASSWORD" | ./bin/mkv_client.sh change
    if [ $? -eq 0 ]; then
        print_status "PASS" "Password change test passed"
        echo "$NEW_PASSWORD" > mkv-current-password.txt
        chmod 600 mkv-current-password.txt
    else
        print_status "FAIL" "Password change test failed"
        exit 1
    fi
    
    cd ..
    echo
}

# Step 6: Show deployment summary
step6_show_summary() {
    print_status "SUCCESS" "Production deployment completed successfully!"
    echo
    echo "🎉 === DEPLOYMENT SUMMARY ==="
    echo
    echo "📁 Deployment Location: $(pwd)/deployment/"
    echo "📊 Components:"
    echo "   ✅ Fabric Peer with MKV integration (61MB)"
    echo "   ✅ MKV API Server (8.5MB)"
    echo "   ✅ MKV Client Tools"
    echo "   ✅ MKV Encryption Library (70KB)"
    echo "   ✅ Docker Compose Configuration"
    echo "   ✅ Production Scripts"
    echo
    echo "🔐 Security:"
    echo "   ✅ PBKDF2-HMAC-SHA256 key derivation"
    echo "   ✅ Runtime password management"
    echo "   ✅ Secure API endpoints"
    echo "   ✅ Container isolation"
    echo
    echo "🚀 Next Steps:"
    echo "   1. cd deployment/"
    echo "   2. ./deploy.sh                    # Deploy with Docker"
    echo "   3. Change default password:"
    echo "      ./bin/mkv_client.sh change \"\$(cat mkv-current-password.txt)\" \"your_secure_password\""
    echo
    echo "📊 Management Commands:"
    echo "   • Health check: curl http://localhost:9876/api/v1/health"
    echo "   • System status: ./bin/mkv_client.sh status"
    echo "   • Change password: ./bin/mkv_client.sh change \"old\" \"new\""
    echo "   • Test password: echo \"password\" | ./bin/mkv_client.sh test"
    echo
    echo "📚 Documentation: deployment/README.md"
    echo
    print_status "SUCCESS" "Ready for production deployment!"
}

# Main execution
main() {
    print_status "DEPLOY" "Starting production quick start process..."
    echo
    
    ensure_correct_directory
    
    # Check if deployment already exists
    if ! check_existing_deployment; then
        cd deployment
        step6_show_summary
        return 0
    fi
    
    step1_check_prerequisites
    step2_build_production
    step3_initialize_mkv
    step4_setup_fabric_network
    step5_test_deployment
    step6_show_summary
}

# Show what this script will do
echo "🚀 This script will create a PRODUCTION-READY MKV deployment:"
echo
echo "✅ Prerequisites Check:"
echo "   • Go 1.24.4 installation"
echo "   • Docker & Docker Compose"
echo "   • OpenSSL"
echo
echo "✅ Production Build:"
echo "   • MKV library with CGO integration"
echo "   • MKV API server binary"
echo "   • Fabric peer with MKV embedded"
echo "   • Complete deployment package"
echo
echo "✅ System Initialization:"
echo "   • Secure password generation"
echo "   • MKV API server startup"
echo "   • Key management system"
echo "   • Runtime password testing"
echo
echo "✅ Container Integration:"
echo "   • Docker Compose configuration"
echo "   • Volume persistence setup"
echo "   • Health monitoring"
echo "   • Production deployment scripts"
echo
echo "🔐 Security Features:"
echo "   • PBKDF2-HMAC-SHA256 (10,000 iterations)"
echo "   • Runtime password changes"
echo "   • Secure API endpoints"
echo "   • Container isolation"
echo
echo "⏱️  Estimated time: 5-10 minutes"
echo

read -p "Do you want to continue with production deployment? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    main "$@"
else
    echo "Production quick start cancelled."
    exit 0
fi
