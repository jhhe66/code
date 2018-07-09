#include <stdio.h>
#include <stdlib.h>
#include <list>
#include <vector>
#include <iostream>

using namespace std;

int
main(int argc, char** argv)
{	
	list<int> 				temp;
	vector<int> 			temp_v;
	list<int>::iterator 	it;
	vector<int>::iterator 	vit;
	void 					*iv;

	vit = temp_v.begin();

	iv = static_cast<void*>(vit);

	return 0;
}
