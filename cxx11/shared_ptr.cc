#include <stdio.h>
#include <memory>

int
main(int argc, char** argv)
{
	std::shared_ptr<int> sp1;

	sp1 = std::make_shared<int>(10);

	printf("sp1: %d\n", *sp1);

	return 0;
}
