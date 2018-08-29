#!/bin/sh

g++ -o redis_binary redis_binary.cc redis_helper.cc -lhiredis -g -O2 -Wall
