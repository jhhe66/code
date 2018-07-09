/********************************************
//�ļ���:jw_define.h
//����:ϵͳ�ṹ����ͷ�ļ�
//����:�Ӻ���
//����ʱ��:2010.07.07
//�޸ļ�¼:

*********************************************/
#ifndef _JW_STRUCT_H_
#define _JW_STRUCT_H_

#ifdef _WIN32
typedef unsigned __int64 uint64_t;
#else
#include <stdint.h>
#endif



#define JW_MAX_ERR_MSG_LEN              255 //������Ϣ����
#define JW_MAX_ACCOUNT_LEN              64  //�û��ʺų���
#define JW_MAX_PWD_LEN                  64  //�û����볤��
#define JW_MAX_TICKET_LEN               32  //��֤Ʊ�ݳ���
#define JW_MAX_DCODE_LEN                16  //��̬�볤��
#define JW_MAX_STOCK_CODE_LEN           10  //��Ʊ���볤��
#define JW_MAX_NAME_LEN                 64  //�û�������
#define JW_MAX_BIRTH_LEN                12  //���ճ���
#define JW_MAX_MOBILE_LEN               20  //�ֻ����볤��
#define JW_MAX_VERIFY_CODE_LEN          20  //��֤�볤��
#define JW_MAX_URL_LEN                  255 //URL����
#define JW_MAX_PROVIANCE_LEN            64 //������Ϣ����
#define JW_MAX_CITY_LEN                 64  //����
#define JW_MAX_AREA_CODE_LEN            6   //���� 

#define JW_MAX_GROUP_STOCK_NUM          50  //ÿ��������������ѡ������
#define JW_MAX_RANK_NUM                 50  //����Ʊ��������
#define JW_MAX_BIND_MOBILE_NUM          6   //�����ֻ���������

#define JW_MAX_CONTENT_LEN  4096 //�ϱ����ݳ���

#define JW_MAX_ORDERID_LEN 32 //�����ų���
#define JW_MAX_INCOME_NUM 32 //�������
#define JW_MAX_PRODUCT_ID_LEN 32 //��Ʒ��ų���

#define JW_MAX_MAIL_LEN 64 //�ʼ�����󳤶�

#define JW_MAX_STOCK_BUY_SELL_NUM 10 //���ί��ί����������


#define JW_MAX_STOCK_NUM 5000 //����������
#define JW_MAX_SECUFID_LEN 31 //��������ID����󳤶�
#define JW_MAX_STOCK_NAME_LEN 32 //��Ʊ������󳤶�
#define JW_NEWSMEDIA_LEN 256 //���ų�����󳤶� 256
#define JW_NEW_CONTENT_LEN (64*1024) //����������󳤶�64K
#define JW_MAX_SECU_NEWS_TITLE_LEN 256 //���ű�����󳤶�256
#define JW_MAX_STOCKCODE_LEN 8  //��Ʊ������󳤶�
#define JW_STOCK_SHORT_NAME_LEN 32  //��Ʊ�����󳤶�
#define JW_MAX_MINDATA_NUM 250  //��������ʱ���ݸ���

#define JW_MAX_ALIAS_NUM 10 //���10������
#define JW_MAX_REGION_NUM 10 //����������
#define JW_MAX_OPER_ANALY 512 //����������󳤶�
#define JW_MAX_INDU_CODE_LEN 30 //��ҵ������󳤶�
#define JW_MAX_INDUSTRY_NUM 300 //�����ҵ����
#define JW_MAX_DDE_HIS_NUM 100 //�����ʷDDE����
#define JW_MAX_RISE_FALL_NUM 150 //�ǵ��������������
#define JW_MAX_DIAGNOSE_CONTENT_LEN 1024 //���������󳤶�
#define JW_MAX_PARAMNAME_LEN 100 //ָ��������򳤶�
#define JW_MAX_KLINE_TYPE 30 //K����̬����
#define JW_MAX_KLINE_NUM 1000 //K����̬������
#define JW_MAX_MULTI_EMPTY_GAME_CAMP 100

#define JW_MAX_RANK_LEN 30 //������󳤶�
#define JW_MAX_INDUSTRY_NAME_LEN 10 //�����ҵ���Ƴ���
#define JW_MAX_GLOSSARY_LEN 50 //����������󳤶�
#define JW_MAX_MAIN_COST_CONTENT_NUM 15 //�����ɱ��������������� 
#define JW_MAX_TARGETS_CONTENT_NUM 30//����ָ���������ݸ���

#pragma pack (push, 1)

/* ����������Ϣ�ṹ */
struct JWCommonReqInfo
{
    unsigned int UIN;//�û�ID
};

/* ����Ӧ����Ϣ�ṹ */
struct JWCommonRsp
{
    unsigned int nErrno;//������
    char szErrMsg[JW_MAX_ERR_MSG_LEN + 1];//������Ϣ
};

/*��ѯ�ʺ�״̬*/
struct JWGetAccountStatusInfo
{
    char szAccount[JW_MAX_ACCOUNT_LEN + 1];
};

struct JWGetAccountStatusRsp:public JWCommonRsp
{
    unsigned char nStatus;//�ʺ�״̬��0:��ע���Ѽ��1����ע��δ���2��δע��
};

/* ��ͨ�ʺ���Ϣ�ṹ */
struct JWCreateAccountInfo
{
    char szAccount[JW_MAX_ACCOUNT_LEN + 1];//�ʺ�
    char szPasswd[JW_MAX_PWD_LEN + 1];//����  
    unsigned int nSourceIP;//��ԴIP�����������ȡ����Զ��IP��
    char szName[JW_MAX_NAME_LEN + 1];//�û��ǳ�

};


struct JWCreateAccountRsp:public JWCommonRsp
{    
    unsigned int UIN;//�û�ID
};

/* ��¼��Ϣ�ṹ */
struct JWLoginInfo
{
    char szAccount[JW_MAX_ACCOUNT_LEN + 1];//�ʺ�
    char szTicket[JW_MAX_TICKET_LEN + 1];//��֤Ʊ��
    char szDCode[JW_MAX_DCODE_LEN + 1];//��̬��
    unsigned int nSourceIP;//��ԴIP�����������ȡ����Զ��IP��
    unsigned int nClientIP;//�ͻ��˱���IP������0
};

struct JWLoginRsp:public JWCommonRsp
{    
    unsigned int UIN;//�û�ID
    char szTicket[JW_MAX_TICKET_LEN + 1];//��֤Ʊ��
    unsigned int nStatus;//�ʺ�״̬
};

/* ע����¼��Ϣ�ṹ */
struct JWLogoutInfo:public JWCommonReqInfo
{    
    unsigned int nSourceIP;//��ԴIP
};

struct JWLogoutRsp:public JWCommonRsp
{    
};

/* ��Ʊ���� + �г����� */
struct JWStockKey
{
    char szStockCode[JW_MAX_STOCK_CODE_LEN + 1];//��Ʊ����
    unsigned int nMarketType;//�г�����
};

struct JWGetStockInfo:public JWCommonReqInfo
{
    unsigned int nGroupID;//��ѡ�ɷ���ID
};

struct JWGetStockInfoRsp:public JWCommonRsp
{    
    unsigned int nGroupID;//��ѡ�ɷ���ID
    unsigned short nStockNum;//����������
    JWStockKey StockList[JW_MAX_GROUP_STOCK_NUM];//��Ʊ�б�
};


struct JWAddStockInfo:public JWCommonReqInfo
{
    JWStockKey StockKey;//��Ʊ����+�г�����
    unsigned int nGroupID;//����ID
};

struct JWAddStockRsp:public JWCommonRsp
{
    
};

