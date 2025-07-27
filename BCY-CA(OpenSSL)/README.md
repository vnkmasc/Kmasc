# BCY-CA: Giải Pháp Thay Thế Fabric-CA trong Hyperledger Fabric

## Mục Tiêu

BCY-CA được phát triển nhằm thay thế hoàn toàn cho các công cụ [cryptogen](https://hyperledger-fabric.readthedocs.io/en/latest/commands/cryptogen.html) và [Fabric-CA](https://hyperledger-fabric-ca.readthedocs.io/en/latest/) trong môi trường Hyperledger Fabric, mang lại hệ thống quản lý danh tính linh hoạt, tích hợp cao qua API hoặc shell script.

BCY-CA dùng OpenSSL kết hợp với Go modules và bindings do nhóm giảng viên phát triển, đảm bảo:

- Tạo và quản lý Root CA.
- Đăng ký, xác thực danh tính người dùng.
- Ký CSR để sinh chứng chỉ.
- Hỗ trợ TLS certificate riêng biệt.
- Xuất ra cấu trúc thư mục MSP hợp lệ theo chuẩn Hyperledger Fabric.

---

## Tổng Quan

BCY-CA là một Certificate Authority tuỳ chỉnh thay thế Fabric-CA trong mối trường Hyperledger Fabric. Hệ thống cung cấp:

- Giao diện API HTTP dễ sử dụng.
- Shell script để thiết lập nhanh các tổ chức Peer/Orderer.
- Tích hợp OpenSSL và module Go để xử lý CA.

### Các module chính:

- `internal/api/handlers.go`: Xử lý request CA.
- `internal/ca/openssl.go`: Wrapper gọi lệnh OpenSSL.
- `internal/config/config.go`: Tổng hợp config.

---

## Tính Năng

| Phương Thức | Endpoint     | Chức Năng                         | Mô Tả                                        |
| ----------- | ------------ | --------------------------------- | -------------------------------------------- |
| `POST`      | `/initca`    | Khởi tạo Root CA                  | Sinh cặp ca-cert.pem và ca-key.pem           |
| `POST`      | `/register`  | Đăng ký người dùng                | Trả về ID và Secret dùng cho enroll          |
| `POST`      | `/enroll`    | Ký chứng chỉ người dùng           | Nhận CSR, ký bằng CA, trả về cert            |
| `POST`      | `/tlsenroll` | Ký chứng chỉ TLS cho peer/orderer | Giống enroll, nhưng sinh cert TLS riêng biệt |

---

## Hướng Dẫn Triển Khai

### 1. Cài Đặt

- Mở file `registerEnroll.sh`.
- Chỉnh sửa dòng sau:

```bash
CA_URL="http://127.0.0.1:7054"
```

- Copy script này vào thư mục `fabric-samples/test-network`.

### 2. Khởi Tạo CA & Tổ Chức

Thực thi các lệnh sau trong `fabric-samples/test-network`:

```bash
./registerEnroll.sh initCA                # Tạo Root CA
./registerEnroll.sh setupPeerOrg org1.example.com
./registerEnroll.sh setupPeerOrg org2.example.com
./registerEnroll.sh setupOrdererOrg
```

### 3. Cấu Hình TLS

Copy CA TLS cert vào MSP của orderer:

```bash
mkdir -p organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts

cp organizations/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem \
   organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/CA
```

### 4. Khởi Chạy Network

```bash
./network.sh up createChannel -c mychannel -s couchdb
```

> **Lưu ý:** KHÔNG dùng option `-ca`, vì sẽ kích hoạt Fabric-CA và ghi đè CA do BCY-CA sinh ra.

---

## Cấu Hình và OpenSSL

### 1. File config

- Định nghĩa trong: `internal/config/config.go`
- Bao gồm cấu hình đường dẫn, port server, path cert/key, ...

### 2. Module OpenSSL

- Wrapper OpenSSL trong Go: `internal/ca/openssl.go`

### 3. Kiểm tra chứng chỉ

```bash
openssl x509 -in organizations/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem -text -noout
```

---

## Ưu Điểm BCY-CA

- ✨ **Linh hoạt**: Quản lý danh tính qua API hoặc script.
- ⚙️ **Tùy chỉnh cao**: Module Go & OpenSSL cho phép tích hợp theo nhu cầu.
- 🏛️ **Chuẩn Hyperledger Fabric**: Sinh MSP directory hợp lệ.
- 🔒 **Hỗ trợ TLS**: Sinh TLS cert riêng cho peer/orderer.

---

## Kết Luận

BCY-CA cung cấp giải pháp thay thế toàn diện cho Fabric-CA, phù hợp với những dự án cần quy trình tự động hóa cao, có khả năng tùy biến linh hoạt. Đây là bước quan trọng tiến tới môi trường Hyperledger Fabric "thuần Việt hoá".

