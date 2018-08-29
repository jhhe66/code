#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include <cpuid.h>
#include <x86intrin.h>


/*
res > 1 support	
*/
static int 
check_support_sse4_2() 
{
    int res = 0;
    __asm__ __volatile__(
                        "movl $1,%%eax\n\t"
                        "cpuid\n\t"
                        "test $0x0100000,%%ecx\n\t"
                        "jz 1f\n\t"
                        "movl $1,%0\n\t"
                        "1:\n\t"
                        :"=m"(res)
                        :
                        :"eax","ebx","ecx","edx");
    return res;
}


static int 
__see42__()
{
	unsigned int eax, ebx, ecx, edx;

	eax = ebx = ecx = edx = 0;

	__get_cpuid(1, &eax, &ebx, &ecx, &edx);

	return ecx & bit_SSE4_2 ? 1 : 0;
}


#ifdef __x86_64__
#define ALIGN_SIZE 8
#else
#define ALIGN_SIZE 4
#endif
#define ALIGN_MASK (ALIGN_SIZE - 1)
 
uint32_t extend(uint32_t init_crc, const char *data, size_t n) 
{
    uint32_t res = init_crc ^ 0xffffffff;
    size_t i;
#ifdef __x86_64__
    uint64_t *ptr_u64;
    uint64_t tmp;
#endif
    uint32_t *ptr_u32;
    uint16_t *ptr_u16;
    uint8_t *ptr_u8;
 
    // aligned to machine word's boundary
    for (i = 0; (i < n) && ((intptr_t)(data + i) & ALIGN_MASK); ++i) {
        res = _mm_crc32_u8(res, data[i]);
    }
 
#ifdef __x86_64__
    tmp = res;
    while (n - i >= sizeof(uint64_t)) {
       ptr_u64 = (uint64_t *)&data[i];
       tmp = _mm_crc32_u64(tmp, *ptr_u64);
       i += sizeof(uint64_t); 
    }
    res = (uint32_t)tmp;
#endif
    while (n - i >= sizeof(uint32_t)) {
       ptr_u32 = (uint32_t *)&data[i];
       res = _mm_crc32_u32(res, *ptr_u32);
       i += sizeof(uint32_t); 
    }
    while (n - i >= sizeof(uint16_t)) {
       ptr_u16 = (uint16_t *)&data[i];
       res = _mm_crc32_u16(res, *ptr_u16);
       i += sizeof(uint16_t); 
    }
    while (n - i >= sizeof(uint8_t)) {
       ptr_u8 = (uint8_t *)&data[i];
       res = _mm_crc32_u8(res, *ptr_u8);
       i += sizeof(uint8_t); 
    }
 
    return res ^ 0xffffffff;
}
static inline uint32_t 
crc32c(const char *data, size_t n) 
{
    return extend(0, data, n);
}





int
main(int argc, char** argv)
{
	if (check_support_sse4_2()) {
		printf("see4.2 support.\n");
		printf("crc32: %u\n", crc32c("temp", strlen("temp")));
	}

	if (__see42__()) {
		printf("cpuid method juge support see42.\n");
	}

	return 0;
}

