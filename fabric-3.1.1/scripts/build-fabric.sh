#!/bin/bash

# Build Hyperledger Fabric Binaries & Docker Images Script
# Author: Phong Ngo
# Date: June 28, 2025

set -e

echo "=== Hyperledger Fabric Build Script ==="
echo "Starting build process..."
echo

# Navigate to fabric source directory
FABRIC_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$FABRIC_DIR"

echo "📁 Current directory: $(pwd)"

# Check Go version
echo "🔍 Checking Go version..."
go version

echo "🧹 Cleaning previous build artifacts..."
make clean

echo "🔨 Building all native Fabric binaries..."
make native

# List of expected binaries
BINARIES=(peer orderer configtxgen configtxlator discover cryptogen osnadmin ledgerutil)

# Check if all binaries exist
BIN_DIR="build/bin"
ALL_BINARIES_OK=true
echo "\n📦 Checking built binaries:"
for bin in "${BINARIES[@]}"; do
    if [ -f "$BIN_DIR/$bin" ]; then
        echo "  ✅ $bin"
    else
        echo "  ❌ $bin NOT FOUND!"
        ALL_BINARIES_OK=false
    fi
done

if [ "$ALL_BINARIES_OK" = false ]; then
    echo "❌ Some binaries are missing! Check the build logs above."
    exit 1
fi

echo "\n✅ All required binaries built successfully!"

# Copy binaries to fabric-samples/bin/
FABRIC_SAMPLES_DIR="$FABRIC_DIR/fabric-samples"
if [ -d "$FABRIC_SAMPLES_DIR" ]; then
    mkdir -p "$FABRIC_SAMPLES_DIR/bin"
    echo "\n📋 Copying new binaries to fabric-samples/bin/..."
    for bin in "${BINARIES[@]}"; do
        cp "$BIN_DIR/$bin" "$FABRIC_SAMPLES_DIR/bin/"
    done
    echo "✅ Binaries copied to fabric-samples/bin/"
else
    echo "⚠️  fabric-samples directory not found. Skipping copy."
fi

echo "\n🔍 Verifying copied binaries:"
if [ -d "$FABRIC_SAMPLES_DIR/bin" ]; then
    for bin in "${BINARIES[@]}"; do
        if [ -f "$FABRIC_SAMPLES_DIR/bin/$bin" ]; then
            echo -n "$bin version: "; "$FABRIC_SAMPLES_DIR/bin/$bin" version || true
        fi
    done
fi

echo "\n=== Building Docker images for Fabric components ==="

# Đơn giản hóa: Luôn dùng docker, không kiểm tra quyền
DOCKER_CMD="docker"

DOCKER_IMAGES=(peer orderer tools ccenv baseos)
DOCKERFILES=(images/peer/Dockerfile images/orderer/Dockerfile images/tools/Dockerfile images/ccenv/Dockerfile images/baseos/Dockerfile)

UBUNTU_VER=${UBUNTU_VER:-24.04}
GO_VER=$(go version | awk '{print $3}' | sed 's/go//')
TARGETARCH=$(go env GOARCH)
TARGETOS=$(go env GOOS)
FABRIC_VER=${FABRIC_VER:-3.1.1}

for i in "${!DOCKER_IMAGES[@]}"; do
    IMAGE="hyperledger/fabric-${DOCKER_IMAGES[$i]}:latest"
    DOCKERFILE="${DOCKERFILES[$i]}"
    echo "\n🐳 Building Docker image: $IMAGE ..."
    $DOCKER_CMD build --build-arg UBUNTU_VER=$UBUNTU_VER --build-arg GO_VER=$GO_VER --build-arg TARGETARCH=$TARGETARCH --build-arg TARGETOS=$TARGETOS --build-arg FABRIC_VER=$FABRIC_VER -t $IMAGE -f $DOCKERFILE .
    if [ $? -eq 0 ]; then
        echo "✅ Docker image $IMAGE built successfully!"
    else
        echo "❌ Failed to build Docker image $IMAGE!"
        exit 1
    fi
done

echo "\n🎉 Fabric binaries and Docker images build completed successfully!"
echo "📍 Source binaries: $FABRIC_DIR/build/bin/"
echo "📍 Copied to: $FABRIC_SAMPLES_DIR/bin/"
echo "📍 Docker images: hyperledger/fabric-{peer,orderer,tools,ccenv,baseos}:latest"
echo

echo "Next steps:"
echo "1. Run './start-network.sh' to start the test network with new binaries and images"
echo "2. Or manually navigate to fabric-samples/test-network/"
echo "3. Test network will now use your newly built binaries and images!"
echo