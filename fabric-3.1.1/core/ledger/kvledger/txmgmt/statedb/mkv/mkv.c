#include "mkv.h"
#include "MKV256.h"
#include <string.h>
#include <stdlib.h>
#include <stdio.h>

#define BLOCK_SIZE 32 // 256 bit

// Padding kiểu PKCS#7
static void pkcs7_pad(const unsigned char* in, int in_len, unsigned char* out, int* out_len) {
    int pad = BLOCK_SIZE - (in_len % BLOCK_SIZE);
    memcpy(out, in, in_len);
    for (int i = 0; i < pad; ++i) out[in_len + i] = (unsigned char)pad;
    *out_len = in_len + pad;
}

static int pkcs7_unpad(unsigned char* buf, int buf_len) {
    if (buf_len == 0) return 0;
    int pad = buf[buf_len - 1];
    if (pad <= 0 || pad > BLOCK_SIZE) return -1;
    for (int i = 0; i < pad; ++i) {
        if (buf[buf_len - 1 - i] != pad) return -1;
    }
    return buf_len - pad;
}

int mkv_encrypt(const unsigned char* plaintext, int plaintext_len, unsigned char* ciphertext, int* ciphertext_len, const unsigned char* key, int key_len) {
    if (!plaintext || !ciphertext || !ciphertext_len || !key) return -1;
    int padded_len = ((plaintext_len / BLOCK_SIZE) + 1) * BLOCK_SIZE;
    unsigned char* padded = (unsigned char*)malloc(padded_len);
    int real_padded_len = 0;
    pkcs7_pad(plaintext, plaintext_len, padded, &real_padded_len);
    uint64_t rKey[80];
    int ret = KeyExpansion256(key_len, key, rKey);
    if (ret != 1) {
        free(padded);
        fprintf(stderr, "KeyExpansion256 failed: key_len=%d\n", key_len);
        return -1;
    }
    for (int i = 0; i < real_padded_len; i += BLOCK_SIZE) {
        if (EncryptOneBlock256(key_len, rKey, padded + i, ciphertext + i) != 1) { free(padded); return -1; }
    }
    *ciphertext_len = real_padded_len;
    free(padded);
    return 0;
}

int mkv_decrypt(const unsigned char* ciphertext, int ciphertext_len, unsigned char* plaintext, int* plaintext_len, const unsigned char* key, int key_len) {
    if (!ciphertext || !plaintext || !plaintext_len || !key) return -1;
    if (ciphertext_len % BLOCK_SIZE != 0) return -1;
    uint64_t irKey[80];
    int ret = InvKeyExpansion256(key_len, key, irKey);
    if (ret != 1) return -1;
    for (int i = 0; i < ciphertext_len; i += BLOCK_SIZE) {
        if (DecryptOneBlock256(key_len, irKey, (unsigned char*)ciphertext + i, plaintext + i) != 1) return -1;
    }
    int unpad_len = pkcs7_unpad(plaintext, ciphertext_len);
    if (unpad_len < 0) return -1;
    *plaintext_len = unpad_len;
    return 0;
} 