#!/bin/sh

#g++ -o redis_ben redis_ben.cc redis_helper.cc -I. -Wl,-dn -L/usr/local/lib/ -lhiredis -Wl,-dy -O2 -Wall -pipe


# static link 代码编译成对象文件的时候，不需要link对应的库，只需要头文件即可
g++ -c redis_helper.cc 
g++ -o redis_ben redis_ben.cc -I. -Wl,-dn -L/usr/local/lib/ -L./ -lhiredis -lredis_helper -Wl,-dy -O2 -Wall -pipe 
