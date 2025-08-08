# MKV Key Management System

Hệ thống quản lý khóa 2 tầng cho MKV encryption trong Hyperledger Fabric Statedb.

## Tổng quan

Hệ thống sử dụng **2 tầng khóa** theo sơ đồ:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Password      │    │   K0 (32 bytes) │    │   K1 (32 bytes) │
│   (User Input)  │───▶│   (Derived Key) │───▶│   (Data Key)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │                        │
                              ▼                        ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │ Encrypted K1    │    │ Encrypted Data  │
                       │ (Stored in DB)  │    │ (Stored in DB)  │
                       └─────────────────┘    └─────────────────┘
```

### Khóa K1 (Data Key)
- **Kích thước**: 32 bytes (256 bits)
- **Nguồn gốc**: Sinh ngẫu nhiên khi tạo chain
- **Mục đích**: Mã hóa dữ liệu trong statedb
- **Lưu trữ**: Được mã hóa bằng K0 và lưu trong database

### Khóa K0 (Derived Key)
- **Kích thước**: 32 bytes (256 bits)
- **Nguồn gốc**: Dẫn xuất từ password bằng PBKDF2-HMAC-SHA256
- **Mục đích**: Mã hóa K1
- **Tham số PBKDF2**:
  - **Salt**: 32 bytes ngẫu nhiên (lưu trong `k0_salt.key`)
  - **Iterations**: 10,000 (theo khuyến nghị OWASP)
  - **PRF**: HMAC-SHA256
- **Công thức**: `K0 = PBKDF2(password, salt, 10000, 32)`

## Nội dung thư mục
- `MKV256.c`, `MKV256.h`, `PrecomputedTable256.h`: Thuật toán mã hóa MKV256
- `mkv.c`, `mkv.h`: Hàm mã hóa/giải mã dữ liệu dài, padding PKCS#7
- `mkv.go`: Go wrapper với hệ thống quản lý khóa 2 tầng và PBKDF2
- `mkv_test.go`: Unit test kiểm tra mã hóa/giải mã
- `key_test.go`: Test cho hệ thống quản lý khóa
- `pbkdf2_test.go`: Test riêng cho implementation PBKDF2
- `key_manager.sh`: Script quản lý khóa tương tác
- `cleanup_all.sh`: Script dọn dẹp hoàn chỉnh
- `demo.sh`: Script demo hệ thống
- `Makefile`: Build shared library `libmkv.so`
- `.gitignore`: Loại trừ file nhạy cảm và tạm thời
- `README.md`: File này

## Cài đặt và Build

### 1. Build MKV Library
```bash
make clean && make
```

### 2. Test MKV Functions
```bash
LD_LIBRARY_PATH=. go test -v
```

### 3. Test Key Management
```bash
LD_LIBRARY_PATH=. go test -v -run TestKeyManagementSystem
LD_LIBRARY_PATH=. go test -v -run TestDataEncryptionWithK1
```

## Sử dụng Key Manager Script

### Khởi tạo hệ thống
```bash
./key_manager.sh init
```
Script sẽ yêu cầu nhập password và tạo:
- `k1.key`: K1 ngẫu nhiên (plaintext)
- `k0.key`: K0 dẫn xuất từ password
- `encrypted_k1.key`: K1 đã mã bằng K0

### Thay đổi password
```bash
./key_manager.sh change
```
Script sẽ yêu cầu:
- Password cũ
- Password mới

Sau đó giải mã K1 bằng password cũ và mã lại bằng password mới.

### Kiểm tra trạng thái
```bash
./key_manager.sh status
```
Hiển thị trạng thái các file khóa hiện tại.

### Dọn dẹp file tạm
```bash
./key_manager.sh cleanup
```
Xóa các file tạm thời (k0.key, decrypted_k1.key).

## Sử dụng trong Code Go

### Khởi tạo hệ thống
```go
import "github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb/mkv"

// Khởi tạo với password
err := mkv.InitializeKeyManagement("mysecretpassword")
if err != nil {
    log.Fatalf("Failed to initialize: %v", err)
}
```

### Lấy K1 hiện tại
```go
// Lấy K1 bằng password
k1, err := mkv.GetCurrentK1("mysecretpassword")
if err != nil {
    log.Fatalf("Failed to get K1: %v", err)
}
```

### Mã hóa dữ liệu với K1
```go
// Dữ liệu cần mã hóa
data := []byte("sensitive data")

// Mã hóa với K1
encryptedData := mkv.EncryptValueMKV(data, k1)
if encryptedData == nil {
    log.Fatalf("Encryption failed")
}

