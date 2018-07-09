#include <boost/shared_ptr.hpp>
#include <stdio.h>

class A {
public:
	A() {};
	~A() {printf("%s\n", __FUNCTION__);}
};

int
main(int argc, char* argv[])
{
	{
		boost::shared_ptr<A> pa(new A);
	
		boost::shared_ptr<A> pa1(new A);

		pa1 = pa;

		printf("leave\n");
	}

	return 0;
}
