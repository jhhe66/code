#include <stdlib.h>
#include <stdio.h>
#include <time.h>
#include <sys/time.h>

int 
main(int argc, char** argv)
{
	printf("time_t: %zu\n", sizeof(time_t));
	printf("now: %ld\n", time(NULL));
	return 0;
}
