#ifndef __FUND_HANDLER_H_
#define __FUND_HANDLER_H_

#include "HandlerBase.h"

namespace Handler
{

class CFundHandler: public CHandlerBase
{
public:
    enum FundType
	{
		Today=0,
		ThreeDays,
		FiveDays,
		TenDays,
		TwentyDays
	};

    CFundHandler() {}
    virtual ~CFundHandler() {}
private:
    int SyncToAll();
    string& ToJson(const void* data, string& json);   
};

}
#endif
