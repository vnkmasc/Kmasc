# Hyperledger Fabric StateDB với MKV256 Encryption

## Tổng quan

Module này tích hợp thuật toán mã hóa MKV256 vào cơ sở dữ liệu trạng thái của Hyperledger Fabric bằng cách sử dụng CGO, cung cấp một giải pháp mã hóa hiệu suất cao và bảo mật cho dữ liệu blockchain với hệ thống quản lý khóa tự động.

## Tính năng

- ✅ Tự động mã hóa/giải mã dữ liệu trạng thái
- ✅ Thuật toán MKV256 (256-bit block cipher)
- ✅ **KeyManager singleton** với `sync.Once` để quản lý khóa tự động
- ✅ Tự động đọc password từ file `password.txt`
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
├── mkv/                          # MKV256 encryption module với KeyManager
│   ├── mkv.go                    # Go wrapper với hệ thống quản lý khóa tự động
│   ├── key_manager.go            # **KeyManager singleton** với sync.Once
│   ├── mkv.c                     # C functions cho MKV
│   ├── mkv.h                     # Header file cho MKV
│   ├── MKV256.c                  # MKV256 algorithm implementation
│   ├── MKV256.h                  # MKV256 header
│   ├── PrecomputedTable256.h     # Precomputed tables cho MKV256
│   ├── Makefile                  # Build script cho MKV
│   ├── mkv_test.go               # Unit tests cho MKV
│   ├── key_manager_test.go       # Tests cho KeyManager singleton
│   ├── key_test.go               # Tests cho hệ thống quản lý khóa
│   ├── pbkdf2_test.go            # Tests cho PBKDF2
│   ├── README.md                 # Chi tiết MKV
│   └── KEY_MANAGER_README.md     # Chi tiết KeyManager
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

### 2. Tạo file password (cho MKV)
```bash
# Tạo file password.txt với nội dung "kmasc"
echo "kmasc" > mkv/password.txt
```

### 3. Chạy Tests
```bash
# Test AES encryption
./run_tests.sh

# Test MKV encryption với KeyManager
cd mkv
go test -v
cd ..
```

### 4. Build Fabric với Mã hóa
```bash
cd ../../../../../..
export CGO_ENABLED=1
make clean && make native
```

### 5. Khởi động Test Network
```bash
./start-network.sh
```

### 6. Kiểm tra Logs Mã hóa
```bash
# Giám sát logs peer cho hoạt động mã hóa
docker exec peer0.org1.example.com cat /root/state_encryption.log
docker exec peer0.org1.example.com cat /tmp/state_mkv.log
```

## Thuật toán Mã hóa

### MKV256 với KeyManager
- **Thuật toán**: MKV256 (256-bit block cipher)
- **Block size**: 256 bit (32 bytes)
- **Key size**: 256 bit (32 bytes)
- **Padding**: PKCS#7
- **Mode**: CBC với IV ngẫu nhiên
- **Quản lý khóa**: **KeyManager singleton** với `sync.Once`
- **Password**: Đọc từ file `password.txt` (fallback: "kmasc")
- **Khóa K1**: Sinh ngẫu nhiên 32 bytes
- **Khóa K0**: Dẫn xuất từ password bằng PBKDF2-HMAC-SHA256

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

#### MKV256 với KeyManager
```go
import "github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb/mkv"

// Mã hóa (tự động sử dụng K1 từ KeyManager)
encryptedValue := mkv.EncryptValueMKV([]byte("Hello World"))

// Giải mã (tự động sử dụng K1 từ KeyManager)
decryptedValue := mkv.DecryptValueMKV(encryptedValue)

// Lấy KeyManager instance
keyManager := mkv.GetKeyManager()

// Kiểm tra trạng thái
status := keyManager.GetStatus()
fmt.Printf("Initialized: %v\n", status["initialized"])

// Thay đổi password
err := keyManager.ChangePassword("new_password")

// Làm mới keys
err := keyManager.RefreshKeys()
```

#### AES-256-CBC
```go
// AES encryption
encryptedValue := statedb.EncryptValue([]byte("Hello World"))
decryptedValue := statedb.DecryptValue(encryptedValue)
```

## Hệ thống KeyManager

### Tính năng chính
- **Singleton Pattern**: Sử dụng `sync.Once` để đảm bảo chỉ khởi tạo một lần
- **Tự động khởi tạo**: Tự động tạo hệ thống khóa khi lần đầu được gọi
- **Quản lý khóa tự động**: Tự động tạo, lưu trữ và quản lý các khóa mã hóa
- **Thread-safe**: Sử dụng `sync.RWMutex` để đảm bảo an toàn khi truy cập đồng thời
- **Fallback mechanism**: Có cơ chế dự phòng khi khởi tạo thất bại
- **Logging**: Ghi log tất cả các hoạt động vào file `/tmp/state_mkv.log`

### Cách hoạt động
1. **Lần đầu gọi**: `mkv.GetKeyManager()` tự động khởi tạo hệ thống khóa
2. **Tạo K1**: Khóa mã hóa ngẫu nhiên 32 bytes
3. **Tạo K0**: Từ password (đọc từ file) bằng PBKDF2
4. **Mã K1**: Bằng K0 và lưu vào `encrypted_k1.key`
5. **Lần sau**: Sử dụng instance đã có, không khởi tạo lại

### Quản lý password
- **File ưu tiên**: `password.txt` trong thư mục hiện tại
- **File fallback**: `/tmp/mkv/password.txt`
- **Password mặc định**: "kmasc" nếu không đọc được file

