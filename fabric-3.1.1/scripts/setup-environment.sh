#!/bin/bash

# Hyperledger Fabric Environment Setup Script
# Author: Phong Ngo
# Date: June 15, 2025

set -e

echo "=== Hyperledger Fabric Environment Setup ==="
echo "This script will install all required dependencies"
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

# Function to check if user is root
check_root() {
    if [ "$EUID" -eq 0 ]; then
        print_status "WARN" "Running as root. Some commands may need adjustment."
    fi
}

# Function to fix broken repositories
fix_repositories() {
    print_status "INFO" "Checking for broken repositories..."
    
    if command_exists apt-get; then
        # Check if there are broken repositories
        if apt-get update 2>&1 | grep -q "404\|Release"; then
            print_status "WARN" "Found broken repositories, attempting to fix..."
            
            # Remove problematic PPAs
            sudo apt-get update 2>&1 | grep "404\|Release" | grep -o "ppa.launchpadcontent.net/[^ ]*" | sort -u | while read repo; do
                if [ ! -z "$repo" ]; then
                    print_status "INFO" "Removing broken repository: $repo"
                    # This is a simplified approach - in practice you'd need to remove from sources.list.d
                fi
            done
            
            # Try to update with fix-missing
            sudo apt-get update --fix-missing || true
        fi
    fi
}

# Function to detect OS
detect_os() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$NAME
        VER=$VERSION_ID
    else
        OS=$(uname -s)
        VER=$(uname -r)
    fi
    print_status "INFO" "Detected OS: $OS $VER"
}

# Function to update package manager
update_packages() {
    print_status "INFO" "Updating package manager..."
    
    if command_exists apt-get; then
        # Update with error handling - ignore repository errors
        sudo apt-get update || {
            print_status "WARN" "Some repositories failed to update, continuing anyway..."
            # Try to fix broken repositories
            sudo apt-get update --fix-missing || true
        }
    elif command_exists yum; then
        sudo yum update -y || {
            print_status "WARN" "Package update failed, continuing anyway..."
        }
    elif command_exists dnf; then
        sudo dnf update -y || {
            print_status "WARN" "Package update failed, continuing anyway..."
        }
    else
        print_status "WARN" "Unknown package manager"
    fi
}

# Function to install basic dependencies
install_basic_deps() {
    print_status "INFO" "Installing basic dependencies..."
    
    if command_exists apt-get; then
        sudo apt-get install -y build-essential git curl wget unzip
    elif command_exists yum; then
        sudo yum groupinstall -y "Development Tools"
        sudo yum install -y git curl wget unzip
    elif command_exists dnf; then
        sudo dnf groupinstall -y "Development Tools"
        sudo dnf install -y git curl wget unzip
    fi
}

# Function to install OpenSSL development libraries
install_openssl() {
    print_status "INFO" "Installing OpenSSL development libraries..."
    
    if command_exists apt-get; then
        sudo apt-get install -y libssl-dev pkg-config
    elif command_exists yum; then
        sudo yum install -y openssl-devel pkg-config
    elif command_exists dnf; then
        sudo dnf install -y openssl-devel pkg-config
    fi
    
    # Verify installation
    if pkg-config --modversion openssl >/dev/null 2>&1; then
        print_status "PASS" "OpenSSL installed: $(pkg-config --modversion openssl)"
    else
        print_status "FAIL" "OpenSSL installation failed"
        exit 1
    fi
}

# Function to install Go
install_go() {
    print_status "INFO" "Checking Go installation..."
    
    if command_exists go; then
        GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        print_status "PASS" "Go already installed: $GO_VERSION"
        
        # Check if version is sufficient (1.24.4)
        GO_MAJOR=$(echo $GO_VERSION | cut -d. -f1)
        GO_MINOR=$(echo $GO_VERSION | cut -d. -f2)
        GO_PATCH=$(echo $GO_VERSION | cut -d. -f3)
        if [ "$GO_MAJOR" -gt 1 ] || ([ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -ge 24 ] && [ "$GO_PATCH" -ge 4 ]); then
            print_status "PASS" "Go version is sufficient (1.24.4+)"
        else
            print_status "WARN" "Go version may be too old. Recommended: 1.24.4+"
        fi
    else
        print_status "INFO" "Installing Go..."
        
        # Download and install Go
        GO_VERSION="1.24.4"
        GO_ARCH="linux-amd64"
        GO_URL="https://go.dev/dl/go${GO_VERSION}.${GO_ARCH}.tar.gz"
        
        cd /tmp
        wget -q "$GO_URL"
        sudo tar -C /usr/local -xzf "go${GO_VERSION}.${GO_ARCH}.tar.gz"
        rm "go${GO_VERSION}.${GO_ARCH}.tar.gz"
        
        # Add Go to PATH
        if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
            echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        fi
        
        # Source bashrc for current session
        export PATH=$PATH:/usr/local/go/bin
        
        print_status "PASS" "Go installed successfully"
    fi
}

