#!/bin/bash

# Build MKV-enabled Peer Container Script
# This script builds a custom Fabric peer image with MKV integration

set -e

echo "🐳 Building MKV-enabled Peer Container..."

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -d "core/ledger/kvledger/txmgmt/statedb/mkv" ]; then
    echo "❌ Please run this script from the fabric-3.1.1 root directory"
    exit 1
fi

# Check if MKV library is built
if [ ! -f "core/ledger/kvledger/txmgmt/statedb/mkv/libmkv.so" ]; then
    echo "❌ MKV library not found. Please run build-mkv.sh first"
    exit 1
fi

# 1. Build Fabric peer binary with MKV
echo "🏗️ Building Fabric peer with MKV integration..."
export CGO_ENABLED=1
export CGO_CFLAGS="-I$(pwd)/core/ledger/kvledger/txmgmt/statedb/mkv"
export CGO_LDFLAGS="-L$(pwd)/core/ledger/kvledger/txmgmt/statedb/mkv -lmkv"

# Build peer binary
go build -o build/bin/peer ./cmd/peer

if [ $? -ne 0 ]; then
    echo "❌ Failed to build peer binary"
    exit 1
fi

echo "✅ Peer binary built successfully"

# 2. Build MKV API server
echo "🔧 Building MKV API server..."
cd core/ledger/kvledger/txmgmt/statedb/mkv
LD_LIBRARY_PATH=. go build -o mkv-api-server ../mkv-api-server/mkv_api_server.go

if [ $? -ne 0 ]; then
    echo "❌ Failed to build MKV API server"
    exit 1
fi

cd - > /dev/null
echo "✅ MKV API server built successfully"

# 3. Create Dockerfile for MKV-enabled peer
echo "📦 Creating MKV-enabled peer Dockerfile..."
cat > docker/peer/Dockerfile.mkv << 'EOF'
# MKV-enabled Hyperledger Fabric Peer
FROM hyperledger/fabric-peer:2.5

# Install dependencies
USER root
RUN apt-get update && apt-get install -y \
    gcc \
    libc6-dev \
    libssl-dev \
    curl \
    jq \
    && rm -rf /var/lib/apt/lists/*

# Copy MKV-enabled peer binary
COPY build/bin/peer /usr/local/bin/peer

# Copy MKV library and API server
COPY core/ledger/kvledger/txmgmt/statedb/mkv/libmkv.so /usr/local/lib/
COPY core/ledger/kvledger/txmgmt/statedb/mkv/mkv-api-server /usr/local/bin/
COPY core/ledger/kvledger/txmgmt/statedb/mkv/mkv_client.sh /usr/local/bin/

# Create MKV directories
RUN mkdir -p /opt/mkv /var/log/mkv && \
    chown -R fabric:fabric /opt/mkv /var/log/mkv

# Set permissions
RUN chmod +x /usr/local/bin/peer && \
    chmod +x /usr/local/bin/mkv-api-server && \
    chmod +x /usr/local/bin/mkv_client.sh

# Environment variables
ENV LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH
ENV MKV_API_PORT=9876
ENV MKV_KEY_PATH=/opt/mkv

# Copy entrypoint script
COPY scripts/mkv-peer-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/mkv-peer-entrypoint.sh

# Expose MKV API port
EXPOSE 9876

# Switch back to fabric user
USER fabric

# Use custom entrypoint
ENTRYPOINT ["/usr/local/bin/mkv-peer-entrypoint.sh"]
EOF

# 4. Create entrypoint script
echo "📝 Creating peer entrypoint script..."
cat > scripts/mkv-peer-entrypoint.sh << 'EOF'
#!/bin/bash

# MKV-enabled Peer Entrypoint Script
set -e

echo "🚀 Starting MKV-enabled Fabric Peer..."

# Create necessary directories
mkdir -p /opt/mkv /var/log/mkv
chmod 755 /opt/mkv /var/log/mkv

# Check if MKV is initialized
if [ ! -f "/opt/mkv/k1.key" ]; then
    echo "⚠️  MKV not initialized. Initializing with default password..."
    cd /opt/mkv
    echo "fabric_mkv_$(date +%Y%m%d)_$(openssl rand -hex 8)" > /tmp/mkv_password.txt
    echo "$(cat /tmp/mkv_password.txt)" | LD_LIBRARY_PATH=/usr/local/lib mkv_client.sh init
    echo "✅ MKV initialized with password: $(cat /tmp/mkv_password.txt)"
    rm /tmp/mkv_password.txt
fi

# Start MKV API Server in background
echo "🔧 Starting MKV API Server..."
cd /opt/mkv
LD_LIBRARY_PATH=/usr/local/lib mkv-api-server > /var/log/mkv/api.log 2>&1 &
MKV_PID=$!

# Wait for API server to start
sleep 2
if kill -0 $MKV_PID 2>/dev/null; then
    echo "✅ MKV API Server started (PID: $MKV_PID)"
else
    echo "❌ Failed to start MKV API Server"
    exit 1
fi

# Start the original peer process
echo "🏃 Starting Fabric Peer..."
exec peer "$@"
EOF

chmod +x scripts/mkv-peer-entrypoint.sh

# 5. Build the Docker image
echo "🐳 Building MKV-enabled peer Docker image..."
docker build -f docker/peer/Dockerfile.mkv -t hyperledger/fabric-peer:mkv-latest .

if [ $? -eq 0 ]; then
    echo "✅ MKV-enabled peer container built successfully!"
    echo "📦 Docker image: hyperledger/fabric-peer:mkv-latest"
    echo ""
    echo "🔍 Image details:"
    docker images hyperledger/fabric-peer:mkv-latest
else
    echo "❌ Failed to build MKV-enabled peer container"
    exit 1
fi

echo ""
echo "🎉 MKV-enabled peer container is ready!"
echo "   • Image: hyperledger/fabric-peer:mkv-latest"
echo "   • MKV API Port: 9876"
echo "   • Ready for fabric-samples deployment"
