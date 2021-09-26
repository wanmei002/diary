## 先新建两个服务
goctl rpc proto -src transform.proto -dir transformsvr/

goctl rpc proto -src shorten.proto -dir shortensvr/

## 修改配置文件
### 修改 shorten 服务的配置文件
 - 修改在 etc/ 目录下的 yaml 文件，在其中添加 transform 服务的 etcd 信息，添加内容如下:
```json
Transform: // 名字可以自己取
  Etcd:
    Hosts:
      - 127.0.0.1:2379  // transform 服务注册的 etcd
    Key: transform.rpc  // transform 服务注册的 服务名
```
 - 在 shorten 服务目录下的 internal/svc/servicecontext.go 文件中添加 transform 客户端连接
```go
    type ServiceContext struct {
    	Config config.Config
    	Transform transformer.Transformer // Transformer服务 客户端接口 位置一般是: transformsvr/transformer/transformer.go
    }
```
创建 transform 服务连接
```go
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		// 创建 transform 客户端连接
		Transform: transformer.NewTransformer(zrpc.MustNewClient(c.Transform)),// c.Transform 是在 yaml 中配置的 transform 服务的配置信息
	}
}
```

逻辑一般都在项目下的 internal/logic 文件夹下, 在 rpc 方法中，在其中调用 其它 gRPC 项目的方法可以
用 l.svcCtx.Transform 对象里的方法

