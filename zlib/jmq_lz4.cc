#include "jmq_lz4.h"
#include <stdio.h>
#include <string.h>
#include <string>

using std::string;

int
main(int argc, char** argv)
{
	int 	ret = 0;
	string 	out;

	ret = lz4_compress(*(argv + 1), strlen(*(argv + 1)), &out);
	if (ret == LZ4_OK) {
		printf("compress: %s\n", out.c_str());
	} else {
		printf("compress error[%d]\n", ret);
	}

	ret = lz4_uncompress(out.data(), out.size(), &out);
	if (ret == LZ4_OK) {
		printf("uncompress: %s\n", out.c_str());
	} else {
		printf("uncompress error[%d]\n", ret);
	}


	return 0;
}
