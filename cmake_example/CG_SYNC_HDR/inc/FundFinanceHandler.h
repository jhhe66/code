#ifndef __FUND_FINANCE_H_
#define __FUND_FINANCE_H_

#include "HandlerBase.h"

namespace Handler
{

class CFundFinanceHandler: public CHandlerBase
{
public:
    CFundFinanceHandler() {}
    ~CFundFinanceHandler() {}
private:
    int SyncToAll();
    string& ToJson(const void* data, string& json);
};

}

#endif
