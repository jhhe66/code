#ifndef __TREND_HANDLER_H_
#define __TREND_HANDLER_H_

#include "HandlerBase.h"

namespace Handler
{

class CTrendMarketHandler: public CHandlerBase
{
public:
    CTrendMarketHandler() {}
    ~CTrendMarketHandler() {}
private:
    int SyncToAll();
    string& ToJson(const void* data, string& json);
};

}


#endif