## Kiểm thử

### Chạy Tất cả Tests
```bash
# Test AES
./run_tests.sh

# Test MKV với KeyManager
cd mkv
go test -v
cd ..
```

### Performance Tests
```bash
# So sánh performance
./compare_performance.sh

# Test MKV performance
cd mkv
go test -bench=. -benchmem
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

### Test KeyManager
```bash
cd mkv

# Test singleton pattern
go test -v -run TestKeyManagerSingleton

# Test mã hóa/giải mã
go test -v -run TestKeyManagerEncryptionDecryption

# Test thay đổi password
go test -v -run TestKeyManagerChangePassword

# Test làm mới keys
go test -v -run TestKeyManagerRefreshKeys

# Test truy cập đồng thời
go test -v -run TestKeyManagerConcurrentAccess

cd ..
```

## Performance

### So sánh Hiệu suất
- **MKV256 với KeyManager**: Thuật toán mới, quản lý khóa tự động, tối ưu cho blockchain
- **AES-256-CBC**: Industry standard, hardware acceleration
- **No encryption**: Baseline performance

### Benchmarks
Xem file `PERFORMANCE_COMPARISON.md` và `PERFORMANCE_EVALUATION_LEVELDB.md` để biết chi tiết về hiệu suất.

## Lưu ý Bảo mật

⚠️ **Quan trọng**: Triển khai hiện tại sử dụng hệ thống quản lý khóa tự động với KeyManager.

Để sử dụng production:
- Tạo file `password.txt` với password mạnh
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
   cd ..
   ```

4. **Go build thất bại**
   ```bash
   go mod tidy
   go build ./...
   ```

5. **KeyManager không khởi tạo được**
   ```bash
   # Kiểm tra file password.txt
   ls -la mkv/password.txt
   
   # Tạo file password nếu chưa có
   echo "kmasc" > mkv/password.txt
   
   # Chạy test để kiểm tra
   cd mkv
   go test -v
   cd ..
   ```

6. **Tests MKV thất bại**
   ```bash
   cd mkv
   
   # Dọn dẹp files cũ
   rm -f *.key *.log
   
   # Tạo password file
   echo "kmasc" > password.txt
   
   # Build lại library
   make clean && make
   
   # Chạy test
   go test -v
   
   cd ..
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

### Tích hợp với KeyManager
```go
// Trong value_encoding.go
encryptedValue := mkv.EncryptValueMKV(v.Value)
encryptedMetadata := mkv.EncryptValueMKV(v.Metadata)

// Trong couchdoc_conv.go
encryptedValue := mkv.EncryptValueMKV(value)
```

## Ghi log

Các thao tác mã hóa/giải mã được ghi log để debug:

### MKV với KeyManager
- Log file: `/tmp/state_mkv.log`
- Format: `TIMESTAMP OPERATION ns=NAMESPACE key=KEY STATUS [ERROR: ERROR_MESSAGE]`
- Ví dụ:
  ```
  2024-01-15T10:30:45.123456Z KEY_MANAGER_INIT ns= key= SUCCESS
  2024-01-15T10:30:45.124567Z ENCRYPT ns= key= SUCCESS
  2024-01-15T10:30:45.125678Z DECRYPT ns= key= SUCCESS
  ```

### AES (Legacy)
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

### Phát triển KeyManager
```bash
cd mkv

# Chạy tests KeyManager
go test -v -run TestKeyManager

# Chạy tests mã hóa/giải mã
go test -v -run TestEncryptDecryptValueMKV

# Chạy tests PBKDF2
go test -v -run TestPBKDF2

cd ..
```

## Scripts Hỗ trợ

Từ thư mục gốc dự án:
- `scripts/build-encryption.sh` - Build AES library
- `scripts/build-mkv.sh` - Build MKV library
- `scripts/test-mkv.sh` - Test MKV library
- `scripts/build-all-libraries.sh` - Build cả hai libraries
- `scripts/quick-start.sh` - Setup hoàn chỉnh

## Ví dụ sử dụng hoàn chỉnh

### 1. Setup ban đầu
```bash
# Build libraries
cd core/ledger/kvledger/txmgmt/statedb
make clean && make

cd mkv
make clean && make

# Tạo password file
echo "kmasc" > password.txt

# Test hệ thống
go test -v

cd ..
```

### 2. Sử dụng trong code
```go
import "github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb/mkv"

// Lấy KeyManager (tự động khởi tạo)
keyManager := mkv.GetKeyManager()

// Mã hóa dữ liệu
data := []byte("Sensitive blockchain data")
encrypted := mkv.EncryptValueMKV(data)

// Giải mã dữ liệu
decrypted := mkv.DecryptValueMKV(encrypted)

// Kiểm tra trạng thái
status := keyManager.GetStatus()
fmt.Printf("System ready: %v\n", status["initialized"])
```

### 3. Quản lý khóa
```go
// Thay đổi password
err := keyManager.ChangePassword("new_strong_password")

// Làm mới keys
err := keyManager.RefreshKeys()

// Kiểm tra trạng thái
status := keyManager.GetStatus()
```

## Liên hệ

Nếu gặp vấn đề hoặc cần hỗ trợ, hãy liên hệ tác giả hoặc mở issue trong repo.

---

**Lưu ý**: Tích hợp mã hóa này với KeyManager tự động chỉ dành cho mục đích trình diễn và phát triển. Để triển khai production, hãy triển khai quản lý khóa và các biện pháp bảo mật đúng cách.
