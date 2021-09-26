##### 反向 RPC

###### 服务端实现
```go
package main

import (
	"net"
	"net/rpc"
	"time"
)

type HelloService struct {}

func (h *HelloService) Hello(r string, rply *string) error {
	*rply = "hello" + r
	return nil
}

func main(){
	rpc.Register(new(HelloService))

	for {
		conn, _ := net.Dial("tcp", "localhost:12345")

		if conn == nil {
			time.Sleep(time.Second)
			continue
		}

		rpc.ServeConn(conn)
		conn.Close()
	}
}

```

###### 客户端实现

```go
package main

import (
	"fmt"
	"net"
	"net/rpc"
)

func main(){
	listener, err := net.Listen("tcp", ":12345")

	if err != nil {
		fmt.Println("listenTCP failed : ", err)
	}

	clientChan := make(chan *rpc.Client)

	go func(){
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("accept error : ", err)
			}

			clientChan <- rpc.NewClient(conn)
		}
	}()

	doClientWork(clientChan)
}

func doClientWork(clientChan <-chan *rpc.Client) {
	client := <-clientChan
	defer client.Close()

	var reply string
	err := client.Call("HelloService.Hello", "hello", &reply)
	if err != nil {
		fmt.Println("get rpc service failed ", err)
	}
	fmt.Println("rpc return data : ", reply)

}

```