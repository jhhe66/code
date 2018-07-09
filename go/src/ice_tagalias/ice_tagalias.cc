#include "ice_tagalias.h"
#include "tagAlias.h"
#include <Ice/Ice.h>
#include <IceUtil/IceUtil.h>
#include <string>
#include <stdio.h>
#include <stdlib.h>

using namespace std;
using namespace TagAlias;
using namespace Ice;


#ifdef __DEBUG__
static void
hex_dump(const char* buff, size_t sz)
{
    printf("[");
	for (size_t idx = 0; idx < sz; idx++) {
		printf("%u ", (unsigned char)*(buff + idx));
	}
	printf("]\n");
}
#endif

#define ICE_PRE(name) ice_tagalias_##name 
typedef struct ice_handle_s ice_handle_t;
struct ice_handle_s {
	ObjectPrx 		_prx;
	CommunicatorPtr _ic;
};


ice_warper_t
ICE_PRE(init)(const char* ice_conf, const char* proxy)
{
	ice_handle_t 	*handle = NULL;
	ObjectPrx 		base;
	char*			argv[] = {const_cast<char*>(ice_conf)};
	int 			argc = 1;
	
	printf("conf: %s\n", argv[0]);
	handle = new ice_handle_t;
	if (handle) {
		try {
			handle->_ic = Ice::initialize(argc, argv);
			base = handle->_ic->stringToProxy(proxy);
			//handle->_prx = QueryUserIPrx::checkedCast(base);
			handle->_prx = base;
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
ICE_PRE(request)(ice_warper_t ice_warper,
				const ice_request_t req,	
				ice_response_t res,
				unsigned int* sz)
{
	ice_handle_t 	*handle;
	bytes 			resp;
	string			str_req, str_res;
	unsigned int 	res_sz;
	TagAliasOpPrx	prx;

	res_sz = *sz;
	str_req.assign(req, strlen(req));
	printf("req: %s\n", req);

	handle = (ice_handle_t*)ice_warper;
	assert(handle);	
	
	if (handle) {
		prx = TagAliasOpPrx::checkedCast(handle->_prx);
		str_req.assign(req, strlen(req));
		::Ice::AsyncResultPtr r = prx->begin_request(str_req);
		prx->end_request(str_res, r);
#ifdef __DEBUG__
		hex_dump(str_res.c_str(), str_res.size());
#endif
		printf("str_res: %zu\n", str_res.size());
		memcpy(res, str_res.c_str(), str_res.size());
		*sz = str_res.size();
	}

	return 0;
}

int
ICE_PRE(validateTags)(ice_warper_t ice_warper,
				const ice_request_t req,	
				ice_response_t res,
				unsigned int* sz)
{
	ice_handle_t 	*handle;
	bytes 			resp;
	string			str_req, str_res;
	unsigned int 	res_sz;
	TagAliasOpPrx	prx;

	res_sz = *sz;
	str_req.assign(req, strlen(req));
	printf("req: %s\n", req);

	handle = (ice_handle_t*)ice_warper;
	assert(handle);	
	
	if (handle) {
		prx = TagAliasOpPrx::checkedCast(handle->_prx);
		str_req.assign(req, strlen(req));
		::Ice::AsyncResultPtr r = prx->begin_validateTags(::TagAlias::bytes(req, req + strlen(req)));
		prx->end_validateTags(resp, r);
		str_res.assign(resp.begin(), resp.end());
		printf("str_res: %zu\n", str_res.size());
		snprintf(res, res_sz, "%s", str_res.c_str());
		*sz = str_res.size();
	}

	return 0;	
}

