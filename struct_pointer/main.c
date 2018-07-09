#include <stdio.h>
#include <string.h>

typedef struct human_s human_t;
struct human_s {
	int 	age;
	int 	sex;
	char 	name[20];
};

#define __ID(h) printf("age: %d sex: %d name: %s\n", h->age, h->sex, h->name)

int
main(int argc, char** argv)
{
	human_t hh;
	human_t *phh;
	int 	*page;
	int 	*psex;
	char	*pname;

	hh.age = 38;
	hh.sex = 1;
	strcpy(hh.name, "chenbo");
	
	phh = &hh;

	__ID(phh);
	
	page = (int*)phh;
	psex = /*(int*)phh + 1;*/ (int*)((void*)phh + sizeof(int)); /* page + 1; */ /* (int*)((char*)phh + sizeof(int))  */
	pname = (char*)((int*)phh + 2); // (char*)((void*)phh + 2 * sizeof(int));


	*page = 40;
	*psex = 2;
	strcpy(pname, "dengyong");

	__ID(phh);
	
	return 0;
}
