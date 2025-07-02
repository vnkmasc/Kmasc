#ifndef ENCRYPT_H
#define ENCRYPT_H

#include <stdint.h>

// Prototype hàm mã hóa AES CBC
int encrypt_aes_cbc(const unsigned char* plaintext, int plaintext_len, unsigned char* ciphertext, int* ciphertext_len);

// Prototype hàm giải mã AES CBC
int decrypt_aes_cbc(const unsigned char* ciphertext, int ciphertext_len, unsigned char* plaintext, int* plaintext_len);

#endif // ENCRYPT_H 