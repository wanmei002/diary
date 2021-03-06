### redis 五大数据类型实现原理

#### 对象的类型和编码
 - redis 底层的数据类型用来标识键和值, 每次在 redis 数据库中创建一个键值对时, 至少会创建两个对象, 一个是键对象, 一个是值对象, 
 而 redis 中的每个对象都是由 redisObject 结构来表示的:
    ```cgo
    typedef struct redisObject {
       unsigned type:4;//类型
       unsigned encoding:4;// 编码
       void *ptr; // 指向底层数据接口的指针
       int refcount;// 引用计数
       unsigned lru:22; // 记录最后一次被程序访问的时间
    }
    ```
    
    + `type` : `string` `list` `hash` `set` `zset`
    + encoding : int embstr raw(这三个是字符串数据结构)
                ht(字典) linkedlist(双端链表)
                ziplist(压缩列表)  intset(整数集合) skiplist(跳跃表和字典)
                
                
#### 字符串编码 string
 - int 编码: 保存的是可以用 long 类型表示的整数值; 用来保存整数值
 - raw 编码: 保存长度大于 44 字节的字符串; 用来保存长字符串
 - embstr 编码: 保存长度小于 44 字节的字符串; 用来保存短字符串
 - `raw` 和 `embstr` 的区别
    + embstr 和 raw 都是使用 redisObject+sds 保存数据, 区别在于, embstr 的使用只分配一次内存空间(因此 redisObject和sds是连续的)
    而 raw 需要分配两次内存空间(分别为 redisObject 和 sds 分配空间). 因此与 raw 相比, embstr 的好处在于创建时少分配一次空间, 删除时少释放一次空间，对象和数据连在一起, 寻找方便
    + `embstr` 坏处: 如果字符串的长度增加需要重新分配内存时, 整个 redisObject 和 sds 都需要重新分配空间, 因此 embstr 实现为只读
    + redis 中对于浮点数类型也是作为字符串保存的, 在需要的时候再将其转换成浮点数类型
 - 编码的转换
    + int 编码保存的值不再是整数, 或大小超过了 long 的范围时, 自动转化为 raw
    + 在对 embstr 对象进行修改时, 都会先转化为 raw 再进行修改, 因此只要修改 embstr 对象最后都会时 raw 对象 不管是否够 44字节
    
#### 列表对象 list 
 - list 列表, 它是简单的字符串列表, 按照插入顺序排序, 你可以添加一个元素到列表的头部或者尾部, 它的底层实际上是哥链表结构
 - 列表对象的编码可以是 ziplist(压缩列表)和linkedlist(双端链表)
 - 当满足下面两个条件时, 使用 ziplist(压缩列表) 编码:
    + 列表元素个数小于 512 个
    + 每个元素长度小于 64 字节
 - 不能满足这两个条件的时候使用 linkedlist 编码
 
#### HASH 列表 hash 
 - 哈希对象的键时一个字符串类型, 值是一个键值对集合
 - 编码 : 哈希对象的编码可以是 ziplist 或者 hashtable
 - 当使用 ziplist 也就是压缩列表作为底层实现时, 新增的键值对时保存到压缩列表的表尾
 - 同时满足一下两个条件时, 使用 ziplist(压缩列表)编码: 条件可以通过配置文件修改
    + 列表元素个数小于 512 个
    + 每个元素长度小于 64 字节

 - 不能满足这两个条件的时候使用 hashtable 编码
 
#### 集合对象 set 
 - 编码 : `intset` `hashtable`
 - intset 编码的集合对象使用整数集合作为底层实现, 集合对象包含的所有元素都被保存在整数集合中
 - hashtable 编码的集合对象使用 字典 作为底层实现, 字典的每个键都是一个字符串对象, 这里每个字符串对象就是一个集合中的元素, 而字典中的值全部为 null
 - 满足以下两个使用 intset 编码:
    + 集合对象中所有元素都是整数
    + 集合对象所有元素数量不超过 512个
 
 
#### 有序集合对象 zset
 - 有序集合为每一个元素设置一个分数(score) 作为排序依据
 - `编码` : 有序集合的编码可以是 ziplist 或者 skiplist 
 - ziplist 编码的有序集合对象使用也锁列表作为底层实现, 每个集合元素使用两个紧挨在一起的压缩列表节点来保存，第一个节点保存元素的成员, 
 第二个节点保存元素的分值, 并且压缩列表内的集合元素按照分值从小到大的顺序进行排列, 小的放置在靠近表头的位置, 大的放置在靠近表尾的位置




