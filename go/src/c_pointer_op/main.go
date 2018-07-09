package main

/*
#include <string.h>

int 
op_copy(char* buff, unsigned int sz)
{
	strncpy(buff, "hello", sz - 1);	
}

*/
import "C"

func main() {
	var pbuff = new([]C.char, 100)
}
