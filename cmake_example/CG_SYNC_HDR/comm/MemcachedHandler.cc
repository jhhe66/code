#include "MemcachedHandler.h"

namespace Handler
{

const unsigned short USER_FLAG = 0;// C# 客户端需要设置该值才能访问

CMemcachedHandler::CMemcachedHandler(const string & host)
	:host_(host)
{

	memcached_server_st* server_ = memcached_servers_parse(host_.c_str());
	memc_ = memcached_create(NULL);
	memcached_server_push(memc_, server_);
	memcached_server_list_free(server_);
}

CMemcachedHandler::~CMemcachedHandler()
{
	memcached_free(memc_);
}

string
CMemcachedHandler::Get(string& key)
{
	size_t vlen;
	uint32_t flags;
	memcached_return rc;
	
	char* result = memcached_get(memc_, key.c_str(), key.size(), &vlen, &flags, &rc);

	string value = result ? result : "";

    if(result != NULL)
		free(result);
	

	if (rc == MEMCACHED_SUCCESS)
	{
		return value;//string(result);
	}
	else
	{
		return string("");
	}
}

int
CMemcachedHandler::Set(const string & key,const string & value)
{
	memcached_return rc;

	rc = memcached_set(memc_, key.c_str(), key.size(), value.c_str(), value.size(), 0, USER_FLAG);// C# 客户端需要设置该值才能访问

	if (rc == MEMCACHED_SUCCESS)
		return 0;
	else
		return -1;
}

}
