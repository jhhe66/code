#!/bin/sh

gcc -o link_tc main.c -ltcmalloc -fno-builtin-malloc -fno-builtin-calloc -fno-builtin-realloc -fno-builtin-free
gcc -o link_je main.c /usr/local/lib/libjemalloc_pic.a -fno-builtin-malloc -fno-builtin-calloc -fno-builtin-realloc -fno-builtin-free -pthread -ldl
gcc -o link_je_so main.c -ljemalloc -fno-builtin-malloc -fno-builtin-calloc -fno-builtin-realloc -fno-builtin-free -pthread -ldl

