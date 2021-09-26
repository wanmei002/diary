### 先搭建一个简单的服务器
 + 参阅 本目录下的 `simple_server.go`

## 需要引入 `"net/http"` 包
### 获取 `get`  参数
 ```
  request 是 *http.Request 实例
  writer 是 http.ResponseWriter 实例
  query := request.URL.Query()
  // 第一种方法
  name := query["name"][0]
  
  // 第二种方法
  name := query.Get("name")
  
  fmt.Printf("GET : id=%s\n", id)
  fmt.Fprintf(writer, `{"code":0}`) // 往页面返回数据
 ```
 
### 获取 `post` 参数
 + `post` 请求有两种
    - `application/json`
    - `application/x-www-form-urlencoded` 
    - `request.Header.Get("Content-Type")` 可以通过方法获取头部信息
 + 需要引入包 `encoding/json`
 ```
  // 根据 body 创建一个 json 解析器实例
  decoder := json.NewDecoder(request.Body)
  // 存放参数
  var params map[string]string
  // 解析参数 存入 map
  decoder.Decode(&params)
  fmt.Printf("POST json: name=%s , age=%s \n",params["name"], params["age"])
  fmt.Fprintf(writer, `{"code":0}`)
 ```
 
 ```
  // 处理 application/x-www-form-urlencoded 类型的POST请求
  request.ParseForm()  // 解析 post 请求
  // 第一种方式
  username := request.Form["name"][0]
  userage := request.Form["age"][0]
  
  // 第二种方式
  username := request.Form.Get("name")
  userage := request.Form.Get("age")
  
  fmt.Println("获取到的 post 数据 : ", request.Form)
  
 ```
 
 ### 返回 `JSON` 格式 
  + 参考 本目录下的 `simple_server.go` 文件
 
 