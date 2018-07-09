#!/bin/sh

gcc -o kafka_c example.c -I/usr/include/librdkafka -Wl,-Bstatic -lrdkafka -Wl,-Bdynamic -pthread -ldl -lrt -lz -llz4 -lssl -lcrypto -g -Wall
