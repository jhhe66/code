#!/bin/sh

go build -buildmode=c-shared -o libfoo.so libfoo.go
