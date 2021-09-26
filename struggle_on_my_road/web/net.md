### net 包里的一些方法
 - ParseIP(string)  // 验证IP地址字符串 返回 IP 对象
 - type IPMask []byte // 子网掩码类型
     ```go
        func IPv4Mask(a, b, c, d byte) IPMask // 用一个四字节的IPv4地址来创建一个掩码
        func (ip IP) DefaultMask() IPMask     // 这是一个IP方法，返回默认的掩码
        func (ip IP) Mask(mask IPMask) IP     // 一个掩码可以使用一个IP地址的方法，找到该IP地址的网络
     ```
 
 - 在 net 包的许多函数和方法会返回一个指向IPAddr 的指针，这不过只是一个包含IP类型的结构体, 这种类型的主要
 用途是通过 IP 主机名执行 DNS 查找
    ```go
       type IPAddr struct {
    	    IP IP
            Zone string // IPv6 scoped addressing zone
       }
       func ResolveIPAddr(net, addr string) (*IPAddr, error)  // 通过域名 返回主机IP
       func LookupHost(host string) (addrs []string, err error) // 一个域名可能有多个地址
    ```
 - 服务运行在主机上，一个主机可以运行多个服务 (TCP, UDP, SCTP ...) 使用端口来加以区分(1 ~ 65535)
    + unix 系统中， /etc/service 文件列出了常用的端口, go 有一个函数可以获取该文件 `func LookupPort(network, service string) (port int, err error)`
   
 - TCPAddr 类型包含一个IP 和 一个 port 的结构 :
    ```go
    type TCPAddr struct {
       IP      IP
       Port    int
    }
    func ResolveTCPAddr(net, addr string) (*TCPAddr, error) // 创建一个 TCPAddr  addr-www.baidu.com:80
    ``` 
    
 - TCP Sockets
  ```go
     func (c *TCPConn) Write(b []byte) (n int, err error)
     func (c *TCPConn) Read(b []byte) (n int, err error)  // 底层还是调用的 ReadFile 函数 把 TCPConn 句柄里的数据 读取到 b 里面
     func DialTCP(net string, laddr, raddr *TCPAddr) // net-(tcp4,tcp6,tcp) laddr-nil 
  ```
 - TCP 通信
    ```go
       func main(){
    	    service := "127.0.0.1:8099"
 	        tcpAddr, err := net.ResolceTCPAddr("ip4", service) // 传入域名:端口 / IP:端口 返回一个TCP类型
            listener, err := net.ListenTCP("tcp", tcpAddr) // 监听端口 也可以直接用 net.Listen("tcp", "127.0.0.1:8099") 来监听端口
            for {//不间断的接受客户端请求
         	    conn, err := listener.Accept() // 这步会阻塞到这 直到有客户端请求
      	        if err != nil {
   	        	    continue
   	            }
      	        go handleClient(conn) // 多线程处理客户端请求
           }
       }
       
       func handleClient(conn net.Conn) {
    	    defer conn.Close()
 	        var buf [512]byte
 	        for {
       	        n, err := conn.Read(buf[0:])
    	        _, err = conn.Write(buf[0:n])
           }
 	    
       }
    ```
    
    - 超时
        + 服务端会断开那些超时的客户端，如果他们响应的不够快，比如没有及时的往服务端写一个请求。
        ` func (c *TCPConn) SetTimeout(nsec int64) error // 用在套接字读写前(Accept)` 
    
    - 存活状态
        + 即使没有任何通信 一个客户端可能希望保持连接到服务器的状态
        ` func (c *TCPConn) SetKeepAlive(keepalive bool) error`
        
 - UDP 通信
    + 主要方法
        ```go
        func ResolveUDPAddr(net, addr string) (*UDPAddr, error) // 解析地址 生成 UDPAddr 类型
        func DialUDP(net string, laddr, raddr *UDPAddr) // 连接 UDP服务器
        func ListenUDP(net string, laddr *UDPAddr) (c *UDPConn, err error) // UDP协议监听端口
        func (c *UDPConn) ReadFromUDP(b []byte) (n int, addr *UDPAddr, err error) // 处理客户端传过来的数据
        func (c *UDPConn) WriteToUDP(b []byte, addr *UDPAddr) (n int, err error) // 发送数据给客户端
        func (c *conn) Write(b []byte) (int, error) // 发送数据给服务端
        func (c *conn) Read(b []byte) (int, error) // 读取服务端发送过来的数据
        ```
 
 - 实际上 TCPConn 和 UDPConn 都实现了 Conn 接口，很大程度上，你可以通过该接口处理而不是这两种类型
    + 客户端可以使用 `func Dial(net, raddr string) (c Conn, err error)` 函数来实现tcp udp 连接
        + `net` 可以是 `tcp`, `tcp4`, `tcp6`, `udp`, `udp4`, `udp6`, `ip`, `ip4`, `ip6`
        ```go
          conn, err := net.Dial("tcp", "127.0.0.1:9999")
        ```
    
    + TCP服务端可以使用 `func Listen(net, laddr string) (l Listener, err error)` , 返回一个 `Listener` 接口对象，该接口有一个方法 ` func (l Listener) Accept() (c Conn, err error) `
     ```go
       listener, err := net.Listen("tcp", "127.0.0.1:9999")
       for {
    	    conn, err := listener.Accept()
 	         go conn.Read(buf[0:])
       }
    ```
    
    + 如果你想写一个 UDP 服务器, 这里有一个 `PacketConn` 的接口和一个实现了该接口的方法:
        ```go
          func ListenPacket(net, laddr string) (c PacketConn, err error)
          PacketConn.ReadFrom()
          PacketConn.WriteTo()
        ```
        
 
##### 套接字编程
 - Go 允许你建立原始套接字，可以使用其它协议通信，或甚至建立你自己的
    Go 提供了最低限度的支持 : 它会连接主机，写入和读取 主机之间的数据包, 

 - 下面让我们做最简单的一个例子: 如何发送一个 ping 消息给主机(Ping 使用 echo 命令的 `ICMP` 协议)
   这是一个面向字节协议, 客户端发送一个字节流到另一个主机, 并等待主机的答复，格式如下:
    + 首字节是 8, 表示 echo 消息
    + 第二个字节是 0
    + 第三 和 第四字节是整个消息的校验
    + 第五 和 第六字节是一个任意标识
    + 第七 和 第八字节是一个任意的序列号
    + 该数据包的其余部分是用户数据
 
 
 
 
 
 
 
 
 
 
 
 
 
 
 
 
 
 
 