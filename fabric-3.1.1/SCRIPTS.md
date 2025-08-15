# Hướng dẫn Scripts Hyperledger Fabric

## Tổng quan

Tài liệu này mô tả tất cả các script có sẵn trong dự án tích hợp mã hóa Hyperledger Fabric. Mỗi script có thể chạy độc lập hoặc như một phần của quy trình thiết lập hoàn chỉnh.

## Khởi động Nhanh ⭐

Để thiết lập hoàn chỉnh từ đầu, chạy:
```bash
chmod +x scripts/quick-start.sh
./scripts/quick-start.sh
```

## Scripts Khởi động Nhanh

### `quick-start.sh` ⭐ (Demo và Development)
**Mục đích**: Thiết lập hoàn toàn tự động cho demo và development
**Chức năng**:
- Sửa lỗi repository
- Thiết lập môi trường (Go, OpenSSL, Docker)
- Tải xuống fabric-samples
- Kiểm tra môi trường
- Xây dựng thư viện MKV encryption
- Kiểm tra hệ thống MKV
- Xây dựng Fabric với MKV encryption
- Khởi động mạng thử nghiệm
- Cung cấp các bước tiếp theo

### `startup.sh` 🚀 (Production)
**Mục đích**: Khởi động network thật cho production
**Chức năng**:
- Sửa lỗi repository
- Thiết lập môi trường (Go, OpenSSL, Docker)
- Tải xuống fabric-samples
- Kiểm tra môi trường
- Xây dựng thư viện MKV encryption
- Xây dựng Fabric với MKV encryption

**Cách sử dụng**:
```bash
# Demo và Development
chmod +x scripts/quick-start.sh
./scripts/quick-start.sh

# Production
chmod +x scripts/startup.sh
./scripts/startup.sh
```

## Scripts Thiết lập Môi trường

### `setup-environment.sh`
**Mục đích**: Cài đặt tất cả các phụ thuộc cần thiết
**Chức năng**:
- Phát hiện hệ điều hành và trình quản lý gói
- Cài đặt công cụ build, OpenSSL, Go, Docker
- Cấu hình môi trường Go (CGO_ENABLED=1)
- Tạo script thử nghiệm

**Cách sử dụng**:
```bash
chmod +x scripts/setup-environment.sh
./scripts/setup-environment.sh
```

### `fix-repositories.sh`
**Mục đích**: Sửa lỗi repository Ubuntu bị hỏng
**Chức năng**:
- Xóa các PPA có vấn đề
- Dọn dẹp bộ nhớ cache gói
- Cập nhật danh sách gói

**Cách sử dụng**:
```bash
chmod +x scripts/fix-repositories.sh
./scripts/fix-repositories.sh
```

### `check-environment.sh`
**Mục đích**: Kiểm tra môi trường toàn diện
**Chức năng**:
- Kiểm tra Go, GCC, OpenSSL, Docker
- Xác minh cài đặt CGO
- Kiểm tra file Fabric
- Chạy thử nghiệm nhanh
- Cung cấp khuyến nghị

**Cách sử dụng**:
```bash
chmod +x scripts/check-environment.sh
./scripts/check-environment.sh
```

### `test_environment.sh`
**Mục đích**: Kiểm tra môi trường đơn giản
**Chức năng**:
- Kiểm tra nhanh các công cụ cơ bản
- Báo cáo phiên bản và trạng thái

**Cách sử dụng**:
```bash
./scripts/test_environment.sh
```

## Scripts Build

### `build-fabric.sh`
**Mục đích**: Build Fabric với tích hợp mã hóa
**Chức năng**:
- Dọn dẹp build trước đó
- Build binary gốc
- Sao chép vào fabric-samples/bin/
- Build Docker image

**Cách sử dụng**:
```bash
export CGO_ENABLED=1
chmod +x scripts/build-fabric.sh
./scripts/build-fabric.sh
```

### `fabric-samples-install.sh`
**Mục đích**: Cài đặt Fabric samples
**Chức năng**:
- Clone repository fabric-samples
- Tải xuống install-fabric.sh
- Cài đặt Docker samples và binary

**Cách sử dụng**:
```bash
chmod +x scripts/fabric-samples-install.sh
./scripts/fabric-samples-install.sh
```

