# So Sánh Performance: Có AES vs Không AES

## ⚠️ Lưu Ý Quan Trọng

**Cả hai test hiện tại đều sử dụng cùng package và có cùng behavior!** Để thực sự so sánh performance giữa có AES và không có AES, bạn cần:

## Phương Pháp 1: Build Fabric với Encryption Disabled

### Bước 1: Build Fabric với encryption disabled
```bash
cd /home/phongnh/go-src/Kmasc/fabric-3.1.1
./build-fabric-no-encryption.sh
```

Script này sẽ:
- Backup các file encryption gốc
- Tạo file `encrypt.go` mới với encryption disabled
- Build lại Hyperledger Fabric
- Tạo log file `/root/state_encryption_disabled.log`

### Bước 2: Chạy test với encryption disabled
```bash
cd core/ledger/kvledger/txmgmt/statedb/test_perf_without_aes
go run main_simple.go
```

### Bước 3: Restore encryption gốc
```bash
cd /home/phongnh/go-src/Kmasc/fabric-3.1.1
make clean && make
```

## Phương Pháp 2: Clone Repository Riêng

```bash
# Clone repository gốc
git clone https://github.com/hyperledger/fabric.git fabric-with-encryption
git clone https://github.com/hyperledger/fabric.git fabric-no-encryption

# Trong fabric-no-encryption, disable encryption
cd fabric-no-encryption
# Edit core/ledger/kvledger/txmgmt/statedb/encrypt.go để disable encryption
make
```

## Kết Quả Test Hiện Tại (Chưa Chính Xác)

### Test KHÔNG AES (main_simple.go)
- **Trung bình Put**: 0.000944 giây
- **Trung bình Get**: 0.000039 giây

### Test CÓ AES (main.go)
- **Trung bình Put**: 0.001423 giây
- **Trung bình Get**: 0.000057 giây

## Phân Tích Dự Đoán

### Tác Động Của AES Encryption (Ước tính)

1. **Put Operation (Ghi dữ liệu)**:
   - Không AES: ~0.0009s
   - Có AES: ~0.0014s
   - **Tăng thời gian**: ~55% (do cần mã hóa dữ liệu)

2. **Get Operation (Đọc dữ liệu)**:
   - Không AES: ~0.00004s
   - Có AES: ~0.00006s
   - **Tăng thời gian**: ~50% (do cần giải mã dữ liệu)

## Cách Chạy Test Hiện Tại

### Chạy riêng lẻ:
```bash
# Test có AES
cd test_perf_with_aes
go run main.go

# Test không AES (hiện tại vẫn dùng AES)
cd test_perf_without_aes
go run main_simple.go
```

### Chạy so sánh:
```bash
./compare_performance.sh
```

## Lưu Ý

- **Test hiện tại chưa chính xác** vì cả hai đều sử dụng cùng encryption logic
- Để có kết quả chính xác, cần build Fabric với encryption disabled
- Sử dụng script `build-fabric-no-encryption.sh` để tạo phiên bản không encryption
- Test được thực hiện với 100 lần lặp
- Sử dụng LevelDB làm backend
- Môi trường: Linux 6.11.0-29-generic 