#ifndef ATTR_API_H
#define ATTR_API_H

#include <stdint.h>
//#include "oi_shm.h"
#define ATTR_SHM_ID_HEX  	0x9527  //�����ڴ��id
#define ATTR_SHM_ID_DEC  	38183   //ʮ���Ʊ�ʾ
#define MAX_ATTR_NODE		1000	//�������֧�����Ը���
#define ATTR_NAME_LEN		64		//�������Ƶ���󳤶�

#define VERSION			"1.1(���Ӱ汾��Ϣ��ȡֵ )"


enum MONITOR_ATTR_TYPE
{
    EM_TYPE_COUNTER     =   0,      //�ۼ�
    EM_TYPE_AVG         =   1,      //ƽ��
    EM_TYPE_STAT        =   2,      //״̬
};

#pragma pack(1)
typedef struct
{
	uint8_t ucUse;						//�Ƿ�ʹ��
	uint32_t uiAttrID;					//����ID
	char acName[ATTR_NAME_LEN];			//��������
	uint8_t ucType;						//��������
	int iCurValue;						//��ǰ����ֵ
	int iTimes;							//���������ڼ���ƽ��ֵ
} AttrNode;
#pragma pack()

typedef struct
{
	AttrNode astNode[MAX_ATTR_NODE];
} AttrList;
#ifdef  __cplusplus
extern "C" {
#endif

int AttrAdd(const char* pName,int iValue);//iValueΪ�ۼ�ֵ
int AttrAddAvg(const char* pName,int iValue);//iValueΪ�ۼ�ֵ
int AttrSet(const char* pName,int iValue);//ֱ�ӽ�iValue���������ڴ�
int GetAttrValue(const char* pName,int *iValue);//�������IDΪattr��ֵ����������ڷ���-1,���ڷ���0������iValue;
void GetAPIversion();

#ifdef  __cplusplus
}
#endif

#endif
