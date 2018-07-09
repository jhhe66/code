#ifndef __REDIS_HELPER_H_
#define __REDIS_HELPER_H_

#include <string>
#include <vector>
#include "hiredis/hiredis.h"


using namespace std;

class CRedisHelper {
public:
	CRedisHelper(const string& host, 
				 unsigned short port, 
				 unsigned short second = 5)
		:reply(0), 
		host_(host), 
		port_(port)
	{
		timeout.tv_sec = second;
		timeout.tv_usec = 0;
		
		context = redisConnectWithTimeout(host_.c_str(), port_, timeout);
	}

	CRedisHelper():context(NULL),
				   reply(NULL) {}
	
	~CRedisHelper() 
	{
		if(reply) freeReply();
		if(context) redisFree(context); 
	}
	
public:

	int Connect(const char* host, 
				const unsigned short& port, 
				const unsigned short second = 5);
	int Connect();

	int Set(const string& key, const string& value);
	int Get(const string& key, string& value);

	int Set(const char* key, const char* value);
	int Get(const char* key, char* value, unsigned int& sz);

	int SetBin(const char* key, const char* value, unsigned int sz);
	int GetBin(const char* key, char* value, unsigned int sz);

	bool Exists(const char* key);
	bool Exists(const string& key);

	int Enqueue(const string& queue, const string& value);
	string& Dequeue(const string& queue, string& value);

	int Push(const string& stack, const string& value);
	string& Pop(const string& stack, string& value);

	int Del(const string& key);
	int Del(const char* key);

	bool IsAlived() { return ping(); }

	string GetHostInfo();
	string ErrMsg() { return _err_msg; }

private:
	redisContext	*context;
	redisReply		*reply;
	struct timeval 	timeout;

	string 			host_;
	unsigned short 	port_;
	string 			_err_msg;

private:
	void SetErrMsg();

	void freeReply() 
	{ 
		if(reply) 
		{
			freeReplyObject(reply);
			reply = NULL;
		}
	}

	CRedisHelper(const CRedisHelper& rhs) {}  //����copy����
	CRedisHelper& operator=(const CRedisHelper& rsh) { return *this; } // ���ܸ�ֵ

	bool ping();
};

#endif

