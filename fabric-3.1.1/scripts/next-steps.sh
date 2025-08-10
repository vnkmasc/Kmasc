#!/bin/bash

# Step 9: Next steps
# This script shows the next steps after the Fabric network is ready

# Source common functions
if [ -f "scripts/functions.sh" ]; then
    source scripts/functions.sh
else
    # Fallback function if functions.sh is not available
    print_status() {
        local level=$1
        local message=$2
        case $level in
            "INFO") echo "ℹ️  $message" ;;
            "WARN") echo "⚠️  $message" ;;
            "ERROR") echo "❌ $message" ;;
            "PASS") echo "✅ $message" ;;
            *) echo "$message" ;;
        esac
    }
fi

step9_next_steps() {
    print_status "INFO" "Step 9: Next steps..."
    
    echo
    echo "🎉 Hyperledger Fabric network with MKV encryption is ready!"
    echo
    echo "📋 Network Information:"
    echo "   - Network: test-network"
    echo "   - Channel: mychannel"
    echo "   - Chaincode: basic"
    echo "   - MKV Encryption: Enabled"
    echo
    echo "🔑 MKV Library Information:"
    echo "   - Library location: core/ledger/kvledger/txmgmt/statedb/mkv/"
    echo "   - Library file: libmkv.so"
    echo "   - Status: Built and ready for use"
    echo "   - Default password: fabric123"
    echo "   - Keys initialized: k1.key, k0.key, encrypted_k1.key"
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
    echo "📝 Notes:"
    echo "   - MKV library is built and ready for use"
    echo "   - MKV keys initialized with default password: fabric123"
    echo "   - Network is ready for app connection"
    echo "   - For production, consider using persistent volumes"
    echo "   - To change password: cd core/ledger/kvledger/txmgmt/statedb/mkv && ./key_manager.sh change"
}

# Main execution
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    step9_next_steps
fi
