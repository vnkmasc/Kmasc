# MKV Password Management - Production Deployment Guide

## 🎯 **Tổng quan hệ thống**

Hệ thống MKV Password Management cung cấp khả năng **đổi password trong runtime** cho Hyperledger Fabric network đang chạy, sử dụng:

- **PBKDF2-HMAC-SHA256** với 10,000 iterations
- **REST API Server** cho quản lý từ xa
- **Client tools** cho automation
- **Salt ngẫu nhiên** cho mỗi deployment

---

## 🚀 **Quick Start Demo**

```bash
# 1. Khởi tạo hệ thống
./mkv_client.sh init "your_initial_password"

# 2. Kiểm tra trạng thái
./mkv_client.sh status

# 3. Test password hiện tại
./mkv_client.sh test "your_initial_password"

# 4. Đổi password (TRONG KHI HỆ THỐNG ĐANG CHẠY!)
./mkv_client.sh change "your_initial_password" "new_secure_password"

# 5. Verify password mới
./mkv_client.sh test "new_secure_password"
```

---

## 🏭 **Production Deployment**

### Step 1: Chuẩn bị môi trường

```bash
# Build MKV library
cd /path/to/mkv
make clean && make

# Test hệ thống
LD_LIBRARY_PATH=. go test -v

# Khởi tạo với password mạnh
./init_with_pbkdf2.sh "$(openssl rand -base64 32)"
```

### Step 2: Khởi động API Server

```bash
# Set environment variables
export MKV_API_PORT=9876
export MKV_API_KEY="$(openssl rand -hex 32)"

# Start server trong production
LD_LIBRARY_PATH=. nohup go run ../mkv_api_server.go > /var/log/mkv_api.log 2>&1 &

# Hoặc sử dụng systemd service
sudo systemctl start mkv-api
```

### Step 3: Security Configuration

```bash
# Firewall rules (chỉ cho phép internal network)
sudo ufw allow from 10.0.0.0/8 to any port 9876
sudo ufw allow from 172.16.0.0/12 to any port 9876
sudo ufw allow from 192.168.0.0/16 to any port 9876

# File permissions
chmod 600 *.key
chmod 644 *.log
chown fabric:fabric *.key
```

---

## 🔧 **API Endpoints**

### Authentication
Tất cả endpoints (trừ `/health`) yêu cầu header:
```
X-API-Key: your_api_key_here
```

### Endpoints

#### 1. Health Check
```bash
GET /api/v1/health
# Response: {"status":"healthy","time":"2025-08-09T03:47:54Z"}
```

#### 2. System Status
```bash
GET /api/v1/status
# Response: {
#   "status":"success",
#   "message":"System is initialized and ready",
#   "k1_exists":true,
#   "salt_exists":true,
#   "encrypted_k1_exists":true
# }
```

#### 3. Initialize System
```bash
POST /api/v1/initialize
Content-Type: application/json
{
  "password": "your_secure_password"
}
```

#### 4. Change Password (🔥 RUNTIME!)
```bash
POST /api/v1/change-password
Content-Type: application/json
{
  "old_password": "current_password",
  "new_password": "new_secure_password"
}
```

#### 5. Test Password
```bash
POST /api/v1/test-password
Content-Type: application/json
{
  "password": "password_to_test"
}
# Response: {"valid":true} or {"valid":false,"error":"..."}
```

---

## 🔄 **Integration với Fabric Network**

### 1. Trong Peer Configuration

```yaml
# core.yaml
ledger:
  state:
    stateDatabase: mkv
    mkv:
      apiEndpoint: "http://localhost:9876"
      apiKey: "${MKV_API_KEY}"
      passwordRotationInterval: "24h"
```

### 2. Password Rotation Script

```bash
#!/bin/bash
# /opt/fabric/scripts/rotate_mkv_password.sh

NEW_PASSWORD=$(openssl rand -base64 32)
API_KEY="${MKV_API_KEY}"

# Change password via API
curl -X POST \
  -H "Content-Type: application/json" \
  -H "X-API-Key: $API_KEY" \
  -d "{\"old_password\":\"$OLD_PASSWORD\",\"new_password\":\"$NEW_PASSWORD\"}" \
  http://localhost:9876/api/v1/change-password

# Update environment variable
echo "MKV_PASSWORD=$NEW_PASSWORD" > /opt/fabric/config/mkv.env

# Restart peers (nếu cần)
systemctl restart peer0 peer1
```

