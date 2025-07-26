#ifndef AES_ENCRYPTOR_H
#define AES_ENCRYPTOR_H

#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

unsigned char *aes_encrypt_pbkdf2(
    const unsigned char *plaintext,
    size_t plaintext_len,
    const char *password,
    size_t *ciphertext_len_out
);

unsigned char *aes_decrypt_pbkdf2(
    const unsigned char *ciphertext_with_salt,
    size_t ciphertext_len,
    const char *password,
    size_t *plaintext_len_out
);

#ifdef __cplusplus
}
#endif

#endif
