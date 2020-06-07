### time 包
 - 时间格式 `2016-01-02 15:04:05`
 
#### 获取 time 包里 Time 对象
 - 获取当前时间 `time.Now()`
 
 - 获取指定时间 `func ParseInLocation(layout, value string, loc *Location)(Time, error)`
    ```go
    time.ParseInLocation("2006-01-02 15:04:05", "2020-06-06 10:59:59", time.Local)
    ```
    
 - 获取时间戳
    ```go
    fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
    ```
    
 - 时间段
    ```go
    dt, _ := time.ParseDuration("1m50s") // 如果是s 最大是 60 h 最大是24 依次类推
    now := time.Now()
    newTime := now.Add(dt)
    fmt.Println(newTime) // 1分50秒后的时间
    ```
    
 - 时间差
    ```go
    now := time.Now()
    now.Sub(Time) // 参数是另一个时间对象， 返回 Duration 类型
    now.Add(Duration) // 时间最小单位
    ```
    
 - 比较两个时间点
    ```go
    now := time.Now()
    now.After(Time) // 参数是一个 Time 结构体的实例 Before 也是一样的
    now.Before(Time)
    ```
    
#### 定时相关的方法
 - `time.After(d Duration)` 表示多少时间之后, 但是在取出 channel 内容之前不阻塞, 后续程序可以继续执行
    `After`通常被用来处理程序超时问题 
    ```go
    select {
    case m := <-c:
 	    fmt.Println("hello world")
    case <- time.After(5 * time.Minute) :
 	    fmt.Println("timed out")
    }
    ```
 - `time.Sleep(d Duration)`  程序休眠指定的时间, 然后才继续执行
 
 - `time.Tick(d Duration) <-chan Time` 用法跟 `time.After` 差不多, 但是它是表示每隔多少时间之后, 是一个重复的过程(可以当心跳使用)，其它与 `After` 一致
 
 - `ticker := time.NewTicker(1 * time.Second)` 跟 Tick 用法一样, 但是可以调用 `ticker.Stop()` 来停止定时  
 
 