#ifndef __ICE_WARPER_H_
#define __ICE_WARPER_H_

#ifdef __cplusplus
extern "C" {
#endif

typedef void* ice_warper_t;
typedef char* ice_response_t; // string*
typedef char* ice_request_t;  // stirng*

ice_warper_t 
ice_warper_init(const char* ice_conf);

void 
ice_warper_free(ice_warper_t ice_warper);

int
ice_warper_valid_appkeys(ice_warper_t ice_warper, 
						 const ice_request_t req, 
						 ice_response_t res, 
						 unsigned int* rsz);

int
ice_warper_valid_users(ice_warper_t ice_warper, 
					   const ice_request_t req, 
					   ice_response_t res, 
					   unsigned int* rsz);

#ifdef __cplusplus
}
#endif

#endif 

