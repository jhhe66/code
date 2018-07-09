/********************************************
//文件名:jw_define.h
//功能:系统结构定义头文件
//作者:钟何明
//创建时间:2010.07.07
//修改记录:

*********************************************/
#ifndef _JW_STRUCT_H_
#define _JW_STRUCT_H_

#ifdef _WIN32
typedef unsigned __int64 uint64_t;
#else
#include <stdint.h>
#endif



#define JW_MAX_ERR_MSG_LEN              255 //错误消息长度
#define JW_MAX_ACCOUNT_LEN              64  //用户帐号长度
#define JW_MAX_PWD_LEN                  64  //用户密码长度
#define JW_MAX_TICKET_LEN               32  //认证票据长度
#define JW_MAX_DCODE_LEN                16  //动态码长度
#define JW_MAX_STOCK_CODE_LEN           10  //股票代码长度
#define JW_MAX_NAME_LEN                 64  //用户名长度
#define JW_MAX_BIRTH_LEN                12  //生日长度
#define JW_MAX_MOBILE_LEN               20  //手机号码长度
#define JW_MAX_VERIFY_CODE_LEN          20  //验证码长度
#define JW_MAX_URL_LEN                  255 //URL长度
#define JW_MAX_PROVIANCE_LEN            64 //区域信息长度
#define JW_MAX_CITY_LEN                 64  //城市
#define JW_MAX_AREA_CODE_LEN            6   //区号 

#define JW_MAX_GROUP_STOCK_NUM          50  //每个分组中最多的自选股数量
#define JW_MAX_RANK_NUM                 50  //最多股票排名数量
#define JW_MAX_BIND_MOBILE_NUM          6   //最多绑定手机号码数量

#define JW_MAX_CONTENT_LEN  4096 //上报内容长度

#define JW_MAX_ORDERID_LEN 32 //订单号长度
#define JW_MAX_INCOME_NUM 32 //收益个数
#define JW_MAX_PRODUCT_ID_LEN 32 //商品编号长度

#define JW_MAX_MAIL_LEN 64 //邮件的最大长度

#define JW_MAX_STOCK_BUY_SELL_NUM 10 //最大委买，委卖个股数量


#define JW_MAX_STOCK_NUM 5000 //最大个股数量
#define JW_MAX_SECUFID_LEN 31 //新闻内容ID串最大长度
#define JW_MAX_STOCK_NAME_LEN 32 //股票名称最大长度
#define JW_NEWSMEDIA_LEN 256 //新闻出处最大长度 256
#define JW_NEW_CONTENT_LEN (64*1024) //新闻内容最大长度64K
#define JW_MAX_SECU_NEWS_TITLE_LEN 256 //新闻标题最大长度256
#define JW_MAX_STOCKCODE_LEN 8  //股票代码最大长度
#define JW_STOCK_SHORT_NAME_LEN 32  //股票简称最大长度
#define JW_MAX_MINDATA_NUM 250  //个股最大分时数据个数

#define JW_MAX_ALIAS_NUM 10 //最多10个别名
#define JW_MAX_REGION_NUM 10 //最大区间个数
#define JW_MAX_OPER_ANALY 512 //操作分析最大长度
#define JW_MAX_INDU_CODE_LEN 30 //行业代码最大长度
#define JW_MAX_INDUSTRY_NUM 300 //最大行业个数
#define JW_MAX_DDE_HIS_NUM 100 //最大历史DDE数据
#define JW_MAX_RISE_FALL_NUM 150 //涨跌幅数据最大条数
#define JW_MAX_DIAGNOSE_CONTENT_LEN 1024 //诊断内容最大长度
#define JW_MAX_PARAMNAME_LEN 100 //指标名称最打长度
#define JW_MAX_KLINE_TYPE 30 //K线形态类型
#define JW_MAX_KLINE_NUM 1000 //K线形态最大个数
#define JW_MAX_MULTI_EMPTY_GAME_CAMP 100

#define JW_MAX_RANK_LEN 30 //排行最大长度
#define JW_MAX_INDUSTRY_NAME_LEN 10 //最大行业名称长度
#define JW_MAX_GLOSSARY_LEN 50 //分析术语最大长度
#define JW_MAX_MAIN_COST_CONTENT_NUM 15 //主力成本数据内容最大个数 
#define JW_MAX_TARGETS_CONTENT_NUM 30//利润指标数据内容个数

#pragma pack (push, 1)

/* 公共请求消息结构 */
struct JWCommonReqInfo
{
    unsigned int UIN;//用户ID
};

/* 公共应答消息结构 */
struct JWCommonRsp
{
    unsigned int nErrno;//错误码
    char szErrMsg[JW_MAX_ERR_MSG_LEN + 1];//错误信息
};

/*查询帐号状态*/
struct JWGetAccountStatusInfo
{
    char szAccount[JW_MAX_ACCOUNT_LEN + 1];
};

struct JWGetAccountStatusRsp:public JWCommonRsp
{
    unsigned char nStatus;//帐号状态，0:已注册已激活，1：已注册未激活，2：未注册
};

/* 开通帐号信息结构 */
struct JWCreateAccountInfo
{
    char szAccount[JW_MAX_ACCOUNT_LEN + 1];//帐号
    char szPasswd[JW_MAX_PWD_LEN + 1];//密码  
    unsigned int nSourceIP;//来源IP（接入服务器取到的远端IP）
    char szName[JW_MAX_NAME_LEN + 1];//用户昵称

};


struct JWCreateAccountRsp:public JWCommonRsp
{    
    unsigned int UIN;//用户ID
};

/* 登录信息结构 */
struct JWLoginInfo
{
    char szAccount[JW_MAX_ACCOUNT_LEN + 1];//帐号
    char szTicket[JW_MAX_TICKET_LEN + 1];//认证票据
    char szDCode[JW_MAX_DCODE_LEN + 1];//动态码
    unsigned int nSourceIP;//来源IP（接入服务器取到的远端IP）
    unsigned int nClientIP;//客户端本机IP，无填0
};

struct JWLoginRsp:public JWCommonRsp
{    
    unsigned int UIN;//用户ID
    char szTicket[JW_MAX_TICKET_LEN + 1];//认证票据
    unsigned int nStatus;//帐号状态
};

/* 注销登录信息结构 */
struct JWLogoutInfo:public JWCommonReqInfo
{    
    unsigned int nSourceIP;//来源IP
};

struct JWLogoutRsp:public JWCommonRsp
{    
};

/* 股票代码 + 市场类型 */
struct JWStockKey
{
    char szStockCode[JW_MAX_STOCK_CODE_LEN + 1];//股票代码
    unsigned int nMarketType;//市场类型
};

struct JWGetStockInfo:public JWCommonReqInfo
{
    unsigned int nGroupID;//自选股分组ID
};

struct JWGetStockInfoRsp:public JWCommonRsp
{    
    unsigned int nGroupID;//自选股分组ID
    unsigned short nStockNum;//拉到的数量
    JWStockKey StockList[JW_MAX_GROUP_STOCK_NUM];//股票列表
};


struct JWAddStockInfo:public JWCommonReqInfo
{
    JWStockKey StockKey;//股票代码+市场类型
    unsigned int nGroupID;//分组ID
};

struct JWAddStockRsp:public JWCommonRsp
{
    
};