// Giải mã với K1
decryptedData := mkv.DecryptValueMKV(encryptedData, k1)
if decryptedData == nil {
    log.Fatalf("Decryption failed")
}
```

### Thay đổi password
```go
err := mkv.ChangePassword("oldpassword", "newpassword")
if err != nil {
    log.Fatalf("Failed to change password: %v", err)
}
```

## Tích hợp với Statedb

### LevelDB
Trong `value_encoding.go`:
```go
// Lấy K1 từ password
k1, err := mkv.GetCurrentK1("statedb_password")
if err != nil {
    return nil, err
}

// Mã hóa value và metadata
encryptedValue := mkv.EncryptValueMKV(v.Value, k1)
encryptedMetadata := mkv.EncryptValueMKV(v.Metadata, k1)
```

### Private Data Storage
Trong `store.go`:
```go
// Lấy K1 cho private data
k1, err := mkv.GetCurrentK1("private_data_password")
if err != nil {
    return err
}

// Mã hóa private data
encryptedData := mkv.EncryptValueMKV(data, k1)
```

### CouchDB
Trong `couchdoc_conv.go`:
```go
// Lấy K1 cho CouchDB
k1, err := mkv.GetCurrentK1("couchdb_password")
if err != nil {
    return nil, err
}

// Mã hóa document
encryptedValue := mkv.EncryptValueMKV(value, k1)
```

## PBKDF2 Implementation

### Tổng quan
Hệ thống sử dụng PBKDF2 (Password-Based Key Derivation Function 2) để tạo khóa K0 từ password, thay vì SHA256 đơn giản như trước đây.

### Tham số PBKDF2
- **PRF**: HMAC-SHA256
- **Salt**: 32 bytes ngẫu nhiên
- **Iterations**: 10,000 (theo khuyến nghị OWASP)
- **Key length**: 32 bytes (256 bits)

### Ưu điểm so với SHA256 đơn giản
1. **Salt**: Chống rainbow table attacks
2. **Iterations**: Làm chậm brute force attacks
3. **HMAC**: Tăng tính bảo mật so với hash đơn giản
4. **Tiêu chuẩn**: PBKDF2 là tiêu chuẩn được chấp nhận rộng rãi

### Backward Compatibility
- Nếu không tìm thấy file `k0_salt.key`, hệ thống sẽ sử dụng salt cố định
- Điều này đảm bảo tương thích với dữ liệu cũ

## Bảo mật

### Lưu trữ khóa
- **K1**: Được mã hóa bằng K0 trước khi lưu
- **K0**: Không lưu trữ, chỉ dẫn xuất từ password khi cần
- **Salt**: Lưu trữ trong `k0_salt.key` để tái tạo K0
- **Password**: Không lưu trữ, chỉ dùng để sinh K0 qua PBKDF2

### Quyền truy cập file
- Tất cả file khóa có quyền 600 (chỉ owner đọc/ghi)
- Log file có quyền 644 (owner đọc/ghi, group/other đọc)

### Xóa file tạm
- `k0.key`: Xóa sau khi sử dụng
- `decrypted_k1.key`: Xóa sau khi mã lại
- `k1.key`: Có thể xóa sau khi đã mã bằng K0

### Git Security (.gitignore)
File `.gitignore` đã được cấu hình để **KHÔNG BAO GIỜ** commit các file nhạy cảm:
- `*.key` - Tất cả file khóa
- `*.log` - File log có thể chứa thông tin nhạy cảm
- `*.o`, `*.so` - Build artifacts
- File tạm thời và backup

**⚠️ QUAN TRỌNG**: Luôn kiểm tra `git status` trước khi commit để đảm bảo không có file nhạy cảm nào bị đưa vào repository.

## Logging

Tất cả thao tác được log vào `/tmp/state_mkv.log`:
```
2024-01-15T10:30:45.123456Z ENCRYPT ns= ns= SUCCESS
2024-01-15T10:30:45.234567Z DECRYPT ns= ns= SUCCESS
2024-01-15T10:30:45.345678Z GENERATE_K1 ns= ns= SUCCESS
2024-01-15T10:30:45.456789Z INIT_KEYS ns= ns= SUCCESS
```

## Troubleshooting

### Lỗi "libmkv.so not found"
```bash
# Build lại thư viện
make clean && make
```

### Lỗi "Failed to decrypt K1"
- Kiểm tra password có đúng không
- Kiểm tra file `encrypted_k1.key` có tồn tại không
- Chạy lại: `./key_manager.sh init`

### Lỗi "Permission denied"
```bash
# Sửa quyền file khóa
chmod 600 *.key
chmod 644 *.log
```

### Lỗi "No K1 found"
```bash
# Khởi tạo lại hệ thống
./key_manager.sh init
```

### Lỗi "Test failed"
```bash
# Dọn dẹp và build lại
./cleanup_all.sh
make clean && make
LD_LIBRARY_PATH=. go test -v
```

### Lỗi "Demo script failed"
```bash
# Kiểm tra trạng thái hệ thống
./key_manager.sh status

