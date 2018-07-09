#include "foo.h"
#include <stdlib.h>
#include <stdio.h>
#include <Ice/Ice.h>
#include <IceUtil/Handle.h>

typedef IceUtil::Handle<Father> FatherPtr;		

int
main(int argc, char* argv[])
{
	FatherPtr father, son;
	
	// 判断handle 是否为空	
	printf("father is created[%s]\n", father ? "true" : "false");

	father = new Father();
	son = new Son();

	printf("Age: %d\n", father->Age());
	printf("Age: %d\n", son->Age());
	//((Son*)son)->GetSon();

	return 0;
}
