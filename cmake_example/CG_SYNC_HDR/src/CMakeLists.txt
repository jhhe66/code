INCLUDE_DIRECTORIES(../inc ../comm ../lib/libprotocol/include)

AUX_SOURCE_DIRECTORY(../comm comm_SRCS)
AUX_SOURCE_DIRECTORY(. src_SRCS)

add_executable(cg_sync_hdr ${src_SRCS} ${comm_SRCS})
#INCLUDE_DIRECTORIES(../inc ../comm ../lib/libprotocol/include)
TARGET_LINK_LIBRARIES(cg_sync_hdr jwprotocol json memcached)
