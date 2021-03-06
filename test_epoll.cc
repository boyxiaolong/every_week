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

        int read_data(){
            int readn = 0;
            bool is_read_error = false;
            while (true)
            {
                int nread = read(fd_, get_data(), get_left_length());
                if (nread < 0)
                {
                    if (errno == EAGAIN)
                    {
                        printf("fd %d read end!", fd_);
                        break;
                    }
                    
                    is_read_error = true;
                    break;
                }
                if (nread == 0)
                {
                    if (readn == 0)
                    {
                        is_read_error = true;
                    }
                    else
                    {
                        printf("readnum 0 so read end");
                    }
                    
                    break;
                }
                
                add_pos(nread);
                readn += nread;
                printf("readnum:%d\n", nread);
            }
                        
            if (is_read_error)
            {
                printf("is_read_error\n");
                return -1;
            }

            process_data();
        }

    private:
        int fd_;
        int state_;
        char* buf_;
        int max_length;
        int cur_pos;
};

class Server
{
    public:
        typedef std::map<int, Socket*> socket_map;
        Server(int ae_fd):ae_fd_(ae_fd){}
        ~Server(){

        }
        virtual int create_server_sock(){
            listen_fd_ = socket(AF_INET, SOCK_STREAM, 0);
            if (listen_fd_ < 1) 
            {
                perror("socket create");
                return -1;    
            }

            set_sock_noblock(listen_fd_);
            struct sockaddr_in serv_addr;
            memset(&serv_addr, 0, sizeof(serv_addr));
            serv_addr.sin_family = AF_INET;
            serv_addr.sin_port = htons (SERVER_PORT);
            serv_addr.sin_addr.s_addr = htonl (INADDR_ANY);
            
            printf("begin bind\n");
            int ret = bind(listen_fd_, (struct sockaddr *) (&serv_addr), sizeof(serv_addr));
            if (ret < 0)
            {
                perror("bind");
                return -1;
            }
            printf("bind success\n");
            
            ret = listen(listen_fd_, max_events);
            if (ret < 0)
            {
                perror("listen");
                return -1;
            }
            printf("listen success %d\n", listen_fd_);

            struct epoll_event ev;
            ev.events = EPOLLIN;
            ev.data.fd = listen_fd_;
            if (epoll_ctl(ae_fd_, EPOLL_CTL_ADD, listen_fd_, &ev) < 0)
            {
                perror("epoll_ctl");
                return -1;
            }
        }

        virtual int ae_accept(){
            printf("begin ae_accept\n");
            int res = 0;
            do
            {
                struct sockaddr_in client_addr;
                int addrlen = 0;
                int new_socket = accept(listen_fd_, (struct sockaddr *)&client_addr,(socklen_t*)&addrlen);
                if (new_socket < 0)
                {
                    if (errno == EAGAIN)
                    {
                        printf("fd %d acceptc nonblock!", ae_fd_);
                        break;
                    }
                    res = -1;
                    break;
                }
                set_sock_noblock(new_socket);
                struct epoll_event client_ev;
                client_ev.events = EPOLLIN;
                client_ev.data.fd = new_socket;
                if (epoll_ctl(ae_fd_, EPOLL_CTL_ADD, new_socket, &client_ev) < 0)
                {
                    res = -1;
                    break;
                }
                printf("accept new socket %d and addto epoll\n", new_socket);
                Socket* ps = new Socket(new_socket);
                socket_map_.insert(std::make_pair(new_socket, ps));
            } while (true);
            printf("end ae_accept\n");
            return res;
        }

        void set_sock_noblock(int& sock)
        {
            int flags = fcntl(sock, F_GETFL);
            flags |= O_NONBLOCK;
            fcntl(sock, F_SETFL, flags);
        }

        Socket* get_sock_ps(int cur_fd){
            socket_map::iterator iter = socket_map_.find(cur_fd);
            if (iter == socket_map_.end())
            {
                return NULL;
            }

            return iter->second;
        }

        virtual void rm_client_sock(int cur_fd){
            socket_map::iterator iter = socket_map_.find(cur_fd);
            if (iter == socket_map_.end())
            {
                return;
            }
            
            Socket* ps = iter->second;
            if (ps)
            {
                delete ps;
            }
            socket_map_.erase(iter);

            epoll_ctl(ae_fd_, EPOLL_CTL_DEL, cur_fd, NULL);
            printf("delete fd:%d client close socket\n", cur_fd);
        }

        virtual int ae_poll(){
            struct epoll_event events[max_events];
            while (true)
            {
                printf("begin epoll\n");
                int nfds = epoll_wait(ae_fd_, events, max_events, -1);
                if (nfds == -1)
                {
                    if (errno == EINTR)
                    {
                        continue;
                    }
                    
                    perror("wait");
                    exit(-1);
                }
                printf("get epoll data %d\n", nfds);
                for (size_t i = 0; i < nfds; i++)
                {
                    epoll_event& cur_ev = events[i];
                    int cur_fd = cur_ev.data.fd;
                    printf("epoll cur_fd %d\n", cur_fd);
                    if (cur_fd == listen_fd_)
                    {
                        ae_accept();
                    }
                    else
                    {
                        Socket* ps = get_sock_ps(cur_fd);
                        if (NULL == ps)
                        {
                            continue;
                        }
                        int res = ps->read_data();
                        if (res < 0)
                        {
                            rm_client_sock(cur_fd);
                        }
                    }
                }
            }
        }
    private:
        int ae_fd_;
        int listen_fd_;
            
        socket_map socket_map_;
};

int main()
{
    int epoll_fd = epoll_create(max_events);
    if (epoll_fd < 0)
    {
        perror("epoll_create");
        exit(-1);
    }
    
    Server ser(epoll_fd);
    int res = ser.create_server_sock();
    if (res < 0)
    {
        printf("create error\n");
        return res;
    }

    ser.ae_poll();
}
