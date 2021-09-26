#### golang 启动 web 服务的方法
 - 最原始的方法
    ```go
    ListenAndServe(addr string, handler Handler) error {
       server := &Server{Addr:addr, Handler: handler}
       return server.ListenAndServe()
    }
    ```
    > 上面的方法是最原始的方法，直接写一个地址，就可以启动web服务, 也可以写一个继承 Handler 接口的结构体
    来实现自定义的路由处理逻辑
    
 - 从上面的代码可以看出 底层是实例了一个 `Server`，然后调用 `Server` 的 `ListenAndServe()` 方法来启动服务
    那让我们来看下 `Server` 结构体 都有哪些属性吧
    ```go
       // 只看我们可以修改的属性
       type Server struct {
           Addr string // 监听的地址
           Handler Handler // 处理路由的程序 不填默认调用 http.DefaultServeMux 来处理路由
           // 看字面意思就知道是处理 https 的
           // tls.Config 有个属性 Certificates []Certificate
           // Certificate 里有属性 Certificate PrivateKey 分别保存 certFile keyFile 证书的内容
           TLSConfig *tls.Config
           // 读取内容超时时间
           ReadTimeout time.Duration
           // 读取 request headers 超时时间 一般用这个 不用上面那个
           ReadHeaderTimeout time.Duration
           // 写入响应超时之前的最大时间
           WriteTimeout time.Duration
           // 启用 keep-alives 时等待下一个请求的最大时间量。如果 IdieTimeout=0 则使用 ReadTimeout, 
           // 如果 ReadTimeout=0则使用 ReadHeaderTimeout
           // golang server 是默认支持 keep-alives 需要客户端请求时说明是否需要 keep-alives
           // 如果在高频高并发的场景下, 有很多请求是可以复用的时候 最好开启 keep-alives 减少三次握手 tcp 销毁连接时有个 timewait 时间
           IdieTimeout time.Duration
           // 最大 header 内容 如果 为零，使用 DefaultMaxHeaderBytes 的值
           MaxHeaderBytes int
           // 是否启用 keep-alives
           SetKeepAlivesEnabled(v bool)
           // 启用 http 服务
           ListenAndServe() error
           // 启用 https 服务
           ListenAndServeTLS(certFile, keyFile string) error
           // 关闭监听
           Close() error
           // 优雅的关闭服务
           Shutdown(ctx context.Context) error
           RegisterOnShutdown(f func())
       }
    ```
    
#### gin 框架中 是这么示例的

```go
    Router := gin.Default()
    .....
    	
    s := &http.Server{
		Addr:           address,
		Handler:        Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
    // 启动服务
    s.ListenAndServe()
```
#### 内部简单的工作流程
```go
    // 1. 监听端口
    ln, err := net.Listen("tcp", addr)
    defer ln.Close()
    // 接收请求
    for {
    	c, err := ln.Accept()
    	go serve(c)
    }
    
    func serve(conn) {
        ServeHTTP(ResponseWriter, *Request)
    }
    
    func ServeHTTP(rw ResponseWriter, req *Request) {
    	// s 是上面实例 Server 得到的
    	// Hander 是一个有 ServeHTTP 方法的接口
    	// 所以 只要你实现了 ServeHTTP 方法就能
    	handler := s.Handler
        if handler == nil {
            handler = DefaultServeMux
        }
    	handler.ServeHTTP(rw, req)
    }
    
    
```
 