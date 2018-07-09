#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/socket.h>

int
main(int argc, char** argv)
{
	int 			sock;
	int 			ssz, nssz;
	int 			rsz, nrsz;
	unsigned int 	sszl = sizeof ssz;
	unsigned int 	rszl = sizeof rsz;

	sock = socket(AF_INET, SOCK_STREAM, 0);
	
	getsockopt(sock, SOL_SOCKET, SO_SNDBUF, &ssz, &sszl);
	nssz = 0;
	setsockopt(sock, SOL_SOCKET, SO_SNDBUF, &nssz, sizeof nssz);
	getsockopt(sock, SOL_SOCKET, SO_SNDBUF, &ssz, &sszl);

	getsockopt(sock, SOL_SOCKET, SO_RCVBUF, &rsz, &rszl);
	nssz = 0;
	setsockopt(sock, SOL_SOCKET, SO_RCVBUF, &nrsz, sizeof nrsz);
	getsockopt(sock, SOL_SOCKET, SO_RCVBUF, &rsz, &rszl);


	printf("SO_SNDBUF: %d\n", ssz);
	printf("SO_RCVBUF: %d\n", rsz);
	
	return 0;
}

