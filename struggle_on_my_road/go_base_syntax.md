## 语法第一
 - 值类型 和 引用类型   `变量的本质对一块内存空间的命名`
    + 基本数据类型 `int float bool string` 和 数组 struct 都是值类型
    + slice map chan 是引用类型
    + 空接口 interface{} 可以实现任何类型
        ```go
            var int1 int
            var bl bool
                bl = (1!=0)
            type Student struct{}
            var arr1 [5]string
            
        ```
    + new 声明值类型的指针类型  make 声明引用类型 并分配内存
        
 - int 类型
    + int 类型 默认是32位 或 64位，与具体平台有关
    
 - float 类型
    + float32 4字节 值范围 1.1755E-38 ~ 3.4028E38  -3.4028E38 ~ -1.1755E-38
        - sign 最高位表示符号域（bit31）1bit
        - exponent 8bit的指数域（bit23-bit30）8bit
        - fraction 23bit的小数域（bit0-bit22）23bit
        
    + float64 8字节  值范围 ~~~~
        - sign 最高位表示符号域（bit63）1bit
        - exponent 8bit的指数域（bit52-bit62）11bit
        - fraction 23bit的小数域（bit0-bit51）52bit
        
 - 字符类型
    + 获取字符串的某个字符 `ch := str[1]` 可以通过下标的方式获取，但是不能赋值
    + 字符串的相加 用 + , 如果换行 + 必须出现上一行末尾 
    + 字符串切片 `str_1 := str[:5]`
    + 汉字字符串遍历
        ```go
            str := "你好吗？ yes"
            slice1 := []rune(str)
            for _, v := range slice1 {
                fmt.Println(string(v))
            }
        ``` 
        
 - 类型转化
    + strconv.Itoa(int_val)  
    + strconv.Atoi(str_val)
    + 注意 字节数不匹配 会从后面开始截取
    
 - 数组
    + 声明 `var b = [3][3]int  var e = new([3]string)  a := [...]int{1,2,3}`
    
 - 切片
    + 创建切片有三种形式，基于数组、数组切片 和 直接创建  `make([]int, 5[, 10])`
    + 数组切片底层引用了一个数组，由三个部分构成：指针、长度和容量，指针指向数组起始下标，长度对应切片中元素的个数，容量则是切片起始位置到底层数组结尾的位置
    + append(slice, ...interface{})
    + copy(slice1, slice2)
    + slice2 := slice1\[:3\] // 动态删除
    + slice\[index\] = 123 // index 不能超过切片的长度, 只有 slice = append(slice, 123) 这样赋值才能超过切片的长度, 实际上已经不是同一个切片了 地址变了
    
 - 字典
    + var testMap map[string]int ; 字典是引用类型，声明后 使用前必须初始化 分配内存 
    + 初始化 
        ```
            // 初始化
            testMap = make(map[string]int)
            // 声明并且初始化
            testMap := map[string]int {
                "one" : 1,
                "two" : 2,  // 最后一个元素必须加上 ,
            }
        ```
    + 查找元素  `value, ok := testMap["one"]`
    + 删除元素 ` delete(testMap, "one") `
    
    
 - 函数
    + 变长参数  *实际底层是一个切片*
        ```go
           func myFunc(numbers ...int){
                for _, number := range numbers {
                    fmt.Println(number)
                }
           }
           slice1 := []int{1,2,3,4}
           myFunc(slice1...)
           // 任意类型的变长参数
           func myFunc(args ...interface{}){
                // 可以用反射机制来判断数据 和 操作数据
           }
        ```
        
        
 - 系统内置函数
    close
    len() 用于 字符串 数组 切片 字段 管道
    cap() 数组 切片 管道
    new() make()  // 用于分配内存 一个分配值类型 一个分配引用类型
    cap() append() // 用于切片
    panic recover // 错误处理机制
    print() println() 建议用 fmt 包
    ... 参考 built 包
    
    
 - struct 结构体
    + String() 方法 相当于 PHP中的 toString() 方法
    
    
 - for 循环
    + break 可以跳出循环体 当后面没有跟标签时跳出当前循环 (标签必须在循环语句上方)，有标签跳到 标签处 ，则循环不会第二次执行
    + for 条件体可以初始化语句，但是只能在for 模块中有效
    
 - goto 跳转
    + 跳转到指定的标签出，如果标签在 goto 上方，如果条件满足 还会继续执行 goto 跳转
    
 - switch 
    + case 结尾不用写 break ,如果要执行 下一个 case 用 fallthrough 会执行下一个 case
       
       
 - import 导入包
    + `.` 点操作的含义是这个包导入之后在你调用这个包的函数时可以省略前缀的包名
    + 别名操作 : 可以把包命名成另一个我们用起来容易记忆的名字
    + `_` 操作 : `_` 操作其实是引入该包，而不直接使用包里面的函数, 而是调用包里面的 `init` 函数
    ```golang
       import(
            . "fmt"   // 可以直接用 Println()函数 省略包名
            f "fmt"   // f.Println()
            _ "fmt"   // 调用 包里的 init 函数
        )
    ```
    
 - runtime 包
    + `Goexit` 退出当前运行的 goroutine, 但 defer 还是会继续调用
    + `Gosched` 让出 goroutine 的执行权限， 调度器安排其它等待的任务运行(一般单核会用)
    + `NumCPU` 返回 CPU 核心数
    + `NumGoroutine` 返回正在执行和排队的任务总数
    + `GOMAXPROCS` 用来设置可以运行的CPU核数
    
### import package 的流程
```go
  —> main
    import package1  —>       package1
                                  | 
	  const... <—|          const...
            |        |            |
	   var...    |           var...
            |        |            |
	  init()         ——————init()
            |
	  main()

```
