g++ -o libice_tagalias.so -shared -fPIC -O2 -Wall -g -I. ice_tagalias.cc tagAlias.cpp


g++ -o tag_test tag_test.cc include/*.cc -lice_tagalias -lIce -lIceUtil -L. -I. -lprotobuf -O2 -Wall -Wl,-rpath,.
