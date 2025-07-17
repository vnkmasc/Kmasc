#!/bin/bash

SCRIPT_DIR=$(dirname "$0")
PROJECT_ROOT=$(realpath "$SCRIPT_DIR/..")
CRYPTO_DIR="$PROJECT_ROOT/crypto"

# Biên dịch file trong thư mục crypto
gcc -shared -fPIC -o "$CRYPTO_DIR/libaes_encryptor.so" "$CRYPTO_DIR/aes_encryptor.c" -lcrypto

# Kiểm tra thành công không
if [ $? -ne 0 ]; then
  echo "❌ Build thất bại ."
  exit 1
fi

# Copy vào /usr/local/lib
sudo cp "$CRYPTO_DIR/libaes_encryptor.so" /usr/local/lib/
sudo ldconfig

echo "✅ Đã build và cài đặt libaes_encryptor.so thành công!"
