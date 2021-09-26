### Elasticsearch
 - 解决了分布式(搜索、索引)，高性能(近实时)，高可用(海量数据，堆机器)的检索服务，可以服务数百台机器处理大数据
 
#### 一次简单的搜索流程
 - 索引-> 分析 -> 检索
    + 索引 : 将收集到的数据建立倒排索引并存储
    + 分析 : 将用户输入的 keyword 分解为索引服务可识别的词缀
    + 检索 : 将对应 token 与索引库中的倒排索引进行对比，并返回检索结果。
    
    
#### es 数据分布式
 - 数据的存储访问在多个节点(node)的多个分片(shard)上。es 自动将分片数据分布在es集群不同的节点上，以实现数据的分布(同分片的主从分片不能在同节点上，防止节点故障数据丢失)
 
 
#### es基本概念
 - Inverted Index，相当于MySQL中的database， document集合，index数据可以分为多个 shard 分片
 - node ,一个 es实例就是一个node, 1个node不等同于1台服务器。一个node下拥有多个 index
 - shard
    + 相当于一个索引的数据用多个杯子(shard)来装, 这样可以将数据分布在多个node上，完成数据的分布式部署，每个分片都是一个 lucene index(全文搜索引擎)的实例
    + primary shard : 相当于主库，可读可写
    + replica shard : 数据备份，相当于从库读，提高读吞吐量。当主分片失效时，同坐自动的election, replica可以成为 primary
    
 - index-mysql_database; type-mysql_table; document-mysql_row; field-mysql_field
 
 
#### 节点的分布与一致性
 - 数据通过切片分布在多个节点上，在 es 中每个节点都可以完成用户的检索请求，安装不同的服务功能划分为3类职责：master/data/client
    + masterNode 创建/删除/分配切片
    + dataNode 数据查询/数据写入
    + clientNode 纯转发
    + clientNode 可以增加集群的访问吞吐量，但是 masterNode 掌握了集群中其它节点的数据元信息(及分片的分布情况)
    
 - es 通过 discovery 模块来完成新加入节点的发现和已有的节点宕机发现，新节点加入前已通过配置确认新节点的身份。宕机节点通过集群中其它节点的 unicast 机制来发现节点宕机情况，并自动剔除集权
    同时 es 集群通过自动选举解决了集群首脑masterNode 宕机时服务不可用的问题
    
    
    