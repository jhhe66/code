#ifndef __FUND_PROFIT_HANDLER_H_
#define __FUND_PROFIT_HANDLER_H_

#include "HandlerBase.h"

namespace Handler
{

class CFundProfitHandler: public CHandlerBase
{
public:
    CFundProfitHandler() {}
    ~CFundProfitHandler() {}
private:
    int SyncToAll();
    string& ToJson(const void* data, string& json);
};

}

#endif
