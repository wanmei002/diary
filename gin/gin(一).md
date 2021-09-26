#### 安装
 - `go get -u github.com/gin-gonic/gin`
 - 有依赖包依次是:
 	+ github.com/golang/protobuf
    + github.com/modern-go/concurrent
    + github.com/modern-go/reflect2
    + golang.org/x/sys
    + golang.org/x/text #windows需要

#### hello world
 ```go
  package main
  import "github.com/gin-gonic/gin"
  func main() {
  	r := gin.Default()
  	r.GET("ping", func(c *gin.Context){
  		c.String(http.StatusOK, "hello world")
    })
    r.Run()
  }
 ```

#### http RUSTFULL 请求风格
 - HEAD 只请求页面的首部
 - PUT 从客户端向服务器传送的数据取代指定的文档的内容
 - OPTIONS 允许客户端查看服务器的性能
 - Any 匹配 gin 支持的所有状态请求

#### 路由匹配
 - /user/:name 可以匹配/user/xxx  但是匹配不上 /user/xxxx/
 - /user/:name/\*action 可以匹配到 /user/zzh/18 如果没有上面那条路由会匹配到 /user/zzh 并重定向到 /user/zzh/ \*action 会连 `/` 也会匹配上的

#### 参数获取
 - url 上正则匹配上的参数 用 gin.Context.Param("name") 获取
 - querystring 参数(?后面的参数) 用 gin.Context.DefaultQuery("key", "defaultValue")  gin.Context.Query("key") 获取值，第一个如果值不存在 则设置默认值
 - 获取 post 参数 gin.Context.PostForm("key")
 - 获取 数组 传参 ids[a]=a&ids[b]=b
 	+ gin.Context.QueryMap("key") 获取 url 上的map参数
 	+ gin.Context.PostFormMap("key") 获取 post 传参的map参数
 
#### 上传文件
 - 单文件上传
 ```go
  file, _ := gin.Context.FormFile("input_file_name")
  // file.Filename 文件名称

  // 保存文件
  gin.Context.SaveUploadFile(file, "/path/filename")
 ```
 - 批量文件上传
 ```go
  form, _ := gin.Context.MultipartForm()
  files := form.File["upload[]"]
  for _, file := range files {
  	gin.Context.SaveUploadFile(file, "/path/filename")
  }
 ```

#### 路由组
 ```
  router := gin.Default()
  v1 := router.Group("/v1")
  // 花括号 代表了 块作用域 作用域内的变量名如果跟域外的变量名重复，会替代域外的变量
  {
  	v1.POST("login", loginV1)
  }
 ```

#### 空白的 gin
 - gin.Default() 是有默认中间件的 比如 gin.logger() gin.Recovery()
 - `r := gin.New()` 是空白的 没有中间件, 此时如果想加载中间件 使用 `r.Use(gin.logger())`

#### 怎样把日志写入文件中
 - 直接上代码
 ```go
 // 字体不用加颜色
  gin.DisableConsoleColor()
 // 创建日志文件
  f, _ := os.Create("gin.log")
  // MultiWriter 函数的作用是 把多个 io.Writer 合并成一个 io.Writer ,可以同时写入多个实现了 io.Writer 的文件中
  gin.DefaultWriter = io.MultiWriter(f, os.Stdout) // 这个可以同时写入文件和控制台中
 ```

#### 自定义日志格式
 - 直接上代码
 ```go
  router := gin.New()
  // 开始自定义日志格式
  router.Use(gin.LoggerWithFormatter(func (param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC3339),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
	}))
 ```

