#include <stdlib.h>
#include <stdio.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>


extern int getId();
extern int getGId();

int
main(int argc, char** argv)
{
	printf("starting...[%d]\n", getpid());
	
	pid_t 	pid = 0;
	int 	status = 0;
	
	pid = fork();
	switch (pid) {
		case 0:
			printf("son[%d]: %d %d\n", getpid(), getId(), getGId());
			sleep(10);
			exit(1);
			break;
		case -1:
			printf("fork failed.\n");
			break;
		default:
			printf("father[%d]: %d %d\n", getpid(), getId(), getGId());
			waitpid(pid, &status, 0); 
			if (WIFEXITED(status)) {
				printf("normly exit[%d]\n", WEXITSTATUS(status));	
			}
			break;
	}
	
	//sleep(10);
	
	return 0;
}
