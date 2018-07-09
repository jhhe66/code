#include "ShmMap.h"
#include "ServerDef.h"

#include <sys/ipc.h>
#include <sys/shm.h>


const unsigned int SHM_FLAG = IPC_CREAT | 0666;

const unsigned int SHM_ITEM_NUM = 5000;
const unsigned int SHM_SIZE = SHM_ITEM_NUM * sizeof(MapItem);


int
CShmMap::init()
{
	shm_id_ = shmget(key_, SHM_SIZE, SHM_FLAG);

	if (shm_id_ == -1)
		return -1;

	shm_head_ = (MapItem*) shmat(shm_id_, NULL, 0);

	if(!shm_head_)
		return -2;

	return 0;
}

unsigned short
CShmMap::GetIndexByCode(const string & code)
{
	return FindIndexByCode(code);
}

string
CShmMap::GetCodeByIdx(unsigned short index)
{
	return FindCodeByIndex(index);
}

//private

unsigned short
CShmMap::FindIndexByCode(const string& code)
{
	CODE_KEY_MAP::const_iterator it;

	//char* c = (char*) code.c_str();
	
	it = code_key_list_.find(code);

	if (it == code_key_list_.end())
		return 0;

	return it->second->index;
}

string
CShmMap::FindCodeByIndex(unsigned short index)
{
	INDEX_KEY_MAP::const_iterator it;

	it = index_key_list_.find(index);

	if (it == index_key_list_.end())
		return string("");

	return string(it->second->code);
}

int
CShmMap::Restore()
{
	MapItem* p = shm_head_;

	unsigned int iCurr = 0;

	while(iCurr <= SHM_ITEM_NUM && strlen(p->code) > 0)
	{
		index_key_list_.insert(pair<unsigned short, MapItem*>(p->index, p));
		code_key_list_.insert(pair<string, MapItem*>(p->code, p));

		p++;

		iCurr++;
	}

	offset_ = iCurr;

	return offset_;
}

void
CShmMap::Insert(unsigned short index,const string & code)
{
	int iCurr = offset_;

	MapItem* p = shm_head_ + iCurr;

	strcpy(p->code, code.c_str());
	p->index = index;

	index_key_list_.insert(pair<unsigned short, MapItem*>(p->index, p));
	code_key_list_.insert(pair<string, MapItem*>(p->code, p));

	offset_++;
}

void
CShmMap::Clear()
{
	index_key_list_.clear();
	code_key_list_.clear();

	bzero(shm_head_, SHM_SIZE);
	
	offset_ = 0;
}

string
CShmMap::GetCodeByOffset(unsigned short offset)
{
	return offset > offset_ ? 
		string("") : (shm_head_ + offset)->code;
}

unsigned short
CShmMap::GetIndexByOffset(unsigned short offset)
{
	return offset > offset_ ?
		0 : (shm_head_ + offset)->index;
}
