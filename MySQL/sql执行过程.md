### sql 执行过程

[sql 执行图解](sql_exec.png)

 1. 客户端发送一条查询给服务器
 2. 服务器先检查查询缓存, 如果命中了缓存，则立刻返回存储在缓存中的结果。否则进入下一个阶段
    + 如果查询缓存时打开的, 那么MySQL会优先检查这个查询是否命中查询缓存中的数据。这个检查是通过一个对大小写敏感的哈希查找实现的
    + 查询和缓存中的查询即使只有一个字节不同, 那也不会匹配缓存结果，如果命中缓存, 
 3. 服务器端进行sql解析、预处理，再由优化器生成对应的执行计划
    + 语法解析: mysql 通过关键字将sql语句解析, 并生成一棵对应的`解析树`, 验证语法是否正确
    + 预处理: 进一步检查解析树是否合法, 例如: 检查数据表和数据列是否存在, 还会解析名称和别名，是否有歧义；验证权限
    + 查询优化器: 优化器的作用就是找到这其中最好的执行计划，以下简单的说下优化:
        - COUNT() MIN() MAX() : 索引是从小到大排序的，MAX() MIN() 可能就是取索引的第一个或最后一个; COUNT() mysim 表中存在一个变量存储表的行数
        - 大部分数据库IN() 是转化成 OR 条件的子句，MySQL中将IN()列表中的数据先进行排序，然后进行二分查找的方式来确定列表中的值是否满足条件，这是一个复杂度O(logN)
 4. mysql 根据优化器生成的执行计划, 调用存储引擎的API来执行查询
 5. 将结果返回给客户端
 
 > mysql 内部每秒能够扫描内存中上百万行数据, 相比之下, mysql 响应数据给客户端就慢得多了。
 删除旧的数据就是一个很好的例子,定期清除大量的数据时，如果用一个大的语句一次性完成的话，则可能
 需要一次锁住很多数据、占满整个事务日志、耗尽系统资源、阻塞很多小的但重要的查询。
 一次性删除一万行数据一般来说是一个比较高效而且对服务器影响也最小的做法。同时 如果每次删除数据后，
 都暂停一会儿再做下一次删除，这样也可以将服务器上原本一次性的压力分散到一个很长的时间段中, 就可以大大降低对
 服务器的影响，还可以大大减少删除时锁的持有
 
 
 - 分解关联查询
 `select * from tag join tag_post on tag_post.tag_id=tag.id join post on tag_post.post_id=post.id where tag.tag='mysql'`
 
 可以分解成下面这些查询来代替:
 ```sql
    select * from tag where tag='mysql';
    select * from tag_post where tag_id=1234;
    select * from post where post.id in (123,456,567,9098);
 ```
 这样分解会有以下优势:
    - 让缓存的效率更高。 对于MySQL的查询缓存来说，如果关联中的某个表发生了变化，那么就无法使用查询缓存了，而拆分后，如果某个表很少改变，那么基于该表的查询就可以重复利用查询缓存结果了
    - 将查询分解后，执行单个查询可以减少锁的竞争
    - 查询本身效率也可能会有所提升。这个例子中，使用 IN() 代替关联查询，可以让 MySQL按照 ID顺序进行查询，这可能会比随机的关联要更高效
    - 减少冗余记录的查询。在应用层做关联查询，意味着对于某条记录应用只需要查询一次，而在数据库中做关联查询，则可能需要重复的访问一部分数据。从这点看，这样的重构还可能会减少网络和内存的消耗
    
    
    
    
#### mysql 的关联查询
 - JOIN 
    + 一张表(关联表) JOIN 另一张表(被关联表), 则: 从关联表中查询出来一条数据然后再去被关联表中查询出这条数据关联的数据; 然后再跳到关联表中查询下一条数据 这样循环查询


 - 子查询 | UNION
    + MySQL在 from 子句中遇到子查询时，先执行子查询并将结果放到一个临时表中, 然后将这个临时表当作一个普通表对待(派生表); UNION 也是使用类似的临时表
        > 临时表时是没有任何索引的, 在编写复杂的子查询和关联查询的时候需要注意这一点
    
#### 排序
 - 无论如何排序都是一个成本很高的操作, 所以从性能角度考虑, 应尽可能避免排序或者尽可能避免对大量数据进行排序
 - 当不使用索引生成排序结果的时候, MySQL需要自己进行排序，如果数据量小与排序缓冲区则在内存中排序，数据量大则需要使用磁盘，MySQL将这个过程统一称为文件排序(filesort)
 - 在关联查询的时候如果需要排序, MySQL会分两种情况来处理这样的文件排序
    + order by 子句中的所有列都来自关联的第一个表, 那么MySQL在关联处理第一个表的时候就进行文件排序
    + 其他情况 MySQL都会先将关联的结果存放到一个临时表中, 然后再所有的关联都结束后，再进行文件排序
 - 如果查询中有 limit 的话，limit 也会在排序之后应用，所以即使需要返回较少的数据，临时表和需要排序的数据量仍然会非常大
```sql
 CREATE TABLE `tb1`( 
   `id` INT NOT NULL PRIMARY KEY AUTO_INCREMENT,  
    `age` INT(11) DEFAULT NULL,   
    `birth` timestamp not NULL 
 )ENGINE=INNODB DEFAULT CHARSET=utf8;
  
 insert into tblA(age,birth) values (12,now());
 insert into tblA(age,birth) values (13,now());
 insert into tblA(age,birth) values (14,now());
  
 create index age_birth on tb1(age,birth);
 select * from tb1;
```

##### 列举下会用上索引的, 只有 where 用上索引的最前列，order by 里也是索引字段的排序顺序
 - `select * from tb1 where age>13 order by age;` order by 会用上索引
 - `select * from tb1 where age>13 order by age,birth;` order by 会用上索引, 如果 age 相等的，在索引里  birth 也是按升序排序
 - `select * from tb1 where age>13 order by age desc,birth desc;` order by 会用上索引
 
