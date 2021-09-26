# search DSL
## compound queries
### bool query
布尔查询是使用一个或多个布尔子句构建的，每个子句都具有类型的出现。

类型有以下这些:
 - `must` 必须出现在匹配的文档中，并将有助于得分。
 - `filter` 子句(查询) 必须出现在匹配的文档中。 filter 子句在filter上下文中执行，这意味着计分被忽略，
    并且子句被考虑用于缓存
 - `should` 子句(查询)应出现在匹配的文档中
 - `must_not` 子句(查询)不得出现在匹配的文档中。子句在过滤器(filter)上下文中执行，这意味着计分被忽略，并且
 子句被视为用于缓存。
 ```json
curl -X POST "localhost:9200/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "bool" : {
      "must" : {
        "term" : { "user.id" : "kimchy" }
      },
      "filter": {
        "term" : { "tags" : "production" }
      },
      "must_not" : {
        "range" : {
          "age" : { "gte" : 10, "lte" : 20 }
        }
      },
      "should" : [
        { "term" : { "tags" : "env1" } },
        { "term" : { "tags" : "deployed" } }
      ],
      "minimum_should_match" : 1,
      "boost" : 1.0
    }
  }
}
'
```
#### minimum_should_match
minimum_should_match 参数指定 should 返回的文档必须匹配的子句的数量或百分比。
> =1 必须使用一个 should 条件，不能不使用

如果 bool 查询包含至少一个 should 子句，而没有 must 或 filter 子句，则默认值为1. 否则默认为0

### Boosting query(助推查询)
```json
curl -X GET "localhost:9200/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "boosting": {
      "positive": {
        "term": {
          "text": "apple"
        }
      },
      "negative": {
        "term": {
          "text": "pie tart fruit crumble tree"
        }
      },
      "negative_boost": 0.5
    }
  }
}
'
```
 - `positive` : 返回的文档必须包含此匹配
 - `negative` : 如果 negative 查询出的结果在positive查询出的结果里出现，将计算文档的最终相关性得分:
    + 从 positive 查询中获取原始的相关性分数。
    + 将分数乘以 negative_boost 值。
    
 - `negative_boost`: *必需项* ，值介于 0~1
 
 
## full text query (全文查询)
匹配查询返回与提供的文本，数字，日期或布尔值匹配的文档。匹配之前对提供的文本进行分析。
```json
curl -X GET "localhost:9200/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "match": {
      "message": {
        "query": "this is a test"
      }
    }
  }
}
'
```
### match的参数<field>
(必须，对象) 你希望查询的字段

### <field>的参数
 - `query` (必须), 可以是 文本 数字 布尔或者日期等你希望在文档中找到的值
 - `analyzer` (可选，字符串) 分析器
 - `auto_generate_synonyms_phrase_query` (可选，布尔值)，如果为 true,则会自动为多个术语同义词创建匹配词组查询，默认为true.
 - `fuzziness` (可选，字符串) 匹配允许的最大编辑据力
 - `max_expansions` (可选，整数)查询将扩展到的最大术语数。默认`50`
 - `prefix_length` (可选，整数) 为模糊匹配保留的起始字符数，默认为0
 - `fuzzy_transpositions` (可选，布尔值) 如果为true, 则模糊匹配的编辑内容包括两个相邻字符的变化(ab->ba), 默认为 true
 - `fuzzy_rewrite` (可选，字符串)用于重写查询的方法。
 - `lenient` (可选，布尔值) 如果为 true,则将忽略基于格式的错误，例如为数字字段提供文本 query 值。默认为 false
 - `operator` (可选，布尔值) 布尔逻辑，用于解释 query 值中的文本。有效值为: 
    `OR`(默认) : 例如 query 值 capital of Hungary 解释为 capital OR of OR Hungary
    `AND` : 例如 query 值 capital of Hungary 解释为 capital AND of AND Hungary

 - `minimum_should_match` : (可选，字符串) 指示如果 analyzer 删除所有标记(例如使用 stop 过滤器时)，是否不返回任何文档。
 有效值为: 
    `none`(默认) : 如果 analyzer 删除所有标记，则不会返回任何文档。
    `all` : 返回所有文档，类似于 match_all 查询

### match_bool_prefix
```shell
curl -X GET "localhost:9200/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "match_bool_prefix" : {
      "message" : "quick brown f"
    }
  }
}
'
# 上面的可以解析成下面的
curl -X GET "localhost:9200/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "bool" : {
      "should": [
        { "term": { "message": "quick" }},
        { "term": { "message": "brown" }},
        { "prefix": { "message": "f"}}
      ]
    }
  }
}
'
```
> 匹配 text文档中有 quick 或 brown 或 前缀是 f 的文档

### match phrase query 
```json
curl -X GET "localhost:9200/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "match_phrase": {
      "message": "this is a test"
    }
  }
}
'
```

> 匹配this is a test ，把它当成一个整体

### match phrase prefix query
```json
curl -X GET "localhost:9200/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "match_phrase_prefix": {
      "message": {
        "query": "quick brown f"
      }
    }
  }
}
'
```

> 把 quick brown f 当成一个整体，匹配文档内容有前缀是  `quick brown f` 的文档


### multi-match query 
```json
curl -X GET "localhost:9200/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "multi_match" : {
      "query":    "this is a test",// 要查询的数据
      "fields": [ "subject", "message" ] // 要查询的字段
    }
  }
}
'
```

```json
// 查询字段通配符
curl -X GET "localhost:9200/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "multi_match" : {
      "query":    "Will Smith",
      "fields": [ "title", "*_name" ] // 匹配 title 和 以 _name 结尾的字段
    }
  }
}
'
```

```json
// 匹配
curl -X GET "localhost:9200/_search?pretty" -H'Content-Type: application/json' -d'
{
  "query": {
    "multi_match" : {
      "query" : "this is a test",
      "fields" : [ "subject^3", "message" ] //将 subject 字段的分数*3 
    }
  }
}
'
```

#### multi_match 查询类型: 
multi_match 内部执行查询的方式取决于 type 参数，可以将其设置为: 
 - `best_fields`  : (默认), 查找与任何字段匹配但使用 _score 最佳字段中的文档
 - `most_fields`  : 查找与任何字段匹配的文档，并将每个字段中 _score 的合并
 - `cross_fields` : 在任何字段中查找每个单词
 - `phrase`       : `match_phrase` 在每个字段上运行查询，并使用 `_score` 的最佳字段
 - `phrase_prefix`: `match_phrase_prefix` 在每个字段上运行查询，并使用 _score 最佳的字段
 - `bool_prefix`  : `match_bool_prefix` 在每个字段上创建查询，并将每个字段中的_score 合并







