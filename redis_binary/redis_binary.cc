#include "redis_helper.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <string>
#include <stddef.h>

typedef struct data_s data_t;
struct data_s {
	unsigned int 	_id;
	char			_name[64];		
};

static void 
__set_bin_test(CRedisHelper* redis, unsigned int loop)
{
	data_t *data = NULL;	


	if (redis == NULL) {
		return;
	}

	for (unsigned int idx = 0; idx < loop; idx++) {
		data = (data_t*)malloc(sizeof(data_t));
		if (data == NULL) break;

		data->_id = idx;
		snprintf(data->_name, sizeof data->_name, "%063u", idx);
		printf("Name: %s\n", data->_name);

		redis->Enqueue("binary_test", string((char*)data, sizeof *data));
	
		free(data);
	}
}

int
main(int argc, char** argv)
{
	CRedisHelper redis;

	if (argc < 4) {
		printf("usage: redis_binary host port loop\n");
		return 0;
	}

	if (redis.Connect(*(argv + 1), atoi(*(argv + 2))) != REDIS_OK) {
		printf("redis connected failed.\n");
		return -1;	
	}

	__set_bin_test(&redis, atoi(*(argv + 3)));


	return 0;
}
