#!/bin/sh

gcc -o redis_libevent redis_libevent.c -levent -L/usr/local/lib -lhiredis -std=gnu99 -Wall -ggdb3
gcc -o redis_libevent_loop redis_libevent_loop.c -levent -L/usr/local/lib -lhiredis -std=gnu99 -Wall -ggdb3


