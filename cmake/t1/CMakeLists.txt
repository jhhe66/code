project(t1)

include_directories(liba libb)
link_directories(liba libb)

add_executable(t1 main.c)

add_definitions(-g -O2 -Wall)

target_link_libraries(t1 a b)
