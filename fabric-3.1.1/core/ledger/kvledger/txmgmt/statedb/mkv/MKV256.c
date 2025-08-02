#include "stdio.h"
#include <stdint.h>

typedef unsigned char u8;
//typedef unsigned int u32;
typedef uint64_t u64;  // Changed from uint64_t to uint64_t


const uint64_t CONSTANT_0 = 0x9302ee911a2ad98c;
const uint64_t CONSTANT_1 = 0xad13e7948ad8b3b2;
const uint64_t CONSTANT_2 = 0xd4da00f33f11fd88;
const uint64_t CONSTANT_3 = 0x22166bb9cd187c55;

const uint64_t CONSTANT_4 = 0x43c853a3f90f70ae;
const uint64_t CONSTANT_5 = 0x8d72ad0e7aac0c71;
const uint64_t CONSTANT_6 = 0xafad739c17a14bd6;
const uint64_t CONSTANT_7 = 0x56c9436300def1e5;

#include "PrecomputedTable256.h"
#include "MKV256.h"


void setInitialKeyState256(int keyLen, const u8 *K, u64 *Block1, u64 *Block2) {
	Block1[0] = GETU64(K);
	Block1[1] = GETU64(K + 8);
	Block1[2] = GETU64(K + 16);
	Block1[3] = GETU64(K + 24);
	if (keyLen == KeyLen256) { //Khoa 128 bit
		Block2[0] = Block1[0] ^ 0xFFFFFFFFFFFFFFFF;  //Nua con lai lay phu dinh
		Block2[1] = Block1[1] ^ 0xFFFFFFFFFFFFFFFF;  //Nua con lai lay phu dinh
		Block2[2] = Block1[2] ^ 0xFFFFFFFFFFFFFFFF;  //Nua con lai lay phu dinh
		Block2[3] = Block1[3] ^ 0xFFFFFFFFFFFFFFFF;  //Nua con lai lay phu dinh
	}
	if (keyLen == KeyLen384) {
		Block2[0] = GETU64(K + 32);
		Block2[1] = GETU64(K + 40);
		Block2[2] = Block1[2] ^ 0xFFFFFFFFFFFFFFFF;
		Block2[3] = Block1[3] ^ 0xFFFFFFFFFFFFFFFF;
	}
	if (keyLen == KeyLen512) {
		Block2[0] = GETU64(K + 32);
		Block2[1] = GETU64(K + 40);
		Block2[2] = GETU64(K + 48);
		Block2[3] = GETU64(K + 56);
	}
}

u64 F(u64 t) {
	return L0[(t >> 56) & 0xFF] ^ L1[(t >> 48) & 0xFF] ^ L2[(t >> 40) & 0xFF] ^ L3[(t >> 32) & 0xFF] ^ L4[(t >> 24) & 0xFF] ^ L5[(t >> 16) & 0xFF] ^ L6[(t >> 8) & 0xFF] ^ L7[t & 0xFF];
}
u64 F1(u64 t) {
	return L0[(t >> 40) & 0xFF] ^ L1[(t >> 32) & 0xFF] ^ L2[(t >> 24) & 0xFF] ^ L3[(t >> 16) & 0xFF] ^ L4[(t >> 8) & 0xFF] ^ L5[t & 0xFF] ^ L6[(t >> 56) & 0xFF] ^ L7[(t >> 48) & 0xFF];
}
u64 F2(u64 t) {
	return L0[(t >> 24) & 0xFF] ^ L1[(t >> 16) & 0xFF] ^ L2[(t >> 8) & 0xFF] ^ L3[t & 0xFF] ^ L4[(t >> 56) & 0xFF] ^ L5[(t >> 48) & 0xFF] ^ L6[(t >> 40) & 0xFF] ^ L7[(t >> 32) & 0xFF];
}
u64 F3(u64 t) {
	return L0[(t >> 8) & 0xFF] ^ L1[t & 0xFF] ^ L2[(t >> 56) & 0xFF] ^ L3[(t >> 48) & 0xFF] ^ L4[(t >> 40) & 0xFF] ^ L5[(t >> 32) & 0xFF] ^ L6[(t >> 24) & 0xFF] ^ L7[(t >> 16) & 0xFF];
}

u64 iF(u64 t) {
	return iL0[(t >> 56) & 0xFF] ^ iL1[(t >> 48) & 0xFF] ^ iL2[(t >> 40) & 0xFF] ^ iL3[(t >> 32) & 0xFF] ^ iL4[(t >> 24) & 0xFF] ^ iL5[(t >> 16) & 0xFF] ^ iL6[(t >> 8) & 0xFF] ^ iL7[t & 0xFF];
}

