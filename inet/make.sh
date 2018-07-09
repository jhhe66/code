#!/bin/sh

gcc -o inet_ntop inet_ntop.c -g -Wall 
gcc -o inet_pton inet_pton.c -g -Wall 

inet_ntop www.baidu.com
inet_pton i4 192.168.1.1
