#include <stdio.h>
#include <utility>
#include <string>

#define __SELF__ printf("%s\n", __PRETTY_FUNCTION__)

class Name {
public:
	Name(const std::string& data):_data(data) {}
	Name(const char* data):_data(data) {}
	Name(const Name& name):_data(name._data) {}
	Name(Name&& name):_data(std::move(name._data)) { __SELF__; }

	Name& operator= (const Name& rsh) { _data = rsh._data; return *this; }
	Name& operator= (Name&& rsh) { _data = rsh._data; return *this; }
	
	std::string _data;
};


class A {
public:
	A():_name("") {}
	A(const char* name):_name(name) {}
	A(A&& a):_name(std::move(a.name())) { __SELF__; }
	
	~A() {} 

	A& operator= (const A& rsh) { __SELF__; _name = rsh.name(); return *this; }
	A& operator= (A&& rsh) { __SELF__; _name = std::move(rsh.name()); return *this; } 

	Name name() const { return _name; }
private:
	Name _name;
};


class B : public A {
private:
	std::string _name2;
};


static A
MakeA(A a) 
{
	return a;
}

static A 
MakeA2()
{
	return A();
}

int 
main(int argc, char** argv)
{
	A a("1");

	A a1 = std::move(a);
	A a2 = MakeA(A());
	A a3 = std::move(MakeA2());

	a3 = std::move(MakeA2());
	a2 = MakeA2();

	printf("Name: %s\n", a1.name()._data.c_str());

	B b1;
	B b2 = std::move(b1);

	
	return 0;	
}



