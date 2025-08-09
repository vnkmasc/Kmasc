#!/bin/bash
pm2 stop nextjs-fe
nvm use 24
cd /root/Kmasc/app/frontend || exit 1

if [ "$1" = "reinstall" ]; then
    npm install
fi

npm run build
HOST=0.0.0.0 pm2 start npm --name nextjs-fe -- start
pm2 save
cp -r ./.next/* /var/www/frontend/
chown -R www-data:www-data /var/www/frontend
systemctl reload nginx
