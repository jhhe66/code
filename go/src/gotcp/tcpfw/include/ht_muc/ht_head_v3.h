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
 *  HTHead 大小固定为48字节 
 *  一个完整的报文：HTV3MAGIC_BEGIN + HTHeadV3 + PayLoad + HTV3MAGIC_END
 *  其中 len_：为整个完整不报文的长度=1 + sizeof(HTHeadV3) + sizeof(PayLoad) + 1
 */
typedef struct HTHeadV3
{
	uint8_t ucFlag;			// 0xF0客户端请求，0xF1 服务器应答, 0xF2  服务器主动发包, 0xF3  客户端应答, 0xF4 服务器之间的包
	uint8_t ucVersion;		// 版本号  VER_MMEDIA = 4
	uint8_t ucKey;			//加密类型  E_NONE_KEY = 0, E_SESSION_KEY = 1, E_RAND_KEY = 2, E_SERV_KEY= 3
	uint8_t ucReserved;		//保留字节, 这个是为了兼容XT_head
	uint16_t usCmd;			//命令字
	uint16_t usSeq;			//序列号
	uint32_t uiFrom;		//源UID  FROM_SERVER = 0	
	uint32_t uiTo;			//目的UID TO_SERVER = 0	
	uint32_t uiLen;			//包总长度
	uint16_t usRet;			//返回码
	uint16_t usSysType;		//包来源
	char acEcho[8];			//回带字段
	char acReserved[16];	//保留字段

    
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

