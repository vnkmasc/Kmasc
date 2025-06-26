# Hyperledger Fabric Scripts Guide

## Overview

This document describes all the scripts available in the Hyperledger Fabric encryption integration project.

## Quick Start Scripts

### `quick-start.sh` ⭐ (Recommended)
**Purpose**: Complete automated setup from scratch
**What it does**:
- Fixes repository issues
- Sets up environment (Go, OpenSSL, Docker)
- Builds encryption library
- Tests encryption integration
- Builds Fabric with encryption
- Starts test network
- Provides next steps

**Usage**:
```bash
chmod +x quick-start.sh
./quick-start.sh
```

## Environment Setup Scripts

### `setup-environment.sh`
**Purpose**: Install all required dependencies
**What it does**:
- Detects OS and package manager
- Installs build tools, OpenSSL, Go, Docker
- Configures Go environment (CGO_ENABLED=1)
- Creates test script

**Usage**:
```bash
chmod +x setup-environment.sh
./setup-environment.sh
```

### `fix-repositories.sh`
**Purpose**: Fix broken Ubuntu repositories
**What it does**:
- Removes problematic PPAs
- Cleans package cache
- Updates package lists

**Usage**:
```bash
chmod +x fix-repositories.sh
./fix-repositories.sh
```

### `check-environment.sh`
**Purpose**: Comprehensive environment verification
**What it does**:
- Checks Go, GCC, OpenSSL, Docker
- Verifies CGO settings
- Tests Fabric files
- Runs quick tests
- Provides recommendations

**Usage**:
```bash
chmod +x check-environment.sh
./check-environment.sh
```

### `test_environment.sh`
**Purpose**: Simple environment test
**What it does**:
- Quick check of basic tools
- Reports versions and status

**Usage**:
```bash
./test_environment.sh
```

## Build Scripts

### `build-fabric.sh`
**Purpose**: Build Fabric with encryption integration
**What it does**:
- Cleans previous builds
- Builds native binaries
- Copies to fabric-samples/bin/
- Builds Docker image

**Usage**:
```bash
export CGO_ENABLED=1
chmod +x build-fabric.sh
./build-fabric.sh
```

### `fabric-samples-install.sh`
**Purpose**: Install Fabric samples
**What it does**:
- Clones fabric-samples repository
- Downloads install-fabric.sh
- Installs Docker samples and binaries

**Usage**:
```bash
chmod +x fabric-samples-install.sh
./fabric-samples-install.sh
```

## Network Scripts

### `start-network.sh`
**Purpose**: Start test network with encryption
**What it does**:
- Starts test network
- Creates channel
- Deploys basic chaincode
- Tests chaincode functionality

**Usage**:
```bash
chmod +x start-network.sh
./start-network.sh
```

## Encryption-Specific Scripts

### `core/ledger/kvledger/txmgmt/statedb/run_tests.sh`
**Purpose**: Test encryption integration
**What it does**:
- Builds C library
- Runs Go tests
- Performs integration tests
- Runs benchmarks

**Usage**:
```bash
cd core/ledger/kvledger/txmgmt/statedb
chmod +x run_tests.sh
./run_tests.sh
```

### `core/ledger/kvledger/txmgmt/statedb/Makefile`
**Purpose**: Build encryption C library
**What it does**:
- Compiles encrypt.c to libencryption.so
- Links with OpenSSL libraries

**Usage**:
```bash
cd core/ledger/kvledger/txmgmt/statedb
make clean && make
```

## Script Execution Order

### For New Installation
1. `quick-start.sh` (recommended)
   OR
2. `fix-repositories.sh` → `setup-environment.sh` → `build-fabric.sh` → `start-network.sh`

### For Development
1. `check-environment.sh` (verify setup)
2. `core/ledger/kvledger/txmgmt/statedb/run_tests.sh` (test encryption)
3. `build-fabric.sh` (rebuild if needed)

### For Troubleshooting
1. `test_environment.sh` (quick check)
2. `check-environment.sh` (detailed diagnostics)
3. `fix-repositories.sh` (if repository issues)

## Environment Variables

### Required
- `CGO_ENABLED=1` - Enable CGO for encryption library

### Optional
- `GOPATH` - Go workspace path
- `GOROOT` - Go installation path

## Common Issues and Solutions

### Repository Errors
```bash
./fix-repositories.sh
```

### CGO Not Enabled
```bash
export CGO_ENABLED=1
```

### Build Failures
```bash
make clean
go mod tidy
./build-fabric.sh
```

### Docker Issues
```bash
sudo systemctl start docker
sudo usermod -aG docker $USER
newgrp docker
```

## Script Dependencies

### System Requirements
- Ubuntu 20.04+ or equivalent
- sudo privileges
- Internet connection

### External Dependencies
- Git
- curl/wget
- Docker (installed by setup script)

## Logging and Monitoring

### Check Encryption Activity
```bash
docker logs -f peer0.org1.example.com | grep -i encrypt
docker logs -f peer0.org1.example.com | grep -i decrypt
```

### Monitor Network
```bash
cd fabric-samples/test-network
./monitordocker.sh
```

## Security Notes

- All scripts are designed for development/demo use
- Production deployment requires additional security measures
- Encryption keys are hardcoded for demonstration only
- Review scripts before running in production environments

---

**Note**: All scripts include error handling and will provide clear feedback on success or failure. Check the output for any warnings or errors that need attention. 