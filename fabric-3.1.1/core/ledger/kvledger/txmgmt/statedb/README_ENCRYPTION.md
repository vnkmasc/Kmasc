# OpenSSL Encryption Integration for Hyperledger Fabric StateDB

## Overview

This module integrates OpenSSL encryption/decryption into Hyperledger Fabric's state database using CGO, replacing Go crypto libraries to leverage OpenSSL's high performance and proven algorithms.

## Features

- ✅ Automatic encryption/decryption of state data
- ✅ OpenSSL AES-256-CBC encryption
- ✅ CGO integration with custom C library
- ✅ Compatible with Hyperledger Fabric 3.1.1
- ✅ Transparent to existing chaincode

## File Structure

```
statedb/
├── statedb.go          # Go wrapper with CGO integration
├── encrypt.c           # C encryption functions
├── encrypt.h           # C header file
├── encrypt_prod.go     # Production encryption wrapper
├── Makefile            # Build script for C library
├── run_tests.sh        # Test runner script
├── statedb_test.go     # Unit tests
├── libencryption.so    # Compiled shared library
└── README_ENCRYPTION.md # This file
```

## Quick Start

### 1. Build Encryption Library
```bash
cd core/ledger/kvledger/txmgmt/statedb
make clean && make
```

### 2. Run Tests
```bash
./run_tests.sh
```

### 3. Build Fabric with Encryption
```bash
cd ../../../../../..
export CGO_ENABLED=1
make clean && make native
```

### 4. Start Test Network
```bash
./start-network.sh
```

### 5. Check Encryption Logs
```bash
# Monitor peer logs for encryption activity
docker logs -f peer0.org1.example.com | grep -i encrypt
docker logs -f peer0.org1.example.com | grep -i decrypt
```

## Usage

### Automatic Encryption/Decryption
The encryption is transparent to applications. Data is automatically encrypted when stored and decrypted when retrieved:

```go
// Data is automatically encrypted when stored
batch.Put("namespace", "key", []byte("sensitive data"), version)

// Data is automatically decrypted when retrieved
value := batch.Get("namespace", "key")
```

### Manual Encryption/Decryption
You can also encrypt/decrypt data manually:

```go
encryptedValue := statedb.EncryptValue([]byte("Hello World"))
decryptedValue := statedb.DecryptValue(encryptedValue)
```

## Testing

### Run All Tests
```bash
./run_tests.sh
```

### Individual Test Commands
```bash
# Build C library
make clean && make

# Run Go tests
go test ./...

# Run with verbose output
go test -v ./...

# Run encryption-specific tests
go test -run TestEncryption ./...

# Run benchmarks
go test -bench=. ./...
```

### Test Results
The test script will verify:
- ✅ C library compilation
- ✅ OpenSSL linking
- ✅ Go package building
- ✅ Unit tests
- ✅ Integration tests
- ✅ Performance benchmarks

## Encryption Algorithm

- **Algorithm**: AES-256-CBC
- **Padding**: PKCS7
- **Key**: 32-byte fixed key (demo only)
- **IV**: Randomly generated for each encryption

## Security Notes

⚠️ **Important**: The current implementation uses a hardcoded encryption key for demonstration purposes only.

For production use:
- Store keys in HSM or key management system
- Use dynamic keys instead of fixed keys
- Implement proper key rotation
- Add random IV for each encryption

## Troubleshooting

### Common Issues

1. **CGO not enabled**
   ```bash
   export CGO_ENABLED=1
   ```

2. **OpenSSL not found**
   ```bash
   sudo apt-get install libssl-dev
   pkg-config --modversion openssl
   ```

3. **Library build fails**
   ```bash
   make clean && make
   ldd libencryption.so
   ```

4. **Go build fails**
   ```bash
   go mod tidy
   go build ./...
   ```

### Environment Check
Run the environment check script from the project root:
```bash
./check-environment.sh
```

## Performance

- OpenSSL is highly optimized for performance
- CGO overhead is minimal
- Compatible with all existing chaincode
- Automatic encryption/decryption with minimal impact

## Integration

The encryption is integrated into the existing Fabric state database interface:
- `UpdateBatch.Put()` - automatically encrypts data
- `UpdateBatch.Get()` - automatically decrypts data
- `VersionedDB` operations - transparent encryption/decryption

## Logging

Encryption/decryption operations are logged for debugging:
- Look for `[ENCRYPT]` and `[DECRYPT]` in peer logs
- Logs include input/output lengths and operation status
- Failed operations fall back to original values

## Development

### Adding New Encryption Algorithms
1. Add functions to `encrypt.c`
2. Update `encrypt.h`
3. Add wrappers in `statedb.go`
4. Update calling logic

### Building for Development
```bash
# Build C library
make clean && make

# Build Go package
go build ./...

# Run tests
go test ./...
```

---

**Note**: This encryption integration is for demonstration and development purposes. For production deployment, implement proper key management and security measures. 