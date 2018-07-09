#include "ice_warper.h"
#include "QueryUserApi.h"
#include <Ice/Ice.h>
#include <IceUtil/IceUtil.h>
#include <stdlib.h>
#include <stdio.h>
#include <string>
#include <assert.h>

using namespace QueryUserApi;
using namespace Ice;
using std::string;

#define ICE_PRE(name) ice_warper_##name

typedef struct ice_handle_s ice_handle_t;
struct ice_handle_s {
	//QueryUserIPrx 	_prx;
	ObjectPrx		_prx;
	CommunicatorPtr _ic;
};

ice_warper_t
ICE_PRE(init)(const char* ice_conf, const char* proxy)
{
	ice_handle_t 	*handle = NULL;
	Ice::ObjectPrx 	base;
	char*			argv[] = {const_cast<char*>(ice_conf)};
	int 			argc = 1;
	
	printf("conf: %s\n", argv[0]);
	handle = new ice_handle_t;
	if (handle) {
		try {
			handle->_ic = Ice::initialize(argc, argv);
			base = handle->_ic->stringToProxy(proxy);
			//handle->_prx = QueryUserIPrx::checkedCast(base);
		} catch (const Ice::Exception& e) {
			printf("error: %s\n", e.what());
			delete handle;
			handle = NULL;
		}
	}

	return (ice_warper_t)handle;
}

void
ICE_PRE(free)(ice_warper_t ice_warper)
{
	delete (ice_handle_t*)ice_warper;
}

int
ICE_PRE(valid_appkeys)(ice_warper_t ice_warper, 
					   ice_request_t req, 
					   ice_response_t res, 
					   unsigned int* rsz)
{
	
	ice_handle_t 	*handle;
	bytes 			resp;
	string 			str_res;
	string 			str_req;
	unsigned int	res_sz;
	QueryUserIPrx   prx;
	
	res_sz = *rsz;
	str_req.assign(req, strlen(req));
	printf("req: %s\n", req);

	handle = (ice_handle_t*)ice_warper;
	assert(handle != NULL);
	if (handle) {
		prx = QueryUserIPrx::checkedCast(handle->_prx);
		Ice::AsyncResultPtr r = prx->begin_ValidAppkeys(::QueryUserApi::bytes(req, req + strlen(req)));
		prx->end_ValidAppkeys(resp, r);
		str_res.assign(resp.begin(), resp.end());
		printf("str_res: %zu\n", str_res.size());
		//snprintf(res, res_sz, "%s", str_res.c_str());
		memcpy(res, str_res.c_str(), str_res.size());
		*rsz = str_res.size();
	}
	
	return 0;
}

int
ICE_PRE(valid_users)(ice_warper_t ice_warper, 
					 ice_request_t req, 
					 ice_response_t res, 
					 unsigned int* rsz)
{
	
	ice_handle_t 	*handle;
	bytes 			resp;
	string 			str_res;
	string 			str_req;
	unsigned int	res_sz;
	QueryUserIPrx   prx;
	
	res_sz = *rsz;
	str_req.assign(req, strlen(req));
	printf("req: %s\n", req);

	handle = (ice_handle_t*)ice_warper;
	assert(handle != NULL);
	if (handle) {
		prx = QueryUserIPrx::checkedCast(handle->_prx);
		Ice::AsyncResultPtr r = prx->begin_ValidUsers(::QueryUserApi::bytes(req, req + strlen(req)));
		prx->end_ValidUsers(resp, r);
		str_res.assign(resp.begin(), resp.end());
		printf("str_res: %zu\n", str_res.size());
		memcpy(res, str_res.c_str(), str_res.size());
		//snprintf(res, res_sz, "%s", str_res.c_str());
		*rsz = str_res.size();
	}
	
	return 0;
}
