#include <stdio.h>
#include <stdlib.h>
#include "hiredis-vip/hircluster.h"

int 
main(int argc, char* argv[])
{
    const char *key		= "key-a";
    const char *field	= "field-1";
    const char *key1	= "key1";
    const char *value1	= "value-1";
    const char *key2	= "key2";
    const char *value2	= "value-2";
	const char *key3	= "key3";
	const char *value3	= "value-3";

	if (argc < 2) {
		printf("hicluster_demo hosts\n");
		exit(-1);
	}

    redisClusterContext *cc = redisClusterConnect(argv[1], HIRCLUSTER_FLAG_NULL);
    if(cc == NULL || cc->err)
    {
        printf("connect error : %s\n", cc == NULL ? "NULL" : cc->errstr);
        return -1;
    }

    redisReply *reply = redisClusterCommand(cc, "hmget %s %s", key, field);
    if(reply == NULL)
    {
        printf("reply is null[%s]\n", cc->errstr);
        redisClusterFree(cc);
        return -1;
    }

    printf("reply->type:%d\n", reply->type);

    freeReplyObject(reply);

    reply = redisClusterCommand(cc, "mset %s %s %s %s %s %s", key1, value1, key2, value2, key3, value3);
    if(reply == NULL)
    {
        printf("reply is null[%s]\n", cc->errstr);
        redisClusterFree(cc);
        return -1;
    }

    printf("reply->str:%s\n", reply->str);

    freeReplyObject(reply);
    redisClusterFree(cc);
    return 0;
}