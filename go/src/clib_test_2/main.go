// clib_test_2 project main.go
package main

/*
#include <stdio.h>

void
myprintf(const char* msg) 
{
	printf("%s\n", msg);
}

*/
import "C"

func main() {
	C.myprintf(C.CString("hello"))
}
