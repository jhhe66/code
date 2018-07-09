#ifndef __TO_SIGNED_H_
#define __TO_SIGNED_H_

#include <stdint.h>

inline double 
ToSigned(uint64_t n, unsigned int val)
{	
	double dRet = val == 0 ? 
		(int64_t)n : (double)((int64_t)n) / val; 
	
	return dRet;
}

inline double 
ToSigned(unsigned int n, unsigned int val)
{
	double dRet = val == 0 ? 
		(int)n : (double)((int)n) / val; 

	return dRet;
}


#endif
