# channel
#### select 配合 case default 使用， case 条件必须是一个 chan
    ```
           chan1 := make(chan int)
           t := time.Ticker(time.Second * 3)
           select {
               case <- chan1 :
                   fmt.Println("在走这条分支")
               case <- t.C :
                   fmt.Println("3 秒超时走这个分支")
           }
           t.Stop()
    ```

#### 共享资源 通过加锁的方式 实现共享内存通信模型
> 我们所说的同步其实就是在控制多个线程对共享资源的访问：
    一个线程在想要访问某一个共享资源的时候，需要先申请对该资源的访问权限，
    并且只有在申请成功之后，访问才能真正开始；而当线程对共享资源的访问结束时，
    它还必须归还对该资源的访问权限，若要再次访问仍需申请.
    
> 一但某个资源被确定为共享资源，则 不管是在子线程 或是 主线程，访问共享资源，一定要加锁，访问结束后 一定要解锁，
    不然很容易引起资源竞争
    
 - sync.Mutex  sync.RWMutex
        + `sync.Mutex` 互斥锁 : 不管是读操作还是写操作都会独占锁 一旦锁定，会阻止其它线程加锁，从而阻塞；
        + `sync.RWMutex` 读写锁 : 读锁调用 `RLock()` `RUnlock ` , 写锁调用 `Lock()` `Unlock()`
        + 为了提升性能，读操作往往是不需要阻塞的, 如果读多写少 建议用 `sync.RWMutex` , 实际 `RWMutex` 底层继承了 `Mutex`
        
        
#### sync.Cond 共享资源状态发生变化时 通知其它因此而阻塞的线程
 - `var mutex = new(sync.RWMutex)` 新建 `var cond = sync.NewCond(mutex)`  // 必须传入一个锁
 - `cond.Wait()` // 阻塞等待通知
 - `cond.Signal()` // 通知一个 wait 阻塞的线程
 - `cond.Broadcast()` // 广播通知线程
 