struct JWDelStockInfo:public JWCommonReqInfo
{
    JWStockKey StockKey;//��Ʊ����+�г�����
    unsigned int nGroupID;//��Ʊ����ID
};

struct JWDelStockRsp:public JWCommonRsp
{
    
};


struct JWSetStockSortInfo:public JWCommonReqInfo
{
    unsigned int nGroupID;//��ѡ�ɷ���ID
    unsigned short nStockNum;//��ѡ������
    JWStockKey StockList[JW_MAX_GROUP_STOCK_NUM];//��Ʊ�б�

};

struct JWSetStockSortRsp:public JWCommonRsp
{
    
};

struct JWGetStockSortInfo:public JWCommonReqInfo
{
    unsigned int nGroupID;//��ѡ�ɷ���ID
};

struct JWGetStockSortRsp:public JWCommonRsp
{    
    unsigned int nGroupID;//��ѡ�ɷ���ID
    unsigned short nStockNum;//��ѡ������
    JWStockKey StockList[JW_MAX_GROUP_STOCK_NUM];//��Ʊ�б�
};

struct JWGetStockRankInfo:public JWCommonReqInfo
{
    unsigned short Type;//�������ͣ�1�������ѡ������2����ѡ�������
};

struct JWStockRankItem
{
    JWStockKey StockKey;//��ƱKEY
    unsigned int nAttentionCount;//��ע����
    unsigned char nStatus;//����״̬: -1Ϊ�����½���0Ϊ�������䣻1Ϊ����������2Ϊ����
};
struct JWGetStockRankRsp:public JWCommonRsp
{    
    unsigned short Type;//�������ͣ�1�������ѡ������2����ѡ�������
    unsigned short nStockNum;//��ѡ������
    JWStockRankItem RankList[JW_MAX_RANK_NUM];
};

#if 0
//2010-11-26 zhongheming ҵ�������������ɽӿ� 
struct JWBindMobileInfo:public JWCommonReqInfo
{
    char szMobile[JW_MAX_MOBILE_LEN + 1];//�ֻ�����
    unsigned int nFlag;//�󶨱�־
    unsigned char opType;//�������� 0.��, 1.�����
};

struct JWBindMobileRsp:public JWCommonRsp
{    
    unsigned int nExpire;//��֤��ʧЧʱ�䣬��λ��
};

struct JWBindVerifyInfo:public JWCommonReqInfo
{
    char szMobile[JW_MAX_MOBILE_LEN + 1];//�ֻ�����
    char szVerifyCode[JW_MAX_VERIFY_CODE_LEN + 1];//��֤��
    unsigned char opType;//�������� 0.��, 1.�����
};

struct JWBindVerifyRsp:public JWCommonRsp
{    

};

struct JWGetBindInfo:public JWCommonReqInfo
{

};

struct JWBindItem
{
    
    char szMobile[JW_MAX_MOBILE_LEN + 1];//�ֻ�����
    unsigned int nFlag;//�󶨱�־
};

struct JWGetBindInfoRsp:public JWCommonRsp
{    
    unsigned char nNumber;//����
    JWBindItem BindList[JW_MAX_BIND_MOBILE_NUM];

};
#endif

////////////////////////////////////////////////////////////////////////////////////////////////////////////
//�µİ��ֻ��ʺŽӿ����õĽṹ
////////////////////////////////////////////////////////////////////////////////////////////////////////////

struct JWBindMobileInfoReq:public JWCommonReqInfo
{
    char szMobile[JW_MAX_MOBILE_LEN + 1];//�ֻ�����
    unsigned char BindType;//�󶨵�ҵ������,��jw_define.h�� jw_bind_type ����
    unsigned char opType;//�������� 0.��, 1.�����    
};

struct JWBindMobileInfoRsp:public JWCommonRsp
{    
    unsigned int nExpire;//��֤��ʧЧʱ�䣬��λ��
};



struct JWBindVerifyInfoReq:public JWCommonReqInfo
{
    char szMobile[JW_MAX_MOBILE_LEN + 1];//�ֻ�����
    unsigned char BindType;//�󶨵�ҵ������,��jw_define.h�� jw_bind_type ����
    char szVerifyCode[JW_MAX_VERIFY_CODE_LEN + 1];//��֤��
    unsigned char opType;//�������� 0.��, 1.�����
};

struct JWBindVerifyInfoRsp:public JWCommonRsp
{    

};

struct JWGetBindMobileInfoReq:public JWCommonReqInfo
{
    unsigned char BindType;//�󶨵�ҵ������,��jw_define.h�� jw_bind_type ����
};

struct JWBindItemInfo
{

    char szMobile[JW_MAX_MOBILE_LEN + 1];//�ֻ�����
    unsigned char BindType;//�󶨵�ҵ������,��jw_define.h�� jw_bind_type ����
};

struct JWGetBindMobileInfoRsp:public JWCommonRsp
{    
    unsigned char nNumber;//����
    JWBindItemInfo BindList[JW_MAX_BIND_MOBILE_NUM];
};

////////////////////////////////////////////////////////////////////////////////////////////////////////////

struct JWGetBasicInfo:public JWCommonReqInfo
{    
    unsigned char nFaceType;//ͷ������,1:��ͷ��2��Сͷ��
};

struct JWGetBasicInfoRsp:public JWCommonRsp
{
    char szName[JW_MAX_NAME_LEN + 1];//����
    unsigned char nGender;//�Ա�
    unsigned int nBirthday;//����,��19800101��ʽ�����ִ�
    char szFaceURL[JW_MAX_URL_LEN + 1];//ͷ��URL
};

struct JWSetBasicInfoReq:public JWCommonReqInfo
{
    char szName[JW_MAX_NAME_LEN + 1];//����
    unsigned char nGender;//�Ա�
    unsigned int nBirthday;//����,��19800101��ʽ�����ִ�
};

struct JWSetBasicInfoRsp:public JWCommonRsp
{

};

/*ȡ���ֻ����޵���Ϣ*/
struct JWGetMobileAreaInfo:public JWCommonReqInfo
{

};

struct JWGetMobileAreaInfoRsp:public JWCommonRsp
{
    char szProviance[JW_MAX_PROVIANCE_LEN + 1];//�ֻ�����ʡ��
    char szCity[JW_MAX_CITY_LEN + 1];//�ֻ���������
    char szAreaCode[JW_MAX_AREA_CODE_LEN + 1];//�ֻ���������
};


struct JWResetPwdInfo
{
    char szAccount[JW_MAX_ACCOUNT_LEN + 1];//�ʺ�
};

struct JWResetPwdRsp:public JWCommonRsp
{
    char szNewPwd[JW_MAX_PWD_LEN + 1];//�û�������
};

struct JWGetMobileAreaByNumber
{
    char szMobileNumber[JW_MAX_MOBILE_LEN + 1];
};


/* ��֤Ʊ�� */
struct JWCheckTicketInfo:public JWCommonReqInfo
{
    char szTicket[JW_MAX_TICKET_LEN + 1];//��֤Ʊ��
    unsigned int nSourceIP;//��ԴIP������IP��
    unsigned int nClientIP;//�ͻ���IP
};

struct JWCheckTicketRsp:public JWCommonRsp
{

};


/*ȡ���ʺŵ�ע��ʱ��*/
struct JWGetAccountRegTimeInfo
{
    char szAccount[JW_MAX_ACCOUNT_LEN + 1];//�û��ʺ�
};

struct JWGetAccountRegTimeRsp:public JWCommonRsp
{
    unsigned int nRegTime;//ע��ʱ��(UNIX TimeStamp)
    unsigned int UIN;//�û�ID
};

