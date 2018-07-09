#include <stdio.h>
#include <stdlib.h>

typedef struct human_s human_t;
struct human_s {
	int 	age;
	char 	name[20];
	char 	color;
};

typedef struct people_s people_t;
struct people_s {
	human_t h;
};

typedef struct white_s white_t;
struct white_s {
	human_t h;
	int 	hight;
};

typedef struct yellow_s yellow_t;
struct yellow_s {
	human_t h;
	int 	weight;
};

#define __LOC(h) do { 															\
	if (h) {																	\
		human_t *hm = (human_t*)h;												\
		printf("age: %d name: %s color: %d\n", hm->age, hm->name, hm->color);	\
	}																			\
} while(0)


static void
__SEARCH(human_t* h)
{
	printf("age: %d name: %s color: %d\n", h->age, h->name, h->color);	
}

int
main(int argc, char** argv)
{
	human_t *h = NULL;

	white_t *w = calloc(1, sizeof *w);
	yellow_t *y = calloc(1, sizeof *y);

	w->h.age = 30;
	y->h.age = 40;

	__LOC(w);
	__LOC(y);
	__LOC((white_t*)y);
	__SEARCH((human_t*)w);
	__SEARCH((human_t*)y);

	return 0;
}
