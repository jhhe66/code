#include <stdio.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <net/if_dl.h>
#include <ifaddrs.h>

int listmacaddrs(void) {
    struct ifaddrs *ifap, *ifaptr;
    unsigned char *ptr;

    if (getifaddrs(&ifap) == 0) {
        for(ifaptr = ifap; ifaptr != NULL; ifaptr = (ifaptr)->ifa_next) {
            if (((ifaptr)->ifa_addr)->sa_family == AF_LINK) {
                ptr = (unsigned char *)LLADDR((struct sockaddr_dl *)(ifaptr)->ifa_addr);
                printf("%s: %02x:%02x:%02x:%02x:%02x:%02x\n",
                                    (ifaptr)->ifa_name,
                                    *ptr, *(ptr+1), *(ptr+2), *(ptr+3), *(ptr+4), *(ptr+5));
            }
        }
        freeifaddrs(ifap);
        return 1;
    } else {
        return 0;
    }   
}

int macaddr(char *ifname, char *macaddrstr) {
    struct ifaddrs *ifap, *ifaptr;
    unsigned char *ptr;

    if (getifaddrs(&ifap) == 0) {
        for(ifaptr = ifap; ifaptr != NULL; ifaptr = (ifaptr)->ifa_next) {
            if (!strcmp((ifaptr)->ifa_name, ifname) && (((ifaptr)->ifa_addr)->sa_family == AF_LINK)) {
                ptr = (unsigned char *)LLADDR((struct sockaddr_dl *)(ifaptr)->ifa_addr);
                sprintf(macaddrstr, "%02x:%02x:%02x:%02x:%02x:%02x",
                                    *ptr, *(ptr+1), *(ptr+2), *(ptr+3), *(ptr+4), *(ptr+5));
                break;
            }
        }
        freeifaddrs(ifap);
        return ifaptr != NULL;
    } else {
        return 0;
    }
}

extern int
main(int argc, char* argv[]) {

    char macaddrstr[18], *ifname;

    if (argc == 2) {
        ifname = argv[1];
        if (!strcmp(ifname,"-l")) {
            return listmacaddrs();
        } else {
            if (macaddr(ifname, macaddrstr)) {
                printf("%s: %s\n", ifname, macaddrstr);
                return 0;
            } else {
                printf("%s: not found\n", ifname);
                return 1;
            }
        }
    } else {
        printf("list all interfaces: %s -l\n", argv[0]);
        printf("single interface: %s interface_name\n", argv[0]);
        return 2;
    }
}