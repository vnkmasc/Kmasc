# Hyperledger Fabric with OpenSSL Encryption Integration

## Tổng quan

Dự án này tích hợp mã hóa/giải mã sử dụng OpenSSL thông qua CGO vào Hyperledger Fabric, thay thế các thư viện mã hóa Go để tận dụng hiệu suất cao của OpenSSL.

## Tính năng

- ✅ Mã hóa/giải mã tự động cho state database
- ✅ Sử dụng OpenSSL AES-256-CBC
- ✅ Tích hợp CGO với thư viện C tùy chỉnh
- ✅ Tương thích với Hyperledger Fabric 3.1.1
- ✅ Test network hoạt động đầy đủ

## Yêu cầu hệ thống

- Ubuntu 20.04+ hoặc tương đương
- Docker và Docker Compose
- Go 1.21+
- GCC và OpenSSL development libraries
- Git

## Cài đặt nhanh

### 1. Clone repository
```bash
git clone <your-repo-url>
cd fabric-3.1.1
```

### 2. Chạy script cài đặt tự động (Khuyến nghị)
```bash
chmod +x quick-start.sh
./quick-start.sh
```

### 3. Hoặc chạy từng bước thủ công
```bash
# Fix repositories (nếu cần)
chmod +x fix-repositories.sh
./fix-repositories.sh

# Setup environment
chmod +x setup-environment.sh
./setup-environment.sh

# Build encryption library
cd core/ledger/kvledger/txmgmt/statedb
make clean && make
cd ../../../../../..

# Build Fabric với encryption
export CGO_ENABLED=1
chmod +x build-fabric.sh
./build-fabric.sh

# Khởi động test network
chmod +x start-network.sh
./start-network.sh
```

## Cài đặt thủ công

### Bước 1: Cài đặt dependencies
```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y build-essential libssl-dev git curl

# Cài đặt Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Cài đặt Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

### Bước 2: Cài đặt Fabric samples
```bash
chmod +x fabric-samples-install.sh
./fabric-samples-install.sh
```

### Bước 3: Build encryption library
```bash
cd core/ledger/kvledger/txmgmt/statedb
make clean && make
cd ../../../../../..
```

### Bước 4: Build Fabric
```bash
make clean
make native
```

## Kiểm tra cài đặt

### Kiểm tra môi trường
```bash
chmod +x check-environment.sh
./check-environment.sh
```

### Test encryption integration
```bash
cd core/ledger/kvledger/txmgmt/statedb
chmod +x run_tests.sh
./run_tests.sh
```

## Sử dụng

### Khởi động network
```bash
./start-network.sh
```

### Kiểm tra logs encryption
```bash
# Xem logs peer với filter encryption
docker logs -f peer0.org1.example.com | grep -i encrypt
docker logs -f peer0.org1.example.com | grep -i decrypt

# Hoặc xem toàn bộ logs
docker logs -f peer0.org1.example.com

# Check log trên peer container
docker exec peer0.org1.example.com cat /root/state_encryption.log
```

### Test chaincode
```bash
# Setup environment
export PATH=${PWD}/fabric-samples/bin:$PATH
export FABRIC_CFG_PATH=$PWD/fabric-samples/config

# Query asset
peer chaincode query -C mychannel -n basic -c '{"function":"ReadAsset","Args":["asset1"]}'
```

## Cấu trúc project

```
fabric-3.1.1/
├── core/ledger/kvledger/txmgmt/statedb/
│   ├── statedb.go          # Go wrapper với CGO
│   ├── encrypt.c           # C functions cho encryption
│   ├── encrypt.h           # Header file
│   ├── Makefile            # Build script
│   └── README_ENCRYPTION.md # Chi tiết encryption
├── build-fabric.sh         # Build script
├── start-network.sh        # Network startup script
├── setup-environment.sh    # Environment setup
├── check-environment.sh    # Environment check
└── README.md              # Tài liệu này
```

## Troubleshooting

### Lỗi Repository (Ubuntu)
Nếu gặp lỗi repository khi chạy `setup-environment.sh`:
```bash
# Fix broken repositories
chmod +x fix-repositories.sh
./fix-repositories.sh

# Sau đó chạy lại setup
./setup-environment.sh
```

### Lỗi CGO
```bash
export CGO_ENABLED=1
go env CGO_ENABLED
```

### Lỗi OpenSSL
```bash
sudo apt-get install libssl-dev
pkg-config --modversion openssl
```

### Lỗi Docker
```bash
sudo systemctl start docker
sudo usermod -aG docker $USER
newgrp docker
```

### Lỗi build
```bash
make clean
go mod tidy
make native
```

## Performance

- OpenSSL AES-256-CBC encryption
- Tự động mã hóa/giải mã khi lưu/đọc state
- Overhead CGO tối thiểu
- Tương thích với tất cả chaincode hiện có

## Bảo mật

⚠️ **Lưu ý**: Khóa mã hóa hiện tại được hardcode cho demo. Trong production cần:
- Sử dụng HSM hoặc key management system
- Khóa động thay vì khóa cố định
- IV ngẫu nhiên cho mỗi lần mã hóa

## Đóng góp

1. Fork repository
2. Tạo feature branch
3. Commit changes
4. Push to branch
5. Tạo Pull Request

## License

Apache 2.0 License

## Hỗ trợ

- Tạo issue trên GitHub
- Kiểm tra README_ENCRYPTION.md cho chi tiết kỹ thuật
- Chạy `./check-environment.sh` để debug môi trường
