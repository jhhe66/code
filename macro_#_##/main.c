#include <stdio.h>

#define test_p(v) printf("%s\n", #v)
#define test_pp(v, n) printf("%d\n", v##n)

int
main(int argc, char** argv)
{
	int i = 10;

	int i10 = 100;

	test_p(i);
	test_pp(i, 10);

	return 0;
}
