#!/bin/sh

g++ -o jmq_gz jmq_gz.cc -lz -Wall -g -pipe
g++ -o jmq_lz4 jmq_lz4.cc -llz4 -Wall -g -pipe