/*ͨ��UINȡ���ʺ���Ϣ*/
struct JWGetUserAccountInfo:public JWCommonReqInfo
{

};

struct JWGetUserAccountRsp:public JWCommonRsp
{
    char szAccount[JW_MAX_ACCOUNT_LEN + 1];//�û��ʺ���Ϣ
};

/*�û������ϱ�*/
struct JWStatReportInfo:public JWCommonReqInfo
{
    char szAccount[JW_MAX_ACCOUNT_LEN + 1];//�û��ʺ���Ϣ
    unsigned char nSrcType;//�ϱ����ݵĲ�Ʒ��Դ
    unsigned int nSrcIP;//�û���ԴIP��ַ
    unsigned int nDataType;//�ϱ���������,�μ�<�û������ϱ���ʽ˵��.docx>
    char szContent[JW_MAX_CONTENT_LEN + 1];//�ϱ�����,��������Ϊ4K
};

struct JWStatReportRsp:public JWCommonRsp
{
};


struct JWModifyPwdInfo
{
    char szAccount[JW_MAX_ACCOUNT_LEN + 1];//�ʺ�
    char szTicket[JW_MAX_TICKET_LEN + 1];//��֤Ʊ��
    char szDCode[JW_MAX_DCODE_LEN + 1];//��̬��
    char szNewPwd[JW_MAX_PWD_LEN + 1];//�û�����(md5 32)
    unsigned int nSourceIP;//��ԴIP�����������ȡ����Զ��IP��
    unsigned int nClientIP;//�ͻ��˱���IP������0
    
};

struct JWModifyUserPwdRsp:public JWCommonRsp
{
};


struct JWIncomeItem
{
    unsigned int dwAmount;//�������
    unsigned int dwSaleUIN;//�����û�UIN
    char szProductID[JW_MAX_PRODUCT_ID_LEN + 1];//��Ʒ���
    unsigned int dwPayAmount;//��Ʒ֧�����
};

struct JWPayTicketIncomeInfo:public JWCommonReqInfo
{
    uint64_t qwSeq;//��ˮ��
    char szOrderID[JW_MAX_ORDERID_LEN +1];//������
    unsigned int dwOrderAmount;//�������
    unsigned int dwOrderResult;//����״̬
    unsigned int dwSubmitTime;//�����ύʱ��
    unsigned int dwIncomeNum;//�������
    JWIncomeItem ItemList[JW_MAX_INCOME_NUM];
};

struct JWPayTicketIncomeRsp:public JWCommonRsp
{
    
};


struct JWSetUserAliasInfo:public JWCommonReqInfo
{
    unsigned char OpType;//�������ͣ�1:���,2:ɾ��
    char szAlias[JW_MAX_ACCOUNT_LEN + 1];//�û���¼����
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
    char szAlias[JW_MAX_ALIAS_NUM][JW_MAX_ACCOUNT_LEN + 1];//�û���¼����
};

struct JWUpdTicketInfo:public JWCommonReqInfo
{
    unsigned char Type;//Ʊ������
};

struct JWUpdTicketInfoRsp:public JWCommonRsp
{
    char szTicket[JW_MAX_TICKET_LEN + 1];//Ʊ��
    unsigned int nGetInterval;//�´���ȡʱ����
};

//�û�ע����Ϣ
struct JWRegAccountInfo 
{
    unsigned char IsActivate;//�Ƿ�ע��ͬʱ����,0:������,1:��������
    char verify_code[JW_MAX_DCODE_LEN+1];//������,utf-8����,\0����
    char account[JW_MAX_ACCOUNT_LEN+1];//�˺�,utf-8����,\0����
    char password[JW_MAX_PWD_LEN+1];// Сд��md5����,32�ֽ�
    char name[JW_MAX_NAME_LEN+1];//�û�������,utf-8����,\0����
    unsigned char gender;//�Ա�0:δ֪,1:��,2:Ů
    char birth[JW_MAX_BIRTH_LEN+1];//����,��ʽΪYYYY/MM/DD,\0����

    unsigned short investLevel;//Ͷ������ 0:�е���Ǯ���뿴��   1����������һ��  2���в��׼� 3���Ƚϸ�ԣ  4:ʲô��û��
    unsigned int attention;//���̷��.
    unsigned int operatorStyle;//���̷�� ,ÿһbit����һ��  0������ 1������ 2������ 
    
    unsigned int city;//��ס��
    unsigned int ip_addr;//�û�ע��ʱ��IP
    unsigned int internalID;//�ڲ�ID,��ע����ԴURL��
    unsigned int advid;//�����ԴID
    unsigned int productType;//Ƕ��ҳ�����ͣ�
};

struct JWRegAccountRsp:public JWCommonRsp
{
    char szEmail[JW_MAX_MAIL_LEN+1];//�ʼ���¼��ַ
};

struct JWAccountInfo
{
    char account[JW_MAX_ACCOUNT_LEN+1];//�˺�,utf-8����,\0����

};

//��ȡ���û�ע����Ϣ
struct JWAccountInfoRsp:public JWCommonRsp
{
    char account[JW_MAX_ACCOUNT_LEN+1];//�˺�,utf-8����,\0����
    char verify_code[JW_MAX_VERIFY_CODE_LEN+1];//������,utf-8����,\0����
    unsigned int reg_time;//ע��ʱ�䣬YYYYMMDD��ʽ��
    char name[JW_MAX_NAME_LEN+1];//�û�������,utf-8����,\0����
    unsigned char gender;//�Ա�0:δ֪,1:��,2:Ů
    char birth[JW_MAX_BIRTH_LEN+1];//���ո�ʽ��YYYYMMDD  ���磺19810306
    unsigned int invest_level;//Ͷ������
    unsigned int operatorStyle;//���̷�� ,ÿһbit����һ��  0������ 1������ 2������ 
    unsigned int city;//��ס��
    unsigned int ip_addr;//ע��IP
};



//�˻�����
struct JWActiveAccountInfo
{
    char account[JW_MAX_ACCOUNT_LEN+1];//�˺�,utf-8����,\0����
    char verify_code[JW_MAX_DCODE_LEN+1];//������,utf-8����,\0����
};

struct JWActiveAccountInfoRsp:public JWCommonRsp
{
};
//
struct JWAuthInfo:public JWCommonReqInfo
{
    unsigned char type; //Ʊ������
};

struct JWAuthInfoRsp:public JWCommonRsp
{
    char ticket[JW_MAX_TICKET_LEN+1];
};

struct JWXXXXRsp:public JWCommonRsp
{    
    unsigned short nStockNum;//����������
    char a[1];//��Ʊ�б�
};

///////////////////////////////////0x0a09/////////////////////////////////
struct JWBuyLevel
{
	unsigned int nBuyPx;           //ί���
	unsigned int nBuyVolume;       //ί����
};

struct JWSellLevel
{
	unsigned int nSellPx;          //ί����
	unsigned int nSellVolume;      //ί����
};


struct JWStockDetail //��������
{
	unsigned short nStockIndex;                   //��Ʊ����
	unsigned int nTxTime;                         //����ʱ��
	unsigned int nPreClosePx;                     //���ռ�
	unsigned int nOpenPx;                         //���̼�
	unsigned int nHighPx;                         //��߼�
	unsigned int nLowPx;                          //��ͼ�
	unsigned int nLastPx;                         //�ּ�
	uint64_t lTotalVolume;                        //�ɽ�����
	uint64_t lTotalAmount;                        //�ɽ��ܽ��
	uint64_t lFloatShare;                         //��ͨ��
	uint64_t lTotalShare;                         //�ܹɱ�
	int nPanCha;                                  //�̲�
	unsigned int nVolumeRatio;                    //����
	int nWeiBi;                                   //ί��
	int nWeiCha;                                  //ί��
	unsigned char cBuyLevel;                      //ί�����
	JWBuyLevel BuyInfo[JW_MAX_STOCK_BUY_SELL_NUM];//ί������
	unsigned char cSellLevel;                     //ί������
	JWSellLevel SellInfo[JW_MAX_STOCK_BUY_SELL_NUM];//ί������
			
};

