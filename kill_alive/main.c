#include <stdio.h>
#include <signal.h>

#define SIG_ALIVE 0

int
main(int argc, char** argv)
{
	int pid;
	
	if (argc < 2) {
		printf("usage: alive pid");
		return 0;
	}

	pid = atoi(*(argv + 1));

	if (kill(pid, SIG_ALIVE)) {
		printf("process[%d] not found.\n", pid);
	} else {
		printf("process[%d] is alived.\n", pid);
	}

	return 0;
}
