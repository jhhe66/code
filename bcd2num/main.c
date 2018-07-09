#include <stdlib.h>
#include <stdio.h>
#include <arpa/inet.h>


static void
b2n()
{
	unsigned int temp = 0;

	sscanf("1a2b3c4d", "%08x", &temp);

	printf("temp: %#x\n", temp);
}

static void
b2n_2()
{
	unsigned int temp = 0;
	
	temp = strtoul("1a2b3c4d", NULL, 16);
	
	printf("temp: %#x\n", temp);
}

int 
main(int argc, char** argv)
{
	char bcd[9] = "1a2b3c4d";
	
	unsigned int temp = 0;
	unsigned char* p = (unsigned char*)&temp;
	unsigned int uc1, uc2, uc3, uc4;
	//sscanf(bcd, "%02x%02x%02x%02x", p + 3, p + 2, p + 1, p);
	//sscanf(bcd, "%02x%02x%02x%02x", p, p + 1, p + 2, p + 3);
	//sscanf(bcd, "%02x%02x%02x%02x", (unsigned char*)&temp, (unsigned char*)&temp + 1, (unsigned char*)&temp + 2, (unsigned char*)&temp + 3);
	sscanf(bcd, "%02x%02x%02x%02x", &uc1, &uc2, &uc3, &uc4);

	printf("temp: %u\n", temp);
	printf("temp: %x\n", htonl(temp));
	printf("temp: %x\n", *p);
	printf("temp: %x\n",*(p + 1));
	printf("temp: %x\n", *(p + 2));
	printf("temp: %x\n", *(p + 3));

	printf("temp: %x\n", uc1);
	printf("temp: %x\n", uc2);
	printf("temp: %x\n", uc3); 
	printf("temp: %x\n", uc4);

	//b2n();
	b2n_2();

	return 0;
}
