package main

import (
	"spider/engine"
	"spider/scheduler"
	"spider/zhenai/parser"
)

func main(){
	// 运行爬虫的起始条件 当内部循环一次结束后 里面信息属于无效
	e := engine.ConcurrentEngine{
		// 初始
		Scheduler: &scheduler.SimpleScheduler{},
		// 设定并发数量
		WorkerCount: 2,
	}

	// 由于有两个run 这里我们直接区分开
	e.Run(engine.Request{
		Url:"http://www.zhenai.com/zhenghun",
		ParseFunc: parser.PrintCityList,
	})
}
