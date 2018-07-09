/********************************************
//文件名:jw_define.h
//功能:系统定义头文件
//作者:钟何明
//创建时间:2010.07.07
//修改记录:

*********************************************/

#ifndef __JW_DEFINE_H__
#define __JW_DEFINE_H__

namespace jw_sourcetype
{    
    /*来源类型定义7bit,高1bit为Server内部使用*/
    //========Client 1-31,0暂时保留=============
    const unsigned char SOURCE_TYPE_MIN_IM = 1;
    const unsigned char SOURCE_TYPE_MAX_IM = 31;

    const unsigned char SOURCE_TYPE_IM_CLIENT = 1;//IM
    const unsigned char SOURCE_TYPE_HUAGUJIE = 2;//华股老街
    const unsigned char SOURCE_TYPE_WEB_GBS_IMPORT = 3;//WEB 股博士导入的帐号

    //========Web 32-63=======================
    const unsigned char SOURCE_TYPE_MIN_WEB = 32;
    const unsigned char SOURCE_TYPE_MAX_WEB = 63;

    const unsigned char SOURCE_TYPE_WEB_SNS = 32;//SNS社区
    const unsigned char SOURCE_TYPE_WEB_HOME = 33;//788111门户
    const unsigned char SOURCE_TYPE_WEB_TAO = 34;//淘股天堂
    const unsigned char SOURCE_TYPE_WEB_RT = 35;//炒股大赛
    const unsigned char SOURCE_TYPE_WEB_GBS = 36;//WEB股博士

    //=========DYJ Client 64-127==================
    const unsigned char SOURCE_TYPE_DYJ_MIN_CLIENT = 64;
    const unsigned char SOURCE_TYPE_DYJ_MAX_CLIENT = 127;

    const unsigned char SOURCE_TYPE_DYJ_CLIENT_XGW = 64;//选股王
    const unsigned char SOURCE_TYPE_DYJ_CLIENT_PCW = 66;//评测王
    const unsigned char SOURCE_TYPE_DYJ_CLIENT_GGW = 68;//股票管家
    const unsigned char SOURCE_TYPE_DYJ_CLIENT_GBS = 70;//股博士
    const unsigned char SOURCE_TYPE_DYJ_CLIENT_MOBILE = 72;//手机终端
    const unsigned char SOURCE_TYPE_DYJ_CLIENT_TAO = 74;//淘股天堂等金融终端

    const unsigned char SOURCE_TYPE_STOCK_GAME_MANUAL = 126;//炒股大赛手工开通的帐号
    const unsigned char SOURCE_TYPE_DYJ_CLIENT_MANUAL = 127;//手工开通的帐号

    //====================================
}

namespace jw_command
{

    const unsigned short JW_SUB_DEFAULT     = 0x0000;//默认子命令字

    //注意：以下没写子命令字的，子命令字默认是0
    const unsigned short JW_KEEP_ALIVE      = 0x0005;//心跳
    const unsigned short JW_REG_PRE_REG = 0x0801;   //注册账户
    const unsigned short JW_REG_ACTIVATE = 0x0802; //激活账户
    const unsigned short JW_ACCOUNT_STATUS  = 0x0803;//查询帐号状态
    const unsigned short JW_REG_GET_REGINFO = 0x0804;//取用户注册信息
    const unsigned short JW_ACCOUNT_REG_TIME =0x0806;//取得帐号注册时间
    const unsigned short JW_CREATE_ACCOUNT  = 0x0809;//开通帐号
    const unsigned short JW_STAT_REPORT = 0x0810;//用户数据上报

    const unsigned short JW_AUTH_LOGIN      = 0x0400;//认证登录
    const unsigned short JW_AUTH_CHECK_TICKET = 0x0401;//验证票据
    const unsigned short JW_AUTH_LOGOUT     = 0x0402;//注销登陆
    const unsigned short JW_AUTH_MODI_PWD = 0x0403;//修改用户密码
    const unsigned short JW_AUTH_UPDATE_TICKET = 0x0405;//更新用户Ticket
    const unsigned short JW_AUTH_GET_TICKET_UIN = 0x0406;//通过UIN取得认证票据
    const unsigned short JW_AUTH_RESET_PWD     = 0x0409;//重置用户密码
    