struct JWPullMarketDataInfo
{
    unsigned short nStockNumber;                   //��������
    unsigned short nStockIndex[JW_MAX_STOCK_NUM];  //��������

};

struct JWPullMarketDataRsp:public JWCommonRsp
{
    unsigned short nStockNumber;                   //����������
    JWStockDetail StockDetail[JW_MAX_STOCK_NUM];   //������������
};

///////////////////////////////////////0x0A0A//////////////////////

struct JWPullMinuteDataInfo:public JWPullMarketDataInfo
{

};

//��ʱ���ݽṹ
struct JWMinData
{
    unsigned int nTxTime;                   //����ʱ��
    unsigned int nHighPx;                   //��߼ۣ������ x1000���ͻ��� /1000
    unsigned int nLowPx;                    //��ͼۣ������ x1000���ͻ��� /1000
    unsigned int nLastPx;                   //�ּۣ������ x1000���ͻ��� /1000

    uint64_t nTotalVolume;                  //�ɽ�����, �����/100000����λ���֣��ͻ��� x100
    uint64_t nTotalAmount;                  //�ɽ��ܽ���λ��Ԫ, ����� x1000���ͻ��� /1000
};

struct JWStockMinData
{
    unsigned short nStockIndex;              //��Ʊ����
    char szStockCode[JW_MAX_STOCKCODE_LEN]; //��Ʊ���루utf-8���룩
    unsigned int nStockType;                  //�г�����
    unsigned int nPreClosePx;                //���ռۣ������ x1000���ͻ��� /1000
    unsigned int nOpenPx;                    //���̼ۣ������ x1000���ͻ��� /1000
    unsigned short nMinDataNum;              //��ʱ���ݸ���
    JWMinData MinData[JW_MAX_MINDATA_NUM];  //��ʱ����
};

struct JWPullMinuteDataInfoRsp:public JWCommonRsp
{
   unsigned short nStockNumber;
   JWStockMinData StockMinData[JW_MAX_STOCK_NUM];

};

///////////////////////////////////////0x0a10///////////////////////
struct JWStockRadarIndicatorInfo
{
    unsigned short nStockIndex;                    //��������
};

struct JWStockRadarIndicatorRsp:public JWCommonRsp
{
    unsigned short nStockIndex;              //��������
    unsigned char cFinanceRisk;              //�������(����� +50 ���>100����ֱ��=100�����<0����ֱ��=0)��Χ��0 ~ 100) (ZHPC|11
    unsigned char cPriceRisk;                //�۸���գ�(��Χ��0 ~ 100) (ZHPC|13)
    unsigned char cAgencyRating;             //����������(��Χ��0 ~ 100) (ZHPC|9)
    unsigned char cCompanyOperation;         //��˾��Ӫ��(��Χ��0 ~ 100) (ZHPC|15)
    unsigned char cMoneyFlow;                //�ʽ�����(��Χ��0 ~ 100) (ZHPC|7)

};


/////////////////////////////0x0a1d///////////////////////////////////
struct JWEarlyWarningInfo
{
    unsigned char cMarketType;      //�г�����,
    unsigned int nVerCounter;       //�汾����������Ʊ���������ݵ�ֵ����Ҫ�ش�����Ʊ�����������ڼ����Ƿ����µ�����
};


struct JWMarketDataInfo  //�����������ݽṹ
{
    unsigned short nStockIndex;         //��Ʊ����
    unsigned int nTxTime;               //����ʱ��
    unsigned int nPreClosePx;           //���ռ�
    unsigned int nLastPx;               //�ּ�
    uint64_t lTotalVolume;              //�ɽ�����
    uint64_t lFloatShare;               //��ͨ��
};

struct JWEarlyWarningRsp:public JWCommonRsp
{
    unsigned int nVerCounter;                         //�汾��������
    unsigned short nStockNumber;                      //��Ʊ����
    JWMarketDataInfo MarketDataInfo[JW_MAX_STOCK_NUM];//�����������ݽṹ
};


/////////////////////////////0x0a22///////////////////////////
struct JWStockBriefMacketData  //��Ʊ�����������ݽṹ
{
    unsigned short nStockIndex;    //��Ʊ����
    unsigned int nTxTime;          //����ʱ��
    unsigned int nPreClosePx;      //���ռۣ������ x1000���ͻ��� /1000
    unsigned int nOpenPx;          //���̼ۣ������ x1000���ͻ��� /1000
    unsigned int nLastPx;          //�ּۣ������ x1000���ͻ��� /1000
    uint64_t nTotalVolume;         //�ɽ���������λ���֣��ͻ��� x100
    uint64_t nTotalAmount;         //�ɽ��ܽ���λ��Ԫ
    uint64_t nFloatShare;          //��ͨ��
    int nPanCha;                   //�̲����<0
    unsigned int nVolumeRatio;      //���ȣ������ /10 (��Ϊ��Ӯ�Ҵ������� �Ѿ�x1000)���ͻ��� /100

};

struct JWBatchBriefMarketDataInfo :public JWPullMarketDataInfo //������Ҫ����������������Ϣ�ṹ����������������������Ϣ�ṹ��һ��
{

};

struct JWBatchBriefMarketDataInfoRsp:public JWCommonRsp  //������ȡ��Ʊ������������Ӧ����Ϣ��ṹ
{
    unsigned short nStockNumber;                                //��Ʊ����
    JWStockBriefMacketData BriefMacketData[JW_MAX_STOCK_NUM];   //��Ʊ�����������ݽṹ

};

//////////////////////////////////0x0a28///////////////////////////////////////////////////
struct JWStkNewsContentInfo
{
    unsigned short nStockIndex;                       //��Ʊ����
    char szSecuFid[JW_MAX_SECUFID_LEN];               //��������ID������ֵΪ��ȡ��ѡ��ʱ��ӷ�������ȡ����SecuFid��
};

struct JWStkNewsContentInfoRsp:public JWCommonRsp
{
    char szStockName[JW_MAX_STOCK_NAME_LEN];               //��Ʊ����
    unsigned int nSecuNewsTime;                            //��������ʱ��(Mine/GP, unixʱ��)
    char szSecuNewsTitle[JW_MAX_SECU_NEWS_TITLE_LEN];      //�������ű���(Mine/GP)
    char szNewsMedia[JW_NEWSMEDIA_LEN];                    //���ų���
    char szNewsContent[JW_NEW_CONTENT_LEN];                //��������

};

//////////////////////////////////0x0a29/������ȡ��ֻ��Ʊ�������ű���ʹ�õĽṹ����////////////////////////////////
struct JWStockNewsNumInfo
{
    unsigned short nStockIndex;    //��Ʊ����
    unsigned short nNewsNum;       // ��ȡ����������������1~20

    JWStockNewsNumInfo& operator=(const JWStockNewsNumInfo &rs)
    {
        this->nStockIndex = rs.nStockIndex;
        this->nNewsNum = rs.nNewsNum;
        return *this;
    }
};

struct JWBatchStockNewsNumInfo  //������ȡ��Ʊ�������ű�������ṹ��
{
    unsigned short nStockNumber;
    JWStockNewsNumInfo StockNews[JW_MAX_STOCK_NUM];

};