#### 参数绑定 并 校验
 - 直接上代码 , 在代码中解释
 ```go
 // tag 说明 binding:"required" 有这个tag 的, 则是必传字段
 // form:"form_user" 如果是 form 表单提交, 则有个字段必须是 form_user 的值才能绑定到 Login.User 属性上
 // User Password 都是必须的 则上传的参数中 必须有对应的字段
  type Login struct {
  	User 		string `form:"form_user" json:"json_user" binding:"required"`
  	Password	string `form:"form_psd" json:"json_psd" binding:"required"`
  }

  func main(){
  	r := gin.Default()
  	r.POST("/login", func(c *gin.Context){
  		var json Login
  		// 如果是 form 表单提交 则是 ShouldBind
  		if err := c.ShouldBindJSON(&json); err != nil {
  			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
  			return
  		}

  		if json.User != "zzh" || json.Password != "123456" {
  			c.JSON(http.StatusUnauthorized, gin.H{"status":"unauthorized"})
  			return
  		}

  		c.JSON(http.StatusOK, gin.H{"status":"you are logged in"})
	})

	r.Run()
  }
 ```

#### 自定义验证规则
 - 直接上代码验证
 ```go
 // 创建结构体
 type Booking struct {
	CheckIn time.Time `form:"check_in" binding:"required" time_format:"2006-01-02"`
	// binding tag 中 bookabledate 函数是验证这个字段的函数, gtfield 字段不知道是干什么的 去掉也可以正常运行, 如有知道的大佬请回复
	CheckOut time.Time `form:"check_out" binding:"required,gtfield=CheckIn,bookabledate" time_format:"2006-01-02"`
 }

 // 自定义的验证函数
 var bookableDate validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	fmt.Println("date : ", date)
	if ok {
		today := time.Now()
		fmt.Println("today : ",today)
		if today.After(date) {
			return false
		}
	}
	return true
}

 func main(){
 	r := gin.Default()
 	// 注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 把 bookableDate 函数 注册到验证器里 验证器里的名字为 bookabledate 跟上面的 booking 结构体中 的 binding tag 里的 函数名称一致
		err := v.RegisterValidation("bookabledate", bookableDate)
		if err != nil {
			fmt.Println("bind func error : ", err)
		}
	}
	// 注册路由
	r.GET("/book", getBookable)
	r.Run()
 }

 // 路由处理函数
 func getBookable(c *gin.Context) {
	var b Booking
	if err := c.ShouldBindWith(&b, binding.Query); err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
 ```