int KeyExpansion256(
	int keyLen,						//do dai khoa 256 or 384 or 512 bit
	const unsigned char* MasterKey, //Khoa chinh
	uint64_t* rKey)			//Khoa vong
{
	if (!keyLen || !MasterKey) return 0;
	//printf("\n");
	rKey[0] = GETU64(MasterKey);
	rKey[1] = GETU64(MasterKey + 8);
	rKey[2] = GETU64(MasterKey + 16);
	rKey[3] = GETU64(MasterKey + 24);

	u64 Block1[4], Block2[4];
	u64 tmp[4];
	setInitialKeyState256(keyLen, MasterKey, Block1, Block2);
	int i;
	int round = Round256;
	if (keyLen == KeyLen384) round = Round384;
	if (keyLen == KeyLen512) round = Round512;
	for (i = 0; i < round; i++) {
		//Ham update
		//Ap dung ham F len khoi 2

		//Cong hang so C1
		Block1[0] ^= CONSTANT_0;
		Block1[1] ^= CONSTANT_1;
		Block1[2] ^= CONSTANT_2;
		Block1[3] ^= (CONSTANT_3 ^ (2 * i + 1));
		//L1S
		tmp[0] = F(Block1[0]);
		tmp[1] = F1(Block1[1]);
		tmp[2] = F2(Block1[2]);
		tmp[3] = F3(Block1[3]);

		//S2
		Block1[0] = L8[(tmp[0] >> 56) & 0xFF] ^ (L8[(tmp[0] >> 48) & 0xFF] >> 8) ^ (L8[(tmp[0] >> 40) & 0xFF] >> 16) ^ (L8[(tmp[0] >> 32) & 0xFF] >> 24) ^
			(L8[(tmp[0] >> 24) & 0xFF] >> 32) ^ (L8[(tmp[0] >> 16) & 0xFF] >> 40) ^ (L8[(tmp[0] >> 8) & 0xFF] >> 48) ^ (L8[tmp[0] & 0xFF] >> 56);
		Block1[1] = L8[(tmp[1] >> 56) & 0xFF] ^ (L8[(tmp[1] >> 48) & 0xFF] >> 8) ^ (L8[(tmp[1] >> 40) & 0xFF] >> 16) ^ (L8[(tmp[1] >> 32) & 0xFF] >> 24) ^
			(L8[(tmp[1] >> 24) & 0xFF] >> 32) ^ (L8[(tmp[1] >> 16) & 0xFF] >> 40) ^ (L8[(tmp[1] >> 8) & 0xFF] >> 48) ^ (L8[tmp[1] & 0xFF] >> 56);

		Block1[2] = L8[(tmp[2] >> 56) & 0xFF] ^ (L8[(tmp[2] >> 48) & 0xFF] >> 8) ^ (L8[(tmp[2] >> 40) & 0xFF] >> 16) ^ (L8[(tmp[2] >> 32) & 0xFF] >> 24) ^
			(L8[(tmp[2] >> 24) & 0xFF] >> 32) ^ (L8[(tmp[2] >> 16) & 0xFF] >> 40) ^ (L8[(tmp[2] >> 8) & 0xFF] >> 48) ^ (L8[tmp[2] & 0xFF] >> 56);
		Block1[3] = L8[(tmp[3] >> 56) & 0xFF] ^ (L8[(tmp[3] >> 48) & 0xFF] >> 8) ^ (L8[(tmp[3] >> 40) & 0xFF] >> 16) ^ (L8[(tmp[3] >> 32) & 0xFF] >> 24) ^
			(L8[(tmp[3] >> 24) & 0xFF] >> 32) ^ (L8[(tmp[3] >> 16) & 0xFF] >> 40) ^ (L8[(tmp[3] >> 8) & 0xFF] >> 48) ^ (L8[tmp[3] & 0xFF] >> 56);
		//L2
		tmp[0] = Block1[1] ^ Block1[2] ^ Block1[3];
		tmp[1] = Block1[0] ^ Block1[2] ^ Block1[3];
		tmp[2] = Block1[0] ^ Block1[1] ^ Block1[3];
		tmp[3] = Block1[0] ^ Block1[1] ^ Block1[2];


		//L1S
		Block1[0] = F(tmp[0]);
		Block1[1] = F1(tmp[1]);
		Block1[2] = F2(tmp[2]);
		Block1[3] = F3(tmp[3]);

		//S2
		tmp[0] = L8[(Block1[0] >> 56) & 0xFF] ^ (L8[(Block1[0] >> 48) & 0xFF] >> 8) ^ (L8[(Block1[0] >> 40) & 0xFF] >> 16) ^ (L8[(Block1[0] >> 32) & 0xFF] >> 24) ^
			(L8[(Block1[0] >> 24) & 0xFF] >> 32) ^ (L8[(Block1[0] >> 16) & 0xFF] >> 40) ^ (L8[(Block1[0] >> 8) & 0xFF] >> 48) ^ (L8[Block1[0] & 0xFF] >> 56);
		tmp[1] = L8[(Block1[1] >> 56) & 0xFF] ^ (L8[(Block1[1] >> 48) & 0xFF] >> 8) ^ (L8[(Block1[1] >> 40) & 0xFF] >> 16) ^ (L8[(Block1[1] >> 32) & 0xFF] >> 24) ^
			(L8[(Block1[1] >> 24) & 0xFF] >> 32) ^ (L8[(Block1[1] >> 16) & 0xFF] >> 40) ^ (L8[(Block1[1] >> 8) & 0xFF] >> 48) ^ (L8[Block1[1] & 0xFF] >> 56);

		tmp[2] = L8[(Block1[2] >> 56) & 0xFF] ^ (L8[(Block1[2] >> 48) & 0xFF] >> 8) ^ (L8[(Block1[2] >> 40) & 0xFF] >> 16) ^ (L8[(Block1[2] >> 32) & 0xFF] >> 24) ^
			(L8[(Block1[2] >> 24) & 0xFF] >> 32) ^ (L8[(Block1[2] >> 16) & 0xFF] >> 40) ^ (L8[(Block1[2] >> 8) & 0xFF] >> 48) ^ (L8[Block1[2] & 0xFF] >> 56);
		tmp[3] = L8[(Block1[3] >> 56) & 0xFF] ^ (L8[(Block1[3] >> 48) & 0xFF] >> 8) ^ (L8[(Block1[3] >> 40) & 0xFF] >> 16) ^ (L8[(Block1[3] >> 32) & 0xFF] >> 24) ^
			(L8[(Block1[3] >> 24) & 0xFF] >> 32) ^ (L8[(Block1[3] >> 16) & 0xFF] >> 40) ^ (L8[(Block1[3] >> 8) & 0xFF] >> 48) ^ (L8[Block1[3] & 0xFF] >> 56);
		//L2
		Block1[0] = tmp[1] ^ tmp[2] ^ tmp[3];
		Block1[1] = tmp[0] ^ tmp[2] ^ tmp[3];
		Block1[2] = tmp[0] ^ tmp[1] ^ tmp[3];
		Block1[3] = tmp[0] ^ tmp[1] ^ tmp[2];

		
		//Cong hang so C2
		Block2[0] ^= CONSTANT_4;
		Block2[1] ^= CONSTANT_5;
		Block2[2] ^= CONSTANT_6;
		Block2[3] ^= (CONSTANT_7 ^ (2 * i + 2));

		//L1S
		tmp[0] = F(Block2[0]);
		tmp[1] = F1(Block2[1]);
		tmp[2] = F2(Block2[2]);
		tmp[3] = F3(Block2[3]);

		//S2
		Block2[0] = L8[(tmp[0] >> 56) & 0xFF] ^ (L8[(tmp[0] >> 48) & 0xFF] >> 8) ^ (L8[(tmp[0] >> 40) & 0xFF] >> 16) ^ (L8[(tmp[0] >> 32) & 0xFF] >> 24) ^
			(L8[(tmp[0] >> 24) & 0xFF] >> 32) ^ (L8[(tmp[0] >> 16) & 0xFF] >> 40) ^ (L8[(tmp[0] >> 8) & 0xFF] >> 48) ^ (L8[tmp[0] & 0xFF] >> 56);
		Block2[1] = L8[(tmp[1] >> 56) & 0xFF] ^ (L8[(tmp[1] >> 48) & 0xFF] >> 8) ^ (L8[(tmp[1] >> 40) & 0xFF] >> 16) ^ (L8[(tmp[1] >> 32) & 0xFF] >> 24) ^
			(L8[(tmp[1] >> 24) & 0xFF] >> 32) ^ (L8[(tmp[1] >> 16) & 0xFF] >> 40) ^ (L8[(tmp[1] >> 8) & 0xFF] >> 48) ^ (L8[tmp[1] & 0xFF] >> 56);

		Block2[2] = L8[(tmp[2] >> 56) & 0xFF] ^ (L8[(tmp[2] >> 48) & 0xFF] >> 8) ^ (L8[(tmp[2] >> 40) & 0xFF] >> 16) ^ (L8[(tmp[2] >> 32) & 0xFF] >> 24) ^
			(L8[(tmp[2] >> 24) & 0xFF] >> 32) ^ (L8[(tmp[2] >> 16) & 0xFF] >> 40) ^ (L8[(tmp[2] >> 8) & 0xFF] >> 48) ^ (L8[tmp[2] & 0xFF] >> 56);
		Block2[3] = L8[(tmp[3] >> 56) & 0xFF] ^ (L8[(tmp[3] >> 48) & 0xFF] >> 8) ^ (L8[(tmp[3] >> 40) & 0xFF] >> 16) ^ (L8[(tmp[3] >> 32) & 0xFF] >> 24) ^
			(L8[(tmp[3] >> 24) & 0xFF] >> 32) ^ (L8[(tmp[3] >> 16) & 0xFF] >> 40) ^ (L8[(tmp[3] >> 8) & 0xFF] >> 48) ^ (L8[tmp[3] & 0xFF] >> 56);
		//L2
		tmp[0] = Block2[1] ^ Block2[2] ^ Block2[3];
		tmp[1] = Block2[0] ^ Block2[2] ^ Block2[3];
		tmp[2] = Block2[0] ^ Block2[1] ^ Block2[3];
		tmp[3] = Block2[0] ^ Block2[1] ^ Block2[2];

		//L1S
		Block2[0] = F(tmp[0]);
		Block2[1] = F1(tmp[1]);
		Block2[2] = F2(tmp[2]);
		Block2[3] = F3(tmp[3]);

		//S2
		tmp[0] = L8[(Block2[0] >> 56) & 0xFF] ^ (L8[(Block2[0] >> 48) & 0xFF] >> 8) ^ (L8[(Block2[0] >> 40) & 0xFF] >> 16) ^ (L8[(Block2[0] >> 32) & 0xFF] >> 24) ^
			(L8[(Block2[0] >> 24) & 0xFF] >> 32) ^ (L8[(Block2[0] >> 16) & 0xFF] >> 40) ^ (L8[(Block2[0] >> 8) & 0xFF] >> 48) ^ (L8[Block2[0] & 0xFF] >> 56);
		tmp[1] = L8[(Block2[1] >> 56) & 0xFF] ^ (L8[(Block2[1] >> 48) & 0xFF] >> 8) ^ (L8[(Block2[1] >> 40) & 0xFF] >> 16) ^ (L8[(Block2[1] >> 32) & 0xFF] >> 24) ^
			(L8[(Block2[1] >> 24) & 0xFF] >> 32) ^ (L8[(Block2[1] >> 16) & 0xFF] >> 40) ^ (L8[(Block2[1] >> 8) & 0xFF] >> 48) ^ (L8[Block2[1] & 0xFF] >> 56);

		tmp[2] = L8[(Block2[2] >> 56) & 0xFF] ^ (L8[(Block2[2] >> 48) & 0xFF] >> 8) ^ (L8[(Block2[2] >> 40) & 0xFF] >> 16) ^ (L8[(Block2[2] >> 32) & 0xFF] >> 24) ^
			(L8[(Block2[2] >> 24) & 0xFF] >> 32) ^ (L8[(Block2[2] >> 16) & 0xFF] >> 40) ^ (L8[(Block2[2] >> 8) & 0xFF] >> 48) ^ (L8[Block2[2] & 0xFF] >> 56);
		tmp[3] = L8[(Block2[3] >> 56) & 0xFF] ^ (L8[(Block2[3] >> 48) & 0xFF] >> 8) ^ (L8[(Block2[3] >> 40) & 0xFF] >> 16) ^ (L8[(Block2[3] >> 32) & 0xFF] >> 24) ^
			(L8[(Block2[3] >> 24) & 0xFF] >> 32) ^ (L8[(Block2[3] >> 16) & 0xFF] >> 40) ^ (L8[(Block2[3] >> 8) & 0xFF] >> 48) ^ (L8[Block2[3] & 0xFF] >> 56);
		//L2
		Block2[0] = tmp[1] ^ tmp[2] ^ tmp[3];
		Block2[1] = tmp[0] ^ tmp[2] ^ tmp[3];
		Block2[2] = tmp[0] ^ tmp[1] ^ tmp[3];
		Block2[3] = tmp[0] ^ tmp[1] ^ tmp[2];

		//Tinh 2 nhanh theo luoc do V

		tmp[0] = Block2[0];
		tmp[1] = Block2[1];
		tmp[2] = Block2[2];
		tmp[3] = Block2[3];

		rKey[4 * (2 * i + 1) + 4] = Block2[0] = tmp[0] ^ Block1[0];
		rKey[4 * (2 * i + 1) + 5] = Block2[1] = tmp[1] ^ Block1[1];
		rKey[4 * (2 * i + 1) + 6] = Block2[2] = tmp[2] ^ Block1[2];
		rKey[4 * (2 * i + 1) + 7] = Block2[3] = tmp[3] ^ Block1[3];
		//
		rKey[4 * (2 * i + 1)] = Block1[0] = tmp[0];
		rKey[4 * (2 * i + 1) + 1] = Block1[1] = tmp[1];
		rKey[4 * (2 * i + 1) + 2] = Block1[2] = tmp[2];
		rKey[4 * (2 * i + 1) + 3] = Block1[3] = tmp[3];
	}
	return 1;
}

u64 InvKey(u64 in)
{
	u64 t = L8[(in >> 56) & 0xFF] ^ (L8[(in >> 48) & 0xFF] >> 8) ^ (L8[(in >> 40) & 0xFF] >> 16) ^ (L8[(in >> 32) & 0xFF] >> 24) ^
		(L8[(in >> 24) & 0xFF] >> 32) ^ (L8[(in >> 16) & 0xFF] >> 40) ^ (L8[(in >> 8) & 0xFF] >> 48) ^ (L8[in & 0xFF] >> 56);
	return iF(t);
}

int InvKeyExpansion256(
	int keyLen,						//do dai khoa 256 or 384 or 512 bit
	const unsigned char* MasterKey, //Khoa chinh
	uint64_t* irKey)		//Khoa vong cho giai ma
{
	if (!keyLen || !MasterKey) return 0;
	int i;

	KeyExpansion256(keyLen, MasterKey, irKey);

	int round = Round256;
	if (keyLen == KeyLen384) round = Round384;
	if (keyLen == KeyLen512) round = Round512;

	for (i = 1; i < 2 * round; i = i + 2) {
		irKey[4 * i] = InvKey(irKey[4 * i]);
		irKey[4 * i + 1] = InvKey(irKey[4 * i + 1]);
		irKey[4 * i + 2] = InvKey(irKey[4 * i + 2]);
		irKey[4 * i + 3] = InvKey(irKey[4 * i + 3]);
	}
	return 1;
}

