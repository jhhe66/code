#include <stdio.h>
#include <threads.h>

static int
run(void* arg)
{
	int stv = *(int*)arg;

	printf("begin sleep...\n");
	sleep(stv);
}

int
main(int argc, char** argv)
{
	thrd_t t1;

	thrd_create(&t1, run, &100);
	

	thrd_join(t1, NULL);

	return 0;
}

