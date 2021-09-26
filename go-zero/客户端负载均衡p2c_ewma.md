### go-zero 客户端负载均衡算法 P2C_EWMA
#### p2c(pick of 2 choices)
2选1
#### EWMA(exponentially weighted averages)
指数移动加权平均, 体现的是一段时间内的平均值(范围可以通过调整参数进行修改, 从而更加灵活)

### 入口
找到实现 PickerBuilder 接口的struct
```go
type PickerBuilder interface {
	// Build takes a slice of ready SubConns, and returns a picker that will be
	// used by gRPC to pick a SubConn.
	Build(readySCs map[resolver.Address]balancer.SubConn) balancer.Picker
}
```
github.com/tal-tech/go-zero/zrpc/internal/balancer/p2c/p2c.go

主要逻辑是实现在 Pick此方法中的
#### 只有一个服务只有一个在提供服务
更新这个请求的 pick属性为当前时间
更新 inflight+1
更新 requests+1
使用此连接
#### 如果一个服务有两个在提供服务
 - 比较这两个请求的EWMA值*(inflight+1) `加1是为了避免 inflight=0`(增大正在使用的连接的EWMA值);
 如果当前没有使用这个连接，那么这个连接的 inflight=0, 如果有好几个请求在使用这个连接, 那么inflight=使用这个连接的请求数量, 当前使用这个连接的越多这个计算的它的值就越大, 越不容易选中到它
    如果 conn1>conn2 则这两个请求连接交换位置;
 - 如果当前时间-上次被选中的时间>1s 此时正好没有其它请求刚使用这个conn(就是同时有几个请求走到这个判断这了，这个请求抢占到它的使用权了, 让同时抢占这个连接的其它请求抢占失败);
    返回 EWMA*(inflight+1) 大的(如果这两个没有在使用的话 就是返回EWMA大的)
 - 否则返回 EWMA*(inflight+1) 小的
 更新 inflight+1
 更新 requests+1
 
#### 如果一个服务不仅仅有两个再提供服务
随机选择两个连接，把这两个连接跟上面 只有两个提供服务的逻辑再走一遍，开始使用此连接

### 连接使用完成后会做 buildDoneFunc 处理
#### buildDoneFunc 方法
 - 给选中的连接的 inflight-1
 - 更新选中连接的 last属性为 当前时间距离固定时间的差值
 - td := lastNow值 - lastPrev
 - w := math.Exp(-td/10s) // 这个连接最近两个请求完时间的差值除以10秒, 牛顿冷却定律中的衰减函数模型
 - lag := int64(now) - start // 当前请求的耗时
 - olag := atomic.LoadUint64(&c.lag)
 - if olag == 0 { w=0 }  // 如果是第一次请求重置 w=0
 - c.lag = olag\*w+lag\*(1-w) // EWMA 指数移动加权平均值
 - c.success = osucc\*w + success\*(1-w)
 - now-stamp>1min && swap(stamp now)
