#include "FundInduHandler.h"
#include "global.h"

#include <json/json.h>
#include <string>

using namespace Json;
using namespace std;

namespace Handler
{

string& 
CFundInduHandler::ToJson(const void* data, string& json)
{
    const JWFundIndustryRanksData* p = (const JWFundIndustryRanksData*)data;

    static const unsigned short UNIT = 1000;
    static const unsigned short UNIT2 = 10000;

    Value v;
    
    v[Code] = map_->GetCodeByIdx(p->StockIndex);
    v[UpdateDate] = p->UpdateDate;
    v[IndustryName] = p->IndustryName;
    v[ZZC] = ToSigned(p->ZZC, UNIT);
    v[ZZC_HY] = ToSigned(p->ZZC_HY, UNIT);
    v[ZZC_PM] = p->ZZC_PM;
    v[JLRL] = ToSigned(p->JLRL, UNIT2);
    v[JLRL_HY] = ToSigned(p->JLRL_HY, UNIT2);
    v[JLRL_PM] = p->JLRL_PM;
    v[JZCSYL] = ToSigned(p->JZCSYL, UNIT2);
    v[JZCSYL_HY] = ToSigned(p->JZCSYL_HY, UNIT2);
    v[JZCSYL_PM] = p->JZCSYL_PM;
    v[XSMLL] = ToSigned(p->XSMLL, UNIT2);
    v[XSMLL_HY] = ToSigned(p->XSMLL_HY, UNIT2);
    v[XSMLL_PM] = p->XSMLL_PM;
    v[InduElementNum] = p->InduElementNum;

    FastWriter writer;

    json = writer.write(v);

    return json; 
}

int
CFundInduHandler::SyncToAll()
{
    client_->SetRecvTimeout(1);
    client_->Connect();
    
    JWFundIndustryRanksReq ReqFlow;

    ReqFlow.StockNumber = map_->size();

    for ( unsigned short index = 0; index < ReqFlow.StockNumber; index++)
    {
        ReqFlow.StockIndex[index] = map_->GetIndexByOffset(index);
    }

    char ReqBuff[PACK_HEADER_LEN + 1 + sizeof ReqFlow] = {0};

    int ReqLen = JWPackGetFundIndustryRanksDataReq(ReqBuff, sizeof ReqBuff, &ReqFlow);

    if (!ReqLen)
        return -1;

    int sLen = 0;

    while( sLen < ReqLen)
    {
        int sVal = client_->Send(ReqBuff + sLen, ReqLen - sLen);

        if(!sVal)
            return -2;

        sLen += sVal;
    }

    HexDumpImp(ReqBuff, ReqLen);

    char RspBuff[PACK_HEADER_LEN + 1 + sizeof(JWFundIndustryRanksRsp)] = {0};

    int rLen = 0;

    while(JWHandleInput(RspBuff, rLen) == 0)
    {
        int rVal = client_->Recv(RspBuff + rLen, sizeof RspBuff - rLen);

        if (!rVal)
            return -3;

        rLen += rVal;
    }

    HexDumpImp(RspBuff, rLen);

    JWFundIndustryRanksRsp RspFlow;

    if(JWUnPackGetFundIndustryRanksDataRsp(RspBuff, rLen, &RspFlow) < 0)
        return -4;
    
    string json;
    string code;

    for ( unsigned short i = 0; i < RspFlow.StockNumber; i++)
    {
        code = map_->GetCodeByIdx(RspFlow.FundIndustryRanksData[i].StockIndex);

        json = ToJson(&RspFlow.FundIndustryRanksData[i], json);
        
        cout << "JSON: " << json << endl; 

        mc_->Set(code + MEMC_KEY_LIST[MC_KEY_FUND_INDU_RANK], json);
    }

    return RspFlow.StockNumber;
}

}

