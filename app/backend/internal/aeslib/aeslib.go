package aeslib

/*
#cgo CFLAGS: -I../../crypto
#cgo LDFLAGS: -L../../crypto -laes_encryptor -lcrypto -Wl,-rpath=../../crypto
#include "aes_encryptor.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func EncryptPBKDF2(data []byte, password string) ([]byte, error) {
	var outLen C.size_t
	cPassword := C.CString(password)
	defer C.free(unsafe.Pointer(cPassword))

	cipher := C.aes_encrypt_pbkdf2(
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.size_t(len(data)),
		cPassword,
		&outLen,
	)

	if cipher == nil {
		return nil, fmt.Errorf("encryption failed")
	}
	defer C.free(unsafe.Pointer(cipher))

	return C.GoBytes(unsafe.Pointer(cipher), C.int(outLen)), nil
}

func DecryptPBKDF2(data []byte, password string) ([]byte, error) {
	var outLen C.size_t
	cPassword := C.CString(password)
	defer C.free(unsafe.Pointer(cPassword))

	plain := C.aes_decrypt_pbkdf2(
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.size_t(len(data)),
		cPassword,
		&outLen,
	)

	if plain == nil {
		return nil, fmt.Errorf("decryption failed")
	}
	defer C.free(unsafe.Pointer(plain))

	return C.GoBytes(unsafe.Pointer(plain), C.int(outLen)), nil
}