struct JWDelStockInfo:public JWCommonReqInfo
{
    JWStockKey StockKey;//股票代码+市场类型
    unsigned int nGroupID;//股票分组ID
};

struct JWDelStockRsp:public JWCommonRsp
{
    
};


struct JWSetStockSortInfo:public JWCommonReqInfo
{
    unsigned int nGroupID;//自选股分组ID
    unsigned short nStockNum;//自选股数量
    JWStockKey StockList[JW_MAX_GROUP_STOCK_NUM];//股票列表

};

struct JWSetStockSortRsp:public JWCommonRsp
{
    
};

struct JWGetStockSortInfo:public JWCommonReqInfo
{
    unsigned int nGroupID;//自选股分组ID
};

struct JWGetStockSortRsp:public JWCommonRsp
{    
    unsigned int nGroupID;//自选股分组ID
    unsigned short nStockNum;//自选股数量
    JWStockKey StockList[JW_MAX_GROUP_STOCK_NUM];//股票列表
};

struct JWGetStockRankInfo:public JWCommonReqInfo
{
    unsigned short Type;//排序类型，1：最多自选排名，2：自选飚升排名
};

struct JWStockRankItem
{
    JWStockKey StockKey;//股票KEY
    unsigned int nAttentionCount;//关注次数
    unsigned char nStatus;//排名状态: -1为排名下降；0为排名不变；1为排名上升；2为新增
};
struct JWGetStockRankRsp:public JWCommonRsp
{    
    unsigned short Type;//排序类型，1：最多自选排名，2：自选飚升排名
    unsigned short nStockNum;//自选股数量
    JWStockRankItem RankList[JW_MAX_RANK_NUM];
};

#if 0
//2010-11-26 zhongheming 业务变更，废弃掉旧接口 
struct JWBindMobileInfo:public JWCommonReqInfo
{
    char szMobile[JW_MAX_MOBILE_LEN + 1];//手机号码
    unsigned int nFlag;//绑定标志
    unsigned char opType;//操作类型 0.绑定, 1.解除绑定
};

struct JWBindMobileRsp:public JWCommonRsp
{    
    unsigned int nExpire;//验证码失效时间，单位秒
};

struct JWBindVerifyInfo:public JWCommonReqInfo
{
    char szMobile[JW_MAX_MOBILE_LEN + 1];//手机号码
    char szVerifyCode[JW_MAX_VERIFY_CODE_LEN + 1];//验证码
    unsigned char opType;//操作类型 0.绑定, 1.解除绑定
};

struct JWBindVerifyRsp:public JWCommonRsp
{    

};

struct JWGetBindInfo:public JWCommonReqInfo
{

};

struct JWBindItem
{
    
    char szMobile[JW_MAX_MOBILE_LEN + 1];//手机号码
    unsigned int nFlag;//绑定标志
};

struct JWGetBindInfoRsp:public JWCommonRsp
{    
    unsigned char nNumber;//数量
    JWBindItem BindList[JW_MAX_BIND_MOBILE_NUM];

};
#endif

////////////////////////////////////////////////////////////////////////////////////////////////////////////
//新的绑定手机帐号接口所用的结构
////////////////////////////////////////////////////////////////////////////////////////////////////////////

struct JWBindMobileInfoReq:public JWCommonReqInfo
{
    char szMobile[JW_MAX_MOBILE_LEN + 1];//手机号码
    unsigned char BindType;//绑定的业务类型,见jw_define.h中 jw_bind_type 定义
    unsigned char opType;//操作类型 0.绑定, 1.解除绑定    
};

struct JWBindMobileInfoRsp:public JWCommonRsp
{    
    unsigned int nExpire;//验证码失效时间，单位秒
};



struct JWBindVerifyInfoReq:public JWCommonReqInfo
{
    char szMobile[JW_MAX_MOBILE_LEN + 1];//手机号码
    unsigned char BindType;//绑定的业务类型,见jw_define.h中 jw_bind_type 定义
    char szVerifyCode[JW_MAX_VERIFY_CODE_LEN + 1];//验证码
    unsigned char opType;//操作类型 0.绑定, 1.解除绑定
};

struct JWBindVerifyInfoRsp:public JWCommonRsp
{    

};

struct JWGetBindMobileInfoReq:public JWCommonReqInfo
{
    unsigned char BindType;//绑定的业务类型,见jw_define.h中 jw_bind_type 定义
};

struct JWBindItemInfo
{

    char szMobile[JW_MAX_MOBILE_LEN + 1];//手机号码
    unsigned char BindType;//绑定的业务类型,见jw_define.h中 jw_bind_type 定义
};

struct JWGetBindMobileInfoRsp:public JWCommonRsp
{    
    unsigned char nNumber;//数量
    JWBindItemInfo BindList[JW_MAX_BIND_MOBILE_NUM];
};

////////////////////////////////////////////////////////////////////////////////////////////////////////////

struct JWGetBasicInfo:public JWCommonReqInfo
{    
    unsigned char nFaceType;//头像类型,1:大头像，2：小头像
};

struct JWGetBasicInfoRsp:public JWCommonRsp
{
    char szName[JW_MAX_NAME_LEN + 1];//姓名
    unsigned char nGender;//性别
    unsigned int nBirthday;//生日,如19800101格式的数字串
    char szFaceURL[JW_MAX_URL_LEN + 1];//头像URL
};

struct JWSetBasicInfoReq:public JWCommonReqInfo
{
    char szName[JW_MAX_NAME_LEN + 1];//姓名
    unsigned char nGender;//性别
    unsigned int nBirthday;//生日,如19800101格式的数字串
};

struct JWSetBasicInfoRsp:public JWCommonRsp
{

};

/*取得手机归宿地信息*/
struct JWGetMobileAreaInfo:public JWCommonReqInfo
{

};

struct JWGetMobileAreaInfoRsp:public JWCommonRsp
{
    char szProviance[JW_MAX_PROVIANCE_LEN + 1];//手机归属省份
    char szCity[JW_MAX_CITY_LEN + 1];//手机归属城市
    char szAreaCode[JW_MAX_AREA_CODE_LEN + 1];//手机归属区号
};


struct JWResetPwdInfo
{
    char szAccount[JW_MAX_ACCOUNT_LEN + 1];//帐号
};

struct JWResetPwdRsp:public JWCommonRsp
{
    char szNewPwd[JW_MAX_PWD_LEN + 1];//用户新密码
};

struct JWGetMobileAreaByNumber
{
    char szMobileNumber[JW_MAX_MOBILE_LEN + 1];
};


/* 验证票据 */
struct JWCheckTicketInfo:public JWCommonReqInfo
{
    char szTicket[JW_MAX_TICKET_LEN + 1];//认证票据
    unsigned int nSourceIP;//来源IP（外网IP）
    unsigned int nClientIP;//客户端IP
};

struct JWCheckTicketRsp:public JWCommonRsp
{

};


/*取得帐号的注册时间*/
struct JWGetAccountRegTimeInfo
{
    char szAccount[JW_MAX_ACCOUNT_LEN + 1];//用户帐号
};

struct JWGetAccountRegTimeRsp:public JWCommonRsp
{
    unsigned int nRegTime;//注册时间(UNIX TimeStamp)
    unsigned int UIN;//用户ID
};

/*通过UIN取得帐号信息*/
struct JWGetUserAccountInfo:public JWCommonReqInfo
{

};

