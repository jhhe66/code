#include "FundProfitHandler.h"
#include "global.h"

#include <json/json.h>
#include <string>

using namespace Json;
using namespace std;

namespace Handler
{

inline string&
CFundProfitHandler::ToJson(const void* data, string& json)
{
    static const unsigned short UNIT = 10000;
    
    Value array, root;

    const JWFundAmentalsTargetsContentData* p = 
        ((const JWFundAmentalsTargetsData*)data)->FundAmentalsTargetsContentData;

    for (unsigned short i = 0; i <  ((const JWFundAmentalsTargetsData*)data)->Num;i++)
    {
        Value v;

        v[UpdateDate] = (p + i)->UpdateDate;
        v[YYSR_TB_ZZL] = ToSigned((p + i)->YYSR_TB_ZZL, UNIT);
        v[JLR_TB_ZZL] = ToSigned((p + i)->JLR_TB_ZZL, UNIT);

        array.append(v);
    }

    root[Code] = ((const JWFundAmentalsTargetsData*)data)->StockIndex;
    root[Num] =  ((const JWFundAmentalsTargetsData*)data)->Num;
    root[Data] = array;

    FastWriter writer;

    json = writer.write(root);

    return json;
}

inline int
CFundProfitHandler::SyncToAll()
{
    client_->SetRecvTimeout(1);
    client_->Connect();
    
    JWFundAmentalsTargetsReq ReqFlow;

    ReqFlow.StockNumber = map_->size();

    for(unsigned short index = 0;index < ReqFlow.StockNumber; index++)
    {
        ReqFlow.StockIndex[index] = map_->GetIndexByOffset(index);
    }

    char ReqBuff[PACK_HEADER_LEN + 1 + sizeof ReqFlow] = {0};

    int ReqLen = 
        JWPackGetFundAmentalsTargetsDataReq(ReqBuff, sizeof ReqBuff, &ReqFlow);

    if(!ReqLen)
        return -1;

    int sLen = 0;

    while (sLen < ReqLen)
    {
        int sVal = client_->Send(ReqBuff + sLen, ReqLen - sLen);

        if(!sVal)
            return -2;

        sLen += sVal;
    }
    
    HexDumpImp(ReqBuff, ReqLen);

    char RspBuff[PACK_HEADER_LEN + 1 + sizeof(JWFundAmentalsTargetsRsp)] = {0};

    int rLen = 0;

    while(JWHandleInput(RspBuff, rLen) == 0)
    {
        int rVal = client_->Recv(RspBuff + rLen, sizeof RspBuff - rLen);

        if (rVal <= 0)
            return -3;

        rLen += rVal;
    }
    
    JWFundAmentalsTargetsRsp RspFlow;

    if (JWUnPackGetFundAmentalsTargetsDataRsp(RspBuff, rLen, &RspFlow) < 0)
        return -4;  
    
    HexDumpImp(RspBuff, rLen);
    
    string json;
    string code;

    for(unsigned short i = 0; i < RspFlow.StockNumber; i++)
    {
        code = map_->GetCodeByIdx(RspFlow.FundAmentalsTargetsData[i].StockIndex);
        
        json = ToJson(&RspFlow.FundAmentalsTargetsData[i], json);
        
        cout << "JSON: " << json << endl;

        mc_->Set(code + MEMC_KEY_LIST[MC_KEY_FUND_PROFIT], json);
    }
    
    return RspFlow.StockNumber;
}


}
