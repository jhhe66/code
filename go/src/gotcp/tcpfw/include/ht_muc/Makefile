C_ARGS    = -O2 -Wall

PROTO=$(wildcard *.proto)
PB_CC=$(patsubst %.proto,%.pb.cc,$(PROTO))
PB_H=$(patsubst %.proto,%.pb.h,$(PROTO))
BINARY=libmuc_proto.a

CXX = g++
CC = gcc
#INC = -I /usr/local/protobuf
#LIB = /usr/local/protobuf/lib/libprotobuf.so.8.0.0

all: $(PB_CC) $(PB_H) $(BINARY)

.cpp.o:
	$(CXX) $(C_ARGS) $(INC) -c $^
	@echo ""

.c.o:
	$(CC) $(C_ARGS) $(INC) -c $^
	@echo ""

%.pb.cc %.pb.h:%.proto
	protoc --cpp_out=. $<
	@echo ""

$(BINARY): $(patsubst %.cc,%.o,$(PB_CC))
	ar rcs $@ $^
	cp $@ ../../lib/
	@echo ""


#install:
#	cp msg_center_proto.a ../../dev/lib
#	cp *.h ../../dev/include/msg_center_proto

clean:
	rm *.o *.pb.cc *.pb.h $(BINARY) ../../lib/$(BINARY)
