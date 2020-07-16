## redis主逻辑流程

阅读源码首先找到main函数，redis的main函数在server.c。这里忽略与主逻辑影响不大的代码，主要弄清楚redis server如何建立tcp server，接受客户端连接并处理相应请求的。

首先看下main函数中最重要的几个函数，可以大概明白server运行的流程：

~~~c
```
    //从代码中定义的宏变量初始化全局变量server，将所有客户端命令set/get等放进命令字典server.orig_commands中
    initServerConfig();
	//根据配置文件设置server的变量
    loadServerConfigFromString(config);
	//初始化部分数据结构
    initServer();
	//初始化后台线程
    InitServerLast();
	//从rdb或者aof文件加载数据
	loadDataFromDisk();
	//运行事件循环（处理读写事件，timer等）
	aeMain(server.el);
```
~~~

接下来分别看下这几个函数的内部实现，可以了解下实现细节,这里依然忽略了不影响主逻辑的部分细节（比如主从功能，异常处理等）：

* 初始化Server

  ```c
  void initServer(void) {
      int j;
  	//注册进程中止和结束的信号处理函数（比如在shell下按ctrl+c redis进程就会结束并打印结束信息）
      setupSignalHandlers();
  
      //根据配置初始化
      
      //是否开启aof功能
      server.aof_state = server.aof_enabled ? AOF_ON : AOF_OFF;
      //定时函数的执行频率
      server.hz = server.config_hz;
      server.pid = getpid();
      server.current_client = NULL;
      server.fixed_time_expire = 0;
      //接受的客户端链表集合
      server.clients = listCreate();
      server.clients_index = raxNew();
      //需要异步关闭的客户端集合
      server.clients_to_close = listCreate();
      server.clients_pending_write = listCreate();
      server.clients_pending_read = listCreate();
      server.clients_timeout_table = raxNew();
      server.unblocked_clients = listCreate();
      server.ready_keys = listCreate();
      server.clients_waiting_acks = listCreate();
  
      //创建事件循环实例,创建maxclients+CONFIG_FDSET_INCR个tcp socket
      server.el = aeCreateEventLoop(server.maxclients+CONFIG_FDSET_INCR);
      
      //创建配置数量的db实例的内存（默认dbnum是16）
      server.db = zmalloc(sizeof(redisDb)*server.dbnum);
  
      //创建指定ip，port的tcp server ipv4和ipv6 非阻塞socket
      listenToPort(server.port,server.ipfd,&server.ipfd_count)
  
      //初始化db实例的状态
      for (j = 0; j < server.dbnum; j++) {
          server.db[j].dict = dictCreate(&dbDictType,NULL);
          server.db[j].expires = dictCreate(&keyptrDictType,NULL);
          server.db[j].expires_cursor = 0;
          server.db[j].blocking_keys = dictCreate(&keylistDictType,NULL);
          server.db[j].ready_keys = dictCreate(&objectKeyPointerValueDictType,NULL);
          server.db[j].watched_keys = dictCreate(&keylistDictType,NULL);
          listSetFreeMethod(server.db[j].defrag_later,(void (*)(void*))sdsfree);
      }
      //初始化lru键值池
      evictionPoolAlloc();
  
      //创建处理后台timer事件
      aeCreateTimeEvent(server.el, 1, serverCron, NULL, NULL)
  
      //将已创建好的tcp server socket绑定到event_loop实例的可读事件上，一旦有新的连接会调用回调函数acceptTcpHandler
      for (j = 0; j < server.ipfd_count; j++) {
          aeCreateFileEvent(server.el, server.ipfd[j], AE_READABLE,acceptTcpHandler,NULL)
      }
  
      //设置event loop被触发前后的处理函数
      aeSetBeforeSleepProc(server.el,beforeSleep);
      aeSetAfterSleepProc(server.el,afterSleep);
  }
  ```

  

那server tcp socket的监听函数如何被触发呢，需要了解事件的创建和eventloop事件循环：

