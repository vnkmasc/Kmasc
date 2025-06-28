Certificate Management System
Hệ thống quản lý văn bằng, chứng chỉ và hồ sơ sinh viên, cho phép nhà trường cấp phát, lưu trữ và xác minh văn bằng điện tử thông qua công nghệ blockchain. Dự án được xây dựng với kiến trúc hiện đại, tích hợp các công nghệ như Go, React.js, Hyperledger Fabric và nhiều công cụ khác để đảm bảo tính bảo mật, minh bạch và dễ sử dụng.
Tính năng chính

Quản lý văn bằng và chứng chỉ: Cấp phát, lưu trữ và xác minh văn bằng điện tử.
Tích hợp blockchain: Sử dụng Hyperledger Fabric để đảm bảo tính bất biến và minh bạch của dữ liệu.
Xác thực người dùng: Hỗ trợ đăng nhập/đăng ký với JWT và quản lý vai trò.
Lưu trữ tệp: Sử dụng MinIO để lưu trữ tệp văn bằng/chứng chỉ an toàn.
Giao diện người dùng: Frontend thân thiện với người dùng, được xây dựng bằng React.js và Tailwind CSS.
Quản lý khoa và trường: Hỗ trợ quản lý thông tin khoa và trường đại học.
Xác minh mã: Tạo và xác minh mã xác thực cho văn bằng/chứng chỉ.

Công nghệ sử dụng
Backend

Go: Ngôn ngữ lập trình chính, sử dụng framework Gin để xây dựng API.
MongoDB: Cơ sở dữ liệu NoSQL để lưu trữ thông tin người dùng, văn bằng, khoa, trường, v.v.
MinIO: Hệ thống lưu trữ tệp phân tán cho văn bằng và chứng chỉ.
JWT: Xác thực và phân quyền người dùng.
Hyperledger Fabric: Blockchain để lưu trữ và xác minh văn bằng.

Frontend

React.js: Thư viện JavaScript để xây dựng giao diện người dùng.
Tailwind CSS: Framework CSS để thiết kế giao diện hiện đại và responsive.

DevOps

Docker & Docker Compose: Container hóa ứng dụng và các dịch vụ liên quan.
Makefile: Tự động hóa các tác vụ xây dựng và triển khai.

Khác

Digital Signature: Đảm bảo tính xác thực và toàn vẹn của văn bằng.

Cấu trúc thư mục
├── cmd
│   └── server
│       ├── admin_seeder.go        # Khởi tạo tài khoản admin
│       ├── main.go                # Điểm vào chính của ứng dụng
│       └── validator.go           # Xử lý xác thực dữ liệu đầu vào
├── docker-compose.yml             # Cấu hình Docker Compose
├── Dockerfile                     # Cấu hình Docker cho ứng dụng
├── go.mod                         # Quản lý phụ thuộc Go
├── go.sum                         # Checksum của phụ thuộc
├── internal
│   ├── common                     # Các tiện ích chung
│   ├── handlers                   # Xử lý các yêu cầu HTTP
│   ├── middleware                 # Middleware
│   ├── models                     # Định nghĩa các mô hình dữ liệu
│   ├── repository                 # Tương tác với cơ sở dữ liệu
│   └── service                    # Logic nghiệp vụ
├── pkg
│   └── database
│       ├── minio.go               # Kết nối và tương tác với MinIO
│       └── mongodb.go             # Kết nối và tương tác với MongoDB
├── routes
│   └── router.go                  # Định nghĩa các tuyến API
├── utils                          # Các tiện ích hỗ trợ
└── web                            # Mã nguồn frontend (React.js)

Yêu cầu cài đặt

Go: Phiên bản sử dụng: 1.24.4.
Docker & Docker Compose: Để chạy các dịch vụ MongoDB, MinIO.
Git: Để clone mã nguồn từ GitHub.
Node.js: Phiên bản 18.x trở lên.
Hyperledger Fabric: Yêu cầu thiết lập môi trường blockchain (xem tài liệu chính thức của Hyperledger Fabric).

Hướng dẫn cài đặt và chạy
1. Clone dự án
git clone https://github.com/tuyenngduc/certificate-management-system.git
cd certificate-management-system

2. Thiết lập môi trường

Đảm bảo bạn đã cài đặt Go, Docker, Docker Compose và Git.
Cài đặt các phụ thuộc Go:
go mod tidy

3. Thiết lập file môi trường (.env)
Tạo một file .env trong thư mục gốc của dự án.
Sao chép các biến môi trường cần thiết từ mẫu dưới đây và thay thế bằng các giá trị phù hợp với môi trường của bạn:

# MongoDB Configuration
MONGODB_URI=mongodb://<username>:<password>@<host>:<port>
DB_NAME=<database_name>
MONGO_INITDB_ROOT_USERNAME=<mongo_username>
MONGO_INITDB_ROOT_PASSWORD=<mongo_password>

# Admin Account Configuration
ADMIN_EMAIL=<admin_email>
ADMIN_PASSWORD=<admin_password>

# Email Configuration
EMAIL_FROM=<email_address>
EMAIL_PASSWORD=<email_password_or_app_password>
EMAIL_HOST=<smtp_host>
EMAIL_PORT=<smtp_port>

# JWT Configuration
JWT_SECRET=<random_secure_string>

# MinIO Configuration
MINIO_ROOT_USER=<minio_username>
MINIO_ROOT_PASSWORD=<minio_password>
MINIO_ENDPOINT=<minio_host>:<minio_port>
MINIO_ACCESS_KEY=<minio_access_key>
MINIO_SECRET_KEY=<minio_secret_key>
MINIO_BUCKET=<bucket_name>
MINIO_USE_SSL=<true_or_false>

4. Khởi động các dịch vụ

Sử dụng Docker Compose để khởi động MongoDB và MinIO:

docker-compose up --build -d

5. Chạy server backend

Chạy ứng dụng Go:

go run cmd/server/main.go

Server sẽ chạy mặc định trên http://localhost:8080.

6. Thiết lập Hyperledger Fabric

Tham khảo tài liệu chính thức của Hyperledger Fabric để thiết lập mạng blockchain.
Đảm bảo tích hợp các chaincode cần thiết để lưu trữ và xác minh văn bằng.

Tác giả: Tuyen Nguyen Duc
Email: tuyenngduc12@gmail.com
GitHub: tuyenngduc

