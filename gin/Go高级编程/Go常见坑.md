##### Go常见的坑
 - 独占CPU导致其它 Goroutine 饿死
 ```go
 func main() {
 	runtime.GOMAXPROCS(1)
 	
 	go func(){
 		for i:=0; i < 10; i++ {
 			fmt.Println(i)
 		}
 	}()
 	// for{} 占用CPU 这样其它Goroutine 就得不到执行
 	// 要想避免这种情况 最好 for { runtime.Gosched() }
 	// 或者 for 换成 select
 	for {} 
 }
 ``` 
 
 - 在循环内部执行 defer 语句
 ```go
 func main(){
 	for i:=0; i<5; i++ {
 		f,err := os.Open("/path/to/file")
 		if err != nil {
 			fmt.Println(err)
 		}
 		// defer 在函数退出时才执行 在 for 执行 defer 会导致资源延迟释放
 		defer f.Close()
 	}
 }
 // 解决的方法可以在 for 中构造一个局部函数, 在局部函数内部执行 defer :
 func main(){
 	for i:=0; i<5; i++ {
 		func(){
 			f, err := os.Open("/path/to/file")
 			if err != nil {
 				fmt.Println(err)
 			}
 			defer f.Close()
 		}()
 	}
 }
 ```
 
 - 切片会导致整个底层数组被锁定，底层数组无法释放内存，如果底层数组较大会对内存产生很大的压力
 
 - 禁止 main 函数退出的方法
 ```go
 func main() {
 	defer func(){ for{} }()
 	// 或者下面这个
 	defer func(){ select {} }()
 	// 或者下面这个
 	defer func(){ <-make(chan bool) }()
 }
 ```
 
 
 
 
 
 