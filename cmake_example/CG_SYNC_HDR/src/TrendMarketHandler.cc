#include "TrendMarketHandler.h"
#include "global.h"

#include <json/json.h>
#include <string>

using namespace std;
using namespace Json;

namespace Handler
{

inline string&
CTrendMarketHandler::ToJson(const void* data, string& json)
{
    const JWTechniccalAndMarketData* p = (const JWTechniccalAndMarketData*)data;

    Value v;
    
    static const unsigned short UNIT = 10000;

    v[Code] = map_->GetCodeByIdx(p->StockIndex);
    v[UpdateDate] = p->UpdateDate;
    v[AccumulativeMarkup] = ToSigned(p->AccumulativeMarkup, UNIT);
    v[AccumulativeMarkupHY] = ToSigned(p->AccumulativeMarkupHY, UNIT);
    v[AccumulativeMarkupSHZS] = ToSigned(p->AccumulativeMarkupSHZS, UNIT);
    v[TurnOverRate1] = ToSigned(p->TurnOverRate1, UNIT);
    v[TurnOverRate3] = ToSigned(p->TurnOverRate3, UNIT);
    v[TurnOverRate5] = ToSigned(p->TurnOverRate5, UNIT);
    
    FastWriter writer;

    json = writer.write(v);

    return json;
}

inline int
CTrendMarketHandler::SyncToAll()
{
    client_->SetRecvTimeout(1);
    client_->Connect();
    
    JWTechniccalAndMarketReq ReqFlow;
    
    ReqFlow.StockNumber = map_->size();

    for ( unsigned short index = 0; index < ReqFlow.StockNumber; index++)
    {
        ReqFlow.StockIndex[index] = map_->GetIndexByOffset(index);
    }

    char ReqBuff[PACK_HEADER_LEN + 1 + sizeof ReqFlow] = {0};

    int ReqLen = 
        JWPackGetTechniccalAndMarketDataReq(ReqBuff, sizeof ReqBuff, &ReqFlow);

    if ( ReqLen <= 0)
        return -1;

    int sLen = 0;

    while (sLen < ReqLen)
    {
        int sVal = client_->Send(ReqBuff + sLen, ReqLen - sLen);

        if (sVal <= 0)
            return -2;

        sLen += sVal;
    }

    HexDumpImp(ReqBuff, sLen);
    
    char RspBuff[PACK_HEADER_LEN + 1 + sizeof(JWTechniccalAndMarketRsp)] = {0};
    
    int rLen = 0;

    while(JWHandleInput(RspBuff, rLen) == 0)
    {
        int rVal = client_->Recv(RspBuff + rLen, sizeof RspBuff - rLen);
        
        if ( rVal <= 0)
            return -3;

        rLen += rVal;
    }

    HexDumpImp(RspBuff, rLen);
    
    JWTechniccalAndMarketRsp RspFlow;
    
    if(JWUnPackGetTechniccalAndMarketDataRsp(RspBuff, rLen, &RspFlow) < 0)
        return -4;
    
    string json;
    string code;

    for (unsigned short i = 0; i < RspFlow.StockNumber; i++)
    {
        code = map_->GetCodeByIdx(RspFlow.TechniccalAndMarketData[i].StockIndex);

        json = ToJson(&RspFlow.TechniccalAndMarketData[i], json);
        
        cout << "JSON: " << json << endl;

        mc_->Set(code + MEMC_KEY_LIST[MC_KEY_TREND_MARKET], json);
    }

    return RspFlow.StockNumber;
}

}

