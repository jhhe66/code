#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "hiredis-vip/hircluster.h"

int 
main(int argc, char** argv)
{
	redisClusterContext *rc;
	redisReply			*reply;
	struct timeval		timeout;
	
	timeout.tv_sec = 5;
	timeout.tv_usec = 0;

	rc = redisClusterConnectWithTimeout(*(argv + 1), timeout, HIRCLUSTER_FLAG_ROUTE_USE_SLOTS);
	if (rc && rc->err == 0) {
		printf("Connected Cluster Succ.\n");		
	} else {
		printf("Connected Cluster failed.\n");
		return -1;
	}

	reply = redisClusterCommand(rc, "scan %d MATCH %s COUNT %d", 0, "key*", 10);
	if (reply) {
		printf("Type: %d\n", reply->type);
	} else {
		printf("reply is null.\n");		
	}

	freeReplyObject(reply);
	redisClusterFree(rc);

	return 0;
}

