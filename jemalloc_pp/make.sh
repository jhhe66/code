#!/bin/sh

g++ -o jemalloc_tc_a main.cc /usr/local/lib/libtcmalloc_minimal.a -lunwind
g++ -o jemalloc_tc_so main.cc -ltcmalloc -fno-builtin-malloc -fno-builtin-calloc -fno-builtin-realloc -fno-builtin-free
g++ -o jemalloc_je_a main.cc /usr/local/lib/libjemalloc_pic.a -pthread -ldl
g++ -o jemalloc_je_so main.cc -ljemalloc -pthread -ldl

