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

##Luồng hoạt động của BCY-CA

Luồng hoạt động của BCY-CA được thiết kế để quản lý danh tính và chứng chỉ một cách tuần tự và hiệu quả. Dưới đây là mô tả chi tiết về từng bước xử lý:

1. Khởi tạo Root CA (/initca)

Input: Không yêu cầu input cụ thể, chỉ cần gọi API.

Quy trình:

Tạo cặp khóa RSA (private key và public key) bằng OpenSSL.

Tạo chứng chỉ tự ký (self-signed certificate) cho Root CA.

Lưu ca-cert.pem và ca-key.pem vào thư mục cấu hình (ví dụ: organizations/ca/).

Output: File ca-cert.pem và ca-key.pem được lưu trong thư mục chỉ định.

Script tương ứng: ./registerEnroll.sh initCA

2. Đăng ký người dùng (/register)

Input: ID người dùng (ví dụ: admin, user1) và thông tin tổ chức (ví dụ: org1.example.com).

Quy trình:

Xác thực yêu cầu đăng ký.

Tạo một secret (mật khẩu tạm thời).

Lưu ID và secret vào cơ sở dữ liệu nội bộ hoặc file cấu hình.

Output: Trả về ID và secret cho người dùng.

Script tương ứng: Được gọi trong ./registerEnroll.sh setupPeerOrg hoặc setupOrdererOrg

3. Ký chứng chỉ người dùng (/enroll)

Input: ID, secret, và CSR (Certificate Signing Request).

Quy trình:

Xác thực ID và secret.

Kiểm tra CSR hợp lệ.

Ký CSR bằng ca-key.pem.

Tạo chứng chỉ người dùng và cấu trúc thư mục MSP.

Output: Chứng chỉ người dùng và thư mục MSP (ví dụ: organizations/peerOrganizations/org1.example.com/users/).

Script tương ứng: ./registerEnroll.sh setupPeerOrg

4. Ký chứng chỉ TLS (/enroll/tls)

Input: CSR TLS từ peer hoặc orderer.

Quy trình:

Xác thực yêu cầu TLS.

Ký CSR bằng CA để tạo TLS cert.

Đảm bảo chứa các extension như subjectAltName.

Lưu chứng chỉ vào thư mục TLS (ví dụ: tls/peer0.org1.example.com/).

Output: Chứng chỉ TLS và thư mục TLS MSP.

Script tương ứng: ./registerEnroll.sh setupPeerOrg hoặc setupOrdererOrg

5. Tạo cấu trúc MSP

Quy trình:

Sau khi ký, BCY-CA tự động tạo cấu trúc MSP đúng chuẩn Hyperledger Fabric.

Bao gồm:

cacerts/: Chứng chỉ CA.

tlscacerts/: Chứng chỉ TLS CA.

keystore/: Khóa riêng.

signcerts/: Chứng chỉ đã ký.

config.yaml: Cấu hình tổ chức.

Output: Thư mục MSP hoàn chỉnh.

## Hướng Dẫn Triển Khai
Chạy customCA server: go run ./cmd/main.go

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
   organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts
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

