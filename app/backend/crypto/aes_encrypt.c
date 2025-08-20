// aeslib.c
#include "aeslib.h"
#include <string.h>
#include <stdlib.h>
#include <openssl/evp.h>
#include <openssl/rand.h>

// ---------- helpers ----------
static int gcm_encrypt(const uint8_t *key, const uint8_t *iv, size_t iv_len,
                       const uint8_t *pt, size_t pt_len,
                       uint8_t **out_ct, size_t *out_ct_len,
                       uint8_t tag[GCM_TAG_LEN]) {
    int rc = -1;
    EVP_CIPHER_CTX *ctx = NULL;
    const EVP_CIPHER *cipher = EVP_aes_256_gcm();
    int len = 0, ct_len = 0;

    *out_ct = NULL; *out_ct_len = 0;

    ctx = EVP_CIPHER_CTX_new();
    if (!ctx) goto done;

    if (EVP_EncryptInit_ex(ctx, cipher, NULL, NULL, NULL) != 1) goto done;
    if (EVP_CIPHER_CTX_ctrl(ctx, EVP_CTRL_GCM_SET_IVLEN, (int)iv_len, NULL) != 1) goto done;
    if (EVP_EncryptInit_ex(ctx, NULL, NULL, key, iv) != 1) goto done;

    uint8_t *ct = (uint8_t*)malloc(pt_len);
    if (!ct) goto done;

    if (EVP_EncryptUpdate(ctx, ct, &len, pt, (int)pt_len) != 1) { free(ct); goto done; }
    ct_len = len;

    if (EVP_EncryptFinal_ex(ctx, ct + ct_len, &len) != 1) { free(ct); goto done; }
    ct_len += len;

    if (EVP_CIPHER_CTX_ctrl(ctx, EVP_CTRL_GCM_GET_TAG, GCM_TAG_LEN, tag) != 1) { free(ct); goto done; }

    *out_ct = ct;
    *out_ct_len = (size_t)ct_len;
    rc = 0;

done:
    if (ctx) EVP_CIPHER_CTX_free(ctx);
    return rc;
}

static int gcm_decrypt(const uint8_t *key, const uint8_t *iv, size_t iv_len,
                       const uint8_t *ct, size_t ct_len,
                       const uint8_t tag[GCM_TAG_LEN],
                       uint8_t **out_pt, size_t *out_pt_len) {
    int rc = -1;
    EVP_CIPHER_CTX *ctx = NULL;
    const EVP_CIPHER *cipher = EVP_aes_256_gcm();
    int len = 0, pt_len = 0;

    *out_pt = NULL; *out_pt_len = 0;

    ctx = EVP_CIPHER_CTX_new();
    if (!ctx) goto done;

    if (EVP_DecryptInit_ex(ctx, cipher, NULL, NULL, NULL) != 1) goto done;
    if (EVP_CIPHER_CTX_ctrl(ctx, EVP_CTRL_GCM_SET_IVLEN, (int)iv_len, NULL) != 1) goto done;
    if (EVP_DecryptInit_ex(ctx, NULL, NULL, key, iv) != 1) goto done;

    uint8_t *pt = (uint8_t*)malloc(ct_len);
    if (!pt) goto done;

    if (EVP_DecryptUpdate(ctx, pt, &len, ct, (int)ct_len) != 1) { free(pt); goto done; }
    pt_len = len;

    if (EVP_CIPHER_CTX_ctrl(ctx, EVP_CTRL_GCM_SET_TAG, GCM_TAG_LEN, (void*)tag) != 1) { free(pt); goto done; }

    if (EVP_DecryptFinal_ex(ctx, pt + pt_len, &len) != 1) { free(pt); goto done; }
    pt_len += len;

    *out_pt = pt;
    *out_pt_len = (size_t)pt_len;
    rc = 0;

done:
    if (ctx) EVP_CIPHER_CTX_free(ctx);
    return rc;
}

static int pbkdf2_derive(const char *password,
                         const uint8_t *salt, size_t salt_len,
                         uint8_t out_key[AES_KEY_LEN]) {
    if (PKCS5_PBKDF2_HMAC(password, (int)strlen(password),
                          salt, (int)salt_len,
                          PBKDF2_ITERS,
                          EVP_sha256(),
                          AES_KEY_LEN, out_key) != 1) {
        return -1;
    }
    return 0;
}

static void secure_free(uint8_t *p, size_t n) {
    if (!p) return;
    OPENSSL_cleanse(p, n);
    free(p);
}

