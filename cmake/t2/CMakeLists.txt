cmake_minimum_required(VERSION 2.8)

project(t2)

set(LIBRARY_OUTPUT_PATH lib)
set(EXECUTABLE_OUTPUT_PATH bin)

cmake_policy(SET CMP0015 NEW)

set(LIBA_SRC 
	liba/liba.c
)
add_library(a SHARED ${LIBA_SRC})
set_target_properties(a PROPERTIES VERSION 1.0 SOVERSION 1) #这里设置版本与soname 并生成连接
add_library(a_static STATIC ${LIBA_SRC})
set_target_properties(a_static PROPERTIES OUTPUT_NAME a) # 这里通过修改输出的名称完成静态库的生成

set(LIBB_SRC 
	libb/libb.c
)
add_library(b SHARED ${LIBB_SRC})
add_library(b_static STATIC ${LIBB_SRC})
set_target_properties(b_static PROPERTIES OUTPUT_NAME b)


include_directories(liba libb)
link_directories(lib)

add_executable(t2 main.c)

add_definitions(-g -O2 -Wall)

target_link_libraries(t2 a b)
