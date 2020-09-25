#include "stdio.h"
#include <sys/epoll.h>
#include <sys/types.h>
#include <sys/socket.h>
#include "stdlib.h"
#include "netinet/in.h"
#include "string.h"
#include "unistd.h"
#include <fcntl.h>
#include <errno.h>

#define max_events 100
#define max_buff 1024
#define SERVER_PORT 9999

int main()
{
    int epoll_fd = epoll_create(max_events);
    if (epoll_fd < 0)
    {
        perror("epoll_create");
        exit(-1);
    }
    
    int listen_sock = socket(AF_INET, SOCK_STREAM, 0);
    if (listen_sock < 1) {
	perror("socket create");
	exit(-1);    
}
    struct sockaddr_in serv_addr;
    memset(&serv_addr, 0, sizeof(serv_addr));
    serv_addr.sin_family = AF_INET;
    serv_addr.sin_port = htons (SERVER_PORT);
    serv_addr.sin_addr.s_addr = htonl (INADDR_ANY);
    
    printf("begin bind\n");
    int ret = bind(listen_sock, (struct sockaddr *) (&serv_addr), sizeof(serv_addr));
    if (ret < 0)
    {
        perror("bind");
        exit(-1);
    }
    printf("bind success\n");
	
    ret = listen(listen_sock, max_events);
    if (ret < 0)
    {
        perror("listen");
        exit(-1);
    }
    printf("listen success\n");

    struct epoll_event ev;
    ev.events = EPOLLIN;
    ev.data.fd = listen_sock;
    if (epoll_ctl(epoll_fd, EPOLL_CTL_ADD, listen_sock, &ev) < 0)
    {
        perror("epoll_ctl");
        exit(-1);
    }
    
    struct epoll_event events[max_events];
    while (true)
    {
	    printf("begin epoll\n");
        int nfds = epoll_wait(epoll_fd, events, max_events, -1);
        if (nfds == -1)
        {
            perror("wait");
            exit(-1);
        }
	    printf("get epoll data %d\n", nfds);
        for (size_t i = 0; i < nfds; i++)
        {
            epoll_event& cur_ev = events[i];
            int cur_fd = cur_ev.data.fd;
            if (cur_fd == listen_sock)
            {
                /* code */
                struct sockaddr_in client_addr;
                int addrlen = 0;
                int new_socket = accept(listen_sock, (struct sockaddr *)&client_addr,(socklen_t*)&addrlen);
                if (new_socket < 0)
                {
                    perror("accept");
                    exit(-1);
                }
                int flags = fcntl(new_socket, F_GETFL);
                flags |= O_NONBLOCK;
                fcntl(new_socket, F_SETFL, flags);
                printf("accept new socket %d and addto epoll\n", new_socket);
                struct epoll_event client_ev;
                client_ev.events = EPOLLIN;
                client_ev.data.fd = new_socket;
                if (epoll_ctl(epoll_fd, EPOLL_CTL_ADD, new_socket, &client_ev) < 0)
                {
                    perror("epoll add clientfd");
                    exit(-1);
                }
            }
            else
            {
                char buf[max_buff];
                int readn = 0;
                bool is_read_error = false;
                while (true)
                {
                    int nread = read(cur_fd, buf + readn, max_buff-readn);
                    if (nread < 0)
                    {
                        if (errno == EAGAIN)
                        {
                            printf("fd %d read end!", cur_fd);
                            break;
                        }
                        
                        epoll_ctl(epoll_fd, EPOLL_CTL_DEL, cur_fd, NULL);
                        printf("delete fd:%d\n", cur_fd);
                        is_read_error = true;
                        break;
                    }
                    if (nread == 0)
                    {
                        if (readn == 0)
                        {
                            epoll_ctl(epoll_fd, EPOLL_CTL_DEL, cur_fd, NULL);
                            printf("delete fd:%d client close socket\n", cur_fd);
                        }
                        else
                        {
                            printf("readnum 0 so read end");
                        }
                        
                        break;
                    }
                    
                    readn += nread;
                    printf("readnum:%d\n", nread);
                }
                
                if (is_read_error)
                {
                    printf("is_read_error\n");
                    continue;
                }
                
                buf[readn] = '\0';
                printf("client fd:%d readn:%d read %s\n", cur_fd, readn, buf);
	        }
        }
    }
    
}
