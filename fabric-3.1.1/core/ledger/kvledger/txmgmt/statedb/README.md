# Hyperledger Fabric StateDB với MKV256 Encryption

## Tổng quan

Module này tích hợp thuật toán mã hóa MKV256 vào cơ sở dữ liệu trạng thái của Hyperledger Fabric bằng cách sử dụng CGO, cung cấp một giải pháp mã hóa hiệu suất cao và bảo mật cho dữ liệu blockchain.

## Tính năng

- ✅ Tự động mã hóa/giải mã dữ liệu trạng thái
- ✅ Thuật toán MKV256 (256-bit block cipher)
- ✅ Tích hợp CGO với thư viện C tùy chỉnh
- ✅ Tương thích với Hyperledger Fabric 3.1.1
- ✅ Minh bạch với chaincode hiện có
- ✅ Performance benchmarks và tests

## Cấu trúc File

```
statedb/
├── statedb.go                    # Go wrapper với tích hợp CGO
├── encrypt.c                     # Các hàm mã hóa AES (legacy)
├── encrypt.h                     # File header AES
├── encrypt.go                    # Wrapper AES
├── Makefile                      # Script build cho thư viện AES
├── run_tests.sh                  # Script chạy test
├── statedb_test.go               # Unit tests
├── libencryption.so              # Thư viện AES đã biên dịch
├── mkv/                          # MKV256 encryption module
│   ├── mkv.go                    # Go wrapper cho MKV
│   ├── mkv.c                     # C functions cho MKV
│   ├── mkv.h                     # Header file cho MKV
│   ├── MKV256.c                  # MKV256 algorithm implementation
│   ├── MKV256.h                  # MKV256 header
│   ├── PrecomputedTable256.h     # Precomputed tables cho MKV256
│   ├── Makefile                  # Build script cho MKV
│   ├── mkv_test.go               # Unit tests cho MKV
│   └── README.md                 # Chi tiết MKV
├── commontests/                  # Common test utilities
├── mock/                         # Mock implementations
├── PERFORMANCE_COMPARISON.md     # So sánh hiệu suất
├── PERFORMANCE_EVALUATION_LEVELDB.md # Đánh giá hiệu suất LevelDB
└── compare_performance.sh        # Script so sánh performance
```

## Bắt đầu nhanh

### 1. Build Thư viện Mã hóa
```bash
# Build AES library
cd core/ledger/kvledger/txmgmt/statedb
make clean && make

# Build MKV library
cd mkv
make clean && make
cd ..
```

### 2. Chạy Tests
```bash
# Test AES encryption
./run_tests.sh

# Test MKV encryption
cd mkv
LD_LIBRARY_PATH=. go test -v
cd ..
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
docker exec peer0.org1.example.com cat /root/state_mkv.log
```

## Thuật toán Mã hóa

### MKV256
- **Thuật toán**: MKV256 (256-bit block cipher)
- **Block size**: 256 bit (32 bytes)
- **Key size**: 256 bit (32 bytes)
- **Padding**: PKCS#7
- **Mode**: CBC với IV ngẫu nhiên

### AES-256-CBC (Legacy)
- **Thuật toán**: AES-256-CBC
- **Block size**: 128 bit (16 bytes)
- **Key size**: 256 bit (32 bytes)
- **Padding**: PKCS7
- **IV**: Được tạo ngẫu nhiên cho mỗi lần mã hóa

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

#### MKV256
```go
// MKV encryption
encryptedValue := mkv.EncryptValueMKV([]byte("Hello World"), key)
decryptedValue := mkv.DecryptValueMKV(encryptedValue, key)
```

#### AES-256-CBC
```go
// AES encryption
encryptedValue := statedb.EncryptValue([]byte("Hello World"))
decryptedValue := statedb.DecryptValue(encryptedValue)
```

## Kiểm thử

### Chạy Tất cả Tests
```bash
# Test AES
./run_tests.sh

# Test MKV
cd mkv
LD_LIBRARY_PATH=. go test -v
cd ..
```

### Performance Tests
```bash
# So sánh performance
./compare_performance.sh

# Test MKV performance
cd mkv
LD_LIBRARY_PATH=. go test -bench=. -benchmem
cd ..
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

## Performance

### So sánh Hiệu suất
- **MKV256**: Thuật toán mới, tối ưu cho blockchain
- **AES-256-CBC**: Industry standard, hardware acceleration
- **No encryption**: Baseline performance

### Benchmarks
Xem file `PERFORMANCE_COMPARISON.md` và `PERFORMANCE_EVALUATION_LEVELDB.md` để biết chi tiết về hiệu suất.

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

2. **Không tìm thấy OpenSSL (cho AES)**
   ```bash
   sudo apt-get install libssl-dev
   pkg-config --modversion openssl
   ```

3. **Build thư viện thất bại**
   ```bash
   # AES library
   make clean && make
   ldd libencryption.so
   
   # MKV library
   cd mkv
   make clean && make
   ldd libmkv.so
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
1. Thêm functions vào file C tương ứng
2. Cập nhật header file
3. Thêm wrappers trong Go
4. Cập nhật logic gọi

### Build cho Phát triển
```bash
# Build AES library
make clean && make

# Build MKV library
cd mkv
make clean && make
cd ..

# Build Go package
go build ./...

# Chạy tests
go test ./...
```

## Scripts Hỗ trợ

Từ thư mục gốc dự án:
- `scripts/build-encryption.sh` - Build AES library
- `scripts/build-mkv.sh` - Build MKV library
- `scripts/test-mkv.sh` - Test MKV library
- `scripts/build-all-libraries.sh` - Build cả hai libraries
- `scripts/quick-start.sh` - Setup hoàn chỉnh

## Liên hệ

Nếu gặp vấn đề hoặc cần hỗ trợ, hãy liên hệ tác giả hoặc mở issue trong repo.

---

**Lưu ý**: Tích hợp mã hóa này chỉ dành cho mục đích trình diễn và phát triển. Để triển khai production, hãy triển khai quản lý khóa và các biện pháp bảo mật đúng cách.