// ---------- public API ----------
int aes_ediploma_encrypt(
    const uint8_t *plaintext, size_t plaintext_len,
    const char *password,
    encrypt_result_t *out
) {
    if (!plaintext || !password || !out) return -1;
    memset(out, 0, sizeof(*out));

    int rc = -1;
    uint8_t K2[AES_KEY_LEN], IV2[GCM_IV_LEN], tag2[GCM_TAG_LEN];
    uint8_t K1[AES_KEY_LEN], IV1[GCM_IV_LEN], tag_wrap[GCM_TAG_LEN];
    uint8_t K0[AES_KEY_LEN], IV0[GCM_IV_LEN], tag0[GCM_TAG_LEN];
    uint8_t salt[PBKDF2_SALT_LEN];

    uint8_t *file_ct = NULL; size_t file_ct_len = 0;
    uint8_t *wrap_ct = NULL; size_t wrap_ct_len = 0;
    uint8_t *k1iv1_ct = NULL; size_t k1iv1_ct_len = 0;

    // 1) Sinh K2, IV2
    if (RAND_bytes(K2, sizeof K2) != 1) goto done;
    if (RAND_bytes(IV2, sizeof IV2) != 1) goto done;

    // 2) Mã hoá file = AES-GCM(K2, IV2, plaintext)
    if (gcm_encrypt(K2, IV2, sizeof IV2, plaintext, plaintext_len,
                    &file_ct, &file_ct_len, tag2) != 0) goto done;

    // 3) Sinh K1, IV1
    if (RAND_bytes(K1, sizeof K1) != 1) goto done;
    if (RAND_bytes(IV1, sizeof IV1) != 1) goto done;

    // 4) Bọc (K2||IV2) bằng AES-GCM(K1, IV1)
    uint8_t key2_iv2[ AES_KEY_LEN + GCM_IV_LEN ];
    memcpy(key2_iv2, K2, AES_KEY_LEN);
    memcpy(key2_iv2 + AES_KEY_LEN, IV2, GCM_IV_LEN);

    if (gcm_encrypt(K1, IV1, sizeof IV1,
                    key2_iv2, sizeof key2_iv2,
                    &wrap_ct, &wrap_ct_len, tag_wrap) != 0) goto done;

    // 6) K0 = PBKDF2(password, salt)
    if (RAND_bytes(salt, sizeof salt) != 1) goto done;
    if (pbkdf2_derive(password, salt, sizeof salt, K0) != 0) goto done;

    // 7) Mã hoá (K1||IV1) bằng AES-GCM(K0, IV0 ngẫu nhiên)
    if (RAND_bytes(IV0, sizeof IV0) != 1) goto done;

    uint8_t key1_iv1[ AES_KEY_LEN + GCM_IV_LEN ];
    memcpy(key1_iv1, K1, AES_KEY_LEN);
    memcpy(key1_iv1 + AES_KEY_LEN, IV1, GCM_IV_LEN);

    if (gcm_encrypt(K0, IV0, sizeof IV0,
                    key1_iv1, sizeof key1_iv1,
                    &k1iv1_ct, &k1iv1_ct_len, tag0) != 0) goto done;

    // 5) Tạo file_blob = file_ct || tag2 || wrap_ct || tag_wrap
    {
        size_t out_len = file_ct_len + GCM_TAG_LEN + wrap_ct_len + GCM_TAG_LEN;
        uint8_t *blob = (uint8_t*)malloc(out_len);
        if (!blob) goto done;

        size_t off = 0;
        memcpy(blob + off, file_ct, file_ct_len); off += file_ct_len;
        memcpy(blob + off, tag2, GCM_TAG_LEN);    off += GCM_TAG_LEN;
        memcpy(blob + off, wrap_ct, wrap_ct_len); off += wrap_ct_len;
        memcpy(blob + off, tag_wrap, GCM_TAG_LEN);

        out->file_blob.data = blob;
        out->file_blob.len  = out_len;
    }

    // 8) DB blob = salt || IV0 || tag0 || k1iv1_ct
    {
        size_t out_len = PBKDF2_SALT_LEN + GCM_IV_LEN + GCM_TAG_LEN + k1iv1_ct_len;
        uint8_t *blob = (uint8_t*)malloc(out_len);
        if (!blob) goto done;

        size_t off = 0;
        memcpy(blob + off, salt, PBKDF2_SALT_LEN);         off += PBKDF2_SALT_LEN;
        memcpy(blob + off, IV0, GCM_IV_LEN);               off += GCM_IV_LEN;
        memcpy(blob + off, tag0, GCM_TAG_LEN);             off += GCM_TAG_LEN;
        memcpy(blob + off, k1iv1_ct, k1iv1_ct_len);        // off += ...

        out->k1iv1_db_blob.data = blob;
        out->k1iv1_db_blob.len  = out_len;
    }

    rc = 0;

done:
    secure_free(file_ct, file_ct_len);
    secure_free(wrap_ct, wrap_ct_len);
    secure_free(k1iv1_ct, k1iv1_ct_len);
    OPENSSL_cleanse(K0, sizeof K0);
    OPENSSL_cleanse(K1, sizeof K1);
    OPENSSL_cleanse(IV1, sizeof IV1);
    OPENSSL_cleanse(K2, sizeof K2);
    OPENSSL_cleanse(IV2, sizeof IV2);

    if (rc != 0) {
        aes_free_result(out);
    }
    return rc;
}

void aes_free_result(encrypt_result_t *r) {
    if (!r) return;
    if (r->file_blob.data) { secure_free(r->file_blob.data, r->file_blob.len); r->file_blob.data = NULL; r->file_blob.len = 0; }
    if (r->k1iv1_db_blob.data) { secure_free(r->k1iv1_db_blob.data, r->k1iv1_db_blob.len); r->k1iv1_db_blob.data = NULL; r->k1iv1_db_blob.len = 0; }
}

