#ifndef MKV_H
#define MKV_H

#include <stdint.h>

// Mã hóa dữ liệu dài với MKV256, có padding kiểu PKCS#7
int mkv_encrypt(const unsigned char* plaintext, int plaintext_len, unsigned char* ciphertext, int* ciphertext_len, const unsigned char* key, int key_len);

// Giải mã dữ liệu dài với MKV256, có xử lý bỏ padding
int mkv_decrypt(const unsigned char* ciphertext, int ciphertext_len, unsigned char* plaintext, int* plaintext_len, const unsigned char* key, int key_len);

#endif // MKV_H 