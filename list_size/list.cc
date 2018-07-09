#include <stdio.h>
#include <list>
#include <string.h>
#include <string>

using namespace std;

int
main(int argc, char* argv[])
{
	list<string> str_list;

	for (int i = 0; i < 5000000; i++) {
		str_list.push_back("b500cb8d7b70984eca9e634b_864231033551389_50:01:d9:53:4f:7c");
	}

	printf("sz: %zu\n", str_list.size());

	return 0;
}
