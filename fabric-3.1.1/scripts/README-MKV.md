# MKV Key Management System

Hệ thống quản lý khóa MKV cho Hyperledger Fabric với mã hóa 2 tầng.

## Tổng quan

Hệ thống sử dụng **2 tầng khóa**:
- **K1 (Data Key)**: 32 bytes ngẫu nhiên để mã hóa dữ liệu
- **K0 (Derived Key)**: 32 bytes dẫn xuất từ password để mã hóa K1

## Scripts có sẵn

### 1. `init-mkv-keys-standalone.sh` - Khởi tạo độc lập
```bash
# Khởi tạo với password mặc định
./scripts/init-mkv-keys-standalone.sh

# Khởi tạo với password tùy chỉnh
./scripts/init-mkv-keys-standalone.sh -p mypassword

# Chỉ khởi tạo, không test
./scripts/init-mkv-keys-standalone.sh --no-test
```

### 2. `init-mkv-keys.sh` - Khởi tạo chi tiết
```bash
# Khởi tạo tự động trong tất cả containers
./scripts/init-mkv-keys.sh auto

# Khởi tạo chỉ trong peer containers
./scripts/init-mkv-keys.sh peers

# Khởi tạo chỉ trong orderer container
./scripts/init-mkv-keys.sh orderer

# Khởi tạo với password tùy chỉnh
./scripts/init-mkv-keys.sh auto -p mypassword
```

### 3. `test-mkv-docker.sh` - Test MKV trong Docker
```bash
# Test trong tất cả containers
./scripts/test-mkv-docker.sh all

# Test chỉ trong peer containers
./scripts/test-mkv-docker.sh peers

# Test chỉ trong orderer container
./scripts/test-mkv-docker.sh orderer
```

## Quy trình sử dụng

### Bước 1: Đảm bảo network đang chạy
```bash
cd fabric-samples/test-network
./network.sh up
```

### Bước 2: Khởi tạo MKV keys
```bash
cd fabric-3.1.1
./scripts/init-mkv-keys-standalone.sh
```

### Bước 3: Kiểm tra logs
```bash
# Kiểm tra MKV logs
docker exec peer0.org1.example.com cat /root/state_mkv.log

# Kiểm tra key files trong container
docker exec peer0.org1.example.com ls -la /root/mkv/
```

### Bước 4: Test chaincode
```bash
cd fabric-samples/test-network
peer chaincode query -C mychannel -n basic -c '{"function":"ReadAsset","Args":["asset1"]}'
```

## Cấu trúc file trong container

Sau khi khởi tạo, mỗi container sẽ có:
```
/root/mkv/
├── libmkv.so              # MKV library
├── mkv.h                  # Header file
├── key_manager.sh         # Key management script
├── k1.key                 # K1 (Data Key) - plaintext
├── k0.key                 # K0 (Derived Key) - plaintext
├── encrypted_k1.key       # K1 encrypted with K0
└── test_*                 # Test files (nếu có)
```

## Troubleshooting

### Lỗi "Container not running"
```bash
# Kiểm tra containers
docker ps

# Start network nếu cần
cd fabric-samples/test-network
./network.sh up
```

### Lỗi "libmkv.so not found"
```bash
# Build lại MKV library
cd fabric-3.1.1/core/ledger/kvledger/txmgmt/statedb/mkv
make clean && make
```

### Lỗi "Permission denied"
```bash
# Chạy với sudo hoặc thêm user vào docker group
sudo usermod -aG docker $USER
# Logout và login lại
```

### Lỗi "Password incorrect"
```bash
# Xóa keys cũ và khởi tạo lại
docker exec peer0.org1.example.com rm -rf /root/mkv/*
./scripts/init-mkv-keys-standalone.sh
```

## Tích hợp với Quick Start

Script `quick-start.sh` đã được cập nhật để tự động:
1. Build MKV library
2. Test MKV library
3. Khởi tạo MKV keys sau khi start network
4. Test MKV trong Docker containers

## Bảo mật

- **Password**: Không lưu trữ, chỉ dùng để sinh K0
- **K0**: Không lưu trữ, chỉ dẫn xuất từ password khi cần
- **K1**: Được mã hóa bằng K0 trước khi lưu
- **File permissions**: 600 cho key files, 644 cho logs

## Logs

MKV logs được lưu tại:
- `/root/state_mkv.log` trong mỗi container
- Format: `timestamp operation namespace key status [error]`

## Liên hệ

Nếu có vấn đề, kiểm tra:
1. Docker containers đang chạy
2. MKV library đã được build
3. Network đã được start
4. Logs trong containers
