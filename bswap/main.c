#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <inttypes.h>
#include <byteswap.h>

int
main(int argc, char *argv[])
{
    uint64_t x;

    if (argc != 2) {
       fprintf(stderr, "Usage: %s <num>\n", argv[0]);
       exit(EXIT_FAILURE);
    }

    x = strtoul(argv[1], NULL, 0);
    printf("0x%" PRIx64 " ==> 0x%" PRIx64 "\n", x, bswap_64(x));

    exit(EXIT_SUCCESS);
}
