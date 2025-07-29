#!/bin/bash
#set -e

CA_URL="http://127.0.0.1:7054"
ORG_BASE="${PWD}/organizations"

function infoln() {
  echo -e "\033[1;32m[INFO]\033[0m $1"
}

function extract_pem() {
  JSON_FILE=$1
  OUT_DIR=$2

  mkdir -p "$OUT_DIR/signcerts" "$OUT_DIR/keystore" "$OUT_DIR/cacerts"

  # Lấy root CA
  jq -r '.caCert' "$JSON_FILE" > "$OUT_DIR/cacerts/ca.pem"

  if jq -e '.signCert' "$JSON_FILE" > /dev/null; then
    jq -r '.signCert' "$JSON_FILE" > "$OUT_DIR/signcerts/cert.pem"
    jq -r '.key' "$JSON_FILE" > "$OUT_DIR/keystore/key.pem"
  elif jq -e '.tlsCert' "$JSON_FILE" > /dev/null; then
    jq -r '.tlsCert' "$JSON_FILE" > "$OUT_DIR/signcerts/cert.pem"
    jq -r '.tlsKey' "$JSON_FILE" > "$OUT_DIR/keystore/key.pem"
  else
    echo "[ERROR] JSON không chứa signCert/tlsCert"
    exit 1
  fi
}

function initCA() {
  infoln "Initializing Root CA..."
  mkdir -p "$ORG_BASE/fabric-ca"
  curl -s -X POST "$CA_URL/initca" -o "$ORG_BASE/fabric-ca/ca.json"
  jq -r '.caCert' "$ORG_BASE/fabric-ca/ca.json" > "$ORG_BASE/fabric-ca/ca-cert.pem"
  jq -r '.caKey' "$ORG_BASE/fabric-ca/ca.json" > "$ORG_BASE/fabric-ca/ca-key.pem"
  infoln "Root CA initialized -> $ORG_BASE/fabric-ca/ca-cert.pem"
}

function registerIdentity() {
  ID=$1
  SECRET=$2
  TYPE=$3
  curl -s -X POST "$CA_URL/register" \
    -H "Content-Type: application/json" \
    -d "{\"id\":\"$ID\",\"secret\":\"$SECRET\",\"type\":\"$TYPE\"}" > /dev/null
  infoln "Registered $TYPE identity: $ID"
}

function enrollIdentity() {
  ID=$1
  SECRET=$2
  OUTDIR=$3
  PROFILE=$4
  HOSTS=$5
  OU=$6

  mkdir -p "$OUTDIR"
  if [ "$PROFILE" == "tls" ]; then
    curl -s -X POST "$CA_URL/tlsenroll" \
      -H "Content-Type: application/json" \
      -d "{\"id\":\"$ID\",\"secret\":\"$SECRET\",\"hosts\":\"$HOSTS\",\"ou\":\"$OU\"}" \
      -o "$OUTDIR/$ID-tls.json"
    extract_pem "$OUTDIR/$ID-tls.json" "$OUTDIR"
  else
    curl -s -X POST "$CA_URL/enroll" \
      -H "Content-Type: application/json" \
      -d "{\"id\":\"$ID\",\"secret\":\"$SECRET\",\"ou\":\"$OU\"}" \
      -o "$OUTDIR/$ID.json"
    extract_pem "$OUTDIR/$ID.json" "$OUTDIR"
  fi
}

function writeNodeOUs() {
  cat > "$1/msp/config.yaml" <<EOF
NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/ca.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/ca.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/ca.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/ca.pem
    OrganizationalUnitIdentifier: orderer
EOF
}

function copyCACert() {
  MSP_DIR=$1
  mkdir -p "$MSP_DIR/cacerts"
  cp "$ORG_BASE/fabric-ca/ca-cert.pem" "$MSP_DIR/cacerts/ca.pem"
}