#### 原子操作
 > 原子操作通常是 CPU 和 操作系统提供支持的，由于执行过程中不会中断，所以可以完全消除竞态条件，从而绝对保证
 并发安全性，此外，由于不会中断，所以原子操作本身要求也很高，既要简单，又要快速。Go语言的原子操作也是基于CPU和
 操作系统的，由于简单和快速的要求，只针对少数数据类型的值提供了原子操作函数，这些函数都位于标准库代码包 `sync/atomic`中。
 这些原子操作包括加法(Add)、比较并交换(Compare And Swap ，简称CAS)、加载(Load)、 存储(Store) 和 交换 (Swap)
 
 - 加减法
    + 我们可以通过 atomic 包提供的下列函数实现加减法的原子操作，第一个参数是操作数对应的指针，第二个参数是 加/减 值:
    ```
        func AddInt32(addr *int32, delta int32) (new int32)
        func AddInt64(addr *int64, delta int64) (new int64)
        func AddUint32(addr *uint32, delta uint32) (new uint32)
        func AddUint64(addr *uint64, delta uint64)(new uint64)
        func AddUintptr(addr *uintptr, delta uintptr) (new uintptr)
        前两个函数 可以传递负数实现 减法操作
    ```
 - 比较并交换
    + 下面函数 第一个参数是操作数对应的指针，第二、第三个参数是待比较和交换的旧值 和 新值
    + 这些函数会在交换之前先判断 old 和 new 对应的值是否相等，如果不相等才会交换
    ```
    func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
    func CompareAndSwapInt64(addr *int64, old, new int64) (swapped bool)
    func CompareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool)
    func CompareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool)
    func CompareAndSwapUint64(addr *uint64, old, new uint64) (swapped bool)
    func CompareAndSwapUintptr(addr *uintptr, old, new uintptr) (swapped bool)
    ```
    
 - 加载
    + 操作函数仅传递一个参数，即待操作数对应的指针，并且有一个返回值，返回传入指针指向的值:
    ```
    func LoadInt32(addr *int32) (val int32)
    func LoadInt64(addr *int64) (val int64)
    func LoadPointer(addr *unsafe.Pointer) (val unsafe.Ponter)
    func LoadUint32(addr *uint32) (val uint32)
    func LoadUint64(addr *uint64) (val uint64)
    func LoadUintptr(addr *uintptr) (val uintptr)
    ```
    
 - 存储
    + 存储相关的原子函数 第一个参数表示待操作变量对应的指针，第二个参数表示要存储到待操作变量的数值:
    ```
        func StoreInt32(addr *int32, val int32)
        func StoreInt64(addr *int64, val int64)
        func StorePointer(addr *unsafe.Pointer, val unsafe.Pointer)
        func StoreUint32(addr *uint32, val uint32)
        func StoreUint64(addr *uint64, val uint64)
        func StoreUintptr(addr *uintptr, val uintptr)
    ```
    + 该操作可以看作是加载操作的逆向操作，一个用于读取，一个用于写入，通过上述原子函数存储数值的时候，不会出现存储流程进行到一半被中断的情况
    比如我们可以通过 `StoreInt32` 函数改写上述设置 `y` 变量的操作代码
    
 - 交换
    + 不管旧值跟新值是否相等，都会通过新值替换旧值，有一个返回值，会返回旧值
        ```
         func SwapInt32(addr *int32, new int32)(old int32)
         func SwapInt64(addr *int64, new int64)(old int64)
         func SwapPointer(addr *unsafe.Pointer, new unsafe.Pointer)(old unsafe.Pointer)
         func SwapUint32(addr *uint32, new uint32)(old uint32)
         func SwapUint64(addr *uint64, new uint64)(old uint64)
         func SwapUintptr(addr *uintptr, new uintptr)(old uintptr)
        ```
        
 - 原子类型 you Store 和 Load 两个指针方法
    ```
    type Value struct {
        v interface{}
    }
    var v atomic.Value
    v.Store(100)
    fmt.Println("v:", v.Load())
    ```
    > 存储值不能是 nil；其次，我们向原子类型存储的第一个值，
    决定了它今后能且只能存储该类型的值。如果违背这两条，编译时会抛出 panic
    
 - sync.WaitGroup 开箱即用，并发安全
    + Add : WaitGroup 类型有一个计数器，默认值是 0, 我们可以通过 Add 方法来增加这个计数器的值, 通常我们可以通过这个方法
    来标记需要等待的子协程数量；
    + Done : 当某个子协程执行完毕后，可以通过 Done 方法标记已完成，该方法会将所属 WaitGroup 类型实例计数器减一, 通常可以
    通过 defer 语言来调用它;
    + Wait : Wait 方法的作用是阻塞当前协程, 直到对应 WaitGroup 类型实例的计算器值归零, 如果在该方法被调用的时候, 对应计数
    器的值已经是 0, 那么它将不会做任何事情.
    
       ```
        func add_num(a, b int, deferFunc func()){
            defer func(){
                deferFunc()
            }()
            
            c := a + b
            fmt.Printf("%d + %d = %d\n", a, b, c)
        }
        
        func main(){
            var wg sync.WaitGroup
            wg.Add(10)
            for i := 0; i < 10; i++ {
                go add_num(i, 1, wg.Done)
            }
            wg.Wait()
            fmt.Println("程序结束运行了")
        }
       ```
       
 - sync.Once 类型
    + 其主要用途是保证指定函数代码只执行一次, 类似于单列模式，常用于应用启动时的一些全局初始化操作
    它只提供了一个 Do 方法, 该方法只接受一个参数, 且这个参数的类型必须是 `func`
    ，即无参数无返回值的函数类型
        ```
        func dosomething(o *sync.Once) {
            fmt.Println("start : ")
            o.Do(func(){
                fmt.Println("Do Something...")
            })
            fmt.Println("Finished.")
        }
        func main(){
            o := &sync.Once{}
            go dosomething(o)
            go dosomething(o)
            time.Sleep(time.Second * 1)
        }
        ```
#### 通过 context 包提供的函数实现多协程之间的协作
 + sync.WaitGroup 类型优化通道对多协程协调的处理，但是现在有一个问题：
 就是我们在启动子协程之前都已经明确知道子协程的总量，如果不知道的话，可以用context 来实现
 + 概念就不多说了，直接写代码
 ```
  func AddNum(a *int32, b int, deferFunc func()){
    defer func(){
        deferFunc()
    }()
    for i := 0; ; i++ {
        curNum := atomic.LoadInt32(a)
        newNum := curNum + 1
        time.Sleep(time.Millisecond * 200)
        if atomic.CompareAndSwapInt32(a, curNum, newNum){
           fmt.Printf("number当前值: %d [%d-%d]\n", *a, b, i)
           break
        } else {
            fmt.Println("error")
        }
    }
  }
  
  func main(){
    total := 10
    var num int32
    fmt.Printf("start num : %d \n", num)
    fmt.Println("start child goroute")
    ctx, cancelFunc := context.WithCancel(context.Backgroud())
    for i := 0; i < total; i++ {
        go AddNum(&num, i, func(){
            if atomic.LoadInt32(&num) == int32(total) {
                cancelFunc()
            }
        })
    }
    
    <- ctx.Done()
    fmt.Println("done")
  }
 ```
 
  
        
 
    

        
        