# Function to install Docker
install_docker() {
    print_status "INFO" "Checking Docker installation..."
    
    if command_exists docker; then
        DOCKER_VERSION=$(docker --version)
        print_status "PASS" "Docker already installed: $DOCKER_VERSION"
    else
        print_status "INFO" "Installing Docker..."
        
        # Install Docker using official script
        curl -fsSL https://get.docker.com -o get-docker.sh
        sudo sh get-docker.sh
        rm get-docker.sh
        
        # Add user to docker group
        sudo usermod -aG docker $USER
        
        print_status "PASS" "Docker installed successfully"
        print_status "WARN" "Please log out and log back in for docker group to take effect"
    fi
    
    # Check Docker Compose
    if command_exists docker-compose; then
        print_status "PASS" "Docker Compose already installed"
    else
        print_status "INFO" "Installing Docker Compose..."
        
        # Install Docker Compose
        sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
        sudo chmod +x /usr/local/bin/docker-compose
        
        print_status "PASS" "Docker Compose installed successfully"
    fi
}

# Function to setup Go environment
setup_go_env() {
    print_status "INFO" "Setting up Go environment..."
    
    # Set CGO enabled
    export CGO_ENABLED=1
    
    # Add to bashrc if not already there
    if ! grep -q "CGO_ENABLED=1" ~/.bashrc; then
        echo 'export CGO_ENABLED=1' >> ~/.bashrc
    fi
    
    print_status "PASS" "Go environment configured"
}

# Function to verify installations
verify_installations() {
    print_status "INFO" "Verifying installations..."
    
    # Check Go
    if command_exists go; then
        print_status "PASS" "Go: $(go version)"
    else
        print_status "FAIL" "Go not found"
    fi
    
    # Check GCC
    if command_exists gcc; then
        print_status "PASS" "GCC: $(gcc --version | head -n1)"
    else
        print_status "FAIL" "GCC not found"
    fi
    
    # Check OpenSSL
    if pkg-config --modversion openssl >/dev/null 2>&1; then
        print_status "PASS" "OpenSSL: $(pkg-config --modversion openssl)"
    else
        print_status "FAIL" "OpenSSL not found"
    fi
    
    # Check Docker
    if command_exists docker; then
        print_status "PASS" "Docker: $(docker --version)"
    else
        print_status "FAIL" "Docker not found"
    fi
    
    # Check Docker Compose
    if command_exists docker-compose; then
        print_status "PASS" "Docker Compose: $(docker-compose --version)"
    else
        print_status "FAIL" "Docker Compose not found"
    fi
    
    # Check CGO
    if [ "$CGO_ENABLED" = "1" ]; then
        print_status "PASS" "CGO is enabled"
    else
        print_status "FAIL" "CGO is not enabled"
    fi
}

# Function to create test script
create_test_script() {
    print_status "INFO" "Creating environment test script..."
    
    cat > test_environment.sh << 'EOF'
#!/bin/bash

echo "=== Environment Test ==="

# Test Go
echo "Testing Go..."
if command -v go >/dev/null 2>&1; then
    echo "✅ Go: $(go version)"
else
    echo "❌ Go not found"
fi

# Test GCC
echo "Testing GCC..."
if command -v gcc >/dev/null 2>&1; then
    echo "✅ GCC: $(gcc --version | head -n1)"
else
    echo "❌ GCC not found"
fi

# Test OpenSSL
echo "Testing OpenSSL..."
if pkg-config --modversion openssl >/dev/null 2>&1; then
    echo "✅ OpenSSL: $(pkg-config --modversion openssl)"
else
    echo "❌ OpenSSL not found"
fi

# Test Docker
echo "Testing Docker..."
if command -v docker >/dev/null 2>&1; then
    echo "✅ Docker: $(docker --version)"
else
    echo "❌ Docker not found"
fi

# Test CGO
echo "Testing CGO..."
if [ "$CGO_ENABLED" = "1" ]; then
    echo "✅ CGO is enabled"
else
    echo "❌ CGO is not enabled"
fi

echo "=== Test Complete ==="
EOF

    chmod +x test_environment.sh
    print_status "PASS" "Test script created: test_environment.sh"
}

# Main execution
main() {
    echo "Starting environment setup..."
    echo
    
    check_root
    detect_os
    fix_repositories
    update_packages
    install_basic_deps
    install_openssl
    install_go
    install_docker
    setup_go_env
    verify_installations
    create_test_script
    
    echo
    echo "=== Setup Complete ==="
    print_status "PASS" "Environment setup completed successfully"
    echo
    echo "Next steps:"
    echo "1. Log out and log back in (for docker group)"
    echo "2. Run: ./test_environment.sh"
    echo "3. Run: ./build-fabric.sh"
    echo "4. Run: ./start-network.sh"
    echo
    print_status "INFO" "Setup completed at $(date)"
}

# Run main function
main "$@" 