#include <stdio.h>

#define test_p(v) printf("%s\n", #v)
#define test_pp(v, n) printf("%d\n", v##n)


#define __BASE__ 1000
#define __CMD__(name, value) const unsigned int name = value + __BASE__

int
main(int argc, char** argv)
{
	int i = 10;

	int i10 = 100;

	test_p(i);
	test_pp(i, 10);
	
	__CMD__(ii, 1);
	printf("cmd: %u\n", ii);

	return 0;
}
