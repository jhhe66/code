#!/bin/sh

go build -buildmode=c-archive -o libgo.a lib.go
gcc test_lib.c -o test_lib libgo.a -I. -pthread