int EncryptOneBlock256(
	int keyLen,				//do dai khoa 256 or 384 or 512 bit
	uint64_t* rKey, //Khoa vong
	unsigned char* in,		//Ban ro
	unsigned char* out)		//Ban ma
{
	u64 s0, s1, s2, s3, t0, t1, t2, t3;
	if (!keyLen || !in || !out || !rKey) return 0;

	s0 = GETU64(in) ^ rKey[0];
	s1 = GETU64(in + 8) ^ rKey[1];
	s2 = GETU64(in + 16) ^ rKey[2];
	s3 = GETU64(in + 24) ^ rKey[3];

	//printState32(rKey[0], rKey[1], rKey[2], rKey[3]);
	//printState32(t0, t1, t2, t3);
	//printf("Vong 1\n");
	//r1
	//L1S
	s0 = F(s0) ^ rKey[4];
	s1 = F1(s1) ^ rKey[5];
	s2 = F2(s2) ^ rKey[6];
	s3 = F3(s3) ^ rKey[7];
	//s0 = L0[(t0 >> 56) & 0xFF] ^ L1[(t0 >> 48) & 0xFF] ^ L2[(t0 >> 40) & 0xFF] ^ L3[(t0 >> 32) & 0xFF] ^ L8[(t0 >> 24) & 0xFF] ^ L5[(t0 >> 16) & 0xFF] ^ L6[(t0 >> 8) & 0xFF] ^ L7[t0 & 0xFF] ^ rKey[4];
	//s1 = L0[(t1 >> 56) & 0xFF] ^ L1[(t1 >> 48) & 0xFF] ^ L2[(t1 >> 40) & 0xFF] ^ L3[(t1 >> 32) & 0xFF] ^ L8[(t1 >> 24) & 0xFF] ^ L5[(t1 >> 16) & 0xFF] ^ L6[(t1 >> 8) & 0xFF] ^ L7[t1 & 0xFF] ^ rKey[4];
	//s2 = L0[(t2 >> 56) & 0xFF] ^ L1[(t2 >> 48) & 0xFF] ^ L2[(t2 >> 40) & 0xFF] ^ L3[(t2 >> 32) & 0xFF] ^ L8[(t2 >> 24) & 0xFF] ^ L5[(t2 >> 16) & 0xFF] ^ L6[(t2 >> 8) & 0xFF] ^ L7[t2 & 0xFF] ^ rKey[4];
	//s3 = L0[(t3 >> 56) & 0xFF] ^ L1[(t3 >> 48) & 0xFF] ^ L2[(t3 >> 40) & 0xFF] ^ L3[(t3 >> 32) & 0xFF] ^ L8[(t3 >> 24) & 0xFF] ^ L5[(t3 >> 16) & 0xFF] ^ L6[(t3 >> 8) & 0xFF] ^ L7[t3 & 0xFF] ^ rKey[4];
	//printState32(rKey[4], rKey[5], rKey[6], rKey[7]);
	//printState32(s0, s1, s2, s3);

	//S
	t0 = L8[(s0 >> 56) & 0xFF] ^ (L8[(s0 >> 48) & 0xFF] >> 8) ^ (L8[(s0 >> 40) & 0xFF] >> 16) ^ (L8[(s0 >> 32) & 0xFF] >> 24) ^ (L8[(s0 >> 24) & 0xFF] >> 32) ^ (L8[(s0 >> 16) & 0xFF] >> 40) ^ (L8[(s0 >> 8) & 0xFF] >> 48) ^ (L8[s0 & 0xFF] >> 56);
	t1 = L8[(s1 >> 56) & 0xFF] ^ (L8[(s1 >> 48) & 0xFF] >> 8) ^ (L8[(s1 >> 40) & 0xFF] >> 16) ^ (L8[(s1 >> 32) & 0xFF] >> 24) ^ (L8[(s1 >> 24) & 0xFF] >> 32) ^ (L8[(s1 >> 16) & 0xFF] >> 40) ^ (L8[(s1 >> 8) & 0xFF] >> 48) ^ (L8[s1 & 0xFF] >> 56);
	t2 = L8[(s2 >> 56) & 0xFF] ^ (L8[(s2 >> 48) & 0xFF] >> 8) ^ (L8[(s2 >> 40) & 0xFF] >> 16) ^ (L8[(s2 >> 32) & 0xFF] >> 24) ^ (L8[(s2 >> 24) & 0xFF] >> 32) ^ (L8[(s2 >> 16) & 0xFF] >> 40) ^ (L8[(s2 >> 8) & 0xFF] >> 48) ^ (L8[s2 & 0xFF] >> 56);
	t3 = L8[(s3 >> 56) & 0xFF] ^ (L8[(s3 >> 48) & 0xFF] >> 8) ^ (L8[(s3 >> 40) & 0xFF] >> 16) ^ (L8[(s3 >> 32) & 0xFF] >> 24) ^ (L8[(s3 >> 24) & 0xFF] >> 32) ^ (L8[(s3 >> 16) & 0xFF] >> 40) ^ (L8[(s3 >> 8) & 0xFF] >> 48) ^ (L8[s3 & 0xFF] >> 56);
	//printState32(t0, t1, t2, t3);
	//L2
	s0 = t1 ^ t2 ^ t3 ^ rKey[8];
	s1 = t0 ^ t2 ^ t3 ^ rKey[9];
	s2 = t0 ^ t1 ^ t3 ^ rKey[10];
	s3 = t0 ^ t1 ^ t2 ^ rKey[11];
	//printState32(rKey[8], rKey[9], rKey[10], rKey[11]);
	//printState32(s0, s1, s2, s3);
	//printf("Vong 2\n");
	//r2
	//L1S
	s0 = F(s0) ^ rKey[12];
	s1 = F1(s1) ^ rKey[13];
	s2 = F2(s2) ^ rKey[14];
	s3 = F3(s3) ^ rKey[15];
	//printState32(rKey[4], rKey[5], rKey[6], rKey[7]);
	//printState32(s0, s1, s2, s3);

	//S
	t0 = L8[(s0 >> 56) & 0xFF] ^ (L8[(s0 >> 48) & 0xFF] >> 8) ^ (L8[(s0 >> 40) & 0xFF] >> 16) ^ (L8[(s0 >> 32) & 0xFF] >> 24) ^ (L8[(s0 >> 24) & 0xFF] >> 32) ^ (L8[(s0 >> 16) & 0xFF] >> 40) ^ (L8[(s0 >> 8) & 0xFF] >> 48) ^ (L8[s0 & 0xFF] >> 56);
	t1 = L8[(s1 >> 56) & 0xFF] ^ (L8[(s1 >> 48) & 0xFF] >> 8) ^ (L8[(s1 >> 40) & 0xFF] >> 16) ^ (L8[(s1 >> 32) & 0xFF] >> 24) ^ (L8[(s1 >> 24) & 0xFF] >> 32) ^ (L8[(s1 >> 16) & 0xFF] >> 40) ^ (L8[(s1 >> 8) & 0xFF] >> 48) ^ (L8[s1 & 0xFF] >> 56);
	t2 = L8[(s2 >> 56) & 0xFF] ^ (L8[(s2 >> 48) & 0xFF] >> 8) ^ (L8[(s2 >> 40) & 0xFF] >> 16) ^ (L8[(s2 >> 32) & 0xFF] >> 24) ^ (L8[(s2 >> 24) & 0xFF] >> 32) ^ (L8[(s2 >> 16) & 0xFF] >> 40) ^ (L8[(s2 >> 8) & 0xFF] >> 48) ^ (L8[s2 & 0xFF] >> 56);
	t3 = L8[(s3 >> 56) & 0xFF] ^ (L8[(s3 >> 48) & 0xFF] >> 8) ^ (L8[(s3 >> 40) & 0xFF] >> 16) ^ (L8[(s3 >> 32) & 0xFF] >> 24) ^ (L8[(s3 >> 24) & 0xFF] >> 32) ^ (L8[(s3 >> 16) & 0xFF] >> 40) ^ (L8[(s3 >> 8) & 0xFF] >> 48) ^ (L8[s3 & 0xFF] >> 56);
	//printState32(t0, t1, t2, t3);
	//L2
	s0 = t1 ^ t2 ^ t3 ^ rKey[16];
	s1 = t0 ^ t2 ^ t3 ^ rKey[17];
	s2 = t0 ^ t1 ^ t3 ^ rKey[18];
	s3 = t0 ^ t1 ^ t2 ^ rKey[19];
	//printState32(rKey[16], rKey[17], rKey[18], rKey[19]);
	//printState32(t0, t1, t2, t3);

	//r3
	//L1S
	//L1S
	s0 = F(s0) ^ rKey[20];
	s1 = F1(s1) ^ rKey[21];
	s2 = F2(s2) ^ rKey[22];
	s3 = F3(s3) ^ rKey[23];
	//printState32(rKey[4], rKey[5], rKey[6], rKey[7]);
	//printState32(s0, s1, s2, s3);

	//S
	t0 = L8[(s0 >> 56) & 0xFF] ^ (L8[(s0 >> 48) & 0xFF] >> 8) ^ (L8[(s0 >> 40) & 0xFF] >> 16) ^ (L8[(s0 >> 32) & 0xFF] >> 24) ^ (L8[(s0 >> 24) & 0xFF] >> 32) ^ (L8[(s0 >> 16) & 0xFF] >> 40) ^ (L8[(s0 >> 8) & 0xFF] >> 48) ^ (L8[s0 & 0xFF] >> 56);
	t1 = L8[(s1 >> 56) & 0xFF] ^ (L8[(s1 >> 48) & 0xFF] >> 8) ^ (L8[(s1 >> 40) & 0xFF] >> 16) ^ (L8[(s1 >> 32) & 0xFF] >> 24) ^ (L8[(s1 >> 24) & 0xFF] >> 32) ^ (L8[(s1 >> 16) & 0xFF] >> 40) ^ (L8[(s1 >> 8) & 0xFF] >> 48) ^ (L8[s1 & 0xFF] >> 56);
	t2 = L8[(s2 >> 56) & 0xFF] ^ (L8[(s2 >> 48) & 0xFF] >> 8) ^ (L8[(s2 >> 40) & 0xFF] >> 16) ^ (L8[(s2 >> 32) & 0xFF] >> 24) ^ (L8[(s2 >> 24) & 0xFF] >> 32) ^ (L8[(s2 >> 16) & 0xFF] >> 40) ^ (L8[(s2 >> 8) & 0xFF] >> 48) ^ (L8[s2 & 0xFF] >> 56);
	t3 = L8[(s3 >> 56) & 0xFF] ^ (L8[(s3 >> 48) & 0xFF] >> 8) ^ (L8[(s3 >> 40) & 0xFF] >> 16) ^ (L8[(s3 >> 32) & 0xFF] >> 24) ^ (L8[(s3 >> 24) & 0xFF] >> 32) ^ (L8[(s3 >> 16) & 0xFF] >> 40) ^ (L8[(s3 >> 8) & 0xFF] >> 48) ^ (L8[s3 & 0xFF] >> 56);
	//printState32(t0, t1, t2, t3);
	//L2
	s0 = t1 ^ t2 ^ t3 ^ rKey[24];
	s1 = t0 ^ t2 ^ t3 ^ rKey[25];
	s2 = t0 ^ t1 ^ t3 ^ rKey[26];
	s3 = t0 ^ t1 ^ t2 ^ rKey[27];

	//r4
	//L1S
	//L1S
	s0 = F(s0) ^ rKey[28];
	s1 = F1(s1) ^ rKey[29];
	s2 = F2(s2) ^ rKey[30];
	s3 = F3(s3) ^ rKey[31];
	//printState32(rKey[4], rKey[5], rKey[6], rKey[7]);
	//printState32(s0, s1, s2, s3);

	//S
	t0 = L8[(s0 >> 56) & 0xFF] ^ (L8[(s0 >> 48) & 0xFF] >> 8) ^ (L8[(s0 >> 40) & 0xFF] >> 16) ^ (L8[(s0 >> 32) & 0xFF] >> 24) ^ (L8[(s0 >> 24) & 0xFF] >> 32) ^ (L8[(s0 >> 16) & 0xFF] >> 40) ^ (L8[(s0 >> 8) & 0xFF] >> 48) ^ (L8[s0 & 0xFF] >> 56);
	t1 = L8[(s1 >> 56) & 0xFF] ^ (L8[(s1 >> 48) & 0xFF] >> 8) ^ (L8[(s1 >> 40) & 0xFF] >> 16) ^ (L8[(s1 >> 32) & 0xFF] >> 24) ^ (L8[(s1 >> 24) & 0xFF] >> 32) ^ (L8[(s1 >> 16) & 0xFF] >> 40) ^ (L8[(s1 >> 8) & 0xFF] >> 48) ^ (L8[s1 & 0xFF] >> 56);
	t2 = L8[(s2 >> 56) & 0xFF] ^ (L8[(s2 >> 48) & 0xFF] >> 8) ^ (L8[(s2 >> 40) & 0xFF] >> 16) ^ (L8[(s2 >> 32) & 0xFF] >> 24) ^ (L8[(s2 >> 24) & 0xFF] >> 32) ^ (L8[(s2 >> 16) & 0xFF] >> 40) ^ (L8[(s2 >> 8) & 0xFF] >> 48) ^ (L8[s2 & 0xFF] >> 56);
	t3 = L8[(s3 >> 56) & 0xFF] ^ (L8[(s3 >> 48) & 0xFF] >> 8) ^ (L8[(s3 >> 40) & 0xFF] >> 16) ^ (L8[(s3 >> 32) & 0xFF] >> 24) ^ (L8[(s3 >> 24) & 0xFF] >> 32) ^ (L8[(s3 >> 16) & 0xFF] >> 40) ^ (L8[(s3 >> 8) & 0xFF] >> 48) ^ (L8[s3 & 0xFF] >> 56);
	//printState32(t0, t1, t2, t3);
	//L2
	s0 = t1 ^ t2 ^ t3 ^ rKey[32];
	s1 = t0 ^ t2 ^ t3 ^ rKey[33];
	s2 = t0 ^ t1 ^ t3 ^ rKey[34];
	s3 = t0 ^ t1 ^ t2 ^ rKey[35];

	//r5
	//L1S
	//L1S
	s0 = F(s0) ^ rKey[36];
	s1 = F1(s1) ^ rKey[37];
	s2 = F2(s2) ^ rKey[38];
	s3 = F3(s3) ^ rKey[39];
	//printState32(rKey[4], rKey[5], rKey[6], rKey[7]);
	//printState32(s0, s1, s2, s3);

	//S
	t0 = L8[(s0 >> 56) & 0xFF] ^ (L8[(s0 >> 48) & 0xFF] >> 8) ^ (L8[(s0 >> 40) & 0xFF] >> 16) ^ (L8[(s0 >> 32) & 0xFF] >> 24) ^ (L8[(s0 >> 24) & 0xFF] >> 32) ^ (L8[(s0 >> 16) & 0xFF] >> 40) ^ (L8[(s0 >> 8) & 0xFF] >> 48) ^ (L8[s0 & 0xFF] >> 56);
	t1 = L8[(s1 >> 56) & 0xFF] ^ (L8[(s1 >> 48) & 0xFF] >> 8) ^ (L8[(s1 >> 40) & 0xFF] >> 16) ^ (L8[(s1 >> 32) & 0xFF] >> 24) ^ (L8[(s1 >> 24) & 0xFF] >> 32) ^ (L8[(s1 >> 16) & 0xFF] >> 40) ^ (L8[(s1 >> 8) & 0xFF] >> 48) ^ (L8[s1 & 0xFF] >> 56);
	t2 = L8[(s2 >> 56) & 0xFF] ^ (L8[(s2 >> 48) & 0xFF] >> 8) ^ (L8[(s2 >> 40) & 0xFF] >> 16) ^ (L8[(s2 >> 32) & 0xFF] >> 24) ^ (L8[(s2 >> 24) & 0xFF] >> 32) ^ (L8[(s2 >> 16) & 0xFF] >> 40) ^ (L8[(s2 >> 8) & 0xFF] >> 48) ^ (L8[s2 & 0xFF] >> 56);
	t3 = L8[(s3 >> 56) & 0xFF] ^ (L8[(s3 >> 48) & 0xFF] >> 8) ^ (L8[(s3 >> 40) & 0xFF] >> 16) ^ (L8[(s3 >> 32) & 0xFF] >> 24) ^ (L8[(s3 >> 24) & 0xFF] >> 32) ^ (L8[(s3 >> 16) & 0xFF] >> 40) ^ (L8[(s3 >> 8) & 0xFF] >> 48) ^ (L8[s3 & 0xFF] >> 56);
	//printState32(t0, t1, t2, t3);
	//L2
	s0 = t1 ^ t2 ^ t3 ^ rKey[40];
	s1 = t0 ^ t2 ^ t3 ^ rKey[41];
	s2 = t0 ^ t1 ^ t3 ^ rKey[42];
	s3 = t0 ^ t1 ^ t2 ^ rKey[43];

	//r6
	//L1S
	//L1S
	s0 = F(s0) ^ rKey[44];
	s1 = F1(s1) ^ rKey[45];
	s2 = F2(s2) ^ rKey[46];
	s3 = F3(s3) ^ rKey[47];
	//printState32(rKey[4], rKey[5], rKey[6], rKey[7]);
	//printState32(s0, s1, s2, s3);

	//S
	t0 = L8[(s0 >> 56) & 0xFF] ^ (L8[(s0 >> 48) & 0xFF] >> 8) ^ (L8[(s0 >> 40) & 0xFF] >> 16) ^ (L8[(s0 >> 32) & 0xFF] >> 24) ^ (L8[(s0 >> 24) & 0xFF] >> 32) ^ (L8[(s0 >> 16) & 0xFF] >> 40) ^ (L8[(s0 >> 8) & 0xFF] >> 48) ^ (L8[s0 & 0xFF] >> 56);
	t1 = L8[(s1 >> 56) & 0xFF] ^ (L8[(s1 >> 48) & 0xFF] >> 8) ^ (L8[(s1 >> 40) & 0xFF] >> 16) ^ (L8[(s1 >> 32) & 0xFF] >> 24) ^ (L8[(s1 >> 24) & 0xFF] >> 32) ^ (L8[(s1 >> 16) & 0xFF] >> 40) ^ (L8[(s1 >> 8) & 0xFF] >> 48) ^ (L8[s1 & 0xFF] >> 56);
	t2 = L8[(s2 >> 56) & 0xFF] ^ (L8[(s2 >> 48) & 0xFF] >> 8) ^ (L8[(s2 >> 40) & 0xFF] >> 16) ^ (L8[(s2 >> 32) & 0xFF] >> 24) ^ (L8[(s2 >> 24) & 0xFF] >> 32) ^ (L8[(s2 >> 16) & 0xFF] >> 40) ^ (L8[(s2 >> 8) & 0xFF] >> 48) ^ (L8[s2 & 0xFF] >> 56);
	t3 = L8[(s3 >> 56) & 0xFF] ^ (L8[(s3 >> 48) & 0xFF] >> 8) ^ (L8[(s3 >> 40) & 0xFF] >> 16) ^ (L8[(s3 >> 32) & 0xFF] >> 24) ^ (L8[(s3 >> 24) & 0xFF] >> 32) ^ (L8[(s3 >> 16) & 0xFF] >> 40) ^ (L8[(s3 >> 8) & 0xFF] >> 48) ^ (L8[s3 & 0xFF] >> 56);
	//printState32(t0, t1, t2, t3);
	//L2
	s0 = t1 ^ t2 ^ t3 ^ rKey[48];
	s1 = t0 ^ t2 ^ t3 ^ rKey[49];
	s2 = t0 ^ t1 ^ t3 ^ rKey[50];
	s3 = t0 ^ t1 ^ t2 ^ rKey[51];
	//printf("\cong khoa cuoi\n");
	//printState32(t0, t1, t2, t3);

	//r7
	//L1S
	s0 = F(s0) ^ rKey[52];
	s1 = F1(s1) ^ rKey[53];
	s2 = F2(s2) ^ rKey[54];
	s3 = F3(s3) ^ rKey[55];
	//printState32(rKey[4], rKey[5], rKey[6], rKey[7]);
	//printState32(s0, s1, s2, s3);

	//S
	t0 = L8[(s0 >> 56) & 0xFF] ^ (L8[(s0 >> 48) & 0xFF] >> 8) ^ (L8[(s0 >> 40) & 0xFF] >> 16) ^ (L8[(s0 >> 32) & 0xFF] >> 24) ^ (L8[(s0 >> 24) & 0xFF] >> 32) ^ (L8[(s0 >> 16) & 0xFF] >> 40) ^ (L8[(s0 >> 8) & 0xFF] >> 48) ^ (L8[s0 & 0xFF] >> 56);
	t1 = L8[(s1 >> 56) & 0xFF] ^ (L8[(s1 >> 48) & 0xFF] >> 8) ^ (L8[(s1 >> 40) & 0xFF] >> 16) ^ (L8[(s1 >> 32) & 0xFF] >> 24) ^ (L8[(s1 >> 24) & 0xFF] >> 32) ^ (L8[(s1 >> 16) & 0xFF] >> 40) ^ (L8[(s1 >> 8) & 0xFF] >> 48) ^ (L8[s1 & 0xFF] >> 56);
	t2 = L8[(s2 >> 56) & 0xFF] ^ (L8[(s2 >> 48) & 0xFF] >> 8) ^ (L8[(s2 >> 40) & 0xFF] >> 16) ^ (L8[(s2 >> 32) & 0xFF] >> 24) ^ (L8[(s2 >> 24) & 0xFF] >> 32) ^ (L8[(s2 >> 16) & 0xFF] >> 40) ^ (L8[(s2 >> 8) & 0xFF] >> 48) ^ (L8[s2 & 0xFF] >> 56);
	t3 = L8[(s3 >> 56) & 0xFF] ^ (L8[(s3 >> 48) & 0xFF] >> 8) ^ (L8[(s3 >> 40) & 0xFF] >> 16) ^ (L8[(s3 >> 32) & 0xFF] >> 24) ^ (L8[(s3 >> 24) & 0xFF] >> 32) ^ (L8[(s3 >> 16) & 0xFF] >> 40) ^ (L8[(s3 >> 8) & 0xFF] >> 48) ^ (L8[s3 & 0xFF] >> 56);
	//printState32(t0, t1, t2, t3);
	//L2
	s0 = t1 ^ t2 ^ t3 ^ rKey[56];
	s1 = t0 ^ t2 ^ t3 ^ rKey[57];
	s2 = t0 ^ t1 ^ t3 ^ rKey[58];
	s3 = t0 ^ t1 ^ t2 ^ rKey[59];

	if (keyLen >= KeyLen384) {
		//r8
		//L1S
		s0 = F(s0) ^ rKey[60];
		s1 = F1(s1) ^ rKey[61];
		s2 = F2(s2) ^ rKey[62];
		s3 = F3(s3) ^ rKey[63];
		//printState32(rKey[4], rKey[5], rKey[6], rKey[7]);
		//printState32(s0, s1, s2, s3);

		//S
		t0 = L8[(s0 >> 56) & 0xFF] ^ (L8[(s0 >> 48) & 0xFF] >> 8) ^ (L8[(s0 >> 40) & 0xFF] >> 16) ^ (L8[(s0 >> 32) & 0xFF] >> 24) ^ (L8[(s0 >> 24) & 0xFF] >> 32) ^ (L8[(s0 >> 16) & 0xFF] >> 40) ^ (L8[(s0 >> 8) & 0xFF] >> 48) ^ (L8[s0 & 0xFF] >> 56);
		t1 = L8[(s1 >> 56) & 0xFF] ^ (L8[(s1 >> 48) & 0xFF] >> 8) ^ (L8[(s1 >> 40) & 0xFF] >> 16) ^ (L8[(s1 >> 32) & 0xFF] >> 24) ^ (L8[(s1 >> 24) & 0xFF] >> 32) ^ (L8[(s1 >> 16) & 0xFF] >> 40) ^ (L8[(s1 >> 8) & 0xFF] >> 48) ^ (L8[s1 & 0xFF] >> 56);
		t2 = L8[(s2 >> 56) & 0xFF] ^ (L8[(s2 >> 48) & 0xFF] >> 8) ^ (L8[(s2 >> 40) & 0xFF] >> 16) ^ (L8[(s2 >> 32) & 0xFF] >> 24) ^ (L8[(s2 >> 24) & 0xFF] >> 32) ^ (L8[(s2 >> 16) & 0xFF] >> 40) ^ (L8[(s2 >> 8) & 0xFF] >> 48) ^ (L8[s2 & 0xFF] >> 56);
		t3 = L8[(s3 >> 56) & 0xFF] ^ (L8[(s3 >> 48) & 0xFF] >> 8) ^ (L8[(s3 >> 40) & 0xFF] >> 16) ^ (L8[(s3 >> 32) & 0xFF] >> 24) ^ (L8[(s3 >> 24) & 0xFF] >> 32) ^ (L8[(s3 >> 16) & 0xFF] >> 40) ^ (L8[(s3 >> 8) & 0xFF] >> 48) ^ (L8[s3 & 0xFF] >> 56);
		//printState32(t0, t1, t2, t3);
		//L2
		s0 = t1 ^ t2 ^ t3 ^ rKey[64];
		s1 = t0 ^ t2 ^ t3 ^ rKey[65];
		s2 = t0 ^ t1 ^ t3 ^ rKey[66];
		s3 = t0 ^ t1 ^ t2 ^ rKey[67];
		
		if (keyLen == KeyLen512) {
			//r8
			//L1S
			s0 = F(s0) ^ rKey[68];
			s1 = F1(s1) ^ rKey[69];
			s2 = F2(s2) ^ rKey[70];
			s3 = F3(s3) ^ rKey[71];
			//printState32(rKey[4], rKey[5], rKey[6], rKey[7]);
			//printState32(s0, s1, s2, s3);

			//S
			t0 = L8[(s0 >> 56) & 0xFF] ^ (L8[(s0 >> 48) & 0xFF] >> 8) ^ (L8[(s0 >> 40) & 0xFF] >> 16) ^ (L8[(s0 >> 32) & 0xFF] >> 24) ^ (L8[(s0 >> 24) & 0xFF] >> 32) ^ (L8[(s0 >> 16) & 0xFF] >> 40) ^ (L8[(s0 >> 8) & 0xFF] >> 48) ^ (L8[s0 & 0xFF] >> 56);
			t1 = L8[(s1 >> 56) & 0xFF] ^ (L8[(s1 >> 48) & 0xFF] >> 8) ^ (L8[(s1 >> 40) & 0xFF] >> 16) ^ (L8[(s1 >> 32) & 0xFF] >> 24) ^ (L8[(s1 >> 24) & 0xFF] >> 32) ^ (L8[(s1 >> 16) & 0xFF] >> 40) ^ (L8[(s1 >> 8) & 0xFF] >> 48) ^ (L8[s1 & 0xFF] >> 56);
			t2 = L8[(s2 >> 56) & 0xFF] ^ (L8[(s2 >> 48) & 0xFF] >> 8) ^ (L8[(s2 >> 40) & 0xFF] >> 16) ^ (L8[(s2 >> 32) & 0xFF] >> 24) ^ (L8[(s2 >> 24) & 0xFF] >> 32) ^ (L8[(s2 >> 16) & 0xFF] >> 40) ^ (L8[(s2 >> 8) & 0xFF] >> 48) ^ (L8[s2 & 0xFF] >> 56);
			t3 = L8[(s3 >> 56) & 0xFF] ^ (L8[(s3 >> 48) & 0xFF] >> 8) ^ (L8[(s3 >> 40) & 0xFF] >> 16) ^ (L8[(s3 >> 32) & 0xFF] >> 24) ^ (L8[(s3 >> 24) & 0xFF] >> 32) ^ (L8[(s3 >> 16) & 0xFF] >> 40) ^ (L8[(s3 >> 8) & 0xFF] >> 48) ^ (L8[s3 & 0xFF] >> 56);
			//printState32(t0, t1, t2, t3);
			//L2
			s0 = t1 ^ t2 ^ t3 ^ rKey[72];
			s1 = t0 ^ t2 ^ t3 ^ rKey[73];
			s2 = t0 ^ t1 ^ t3 ^ rKey[74];
			s3 = t0 ^ t1 ^ t2 ^ rKey[75];
		}
	}
	PUTU64(out, s0);
	PUTU64(out + 8, s1);
	PUTU64(out + 16, s2);
	PUTU64(out + 24, s3);

	return 1;
}

