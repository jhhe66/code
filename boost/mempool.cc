#include <stdio.h>
#include <boost/pool/pool.hpp>

int 
main(int argc, char* argv[])
{
	boost::pool<> 	pool(100);
	void 			*temp = NULL;
	
	temp = pool.malloc();

	printf("alloc_siz: %d\n", pool.alloc_size());

	pool.free(temp);

	temp = pool.malloc();

	printf("alloc_siz: %d\n", pool.alloc_size());

	pool.free(temp);

	return 0;
}
