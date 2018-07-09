#ifndef __ICE_TAGALIAS_H_
#define __ICE_TAGALIAS_H_

#ifdef __cplusplus 
extern "C" {
#endif 

typedef void* ice_warper_t;
typedef char* ice_response_t;
typedef char* ice_request_t;

ice_warper_t
ice_tagalias_init(const char* ice_conf, const char* proxy);

void
ice_tagalias_free(ice_warper_t ice_warper);

int
ice_tagalias_request(ice_warper_t ice_warper,
					 const ice_request_t req,
					 ice_response_t res,
					 unsigned int* sz);

int
ice_tagalias_validateTags(ice_warper_t ice_warper,
						  const ice_request_t req,
						  ice_response_t res,
						  unsigned int* sz);


#ifdef __cplusplus 
}
#endif

#endif

