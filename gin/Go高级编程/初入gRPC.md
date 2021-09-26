#### ���� gRPC

##### ǰ��׼��
 - ��װ [protoc������](https://github.com/google/protobuf/releases)
 - ��װ protoc-gen-go `go get github.com/golang/protobuf/protoc-gen-go`
 - ��װ grpc `go get google.golang.org/grpc`
 
##### ��д���������ļ� protobuf
```proto
// path/filename :  hello/hello.proto
syntax = "proto3";

package hello;

message String {
    string value = 1;
}

// �������� HelloService ���� 2���ӿ�
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

##### �� proto �ļ����ɴ���
 > `protoc --go_out=plugins=grpc:. hello.proto`  ���ڵ�ǰ�ļ������� hello.pb.go �ļ�
 1. `protoc` ������� ��װ `protoc` ������ 
 2. `--go_out` ִ�����ѡ�� ���밲װ `protoc-gen-go`
 
##### ��д����˴���
```go
package main

import (
	"context"
	"fmt"
	"netrpc/hello" // netrpc ���ҵ�ģ���� �������ɵ� hello.pb.go �ļ������ hello �ļ�����
	"google.golang.org/grpc"
	"net"
)
// ʵ�� HelloServiceServer �ӿ�
type HelloServiceImpl struct {}

func (p *HelloServiceImpl) Hello(
	ctx context.Context, args *String,
	) (*String, error) {
		reply := &String{Value:"Hello : " + args.GetValue()}
		return reply, nil
}

func main() {
	// ����һ�� grpc �������
	grpcServer := grpc.NewServer()
	// ����ע��
	hellO.RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))

	fmt.Println("start listen port 12345")
	// �����˿�
	lis , err := net.Listen("tcp", ":12345")

	if err != nil {
		fmt.Println("tcp : 12345 failed : ", err)
	}
	// �ڼ����Ķ˿��� �ṩ gRPC ����
	grpcServer.Serve(lis)
}

```

##### �ͻ��˴���ʵ��

```go
package main

import (
	"fmt"
	"netrpc/hello" // netrpc ���ҵ�ģ���� �������ɵ� hello.pb.go �ļ������ hello �ļ�����
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
    // ��΢�����Ϸ�������
	reply, err := client.Hello(context.Background(), &String{Value:"hello"})

	if err != nil {
		fmt.Println("client failed : ", err)
	}


    // ��ӡ���յ�����Ϣ
	fmt.Println(reply.GetValue())


}
```