package scheduler

import "spider/engine"

type SimpleScheduler struct {
	// chan 类型 可以 无限制的存放数据
	workerChan chan engine.Request
}

// 注意是指针类型 c 是创建的 chan 赋值给workerChan 后 workerChan 可以使用
func (s *SimpleScheduler) ConfigureMasterWorkerChan(c chan engine.Request){
	s.workerChan = c
}

// 将数据放入 chan 中 无限存放数据
func (s *SimpleScheduler) Submit(r engine.Request) {
	// 这里开并发的原因是 另一头想要执行 必须要接受 chan 而这边必须要将有效 chan 值发送出去次才能继续
	// 两头一个要接受值 但是没有值可接收 这个是要发送值没有值可发送
	// 直接就死程序了 所以要并发 我们还有另一个解决方法 这里先使用这个
	go func(){
		s.workerChan <- r
	}()
}

