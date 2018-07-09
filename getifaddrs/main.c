#include <unistd.h>
#include <stdio.h>
#include <sys/types.h>
#include <ifaddrs.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <errno.h>
#include <string.h>
 
int 
getInterfaceIP(const char* interface, char* ip)
{
	if (ip == NULL) {
		return -1;
	}

	struct ifaddrs* ifa = NULL;
	int rec = getifaddrs(&ifa);
	if (rec) {
		printf("getifaddrs error : %d\n", errno);
		return -1;
	}

	int code = 0;
	struct ifaddrs *ifp = ifa;
	for (; ifp != NULL; ifp = ifp->ifa_next) {
		if (ifp->ifa_addr && ifp->ifa_addr->sa_family == AF_INET
			&& !strcmp (interface, ifp->ifa_name)) {
			strcpy(ip, inet_ntoa(((struct sockaddr_in*) ifp->ifa_addr)->sin_addr));
			code = 1;
			break;
		}
	}

	freeifaddrs(ifa);

	return code;
}
 
 
int 
main(int argc, char** argv)
{
	if (argc < 2) {
		printf("usage : %s [eth0|eth1|lo]\n", argv[0]);
		return -1;
	}

	char ip[32] = {0};
	int rec = getInterfaceIP(argv[1], ip);
	if (rec > 0) {
		printf("interface(%s)'s ip is %s\n", argv[1], ip);
	}

	printf("good bye and good luck!\n");
	return 0;
}