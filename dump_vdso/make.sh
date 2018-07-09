#!/bin/sh

gcc -o dump_vdso dump_vdso.c

dump_vdso > vdso.so

objdump -T vdso.so