struct JWStockSingleNewsData  //���ɵ������Žṹ
{
    unsigned int nSecuNewsTime;
    char szSecuNewsTitle[JW_MAX_SECU_NEWS_TITLE_LEN];
    char szSecuFid[JW_MAX_SECUFID_LEN];
};

struct JWStockAllNewsData
{
    unsigned short nStockIndex;
    unsigned nNewsNum;
    JWStockSingleNewsData NewsData[20]; //��ֻ��Ʊ�����ȡ20������

};

struct JWBatchStockNewsNumInfoRsp:public JWCommonRsp
{
    unsigned short nStockNumber;
    JWStockAllNewsData StockAllNewsData[JW_MAX_STOCK_NUM];

};

//////////////////////////ȡreference data, 0x0a60//////////////////////////
//�������������Ϣ��Ϊ��

struct JWStockInfoData
{
    unsigned short nStockIndex;                           //��Ʊ����
    unsigned int nStockType;                              //��Ʊ����
    char szStockCode[JW_MAX_STOCKCODE_LEN];               //��Ʊ����
    char szStockName[JW_MAX_STOCK_NAME_LEN];              //��Ʊ����
    char szStockShortName[JW_STOCK_SHORT_NAME_LEN];       //��Ʊ���
    unsigned int nPreClosePx;                             //���ռ�
    uint64_t nFloatShare;                                 //��ͨ��
    uint64_t nTotalShare;                                 //�ܹɱ�
};

struct JWStockGetRefereceDataRsp : public JWCommonRsp
{
    unsigned int nRefVersion; //codelist�汾��
    unsigned int nStockNumber;
    JWStockInfoData StockInfoData[JW_MAX_STOCK_NUM];

};

//////////////////////////���ɳɱ��ֲ�ָ��, 0x0a17//////////////////////////
struct  JWRegionData
{
       unsigned char RegionType; //�������ͣ�(CBFB)
       unsigned int RegionBegin;  //������㣨�ͻ���  /1000��,��λ��Ԫ(CBFB)
       unsigned int RegionEnd;  //�����յ㣨�ͻ���  /1000��,��λ��Ԫ(CBFB)
       unsigned int HangupNumber;//�����������ͻ���  /1000����,��λ����(CBFB)
       unsigned int HangupAmount; //���ν��ͻ���  /1000������λ����Ԫ(CBFB)
};

struct JWCostDistributionRsp : public JWCommonRsp
{
    unsigned short StockIndex; //��Ʊ����
    unsigned char RegionNumber;//�������
    JWRegionData Region[JW_MAX_REGION_NUM]; //�ֲ�ָ������
    char OperatorAnalysis[JW_MAX_OPER_ANALY];//������������UTF8��

};



//////////////////////////���ɳɱ��ֲ���ϸ, 0x0a18//////////////////////////
struct  JWRegionData_0A18
{
    unsigned char RegionType; //�������ͣ�
    unsigned int RegionBegin; //�������
    unsigned int RegionEnd;   //�����յ�
};

struct JWRegionDetailData
{
    unsigned int DealPx;  //�ɽ���
    unsigned int ChipsNumber; //������
};

struct JWCostDistributionDetailRsp : public JWCommonRsp
{
    unsigned short StockIndex; //��Ʊ����
    unsigned char RegionNumber;//�������
    JWRegionData_0A18 Region[JW_MAX_REGION_NUM]; //������������
    unsigned short DetailNumber; //��ϸ����
    JWRegionDetailData DetailData[1024]; //��ϸ����

};

//////////////////////////�������ʽ���������, 0x0a50//////////////////////////
struct JWStkFundFlowInfo
{
    unsigned short DayCount; //��������(0=���ա�1=3�ա�2=5�ա�3=10�ա�4=20��)
    unsigned short StockNumber; //��Ʊ����
    unsigned short StockIndex[JW_MAX_STOCK_NUM]; //��Ʊ��������
};

struct JWStkFundFlowData
{
    unsigned short StockIndex; //��Ʊ����
    unsigned int TimeStamp; //ʱ���
    unsigned short DayCount; //��������(0=���ա�1=3�ա�2=5�ա�3=10�ա�4=20��)
    char IndustryCode[JW_MAX_INDU_CODE_LEN];//��ҵ����
    unsigned int NowPrice; //�ּۣ������ x1000���ͻ��� /1000 ����Ӯ�Ҵ������ĵ����ݾ��Ѿ�x1000�� ��Ʊ������͸�����ͻ�����Ҫ/1000��
    unsigned int DeltaPercent; //�ǵ���
    unsigned int ChangeRate;//������
    uint64_t TotalFlowIn; //������
    uint64_t TotalFlowOut;//������
    uint64_t NetFlowIn;//������
    uint64_t NetFlowOut;//������
    unsigned int NetFlowInPower;//���������� 
    unsigned int NetFlowOutPower;//���������� 
    uint64_t NetBigBill;//�󵥾���
    unsigned int ImpetusBill;//�󵥶��� 
    uint64_t MainNetFlowIn; //����������
    uint64_t MainNetFlowOut;//����������
    unsigned int MainNetInRate;//��������������
    unsigned int MainNetOutRate;//��������������
    unsigned int SeriesAddDays;//������������
    uint64_t SeriesNetIn; //�������ֻ��������������/1000
    uint64_t SeriesNetOut;//�������ֻ���������������/1000
    uint64_t SeriesNetInPower;//�������ֻ�������������
    uint64_t SeriesNetOutPower;//�������ֻ�������������
    unsigned int AreaAmountRate;//ռ�����ܳɽ�����
    unsigned int AreaChangeRate;//���任���� 
    unsigned int AreaDeltaPercent;//�����ǵ���
    unsigned int AmountInRate;//������ռ�ɽ����
    unsigned int AmountOutRate;//������ռ�ɽ����
    unsigned int AreaClose; //�������̼�
    unsigned int nDayClose;//3�ջ���5�ջ����������ڵ����̼�
    unsigned int AreaDate;
    unsigned int nDayDate;
};

struct JWStkFundFlowDataRsp:public JWCommonRsp
{
    unsigned short StockNumber; //��Ʊ����
    JWStkFundFlowData StkFundFlowData[JW_MAX_STOCK_NUM];//�ʽ��������ݽṹ����
};

//////////////////////////����ҵ�ʽ���������, 0x0a53//////////////////////////
struct JWIndustryFundFlowInfo
{
    unsigned short DayCount; //��������(0=���ա�1=3�ա�2=5�ա�3=10�ա�4=20��)
    unsigned short IndustryType;//��ҵ���ͣ�0=֤�����ҵ; 1=�²Ƹ���ҵ��
    unsigned short IndustryNumber;//��ҵ����
    unsigned short IndustryIndex[JW_MAX_INDUSTRY_NUM];//��ҵ��������
};

struct JWIndustryFundFlowData
{
    unsigned short IndustryIndex; //��ҵ����
    unsigned int TimeStamp; //ʱ���
    unsigned int DeltaPercent;//�ǵ���
    unsigned int DownCount; //�µ�����
    unsigned int UpCount; //���Ǽ���
    unsigned int ChangeRate; //��ҵ������
    uint64_t TotalFlowIn;//��ҵ�ʽ�������
    uint64_t TotalFlowOut;//��ҵ�ʽ�������
    uint64_t NetFlowIn;//��ҵ�ʽ�����
    unsigned int NetFlowInPower;//��ҵ�ʽ���������
    unsigned int LeaderStock;//��ҵ�����ǹ�
    unsigned int LeaderZDF;//���ǹ��ǵ���
};

