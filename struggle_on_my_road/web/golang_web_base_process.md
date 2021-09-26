### 读 Go Web 一书有感
#### 简单的 客户端服务端请求流程
1. 客户端通过 TCP/IP 协议建立到服务器的TCP连接
    1. 客户端通过 URL(Uniform Resource Locator 统一资源定位符)描述一个网络资源，基本格式如下
        schema:\/\/host\[:port#\]/path/.../\[?query-string\]\[#anchor\] 
        1. schema 指定底层使用的协议 (http https ftp)
        2. host   http服务器的 IP 地址或者域名
        3. port
        4. path   访问资源的路径
        5. query-string 发送给 http 服务器的数据
        6. anchor  锚
    2. 客户端会拿着域名 先去本地host 文件 查看IP是否存在，不存在查找贝蒂 DNS 解析器缓存，不存在去 DNS服务器上查询。。。
    3. 获取 IP 客户端拿着IP 找到对应的服务器 发送给服务器请求数据
    
2. GO http包源码  http.ListenAndServe(":9090", nil) 方法主要的内部实现
    ```
    func (srv *Server) Server(l net.Listener) error {
        defer l.Close()
        var tempDelay time.Duration
        for {
            re, e := l.Accept()
            if e != nil {
                if ne, ok := e.(net.Error); ok && ne.Temporary(){
                    if tempDelay == 0 {
                        tempDelay = 5 * time.Millisecond
                    } else {
                        tempDelay *= 2
                    }
                    
                    if max := 1 * time.Second; tempDelay > max {
                        tempDelay = max
                    }
                    
                    log.Printf("http: Accept error:%v; retrying in %v", e, tempDelay)
                    time.Sleep(tempDelay)
                    continue
                }
                return e
            }
            
            tempDelay = 0
            if srv.ReadTimeout != 0 {
                rw.SetReadDeadline(time.Now().Add(srv.ReadTimeout))
            }
            if srv.WriteTimeout != 0 {
                rw.SetWriteDeadline(time.Now().Add(srv.WriteTimeout))
            }
            c, err := srv.newConn(rw)
            if err != nil {
                continue
            }
            go c.serve()     // 对每一个请求 开一个协程 进行处理
        }
        panic("not reached")
    }
    ```
    
    
3. Go 的 http 包详解
    1. 客户端的每次请求都会创建一个 Conn, 这个Conn 里面保存了该次请求的信息，然后再传递到
    对应的 handler , 该 handler 中便可以读取到相应的 header 信息, 这样可以保证每个请求的独立性
    2. ServeMux 
        1. Conn.serve() 的时候内部调用 http 包默认的路由器, 通过路由 传递请求信息 到 handler,
        2. 路由的实现
        ```  
            type ServeMux struct {
                mu sync.RWMutex         // 请求涉及到并发处理, 需要一个锁机制
                m  map[string]muxEntry  // 路由规则 一个 string 对应一个 mux 实体, 这里 string 就是注册的路由
            }
            type muxEntry struct {
                explicit bool       // 是否精确匹配
                h        Handler    // 这个路由表达式对应哪个 handler
            }
            // Handler 定义
            type Handler interface {
                ServeHTTP(ResponseWriter, *Request)     // 路由实现器
            }
            type HandlerFunc func(ResponseWriter, *Request)
            func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
                f(w, r)
            }
            
            // 存储好 路由后 ，路由接收到请求之后调用 mux.handler(r).ServeHTTP(w, r)
            func (mux *ServeNux) handler(r *Request) Handler {
                nux.mu.RLock()
                defer nux.mu.RUnlock()
                h := mux.match(r.Host + r.URL.Path)
                if h == nil {
                    h = mux.match(r.URL.Path)
                }
                if h == nil {
                    h = NotFundHandler()
                }
                return h
            }
            根据用户请求的URL和路由器里面存储的map去匹配, 匹配到后返回存储的 handler , 调用这个 handler
            的 ServeHTTP
        ```
        
4. 通过 上面的介绍我们了解了 整个路由过程, Go其实支持 外部 实现的路由器 ListenAndServe 的第二个参数就是用以
配置外部路由器的， 它是一个 Handler接口，即外部路由器只要实现和 Handler 接口就可以, 我们可以在自己是实现的路由器
ServeHTTP 里面实现自定义路由功能 
    ```
    package main
    
    import (
    	"fmt"
    	"net/http"
    )
    
    type MyMux struct {
    
    }
    
    func (m *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request){
    	if r.URL.Path == "/" {
    		fmt.Fprintf(w, "hello MyMux")
    		return
    	}
    	http.NotFound(w, r)
    	return
    }
    
    func main(){
    	mux := &MyMux{}
    
    	http.ListenAndServe("127.0.0.1:9090", mux)
    }
    ```
    
5. 总结下上面执行过程
    1. 调用 http.HandleFunc
        1. 调用了 DefaultServeMux 的 HandleFunc
        2. 调用了 DefaultServeMux 的 Handle
        3. 往 DefaultServeMux 的 map\[string\]muxEntry 中增加对应的 handler 和 路由规则
    2. 调用http.ListenAndServe(":9090", nil)
        1. 实例化Serve
        2. 调用 Serve 的 ListenAndServe()
        3. 调用 net,Listen("tcp", addr) 监听端口
        4. 启动一个 for 循环，在循环体内Accept请求
        5. 对每个请求实例化一个 Conn, 并开启一个 goroutine 为这个请求进行服务 `go c.serve()`
        6. 读取每个请求的内容， w, err := c.readRequest()
        7. 判断 handler 是否为空, 如果没有设置 handler  handler 就设置为 defaultservemux
        8. 调用 handler 的 ServeHttp
        9. 下面 进入默认  DefaultServeMux.ServeHTTP
        10. 根据 request 选择 handler, 并且进入到 这个 handler 的 ServeHTTP mux.handler(r).ServeHTTP(w, r)
        
        
#### 表单
`r *http.Request`
1. r.ParseForm() // 解析传递的参数，如果没有调用 ParseForm 方法，则 无法获取表单的数据
2. r.URL.Path
3. r.URL.Scheme
4. r.Form   // 这是一个map 保存着 form 表单发送的数据 r.Form\[key\]\[0\] 这么获取表单数据
5. r.Method  // 获取请求的方法 GET POST PUT ...
6. r.FormValue(key) // 获取表单数据, 相当于 r.ParseForm 和 r.Form 的组合

##### 表单验证
1. r.Form.Get() 可以获取值 如果值不存在则为空
2. 数字
```
1.
getInt, err := strconv.Atoi(r.Form.Get("age"))
if err != nil {
    // 数字转换出错 
}
2. if m, _ := regexp.MatchString("^[0-9]+$", r.Form.Get("age")); !m {
    return false
}
```
3. 中文
```
if m, _ := regexp.MatchString("^[\\x{4e00}-\\x{9fa5}]$", r.Form.Get("name")); !m {
    return false
}
```
4. 英文正则 `^[a-zA-Z]+$`
5. 电子邮件正则 `^([\w\.\_]{2,10})@(\w{1,}).([a-z{2,4}])$`
6. 手机号 `^1[3|4|5|8][0-9]\d{4,8}$`  

##### XSS 防护
1. `template.HTMLEscapeString(s string) string `
2. `template.HTMLEscape(w io.Writer, b []byte)`
3. `HTMLEscaper(args ...interface{}) string`


##### cookie
1. cookie 的写入
    1. `cookie := http.Cookie{Name:"username", Value:"zzh", Expires: 1e9}`
    2. `http.SetCookie(w, &cookie)`
    
2. cookie 的读取
    1. `cookie, _ := r.Cookie("username")`
    2. `for _, cookie := range r.Cookies() {}`
    
3. `session` 可以根据 `cookie` 自己实现


#### 正则
1. 正则匹配 `regexp` 包
```
func Match(pattern string, b []byte) (matched bool, error error)
func MatchReader(pattern string, r io.RunReader) (matched bool, error error)
func MatchString(pattern string, s string) (matched bool, error error)
```
2. 正则获取内容
    1. 先解析正则是否合法
        ```
        func Compile(expr string) (*Regexp, error)
        func CompilePOSIX(expr string) (*Regexp, error)
        func MustCompile(str string) *Regexp
        ```
    2. 获取正则匹配内容
        ```
        func (re *Regexp) Find(b []byte) []byte
        func (re *Regexp) FindAll(b []byte, n int) [][]byte
        func (re *Regexp) FindAllIndex(b []byte, n int) [][]int
        func (re *Regexp) FindAllSubmatch(b []byte, n int) [][][]byte
        func (re *Regexp) FindAllSubmatchIndex(b []byte, n int) [][]int
        func (re *Regexp) FindIndex(b []byte) (loc []int)
        func (re *Regexp) FindSubmatch(b []byte) [][]byte
        func (re *Regexp) FindSubmatchIndex(b []byte) []int
        // 还有string io.RuneReader
        ```
