#include "ice_warper.h"
#include <stdio.h>

int
main(int argc, char* argv[])
{
	printf("Result: %d\n", ice_warper_valid_appkeys(NULL, NULL, 1000));
	return 0;
}

