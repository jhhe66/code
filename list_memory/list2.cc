#include <stdio.h>
#include <stdlib.h>
#include <list>
#include <unistd.h>

using namespace std;

int 
main(int argc, char* argv[])
{
	list<unsigned long> long_list;

	for (unsigned long i = 0; i < 100000; i++) {
		long_list.push_back(i);
	}
	
	long_list.clear();

	{
		long_list.resize(0);
		list<unsigned long> tmp;
		tmp.swap(long_list);
	}

	for (unsigned long i = 0; i < 10000000; i++) {
		long_list.push_back(i);
	}
	
	long_list.clear();
	{
		long_list.resize(0);
		list<unsigned long> tmp;
		tmp.swap(long_list);
	}
	

	sleep(300);

	return 0;
}