* 首先`aeCreateFileEvent`

  ```c
  int aeCreateFileEvent(aeEventLoop *eventLoop, int fd, int mask,
          aeFileProc *proc, void *clientData)
  {
      if (fd >= eventLoop->setsize) {
          errno = ERANGE;
          return AE_ERR;
      }
      aeFileEvent *fe = &eventLoop->events[fd];
  
      if (aeApiAddEvent(eventLoop, fd, mask) == -1)
          return AE_ERR;
      //根据标志位设置读写回调
      fe->mask |= mask;
      if (mask & AE_READABLE) fe->rfileProc = proc;
      if (mask & AE_WRITABLE) fe->wfileProc = proc;
      fe->clientData = clientData;
      if (fd > eventLoop->maxfd)
          eventLoop->maxfd = fd;
      return AE_OK;
  }
  ```

* `aeMain`如何处理对应客户端的读事件的：

  ```c
  void aeMain(aeEventLoop *eventLoop) {
      eventLoop->stop = 0;
      //事件循环
      while (!eventLoop->stop) {
          aeProcessEvents(eventLoop, AE_ALL_EVENTS|
                                     AE_CALL_BEFORE_SLEEP|
                                     AE_CALL_AFTER_SLEEP);
      }
  }
  //处理读写、timer事件
  int aeProcessEvents(aeEventLoop *eventLoop, int flags)
  {
      int processed = 0, numevents;
  
      if (eventLoop->maxfd != -1 ||
          ((flags & AE_TIME_EVENTS) && !(flags & AE_DONT_WAIT))) {
          struct timeval tv, *tvp;
          
          if (eventLoop->beforesleep != NULL && flags & AE_CALL_BEFORE_SLEEP)
              eventLoop->beforesleep(eventLoop);
  		
          //这里调用封装的多路复用的API（linux下是epoll），如果超时或者有事件则返回
          numevents = aeApiPoll(eventLoop, tvp);
  
          if (eventLoop->aftersleep != NULL && flags & AE_CALL_AFTER_SLEEP)
              eventLoop->aftersleep(eventLoop);
  		
          //处理对应fd上的IO事件
          for (j = 0; j < numevents; j++) {
              aeFileEvent *fe = &eventLoop->events[eventLoop->fired[j].fd];
              int mask = eventLoop->fired[j].mask;
              int fd = eventLoop->fired[j].fd;
              int fired = 0;
  			
              //处理可读事件
              if (fe->mask & mask & AE_READABLE) {
                  fe->rfileProc(eventLoop,fd,fe->clientData,mask);
                  fired++;
                  fe = &eventLoop->events[fd]; /* Refresh in case of resize. */
              }
  
              //处理可写事件
              if (fe->mask & mask & AE_WRITABLE) {
                  if (!fired || fe->wfileProc != fe->rfileProc) {
                      fe->wfileProc(eventLoop,fd,fe->clientData,mask);
                      fired++;
                  }
              }
  
              processed++;
          }
      }
  
      return processed;
  }
  ```

  

因此一旦事件循环监测到server socket上有可读事件就会调用`acceptTcpHandler`，具体实现如下：

```c
void acceptTcpHandler(aeEventLoop *el, int fd, void *privdata, int mask) {
    int cport, cfd, max = MAX_ACCEPTS_PER_CALL;
    char cip[NET_IP_STR_LEN];
	//可能一次触发有多个连接，所以这里使用循环;如果因为server socke设置的非阻塞，一旦accept失败马上返回失败则跳出循环
    while(max--) {
        //新连接为cfd
        cfd = anetTcpAccept(server.neterr, fd, cip, sizeof(cip), &cport);
        if (cfd == ANET_ERR) {
            if (errno != EWOULDBLOCK)
                serverLog(LL_WARNING,
                    "Accepting client connection: %s", server.neterr);
            return;
        }
        serverLog(LL_VERBOSE,"Accepted %s:%d", cip, cport);
        //连接后续处理
        //connCreateAcceptedSocket创建一个client的connection连接
        acceptCommonHandler(connCreateAcceptedSocket(cfd),0,cip);
    }
}
```

`acceptCommonHandler`如何使得客户端连接处理读写事件的，真正处理的函数在下面：

```c
//根据连接信息创建client
client *createClient(connection *conn) {
    if (NULL == conn) return NULL;
    client *c = zmalloc(sizeof(client));
    selectDb(c,0);
    uint64_t client_id = ++server.next_client_id;
    c->id = client_id;
    c->resp = 2;
    c->conn = conn;
    c->name = NULL;
    connNonBlock(conn);
    connEnableTcpNoDelay(conn);
    //设置可读回调
    connSetReadHandler(conn, readQueryFromClient);
    connSetPrivateData(conn, c);
	//将client放在全局clients链表中
    linkClient(c);
    return c;
}
```

