### 先来一个入口
```go
func checkLegal(uids []int64) ([]int64, error) {
    r, err := mr.MapReduce(func(source chan<- interface{}) {
        for _, uid := range uids {
            source <- uid
        }
    }, func(item interface{}, writer mr.Writer, cancel func(error)) {
        uid := item.(int64)
        ok, err := check(uid)
        if err != nil {
            cancel(err)
        }
        if ok {
            writer.Write(uid)
        }
    }, func(pipe <-chan interface{}, writer mr.Writer, cancel func(error)) {
        var uids []int64
        for p := range pipe {
            uids = append(uids, p.(int64))
        }
        writer.Write(uids)
    })
    if err != nil {
        log.Printf("check error: %v", err)
        return nil, err
    }
    
    return r.([]int64), nil
}
```
把传入的方法单独都拿出来
```go
a := func(source chan<- interface{}) {
         for _, uid := range uids {
             source <- uid
         }
     }
```
```go
b := func(item interface{}, writer mr.Writer, cancel func(error)) {
         uid := item.(int64)
         ok, err := check(uid)
         if err != nil {
             cancel(err)
         }
         if ok {
             writer.Write(uid)
         }
     }
```
```go
c := func(pipe <-chan interface{}, writer mr.Writer, cancel func(error)) {
         var uids []int64
         for p := range pipe {
             uids = append(uids, p.(int64))
         }
         writer.Write(uids)
     }
```

### 先看看 MapReduce 方法的执行流程
```go
func MapReduce(a, b, c, opts ...Option) (interface{}, error) {
	// 启动一个无缓存的chan, 用于保存要执行的方法
    source := buildSource(a)
	return MapReduceWithSource(source, b, c, opts...)
}
```
