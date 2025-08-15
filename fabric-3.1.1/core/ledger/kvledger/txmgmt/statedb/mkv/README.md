# MKV Encryption System

Hệ thống mã hóa MKV với quản lý khóa tự động trong Hyperledger Fabric Statedb.

## Tổng quan

Hệ thống sử dụng **KeyManager singleton** với `sync.Once` để quản lý khóa mã hóa tự động:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Password      │    │   K0 (32 bytes) │    │   K1 (32 bytes) │
│   (from file)   │───▶│   (PBKDF2)      │───▶│   (Random)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │                        │
                              ▼                        ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │ Encrypted K1    │    │ Encrypted Data  │
                       │ (Stored)        │    │ (Stored)        │
                       └─────────────────┘    └─────────────────┘
```

### Khóa K1 (Data Key)
- **Kích thước**: 32 bytes (256 bits)
- **Nguồn gốc**: Sinh ngẫu nhiên khi khởi tạo KeyManager
- **Mục đích**: Mã hóa dữ liệu trong statedb
- **Lưu trữ**: Được mã hóa bằng K0 và lưu trong `encrypted_k1.key`

### Khóa K0 (Derived Key)
- **Kích thước**: 32 bytes (256 bits)
- **Nguồn gốc**: Dẫn xuất từ password bằng PBKDF2-HMAC-SHA256
- **Mục đích**: Mã hóa K1
- **Tham số PBKDF2**:
  - **Salt**: 32 bytes ngẫu nhiên (lưu trong `k0_salt.key`)
  - **Iterations**: 10,000
  - **PRF**: HMAC-SHA256
- **Công thức**: `K0 = PBKDF2(password, salt, 10000, 32)`

## Nội dung thư mục

### Core Files
- `MKV256.c`, `MKV256.h`, `PrecomputedTable256.h`: Thuật toán mã hóa MKV256
- `mkv.c`, `mkv.h`: Hàm mã hóa/giải mã dữ liệu dài, padding PKCS#7
- `mkv.go`: Go wrapper với hệ thống quản lý khóa tự động
- `key_manager.go`: **KeyManager singleton** với `sync.Once`

### Test Files
- `mkv_test.go`: Unit test kiểm tra mã hóa/giải mã
- `key_test.go`: Test cho hệ thống quản lý khóa
- `pbkdf2_test.go`: Test riêng cho implementation PBKDF2
- `key_manager_test.go`: Test cho KeyManager singleton

### Scripts
- `key_manager.sh`: Script quản lý khóa tương tác
- `cleanup_all.sh`: Script dọn dẹp hoàn chỉnh
- `init_with_pbkdf2.sh`: Script khởi tạo với PBKDF2

### Build Files
- `Makefile`: Build shared library `libmkv.so`
- `.gitignore`: Loại trừ file nhạy cảm và tạm thời

## Cài đặt và Build

### 1. Build MKV Library
```bash
make clean && make
```

### 2. Test toàn bộ hệ thống
```bash
go test -v
```

### 3. Test riêng từng phần
```bash
# Test KeyManager
go test -v -run TestKeyManager

# Test mã hóa/giải mã
go test -v -run TestEncryptDecryptValueMKV

# Test PBKDF2
go test -v -run TestPBKDF2
```

## Sử dụng KeyManager

### Khởi tạo tự động
```go
import "github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb/mkv"

// Lần đầu gọi sẽ tự động khởi tạo
keyManager := mkv.GetKeyManager()

// Kiểm tra trạng thái
status := keyManager.GetStatus()
fmt.Printf("Initialized: %v\n", status["initialized"])
```

### Mã hóa/giải mã dữ liệu
```go
// Dữ liệu cần mã hóa
data := []byte("sensitive data")

// Mã hóa (tự động sử dụng K1 từ KeyManager)
encryptedData := mkv.EncryptValueMKV(data)
if encryptedData == nil {
    log.Fatalf("Encryption failed")
}

