#include <stdio.h>
#include <stdlib.h>
#include <vector>
#include <algorithm>

using std::vector;

int
main(int argc, char** argv)
{
	vector<int> v1(10000000, 1);	
	vector<int> v2(10000000, 0);

	std::copy(v1.begin(), v1.end(), v2.begin());

	return 0;
}
