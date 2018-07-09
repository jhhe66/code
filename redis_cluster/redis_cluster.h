#ifndef __REDIS_CLUSTER_H_
#define __REDIS_CLUSTER_H_

#include <vector>
#include <string>

class CRedisHelper;

typedef struct redis_conf_s redis_conf_t;
struct redis_conf_s {
	char 			ip[16];
	unsigned int 	port;
};

class CRedisCluster {
public:
	CRedisCluster();
	~CRedisCluster();
public:
	int Init(redis_conf_t* redis_confs, unsigned int nums); /* create redis instance */
	
	int Set(const std::string& key, const std::string& value);
	int Set(const char* key, const char* value);
	int Get(const std::string& key, std::string& value);
	int Get(const char* key, char* value, unsigned int& sz);

private:
	int _set(const std::string& key, const std::string& value, unsigned short idx);
	int _get(const std::string& key, std::string& value, unsigned short idx);

	int _set(const char* key, const char* value, unsigned short idx);
	int _get(const char* key, char* value, unsigned int& sz, unsigned short idx); 

	int _get_idx(const std::string& key);
	int _get_idx(const char* key);
private:
	std::vector<CRedisHelper*> 	_redis_list;
	unsigned int 				_redis_nums;
};

#endif 
