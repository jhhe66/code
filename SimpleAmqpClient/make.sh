#!/bin/sh

g++ main.cc -I../../job/jpush-server/server/external-el7/src/SimpleAmqpClient/src/ -L/opt/push/lib64 -lboost_chrono-mt -L../../job/jpush-server/server/external-el7/lib/ -lSimpleAmqpClient -lrabbitmq