struct JWGetUserAccountRsp:public JWCommonRsp
{
    char szAccount[JW_MAX_ACCOUNT_LEN + 1];//用户帐号信息
};

/*用户数据上报*/
struct JWStatReportInfo:public JWCommonReqInfo
{
    char szAccount[JW_MAX_ACCOUNT_LEN + 1];//用户帐号信息
    unsigned char nSrcType;//上报数据的产品来源
    unsigned int nSrcIP;//用户来源IP地址
    unsigned int nDataType;//上报数据类型,参见<用户数据上报格式说明.docx>
    char szContent[JW_MAX_CONTENT_LEN + 1];//上报内容,长度限制为4K
};

struct JWStatReportRsp:public JWCommonRsp
{
};


struct JWModifyPwdInfo
{
    char szAccount[JW_MAX_ACCOUNT_LEN + 1];//帐号
    char szTicket[JW_MAX_TICKET_LEN + 1];//认证票据
    char szDCode[JW_MAX_DCODE_LEN + 1];//动态码
    char szNewPwd[JW_MAX_PWD_LEN + 1];//用户密码(md5 32)
    unsigned int nSourceIP;//来源IP（接入服务器取到的远端IP）
    unsigned int nClientIP;//客户端本机IP，无填0
    
};

struct JWModifyUserPwdRsp:public JWCommonRsp
{
};


struct JWIncomeItem
{
    unsigned int dwAmount;//收益点数
    unsigned int dwSaleUIN;//收益用户UIN
    char szProductID[JW_MAX_PRODUCT_ID_LEN + 1];//商品编号
    unsigned int dwPayAmount;//商品支付金额
};

struct JWPayTicketIncomeInfo:public JWCommonReqInfo
{
    uint64_t qwSeq;//流水号
    char szOrderID[JW_MAX_ORDERID_LEN +1];//订单号
    unsigned int dwOrderAmount;//订单金额
    unsigned int dwOrderResult;//订单状态
    unsigned int dwSubmitTime;//订单提交时间
    unsigned int dwIncomeNum;//收益个数
    JWIncomeItem ItemList[JW_MAX_INCOME_NUM];
};

struct JWPayTicketIncomeRsp:public JWCommonRsp
{
    
};


struct JWSetUserAliasInfo:public JWCommonReqInfo
{
    unsigned char OpType;//操作类型，1:添加,2:删除
    char szAlias[JW_MAX_ACCOUNT_LEN + 1];//用户登录别名
};

struct JWSetUserAliasInfoRsp:public JWCommonRsp
{
};

struct JWGetUserAliasInfo:public JWCommonReqInfo
{
};

struct JWGetUserAliasInfoRsp:public JWCommonRsp
{
    unsigned short alias_num;
    char szAlias[JW_MAX_ALIAS_NUM][JW_MAX_ACCOUNT_LEN + 1];//用户登录别名
};

struct JWUpdTicketInfo:public JWCommonReqInfo
{
    unsigned char Type;//票据类型
};

struct JWUpdTicketInfoRsp:public JWCommonRsp
{
    char szTicket[JW_MAX_TICKET_LEN + 1];//票据
    unsigned int nGetInterval;//下次拉取时间间隔
};

//用户注册信息
struct JWRegAccountInfo 
{
    unsigned char IsActivate;//是否注册同时激活,0:不激活,1:立即激活
    char verify_code[JW_MAX_DCODE_LEN+1];//激活码,utf-8编码,\0结束
    char account[JW_MAX_ACCOUNT_LEN+1];//账号,utf-8编码,\0结束
    char password[JW_MAX_PWD_LEN+1];// 小写的md5密码,32字节
    char name[JW_MAX_NAME_LEN+1];//用户姓名，,utf-8编码,\0结束
    unsigned char gender;//性别0:未知,1:男,2:女
    char birth[JW_MAX_BIRTH_LEN+1];//生日,格式为YYYY/MM/DD,\0结束

    unsigned short investLevel;//投资能力 0:有点余钱，想看看   1：经济能力一般  2：中产阶级 3：比较富裕  4:什么都没有
    unsigned int attention;//操盘风格.
    unsigned int operatorStyle;//操盘风格 ,每一bit代表一项  0：短线 1：中线 2：长线 
    
    unsigned int city;//居住地
    unsigned int ip_addr;//用户注册时的IP
    unsigned int internalID;//内部ID,如注册来源URL等
    unsigned int advid;//广告来源ID
    unsigned int productType;//嵌入页面类型，
};

struct JWRegAccountRsp:public JWCommonRsp
{
    char szEmail[JW_MAX_MAIL_LEN+1];//邮件登录地址
};

struct JWAccountInfo
{
    char account[JW_MAX_ACCOUNT_LEN+1];//账号,utf-8编码,\0结束

};

//获取的用户注册信息
struct JWAccountInfoRsp:public JWCommonRsp
{
    char account[JW_MAX_ACCOUNT_LEN+1];//账号,utf-8编码,\0结束
    char verify_code[JW_MAX_VERIFY_CODE_LEN+1];//激活码,utf-8编码,\0结束
    unsigned int reg_time;//注册时间，YYYYMMDD格式。
    char name[JW_MAX_NAME_LEN+1];//用户姓名，,utf-8编码,\0结束
    unsigned char gender;//性别0:未知,1:男,2:女
    char birth[JW_MAX_BIRTH_LEN+1];//生日格式：YYYYMMDD  例如：19810306
    unsigned int invest_level;//投资能力
    unsigned int operatorStyle;//操盘风格 ,每一bit代表一项  0：短线 1：中线 2：长线 
    unsigned int city;//居住地
    unsigned int ip_addr;//注册IP
};



//账户激活
struct JWActiveAccountInfo
{
    char account[JW_MAX_ACCOUNT_LEN+1];//账号,utf-8编码,\0结束
    char verify_code[JW_MAX_DCODE_LEN+1];//激活码,utf-8编码,\0结束
};

struct JWActiveAccountInfoRsp:public JWCommonRsp
{
};
//
struct JWAuthInfo:public JWCommonReqInfo
{
    unsigned char type; //票据类型
};

struct JWAuthInfoRsp:public JWCommonRsp
{
    char ticket[JW_MAX_TICKET_LEN+1];
};

struct JWXXXXRsp:public JWCommonRsp
{    
    unsigned short nStockNum;//拉到的数量
    char a[1];//股票列表
};

///////////////////////////////////0x0a09/////////////////////////////////
struct JWBuyLevel
{
	unsigned int nBuyPx;           //委买价
	unsigned int nBuyVolume;       //委买量
};

struct JWSellLevel
{
	unsigned int nSellPx;          //委卖价
	unsigned int nSellVolume;      //委卖量
};


struct JWStockDetail //个股详情
{
	unsigned short nStockIndex;                   //股票索引
	unsigned int nTxTime;                         //交易时间
	unsigned int nPreClosePx;                     //昨收价
	unsigned int nOpenPx;                         //开盘价
	unsigned int nHighPx;                         //最高价
	unsigned int nLowPx;                          //最低价
	unsigned int nLastPx;                         //现价
	uint64_t lTotalVolume;                        //成交总量
	uint64_t lTotalAmount;                        //成交总金额
	uint64_t lFloatShare;                         //流通股
	uint64_t lTotalShare;                         //总股本
	int nPanCha;                                  //盘差
	unsigned int nVolumeRatio;                    //量比
	int nWeiBi;                                   //委比
	int nWeiCha;                                  //委差
	unsigned char cBuyLevel;                      //委买个数
	JWBuyLevel BuyInfo[JW_MAX_STOCK_BUY_SELL_NUM];//委买详情
	unsigned char cSellLevel;                     //委卖个数
	JWSellLevel SellInfo[JW_MAX_STOCK_BUY_SELL_NUM];//委卖详情
			
};

