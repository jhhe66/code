#ifndef __FUND_ALU_HANDLER_H_
#define __FUND_ALU_HANDLER_H_

#include "HandlerBase.h"

namespace Handler
{

class CFundAluHandler: public CHandlerBase
{
public:
    CFundAluHandler() {}
    ~CFundAluHandler() {}
private:    
    int SyncToAll();
    string& ToJson(const void* data, string& json);
};

}

#endif
