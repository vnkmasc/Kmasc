#!/bin/bash

# MKV Entrypoint Script for Peer Container
set -e

echo "🚀 Starting Peer with MKV Integration..."

# Create necessary directories
mkdir -p /opt/mkv /tmp/mkv /var/log/mkv

# Set proper permissions
chmod 755 /opt/mkv /tmp/mkv /var/log/mkv

# Check if MKV is initialized
if [ ! -f "/opt/mkv/k1.key" ]; then
    echo "⚠️  MKV not initialized. Please initialize with:"
    echo "   docker exec <container> mkv_client.sh init \"your_password\""
    echo "   Or mount pre-initialized keys to /opt/mkv/"
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

# Health check
if curl -f http://localhost:9876/api/v1/health >/dev/null 2>&1; then
    echo "✅ MKV API Server is healthy"
else
    echo "⚠️  MKV API Server health check failed"
fi

# Start peer with original arguments
echo "🚀 Starting Hyperledger Fabric Peer..."
exec "$@"
