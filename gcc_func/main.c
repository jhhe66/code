#include <stdio.h>

extern int id;

int
main(int argc, char* argv[])
{
	get_id_2();
	id = 200;
	return 0;
}
