### RDB 快照
把当前内存中的数据集快照写入磁盘, 也就是 snapshot 快照(数据库中所有键值对数据). 恢复时是将快照文件直接读取到内存里

#### RDB 的触发方式
##### 自动触发
save m n : 表示 m 秒内数据集存在 n 次修改时，自动触发 bgsave
```json
save 900 1 // 表示 900 秒内至少有 1 个 key 的值变化, 则保存
save 300 10 // 表示 300 秒内 如果至少有 10 个 key 的值变化，则保存
save 60  10000  // 表示 60秒内如果至少有 10000 个 key 的值变化, 则保存
```
> 如果不需要持久化, 那么可以注释掉所有的 save 行来停用保存功能，或者 save ""

##### 手动触发
 1. save: 该命令会阻塞当前 redis 服务器, 执行 save 命令期间, redis 不能处理其它命令, 直到 RDB 过程完成为止
 2. bgsave: 执行该命令时, redis 会在后台异步进行快照操作, 快照同时还可以响应客户端请求。
 > save 命令对于内存比较大的实例, 会造成长时间阻塞, 这是致命的缺陷; 
 bgsave redis 进程执行 fork 操作创建子进程, RDB 持久化过程由子进程负责, 完成后自动结束，阻塞只发生在 fork 阶段, 一般很短
 
##### 停止 RDB 持久化
redis-cli config set save " "

##### RDB 的优势和劣势
 1. 优势
    1. 内容比较紧凑，保存了某个时间点上的数据集. 这种文件非常适合用于进行备份和灾难恢复
    2. 生成 RDB 文件的时候, redis 主进程 fork() 子进程来处理保存工作, 不影响主进程
    3. RDB 在回复大数据集时的速度比 AOF 的恢复速度要快
 2. 劣势
    1. RDB没办法做到实时持久化/秒级持久化. 因为 bgsave 每次运行都要执行 fork 操作创建子进程，属于
    重量级操作，如果不采用压缩算法，那么内存中的数据被克隆一份(相当于在数据保存到 RDB 文件里的时候，内容里有两份一样的数据),
    大致 2 倍的膨胀性需要考虑, 频繁执行成本过高( 影响性能)
    2. RDB 文件使用特定二进制格式保存, redis 版本演进过程中有多个格式的 RDB 版本, 存在不兼容问题
    3. 在一定间隔时间做一次备份, 所以如果redis意外 down 掉的话, 就会丢失最后一次快照后的所有修改(数据有丢失)
    
### RDB 自动保存的原理
redis 有个结构
```c
struct redisService{
    // 1. 记录保存 save 条件的数组
    struct saveparam *saveparams;
    // 2. 修改计数器, redis 每次修改dirty+1，如果执行了 bgsave, dirty=0
    long long dirty;
    // 3. 上一次执行保存的时间, 用于计算距离当前时间 
    time_t lastsave;
}
```
saveparam 结构
```c
struct saveparam{
    // 秒数
    time_t seconds;
    // 修改数
    int changes;
}
```
 redis 每 100 毫秒会执行一次 severCron 函数, 这个函数会遍历 saveparams 数组中保存的所有条件，只要有
 一个条件被满足, 那么就会执行 bgsave 命令. 执行完后 dirty 计数器更新为0, lastsave 也更新为执行命令的完成时间
 





