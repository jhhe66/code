#ifndef __HANDLER_H_
#define __HANDLER_H_

#include <boost/scoped_ptr.hpp>
#include <string>

using namespace std;

class CTcpClient;

class CShmMap;

namespace Handler
{

class CMemcachedHandler;

class CHandlerBase
{
public:
    CHandlerBase();
    virtual ~CHandlerBase();

    int Sync();
protected:
    boost::scoped_ptr<CTcpClient> client_;
    boost::scoped_ptr<CShmMap> map_;
    boost::scoped_ptr<CMemcachedHandler> mc_;
private:
    virtual int SyncToAll() =0;
    virtual string& ToJson(const void* data, string& json) =0 ;
};

}


#endif
