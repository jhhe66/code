#include "hiredis/hiredis.h"
#include <stdlib.h>
#include <stdio.h>

#define __REDIS_HOST__ "127.0.0.1"
#define __REDIS_PORT__ 6379

#define OP_COUNT 1000

int
main(int argc, char* argv[])
{
	redisContext 	*redis;
	redisReply		*reply;
	char 			key[256] = {0};
	int 			cursor = 0;

	redis = redisConnect(__REDIS_HOST__, __REDIS_PORT__);
	if (redis) {
		printf("redis connect success.\n");
	}
	
	// sadd
	
	for (unsigned int idx = 0; idx < OP_COUNT; idx++) {
		reply = redisCommand(redis, "sadd %s %d", "aaaaaa", idx);
		free(reply);
	}

	do {
		reply = redisCommand(redis, "sscan %s %d count %d", "aaaaa", cursor, 20);
		if (reply && reply->type == REDIS_REPLY_ARRAY) {
			
			printf("Type: %d\n", (*(reply->element))->type);
			if ((*(reply->element))->type == REDIS_REPLY_STRING) {
				printf("CURSOR: %s\n", (*(reply->element))->str);
				cursor = atoi((*(reply->element))->str);
			}
			for (unsigned int idx = 0; idx < (*(reply->element + 1))->elements; idx++) {
				printf("ELE[%u]: %s\n", idx, (*((*(reply->element + 1))->element + idx))->str);
			}
		} else {
			printf("query failed.\n");
		}
		free(reply);
	} while ( cursor > 0 );

#if 0
	reply = redisCommand(redis, "sscan %s %d count %d", "aaaaaa", cursor, 1);
	if (reply && reply->type == REDIS_REPLY_ARRAY) {
		printf("CURSOR: %d\n", (*(reply->element))->integer);
		for (unsigned int idx = 0; idx < (*(reply->element + 1))->elements; idx++) {
			printf("ELE[%u]: %s\n", idx, (*((*(reply->element + 1))->element + idx))->str);
		}
	} else {
		printf("query failed.\n");
	}
#endif
	
	return 0;
}