    const unsigned short JW_STOCK           = 0x000F;//自选股主命令字
        const unsigned short JW_SUB_STOCK_ADD                   = 0x0001;//增加自选股
        const unsigned short JW_SUB_STOCK_DEL                   = 0x0002;//删除自选股
        const unsigned short JW_SUB_STOCK_LIST_BY_GROUPID       = 0x0010;//拉取自选股列表
        const unsigned short JW_SUB_SET_STOCK_SORT              = 0x0011;//设置自选股排序规则
        const unsigned short JW_SUB_GET_STOCK_SORT              = 0x0012;//拉取自选股排序规则
        const unsigned short JW_SUB_GET_STOCK_RANK              = 0x0013;//拉取股票排名顺序

    const unsigned short JW_BIND_MOBILE = 0x5030;//绑定手机业务主命令字
#if 0
        const unsigned short JW_SUB_BIND_UNBIND_APPLY           = 0x0001;//绑定/解除手机绑定业务
        const unsigned short JW_SUB_CHECK_BIND_UNBIND           = 0x0002;//校验绑定(解除绑定)手机
        const unsigned short JW_SUB_GET_BIND_LIST               = 0x0003;//拉取用户绑定的手机信息
#endif
        const unsigned short JW_SUB_BIND_UNBIND                 = 0X0004;//绑定/解除手机绑定业务
        const unsigned short JW_SUB_CHECK_VERIFY_CODE           = 0x0005;//校验绑定(解除绑定)手机验证码
        const unsigned short JW_SUB_GET_BIND_INFO               = 0x0006;//拉取用户绑定的手机信息

        const unsigned short JW_INFO            =   0x5000;//用户资料
        const unsigned short JW_SUB_GET_COMM_INFO = 0x0010;//取用户公共资料信息及头像URL
        const unsigned short JW_SUB_SET_COMM_INFO = 0x0011;//设置用户公共资料信息
        const unsigned short JW_SUB_GET_MOBILE_AREA=0x0012;//通过用户ID查询手机帐号的归宿地
        const unsigned short JW_SUB_GET_MOBILE_AREA_BY_NUMBER = 0x0013;//通过手机号码查询归宿地
        const unsigned short JW_SUB_GET_USER_ACCOUNT_INFO = 0x0014;//通过UIN取得用户帐号信息
        const unsigned short JW_SUB_SET_ALIAS = 0x0018;//设置用户的别名信息
        const unsigned short JW_SUB_GET_ALIAS = 0x0019;//取得用户的别名信息
        
        const unsigned short JW_PAY_MNG             = 0x5060;
        const unsigned short JW_SUB_PAY_TICKETINCOME_NOTIFY = 0x001B;


////////////////////////////////////// 股票服务器命令字定义/////////////////////////////////
        const unsigned short JW_STK_WEB_PULL_MARKET_DATA = 0x0A09;  //拉取最新行情数据
        const unsigned short JW_STK_WEB_PULL_MINUTE_DATA = 0x0A0A;  //拉分时数据线
        const unsigned short JW_STK_WEB_RADAR_INDICATOR  = 0x0A10; //个股雷达图指标
        const unsigned short JW_STK_WEB_EARLY_WARNING  =  0X0A1D;  //预警行情数据
        const unsigned short JW_STK_WEB_BATCH_BRIEF_MARKET_DATA = 0x0A22; //批量简要的行情数据
        const unsigned short JW_STK_WEB_NEWS_CONTENT  =  0x0A28;  //单只股票新闻内容数据
        const unsigned short JW_STK_WEB_BATCH_NEWS_CONTENT  = 0x0A29; //批量拉取单只股票最新新闻标题
        const unsigned short JW_STK_HOME_GET_REFERENCE_DATA = 0x0A60; //取reference data
        const unsigned short JW_STK_WEB_COST_DISTRIBUTION_INDICATOR = 0x0A17;//个股成本分布指标（? Business 服务器）0x0A17

