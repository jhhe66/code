#include <hiredis/hiredis.h>
#include <hiredis/async.h>
#include <hiredis/adapters/libevent.h>
#include <event2/event.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

static unsigned int need_cnt = 0;
static unsigned int done_cnt = 0;

static void
__getCB(redisAsyncContext *c, void *r, void *privdata)
{
}

static void
__setCB(redisAsyncContext *c, void *r, void *privdata)
{
#if 0
	redisReply *reply = r;
	if (reply == NULL) {
		printf("set[%u] op failed.\n", done_cnt);
		return;
	} else {
		switch (reply->type) {
			case REDIS_REPLY_INTEGER:
				printf("set op [%d] \n", reply->integer);
				break;
			case REDIS_REPLY_STRING:
				printf("set op [%s] \n", reply->str);
				break;
			case REDIS_REPLY_STATUS:
				printf("set op [%d] \n", reply->integer);
				break;
			
			default:
				printf("unknow error.\n");
				break;
		} }
#endif

	printf("%s %u\n", __PRETTY_FUNCTION__, done_cnt);

#if 1
	if (++done_cnt >= need_cnt) {
		//if (c->ev.cleanup) c->ev.cleanup(c->ev.data);	
		struct event_base *base = (struct event_base*)privdata;	
		event_base_loopbreak(base);
		printf("%s %u\n", __PRETTY_FUNCTION__, done_cnt);
	}

	printf("data: %s\n", (char*)(c->data));
#endif
}

static void
__connectCB(const struct redisAsyncContext* c, int status) {
	if (status != REDIS_OK) {
		printf("error: %s\n", c->errstr);
		return;
	}

	//if (c->ev.cleanup) c->ev.cleanup(c->ev.data);	
	printf("connected redis succ.\n");
}

static void
__disconnectCB(const struct redisAsyncContext* c, int status)
{
	if (status != REDIS_OK) {
		printf("error: %s\n", c->errstr);
		return;
	}

	printf("disconnectd redis.\n");
}

int
main(int argc, char** argv)
{
	char 				*host;
	unsigned short 		port;
	char 				data[32];
	unsigned int 		loop = 0;
	struct event_base 	*base = NULL;
	
	
	
	if (argc < 4) {
		printf("usage: host port");
		return -1;
	}

	host = *(argv + 1);
	port = atoi(*(argv + 2));
	loop = (unsigned int)atoi(*(argv + 3));

	base = event_base_new();

	redisAsyncContext *c = redisAsyncConnect(host, port);	
	if (c->err) {
		printf("error: %s\n", c->errstr);
		return -1;
	}
	
	snprintf(data, sizeof data, "%s", "hello");	
	c->data = data;

	redisLibeventAttach(c, base);
	redisAsyncSetConnectCallback(c, __connectCB);
	event_base_loop(base, EVLOOP_ONCE);

	printf("start op...\n");	
	redisLibeventAttach(c, base);
	need_cnt = done_cnt = 0;
	for (unsigned int idx = 0; idx < loop; idx++) {
		redisAsyncCommand(c, __setCB, base, "set libev_%d %d", idx, idx);
		need_cnt++;
	}

	redisAsyncSetDisconnectCallback(c, __disconnectCB);
	event_base_loop(base, 0);

	return 0;
}
