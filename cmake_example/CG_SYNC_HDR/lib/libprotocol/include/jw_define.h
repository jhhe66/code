/********************************************
//�ļ���:jw_define.h
//����:ϵͳ����ͷ�ļ�
//����:�Ӻ���
//����ʱ��:2010.07.07
//�޸ļ�¼:

*********************************************/

#ifndef __JW_DEFINE_H__
#define __JW_DEFINE_H__

namespace jw_sourcetype
{    
    /*��Դ���Ͷ���7bit,��1bitΪServer�ڲ�ʹ��*/
    //========Client 1-31,0��ʱ����=============
    const unsigned char SOURCE_TYPE_MIN_IM = 1;
    const unsigned char SOURCE_TYPE_MAX_IM = 31;

    const unsigned char SOURCE_TYPE_IM_CLIENT = 1;//IM
    const unsigned char SOURCE_TYPE_HUAGUJIE = 2;//�����Ͻ�
    const unsigned char SOURCE_TYPE_WEB_GBS_IMPORT = 3;//WEB �ɲ�ʿ������ʺ�

    //========Web 32-63=======================
    const unsigned char SOURCE_TYPE_MIN_WEB = 32;
    const unsigned char SOURCE_TYPE_MAX_WEB = 63;

    const unsigned char SOURCE_TYPE_WEB_SNS = 32;//SNS����
    const unsigned char SOURCE_TYPE_WEB_HOME = 33;//788111�Ż�
    const unsigned char SOURCE_TYPE_WEB_TAO = 34;//�Թ�����
    const unsigned char SOURCE_TYPE_WEB_RT = 35;//���ɴ���
    const unsigned char SOURCE_TYPE_WEB_GBS = 36;//WEB�ɲ�ʿ

    //=========DYJ Client 64-127==================
    const unsigned char SOURCE_TYPE_DYJ_MIN_CLIENT = 64;
    const unsigned char SOURCE_TYPE_DYJ_MAX_CLIENT = 127;

    const unsigned char SOURCE_TYPE_DYJ_CLIENT_XGW = 64;//ѡ����
    const unsigned char SOURCE_TYPE_DYJ_CLIENT_PCW = 66;//������
    const unsigned char SOURCE_TYPE_DYJ_CLIENT_GGW = 68;//��Ʊ�ܼ�
    const unsigned char SOURCE_TYPE_DYJ_CLIENT_GBS = 70;//�ɲ�ʿ
    const unsigned char SOURCE_TYPE_DYJ_CLIENT_MOBILE = 72;//�ֻ��ն�
    const unsigned char SOURCE_TYPE_DYJ_CLIENT_TAO = 74;//�Թ����õȽ����ն�

    const unsigned char SOURCE_TYPE_STOCK_GAME_MANUAL = 126;//���ɴ����ֹ���ͨ���ʺ�
    const unsigned char SOURCE_TYPE_DYJ_CLIENT_MANUAL = 127;//�ֹ���ͨ���ʺ�

    //====================================
}

namespace jw_command
{

    const unsigned short JW_SUB_DEFAULT     = 0x0000;//Ĭ����������

    //ע�⣺����ûд�������ֵģ���������Ĭ����0
    const unsigned short JW_KEEP_ALIVE      = 0x0005;//����
    const unsigned short JW_REG_PRE_REG = 0x0801;   //ע���˻�
    const unsigned short JW_REG_ACTIVATE = 0x0802; //�����˻�
    const unsigned short JW_ACCOUNT_STATUS  = 0x0803;//��ѯ�ʺ�״̬
    const unsigned short JW_REG_GET_REGINFO = 0x0804;//ȡ�û�ע����Ϣ
    const unsigned short JW_ACCOUNT_REG_TIME =0x0806;//ȡ���ʺ�ע��ʱ��
    const unsigned short JW_CREATE_ACCOUNT  = 0x0809;//��ͨ�ʺ�
    const unsigned short JW_STAT_REPORT = 0x0810;//�û������ϱ�

