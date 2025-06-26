#include "encrypt.h"
#include <string.h>

// Dummy AES CBC encrypt: chỉ copy plaintext sang ciphertext
int encrypt_aes_cbc(const unsigned char* plaintext, int plaintext_len, unsigned char* ciphertext, int* ciphertext_len) {
    if (!plaintext || !ciphertext || !ciphertext_len) return -1;
    memcpy(ciphertext, plaintext, plaintext_len);
    *ciphertext_len = plaintext_len;
    return 0; // 0 = success
}

// Dummy AES CBC decrypt: chỉ copy ciphertext sang plaintext
int decrypt_aes_cbc(const unsigned char* ciphertext, int ciphertext_len, unsigned char* plaintext, int* plaintext_len) {
    if (!ciphertext || !plaintext || !plaintext_len) return -1;
    memcpy(plaintext, ciphertext, ciphertext_len);
    *plaintext_len = ciphertext_len;
    return 0; // 0 = success
} 