function setupPeerOrg() {
  ORG=$1
  infoln "=== Setting up $ORG ==="

  ORG_DIR="$ORG_BASE/peerOrganizations/$ORG"
  mkdir -p "$ORG_DIR/msp/cacerts"
  mkdir -p "$ORG_DIR/msp/tlscacerts"
  mkdir -p "$ORG_DIR/tlsca"
  mkdir -p "$ORG_DIR/ca"
  cp "$ORG_BASE/fabric-ca/ca-cert.pem" "$ORG_DIR/msp/cacerts/ca.pem"
  cp "$ORG_BASE/fabric-ca/ca-cert.pem" "$ORG_DIR/msp/tlscacerts/ca.crt"
  cp "$ORG_BASE/fabric-ca/ca-cert.pem" "$ORG_DIR/tlsca/tlsca.$ORG-cert.pem"
  cp "$ORG_BASE/fabric-ca/ca-cert.pem" "$ORG_DIR/ca/ca.$ORG-cert.pem"
  writeNodeOUs "$ORG_DIR"

  registerIdentity peer0 peer0pw peer
  registerIdentity user1 user1pw client
  registerIdentity admin adminpw admin

  # Peer0 MSP
  PEER_DIR="$ORG_DIR/peers/peer0.$ORG"
  mkdir -p "$PEER_DIR/msp"
  enrollIdentity peer0 peer0pw "$PEER_DIR/msp" "" "" "peer"
  copyCACert "$PEER_DIR/msp"
  cp "$ORG_DIR/msp/config.yaml" "$PEER_DIR/msp/config.yaml"

  # Peer0 TLS
  mkdir -p "$PEER_DIR/tls"
  enrollIdentity peer0 peer0pw "$PEER_DIR/tls" "tls" "peer0.$ORG,127.0.0.1,localhost" "peer"
  cp "$PEER_DIR/tls/cacerts/ca.pem" "$PEER_DIR/tls/peer0-tls-ca.pem"
  cp "$PEER_DIR/tls/signcerts/cert.pem" "$PEER_DIR/tls/peer0-tls-cert.pem"
  cp "$PEER_DIR/tls/keystore/key.pem" "$PEER_DIR/tls/peer0-tls-key.pem"
  cp "$PEER_DIR/tls/peer0-tls-ca.pem" "$PEER_DIR/tls/ca.crt"
  cp "$PEER_DIR/tls/peer0-tls-cert.pem" "$PEER_DIR/tls/server.crt"
  cp "$PEER_DIR/tls/peer0-tls-key.pem" "$PEER_DIR/tls/server.key"

  # User MSP
  USER_DIR="$ORG_DIR/users/User1@$ORG"
  mkdir -p "$USER_DIR/msp"
  enrollIdentity user1 user1pw "$USER_DIR/msp" "" "" "client"
  copyCACert "$USER_DIR/msp"
  cp "$ORG_DIR/msp/config.yaml" "$USER_DIR/msp/config.yaml"

  # Admin MSP
  ADMIN_DIR="$ORG_DIR/users/Admin@$ORG"
  mkdir -p "$ADMIN_DIR/msp"
  enrollIdentity admin adminpw "$ADMIN_DIR/msp" "" "" "admin"
  copyCACert "$ADMIN_DIR/msp"
  cp "$ORG_DIR/msp/config.yaml" "$ADMIN_DIR/msp/config.yaml"
}

function setupOrdererOrg() {
  ORG="example.com"
  infoln "=== Setting up orderer.$ORG ==="

  ORG_DIR="$ORG_BASE/ordererOrganizations/$ORG"
  mkdir -p "$ORG_DIR/msp/cacerts"
  mkdir -p "$ORG_DIR/msp/tlscacerts"
  mkdir -p "$ORG_DIR/tlsca"
  cp "$ORG_BASE/fabric-ca/ca-cert.pem" "$ORG_DIR/msp/cacerts/ca.pem"
  cp "$ORG_BASE/fabric-ca/ca-cert.pem" "$ORG_DIR/msp/tlscacerts/tlsca.$ORG-cert.pem"
  cp "$ORG_BASE/fabric-ca/ca-cert.pem" "$ORG_DIR/tlsca/tlsca.$ORG-cert.pem"
  writeNodeOUs "$ORG_DIR"

  registerIdentity orderer ordererpw orderer
  registerIdentity ordererAdmin ordererAdminpw admin

  ORDERER_DIR="$ORG_DIR/orderers/orderer.$ORG"
  mkdir -p "$ORDERER_DIR/msp"
  enrollIdentity orderer ordererpw "$ORDERER_DIR/msp" "" "" "orderer"
  copyCACert "$ORDERER_DIR/msp"
  cp "$ORG_DIR/msp/config.yaml" "$ORDERER_DIR/msp/config.yaml"

  mkdir -p "$ORDERER_DIR/tls"
  enrollIdentity orderer ordererpw "$ORDERER_DIR/tls" "tls" "orderer.$ORG,localhost,127.0.0.1" "orderer"
  cp "$ORDERER_DIR/tls/cacerts/ca.pem" "$ORDERER_DIR/tls/orderer-tls-ca.pem"
  cp "$ORDERER_DIR/tls/signcerts/cert.pem" "$ORDERER_DIR/tls/orderer-tls-cert.pem"
  cp "$ORDERER_DIR/tls/keystore/key.pem" "$ORDERER_DIR/tls/orderer-tls-key.pem"
  cp "$ORDERER_DIR/tls/orderer-tls-ca.pem" "$ORDERER_DIR/tls/ca.crt"
  cp "$ORDERER_DIR/tls/orderer-tls-cert.pem" "$ORDERER_DIR/tls/server.crt"
  cp "$ORDERER_DIR/tls/orderer-tls-key.pem" "$ORDERER_DIR/tls/server.key"

  ADMIN_DIR="$ORG_DIR/users/Admin@$ORG"
  mkdir -p "$ADMIN_DIR/msp"
  enrollIdentity ordererAdmin ordererAdminpw "$ADMIN_DIR/msp" "" "" "admin"
  copyCACert "$ADMIN_DIR/msp"
  cp "$ORG_DIR/msp/config.yaml" "$ADMIN_DIR/msp/config.yaml"
}

# CLI args
if [ $# -gt 0 ]; then
  "$@"
else
  echo "Usage: $0 {initCA|setupPeerOrg org1.example.com|setupPeerOrg org2.example.com|setupOrdererOrg}"
fi
 