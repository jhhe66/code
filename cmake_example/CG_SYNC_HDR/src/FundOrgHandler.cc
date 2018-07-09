#include "FundOrgHandler.h"
#include "global.h"

#include <json/json.h>
#include <string>

using namespace std;
using namespace Json;

namespace Handler
{

inline string&
CFundOrgHandler::ToJson(const void* data, string& json)
{
    const JWFundOrganizationViewData* p = (const JWFundOrganizationViewData*)data;

    static const unsigned int UNIT = 100000;

    Value v;

    v[Code] = map_->GetCodeByIdx(p->StockIndex);
    v[UpdateDate] = p->UpdateDate;
    v[OrganNum] = p->OrganNum;
    v[StatDays] = p->StatDays;
    v[Rating] = ToSigned(p->Rating, UNIT);
    v[Attention] = ToSigned(p->Attention, UNIT);
    v[BuyNum] = p->BuyNum;
    v[HoldingsNum] = p->HoldingsNum;
    v[NeutralNum] = p->NeutralNum;
    v[ReductionNum] = p->ReductionNum;
    v[SellNum] = p->SellNum;
    v[PE_TTM_RANKING] = p->PE_TTM_RANKING;
    v[PB] = ToSigned(p->PB, UNIT);
    v[PB_HY] = ToSigned(p->PB_HY, UNIT);
    v[PB_RANKING] = p->PB_RANKING;
    v[PEG] = ToSigned(p->PEG, UNIT);
    v[PEG_HY] = ToSigned(p->PEG_HY, UNIT);
    v[PEG_RANKING] = p->PEG_RANKING;

    FastWriter writer;

    json = writer.write(v);

    return json;
}

inline int
CFundOrgHandler::SyncToAll()
{
   client_->SetRecvTimeout(1);
   client_->Connect();

    JWFundOrganizationViewReq ReqFlow;

    ReqFlow.StockNumber = map_->size();

    for(unsigned short index = 0; index < ReqFlow.StockNumber; index++)
    {
        ReqFlow.StockIndex[index] = map_->GetIndexByOffset(index);
    }

    char ReqBuff[PACK_HEADER_LEN + 1 + sizeof ReqFlow] = {0};

    int ReqLen = 
        JWPackGetFundOrganizationViewDataReq(ReqBuff, sizeof ReqBuff, &ReqFlow);

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

    //HexDumpImp(ReqBuff, sLen);
    
    char RspBuff[PACK_HEADER_LEN + 1 + sizeof(JWFundOrganizationViewRsp)] = {0};

    int rLen = 0;

    while(JWHandleInput(RspBuff, rLen) == 0)
    {
        int rVal = client_->Recv(RspBuff + rLen, sizeof RspBuff - rLen);

        if(!rVal)
            return -3;

        rLen += rVal;
    }

    //HexDumpImp(RspBuff, rLen);

    JWFundOrganizationViewRsp RspFlow;

    if(JWUnPackGetFundOrganizationViewDataRsp(RspBuff, rLen, &RspFlow) < 0)
        return -4;
    
    string json;
    string code;

    for(unsigned short i = 0; i <  RspFlow.StockNumber; i++)
    {
        code = map_->GetCodeByIdx(RspFlow.FundOrganizationViewData[i].StockIndex);

        json = ToJson(&RspFlow.FundOrganizationViewData[i], json);

        cout << "JSON: " << json << endl;

        mc_->Set(code + MEMC_KEY_LIST[MC_KEY_FUND_ORG], json);
    }

    return RspFlow.StockNumber;
}

}


