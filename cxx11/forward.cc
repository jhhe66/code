#include <stdio.h>
#include <utility>

using namespace std;

#define __SELF__ 	printf("%s:%d\n", __PRETTY_FUNCTION__, __LINE__)
#define __SPACE__(v)	printf("%s\n", #v)


void
receiver(int&& t)
{
	__SELF__;
}

void
receiver(const int& t)
{
	__SELF__;
}

template <typename T>
void 
receiver2(T&& t)
{
	__SELF__;
}

template <typename T>
void 
receiver2(const T& t)
{
	__SELF__;
}


template<typename T> 
void wrapper(T&& t) 
{
	//receiver2(forward<T>(t));	
	receiver(forward<T>(t));	
	receiver(t);	
}

int 
main(int argc, char** argv)
{
	int i = 10;

	wrapper(i);
	
	__SPACE__(i);	
	
	wrapper(10);

	return 0;	
}
