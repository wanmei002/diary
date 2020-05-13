```
package main

import "fmt"

func calc(intChan chan int, resChan chan int, exitChan chan bool){
	fmt.Println("calc start")
	count := 0
	for v := range intChan {
		flag := true
		for i := 2; i < v; i++ {
			if v % i == 0 {
				flag = false
				break
			}
		}

		if flag {
			resChan <- v
		}
		count++
	}

	fmt.Println("calc 运行了 ", count, " 次")

	exitChan <- true
}

func main(){
	intChan := make(chan int, 1000)
	resChan := make(chan int, 1000)
	exitChan := make(chan bool, 8)   // 检查其它 goroute 是否结束运行

	for i := 0; i < 1000; i++ {
		intChan <- i
	}

	close(intChan)   // 如果要用 range 遍历 channel ， 则必须先 关闭 

	for i := 0; i < 8; i++ {
		go calc(intChan, resChan, exitChan)
	}

	for i := 0; i < 8; i++ {
		<- exitChan  // 通过阻塞的方式 等待上面的 协程都运行完
	}

	close(resChan)

	for v := range resChan {
		fmt.Println("resChan 里的数 :", v)
	}
}

```