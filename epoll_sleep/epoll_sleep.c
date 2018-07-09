#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <sys/epoll.h>

int
main(int argc, char** argv)
{	
	int efd;

	efd = epoll_create(256);
	if (!efd) {
		printf("epoll create failed.\n");
		return -1;
	}
	printf("sleep begin.\n");

	epoll_wait(efd, NULL, 1, 1000);

	printf("sleep end.\n");

	return 0;
}
