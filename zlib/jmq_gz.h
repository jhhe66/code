#ifndef __JMQ_GZ_H_
#define __JMQ_GZ_H_

#include <zlib.h>
#include <stdlib.h>
#include <string>

using std::string;

#define GZ(name) 	gz_##name
#define __RATIO__ 	10

//gz_compress
inline int
GZ(compress)(const char* in, int isz, string* out)
{
	unsigned char 	*obuff = NULL;
	unsigned long 	osz = 0;
	int 			ret = 0;

	osz = compressBound(isz);
	obuff = (unsigned char*)malloc(osz);
	if (obuff == NULL) return Z_MEM_ERROR;

	ret = compress(obuff, &osz, (const unsigned char*)in, isz);
	if (ret == Z_OK) {
		out->assign((char*)obuff, osz);
	}
		
	free(obuff);

	return ret;
}

//gz_uncompress
inline int
GZ(uncompress)(const char* in, int isz, string* out)
{
	unsigned char	*obuff = NULL;
	unsigned long 	osz = 0;
	int 			ret = 0;
	
	osz = __RATIO__ * isz;
	obuff = (unsigned char*)malloc(osz);
	if (obuff == NULL) return Z_MEM_ERROR;

	ret = uncompress(obuff, &osz, (const unsigned char*)in, isz);
	if (ret == Z_OK) {
		out->assign((char*)obuff, osz);
	}

	free(obuff);

	return ret;
}


#endif

