# MKV256 Encryption Module for Go (cgo)

Thư mục này chứa mã nguồn thư viện mã hóa MKV256 (thuần C) và Go wrapper để sử dụng trong các dự án Go (ví dụ: Hyperledger Fabric statedb).

## Nội dung
- `MKV256.c`, `MKV256.h`, `PrecomputedTable256.h`: Thuật toán mã hóa MKV256.
- `mkv.c`, `mkv.h`: Hàm mã hóa/giải mã dữ liệu dài, padding PKCS#7, interface tương thích statedb.
- `mkv.go`: Go wrapper sử dụng cgo để gọi thư viện C.
- `mkv_test.go`: Unit test kiểm tra mã hóa/giải mã từ Go.
- `Makefile`: Build shared library `libmkv.so`.

## Hướng dẫn build
1. **Build thư viện C:**
   ```bash
   make clean && make
   ```
   Kết quả sẽ tạo ra file `libmkv.so` trong thư mục này.

2. **Chạy unit test Go:**
   ```bash
   LD_LIBRARY_PATH=. go test -v
   ```
   > Lưu ý: Biến môi trường `LD_LIBRARY_PATH` cần trỏ tới thư mục chứa `libmkv.so`.

## Tích hợp vào Go project khác
- Import package Go này (ví dụ: `import "github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb/mkv"`).
- Đảm bảo file `libmkv.so` nằm trong thư mục runtime hoặc thư viện hệ thống (`/usr/local/lib`).
- Sử dụng các hàm:
  - `EncryptValueMKV(value []byte, key []byte) []byte`
  - `DecryptValueMKV(value []byte, key []byte) []byte`

## Lưu ý
- Không để file C nào có hàm `main` trong thư mục này khi build với Go/cgo.
- Key phải đúng độ dài (32 bytes cho 256 bit).
- Chỉ hỗ trợ block size 256 bit (32 bytes) cho MKV256.

## Liên hệ
Nếu gặp vấn đề hoặc cần hỗ trợ, hãy liên hệ tác giả hoặc mở issue trong repo. 