#include <stdio.h>
#include <stdlib.h>

constexpr static int 
__ARRAY_LEN__(int sz) 
{
	return sz * 2;
}

class A {
public:
	constexpr A(int id):_id(id) {}
	// ~A() {} 这里不能声明定义析构函数否则编译错误

	void setId(int id) { _id = id; }
	constexpr int Id() { return _id; }
private:
	int _id;
};

struct B {
	constexpr B(int id):_id(id) {}
	int _id;
};


int
main(int argc, char** argv)
{
	int temp[__ARRAY_LEN__(5)];

	printf("sizeof: %d\n", sizeof temp / sizeof(int));


	constexpr A a(10);

	enum EIDA { ID1 = a.Id() };
	
	constexpr B b(10);

	enum EID { ID_1 = b._id };

	printf("ID_1: %d\n", ID_1);

	return 0;
}
