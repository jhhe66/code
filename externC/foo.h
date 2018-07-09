#ifndef __FOO_H_
#define __FOO_H_

#if 1 // 编译不通过
extern "C" {
	int a;
}
#endif


#if 0 // 编译通过
extern "C" { 
	extern int a;
}
#endif

#if 0 // 编译通过
extern "C" int a;
#endif


#endif
