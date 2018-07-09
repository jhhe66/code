#include <stdio.h>
#include <stdlib.h>
#include <vector>
#include <algorithm>

using std::vector;

int
main(int argc, char** argv)
{
	vector<int> 	v1(10000000, 1);	
	vector<int> 	v2(10000000, 0);
	unsigned int 	vsz = v1.size();

	for (vector<int>::iterator it = v1.begin();it != v1.end(); it++) {
		v2.push_back(*it);
	}

	for (unsigned int idx = 0; idx < vsz; idx++) {
		v2[idx] = v1[idx];
	}

	return 0;
}
