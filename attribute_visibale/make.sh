#!/bin/sh

g++ -shared -fPIC -o libfunc.so func.cc

g++ -o v_1 main.cc -lfunc -I . -L .

g++ main_2.cc func.cc -I . -o v_2
