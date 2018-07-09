#include <stdio.h>
#include  <stdlib.h>

int
main(int argc, char** argv)
{
	bool 	b;
	int 	i;
	enum {} EType;


	printf("size of bool: %zu\n", sizeof b);
	printf("size of int: %zu\n", sizeof i);
	printf("size of enum: %zu\n", sizeof(EType));

	b = 1;

	printf("bool is %d\n", b);

	return 0;
}
