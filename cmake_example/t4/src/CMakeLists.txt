include(${CMAKE_SOURCE_DIR}/../cmake/libhello.cmake)

add_executable(main main.c)
TARGET_INCLUDE_DIRECTORIES(main PUBLIC /usr/local/include/hello)
TARGET_LINK_LIBRARIES(main hello)
