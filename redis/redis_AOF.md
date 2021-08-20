## redis AOF持久化
AOF 是通过保存 redis 服务器所执行的写命令来记录数据库状态

### 开启 AOF 持久化
```ini
# appendonly 参数在 APPEND ONLY MODE 下
# appendonly 默认是 no, 不开启 AOF 持久化
appendonly yes
# aof 文件名, 默认是 "appendonly.aof"
appendfilename "appendonly.aof"

# appendfsync: AOF持久化策略的配置
appendfsync no|always|everysec
```
#### appendfsync 参数讲解
 - no: 表示不执行 fsync, 由操作系统保证数据同步到磁盘, 速度快, 但是不太安全;
 - always: 表示每次写入都执行 fsync, 以保证数据同步到磁盘, 效率很低;
 - everysec: 表示每秒执行一次 fsync, 可能会导致丢失这 1s 数据。通常选择 everysec, 兼顾安全性和效率。

#### no-appendfsync-on-rewrite
在 AOF 重写或者写入 RDB 文件的时候, 会执行大量IO, 此时对于 AOF:appendfsync:everysec|always 的模式来说,
执行 fsync 会造成阻塞过长时间(大量的IO操作), 值为 `yes` 表示 rewrite 期间对新写操作不 fsync, 暂时存在内存中，
等 rewrite 完成后再写入, 默认值为 `no`, 建议 `yes`

#### auto-aof-rewrite-percentage
什么时候重写AOF文件, 默认值为100。当前AOF文件大小是上次日志重写得到AOF文件大小的二倍（设置为100）时，自动启动新的日志重写过程。

#### auto-aof-rewrite-min-size
64mb。设置允许重写的最小aof文件大小，避免了达到约定百分比但尺寸仍然很小的情况还要重写。

#### aof-load-truncated
aof文件可能在尾部是不完整的，

### AOF 文件恢复
重启 redis 之后就会进行 AOF 文件的载入。

异常修复命令: redis-check-aof --fix

### AOF 重写
 AOF是不断将写命令记录到 AOF 文件中, 随着 redis 不断的进行, AOF 会越来越大, 占用服务器资源越多以及 AOF 
 恢复要求时间越长. 为了解决这个问题, redis 新增了重写机制, 当 AOF 文件的大小超过所设定的阙值时, redis 就会
 启动 AOF 文件压缩，合并命令, 只保存可以恢复数据的最小指令集, 可以使用命令 `bgrewriteaof` 来重写
 
 也就是说 AOF 文件重写并不是对原文件进行重新整理，而是直接读取服务器现有的键值对，
 然后用一条命令去代替之前记录这个键值对的多条命令，生成一个新的文件后去替换原来的 AOF 文件。
 




