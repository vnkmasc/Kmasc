#pragma once
#include <stdint.h>

// Define rotation functions for Linux compatibility
#define _rotl64(x, r) (((x) << (r)) | ((x) >> (64 - (r))))
#define _rotr64(x, r) (((x) >> (r)) | ((x) << (64 - (r))))

#define Round256 7 //MKV256_256
#define Round384 8 //MKV256_384
#define Round512 9 //MKV256_512
#define BlockLength 32

#define KeyLen256 256
#define KeyLen384 384
#define KeyLen512 512

#define SWAP(x) (_rotl64(x,  8) & 0x000000ff000000ff | _rotl64(x, 24) & 0x0000ff000000ff00 | _rotr64(x,  8) & 0xff000000ff000000 | _rotr64(x, 24) & 0x00ff000000ff0000)

#define GETU64(p) SWAP(*((uint64_t *)(p)))
#define PUTU64(ct, st) { *((uint64_t *)(ct)) = SWAP((st)); }


//Ham mo rong khoa
int KeyExpansion256(
	int keyLen,						//do dai khoa 256 or 384 or 512 bit
	const unsigned char *MasterKey, //Khoa chinh
	uint64_t *rKey);		//Khoa vong

//Ham Tinh khoa tuong duong phuc vu giai ma
int InvKeyExpansion256(
	int keyLen,						//do dai khoa 256 or 384 or 512 bit
	const unsigned char *MasterKey, //Khoa chinh
	uint64_t *irKey);		//Khoa vong cho giai ma

//Ham ma hoa giai ma mot khoi cho phien ban khoi 256, khoa 256
int EncryptOneBlock256(
	int keyLen,				//do dai khoa 256 or 384 or 512 bit
	uint64_t *rKey, //Khoa vong
	unsigned char *in,		//Ban ro
	unsigned char *out);	//Ban ma
int DecryptOneBlock256(
	int keyLen,				//do dai khoa 256 or 384 or 512 bit
	uint64_t *rKey, //Khoa vong giai ma
	unsigned char *in,		//Ban ma dau vao
	unsigned char *out);	//Khoi ban ro sau giai ma

