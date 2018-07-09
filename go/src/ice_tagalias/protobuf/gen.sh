#!/bin/sh

protoc --go_out=. tagalias_common/tagalias_common.proto 
protoc --go_out=. TagAliasBatchQuery/tagalias_query_interface.proto

sed -i "s/\"tagalias_common\"/\"ice_tagalias\/protobuf\/tagalias_common\"/g" TagAliasBatchQuery/tagalias_query_interface.pb.go

cd tagalias_common
go install

cd ../TagAliasBatchQuery
go install

cd ..
