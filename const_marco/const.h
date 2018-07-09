#define __CONST__(T, V, v) (const T V = v)
#define __CONST_INT__(V, v) __CONST__(int, V, v)

__CONST_INT__(RET_SUCC, 0);