### `download-fabric-samples.sh`
**Mục đích**: Tải xuống repository fabric-samples
**Chức năng**:
- Clone repository fabric-samples
- Tải xuống install-fabric.sh
- Cài đặt Docker samples và binary

**Cách sử dụng**:
```bash
chmod +x scripts/download-fabric-samples.sh
./scripts/download-fabric-samples.sh
```

### `build-encryption.sh`
**Mục đích**: Build thư viện mã hóa AES (libencryption.so)
**Chức năng**:
- Biên dịch encrypt.c thành libencryption.so
- Liên kết với thư viện OpenSSL

**Cách sử dụng**:
```bash
chmod +x scripts/build-encryption.sh
./scripts/build-encryption.sh
```

### `build-mkv-encryption.sh` ⭐
**Mục đích**: Build thư viện MKV encryption (libmkv.so)
**Chức năng**:
- Biên dịch mkv.c và MKV256.c thành libmkv.so
- Tích hợp với KeyManager singleton
- Tự động khởi tạo khóa mã hóa

**Cách sử dụng**:
```bash
chmod +x scripts/build-mkv-encryption.sh
./scripts/build-mkv-encryption.sh
```

## Scripts Mạng

### `start-network.sh`
**Mục đích**: Khởi động mạng thử nghiệm với mã hóa
**Chức năng**:
- Khởi động mạng thử nghiệm
- Tạo channel
- Triển khai chaincode cơ bản
- Kiểm tra chức năng chaincode

**Cách sử dụng**:
```bash
chmod +x scripts/start-network.sh
./scripts/start-network.sh
```

## Scripts Mã hóa Đặc biệt

### `test-encryption.sh`
**Mục đích**: Kiểm tra tích hợp mã hóa AES
**Chức năng**:
- Build thư viện C
- Chạy thử nghiệm Go
- Thực hiện thử nghiệm tích hợp
- Chạy benchmark

**Cách sử dụng**:
```bash
chmod +x scripts/test-encryption.sh
./scripts/test-encryption.sh
```

### `test-mkv-system.sh` ⭐
**Mục đích**: Kiểm tra hệ thống MKV encryption
**Chức năng**:
- Kiểm tra file khóa MKV
- Kiểm tra file thư viện MKV
- Chạy Go tests cho MKV
- Kiểm tra tích hợp LevelDB
- Xác minh KeyManager hoạt động

**Cách sử dụng**:
```bash
chmod +x scripts/test-mkv-system.sh
./scripts/test-mkv-system.sh
```

### `core/ledger/kvledger/txmgmt/statedb/run_tests.sh`
**Mục đích**: Kiểm tra tích hợp mã hóa
**Chức năng**:
- Build thư viện C
- Chạy thử nghiệm Go
- Thực hiện thử nghiệm tích hợp
- Chạy benchmark

**Cách sử dụng**:
```bash
cd core/ledger/kvledger/txmgmt/statedb
chmod +x scripts/run_tests.sh
./scripts/run_tests.sh
```

### `core/ledger/kvledger/txmgmt/statedb/Makefile`
**Mục đích**: Build thư viện mã hóa C
**Chức năng**:
- Biên dịch encrypt.c thành libencryption.so
- Liên kết với thư viện OpenSSL

**Cách sử dụng**:
```bash
cd core/ledger/kvledger/txmgmt/statedb
make clean && make
```

## Scripts Tiện ích

### `list-scripts.sh`
**Mục đích**: Liệt kê tất cả script có sẵn
**Chức năng**:
- Hiển thị danh sách tất cả script
- Mô tả chức năng từng script
- Phân loại theo mục đích sử dụng

**Cách sử dụng**:
```bash
chmod +x scripts/list-scripts.sh
./scripts/list-scripts.sh
```

## Scripts Thử nghiệm & Demo

### `test-quick-start.sh`
**Mục đích**: Kiểm tra tính toàn vẹn và cấu hình script
**Chức năng**:
- Xác minh tất cả script được cấu hình đúng
- Kiểm tra quyền thực thi
- Xác minh cấu trúc thư mục

**Cách sử dụng**:
```bash
chmod +x scripts/test-quick-start.sh
./scripts/test-quick-start.sh
```

