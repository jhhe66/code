#include <stdio.h>
#include <vector>
#include <utility>

using namespace std;

#define __SELF__ 	printf("%s\n", __PRETTY_FUNCTION__)
#define __NNN__(v)	printf("%s\n", #v)	

class A {
public:
	A(const A& a) { __SELF__; }
	A()	{ __SELF__; }
	A(A&& a) { __SELF__; }
};

class B {
public:
	B(const B& b) { __SELF__; *this = b; }
	B(int i):_i(i) { __SELF__; }
	B(B&& b):_i(std::move(b._i)) { __SELF__; }

	B& operator=(const B& rsh) { __SELF__; this->_i = rsh._i; return *this; }
	B& operator=(B&& rsh) { __SELF__; *this = move(rsh); return *this; }
public:
	int _i;
};

int
main(int argc, char** argv)
{
	vector<A> temp;
	vector<B> vb;
	
	A a();
	B b(20);

	__NNN__(000);

	temp.push_back(A());
	
	__NNN__(100);

	temp.emplace_back(A());

	__NNN__(200);

	vb.push_back(B(10));
	
	__NNN__(300);

	vb.emplace_back(std::move(10));
	
	__NNN__(400);
	vb.push_back(b);

	__NNN__(500);
	vb.emplace_back(b);
}