struct JWIndustryFundFlowDataRsp:public JWCommonRsp
{
    unsigned short DayCount;//��������(0=���ա�1=3�ա�2=5�ա�3=10�ա�4=20��)
    unsigned short IndustryType;//��ҵ���ͣ�0=֤�����ҵ; 1=�²Ƹ���ҵ��
    unsigned short IndustryNumber;//��ҵ����
    JWIndustryFundFlowData IndustryFundFlowData[JW_MAX_INDUSTRY_NUM];//��ҵ�ʽ��������ݽṹ����
};

//////////////////////////��DDE����, 0x0a61//////////////////////////
struct JWDdeDataInfo
{
    unsigned short StockNumber;//��Ʊ����
    unsigned short StockIndex[JW_MAX_STOCK_NUM];//��Ʊ��������
};

struct JWDDEData
{
    unsigned short StockIndex; //��Ʊ����
    unsigned int Date; //���ڣ���ʽ��Unix��ʽʱ����ֵ��
    uint64_t MainHold;//�����ֲ�
    uint64_t DisperseHold; //ɢ���ֲ�
    uint64_t BuyBigOrder; //����� ��/���
    uint64_t SelBigOrder; //������ ��/���
    uint64_t BuyMinOrder; //����С�� ��/���
    uint64_t SelMinOrder; //����С�� ��/���
    uint64_t BuyMidOrder;//�����е� ��/���
    uint64_t SelMidOrder;//�����е� ��/���
    uint64_t BuyLargeOrder;//�����ش� ��/���
    uint64_t SelLargeOrder; //�����ش� ��/���
    uint64_t BuyCount; //����ɽ�����
    uint64_t SellCount; //�����ɽ�����
};

struct JWDDEDataRsp:public JWCommonRsp
{
    unsigned short DDEDataNums; //DDE���ݸ���
    JWDDEData DDEData[JW_MAX_STOCK_NUM];//DDE���ݽṹ����
};

//////////////////////////��DDSort����, 0x0a62//////////////////////////
struct JWDdsortDataInfo:public JWDdeDataInfo
{
   
};

struct JWDdsortData
{
    unsigned short StockIndex; //��Ʊ����
    unsigned int Time;//����ʱ��
    unsigned int Price;//�ּ�  * 1000 ȡ��
    unsigned int DeltaPercent;//�ǵ��� * 10000 ȡ��, ��ԭ����������
    unsigned int DDX;//DDX * 10000 ȡ��, ��ԭ����������
    unsigned int DDY;//DDY * 10000 ȡ��, ��ԭ����������
    unsigned int DDZ;//DDZ * 10000 ȡ��, ��ԭ����������
    unsigned int DDX60;//60��DDX * 10000 ȡ��, ��ԭ����������
    unsigned int DDY60;//60��DDY * 10000 ȡ��, ��ԭ����������
    unsigned int Up10;//10����Ʈ��
    unsigned int UpDays;//����Ʈ��
    unsigned int LagerBuy;//�ش���* 10000 ȡ��, ��ԭ����������
    unsigned int BigBuy;//����* 10000 ȡ��, ��ԭ����������
    unsigned int LagerSell;//�ش���* 10000 ȡ��, ��ԭ����������
    unsigned int BigSell;//����* 10000 ȡ��, ��ԭ����������
};

struct JWDDSortDataRsp:public JWCommonRsp
{
    unsigned short StockNumber;//��Ʊ����
    JWDdsortData DdsortData[JW_MAX_STOCK_NUM];//����DDSort���ݽṹ����
};

////////////��DDE��ʷ���ݣ�?Business �������� 0x0A64/////////////////
struct JWDdeHisDataInfo
{
    unsigned short StockIndex;//��Ʊ����
    unsigned short Number; //�������������ش�������ǰ�����ݡ�һ����������ṩ60�����ݡ������1����ô�᷵�����µ�һ������
};

struct JWDDEHisData
{
    unsigned short StopBit; //ֹͣλ��0��ʾֹͣ��>0��ʾ��������
    unsigned int Date; //���ڣ���ʽ��Unix��ʽʱ����ֵ��
    uint64_t BuyBigOrder; //����� ��/������� x1000���ͻ��� /1000
    uint64_t SelBigOrder; //������ ��/������� x1000���ͻ��� /1000
    uint64_t BuyMinOrder; //����С�� ��/������� x1000���ͻ��� /1000
    uint64_t SelMinOrder; //����С�� ��/������� x1000���ͻ��� /1000
    uint64_t BuyMidOrder; //�����е� ��/������� x1000���ͻ��� /1000
    uint64_t SelMidOrder; //�����е� ��/������� x1000���ͻ��� /1000
    uint64_t BuyLargeOrder; //�����ش� ��/������� x1000���ͻ��� /1000
    uint64_t SelLargeOrder; //�����ش� ��/������� x1000���ͻ��� /1000
    uint64_t BuyCount; //����ɽ������������ x1000���ͻ��� /1000
    uint64_t SellCount; //�����ɽ������������ x1000���ͻ��� /1000

};

struct JWDDEHisDataRsp:public JWCommonRsp
{
    unsigned short StockIndex; //��Ʊ����
    unsigned int Number; //��������
    JWDDEHisData DDEHisData[JW_MAX_DDE_HIS_NUM];//��ʷDDE��������
};

////////////��ʵʱ�����ǵ������ݣ�?Business �������� 0x0A65/////////////////
struct JWRiseFallBreadthDataInfo
{
    unsigned short Type;//��Ʊ����
    unsigned short Number; //�������������ش�������ǰ�����ݡ�һ����������ṩ60�����ݡ������1����ô�᷵�����µ�һ������
};

struct JWRiseFallBreadthData
{
    unsigned short StockIndex;//��Ʊ����
    unsigned int LastPx;//�ּۣ������ x1000���ͻ��� /1000
    unsigned int PreClosePx;//���ռۣ������ x1000���ͻ��� /1000
    unsigned int BreadthValue;//����ֵ�� ����ֵ=((�ּ�-����)/����)*100000 ���������� x100000 �ͻ���/100000

};

struct JWRiseFallBreadthDataRsp:public JWCommonRsp
{
    unsigned short Type; //���ͣ�1Ϊ�Ƿ���2Ϊ����
    unsigned int Number; //�������������100��
    JWRiseFallBreadthData RiseFallBreadthData[JW_MAX_RISE_FALL_NUM];//�����ǵ������ݽṹ����
};


/////////////��ȡ�ۺ������Ϣ 0x0A40/////////////////
struct JWDiagnoseDataRsp:public JWCommonRsp
{
    unsigned short StockIndex; //��Ʊ����
    char Content[JW_MAX_DIAGNOSE_CONTENT_LEN]; //������ݣ�utf-8���룩
};


//////////�����ʽ�������ϸ����Business ������0x0A66////////

struct JWStkFundFlowDataDetailReq:public JWDdeDataInfo
{

};

struct JWStkFundFlowDetailData
{
    unsigned short StockIndex;//��Ʊ����
    unsigned int Date;//���ڣ���ʽ��Unix��ʽʱ����ֵ��
    uint64_t BuyBigOrderVolume;//�������������� x1000���ͻ��� /1000
    uint64_t BuyBigOrderAmount;//����󵥶����� x1000���ͻ��� /1000
    uint64_t SelBigOrderVolume;//��������������� x1000���ͻ��� /1000
    uint64_t SelBigOrderAmount;//�����󵥶����� x1000���ͻ��� /1000
    uint64_t BuyMinOrderVolume;//����С����������� x1000���ͻ��� /1000
    uint64_t BuyMinOrderAmount;//����С�������� x1000���ͻ��� /1000
    uint64_t SelMinOrderVolume;//����С����������� x1000���ͻ��� /1000
    uint64_t SelMinOrderAmount;//����С�������� x1000���ͻ��� /1000
    uint64_t BuyMidOrderVolume;//�����е���������� x1000���ͻ��� /1000
    uint64_t BuyMidOrderAmount;//�����е������� x1000���ͻ��� /1000
    uint64_t SelMidOrderVolume;//�����е���������� x1000���ͻ��� /1000
    uint64_t SelMidOrderAmount;//�����е������� x1000���ͻ��� /1000
    uint64_t BuyLargeOrderVolume;//�����ش���������� x1000���ͻ��� /1000
    uint64_t BuyLargeOrderAmount;//�����ش󵥶����� x1000���ͻ��� /1000
    uint64_t SelLargeOrderVolume;//�����ش���������� x1000���ͻ��� /1000
    uint64_t SelLargeOrderAmount;//�����ش󵥶����� x1000���ͻ��� /1000
    uint64_t DealCountVolume;//�ɽ�������������� x1000���ͻ��� /1000
    uint64_t DealCountAmount;//�ɽ������� x1000���ͻ��� /1000
};

