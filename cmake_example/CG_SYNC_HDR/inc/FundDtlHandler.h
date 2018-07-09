#ifndef __FUND_DTL_HANDLER_H_
#define __FUND_DTL_HANDLER_H_

#include "HandlerBase.h"

namespace Handler
{

class CFundDtlHandler: public CHandlerBase
{
public:
    CFundDtlHandler() {}
    ~CFundDtlHandler() {}
private:
    int SyncToAll();
    string& ToJson(const void* data, string& json);
};

}



#endif
