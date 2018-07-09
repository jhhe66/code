#include "FundDtlHandler.h"
#include "global.h"

#include <string>
#include <json/json.h>

using namespace std;
using namespace Json;

namespace Handler
{

inline int
CFundDtlHandler::SyncToAll()
{
	client_->SetRecvTimeout(5);
	client_->Connect();
	
	
	JWStkFundFlowDataDetailReq ReqFlow;

	ReqFlow.StockNumber = map_->size();

	for(unsigned short i = 0; i < map_->size(); i++)
		ReqFlow.StockIndex[i] = map_->GetIndexByOffset(i); 
	
	char ReqBuff[29 + sizeof ReqFlow + 1] = {0};

	int ReqLen = JWPackGetStkFundFlowDetailDataReq(ReqBuff, sizeof ReqBuff, &ReqFlow);
	
	if ( ReqLen < 0)
		return -1;

	HexDumpImp(ReqBuff, ReqLen);

	int sLen = 0;

	while(sLen < ReqLen)
	{
		int sVal = client_->Send(ReqBuff + sLen, ReqLen - sLen);

		if (sVal <= 0)
			return -2;

		sLen += sVal;
	}

	char RspBuff[1024 * 1024] = {0};

	int rLen = 0;
	int rVal = 0;
	
	JWStkFundFlowDetailDataRsp RspFlow;
	
	while(JWHandleInput(RspBuff, rLen) == 0)
	{
		rVal = client_->Recv(RspBuff + rLen, sizeof(RspBuff) - rLen);

		if(rVal < 0)
			return -3;
		
		rLen += rVal;
	}

	if (JWUnPackGetStkFundFlowDetailDataRsp(RspBuff, rLen, &RspFlow))
		return -4;

	HexDumpImp(RspBuff, rLen);

	cout << "rLen: " << rLen << endl;

	string json;
	string code;
	
	for(int i = 0; i < RspFlow.StockNumber;i++)
	{
		json = ToJson(&(RspFlow.StkFundFlowData[i]), json);
        code = map_->GetCodeByIdx(RspFlow.StkFundFlowData[i].StockIndex);

		cout << "JSON: " << json << endl;
        cout << "KEY: " << code + MEMC_KEY_LIST[MC_KEY_FUNDS_DTL] << endl;
		mc_->Set(code + MEMC_KEY_LIST[MC_KEY_FUNDS_DTL], json);
	}

    return RspFlow.StockNumber;
}

inline string&
CFundDtlHandler::ToJson(const void* data, string& json)
{
    const JWStkFundFlowDetailData* p = (const JWStkFundFlowDetailData*)data;

	Value v;

	const unsigned short UNIT = 1000;

	v[Code] = map_->GetCodeByIdx(p->StockIndex);
	v[Date] = p->Date;
	v[BuyBigOrderAmount] = ToSigned(p->BuyBigOrderAmount, UNIT);
	v[BuyBigOrderVolume] = ToSigned(p->BuyBigOrderVolume, UNIT);
	v[BuyLargeOrderAmount] = ToSigned(p->BuyLargeOrderAmount, UNIT);
	v[BuyLargeOrderVolume] = ToSigned(p->BuyLargeOrderVolume, UNIT);
	v[BuyMidOrderAmount] = ToSigned(p->BuyMidOrderAmount, UNIT);
	v[BuyMidOrderVolume] = ToSigned(p->BuyMidOrderVolume, UNIT);
	v[BuyMinOrderAmount] = ToSigned(p->BuyMinOrderAmount, UNIT);
	v[BuyMinOrderVolume] = ToSigned(p->BuyMinOrderVolume, UNIT);

	v[DealCountAmount] = ToSigned(p->DealCountAmount, UNIT);
	v[DealCountVolume] = ToSigned(p->DealCountVolume, UNIT);

	v[SelBigOrderAmount] = ToSigned(p->SelBigOrderAmount, UNIT);
	v[SelBigOrderVolume] = ToSigned(p->SelBigOrderVolume, UNIT);
	v[SelLargeOrderAmount] = ToSigned(p->SelLargeOrderAmount, UNIT);
	v[SelLargeOrderVolume] = ToSigned(p->SelLargeOrderVolume, UNIT);
	v[SelMidOrderAmount] = ToSigned(p->SelMidOrderAmount, UNIT);
	v[SelMidOrderVolume] = ToSigned(p->SelMidOrderVolume, UNIT);
	v[SelMinOrderAmount] = ToSigned(p->SelMinOrderAmount, UNIT);
	v[SelMinOrderVolume] = ToSigned(p->SelMinOrderVolume, UNIT);

	FastWriter writer;

	json = writer.write(v);

	return json;
}

}
