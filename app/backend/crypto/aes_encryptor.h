#ifndef AES_ENCRYPTOR_H
#define AES_ENCRYPTOR_H

#include <stdint.h>
#include <stddef.h>

#define AES_KEY_SIZE 32  // 256-bit
#define AES_IV_SIZE 16   // 128-bit

#ifdef __cplusplus
extern "C" {
#endif

void generate_aes_key(uint8_t key[AES_KEY_SIZE]);
void generate_iv(uint8_t iv[AES_IV_SIZE]);


uint8_t* aes_encrypt(
    const uint8_t* data,
    size_t data_len,
    const uint8_t key[AES_KEY_SIZE],
    const uint8_t iv[AES_IV_SIZE],
    size_t* output_len
);

uint8_t* aes_decrypt(
    const uint8_t* enc_data,
    size_t enc_len,
    const uint8_t key[AES_KEY_SIZE],
    const uint8_t iv[AES_IV_SIZE],
    size_t* output_len
);

#ifdef __cplusplus
}
#endif

#endif 
