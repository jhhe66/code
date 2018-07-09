#include <unistd.h>

static int id = 0;

int g_id = 0;

int 
getId()
{
	return ++id;
}

int
getGId()
{
	return ++g_id;
}

int self()
{
	return getpid();
}
