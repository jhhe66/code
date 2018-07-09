#!/bin/sh

#go build -gcflags "-m -m" main.go
go tool compile -m main.go
