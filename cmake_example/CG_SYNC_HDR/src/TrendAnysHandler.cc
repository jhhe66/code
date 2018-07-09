#include "TrendAnysHandler.h"
#include "global.h"

#include <iostream>
#include <json/json.h>

using namespace Json;
using namespace std;

namespace Handler
{

string&
CTrendAnysHandler::ToJson(const void* data, string& json)
{
    const JWTechniccalTrendAnaysisData* p = (const JWTechniccalTrendAnaysisData*)data;
    //const JWTechniccalTrendAnaysisData* p = static_cast<const JWTechniccalTrendAnaysisData*>(p);

    Value v;
    
    static const unsigned short UNIT = 1000;
    
    cout << "index: " << p->StockIndex;

    v[Code] = map_->GetCodeByIdx(p->StockIndex);;
    v[UpdateDate] = p->UpdateDate;
    v[ShortSupportPrice] = ToSigned(p->ShortSupportPrice, UNIT);
    v[ShortResistancePrice] = ToSigned(p->ShortResistancePrice, UNIT);
    v[MidSupportPrice] = ToSigned(p->MidSupportPrice, UNIT);
    v[MidResistancePrice] = ToSigned(p->MidResistancePrice, UNIT);
    v[Glossary] = p->Glossary;
    
    FastWriter writer;

    json = writer.write(v);

    return json;
}

int
CTrendAnysHandler::SyncToAll()
{
    client_->SetRecvTimeout(1);
    client_->Connect();

    JWTechniccalTrendAnaysisReq ReqFlow;

    ReqFlow.StockNumber = map_->size();

    for ( unsigned short index = 0; index < ReqFlow.StockNumber; index++ )
    {
        ReqFlow.StockIndex[index] = map_->GetIndexByOffset(index); 
    }

    char ReqBuff[PACK_HEADER_LEN + 1 + sizeof(ReqFlow)] = {0};

    int ReqLen = 
        JWPackGetTechniccalTrendAnaysisDataReq(ReqBuff, sizeof(ReqFlow), &ReqFlow);

    if ( ReqLen <= 0)
        return -1;

    int sLen = 0;

    while(sLen < ReqLen)
    {
        int sVal = client_->Send(ReqBuff + sLen, ReqLen - sLen);
        
        if ( sVal <= 0)
            return -2;

        sLen += sVal;
    }
    
    HexDumpImp(ReqBuff, sLen);
    

    cout << "Recving..." << endl;

    char RspBuff[PACK_HEADER_LEN + 1 + sizeof(JWTechniccalTrendAnaysisRsp)] = {0};

    int rLen = 0;

    while(JWHandleInput(RspBuff, rLen) == 0)
    {
        int rVal = client_->Recv(RspBuff + rLen, sizeof RspBuff - rLen);

        if (rVal <= 0)
        {
            cout << "rVal: " << rVal << endl;
            cout << "rLen: " << rLen << endl;

            return -3;
        }
        rLen += rVal;
    }

    HexDumpImp(RspBuff, rLen);
    
    JWTechniccalTrendAnaysisRsp RspFlow;

    int iRet = JWUnPackTechniccalTrendAnaysisDataRsp(RspBuff, rLen, &RspFlow);
    
    if (iRet < 0)
       return -4;
    
    string json("");
    string code("");

    for ( unsigned short i = 0; i < RspFlow.StockNumber; i++)
    {
        code = map_->GetCodeByIdx(RspFlow.TechniccalTrendAnaysisData[i].StockIndex);

        json = ToJson(&(RspFlow.TechniccalTrendAnaysisData[i]), json);
        
        cout << "JSON: " << json << endl;

        mc_->Set(code + MEMC_KEY_LIST[MC_KEY_TREND_ANYS], json);
    }

    return RspFlow.StockNumber;
}

}
