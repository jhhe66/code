#include <stdio.h>
#include <map>
#include <boost/shared_ptr.hpp>

using namespace std;
using namespace boost;

class A {
public:
	~A() { printf("%s\n", __FUNCTION__); }
};

int
main(int argc, char* argv[])
{
	
	map<int, shared_ptr<A> > a_map;

	{
		shared_ptr<A> si(new A);

		a_map.insert(make_pair<int, shared_ptr<A> >(1, si));
		a_map.erase(1);
	}

	printf("size: %d\n", a_map.size());

	return 0;
}
