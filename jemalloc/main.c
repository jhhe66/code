#include <stdlib.h>
#include <stdio.h>
#include <jemalloc/jemalloc.h>

int
main(int argc, char** argv)
{
	int *p = NULL;

	p = malloc(sizeof(int) * 100);

	free(p);

	return 0;
}
