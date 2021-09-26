

## 入口
```go 
// 第二方法会赋值给 s 的register属性，然后 调用 s.Start() 方法, 注册服务
s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		transform.RegisterTransformerServer(grpcServer, srv)
	})
```

## Start()
```go 
func (rs *RpcServer) Start() {
	if err := rs.server.Start(rs.register); err != nil {
		logx.Error(err)
		panic(err)
	}
}
```

## new server
### 第一步
```go 
func MustNewServer(c RpcServerConf, register internal.RegisterFn) *RpcServer {
	server, err := NewServer(c, register)
	if err != nil {
		log.Fatal(err)
	}

	return server
}
```
### start new server
#### 第一步
```go 
server, err = internal.NewRpcPubServer(c.Etcd.Hosts, c.Etcd.Key, listenOn, internal.WithMetrics(metrics))
```

#### 第二步
```go 
func NewRpcPubServer(etcdEndpoints []string, etcdKey, listenOn string, opts ...ServerOption) (Server, error) {
	registerEtcd := func() error {
	    // 这一步是连接etcd, 并返回etcd 的keeplive
		pubClient := discov.NewPublisher(etcdEndpoints, etcdKey, listenOn)
		return pubClient.KeepAlive()
	}
	server := keepAliveServer{
		registerEtcd: registerEtcd,
		// 这个 rpcServer 有一个方法 就是注册 grpc server 
		Server:       NewRpcServer(listenOn, opts...),
	}

	return server, nil
}
```
#### 注册 rpc 服务
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

### 最最主要的 
```go 
func NewServer(c RpcServerConf, register internal.RegisterFn) (*RpcServer, error) {
	var err error
	if err = c.Validate(); err != nil {
		return nil, err
	}

	var server internal.Server
	metrics := stat.NewMetrics(c.ListenOn)
	if c.HasEtcd() {
		listenOn := figureOutListenOn(c.ListenOn)
		// 上面有介绍 主要是创建 etcd 连接，
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