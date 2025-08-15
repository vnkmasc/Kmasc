# MKV Key Manager

## Tổng quan

KeyManager là một hệ thống quản lý khóa mã hóa sử dụng pattern Singleton với `sync.Once` để đảm bảo chỉ khởi tạo một lần và sử dụng lại instance đó trong suốt vòng đời của ứng dụng.

## Tính năng chính

- **Singleton Pattern**: Sử dụng `sync.Once` để đảm bảo chỉ tạo một instance
- **Tự động khởi tạo**: Tự động tạo hệ thống khóa khi lần đầu được gọi
- **Quản lý khóa tự động**: Tự động tạo, lưu trữ và quản lý các khóa mã hóa
- **Thread-safe**: Sử dụng `sync.RWMutex` để đảm bảo an toàn khi truy cập đồng thời
- **Fallback mechanism**: Có cơ chế dự phòng khi khởi tạo thất bại
- **Logging**: Ghi log tất cả các hoạt động vào file `/tmp/state_mkv.log`

## Cách sử dụng

### 1. Lấy instance KeyManager

```go
import "github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb/mkv"

// Lấy instance duy nhất của KeyManager
keyManager := mkv.GetKeyManager()
```

### 2. Sử dụng trực tiếp cho mã hóa/giải mã

```go
// Mã hóa dữ liệu (tự động sử dụng K1 từ KeyManager)
data := []byte("Hello, World!")
encrypted := mkv.EncryptValueMKV(data)

// Giải mã dữ liệu (tự động sử dụng K1 từ KeyManager)
decrypted := mkv.DecryptValueMKV(encrypted)
```

### 3. Truy cập khóa mã hóa

```go
keyManager := mkv.GetKeyManager()
encryptionKey := keyManager.GetEncryptionKey()
```

### 4. Kiểm tra trạng thái

```go
status := keyManager.GetStatus()
fmt.Printf("Initialized: %v\n", status["initialized"])
fmt.Printf("Has Key: %v\n", status["has_key"])
fmt.Printf("Key Length: %d\n", status["key_length"])
fmt.Printf("Password Set: %v\n", status["password_set"])
fmt.Printf("Timestamp: %s\n", status["timestamp"])
```

### 5. Thay đổi password

```go
err := keyManager.ChangePassword("new_password")
if err != nil {
    log.Printf("Failed to change password: %v", err)
}
```

### 6. Làm mới keys

```go
err := keyManager.RefreshKeys()
if err != nil {
    log.Printf("Failed to refresh keys: %v", err)
}
```

## Cấu trúc file

Khi khởi tạo, KeyManager sẽ tạo các file sau trong các thư mục tìm kiếm:

- `k1.key`: Khóa mã hóa chính (K1) dạng plaintext
- `k0.key`: Khóa dẫn xuất từ password (K0)
- `encrypted_k1.key`: K1 đã được mã hóa bằng K0
- `k0_salt.key`: Salt cho PBKDF2
- `password.txt`: Password (để tham khảo)

## Các thư mục tìm kiếm

KeyManager sẽ tìm kiếm và lưu trữ keys trong các thư mục sau:

1. `.` (thư mục hiện tại)
2. `/tmp`
3. `/tmp/mkv`
4. `/opt/mkv`
5. `/home/chaincode/mkv`
6. `/root/mkv`

## Cơ chế hoạt động

### 1. Khởi tạo lần đầu

```go
keyManager := mkv.GetKeyManager() // Tự động khởi tạo
```

- Tạo K1 ngẫu nhiên 32 bytes
- Tạo salt ngẫu nhiên 32 bytes
- Tạo K0 từ password (đọc từ file `password.txt`) với PBKDF2
- Mã hóa K1 bằng K0
- Lưu tất cả keys vào các thư mục

### 2. Sử dụng sau khởi tạo

```go
keyManager := mkv.GetKeyManager() // Trả về instance đã có
```

- Load K1 đã mã từ file
- Giải mã K1 bằng password
- Sử dụng K1 để mã hóa/giải mã dữ liệu

### 3. Fallback mechanism

Nếu khởi tạo thất bại, hệ thống sẽ sử dụng fallback key:
```go
fallbackKey := []byte("1234567890abcdef1234567890abcdef")
```

## Thread Safety

KeyManager sử dụng `sync.RWMutex` để đảm bảo an toàn khi truy cập đồng thời:

- **Read operations** (GetEncryptionKey, GetStatus): Sử dụng `RLock()`
- **Write operations** (ChangePassword, RefreshKeys): Sử dụng `Lock()`

## Logging

