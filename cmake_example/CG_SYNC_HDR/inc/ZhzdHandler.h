#ifndef __ZHZD_HANDLER_H_
#define __ZHZD_HANDLER_H_

#include "HandlerBase.h"

namespace Handler
{

class CZhzdHandler: public CHandlerBase
{
public:
    CZhzdHandler() {}
    virtual ~CZhzdHandler() {}
private:
    int SyncToAll();
    string& ToJson(const void* data, string& json); 
};

}

#endif
