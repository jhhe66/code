
#ifndef _JW_PROTOCOL_STOCK_H_
#define _JW_PROTOCOL_STOCK_H_

#include "jw_os.h"
#include "jw_struct.h"
#include "jw_protocol.h"
//#include "jw_define.h"


/******************************************************************
*函数名:JWPackCostDistributionDataReq 
*功能说明：生成拉个股成本分布指标的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN StockIndex:个股索引
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackCostDistributionDataReq(char *buf, int len, const unsigned int StockIndex);

/******************************************************************
*函数名:JWUnPackCostDistributionDataRsp
*功能说明：解析拉个股成本分布指标的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:个股成本分布指标的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackCostDistributionDataRsp(const char * p, int len, JWCostDistributionRsp * pRsp);

/******************************************************************
*函数名:JWPackCostDistributionDetailReq 
*功能说明：生成拉个股成本明细指标的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN StockIndex:个股索引
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackCostDistributionDetailReq(char *buf, int len, const unsigned int StockIndex);

/******************************************************************
*函数名:JWUnPackCostDistributionDetailRsp
*功能说明：解析拉个股成本明细指标的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:个股成本明细的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackCostDistributionDetailRsp(const char * p, int len, JWCostDistributionDetailRsp * pRsp); 

/******************************************************************
*函数名:JWPackStkFundFlowDataReq 
*功能说明：生成拉个股资金流向数据的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:请求消息结构体
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackStkFundFlowDataReq(char *buf, int len, const JWStkFundFlowInfo *pInfo);

/******************************************************************
*函数名:JWUnPackCostDistributionDetailRsp
*功能说明：解析拉个股成本分布指标的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:个股资金流向明细的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackStkFundFlowDataRsp(const char * p, int len, JWStkFundFlowDataRsp * pRsp); 

/******************************************************************
*函数名:JWPackIndustryFundFlowDataReq 
*功能说明：生成拉行业资金流向数据的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:行业资金流向请求消息结构体
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackIndustryFundFlowDataReq(char *buf, int len, const JWIndustryFundFlowInfo *pInfo);

/******************************************************************
*函数名:JWUnPackCostDistributionDetailRsp
*功能说明：解析拉行业资金流向的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:行业资金流向明细的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackIndustryFundFlowDataRsp(const char * p, int len, JWIndustryFundFlowDataRsp * pRsp); 

/******************************************************************
*函数名:JWPackGetDDEDataReq 
*功能说明：生成拉DDE数据的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:DDE数据请求消息结构体
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackGetDDEDataReq(char *buf, int len, const JWDdeDataInfo *pInfo);

/******************************************************************
*函数名:JWUnPackGetDDEDataRsp
*功能说明：解析拉DDE数据的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:DDE数据的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetDDEDataRsp(const char * p, int len, JWDDEDataRsp * pRsp); 


/******************************************************************
*函数名:JWPackGetDDSortDataReq 
*功能说明：生成拉DDSort数据的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:DDE数据请求消息结构体
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackGetDDSortDataReq(char *buf, int len, const JWDdsortDataInfo *pInfo);

/******************************************************************
*函数名:JWUnPackGetDDSortDataRsp
*功能说明：解析拉DDSort数据的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:DDSort数据的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetDDSortDataRsp(const char * p, int len, JWDDSortDataRsp * pRsp); 

/******************************************************************
*函数名:JWPackGetDDEHisDataReq 
*功能说明：生成拉历史DDE数据的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:历史DDE数据请求消息结构体
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackGetDDEHisDataReq(char *buf, int len, const JWDdeHisDataInfo *pInfo);

/******************************************************************
*函数名:JWUnPackGetDDEHisDataRsp
*功能说明：解析拉历史DDE数据的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:历史DDE数据的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetDDEHisDataRsp(const char * p, int len, JWDDEHisDataRsp * pRsp); 

/******************************************************************
*函数名:JWPackRiseFallBreadthDataReq 
*功能说明：生成拉实时行情涨跌幅数据的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:拉实时行情涨跌幅数据请求消息结构体
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackRiseFallBreadthDataReq(char *buf, int len, const JWRiseFallBreadthDataInfo *pInfo);

/******************************************************************
*函数名:JWUnPackGetDDEHisDataRsp
*功能说明：解析拉实时行情涨跌幅数据的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:实时行情涨跌幅数据的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackRiseFallBreadthDataRsp(const char * p, int len, JWRiseFallBreadthDataRsp * pRsp); 

/******************************************************************
*函数名:JWPackRiseFallBreadthDataReq 
*功能说明：生成拉取综合诊断信息的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN stockindex:要拉取的个股序号
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackGetDiagnoseDataReq(char *buf, int len, const unsigned short stockindex);


/******************************************************************
*函数名:JWUnPackGetDDEHisDataRsp
*功能说明：解析拉取综合诊断信息的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:拉取综合诊断信息的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetDiagnoseDataRsp(const char * p, int len, JWDiagnoseDataRsp * pRsp);


/******************************************************************
*函数名:JWPackGetStkFundFlowDetailDataReq 
*功能说明：生成拉取个股资金流向明细数据的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN reqInfo:请求参数结构体
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackGetStkFundFlowDetailDataReq(char *buf, int len, const JWStkFundFlowDataDetailReq *reqInfo);

/******************************************************************
*函数名:JWUnPackGetStkFundFlowDetailDataRsp
*功能说明：解析拉个股资金流向明细数据的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:拉取个股资金流向明细数据结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetStkFundFlowDetailDataRsp(const char * p, int len, JWStkFundFlowDetailDataRsp *pRsp);


/******************************************************************
*函数名:JWPackGetJWKlineFormDataReq 
*功能说明：生成拉取个股K线形态引擎数据 的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN reqInfo:请求参数结构体
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackGetJWKlineFormDataReq(char *buf, int len, const JWKlineFormDataReq *reqInfo);

/******************************************************************
*函数名:JWUnPackGetJWKlineFormDataRsp
*功能说明：解析拉个股K线形态引擎数据 的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:拉取个股K线形态引擎数据 据结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetJWKlineFormDataRsp(const char * p, int len, JWKlineFormDataRsp *pRsp);


/******************************************************************
*函数名:JWPackGetMultiEmptyGameDataReq 
*功能说明：生成拉取多空博弈阵营数据的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN stockindex:个股索引
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackGetMultiEmptyGameDataReq(char *buf, int len, const unsigned short stockindex);

/******************************************************************
*函数名:JWUnPackGetMultiEmptyGameDataRsp
*功能说明：解析拉多空博弈阵营数据的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:拉取多空博弈阵营数据据结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetMultiEmptyGameDataRsp(const char * p, int len, JWMultiEmptyGameDataRsp *pRsp);


/******************************************************************
*函数名:JWPackGetTechniccalTrendAnaysisDataReq 
*功能说明：生成拉技术面-趋势分析数据的请求包  0x0A69
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN reqInfo:请求参数结构体
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackGetTechniccalTrendAnaysisDataReq(char *buf, int len, const JWTechniccalTrendAnaysisReq *reqInfo);

/******************************************************************
*函数名:JWUnPackTechniccalTrendAnaysisDataRsp
*功能说明：解析技术面-趋势分析数据的应答包  0x0A69
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:拉取技术面-趋势分析数据据据结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackTechniccalTrendAnaysisDataRsp(const char * p, int len, JWTechniccalTrendAnaysisRsp *pRsp);


/******************************************************************
*函数名:JWPackGetFinancialMainCostDataReq 
*功能说明：生成拉资金面-主力分析-主力成本数据的请求包  0x0A70
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN reqInfo:请求参数结构体
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackGetFinancialMainCostDataReq(char *buf, int len, const JWFinancialMainCostReq *reqInfo);

/******************************************************************
*函数名:JWUnPackGetFinancialMainCostDataRsp
*功能说明：解析资金面-主力分析-主力成本数据的应答包   0x0A70
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:拉取资金面-主力分析-主力成本数据结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetFinancialMainCostDataRsp(const char * p, int len, JWFinancialMainCostRsp *pRsp);


/******************************************************************
*函数名:JWPackGetTechniccalAndMarketDataReq 
*功能说明：生成拉技术面-市场表现数据的请求包  0x0A71
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN reqInfo:请求参数结构体
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackGetTechniccalAndMarketDataReq(char *buf, int len, const JWTechniccalAndMarketReq *reqInfo);

/******************************************************************
*函数名:JWUnPackGetTechniccalAndMarketDataRsp
*功能说明：解析技术面-市场表现数据的应答包  0x0A71
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:拉取技术面-市场表现数据结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetTechniccalAndMarketDataRsp(const char * p, int len, JWTechniccalAndMarketRsp *pRsp);


/******************************************************************
*函数名:JWPackGetFundAmentalsAnalysisDataReq 
*功能说明：生成基本面-财务分析数据的请求包  0x0A72
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN reqInfo:请求参数结构体
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackGetFundAmentalsAnalysisDataReq(char *buf, int len, const JWFundAmentalsAnalysisReq *reqInfo);

/******************************************************************
*函数名:JWUnPackGetFundAmentalsAnalysisDataRsp
*功能说明：解析基本面-财务分析数据的应答包  0x0A72
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:拉取基本面-财务分析数据结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetFundAmentalsAnalysisDataRsp(const char * p, int len, JWFundAmentalsAnalysisRsp *pRsp);


/******************************************************************
*函数名:JWPackGetFundAmentalsTargetsDataReq 
*功能说明：生成基本面-财务分析-利润指标数据的请求包  0x0A73
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN reqInfo:请求参数结构体
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackGetFundAmentalsTargetsDataReq(char *buf, int len, const JWFundAmentalsTargetsReq *reqInfo);

/******************************************************************
*函数名:JWUnPackGetFundAmentalsTargetsDataRsp
*功能说明：解析基本面-财务分析-利润指标数据的应答包  0x0A73
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:拉取基本面-财务分析-利润指标数据结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetFundAmentalsTargetsDataRsp(const char * p, int len, JWFundAmentalsTargetsRsp *pRsp);



/******************************************************************
*函数名:JWPackGetFundIndustryRanksDataReq 
*功能说明：生成基本面-财务分析-行业排名数据的请求包  0x0A74
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN reqInfo:请求参数结构体
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackGetFundIndustryRanksDataReq(char *buf, int len, const JWFundIndustryRanksReq *reqInfo);

/******************************************************************
*函数名:JWUnPackGetFundIndustryRanksDataRsp
*功能说明：解析基本面-财务分析-行业排名数据的应答包  0x0A74
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:拉取基本面-财务分析-行业排名数据结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetFundIndustryRanksDataRsp(const char * p, int len, JWFundIndustryRanksRsp *pRsp);



/******************************************************************
*函数名:JWPackGetFundAluationReserchDataReq   0x0A75
*功能说明：生成基本面-估值研究数据的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN reqInfo:请求参数结构体
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackGetFundAluationReserchDataReq(char *buf, int len, const JWFundAluationReserchReq *reqInfo);

/******************************************************************
*函数名:JWUnPackGetFundAluationReserchDataRsp  0x0A75
*功能说明：解析基本面-估值研究数据的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:拉取基本面-估值研究数据结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetFundAluationReserchDataRsp(const char * p, int len, JWFundAluationReserchRsp *pRsp);



/******************************************************************
*函数名:JWPackGetFundOrganizationViewDataReq   0x0A76
*功能说明：生成机构观点数据的请求包  
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN reqInfo:请求参数结构体
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackGetFundOrganizationViewDataReq(char *buf, int len, const JWFundOrganizationViewReq *reqInfo);

/******************************************************************
*函数名:JWUnPackGetFundOrganizationViewDataRsp  0x0A76
*功能说明：解析机构观点数据的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:拉取机构观点数据结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetFundOrganizationViewDataRsp(const char * p, int len, JWFundOrganizationViewRsp *pRsp);


#endif
