#ifndef ATTR_API_H
#define ATTR_API_H

#include <stdint.h>
//#include "oi_shm.h"
#define ATTR_SHM_ID_HEX  	0x9527  //共享内存的id
#define ATTR_SHM_ID_DEC  	38183   //十进制表示
#define MAX_ATTR_NODE		1000	//单机最大支持属性个数
#define ATTR_NAME_LEN		64		//属性名称的最大长度

#define VERSION			"1.1(增加版本信息与取值 )"


enum MONITOR_ATTR_TYPE
{
    EM_TYPE_COUNTER     =   0,      //累加
    EM_TYPE_AVG         =   1,      //平均
    EM_TYPE_STAT        =   2,      //状态
};

#pragma pack(1)
typedef struct
{
	uint8_t ucUse;						//是否使用
	uint32_t uiAttrID;					//属性ID
	char acName[ATTR_NAME_LEN];			//属性名称
	uint8_t ucType;						//属性类型
	int iCurValue;						//当前计数值
	int iTimes;							//次数，用于计算平均值
} AttrNode;
#pragma pack()

typedef struct
{
	AttrNode astNode[MAX_ATTR_NODE];
} AttrList;
#ifdef  __cplusplus
extern "C" {
#endif

int AttrAdd(const char* pName,int iValue);//iValue为累加值
int AttrAddAvg(const char* pName,int iValue);//iValue为累加值
int AttrSet(const char* pName,int iValue);//直接将iValue赋给共享内存
int GetAttrValue(const char* pName,int *iValue);//获得属性ID为attr的值，如果不存在返回-1,存在返回0并附给iValue;
void GetAPIversion();

#ifdef  __cplusplus
}
#endif

#endif
