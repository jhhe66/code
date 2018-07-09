#include <stdlib.h>
#include <stdio.h>
#include <sys/time.h>
#include <time.h>

int
main(int argc, char** argv)
{
	struct timeval 	now;
	unsigned long 	id;

	gettimeofday(&now, NULL);
	
	id = now.tv_sec << 10 | now.tv_usec >> 10;

	printf("id: %x now.tv_sec: %x now.tv_usec: %x \n", id, now.tv_sec, now.tv_usec);

	return 0;	
}
