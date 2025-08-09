#!/bin/bash

# Build Production-Ready Peer with MKV Integration
set -e

echo "🚀 Building Production-Ready Peer with MKV Integration..."

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -d "core/ledger/kvledger/txmgmt/statedb/mkv" ]; then
    echo "❌ Please run this script from the fabric-3.1.1 root directory"
    exit 1
fi

# 1. Build MKV library
echo "📦 Building MKV library..."
cd core/ledger/kvledger/txmgmt/statedb/mkv
make clean && make
if [ $? -ne 0 ]; then
    echo "❌ Failed to build MKV library"
    exit 1
fi

# 2. Build MKV API server
echo "🔧 Building MKV API server..."
LD_LIBRARY_PATH=. go build -o mkv-api-server ../mkv-api-server/mkv_api_server.go
if [ $? -ne 0 ]; then
    echo "❌ Failed to build MKV API server"
    exit 1
fi

# 3. Go back to root and build peer
cd /home/phongnh/go-src/Kmasc/fabric-3.1.1
echo "🏗️ Building Fabric peer with MKV integration..."

# Set CGO flags for MKV integration
export CGO_ENABLED=1
export CGO_CFLAGS="-I$(pwd)/core/ledger/kvledger/txmgmt/statedb/mkv"
export CGO_LDFLAGS="-L$(pwd)/core/ledger/kvledger/txmgmt/statedb/mkv -lmkv"

# Build peer
make clean
make peer

if [ $? -ne 0 ]; then
    echo "❌ Failed to build peer"
    exit 1
fi

# 4. Create deployment directory structure
echo "📁 Creating deployment structure..."
mkdir -p deployment/{bin,lib,config,scripts}

# 5. Copy binaries and libraries
echo "📋 Copying binaries and libraries..."
cp build/bin/peer deployment/bin/
cp core/ledger/kvledger/txmgmt/statedb/mkv/libmkv.so deployment/lib/
cp core/ledger/kvledger/txmgmt/statedb/mkv/mkv-api-server deployment/bin/
cp core/ledger/kvledger/txmgmt/statedb/mkv/mkv_client.sh deployment/bin/
chmod +x deployment/bin/*

# 6. Copy configuration and scripts
echo "📝 Copying configuration files..."
cp scripts/mkv-entrypoint.sh deployment/scripts/
cp docker-compose-mkv.yml deployment/
cp docker/peer/Dockerfile.mkv deployment/

# 7. Create production environment file
echo "⚙️ Creating production environment..."
cat > deployment/.env << EOF
# MKV Configuration
MKV_API_PORT=9876
MKV_API_KEY=mkv_api_secret_production_$(openssl rand -hex 16)
MKV_KEY_PATH=/opt/mkv
LD_LIBRARY_PATH=/usr/local/lib

# Fabric Configuration
FABRIC_LOGGING_SPEC=INFO
CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
CORE_PEER_TLS_ENABLED=true
CORE_PEER_PROFILE_ENABLED=true

# Network Configuration
COMPOSE_PROJECT_NAME=fabric-mkv
FABRIC_CFG_PATH=/etc/hyperledger/fabric
EOF

# 8. Create quick deployment script
echo "🚀 Creating quick deployment script..."
cat > deployment/deploy.sh << 'EOF'
#!/bin/bash

echo "🚀 Deploying Fabric with MKV Integration..."

# Load environment variables
set -a
source .env
set +a

# Create persistent directories
mkdir -p mkv-keys logs

# Check if MKV is initialized
if [ ! -f "mkv-keys/k1.key" ]; then
    echo "⚠️  MKV not initialized. Initializing with default password..."
    echo "🔐 Please change this password in production!"
    
    # Copy MKV library to keys directory
    cp lib/libmkv.so mkv-keys/
    
    # Initialize with default password (CHANGE THIS!)
    cd mkv-keys
    echo "fabric_production_password_$(date +%Y%m%d)" | LD_LIBRARY_PATH=. ../bin/mkv_client.sh init
    cd ..
fi

# Build and start containers
echo "🐳 Building and starting containers..."
docker-compose -f docker-compose-mkv.yml build
docker-compose -f docker-compose-mkv.yml up -d

# Wait for services to start
echo "⏳ Waiting for services to start..."
sleep 10

# Health check
echo "🏥 Performing health checks..."
if curl -f http://localhost:9876/api/v1/health >/dev/null 2>&1; then
    echo "✅ MKV API Server is healthy"
else
    echo "❌ MKV API Server health check failed"
    docker-compose -f docker-compose-mkv.yml logs peer0.org1.example.com
    exit 1
fi

echo ""
echo "🎉 Deployment Complete!"
echo ""
echo "📊 Services:"
echo "- Peer: localhost:7051"
echo "- MKV API: localhost:9876"
echo ""
echo "🔧 Management:"
echo "- Change password: ./bin/mkv_client.sh change \"old_pass\" \"new_pass\""
echo "- Check status: ./bin/mkv_client.sh status"
echo "- View logs: docker-compose -f docker-compose-mkv.yml logs -f"
EOF

chmod +x deployment/deploy.sh

# 9. Create production guide
echo "📚 Creating production guide..."
cat > deployment/README.md << 'EOF'
# MKV Production Deployment

## Quick Start

```bash
# 1. Deploy the system
./deploy.sh

# 2. Change default password (IMPORTANT!)
./bin/mkv_client.sh change "fabric_production_password_$(date +%Y%m%d)" "your_secure_password"

# 3. Test the system
./bin/mkv_client.sh test "your_secure_password"
```

## Directory Structure

```
deployment/
├── bin/                    # Binaries
│   ├── peer               # Fabric peer with MKV
│   ├── mkv-api-server     # MKV API server
│   └── mkv_client.sh      # MKV client tools
├── lib/                   # Libraries
│   └── libmkv.so         # MKV encryption library
├── config/               # Configuration files
├── scripts/              # Helper scripts
├── mkv-keys/            # MKV key files (persistent)
├── logs/                # Application logs
├── .env                 # Environment variables
├── docker-compose-mkv.yml # Docker compose file
└── deploy.sh            # Quick deployment script
```

## Security Notes

1. **Change default passwords immediately**
2. **Secure the MKV API endpoint**
3. **Regular password rotation**
4. **Monitor logs for security events**
5. **Backup key files regularly**

## Monitoring

```bash
# Health check
curl http://localhost:9876/api/v1/health

# System status
./bin/mkv_client.sh status

# View logs
docker-compose -f docker-compose-mkv.yml logs -f
```
EOF

echo ""
echo "✅ Production build completed successfully!"
echo ""
echo "📁 Deployment files created in: deployment/"
echo "📋 Next steps:"
echo "  1. cd deployment/"
echo "  2. ./deploy.sh"
echo "  3. Change default password!"
echo ""
echo "🔐 Security: Remember to change the default password immediately!"
