#include "jmq_gz.h"
#include <stdio.h>
#include <string.h>

int
main(int argc, char** argv)
{
	char 	in[] = "我们的明天是gzip";
	string 	out;
	int  	ret;

	ret = gz_compress(*(argv + 1), strlen(*(argv + 1)), &out);
	if (ret) {
		printf("compress error[%d]\n", ret);

		return -1;
	}

	printf("compress: %s\n", out.c_str());

	ret = gz_uncompress(out.data(), out.size(), &out);
	if (ret) {
		printf("uncompress error[%d]\n", ret);

		return -1;
	}
	
	printf("uncompress: %s\n", out.c_str());

	return 0;
}
