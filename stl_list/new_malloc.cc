#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <malloc.h>

static void
__use_mmap()
{
    mallopt(M_MMAP_THRESHOLD, 1 * 1024);
}

int
main(int argc, char** argv)
{
	int *pi = NULL;

	__use_mmap();

	//pi = (int*)malloc(8 * 1024);
	pi = new int[8 * 1024];

	sleep(1000);
}
