
#ifndef _JW_PROTOCOL_H_
#define _JW_PROTOCOL_H_

#include "jw_os.h"
#include "jw_struct.h"



/******************************************************************
*������:JWInit
*����˵������ʼ��
*����˵����
*   ��
*����ֵ��
*   ��
*****************************************************************/
JW_PROTOCOL_API void  JWInit();

/******************************************************************
*������:JWUnInit
*����˵������Դ����
*����˵����
*   ��
*����ֵ��
*   ��
*****************************************************************/
JW_PROTOCOL_API void  JWUnInit();


/******************************************************************
*������:JWGetModuleVer
*����˵����ȡ��ģ��汾��
*����˵����
*   ��
*����ֵ��
*   ģ��汾��Ϣ(��2�ֽڱ�ʾ���汾�ţ���2�ֽڱ�ʾ�Ӱ汾��)
*****************************************************************/
JW_PROTOCOL_API unsigned int  JWGetModuleVer();



/******************************************************************
*������:JWHandleInput
*����˵������ȡһ��������Э���
*����˵����
*   IN buf:���յ�������ָ��
*   IN len:����������
*����ֵ��
*   >0:��ʾ���ݰ��ѽ��������Ҹ�ֵ��ʾ���ݰ��ĳ���
*   0����ʾ���ݰ���δ��������
*   <0����ʾ����
*****************************************************************/
JW_PROTOCOL_API int  JWHandleInput(const char * buf, unsigned int len);


/******************************************************************
*������:JWGetPackLen
*����˵����ȡ�����ݰ��ĳ���
*����˵����
*   IN p:�������ݰ�ָ��
*����ֵ��
*   ���ݰ��ĳ���
*****************************************************************/
JW_PROTOCOL_API unsigned int  JWGetPackLen(const char * p);

/******************************************************************
*������:JWGetPackCmd
*����˵����ȡ�����ݰ���������
*����˵����
*   IN p:�������ݰ�ָ��
*   OUT main_cmd:��������
*   OUT sub_cmd:��������
*����ֵ��
*   0:�ɹ�
*   <0����ʾ����
*****************************************************************/
JW_PROTOCOL_API int  JWGetPackCmd(const char * p, unsigned short * main_cmd, unsigned short * sub_cmd);

/******************************************************************
*������:JWGetPackSeq
*����˵����ȡ�����ݰ������
*����˵����
*   IN p:�������ݰ�ָ��
*����ֵ��
*   ���ݰ������
*****************************************************************/
JW_PROTOCOL_API int  JWGetPackSeq(const char * p);

/******************************************************************
*������:JWSetPackSeq
*����˵�����������ݰ������
*����˵����
*   IN p:�������ݰ�ָ��
*   IN seq:���ݰ����
*����ֵ��
*   ��
*****************************************************************/
JW_PROTOCOL_API void  JWSetPackSeq(const char * p, unsigned int seq);


/******************************************************************
*������:JWGetSourceType
*����˵����ȡ�����ݰ�����Դ���ͱ�־
*����˵����
*   IN p:�������ݰ�ָ��
*����ֵ��
*   ���ݰ�����Դ���ͱ�־
*****************************************************************/
JW_PROTOCOL_API unsigned char  JWGetSourceType(const char * p);

/******************************************************************
*������:JWSetSourceType
*����˵�����������ݰ�����Դ���ͱ�־
*����˵����
*   IN p:�������ݰ�ָ��
*   IN source_type:���ݰ���Դ���ͱ�־
*����ֵ��
*   ��
*****************************************************************/
JW_PROTOCOL_API void  JWSetSourceType(const char * p, unsigned char source_type);

/******************************************************************
*������:JWEncryptPack
*����˵�����������ݰ�
*����˵����
*   IN OUT p:�������ݰ�ָ��
*   IN len:���ݰ�����
*����ֵ��
*   ��
*****************************************************************/
JW_PROTOCOL_API void  JWEncryptPack(char * p, unsigned int len);