### 3. Monitoring Integration

```bash
# Prometheus metrics endpoint
GET /api/v1/metrics

# Health check cho monitoring
*/5 * * * * curl -f http://localhost:9876/api/v1/health || alert-manager
```

---

## 🛡️ **Security Best Practices**

### 1. Password Policy
```bash
# Generate strong passwords
openssl rand -base64 32

# Rotate passwords định kỳ
0 2 * * 0 /opt/fabric/scripts/rotate_mkv_password.sh
```

### 2. Network Security
```bash
# TLS termination với nginx
server {
    listen 443 ssl;
    server_name mkv-api.internal;
    
    location /api/ {
        proxy_pass http://localhost:9876;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 3. Audit Logging
```bash
# Log all API calls
tail -f /var/log/mkv_api.log | grep -E "(CHANGE_PASSWORD|INIT_KEYS)"
```

---

## 📊 **Monitoring & Alerting**

### 1. Health Monitoring
```bash
# Nagios/Zabbix check
#!/bin/bash
response=$(curl -s http://localhost:9876/api/v1/health)
if [[ $? -eq 0 && "$response" == *"healthy"* ]]; then
    echo "OK - MKV API is healthy"
    exit 0
else
    echo "CRITICAL - MKV API is down"
    exit 2
fi
```

### 2. Performance Metrics
```bash
# Response time monitoring
curl -w "@curl-format.txt" -s -o /dev/null http://localhost:9876/api/v1/status
```

### 3. Security Alerts
```bash
# Failed authentication attempts
grep "Unauthorized" /var/log/mkv_api.log | wc -l
```

---

## 🚨 **Disaster Recovery**

### 1. Backup Strategy
```bash
# Backup keys và salt
tar -czf mkv_backup_$(date +%Y%m%d).tar.gz *.key *.log

# Encrypted backup
gpg --cipher-algo AES256 --compress-algo 1 --symmetric mkv_backup.tar.gz
```

### 2. Recovery Process
```bash
# Restore từ backup
tar -xzf mkv_backup_20250809.tar.gz

# Verify integrity
./mkv_client.sh test "recovery_password"

# Restart services
systemctl restart mkv-api peer0 peer1
```

### 3. Emergency Password Reset
```bash
# Nếu quên password, sử dụng K1 plaintext (nếu có)
./mkv_client.sh init "new_emergency_password"

# Hoặc restore từ backup
```

---

## 📈 **Scaling & High Availability**

### 1. Load Balancer Setup
```nginx
upstream mkv_api {
    server 10.0.1.10:9876;
    server 10.0.1.11:9876;
    server 10.0.1.12:9876;
}

server {
    location /api/ {
        proxy_pass http://mkv_api;
    }
}
```

### 2. Database Replication
```bash
# Sync keys across nodes
rsync -avz *.key user@backup-node:/opt/mkv/
```

---

## 🔍 **Troubleshooting**

### Common Issues

#### 1. "Password test failed"
```bash
# Check salt file exists
ls -la k0_salt.key

# Verify API server is running
curl http://localhost:9876/api/v1/health

# Check logs
tail -f /var/log/mkv_api.log
```

#### 2. "API Server not responding"
```bash
# Check process
ps aux | grep mkv_api_server

# Check port binding
netstat -tlnp | grep 9876

# Restart server
pkill -f mkv_api_server
LD_LIBRARY_PATH=. go run ../mkv_api_server.go &
```

#### 3. "Permission denied"
```bash
# Fix file permissions
chmod 600 *.key
chown fabric:fabric *.key

# Check SELinux context
ls -Z *.key
```

---

## 🎉 **Success Criteria**

✅ **Hệ thống hoạt động nếu:**
- Health check returns `{"status":"healthy"}`
- Password test với current password returns `{"valid":true}`
- Có thể đổi password thành công trong runtime
- Old password fails, new password works
- All files có correct permissions
- API server responds trong < 1s

---

## 📞 **Support & Maintenance**

### Regular Tasks
- [ ] Weekly password rotation
- [ ] Monthly backup verification
- [ ] Quarterly security audit
- [ ] Annual disaster recovery test

### Contacts
- **Technical Lead**: [Your Name]
- **Security Team**: security@company.com
- **Operations**: ops@company.com

---

**🚀 Hệ thống MKV Password Management đã sẵn sàng cho production!**
