### search you data
```json
// 开场白
{
    "query": {
        "match": {
            "user.id": "kimchy"
        }
    }
}
// 查询出来的数据在 hits.hits 属性里


// ------想指定查询超时时间 (timeout)
// GET http://127.0.0.1:9200/my_index/_search?pretty
{
    "timeout": "2s",
    "query": {
        "match": {
            "user.id": "kimchy"
        }
    }
}

// ----------控制 hits.total 是否显示
{
    "track_total_hits": false,// false-不显示  默认为true-显示
    "query": {
        "match": {
            "user.id": "elkbee"
        }
    }
}
  // track_total_hits 也可以设置成整数,  for instance 
  "track_total_hits": 100
  // return data 
{
  "_shards": ...
  "timed_out": false,
  "took": 30,
  "hits": {
    "max_score": 1.0,
    "total": {
      "value": 42,     // 总共有 42 条匹配的数据    
      "relation": "eq" // count is eq 42
    },
    "hits": ...
  }
}

// 也可能是返回下面的数据
{
  "_shards": ...
  "hits": {
    "max_score": 1.0,
    "total": {
      "value": 100,    //      
      "relation": "gte"   // 说明总条数是 大于等于 100 的  
    },
    "hits": ...
  }
}
```