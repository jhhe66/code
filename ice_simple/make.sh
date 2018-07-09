#!/bin/sh

g++ -o server server.cc Printer.cpp -I. -lIce -lIceUtil
g++ -o client client.cc Printer.cpp -I. -lIce -lIceUtil