struct JWPullMarketDataInfo
{
    unsigned short nStockNumber;                   //个股数量
    unsigned short nStockIndex[JW_MAX_STOCK_NUM];  //个股索引

};

struct JWPullMarketDataRsp:public JWCommonRsp
{
    unsigned short nStockNumber;                   //个股行情数
    JWStockDetail StockDetail[JW_MAX_STOCK_NUM];   //个股详情数组
};

///////////////////////////////////////0x0A0A//////////////////////

struct JWPullMinuteDataInfo:public JWPullMarketDataInfo
{

};

//分时数据结构
struct JWMinData
{
    unsigned int nTxTime;                   //交易时间
    unsigned int nHighPx;                   //最高价，服务端 x1000，客户端 /1000
    unsigned int nLowPx;                    //最低价，服务端 x1000，客户端 /1000
    unsigned int nLastPx;                   //现价，服务端 x1000，客户端 /1000

    uint64_t nTotalVolume;                  //成交总量, 服务端/100000，单位是手，客户端 x100
    uint64_t nTotalAmount;                  //成交总金额，单位是元, 服务端 x1000，客户端 /1000
};

struct JWStockMinData
{
    unsigned short nStockIndex;              //股票索引
    char szStockCode[JW_MAX_STOCKCODE_LEN]; //股票代码（utf-8编码）
    unsigned int nStockType;                  //市场类型
    unsigned int nPreClosePx;                //昨收价，服务端 x1000，客户端 /1000
    unsigned int nOpenPx;                    //开盘价，服务端 x1000，客户端 /1000
    unsigned short nMinDataNum;              //分时数据个数
    JWMinData MinData[JW_MAX_MINDATA_NUM];  //分时数据
};

struct JWPullMinuteDataInfoRsp:public JWCommonRsp
{
   unsigned short nStockNumber;
   JWStockMinData StockMinData[JW_MAX_STOCK_NUM];

};

///////////////////////////////////////0x0a10///////////////////////
struct JWStockRadarIndicatorInfo
{
    unsigned short nStockIndex;                    //个股索引
};

struct JWStockRadarIndicatorRsp:public JWCommonRsp
{
    unsigned short nStockIndex;              //个股索引
    unsigned char cFinanceRisk;              //财务风险(服务端 +50 如果>100，则直接=100；如果<0，则直接=0)范围：0 ~ 100) (ZHPC|11
    unsigned char cPriceRisk;                //价格风险，(范围：0 ~ 100) (ZHPC|13)
    unsigned char cAgencyRating;             //机构评级，(范围：0 ~ 100) (ZHPC|9)
    unsigned char cCompanyOperation;         //公司运营，(范围：0 ~ 100) (ZHPC|15)
    unsigned char cMoneyFlow;                //资金流向，(范围：0 ~ 100) (ZHPC|7)

};


/////////////////////////////0x0a1d///////////////////////////////////
struct JWEarlyWarningInfo
{
    unsigned char cMarketType;      //市场类型,
    unsigned int nVerCounter;       //版本计数器，股票服务器传递的值，需要回传给股票服务器，用于鉴别是否有新的数据
};


struct JWMarketDataInfo  //个股行情数据结构
{
    unsigned short nStockIndex;         //股票索引
    unsigned int nTxTime;               //交易时间
    unsigned int nPreClosePx;           //昨收价
    unsigned int nLastPx;               //现价
    uint64_t lTotalVolume;              //成交总量
    uint64_t lFloatShare;               //流通股
};

struct JWEarlyWarningRsp:public JWCommonRsp
{
    unsigned int nVerCounter;                         //版本计数器。
    unsigned short nStockNumber;                      //股票个数
    JWMarketDataInfo MarketDataInfo[JW_MAX_STOCK_NUM];//个股行情数据结构
};


/////////////////////////////0x0a22///////////////////////////
struct JWStockBriefMacketData  //股票简易行情数据结构
{
    unsigned short nStockIndex;    //股票索引
    unsigned int nTxTime;          //交易时间
    unsigned int nPreClosePx;      //昨收价，服务端 x1000，客户端 /1000
    unsigned int nOpenPx;          //开盘价，服务端 x1000，客户端 /1000
    unsigned int nLastPx;          //现价，服务端 x1000，客户端 /1000
    uint64_t nTotalVolume;         //成交总量，单位是手，客户端 x100
    uint64_t nTotalAmount;         //成交总金额，单位是元
    uint64_t nFloatShare;          //流通股
    int nPanCha;                   //盘差，可能<0
    unsigned int nVolumeRatio;      //量比，服务端 /10 (因为大赢家传的数据 已经x1000)，客户端 /100

};

struct JWBatchBriefMarketDataInfo :public JWPullMarketDataInfo //批量简要的行情数据请求消息结构与拉最新行情数据请求消息结构体一样
{

};

struct JWBatchBriefMarketDataInfoRsp:public JWCommonRsp  //批量拉取股票简易行情数据应答消息体结构
{
    unsigned short nStockNumber;                                //股票个数
    JWStockBriefMacketData BriefMacketData[JW_MAX_STOCK_NUM];   //股票简易行情数据结构

};

//////////////////////////////////0x0a28///////////////////////////////////////////////////
struct JWStkNewsContentInfo
{
    unsigned short nStockIndex;                       //股票索引
    char szSecuFid[JW_MAX_SECUFID_LEN];               //新闻内容ID串（该值为拉取自选股时候从服务器中取到的SecuFid）
};

struct JWStkNewsContentInfoRsp:public JWCommonRsp
{
    char szStockName[JW_MAX_STOCK_NAME_LEN];               //股票名称
    unsigned int nSecuNewsTime;                            //个股新闻时间(Mine/GP, unix时间)
    char szSecuNewsTitle[JW_MAX_SECU_NEWS_TITLE_LEN];      //个股新闻标题(Mine/GP)
    char szNewsMedia[JW_NEWSMEDIA_LEN];                    //新闻出处
    char szNewsContent[JW_NEW_CONTENT_LEN];                //新闻内容

};

//////////////////////////////////0x0a29/批量拉取单只股票最新新闻标题使用的结构定义////////////////////////////////
struct JWStockNewsNumInfo
{
    unsigned short nStockIndex;    //股票索引
    unsigned short nNewsNum;       // 拉取的新闻内容条数：1~20

    JWStockNewsNumInfo& operator=(const JWStockNewsNumInfo &rs)
    {
        this->nStockIndex = rs.nStockIndex;
        this->nNewsNum = rs.nNewsNum;
        return *this;
    }
};

struct JWBatchStockNewsNumInfo  //批量拉取股票最新新闻标题请求结构体
{
    unsigned short nStockNumber;
    JWStockNewsNumInfo StockNews[JW_MAX_STOCK_NUM];

};

