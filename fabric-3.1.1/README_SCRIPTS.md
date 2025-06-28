# Hyperledger Fabric Scripts Guide

## Overview
This directory contains modular scripts for setting up and managing Hyperledger Fabric with encryption integration. Each script can be run independently or as part of the complete setup process.

## Quick Start
For a complete setup from scratch, run:
```bash
./quick-start.sh
```

## Available Scripts

### Main Scripts
- **`quick-start.sh`** - Complete setup from scratch (recommended)
- **`setup-environment.sh`** - Install all dependencies (Go, Docker, OpenSSL)
- **`build-fabric.sh`** - Build Fabric binaries with encryption
- **`start-network.sh`** - Start the test network

### Utility Scripts
- **`download-fabric-samples.sh`** - Download fabric-samples repository
- **`build-encryption.sh`** - Build encryption library (libencryption.so)
- **`test-encryption.sh`** - Test encryption integration
- **`check-environment.sh`** - Comprehensive environment check
- **`fix-repositories.sh`** - Fix broken package repositories
- **`test_environment.sh`** - Quick environment test
- **`list-scripts.sh`** - List all available scripts

### Testing & Demo Scripts
- **`test-quick-start.sh`** - Test script integrity and configuration
- **`demo-scripts.sh`** - Interactive demo of script functionality

## Individual Usage

### 1. Environment Setup
```bash
# Install all dependencies
./setup-environment.sh

# Check environment
./check-environment.sh

# Quick test
./test_environment.sh
```

### 2. Download Fabric Samples
```bash
./download-fabric-samples.sh
```

### 3. Build Encryption Library
```bash
./build-encryption.sh
```

### 4. Test Encryption
```bash
./test-encryption.sh
```

### 5. Build Fabric
```bash
./build-fabric.sh
```

### 6. Start Network
```bash
./start-network.sh
```

### 7. Testing & Demo
```bash
# Test script integrity
./test-quick-start.sh

# Run interactive demo
./demo-scripts.sh

# List all available scripts
./list-scripts.sh
```

## Script Features

### Automatic Permissions
All scripts automatically set executable permissions when called from `quick-start.sh`.

### Error Handling
Each script includes comprehensive error handling and colored output for better user experience.

### Independent Execution
Every script can be run independently without dependencies on other scripts.

### Environment Checks
Scripts verify the environment before execution and provide helpful error messages.

### Testing & Validation
- `test-quick-start.sh` validates all scripts are properly configured
- `demo-scripts.sh` provides interactive demonstration of functionality

## Troubleshooting

### Common Issues

1. **Permission Denied**
   ```bash
   chmod +x *.sh
   ```

2. **Go Version Issues**
   ```bash
   ./setup-environment.sh  # Will install correct Go version
   ```

3. **Docker Not Running**
   ```bash
   sudo systemctl start docker
   sudo usermod -aG docker $USER  # Then log out and back in
   ```

4. **Repository Issues**
   ```bash
   ./fix-repositories.sh
   ```

### Environment Check
Run the comprehensive environment check:
```bash
./check-environment.sh
```

### View All Scripts
List all available scripts with descriptions:
```bash
./list-scripts.sh
```

### Test Script Integrity
Verify all scripts are properly configured:
```bash
./test-quick-start.sh
```

### Interactive Demo
Run an interactive demo to see scripts in action:
```bash
./demo-scripts.sh
```

## Script Dependencies

### Required Tools
- bash
- curl or wget
- sudo access
- Internet connection

### Optional Tools
- ldd (for library dependency checking)
- pkg-config (for OpenSSL verification)
- timeout (for demo script)

## File Structure
```
fabric-3.1.1/
├── quick-start.sh              # Main setup script
├── setup-environment.sh        # Environment setup
├── build-fabric.sh            # Fabric build
├── start-network.sh           # Network startup
├── download-fabric-samples.sh # Download samples
├── build-encryption.sh        # Build encryption
├── test-encryption.sh         # Test encryption
├── check-environment.sh       # Environment check
├── fix-repositories.sh        # Fix repositories
├── test_environment.sh        # Quick test
├── list-scripts.sh            # List scripts
├── test-quick-start.sh        # Test script integrity
├── demo-scripts.sh            # Interactive demo
└── README_SCRIPTS.md          # This file
```

## Contributing
When adding new scripts:
1. Include proper error handling
2. Add colored output using the standard format
3. Make scripts executable
4. Update this README
5. Test independently and as part of quick-start.sh
6. Add to test-quick-start.sh validation

## Support
For issues or questions:
1. Run `./check-environment.sh` for diagnostics
2. Run `./test-quick-start.sh` to verify script integrity
3. Run `./demo-scripts.sh` to see scripts in action
4. Check individual script logs
5. Ensure all dependencies are installed
6. Verify you're in the correct directory (fabric-3.1.1/) 