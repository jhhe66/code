#!/bin/sh

g++ -o libice_warper.so ice_warper.cc QueryUserApi.cpp  -shared -fPIC -I./
#g++ -o libice_warper.so ice_warper.cc -shared -fPIC -I./

#ar -rcu libice_warper.a  ice_warper.o
#gcc -o test_warper test_warper.c -L. -lice_warper -lIce -lIceUtil


#可以使用 go build . 直接编译go 和 相关c库代码 这样的编译方式是采用静态库的连接方式, 执行文件体积大
#go build . 

# 动态库必须提前存在，否则还是静态连接方式
#go build ice_warper.go 