### `demo-scripts.sh`
**Mục đích**: Demo tương tác về chức năng script
**Chức năng**:
- Cung cấp demo tương tác về chức năng
- Hiển thị cách sử dụng các script
- Giải thích từng bước thực hiện

**Cách sử dụng**:
```bash
chmod +x scripts/demo-scripts.sh
./scripts/demo-scripts.sh
```

## Sử dụng Cá nhân

### 1. Thiết lập Môi trường
```bash
# Cài đặt tất cả phụ thuộc
./scripts/setup-environment.sh

# Kiểm tra môi trường
./scripts/check-environment.sh

# Thử nghiệm nhanh
./scripts/test_environment.sh
```

### 2. Tải xuống Fabric Samples
```bash
./scripts/download-fabric-samples.sh
```

### 3. Build Thư viện Mã hóa
```bash
./scripts/build-encryption.sh
```

### 4. Kiểm tra Mã hóa
```bash
./scripts/test-encryption.sh
```

### 5. Build Fabric
```bash
./scripts/build-fabric.sh
```

### 6. Khởi động Mạng
```bash
./scripts/start-network.sh
```

### 7. Thử nghiệm & Demo
```bash
# Kiểm tra tính toàn vẹn script
./scripts/test-quick-start.sh

# Chạy demo tương tác
./scripts/demo-scripts.sh

# Liệt kê tất cả script có sẵn
./scripts/list-scripts.sh
```

## Thứ tự Thực thi Script

### Cho Cài đặt Mới
1. **Demo và Development**: `quick-start.sh` (khuyến nghị)
   HOẶC
2. **Production**: `startup.sh`
   HOẶC
3. **Thủ công**: `fix-repositories.sh` → `setup-environment.sh` → `build-mkv-encryption.sh` → `build-fabric.sh` → `start-network.sh`

### Cho Phát triển
1. `check-environment.sh` (xác minh thiết lập)
2. `build-mkv-encryption.sh` (build MKV library)
3. `test-mkv-system.sh` (kiểm tra hệ thống MKV)
4. `core/ledger/kvledger/txmgmt/statedb/run_tests.sh` (kiểm tra mã hóa AES)
5. `build-fabric.sh` (build lại nếu cần)

### Cho Xử lý Sự cố
1. `test_environment.sh` (kiểm tra nhanh)
2. `check-environment.sh` (chẩn đoán chi tiết)
3. `fix-repositories.sh` (nếu có vấn đề repository)

## Tính năng Script

### Quyền Tự động
Tất cả script tự động thiết lập quyền thực thi khi được gọi từ `quick-start.sh`.

### Xử lý Lỗi
Mỗi script bao gồm xử lý lỗi toàn diện và đầu ra có màu để trải nghiệm người dùng tốt hơn.

### Thực thi Độc lập
Mọi script có thể chạy độc lập mà không phụ thuộc vào script khác.

### Kiểm tra Môi trường
Scripts xác minh môi trường trước khi thực thi và cung cấp thông báo lỗi hữu ích.

### Thử nghiệm & Xác minh
- `test-quick-start.sh` xác minh tất cả script được cấu hình đúng
- `demo-scripts.sh` cung cấp demo tương tác về chức năng

## Biến Môi trường

### Bắt buộc
- `CGO_ENABLED=1` - Bật CGO cho thư viện mã hóa

### Tùy chọn
- `GOPATH` - Đường dẫn workspace Go
- `GOROOT` - Đường dẫn cài đặt Go

## Vấn đề Thường gặp và Giải pháp

### Lỗi Repository
```bash
./fix-repositories.sh
```

### CGO Không được Bật
```bash
export CGO_ENABLED=1
```

### Lỗi Build
```bash
make clean
go mod tidy
./build-fabric.sh
```

### Vấn đề Docker
```bash
sudo systemctl start docker
sudo usermod -aG docker $USER
newgrp docker
```

### Quyền Bị Từ chối
```bash
chmod +x *.sh
```

### Vấn đề Phiên bản Go
```bash
./scripts/setup-environment.sh  # Sẽ cài đặt phiên bản Go đúng
```

### Docker Không Chạy
```bash
sudo systemctl start docker
sudo usermod -aG docker $USER  # Sau đó đăng xuất và đăng nhập lại
```

## Phụ thuộc Script

### Yêu cầu Hệ thống
- Ubuntu 20.04+ hoặc tương đương
- Quyền sudo
- Kết nối Internet

