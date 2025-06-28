#!/bin/bash

# List Available Scripts
# Author: Phong Ngo
# Date: June 15, 2025

echo "=== Available Fabric Scripts ==="
echo "Date: $(date)"
echo "Directory: $(pwd)"
echo

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}📋 Main Scripts:${NC}"
echo "  🚀 quick-start.sh          - Complete setup from scratch (recommended)"
echo "  🔧 setup-environment.sh    - Install dependencies (Go, Docker, OpenSSL)"
echo "  🏗️  build-fabric.sh        - Build Fabric binaries with encryption"
echo "  🌐 start-network.sh        - Start test network"
echo

echo -e "${GREEN}🔧 Utility Scripts:${NC}"
echo "  📥 download-fabric-samples.sh - Download fabric-samples"
echo "  🔐 build-encryption.sh     - Build encryption library"
echo "  🧪 test-encryption.sh      - Test encryption integration"
echo "  🔍 check-environment.sh    - Check environment setup"
echo "  🧹 fix-repositories.sh     - Fix broken repositories"
echo "  ✅ test_environment.sh     - Quick environment test"
echo "  📋 list-scripts.sh         - List all available scripts"
echo "  🧪 test-quick-start.sh     - Test quick start script integrity"
echo

echo -e "${GREEN}📖 Usage Examples:${NC}"
echo "  # Complete setup:"
echo "  ./quick-start.sh"
echo
echo "  # Individual steps:"
echo "  ./setup-environment.sh"
echo "  ./download-fabric-samples.sh"
echo "  ./build-encryption.sh"
echo "  ./build-fabric.sh"
echo "  ./start-network.sh"
echo
echo "  # Testing:"
echo "  ./test_environment.sh"
echo "  ./test-encryption.sh"
echo "  ./check-environment.sh"
echo "  ./test-quick-start.sh"
echo

echo -e "${YELLOW}💡 Tips:${NC}"
echo "  • All scripts can be run independently"
echo "  • Scripts automatically set executable permissions"
echo "  • Check logs for detailed information"
echo "  • Use check-environment.sh to diagnose issues"
echo "  • Run test-quick-start.sh to verify script integrity"
echo

echo -e "${BLUE}📁 Current Directory Contents:${NC}"
ls -la *.sh | grep -E "\.sh$" | awk '{print "  " $9 " (" $5 " bytes)"}'
echo

echo "=== Script List Complete ===" 