    const unsigned short JW_AUTH_LOGIN      = 0x0400;//��֤��¼
    const unsigned short JW_AUTH_CHECK_TICKET = 0x0401;//��֤Ʊ��
    const unsigned short JW_AUTH_LOGOUT     = 0x0402;//ע����½
    const unsigned short JW_AUTH_MODI_PWD = 0x0403;//�޸��û�����
    const unsigned short JW_AUTH_UPDATE_TICKET = 0x0405;//�����û�Ticket
    const unsigned short JW_AUTH_GET_TICKET_UIN = 0x0406;//ͨ��UINȡ����֤Ʊ��
    const unsigned short JW_AUTH_RESET_PWD     = 0x0409;//�����û�����
    

    const unsigned short JW_STOCK           = 0x000F;//��ѡ����������
        const unsigned short JW_SUB_STOCK_ADD                   = 0x0001;//������ѡ��
        const unsigned short JW_SUB_STOCK_DEL                   = 0x0002;//ɾ����ѡ��
        const unsigned short JW_SUB_STOCK_LIST_BY_GROUPID       = 0x0010;//��ȡ��ѡ���б�
        const unsigned short JW_SUB_SET_STOCK_SORT              = 0x0011;//������ѡ���������
        const unsigned short JW_SUB_GET_STOCK_SORT              = 0x0012;//��ȡ��ѡ���������
        const unsigned short JW_SUB_GET_STOCK_RANK              = 0x0013;//��ȡ��Ʊ����˳��

    const unsigned short JW_BIND_MOBILE = 0x5030;//���ֻ�ҵ����������
#if 0
        const unsigned short JW_SUB_BIND_UNBIND_APPLY           = 0x0001;//��/����ֻ���ҵ��
        const unsigned short JW_SUB_CHECK_BIND_UNBIND           = 0x0002;//У���(�����)�ֻ�
        const unsigned short JW_SUB_GET_BIND_LIST               = 0x0003;//��ȡ�û��󶨵��ֻ���Ϣ
#endif
        const unsigned short JW_SUB_BIND_UNBIND                 = 0X0004;//��/����ֻ���ҵ��
        const unsigned short JW_SUB_CHECK_VERIFY_CODE           = 0x0005;//У���(�����)�ֻ���֤��
        const unsigned short JW_SUB_GET_BIND_INFO               = 0x0006;//��ȡ�û��󶨵��ֻ���Ϣ

        const unsigned short JW_INFO            =   0x5000;//�û�����
        const unsigned short JW_SUB_GET_COMM_INFO = 0x0010;//ȡ�û�����������Ϣ��ͷ��URL
        const unsigned short JW_SUB_SET_COMM_INFO = 0x0011;//�����û�����������Ϣ
        const unsigned short JW_SUB_GET_MOBILE_AREA=0x0012;//ͨ���û�ID��ѯ�ֻ��ʺŵĹ��޵�
        const unsigned short JW_SUB_GET_MOBILE_AREA_BY_NUMBER = 0x0013;//ͨ���ֻ������ѯ���޵�
        const unsigned short JW_SUB_GET_USER_ACCOUNT_INFO = 0x0014;//ͨ��UINȡ���û��ʺ���Ϣ
        const unsigned short JW_SUB_SET_ALIAS = 0x0018;//�����û��ı�����Ϣ
        const unsigned short JW_SUB_GET_ALIAS = 0x0019;//ȡ���û��ı�����Ϣ
        
