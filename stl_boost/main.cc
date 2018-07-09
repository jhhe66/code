#include <stdio.h>
#include <stdlib.h>
#include <boost/unordered_map.hpp>
#include <unordered_map>
#include <utility>
#include <unistd.h>
#include <sys/time.h>
#include <time.h>

using namespace std;
using namespace boost;

#define __BEGIN__(t) gettimeofday(&t, NULL)
#define __END__(b, e) {                             	\
    gettimeofday(&e, NULL);                         	\
	unsigned long bms = b.tv_sec * 1000000 + b.tv_usec;    \
	unsigned long ems = e.tv_sec * 1000000 + e.tv_usec;    \
	printf("elsp: %lu %lu %luus\n", bms, ems, ems - bms); \
}

static void
__proof_stl(unsigned int max)
{
	std::unordered_map<unsigned long, unsigned long> lls;
	
	for (unsigned int idx = 0; idx < max; idx++) {
		lls.insert(make_pair((unsigned long)idx, (unsigned long)idx));
	}
}

static void
__proof_boost(unsigned int max)
{
	boost::unordered_map<unsigned long, unsigned long> lls;

	for (unsigned int idx = 0; idx < max; idx++) {
		lls.insert(make_pair((unsigned long)idx, (unsigned long)idx));
	}
}

int 
main(int argc, char** argv)
{
	struct timeval begin, end;
	printf("std::unordered_map proof...\n");
	__BEGIN__(begin);
	__proof_stl(5000000);
	__END__(begin, end)

	printf("boost::unordered_map proof...\n");
	__BEGIN__(begin);
	__proof_boost(5000000);
	__END__(begin, end)
	
	return 0;
}
