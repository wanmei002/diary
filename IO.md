# 大话 阻塞IO | 非阻塞IO | 同步IO | 异步IO
> 以下的理解接来源与对网络IO的理解

介绍以上的概念之前，我们先来简单说下一次网络的请求的流程:
sock()->bind()->listen()->accept()->recv()->send()->off

## 阻塞IO
有好几个客户端与服务端建成了连接如 fd1 fd2 fd3....., 此时用户态调用 fd1的recv() 问内核态(kernel) fd1 的数据准备好了吗？
此时阻塞，一直等内核态(kernel)说准备好了，然后把数据发送给用户态，开始问fd2 (大话西游的话: 你妈好吗).....
> 此时最早的优化模型是过来一个连接开启一个线程处理这个连接

## 非阻塞IO
有好几个客户端与服务端建成了连接如  fd1 fd2 fd3....., 此时用户态调用 fd1 的recv() 问内核态(kernel) fd1发送过来数据了吗？
内核态(kernel) 说没准备好; 用户态在调用 fd2的 recv() 数据准备好了吗 ......; 直到有个 fd正在接受数据，开始阻塞，直到数据接收完成 返还给用户态

> 阻塞IO 与 非阻塞IO 的区别体现在问内核数据准备的怎么样，阻塞IO会阻塞直到内核态说数据准备好了; 而非阻塞IO