struct JWStkFundFlowDetailDataRsp:public JWCommonRsp
{
    unsigned short StockNumber;//��Ʊ����
    JWStkFundFlowDetailData StkFundFlowData[JW_MAX_STOCK_NUM];
};

/////////////////K����̬�������� 0x0A67///////////////////
struct JWKlineFormData
{
    unsigned short StockIndex;//��Ʊ����
    char KLineType[JW_MAX_KLINE_TYPE];//K����̬����
    unsigned int CalculateDate; //��������
    unsigned int AppearDate; //�������ڣ���ʽ��Unix��ʽʱ����ֵ��
    unsigned int FormTurnover;//��̬�����ʣ�������*1000���ͻ���/1000��ԭ
};

struct JWKlineFormDataRsp:public JWCommonRsp
{
    unsigned short Type;//���ͣ�0�����¿�����̬��1�����ڿ�����Ч��̬��2�����¿�����̬��3�����ڿ�����Ч��̬
    unsigned short Number;//���������һ����200���������800����
    JWKlineFormData KlineFormData[JW_MAX_KLINE_NUM];//

};

struct JWKlineFormDataReq
{
    unsigned short Type;//���ͣ�0�����¿�����̬��1�����ڿ�����Ч��̬��2�����¿�����̬��3�����ڿ�����Ч��̬
    unsigned short Number;//���������һ����200���������800��
};

////////////////��ղ�����Ӫ���� 0x0A68///////////////
struct JWMultiEmptyGameData
{
    char ParamName[JW_MAX_PARAMNAME_LEN];//ָ������
    unsigned int Yield;//�����ʣ������ x1000���ͻ��� /1000
    unsigned int SuccessRate;//�ɹ��ʣ������ x1000���ͻ��� /1000
    unsigned short Type;//���ͣ� 1�����/2������
    unsigned int Date;//���ڣ���ʽ��Unix��ʽʱ����ֵ��
    unsigned int PreClosePx;//���̼ۣ������ x1000���ͻ��� /1000
    unsigned int OpenPx;//���̼ۣ������ x1000���ͻ��� /1000
};

struct JWMultiEmptyGameDataRsp:public JWCommonRsp
{
    unsigned short StockIndex;//��Ʊ����
    unsigned int StartCalculateDate;//��ʼ�������ڣ���ʽ��Unix��ʽʱ����ֵ��
    unsigned short CalculateDateNum;//��������
    unsigned short Num;//���ݵĸ���
    JWMultiEmptyGameData MultiEmptyGameData[JW_MAX_MULTI_EMPTY_GAME_CAMP];//���ݵĸ���
};

//////////////////0x0A69 ������-���Ʒ�������(Business ������)0x0A69///////////////////////////////////
struct JWTechniccalTrendAnaysisReq
{
    unsigned short StockNumber; //��Ʊ����
    unsigned short StockIndex[JW_MAX_STOCK_NUM];//��Ʊ��������
};

struct JWTechniccalTrendAnaysisData
{
    unsigned short StockIndex;          //��Ʊ����
    unsigned int UpdateDate;            //���ڣ���ʽ��Unix��ʽʱ����ֵ��
    unsigned int ShortSupportPrice;     //����֧��λ������� x1000���ͻ��� /1000
    unsigned int ShortResistancePrice;  //��������λ������� x1000���ͻ��� /1000
    unsigned int MidSupportPrice;       //����֧��λ������� x1000���ͻ��� /1000
    unsigned int MidResistancePrice;    //��������λ������� x1000���ͻ��� /1000
    char Glossary[JW_MAX_GLOSSARY_LEN]; //��������(MACD����ָ��)��utf-8���룩
};

struct JWTechniccalTrendAnaysisRsp:public JWCommonRsp
{
    unsigned short StockNumber;//��Ʊ����
    JWTechniccalTrendAnaysisData TechniccalTrendAnaysisData[JW_MAX_STOCK_NUM];
};

////////////////////// 0x0A70 �ʽ���-��������-�����ɱ�����(Business ������)0x0A70//////////////////////////
struct JWFinancialMainCostReq:public JWTechniccalTrendAnaysisReq
{

};

struct JWFinancialMainCostContentData
{
    unsigned int UpdateDate; //���ڣ���ʽ��Unix��ʽʱ����ֵ��
    unsigned int ClosePrice; //���̼ۣ������ x1000���ͻ��� /1000
    unsigned int MarketPrice;//10���г�ƽ�����׳ɱ�������� x1000���ͻ��� /1000
    unsigned int MainPrice;//10������ƽ�����׳ɱ�������� x1000���ͻ��� /1000

};

struct JWFinancialMainCostData
{   
    unsigned short StockIndex; //��Ʊ����
    unsigned short Num; //���ݵĸ���
    JWFinancialMainCostContentData FinancialMainCostContentData[JW_MAX_MAIN_COST_CONTENT_NUM];//���ݵ�����

};

struct JWFinancialMainCostRsp:public JWCommonRsp
{
    unsigned short StockNumber; //��Ʊ����
    JWFinancialMainCostData FinancialMainCostData[JW_MAX_STOCK_NUM];//������������-�����ɱ���������
};


/////////////////////////////0x0a71 ��Ʊ WEB �C������-�г���������(Business ������)0x0A71

struct JWTechniccalAndMarketReq:public JWTechniccalTrendAnaysisReq
{

};

struct JWTechniccalAndMarketData
{
    unsigned short StockIndex;              //��Ʊ����
    unsigned int UpdateDate;                //���ڣ���ʽ��Unix��ʽʱ����ֵ��
    unsigned int AccumulativeMarkup;        //20���ۼ��Ƿ�������� x100000���ͻ��� /100000
    unsigned int AccumulativeMarkupHY;      //20����ҵ�ۼ��Ƿ�������� x100000���ͻ��� /100000
    unsigned int AccumulativeMarkupSHZS;    //20����ָ֤���ۼ��Ƿ�������� x100000���ͻ��� /100000
    unsigned int TurnOverRate1;             //1�ջ����ʣ������ x100000���ͻ��� /100000
    unsigned int TurnOverRate3;             //3�ջ����ʣ������ x100000���ͻ��� /100000
    unsigned int TurnOverRate5;             //5�ջ����ʣ������ x100000���ͻ��� /100000
};

struct JWTechniccalAndMarketRsp:public JWCommonRsp
{
    unsigned short StockNumber;
    JWTechniccalAndMarketData TechniccalAndMarketData[JW_MAX_STOCK_NUM];

};


//////////////////0x0a72 ������-�����������(Business ������)0x0A72//////////

struct JWFundAmentalsAnalysisReq:public JWTechniccalTrendAnaysisReq
{

};

