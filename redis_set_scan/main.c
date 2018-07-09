#include "hiredis/hiredis.h"
#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <sys/time.h>
#include <unistd.h>

int
main(int argc, char* argv[])
{
	redisContext 	*redis;
	redisReply 		*reply;
	struct timeval	tv;
	
	tv.tv_sec = 5;
	tv.tv_usec = 0;
	
	redis = redisConnectWithTimeout("127.0.0.1", 6379, tv);
	if (redis) {
		reply = redisCommand(redis, "sscan %s %d count %d", "chenbo_set", 0, 100);
		if (reply) {
			if (reply->type ==  REDIS_REPLY_ARRAY) {
				for (unsigned int idx = 0; idx < reply->elements; idx++) {
					switch (reply->element[idx]->type) {
						case REDIS_REPLY_STRING:
							printf("str: %s\n", reply->element[idx]->str);
							break;
						case REDIS_REPLY_ARRAY:
							printf("array: %d\n", reply->element[idx]->elements);
							for (unsigned int inx = 0; inx < reply->element[idx]->elements; inx++) {
								printf("array[%d]: %s\n", inx, reply->element[idx]->element[inx]->str);
							}
							break;
						default:
							printf("no more\n");
							break;
					}
				}
			}
		}
	}
	
	return 0;
}

