#include <stdlib.h>
#include <stdio.h>
#include <time.h>
#include <sys/time.h>
#include <string.h>

#define __MOD__(n, m) ((n) % (m))

int
main(int argc, char* argv[])
{
	int a[10] = {0, 1, 2, 3, 4, 5, 6, 7, 8, 9};
	int bucket[10];
	
	for (int i = 0; i < 10; i++) {
		*(bucket + i) = -1;
		printf("bucket[%d]: %d\n", i, bucket[i]);
	}
	
	srand(time(NULL));
	for (int i = 0; i < 10; i++) {
		while (1) {
			int seq = __MOD__(rand(), 10);
			printf("seq: %d %d\n", seq, bucket[seq]);
			if (bucket[seq] >= 0) {
				continue;
			} else {
				bucket[seq] = a[i];
				break;
			}
		}
	}
	
	for (int i = 0; i < 10; i++) {
		printf("bucket[%d]: %d\n", i, bucket[i]);
	}

	return 0;
}
