#include <iostream>
#include <list>
#include <vector>
#include <string>
#include <sstream>
#include <unistd.h>
#include <malloc.h>
#include <ext/malloc_allocator.h>

using namespace std;

static void
__use_mmap()
{
	mallopt(M_MMAP_THRESHOLD, 4096);
}

int 
main(int argc, char* argv[]) 
{
    // list<string> mylist;

	//__use_mmap();

    for (int i = 0; i<10; i++) {
		vector<string> vec;

        for (int j = 0; j < 10000000; j++) {
            //stringstream ss;
            //ss << "test_" << j;
           // mylist.push_back(ss.str());
		   vec.push_back("aaaa");
        }

		//sleep(60);
        vec.clear();
		{
			vector<string>().swap(vec);
		}
        cout << "clear" << endl;
    }
	
	//malloc_stats();
	//malloc_trim(0);
    cout << "finish, sleep" << endl;
    sleep(10000);
    return 0;
}
