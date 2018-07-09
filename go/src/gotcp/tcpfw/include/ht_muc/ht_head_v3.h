/***********************************************************************
 * File : ht_head_v3.h
 * Brief: 
 * 
 * History
 * ---------------------------------------------------------------------
 * 2015-12-12     meepowang   1.0    created
 * 
 ***********************************************************************
 */


#ifndef HT_HEAD_V3_H
#define HT_HEAD_V3_H

#ifndef __STDC_FORMAT_MACROS // add for PRIU64
#define __STDC_FORMAT_MACROS
#endif
#include <inttypes.h>
#include <stdint.h>
#include <netinet/in.h>

#if __BYTE_ORDER == __BIG_ENDIAN
#define ntohll(x)       (x)
#define htonll(x)       (x)
#else
#if __BYTE_ORDER == __LITTLE_ENDIAN
#define ntohll(x)     __bswap_64 (x)
#define htonll(x)     __bswap_64 (x)
#endif
#endif


#define HTV3MAGIC_BEGIN 0x0a
#define HTV3MAGIC_END 0x0b

#pragma pack(1)
/**
 *  HTHead ��С�̶�Ϊ48�ֽ� 
 *  һ�������ı��ģ�HTV3MAGIC_BEGIN + HTHeadV3 + PayLoad + HTV3MAGIC_END
 *  ���� len_��Ϊ�������������ĵĳ���=1 + sizeof(HTHeadV3) + sizeof(PayLoad) + 1
 */
typedef struct HTHeadV3
{
	uint8_t ucFlag;			// 0xF0�ͻ�������0xF1 ������Ӧ��, 0xF2  ��������������, 0xF3  �ͻ���Ӧ��, 0xF4 ������֮��İ�
	uint8_t ucVersion;		// �汾��  VER_MMEDIA = 4
	uint8_t ucKey;			//��������  E_NONE_KEY = 0, E_SESSION_KEY = 1, E_RAND_KEY = 2, E_SERV_KEY= 3
	uint8_t ucReserved;		//�����ֽ�, �����Ϊ�˼���XT_head
	uint16_t usCmd;			//������
	uint16_t usSeq;			//���к�
	uint32_t uiFrom;		//ԴUID  FROM_SERVER = 0	
	uint32_t uiTo;			//Ŀ��UID TO_SERVER = 0	
	uint32_t uiLen;			//���ܳ���
	uint16_t usRet;			//������
	uint16_t usSysType;		//����Դ
	char acEcho[8];			//�ش��ֶ�
	char acReserved[16];	//�����ֶ�

    
    HTHeadV3()
    {
    	memset(this, 0, sizeof(HTHeadV3));
    }

    void Reset()
    {
		memset(this, 0, sizeof(HTHeadV3));
    }

	void Hton()
	{
		usCmd = htons(usCmd);
		usSeq = htons(usSeq);
		uiFrom = htonl(uiFrom);
		uiTo = htonl(uiTo);
		uiLen = htonl(uiLen);
		usRet = htons(usRet);
		usSysType = htons(usSysType);
	}

	void Ntoh()
	{
		usCmd = ntohs(usCmd);
		usSeq = ntohs(usSeq);
		uiFrom = ntohl(uiFrom);
		uiTo = ntohl(uiTo);
		uiLen = ntohl(uiLen);
		usRet = ntohs(usRet);
		usSysType = ntohs(usSysType);
	}
}HTHeadV3;
#pragma pack()

#endif

