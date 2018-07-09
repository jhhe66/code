#include <stdio.h>
#include <vector>

using namespace std;

class Data {
public:
	Data(int id):_id(id) { printf("%s\n", __PRETTY_FUNCTION__); }
	Data() {}
	~Data() { printf("%s:%d\n", __PRETTY_FUNCTION__, _id); }
	
	int ID() { return _id; }
private:
	int _id;
};



int 
main(int argc, char **argv)
{
	vector<Data> v1;
	
	for (unsigned int idx = 0; idx < 10; idx++) {
		v1.push_back(std::move(Data(idx)));
	}
	
	printf("\n");
	
	v1.resize(5); //resize 的操作是清空size-newsize 的item 9 ~ 5 保留 0 ~ 4
	
	printf("\n");
	
	printf("size: %zu\n", v1.size());
	
	for (unsigned int idx = 0; idx < 5; idx++) {
		printf("v1[%d]:%d ", idx, v1[idx].ID());
	}	
	printf("\n");
	return 0;
}
