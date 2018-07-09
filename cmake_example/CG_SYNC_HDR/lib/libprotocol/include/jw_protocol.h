
#ifndef _JW_PROTOCOL_H_
#define _JW_PROTOCOL_H_

#include "jw_os.h"
#include "jw_struct.h"



/******************************************************************
*函数名:JWInit
*功能说明：初始化
*参数说明：
*   无
*返回值：
*   无
*****************************************************************/
JW_PROTOCOL_API void  JWInit();

/******************************************************************
*函数名:JWUnInit
*功能说明：资源析构
*参数说明：
*   无
*返回值：
*   无
*****************************************************************/
JW_PROTOCOL_API void  JWUnInit();


/******************************************************************
*函数名:JWGetModuleVer
*功能说明：取得模块版本号
*参数说明：
*   无
*返回值：
*   模块版本信息(高2字节表示主版本号，低2字节表示子版本号)
*****************************************************************/
JW_PROTOCOL_API unsigned int  JWGetModuleVer();



/******************************************************************
*函数名:JWHandleInput
*功能说明：截取一个完整的协议包
*参数说明：
*   IN buf:接收的数据流指针
*   IN len:数据流长度
*返回值：
*   >0:表示数据包已接收完整且该值表示数据包的长度
*   0：表示数据包还未接收完整
*   <0：表示出错
*****************************************************************/
JW_PROTOCOL_API int  JWHandleInput(const char * buf, unsigned int len);


/******************************************************************
*函数名:JWGetPackLen
*功能说明：取得数据包的长度
*参数说明：
*   IN p:完整数据包指针
*返回值：
*   数据包的长度
*****************************************************************/
JW_PROTOCOL_API unsigned int  JWGetPackLen(const char * p);

/******************************************************************
*函数名:JWGetPackCmd
*功能说明：取得数据包的命令字
*参数说明：
*   IN p:完整数据包指针
*   OUT main_cmd:主命令字
*   OUT sub_cmd:子命令字
*返回值：
*   0:成功
*   <0：表示出错
*****************************************************************/
JW_PROTOCOL_API int  JWGetPackCmd(const char * p, unsigned short * main_cmd, unsigned short * sub_cmd);

/******************************************************************
*函数名:JWGetPackSeq
*功能说明：取得数据包的序号
*参数说明：
*   IN p:完整数据包指针
*返回值：
*   数据包的序号
*****************************************************************/
JW_PROTOCOL_API int  JWGetPackSeq(const char * p);

/******************************************************************
*函数名:JWSetPackSeq
*功能说明：设置数据包的序号
*参数说明：
*   IN p:完整数据包指针
*   IN seq:数据包序号
*返回值：
*   无
*****************************************************************/
JW_PROTOCOL_API void  JWSetPackSeq(const char * p, unsigned int seq);


/******************************************************************
*函数名:JWGetSourceType
*功能说明：取得数据包的来源类型标志
*参数说明：
*   IN p:完整数据包指针
*返回值：
*   数据包的来源类型标志
*****************************************************************/
JW_PROTOCOL_API unsigned char  JWGetSourceType(const char * p);

/******************************************************************
*函数名:JWSetSourceType
*功能说明：设置数据包的来源类型标志
*参数说明：
*   IN p:完整数据包指针
*   IN source_type:数据包来源类型标志
*返回值：
*   无
*****************************************************************/
JW_PROTOCOL_API void  JWSetSourceType(const char * p, unsigned char source_type);

/******************************************************************
*函数名:JWEncryptPack
*功能说明：加密数据包
*参数说明：
*   IN OUT p:完整数据包指针
*   IN len:数据包长度
*返回值：
*   无
*****************************************************************/
JW_PROTOCOL_API void  JWEncryptPack(char * p, unsigned int len);


/******************************************************************
*函数名:JWDecryptPack
*功能说明：解密数据包
*参数说明：
*   IN OUT p:完整数据包指针
*   IN len:数据包长度
*返回值：
*   无
*****************************************************************/
JW_PROTOCOL_API void  JWDecryptPack(char * p,unsigned int len);



