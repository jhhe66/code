#include <stdlib.h>
#include <stdio.h>
#include <list>

using namespace std;

int
main(int argc, char* argv[])
{
	list<int> ilist;

	for (int i = 0; i < 50000000; i++) {
		ilist.push_back(i);
	}
	

	printf("sz: %d\n", ilist.size());
	return 0;
}
