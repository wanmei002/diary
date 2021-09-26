##### ����˴���
```go
package main

import (
	"fmt"
	"net"
	"net/rpc"
)

type HelloService struct {
	conn net.Conn
	isLogin bool
}

func (p *HelloService) Hello(request string, reply *string) error {
	p.Login(request, reply)
	if !p.isLogin {
		return fmt.Errorf("please login")
	}
	*reply = "hello : " + request + ", from " + p.conn.RemoteAddr().String()
	return nil
}

func (p *HelloService) Login(request string, reply *string) error {
	//fmt.Println("login : request : ", request)
	if request != "zzh:zyn" {
		return fmt.Errorf("auth failed")
	}
	//fmt.Println("login ok")
	p.isLogin = true
	return nil
}


// Ϊÿ���������������� RPC ����
func main(){
	fmt.Println("start listen tcp port-12345")
	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		fmt.Println("listenTCP failed : ", err)
	}

	fmt.Println("listen access")

	for {
		fmt.Println("wait request")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept failed : ", err)
		}

		p := rpc.NewServer()
		err = p.Register(&HelloService{conn:conn})
		if err != nil {
			fmt.Println("register rpc failed : ", err)
		}
		p.ServeConn(conn)
		conn.Close()

	}


}

```


##### �ͻ���
```go
package main

import (
	"fmt"
	"net/rpc"
)

func main(){
	client, err := rpc.Dial("tcp", "localhost:12345")

	if err != nil {
		fmt.Println("dial failed : ", err)
	}

	var reply string
	err = client.Call("HelloService.Hello", "zzh:zyn", &reply)
	if err != nil {
		fmt.Println("call server failed : ", err)
	}
	fmt.Println("rpc service send data : ", reply)
}

```