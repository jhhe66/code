#ifndef __JMQ_LZ4_H_
#define __JMQ_LZ4_H_

#include <lz4.h>
#include <string>
#include <stdlib.h>

using std::string;

#define LZ4(name) lz4_##name
#define __RATIO__ 10

#define LZ4_OK 		0
#define LZ4_BUF_ERR -1
#define LZ4_ERR 	-2

inline int
LZ4(compress)(const char* in, int isz, string* out)
{
	char 	*obuff = NULL;
	int 	osz = 0;
	int 	ret = 0;

	osz = LZ4_compressBound(isz);
	obuff = (char*)malloc(osz);
	if (obuff == NULL) return LZ4_BUF_ERR;
	
	ret = LZ4_compress_default(in, obuff, isz, osz);
	if (ret) {
		out->assign(obuff, ret);
	}

	free(obuff);

	return ret ? LZ4_OK : LZ4_ERR;
}

inline int
LZ4(uncompress)(const char* in, int isz, string* out)
{
	char 	*obuff = NULL;
	int 	osz = 0;
	int 	ret = 0;

	osz = __RATIO__ * isz;
	obuff = (char*)malloc(osz);
	if (obuff == NULL) return LZ4_BUF_ERR;
	
	ret = LZ4_decompress_safe(in, obuff, isz, osz);
	if (ret) {
		out->assign(obuff, ret);
	}

	free(obuff);

	return ret ? LZ4_OK : LZ4_ERR;
}

#endif 