struct JWStockSingleNewsData  //个股单条新闻结构
{
    unsigned int nSecuNewsTime;
    char szSecuNewsTitle[JW_MAX_SECU_NEWS_TITLE_LEN];
    char szSecuFid[JW_MAX_SECUFID_LEN];
};

struct JWStockAllNewsData
{
    unsigned short nStockIndex;
    unsigned nNewsNum;
    JWStockSingleNewsData NewsData[20]; //单只股票最大拉取20条新闻

};

struct JWBatchStockNewsNumInfoRsp:public JWCommonRsp
{
    unsigned short nStockNumber;
    JWStockAllNewsData StockAllNewsData[JW_MAX_STOCK_NUM];

};

//////////////////////////取reference data, 0x0a60//////////////////////////
//该命令的请求消息体为空

struct JWStockInfoData
{
    unsigned short nStockIndex;                           //股票索引
    unsigned int nStockType;                              //股票类型
    char szStockCode[JW_MAX_STOCKCODE_LEN];               //股票代码
    char szStockName[JW_MAX_STOCK_NAME_LEN];              //股票名称
    char szStockShortName[JW_STOCK_SHORT_NAME_LEN];       //股票简称
    unsigned int nPreClosePx;                             //昨收价
    uint64_t nFloatShare;                                 //流通股
    uint64_t nTotalShare;                                 //总股本
};

struct JWStockGetRefereceDataRsp : public JWCommonRsp
{
    unsigned int nRefVersion; //codelist版本号
    unsigned int nStockNumber;
    JWStockInfoData StockInfoData[JW_MAX_STOCK_NUM];

};

//////////////////////////个股成本分布指标, 0x0a17//////////////////////////
struct  JWRegionData
{
       unsigned char RegionType; //区间类型：(CBFB)
       unsigned int RegionBegin;  //区间起点（客户端  /1000）,单位：元(CBFB)
       unsigned int RegionEnd;  //区间终点（客户端  /1000）,单位：元(CBFB)
       unsigned int HangupNumber;//套牢手数（客户端  /1000），,单位：手(CBFB)
       unsigned int HangupAmount; //套牢金额（客户端  /1000），单位：万元(CBFB)
};

struct JWCostDistributionRsp : public JWCommonRsp
{
    unsigned short StockIndex; //股票索引
    unsigned char RegionNumber;//区间个数
    JWRegionData Region[JW_MAX_REGION_NUM]; //分部指标数组
    char OperatorAnalysis[JW_MAX_OPER_ANALY];//操作分析，（UTF8）

};



//////////////////////////个股成本分布明细, 0x0a18//////////////////////////
struct  JWRegionData_0A18
{
    unsigned char RegionType; //区间类型：
    unsigned int RegionBegin; //区间起点
    unsigned int RegionEnd;   //区间终点
};

struct JWRegionDetailData
{
    unsigned int DealPx;  //成交价
    unsigned int ChipsNumber; //筹码数
};

struct JWCostDistributionDetailRsp : public JWCommonRsp
{
    unsigned short StockIndex; //股票索引
    unsigned char RegionNumber;//区间个数
    JWRegionData_0A18 Region[JW_MAX_REGION_NUM]; //区间数据数组
    unsigned short DetailNumber; //明细个数
    JWRegionDetailData DetailData[1024]; //明细数组

};

//////////////////////////拉个股资金流向数据, 0x0a50//////////////////////////
struct JWStkFundFlowInfo
{
    unsigned short DayCount; //日期类型(0=今日、1=3日、2=5日、3=10日、4=20日)
    unsigned short StockNumber; //股票个数
    unsigned short StockIndex[JW_MAX_STOCK_NUM]; //股票索引数组
};

struct JWStkFundFlowData
{
    unsigned short StockIndex; //股票索引
    unsigned int TimeStamp; //时间戳
    unsigned short DayCount; //日期类型(0=今日、1=3日、2=5日、3=10日、4=20日)
    char IndustryCode[JW_MAX_INDU_CODE_LEN];//行业代码
    unsigned int NowPrice; //现价，服务端 x1000，客户端 /1000 （大赢家传过来的的数据就已经x1000， 股票服务器透传，客户端需要/1000）
    unsigned int DeltaPercent; //涨跌幅
    unsigned int ChangeRate;//换手率
    uint64_t TotalFlowIn; //总流入
    uint64_t TotalFlowOut;//总流出
    uint64_t NetFlowIn;//净流入
    uint64_t NetFlowOut;//净流出
    unsigned int NetFlowInPower;//净流入力度 
    unsigned int NetFlowOutPower;//净流出力度 
    uint64_t NetBigBill;//大单净量
    unsigned int ImpetusBill;//大单动能 
    uint64_t MainNetFlowIn; //主力净流入
    uint64_t MainNetFlowOut;//主力净流出
    unsigned int MainNetInRate;//主力净流入力度
    unsigned int MainNetOutRate;//主力净流出力度
    unsigned int SeriesAddDays;//连续增仓天数
    uint64_t SeriesNetIn; //连续增仓机构净流入服务器/1000
    uint64_t SeriesNetOut;//连续增仓机构净流出服务器/1000
    uint64_t SeriesNetInPower;//连续减仓机构净流入力度
    uint64_t SeriesNetOutPower;//连续减仓机构净流出力度
    unsigned int AreaAmountRate;//占区间总成交比例
    unsigned int AreaChangeRate;//区间换手率 
    unsigned int AreaDeltaPercent;//区间涨跌幅
    unsigned int AmountInRate;//净流入占成交额比
    unsigned int AmountOutRate;//净流出占成交额比
    unsigned int AreaClose; //区间收盘价
    unsigned int nDayClose;//3日或者5日或者其它日期的收盘价
    unsigned int AreaDate;
    unsigned int nDayDate;
};

struct JWStkFundFlowDataRsp:public JWCommonRsp
{
    unsigned short StockNumber; //股票个数
    JWStkFundFlowData StkFundFlowData[JW_MAX_STOCK_NUM];//资金流向数据结构数组
};

//////////////////////////拉行业资金流向数据, 0x0a53//////////////////////////
struct JWIndustryFundFlowInfo
{
    unsigned short DayCount; //日期类型(0=今日、1=3日、2=5日、3=10日、4=20日)
    unsigned short IndustryType;//行业类型（0=证监会行业; 1=新财富行业）
    unsigned short IndustryNumber;//行业个数
    unsigned short IndustryIndex[JW_MAX_INDUSTRY_NUM];//行业索引数组
};

struct JWIndustryFundFlowData
{
    unsigned short IndustryIndex; //行业索引
    unsigned int TimeStamp; //时间戳
    unsigned int DeltaPercent;//涨跌幅
    unsigned int DownCount; //下跌家数
    unsigned int UpCount; //上涨家数
    unsigned int ChangeRate; //行业换手率
    uint64_t TotalFlowIn;//行业资金总流入
    uint64_t TotalFlowOut;//行业资金总流出
    uint64_t NetFlowIn;//行业资金净流入
    unsigned int NetFlowInPower;//行业资金净流入力度
    unsigned int LeaderStock;//行业资领涨股
    unsigned int LeaderZDF;//领涨股涨跌幅
};

struct JWIndustryFundFlowDataRsp:public JWCommonRsp
{
    unsigned short DayCount;//日期类型(0=今日、1=3日、2=5日、3=10日、4=20日)
    unsigned short IndustryType;//行业类型（0=证监会行业; 1=新财富行业）
    unsigned short IndustryNumber;//行业个数
    JWIndustryFundFlowData IndustryFundFlowData[JW_MAX_INDUSTRY_NUM];//行业资金流向数据结构数组
};

