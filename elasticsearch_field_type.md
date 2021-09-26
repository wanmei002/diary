### keyword family
keyword 通常用于 排序 汇总 和 术语级查询, 如果要用于全文搜索 请用 `text` 类型
#### keyword
keyword 可以保存像 id email 域名 状态码 标签等类型的数据

```json
PUT my_index
{
  "mappings": {
    "properties": {
      "tag": {
        "type": "keyword"
      }
    }
  }
}

```

numeric 类型主要用于 range 查询, 但是有些 numeric 不 range 查询，可以转换成 keyword, 下面这些情况可以尝试转成 keyword
 - 不计划用 range 查询的
 - term 在 keyword 上查询比用 range 查询快
 
#### constant_keyword
constant_keyword 保存值一定的数据
```json
PUT my_index
{
  "mappings":{
    "properties":{
      "@timestamp":{
        "type":"date"
      },
      "message":{
        "type":"keywrod"
      },
      "level":{
        "type":"constant_keyword",
        "value":"debug"
      }
    }
  }
}
```
level 不传的话也会默认为 debug 
```json
POST my_index
{
  "date": "2019-12-12",
  "message": "Starting up Elasticsearch",
  "level": "debug"
}

POST my_idnex
{
  "date": "2019-12-12",
  "message": "Starting up Elasticearch"
}

```

#### wildcard
wildcard 可以用于 grep 正则查询
```json
PUT my_index
{
  "mappings": {
    "properties": {
      "my_wildcard": {
        "type": "wildcard"
      }
    }
  }
}

PUT my_index/_doc/1
{
  "my_wildcard": "this string can be quite lengthy"
}

// search 查询
GET my_index
{
  "query": {
    "wildcard": {
      "my_wildcard": {
        "value": "*quite*lengthy"
      }
    }
  }
}
```

### Nested 
 `nested` 是一个特殊的对象数组

```json
PUT my_index
{
  "mappings": {
    "properties": {
      "user": {
        "type": "nested"
      }
    }
  }
}
```

### object
 json 文档本质是分层级的； json文档可以包含object
```json
// PUT my_index
{
  "region": "china",
  "manager": {
    "age": 30,
    "name": {
      "first": "zzh",
      "last" : "zyn"
    }
  }
}
```
 - 最外层永远是一个 `json` 对象
 - 最外层的 json 对象包含了一个对象叫 `manager`
 - `manager` 对象里面又包含了一个对象叫 `name`

#### 我们看一下上面 document 的数据结构
```json
// curl http://127.0.0.1:9200/zzh_1/_mappings?pretty
{
  "mappings": {
    "properties": {
      "region": {
        "type": "keyword"
      },
      "manager": {
        "properties": {
          "age": {"type": "integer"},
          "name": {
            "properties": {
              "first": {"type": "keyword"},
              "last" : {"type": "keyword"}
            }
          }
        }
      }
    }
  }
}
```

### rank feature 字段
 type 是 `rank feature`  数据一般是 `object`
#### mappings
```json
{
  "mappings": {
    "properties": {
      "pagerank": {
        "type": "rank_feature"
      },
      "url_len": {
        "type": "rank_feature",
        "positive_score_impact": false
      }
    }
  }
}
```

#### add data
```json
// PUT http://127.0.0.1:9200/zzh_1/_doc/1
{
  "topics": {
    "politics": 20,
    "economics": 50.8
  }
}
```
#### search  数据 会根据 `_score` 大小从大到小排序
```json
curl -X GET http://127.0.0.1:9200/zzh_1/_search?pretty -H'Content-Type: application/json' -d'
{
  "query": {
    "rank_feature": {
      "field": "pagerank"
    }
  }
}
'
```

### text
 text 适合保存像 邮件内容或产品描述这些大文本的，文本保存， 在索引之前它们通过 分析器 转化成列表


### unsigned_long 无符号整型
```json
// PUT my_index
{
  "mappings": {
    "properties": {
      "my_counter": {
        "type": "unsigned_long"
      }
    }
  }
}

// GET my_index/_search?pretty
{
    "query": {
        "term" : {
            "my_counter" : 18446744073709551615
        }
    }
}

{
    "query": {
        "range" : {
            "my_counter" : {
                "gte" : "9223372036854775808.5",
                "lte" : "18446744073709551615"
            }
        }
    }
}
```


 