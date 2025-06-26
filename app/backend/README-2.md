# Hệ thống Quản lý Văn bằng - Certificate Management System

## Tổng quan

Hệ thống quản lý văn bằng, chứng chỉ và hồ sơ sinh viên, cho phép nhà trường cấp phát, lưu trữ và xác minh văn bằng điện tử thông qua công nghệ blockchain. Dự án được xây dựng với kiến trúc hiện đại, tích hợp các công nghệ như Go, Next.js, Hyperledger Fabric và nhiều công cụ khác để đảm bảo tính bảo mật, minh bạch và dễ sử dụng.

## Tính năng chính

- **Quản lý văn bằng và chứng chỉ**: Cấp phát, lưu trữ và xác minh văn bằng điện tử.
- **Tích hợp blockchain**: Sử dụng Hyperledger Fabric để đảm bảo tính bất biến và minh bạch của dữ liệu.
- **Xác thực người dùng**: Hỗ trợ đăng nhập/đăng ký với JWT để quản lý authenciator/authorization.
- **Lưu trữ tệp**: Sử dụng MinIO để lưu trữ tệp văn bằng/chứng chỉ an toàn.
- **Giao diện người dùng**: Frontend được xây dựng bằng NextJS, Shadcn và TailwindCSS.
- **Quản lý khoa và trường**: Hỗ trợ quản lý thông tin khoa và trường đại học.
- **Xác minh mã**: Tạo và xác minh mã xác thực cho văn bằng/chứng chỉ.

## Kiến trúc Cơ sở Dữ liệu

### 1. Model Account (Tài khoản)

- **Mô tả**: Quản lý thông tin đăng nhập của người dùng
- **Bảng**: `accounts`
- **Trường chính**:
  - `_id`: ID tài khoản
  - `student_id`: ID sinh viên (nếu là sinh viên)
  - `university_id`: ID trường đại học (nếu là admin trường)
  - `student_email`: Email trường (@domain.edu.vn)
  - `personal_email`: Email cá nhân (dùng để đăng nhập)
  - `password_hash`: Mật khẩu đã mã hóa
  - `role`: Vai trò (student, university_admin, admin)

### 2. Model User (Sinh viên)

- **Mô tả**: Thông tin chi tiết sinh viên
- **Bảng**: `users`
- **Trường chính**:
  - `_id`: ID sinh viên
  - `student_code`: Mã sinh viên
  - `full_name`: Họ tên đầy đủ
  - `email`: Email sinh viên
  - `faculty_id`: ID khoa
  - `university_id`: ID trường đại học
  - `course`: Khóa học (K47, K48, ...)
  - `status`: Trạng thái tốt nghiệp (0: chưa, 1: cử nhân, 2: kỹ sư, 3: thạc sĩ, 4: tiến sĩ)

### 3. Model University (Trường Đại học)

- **Mô tả**: Thông tin các trường đại học
- **Bảng**: `universities`
- **Trường chính**:
  - `_id`: ID trường
  - `university_name`: Tên trường
  - `university_code`: Mã trường
  - `address`: Địa chỉ
  - `email_domain`: Tên miền email (@actvn.edu.vn)
  - `status`: Trạng thái (pending, approved, rejected)

### 4. Model Faculty (Khoa)

- **Mô tả**: Thông tin các khoa trong trường
- **Bảng**: `faculties`
- **Trường chính**:
  - `_id`: ID khoa
  - `faculty_code`: Mã khoa
  - `faculty_name`: Tên khoa
  - `university_id`: ID trường đại học

### 5. Model Certificate (Văn bằng/Chứng chỉ)

- **Mô tả**: Thông tin văn bằng và chứng chỉ
- **Bảng**: `certificates`
- **Trường chính**:
  - `_id`: ID văn bằng
  - `user_id`: ID sinh viên
  - `faculty_id`: ID khoa
  - `university_id`: ID trường
  - `student_code`: Mã sinh viên
  - `is_degree`: Có phải văn bằng không (true: văn bằng, false: chứng chỉ)
  - `certificate_type`: Loại văn bằng (Cử nhân, Kỹ sư, Thạc sĩ, Tiến sĩ)
  - `name`: Tên chứng chỉ (nếu không phải văn bằng)
  - `serial_number`: Số hiệu văn bằng
  - `registration_number`: Số vào sổ
  - `path`: Đường dẫn tệp PDF
  - `issue_date`: Ngày cấp
  - `signed`: Đã ký số hay chưa

### 6. Model VerificationCode (Mã xác thực)