// Giải mã (tự động sử dụng K1 từ KeyManager)
decryptedData := mkv.DecryptValueMKV(encryptedData)
if decryptedData == nil {
    log.Fatalf("Decryption failed")
}
```

### Thay đổi password
```go
keyManager := mkv.GetKeyManager()
err := keyManager.ChangePassword("new_password")
if err != nil {
    log.Printf("Failed to change password: %v", err)
}
```

### Làm mới keys
```go
keyManager := mkv.GetKeyManager()
err := keyManager.RefreshKeys()
if err != nil {
    log.Printf("Failed to refresh keys: %v", err)
}
```

## Cách hoạt động của KeyManager

### 1. Singleton Pattern với sync.Once
```go
var (
    keyManagerInstance *KeyManager
    keyManagerOnce     sync.Once
)

func GetKeyManager() *KeyManager {
    keyManagerOnce.Do(func() {
        // Chỉ khởi tạo một lần
        keyManagerInstance = &KeyManager{...}
    })
    return keyManagerInstance
}
```

### 2. Tự động khởi tạo
- **Lần đầu gọi**: Tự động tạo hệ thống khóa
- **Lần sau**: Sử dụng instance đã có
- **Thread-safe**: Sử dụng `sync.RWMutex`

### 3. Đọc password từ file
- **Ưu tiên**: `password.txt` trong thư mục hiện tại
- **Fallback**: `/tmp/mkv/password.txt`
- **Default**: "kmasc" nếu không đọc được file

### 4. Quản lý khóa tự động
- **Tạo K1**: Ngẫu nhiên 32 bytes
- **Tạo K0**: Từ password bằng PBKDF2
- **Mã K1**: Bằng K0 và lưu vào `encrypted_k1.key`
- **Lưu trữ**: Trong nhiều thư mục để đảm bảo khả năng truy cập

## Tích hợp với Statedb

### Sử dụng trực tiếp
```go
// Mã hóa value
encryptedValue := mkv.EncryptValueMKV(v.Value)

// Giải mã value
decryptedValue := mkv.DecryptValueMKV(encryptedValue)
```

### Tích hợp với LevelDB
```go
// Trong value_encoding.go
encryptedValue := mkv.EncryptValueMKV(v.Value)
encryptedMetadata := mkv.EncryptValueMKV(v.Metadata)
```

### Tích hợp với CouchDB
```go
// Trong couchdoc_conv.go
encryptedValue := mkv.EncryptValueMKV(value)
```

## PBKDF2 Implementation

### Tổng quan
Hệ thống sử dụng PBKDF2 để tạo khóa K0 từ password, thay vì SHA256 đơn giản.

### Tham số PBKDF2
- **PRF**: HMAC-SHA256
- **Salt**: 32 bytes ngẫu nhiên
- **Iterations**: 10,000
- **Key length**: 32 bytes (256 bits)

### Backward Compatibility
- Nếu không tìm thấy file `k0_salt.key`, hệ thống sẽ sử dụng salt cố định
- Điều này đảm bảo tương thích với dữ liệu cũ

## Bảo mật

### Lưu trữ khóa
- **K1**: Được mã hóa bằng K0 trước khi lưu
- **K0**: Không lưu trữ, chỉ dẫn xuất từ password khi cần
- **Salt**: Lưu trữ trong `k0_salt.key` để tái tạo K0
- **Password**: Đọc từ file `password.txt`

### Quyền truy cập file
- Tất cả file khóa có quyền 600 (chỉ owner đọc/ghi)
- Log file có quyền 644 (owner đọc/ghi, group/other đọc)

### Git Security (.gitignore)
File `.gitignore` đã được cấu hình để **KHÔNG BAO GIỜ** commit các file nhạy cảm:
- `*.key` - Tất cả file khóa
- `*.log` - File log có thể chứa thông tin nhạy cảm
- `*.o`, `*.so` - Build artifacts
- File tạm thời và backup

## Logging

Tất cả thao tác được log vào `/tmp/state_mkv.log`:
```
2024-01-15T10:30:45.123456Z ENCRYPT ns= key= SUCCESS
2024-01-15T10:30:45.234567Z DECRYPT ns= key= SUCCESS
2024-01-15T10:30:45.345678Z KEY_MANAGER_INIT ns= key= SUCCESS
```

## Troubleshooting

### Lỗi "libmkv.so not found"
```bash
# Build lại thư viện
make clean && make
```

### Lỗi "Failed to decrypt K1"
- Kiểm tra file `password.txt` có tồn tại không
- Kiểm tra file `encrypted_k1.key` có tồn tại không
- Chạy lại: `go test -v`

### Lỗi "Permission denied"
```bash
# Sửa quyền file khóa
chmod 600 *.key
chmod 644 *.log
```

### Lỗi "Test failed"
```bash
# Dọn dẹp và build lại
./cleanup_all.sh
make clean && make
go test -v
```

## Hướng dẫn sử dụng từng bước

### Bước 1: Build và Test
```bash
# Build thư viện MKV
make clean && make

