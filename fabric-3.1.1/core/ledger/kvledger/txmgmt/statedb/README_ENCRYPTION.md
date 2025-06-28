# Tích hợp Mã hóa OpenSSL cho Hyperledger Fabric StateDB

## Tổng quan

Module này tích hợp mã hóa/giải mã OpenSSL vào cơ sở dữ liệu trạng thái của Hyperledger Fabric bằng cách sử dụng CGO, thay thế các thư viện crypto của Go để tận dụng hiệu suất cao và các thuật toán đã được chứng minh của OpenSSL.

## Tính năng

- ✅ Tự động mã hóa/giải mã dữ liệu trạng thái
- ✅ Mã hóa OpenSSL AES-256-CBC
- ✅ Tích hợp CGO với thư viện C tùy chỉnh
- ✅ Tương thích với Hyperledger Fabric 3.1.1
- ✅ Minh bạch với chaincode hiện có

## Cấu trúc File

```
statedb/
├── statedb.go          # Go wrapper với tích hợp CGO
├── encrypt.c           # Các hàm mã hóa C
├── encrypt.h           # File header C
├── encrypt_prod.go     # Wrapper mã hóa production
├── Makefile            # Script build cho thư viện C
├── run_tests.sh        # Script chạy test
├── statedb_test.go     # Unit tests
├── libencryption.so    # Thư viện chia sẻ đã biên dịch
└── README_ENCRYPTION.md # File này
```

## Bắt đầu nhanh

### 1. Build Thư viện Mã hóa
```bash
cd core/ledger/kvledger/txmgmt/statedb
make clean && make
```

### 2. Chạy Tests
```bash
./run_tests.sh
```

### 3. Build Fabric với Mã hóa
```bash
cd ../../../../../..
export CGO_ENABLED=1
make clean && make native
```

### 4. Khởi động Test Network
```bash
./start-network.sh
```

### 5. Kiểm tra Logs Mã hóa
```bash
# Giám sát logs peer cho hoạt động mã hóa
docker exec peer0.org1.example.com cat /root/state_encryption.log
```

## Sử dụng

### Tự động Mã hóa/Giải mã
Mã hóa minh bạch với các ứng dụng. Dữ liệu được tự động mã hóa khi lưu trữ và giải mã khi truy xuất:

```go
// Dữ liệu được tự động mã hóa khi lưu trữ
batch.Put("namespace", "key", []byte("sensitive data"), version)

// Dữ liệu được tự động giải mã khi truy xuất
value := batch.Get("namespace", "key")
```

### Mã hóa/Giải mã Thủ công
Bạn cũng có thể mã hóa/giải mã dữ liệu thủ công:

```go
encryptedValue := statedb.EncryptValue([]byte("Hello World"))
decryptedValue := statedb.DecryptValue(encryptedValue)
```

## Kiểm thử

### Chạy Tất cả Tests
```bash
./run_tests.sh
```

### Lệnh Test Riêng lẻ
```bash
# Build thư viện C
make clean && make

# Chạy Go tests
go test ./...

# Chạy với output chi tiết
go test -v ./...

# Chạy tests mã hóa cụ thể
go test -run TestEncryption ./...

# Chạy benchmarks
go test -bench=. ./...
```

### Kết quả Test
Script test sẽ xác minh:
- ✅ Biên dịch thư viện C
- ✅ Liên kết OpenSSL
- ✅ Build package Go
- ✅ Unit tests
- ✅ Integration tests
- ✅ Performance benchmarks

## Thuật toán Mã hóa

- **Thuật toán**: AES-256-CBC
- **Padding**: PKCS7
- **Khóa**: Khóa cố định 32-byte (chỉ demo)
- **IV**: Được tạo ngẫu nhiên cho mỗi lần mã hóa

## Lưu ý Bảo mật

⚠️ **Quan trọng**: Triển khai hiện tại chỉ sử dụng khóa mã hóa cố định cho mục đích trình diễn.

Để sử dụng production:
- Lưu trữ khóa trong HSM hoặc hệ thống quản lý khóa
- Sử dụng khóa động thay vì khóa cố định
- Triển khai xoay khóa đúng cách
- Thêm IV ngẫu nhiên cho mỗi lần mã hóa

## Xử lý sự cố

### Các vấn đề thường gặp

1. **CGO chưa được bật**
   ```bash
   export CGO_ENABLED=1
   ```

2. **Không tìm thấy OpenSSL**
   ```bash
   sudo apt-get install libssl-dev
   pkg-config --modversion openssl
   ```

3. **Build thư viện thất bại**
   ```bash
   make clean && make
   ldd libencryption.so
   ```

4. **Go build thất bại**
   ```bash
   go mod tidy
   go build ./...
   ```

### Kiểm tra Môi trường
Chạy script kiểm tra môi trường từ thư mục gốc dự án:
```bash
./check-environment.sh
```

## Hiệu suất

- OpenSSL được tối ưu hóa cao cho hiệu suất
- Overhead CGO tối thiểu
- Tương thích với tất cả chaincode hiện có
- Tự động mã hóa/giải mã với tác động tối thiểu

## Tích hợp

Mã hóa được tích hợp vào interface cơ sở dữ liệu trạng thái Fabric hiện có:
- `UpdateBatch.Put()` - tự động mã hóa dữ liệu
- `UpdateBatch.Get()` - tự động giải mã dữ liệu
- Các thao tác `VersionedDB` - mã hóa/giải mã minh bạch

## Ghi log

Các thao tác mã hóa/giải mã được ghi log để debug:
- Tìm `[ENCRYPT]` và `[DECRYPT]` trong logs peer
- Logs bao gồm độ dài input/output và trạng thái thao tác
- Các thao tác thất bại sẽ fallback về giá trị gốc

## Phát triển

### Thêm Thuật toán Mã hóa Mới
1. Thêm functions vào `encrypt.c`
2. Cập nhật `encrypt.h`
3. Thêm wrappers trong `statedb.go`
4. Cập nhật logic gọi

### Build cho Phát triển
```bash
# Build thư viện C
make clean && make

# Build Go package
go build ./...

# Chạy tests
go test ./...
```

---

**Lưu ý**: Tích hợp mã hóa này chỉ dành cho mục đích trình diễn và phát triển. Để triển khai production, hãy triển khai quản lý khóa và các biện pháp bảo mật đúng cách. 