#include <stdlib.h>
#include "def.h"
#include "tag.h"

int
main(int argc, char** argv)
{
	tag_t tag;

	tag.id = 10;

	getID(&tag);	
	
	return 0;
}
