#include "FinanceAnysHandler.h"

#include "global.h"

#include <json/json.h>
#include <string>

using namespace std;
using namespace Json;

namespace Handler
{


string&
CFinanceAnysHandler::ToJson(const void* data, string& json)
{
    Value array;
    
    const JWFinancialMainCostContentData* p = ((const JWFinancialMainCostData*)data)->FinancialMainCostContentData;

    for ( unsigned short i = 0; i < ((const JWFinancialMainCostData*)data)->Num; i++)
    {
        Value v;
        
        v[UpdateDate] = (p + i)->UpdateDate;
        v[ClosePrice] = ToSigned((p + i)->ClosePrice, 1000);
        v[MarketPrice] = ToSigned((p + i)->MarketPrice, 1000);
        v[MainPrice] = ToSigned((p + i)->MainPrice, 1000);
        
        array.append(v);
    }

    Value root;

    root[Code] = map_->GetCodeByIdx(((const JWFinancialMainCostData*)data)->StockIndex);
    root[Num] = ((const JWFinancialMainCostData*)data)->Num;
    root[Data] = array;

    FastWriter writer;

    json = writer.write(root);

    return json;
}

int
CFinanceAnysHandler::SyncToAll()
{
    JWFinancialMainCostReq ReqFlow;
    
    ReqFlow.StockNumber = map_->size();

    for (unsigned short index = 0; index < ReqFlow.StockNumber;index++)
    {
        ReqFlow.StockIndex[index] = map_->GetIndexByOffset(index);
    }
    
    client_->SetRecvTimeout(1);
    //client_->SetRecvBuffSize(1024 * 100);
    client_->Connect();

    char ReqBuff[PACK_HEADER_LEN + 1 + sizeof ReqFlow] = {0};
    
    int ReqLen = JWPackGetFinancialMainCostDataReq(ReqBuff, 
            sizeof ReqBuff, &ReqFlow);

    if (ReqLen < 0)
        return -1;

    int sLen = 0;
    
    while(sLen < ReqLen)
    {
        int sVal = client_->Send(ReqBuff + sLen, ReqLen - sLen);

        if (sVal <= 0)
        {
            return -2;
        }

        sLen += sVal;
    }

    
    HexDumpImp(ReqBuff, sLen);

    char RspBuff[PACK_HEADER_LEN + 1 + sizeof(JWFinancialMainCostRsp)] = {0};

    int rLen = 0;

    while(JWHandleInput(RspBuff, rLen) == 0)
    {
        int rVal = client_->Recv(RspBuff + rLen, 
                    sizeof RspBuff - rLen);
        
        if ( rVal <= 0)
            return -3;

        rLen += rVal;
    }
    
    HexDumpImp(RspBuff, rLen);

    JWFinancialMainCostRsp RspFlow;

    if(JWUnPackGetFinancialMainCostDataRsp(RspBuff, rLen, &RspFlow) < 0 )
        return -4;
    
    string json("");
    string code("");
    
    for ( unsigned short i = 0; i < RspFlow.StockNumber; i++)
    {
        code = map_->GetCodeByIdx(RspFlow.FinancialMainCostData[i].StockIndex); 
            
        json = ToJson(&RspFlow.FinancialMainCostData[i], json);
            
        cout << "JSON: " << json << endl;

        mc_->Set( code + MEMC_KEY_LIST[MC_KEY_FIN_MAIN_COST], json );
    }
    
    return map_->size();
}

}
