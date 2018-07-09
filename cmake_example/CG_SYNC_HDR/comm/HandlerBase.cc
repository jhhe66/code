#include "HandlerBase.h"
#include "global.h"

using namespace std;

namespace Handler
{

CHandlerBase::CHandlerBase()
    :client_(new CTcpClient(HOST, PORT)),
    map_(new CShmMap(SHM_KEY)),
    mc_(new CMemcachedHandler(MEMC_HOST))
{
    map_->init();
    map_->Restore();
}

CHandlerBase::~CHandlerBase() {}

int
CHandlerBase::Sync()
{
    return SyncToAll();
}

}

