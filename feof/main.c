#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <errno.h>

static int 
write_file(const char* path)
{
	FILE *file;

	file = fopen("data", "a+");
	if (file == NULL) {
		printf("open file failed %d.", errno);
		return -1;
	}
	
	for (int idx = 0; idx < 10; idx++) {
		fprintf(file, "%d", idx);
	}
	
	return 0;	
}

static int
read_file(const char* path)
{
	FILE *file;
	char c;
	
	file = fopen("data", "r+");
	if (file == NULL) {
		printf("open file failed %d.", errno);
		return -1;
	}
	
	while ((c = fgetc(file)) != EOF) {
		printf("C: %c\n", c);
	}

	return 0;
}

static int
read_file_2(const char* path)
{	
	FILE *file;
	char c;
	
	file = fopen("data", "r+");
	if (file == NULL) {
		printf("open file failed %d.", errno);
		return -1;
	}
	
	while (!feof(file)) { //feof 是判断当前的读取是不是结束的，所以会多读一次
		printf("C: %c\n", fgetc(file));
	}

	return 0;
}

int
main(int argc, char** argv)
{
	write_file("data");
	read_file("data");
	read_file_2("data");


	return 0;
}
