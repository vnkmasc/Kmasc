#include "stdio.h"
#include "time.h"
#include <stdint.h>

typedef unsigned char u8;
typedef unsigned int u32;
typedef uint64_t u64;  // Changed from unsigned __int64 to uint64_t

#include "MKV256.h"

const u8 MasterKey[] = {
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27,
		0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F,
		0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
		0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F
};
u8 Plainttext[32] = {
		0xFF, 0xEE, 0xDD, 0xCC, 0xBB, 0xAA, 0x99, 0x88,
		0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x00,
		0xFF, 0xEE, 0xDD, 0xCC, 0xBB, 0xAA, 0x99, 0x88,
		0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x00
};

void print64(u64 x) {
	u32 a = (x >> 32) & 0xFFFFFFFF;
	u32 b = x & 0xFFFFFFFF;
	printf("%08X%08X, ", a, b);

}
void printState256(u64* A) {
	print64(A[0]);
	print64(A[1]);
	print64(A[2]);
	print64(A[3]);
	printf("\n");
}

void Test_MKV256(int keyLen, const u8 *key, u8 *data) {
	int r;
	if (keyLen == KeyLen256) {
		printf("\n ================= MKV 256-256 (CAI 64 BIT) =====================\n\n");
		r = Round256;
	}
	if (keyLen == KeyLen384) {
		printf("\n ================= MKV 256-384 (CAI 64 BIT) =====================\n\n");
		r = Round384;
	}
	if (keyLen == KeyLen512) {
		printf("\n ================= MKV 256-512 (CAI 64 BIT) =====================\n\n");
		r = Round512;
	}
	int i;
	u8 X[32];
	for (int i = 0; i < BlockLength; i++) X[i] = data[i];

	printf("Plainttext: ");
	for (i = 0; i < BlockLength; i++) printf("%02X ", X[i]);
	printf("\nMaster key: ");
	for (i = 0; i < keyLen/8; i++) printf("%02X ", key[i]);

	u64 rKey[(2 * Round512 + 1) * 4], irKey[(2 * Round512 + 1) * 4];
	printf("\nRound key:\n");
	KeyExpansion256(keyLen, key, rKey);
	
	int t = 0;
	for (i = 0; i < (2 * r + 1) * 4; i += 4) {
		int j = i / 4;
		if (j % 2 == 0) j = 0;
		else j = 1;
		printf("K%d%d: ", t, j);
		printState256(rKey + i);
		if (j == 1)	t++;
	}

	InvKeyExpansion256(keyLen, key, irKey);


	u8 Y[32];
	EncryptOneBlock256(keyLen, rKey, X, Y);
	printf("\nCiphertext: ");
	for (i = 0; i < BlockLength; i++) printf("%02X ", Y[i]);
	printf("\n\n");
	DecryptOneBlock256(keyLen, irKey, Y, X);
	printf("Plainttext*: ");
	for (i = 0; i < BlockLength; i++) printf("%02X ", X[i]);
	printf("\n\n");

	printf("\n ================= SPEED TEST =====================\n");

	u32 k;
	clock_t tim;
	tim = clock();
	for (k = 0; k < 20000000; k++) {
		EncryptOneBlock256(keyLen, rKey, X, X);
	}
	tim = clock() - tim;
	double duration = (double)(tim) / CLOCKS_PER_SEC;
	duration = (float)4882 / duration;
	printf(" Encryption speed of MKV256-%d: %f Mbit/seconds", keyLen, duration);

	tim = clock();
	for (k = 0; k < 20000000; k++) {
		DecryptOneBlock256(keyLen, irKey, X, X);
	}
	tim = clock() - tim;
	duration = (double)(tim) / CLOCKS_PER_SEC;
	duration = (float)4882 / duration;
	printf("\n Decryption speed of MKV256-%d: %f Mbit/seconds", keyLen, duration);
	printf("\nEnd\n\n\n");
}

int main() {
	Test_MKV256(KeyLen256, MasterKey, Plainttext);
	Test_MKV256(KeyLen384, MasterKey, Plainttext);
	Test_MKV256(KeyLen512, MasterKey, Plainttext);

	getchar();
	return 1;
}