一旦有来自客户端的命令，就会触发`readQueryFromClient`,接收client数据并解析：

```c
void readQueryFromClient(connection *conn) {
    client *c = connGetPrivateData(conn);
    int nread, readlen;
    size_t qblen;

    readlen = PROTO_IOBUF_LEN;
    if (c->reqtype == PROTO_REQ_MULTIBULK && c->multibulklen && c->bulklen != -1
        && c->bulklen >= PROTO_MBULK_BIG_ARG)
    {
        ssize_t remaining = (size_t)(c->bulklen+2)-sdslen(c->querybuf);
        if (remaining > 0 && remaining < readlen) readlen = remaining;
    }

    qblen = sdslen(c->querybuf);
    if (c->querybuf_peak < qblen) c->querybuf_peak = qblen;
    c->querybuf = sdsMakeRoomFor(c->querybuf, readlen);
    //调用读函数，从socket读取数据
    nread = connRead(c->conn, c->querybuf+qblen, readlen);
    if (nread == -1) {
        if (connGetState(conn) == CONN_STATE_CONNECTED) {
            return;
        } else {
            serverLog(LL_VERBOSE, "Reading from client: %s",connGetLastError(c->conn));
            freeClientAsync(c);
            return;
        }
    } else if (nread == 0) {
        serverLog(LL_VERBOSE, "Client closed connection");
        freeClientAsync(c);
        return;
    } else if (c->flags & CLIENT_MASTER) {
        c->pending_querybuf = sdscatlen(c->pending_querybuf,
                                        c->querybuf+qblen,nread);
    }

    sdsIncrLen(c->querybuf,nread);
    //准备好raw数据来解析并处理
     processInputBuffer(c);
}
```

```c
void processInputBuffer(client *c) {
    while(c->qb_pos < sdslen(c->querybuf)) {
        if (!c->reqtype) {
            if (c->querybuf[c->qb_pos] == '*') {
                c->reqtype = PROTO_REQ_MULTIBULK;
            } else {
                c->reqtype = PROTO_REQ_INLINE;
            }
        }
		//普通命令以*开头，请求类型为 PROTO_REQ_MULTIBULK；管道命令类型为PROTO_REQ_INLINE
        if (c->reqtype == PROTO_REQ_INLINE) {
            if (processInlineBuffer(c) != C_OK) break;
            if (server.gopher_enabled &&
                ((c->argc == 1 && ((char*)(c->argv[0]->ptr))[0] == '/') ||
                  c->argc == 0))
            {
                processGopherRequest(c);
                resetClient(c);
                c->flags |= CLIENT_CLOSE_AFTER_REPLY;
                break;
            }
        } else if (c->reqtype == PROTO_REQ_MULTIBULK) {
            //解析redis消息协议
            if (processMultibulkBuffer(c) != C_OK) break;
        } else {
            serverPanic("Unknown request type");
        }

        if (c->argc == 0) {
            resetClient(c);
        } else {
            if (c->flags & CLIENT_PENDING_READ) {
                c->flags |= CLIENT_PENDING_COMMAND;
                break;
            }

            //这里开始执行命令
            if (processCommandAndResetClient(c) == C_ERR) {
                return;
            }
        }
    }

    if (c->qb_pos) {
        sdsrange(c->querybuf,c->qb_pos,-1);
        c->qb_pos = 0;
    }
}
```

正常情况下处理客户端命令处理比较简单，这里只列出关键步骤：

```c
//根据解析出的命令字符串查到对应的cmd
c->cmd = c->lastcmd = lookupCommand(c->argv[0]->ptr);
//检测当前命令参数合法性
if ((c->cmd->arity > 0 && c->cmd->arity != c->argc) ||
               (c->argc < -c->cmd->arity)) {
        rejectCommandFormat(c,"wrong number of arguments for '%s' command",
            c->cmd->name);
        return C_OK;
}
//这里真正执行了命令并返回消息给客户端
c->cmd->proc(c);
```

比如"hset name allen 1"就会查到`hsetCommand`。至此，我们了解了redis-server是如何处理客户端命令的基本流程。
