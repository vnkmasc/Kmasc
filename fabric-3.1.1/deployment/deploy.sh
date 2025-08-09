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
    echo "fabric_production_password" | LD_LIBRARY_PATH=. ../bin/mkv_client.sh init
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
