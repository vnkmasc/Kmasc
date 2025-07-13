#!/bin/bash

# Start Hyperledger Fabric Test Network Script
# Author: Phong Ngo  
# Date: June 15, 2025

set -e

echo "=== Hyperledger Fabric Test Network Startup Script ==="

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
TEST_NETWORK_DIR="$(cd "$SCRIPT_DIR/.." && pwd)/fabric-samples/test-network"
cd "$TEST_NETWORK_DIR"

echo "📁 Current directory: $(pwd)"

# Check Docker
echo "🐳 Checking Docker..."
if ! docker --version > /dev/null 2>&1; then
    echo "❌ Docker is not installed or not running!"
    exit 1
fi

if ! docker-compose --version > /dev/null 2>&1; then
    echo "❌ Docker Compose is not installed!"
    exit 1
fi

echo "✅ Docker is ready"

# --- Tự động chuẩn bị config cho test-network ---
CONFIG_DIR="$(cd "$SCRIPT_DIR/.." && pwd)/fabric-samples/config"
CORE_YAML_SRC="$(cd "$SCRIPT_DIR/.." && pwd)/fabric-samples/test-network/addOrg3/compose/docker/peercfg/core.yaml"
CORE_YAML_DEST="$CONFIG_DIR/core.yaml"

if [ ! -d "$CONFIG_DIR" ]; then
  mkdir -p "$CONFIG_DIR"
  echo "[start-network] Đã tạo thư mục config."
fi

if [ ! -f "$CORE_YAML_DEST" ]; then
  cp "$CORE_YAML_SRC" "$CORE_YAML_DEST"
  echo "[start-network] Đã copy core.yaml mẫu vào config."
fi

# Stop any existing network
echo "🛑 Stopping any existing network..."
./network.sh down

# Start the network with channel
echo "🚀 Starting test network..."
./network.sh up createChannel

# Check if network started successfully
if [ $? -eq 0 ]; then
    echo "✅ Test network started successfully!"
    
    echo ""
    echo "📊 Running containers:"
    docker ps --format "table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}"
    
    # Deploy basic chaincode
    echo ""
    echo "📦 Deploying basic asset transfer chaincode..."
    ./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go
    
    if [ $? -eq 0 ]; then
        echo "✅ Chaincode deployed successfully!"
        
        # Test the chaincode
        echo ""
        echo "🧪 Testing chaincode..."
        
        # Setup environment
        export PATH=${PWD}/../bin:$PATH
        export FABRIC_CFG_PATH=$PWD/../config
        export CORE_PEER_TLS_ENABLED=true
        export CORE_PEER_LOCALMSPID="Org1MSP"
        export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
        export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
        export CORE_PEER_ADDRESS=localhost:7051
        
        # Initialize ledger
        echo "Initializing ledger with sample data..."
        peer chaincode invoke -o localhost:7050 \
            --ordererTLSHostnameOverride orderer.example.com \
            --tls \
            --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
            -C mychannel \
            -n basic \
            --peerAddresses localhost:7051 \
            --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
            --peerAddresses localhost:9051 \
            --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
            -c '{"function":"InitLedger","Args":[]}'
        
        sleep 2
        
        # Query an asset
        echo ""
        echo "Querying asset1..."
        RESULT=$(peer chaincode query -C mychannel -n basic -c '{"function":"ReadAsset","Args":["asset1"]}' 2>/dev/null)
        
        if [ $? -eq 0 ] && [ ! -z "$RESULT" ]; then
            echo "✅ Chaincode test successful!"
            echo "📋 Asset1 data: $RESULT"
        else
            echo "⚠️  Chaincode query failed or returned empty result"
        fi
        
        echo ""
        echo "🎉 Test network is ready!"
        echo ""
        echo "📖 Network Information:"
        echo "   • Channel: mychannel"
        echo "   • Chaincode: basic (Go)"
        echo "   • Orderer: localhost:7050"
        echo "   • Peer Org1: localhost:7051"
        echo "   • Peer Org2: localhost:9051"
        echo ""
        echo "📝 Quick Commands:"
        echo "   • Monitor logs: ./monitordocker.sh"
        echo "   • Query asset: peer chaincode query -C mychannel -n basic -c '{\"function\":\"ReadAsset\",\"Args\":[\"asset1\"]}'"
        echo "   • Stop network: ./network.sh down"
        
    else
        echo "❌ Chaincode deployment failed!"
        exit 1
    fi
    
else
    echo "❌ Failed to start test network!"
    exit 1
fi