#!/bin/bash

# Build Hyperledger Fabric Binaries Script
# Author: Phong Ngo
# Date: June 15, 2025

set -e

echo "=== Hyperledger Fabric Build Script ==="
echo "Starting build process..."

# Navigate to fabric source directory
FABRIC_DIR="/home/phongnh/Kmasc/fabric-3.1.1"
cd "$FABRIC_DIR"

echo "📁 Current directory: $(pwd)"

# Check Go version
echo "🔍 Checking Go version..."
go version

# Clean previous build artifacts
echo "🧹 Cleaning previous build artifacts..."
make clean

# Build native binaries
echo "🔨 Building native binaries..."
make native

# Check if build was successful
if [ -d "build/bin" ] && [ "$(ls -A build/bin)" ]; then
    echo "✅ Build completed successfully!"
    echo "📦 Built binaries:"
    ls -la build/bin/
    
    echo ""
    echo "🔍 Binary versions:"
    echo "Peer version:"
    ./build/bin/peer version
    echo ""
    echo "Orderer version:"
    ./build/bin/orderer version
    echo ""
    echo "ConfigTxGen help:"
    ./build/bin/configtxgen --help | head -5
    
    # Copy binaries to fabric-samples/bin/
    FABRIC_SAMPLES_DIR="/home/phongnh/Kmasc/fabric-3.1.1/fabric-samples"
    if [ -d "$FABRIC_SAMPLES_DIR/bin" ]; then
        echo ""
        echo "📋 Copying new binaries to fabric-samples/bin/..."
        
        # Backup existing binaries
        BACKUP_DIR="$FABRIC_SAMPLES_DIR/bin.backup.$(date +%Y%m%d_%H%M%S)"
        cp -r "$FABRIC_SAMPLES_DIR/bin" "$BACKUP_DIR"
        echo "📁 Backed up old binaries to: $BACKUP_DIR"
        
        # Copy new binaries
        cp build/bin/* "$FABRIC_SAMPLES_DIR/bin/"
        echo "✅ New binaries copied successfully!"
        
        # Verify copied binaries
        echo ""
        echo "🔍 Verifying copied binaries:"
        cd "$FABRIC_SAMPLES_DIR"
        echo "Fabric-samples peer version:"
        ./bin/peer version
        echo "Fabric-samples orderer version:"
        ./bin/orderer version
        
        cd "$FABRIC_DIR"
    else
        echo "⚠️  fabric-samples/bin directory not found. Skipping copy."
    fi
    
else
    echo "❌ Build failed! No binaries found in build/bin/"
    exit 1
fi

echo ""
echo "🎉 Fabric binaries build and copy completed successfully!"
echo "📍 Source binaries: $FABRIC_DIR/build/bin/"
echo "📍 Copied to: $FABRIC_SAMPLES_DIR/bin/"
echo ""
echo "Next steps:"
echo "1. Run './start-network.sh' to start the test network with new binaries"
echo "2. Or manually navigate to fabric-samples/test-network/"
echo "3. Test network will now use your newly built binaries!"

echo ""
echo "🐳 Building Docker image for peer (hyperledger/fabric-peer:latest) ..."

# Đặt các biến version từ Makefile hoặc lệnh hệ thống
UBUNTU_VER=${UBUNTU_VER:-24.04}
GO_VER=$(go version | awk '{print $3}' | sed 's/go//')
TARGETARCH=$(go env GOARCH)
TARGETOS=$(go env GOOS)
FABRIC_VER=${FABRIC_VER:-3.1.1}

docker build --build-arg UBUNTU_VER=$UBUNTU_VER --build-arg GO_VER=$GO_VER --build-arg TARGETARCH=$TARGETARCH --build-arg TARGETOS=$TARGETOS --build-arg FABRIC_VER=$FABRIC_VER -t hyperledger/fabric-peer:latest -f images/peer/Dockerfile .
if [ $? -eq 0 ]; then
    echo "✅ Docker image hyperledger/fabric-peer:latest built successfully!"
else
    echo "❌ Failed to build Docker image hyperledger/fabric-peer:latest!"
    exit 1
fi