int DecryptOneBlock256(
	int keyLen,				//do dai khoa 256 or 384 or 512 bit
	uint64_t* rKey, //Khoa vong giai ma
	unsigned char* in,		//Ban ma dau vao
	unsigned char* out)		//Khoi ban ro sau giai ma
{
	u64 s0 = 0, s1 = 0, s2 = 0, s3 = 0, t0, t1, t2, t3;
	if (!keyLen || !in || !out || !rKey) return 0;
	//printf("\n Giai ma\n");
	//for (int i = 0; i < 16; i++) printf("%02X ", in[i]);
	//printf("\n");

	if (keyLen == KeyLen512) {
		s0 = GETU64(in) ^ rKey[72];
		s1 = GETU64(in + 8) ^ rKey[73];
		s2 = GETU64(in + 16) ^ rKey[74];
		s3 = GETU64(in + 24) ^ rKey[75];
		//printState32(rKey[48], rKey[49], rKey[50], rKey[51]);
		//printState32(t0, t1, t2, t3);
		//r9
		//L2
		t0 = s1 ^ s2 ^ s3;
		t1 = s0 ^ s2 ^ s3;
		t2 = s0 ^ s1 ^ s3;
		t3 = s0 ^ s1 ^ s2;
		//printState32(s0, s1, s2, s3);
		//iL1S
		t0 = iF(t0) ^ rKey[68];
		t1 = iF(t1) ^ rKey[69];
		t2 = iF(t2) ^ rKey[70];
		t3 = iF(t3) ^ rKey[71];

		//printState32(t0, t1, t2, t3);
		//S
		s0 = iL8[(t0 >> 56) & 0xFF] ^ (iL8[(t0 >> 48) & 0xFF] >> 8) ^ (iL8[(t0 >> 40) & 0xFF] >> 16) ^ (iL8[(t0 >> 32) & 0xFF] >> 24) ^ (iL8[(t0 >> 24) & 0xFF] >> 32) ^ (iL8[(t0 >> 16) & 0xFF] >> 40) ^ (iL8[(t0 >>  8) & 0xFF] >> 48) ^ (iL8[ t0        & 0xFF] >> 56) ^ rKey[64];
		s1 = iL8[(t1 >>  8) & 0xFF] ^ (iL8[(t1      ) & 0xFF] >> 8) ^ (iL8[(t1 >> 56) & 0xFF] >> 16) ^ (iL8[(t1 >> 48) & 0xFF] >> 24) ^ (iL8[(t1 >> 40) & 0xFF] >> 32) ^ (iL8[(t1 >> 32) & 0xFF] >> 40) ^ (iL8[(t1 >> 24) & 0xFF] >> 48) ^ (iL8[(t1 >> 16) & 0xFF] >> 56) ^ rKey[65];
		s2 = iL8[(t2 >> 24) & 0xFF] ^ (iL8[(t2 >> 16) & 0xFF] >> 8) ^ (iL8[(t2 >>  8) & 0xFF] >> 16) ^ (iL8[(t2      ) & 0xFF] >> 24) ^ (iL8[(t2 >> 56) & 0xFF] >> 32) ^ (iL8[(t2 >> 48) & 0xFF] >> 40) ^ (iL8[(t2 >> 40) & 0xFF] >> 48) ^ (iL8[(t2 >> 32) & 0xFF] >> 56) ^ rKey[66];
		s3 = iL8[(t3 >> 40) & 0xFF] ^ (iL8[(t3 >> 32) & 0xFF] >> 8) ^ (iL8[(t3 >> 24) & 0xFF] >> 16) ^ (iL8[(t3 >> 16) & 0xFF] >> 24) ^ (iL8[(t3 >>  8) & 0xFF] >> 32) ^ (iL8[(t3      ) & 0xFF] >> 40) ^ (iL8[(t3 >> 56) & 0xFF] >> 48) ^ (iL8[(t3 >> 48) & 0xFF] >> 56) ^ rKey[67];

		//r8
		//L2
		t0 = s1 ^ s2 ^ s3;
		t1 = s0 ^ s2 ^ s3;
		t2 = s0 ^ s1 ^ s3;
		t3 = s0 ^ s1 ^ s2;
		//printState32(s0, s1, s2, s3);
		//iL1S
		t0 = iF(t0) ^ rKey[60];
		t1 = iF(t1) ^ rKey[61];
		t2 = iF(t2) ^ rKey[62];
		t3 = iF(t3) ^ rKey[63];

		//printState32(t0, t1, t2, t3);
		//S
		s0 = iL8[(t0 >> 56) & 0xFF] ^ (iL8[(t0 >> 48) & 0xFF] >> 8) ^ (iL8[(t0 >> 40) & 0xFF] >> 16) ^ (iL8[(t0 >> 32) & 0xFF] >> 24) ^ (iL8[(t0 >> 24) & 0xFF] >> 32) ^ (iL8[(t0 >> 16) & 0xFF] >> 40) ^ (iL8[(t0 >> 8) & 0xFF] >> 48) ^ (iL8[t0 & 0xFF] >> 56) ^ rKey[56];
		s1 = iL8[(t1 >> 8) & 0xFF] ^ (iL8[(t1) & 0xFF] >> 8) ^ (iL8[(t1 >> 56) & 0xFF] >> 16) ^ (iL8[(t1 >> 48) & 0xFF] >> 24) ^ (iL8[(t1 >> 40) & 0xFF] >> 32) ^ (iL8[(t1 >> 32) & 0xFF] >> 40) ^ (iL8[(t1 >> 24) & 0xFF] >> 48) ^ (iL8[(t1 >> 16) & 0xFF] >> 56) ^ rKey[57];
		s2 = iL8[(t2 >> 24) & 0xFF] ^ (iL8[(t2 >> 16) & 0xFF] >> 8) ^ (iL8[(t2 >> 8) & 0xFF] >> 16) ^ (iL8[(t2) & 0xFF] >> 24) ^ (iL8[(t2 >> 56) & 0xFF] >> 32) ^ (iL8[(t2 >> 48) & 0xFF] >> 40) ^ (iL8[(t2 >> 40) & 0xFF] >> 48) ^ (iL8[(t2 >> 32) & 0xFF] >> 56) ^ rKey[58];
		s3 = iL8[(t3 >> 40) & 0xFF] ^ (iL8[(t3 >> 32) & 0xFF] >> 8) ^ (iL8[(t3 >> 24) & 0xFF] >> 16) ^ (iL8[(t3 >> 16) & 0xFF] >> 24) ^ (iL8[(t3 >> 8) & 0xFF] >> 32) ^ (iL8[(t3) & 0xFF] >> 40) ^ (iL8[(t3 >> 56) & 0xFF] >> 48) ^ (iL8[(t3 >> 48) & 0xFF] >> 56) ^ rKey[59];

	}
	if (keyLen == KeyLen384) {
		s0 = GETU64(in) ^ rKey[64];
		s1 = GETU64(in + 8) ^ rKey[65];
		s2 = GETU64(in + 16) ^ rKey[66];
		s3 = GETU64(in + 24) ^ rKey[67];
		//printState32(rKey[48], rKey[49], rKey[50], rKey[51]);
		//printState32(t0, t1, t2, t3);
		//r8
		//L2
		t0 = s1 ^ s2 ^ s3;
		t1 = s0 ^ s2 ^ s3;
		t2 = s0 ^ s1 ^ s3;
		t3 = s0 ^ s1 ^ s2;
		//printState32(s0, s1, s2, s3);
		//iL1S
		t0 = iF(t0) ^ rKey[60];
		t1 = iF(t1) ^ rKey[61];
		t2 = iF(t2) ^ rKey[62];
		t3 = iF(t3) ^ rKey[63];

		//printState32(t0, t1, t2, t3);
		//S
		s0 = iL8[(t0 >> 56) & 0xFF] ^ (iL8[(t0 >> 48) & 0xFF] >> 8) ^ (iL8[(t0 >> 40) & 0xFF] >> 16) ^ (iL8[(t0 >> 32) & 0xFF] >> 24) ^ (iL8[(t0 >> 24) & 0xFF] >> 32) ^ (iL8[(t0 >> 16) & 0xFF] >> 40) ^ (iL8[(t0 >> 8) & 0xFF] >> 48) ^ (iL8[t0 & 0xFF] >> 56) ^ rKey[56];
		s1 = iL8[(t1 >> 8) & 0xFF] ^ (iL8[(t1) & 0xFF] >> 8) ^ (iL8[(t1 >> 56) & 0xFF] >> 16) ^ (iL8[(t1 >> 48) & 0xFF] >> 24) ^ (iL8[(t1 >> 40) & 0xFF] >> 32) ^ (iL8[(t1 >> 32) & 0xFF] >> 40) ^ (iL8[(t1 >> 24) & 0xFF] >> 48) ^ (iL8[(t1 >> 16) & 0xFF] >> 56) ^ rKey[57];
		s2 = iL8[(t2 >> 24) & 0xFF] ^ (iL8[(t2 >> 16) & 0xFF] >> 8) ^ (iL8[(t2 >> 8) & 0xFF] >> 16) ^ (iL8[(t2) & 0xFF] >> 24) ^ (iL8[(t2 >> 56) & 0xFF] >> 32) ^ (iL8[(t2 >> 48) & 0xFF] >> 40) ^ (iL8[(t2 >> 40) & 0xFF] >> 48) ^ (iL8[(t2 >> 32) & 0xFF] >> 56) ^ rKey[58];
		s3 = iL8[(t3 >> 40) & 0xFF] ^ (iL8[(t3 >> 32) & 0xFF] >> 8) ^ (iL8[(t3 >> 24) & 0xFF] >> 16) ^ (iL8[(t3 >> 16) & 0xFF] >> 24) ^ (iL8[(t3 >> 8) & 0xFF] >> 32) ^ (iL8[(t3) & 0xFF] >> 40) ^ (iL8[(t3 >> 56) & 0xFF] >> 48) ^ (iL8[(t3 >> 48) & 0xFF] >> 56) ^ rKey[59];

	}

	if (keyLen == KeyLen256) {
		s0 = GETU64(in) ^ rKey[56];
		s1 = GETU64(in + 8) ^ rKey[57];
		s2 = GETU64(in + 16) ^ rKey[58];
		s3 = GETU64(in + 24) ^ rKey[59];
	}
	//printState32(rKey[48], rKey[49], rKey[50], rKey[51]);
	//printState32(t0, t1, t2, t3);

	//r7
	//L2
	t0 = s1 ^ s2 ^ s3;
	t1 = s0 ^ s2 ^ s3;
	t2 = s0 ^ s1 ^ s3;
	t3 = s0 ^ s1 ^ s2;
	//printState32(s0, s1, s2, s3);
	//iL1S
	t0 = iF(t0) ^ rKey[52];
	t1 = iF(t1) ^ rKey[53];
	t2 = iF(t2) ^ rKey[54];
	t3 = iF(t3) ^ rKey[55];

	//printState32(t0, t1, t2, t3);
	//S
	s0 = iL8[(t0 >> 56) & 0xFF] ^ (iL8[(t0 >> 48) & 0xFF] >> 8) ^ (iL8[(t0 >> 40) & 0xFF] >> 16) ^ (iL8[(t0 >> 32) & 0xFF] >> 24) ^ (iL8[(t0 >> 24) & 0xFF] >> 32) ^ (iL8[(t0 >> 16) & 0xFF] >> 40) ^ (iL8[(t0 >> 8) & 0xFF] >> 48) ^ (iL8[t0 & 0xFF] >> 56) ^ rKey[48];
	s1 = iL8[(t1 >> 8) & 0xFF] ^ (iL8[(t1) & 0xFF] >> 8) ^ (iL8[(t1 >> 56) & 0xFF] >> 16) ^ (iL8[(t1 >> 48) & 0xFF] >> 24) ^ (iL8[(t1 >> 40) & 0xFF] >> 32) ^ (iL8[(t1 >> 32) & 0xFF] >> 40) ^ (iL8[(t1 >> 24) & 0xFF] >> 48) ^ (iL8[(t1 >> 16) & 0xFF] >> 56) ^ rKey[49];
	s2 = iL8[(t2 >> 24) & 0xFF] ^ (iL8[(t2 >> 16) & 0xFF] >> 8) ^ (iL8[(t2 >> 8) & 0xFF] >> 16) ^ (iL8[(t2) & 0xFF] >> 24) ^ (iL8[(t2 >> 56) & 0xFF] >> 32) ^ (iL8[(t2 >> 48) & 0xFF] >> 40) ^ (iL8[(t2 >> 40) & 0xFF] >> 48) ^ (iL8[(t2 >> 32) & 0xFF] >> 56) ^ rKey[50];
	s3 = iL8[(t3 >> 40) & 0xFF] ^ (iL8[(t3 >> 32) & 0xFF] >> 8) ^ (iL8[(t3 >> 24) & 0xFF] >> 16) ^ (iL8[(t3 >> 16) & 0xFF] >> 24) ^ (iL8[(t3 >> 8) & 0xFF] >> 32) ^ (iL8[(t3) & 0xFF] >> 40) ^ (iL8[(t3 >> 56) & 0xFF] >> 48) ^ (iL8[(t3 >> 48) & 0xFF] >> 56) ^ rKey[51];


	//r6
	//L2
	t0 = s1 ^ s2 ^ s3;
	t1 = s0 ^ s2 ^ s3;
	t2 = s0 ^ s1 ^ s3;
	t3 = s0 ^ s1 ^ s2;
	//printState32(s0, s1, s2, s3);
	//iL1S
	t0 = iF(t0) ^ rKey[44];
	t1 = iF(t1) ^ rKey[45];
	t2 = iF(t2) ^ rKey[46];
	t3 = iF(t3) ^ rKey[47];

	//printState32(t0, t1, t2, t3);
	//S
	s0 = iL8[(t0 >> 56) & 0xFF] ^ (iL8[(t0 >> 48) & 0xFF] >> 8) ^ (iL8[(t0 >> 40) & 0xFF] >> 16) ^ (iL8[(t0 >> 32) & 0xFF] >> 24) ^ (iL8[(t0 >> 24) & 0xFF] >> 32) ^ (iL8[(t0 >> 16) & 0xFF] >> 40) ^ (iL8[(t0 >> 8) & 0xFF] >> 48) ^ (iL8[t0 & 0xFF] >> 56) ^ rKey[40];
	s1 = iL8[(t1 >> 8) & 0xFF] ^ (iL8[(t1) & 0xFF] >> 8) ^ (iL8[(t1 >> 56) & 0xFF] >> 16) ^ (iL8[(t1 >> 48) & 0xFF] >> 24) ^ (iL8[(t1 >> 40) & 0xFF] >> 32) ^ (iL8[(t1 >> 32) & 0xFF] >> 40) ^ (iL8[(t1 >> 24) & 0xFF] >> 48) ^ (iL8[(t1 >> 16) & 0xFF] >> 56) ^ rKey[41];
	s2 = iL8[(t2 >> 24) & 0xFF] ^ (iL8[(t2 >> 16) & 0xFF] >> 8) ^ (iL8[(t2 >> 8) & 0xFF] >> 16) ^ (iL8[(t2) & 0xFF] >> 24) ^ (iL8[(t2 >> 56) & 0xFF] >> 32) ^ (iL8[(t2 >> 48) & 0xFF] >> 40) ^ (iL8[(t2 >> 40) & 0xFF] >> 48) ^ (iL8[(t2 >> 32) & 0xFF] >> 56) ^ rKey[42];
	s3 = iL8[(t3 >> 40) & 0xFF] ^ (iL8[(t3 >> 32) & 0xFF] >> 8) ^ (iL8[(t3 >> 24) & 0xFF] >> 16) ^ (iL8[(t3 >> 16) & 0xFF] >> 24) ^ (iL8[(t3 >> 8) & 0xFF] >> 32) ^ (iL8[(t3) & 0xFF] >> 40) ^ (iL8[(t3 >> 56) & 0xFF] >> 48) ^ (iL8[(t3 >> 48) & 0xFF] >> 56) ^ rKey[43];

	//r5
	//L2
	t0 = s1 ^ s2 ^ s3;
	t1 = s0 ^ s2 ^ s3;
	t2 = s0 ^ s1 ^ s3;
	t3 = s0 ^ s1 ^ s2;
	//printState32(s0, s1, s2, s3);
	//iL1S
	t0 = iF(t0) ^ rKey[36];
	t1 = iF(t1) ^ rKey[37];
	t2 = iF(t2) ^ rKey[38];
	t3 = iF(t3) ^ rKey[39];

	//printState32(t0, t1, t2, t3);
	//S
	s0 = iL8[(t0 >> 56) & 0xFF] ^ (iL8[(t0 >> 48) & 0xFF] >> 8) ^ (iL8[(t0 >> 40) & 0xFF] >> 16) ^ (iL8[(t0 >> 32) & 0xFF] >> 24) ^ (iL8[(t0 >> 24) & 0xFF] >> 32) ^ (iL8[(t0 >> 16) & 0xFF] >> 40) ^ (iL8[(t0 >> 8) & 0xFF] >> 48) ^ (iL8[t0 & 0xFF] >> 56) ^ rKey[32];
	s1 = iL8[(t1 >> 8) & 0xFF] ^ (iL8[(t1) & 0xFF] >> 8) ^ (iL8[(t1 >> 56) & 0xFF] >> 16) ^ (iL8[(t1 >> 48) & 0xFF] >> 24) ^ (iL8[(t1 >> 40) & 0xFF] >> 32) ^ (iL8[(t1 >> 32) & 0xFF] >> 40) ^ (iL8[(t1 >> 24) & 0xFF] >> 48) ^ (iL8[(t1 >> 16) & 0xFF] >> 56) ^ rKey[33];
	s2 = iL8[(t2 >> 24) & 0xFF] ^ (iL8[(t2 >> 16) & 0xFF] >> 8) ^ (iL8[(t2 >> 8) & 0xFF] >> 16) ^ (iL8[(t2) & 0xFF] >> 24) ^ (iL8[(t2 >> 56) & 0xFF] >> 32) ^ (iL8[(t2 >> 48) & 0xFF] >> 40) ^ (iL8[(t2 >> 40) & 0xFF] >> 48) ^ (iL8[(t2 >> 32) & 0xFF] >> 56) ^ rKey[34];
	s3 = iL8[(t3 >> 40) & 0xFF] ^ (iL8[(t3 >> 32) & 0xFF] >> 8) ^ (iL8[(t3 >> 24) & 0xFF] >> 16) ^ (iL8[(t3 >> 16) & 0xFF] >> 24) ^ (iL8[(t3 >> 8) & 0xFF] >> 32) ^ (iL8[(t3) & 0xFF] >> 40) ^ (iL8[(t3 >> 56) & 0xFF] >> 48) ^ (iL8[(t3 >> 48) & 0xFF] >> 56) ^ rKey[35];


	//r4
	//L2
	t0 = s1 ^ s2 ^ s3;
	t1 = s0 ^ s2 ^ s3;
	t2 = s0 ^ s1 ^ s3;
	t3 = s0 ^ s1 ^ s2;
	//printState32(s0, s1, s2, s3);
	//iL1S
	t0 = iF(t0) ^ rKey[28];
	t1 = iF(t1) ^ rKey[29];
	t2 = iF(t2) ^ rKey[30];
	t3 = iF(t3) ^ rKey[31];

	//printState32(t0, t1, t2, t3);
	//S
	s0 = iL8[(t0 >> 56) & 0xFF] ^ (iL8[(t0 >> 48) & 0xFF] >> 8) ^ (iL8[(t0 >> 40) & 0xFF] >> 16) ^ (iL8[(t0 >> 32) & 0xFF] >> 24) ^ (iL8[(t0 >> 24) & 0xFF] >> 32) ^ (iL8[(t0 >> 16) & 0xFF] >> 40) ^ (iL8[(t0 >> 8) & 0xFF] >> 48) ^ (iL8[t0 & 0xFF] >> 56) ^ rKey[24];
	s1 = iL8[(t1 >> 8) & 0xFF] ^ (iL8[(t1) & 0xFF] >> 8) ^ (iL8[(t1 >> 56) & 0xFF] >> 16) ^ (iL8[(t1 >> 48) & 0xFF] >> 24) ^ (iL8[(t1 >> 40) & 0xFF] >> 32) ^ (iL8[(t1 >> 32) & 0xFF] >> 40) ^ (iL8[(t1 >> 24) & 0xFF] >> 48) ^ (iL8[(t1 >> 16) & 0xFF] >> 56) ^ rKey[25];
	s2 = iL8[(t2 >> 24) & 0xFF] ^ (iL8[(t2 >> 16) & 0xFF] >> 8) ^ (iL8[(t2 >> 8) & 0xFF] >> 16) ^ (iL8[(t2) & 0xFF] >> 24) ^ (iL8[(t2 >> 56) & 0xFF] >> 32) ^ (iL8[(t2 >> 48) & 0xFF] >> 40) ^ (iL8[(t2 >> 40) & 0xFF] >> 48) ^ (iL8[(t2 >> 32) & 0xFF] >> 56) ^ rKey[26];
	s3 = iL8[(t3 >> 40) & 0xFF] ^ (iL8[(t3 >> 32) & 0xFF] >> 8) ^ (iL8[(t3 >> 24) & 0xFF] >> 16) ^ (iL8[(t3 >> 16) & 0xFF] >> 24) ^ (iL8[(t3 >> 8) & 0xFF] >> 32) ^ (iL8[(t3) & 0xFF] >> 40) ^ (iL8[(t3 >> 56) & 0xFF] >> 48) ^ (iL8[(t3 >> 48) & 0xFF] >> 56) ^ rKey[27];


	//r3
	//L2
	t0 = s1 ^ s2 ^ s3;
	t1 = s0 ^ s2 ^ s3;
	t2 = s0 ^ s1 ^ s3;
	t3 = s0 ^ s1 ^ s2;
	//printState32(s0, s1, s2, s3);
	//iL1S
	t0 = iF(t0) ^ rKey[20];
	t1 = iF(t1) ^ rKey[21];
	t2 = iF(t2) ^ rKey[22];
	t3 = iF(t3) ^ rKey[23];

	//printState32(t0, t1, t2, t3);
	//S
	s0 = iL8[(t0 >> 56) & 0xFF] ^ (iL8[(t0 >> 48) & 0xFF] >> 8) ^ (iL8[(t0 >> 40) & 0xFF] >> 16) ^ (iL8[(t0 >> 32) & 0xFF] >> 24) ^ (iL8[(t0 >> 24) & 0xFF] >> 32) ^ (iL8[(t0 >> 16) & 0xFF] >> 40) ^ (iL8[(t0 >> 8) & 0xFF] >> 48) ^ (iL8[t0 & 0xFF] >> 56) ^ rKey[16];
	s1 = iL8[(t1 >> 8) & 0xFF] ^ (iL8[(t1) & 0xFF] >> 8) ^ (iL8[(t1 >> 56) & 0xFF] >> 16) ^ (iL8[(t1 >> 48) & 0xFF] >> 24) ^ (iL8[(t1 >> 40) & 0xFF] >> 32) ^ (iL8[(t1 >> 32) & 0xFF] >> 40) ^ (iL8[(t1 >> 24) & 0xFF] >> 48) ^ (iL8[(t1 >> 16) & 0xFF] >> 56) ^ rKey[17];
	s2 = iL8[(t2 >> 24) & 0xFF] ^ (iL8[(t2 >> 16) & 0xFF] >> 8) ^ (iL8[(t2 >> 8) & 0xFF] >> 16) ^ (iL8[(t2) & 0xFF] >> 24) ^ (iL8[(t2 >> 56) & 0xFF] >> 32) ^ (iL8[(t2 >> 48) & 0xFF] >> 40) ^ (iL8[(t2 >> 40) & 0xFF] >> 48) ^ (iL8[(t2 >> 32) & 0xFF] >> 56) ^ rKey[18];
	s3 = iL8[(t3 >> 40) & 0xFF] ^ (iL8[(t3 >> 32) & 0xFF] >> 8) ^ (iL8[(t3 >> 24) & 0xFF] >> 16) ^ (iL8[(t3 >> 16) & 0xFF] >> 24) ^ (iL8[(t3 >> 8) & 0xFF] >> 32) ^ (iL8[(t3) & 0xFF] >> 40) ^ (iL8[(t3 >> 56) & 0xFF] >> 48) ^ (iL8[(t3 >> 48) & 0xFF] >> 56) ^ rKey[19];


	//r2
	//L2
	t0 = s1 ^ s2 ^ s3;
	t1 = s0 ^ s2 ^ s3;
	t2 = s0 ^ s1 ^ s3;
	t3 = s0 ^ s1 ^ s2;
	//printState32(s0, s1, s2, s3);
	//iL1S
	t0 = iF(t0) ^ rKey[12];
	t1 = iF(t1) ^ rKey[13];
	t2 = iF(t2) ^ rKey[14];
	t3 = iF(t3) ^ rKey[15];

	//printState32(t0, t1, t2, t3);
	//S
	s0 = iL8[(t0 >> 56) & 0xFF] ^ (iL8[(t0 >> 48) & 0xFF] >> 8) ^ (iL8[(t0 >> 40) & 0xFF] >> 16) ^ (iL8[(t0 >> 32) & 0xFF] >> 24) ^ (iL8[(t0 >> 24) & 0xFF] >> 32) ^ (iL8[(t0 >> 16) & 0xFF] >> 40) ^ (iL8[(t0 >> 8) & 0xFF] >> 48) ^ (iL8[t0 & 0xFF] >> 56) ^ rKey[8];
	s1 = iL8[(t1 >> 8) & 0xFF] ^ (iL8[(t1) & 0xFF] >> 8) ^ (iL8[(t1 >> 56) & 0xFF] >> 16) ^ (iL8[(t1 >> 48) & 0xFF] >> 24) ^ (iL8[(t1 >> 40) & 0xFF] >> 32) ^ (iL8[(t1 >> 32) & 0xFF] >> 40) ^ (iL8[(t1 >> 24) & 0xFF] >> 48) ^ (iL8[(t1 >> 16) & 0xFF] >> 56) ^ rKey[9];
	s2 = iL8[(t2 >> 24) & 0xFF] ^ (iL8[(t2 >> 16) & 0xFF] >> 8) ^ (iL8[(t2 >> 8) & 0xFF] >> 16) ^ (iL8[(t2) & 0xFF] >> 24) ^ (iL8[(t2 >> 56) & 0xFF] >> 32) ^ (iL8[(t2 >> 48) & 0xFF] >> 40) ^ (iL8[(t2 >> 40) & 0xFF] >> 48) ^ (iL8[(t2 >> 32) & 0xFF] >> 56) ^ rKey[10];
	s3 = iL8[(t3 >> 40) & 0xFF] ^ (iL8[(t3 >> 32) & 0xFF] >> 8) ^ (iL8[(t3 >> 24) & 0xFF] >> 16) ^ (iL8[(t3 >> 16) & 0xFF] >> 24) ^ (iL8[(t3 >> 8) & 0xFF] >> 32) ^ (iL8[(t3) & 0xFF] >> 40) ^ (iL8[(t3 >> 56) & 0xFF] >> 48) ^ (iL8[(t3 >> 48) & 0xFF] >> 56) ^ rKey[11];


	//r1
	//L2
	t0 = s1 ^ s2 ^ s3;
	t1 = s0 ^ s2 ^ s3;
	t2 = s0 ^ s1 ^ s3;
	t3 = s0 ^ s1 ^ s2;
	//printState32(s0, s1, s2, s3);
	//iL1S
	t0 = iF(t0) ^ rKey[4];
	t1 = iF(t1) ^ rKey[5];
	t2 = iF(t2) ^ rKey[6];
	t3 = iF(t3) ^ rKey[7];

	//printState32(t0, t1, t2, t3);
	//S
	s0 = iL8[(t0 >> 56) & 0xFF] ^ (iL8[(t0 >> 48) & 0xFF] >> 8) ^ (iL8[(t0 >> 40) & 0xFF] >> 16) ^ (iL8[(t0 >> 32) & 0xFF] >> 24) ^ (iL8[(t0 >> 24) & 0xFF] >> 32) ^ (iL8[(t0 >> 16) & 0xFF] >> 40) ^ (iL8[(t0 >> 8) & 0xFF] >> 48) ^ (iL8[t0 & 0xFF] >> 56) ^ rKey[0];
	s1 = iL8[(t1 >> 8) & 0xFF] ^ (iL8[(t1) & 0xFF] >> 8) ^ (iL8[(t1 >> 56) & 0xFF] >> 16) ^ (iL8[(t1 >> 48) & 0xFF] >> 24) ^ (iL8[(t1 >> 40) & 0xFF] >> 32) ^ (iL8[(t1 >> 32) & 0xFF] >> 40) ^ (iL8[(t1 >> 24) & 0xFF] >> 48) ^ (iL8[(t1 >> 16) & 0xFF] >> 56) ^ rKey[1];
	s2 = iL8[(t2 >> 24) & 0xFF] ^ (iL8[(t2 >> 16) & 0xFF] >> 8) ^ (iL8[(t2 >> 8) & 0xFF] >> 16) ^ (iL8[(t2) & 0xFF] >> 24) ^ (iL8[(t2 >> 56) & 0xFF] >> 32) ^ (iL8[(t2 >> 48) & 0xFF] >> 40) ^ (iL8[(t2 >> 40) & 0xFF] >> 48) ^ (iL8[(t2 >> 32) & 0xFF] >> 56) ^ rKey[2];
	s3 = iL8[(t3 >> 40) & 0xFF] ^ (iL8[(t3 >> 32) & 0xFF] >> 8) ^ (iL8[(t3 >> 24) & 0xFF] >> 16) ^ (iL8[(t3 >> 16) & 0xFF] >> 24) ^ (iL8[(t3 >> 8) & 0xFF] >> 32) ^ (iL8[(t3) & 0xFF] >> 40) ^ (iL8[(t3 >> 56) & 0xFF] >> 48) ^ (iL8[(t3 >> 48) & 0xFF] >> 56) ^ rKey[3];


	PUTU64(out, s0);
	PUTU64(out + 8, s1);
	PUTU64(out + 16, s2);
	PUTU64(out + 24, s3);

	return 1;
}
