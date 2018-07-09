cmake_minimum_required(VERSION 2.8)

PROJECT(CG_SYNC_HDR)

INCLUDE_DIRECTORIES(inc comm lib/libprotocol/include)

ADD_SUBDIRECTORY(src bin)
#ADD_SUBDIRECTORY(comm)