/******************************************************************
*������:JWDecryptPack
*����˵�����������ݰ�
*����˵����
*   IN OUT p:�������ݰ�ָ��
*   IN len:���ݰ�����
*����ֵ��
*   ��
*****************************************************************/
JW_PROTOCOL_API void  JWDecryptPack(char * p,unsigned int len);



/******************************************************************
*������:JWPackKeepAliveReq
*����˵��������������Ϣ�������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:����Ĳ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackKeepAliveReq(char * buf, int buf_len, const JWCommonReqInfo * pInfo);


/******************************************************************
*������:JWUnPackKeepAliveRsp
*����˵��������������Ϣ��Ӧ���
*����˵����
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackKeepAliveRsp(const char * p, int len, JWCommonRsp * pRsp);


/******************************************************************
*������:JWPackGetAccountStatusReq
*����˵�������ɲ�ѯ�ʺ�״̬�������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:��ѯ�ʺŵĲ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetAccountStatusReq(char * buf, int buf_len, const JWGetAccountStatusInfo * pInfo);

/******************************************************************
*������:JWUnPackGetAccountStatusRsp
*����˵����������ѯ�ʺ�״̬��Ӧ���
*����˵����
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ѯ�ʺ�Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetAccountStatusRsp(const char * p, int len, JWGetAccountStatusRsp * pRsp);


/******************************************************************
*������:JWPackCreateAccountReq
*����˵�������ɿ�ͨ�ʺŵ������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:��ͨ�ʺŵĲ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackCreateAccountReq(char * buf, int buf_len, const JWCreateAccountInfo * pInfo);

/******************************************************************
*������:JWUnPackCreateAccountRsp
*����˵����������ͨ�ʺŵ�Ӧ���
*����˵����
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ͨ�ʺ�Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackCreateAccountRsp(const char * p, int len, JWCreateAccountRsp * pRsp);


/******************************************************************
*������:JWPackLoginReq
*����˵�������ɵ�¼�������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:��½����
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackLoginReq(char * buf, int buf_len, const JWLoginInfo * pInfo);


/******************************************************************
*������:JWUnPackLoginRsp
*����˵����������¼��Ӧ���
*����˵����
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��½Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackLoginRsp(const char * p, int len, JWLoginRsp * pRsp);

/******************************************************************
*������:JWPackLogoutReq
*����˵��������ע���������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:ע������
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackLogoutReq(char * buf, int buf_len, const JWLogoutInfo * pInfo);

/******************************************************************
*������:JWUnPackLogoutRsp
*����˵��������ע����Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:ע��Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackLogoutRsp(const char * p, int len, JWLogoutRsp * pRsp);



/******************************************************************
*������:JWPackAddStockReq
*����˵��������������ѡ�ɵ������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:������ѡ�ɵĲ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackAddStockReq(char * buf, int buf_len, const JWAddStockInfo * pInfo);

/******************************************************************
*������:JWUnPackAddStockRsp
*����˵��������������ѡ�ɵ�Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:������ѡ�ɵ�Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackAddStockRsp(const char * p, int len, JWAddStockRsp * pRsp);


/******************************************************************
*������:JWPackDelStockReq
*����˵��������ɾ����ѡ�ɵ������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:ɾ����ѡ�ɵĲ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackDelStockReq(char * buf, int buf_len, const JWDelStockInfo * pInfo);

/******************************************************************
*������:JWUnPackDelStockRsp
*����˵��������ɾ����ѡ�ɵ�Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:ɾ����ѡ�ɵ�Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackDelStockRsp(const char * p, int len, JWDelStockRsp * pRsp);

/******************************************************************
*������:JWPackGetStockListReq
*����˵����������ȡ��ѡ���б�������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:��ȡ��ѡ���б�Ĳ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetStockListReq(char * buf, int buf_len, const JWGetStockInfo * pInfo);

/******************************************************************
*������:JWUnPackGetStockListRsp
*����˵����������ȡ��ѡ���б��Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ��ѡ���б��Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetStockListRsp(const char * p, int len, JWGetStockInfoRsp * pRsp);

/******************************************************************
*������:JWPackSetStockSortReq
*����˵��������������ѡ������������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:������ѡ������Ĳ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackSetStockSortReq(char * buf, int buf_len, const JWSetStockSortInfo * pInfo);

/******************************************************************
*������:JWUnPackSetStockSortRsp
*����˵��������������ѡ�������Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:������ѡ�������Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackSetStockSortRsp(const char * p, int len, JWSetStockSortRsp * pRsp);

/******************************************************************
*������:JWPackGetStockSortReq
*����˵����������ȡ��ѡ������������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:��ȡ��ѡ������Ĳ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetStockSortReq(char * buf, int buf_len, const JWGetStockSortInfo * pInfo);

/******************************************************************
*������:JWUnPackGetStockSortRsp
*����˵����������ȡ��ѡ�������Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ��ѡ�������Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetStockSortRsp(const char * p, int len, JWGetStockSortRsp * pRsp);

/******************************************************************
*������:JWPackGetStockRankReq
*����˵����������ȡ��ѡ�������������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:��ȡ��ѡ�������Ĳ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetStockRankReq(char * buf, int buf_len, const JWGetStockRankInfo * pInfo);

/******************************************************************
*������:JWUnPackGetStockRankRsp
*����˵����������ȡ��ѡ��������Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ��ѡ��������Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetStockRankRsp(const char * p, int len, JWGetStockRankRsp * pRsp);

#if 0
//2010-11-26 zhongheming ҵ�������������ɽӿ� 
/******************************************************************
*������:JWPackBindMobileReq
*����˵�������ɰ�/������ֻ�����������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:��/������ֻ�����Ĳ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackBindMobileReq(char * buf, int buf_len, const JWBindMobileInfo * pInfo);

/******************************************************************
*������:JWUnPackBindMobileRsp
*����˵����������/������ֻ������Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��/������ֻ������Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackBindMobileRsp(const char * p, int len, JWBindMobileRsp * pRsp);

/******************************************************************
*������:JWPackBindVerifyReq
*����˵��������У�� ��/����� ��֤��������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:��/�������֤��Ĳ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackBindVerifyReq(char * buf, int buf_len, const JWBindVerifyInfo * pInfo);

/******************************************************************
*������:JWUnPackBindVerifyRsp
*����˵��������У���/������ֻ�������֤���Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:У���/������ֻ�������֤���Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackBindVerifyRsp(const char * p, int len, JWBindVerifyRsp * pRsp);



/******************************************************************
*������:JWPackGetBindListReq
*����˵��������ȡ�ð󶨵��ֻ���Ϣ�����
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:ȡ�ð󶨵��ֻ���Ϣ�Ĳ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetBindListReq(char * buf, int buf_len, const JWGetBindInfo * pInfo);

/******************************************************************
*������:JWUnPackGetBindListRsp
*����˵��������ȡ�ð󶨵��ֻ���Ϣ��Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:ȡ�ð󶨵��ֻ���Ϣ��Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetBindListRsp(const char * p, int len, JWGetBindInfoRsp * pRsp);

#endif

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/******************************************************************
*������:JWPackBindMobileInfoReq
*����˵�������ɰ�/������ֻ�����������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:��/������ֻ�����Ĳ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackBindMobileInfoReq(char * buf, int buf_len, const JWBindMobileInfoReq * pInfo);

/******************************************************************
*������:JWUnPackBindMobileInfoRsp
*����˵����������/������ֻ������Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��/������ֻ������Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackBindMobileInfoRsp(const char * p, int len, JWBindMobileInfoRsp * pRsp);

/******************************************************************
*������:JWPackBindVerifyInfoReq
*����˵��������У�� ��/����� ��֤��������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:��/�������֤��Ĳ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackBindVerifyInfoReq(char * buf, int buf_len, const JWBindVerifyInfoReq * pInfo);

/******************************************************************
*������:JWUnPackBindVerifyInfoRsp
*����˵��������У���/������ֻ�������֤���Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:У���/������ֻ�������֤���Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackBindVerifyInfoRsp(const char * p, int len, JWBindVerifyInfoRsp * pRsp);



/******************************************************************
*������:JWPackGetBindMobileInfoReq
*����˵��������ȡ�ð󶨵��ֻ���Ϣ�����
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:ȡ�ð󶨵��ֻ���Ϣ�Ĳ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetBindMobileInfoReq(char * buf, int buf_len, const JWGetBindMobileInfoReq * pInfo);

/******************************************************************
*������:JWUnPackGetBindMobileInfoRsp
*����˵��������ȡ�ð󶨵��ֻ���Ϣ��Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:ȡ�ð󶨵��ֻ���Ϣ��Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetBindMobileInfoRsp(const char * p, int len, JWGetBindMobileInfoRsp * pRsp);


/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/******************************************************************
*������:JWPackGetBasicInfoReq
*����˵��������ȡ���û�������Ϣ�������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:ȡ���û�������Ϣ�Ĳ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetBasicInfoReq(char * buf, int buf_len, const JWGetBasicInfo * pInfo);

/******************************************************************
*������:JWUnPackGetBasicInfoRsp
*����˵��������ȡ���û�������Ϣ��Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:ȡ���û�������Ϣ��Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetBasicInfoRsp(const char * p, int len, JWGetBasicInfoRsp * pRsp);

/******************************************************************
*������:JWPackSetBasicInfoReq
*����˵�������������û�������Ϣ�������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:�û�������Ϣ�Ĳ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackSetBasicInfoReq(char * buf, int buf_len, const JWSetBasicInfoReq * pInfo);

/******************************************************************
*������:JWUnPackSetBasicInfoRsp
*����˵�������������û�������Ϣ��Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:�����û�������Ϣ��Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackSetBasicInfoRsp(const char * p, int len, JWSetBasicInfoRsp * pRsp);

/******************************************************************
*������:JWPackGetMobileAreaInfoReq
*����˵�������ɻ�ȡ�û��ֻ��ʺ���������������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:��ȡ�û��ֻ��ʺ���������Ĳ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetMobileAreaInfoReq(char * buf, int buf_len, const JWGetMobileAreaInfo * pInfo);


/******************************************************************
*������:JWPackGetMobileAreaByNumberReq
*����˵�������ɻ�ȡ�û��ֻ��ʺ���������������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:��ȡ�û��ֻ��ʺ���������Ĳ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetMobileAreaByNumberReq(char * buf, int buf_len, const JWGetMobileAreaByNumber * pInfo);

/******************************************************************
*������:JWUnPackGetMobileAreaInfoRsp
*����˵����������ȡ�û��ֻ��ʺ�����������Ϣ��Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ�û��ֻ��ʺ����������Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetMobileAreaInfoRsp(const char * p, int len, JWGetMobileAreaInfoRsp * pRsp);



/******************************************************************
*������:JWPackResetPwdReq
*����˵�������������û�����������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:�����û�����Ĳ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackResetPwdReq(char * buf, int buf_len, const JWResetPwdInfo * pInfo);

/******************************************************************
*������:JWUnPackResetPwdRsp
*����˵�������������û������Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:�����û������Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackResetPwdRsp(const char * p, int len, JWResetPwdRsp * pRsp);


/******************************************************************
*������:JWPackCheckTicketReq
*����˵����������֤��¼Ʊ�ݵ������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:��֤��¼Ʊ�ݵĲ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackCheckTicketReq(char * buf, int buf_len, const JWCheckTicketInfo * pInfo);

/******************************************************************
*������:JWUnPackCheckTicketRsp
*����˵����������֤��¼Ʊ�ݵ�Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��֤��¼Ʊ�ݵ�Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackCheckTicketRsp(const char * p, int len, JWCheckTicketRsp * pRsp);



/******************************************************************
*������:JWPackGetAccountRegTimeReq
*����˵��������ȡ���ʺ�ע��ʱ��������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:ȡ���ʺ�ע��ʱ��Ĳ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetAccountRegTimeReq(char * buf, int buf_len, const JWGetAccountRegTimeInfo * pInfo);

/******************************************************************
*������:JWUnPackGetAccountRegTimeRsp
*����˵��������ȡ���ʺ�ע��ʱ���Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:ȡ���ʺ�ע��ʱ���Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetAccountRegTimeRsp(const char * p, int len, JWGetAccountRegTimeRsp * pRsp);


/******************************************************************
*������:JWPackGetUserAccountReq
*����˵��������ȡ���û��ʺ���Ϣ�������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:ȡ���û��ʺ���Ϣ�Ĳ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetUserAccountReq(char * buf, int buf_len, const JWGetUserAccountInfo * pInfo);

/******************************************************************
*������:JWUnPackGetUserAccountInfoRsp
*����˵��������ȡ���û��ʺ���Ϣ��Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:ȡ���û��ʺ���Ϣ��Ӧ�����
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetUserAccountInfoRsp(const char * p, int len, JWGetUserAccountRsp * pRsp);


/******************************************************************
*������:JWPackStatReportInfoReq
*����˵���������ϱ��û���Ϣ�������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:���ϱ��û���Ϣ
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackStatReportInfoReq(char * buf, int buf_len, const JWStatReportInfo * pInfo);

/******************************************************************
*������:JWUnPackStatReportInfoRsp
*����˵���������ϱ��û���Ϣ��Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:�ϱ��û���Ϣ��Ӧ���
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackStatReportInfoRsp(const char * p, int len, JWStatReportRsp * pRsp);


/******************************************************************
*������:JWPackModifyUserPwdReq
*����˵���������޸��û�����������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:�޸��û�����Ĳ�����Ϣ
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackModifyUserPwdReq(char * buf, int buf_len, const JWModifyPwdInfo * pInfo);

/******************************************************************
*������:JWUnPackModifyUserPwdRsp
*����˵���������޸��û������Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:�޸��û������Ӧ���
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackModifyUserPwdRsp(const char * p, int len, JWModifyUserPwdRsp * pRsp);

/******************************************************************
*������:JWPackPayTicketIncomeNotifyReq
*����˵��������ţƱ�������֪ͨ����(������ֳ�)
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:ţƱ�������֪ͨ���������Ϣ
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackPayTicketIncomeNotifyReq(char * buf, int buf_len, const JWPayTicketIncomeInfo * pInfo);

/******************************************************************
*������:JWUnPackPayTicketIncomeNotifyRsp
*����˵��������ţƱ�������֪ͨӦ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:ţƱ�������Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackPayTicketIncomeNotifyRsp(const char * p, int len, JWPayTicketIncomeRsp * pRsp);

/******************************************************************
*������:JWPackSetUserAliasReq
*����˵�������������û���¼����
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:�����û���¼����������Ϣ
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackSetUserAliasReq(char * buf, int buf_len, const JWSetUserAliasInfo * pInfo);

/******************************************************************
*������:JWUnPackSetUserAliasRsp
*����˵�������������û�������Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:�����û���¼������Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackSetUserAliasRsp(const char * p, int len, JWSetUserAliasInfoRsp * pRsp);

/******************************************************************
*������:JWPackGetUserAliasReq
*����˵�������ɻ�ȡ�û���¼����
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:��ȡ�û���¼����������Ϣ
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackGetUserAliasReq(char * buf, int buf_len, const JWGetUserAliasInfo * pInfo);

/******************************************************************
*������:JWUnPackGetUserAliasRsp
*����˵����������ȡ�û�������Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ�û���¼������Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetUserAliasRsp(const char * p, int len, JWGetUserAliasInfoRsp * pRsp);

/******************************************************************
*������:JWPackUpdTicketReq
*����˵�������ɸ����û�Ticket�������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:�����û�Ticket �Ĳ�����Ϣ
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWPackUpdTicketReq(char * buf, int buf_len, const JWUpdTicketInfo * pInfo);

/******************************************************************
*������:JWUnPackUpdTicketRsp
*����˵����������ȡ�û�Ticket��Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ�û�Ticket��Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackUpdTicketRsp(const char * p, int len, JWUpdTicketInfoRsp * pRsp);

/******************************************************************
*������:JWPackRegReq
*����˵���������û�ע��������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:�û�ע��Ĳ�����Ϣ
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackRegReq(char *buf, int len, const JWRegAccountInfo *pInfo);


/******************************************************************
*������:JWUnPackRegRsp
*����˵���������û�ע���Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ�û�ע���Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackRegRsp(const char * p, int len, JWRegAccountRsp * pRsp);


/******************************************************************
*������:JWPackActiveReq
*����˵���������˺ż���������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:�˺ż���Ĳ�����Ϣ
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackActiveReq(char *buf, int len, const JWActiveAccountInfo *pInfo);


/******************************************************************
*������:JWUnPackActiveRsp
*����˵���������˺ż����Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ�˺ż����Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackActiveRsp(const char * p, int len, JWActiveAccountInfoRsp * pRsp);

/******************************************************************
*������:JWPackGetTicketReq
*����˵�������ɻ�ȡ��֤��Ϣ�������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:��֤��Ϣ������Ϣ
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackGetTicketReq(char *buf, int len, const JWAuthInfo *pInfo);

/******************************************************************
*������:JWUnPackGetTicketRsp
*����˵������ȡ��֤Ʊ�ݵ�Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ��֤Ʊ�ݵ�Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetTicketRsp(const char * p, int len, JWAuthInfoRsp * pRsp);

/******************************************************************
*������:JWPackGetRegInfoReq
*����˵�������ɻ�ȡ�˻���Ϣ�������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:�˻����������Ϣ
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackGetRegInfoReq(char *buf, int len, const JWAccountInfo *pInfo);

/******************************************************************
*������:JWUnPackGetRegInfoRsp
*����˵������ȡ�˻���Ϣ��Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ�˻���Ϣ��Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int  JWUnPackGetRegInfoRsp(const char * p, int len, JWAccountInfoRsp * pRsp);


/******************************************************************
*������:JWPackPullMarketDataReq 
*����˵���������������������ݵ������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:�����������������Ϣ
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackPullMarketDataReq(char *buf, int len, const JWPullMarketDataInfo *pInfo);

/******************************************************************
*������:JWUnPackPullMarketDataRsp
*����˵�����������������ݵ�Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:�������������ݵ�Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackPullMarketDataRsp(const char * p, int len, JWPullMarketDataRsp * pRsp);


/******************************************************************
*������:JWPackRadarIndicatorReq 
*����˵�������ɸ����״�ͼָ��������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:�����״�ͼָ�����������Ϣ
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackRadarIndicatorReq(char *buf, int len, const JWStockRadarIndicatorInfo *pInfo);


/******************************************************************
*������:JWUnPackRadarIndicatorRsp
*����˵�������������״�ͼָ���Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:�����״�ͼָ���Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackRadarIndicatorRsp(const char * p, int len, JWStockRadarIndicatorRsp *pRsp);


/******************************************************************
*������:JWPackEarlyWarningReq 
*����˵��������Ԥ���������ݵ������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:Ԥ�������������������Ϣ
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackEarlyWarningReq(char *buf, int len, const JWEarlyWarningInfo *pInfo);

/******************************************************************
*������:JWUnPackEarlyWarningRsp
*����˵��������Ԥ���������ݵ�Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:Ԥ���������ݵ�Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackEarlyWarningRsp(const char * p, int len, JWEarlyWarningRsp *pRsp);


/******************************************************************
*������:JWPackEarlyWarningReq 
*����˵��������������ȡ��Ʊ�����������ݵ������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:������ȡ��Ʊ���������������������Ϣ
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackBatchBriefMarketDataReq(char *buf, int len, const JWBatchBriefMarketDataInfo *pInfo);

/******************************************************************
*������:JWUnPackEarlyWarningRsp
*����˵��������������ȡ��Ʊ���������Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:������ȡ��Ʊ�����������ݵ�Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackBatchBriefMarketDataRsp(const char * p, int len, JWBatchBriefMarketDataInfoRsp *pRsp);


/******************************************************************
*������:JWPackStockNewsContentReq 
*����˵�������ɻ�ȡ��ֻ��Ʊ�����������ݵ������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:��ȡ��ֻ��Ʊ���������������������Ϣ
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackStockNewsContentReq(char *buf, int len, const JWStkNewsContentInfo *pInfo);

/******************************************************************
*������:JWUnPackStockNewsContentRsp
*����˵����������ȡ��ֻ��Ʊ�����������ݵ�Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ��ֻ��Ʊ�����������ݵ�Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackStockNewsContentRsp(const char * p, int len, JWStkNewsContentInfoRsp *pRsp);

/******************************************************************
*������:JWPackEarlyWarningReq 
*����˵��������������ȡ��Ʊ�������ݵ������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:������ȡ��Ʊ�����������������Ϣ
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackBatchStockNewsNumInfoReq(char *buf, int len, const JWBatchStockNewsNumInfo *pInfo);

/******************************************************************
*������:JWUnPackBatchStockNewsNumInfoRsp
*����˵����������������ȡ��Ʊ�������ݵ�Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:������ȡ��Ʊ�������ݵ�Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackBatchStockNewsNumInfoRsp(const char * p, int len, JWBatchStockNewsNumInfoRsp *pRsp);

/******************************************************************
*������:JWPackGetRefereceDataReq 
*����˵�������ɻ�ȡ���и�����Ϣ�������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*ע�⣺������������Ϣ��Ϊ��
*****************************************************************/
JW_PROTOCOL_API int JWPackGetRefereceDataReq(char *buf, int len);

/******************************************************************
*������:JWUnPackGetRefereceDataRsp
*����˵����������ȡ���и�����Ϣ��Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:��ȡ���и�����Ϣ��Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackGetRefereceDataRsp(const char * p, int len, JWStockGetRefereceDataRsp *pRsp);


/******************************************************************
*������:JWPackPullMinDataReq 
*����˵������������ʱ�����ݵ������
*����˵����
*   OUT buf:���������
*   IN buf_len:�������ĳ���
*   IN pInfo:��ʱ�������������Ϣ
*����ֵ��
*   >0:���ɵİ�����,<=0 ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWPackPullMinDataReq(char *buf, int len, const JWPullMinuteDataInfo *pInfo);

/******************************************************************
*������:JWUnPackPullMinDataRsp
*����˵������������ʱ�����ݵ�Ӧ���
*   IN p:�������ݰ�ָ��
*   IN len:���ݰ�����
*   OUT pRsp:����ʱ�����ݵ�Ӧ����Ϣ�ṹ
*����ֵ��
*   0:�ɹ�,<0ʧ��
*****************************************************************/
JW_PROTOCOL_API int JWUnPackPullMinDataRsp(const char * p, int len, JWPullMinuteDataInfoRsp * pRsp);


#endif
