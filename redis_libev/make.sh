#!/bin/sh

#link 0.13.3
gcc -o redis_libev redis_libev.c -I/usr/local/include -L/usr/local/lib -lhiredis -lev -std=gnu99 -g -Wall
#link 0.12
#gcc -o redis_libev redis_libev.c -lhiredis -lev -std=gnu99 -g -Wall
gcc -o redis_libev_loop redis_libev_loop.c -I/usr/local/include -L/usr/local/lib -lhiredis -lev -std=gnu99 -g -Wall
