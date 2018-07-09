
#ifndef _JW_PROTOCOL_STOCK_H_
#define _JW_PROTOCOL_STOCK_H_

#include "jw_os.h"
#include "jw_struct.h"
#include "jw_protocol.h"
//#include "jw_define.h"


/******************************************************************
*������:JWPackCostDistributionDataReq 
*����˵�������������ɳɱ��ֲ�ָ��������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN StockIndex:��������
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackCostDistributionDataReq(char *buf, int len, const unsigned int StockIndex);

/******************************************************************
*������:JWUnPackCostDistributionDataRsp
*����˵�������������ɳɱ��ֲ�ָ���Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:���ɳɱ��ֲ�ָ���Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackCostDistributionDataRsp(const char * p, int len, JWCostDistributionRsp * pRsp);

/******************************************************************
*������:JWPackCostDistributionDetailReq 
*����˵�������������ɳɱ���ϸָ��������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN StockIndex:��������
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackCostDistributionDetailReq(char *buf, int len, const unsigned int StockIndex);

/******************************************************************
*������:JWUnPackCostDistributionDetailRsp
*����˵�������������ɳɱ���ϸָ���Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:���ɳɱ���ϸ��Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackCostDistributionDetailRsp(const char * p, int len, JWCostDistributionDetailRsp * pRsp); 

/******************************************************************
*������:JWPackStkFundFlowDataReq 
*����˵���������������ʽ��������ݵ������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:������Ϣ�ṹ��
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackStkFundFlowDataReq(char *buf, int len, const JWStkFundFlowInfo *pInfo);

/******************************************************************
*������:JWUnPackCostDistributionDetailRsp
*����˵�������������ɳɱ��ֲ�ָ���Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:�����ʽ�������ϸ��Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackStkFundFlowDataRsp(const char * p, int len, JWStkFundFlowDataRsp * pRsp); 

/******************************************************************
*������:JWPackIndustryFundFlowDataReq 
*����˵������������ҵ�ʽ��������ݵ������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:��ҵ�ʽ�����������Ϣ�ṹ��
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackIndustryFundFlowDataReq(char *buf, int len, const JWIndustryFundFlowInfo *pInfo);

/******************************************************************
*������:JWUnPackCostDistributionDetailRsp
*����˵������������ҵ�ʽ������Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ҵ�ʽ�������ϸ��Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackIndustryFundFlowDataRsp(const char * p, int len, JWIndustryFundFlowDataRsp * pRsp); 

/******************************************************************
*������:JWPackGetDDEDataReq 
*����˵����������DDE���ݵ������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:DDE����������Ϣ�ṹ��
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackGetDDEDataReq(char *buf, int len, const JWDdeDataInfo *pInfo);

/******************************************************************
*������:JWUnPackGetDDEDataRsp
*����˵����������DDE���ݵ�Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:DDE���ݵ�Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetDDEDataRsp(const char * p, int len, JWDDEDataRsp * pRsp); 


/******************************************************************
*������:JWPackGetDDSortDataReq 
*����˵����������DDSort���ݵ������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:DDE����������Ϣ�ṹ��
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackGetDDSortDataReq(char *buf, int len, const JWDdsortDataInfo *pInfo);

/******************************************************************
*������:JWUnPackGetDDSortDataRsp
*����˵����������DDSort���ݵ�Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:DDSort���ݵ�Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetDDSortDataRsp(const char * p, int len, JWDDSortDataRsp * pRsp); 

/******************************************************************
*������:JWPackGetDDEHisDataReq 
*����˵������������ʷDDE���ݵ������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:��ʷDDE����������Ϣ�ṹ��
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackGetDDEHisDataReq(char *buf, int len, const JWDdeHisDataInfo *pInfo);

/******************************************************************
*������:JWUnPackGetDDEHisDataRsp
*����˵������������ʷDDE���ݵ�Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ʷDDE���ݵ�Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetDDEHisDataRsp(const char * p, int len, JWDDEHisDataRsp * pRsp); 

/******************************************************************
*������:JWPackRiseFallBreadthDataReq 
*����˵����������ʵʱ�����ǵ������ݵ������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:��ʵʱ�����ǵ�������������Ϣ�ṹ��
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackRiseFallBreadthDataReq(char *buf, int len, const JWRiseFallBreadthDataInfo *pInfo);

/******************************************************************
*������:JWUnPackGetDDEHisDataRsp
*����˵����������ʵʱ�����ǵ������ݵ�Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:ʵʱ�����ǵ������ݵ�Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackRiseFallBreadthDataRsp(const char * p, int len, JWRiseFallBreadthDataRsp * pRsp); 

/******************************************************************
*������:JWPackRiseFallBreadthDataReq 
*����˵����������ȡ�ۺ������Ϣ�������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN stockindex:Ҫ��ȡ�ĸ������
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackGetDiagnoseDataReq(char *buf, int len, const unsigned short stockindex);


/******************************************************************
*������:JWUnPackGetDDEHisDataRsp
*����˵����������ȡ�ۺ������Ϣ��Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ�ۺ������Ϣ��Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetDiagnoseDataRsp(const char * p, int len, JWDiagnoseDataRsp * pRsp);


/******************************************************************
*������:JWPackGetStkFundFlowDetailDataReq 
*����˵����������ȡ�����ʽ�������ϸ���ݵ������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN reqInfo:��������ṹ��
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackGetStkFundFlowDetailDataReq(char *buf, int len, const JWStkFundFlowDataDetailReq *reqInfo);

/******************************************************************
*������:JWUnPackGetStkFundFlowDetailDataRsp
*����˵���������������ʽ�������ϸ���ݵ�Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ�����ʽ�������ϸ���ݽṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetStkFundFlowDetailDataRsp(const char * p, int len, JWStkFundFlowDetailDataRsp *pRsp);


/******************************************************************
*������:JWPackGetJWKlineFormDataReq 
*����˵����������ȡ����K����̬�������� �������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN reqInfo:��������ṹ��
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackGetJWKlineFormDataReq(char *buf, int len, const JWKlineFormDataReq *reqInfo);

/******************************************************************
*������:JWUnPackGetJWKlineFormDataRsp
*����˵��������������K����̬�������� ��Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ����K����̬�������� �ݽṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetJWKlineFormDataRsp(const char * p, int len, JWKlineFormDataRsp *pRsp);


/******************************************************************
*������:JWPackGetMultiEmptyGameDataReq 
*����˵����������ȡ��ղ�����Ӫ���ݵ������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN stockindex:��������
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackGetMultiEmptyGameDataReq(char *buf, int len, const unsigned short stockindex);

/******************************************************************
*������:JWUnPackGetMultiEmptyGameDataRsp
*����˵������������ղ�����Ӫ���ݵ�Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ��ղ�����Ӫ���ݾݽṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetMultiEmptyGameDataRsp(const char * p, int len, JWMultiEmptyGameDataRsp *pRsp);


/******************************************************************
*������:JWPackGetTechniccalTrendAnaysisDataReq 
*����˵����������������-���Ʒ������ݵ������  0x0A69
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN reqInfo:��������ṹ��
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackGetTechniccalTrendAnaysisDataReq(char *buf, int len, const JWTechniccalTrendAnaysisReq *reqInfo);

/******************************************************************
*������:JWUnPackTechniccalTrendAnaysisDataRsp
*����˵��������������-���Ʒ������ݵ�Ӧ���  0x0A69
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ������-���Ʒ������ݾݾݽṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackTechniccalTrendAnaysisDataRsp(const char * p, int len, JWTechniccalTrendAnaysisRsp *pRsp);


/******************************************************************
*������:JWPackGetFinancialMainCostDataReq 
*����˵�����������ʽ���-��������-�����ɱ����ݵ������  0x0A70
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN reqInfo:��������ṹ��
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackGetFinancialMainCostDataReq(char *buf, int len, const JWFinancialMainCostReq *reqInfo);

/******************************************************************
*������:JWUnPackGetFinancialMainCostDataRsp
*����˵���������ʽ���-��������-�����ɱ����ݵ�Ӧ���   0x0A70
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ�ʽ���-��������-�����ɱ����ݽṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetFinancialMainCostDataRsp(const char * p, int len, JWFinancialMainCostRsp *pRsp);


/******************************************************************
*������:JWPackGetTechniccalAndMarketDataReq 
*����˵����������������-�г��������ݵ������  0x0A71
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN reqInfo:��������ṹ��
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackGetTechniccalAndMarketDataReq(char *buf, int len, const JWTechniccalAndMarketReq *reqInfo);

/******************************************************************
*������:JWUnPackGetTechniccalAndMarketDataRsp
*����˵��������������-�г��������ݵ�Ӧ���  0x0A71
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ������-�г��������ݽṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetTechniccalAndMarketDataRsp(const char * p, int len, JWTechniccalAndMarketRsp *pRsp);


/******************************************************************
*������:JWPackGetFundAmentalsAnalysisDataReq 
*����˵�������ɻ�����-����������ݵ������  0x0A72
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN reqInfo:��������ṹ��
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackGetFundAmentalsAnalysisDataReq(char *buf, int len, const JWFundAmentalsAnalysisReq *reqInfo);

/******************************************************************
*������:JWUnPackGetFundAmentalsAnalysisDataRsp
*����˵��������������-����������ݵ�Ӧ���  0x0A72
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ������-����������ݽṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetFundAmentalsAnalysisDataRsp(const char * p, int len, JWFundAmentalsAnalysisRsp *pRsp);


/******************************************************************
*������:JWPackGetFundAmentalsTargetsDataReq 
*����˵�������ɻ�����-�������-����ָ�����ݵ������  0x0A73
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN reqInfo:��������ṹ��
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackGetFundAmentalsTargetsDataReq(char *buf, int len, const JWFundAmentalsTargetsReq *reqInfo);

/******************************************************************
*������:JWUnPackGetFundAmentalsTargetsDataRsp
*����˵��������������-�������-����ָ�����ݵ�Ӧ���  0x0A73
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ������-�������-����ָ�����ݽṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetFundAmentalsTargetsDataRsp(const char * p, int len, JWFundAmentalsTargetsRsp *pRsp);



/******************************************************************
*������:JWPackGetFundIndustryRanksDataReq 
*����˵�������ɻ�����-�������-��ҵ�������ݵ������  0x0A74
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN reqInfo:��������ṹ��
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackGetFundIndustryRanksDataReq(char *buf, int len, const JWFundIndustryRanksReq *reqInfo);

/******************************************************************
*������:JWUnPackGetFundIndustryRanksDataRsp
*����˵��������������-�������-��ҵ�������ݵ�Ӧ���  0x0A74
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ������-�������-��ҵ�������ݽṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetFundIndustryRanksDataRsp(const char * p, int len, JWFundIndustryRanksRsp *pRsp);



/******************************************************************
*������:JWPackGetFundAluationReserchDataReq   0x0A75
*����˵�������ɻ�����-��ֵ�о����ݵ������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN reqInfo:��������ṹ��
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackGetFundAluationReserchDataReq(char *buf, int len, const JWFundAluationReserchReq *reqInfo);

/******************************************************************
*������:JWUnPackGetFundAluationReserchDataRsp  0x0A75
*����˵��������������-��ֵ�о����ݵ�Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ������-��ֵ�о����ݽṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetFundAluationReserchDataRsp(const char * p, int len, JWFundAluationReserchRsp *pRsp);



/******************************************************************
*������:JWPackGetFundOrganizationViewDataReq   0x0A76
*����˵�������ɻ����۵����ݵ������  
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN reqInfo:��������ṹ��
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackGetFundOrganizationViewDataReq(char *buf, int len, const JWFundOrganizationViewReq *reqInfo);

/******************************************************************
*������:JWUnPackGetFundOrganizationViewDataRsp  0x0A76
*����˵�������������۵����ݵ�Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ�����۵����ݽṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetFundOrganizationViewDataRsp(const char * p, int len, JWFundOrganizationViewRsp *pRsp);


#endif