        //add by liudaile at 2010.10.20
        const unsigned short JW_STK_WEB_COST_DISTRIBUTION_DETAIL = 0x0A18; //个股成本分布明细（? Business 服务器）0x0A18
        const unsigned short JW_STK_WEB_STOCK_FUND_FLOW = 0x0A50; //拉个股资金流向数据
        const unsigned short JW_STK_WEB_INDUSTRY_FUND_FLOW = 0x0A53; //拉行业资金流向数据
        const unsigned short JW_STK_WEB_DDE_DATA = 0x0A61; // 拉DDE数据（?Business 服务器） 0x0A61
        const unsigned short JW_STK_WEB_DDSORT_DATA= 0x0A62; //拉DDSort数据（?Business 服务器） 0x0A62
        const unsigned short JW_STK_WEB_DDE_HIS_DATA = 0x0A64; //拉DDE历史数据（?Business 服务器） 0x0A64
        const unsigned short JW_STK_WEB_RISE_FALL_BREADTH= 0x0A65;//拉实时行情涨跌幅数据（?Business 服务器） 0x0A65
        const unsigned short JW_STK_WEB_ZHZD_INFO = 0x0A40; //拉取综合诊断信息(Business 服务器） 0x0A40
        //end add
        //add by liudaile at 2010.11.11
        const unsigned short JW_STK_WEB_STOCK_FUND_DETAIL = 0x0A66;//个股资金流向明细数据  0x0A66
        const unsigned short JW_STK_WEB_KLINE_FORM_ENGINE = 0x0A67;//K线形态引擎数据0x0A67
        const unsigned short JW_STK_WEB_MULTI_EMPTY_GAME_CAMP = 0x0A68;//多空博弈阵营数据(Business 服务器) 0x0A68

        const unsigned short JW_STK_WEB_TECHNICAL_TREND_ANALYSIS = 0x0A69;//技术面-趋势分析数据
        const unsigned short JW_STK_WEB_FINANCIAL_MAIN_ANALYSIS_COST = 0x0A70;//资金面-主力分析-主力成本数据
        const unsigned short JW_STK_WEB_TECHNICAL_MARKET_PERFORMANCE = 0x0A71; //技术面-市场表现数据
        const unsigned short JW_STK_WEB_FUNDAMENTALS_FINANCIAL_ANALYSIS = 0x0A72;//基本面-财务分析数据
        const unsigned short JW_STK_WEB_FUNDAMENTALS_FA_PROFIT_TARGETS = 0x0A73;//基本面-财务分析-利润指标数据
        const unsigned short JW_STK_WEB_FUNDAMENTALS_FA_INDUSTRY_RANKINGS = 0x0A74;//基本面-财务分析-行业排名数据
        const unsigned short JW_STK_WEB_FUNDAMENTALS_ALUATION_RESEARCH = 0x0A75;//基本面-估值研究数据
        const unsigned short JW_STK_WEB_FUNDAMENTALS_ORGANIZATION_VIEW = 0x0A76;//基本面-机构观点数据
////////////////////////////////////// 结束股票服务器命令字定义/////////////////////////////////
}

namespace jw_bind_type
{
    //绑定类型定义
    const unsigned char JW_BIND_TYPE_EW_GBS = 1;//股博士接收预警消息的手机号码

}


namespace jw_errorcode
{
    /*全局错误码定义*/    
    /*    
    const unsigned int RET_OK = 0;//成功   
    const unsigned int ERR_PROTOCOL = 1;//协议错误    
    const unsigned int ERR_ENC_TYPE = 2;//不支持的加密方式
    const unsigned int ERR_PROTOCOL_VER = 3;//不支持的版本号
    const unsigned int ERR_DB = 4;//操作数据库错误
    const unsigned int ERR_NETWORK = 5;//网络错误
    const unsigned int IM_ERR_UNKNOWN = 0x0000FFFF;//未知错误
    */
    
    /* 认证相关错误码定义*/    
    const unsigned int RET_AUTH_ACCT_NOTEXIST = 0x00080001;//帐号不存在
    const unsigned int RET_AUTH_ACCT_DISABLE = 0x00080002;//帐号被禁用
    const unsigned int RET_AUTH_ACCT_NOTACT = 0x00080003;//帐号未激活
    const unsigned int RET_AUTH_ACCT_FREEZ = 0x00080004;//帐号被冻结(临时限制登陆)
    const unsigned int RET_AUTH_PWD_ERROR = 0x00080005;//密码错误
    const unsigned int RET_AUTH_TICKET_INVALID = 0x00080006;//票据不正确
    const unsigned int RET_AUTH_TICKET_EXPIRE = 0x00080007;//票据已过期
    const unsigned int RET_AUTH_SYSTEM_ERROR = 0x00080008;//认证服务器错误
    const unsigned int RET_AUTH_PROTOCOL_ERROR = 0x00080009;//认证协议不正确    
    const unsigned int RET_AUTH_DP_ERROR = 0x0008000b;//数据平台错误

