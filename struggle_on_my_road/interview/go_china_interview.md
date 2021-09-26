### 看golang中文网面试题总结
1. defer panic 
    ```go
    package main
     import (
         "fmt"
    )
     func main() {
         defer_call()
     }
    
    func defer_call() {
        defer func() { fmt.Println("打印前") }()
        defer func() { fmt.Println("打印中") }()
        defer func() { fmt.Println("打印后") }()
    
        panic("触发异常")
     
        defer func(){ fmt.Println("异常后") }()
    }
    ```
    > 解释: 输出的结果是  打印后 - 打印中 - 打印前 - 触发异常， 在触发异常之前已经进入栈内的 defer 会触发, 但是 异常(panic) 后的defer 是不触发的 ,
    它还没来得及压入栈中
    
2. for range point 问题
    ```go
    slic := []int{0,1,2,3}
    m := make(map[int]*int)
    for key, val := range slic {
        m[key] = &val
    }
    
    for key,val := range m {
        fmt.Println(key, "->", *val)
    }
    // 结果是:
    0->3
    1->3
    2->3
    3->3
    ```
    > 解释: for range 语句块中 key, val 是局部变量, range 的时候会把 slice map 里的键值 赋值给 key, val, 
    但是 key val 的地址是不变的 变的只是地址里存储的值