//////////////////////////拉DDE数据, 0x0a61//////////////////////////
struct JWDdeDataInfo
{
    unsigned short StockNumber;//股票个数
    unsigned short StockIndex[JW_MAX_STOCK_NUM];//股票索引数组
};

struct JWDDEData
{
    unsigned short StockIndex; //股票索引
    unsigned int Date; //日期（格式：Unix格式时间数值）
    uint64_t MainHold;//主力持仓
    uint64_t DisperseHold; //散户持仓
    uint64_t BuyBigOrder; //买入大单 数/金额
    uint64_t SelBigOrder; //卖出大单 数/金额
    uint64_t BuyMinOrder; //买入小单 数/金额
    uint64_t SelMinOrder; //卖出小单 数/金额
    uint64_t BuyMidOrder;//买入中单 数/金额
    uint64_t SelMidOrder;//卖出中单 数/金额
    uint64_t BuyLargeOrder;//买入特大单 数/金额
    uint64_t SelLargeOrder; //卖出特大单 数/金额
    uint64_t BuyCount; //买入成交笔数
    uint64_t SellCount; //卖出成交笔数
};

struct JWDDEDataRsp:public JWCommonRsp
{
    unsigned short DDEDataNums; //DDE数据个数
    JWDDEData DDEData[JW_MAX_STOCK_NUM];//DDE数据结构数组
};

//////////////////////////拉DDSort数据, 0x0a62//////////////////////////
struct JWDdsortDataInfo:public JWDdeDataInfo
{
   
};

struct JWDdsortData
{
    unsigned short StockIndex; //股票索引
    unsigned int Time;//日期时间
    unsigned int Price;//现价  * 1000 取整
    unsigned int DeltaPercent;//涨跌幅 * 10000 取整, 还原后又正负号
    unsigned int DDX;//DDX * 10000 取整, 还原后又正负号
    unsigned int DDY;//DDY * 10000 取整, 还原后又正负号
    unsigned int DDZ;//DDZ * 10000 取整, 还原后又正负号
    unsigned int DDX60;//60日DDX * 10000 取整, 还原后又正负号
    unsigned int DDY60;//60日DDY * 10000 取整, 还原后又正负号
    unsigned int Up10;//10日内飘红
    unsigned int UpDays;//连续飘红
    unsigned int LagerBuy;//特大买* 10000 取整, 还原后又正负号
    unsigned int BigBuy;//大单买* 10000 取整, 还原后又正负号
    unsigned int LagerSell;//特大卖* 10000 取整, 还原后又正负号
    unsigned int BigSell;//大单买* 10000 取整, 还原后又正负号
};

struct JWDDSortDataRsp:public JWCommonRsp
{
    unsigned short StockNumber;//股票个数
    JWDdsortData DdsortData[JW_MAX_STOCK_NUM];//个股DDSort数据结构数组
};

////////////拉DDE历史数据（?Business 服务器） 0x0A64/////////////////
struct JWDdeHisDataInfo
{
    unsigned short StockIndex;//股票索引
    unsigned short Number; //数据条数，返回从最新往前的数据。一般服务器上提供60条数据。如果填1，那么会返回最新的一条数据
};

struct JWDDEHisData
{
    unsigned short StopBit; //停止位：0表示停止，>0表示还有数据
    unsigned int Date; //日期（格式：Unix格式时间数值）
    uint64_t BuyBigOrder; //买入大单 数/金额，服务端 x1000，客户端 /1000
    uint64_t SelBigOrder; //卖出大单 数/金额，服务端 x1000，客户端 /1000
    uint64_t BuyMinOrder; //买入小单 数/金额，服务端 x1000，客户端 /1000
    uint64_t SelMinOrder; //卖出小单 数/金额，服务端 x1000，客户端 /1000
    uint64_t BuyMidOrder; //买入中单 数/金额，服务端 x1000，客户端 /1000
    uint64_t SelMidOrder; //卖出中单 数/金额，服务端 x1000，客户端 /1000
    uint64_t BuyLargeOrder; //买入特大单 数/金额，服务端 x1000，客户端 /1000
    uint64_t SelLargeOrder; //卖出特大单 数/金额，服务端 x1000，客户端 /1000
    uint64_t BuyCount; //买入成交笔数，服务端 x1000，客户端 /1000
    uint64_t SellCount; //卖出成交笔数，服务端 x1000，客户端 /1000

};

struct JWDDEHisDataRsp:public JWCommonRsp
{
    unsigned short StockIndex; //股票索引
    unsigned int Number; //数据条数
    JWDDEHisData DDEHisData[JW_MAX_DDE_HIS_NUM];//历史DDE数据数组
};

////////////拉实时行情涨跌幅数据（?Business 服务器） 0x0A65/////////////////
struct JWRiseFallBreadthDataInfo
{
    unsigned short Type;//股票索引
    unsigned short Number; //数据条数，返回从最新往前的数据。一般服务器上提供60条数据。如果填1，那么会返回最新的一条数据
};

struct JWRiseFallBreadthData
{
    unsigned short StockIndex;//股票索引
    unsigned int LastPx;//现价，服务端 x1000，客户端 /1000
    unsigned int PreClosePx;//昨收价，服务端 x1000，客户端 /1000
    unsigned int BreadthValue;//幅度值（ 幅度值=((现价-昨收)/昨收)*100000 ），服务器 x100000 客户端/100000

};

struct JWRiseFallBreadthDataRsp:public JWCommonRsp
{
    unsigned short Type; //类型：1为涨幅，2为跌幅
    unsigned int Number; //数据条数，最大100条
    JWRiseFallBreadthData RiseFallBreadthData[JW_MAX_RISE_FALL_NUM];//行情涨跌幅数据结构数组
};


/////////////拉取综合诊断信息 0x0A40/////////////////
struct JWDiagnoseDataRsp:public JWCommonRsp
{
    unsigned short StockIndex; //股票索引
    char Content[JW_MAX_DIAGNOSE_CONTENT_LEN]; //诊断内容（utf-8编码）
};


//////////个股资金流向明细数据Business 服务器0x0A66////////

struct JWStkFundFlowDataDetailReq:public JWDdeDataInfo
{

};

