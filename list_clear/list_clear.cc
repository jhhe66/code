#include <stdio.h>
#include <unistd.h>
#include <list>

using namespace std;

class Data {
public:
	Data(Data&& d) 
	{
		printf("%s\n", __PRETTY_FUNCTION__);	
	}

	Data(const Data& d) 
	{
		printf("%s\n", __PRETTY_FUNCTION__);
	}

	Data() { printf("%s\n", __PRETTY_FUNCTION__); }
	~Data() { printf("%s\n", __PRETTY_FUNCTION__); }
};


int
main(int argc, char** argv)
{	
	list<Data> 		dl;
	vector<Data> 	vl;

	Data d;

	for (int idx = 0; idx < 10; idx++) {
		dl.push_back(d);
		//dl.emplace_back(d);
	}

	dl.clear();

	sleep(100);

	return 0;
}
