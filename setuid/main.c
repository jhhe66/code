#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int
main(int argc, char* argv[])
{
	printf("current: %d\n", getuid());
	setuid(1001);
	printf("current: %d\n", getuid());


	sleep(100);

	return 0;
}