        const unsigned short JW_PAY_MNG             = 0x5060;
        const unsigned short JW_SUB_PAY_TICKETINCOME_NOTIFY = 0x001B;


////////////////////////////////////// ��Ʊ�����������ֶ���/////////////////////////////////
        const unsigned short JW_STK_WEB_PULL_MARKET_DATA = 0x0A09;  //��ȡ������������
        const unsigned short JW_STK_WEB_PULL_MINUTE_DATA = 0x0A0A;  //����ʱ������
        const unsigned short JW_STK_WEB_RADAR_INDICATOR  = 0x0A10; //�����״�ͼָ��
        const unsigned short JW_STK_WEB_EARLY_WARNING  =  0X0A1D;  //Ԥ����������
        const unsigned short JW_STK_WEB_BATCH_BRIEF_MARKET_DATA = 0x0A22; //������Ҫ����������
        const unsigned short JW_STK_WEB_NEWS_CONTENT  =  0x0A28;  //��ֻ��Ʊ������������
        const unsigned short JW_STK_WEB_BATCH_NEWS_CONTENT  = 0x0A29; //������ȡ��ֻ��Ʊ�������ű���
        const unsigned short JW_STK_HOME_GET_REFERENCE_DATA = 0x0A60; //ȡreference data
        const unsigned short JW_STK_WEB_COST_DISTRIBUTION_INDICATOR = 0x0A17;//���ɳɱ��ֲ�ָ�꣨? Business ��������0x0A17

        //add by liudaile at 2010.10.20
        const unsigned short JW_STK_WEB_COST_DISTRIBUTION_DETAIL = 0x0A18; //���ɳɱ��ֲ���ϸ��? Business ��������0x0A18
        const unsigned short JW_STK_WEB_STOCK_FUND_FLOW = 0x0A50; //�������ʽ���������
        const unsigned short JW_STK_WEB_INDUSTRY_FUND_FLOW = 0x0A53; //����ҵ�ʽ���������
        const unsigned short JW_STK_WEB_DDE_DATA = 0x0A61; // ��DDE���ݣ�?Business �������� 0x0A61
        const unsigned short JW_STK_WEB_DDSORT_DATA= 0x0A62; //��DDSort���ݣ�?Business �������� 0x0A62
        const unsigned short JW_STK_WEB_DDE_HIS_DATA = 0x0A64; //��DDE��ʷ���ݣ�?Business �������� 0x0A64
        const unsigned short JW_STK_WEB_RISE_FALL_BREADTH= 0x0A65;//��ʵʱ�����ǵ������ݣ�?Business �������� 0x0A65
        const unsigned short JW_STK_WEB_ZHZD_INFO = 0x0A40; //��ȡ�ۺ������Ϣ(Business �������� 0x0A40
        //end add
        //add by liudaile at 2010.11.11
        const unsigned short JW_STK_WEB_STOCK_FUND_DETAIL = 0x0A66;//�����ʽ�������ϸ����  0x0A66
        const unsigned short JW_STK_WEB_KLINE_FORM_ENGINE = 0x0A67;//K����̬��������0x0A67
        const unsigned short JW_STK_WEB_MULTI_EMPTY_GAME_CAMP = 0x0A68;//��ղ�����Ӫ����(Business ������) 0x0A68

        const unsigned short JW_STK_WEB_TECHNICAL_TREND_ANALYSIS = 0x0A69;//������-���Ʒ�������
        const unsigned short JW_STK_WEB_FINANCIAL_MAIN_ANALYSIS_COST = 0x0A70;//�ʽ���-��������-�����ɱ�����
        const unsigned short JW_STK_WEB_TECHNICAL_MARKET_PERFORMANCE = 0x0A71; //������-�г���������
        const unsigned short JW_STK_WEB_FUNDAMENTALS_FINANCIAL_ANALYSIS = 0x0A72;//������-�����������
        const unsigned short JW_STK_WEB_FUNDAMENTALS_FA_PROFIT_TARGETS = 0x0A73;//������-�������-����ָ������
        const unsigned short JW_STK_WEB_FUNDAMENTALS_FA_INDUSTRY_RANKINGS = 0x0A74;//������-�������-��ҵ��������
        const unsigned short JW_STK_WEB_FUNDAMENTALS_ALUATION_RESEARCH = 0x0A75;//������-��ֵ�о�����
        const unsigned short JW_STK_WEB_FUNDAMENTALS_ORGANIZATION_VIEW = 0x0A76;//������-�����۵�����
////////////////////////////////////// ������Ʊ�����������ֶ���/////////////////////////////////
}

