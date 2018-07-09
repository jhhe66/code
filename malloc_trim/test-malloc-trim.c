//#define _GNU_SOURCE  undef use -std=gnu99
#include <stdio.h>
#include <stdlib.h>
#include <strings.h>
#include <malloc.h>
#include <sys/mman.h>
#include <unistd.h>
#include <fcntl.h>

#define ALLOW_BRK  0

static char** p = NULL;
static int array_sz = 409600; // / 8; // / 32 * 2;
static char* q;

static void
myalloc() {
    int i;
    //p = (char** )malloc(sizeof(char*) * array_sz);
    p = (char** )memalign(4096, sizeof(char*) * array_sz);
    //fprintf(stderr, "%p\n", p);

    for (i = 0; i < array_sz; i++) {
        //p[i] = (char*)malloc(1024 * 4); // * 8);
        p[i] = (char*)memalign(4096, 1024 * 4); // * 8);
        //p[i] = (char*)malloc(1024 * 4 * 32 / 2);
        //fprintf(stderr, "%p\n", p[i]);
    }

    q = (char *) memalign(4096, 4096);
}

static void
myfree() {
    if (p) {
        int i;
        for (i = 0; i < array_sz; i++) {
            free(p[i]);
        }

        free(p);
    }

    p = NULL;
}

typedef struct {
    unsigned long size, resident;
} statm_t;

const char* statm_path = "/proc/self/statm";

void
getmem(const char *tag) {
#if 1
    statm_t result;

    FILE *f = fopen(statm_path,"r");
    if (f == NULL){
        perror(statm_path);
        abort();
    }

    if (fscanf(f,"%ld %ld", &result.size, &result.resident) != 2)
    {
        perror(statm_path);
        abort();
    }

    fclose(f);

    fprintf(stderr, "%s RES %lld\n", tag, (long long) result.resident);
#endif
}

int
main (int argc, char** argv) {
    char c;

    //mallopt(M_MMAP_THRESHOLD, 0);
    //mallopt(M_TOP_PAD, 0);
#if 0
    mallopt(M_TRIM_THRESHOLD, 1);
#endif

#if !ALLOW_BRK
    if (mmap(sbrk(0), 1, PROT_READ|PROT_WRITE,
                MAP_FIXED|MAP_ANONYMOUS|MAP_PRIVATE, -1, 0) == MAP_FAILED) {
        perror("mmap");
        return -1;
    }
#endif

    getmem("i");

    for (int i = 0; i < 1; i++) {
        myalloc();

        getmem("a");

        myfree();

        getmem("f");

        int rc = malloc_trim(1);
        fprintf(stderr, "trim rc = %d\n", rc);

        getmem("t");

        rc = malloc_trim(0);
        fprintf(stderr, "trim rc = %d\n", rc);

        getmem("t");

        fprintf(stderr, "---\n");
    }

    return 0;
}
