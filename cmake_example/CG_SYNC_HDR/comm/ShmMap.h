#ifndef __SHMMAP_H_
#define __SHMMAP_H_

#include <string>
#include <string.h>
#include <map>

using namespace std;

//const unsigned int SHM_KEY = 0x10000000;


#pragma pack(push, 1)

typedef struct tagMapItem
{
	tagMapItem()
		:index(0)
	{
		bzero(code, 7);
	}

	~tagMapItem() {}

	unsigned short index;
	char code[7];
}MapItem;

typedef struct tagMapHead
{
	unsigned int item_num;		
}MapHead;

#pragma pack(pop)

typedef map<unsigned short, MapItem*> INDEX_KEY_MAP;
typedef map<string, MapItem*> CODE_KEY_MAP;

class CShmMap
{
public:
	CShmMap(unsigned int key)
		:key_(key), 
		shm_id_(-1),
		shm_head_(0),
		offset_(0)
	{}
	~CShmMap() {}

	int init();
	int Restore();

	string GetCodeByIdx(unsigned short index);
	unsigned short GetIndexByCode(const string& code);

	void Insert(unsigned short index, const string& code);
	void Clear();

	unsigned short size() const { return offset_; }
	unsigned short Index_Map_Size() const {return index_key_list_.size();}
	unsigned short Code_Map_Size() const {return code_key_list_.size();}

	string GetCodeByOffset(unsigned short offset);
	unsigned short GetIndexByOffset(unsigned short offset);

private:
	unsigned int key_;
	int shm_id_;
	MapItem* shm_head_;
	unsigned int offset_;

	INDEX_KEY_MAP index_key_list_;
	CODE_KEY_MAP code_key_list_;

private:
	string FindCodeByIndex(unsigned short index);
	unsigned short FindIndexByCode(const string& code);
};

#endif
 
