#ifndef __TREND_ANYS_HANDLER_H_
#define __TREND_ANYS_HANDLER_H_

#include "HandlerBase.h"

namespace Handler
{   

class CTrendAnysHandler: public CHandlerBase
{
public:
    CTrendAnysHandler() {}
    ~CTrendAnysHandler() {}
private:
    int SyncToAll();
    string& ToJson(const void* data, string& json);
};


}

#endif

