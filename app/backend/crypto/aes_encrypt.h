// aeslib.h
#pragma once
#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

// Độ dài mặc định
#define AES_KEY_LEN     32     // 256-bit
#define GCM_IV_LEN      12     // IV khuyến nghị cho GCM
#define GCM_TAG_LEN     16     // Tag 128-bit
#define PBKDF2_ITERS    100000 // tuỳ chỉnh
#define PBKDF2_SALT_LEN 16

typedef struct {
    uint8_t *data;
    size_t   len;
} buf_t;

typedef struct {
    // Để upload MinIO: file_ct || file_tag || wrap2_ct || wrap2_tag
    buf_t file_blob;
    // Để lưu DB: salt || iv0 || tag0 || ct0
    buf_t k1iv1_db_blob;
} encrypt_result_t;

// Encrypt theo mô tả:
// - password: mật khẩu người dùng (plaintext); module sẽ tự sinh K0 bằng PBKDF2.
// - plaintext, plaintext_len: nội dung file rõ.
//
// Trả về 0 nếu OK, !=0 nếu lỗi.
// Caller chịu trách nhiệm free kết quả bằng aes_free_result().
int aes_ediploma_encrypt(
    const uint8_t *plaintext, size_t plaintext_len,
    const char *password,
    encrypt_result_t *out);

// Giải phóng bộ nhớ kết quả
void aes_free_result(encrypt_result_t *r);

// (Tuỳ chọn) Giải mã ngược lại chỉ để test/verify
// Yêu cầu: có password (để giải blob DB ra K1,IV1) + file_blob đọc từ MinIO.
int aes_ediploma_decrypt(
    const uint8_t *file_blob, size_t file_blob_len,
    const char *password,
    const uint8_t *k1iv1_db_blob, size_t k1iv1_db_blob_len,
    buf_t *out_plain);

#ifdef __cplusplus
}
#endif
