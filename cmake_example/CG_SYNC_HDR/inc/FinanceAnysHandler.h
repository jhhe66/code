#ifndef __FINANCE_ANYS_HANDLER_H_
#define __FINANCE_ANYS_HANDLER_H_

#include "HandlerBase.h"

namespace Handler
{
class CFinanceAnysHandler: public CHandlerBase
{
public:
	CFinanceAnysHandler() {}
	~CFinanceAnysHandler() {}
private:
	int SyncToAll();
    string& ToJson(const void* data, string& json);
};

}

#endif
