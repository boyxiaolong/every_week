#include "stdio.h"
#include <sys/select.h>
#include <sys/types.h>
#include <sys/socket.h>
#include "stdlib.h"
#include "netinet/in.h"
#include "string.h"
#include "unistd.h"
#include <fcntl.h>
#include <errno.h>
#include <map>

#define max_events 100
#define max_buff 1024
#define SERVER_PORT 9999

class Socket
{
    public:
        Socket(int fd):fd_(fd)
        , max_length(1024)
        , cur_pos(0)
        {
            buf_ = new char[max_length];
        }
        ~Socket()
        {
            if (NULL != buf_)
            {
                delete []buf_;
            }
            printf("fd %d dtor\n", fd_);
        }
        
        char* get_data()
        {
            return buf_+cur_pos;
        }

        int get_left_length()
        {
            return max_length-cur_pos;
        }

        void process_data()
        {
            //just f test
            if (get_left_length() > 0)
            {
                *(buf_+cur_pos+1) = '\0';
            }
            
            printf("get data length %d data:%s\n", cur_pos, buf_);
            cur_pos = 0;
        }

        void add_pos(int length)
        {
            cur_pos += length;
        }   

    private:
        int fd_;
        int state_;
        char* buf_;
        int max_length;
        int cur_pos;
};

int main()
{
    fd_set rfds;
    FD_ZERO(&rfds);
    timeval tv;
    tv.tv_sec = 0;
    tv.tv_usec = 0;

    int max_fd = 0;
    
    int listen_sock = socket(AF_INET, SOCK_STREAM, 0);
    if (listen_sock < 1) 
    {
        perror("socket create");
        exit(-1);    
    }
    struct sockaddr_in serv_addr;
    memset(&serv_addr, 0, sizeof(serv_addr));
    serv_addr.sin_family = AF_INET;
    serv_addr.sin_port = htons (SERVER_PORT);
    serv_addr.sin_addr.s_addr = htonl (INADDR_ANY);
    
    printf("begin bind socket %d\n", listen_sock);
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

    max_fd = listen_sock;

    FD_SET(listen_sock, &rfds);

    printf("listen_sock %d isset:%d\n", listen_sock, FD_ISSET(listen_sock, &rfds));

    typedef std::map<int, Socket*> socket_map;
    socket_map socket_map_;
    
    fd_set working_set;
    while (true)
    {
	    printf("begin select\n");
        memcpy(&working_set, &rfds, sizeof(rfds));
        int rc = select(max_fd+1, &working_set, NULL, NULL, NULL);
        if (rc < 0)
        {
            perror("select fail");
            exit(-1);
        }

        if (rc == 0)
        {
            perror("select time out");
            exit(-1);
        }
        
        int total_events = rc;
	    printf("get select data %d max_fd %d total_events %d \n", rc, max_fd, total_events);
        for (size_t i = 0; i <= max_fd && total_events > 0; i++)
        {
            if (!FD_ISSET(i, &rfds))
            {
                printf("i:%d not set\n", i);
                continue;                
            }

            printf("cur fd:%d\n", i);
            --total_events;
            int cur_fd = i;
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
                printf("accept new socket %d and addto select\n", new_socket);
                FD_SET(new_socket, &rfds);
                if (new_socket > max_fd)
                {
                    max_fd = new_socket;
                    printf("max_fd is %d\n", max_fd);
                }
                
                Socket* ps = new Socket(new_socket);
                socket_map_.insert(std::make_pair(new_socket, ps));
            }
            else
            {
                socket_map::iterator iter = socket_map_.find(cur_fd);
                if (iter == socket_map_.end())
                {
                    printf("can find fd %d socket data\n");
                    continue;
                }
                
                Socket* ps = iter->second;
                int readn = 0;
                bool is_read_error = false;
                while (true)
                {
                    int nread = read(cur_fd, ps->get_data(), ps->get_left_length());
                    if (nread < 0)
                    {
                        if (errno == EAGAIN)
                        {
                            printf("fd %d read end!", cur_fd);
                            break;
                        }
                        
                        FD_CLR(cur_fd, &rfds);
                        --max_fd;
                        printf("delete fd:%d max_fd:%d\n", cur_fd, max_fd);
                        socket_map_.erase(iter);
                        delete ps;
                        ps = NULL;
                        is_read_error = true;
                        break;
                    }
                    if (nread == 0)
                    {
                        if (readn == 0)
                        {
                            FD_CLR(cur_fd, &rfds);
                            --max_fd;
                            printf("delete fd:%d max_fd:%d\n", cur_fd, max_fd);
                            socket_map_.erase(iter);
                            delete ps;
                            ps = NULL;
                        }
                        else
                        {
                            printf("readnum 0 so read end");
                        }
                        
                        break;
                    }
                    
                    ps->add_pos(nread);
                    readn += nread;
                    printf("readnum:%d\n", nread);
                }
                
                if (is_read_error)
                {
                    printf("is_read_error\n");
                    continue;
                }
                
                if (ps)
                {
                    ps->process_data();
                }
	        }
        }
    }
    
}
