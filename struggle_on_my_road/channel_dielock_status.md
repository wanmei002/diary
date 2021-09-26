## channel 死锁其实没那么复杂 借鉴自大牛的总结
   `https://www.cnblogs.com/bigdataZJ/p/go-channel-deadlock.html`

## 信道
 - 信道分类
    + 无缓冲信道  `ch := make(chan string)`
    + 缓冲信道    `ch := make(chan string, 2)`
    
 - 区别
    > 无缓冲信道本身不存储信息，它只负责转手，有人传给它，它就必须要传给别人，如果只有进或者只有出的操作，都会造成阻塞。有缓冲的可以存储指定容量个变量，但是超过这个容量再取值也会阻塞。
    
1. 死锁情况1 取完了信道存储的信息再去取信息，也会死锁
    ```
    func main(){
    	chs := make(chan string, 2)
    	chs <- "string1"
    	chs <- "string2"
    	for v := range chs {
    		fmt.Println(v)
    	}
    }
    ```
    > fatal error: all goroutines are asleep - deadlock!
    
    > 原因分析 : chs 信道中数据读取完，chs 此时相当于 无缓存信道，对无缓冲信道 做 只读 或 只写 操作 都会造成死锁
    
2. 死锁情况2 无缓冲信道 只读 或 只写
```
func main(){
	chs := make(chan string)
	chs <- "string"
}
```
> fatal error: all goroutines are asleep - deadlock!

> 原因分析 : 主线程阻塞在 无缓冲信道的写入中，系统一直在等待释放无缓冲信道，结果没有释放，系统判断为死锁
> 在子线程中执行只写操作 则 主线程不会deadlock 

3. 死锁情况3 子线程中阻塞 导致 主线程阻塞
```
func main(){
	ch1 := make(chan string)
	ch2 := make(chan string)
	go func(){
		ch2 <- "string1"
		ch1 <- "string2"
	}()
	<- ch1
}
```
> 原因分析 : 子线程 阻塞在 ch2 <- "string1" 主线程中运行到 <-ch1 结果 ch1中没有数值 导致主线程阻塞