Tất cả các hoạt động được ghi vào file `/tmp/state_mkv.log` với format:

```
TIMESTAMP OPERATION ns=NAMESPACE key=KEY STATUS [ERROR: ERROR_MESSAGE]
```

Ví dụ:
```
2024-01-15T10:30:45.123456Z KEY_MANAGER_INIT ns= key= SUCCESS
2024-01-15T10:30:45.124567Z ENCRYPT ns= key= SUCCESS
2024-01-15T10:30:45.125678Z DECRYPT ns= key= SUCCESS
```

## Testing

Chạy tests:

```bash
cd fabric-3.1.1/core/ledger/kvledger/txmgmt/statedb/mkv
go test -v
```

### Các test functions có sẵn:

1. **TestKeyManagerSingleton** - Test singleton pattern
2. **TestKeyManagerEncryptionDecryption** - Test mã hóa/giải mã
3. **TestKeyManagerChangePassword** - Test thay đổi password
4. **TestKeyManagerRefreshKeys** - Test làm mới keys
5. **TestKeyManagerConcurrentAccess** - Test truy cập đồng thời

## Lưu ý bảo mật

1. **Password mặc định**: Hệ thống sử dụng password mặc định "kmasc" nếu không đọc được file `password.txt`. Trong môi trường production, nên tạo file `password.txt` với password mạnh.

2. **File permissions**: Các file key được tạo với permission 0600 (chỉ owner có thể đọc/ghi).

3. **Key storage**: Keys được lưu trữ ở nhiều vị trí để đảm bảo khả năng truy cập, nhưng điều này có thể tạo ra rủi ro bảo mật.

4. **Fallback key**: Fallback key được hardcode trong code, điều này không an toàn trong môi trường production.

## Troubleshooting

### 1. Không thể tạo thư mục

```
ERROR: Failed to create directory: permission denied
```

**Giải pháp**: Kiểm tra quyền ghi vào các thư mục đích.

### 2. Không thể đọc key files

```
ERROR: key file not found in any search path
```

**Giải pháp**: Kiểm tra xem các file key có tồn tại không và quyền đọc.

### 3. Mã hóa/giải mã thất bại

```
ERROR: EncryptValueMKV error
ERROR: DecryptValueMKV error
```

**Giải pháp**: Kiểm tra log file để xem chi tiết lỗi và đảm bảo keys được khởi tạo đúng.

### 4. Password không đọc được

```
WARN: Failed to read password file: ..., using default password
```

**Giải pháp**: Tạo file `password.txt` với nội dung password mong muốn.

## Tích hợp với Fabric

KeyManager được thiết kế để tích hợp dễ dàng với Hyperledger Fabric:

1. **StateDB encryption**: Sử dụng để mã hóa dữ liệu trong state database
2. **Chaincode security**: Có thể sử dụng trong chaincode để mã hóa dữ liệu nhạy cảm
3. **Configuration**: Có thể cấu hình thông qua file `password.txt`

## Ví dụ hoàn chỉnh

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb/mkv"
)

func main() {
    // Lấy KeyManager instance (tự động khởi tạo)
    keyManager := mkv.GetKeyManager()
    
    // Kiểm tra trạng thái
    status := keyManager.GetStatus()
    fmt.Printf("KeyManager initialized: %v\n", status["initialized"])
    
    // Test mã hóa/giải mã
    testData := []byte("Hello, MKV encryption!")
    
    // Mã hóa
    encrypted := mkv.EncryptValueMKV(testData)
    if encrypted == nil {
        log.Fatal("Encryption failed!")
    }
    fmt.Printf("Encrypted data length: %d\n", len(encrypted))
    
    // Giải mã
    decrypted := mkv.DecryptValueMKV(encrypted)
    if decrypted == nil {
        log.Fatal("Decryption failed!")
    }
    fmt.Printf("Decrypted data: %s\n", string(decrypted))
    
    // So sánh kết quả
    if string(decrypted) == string(testData) {
        fmt.Println("✅ Encryption/Decryption successful!")
    } else {
        fmt.Println("❌ Data mismatch!")
    }
}
```

## Hoàn thành!

Hệ thống KeyManager đã hoạt động hoàn hảo với:

- ✅ **Singleton pattern** với `sync.Once`
- ✅ **Tự động khởi tạo** hệ thống khóa
- ✅ **Thread-safe** với `sync.RWMutex`
- ✅ **Fallback mechanism** khi khởi tạo thất bại
- ✅ **Logging** đầy đủ vào file
- ✅ **Tests** hoàn chỉnh và pass
- ✅ **Tích hợp dễ dàng** với Hyperledger Fabric