    /* 注册相关错误码定义*/
    const unsigned int REG_ERR_VERIFY               = 0x00150001;  //激活验证失败
    const unsigned int REG_ERR_ALREADY_REG          = 0x00150002;  //已注册了
    const unsigned int REG_ERR_NEED_ACTIVATE        = 0x00150003;  //需要激活
    const unsigned int REG_ERR_NO_AVA_ID            = 0x00150004;  //没有可用ID了
    const unsigned int REG_ERR_DB                   = 0x00150005;  //DB错误
    const unsigned int REG_ERR_ACTIVATE_NOT_FOUND   = 0x00150006;  //没有找到激活信息
    const unsigned int REG_ERR_ACTIVATE_VERIFY_ERR  = 0x00150007;  //激活码不正确
    const unsigned int REG_ERR_CHGACT_PWD_ERR       = 0x00150008;  //修改账号时密码错误
    const unsigned int REG_ERR_CHGACT_ACT_ERR       = 0x00150009;  //修改账号时新账号已经被占用
    const unsigned int REG_ERR_CHGACT_TIME_ERR      = 0x0015000A;  //修改账号时密码设置时间未到
    const unsigned int REG_ERR_CHGACT_ARG_ERR       = 0x0015000B;  //修改账号时送过来的参数错误.uin不对
    const unsigned int REG_ERR_CHGACT_CTIME_ERR     = 0x0015000C;  //修改账号时账号修改时间过短
    const unsigned int REG_ERR_REPEAT_REG           = 0x0015000D;//帐号已注册
    const unsigned int REG_ERR_INVALID_ACCOUNT      = 0x0015000E;//无效的帐号


    /*资料相关错误码定义*/
    const unsigned int RET_INFO_SYS_ERR = 0x00020001; //系统错误
    const unsigned int RET_INFO_FIELD_TYPE_ERR = 0x00020002;//字段类型错误
    const unsigned int RET_INFO_FIELD_LEN_ERR = 0x00020003;//字段长度错误
    const unsigned int RET_INFO_MOBILE_ERR = 0x00020004;//错误的手机号码格式
    const unsigned int RET_INFO_MOBILE_NOT_FOUND = 0x00020005;//没有找到对应的归属地

    /*自选股相关错误码定义*/  
    const unsigned int RET_STK_SYS_ERR = 0x00090001; //系统错误
    const unsigned int RET_STK_DB_ERR = 0x00090002;//DB错误
    const unsigned int RET_STK_PROTOCOL_ERR = 0x00090003;//协议错误
    const unsigned int RET_STK_COUNT_MAX_LIMIT_ERR = 0x00090004; //达到系统允许的数量上限
    const unsigned int RET_STK_STRING_LEN_LIMIT = 0x00090005; //超出定义的string字段的长度
    const unsigned int RET_STK_NOTFOUND = 0x00090006;//没找到
    const unsigned int RET_STK_OR_GRPUP_ERROR = 0x00090007;//股票信息或者自选股分组不存在
    const unsigned int RET_STK_REPEATED = 0x00090008;//股票已经存在
    const unsigned int RET_STK_DEFALUT_ERR = 0x00090009;//默认分组不能更改或者删除
    const unsigned int RET_STK_NO_STKRANK_ERR = 0x0009000A;//没有对应的股票排名
    const unsigned int RET_STK_NO_STKSORT = 0x0009000B;//没有自选股排序规则

    /*手机绑定相关错误码定义*/
    const unsigned int RET_BIND_SYSTEM_ERR           = 0x02030001;//系统错误
    const unsigned int RET_BIND_PROTOCOL_ERR         = 0x02030002;//协议错误
    const unsigned int RET_BIND_DB_ERR               = 0x02030003;//DB错误
    const unsigned int RET_BIND_UNSUPPORT_ERR        = 0x02030004;//不支持的数据格式
    const unsigned int RET_BIND_ACCOUNT_ERR          = 0x02030005;//账号不存在
    const unsigned int RET_BIND_INVALID_CODE_ERR     = 0x02030006;//无效的验证码
    const unsigned int RET_BIND_NUMBER_LIMIT_ERR     = 0x02030007;//超过绑定号码个数
    const unsigned int RET_BIND_MOBILE_UNBIND_ERR    = 0x02030008;//要解除绑定的手机没有绑定任何应用
    const unsigned int RET_BIND_MOBILE_BINDED_ERR    = 0x02030009; //该手机号码已经被其它用户绑定
}


namespace jw_ticket_type
{
    const unsigned char TICKET_TYPE_JUMP_HUAGU_URL = 2;//跳转华股URL票据
    const unsigned char TICKET_TYPE_CLIENT_GBS = 128;//股博士认证票据
    const unsigned char TICKET_TYPE_TERMINAL = 127;//终端认证票据
    
}


#endif
