#include <stdio.h>

#define __CXX11 201103LL
#define __CXX98 199711LL


#if defined(__GXX_EXPERIMENTAL_CXX0X__)
#warning "c++ 0X"
#elif defined(__GXX_EXPERIMENTAL_CXX11__)
#warning "c++ 11"
#endif

int
main(int argc, char* argv[])
{
	
	if (__cplusplus <= __CXX98) {
		printf("is c++0x\n");
	} else if (__cplusplus >= __CXX11) {
		printf("is c++11\n");
	}
	return 0;
}
