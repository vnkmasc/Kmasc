
# Đường dẫn gốc test-network
FABRIC_SAMPLES_PATH=~/Kmasc/fabric-3.1.1/fabric-samples/test-network
DEST_PATH=./config/tls

mkdir -p $DEST_PATH

# Copy Orderer CA
cp $FABRIC_SAMPLES_PATH/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt \
   $DEST_PATH/orderer-ca.crt

# Copy Peer0 Org1 CA
cp $FABRIC_SAMPLES_PATH/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
   $DEST_PATH/org1-peer0-ca.crt

# Copy Peer0 Org2 CA
cp $FABRIC_SAMPLES_PATH/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
   $DEST_PATH/org2-peer0-ca.crt

# Copy CA org1 cert
cp $FABRIC_SAMPLES_PATH/organizations/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem \
   $DEST_PATH/ca-org1-cert.pem

echo "✅TLS CA files copied to: $DEST_PATH"
