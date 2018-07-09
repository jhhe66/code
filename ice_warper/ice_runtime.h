#ifndef __ICE_RUNTIME_H_
#define __ICE_RUNTIME_H_

#include "QueryUserApi.h"
#include <Ice/Ice.h>
#include <IceUtil/IceUtil.h>

extern QueryUserIPrx 
ice_new_proxy(const char* ice_conf);


#endif

