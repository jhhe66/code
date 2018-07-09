#/bin/sh

#go run -gcflags '-S -S' main.go
go tool compile -S main.go
