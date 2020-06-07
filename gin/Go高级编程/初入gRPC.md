#### 初用 gRPC

##### 前期准备
 - 安装 [protoc编译器](https://github.com/google/protobuf/releases)
 - 安装 protoc-gen-go `go get github.com/golang/protobuf/protoc-gen-go`
 - 安装 grpc `go get google.golang.org/grpc`
 
##### 编写数据描述文件 protobuf
```proto
// path/filename :  hello/hello.proto
syntax = "proto3";

package hello;

message String {
    string value = 1;
}

// 代码会根据 HelloService 生成 2个接口
// type HelloServiceServer interface {
//     Hello(context.Context, *String) (*String, error)
// }

// type HelloServiceClient interface {
//     Hello(context.Context, *String, ...grpc.CallOption) (*String, error)
// }

service HelloService {
    rpc Hello (String) returns (String);
}
```

##### 把 proto 文件生成代码
 > `protoc --go_out=plugins=grpc:. hello.proto`  会在当前文件夹生成 hello.pb.go 文件
 1. `protoc` 命令必须 安装 `protoc` 解析器 
 2. `--go_out` 执行这个选项 必须安装 `protoc-gen-go`
 
##### 编写服务端代码
```go
package main

import (
	"context"
	"fmt"
	"netrpc/hello" // netrpc 是我的模块名 上面生成的 hello.pb.go 文件存放在 hello 文件夹中
	"google.golang.org/grpc"
	"net"
)
// 实现 HelloServiceServer 接口
type HelloServiceImpl struct {}

func (p *HelloServiceImpl) Hello(
	ctx context.Context, args *String,
	) (*String, error) {
		reply := &String{Value:"Hello : " + args.GetValue()}
		return reply, nil
}

func main() {
	// 构造一个 grpc 服务对象
	grpcServer := grpc.NewServer()
	// 服务注册
	hellO.RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))

	fmt.Println("start listen port 12345")
	// 监听端口
	lis , err := net.Listen("tcp", ":12345")

	if err != nil {
		fmt.Println("tcp : 12345 failed : ", err)
	}
	// 在监听的端口上 提供 gRPC 服务
	grpcServer.Serve(lis)
}

```

##### 客户端代码实现

```go
package main

import (
	"fmt"
	"netrpc/hello" // netrpc 是我的模块名 上面生成的 hello.pb.go 文件存放在 hello 文件夹中
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main(){
	fmt.Println("start dial localhost:12345")
	conn, err := grpc.Dial("localhost:12345", grpc.WithInsecure())

	if err != nil {
		fmt.Println("grpc dial failed : ", err)
	}

	defer conn.Close()

	client := hello.NewHelloServiceClient(conn)
    // 往微服务上发送数据
	reply, err := client.Hello(context.Background(), &String{Value:"hello"})

	if err != nil {
		fmt.Println("client failed : ", err)
	}


    // 打印接收到的信息
	fmt.Println(reply.GetValue())


}
```