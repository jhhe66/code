#include <stdio.h>
#include <stdlib.h>

int
main(int argc, char** argv)
{
	int x = 6944350, y = 8388608;
	float dd = (float)x / (float)y;


	printf("result: %0.6f\n", dd);
	printf("result: %d\n", (int)(dd * 100));

	return 0;
}
