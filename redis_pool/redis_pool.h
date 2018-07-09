#ifndef __REDIS_POOL_H_
#define __REDIS_POOL_H_

#include <set>
#include <string>

class JCRedisPool {
public:
	JCRedisPool();
	~JCRedisPool();
public:
	void* GetObject();
	void Create();
	size_t UsedSize() { return _useds.size(); }
	size_t FreeSize() { return _frees.size(); }
protected:
	void __Alloc();
	void __Free();
	void __AddObject();
private:
	std::set<void*> _useds;
	std::set<void*> _frees;
	std::string		_rid;
	std::string 	_host;
	unsigned short 	_port;
}

#endif


