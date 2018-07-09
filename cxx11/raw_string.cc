#include <string>
#include <iostream>

using namespace std;

int
main(int argc, char** argv)
{
	string s1 = R"(hello
				)";

	
	cout << s1 << endl;
	cout << L"(aaaa)" << endl;
	cout << U"(中国)" << endl;

	return 0;
}
