#include <stdio.h>

class A {
public:
	A(int id):_id(id) {}
	~A() {}
	
	A(const A& rsh) {
		_id = rsh._id;
	}
	
	A& operator =(const A& rsh) { // 拷贝构造函数中可以访问其他对象的私有成员，因为他是统一个类的方法
		this->_id = rsh._id;
	}
private:
	int _id;
};


int main(int argc, char **argv)
{
	A a1(100);
	A a2(200);
	
	a1 = a2;
	
	return 0;
}
