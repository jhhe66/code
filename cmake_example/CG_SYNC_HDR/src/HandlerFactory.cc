#include "HandlerFactory.h"
#include "HandlerBase.h"
#include "AllHandler.h"

namespace Handler
{

CHandlerFactory::CHandlerFactory() 
{
}

CHandlerFactory::~CHandlerFactory() 
{
}

CHandlerBase*
CHandlerFactory::CreateHandler(unsigned short handlerID)
{
    switch(handlerID)
    {
        case 0x0A40:
            return new CZhzdHandler();
            break;
        case 0x0A50:
            return new CFundHandler();
            break;
        case 0x0A66:
            return new CFundDtlHandler();
            break;
        case 0x0A69:
            return new CTrendAnysHandler();
            break;
        case 0x0A70:
            return new CFinanceAnysHandler();
            break;
        case 0x0A71:
            return new CTrendMarketHandler();
            break;
        case 0x0A72:
            return new CFundFinanceHandler();
            break;
        case 0x0A73:
            return new CFundProfitHandler();
            break;
        case 0x0A74:
            return new CFundInduHandler();
            break;
        case 0x0A75:
            return new CFundAluHandler();
            break;
        case 0x0A76:
            return new CFundOrgHandler();
            break;
    }

    return NULL;
}


}

