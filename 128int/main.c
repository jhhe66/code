#include <stdio.h>
#include <stdlib.h>

int
main(int argc, char** argv)
{
	__int128_t i128;
	__uint128_t  ui128 = -1;

	printf("int128 size: %zu\n", sizeof i128);
	printf("max uint128: %llu\n", ui128);
	return 0;
}