namespace jw_bind_type
{
    //�����Ͷ���
    const unsigned char JW_BIND_TYPE_EW_GBS = 1;//�ɲ�ʿ����Ԥ����Ϣ���ֻ�����

}


namespace jw_errorcode
{
    /*ȫ�ִ����붨��*/    
    /*    
    const unsigned int RET_OK = 0;//�ɹ�   
    const unsigned int ERR_PROTOCOL = 1;//Э�����    
    const unsigned int ERR_ENC_TYPE = 2;//��֧�ֵļ��ܷ�ʽ
    const unsigned int ERR_PROTOCOL_VER = 3;//��֧�ֵİ汾��
    const unsigned int ERR_DB = 4;//�������ݿ����
    const unsigned int ERR_NETWORK = 5;//�������
    const unsigned int IM_ERR_UNKNOWN = 0x0000FFFF;//δ֪����
    */
    
    /* ��֤��ش����붨��*/    
    const unsigned int RET_AUTH_ACCT_NOTEXIST = 0x00080001;//�ʺŲ�����
    const unsigned int RET_AUTH_ACCT_DISABLE = 0x00080002;//�ʺű�����
    const unsigned int RET_AUTH_ACCT_NOTACT = 0x00080003;//�ʺ�δ����
    const unsigned int RET_AUTH_ACCT_FREEZ = 0x00080004;//�ʺű�����(��ʱ���Ƶ�½)
    const unsigned int RET_AUTH_PWD_ERROR = 0x00080005;//�������
    const unsigned int RET_AUTH_TICKET_INVALID = 0x00080006;//Ʊ�ݲ���ȷ
    const unsigned int RET_AUTH_TICKET_EXPIRE = 0x00080007;//Ʊ���ѹ���
    const unsigned int RET_AUTH_SYSTEM_ERROR = 0x00080008;//��֤����������
    const unsigned int RET_AUTH_PROTOCOL_ERROR = 0x00080009;//��֤Э�鲻��ȷ    
    const unsigned int RET_AUTH_DP_ERROR = 0x0008000b;//����ƽ̨����

    /* ע����ش����붨��*/
    const unsigned int REG_ERR_VERIFY               = 0x00150001;  //������֤ʧ��
    const unsigned int REG_ERR_ALREADY_REG          = 0x00150002;  //��ע����
    const unsigned int REG_ERR_NEED_ACTIVATE        = 0x00150003;  //��Ҫ����
    const unsigned int REG_ERR_NO_AVA_ID            = 0x00150004;  //û�п���ID��
    const unsigned int REG_ERR_DB                   = 0x00150005;  //DB����
    const unsigned int REG_ERR_ACTIVATE_NOT_FOUND   = 0x00150006;  //û���ҵ�������Ϣ
    const unsigned int REG_ERR_ACTIVATE_VERIFY_ERR  = 0x00150007;  //�����벻��ȷ
    const unsigned int REG_ERR_CHGACT_PWD_ERR       = 0x00150008;  //�޸��˺�ʱ�������
    const unsigned int REG_ERR_CHGACT_ACT_ERR       = 0x00150009;  //�޸��˺�ʱ���˺��Ѿ���ռ��
    const unsigned int REG_ERR_CHGACT_TIME_ERR      = 0x0015000A;  //�޸��˺�ʱ��������ʱ��δ��
    const unsigned int REG_ERR_CHGACT_ARG_ERR       = 0x0015000B;  //�޸��˺�ʱ�͹����Ĳ�������.uin����
    const unsigned int REG_ERR_CHGACT_CTIME_ERR     = 0x0015000C;  //�޸��˺�ʱ�˺��޸�ʱ�����
    const unsigned int REG_ERR_REPEAT_REG           = 0x0015000D;//�ʺ���ע��
    const unsigned int REG_ERR_INVALID_ACCOUNT      = 0x0015000E;//��Ч���ʺ�


