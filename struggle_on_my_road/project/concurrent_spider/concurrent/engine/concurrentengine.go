package engine

import "log"

// 集合成员的作用 使 run 为具体的那个对象 因为 simple 中也定义 run 函数 主要作用是区分
type ConcurrentEngine struct {
	// 将接口作为子集使用
	Scheduler Scheduler
	// 控制 goroutine 数量
	WorkerCount int
}

// 定义接口 我们要用这个接口 具体怎么实现接口由 simple 实现
type Scheduler interface {
	// 接口 chan
	Submit(Request)
	// 将创建的 chan 赋值给该函数的成员变量
	ConfigureMasterWorkerChan(chan Request)
	// 接收每一次chan 只要我发送了chan 才能证明我想要 chan
	// WorkReady(w chan Request)
	// 执行队列函数
	// Run()
}

func (e *ConcurrentEngine) Run(seeds ...Request) {
	// 创建接收者
	in := make(chan Request)
	// 创建发送者
	out := make(chan ParseResult)

	// 将创建的 chan in传给 simple 包的 workerChan 因为 workerChan 只是定义 并没有创建, 这样在外面分配内存 传给结构体里面的属性
	e.Scheduler.ConfigureMasterWorkerChan(in)

	// 由 main 定义 goroutine 数量
	// 创建指定的协程 在协程里面 无限循环
	// 无限循环里从 in 管道里 取出 url, 请求 url 处理获得的数据 添加进 out 管道
	for i := 0; i < e.WorkerCount; i++ {
		// 将接收者 与 放者送入函数
		createWorker(in, out)
	}

	// 取出每一个 Request 放入 simple 中的 workerChan
	for _, r := range seeds {
		// 这一步 把传入的 Request 传入 Scheduler 里的 workerChan 管道，
		// 在 上面 e.Scheduler.ConfigureMasterWorkerChan(in) 这一步中 workerChan 管道跟 in 管段指向的是同一个内存地址
		// 所以 传入(存入) workerChan 就是 存入 in 管道
		e.Scheduler.Submit(r)
	}

	// 为每一位 成员添加序号 起始为 零
	itemCount := 0

	for {
		// chan 传输的值
		// createWorker 创建协程 处理 url 把处理过的数据 存入 out 管道中
		result := <- out
		// 获取 result 的值 注意 range 会判断 result 是否有内容 如果没有内容 不循环
		for _, item := range result.Items {
			// 将 item 中的数据打印出来 这里我们应该明白 item 就是存放我们需要的数据结构体变量
			log.Printf("Got item:#%d %v", itemCount, item)
			// 编号
			itemCount ++
		}

		// 提取地址与正则表达式函数
		for _, request := range result.Requests {
			// 将新的地址传入函数 在函数内的结构体进行 in 值的改变

			// 这一步 存入 in 管道中 createWorker 方法 再从管道中读取 url，请求url, 处理请求到的数据
			e.Scheduler.Submit(request)
		}
	}
}

// 此函数是并发的数据汇集区
func createWorker(in chan Request, out chan ParseResult) {
	// 直接 goroutine 无限接收数据并传送数据 外围 for 循环控制几个这样的 goroutine 同时并发
	go func(){
		// 无限接收数据靠 for 循环
		for {
			// 整个程序中比较绕的在这里
			// 将 in 的值赋值给 request 没问题，问题是同时并发并且无限接收数据的 in 是从哪里来的
			// 注意了 关键的是 in 有一个指针指向in 这个指针也在 goroutine 无限循环接收
			// 所以不要被绕糊涂 in 的值是从哪里来的 in 的值被指针指向了 ！！！
			request := <- in
			// 取得ParseResult 返回值
			result, err := worker(request)
			if err != nil {
				continue
			}
			// 将 返回值 放入 out chan 类型，类似于无限数组 可以放入无限的数据
			out <- result
		}
	}()
}
