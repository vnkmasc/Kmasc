# Hyperledger Fabric with OpenSSL Encryption Integration

## Tổng quan

Dự án này tích hợp mã hóa/giải mã sử dụng OpenSSL thông qua CGO vào Hyperledger Fabric, thay thế các thư viện mã hóa Go để tận dụng hiệu suất cao của OpenSSL.

## Tính năng

- Mã hóa/giải mã tự động cho state database
- Sử dụng OpenSSL AES-256-CBC
- Tích hợp CGO với thư viện C tùy chỉnh
- Tương thích với Hyperledger Fabric 3.1.1
- Test network hoạt động đầy đủ
- Scripts modular và có thể chạy độc lập

## Đã kiểm thử trên

- **Hệ điều hành:** Ubuntu 24.04 LTS
- **Go:** 1.24.4
- **Docker:** 24.x (Docker Engine v24.x, Docker Compose v2.x)
- **Docker hiện tại:** Docker Engine v28.2.2, Docker Compose v2.36.2

Dự án đã được xác nhận chạy thành công trên các phiên bản phần mềm trên.

## Yêu cầu hệ thống

- Ubuntu 20.04+ hoặc tương đương (Khuyến nghị: Ubuntu 24.04)
- Docker và Docker Compose (Khuyến nghị: Docker Engine 24.x, Docker Compose 2.x)
- Go 1.21+ (Khuyến nghị: Go 1.24.4)
- GCC và OpenSSL development libraries
- Git

## Cài đặt nhanh

### Phương pháp 1: Setup tự động (Khuyến nghị)
```bash
git clone <your-repo-url>
cd fabric-3.1.1
./quick-start.sh
```

### Phương pháp 2: Setup từng bước
Chạy các file lệnh sau theo thứ tự:

```bash
# Bước 1: Sửa repositories (nếu cần)
./scripts/fix-repositories.sh

# Bước 2: Cài đặt môi trường
./scripts/setup-environment.sh

# Bước 3: Tải fabric-samples
./scripts/download-fabric-samples.sh

# Bước 4: Kiểm tra môi trường
./scripts/test_environment.sh

# Bước 5: Build thư viện encryption
./scripts/build-encryption.sh

# Bước 6: Build Fabric
./scripts/build-fabric.sh

# Bước 7: Khởi động test network
./scripts/start-network.sh
```

## Scripts có sẵn

Dự án này bao gồm một bộ scripts modular để dễ dàng setup và quản lý:

### Scripts chính
- `quick-start.sh` - Setup hoàn chỉnh từ đầu
- `setup-environment.sh` - Cài đặt dependencies
- `build-fabric.sh` - Build Fabric với encryption
- `start-network.sh` - Khởi động test network

### Scripts tiện ích
- `download-fabric-samples.sh` - Tải fabric-samples
- `build-encryption.sh` - Build encryption library
- `test-encryption.sh` - Test encryption integration
- `check-environment.sh` - Kiểm tra môi trường
- `list-scripts.sh` - Liệt kê tất cả scripts

**📖 Xem [SCRIPTS.md](SCRIPTS.md) để biết chi tiết về tất cả scripts và cách sử dụng.**

## Cài đặt thủ công

### Bước 1: Cài đặt dependencies
```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y build-essential libssl-dev git curl

# Cài đặt Go
wget https://go.dev/dl/go1.24.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.4.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Cài đặt Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

### Bước 2: Cài đặt Fabric samples
```bash
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
export CGO_ENABLED=1
make clean
make native
```

## Kiểm tra cài đặt

### Kiểm tra môi trường
```bash
./check-environment.sh
```

### Test encryption integration
```bash
cd core/ledger/kvledger/txmgmt/statedb
./run_tests.sh
```

## Sử dụng

### Khởi động network
```bash
./start-network.sh
```

### Kiểm tra logs encryption
```bash
# Check log trên peer container
docker exec peer0.org1.example.com cat /root/state_encryption.log

# Xem logs peer với filter encryption
docker logs -f peer0.org1.example.com | grep -i encrypt
docker logs -f peer0.org1.example.com | grep -i decrypt

# Hoặc xem toàn bộ logs
docker logs -f peer0.org1.example.com
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
├── *.sh                    # Scripts setup và quản lý
├── README_SCRIPTS.md       # Hướng dẫn chi tiết scripts
└── README.md              # Tài liệu này
```

## Troubleshooting

### Lỗi Repository (Ubuntu)
```bash
./fix-repositories.sh
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

**🔍 Xem [SCRIPTS.md](SCRIPTS.md) để biết thêm chi tiết về troubleshooting và các script hỗ trợ.**

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
