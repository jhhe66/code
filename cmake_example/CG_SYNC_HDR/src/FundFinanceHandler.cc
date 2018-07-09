#include "FundFinanceHandler.h"
#include "global.h"

#include <json/json.h>
#include <string>

using namespace Json;

namespace Handler
{
string&
CFundFinanceHandler::ToJson(const void* data, string& json)
{
    const JWFundAmentalsAnalysisData* p = (const JWFundAmentalsAnalysisData*)data; 

    static const unsigned short UNIT = 1000;

    Value v;

    v[Code] = map_->GetCodeByIdx(p->StockIndex);
    v[UpdateDate] = p->UpdateDate;
    v[YYSR] = ToSigned(p->YYSR, UNIT);
    v[YYSR_TB_ZZ] = ToSigned(p->YYSR_TB_ZZ, UNIT);
    v[JLR_PARENT] = ToSigned(p->JLR_PARENT, UNIT);
    v[JLR_TB_ZZ] = ToSigned(p->JLR_TB_ZZ, UNIT);
    v[MGSY] = ToSigned(p->MGSY, UNIT);

    FastWriter writer;

    json = writer.write(v);

    return json;
}

int
CFundFinanceHandler::SyncToAll()
{
    client_->SetRecvTimeout(1);
    client_->Connect();

    JWFundAmentalsAnalysisReq ReqFlow;

    ReqFlow.StockNumber = map_->size();

    for ( unsigned short index = 0; index < ReqFlow.StockNumber; index++)
    {
        ReqFlow.StockIndex[index] = map_->GetIndexByOffset(index);
    }

    char ReqBuff[PACK_HEADER_LEN + 1 + sizeof ReqFlow] = {0};

    int ReqLen = 
        JWPackGetFundAmentalsAnalysisDataReq(ReqBuff, sizeof ReqBuff, &ReqFlow);

    if ( ReqLen <= 0 )
        return -1;

    int sLen = 0;

    while( sLen < ReqLen)
    {
        int sVal = client_->Send(ReqBuff + sLen, ReqLen - sLen);

        if (sVal <= 0)
            return -2;

        sLen += sVal;
    }

    HexDumpImp(ReqBuff, ReqLen);

    
    char RspBuff[PACK_HEADER_LEN + 1 + sizeof(JWFundAmentalsAnalysisRsp)] = {0};

    int rLen = 0;

    while(JWHandleInput(RspBuff, rLen) == 0)
    {
        int rVal = client_->Recv(RspBuff + rLen, sizeof RspBuff - rLen);

        if (rVal <= 0)
            return -3;

        rLen += rVal;
    }
    
    JWFundAmentalsAnalysisRsp RspFlow;

    int iRet = JWUnPackGetFundAmentalsAnalysisDataRsp(RspBuff, rLen, &RspFlow);

    if (iRet < 0)
        return -4;
    
    HexDumpImp(RspBuff, rLen);
    
    string json;
    string code;

    for(unsigned short i = 0; i < RspFlow.StockNumber;i++)
    {
        code = map_->GetCodeByIdx(RspFlow.FundAmentalsAnalysisData[i].StockIndex);

        json = ToJson(&RspFlow.FundAmentalsAnalysisData[i], json);

        mc_->Set(code + MEMC_KEY_LIST[MC_KEY_FUND_FINANCE], json);
    }

    return RspFlow.StockNumber;
}

}