# Nếu chưa khởi tạo, chạy:
./key_manager.sh init
```

## Hướng dẫn sử dụng từng bước

### Bước 1: Build và Test
```bash
# Build thư viện MKV
make clean && make

# Test các chức năng cơ bản
LD_LIBRARY_PATH=. go test -v

# Test hệ thống quản lý khóa
LD_LIBRARY_PATH=. go test -v -run TestKeyManagementSystem
LD_LIBRARY_PATH=. go test -v -run TestDataEncryptionWithK1
```

### Bước 2: Khởi tạo hệ thống
```bash
# Khởi tạo hệ thống quản lý khóa
./key_manager.sh init
# Nhập password khi được yêu cầu (ví dụ: mysecret123)
```

### Bước 3: Kiểm tra trạng thái
```bash
# Kiểm tra trạng thái các file khóa
./key_manager.sh status
```

### Bước 4: Sử dụng trong ứng dụng
```bash
# Chạy demo hoàn chỉnh
./demo.sh

# Hoặc test thủ công
LD_LIBRARY_PATH=. go test -v
```

### Bước 5: Quản lý khóa (tùy chọn)
```bash
# Thay đổi password
./key_manager.sh change
# Nhập password cũ và password mới

# Dọn dẹp file tạm
./key_manager.sh cleanup

# Dọn dẹp hoàn toàn (xóa tất cả file khóa)
./cleanup_all.sh
```

## Ví dụ hoàn chỉnh

```bash
# 1. Build thư viện
make clean && make

# 2. Test các chức năng
LD_LIBRARY_PATH=. go test -v

# 3. Khởi tạo hệ thống
./key_manager.sh init
# Nhập password: mysecret123

# 4. Kiểm tra trạng thái
./key_manager.sh status

# 5. Chạy demo
./demo.sh

# 6. Thay đổi password (tùy chọn)
./key_manager.sh change
# Nhập old password: mysecret123
# Nhập new password: newsecret456

# 7. Dọn dẹp
./key_manager.sh cleanup
```

## Lưu ý

- **Backup**: Luôn backup file `encrypted_k1.key` trước khi thay đổi password
- **Password**: Sử dụng password mạnh (ít nhất 12 ký tự, có số, chữ hoa, chữ thường, ký tự đặc biệt)
- **Environment**: Đảm bảo môi trường an toàn khi nhập password
- **Rotation**: Nên thay đổi password định kỳ
- **Recovery**: Có kế hoạch khôi phục khóa trong trường hợp mất password

##  **Hoàn thành!**

Hệ thống quản lý khóa MKV đã hoạt động hoàn hảo! Đây là tóm tắt những gì đã được thực hiện:

### ✅ **Đã sửa lỗi thành công:**

1. **Lỗi thay đổi password**: Đã sửa logic trong `key_manager.sh` để không ghi đè K0 trước khi giải mã K1
2. **Lỗi file tạm**: Đã tạo script `cleanup_all.sh` để dọn dẹp hoàn chỉnh
3. **Lỗi xung đột main function**: Đã dọn dẹp file C tạm thời trước khi test Go

### ✅ **Hệ thống hoạt động hoàn hảo:**

1. **Khởi tạo hệ thống**: ✅
   - Tạo K1 ngẫu nhiên 32 bytes
   - Tạo K0 từ password bằng SHA256
   - Mã K1 bằng K0 và lưu vào file

2. **Thay đổi password**: ✅
   - Giải mã K1 bằng password cũ
   - Mã lại K1 bằng password mới
   - K1 không thay đổi, chỉ thay đổi cách mã hóa

3. **Go integration**: ✅
   - Tất cả tests pass
   - Data encryption/decryption hoạt động
   - File operations hoạt động

### 🚀 **Thứ tự sử dụng chính xác:**

```bash
# 1. Build và test
make clean && make
LD_LIBRARY_PATH=. go test -v

# 2. Khởi tạo hệ thống
./key_manager.sh init

# 3. Kiểm tra trạng thái
./key_manager.sh status

# 4. Sử dụng (demo hoặc tích hợp)
./demo.sh

# 5. Quản lý khóa (tùy chọn)
./key_manager.sh change
./key_manager.sh cleanup
```

### 📁 **Files đã tạo:**

- `cleanup_all.sh` - Script dọn dẹp hoàn chỉnh
- `key_manager.sh` - Script quản lý khóa (đã sửa lỗi)
- `mkv.go` - Go wrapper với hệ thống quản lý khóa 2 tầng
- `key_test.go` - Tests hoàn chỉnh
- `.gitignore` - Bảo vệ file nhạy cảm khỏi Git
- `README.md` - Documentation chi tiết
- `IMPLEMENTATION_SUMMARY.md` - Tóm tắt implementation