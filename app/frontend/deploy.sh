#!/bin/bash

# Dừng tiến trình Next.js đang chạy bằng PM2
pm2 stop nextjs-fe

# Chuyển đến thư mục frontend
cd /root/Kmasc/app/frontend || exit 1

# Build lại dự án
npm run build

# Start lại Next.js bằng PM2 với tên nextjs-fe
HOST=0.0.0.0 pm2 start npm --name nextjs-fe -- start

# Lưu cấu hình PM2
pm2 save

# Copy build sang thư mục web server
cp -r ./.next/* /var/www/frontend/

# Cấp quyền cho www-data
chown -R www-data:www-data /var/www/frontend

# Reload Nginx
systemctl reload nginx
