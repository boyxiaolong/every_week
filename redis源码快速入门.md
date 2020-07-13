
这是一篇了解关于redis是什么的快速开始文档

--------------
Redis常被标记为数据结构服务器。这意味着Redis提供通过一系列命令访问共享数据，发送指令通过建立tcp连接和简单协议来构建*服务器-客户端*模式。因此不同的进程可以通过共享模式来查询和修改同样的数据结构。
Redis实现的数据结构有部分特殊的属性:
* 尽管数据结构经常在服务器内存中被使用和修改，Redis关心如何将数据存储磁盘。这意味着数据在Redis很快但不会遗失。

* 数构的实现对于内存效率压力很大，因此Redis中的数据结构相比更高层次的编程语言会尽量更节省内存。

* Redis提供很多能在一个数据库中自然见到的特性，比如备份，可调整级别的持久性，集群，高可用性。有种很棒的关点认为Redis是memcached的复杂版本，因为Redis不止是SET和GET，如果你想了解的更多，下面有一系列网站可以尝试。

  

* 介绍Redis数据类型  http://redis.io/topics/data-types-intro

* 在你的浏览器试试Redis. http://try.redis.io

* Redis所有指令列表. http://redis.io/commands

* Redis官方文档. http://redis.io/documentation

Redis的安装可以选择源码编译或者linux下安装。

运行Redis
-------------
使用默认配置运行:
    % ./redis-server

也可以提供redic.conf作为额外参数
    % ./redis-server /path/to/redis.conf

也可以提供更多的参数，比如

    % ./redis-server --port 9999 --replicaof 127.0.0.1 6379
    % ./redis-server /etc/redis/6379.conf --loglevel debug

 试试和Redis交互
------------------
使用redis-cli可以连接redis-server实例，在另一终端输入：
    % ./redis-cli
    redis> ping
    PONG
    redis> set foo bar
    OK
    redis> get foo
    "bar"
    redis> incr mycounter
    (integer) 1
    redis> incr mycounter
    (integer) 2
    redis>

深入Redis
===
在此我们解释Redis源码分布，也就是每个文件的大致想法和Redis服务器内部最重要的函数和数据结构。我们始终站在高处而不是深入细节去讨论，不然文档将会很大，并且源码在持续修改。大部分代码会有很多注释也很好理解。



源码分布
---
在根目录有以下重要的目录：

* `src`: contains the Redis implementation, written in C.包含Redis的C代码实现



理解程序如何允许最简单的方法是理解它使用的数据结构。因此我们从重要的头文件 `server.h`开始

所有的服务器配置和所有共享的状态是通过类型是 `struct redisServer`的 `server`这个全局结构来定义的。

* `server.db`是Redis存储数据库的数组

* `server.commands`是命令表(hash table)

* `server.clients`是连接到服务器的client链表。

* `server.master`是一个特殊的client，如果当前实例是副本，那它就是指的master
这里有很多字段，大部分字段都在数据结构定义的地方直接注释了
另一个重要的Redis数据结构是client，它有很多字段，这里我们只展示重要的几个：

    struct client {
        int fd;
        sds querybuf;
        int argc;
        robj **argv;
        redisDb *db;
        int flags;
        list *reply;
        char buf[PROTO_REPLY_CHUNK_BYTES];
        ... many other fields ...
    }

结构定义了一个*已连接客户端*:

* `fd` 字段是socket文件标识符.

* `argc` `argv` 填充了客户端需要执行的指令，因此函数定义了Redis指令如何读取这些参数.

* `querybuf` 累计从客户端发来的请求，被Redis server根据Redis协议解析并通过指令的实现来执行.

* `reply` `buf` 是收集服务器发给客户端响应的动态和静态的缓存.一旦socket可写，这些缓存就会逐渐写入socket。
  在client中arguments 被描述为`robj`结构，具体细节如下：

      typedef struct redisObject {
          unsigned type:4;
          unsigned encoding:4;
          unsigned lru:LRU_BITS; /* lru time (relative to server.lruclock) */
          int refcount;
          void *ptr;
      } robj;
  一般这个数据结构代表了Redis所有的基础数据类型，比如strings, lists, sets, sorted sets等。有趣的是它含有一个 `type`类型，因此可以通过 `type`知道给定对象的类型，通过`refcount`一个同样的对象可以被多个地方引用而不用多次构造。最终`ptr`字段代表了实际数据，但最终要以`encoding`字段来决定。

  server.c

  
---
这是Redis Server的入口，定义了main函数。下面就是启动Redis服务器的关键步骤

* `initServerConfig()` 设置所有`server`结构的初始值.
* `initServer()`分配需要操作数据结构空间，设置监听socket.
* `aeMain()` 开始监听新连接的事件循环.
这里有两个特别的函数会在事件循环里周期性调用:
1. `serverCron()` 被周期性调用，执行与时间有关的任务比如检测客户端超市.
2. `beforeSleep()`当事件循环触发时被调用，处理一些请求，并返回事件循环.
在server.c 中还有其他函数来处理一些事情:
* `call()`用来在指定client中调用指定命令.

* `activeExpireCycle()` 处理通过`EXPIRE`指令设置的过期keys.

* `freeMemoryIfNeeded()`会被调用当一个新的写指令需要执行但redis的内存已经超过指定的最大内存.

The global variable `redisCommandTable` 全局变量 `redisCommandTable`定义所有的指令，包括指令的名字，函数指针，参数数量和其他属性.

  networking.c
---
这个文件定义所有和客户端，主从，集群相关的IO函数，:
* `createClient()` 为新的客户端分配内存并初始化.
* `addReply*`函数家族被用来实现在client数据结构追加数据.
* `writeToClient()` 写入client输出缓存，并在可写的时候调用`sendReplyToClient()`.
* `readQueryFromClient()`是一个可读函数，收集从客户端发送的查询缓存.
* `processInputBuffer()` 是一句Redis协议解析客户端缓存的入口点。一旦一个命令准备被处理会调用`processCommand()` 来实际执行一个指令.
* `freeClient()` 对客户端释放内存，断开并移除连接，.

db.c
---
* 某些指令处理一些特殊数据结构，其他是通用的。比如通用的指令`DEL`和`EXPIRE`。所有通用指令在`db.c`定义.
`db.c`实现了处理某些操作的API而不直接访问内部数据结构。
在`db.c`中最重要的函数是:

* `lookupKeyRead()` 和`lookupKeyWrite()`返回指定key的指针，如果不存在返回NULL`.
* `dbAdd()` 和 `setKey()` create a new key in a Redis database创建一个新key.
* `dbDelete()` 移除key和关联的值.
* `emptyDb()`移除整个数据库或者定义的所有数据库.



其他C 文件
---
* `t_hash.c`, `t_list.c`, `t_set.c`, `t_string.c`, `t_zset.c` `t_stream.c` 包含了Redis数据类型的实现。她们实现访问一个指定数据类型的API和客户端相关的指令.
* `ae.c`实现了Redis事件循环.
* `sds.c`实现了Redis string库.
* `anet.c` 是相比直接使用内核暴露的接口，是更加简单的POSIX网络接口库.
* `dict.c`实现了非阻塞（渐进式rehash）hash表.
* `scripting.c`实现lua脚本.
* `cluster.c`实现Redis集群，读完其他代码会更容易理解这部分.





祝你找到深入Redis代码的道路！

