# Hyperledger Fabric với MKV256 Encryption cho StateDB và Private Data

## Tổng quan

Dự án này tích hợp thuật toán mã hóa MKV256 vào Hyperledger Fabric để mã hóa dữ liệu trong StateDB và Private Data Collections. Hệ thống sử dụng **KeyManager singleton** với `sync.Once` để quản lý khóa mã hóa tự động, cung cấp giải pháp bảo mật cao cấp cho dữ liệu blockchain.

## Tính năng chính

- ✅ **Mã hóa StateDB**: Tự động mã hóa/giải mã dữ liệu trạng thái
- ✅ **Mã hóa Private Data**: Bảo vệ dữ liệu nhạy cảm trong collections
- ✅ **MKV256 Algorithm**: Thuật toán mã hóa 256-bit block cipher
- ✅ **KeyManager tự động**: Quản lý khóa singleton với `sync.Once`
- ✅ **Password từ file**: Đọc password từ `password.txt` (fallback: "kmasc")
- ✅ **Tích hợp CGO**: Hiệu suất cao với thư viện C tùy chỉnh
- ✅ **Tương thích Fabric 3.1.1**: Hoạt động với tất cả chaincode hiện có

## Kiến trúc mã hóa

### StateDB Encryption
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Password      │    │   K0 (32 bytes) │    │   K1 (32 bytes) │
│   (from file)   │───▶│   (PBKDF2)      │───▶│   (Random)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │                        │
                              ▼                        ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │ Encrypted K1    │    │ Encrypted Data  │
                       │ (Stored)        │    │ (StateDB)       │
                       └─────────────────┘    └─────────────────┘
```

### Private Data Encryption
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Password      │    │   K0 (32 bytes) │    │   K1 (32 bytes) │
│   (from file)   │───▶│   (PBKDF2)      │───▶│   (Random)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │                        │
                              ▼                        ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │ Encrypted K1    │    │ Encrypted Data  │
                       │ (Stored)        │    │ (Private Data)  │
                       └─────────────────┘    └─────────────────┘
```

## Scripts chính

### 🎯 **Demo và Development**
- **`quick-start.sh`** - Setup hoàn chỉnh cho demo và development
  - Build tất cả libraries
  - Test MKV encryption
  - Khởi động test network
  - Phù hợp cho testing và development

### 🚀 **Production và Deployment**
- **`startup.sh`** - Khởi động network thật cho production
  - Khởi động network với cấu hình production
  - Sử dụng MKV encryption đã được test
  - Phù hợp cho deployment thực tế

## Cài đặt nhanh

### Phương pháp 1: Demo và Development (Khuyến nghị cho testing)
```bash
git clone <your-repo-url>
cd fabric-3.1.1
./quick-start.sh
```

### Phương pháp 2: Production deployment
```bash
# Setup environment trước
./setup-environment.sh

# Build libraries
./build-all-libraries.sh

# Test encryption
./test-mkv.sh

# Khởi động production network
./startup.sh
```

## Cách hoạt động của MKV Encryption

### 1. **Khởi tạo tự động**
```go
// Lần đầu gọi sẽ tự động khởi tạo KeyManager
keyManager := mkv.GetKeyManager()

// Tự động tạo:
// - K1: Khóa mã hóa ngẫu nhiên 32 bytes
// - K0: Từ password (đọc từ file) bằng PBKDF2
// - Mã K1 bằng K0 và lưu vào encrypted_k1.key
```

### 2. **Mã hóa StateDB**
```go
// Trong value_encoding.go
encryptedValue := mkv.EncryptValueMKV(v.Value)
encryptedMetadata := mkv.EncryptValueMKV(v.Metadata)

// Dữ liệu được tự động mã hóa khi lưu vào StateDB
// và tự động giải mã khi đọc ra
```

### 3. **Mã hóa Private Data**
```go
// Trong store.go
encryptedPrivateData := mkv.EncryptValueMKV(privateData)

// Private data được mã hóa trước khi lưu vào collection
// và giải mã khi truy xuất
```

## Cài đặt chi tiết

### Bước 1: Cài đặt dependencies
```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y build-essential libssl-dev git curl

# Cài đặt Go 1.24.4
wget https://go.dev/dl/go1.24.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.4.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### Bước 2: Build MKV libraries
```bash
# Build MKV library với KeyManager
cd core/ledger/kvledger/txmgmt/statedb/mkv
make clean && make

# Tạo file password
echo "kmasc" > password.txt

# Test hệ thống
go test -v
cd ../../../../../../..
```

### Bước 3: Build Fabric với MKV
```bash
export CGO_ENABLED=1
make clean
make native
```

## Sử dụng

### Demo và Testing (quick-start.sh)
```bash
# Setup hoàn chỉnh cho demo
./quick-start.sh

