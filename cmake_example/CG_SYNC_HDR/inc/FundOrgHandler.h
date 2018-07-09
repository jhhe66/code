#ifndef __FUND_ORG_HANDLER_H_
#define __FUND_ORG_HANDLER_H_

#include "HandlerBase.h"

namespace Handler
{

class CFundOrgHandler: public CHandlerBase
{
public:
    CFundOrgHandler() {}
    ~CFundOrgHandler() {}
private:
    int SyncToAll();
    string& ToJson(const void* data, string& json); 
};

}


#endif