struct JWStkFundFlowDetailData
{
    unsigned short StockIndex;//股票索引
    unsigned int Date;//日期（格式：Unix格式时间数值）
    uint64_t BuyBigOrderVolume;//买入大单量，服务端 x1000，客户端 /1000
    uint64_t BuyBigOrderAmount;//买入大单额，服务端 x1000，客户端 /1000
    uint64_t SelBigOrderVolume;//卖出大单量，服务端 x1000，客户端 /1000
    uint64_t SelBigOrderAmount;//卖出大单额，服务端 x1000，客户端 /1000
    uint64_t BuyMinOrderVolume;//买入小单量，服务端 x1000，客户端 /1000
    uint64_t BuyMinOrderAmount;//买入小单额，服务端 x1000，客户端 /1000
    uint64_t SelMinOrderVolume;//卖出小单量，服务端 x1000，客户端 /1000
    uint64_t SelMinOrderAmount;//卖出小单额，服务端 x1000，客户端 /1000
    uint64_t BuyMidOrderVolume;//买入中单量，服务端 x1000，客户端 /1000
    uint64_t BuyMidOrderAmount;//买入中单额，服务端 x1000，客户端 /1000
    uint64_t SelMidOrderVolume;//卖出中单量，服务端 x1000，客户端 /1000
    uint64_t SelMidOrderAmount;//卖出中单额，服务端 x1000，客户端 /1000
    uint64_t BuyLargeOrderVolume;//买入特大单量，服务端 x1000，客户端 /1000
    uint64_t BuyLargeOrderAmount;//买入特大单额，服务端 x1000，客户端 /1000
    uint64_t SelLargeOrderVolume;//卖出特大单量，服务端 x1000，客户端 /1000
    uint64_t SelLargeOrderAmount;//卖出特大单额，服务端 x1000，客户端 /1000
    uint64_t DealCountVolume;//成交笔数量，服务端 x1000，客户端 /1000
    uint64_t DealCountAmount;//成交额，服务端 x1000，客户端 /1000
};

struct JWStkFundFlowDetailDataRsp:public JWCommonRsp
{
    unsigned short StockNumber;//股票个数
    JWStkFundFlowDetailData StkFundFlowData[JW_MAX_STOCK_NUM];
};

/////////////////K线形态引擎数据 0x0A67///////////////////
struct JWKlineFormData
{
    unsigned short StockIndex;//股票索引
    char KLineType[JW_MAX_KLINE_TYPE];//K线形态类型
    unsigned int CalculateDate; //计算日期
    unsigned int AppearDate; //出现日期（格式：Unix格式时间数值）
    unsigned int FormTurnover;//形态换手率，服务器*1000，客户端/1000还原
};

struct JWKlineFormDataRsp:public JWCommonRsp
{
    unsigned short Type;//类型：0、最新看多形态；1、近期看多有效形态；2、最新看空形态；3、近期看空有效形态
    unsigned short Number;//请求个数，一般是200个，最大是800个；
    JWKlineFormData KlineFormData[JW_MAX_KLINE_NUM];//

};

struct JWKlineFormDataReq
{
    unsigned short Type;//类型：0、最新看多形态；1、近期看多有效形态；2、最新看空形态；3、近期看空有效形态
    unsigned short Number;//请求个数，一般是200个，最大是800个
};

////////////////多空博弈阵营数据 0x0A68///////////////
struct JWMultiEmptyGameData
{
    char ParamName[JW_MAX_PARAMNAME_LEN];//指标名称
    unsigned int Yield;//收益率，服务端 x1000，客户端 /1000
    unsigned int SuccessRate;//成功率，服务端 x1000，客户端 /1000
    unsigned short Type;//类型： 1：金叉/2：死叉
    unsigned int Date;//日期（格式：Unix格式时间数值）
    unsigned int PreClosePx;//收盘价，服务端 x1000，客户端 /1000
    unsigned int OpenPx;//开盘价，服务端 x1000，客户端 /1000
};

struct JWMultiEmptyGameDataRsp:public JWCommonRsp
{
    unsigned short StockIndex;//股票索引
    unsigned int StartCalculateDate;//开始计算日期（格式：Unix格式时间数值）
    unsigned short CalculateDateNum;//计算天数
    unsigned short Num;//内容的个数
    JWMultiEmptyGameData MultiEmptyGameData[JW_MAX_MULTI_EMPTY_GAME_CAMP];//内容的个数
};

//////////////////0x0A69 技术面-趋势分析数据(Business 服务器)0x0A69///////////////////////////////////
struct JWTechniccalTrendAnaysisReq
{
    unsigned short StockNumber; //股票个数
    unsigned short StockIndex[JW_MAX_STOCK_NUM];//股票索引数组
};

struct JWTechniccalTrendAnaysisData
{
    unsigned short StockIndex;          //股票索引
    unsigned int UpdateDate;            //日期（格式：Unix格式时间数值）
    unsigned int ShortSupportPrice;     //短线支撑位，服务端 x1000，客户端 /1000
    unsigned int ShortResistancePrice;  //短线阻力位，服务端 x1000，客户端 /1000
    unsigned int MidSupportPrice;       //中线支撑位，服务端 x1000，客户端 /1000
    unsigned int MidResistancePrice;    //中线阻力位，服务端 x1000，客户端 /1000
    char Glossary[JW_MAX_GLOSSARY_LEN]; //分析术语(MACD技术指标)（utf-8编码）
};

struct JWTechniccalTrendAnaysisRsp:public JWCommonRsp
{
    unsigned short StockNumber;//股票个数
    JWTechniccalTrendAnaysisData TechniccalTrendAnaysisData[JW_MAX_STOCK_NUM];
};

////////////////////// 0x0A70 资金面-主力分析-主力成本数据(Business 服务器)0x0A70//////////////////////////
struct JWFinancialMainCostReq:public JWTechniccalTrendAnaysisReq
{

};

struct JWFinancialMainCostContentData
{
    unsigned int UpdateDate; //日期（格式：Unix格式时间数值）
    unsigned int ClosePrice; //收盘价，服务端 x1000，客户端 /1000
    unsigned int MarketPrice;//10日市场平均交易成本，服务端 x1000，客户端 /1000
    unsigned int MainPrice;//10日主力平均交易成本，服务端 x1000，客户端 /1000

};

struct JWFinancialMainCostData
{   
    unsigned short StockIndex; //股票索引
    unsigned short Num; //内容的个数
    JWFinancialMainCostContentData FinancialMainCostContentData[JW_MAX_MAIN_COST_CONTENT_NUM];//内容的数组

};

struct JWFinancialMainCostRsp:public JWCommonRsp
{
    unsigned short StockNumber; //股票个数
    JWFinancialMainCostData FinancialMainCostData[JW_MAX_STOCK_NUM];//个股主力分析-主力成本数据数组
};


/////////////////////////////0x0a71 股票 WEB –技术面-市场表现数据(Business 服务器)0x0A71

struct JWTechniccalAndMarketReq:public JWTechniccalTrendAnaysisReq
{

};

struct JWTechniccalAndMarketData
{
    unsigned short StockIndex;              //股票索引
    unsigned int UpdateDate;                //日期（格式：Unix格式时间数值）
    unsigned int AccumulativeMarkup;        //20日累计涨幅，服务端 x100000，客户端 /100000
    unsigned int AccumulativeMarkupHY;      //20日行业累计涨幅，服务端 x100000，客户端 /100000
    unsigned int AccumulativeMarkupSHZS;    //20日上证指数累计涨幅，服务端 x100000，客户端 /100000
    unsigned int TurnOverRate1;             //1日换手率，服务端 x100000，客户端 /100000
    unsigned int TurnOverRate3;             //3日换手率，服务端 x100000，客户端 /100000
    unsigned int TurnOverRate5;             //5日换手率，服务端 x100000，客户端 /100000
};

struct JWTechniccalAndMarketRsp:public JWCommonRsp
{
    unsigned short StockNumber;
    JWTechniccalAndMarketData TechniccalAndMarketData[JW_MAX_STOCK_NUM];

};


//////////////////0x0a72 基本面-财务分析数据(Business 服务器)0x0A72//////////

struct JWFundAmentalsAnalysisReq:public JWTechniccalTrendAnaysisReq
{

};

