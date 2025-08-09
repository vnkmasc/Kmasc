#!/bin/bash

# Show Next Steps Script
# This script displays the next steps and usage information after setup completion

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

echo "=== Next Steps and Usage Information ==="
echo "Date: $(date)"
echo

print_status "INFO" "Setup completed successfully!"

echo
echo "🎉 Hyperledger Fabric network with MKV encryption is ready!"
echo

echo "📋 Network Information:"
echo "   - Network: test-network"
echo "   - Channel: mychannel"
echo "   - Chaincode: basic"
echo "   - MKV Encryption: Enabled"
echo

echo "🔑 MKV Keys Information:"
echo "   - Keys location: $(pwd)/"
echo "   - Password: fabric_mkv_password_2025"
echo "   - Files: k1.key, k0.key, encrypted_k1.key"
echo

echo "🌐 MKV API Server:"
echo "   - URL: http://localhost:9876"
echo "   - Status: Running in background"
echo "   - API Key: mkv_api_secret_2025"
echo "   - Log: core/ledger/kvledger/txmgmt/statedb/mkv-api-server/mkv_api.log"
echo

echo "🚀 Quick Test Commands:"
echo "   cd fabric-samples/test-network"
echo "   export PATH=\${PWD}/bin:\${PWD}/../bin:\${PWD}/../../bin:\$PATH"
echo "   export FABRIC_CFG_PATH=\$PWD/../config/"
echo "   export CORE_PEER_TLS_ENABLED=true"
echo "   export CORE_PEER_LOCALMSPID=\"Org1MSP\""
echo "   export CORE_PEER_MSPCONFIGPATH=\${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp"
echo "   export CORE_PEER_TLS_ROOTCERT_FILE=\${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
echo "   export CORE_PEER_ADDRESS=localhost:7051"
echo "   export ORDERER_CA=\${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
echo
echo "   # Test query"
echo "   peer chaincode query -C mychannel -n basic -c '{\"function\":\"ReadAsset\",\"Args\":[\"asset1\"]}'"
echo

echo "🔧 MKV API Management Commands:"
echo "   # Using helper script (recommended)"
echo "   ./scripts/mkv-api.sh health                    # Check server health"
echo "   ./scripts/mkv-api.sh status                    # Get system status"
echo "   ./scripts/mkv-api.sh change OLD_PASS NEW_PASS  # Change password"
echo "   ./scripts/mkv-api.sh test PASSWORD             # Test password"
echo "   ./scripts/mkv-api.sh stop                      # Stop server"
echo
echo "   # Direct curl commands"
echo "   curl http://localhost:9876/api/v1/health"
echo
echo "   curl -H \"X-API-Key: mkv_api_secret_2025\" http://localhost:9876/api/v1/status"
echo
echo "   curl -X POST -H \"Content-Type: application/json\" \\"
echo "        -H \"X-API-Key: mkv_api_secret_2025\" \\"
echo "        -d '{\"old_password\":\"OLD_PASS\",\"new_password\":\"NEW_PASS\"}' \\"
echo "        http://localhost:9876/api/v1/change-password"
echo
echo "   curl -X POST -H \"Content-Type: application/json\" \\"
echo "        -H \"X-API-Key: mkv_api_secret_2025\" \\"
echo "        -d '{\"password\":\"PASS_TO_TEST\"}' \\"
echo "        http://localhost:9876/api/v1/test-password"
echo

echo "🛑 Stop Commands:"
echo "   # Stop API server"
echo "   ./scripts/mkv-api.sh stop"
echo "   # OR manually"
echo "   kill \$(cat core/ledger/kvledger/txmgmt/statedb/mkv-api-server/mkv_api.pid)"
echo
echo "   # Stop Fabric network"
echo "   cd fabric-samples/test-network && ./network.sh down"
echo

echo "📝 Important Notes:"
echo "   - MKV keys are created and ready to use"
echo "   - Network is ready for application connection"
echo "   - API server runs in background for password management"
echo "   - Default password: fabric_mkv_password_2025"
echo "   - For production, change the default password immediately"
echo "   - For production, consider using persistent volumes and secure API keys"
echo

echo "💡 Quick Password Change Example:"
echo "   # Change from default to secure password"
echo "   ./scripts/mkv-api.sh change fabric_mkv_password_2025 \$(openssl rand -hex 16)"
echo

print_status "PASS" "All information displayed. Your Fabric network with MKV encryption is ready to use!"
echo