/******************************************************************
*函数名:JWPackKeepAliveReq
*功能说明：生成心跳信息的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:请求的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackKeepAliveReq(char * buf, int buf_len, const JWCommonReqInfo * pInfo);


/******************************************************************
*函数名:JWUnPackKeepAliveRsp
*功能说明：解析心跳信息的应答包
*参数说明：
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackKeepAliveRsp(const char * p, int len, JWCommonRsp * pRsp);


/******************************************************************
*函数名:JWPackGetAccountStatusReq
*功能说明：生成查询帐号状态的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:查询帐号的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetAccountStatusReq(char * buf, int buf_len, const JWGetAccountStatusInfo * pInfo);

/******************************************************************
*函数名:JWUnPackGetAccountStatusRsp
*功能说明：解析查询帐号状态的应答包
*参数说明：
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:查询帐号应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetAccountStatusRsp(const char * p, int len, JWGetAccountStatusRsp * pRsp);


/******************************************************************
*函数名:JWPackCreateAccountReq
*功能说明：生成开通帐号的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:开通帐号的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackCreateAccountReq(char * buf, int buf_len, const JWCreateAccountInfo * pInfo);

/******************************************************************
*函数名:JWUnPackCreateAccountRsp
*功能说明：解析开通帐号的应答包
*参数说明：
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:开通帐号应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackCreateAccountRsp(const char * p, int len, JWCreateAccountRsp * pRsp);


/******************************************************************
*函数名:JWPackLoginReq
*功能说明：生成登录的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:登陆参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackLoginReq(char * buf, int buf_len, const JWLoginInfo * pInfo);


/******************************************************************
*函数名:JWUnPackLoginRsp
*功能说明：解析登录的应答包
*参数说明：
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:登陆应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackLoginRsp(const char * p, int len, JWLoginRsp * pRsp);

/******************************************************************
*函数名:JWPackLogoutReq
*功能说明：生成注销的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:注销参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackLogoutReq(char * buf, int buf_len, const JWLogoutInfo * pInfo);

/******************************************************************
*函数名:JWUnPackLogoutRsp
*功能说明：解析注销的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:注销应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackLogoutRsp(const char * p, int len, JWLogoutRsp * pRsp);



/******************************************************************
*函数名:JWPackAddStockReq
*功能说明：生成增加自选股的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:增加自选股的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackAddStockReq(char * buf, int buf_len, const JWAddStockInfo * pInfo);

/******************************************************************
*函数名:JWUnPackAddStockRsp
*功能说明：解析增加自选股的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:增加自选股的应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackAddStockRsp(const char * p, int len, JWAddStockRsp * pRsp);


/******************************************************************
*函数名:JWPackDelStockReq
*功能说明：生成删除自选股的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:删除自选股的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackDelStockReq(char * buf, int buf_len, const JWDelStockInfo * pInfo);

/******************************************************************
*函数名:JWUnPackDelStockRsp
*功能说明：解析删除自选股的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:删除自选股的应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackDelStockRsp(const char * p, int len, JWDelStockRsp * pRsp);

/******************************************************************
*函数名:JWPackGetStockListReq
*功能说明：生成拉取自选股列表的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:拉取自选股列表的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetStockListReq(char * buf, int buf_len, const JWGetStockInfo * pInfo);

/******************************************************************
*函数名:JWUnPackGetStockListRsp
*功能说明：解析拉取自选股列表的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:拉取自选股列表的应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetStockListRsp(const char * p, int len, JWGetStockInfoRsp * pRsp);

/******************************************************************
*函数名:JWPackSetStockSortReq
*功能说明：生成设置自选股排序的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:设置自选股排序的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackSetStockSortReq(char * buf, int buf_len, const JWSetStockSortInfo * pInfo);

/******************************************************************
*函数名:JWUnPackSetStockSortRsp
*功能说明：解析设置自选股排序的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:设置自选股排序的应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackSetStockSortRsp(const char * p, int len, JWSetStockSortRsp * pRsp);

/******************************************************************
*函数名:JWPackGetStockSortReq
*功能说明：生成拉取自选股排序的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:拉取自选股排序的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetStockSortReq(char * buf, int buf_len, const JWGetStockSortInfo * pInfo);

/******************************************************************
*函数名:JWUnPackGetStockSortRsp
*功能说明：解析拉取自选股排序的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:拉取自选股排序的应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetStockSortRsp(const char * p, int len, JWGetStockSortRsp * pRsp);

/******************************************************************
*函数名:JWPackGetStockRankReq
*功能说明：生成拉取自选股排名的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:拉取自选股排名的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetStockRankReq(char * buf, int buf_len, const JWGetStockRankInfo * pInfo);

/******************************************************************
*函数名:JWUnPackGetStockRankRsp
*功能说明：解析拉取自选股排名的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:拉取自选股排名的应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetStockRankRsp(const char * p, int len, JWGetStockRankRsp * pRsp);

#if 0
//2010-11-26 zhongheming 业务变更，废弃掉旧接口 
/******************************************************************
*函数名:JWPackBindMobileReq
*功能说明：生成绑定/解除绑定手机号码的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:绑定/解除绑定手机号码的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackBindMobileReq(char * buf, int buf_len, const JWBindMobileInfo * pInfo);

/******************************************************************
*函数名:JWUnPackBindMobileRsp
*功能说明：解析绑定/解除绑定手机号码的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:绑定/解除绑定手机号码的应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackBindMobileRsp(const char * p, int len, JWBindMobileRsp * pRsp);

/******************************************************************
*函数名:JWPackBindVerifyReq
*功能说明：生成校验 绑定/解除绑定 验证码的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:绑定/解除绑定验证码的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackBindVerifyReq(char * buf, int buf_len, const JWBindVerifyInfo * pInfo);

/******************************************************************
*函数名:JWUnPackBindVerifyRsp
*功能说明：解析校验绑定/解除绑定手机号码验证码的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:校验绑定/解除绑定手机号码验证码的应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackBindVerifyRsp(const char * p, int len, JWBindVerifyRsp * pRsp);



/******************************************************************
*函数名:JWPackGetBindListReq
*功能说明：生成取得绑定的手机信息请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:取得绑定的手机信息的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetBindListReq(char * buf, int buf_len, const JWGetBindInfo * pInfo);

/******************************************************************
*函数名:JWUnPackGetBindListRsp
*功能说明：解析取得绑定的手机信息的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:取得绑定的手机信息的应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetBindListRsp(const char * p, int len, JWGetBindInfoRsp * pRsp);

#endif

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/******************************************************************
*函数名:JWPackBindMobileInfoReq
*功能说明：生成绑定/解除绑定手机号码的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:绑定/解除绑定手机号码的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackBindMobileInfoReq(char * buf, int buf_len, const JWBindMobileInfoReq * pInfo);

/******************************************************************
*函数名:JWUnPackBindMobileInfoRsp
*功能说明：解析绑定/解除绑定手机号码的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:绑定/解除绑定手机号码的应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackBindMobileInfoRsp(const char * p, int len, JWBindMobileInfoRsp * pRsp);

/******************************************************************
*函数名:JWPackBindVerifyInfoReq
*功能说明：生成校验 绑定/解除绑定 验证码的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:绑定/解除绑定验证码的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackBindVerifyInfoReq(char * buf, int buf_len, const JWBindVerifyInfoReq * pInfo);

/******************************************************************
*函数名:JWUnPackBindVerifyInfoRsp
*功能说明：解析校验绑定/解除绑定手机号码验证码的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:校验绑定/解除绑定手机号码验证码的应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackBindVerifyInfoRsp(const char * p, int len, JWBindVerifyInfoRsp * pRsp);



/******************************************************************
*函数名:JWPackGetBindMobileInfoReq
*功能说明：生成取得绑定的手机信息请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:取得绑定的手机信息的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetBindMobileInfoReq(char * buf, int buf_len, const JWGetBindMobileInfoReq * pInfo);

/******************************************************************
*函数名:JWUnPackGetBindMobileInfoRsp
*功能说明：解析取得绑定的手机信息的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:取得绑定的手机信息的应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetBindMobileInfoRsp(const char * p, int len, JWGetBindMobileInfoRsp * pRsp);


/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/******************************************************************
*函数名:JWPackGetBasicInfoReq
*功能说明：生成取得用户基本信息的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:取得用户基本信息的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetBasicInfoReq(char * buf, int buf_len, const JWGetBasicInfo * pInfo);

/******************************************************************
*函数名:JWUnPackGetBasicInfoRsp
*功能说明：解析取得用户基本信息的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:取得用户基本信息的应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetBasicInfoRsp(const char * p, int len, JWGetBasicInfoRsp * pRsp);

/******************************************************************
*函数名:JWPackSetBasicInfoReq
*功能说明：生成设置用户基本信息的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:用户基本信息的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackSetBasicInfoReq(char * buf, int buf_len, const JWSetBasicInfoReq * pInfo);

/******************************************************************
*函数名:JWUnPackSetBasicInfoRsp
*功能说明：解析设置用户基本信息的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:设置用户基本信息的应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackSetBasicInfoRsp(const char * p, int len, JWSetBasicInfoRsp * pRsp);

/******************************************************************
*函数名:JWPackGetMobileAreaInfoReq
*功能说明：生成获取用户手机帐号所在区域的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:获取用户手机帐号所在区域的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetMobileAreaInfoReq(char * buf, int buf_len, const JWGetMobileAreaInfo * pInfo);


/******************************************************************
*函数名:JWPackGetMobileAreaByNumberReq
*功能说明：生成获取用户手机帐号所在区域的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:获取用户手机帐号所在区域的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetMobileAreaByNumberReq(char * buf, int buf_len, const JWGetMobileAreaByNumber * pInfo);

/******************************************************************
*函数名:JWUnPackGetMobileAreaInfoRsp
*功能说明：解析获取用户手机帐号所在区域信息的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:获取用户手机帐号所在区域的应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetMobileAreaInfoRsp(const char * p, int len, JWGetMobileAreaInfoRsp * pRsp);



/******************************************************************
*函数名:JWPackResetPwdReq
*功能说明：生成重置用户密码的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:重置用户密码的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackResetPwdReq(char * buf, int buf_len, const JWResetPwdInfo * pInfo);

/******************************************************************
*函数名:JWUnPackResetPwdRsp
*功能说明：解析重置用户密码的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:重置用户密码的应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackResetPwdRsp(const char * p, int len, JWResetPwdRsp * pRsp);


/******************************************************************
*函数名:JWPackCheckTicketReq
*功能说明：生成验证登录票据的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:验证登录票据的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackCheckTicketReq(char * buf, int buf_len, const JWCheckTicketInfo * pInfo);

/******************************************************************
*函数名:JWUnPackCheckTicketRsp
*功能说明：解析验证登录票据的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:验证登录票据的应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackCheckTicketRsp(const char * p, int len, JWCheckTicketRsp * pRsp);



/******************************************************************
*函数名:JWPackGetAccountRegTimeReq
*功能说明：生成取得帐号注册时间的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:取得帐号注册时间的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetAccountRegTimeReq(char * buf, int buf_len, const JWGetAccountRegTimeInfo * pInfo);

/******************************************************************
*函数名:JWUnPackGetAccountRegTimeRsp
*功能说明：解析取得帐号注册时间的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:取得帐号注册时间的应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetAccountRegTimeRsp(const char * p, int len, JWGetAccountRegTimeRsp * pRsp);


/******************************************************************
*函数名:JWPackGetUserAccountReq
*功能说明：生成取得用户帐号信息的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:取得用户帐号信息的参数
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetUserAccountReq(char * buf, int buf_len, const JWGetUserAccountInfo * pInfo);

/******************************************************************
*函数名:JWUnPackGetUserAccountInfoRsp
*功能说明：解析取得用户帐号信息的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:取得用户帐号信息的应答参数
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetUserAccountInfoRsp(const char * p, int len, JWGetUserAccountRsp * pRsp);


/******************************************************************
*函数名:JWPackStatReportInfoReq
*功能说明：生成上报用户信息的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:待上报用户信息
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackStatReportInfoReq(char * buf, int buf_len, const JWStatReportInfo * pInfo);

/******************************************************************
*函数名:JWUnPackStatReportInfoRsp
*功能说明：解析上报用户信息的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:上报用户信息的应答包
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackStatReportInfoRsp(const char * p, int len, JWStatReportRsp * pRsp);


/******************************************************************
*函数名:JWPackModifyUserPwdReq
*功能说明：生成修改用户密码的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:修改用户密码的参数信息
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackModifyUserPwdReq(char * buf, int buf_len, const JWModifyPwdInfo * pInfo);

/******************************************************************
*函数名:JWUnPackModifyUserPwdRsp
*功能说明：解析修改用户密码的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:修改用户密码的应答包
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackModifyUserPwdRsp(const char * p, int len, JWModifyUserPwdRsp * pRsp);

/******************************************************************
*函数名:JWPackPayTicketIncomeNotifyReq
*功能说明：生成牛票购买完成通知请求(含收益分成)
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:牛票购买完成通知请求参数信息
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackPayTicketIncomeNotifyReq(char * buf, int buf_len, const JWPayTicketIncomeInfo * pInfo);

/******************************************************************
*函数名:JWUnPackPayTicketIncomeNotifyRsp
*功能说明：解析牛票购买完成通知应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:牛票购买完成应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackPayTicketIncomeNotifyRsp(const char * p, int len, JWPayTicketIncomeRsp * pRsp);

/******************************************************************
*函数名:JWPackSetUserAliasReq
*功能说明：生成设置用户登录别名
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:设置用户登录别名参数信息
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackSetUserAliasReq(char * buf, int buf_len, const JWSetUserAliasInfo * pInfo);

/******************************************************************
*函数名:JWUnPackSetUserAliasRsp
*功能说明：解析设置用户别名的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:设置用户登录别名的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackSetUserAliasRsp(const char * p, int len, JWSetUserAliasInfoRsp * pRsp);

/******************************************************************
*函数名:JWPackGetUserAliasReq
*功能说明：生成获取用户登录别名
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:获取用户登录别名参数信息
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetUserAliasReq(char * buf, int buf_len, const JWGetUserAliasInfo * pInfo);

/******************************************************************
*函数名:JWUnPackGetUserAliasRsp
*功能说明：解析获取用户别名的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:获取用户登录别名的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetUserAliasRsp(const char * p, int len, JWGetUserAliasInfoRsp * pRsp);

/******************************************************************
*函数名:JWPackUpdTicketReq
*功能说明：生成更新用户Ticket的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:更新用户Ticket 的参数信息
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int  JWPackUpdTicketReq(char * buf, int buf_len, const JWUpdTicketInfo * pInfo);

/******************************************************************
*函数名:JWUnPackUpdTicketRsp
*功能说明：解析获取用户Ticket的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:获取用户Ticket的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackUpdTicketRsp(const char * p, int len, JWUpdTicketInfoRsp * pRsp);

/******************************************************************
*函数名:JWPackRegReq
*功能说明：生成用户注册的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:用户注册的参数信息
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackRegReq(char *buf, int len, const JWRegAccountInfo *pInfo);


/******************************************************************
*函数名:JWUnPackRegRsp
*功能说明：解析用户注册的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:获取用户注册的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackRegRsp(const char * p, int len, JWRegAccountRsp * pRsp);


/******************************************************************
*函数名:JWPackActiveReq
*功能说明：生成账号激活的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:账号激活的参数信息
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackActiveReq(char *buf, int len, const JWActiveAccountInfo *pInfo);


/******************************************************************
*函数名:JWUnPackActiveRsp
*功能说明：解析账号激活的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:获取账号激活的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackActiveRsp(const char * p, int len, JWActiveAccountInfoRsp * pRsp);

/******************************************************************
*函数名:JWPackGetTicketReq
*功能说明：生成获取认证信息的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:认证信息参数信息
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackGetTicketReq(char *buf, int len, const JWAuthInfo *pInfo);

/******************************************************************
*函数名:JWUnPackGetTicketRsp
*功能说明：获取认证票据的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:获取认证票据的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetTicketRsp(const char * p, int len, JWAuthInfoRsp * pRsp);

/******************************************************************
*函数名:JWPackGetRegInfoReq
*功能说明：生成获取账户信息的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:账户请求参数信息
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackGetRegInfoReq(char *buf, int len, const JWAccountInfo *pInfo);

/******************************************************************
*函数名:JWUnPackGetRegInfoRsp
*功能说明：获取账户信息的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:获取账户信息的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetRegInfoRsp(const char * p, int len, JWAccountInfoRsp * pRsp);


/******************************************************************
*函数名:JWPackPullMarketDataReq 
*功能说明：生成拉最新行情数据的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:最新行情请求参数信息
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackPullMarketDataReq(char *buf, int len, const JWPullMarketDataInfo *pInfo);

/******************************************************************
*函数名:JWUnPackPullMarketDataRsp
*功能说明：拉最新行情数据的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:拉最新行情数据的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackPullMarketDataRsp(const char * p, int len, JWPullMarketDataRsp * pRsp);


/******************************************************************
*函数名:JWPackRadarIndicatorReq 
*功能说明：生成个股雷达图指标的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:个股雷达图指标请求参数信息
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackRadarIndicatorReq(char *buf, int len, const JWStockRadarIndicatorInfo *pInfo);


/******************************************************************
*函数名:JWUnPackRadarIndicatorRsp
*功能说明：解析个股雷达图指标的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:个股雷达图指标的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackRadarIndicatorRsp(const char * p, int len, JWStockRadarIndicatorRsp *pRsp);


/******************************************************************
*函数名:JWPackEarlyWarningReq 
*功能说明：生成预警行情数据的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:预警行情数据请求参数信息
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackEarlyWarningReq(char *buf, int len, const JWEarlyWarningInfo *pInfo);

/******************************************************************
*函数名:JWUnPackEarlyWarningRsp
*功能说明：解析预警行情数据的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:预警行情数据的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackEarlyWarningRsp(const char * p, int len, JWEarlyWarningRsp *pRsp);


/******************************************************************
*函数名:JWPackEarlyWarningReq 
*功能说明：生成批量拉取股票简易行情数据的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:批量拉取股票简易行情数据请求参数信息
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackBatchBriefMarketDataReq(char *buf, int len, const JWBatchBriefMarketDataInfo *pInfo);

/******************************************************************
*函数名:JWUnPackEarlyWarningRsp
*功能说明：解析批量拉取股票简易行情的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:批量拉取股票简易行情数据的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackBatchBriefMarketDataRsp(const char * p, int len, JWBatchBriefMarketDataInfoRsp *pRsp);


/******************************************************************
*函数名:JWPackStockNewsContentReq 
*功能说明：生成获取单只股票新闻内容数据的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:获取单只股票新闻内容数据请求参数信息
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackStockNewsContentReq(char *buf, int len, const JWStkNewsContentInfo *pInfo);

/******************************************************************
*函数名:JWUnPackStockNewsContentRsp
*功能说明：解析获取单只股票新闻内容数据的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:获取单只股票新闻内容数据的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackStockNewsContentRsp(const char * p, int len, JWStkNewsContentInfoRsp *pRsp);

/******************************************************************
*函数名:JWPackEarlyWarningReq 
*功能说明：生成批量拉取股票新闻数据的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:批量拉取股票新闻数据请求参数信息
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackBatchStockNewsNumInfoReq(char *buf, int len, const JWBatchStockNewsNumInfo *pInfo);

/******************************************************************
*函数名:JWUnPackBatchStockNewsNumInfoRsp
*功能说明：解析获批量拉取股票新闻数据的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:批量拉取股票新闻数据的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackBatchStockNewsNumInfoRsp(const char * p, int len, JWBatchStockNewsNumInfoRsp *pRsp);

/******************************************************************
*函数名:JWPackGetRefereceDataReq 
*功能说明：生成获取所有个股信息的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*返回值：
*   >0:生成的包长度,<=0 失败
*注意：该命令请求消息体为空
*****************************************************************/
JW_PROTOCOL_API int JWPackGetRefereceDataReq(char *buf, int len);

/******************************************************************
*函数名:JWUnPackGetRefereceDataRsp
*功能说明：解析获取所有个股信息的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:获取所有个股信息的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetRefereceDataRsp(const char * p, int len, JWStockGetRefereceDataRsp *pRsp);


/******************************************************************
*函数名:JWPackPullMinDataReq 
*功能说明：生成拉分时线数据的请求包
*参数说明：
*   OUT buf:输出缓冲区
*   IN buf_len:缓冲区的长度
*   IN pInfo:分时数据请求参数信息
*返回值：
*   >0:生成的包长度,<=0 失败
*****************************************************************/
JW_PROTOCOL_API int JWPackPullMinDataReq(char *buf, int len, const JWPullMinuteDataInfo *pInfo);

/******************************************************************
*函数名:JWUnPackPullMinDataRsp
*功能说明：解析拉分时线数据的应答包
*   IN p:完整数据包指针
*   IN len:数据包长度
*   OUT pRsp:拉分时线数据的应答信息结构
*返回值：
*   0:成功,<0失败
*****************************************************************/
JW_PROTOCOL_API int JWUnPackPullMinDataRsp(const char * p, int len, JWPullMinuteDataInfoRsp * pRsp);


#endif
