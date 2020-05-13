#### Context 包的理解
 - Context 通常被译作上下文(比较抽象). 一般理解为 程序单元 的一个运行状态、现场、快照。 在 Go 语言中，程序单元也就指的是 Goroutine
 - 每个 Goroutine 在执行之前，都要先知道程序当前的执行状态，通常将这些执行状态封装在一个 Context 变量中，传递给要执行的 Goroutine 中。
 上下文则几乎已经成为传递与请求同生存周期变量的标准方法。在网络编程下，当接收到一个网络请求Request. 处理 Request 时，我们可能需要开启不同的
 Goroutine 来获取数据与逻辑处理，即一个请求 Request ，会在多个 Goroutine 中处理。而这些 Goroutine 可能需要共享 Request 的一些信息;
 同时当 Request 被取消或者 超时的时候，所有从这个 Request 创建的所有 Goroutine 也应该被结束
 
 - Context 包不仅实现了在程序单元之间共享状态变量的方法同时能通过简单的方法使我们在被调用程序单元的外部，通过设置 ctx 变量值，将过期或撤销这些
 信号传递给被调用的程序单元 A 调用 B 的 API, B 再调用 C 的API, 若如果 A 调用 B 取消, 那么也要取消 B 调用 C，通过在 A,B,C 
 的 API 调用之间传递 Context, 以及判断其状态，就能解决此问题. 有很多 gRPC 接口函数第一个参数是 context.Context
 
 ```go
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}
```
 - 子 Goroutine 是从 父 Goroutine 中创建的，父 Goroutine 应该有权力关闭子 Goroutine 