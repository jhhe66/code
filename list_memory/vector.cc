#include <stdio.h>
#include <stdlib.h>
#include <list>
#include <vector>
#include <unistd.h>

using namespace std;

int 
main(int argc, char* argv[])
{
	vector<unsigned long> long_list;

	for (unsigned long i = 0; i < 10000000; i++) {
		long_list.push_back(i);
	}
	
	long_list.clear();
	
	{
		vector<unsigned long>().swap(long_list);
	}

	sleep(300);

	return 0;
}
