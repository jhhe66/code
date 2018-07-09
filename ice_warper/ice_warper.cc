#include "ice_warper.h"
#include "QueryUserApi.h"
#include <Ice/Ice.h>
#include <IceUtil/IceUtil.h>
#include <stdlib.h>

using namespace QueryUserApi;
using namespace Ice;

#define ICE_PRE(name) ice_warper_##name

typedef struct ice_handle_s ice_handle_t;
struct ice_handle_s {
	QueryUserIPrx 	_prx;
	CommunicatorPtr _ic;
};

ice_warper_t
ICE_PRE(init)(const char* ice_conf)
{
	ice_handle_t 	*handle;
	Ice::ObjectPrx 	base;
	char*			argv[] = {const_cast<char*>(ice_conf)};
	int 			argc = 1;
	
	handle = new ice_handle_t;
	if (handle) {
		try {
			handle->_ic = Ice::initialize(argc, argv);
			base = handle->_ic->stringToProxy("ValidUsersProxy");
			handle->_prx = QueryUserIPrx::checkedCast(base);
		} catch (const Ice::Exception& e) {
			delete handle;
			handle = NULL;

			return (void*)handle;
		}
	}

	return (void*)handle;
}

void
ICE_PRE(free)(ice_warper_t ice_warper)
{
	delete (ice_handle_t*)ice_warper;
}

int
ICE_PRE(valid_appkeys)(ice_warper_t ice_warper, const char* appkey, unsigned long uid)
{
	return uid;
}
