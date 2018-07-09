#!/bin/sh

g++ -o thread_cpp thread.cc -std=gnu++11 -g -Wall -pthread
g++ -o shared_ptr shared_ptr.cc -std=gnu++11 -g -Wall
g++ -o move_constructor move_constructor.cc -std=gnu++11 -g -Wall
g++ -o forward forward.cc -std=gnu++11 -g -Wall
g++ -o emplace_back emplace_back.cc -std=gnu++11 -g -Wall
g++ -o raw_string raw_string.cc -std=gnu++11 -g -Wall
