#include "encrypt.h"
#include <string.h>
#include <openssl/evp.h>
#include <openssl/err.h>
#include <openssl/rand.h>

// Key và IV cho AES (trong thực tế nên được lưu an toàn)
static unsigned char aes_key[32] = {
    0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
    0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10,
    0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
    0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10
};

static unsigned char aes_iv[16] = {
    0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0,
    0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0
};

// AES CBC encrypt với OpenSSL
int encrypt_aes_cbc(const unsigned char* plaintext, int plaintext_len, unsigned char* ciphertext, int* ciphertext_len) {
    if (!plaintext || !ciphertext || !ciphertext_len) return -1;
    
    EVP_CIPHER_CTX *ctx = EVP_CIPHER_CTX_new();
    if (!ctx) return -1;
    
    int len;
    int ciphertext_len_int = 0;
    
    // Initialize encryption
    if (EVP_EncryptInit_ex(ctx, EVP_aes_256_cbc(), NULL, aes_key, aes_iv) != 1) {
        EVP_CIPHER_CTX_free(ctx);
        return -1;
    }
    
    // Encrypt data
    if (EVP_EncryptUpdate(ctx, ciphertext, &len, plaintext, plaintext_len) != 1) {
        EVP_CIPHER_CTX_free(ctx);
        return -1;
    }
    ciphertext_len_int = len;
    
    // Finalize encryption
    if (EVP_EncryptFinal_ex(ctx, ciphertext + len, &len) != 1) {
        EVP_CIPHER_CTX_free(ctx);
        return -1;
    }
    ciphertext_len_int += len;
    
    EVP_CIPHER_CTX_free(ctx);
    *ciphertext_len = ciphertext_len_int;
    return 0;
}

// AES CBC decrypt với OpenSSL
int decrypt_aes_cbc(const unsigned char* ciphertext, int ciphertext_len, unsigned char* plaintext, int* plaintext_len) {
    if (!ciphertext || !plaintext || !plaintext_len) return -1;
    
    EVP_CIPHER_CTX *ctx = EVP_CIPHER_CTX_new();
    if (!ctx) return -1;
    
    int len;
    int plaintext_len_int = 0;
    
    // Initialize decryption
    if (EVP_DecryptInit_ex(ctx, EVP_aes_256_cbc(), NULL, aes_key, aes_iv) != 1) {
        EVP_CIPHER_CTX_free(ctx);
        return -1;
    }
    
    // Decrypt data
    if (EVP_DecryptUpdate(ctx, plaintext, &len, ciphertext, ciphertext_len) != 1) {
        EVP_CIPHER_CTX_free(ctx);
        return -1;
    }
    plaintext_len_int = len;
    
    // Finalize decryption
    if (EVP_DecryptFinal_ex(ctx, plaintext + len, &len) != 1) {
        EVP_CIPHER_CTX_free(ctx);
        return -1;
    }
    plaintext_len_int += len;
    
    EVP_CIPHER_CTX_free(ctx);
    *plaintext_len = plaintext_len_int;
    return 0;
} 