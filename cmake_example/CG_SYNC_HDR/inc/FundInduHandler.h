#ifndef __FUND_INDU_HANDLER_H_
#define __FUND_INDU_HANDLER_H_

#include "HandlerBase.h"

namespace Handler
{

class CFundInduHandler: public CHandlerBase
{
public:
    CFundInduHandler() {}
    ~CFundInduHandler() {}
private:
    int SyncToAll();
    string& ToJson(const void* data, string& json);
};

}

#endif

