#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <sys/time.h>

int
main(int argc, char** argv)
{
	time_t now;
	struct timeval tv;

	now = time(NULL);

	gettimeofday(&tv, NULL);

	return 0;
}
