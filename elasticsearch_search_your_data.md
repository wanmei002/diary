# search your data
##  collapse search results
```json
GET /index-1/_search
{
  "query": {
    "match": {
      "message": "GET /search"
    }
  },
  "collapse": {// 折叠属性 用 user.id 折叠
    "field": "user.id"
  },
  "sort": ["http.response.bytes"],
  "from": 10
}
```
> the field used for collapsing must be a single valued `keyword` or `numeric` field with
`doc_values` activated

```json
GET /my-index-000001/_search
{
  "query": {
    "match": {
      "message": "GET /search"
    }
  },
  "collapse": {
    "field": "user.id",                       
    "inner_hits": {
      "name": "most_recent",                  
      "size": 5,                              
      "sort": [ { "@timestamp": "asc" } ]     
    },
    "max_concurrent_group_searches": 4        
  },
  "sort": [ "http.response.bytes" ]
}
```

## filter search results
