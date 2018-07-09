#include <functional>
#include <stdlib.h>
#include <stdio.h>

using namespace std;

typedef function<void()> cb_func_t; // cb_func_t可以使用bind 指向任何void 类型的function

typedef struct callback_s callback_t;
struct callback_s {
	callback_s(const cb_func_t& cb):_cb(cb) { printf("%s\n", __PRETTY_FUNCTION__);}
	//callback_s(cb_func_t cb):_cb(cb) { printf("%s\n", __PRETTY_FUNCTION__);}
	cb_func_t _cb;
};

static int 
self(int i)
{
	return i;
}

static void
self2(float f)
{
	return;
}

int
main(int argc, char** argv)
{
	callback_t cb(bind(self, 1));
	callback_t cb2(bind(self2, 1.0));
	//callback_t cb3(boost::bind(self, _1)); // 这是不允许的
	
	cb._cb(); // 这里的函数调用是没有返回值的因为cb_func_t 定义就是void
	cb2._cb();
	//cb3._cb(3);


	return 0;
}