# Kiểm tra logs encryption
docker exec peer0.org1.example.com cat /tmp/state_mkv.log

# Test chaincode với encryption
peer chaincode query -C mychannel -n basic -c '{"function":"ReadAsset","Args":["asset1"]}'
```

### Production (startup.sh)
```bash
# Khởi động production network
./startup.sh

# Kiểm tra encryption hoạt động
docker logs -f peer0.org1.example.com | grep -i "ENCRYPT\|DECRYPT"
```

## Kiểm tra encryption

### Kiểm tra StateDB encryption
```bash
# Xem logs encryption
docker exec peer0.org1.example.com cat /tmp/state_mkv.log

# Kiểm tra dữ liệu đã mã hóa
docker exec peer0.org1.example.com ls -la /var/hyperledger/production/ledgersData/stateLeveldb/
```

### Kiểm tra Private Data encryption
```bash
# Xem private data collections
peer lifecycle chaincode queryinstalled

# Kiểm tra dữ liệu trong collection
peer chaincode query -C mychannel -n basic -c '{"function":"ReadPrivateAsset","Args":["asset1"]}'
```

## Cấu trúc project

```
fabric-3.1.1/
├── core/ledger/kvledger/txmgmt/statedb/
│   ├── statedb.go          # Go wrapper với CGO
│   ├── encrypt.c           # C functions cho AES (legacy)
│   ├── encrypt.h           # Header file AES
│   ├── Makefile            # Build script AES
│   └── mkv/                # MKV encryption module với KeyManager
│       ├── mkv.go          # Go wrapper với hệ thống quản lý khóa tự động
│       ├── key_manager.go  # KeyManager singleton với sync.Once
│       ├── mkv.c           # C functions cho MKV
│       ├── mkv.h           # Header file cho MKV
│       ├── MKV256.c        # MKV256 algorithm implementation
│       ├── MKV256.h        # MKV256 header
│       ├── Makefile        # Build script cho MKV
│       ├── README.md       # Chi tiết MKV
│       └── KEY_MANAGER_README.md # Chi tiết KeyManager
├── quick-start.sh          # Setup demo và development
├── startup.sh              # Khởi động production network
├── build-all-libraries.sh  # Build tất cả libraries
├── test-mkv.sh             # Test MKV encryption
└── README.md               # Tài liệu này
```

## Troubleshooting

### Lỗi KeyManager không khởi tạo
```bash
cd core/ledger/kvledger/txmgmt/statedb/mkv

# Kiểm tra file password
ls -la password.txt

# Tạo file password nếu chưa có
echo "kmasc" > password.txt

# Test hệ thống
go test -v
```

### Lỗi MKV library không build
```bash
cd core/ledger/kvledger/txmgmt/statedb/mkv

# Dọn dẹp và build lại
make clean && make

# Kiểm tra library
ldd libmkv.so
```

### Lỗi encryption trong StateDB
```bash
# Kiểm tra logs
docker exec peer0.org1.example.com cat /tmp/state_mkv.log

# Restart peer nếu cần
docker restart peer0.org1.example.com
```

### Lỗi Private Data encryption
```bash
# Kiểm tra cấu hình collection
peer lifecycle chaincode queryinstalled

# Kiểm tra logs
docker logs peer0.org1.example.com | grep -i "private\|collection"
```

## Performance và Bảo mật

### Performance
- **MKV256**: Thuật toán tối ưu cho blockchain
- **KeyManager**: Quản lý khóa hiệu quả với singleton pattern
- **CGO**: Tối thiểu overhead khi gọi C functions

### Bảo mật
- **Khóa K1**: Sinh ngẫu nhiên 32 bytes cho mỗi instance
- **Khóa K0**: Dẫn xuất từ password bằng PBKDF2-HMAC-SHA256
- **Password**: Đọc từ file `password.txt` (fallback: "kmasc")
- **Salt**: Ngẫu nhiên 32 bytes cho PBKDF2

⚠️ **Lưu ý Production**: 
- Tạo file `password.txt` với password mạnh
- Sử dụng HSM hoặc key management system
- Triển khai xoay khóa định kỳ

## Hỗ trợ

- **Demo và Development**: Sử dụng `./quick-start.sh`
- **Production**: Sử dụng `./startup.sh`
- **Testing**: Chạy `./test-mkv.sh`
- **Troubleshooting**: Xem logs trong `/tmp/state_mkv.log`
- **Documentation**: Xem `core/ledger/kvledger/txmgmt/statedb/mkv/README.md`

---

**🎯 Tóm tắt**: Dự án này cung cấp giải pháp mã hóa MKV256 hoàn chỉnh cho Hyperledger Fabric, với KeyManager tự động quản lý khóa và hỗ trợ cả StateDB và Private Data. Sử dụng `quick-start.sh` cho demo/development và `startup.sh` cho production.