struct JWFundAmentalsAnalysisData
{
    unsigned short StockIndex;       //��Ʊ����
    unsigned int UpdateDate;         //���ڣ���ʽ��Unix��ʽʱ����ֵ��
    uint64_t YYSR;                   //Ӫҵ���룬����� x1000���ͻ��� /1000
    unsigned int YYSR_TB_ZZ;         //Ӫҵ����ͬ������������� x1000���ͻ��� /1000
    uint64_t JLR_PARENT;             //�����󣬷���� x1000���ͻ��� /1000
    unsigned int JLR_TB_ZZ;          //������ͬ������������� x1000���ͻ��� /1000
    unsigned int MGSY;               //ÿ�����棬����� x1000���ͻ��� /1000
};

struct JWFundAmentalsAnalysisRsp:public JWCommonRsp
{
    unsigned short StockNumber;//��Ʊ����
    JWFundAmentalsAnalysisData FundAmentalsAnalysisData[JW_MAX_STOCK_NUM];//���������������

};

/////////////////////0x0A73 ������-�������-����ָ������(Business ������)0x0A73
struct JWFundAmentalsTargetsReq:public JWTechniccalTrendAnaysisReq
{

};

struct JWFundAmentalsTargetsContentData
{
    unsigned int UpdateDate;        //���ڣ���ʽ��Unix��ʽʱ����ֵ��
    unsigned int YYSR_TB_ZZL;       //Ӫҵ����ͬ�������ʣ������ x100000���ͻ��� /100000
    unsigned int JLR_TB_ZZL;        //������ͬ�������ʣ������ x100000���ͻ��� /100000


};

struct JWFundAmentalsTargetsData
{
    unsigned short StockIndex;                                          //��Ʊ����
    unsigned short Num;                                                 //���ݵĸ���
    JWFundAmentalsTargetsContentData FundAmentalsTargetsContentData[JW_MAX_TARGETS_CONTENT_NUM];  //��������
};

struct JWFundAmentalsTargetsRsp:public JWCommonRsp
{
    unsigned short StockNumber;
    JWFundAmentalsTargetsData FundAmentalsTargetsData[JW_MAX_STOCK_NUM];

};

/////////////////////////  0x0A74 ������-�������-��ҵ��������(Business ������) 0x0A74

struct JWFundIndustryRanksReq:public JWTechniccalTrendAnaysisReq
{

};

struct JWFundIndustryRanksData
{
    unsigned short StockIndex;          //��Ʊ����
    unsigned int UpdateDate;            //���ڣ���ʽ��Unix��ʽʱ����ֵ��
    char IndustryName[JW_MAX_INDUSTRY_NAME_LEN];              //������ҵ����
    uint64_t ZZC;                       //�������ʲ�������� x1000���ͻ��� /1000
    uint64_t ZZC_HY;                    //��ҵ���ʲ���ֵ������� x1000���ͻ��� /1000
    unsigned short ZZC_PM;              //�������ʲ�����
    unsigned int JLRL;                  //���ɾ������ʣ������ x100000���ͻ��� /100000
    unsigned int JLRL_HY;               //��ҵ�������ʣ������ x100000���ͻ��� /100000
    unsigned short JLRL_PM;             //���ɾ�����������
    unsigned int JZCSYL;                //���ɾ��ʲ������ʣ������ x100000���ͻ��� /100000
    unsigned int JZCSYL_HY;             //��ҵ���ʲ������ʣ������ x100000���ͻ��� /100000
    unsigned short JZCSYL_PM;           //���ɾ��ʲ�����������
    unsigned int XSMLL;                 //��������ë���ʣ������ x100000���ͻ��� /100000
    unsigned int XSMLL_HY;              //��ҵ����ë���ʣ������ x100000���ͻ��� /100000
    unsigned short XSMLL_PM;            //��������ë��������
    unsigned short InduElementNum;      //������ҵ�ɷݹ�����
};

struct JWFundIndustryRanksRsp:public JWCommonRsp
{
    unsigned short StockNumber;//��Ʊ����
    JWFundIndustryRanksData FundIndustryRanksData[JW_MAX_STOCK_NUM];//��ҵ������������
};

///////////////////////////////0x0A75 ������-��ֵ�о�����(Business ������)0x0A75////////////
struct JWFundAluationReserchReq:public JWTechniccalTrendAnaysisReq
{

};

struct JWFundAluationReserchData
{
    unsigned short StockIndex;      //��Ʊ����
    unsigned int UpdateDate;        //���ڣ���ʽ��Unix��ʽʱ����ֵ��
    unsigned int PE_TTM;            //��ӯ�ʣ������ x100000���ͻ��� /100000
    unsigned int PE_TTM_HY;         //��ҵ��ӯ�ʣ������ x100000���ͻ��� /100000
    unsigned int PE_TTM_MARKET;     //�г���ӯ�ʣ������ x100000���ͻ��� /100000
    char PE_TTM_RANKING[JW_MAX_RANK_LEN];          //��ӯ������
    unsigned int PB;                //�о��ʣ������ x100000���ͻ��� /100000
    unsigned int PB_HY;             //��ҵ�о��ʣ������ x100000���ͻ��� /100000
    unsigned int PB_MARKET;         //�г��о��ʣ������ x100000���ͻ��� /100000
    char PB_RANKING[JW_MAX_RANK_LEN];              //�о�������
    unsigned int PEG;               //PEG������� x100000���ͻ��� /100000
    unsigned int PEG_HY;            //��ҵPEG������� x100000���ͻ��� /100000
    unsigned int PEG_MARKET;        //�г�PEG������� x100000���ͻ��� /100000
    char PEG_RANKING[JW_MAX_RANK_LEN];             //PEG����
};

struct JWFundAluationReserchRsp:public JWCommonRsp
{
    unsigned short StockNumber;                         //��Ʊ����
    JWFundAluationReserchData FundAluationRserchData[JW_MAX_STOCK_NUM];  //��ֵ�о���������
};

////////////////////////////////0x0A76 �����۵�����(Business ������)0x0A76/////////
struct JWFundOrganizationViewReq:public JWTechniccalTrendAnaysisReq
{

};

struct JWFundOrganizationViewData
{
    unsigned short StockIndex;                     //��Ʊ����
    unsigned int UpdateDate;                       //���ڣ���ʽ��Unix��ʽʱ����ֵ��
    unsigned short OrganNum;                       //��������
    unsigned short StatDays;                       //ͳ������
    unsigned int Rating;                           //��������������� x100000���ͻ��� /100000
    unsigned int Attention;                        //������ע�ȣ������ x100000���ͻ��� /100000
    unsigned short BuyNum;                         //��������
    unsigned short HoldingsNum;                    //��������
    unsigned short NeutralNum;                     //��������
    unsigned short ReductionNum;                   //��������
    unsigned short SellNum;                        //��������
    char PE_TTM_RANKING[JW_MAX_RANK_LEN];          //��ӯ������
    unsigned int PB;                               //�о��ʣ������ x100000���ͻ��� /100000
    unsigned int PB_HY;                            //��ҵ�о��ʣ������ x100000���ͻ��� /100000
    char PB_RANKING[JW_MAX_RANK_LEN];              //�о�������
    unsigned int PEG;                              //PEG������� x100000���ͻ��� /100000
    unsigned int PEG_HY;                           //��ҵPEG������� x100000���ͻ��� /100000
    char PEG_RANKING[JW_MAX_RANK_LEN];             //PEG����
};

struct JWFundOrganizationViewRsp:public JWCommonRsp
{
    unsigned short StockNumber;      //��Ʊ����
    JWFundOrganizationViewData FundOrganizationViewData[JW_MAX_STOCK_NUM]; //�����۵���������

};


#pragma pack(pop)



#endif