# Test toàn bộ hệ thống
go test -v
```

### Bước 2: Tạo file password
```bash
# Tạo file password.txt với nội dung "kmasc"
echo "kmasc" > password.txt
```

### Bước 3: Sử dụng trong ứng dụng
```go
import "github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb/mkv"

// Lấy KeyManager (tự động khởi tạo)
keyManager := mkv.GetKeyManager()

// Mã hóa dữ liệu
data := []byte("Hello, World!")
encrypted := mkv.EncryptValueMKV(data)

// Giải mã dữ liệu
decrypted := mkv.DecryptValueMKV(encrypted)
```

### Bước 4: Quản lý khóa (tùy chọn)
```go
// Thay đổi password
err := keyManager.ChangePassword("new_password")

// Làm mới keys
err := keyManager.RefreshKeys()

// Kiểm tra trạng thái
status := keyManager.GetStatus()
```

## Ví dụ hoàn chỉnh

```bash
# 1. Build thư viện
make clean && make

# 2. Tạo file password
echo "kmasc" > password.txt

# 3. Test hệ thống
go test -v

# 4. Sử dụng trong code Go
# (xem ví dụ code ở trên)
```

## Lưu ý

- **Password**: Đặt trong file `password.txt` (mặc định: "kmasc")
- **Backup**: Luôn backup file `encrypted_k1.key` trước khi thay đổi password
- **Environment**: Đảm bảo môi trường an toàn khi tạo file password
- **Rotation**: Nên thay đổi password định kỳ
- **Recovery**: Có kế hoạch khôi phục khóa trong trường hợp mất password

## Hoàn thành!

Hệ thống MKV với KeyManager singleton đã hoạt động hoàn hảo! Đây là tóm tắt những gì đã được thực hiện:

### ✅ **Hệ thống KeyManager hoạt động hoàn hảo:**

1. **Singleton Pattern**: Sử dụng `sync.Once` để đảm bảo chỉ khởi tạo một lần
2. **Tự động khởi tạo**: Tự động tạo hệ thống khóa khi lần đầu được gọi
3. **Quản lý khóa tự động**: Tự động tạo, lưu trữ và quản lý các khóa mã hóa
4. **Thread-safe**: Sử dụng `sync.RWMutex` để đảm bảo an toàn khi truy cập đồng thời
5. **Fallback mechanism**: Có cơ chế dự phòng khi khởi tạo thất bại
6. **Logging**: Ghi log tất cả các hoạt động vào file `/tmp/state_mkv.log`

### 🚀 **Cách sử dụng đơn giản:**

```go
// Lấy KeyManager (tự động khởi tạo)
keyManager := mkv.GetKeyManager()

// Mã hóa/giải mã (tự động sử dụng K1)
encrypted := mkv.EncryptValueMKV(data)
decrypted := mkv.DecryptValueMKV(encrypted)
```

### 📁 **Files chính:**

- `key_manager.go` - KeyManager singleton với `sync.Once`
- `mkv.go` - Go wrapper với hệ thống quản lý khóa tự động
- `key_manager_test.go` - Tests cho KeyManager
- `mkv_test.go` - Tests cho mã hóa/giải mã
- `README.md` - Documentation chi tiết