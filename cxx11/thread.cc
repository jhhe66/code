#include <stdio.h>
#include <thread>
#include <stdlib.h>
#include <unistd.h>

void
run(int stv)
{
	printf("begin sleep...\n");
	sleep(stv);
}

int
main(int argc, char** argv)
{
	std::thread t1(run, 100);
	
	printf("pid: %lu\n", t1.native_handle());

	t1.join();
}
