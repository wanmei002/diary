### Context接口
#### context 基本结构
 - context 是一个基本接口, 所有的 context 对象都要实现该接口, context 的使用者在调用接口中都使用 context 作为参数类型
     ```go
     type Context interface {
  	    // 如果 context 实现了超时控制, 则该方法返回 超时时间(deadline), ok (true | false)
        Deadline() (deadline time.Time, ok bool)
        // goroutine 应该监听该方法返回的 chan , 以便及时释放资源
        Done() <-chan struct{}
        // Done 返回的 chan 收到通知的时候, 才可以访问Err() 获知因为什么原因被取消
        Err() error
        // 可以访问上游 goroutine 传递给下游 goroutine 的值
        Value(key interface{}) interface{}
     }
    
     ```
 - canceler 接口是一个扩展接口, 规定了取消通知的 Context 具体类型需要实现的接口。 context 包中的具体类型 *cancelCtx 和 *timerCtx
   都实现了该接口
   ```go
   type canceler interface {
	   // 拥有 cancel 接口实例的 goroutine 调用 cancel 方法通知后续创建的 goroutine 退出
       cancel(removeFromParent bool, err error)
       // Done 方法返回的 chan 需要后端 goroutine 来监听, 并及时退出
       Done() <-chan struct{}
   } 
   ```
 - emptyCtx 实现了 Context 接口, 但不具备任何功能, 因为其所有的方法都是空实现, 其存在的目的是作为 Context 对象树的根 (root节点). 因为
 Context 包的使用思路就是不停地调用 context 包提供的包装函数来创建具有特殊功能的 Context 实例，每一个 Context 实例的创建都以上一个 Context
 对象为参数, 最终形成一个树状的结构
 
 - context 包定义了两个全局变量和两个封装函数, 反回两个 emptyCtx 实例对象, 实际使用时通过调用这两个封装函数来构造 Context 的 root 节点
 ```go
    var (
    	background = new(emptyCtx)
    	todo       = new(emptyCtx)
    )

    func Background() Context {
    	return background
    }

    func TODO() Context {
    	return todo
    }
 ```
 
 - cancelCtx 是一个实现了 Context 接口的具体类型, 同时实现了 conceler 接口, conceler 具有退出通知方法. 注意退出通知机制不但能通知
 自己，也能逐层通知其 children 节点. 
 
 - timerCtx 是一个实现了 Context 接口的具体类型，内部封装了 cancelCtx 类型实例, 同时有一个 deadline 变量, 用来实现定时退出通知
 
 - valueCtx 是一个实现了 Context 接口的具体类型，内部封装了 Context 接口类型, 同时封装了一个 K/V 的存储变量。valueCtx 可用来传递通知信息
 
 - Background() TODO() 这两个函数是构造 Context 取消树的根节点对象, 根节点对象用作 后续 With 包装函数的实参
   With 包装函数用来构建不同功能的 Context 具体对象
   + 创建一个带有退出通知的 Context 具体对象, 内部创建一个 cancelCtx 的类型实例:
   ```go
   func WithCancel(parent Context) (ctx Context, cancel CancelFunc) 
   ```
   + 创建一个带有超时通知的 Context 具体对象, 内部创建一个 timerCtx 的类型实例:
   ```go
    func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)
    func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
    ```
   + 创建一个能够传递数据的 Context 具体对象, 内部创建一个 valueCtx 的类型实例:
   ```go
    func WithValue(parent Context, key, val interface{}) Context
    ```
    > 通过上面4个函数 返回 context 的子 context, 可以通过调用父级关闭函数来关闭 goroutine
    
    > 使用 context 包主要是解决 goroutine 的通知退出，传递数据是其一个额外功能
    
   + context 包最好传递如下数据
    + 日志信息
    + 调试信息
    + 不影响业务主逻辑的可选数据
    
    ```go
    // 简单的例子
     package main
     
     import (
        "context"
        "fmt"
        "time"
     )
     
     func main(){
        ctx, cancel := context.WithCancel(context.Background())
        // ctx1 就是 ctx 的 children
        ctx1, _ := context.WithTimeout(ctx, 10 * time.Second)
     
        // ctx2, _ := context.WithDeadline(ctx, time.Now().Add(10 * time.Second))
     
        go func(ctx1 context.Context){
            for {
                select {
                default:
                    fmt.Println("i am ctx1")
                    time.Sleep(500 * time.Millisecond)
                case <-ctx1.Done():
                    fmt.Println("ctx1 is done")
                    return
     
                }
            }
        }(ctx1)
     
        go work(ctx1)
        time.Sleep(1 * time.Second)
        // 通过调用 ctx 的关闭函数来关闭 goroutine
        cancel()
        time.Sleep(1 * time.Second)
     }
     
     func work(ctx context.Context) {
        for {
            select {
            default:
                fmt.Println("i am ctx2")
                time.Sleep(500 * time.Millisecond)
            case <-ctx.Done():
                fmt.Println("ctx2 is done")
                return
            }
        }
     }
    ```
 
 
 
 