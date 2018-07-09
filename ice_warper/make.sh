#!/bin/sh

g++ -o libice_warper.so ice_warper.cc QueryUserApi.cpp  -shared -fPIC -I./
#ar -rcu libice_warper.a  ice_warper.o
gcc -o test_warper test_warper.c -L. -lice_warper -lIce -lIceUtil