    /*������ش����붨��*/
    const unsigned int RET_INFO_SYS_ERR = 0x00020001; //ϵͳ����
    const unsigned int RET_INFO_FIELD_TYPE_ERR = 0x00020002;//�ֶ����ʹ���
    const unsigned int RET_INFO_FIELD_LEN_ERR = 0x00020003;//�ֶγ��ȴ���
    const unsigned int RET_INFO_MOBILE_ERR = 0x00020004;//������ֻ������ʽ
    const unsigned int RET_INFO_MOBILE_NOT_FOUND = 0x00020005;//û���ҵ���Ӧ�Ĺ�����

    /*��ѡ����ش����붨��*/  
    const unsigned int RET_STK_SYS_ERR = 0x00090001; //ϵͳ����
    const unsigned int RET_STK_DB_ERR = 0x00090002;//DB����
    const unsigned int RET_STK_PROTOCOL_ERR = 0x00090003;//Э�����
    const unsigned int RET_STK_COUNT_MAX_LIMIT_ERR = 0x00090004; //�ﵽϵͳ�������������
    const unsigned int RET_STK_STRING_LEN_LIMIT = 0x00090005; //���������string�ֶεĳ���
    const unsigned int RET_STK_NOTFOUND = 0x00090006;//û�ҵ�
    const unsigned int RET_STK_OR_GRPUP_ERROR = 0x00090007;//��Ʊ��Ϣ������ѡ�ɷ��鲻����
    const unsigned int RET_STK_REPEATED = 0x00090008;//��Ʊ�Ѿ�����
    const unsigned int RET_STK_DEFALUT_ERR = 0x00090009;//Ĭ�Ϸ��鲻�ܸ��Ļ���ɾ��
    const unsigned int RET_STK_NO_STKRANK_ERR = 0x0009000A;//û�ж�Ӧ�Ĺ�Ʊ����
    const unsigned int RET_STK_NO_STKSORT = 0x0009000B;//û����ѡ���������

    /*�ֻ�����ش����붨��*/
    const unsigned int RET_BIND_SYSTEM_ERR           = 0x02030001;//ϵͳ����
    const unsigned int RET_BIND_PROTOCOL_ERR         = 0x02030002;//Э�����
    const unsigned int RET_BIND_DB_ERR               = 0x02030003;//DB����
    const unsigned int RET_BIND_UNSUPPORT_ERR        = 0x02030004;//��֧�ֵ����ݸ�ʽ
    const unsigned int RET_BIND_ACCOUNT_ERR          = 0x02030005;//�˺Ų�����
    const unsigned int RET_BIND_INVALID_CODE_ERR     = 0x02030006;//��Ч����֤��
    const unsigned int RET_BIND_NUMBER_LIMIT_ERR     = 0x02030007;//�����󶨺������
    const unsigned int RET_BIND_MOBILE_UNBIND_ERR    = 0x02030008;//Ҫ����󶨵��ֻ�û�а��κ�Ӧ��
    const unsigned int RET_BIND_MOBILE_BINDED_ERR    = 0x02030009; //���ֻ������Ѿ��������û���
}


namespace jw_ticket_type
{
    const unsigned char TICKET_TYPE_JUMP_HUAGU_URL = 2;//��ת����URLƱ��
    const unsigned char TICKET_TYPE_CLIENT_GBS = 128;//�ɲ�ʿ��֤Ʊ��
    const unsigned char TICKET_TYPE_TERMINAL = 127;//�ն���֤Ʊ��
    
}


#endif
