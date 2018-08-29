#include <stdio.h>

#define test_p(v) printf("%s\n", #v)
#define test_pp(v, n) printf("%d\n", v##n)


#define __BASE__ 1000
#define __CMD__(name, value) const unsigned int name = value + __BASE__

//#define __HAVE__(map, key, type) map[#key].as<type>()
#define __CON__(v1, v2) v1##v2>

#define __HAVE__(node, elem) (node[#elem].IsDefined())
#define __VALUE__(node, elem, type) node[#elem].as<type>()
#define __EVAL__(node, elem, type, v) do {			\
	if (__HAVE__(node, elem)) {						\
		v = __VALUE__(node, elem, type);			\
	} else {										\
		return RET_ERROR;							\
	}												\
} while (0)


int
main(int argc, char** argv)
{
	int i = 10;

	int i10 = 100;

	test_p(i);
	test_pp(i, 10);
	
	__CMD__(ii, 1);
	printf("cmd: %u\n", ii);


	//__HAVE__(ii, 1, int);
	//__HAVE__(ii, 1, string);
	//__HAVE__(ii, 1, unsigned int);
	__CON__(i, ii);

	__EVAL__(node, elem, int, ii);

	return 0;
}