- **Mô tả**: Mã xác thực để chia sẻ thông tin văn bằng
- **Bảng**: `verification_codes`
- **Trường chính**:
  - `_id`: ID mã xác thực
  - `user_id`: ID sinh viên
  - `code`: Mã xác thực (8 ký tự)
  - `can_view_score`: Có thể xem điểm
  - `can_view_data`: Có thể xem thông tin
  - `can_view_file`: Có thể xem tệp
  - `expired_at`: Thời gian hết hạn

### 7. Model OTP

- **Mô tả**: Mã OTP để xác thực email khi đăng ký
- **Collection**: `otps`
- **Trường chính**:
  - `email`: Email sinh viên
  - `code`: Mã OTP (6 số)
  - `expires_at`: Thời gian hết hạn

## Phân quyền và Chức năng API

### 1. ADMIN SYSTEM (Quản trị viên hệ thống)

#### Quản lý Trường Đại học

- `POST /api/v1/universities` - Tạo yêu cầu đăng ký trường mới
- `GET /api/v1/universities` - Xem danh sách tất cả trường
- `GET /api/v1/universities/status?status=pending` - Xem trường theo trạng thái
- `POST /api/v1/universities/approve-or-reject` - Phê duyệt/từ chối trường

#### Quản lý Tài khoản

- `GET /api/v1/auth/accounts` - Xem tất cả tài khoản
- `DELETE /api/v1/auth/accounts` - Xóa tài khoản
- `GET /api/v1/auth/university-admin-info` - Xem thông tin admin trường
- `GET /api/v1/auth/students-info` - Xem thông tin tài khoản sinh viên

### 2. UNIVERSITY ADMIN (Quản trị viên Trường)

#### Quản lý Khoa

- `POST /api/v1/faculties` - Tạo khoa mới
- `GET /api/v1/faculties` - Xem danh sách khoa
- `PUT /api/v1/faculties/:id` - Cập nhật thông tin khoa
- `DELETE /api/v1/faculties/:id` - Xóa khoa
- `GET /api/v1/faculties/:id` - Xem chi tiết khoa

#### Quản lý Sinh viên

- `POST /api/v1/users` - Tạo sinh viên mới
- `POST /api/v1/users/import-excel` - Import sinh viên từ Excel
- `GET /api/v1/users` - Xem danh sách sinh viên
- `GET /api/v1/users/search` - Tìm kiếm sinh viên
- `GET /api/v1/users/:id` - Xem chi tiết sinh viên
- `PUT /api/v1/users/:id` - Cập nhật thông tin sinh viên
- `DELETE /api/v1/users/:id` - Xóa sinh viên
- `GET /api/v1/users/faculty/:faculty_code` - Xem sinh viên theo khoa

#### Quản lý Văn bằng/Chứng chỉ

- `POST /api/v1/certificates` - Tạo văn bằng/chứng chỉ
- `POST /api/v1/certificates/import-excel` - Import từ Excel
- `GET /api/v1/certificates` - Xem danh sách văn bằng
- `GET /api/v1/certificates/search` - Tìm kiếm văn bằng
- `GET /api/v1/certificates/:id` - Xem chi tiết văn bằng
- `POST /api/v1/certificates/upload-pdf` - Upload tệp PDF
- `GET /api/v1/certificates/tệp/:id` - Tải tệp văn bằng
- `DELETE /api/v1/certificates/:id` - Xóa văn bằng
- `GET /api/v1/certificates/student/:id` - Xem văn bằng theo sinh viên

### 3. STUDENT (Sinh viên)

#### Đăng ký và Xác thực

- `POST /api/v1/auth/request-otp` - Yêu cầu mã OTP
- `POST /api/v1/auth/verify-otp` - Xác thực OTP
- `POST /api/v1/auth/register` - Đăng ký tài khoản

#### Quản lý Tài khoản

- `POST /api/v1/auth/login` - Đăng nhập
- `POST /api/v1/auth/change-password` - Đổi mật khẩu
- `GET /api/v1/users/me` - Xem thông tin cá nhân

#### Quản lý Văn bằng

- `GET /api/v1/certificates/my-certificate` - Xem văn bằng của mình
- `GET /api/v1/certificates/simple` - Xem danh sách tên văn bằng
- `GET /api/v1/certificates/tệp/:id` - Tải tệp văn bằng

#### Tạo Mã xác thực (Chia sẻ)

- `POST /api/v1/verification/create` - Tạo mã xác thực mới
- `GET /api/v1/verification/my-codes` - Xem các mã đã tạo

### 4. PUBLIC (Không cần đăng nhập)

#### Xác thực Văn bằng

