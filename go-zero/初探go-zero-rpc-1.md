## 前期准备
因为 服务发现和服务注册用到了 etcd , 但是最新的 grpc 跟 etcd 不兼容，

所以 protoc-gen-go 跟 grpc 的版本要降级

`go get -u github.com/golang/protobuf/protoc-gen-go@v1.3.2`

`go get google.golang.org/grpc@v1.29.1`


先贴以下我的 go.mod
```go 
require (
	github.com/golang/protobuf v1.4.2
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/tal-tech/go-zero v1.1.4
	google.golang.org/grpc v1.29.1
)
```

## 简单初始化 一个 grpc server
> 以下代码来自 go-zero 文档, 地址: [https://github.com/tal-tech/zero-doc/blob/main/doc/shorturl.md](https://github.com/tal-tech/zero-doc/blob/main/doc/shorturl.md)

 - 建一个文件，进入 初始化: go mod init *****
 - `goctl rpc template -o transform.proto`
 - `goctl rpc proto -src transform.proto -dir .`
 - `go get github.com/golang/protobuf@v1.4.2`
 - `go get google.golang.org/grpc@v1.29.1`
 - `go get `  获得所有的依赖
 - 查看下 etc/transform.yaml 文件里的etcd 地址:端口 是否正确
 - go run transform.go 运行起来程序

## 现在我们来了解下 go-zero 服务是怎么运行了
### 新建的实现 proto 中声明的 rpc 方法
```go 
// 第二个参数的值 会赋值给 s 的register属性，然后 调用 s.Start() 方法, 注册服务
s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		transform.RegisterTransformerServer(grpcServer, srv)
	})
```
### 让我们看看新建的时候都做了什么
#### 第一步
```go 
func MustNewServer(c RpcServerConf, register internal.RegisterFn) *RpcServer {
	server, err := NewServer(c, register)// 这个是第二步
	if err != nil {
		log.Fatal(err)
	}

	return server
}
```
#### 第二步
```go 
func NewServer(c RpcServerConf, register internal.RegisterFn) (*RpcServer, error) {
	var err error
    // 如果要 Auth 认证，需要配置 redis ，这里检查是否配置了 redis
	if err = c.Validate(); err != nil {
		return nil, err
	}

	var server internal.Server
	metrics := stat.NewMetrics(c.ListenOn)
	if c.HasEtcd() {
        // 检查 配置的 ip
		listenOn := figureOutListenOn(c.ListenOn)
		// 重要的一步 连接etcd 并注册服务，并保持续期  这是第三步
        server, err = internal.NewRpcPubServer(c.Etcd.Hosts, c.Etcd.Key, listenOn, internal.WithMetrics(metrics))
		if err != nil {
			return nil, err
		}
	} else {
		server = internal.NewRpcServer(c.ListenOn, internal.WithMetrics(metrics))
	}

	server.SetName(c.Name)
	if err = setupInterceptors(server, c, metrics); err != nil {
		return nil, err
	}

	rpcServer := &RpcServer{
		server:   server,
		register: register,
	}
	if err = c.SetUp(); err != nil {
		return nil, err
	}

	return rpcServer, nil
}
```
#### 第三步
```go 
server, err = internal.NewRpcPubServer(c.Etcd.Hosts, c.Etcd.Key, listenOn, internal.WithMetrics(metrics))
```

上面的 NewRpcPubServer 方法里的代码
```go 
func NewRpcPubServer(etcdEndpoints []string, etcdKey, listenOn string, opts ...ServerOption) (Server, error) {
	registerEtcd := func() error {
        // 连接 etcd 并设置租约
		pubClient := discov.NewPublisher(etcdEndpoints, etcdKey, listenOn)
		return pubClient.KeepAlive()
	}
	server := keepAliveServer{
		registerEtcd: registerEtcd,
                      // 这是第四步
		Server:       NewRpcServer(listenOn, opts...),
	}

	return server, nil
}
```

#### 第四步
```go 
func NewRpcServer(address string, opts ...ServerOption) Server {
	var options rpcServerOptions
	for _, opt := range opts {
		opt(&options)
	}
	if options.metrics == nil {
		options.metrics = stat.NewMetrics(address)
	}
    // 这个 server 是比较重要的 他有个方法
	return &rpcServer{
		baseRpcServer: newBaseRpcServer(address, options.metrics),
	}
}
```

rpcServer 的 Start 方法, 这个方法就是注册grpc 服务 和 注册拦截器
```go 
func (s *rpcServer) Start(register RegisterFn) error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		serverinterceptors.UnaryTracingInterceptor(s.name),
		serverinterceptors.UnaryCrashInterceptor(),
		serverinterceptors.UnaryStatInterceptor(s.metrics),
		serverinterceptors.UnaryPrometheusInterceptor(),
	}
	unaryInterceptors = append(unaryInterceptors, s.unaryInterceptors...)
	streamInterceptors := []grpc.StreamServerInterceptor{
		serverinterceptors.StreamCrashInterceptor,
	}
	streamInterceptors = append(streamInterceptors, s.streamInterceptors...)
	options := append(s.options, WithUnaryServerInterceptors(unaryInterceptors...),
		WithStreamServerInterceptors(streamInterceptors...))
	server := grpc.NewServer(options...)
	register(server)
	// we need to make sure all others are wrapped up
	// so we do graceful stop at shutdown phase instead of wrap up phase
	waitForCalled := proc.AddWrapUpListener(func() {
		server.GracefulStop()
	})
	defer waitForCalled()

	return server.Serve(lis)
}
```

#### 第五步 启动起来
```go 
func (rs *RpcServer) Start() {
	if err := rs.server.Start(rs.register); err != nil {
		logx.Error(err)
		panic(err)
	}
}
```

就这样 grpc 服务就跑起来了
 


