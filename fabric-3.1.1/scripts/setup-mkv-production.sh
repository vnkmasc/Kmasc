#!/bin/bash

# Setup MKV for Production Deployment
set -e

echo "🚀 Setting up MKV for Production..."

# 1. Create persistent directories
sudo mkdir -p /opt/fabric/mkv/{keys,logs,config}
sudo mkdir -p /var/log/fabric/mkv

# 2. Set proper ownership
sudo chown -R $(whoami):$(whoami) /opt/fabric/mkv
sudo chown -R $(whoami):$(whoami) /var/log/fabric/mkv

# 3. Build MKV library
echo "📦 Building MKV library..."
cd core/ledger/kvledger/txmgmt/statedb/mkv
make clean && make

# 4. Copy binaries to system locations
echo "📋 Installing MKV components..."
sudo cp libmkv.so /usr/local/lib/
sudo cp mkv_client.sh /usr/local/bin/
sudo chmod +x /usr/local/bin/mkv_client.sh

# 5. Build MKV API server
echo "🔧 Building MKV API server..."
go build -o mkv-api-server ../mkv_api_server.go
sudo cp mkv-api-server /usr/local/bin/
sudo chmod +x /usr/local/bin/mkv-api-server

# 6. Create systemd service
echo "⚙️  Creating systemd service..."
sudo tee /etc/systemd/system/mkv-api.service > /dev/null <<EOF
[Unit]
Description=MKV API Server for Hyperledger Fabric
After=network.target

[Service]
Type=simple
User=$(whoami)
WorkingDirectory=/opt/fabric/mkv/keys
Environment=LD_LIBRARY_PATH=/usr/local/lib
Environment=MKV_API_PORT=9876
Environment=MKV_API_KEY=mkv_api_secret_2025
ExecStart=/usr/local/bin/mkv-api-server
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

# 7. Initialize MKV system
echo "🔐 Initializing MKV system..."
cd /opt/fabric/mkv/keys

# Copy MKV library here
cp /usr/local/lib/libmkv.so .

# Initialize with secure password
if [ -z "$MKV_PASSWORD" ]; then
    echo "Enter secure password for MKV initialization:"
    read -s MKV_PASSWORD
fi

# Initialize keys
LD_LIBRARY_PATH=. /usr/local/bin/mkv_client.sh init "$MKV_PASSWORD"

# 8. Start and enable service
echo "🚀 Starting MKV API service..."
sudo systemctl daemon-reload
sudo systemctl enable mkv-api
sudo systemctl start mkv-api

# 9. Health check
sleep 3
if curl -f http://localhost:9876/api/v1/health >/dev/null 2>&1; then
    echo "✅ MKV API Server is running and healthy!"
else
    echo "❌ MKV API Server health check failed"
    sudo systemctl status mkv-api
    exit 1
fi

echo ""
echo "🎉 MKV Production Setup Complete!"
echo ""
echo "📋 Next steps:"
echo "1. Update docker-compose.yml to mount /opt/fabric/mkv/keys to peer containers"
echo "2. Set environment variable MKV_API_ENDPOINT=http://host.docker.internal:9876"
echo "3. Test password change: mkv_client.sh change \"old_pass\" \"new_pass\""
echo ""
echo "📊 Monitoring:"
echo "- Service status: sudo systemctl status mkv-api"
echo "- Logs: sudo journalctl -u mkv-api -f"
echo "- Health: curl http://localhost:9876/api/v1/health"
