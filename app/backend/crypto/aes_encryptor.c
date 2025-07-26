#include "aes_encryptor.h"
#include <openssl/evp.h>
#include <openssl/rand.h>
#include <string.h>
#include <stdlib.h>

#define SALT_SIZE 8
#define KEY_SIZE 32
#define IV_SIZE 16
#define PBKDF2_ITERATIONS 100000

unsigned char *aes_encrypt_pbkdf2(
    const unsigned char *plaintext,
    size_t plaintext_len,
    const char *password,
    size_t *ciphertext_len_out
) {
    unsigned char salt[SALT_SIZE];
    if (!RAND_bytes(salt, SALT_SIZE)) return NULL;

    unsigned char key_iv[KEY_SIZE + IV_SIZE];
    if (!PKCS5_PBKDF2_HMAC(password, strlen(password), salt, SALT_SIZE,
                           PBKDF2_ITERATIONS, EVP_sha256(),
                           KEY_SIZE + IV_SIZE, key_iv)) {
        return NULL;
    }

    unsigned char *key = key_iv;
    unsigned char *iv = key_iv + KEY_SIZE;

    EVP_CIPHER_CTX *ctx = EVP_CIPHER_CTX_new();
    if (!ctx) return NULL;

    int outlen1, outlen2;
    size_t max_len = plaintext_len + EVP_MAX_BLOCK_LENGTH;
    unsigned char *ciphertext = malloc(SALT_SIZE + max_len);
    if (!ciphertext) return NULL;

    memcpy(ciphertext, salt, SALT_SIZE);  // Save salt at the beginning

    if (!EVP_EncryptInit_ex(ctx, EVP_aes_256_cbc(), NULL, key, iv)) goto err;

    if (!EVP_EncryptUpdate(ctx, ciphertext + SALT_SIZE, &outlen1, plaintext, plaintext_len)) goto err;

    if (!EVP_EncryptFinal_ex(ctx, ciphertext + SALT_SIZE + outlen1, &outlen2)) goto err;

    *ciphertext_len_out = SALT_SIZE + outlen1 + outlen2;

    EVP_CIPHER_CTX_free(ctx);
    return ciphertext;

err:
    EVP_CIPHER_CTX_free(ctx);
    free(ciphertext);
    return NULL;
}

unsigned char *aes_decrypt_pbkdf2(
    const unsigned char *ciphertext_with_salt,
    size_t ciphertext_len,
    const char *password,
    size_t *plaintext_len_out
) {
    if (ciphertext_len < SALT_SIZE) return NULL;

    const unsigned char *salt = ciphertext_with_salt;
    const unsigned char *ciphertext = ciphertext_with_salt + SALT_SIZE;
    size_t enc_len = ciphertext_len - SALT_SIZE;

    unsigned char key_iv[KEY_SIZE + IV_SIZE];
    if (!PKCS5_PBKDF2_HMAC(password, strlen(password), salt, SALT_SIZE,
                           PBKDF2_ITERATIONS, EVP_sha256(),
                           KEY_SIZE + IV_SIZE, key_iv)) {
        return NULL;
    }

    unsigned char *key = key_iv;
    unsigned char *iv = key_iv + KEY_SIZE;

    EVP_CIPHER_CTX *ctx = EVP_CIPHER_CTX_new();
    if (!ctx) return NULL;

    int outlen1, outlen2;
    unsigned char *plaintext = malloc(enc_len);
    if (!plaintext) return NULL;

    if (!EVP_DecryptInit_ex(ctx, EVP_aes_256_cbc(), NULL, key, iv)) goto err;

    if (!EVP_DecryptUpdate(ctx, plaintext, &outlen1, ciphertext, enc_len)) goto err;

    if (!EVP_DecryptFinal_ex(ctx, plaintext + outlen1, &outlen2)) goto err;

    *plaintext_len_out = outlen1 + outlen2;

    EVP_CIPHER_CTX_free(ctx);
    return plaintext;

err:
    EVP_CIPHER_CTX_free(ctx);
    free(plaintext);
    return NULL;
}
