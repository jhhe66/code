#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <uuid/uuid.h>
#include "xxhash.h"

static int
save2cache(char* keys)
{
	/*生成uuid*/
	uuid_t uu;
	uuid_generate(uu);
	memset(keys, 0x00, sizeof(keys));
	
	uuid_unparse(uu, keys);
	//sprintf(keys, "%x", uu);
	return 0;
}

int
main(int argc, char** argv)
{
	unsigned long hash_value = 0;
	char uuid[64] = {0};
	
	save2cache(uuid);
	printf("uuid: %s\n", uuid);

	hash_value = XXH64(uuid, strlen(uuid), 0);
	
	printf("hash: %llu\n", hash_value);

	return 0;
}