// --------- (tuỳ chọn) decrypt để verify ----------
int aes_ediploma_decrypt(
    const uint8_t *file_blob, size_t file_blob_len,
    const char *password,
    const uint8_t *k1iv1_db_blob, size_t k1iv1_db_blob_len,
    buf_t *out_plain
) {
    if (!file_blob || file_blob_len < (GCM_TAG_LEN*2) || !password || !k1iv1_db_blob || k1iv1_db_blob_len < (PBKDF2_SALT_LEN+GCM_IV_LEN+GCM_TAG_LEN+1) || !out_plain) {
        return -1;
    }
    memset(out_plain, 0, sizeof(*out_plain));

    int rc = -1;
    uint8_t K0[AES_KEY_LEN], IV0[GCM_IV_LEN], tag0[GCM_TAG_LEN], salt[PBKDF2_SALT_LEN];
    const uint8_t *ct0 = NULL; size_t ct0_len = 0;

    // parse DB blob
    size_t off = 0;
    memcpy(salt, k1iv1_db_blob + off, PBKDF2_SALT_LEN); off += PBKDF2_SALT_LEN;
    memcpy(IV0,   k1iv1_db_blob + off, GCM_IV_LEN);     off += GCM_IV_LEN;
    memcpy(tag0,  k1iv1_db_blob + off, GCM_TAG_LEN);    off += GCM_TAG_LEN;
    ct0 = k1iv1_db_blob + off; ct0_len = k1iv1_db_blob_len - off;

    if (pbkdf2_derive(password, salt, sizeof salt, K0) != 0) return -2;

    // decrypt K1||IV1
    uint8_t *key1_iv1 = NULL; size_t key1_iv1_len = 0;
    if (gcm_decrypt(K0, IV0, sizeof IV0, ct0, ct0_len, tag0, &key1_iv1, &key1_iv1_len) != 0) goto done;
    if (key1_iv1_len != AES_KEY_LEN + GCM_IV_LEN) goto done;

    uint8_t K1[AES_KEY_LEN], IV1[GCM_IV_LEN];
    memcpy(K1, key1_iv1, AES_KEY_LEN);
    memcpy(IV1, key1_iv1 + AES_KEY_LEN, GCM_IV_LEN);
    secure_free(key1_iv1, key1_iv1_len);

    // parse file_blob: file_ct || tag2 || wrap_ct || tag_wrap
    if (file_blob_len < 2*GCM_TAG_LEN + 1) goto done;

    // Ta không biết độ dài file_ct / wrap_ct riêng, nên cần truyền kèm từ lớp Go,
    // hoặc quy ước: wrap_ct có độ dài cố định = len( K2||IV2 ) = 32+12 = 44 bytes.
    // => wrap_ct_len = 44; phần còn lại (trừ 2 tags) là file_ct.
    const size_t WRAP_CT_LEN = AES_KEY_LEN + GCM_IV_LEN; // 44
    if (file_blob_len < WRAP_CT_LEN + 2*GCM_TAG_LEN) goto done;

    size_t file_ct_len = file_blob_len - WRAP_CT_LEN - 2*GCM_TAG_LEN;
    const uint8_t *file_ct = file_blob;
    const uint8_t *tag2    = file_blob + file_ct_len;
    const uint8_t *wrap_ct = file_blob + file_ct_len + GCM_TAG_LEN;
    const uint8_t *tagw    = file_blob + file_ct_len + GCM_TAG_LEN + WRAP_CT_LEN;

    // decrypt wrap: (K2||IV2)
    uint8_t *key2_iv2 = NULL; size_t key2_iv2_len = 0;
    if (gcm_decrypt(K1, IV1, sizeof IV1, wrap_ct, WRAP_CT_LEN, tagw, &key2_iv2, &key2_iv2_len) != 0) goto done;
    if (key2_iv2_len != AES_KEY_LEN + GCM_IV_LEN) { secure_free(key2_iv2, key2_iv2_len); goto done; }

    uint8_t K2[AES_KEY_LEN], IV2[GCM_IV_LEN];
    memcpy(K2, key2_iv2, AES_KEY_LEN);
    memcpy(IV2, key2_iv2 + AES_KEY_LEN, GCM_IV_LEN);
    secure_free(key2_iv2, key2_iv2_len);

    // decrypt file
    uint8_t *pt = NULL; size_t pt_len = 0;
    if (gcm_decrypt(K2, IV2, sizeof IV2, file_ct, file_ct_len, tag2, &pt, &pt_len) != 0) goto done;

    out_plain->data = pt;
    out_plain->len  = pt_len;
    rc = 0;

done:
    OPENSSL_cleanse(K0, sizeof K0);
    return rc;
}


s.certificateRepo.UpdateCertificateByID(ctx, certificateID, bson.M{
    "$set": bson.M{
        "k1iv1_db_blob": blob, // k1iv1_db_blob chứa salt || IV0 || tag0 || k1iv1_ct
        "updated_at":    time.Now(),
    },
})
