#!/bin/sh 

gcc -o liba.so.1.0.0 liba.c -Wl,-soname,liba.so.1 -shared -fPIC
gcc -o libaaa.so.1.0.0 liba.c -shared -fPIC

gcc -o liba.so.2.0.0 liba2.c -Wl,-soname,liba.so.2 -shared -fPIC


#invoke

gcc -o user_a_1 main.c -L. -la -static-libgcc 
gcc -o user_a_11 main.c -L. -laaa -static-libgcc 

gcc -o user_a_2 main2.c -L. -l:liba.so.2 -static-libgcc 
#gcc -o user_a_2 main2.c -L. -la-2 -static-libgcc 
