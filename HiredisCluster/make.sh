#!/bin/sh

#gcc -o hicluster_demo main.c -I/usr/local/include/hiredis-vip -Wl,-Bstatic -lhiredis_vip -Wl,-Bdynamic
gcc -o hicluster_demo main.c -Wl,-Bstatic -lhiredis_vip -Wl,-Bdynamic
gcc -o hiscan scan.c -Wl,-dn -lhiredis_vip -Wl,-dy
