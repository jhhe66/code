#ifndef __FOO_H_
#define __FOO_H_

#include <IceUtil/Handle.h>
#include <IceUtil/Shared.h>

class Father : public IceUtil::Shared {
public:
	Father() {}
	virtual ~Father() {}

	int Age() { return _age; }
	void SetAge(int age) { _age = age; }
private:
	int _age;
};


class Son : public Father {
public:
	Son() {}
	virtual ~Son() {}
	
	void GetSon() {}
};


#endif

