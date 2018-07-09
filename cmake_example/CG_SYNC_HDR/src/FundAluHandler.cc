#include "FundAluHandler.h"
#include "global.h"

#include <json/json.h>
#include <string>

using namespace std;
using namespace Json;

namespace Handler
{

inline string& 
CFundAluHandler::ToJson(const void* data, string& json)
{
    const JWFundAluationReserchData* p = (const JWFundAluationReserchData*)data;

    static const unsigned int UNIT = 100000;

    Value v;

    v[Code] = map_->GetCodeByIdx(p->StockIndex); 
    v[UpdateDate] = p->UpdateDate;
    v[PE_TTM] = ToSigned(p->PE_TTM, UNIT);
    v[PE_TTM_HY] = ToSigned(p->PE_TTM_HY, UNIT);
    v[PE_TTM_MARKET] = ToSigned(p->PE_TTM_MARKET, UNIT);
    v[PE_TTM_RANKING] = p->PE_TTM_RANKING;
    v[PB] = ToSigned(p->PB, UNIT);
    v[PB_HY] = ToSigned(p->PB_HY, UNIT);
    v[PB_MARKET] = ToSigned(p->PB_MARKET, UNIT);
    v[PB_RANKING] = p->PB_RANKING;
    v[PEG] = ToSigned(p->PEG, UNIT);
    v[PEG_HY] = ToSigned(p->PEG_HY, UNIT);
    v[PEG_MARKET] = ToSigned(p->PEG_MARKET, UNIT);
    v[PEG_RANKING] = PEG_RANKING;

    FastWriter writer;

    json = writer.write(v);

    return json;
}

int
CFundAluHandler::SyncToAll()
{
    client_->SetRecvTimeout(1);
    client_->Connect();

    JWFundAluationReserchReq ReqFlow;

    ReqFlow.StockNumber = map_->size();

    for (unsigned short index = 0; index < ReqFlow.StockNumber; index++)
    {
        ReqFlow.StockIndex[index] = map_->GetIndexByOffset(index);
    }

    char ReqBuff[PACK_HEADER_LEN + 1 + sizeof ReqFlow] = {0};

    int ReqLen = 
        JWPackGetFundAluationReserchDataReq(ReqBuff, sizeof ReqBuff, &ReqFlow);

    if (!ReqLen)
        return -1;

    int sLen = 0;

    while(sLen < ReqLen)
    {
        int sVal = client_->Send(ReqBuff + sLen, ReqLen - sLen);

        if (!sVal)
            return -2;
        
        sLen += sVal;
    }

    HexDumpImp(ReqBuff, sLen);

    char RspBuff[PACK_HEADER_LEN + 1 + sizeof(JWFundAluationReserchRsp)] = {0};
    
    int rLen = 0;

    while(JWHandleInput(RspBuff, rLen) == 0)
    {
        int rVal = client_->Recv(RspBuff + rLen, sizeof RspBuff - rLen);

        if (!rVal)
            return -3;

        rLen += rVal;
    }
    
    HexDumpImp(RspBuff, rLen);

    JWFundAluationReserchRsp RspFlow;

    if (JWUnPackGetFundAluationReserchDataRsp(RspBuff, rLen, &RspFlow) < 0)
        return -4;
    string json;
    string code;

    for(unsigned short i = 0; i < RspFlow.StockNumber; i++)
    {
        code = map_->GetCodeByIdx(RspFlow.FundAluationRserchData[i].StockIndex);

        json = ToJson(&RspFlow.FundAluationRserchData[i], json);

        cout << "JSON: " << json << endl; 
        
        mc_->Set(code + MEMC_KEY_LIST[MC_KEY_FUND_ALUA], json);
    }
    
    return RspFlow.StockNumber;
}


}

