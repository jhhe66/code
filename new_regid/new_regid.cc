#include "regid.pb.h"
#include "base64.h"
#include <stdlib.h>
#include <stdio.h>

using namespace JPush;
using namespace std;

#define __ENCODE__ 'e'
#define __DECODE__ 'd'

#define __REG_VERSION__ 1

static void
__new_id(const char* appkey, 
		 const unsigned long uid, 
		 const int platform)
{
	regid 	rid;
	string 	wbuf;
	char 	new_id[256];

	rid.set_version(__REG_VERSION__);
	rid.set_appkey(appkey);
	rid.set_uid(uid);
	rid.set_platform(platform);

	rid.SerializeToString(&wbuf);
	printf("wbuf: %zu data: %s\n", wbuf.size(), wbuf.c_str());
	base64encode(wbuf.data(), wbuf.size(), new_id, 256);
	printf("new_id: %s\n", new_id);
}

static void
__parse_id(const char* reg_id)
{
	size_t 			olen;
	unsigned char 	out[256];
	regid			rid;

	base64decode((char*)reg_id, strlen(reg_id), out, &olen);
	rid.ParseFromArray(out, olen);
	printf("version: %lu, appkey: %s uid: %lu platform: %lu\n", 
		  rid.version(), 
		  rid.appkey().c_str(), 
		  rid.uid(), 
		  rid.platform());
}

int
main(int argc, char** argv)
{
	char 			flag;
	unsigned long 	uid;
	unsigned int 	platform;
	
	flag = *(argv[1]);

	switch (flag) {
		case __ENCODE__:
			uid = (unsigned long)atol(argv[3]);
			platform = (unsigned int)atoi(argv[4]);	
			__new_id(argv[2], uid, platform);
			break;
		case __DECODE__:
			__parse_id(argv[2]);
			break;
		default:
			printf("parameters invailed.\n");
			break;
	}

	return 0;
}
