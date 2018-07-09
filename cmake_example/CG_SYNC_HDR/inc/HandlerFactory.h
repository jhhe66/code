#ifndef __HANDLER_FACTORY_H_
#define __HANDLER_FACTORY_H_

namespace Handler
{

class CHandlerBase;

class CHandlerFactory
{
public:
    CHandlerFactory();
    virtual ~CHandlerFactory();

    CHandlerBase* CreateHandler(unsigned short handlerID);
};

}


#endif