- `POST /api/v1/auth/verification` - Xác thực văn bằng bằng mã

## Quy trình Hoạt động

### 1. Đăng ký Trường Đại học

1. Trường gửi yêu cầu đăng ký đến Admin
2. Admin hệ thống phê duyệt
3. Hệ thống tự động tạo tài khoản admin cho trường
4. Gửi email thông báo tài khoản đến trường

### 2. Quản lý Sinh viên

1. Admin trường tạo các khoa
2. Import danh sách sinh viên từ Excel hoặc tạo thủ công
3. Cấp văn bằng cho sinh viên với tên tệp là mã số hiệu
4. Cấp chứng chỉ cho sinh viên cung cấp loại chứng chỉ và tên tệp là mã sinh viên

### 3. Cấp Văn bằng

1. Admin trường tạo văn bằng/chứng chỉ cho sinh viên bằng
2. Upload tệp PDF văn bằng/chứng chỉ
3. Đánh dấu đã ký số
4. Trường đẩy tệp lên blockchain (có thể không đẩy)
5. Có thể xem văn bằng/chứng chỉ và tải về

### 4. Sinh viên sử dụng hệ thống

1. Nhập tài khoản bằng email trường để nhận yêu cầu OTP
2. Đăng ký tài khoản mới
3. Xem thông tin cơ bản và thông tin chứng chỉ từ hệ thống học viện

### 5. Xác thực Văn bằng

1. Sinh viên tạo mã xác thực với thời hạn
2. Chia sẻ mã với bên thứ ba
3. Bên thứ ba sử dụng mã để xác thực thông tin

## Công nghệ Sử dụng

- **Backend**: Go (Gin Framework)
- **Frontend**: NextJS + TailwindCSS + shadcn/ui
- **Database**: MongoDB
- **File Storage**: MinIO
- **Authentication**: JWT
- **Email**: SMTP
- **File Format**: Excel import/export
- **Docker**: Khởi động MongoDB và MinIO

## Bảo mật

- Mã hóa mật khẩu bằng bcrypt
- JWT token cho xác thực
- Phân quyền theo role
- OTP xác thực email
- Mã xác thực có thời hạn
- Upload tệp an toàn

## Hướng dẫn Cài đặt Dự án

### Yêu cầu Hệ thống

- **Go**: version 1.24.3
- **NodeJS**: version 24.0.1
- **Docker**: version 28.2.2
- **Docker Compose**: version 2.36.2

### 1. Clone Dự án

```bash
# Clone frontend
git clone https://github.com/Tuienn/nextjs-common.git

# Clone backend
git clone https://github.com/tuyenngduc/certificate-management-system.git
```

### 2. Cấu hình Environment Variables

#### Frontend (.env.local)

```env
SESSION_SECRET=<SECRET_KEY>
NEXT_PUBLIC_API_URL=http://localhost:8080
```

#### Backend (.env)

```env
# Database
MONGODB_URI=mongodb://<username>:<password>@<host>:<port>
DB_NAME=<database_name>
MONGO_INITDB_ROOT_USERNAME=<mongo_username>
MONGO_INITDB_ROOT_PASSWORD=<mongo_password>

# Admin Account
ADMIN_EMAIL=<admin_email>
ADMIN_PASSWORD=<admin_password>

# Email Configuration
EMAIL_FROM=<email_address>
EMAIL_PASSWORD=<email_password_or_app_password>
EMAIL_HOST=<smtp_host>
EMAIL_PORT=<smtp_port>

# JWT
JWT_SECRET=<random_secure_string>

# MinIO Configuration
MINIO_ROOT_USER=<minio_username>
MINIO_ROOT_PASSWORD=<minio_password>
MINIO_ENDPOINT=<minio_host>:<minio_port>
MINIO_ACCESS_KEY=<minio_access_key>
MINIO_SECRET_KEY=<minio_secret_key>
MINIO_BUCKET=<bucket_name>
MINIO_USE_SSL=<true_or_false>
```

### 3. Khởi động Dự án

#### Frontend

```bash
# Cài đặt dependencies (chỉ chạy lần đầu)
npm install

# Khởi động development server
npm run dev
# Server NextJS sẽ chạy tại: http://localhost:3000
```

#### Backend

```bash
# Khởi động MongoDB và MinIO với Docker
docker-compose up --build -d
# Khởi động MongoDB và MinIO

# Khởi động Go server
go run ./cmd/server
# Server API sẽ chạy tại: http://localhost:8080
```

### 4. Truy cập Ứng dụng

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **MinIO Console**: http://localhost:9001