#### 参数绑定
 - 以下方法 在 路由处理函数中  c 是 \*gin.Context 的实例
 - only bind query string 就是只绑定 url 上 `?` 后面的参数
 ```go
  type Person struct {
  	Name string `form:"name"`
  	Address string `form:"address"`
  }

 	func startBind(c *gin.Context) {
 		var person Person
 		// ShouldBindQuery 方法仅绑定 query string
 		if c.ShouldBindQuery(&person) == nil {

 		}
 	}

 ```

 - bind query string or post data 绑定? 后面的参数 也可以绑定 post 传参
 	+ ShouldBind(&person)

 - Bind Uri use method : ShouldBindUri
 	```go
 		type Person struct {
 			ID string `uri:"id" binding:"required"`
 			Name string `uri:"name" binding:"required"`
 		}
 		// main 函数中注册路由
 		router.GET("/user/:name/:id", startBind)
 		// 在 路由处理函数中
 		c.ShouldBindUri(&person)
 	```

 - Bind Header
 	```go
    	type testHeader struct {
    		Rate	Int 	`header:"Rate"`
    		Domain	string  `header:"Domain"`
    	}
    	// 路由处理函数中
    	c.ShouldBindHeader(&testHeader)
 	```

 - Bind HTML CheckBoxes 这个意思应该是绑定 form 传参的 数组 绑定
 	```go
 		type myForm struct {
 			Colors []string `form:"colors[]"`
 		}
 		// 处理 路由 的函数中
 		c.ShouldBind(&myForm)
 		// HTML 部分
 		<input type="checkbox" name="colors[]" value="red" />
 		<input type="checkbox" name="colors[]" value="blue" />
 	```

 - Multipart/Urlencoded binding 上传文件绑定
 ```go
 type profileForm struct {
	Name string 	`form:"name" binding:"required"`
	Avatar *multipart.FileHeader `form:"avatar" binding:"required"`
 }
 // 路由处理函数中
 c.ShouldBind(&profileForm)  // form 表单中的数据绑定到 profileForm 结构体实例中
 // 保存文件到指定的位置
 c.SaveUploadFile(&profileForm.Avatar, "文件路径" + "文件名")
 ```

 - XML JSON YAML and ProtoBuf AsciiJSON rendering 数据渲染 c 是 \*gin.Context
 	+ json 输出
 	`c.JSON(http.StatusOK, gin.H{"message": "hello world"})`
 	`c.JSON(http.StatusOK, msg)` // msg 是一个结构体 
 	+ XML 输出
 	`c.XML(http.StatusOK, gin.H{"message": "hello world"})`
 	+ YAML 输出
 	`c.YAML(http.StatusOK, gin.H{"message": "hello world"})`
 	+ ProtoBuf 原型输出 (不知道什么意思 如果有哪个大佬知道请留言指教)

 	+ ascii json
 	`c.AsciiJSON(http.StatusOK, gin.H{"lang":"Go语言"})` // 把中文变成 ascii 码, c.JSON不会 还是中文
 	+ pure json 普通的 json 输出会把 html 标签变成实体字符 PureJSON 还是原样输出
 	`c.PureJSON(200, gin.H{"html":"<b>hello world!</b>"})`


 - 静态文件
 	+ r.Static("/static", "F:\\static") // 本地文件系统
 	+ r.StaticFs("/more_static", http.Dir(\`F:\\static\`)) 生成处理网络请求的文件系统 这个更好点
 	+ r.StaticFile("/favicon.ico", \`F:\\favicon.ico\`)

 - 从文件中读取数据到请求中
 	```go
 		r.GET("/local/file", func(c *gin.Context){
 			c.File(`F:\txt.log`)
 		})
 	```

 - 接下来是自定义模板 就不做介绍了
 
 - 重定向
 	```go
 	r.GET("/redirect", func (c *gin.Context){
		c.Redirect(http.StatusMovedPermanently, "http://www.baidu.com/")
	})
 	```
 - 本地重定向
 	```go
 		r.GET("/red_local", func(c *gin.Context){
			c.Request.URL.Path = `/local/file`
			r.HandleContext(c)
		})
 	```

 - 中间件
 	+ 直接上代码
 	```go
	 	func Logger() gin.HandlerFunc {// 返回一个参数是 gin.Context 的函数
			return func(c *gin.Context) {
				fmt.Println("before")
				c.Next()
				fmt.Println("next")
			}
		}
 	```

 - 运用基本的权限验证
  + 直接在代码中验证
  ```go
  // 先设置一个全局变量
  // secrets 是一个 map[string]interface{}  键是用户 值是用的信息
  var secrets = gin.H{
  "foo" : gin.H{"email": "foo@bar.com", "phone": "123433"},
  "austin" : gin.H{"email": "austin@example.com", "phone": "666"},
  "lena" : gin.H{"email":"lena@guapa.com", "phone": "523443"},
  }

  // 在 main 函数中
  authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
    "foo": "bar",// foo 是用户名 bar 是密码
    "austin": "1234",
    "lena": "hello2",
    "manu": "4321",
  }))

  authorized.GET("/secrets", func(c *gin.Context){
    user := c.MustGet(gin.AuthUserKey).(string)
    if secret, ok := secrets[user]; ok {
      c.JSON(http.StatusOK, gin.H{"user":user, "secret":secret})
    } else {
      c.JSON(http.StatusOK, gin.H{"user":user, "secret": "NO SECRET :("})
    }
  })
  ```
  > 在浏览器里 输入 localhost:8080/admin/secrets 会弹出一个输入框让你输入 用户名和密码
  > 如果要用 curl 请求 请添加 header 信息 -H'Authorization:Basic YXVzdGluOjEyMzQ'  注意 `YXVzdGluOjEyMzQ=` 是 austin:1234 的 base64 编码


 - goroutines 中 添加中间件
    + 不能使用初始的 Context 必须拷贝一份
    ```go
    r.GET("/long_async", func(c *gin.Context){
        cCp := c.Copy() // 使用自带的方法拷贝一份
        go func(){
          time.Sleep(5 * time.Second)
          fmt.Println("Done! in path " + cCp.Request.URL.Path)  
        }()
    })
    ```

 - Custom HTTP configuration 自定义 http 请求
 ```go
 func main(){
  router := gin.Default()
  http.ListenAndServe(":8080", r)
 }

 // 或者用下面这种方式
 func main(){
  router := gin.Default()

  s := &http.Server{
    Addr :          ":8080",
    Handler:        router,
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second.
    MaxHeaderBytes: 1 << 20,
  }

  s.ListenAndServe()
 }
 ```

 - gin 开启多个服务
 ```go
    package main

    import (
      "github.com/gin-gonic/gin"
      "golang.org/x/sync/errgroup"
      "log"
      "net/http"
      "time"
    )

    var (
      g errgroup.Group
    )

    func router01() http.Handler {
      e := gin.New()
      e.Use(gin.Recovery())
      e.GET("/", func(c *gin.Context){
        c.JSON(http.StatusOK, gin.H{
          "code" : http.StatusOK,
          "error" : "welcome server 01",
        })
      })
      return e
    }

    func router02() http.Handler {
      e := gin.New()
      e.Use(gin.Recovery())
      e.GET("/", func(c *gin.Context){
        c.JSON(200, gin.H{
          "code" : 200,
          "error" : "welcome server 02",
        })
      })
      return e
    }

    func main(){
      server01 := &http.Server{
        Addr:       ":8082",
        Handler:    router01(),
        ReadTimeout:  5*time.Second,
        WriteTimeout: 10*time.Second,
      }

      // 第二个服务器
      server02 := &http.Server{
        Addr:       ":8081",
        Handler:    router02(),
        ReadTimeout:  5*time.Second,
        WriteTimeout:   10*time.Second,
      }

      g.Go(func() error {
        err := server01.ListenAndServe()
        if err != nil && err != http.ErrServerClosed {
          log.Fatal(err)
        }
        return err
      })

      g.Go(func() error {
        err := server02.ListenAndServe()
        if err != nil && err != http.ErrServerClosed {
          log.Fatal(err)
        }
        return err
      })

      if err := g.Wait(); err != nil {
        log.Fatal(err)
      }

    }
 ```

 - 绑定 form 表单数据 到自定义结构
    ```go
    type StructA struct {
        FieldA string `form:"field_a"`
    }
    
    type StructB struct {
        NestedStruct StructA
        FieldB string `form:"field_b"`
    }
    func GetBataB(c *gin.Context) {
    	var b StructB
    	err := c.Bind(&b)
    	if err != nil {
    		fmt.Println("bind B error : ", err)
    	}
    	c.JSON(200, gin.H{
    		"a" : b.NestedStruct,
    		"b" : b.FieldB,
    	})
    }
    func main(){
    	r := gin.Default()
    
    	r.GET("/getb", GetBataB)
    	r.Run(":8088")
    
    }
    ```
    > localhost:8088/getb?field_a=1&field_b=2  return : {"a":{"FieldA":"1"},"b":"2"}
    
 - 尝试绑定 request.Body 进不同的结构体
    + `c.ShouldBind` 方法只能用于一次绑定 是消耗品
    + `c.ShouldBindBodyWith(&objA, binding.JSON)` 这个方法可以用于多次绑定
    + 如果只用于一次绑定 最好用 `ShouldBind` 
 