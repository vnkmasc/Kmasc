#include "aes_encryptor.h"

#include <stdlib.h>
#include <string.h>
#include <openssl/evp.h>
#include <openssl/rand.h>

void generate_aes_key(uint8_t key[AES_KEY_SIZE]) {
    RAND_bytes(key, AES_KEY_SIZE);
}

void generate_iv(uint8_t iv[AES_IV_SIZE]) {
    RAND_bytes(iv, AES_IV_SIZE);
}

uint8_t* aes_encrypt(
    const uint8_t* data,
    size_t data_len,
    const uint8_t key[AES_KEY_SIZE],
    const uint8_t iv[AES_IV_SIZE],
    size_t* output_len
) {
    EVP_CIPHER_CTX* ctx = EVP_CIPHER_CTX_new();
    if (!ctx) return NULL;

    uint8_t* ciphertext = malloc(data_len + AES_IV_SIZE);
    if (!ciphertext) {
        EVP_CIPHER_CTX_free(ctx);
        return NULL;
    }

    int len = 0;
    int total_len = 0;

    EVP_EncryptInit_ex(ctx, EVP_aes_256_cbc(), NULL, key, iv);

    EVP_EncryptUpdate(ctx, ciphertext, &len, data, (int)data_len);
    total_len = len;

    EVP_EncryptFinal_ex(ctx, ciphertext + len, &len);
    total_len += len;

    EVP_CIPHER_CTX_free(ctx);

    *output_len = total_len;
    return ciphertext;
}

uint8_t* aes_decrypt(
    const uint8_t* enc_data,
    size_t enc_len,
    const uint8_t key[AES_KEY_SIZE],
    const uint8_t iv[AES_IV_SIZE],
    size_t* output_len
) {
    EVP_CIPHER_CTX* ctx = EVP_CIPHER_CTX_new();
    if (!ctx) return NULL;

    uint8_t* plaintext = malloc(enc_len);
    if (!plaintext) {
        EVP_CIPHER_CTX_free(ctx);
        return NULL;
    }

    int len = 0;
    int total_len = 0;

    EVP_DecryptInit_ex(ctx, EVP_aes_256_cbc(), NULL, key, iv);

    EVP_DecryptUpdate(ctx, plaintext, &len, enc_data, (int)enc_len);
    total_len = len;

    if (!EVP_DecryptFinal_ex(ctx, plaintext + len, &len)) {
        free(plaintext);
        EVP_CIPHER_CTX_free(ctx);
        return NULL;
    }

    total_len += len;
    *output_len = total_len;

    EVP_CIPHER_CTX_free(ctx);
    return plaintext;
}