struct JWFundAmentalsAnalysisData
{
    unsigned short StockIndex;       //股票索引
    unsigned int UpdateDate;         //日期（格式：Unix格式时间数值）
    uint64_t YYSR;                   //营业收入，服务端 x1000，客户端 /1000
    unsigned int YYSR_TB_ZZ;         //营业收入同比增长，服务端 x1000，客户端 /1000
    uint64_t JLR_PARENT;             //净利润，服务端 x1000，客户端 /1000
    unsigned int JLR_TB_ZZ;          //净利润同比增长，服务端 x1000，客户端 /1000
    unsigned int MGSY;               //每股收益，服务端 x1000，客户端 /1000
};

struct JWFundAmentalsAnalysisRsp:public JWCommonRsp
{
    unsigned short StockNumber;//股票个数
    JWFundAmentalsAnalysisData FundAmentalsAnalysisData[JW_MAX_STOCK_NUM];//财务分析数据数组

};

/////////////////////0x0A73 基本面-财务分析-利润指标数据(Business 服务器)0x0A73
struct JWFundAmentalsTargetsReq:public JWTechniccalTrendAnaysisReq
{

};

struct JWFundAmentalsTargetsContentData
{
    unsigned int UpdateDate;        //日期（格式：Unix格式时间数值）
    unsigned int YYSR_TB_ZZL;       //营业收入同比增长率，服务端 x100000，客户端 /100000
    unsigned int JLR_TB_ZZL;        //净利润同比增长率，服务端 x100000，客户端 /100000


};

struct JWFundAmentalsTargetsData
{
    unsigned short StockIndex;                                          //股票索引
    unsigned short Num;                                                 //内容的个数
    JWFundAmentalsTargetsContentData FundAmentalsTargetsContentData[JW_MAX_TARGETS_CONTENT_NUM];  //内容数组
};

struct JWFundAmentalsTargetsRsp:public JWCommonRsp
{
    unsigned short StockNumber;
    JWFundAmentalsTargetsData FundAmentalsTargetsData[JW_MAX_STOCK_NUM];

};

/////////////////////////  0x0A74 基本面-财务分析-行业排名数据(Business 服务器) 0x0A74

struct JWFundIndustryRanksReq:public JWTechniccalTrendAnaysisReq
{

};

struct JWFundIndustryRanksData
{
    unsigned short StockIndex;          //股票索引
    unsigned int UpdateDate;            //日期（格式：Unix格式时间数值）
    char IndustryName[JW_MAX_INDUSTRY_NAME_LEN];              //所属行业名称
    uint64_t ZZC;                       //个股总资产，服务端 x1000，客户端 /1000
    uint64_t ZZC_HY;                    //行业总资产均值，服务端 x1000，客户端 /1000
    unsigned short ZZC_PM;              //个股总资产排名
    unsigned int JLRL;                  //个股净利润率，服务端 x100000，客户端 /100000
    unsigned int JLRL_HY;               //行业净利润率，服务端 x100000，客户端 /100000
    unsigned short JLRL_PM;             //个股净利润率排名
    unsigned int JZCSYL;                //个股净资产收益率，服务端 x100000，客户端 /100000
    unsigned int JZCSYL_HY;             //行业净资产收益率，服务端 x100000，客户端 /100000
    unsigned short JZCSYL_PM;           //个股净资产收益率排名
    unsigned int XSMLL;                 //个股销售毛利率，服务端 x100000，客户端 /100000
    unsigned int XSMLL_HY;              //行业销售毛利率，服务端 x100000，客户端 /100000
    unsigned short XSMLL_PM;            //个股销售毛利率排名
    unsigned short InduElementNum;      //所属行业成份股总数
};

struct JWFundIndustryRanksRsp:public JWCommonRsp
{
    unsigned short StockNumber;//股票个数
    JWFundIndustryRanksData FundIndustryRanksData[JW_MAX_STOCK_NUM];//行业排名数据数组
};

///////////////////////////////0x0A75 基本面-估值研究数据(Business 服务器)0x0A75////////////
struct JWFundAluationReserchReq:public JWTechniccalTrendAnaysisReq
{

};

struct JWFundAluationReserchData
{
    unsigned short StockIndex;      //股票索引
    unsigned int UpdateDate;        //日期（格式：Unix格式时间数值）
    unsigned int PE_TTM;            //市盈率，服务端 x100000，客户端 /100000
    unsigned int PE_TTM_HY;         //行业市盈率，服务端 x100000，客户端 /100000
    unsigned int PE_TTM_MARKET;     //市场市盈率，服务端 x100000，客户端 /100000
    char PE_TTM_RANKING[JW_MAX_RANK_LEN];          //市盈率排名
    unsigned int PB;                //市净率，服务端 x100000，客户端 /100000
    unsigned int PB_HY;             //行业市净率，服务端 x100000，客户端 /100000
    unsigned int PB_MARKET;         //市场市净率，服务端 x100000，客户端 /100000
    char PB_RANKING[JW_MAX_RANK_LEN];              //市净率排名
    unsigned int PEG;               //PEG，服务端 x100000，客户端 /100000
    unsigned int PEG_HY;            //行业PEG，服务端 x100000，客户端 /100000
    unsigned int PEG_MARKET;        //市场PEG，服务端 x100000，客户端 /100000
    char PEG_RANKING[JW_MAX_RANK_LEN];             //PEG排名
};

struct JWFundAluationReserchRsp:public JWCommonRsp
{
    unsigned short StockNumber;                         //股票个数
    JWFundAluationReserchData FundAluationRserchData[JW_MAX_STOCK_NUM];  //估值研究数据数组
};

////////////////////////////////0x0A76 机构观点数据(Business 服务器)0x0A76/////////
struct JWFundOrganizationViewReq:public JWTechniccalTrendAnaysisReq
{

};

struct JWFundOrganizationViewData
{
    unsigned short StockIndex;                     //股票索引
    unsigned int UpdateDate;                       //日期（格式：Unix格式时间数值）
    unsigned short OrganNum;                       //机构数量
    unsigned short StatDays;                       //统计天数
    unsigned int Rating;                           //机构评级，服务端 x100000，客户端 /100000
    unsigned int Attention;                        //机构关注度，服务端 x100000，客户端 /100000
    unsigned short BuyNum;                         //买入数量
    unsigned short HoldingsNum;                    //增持数量
    unsigned short NeutralNum;                     //中性数量
    unsigned short ReductionNum;                   //减持数量
    unsigned short SellNum;                        //卖出数量
    char PE_TTM_RANKING[JW_MAX_RANK_LEN];          //市盈率排名
    unsigned int PB;                               //市净率，服务端 x100000，客户端 /100000
    unsigned int PB_HY;                            //行业市净率，服务端 x100000，客户端 /100000
    char PB_RANKING[JW_MAX_RANK_LEN];              //市净率排名
    unsigned int PEG;                              //PEG，服务端 x100000，客户端 /100000
    unsigned int PEG_HY;                           //行业PEG，服务端 x100000，客户端 /100000
    char PEG_RANKING[JW_MAX_RANK_LEN];             //PEG排名
};

struct JWFundOrganizationViewRsp:public JWCommonRsp
{
    unsigned short StockNumber;      //股票个数
    JWFundOrganizationViewData FundOrganizationViewData[JW_MAX_STOCK_NUM]; //机构观点数据数组

};


#pragma pack(pop)



#endif

