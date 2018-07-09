/* Creates a socket and binds it. Purpose: experiment with IPv6-only vs IPv6-IPv4. */

#include <stdio.h>
#include <stdbool.h>
#include <unistd.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <string.h>
#include <errno.h>
#include <arpa/inet.h>

#define PORT 8081

static char    *progname;

static void
usage()
{
    (void) fprintf(stderr,
                   "Usage: %s: [-4] [-6] (only one address family option, default is v6)\n\t[-o] [-b] (-o : v6 only, -b : both v4 and v6, only one option, default depends on the system\n\t",
                   progname);
    exit(EXIT_FAILURE);
}

char           *
text_of(struct sockaddr *address)
{
    char           *text = malloc(INET6_ADDRSTRLEN);
    struct sockaddr_in6 *address_v6;
    struct sockaddr_in *address_v4;
    if (address->sa_family == AF_INET6) {
        address_v6 = (struct sockaddr_in6 *) address;
        inet_ntop(AF_INET6, &address_v6->sin6_addr, text, INET6_ADDRSTRLEN);
    } else if (address->sa_family == AF_INET) {
        address_v4 = (struct sockaddr_in *) address;
        inet_ntop(AF_INET, &address_v4->sin_addr, text, INET_ADDRSTRLEN);
    } else {
        return ("[Unknown family address]");
    }
    return text;                /* Caller has to free it, if necessary */
}

int
main(int argc, char **argv)
{
    int             ch;
    int             family;
    int             status;
    int             fd;
    struct sockaddr_storage me;
    struct sockaddr_storage *me_p;
    struct sockaddr_in *me4;
    struct sockaddr_in6 *me6;
    struct sockaddr_storage peer;
    socklen_t       struct_sockaddr_size;
    bool           *v6only = NULL;
    const int       true_opt = 1;
    const int       false_opt = 0;

    /* Process arguments */
    progname = argv[0];
    /* Default */
    family = AF_INET6;
    while ((ch = getopt(argc, argv, "46ob")) != -1) {
        /* Some things are not checked: if the user specifies -4 and -6, the result
         * is unpredictible. And if the user specifies -4 and -b, it crashes
         * "Protocol option not available". */
        switch ((char) ch) {
        case '4':
            family = AF_INET;
            break;
        case '6':
            family = AF_INET6;
            break;
        case 'o':
            if (v6only != NULL) {
                usage();
            }
            v6only = malloc(sizeof(bool));
            *v6only = true;
            break;
        case 'b':
            if (v6only != NULL) {
                usage();
            }
            v6only = malloc(sizeof(bool));
            *v6only = false;
            break;
        default:
            usage();
        }
    }
    argc -= optind;
    argv += optind;
    if (argc != 0) {
        usage();
    }

    /* Create the socket */
    status = socket(family, SOCK_STREAM, 0);
    if (status < 0) {
        fprintf(stderr, "Cannot create the socket: %s\n", strerror(errno));
        abort();
    }
    fd = status;
    if (v6only != NULL) {
        status =
            setsockopt(fd, IPPROTO_IPV6, IPV6_V6ONLY,
                       (*v6only ? &true_opt : &false_opt), sizeof(int));
        if (status < 0) {
            fprintf(stderr, "Cannot set options on the socket: %s\n",
                    strerror(errno));
            abort();
        }
    }

    /* Bind the socket */
    memset(&me, 0, sizeof(me));
    if (family == AF_INET) {
        struct_sockaddr_size = (socklen_t) sizeof(struct sockaddr_in);
        me4 = (struct sockaddr_in *) &me;
        me4->sin_family = family;
        me4->sin_addr.s_addr = htonl(INADDR_ANY);
        me4->sin_port = htons(PORT);
        me_p = memcpy(&me, me4, sizeof(*me4));
    } else if (family == AF_INET6) {
        struct_sockaddr_size = (socklen_t) sizeof(struct sockaddr_in6);
        me6 = (struct sockaddr_in6 *) &me;
        me6->sin6_family = family;
        me6->sin6_addr = in6addr_any;
        me6->sin6_port = htons(PORT);
        me_p = memcpy(&me, me6, sizeof(*me6));
    } else {
        fprintf(stderr, "Unsupported address family %i\n", family);
        abort();
    }
    status = bind(fd, (struct sockaddr *) me_p, struct_sockaddr_size);
    if (status < 0) {
        fprintf(stderr, "Cannot bind the socket: %s\n", strerror(errno));
        abort();
    }

    /* Start listening on the socket and accept connections */
    status = listen(fd, 1);
    if (status < 0) {
        fprintf(stderr, "Cannot listen to the socket: %s\n", strerror(errno));
        abort();
    }
    fprintf(stdout,
            "Now listening: you can use 'netstat -l -n | grep %i' to see what is going on.\nControl-C to terminate.\n",
            PORT);
    while (1) {
        status = accept(fd, (struct sockaddr *) &peer, &struct_sockaddr_size);
        if (status < 0) {
            fprintf(stderr, "Cannot accept connections on the socket: %s\n",
                    strerror(errno));
            abort();
        }
        fprintf(stdout, "One connection accepted from %s\n",
                text_of((struct sockaddr *) &peer));
    }

    /* Never reached */
    return 0;
}
