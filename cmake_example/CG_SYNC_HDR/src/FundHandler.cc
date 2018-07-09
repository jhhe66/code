#include "FundHandler.h"
#include "global.h"

#include <json/json.h>
#include <string>

using namespace Json;
using namespace std;

namespace Handler
{

inline string&
CFundHandler::ToJson(const void* data, string& json)
{
    Value array;

    for(unsigned short day = 0; day <= TwentyDays; day++)
    {
        const JWStkFundFlowDataRsp* p = (const JWStkFundFlowDataRsp*)data + day;
        const JWStkFundFlowData* pdata = p->StkFundFlowData; // &(p->StkFundFlowData[0]);
        
        Value v;

		v[TimeStamp] = pdata->TimeStamp;
        v[DayCount] = pdata->DayCount;									 
		v[IndustryCode] = pdata->IndustryCode;							 
		v[NowPrice] = ToSigned(pdata->NowPrice, 1000);					 
		v[DeltaPercent] = ToSigned(pdata->DeltaPercent, 1000);			 
		v[ChangeRate] = ToSigned(pdata->ChangeRate, 1000);				 
		v[TotalFlowIn] = ToSigned(pdata->TotalFlowIn, 1000); 			 
		v[TotalFlowOut] = ToSigned(pdata->TotalFlowOut, 1000);			 
		v[NetFlowIn] = ToSigned(pdata->NetFlowIn, 1000); 				 
		v[NetFlowOut] = ToSigned(pdata->NetFlowOut, 1000);				 
		v[NetFlowInPower] = ToSigned(pdata->NetFlowInPower, 100000); 	 
		v[NetFlowOutPower] = ToSigned(pdata->NetFlowOutPower, 100000);	 
		v[NetBigBill] = ToSigned(pdata->NetBigBill, 1000);				 
		v[ImpetusBill] = ToSigned(pdata->ImpetusBill, 100000);			 
		v[MainNetFlowIn] = ToSigned(pdata->MainNetFlowIn, 1000); 		 
		v[MainNetFlowOut] = ToSigned(pdata->MainNetFlowOut, 1000);		 
		v[MainNetInRate] = ToSigned(pdata->MainNetInRate, 100000);		 
		v[MainNetOutRate] = ToSigned(pdata->MainNetOutRate, 100000); 	 
		v[SeriesAddDays] = ToSigned(pdata->SeriesAddDays, 0);					 
		v[SeriesNetIn] = ToSigned(pdata->SeriesNetIn, 1000); 			 
		v[SeriesNetOut] = ToSigned(pdata->SeriesNetOut, 1000);			 
		v[SeriesNetInPower] = ToSigned(pdata->SeriesNetInPower, 100000);  
		v[SeriesNetOutPower] = ToSigned(pdata->SeriesNetOutPower, 100000);
		v[AreaAmountRate] = ToSigned(pdata->AreaAmountRate, 100000); 	 
		v[AreaChangeRate] = ToSigned(pdata->AreaChangeRate, 1000);		 
		v[AreaDeltaPercent] = ToSigned(pdata->AreaDeltaPercent, 1000);	 
		v[AmountInRate] = ToSigned(pdata->AmountOutRate, 100000);		 
		v[AreaClose] = ToSigned(pdata->AreaClose, 1000); 				 
		v[nDayClose] = ToSigned(pdata->nDayClose, 1000); 				 
		v[AreaDate] = pdata->AreaDate;									 
		v[nDayDate] = pdata->nDayDate;  
        
        array.append(v);
    }
    
    Value root;
    
    root[Code] = map_->GetCodeByIdx(((const JWStkFundFlowDataRsp*)data)->StkFundFlowData->StockIndex);
    root[Num] = ((const JWStkFundFlowDataRsp*)data)->StockNumber;
    root[Data] = array;
    
    FastWriter writer;

    json = writer.write(root);
    
    return json;
}

inline int
CFundHandler::SyncToAll()
{
    for(unsigned short index = 0; index < map_->size(); index++)
	{
		JWStkFundFlowDataRsp RspFlow[TwentyDays];
		
		for(unsigned short day = 0; day <= TwentyDays; day++)
		{
			client_->SetRecvTimeout(1);
			client_->Connect();
			
			JWStkFundFlowInfo ReqFlow;
			
			ReqFlow.DayCount = day;
			ReqFlow.StockNumber = 1;
			ReqFlow.StockIndex[0] = map_->GetIndexByOffset(index);
			
			
			char ReqBuff[PACK_HEADER_LEN + sizeof ReqFlow + 1] = {0};
			
			int ReqLen = JWPackStkFundFlowDataReq(ReqBuff, sizeof(ReqBuff), &ReqFlow);
			
			//HexDumpImp(ReqBuff, ReqLen);
			
			
			if (ReqLen <= 0)
			    return -1;
			
			int sLen = 0;
			
			sLen = client_->Send(ReqBuff, ReqLen);
			
			if (sLen <= 0)
			    return -2;
			
			char RspBuff[4096] = {0};
			
			int rLen = 0;
			
			while(JWHandleInput(RspBuff, rLen) == 0)
			{
			    int rVal = client_->Recv(RspBuff + rLen, sizeof(RspBuff) - rLen);
			
			    if (rVal == -1)
				    return -3;
			
			    rLen += rVal;
			}
			
			//JWStkFundFlowDataRsp RspFlow;
			
			JWUnPackStkFundFlowDataRsp(RspBuff, rLen, RspFlow + day);
			
			//HexDumpImp(RspBuff, RspLen);
		}
		
		string json("");
		string code = map_->GetCodeByOffset(index);

         
		
		json = ToJson(RspFlow, json);
		
		mc_->Set(code + MEMC_KEY_LIST[MC_KEY_FUNDS_PATH], json);
    }

    return map_->size();
}


}


