#ifndef __ICE_WARPER_H_
#define __ICE_WARPER_H_

#ifdef __cplusplus
extern "C" {
#endif

typedef void* ice_warper_t;

ice_warper_t 
ice_warper_init(const char* ice_conf);

void 
ice_warper_free(ice_warper_t ice_warper);

int
ice_warper_valid_appkeys(ice_warper_t ice_warper, const char* appkey, unsigned long uid);

#ifdef __cplusplus
}
#endif

#endif 

