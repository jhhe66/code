#include "redis_helper.h"
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#define __PRT_LEN__ 	"1024"
#define __R_BUFF_LEN__	1024 

#define __START__(b) 		gettimeofday(&b, NULL);
#define __STOP__(b, e, t) 	gettimeofday(&e, NULL);printf("%s elsp: %lu\n", t, ((e.tv_sec << 10) + (e.tv_usec >> 10) - ((b.tv_sec << 10) + (b.tv_usec >> 10))));

#define __FREE__(r) do {							\
	if (r) { freeReplyObject(r); r = NULL;}			\
} while (0)

static void 
set_test(CRedisHelper* redis, unsigned long down, unsigned long up)
{
	char key_buff[64];
	char value_buff[__R_BUFF_LEN__ + 1];
	struct timeval  begin, end;
	
	for (unsigned long i = down; i < up; i++) {
		snprintf(key_buff, sizeof key_buff, "%057lu", i);
		i % 2 ? snprintf(value_buff, sizeof value_buff, "%0"__PRT_LEN__"lu", i) : 
				snprintf(value_buff, sizeof value_buff, "%lu", i);
		__START__(begin)
		redis->Set(key_buff, value_buff);
		__STOP__(begin, end, "set")
	
#if 0
	if (redis->IsActived()) {
			redis->Set(key_buff, value_buff);
		}
#endif	
	}		
}

static void
get_test(CRedisHelper* redis, unsigned long down, unsigned long up)
{
	char 			key_buff[64];
	char 			value_buff[__R_BUFF_LEN__ + 1];
	unsigned int	value_len;
	struct timeval  begin, end;
	
	for (unsigned long i = down; i < up; i++) {
		snprintf(key_buff, sizeof key_buff, "%057lu", i);
	
		value_len = sizeof value_buff;
		
		__START__(begin)
		redis->Get(key_buff, value_buff, value_len);
		__STOP__(begin, end, "get")
//		if (strlen(value_buff) == 0) {
//			printf("[%s]get failed.\n", key_buff);
//		}
#if 0
	if (redis->IsActived()) {
			redis->Get(key_buff, value_buff, sizeof value_buff);
		}
#endif
	}		
}

static void
del_test(CRedisHelper* redis, unsigned long down, unsigned long up)
{
	char 			key_buff[64];
	
	for (unsigned long i = down; i < up; i++) {
		snprintf(key_buff, sizeof key_buff, "%057lu", i);
	
		redis->Del(key_buff);
	}	
}

static void
mix_test(CRedisHelper* redis, unsigned long down, unsigned long up)
{
	char 			key_buff[64];
	char 			value_buff[__R_BUFF_LEN__ + 1];
	unsigned int	value_len; 
	struct timeval  begin, end;

	for (unsigned long i = down; i < up; i++) {
		snprintf(key_buff, sizeof key_buff, "%057lu", i);
		
		value_len = sizeof value_buff;
		
		__START__(begin)
		redis->Get(key_buff, value_buff, value_len);
		__STOP__(begin, end, "get")	

		__START__(begin)
		redis->Del(key_buff);
		__STOP__(begin, end, "del")

		snprintf(value_buff, sizeof value_buff, "%01024lu", i);
		
		__START__(begin)
		redis->Set(key_buff, value_buff);
		__STOP__(begin, end, "set")
	}
}

static void
hash_test(const char* host, const unsigned short port, unsigned long down, unsigned long up)
{
	redisContext	*redis;
	redisReply		*reply;
	char 			key_buff[64];
	vector<string>	values;
	struct timeval  begin, end;

	redis = redisConnect(host, port);
	if (!redis) {
		printf("redis connect faild.\n");
		return;
	}
	
	for (unsigned long idx = down; idx < up; idx++ ) {
		snprintf(key_buff, sizeof key_buff, "%lu", idx);

		__START__(begin)
		reply = (redisReply*)redisCommand(redis, "hmset %s %s %s %s %lu %s %d %s %d %lu %s %lu %s", 
							 key_buff, 
							 "appkey", 
			                 "8e5cf4ca06728a452207f701",
							 "ttl",
			                 idx,
							 "max_count",
			                 5,
							 "num",
							 2,
							 1099985724LU,
							 "stime:1496286894,ctime:1496286894,ttl:1496373294",
							 1101575128LU,
			                 "stime:1496376340,ctime:1496376340,ttl:1496289940");

		if (reply && reply->type == REDIS_REPLY_STATUS && !strcmp("OK", reply->str)) {

		} else {
			printf("hmset faild %s\n", key_buff);
		}
		__FREE__(reply);
		__STOP__(begin, end, "hmset")
	
		__START__(begin)
		reply = (redisReply*)redisCommand(redis, "hgetall %s", key_buff);
		if (reply && reply->type == REDIS_REPLY_ARRAY) {

		} else {
			printf("hgetall faild %s.\n", key_buff);
		}
		__FREE__(reply);
		__STOP__(begin, end, "getall")
		
		__START__(begin)
		reply = (redisReply*)redisCommand(redis, "hdel %s %lu", key_buff, 1099985724LU);
		if (reply && reply->type == REDIS_REPLY_INTEGER && reply->integer == 1) {
		
		} else {
		
		}
		__FREE__(reply);
		__STOP__(begin, end, "hdel")
	}

	redisFree(redis);
}

