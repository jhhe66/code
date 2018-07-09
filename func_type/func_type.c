#include <stdio.h>

typedef void (function_t)(int);

static void
show(int i)
{
	printf("i: %d\n", i);
}


static void 
binding(function_t *l, function_t *r)
{
	l = r;
}

int
main(int argc, char** argv)
{
	function_t *pf = NULL;
	function_t ff; //这样的定义没有意义，编译的时候会提示没有实现

	pf = show;

	binding(&ff, &show);

	pf(10);
	ff(20);
	

	return 0;
}
