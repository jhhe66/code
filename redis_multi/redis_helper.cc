#include "redis_helper.h"
#include <string.h>


int
CRedisHelper::Connect(const char* host, 
					  const unsigned short& port, 
					  const unsigned short second)
{
	 if (context) {
		redisFree(context);
	 }
	 
	 struct timeval tv = {second, 0};
	 
	 context = redisConnectWithTimeout(host, port, tv);
	 
	 return context ? 0 : REDIS_ERR;
}

int
CRedisHelper::Connect()
{
	 if (context) {
		redisFree(context);
	 }
	 
	 context = redisConnectWithTimeout(host_.c_str(), port_, timeout);
	 
	 return context ? 0 : REDIS_ERR;
}

int
CRedisHelper::Set(const string& key, const string& value)
{
	 reply = static_cast<redisReply*>(redisCommand(context, "SET %s %s", key.c_str(), value.c_str()));
	 
	 int result = reply ? reply->integer : REDIS_ERR;
	 
	 SetErrMsg();

	 freeReply();
	 
	 return result;
}

int
CRedisHelper::Get(const string& key, string& value)
{
	 int ret = REDIS_ERR;

	 reply = static_cast<redisReply*>(redisCommand(context, "GET %s", key.c_str()));
	 
	 if (reply->type == REDIS_REPLY_STRING) {
	 	value = reply->str;	
		ret = REDIS_OK;
	 } else if (reply && reply->type == REDIS_REPLY_NIL) {
	 	ret = REDIS_OK;
	 } else {
		 SetErrMsg();
	 }
	 
	 freeReply();
	 
	 return ret;
}

int
CRedisHelper::Set(const char* key, const char* value)
{
	 reply = static_cast<redisReply*>(redisCommand(context, "SET %s %s", key, value));
	 
	 int result = reply ? reply->integer : REDIS_ERR;

	 SetErrMsg();
	 
	 freeReply();
	 
	 return result;
}

int
CRedisHelper::Get(const char* key, char* value, unsigned int& sz)
{
	 int ret = REDIS_ERR;

	 reply = static_cast<redisReply*>(redisCommand(context, "GET %s", key));
	 
	 if (reply && reply->type == REDIS_REPLY_STRING) {
	 	memcpy(value, reply->str, reply->len);
	 	sz = reply->len;
		ret = REDIS_OK;
	 } else if (reply && reply->type == REDIS_REPLY_NIL) {
	 	sz = 0;
	 	ret = REDIS_OK;
	 } else {
		 SetErrMsg();
	 }
	 
	 freeReply();
	 
	 return ret;
}

int
CRedisHelper::SetBin(const char* key, const char* value, unsigned int sz)
{
	reply = static_cast<redisReply*>(redisCommand(context, "SET %s %b", key, value, sz));
	
	int result = reply ? reply->integer : REDIS_ERR;

	freeReply();

	return result;	
}

int
CRedisHelper::GetBin(const char* key, char* value, unsigned int sz)
{
	 int ret = 0;

	 reply = static_cast<redisReply*>(redisCommand(context, "GET %s", key));
	 
	 if (reply->type == REDIS_REPLY_STRING) {
	 	memcpy(value, reply->str, reply->len);
	 	ret = reply->len;
	 }
	 
	 freeReply();
	 
	 return ret;
}

int
CRedisHelper::Enqueue(const string& queue, const string& value)
{
	reply = static_cast<redisReply*>(redisCommand(context, "rpush %s %s", queue.c_str(), value.c_str()));

	int result = reply ? reply->integer : 1; 

	freeReply();

	return result;
}

string&
CRedisHelper::Dequeue(const string& queue, string& value)
{
	reply = static_cast<redisReply*>(redisCommand(context, "lpop %s", queue.c_str()));

	//log_debug("KEY: %s\n", queue.c_str());
	if (reply) {
		//log_debug("reply->type: %d\n", reply->type);

		if (reply->type == 1) {
			value = reply->str;		
		}

		freeReply();
	}

	return value;
}

int
CRedisHelper::Push(const string& stack, const string& value)
{
	reply = static_cast<redisReply*>(redisCommand(context, "rpush %s %s", stack.c_str(), value.c_str()));

	int result = reply ? reply->integer : REDIS_ERR; 

	freeReply();

	return result;
}

string&
CRedisHelper::Pop(const string& stack, string& value)
{
	reply = static_cast<redisReply*>(redisCommand(context, "rpop %s", stack.c_str()));

	value = reply->str;

	freeReply();

	return value;
}

bool
CRedisHelper::ping()
{
	if (!context)
	{
		return false;
	}
	
	reply = static_cast<redisReply*>(redisCommand(context, "ping"));

	bool IsActviced = false;

	if (reply) {
		IsActviced = strcmp("PONG", reply->str) == 0 ? true : false;

		freeReply();
	}

	return IsActviced;
}

string
CRedisHelper::GetHostInfo()
{
	char temp[32];

	snprintf(temp, sizeof temp, "%s:%d", host_.c_str(), port_);
	
	return string(temp);
}

void 
CRedisHelper::SetErrMsg()
{
	if (!reply) return;

	switch (reply->type) {
		case REDIS_REPLY_ERROR:
			_err_msg.assign(reply->str);
			break;
		default:
			break;
	}
}

bool
CRedisHelper::Exists(const char* key)
{
	bool IsExists = false;

	reply = static_cast<redisReply*>(redisCommand(context, "exists %s", key));

	if (reply && reply->type == REDIS_REPLY_INTEGER) {
		IsExists = reply->integer == 1 ? true : false;
	} else {
		SetErrMsg();
	}

	freeReply();

	return IsExists;
}

bool
CRedisHelper::Exists(const string& key)
{
	bool IsExists = false;
	
	reply = static_cast<redisReply*>(redisCommand(context, "exists %s", key.c_str()));

	if (reply && reply->type == REDIS_REPLY_INTEGER) {
		IsExists = reply->integer == 1 ? true : false;
	} else {
		SetErrMsg();
	}

	freeReply();

	return IsExists;
}

int 
CRedisHelper::Del(const string& key) 
{
	int result;

	reply = static_cast<redisReply*>(redisCommand(context, "DEL %s", key.c_str()));
	
	result = reply && reply->type == REDIS_REPLY_INTEGER ? REDIS_OK : REDIS_ERR;

	SetErrMsg();

	freeReply();

	return result;		
}

int
CRedisHelper::Del(const char* key) 
{
	int result;

	reply = static_cast<redisReply*>(redisCommand(context, "DEL %s", key));
	
	result = reply && reply->type == REDIS_REPLY_INTEGER ? REDIS_OK : REDIS_ERR;

	SetErrMsg();

	freeReply();

	return result;	
}