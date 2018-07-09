#include "ZhzdHandler.h"
#include "global.h"

#include <string>
#include <json/json.h>
#include <time.h>
#include <sys/time.h>
#include <sstream>
#include <iomanip>

using namespace std;
using namespace Json;

namespace Handler
{

// date format of china
inline string GetDateZH()
{
	struct tm* ptr;
	time_t t;
	
	t = time(&t);
	
	ptr = localtime(&t);
	
	ostringstream oss;
	
	oss << (ptr->tm_year % 100) << "年" << ptr->tm_mon + 1 << "月" 
		<< ptr->tm_mday << "日 " << ptr->tm_hour << ":" << setw(2) << setfill('0') << ptr->tm_min << ":" << setw(2) << setfill('0') << ptr->tm_sec;
	
	return oss.str();
}


inline string&
CZhzdHandler::ToJson(const void* data, string& json)
{
    return json;
}

inline int
CZhzdHandler::SyncToAll()
{
    client_->SetRecvTimeout(1);
    client_->Connect();

    char ReqBuff[PACK_HEADER_LEN + 1 + sizeof(unsigned short)] = {0};

    for( unsigned short index = 0; index < map_->size(); index++)
    {
        unsigned short StockIndex = map_->GetIndexByOffset(index);

        int ReqLen = JWPackGetDiagnoseDataReq(ReqBuff, sizeof ReqBuff, StockIndex);

        if( !ReqLen)
            return -1;

        int sLen = client_->Send(ReqBuff, ReqLen);

        if (!sLen)
            return -2;
        
        char RspBuff[10 * 1024] = {0};
        
        int rLen = 0;

        while(JWHandleInput(RspBuff, rLen) == 0)
        {
            int rVal = client_->Recv(RspBuff + rLen, sizeof RspBuff - rLen);

            if (!rVal)
                return -3;

            rLen += rVal;
        }

        HexDumpImp(RspBuff, rLen);

        JWDiagnoseDataRsp RspFlow;
        
        JWUnPackGetDiagnoseDataRsp(RspBuff, rLen, &RspFlow);

        //if ((JWUnPackGetDiagnoseDataRsp(RspBuff, rLen, &RspFlow)) < 0)
        //    return -4;

        mc_->Set( map_->GetCodeByIdx(StockIndex) + MEMC_KEY_LIST[MC_KEY_ZHHS], 
                string(RspFlow.Content).append( "(" + GetDateZH() + ")"));

        memset(ReqBuff, 0, sizeof ReqBuff);
    }

    return map_->size();
}

}
