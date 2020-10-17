#### 专业属于
 - 脏页 : 当内存数据页和磁盘数据页上的内容不一致时, 我们称这个内存页为脏页
 
#### InnoDB存储引擎内存
##### 缓冲池(buffer pool)   重做日志缓冲池(redo log buffer)   额外的内存池(additional memory pool)

##### 缓冲池(buffer pool)  innodb_buffer_pool_size 指定了缓冲池的大小
 - 缓冲池是占最大块内存的部分, 用来存放各种数据的缓存. 
 - InnoDB存储引擎的工作方式总是将数据库文件按页(每页16K)读取到缓冲池, 然后按最少使用算法(LRU)来保留在缓冲池中的缓存数据
    + LRU算法: 
        - mysql 中数据页存储用链表的方式存储, 每个数据节点称为 block, 链表又分为2部分(young old)
        - 当MySQL没有在缓冲池中找到数据页时, 读取磁盘数据, 然后把读取的 block 插入到 old 表头, 如果 old 已满, 则删除表尾的数据
        - 当 MySQL 读取的页(block) 在 old 区域, 则拿出来放到 young 表头, old 前面的数据依次后移
        - 如果读取的页(block) 在 young 区域, 则拿出来放到 young 表头
        > 被访问的放到最前面, 没有被访问的逐渐后移, 直到被删除, 当表满时, 表尾的数据会被删除
        
 - 如果数据库文件需要修改, 则优先修改缓冲池中的页, 此时缓冲池中的页称为脏页; 然后按照一定的频率将缓冲池的脏页刷新(flush)到文件
 - 缓冲池中主要有索引页和数据页, 但是也有 undo页 插入缓冲(insert buffer) 自适应哈希索引(adaptive hash index)  InnoDB存储的锁信息(lock info)
    数据字典信息(data dictionary)等.  
    
##### 重做日志缓冲池(redo log)  主要负责 事务日志  innodb_flush_log_at_trx_commit 用来控制重做日志刷新到磁盘的策略
 - 当事务提交 commit 时, 必须先将事务的所有日志写入到重做日志文件进行持久化
 - 重做日志由两部分组成 redo log(保证事务的持久性) undo log(帮助事务回滚及MVCC的功能)
 - innodb_flush_log_at_trx_commit 
    +默认为1: 表示事务提交时必须调用一次系统的 fsync操作(强制刷新到重做日志文件中)
     fsync 效率取决磁盘的性能(磁盘的性能决定事务提交的性能, 间接的决定MySQL的新能)
    + 0 : 表示事务提交时不进行写入重做日志操作, 这个操作仅在 `master thread` 中完成, master thread 每一秒进行一次 fsync 操作
    + 2 : 表示事务提交时将重做日志写入重做日志文件, 但仅写入文件系统的缓冲中, 不进行 fsync 操作
    
##### 额外的内存池(additional memory pool)
 - InnoDB实例会申请缓冲池的空间, 但是每个缓冲池中的帧缓冲(frame buffer) 还有对应的缓冲控制对象, 而这些对象记录了
 诸如 LRU、锁、等待等方面的信息, 而这个对象的内存需要从额外内存池中申请. 因此, 当你申请了很大的 InnoDB 缓冲池时, 这个值也应该相应增加
 
#### master thread
 - 刷新脏页到磁盘(可能)
 - 合并至多5个插入缓冲(总是)
 - 将日志缓冲刷新到磁盘(总是)
 - 删除无用的 undo 页(总是)
 - 刷新 100 个 或 10个脏页到磁盘(总是)
 - 产生一个检查点(总是)
 
 