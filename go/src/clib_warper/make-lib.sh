#!/bin/sh

go build -buildmode=c-shared -o libclib_warper.so clib_warper.go

