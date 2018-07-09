#include <boost/function.hpp>
#include <boost/bind.hpp>
#include <boost/ref.hpp>
#include <stdio.h>

using namespace boost;

typedef void anything_t();

typedef boost::function<anything_t> cb_func_t;

typedef struct cb_entry_s cb_entry_t;
struct cb_entry_s {
	cb_entry_s(const cb_func_t& cb):_cb(cb) { printf("%s\n", __PRETTY_FUNCTION__); }
	cb_func_t _cb;	
};

static int 
__self(int i)
{
	printf("i: %d\n", i);
	return i;
}

static void
__self_2()
{
	printf("%s\n", __FUNCTION__);
}


class Foo {
public:
	Foo(int id):_id(id) {}
	~Foo() {}
	int getID() {return _id;}
	void setID(int id) { _id = id; }
private:
	int _id;
};

typedef boost::function<void(int)> cb_method_t;
typedef boost::function<int()> cb_get_method_t;


int 
main(int argc, char** argv)
{
	cb_entry_t cbe(bind(__self, 1));
	cb_entry_t cbe2(bind(__self_2));
	//cb_entry_t cbe3(bind(__self, _1)); // 不能使用占位方式function<void()>

	cbe._cb();
	cbe2._cb();

	{
		Foo foo(10);

		cb_method_t cbmt = bind(&Foo::setID, ref(foo), _1);
		cb_get_method_t cbget = bind(&Foo::getID, &foo);

		cbmt(20);
		printf("id: %d\n", cbget());
	}

	return 0;
}
