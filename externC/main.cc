#include <stdio.h>
#include "foo.h"

extern void call_a();
extern void call_aa();


int
main(int argc, char** argv)
{
	call_a();
	call_aa();

	return 0;
}
