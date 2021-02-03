# elasticSearch query
## bank 数据结构
```json
{
    "account_number": 0,
    "balance": 16623,
    "firstname": "Bradshaw",
    "lastname": "Mckenzie",
    "age": 29,
    "gender": "F",
    "address": "244 Columbus Place",
    "employer": "Euron",
    "email": "bradshawmckenzie@euron.com",
    "city": "Hobucken",
    "state": "CO"
}
```

### 解释下 elasticSearch 返回的字段说明
```json
{
  "took" : 63,
  "timed_out" : false,
  "_shards" : {
    "total" : 5,
    "successful" : 5,
    "skipped" : 0,
    "failed" : 0
  },
  "hits" : {
    "total" : {
        "value": 1,
        "relation": "eq"
    },
    "max_score" : null,
    "hits" : [ {
      "_index" : "bank",
      "_type" : "_doc",
      "_id" : "0",
      "sort": [0],
      "_score" : null,
      "_source" : {"account_number":0,"balance":16623,"firstname":"Bradshaw","lastname":"Mckenzie","age":29,"gender":"F","address":"244 Columbus Place","employer":"Euron","email":"bradshawmckenzie@euron.com","city":"Hobucken","state":"CO"}
    }]
  }
}
```

 - took elasticSearch 运行查询多长时间(毫秒为单位)
 - timed_out 搜索请求是否超时
 - _shards 搜索了多少个分片以及成功、失败和跳过了多少个分片
 - max_score 找到最相关文件的分数
 - hits.total.value 找到了多少个匹配的文档
 - hits.sort 文档排序位置 (不按相关性得分排序时)
 - hits._score 文档的相关性得分 (使用时, 不适用match_all)

### sort 排序
```json
POST /zzh_aa/_search
{
  "sort":[
    {"age":{"order":"asc"}},
    {"type":{"order":"desc"}}
  ]
}

```

### 空格分词匹配
```bash
curl -d '{"query":{"query_string":{"query":"brith:zyn world"}}}' -H'Content-Type: application/json; charset=UTF-8' -X POST http://127.0.0.1:9200/zzh_aa/_search
```
```json
POST /zzh_aa/_search
{
  "query":{
    "query_string":{
      "query":"brith:zyn world"
    }
  }
}
```
> brith 是字段名

### 空格不分词匹配
#### match_phrase 会把 `hello zzh` 整个匹配不会把他们分开匹配
```json
POST /bank/_search
{
  "query": { "match_phrase": { "brith": "hello zzh" } }
}
```

### 必须匹配和排除匹配
```json
POST /bank/_search
{
  "query": {
    "bool": {
      "must": [
        { "match": { "age": "40" } }
      ],
      "must_not": [
        { "match": { "state": "ID" } }
      ]
    }
  }
}
```

