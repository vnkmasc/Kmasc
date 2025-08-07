# So Sánh Performance: Có AES vs Không AES

## Build Fabric với Encryption Disabled

### Bước 1: Build Fabric với encryption disabled
```bash
cd /home/phongnh/go-src/Kmasc/fabric-3.1.1
./build-fabric-no-encryption.sh
```

Script này sẽ:
- Backup các file encryption gốc
- Tạo file `encrypt.go` mới với encryption disabled
- Build lại Hyperledger Fabric

### Bước 2: Chạy test với encryption disabled
```bash
cd core/ledger/kvledger/txmgmt/statedb/test_perf_without_aes
go run main.go
```

### Bước 3: Restore encryption gốc
```bash
cd /home/phongnh/go-src/Kmasc/fabric-3.1.1
make clean && make
```

## Cách Chạy Test Hiện Tại

### Chạy riêng lẻ:
```bash
# Test có AES
cd test_perf_with_aes
go run main.go

# Test không AES (khi đã build lại với encryption disabled)
cd test_perf_without_aes
go run main.go
```

## Lưu Ý

- **Test hiện tại chưa chính xác** vì cả hai đều sử dụng cùng encryption logic
- Để có kết quả chính xác, cần build Fabric với encryption disabled
- Sử dụng script `build-fabric-no-encryption.sh` để tạo phiên bản không encryption
- Test được thực hiện với 100 lần lặp
- Sử dụng LevelDB làm backend