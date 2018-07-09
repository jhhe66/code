#include <stdio.h>
#include <stdlib.h>

int
main(int argc, char** argv)
{
	int 		i;
	char 		key[16] = "world";
	const char 	*p = "hello";
	char 		value[16] = "value";

	snprintf(value, sizeof value, p ? "%s_%s" : "%s_%d", "hello", p ? key : i);

	return 0;
	
}
