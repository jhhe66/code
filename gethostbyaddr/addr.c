#include <netdb.h>
#include <sys/socket.h>
#include <stdlib.h>
#include <stdio.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>

int 
main(int argc, char **argv)
{
    char *ptr;
    char **pptr;
    char str[32];
    char ipaddr[16];
    struct hostent *hptr;
    struct in_addr hipaddr;

    ptr = argv[1];  // 取得命令后第一个参数，即要解析的IP地址
    /* 调用inet_aton()，ptr就是以字符串存放的地方的指针，hipaddr是in_addr形式的地址 */
    if(!inet_aton(ptr, &hipaddr))
    {
        printf("inet_aton error\n");
        return 1;
    }

    /* 调用gethostbyaddr()，调用结果都存在hptr中 */
    if ((hptr = gethostbyaddr(&hipaddr, 4, AF_INET) ) == NULL )
    {
        printf("gethostbyaddr error for addr:%s\n", ptr);
        return 1;
    }
    
    printf("hostname:%s\n",hptr->h_name);  // 将主机的规范名打出来
    for (pptr = hptr->h_aliases; *pptr != NULL; pptr++)  // 将主机所有别名分别打出来
        printf("  alias:%s\n",*pptr);
	
	printf("__LINE__\n");
    /* 根据地址类型，将地址打出来 */
    switch (hptr->h_addrtype)
    {
        case AF_INET:
        case AF_INET6:
            pptr=hptr->h_addr_list;
            /* 将刚才得到的所有地址都打出来。其中调用了inet_ntop()函数 */
            for(; *pptr!=NULL; pptr++)
                printf("address:%s\n", inet_ntop(hptr->h_addrtype, *pptr, str, sizeof(str)));
            break;
        default:
            printf("unknown address type\n");
            break;
    }

    return 0;
}
