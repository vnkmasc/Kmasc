#!/bin/bash

# Start MKV-enabled Test Network Script
# This script starts fabric-samples test network using MKV-enabled peer containers

set -e

echo "🌐 Starting MKV-enabled Test Network..."

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Please run this script from the fabric-3.1.1 root directory"
    exit 1
fi

# Check if fabric-samples exists
if [ ! -d "fabric-samples" ]; then
    echo "❌ fabric-samples directory not found"
    echo "Please run download-fabric-samples.sh first"
    exit 1
fi

# Check if MKV-enabled peer image exists
if ! docker images | grep -q "hyperledger/fabric-peer.*mkv-latest"; then
    echo "❌ MKV-enabled peer image not found"
    echo "Please run build-mkv-container.sh first"
    exit 1
fi

# 1. Go to fabric-samples test-network
echo "📁 Navigating to fabric-samples test-network..."
cd fabric-samples/test-network

# 2. Clean up any existing network
echo "🧹 Cleaning up existing network..."
./network.sh down

# 3. Create docker-compose override for MKV peers
echo "📝 Creating MKV peer configuration..."
cat > docker-compose-mkv-override.yml << 'EOF'
version: '3.8'

services:
  peer0.org1.example.com:
    image: hyperledger/fabric-peer:mkv-latest
    ports:
      - "9876:9876"  # MKV API port
    volumes:
      - ./mkv-data/org1:/opt/mkv
    environment:
      - MKV_API_PORT=9876
      - MKV_API_KEY=mkv_api_secret_2025

  peer0.org2.example.com:
    image: hyperledger/fabric-peer:mkv-latest
    ports:
      - "9877:9876"  # MKV API port (different port for org2)
    volumes:
      - ./mkv-data/org2:/opt/mkv
    environment:
      - MKV_API_PORT=9876
      - MKV_API_KEY=mkv_api_secret_2025
EOF

# 4. Create MKV data directories
echo "📁 Creating MKV data directories..."
mkdir -p mkv-data/org1 mkv-data/org2
chmod 755 mkv-data/org1 mkv-data/org2

# 5. Start the network with MKV-enabled peers
echo "🚀 Starting test network with MKV-enabled peers..."
export COMPOSE_FILE=docker-compose-test-net.yaml:docker-compose-mkv-override.yml

# Start the network
./network.sh up createChannel -ca

if [ $? -ne 0 ]; then
    echo "❌ Failed to start MKV-enabled test network"
    exit 1
fi

# 6. Wait for peers to be ready
echo "⏳ Waiting for MKV API servers to be ready..."
sleep 10

# Check if MKV API servers are responding
echo "🔍 Checking MKV API servers..."
for i in {1..30}; do
    if curl -s http://localhost:9876/api/v1/health >/dev/null 2>&1; then
        echo "✅ Org1 MKV API server is ready"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "⚠️  Org1 MKV API server not responding, but continuing..."
        break
    fi
    sleep 2
done

for i in {1..30}; do
    if curl -s http://localhost:9877/api/v1/health >/dev/null 2>&1; then
        echo "✅ Org2 MKV API server is ready"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "⚠️  Org2 MKV API server not responding, but continuing..."
        break
    fi
    sleep 2
done

# 7. Deploy chaincode
echo "📦 Deploying basic chaincode..."
./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go

if [ $? -eq 0 ]; then
    echo "✅ Chaincode deployed successfully!"
else
    echo "❌ Failed to deploy chaincode"
    exit 1
fi

# 8. Test chaincode
echo "🧪 Testing chaincode..."
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_ADDRESS=localhost:7051

# Initialize ledger
echo "Initializing ledger with sample data..."
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"InitLedger","Args":[]}'

# Query an asset
echo ""
echo "Querying asset1..."
peer chaincode query -C mychannel -n basic -c '{"function":"ReadAsset","Args":["asset1"]}'

if [ $? -eq 0 ]; then
    echo "✅ Chaincode test successful!"
else
    echo "❌ Chaincode test failed"
    exit 1
fi

# 9. Show network information
echo ""
echo "🎉 MKV-enabled test network is ready!"
echo ""
echo "📖 Network Information:"
echo "   • Channel: mychannel"
echo "   • Chaincode: basic (Go)"
echo "   • Orderer: localhost:7050"
echo "   • Peer Org1: localhost:7051 (MKV API: localhost:9876)"
echo "   • Peer Org2: localhost:9051 (MKV API: localhost:9877)"
echo ""
echo "🔐 MKV API Endpoints:"
echo "   • Org1 Health: curl http://localhost:9876/api/v1/health"
echo "   • Org2 Health: curl http://localhost:9877/api/v1/health"
echo "   • API Key: mkv_api_secret_2025"
echo ""
echo "📝 Quick Commands:"
echo "   • Monitor logs: ./monitordocker.sh"
echo "   • Query asset: peer chaincode query -C mychannel -n basic -c '{\"function\":\"ReadAsset\",\"Args\":[\"asset1\"]}'"
echo "   • Stop network: ./network.sh down"
echo ""
echo "🚀 You can now call MKV APIs to manage encryption keys!"
