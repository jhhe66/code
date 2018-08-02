#ifndef __FUNC_H_
#define __FUNC_H_

#ifdef __cplusplus
extern "C"
{
#endif

void show() __attribute__((visibility ("hidden")));

#ifdef __cplusplus
}
#endif

#endif
