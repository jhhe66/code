#include <boost/algorithm/string.hpp>
#include <boost/algorithm/string/iter_find.hpp>
#include <stdio.h>
#include <string.h>
#include <vector>

using namespace std;
using namespace boost::algorithm;

int
main(int argc, char* argv[])
{
	string str("hello");

	vector<string> result;

	boost::iter_split(result, str, first_finder("|"));
	
	printf("sz: %d\n", result.size());
	printf("str: %s\n", result[0].c_str());

	return 0;
}
