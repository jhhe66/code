#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>

#include <set>
#include <map>
#include <tr1/unordered_map>
#include <tr1/unordered_set>
#include <algorithm>
#include <utility>
#include <sys/time.h>
#include <time.h>
#include <vector>

using namespace std;
using namespace std::tr1;

#if 0
#ifndef __USE_UNORDERED__
#define __USE_UNORDERED__
#endif
#endif

#if defined(__USE_UNORDERED__)
typedef unordered_set<unsigned long> 		uids_t;
typedef uids_t::iterator 					uids_itr_t;
#else
typedef set<unsigned long> 					uids_t;
typedef uids_t::iterator 					uids_itr_t;
#endif


#if defined(__USE_UNORDERED__)
typedef unordered_map<unsigned long, unsigned char>	uidms_t;
typedef uidms_t::iterator							uidms_itr_t;
#else
typedef map<unsigned long, unsigned char> 			uidms_t;
typedef uidms_t::iterator 							uidms_itr_t;
#endif

typedef vector<unsigned long> 				uidv_t;
typedef uidv_t::iterator 					uidv_itr_t;


static bool 
__set_find(uids_t* uids, unsigned long key)
{
	uids_itr_t it;
	
	if (uids == NULL) return false;

	it = uids->find(key);
	return it != uids->end() ? true : false;
}


static bool
__map_find(uidms_t* uids, unsigned long key)
{
	uidms_itr_t it;
	if (uids == NULL) return false;

	it = uids->find(key);
	return it != uids->end() ? true : false;
}

static bool
__vector_find(uidv_t* uids, unsigned long key)
{
	uidv_itr_t it;

	it = find(uids->begin(), uids->end(), key);
	
	return it != uids->end() ? true : false;
}

static bool
__array_find(unsigned long* uids, unsigned int sz, unsigned long key)
{
	unsigned long *it = NULL;

	it = find(uids, uids + sz, key);
	
	return it != uids + sz ? true : false;
}

static void
__vector_random(uidv_t* uids)
{
	unsigned long old = 0, curr = 0;
	unsigned long vsz = uids->size();
	uidv_t temp(uids->size(), 0);
	
	srand(time(NULL));
	for (unsigned int idx = 0; idx < vsz; idx++) {
		do {
			curr = rand() % vsz;
			//printf("curr: %lu\n", curr);
		} while (old == curr || curr == idx || temp[curr] != 0);

		old = curr;
		temp[curr] = (*uids)[idx];
		printf("idx: %u\n", idx);
	}

	uids->swap(temp);

#if 0
	printf("v[");
	for (unsigned int idx = 0; idx < vsz; idx++) {
		printf("%lu ", (*uids)[idx]);	
	}
	printf("]\n");
#endif
}

#define __BEGIN__(t) gettimeofday(&t, NULL)
#define __END__(b, e) {                             \	
	gettimeofday(&e, NULL);							\
	unsigned long bms = b.tv_sec * 1000000 + b.tv_usec;	\
	unsigned long ems = e.tv_sec * 1000000 + e.tv_usec;	\
	printf("elsp: %lu %lu %luus\n", bms, ems, ems - bms);\
}

int
main(int argc, char** argv)
{	
	uids_t 			users;
	uidms_t 		userms;
	uidv_t 			uservs;
	unsigned long   *useras = NULL;
	int				opt;
	unsigned int 	user_max;
	unsigned int 	key;
	struct timeval 	begin, end;
	
	if (argc < 4) return 0;

	opt = atoi(argv[1]);
	user_max = (unsigned int)atoi(argv[2]);
	key = (unsigned int)atoi(argv[3]);

	for (unsigned int idx = 0; idx < user_max; idx++) {
		switch (opt) {
			case 0: 
				users.insert(idx);
				break;
			case 1:
				userms.insert(make_pair(idx, 1));
				break;
			case 2:
				uservs.push_back(idx);
				break;
			case 3:
				if (useras == NULL) {
					useras = (unsigned long*)malloc(sizeof(unsigned long) * user_max);
				}

				*(useras + idx) = idx;
				break;
		}
	}

	switch (opt) {
		case 0: 
			__BEGIN__(begin);
			printf("set: %s\n", __set_find(&users, key) ? "true" : "false");
			__END__(begin, end)
			break;
		case 1:
			__BEGIN__(begin);
			printf("map: %s\n", __map_find(&userms, key) ? "true" : "false");
			__END__(begin, end)
			break;
		case 2:
			//__vector_random(&uservs);
			__BEGIN__(begin);
			//sort(uservs.begin(), uservs.end());
			printf("vector: %s\n", __vector_find(&uservs, key) ? "true" : "false");
			__END__(begin, end)
			break;
		case 3:
			__BEGIN__(begin);
			printf("array: %s\n", __array_find(useras, user_max, key) ? "true" : "false");
			__END__(begin, end)
			break;
	}

	sleep(30);

	return 0;
}
