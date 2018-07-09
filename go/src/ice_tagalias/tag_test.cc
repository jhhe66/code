#include "ice_tagalias.h"
#include "include/tagalias_query_interface.pb.h"
#include "include/tagalias_common.pb.h"
#include <stdlib.h>
#include <stdio.h>

using namespace TagAliasBatchQuery;

static void
hex_dump(char* buff, size_t sz)
{
	printf("[");
	for (size_t idx = 0; idx < sz; idx++) {
		printf("%u ", (unsigned char)*(buff + idx));
	}
	printf("]\n");
}

int
main(int argc, char** argv)
{
	char 			request[512], response[512];
	char 			req_buff[512];
	TagAliasQuery  	tagAlsQuery;
	TagAlsResp		tagResp;
	unsigned int 	sz;

	tagAlsQuery.set_reqno(10000);
	tagAlsQuery.set_appkey("1fa5281c6c1d4cf5bb0bbbe0");
	tagAlsQuery.set_platform(IPHONE);

	OpsUnit* ops = tagAlsQuery.add_query();
	ops->set_dev_type(NORMAL_TYPE);
	ops->set_page(1);
	ops->set_cmd(CHECK_TAGS_HAS_USERS);
	ops->add_tags("aaaa");
	
	ice_warper_t ice = ice_tagalias_init("--Ice.Config=./client.conf", "TagAliasAccessProxy");
	if (ice == NULL) {
		printf("ice init faild.\n");
		return -1;
	}

	tagAlsQuery.SerializeToArray(req_buff, sizeof req_buff);
	
	ice_tagalias_request(ice, req_buff, response, &sz);
	
	if (!tagResp.ParseFromArray(response, sz)) {
		printf("response parse faild. [%u]\n", sz);
	}
	
	hex_dump(response, sz);
	
	return 0;	
}