static void
hash_set(const char* host, const unsigned short port, unsigned long down, unsigned long up)
{
	redisContext	*redis;
	redisReply		*reply;
	char 			key_buff[64];
	vector<string>	values;
	struct timeval  begin, end;

	redis = redisConnect(host, port);
	if (!redis) {
		printf("redis connect faild.\n");
		return;
	}
	
	for (unsigned long idx = down; idx < up; idx++ ) {
		snprintf(key_buff, sizeof key_buff, "%lu", idx);

		__START__(begin)
		reply = (redisReply*)redisCommand(redis, "hmset %s %s %s %s %lu %s %d %s %d %lu %s %lu %s", 
							 key_buff, 
							 "appkey", 
			                 "8e5cf4ca06728a452207f701",
							 "ttl",
			                 idx,
							 "max_count",
			                 5,
							 "num",
							 2,
							 1099985724LU,
							 "stime:1496286894,ctime:1496286894,ttl:1496373294",
							 1101575128LU,
			                 "stime:1496376340,ctime:1496376340,ttl:1496289940");

		if (reply && reply->type == REDIS_REPLY_STATUS && !strcmp("OK", reply->str)) {

		} else {
			printf("hmset faild %s\n", key_buff);
		}
		__FREE__(reply);
		__STOP__(begin, end, "hmset")
	}

	redisFree(redis);
}

static void
msg_test(unsigned long down, unsigned long up)
{
	redisContext 	**redis_list;
	redisReply		*reply;
	char 			key_buff[64];
	struct timeval  begin, end;
	unsigned long 	ridx;

	redis_list = (redisContext**)malloc(sizeof *redis_list * 2);
	if (redis_list) {
		*redis_list = redisConnect("172.16.203.112", 9221);
		if (*redis_list == NULL) {
			printf("redis connect failed:[%s:%d]\n", "172.16.203.154", 9221);
			return;
		}
		*(redis_list + 1) = redisConnect("172.16.203.122", 9221);
		if (*(redis_list + 1) == NULL) {
			printf("redis connect failed:[%s:%d]\n", "172.16.203.163", 9221);
			return;
		}

		for (unsigned long idx = down; idx < up; idx++ ) {
			snprintf(key_buff, sizeof key_buff, "%lu", idx);
			
			ridx = idx % 2;

			__START__(begin)
			reply = (redisReply*)redisCommand(*(redis_list + ridx), "hmset %s %s %s %s %lu %s %d %s %d %lu %s %lu %s", 
							 key_buff, 
							 "appkey", 
			                 "8e5cf4ca06728a452207f701",
							 "ttl",
			                 idx,
							 "max_count",
			                 5,
							 "num",
							 2,
							 1099985724LU,
							 "stime:1496286894,ctime:1496286894,ttl:1496373294",
							 1101575128LU,
			                 "stime:1496376340,ctime:1496376340,ttl:1496289940");

			if (reply && reply->type == REDIS_REPLY_STATUS && !strcmp("OK", reply->str)) {

			} else {
				printf("hmset faild %s\n", key_buff);
			}
			__FREE__(reply);
			__STOP__(begin, end, "hmset")
		}
		redisFree(*redis_list);
		redisFree(*(redis_list + 1));
		redisFree(*(redis_list + 2));
		free(redis_list);

		return;
	}

	printf("redis cluster connect faild.\n");
}

int
main(int argc, char* argv[])
{
	unsigned long 	key_down, key_up;
	CRedisHelper 	redis_server;
	char 			*host;
	unsigned short 	port;
	time_t 			begin, end;
	char			*type;

	if (argc != 6) {
		printf("redis_ben host port down up\n");
		exit(1);
	}
	
	host = *(argv + 1);
	port = atoi(*(argv + 2));
	key_down = atol(*(argv + 3));
	key_up = atol(*(argv + 4));
	type = *(argv + 5);
	
	redis_server.Connect(host, port, 5);
	if (!redis_server.IsAlived()) {
		printf("redis server connected failed.\n");
		exit(1);
	}
	
	begin = time(NULL);
	if (strcmp(type, "set") == 0) {
		set_test(&redis_server, key_down, key_up);
	} else if (strcmp(type, "get") == 0) {
		get_test(&redis_server, key_down, key_up);
	} else if (strcmp(type, "del") == 0) {
		del_test(&redis_server, key_down, key_up);
	} else if (strcmp(type, "mix") == 0) {
		mix_test(&redis_server, key_down, key_up);
	} else if (strcmp(type, "hash") == 0) {
		hash_test(host, port, key_down, key_up);
	} else if (strcmp(type, "hset") == 0) {
		hash_set(host, port, key_down, key_up);
	} else if (strcmp(type, "msg") == 0) {
		msg_test(key_down, key_up);
	}
	end = time(NULL);
	
	
	//printf("total: %lu elsp: %lu\n", key_up - key_down, end - begin);
	return 0;
}
