#include "f1.h"
#include "foo.h"
#include "fn.h"
#include <stdio.h>

//extern unsigned int aaa; //上include头文件的时候，使用这种方式，必须加extern

void show_ccc();

int main(int argc, char **argv)
{
	printf("aaa: %u\n", aaa); // aaa 必须引用声明的头文件或者在引用的源文件中声明才可以使用
	
	show_var(aaa); 
	show_var_2(aaa);
	
	show_me();
	
	show_ccc();
	return 0;
}