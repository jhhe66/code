#include "redis_cluster.h"
#include "RedisHelper.h"
#include "xxhash.h"
#include <string>
#include <string.h>

CRedisCluster::CRedisCluster():_redis_nums(0)
{
}

CRedisCluster::~CRedisCluster()
{
    for (std::vector<CRedisHelper*>::iterator it; it != _redis_list.end(); it++) {
        delete *it;
    }
}

int 
CRedisCluster::Init(redis_conf_t* redis_confs, unsigned int nums)
{
    for (unsigned int idx = 0; idx < nums; idx++) {
        CRedisHelper *redis = new CRedisHelper((redis_confs + idx)->ip, (redis_confs + idx)->port);
        if (redis) {
            _redis_list.push_back(redis);
        } else {
            return -1; // error
        }
    }

    _redis_nums = _redis_list.size();

    return 0;
}

int
CRedisCluster::Set(const std::string& key, const std::string& value)
{
    int redis_idx;

    redis_idx = _get_idx(key);
    return _set(key, value, redis_idx);
}

int
CRedisCluster::Get(const std::string& key, std::string& value)
{
    int redis_idx;

    redis_idx = _get_idx(key);
    return _get(key, value, redis_idx);
}

int
CRedisCluster::Set(const char* key, const char* value)
{
    int redis_idx;

    redis_idx = _get_idx(key);
    return _set(key, value, redis_idx);
}

int 
CRedisCluster::Get(const char* key, char* value, unsigned int& sz)
{
    int redis_idx;

    redis_idx = _get_idx(key);
    return _get(key, value, sz, redis_idx);
}

int 
CRedisCluster::_get_idx(const std::string& key)
{
    int idx, hash;

    hash = XXH64(key.c_str(), key.size(), 0);

    idx = hash % _redis_nums;

    return idx;
}

int 
CRedisCluster::_get_idx(const char* key)
{
	int idx, hash;

	hash = XXH64(key, strlen(key), 0);

	idx = hash % _redis_nums;

	return idx;
}

int 
CRedisCluster::_set(const std::string& key, 
                    const std::string& value,
                    unsigned short idx)
{
    CRedisHelper *redis;

    redis = _redis_list[idx];
    if (redis) {
        return redis->Set(key, value);
    }

    return -1;
}

int
CRedisCluster::_get(const std::string& key,
                    std::string& value,
                    unsigned short idx)
{
    CRedisHelper *redis;

    redis = _redis_list[idx];
    if (redis) {
		redis->Get(key, value);
		return 0;
    }

    return -1;
}

int 
CRedisCluster::_set(const char* key,
                    const char* value,
                    unsigned short idx)
{
    CRedisHelper *redis;

    redis = _redis_list[idx];
    if (redis) {
        return redis->Set(key, value);
    }

    return -1;
}

int 
CRedisCluster::_get(const char* key,
                    char* value,
                    unsigned int& sz,
					unsigned short idx)
{
    CRedisHelper *redis;

    redis = _redis_list[idx];
    if (redis) {
        return redis->Get(key, value, sz);
    }

    return -1;
}
