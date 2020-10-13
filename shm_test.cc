#include <stdio.h>
#include <stdlib.h>
#include <sys/ipc.h>
#include <sys/shm.h>
#include "string.h"

#define ipckey 0x366378


struct st_setting
{
        char agen[10];
        unsigned char file_no;
};

int main()
{
        int shm_id = shmget(ipckey, 1028, 0640);
        if(-1 != shm_id){
                system("ipcs -m");
                printf("already has\n");
                st_setting* ps = (st_setting*)(shmat(shm_id, NULL, 0));
                printf("shm data: agen:%s, file_no:%d\n:", ps->agen, ps->file_no);
                int res = shmctl(shm_id, IPC_RMID, NULL);
                printf("del shm res:%d\n", res);
                return 0;
        }
        shm_id = shmget(ipckey, 1028, 0640|IPC_CREAT|IPC_EXCL);
        if(-1 == shm_id){
                printf("create shm fail\n");
                return 0;
        }
        printf("create shm success\n");
        st_setting* ps = (st_setting*)(shmat(shm_id, NULL, 0));
        strncpy(ps->agen, "gate", 10);
        ps->file_no = 1;
        system("ipcs -m");
}