#include <IceUtil/Handle.h>
#include <IceUtil/Shared.h>
#include <Ice/Ice.h>
#include <stdio.h>
#include <unistd.h>

struct A : virtual public IceUtil::Shared {
	A(int id):_id(id) {}
	~A() { printf("~A: %d\n", _id);}
private:
	int _id;
};

typedef IceUtil::Handle<A> APtr;

struct B {
	B(int id):_a(new A(id)) {}
	B() {}
	~B() { printf("~B\n"); }
	APtr _a;
};



int
main(int argc, char** argv)
{
	APtr aaa;
	{
		APtr a(new A(1));

		APtr aa = a;
		aaa = a;
	}

	printf("re binding\n");
	aaa = new A(2);

	aaa = NULL;

	printf("B begin \n");

	{
		B b(3);
		{
			B bb = b;
		}
		printf("back\n");
	}

	sleep(10);
	return 0;
}