### Phụ thuộc Bên ngoài
- Git
- curl/wget
- Docker (được cài đặt bởi script setup)

### Công cụ Bắt buộc
- bash
- curl hoặc wget
- Quyền truy cập sudo
- Kết nối Internet

### Công cụ Tùy chọn
- ldd (để kiểm tra phụ thuộc thư viện)
- pkg-config (để xác minh OpenSSL)
- timeout (cho script demo)

## Ghi log và Giám sát

### Kiểm tra Hoạt động Mã hóa
```bash
# Kiểm tra logs MKV encryption
docker exec peer0.org1.example.com cat /tmp/state_mkv.log

# Kiểm tra logs mã hóa AES
docker logs -f peer0.org1.example.com | grep -i encrypt
docker logs -f peer0.org1.example.com | grep -i decrypt

# Kiểm tra logs MKV
docker logs -f peer0.org1.example.com | grep -i "ENCRYPT\|DECRYPT"
```

### Giám sát Mạng
```bash
cd fabric-samples/test-network
./monitordocker.sh
```

## Cấu trúc File
```
fabric-3.1.1
├── /scripts/quick-start.sh              # Script demo và development
├── /scripts/startup.sh                  # Script production
├── /scripts/setup-environment.sh        # Thiết lập môi trường
├── /scripts/build-fabric.sh            # Build Fabric
├── /scripts/start-network.sh           # Khởi động mạng
├── /scripts/download-fabric-samples.sh # Tải xuống samples
├── /scripts/build-encryption.sh        # Build mã hóa AES
├── /scripts/build-mkv-encryption.sh    # Build MKV encryption ⭐
├── /scripts/test-encryption.sh         # Kiểm tra mã hóa AES
├── /scripts/test-mkv-system.sh         # Kiểm tra hệ thống MKV ⭐
├── /scripts/check-environment.sh       # Kiểm tra môi trường
├── /scripts/fix-repositories.sh        # Sửa repository
├── /scripts/test_environment.sh        # Thử nghiệm nhanh
├── /scripts/list-scripts.sh            # Liệt kê scripts
├── /scripts/test-quick-start.sh        # Kiểm tra tính toàn vẹn
├── /scripts/demo-scripts.sh            # Demo tương tác
└── SCRIPTS.md                          # File này
```

## Ghi chú Bảo mật

- Tất cả script được thiết kế cho mục đích phát triển/demo
- Triển khai sản xuất yêu cầu các biện pháp bảo mật bổ sung
- **MKV Encryption**: Sử dụng KeyManager với password từ file `password.txt`
- **Fallback Password**: "kmasc" nếu không đọc được file password
- **Production**: Sử dụng `startup.sh` thay vì `quick-start.sh`
- Xem xét script trước khi chạy trong môi trường sản xuất

## Đóng góp
Khi thêm script mới:
1. Bao gồm xử lý lỗi phù hợp
2. Thêm đầu ra có màu sử dụng định dạng chuẩn
3. Làm cho script có thể thực thi
4. Cập nhật README này
5. Thử nghiệm độc lập và như một phần của quick-start.sh
6. Thêm vào xác minh test-quick-start.sh

## Hỗ trợ
Cho vấn đề hoặc câu hỏi:
1. Chạy `./scripts/check-environment.sh` để chẩn đoán
2. Chạy `./scripts/test-quick-start.sh` để xác minh tính toàn vẹn script
3. Chạy `./scripts/demo-scripts.sh` để xem scripts hoạt động
4. Kiểm tra log script cá nhân
5. Đảm bảo tất cả phụ thuộc được cài đặt
6. Xác minh bạn đang ở thư mục đúng (fabric-3.1.1/)

---

**Lưu ý**: 
- Tất cả script bao gồm xử lý lỗi và sẽ cung cấp phản hồi rõ ràng về thành công hoặc thất bại
- **Không còn prompt Y/N**: Tất cả script chạy tự động mà không cần xác nhận
- **MKV Encryption**: Hệ thống mã hóa chính với KeyManager tự động
- **Demo vs Production**: Sử dụng `quick-start.sh` cho demo, `startup.sh` cho production
- Kiểm tra đầu ra cho bất kỳ cảnh báo hoặc lỗi nào cần chú ý 