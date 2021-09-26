package main

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
)

// Go 自带的路由器有几个限制:
/*
1. 不支持参数设定，例如 /user/:uid 这种泛类型匹配 (url & 传参对 SEO 不友好, 对用户也不友好)
2. 无法很好的支持 REST 模式 ):_: 我在工作中也只是用到 get post 。。。
3. 一般网站的路由规则太多了，编写繁琐  ( 新手的我无法理解大佬的极致精神
end : 让我跟着大佬的思想设计一个 route register 实现 ，看招~~~
*/

// 储存单元路由的信息
type controllerInfo struct {
	regex  *regexp.Regexp  //保存正则路由
	params map[int]string  // 路由上的单元正则匹配
	controllerType reflect.Type  // 路由对应的控制器
}

type App struct {
	Name string
}

// 路由注册到对应的方法上
type ControllerRegistor struct {
	routers []*controllerInfo  // 保存所有路由的信息
	Application *App  // 应用基本信息
}

// 公共的 Controller 方法 所有的 控制器 都要继承这个 Controller
type ControllerInterface interface {

}

// 测试控制器
type UserController struct {
	ControllerInterface
}
func (u *UserController) String(){
	fmt.Println("i am UserController")
}
func (u *UserController) Get(){
	fmt.Println("user get request")
}

// 添加路由
func (p *ControllerRegistor) Add(pattern string, c ControllerInterface) {
	parts := strings.Split(pattern, "/")  // 根据 / 切割 url
	j := 0
	params := make(map[int]string)  // 保存匹配到的URL参数 :id 之类
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {  // 匹配到 url 上的参数了
			expr := "([^/]+)"  // 匹配所有 的参数值
			// 匹配用户自定义的正则
			if index := strings.Index(part, "("); index != -1 {  // 匹配到了用户自定义的正则
				expr = part[index:]
				part = part[:index]
			}
			params[j] = part  // 保存正则匹配到的键  比如 url 上的 :id 作为键
			parts[i] = expr   // 把 url 上的 :id 参数 替换成正则表达式
			j++
		}
	}

	pattern = strings.Join(parts, "/")  // 用 / 把切片组合起来
	regex, regexErr := regexp.Compile(pattern)  // 校验正则 返回正则实例
	if regexErr != nil {
		panic(regexErr)  // 简单写法 可以根据自己的需求定制
		return
	}
	// 现在开始创建路由
	t := reflect.Indirect(reflect.ValueOf(c)).Type()   // 根据值 找到对应的 类型 也就是 对应的控制器
	fmt.Println("url 匹配到的 controller : ", t)
	route := &controllerInfo{}
	route.regex = regex
	route.params = params
	route.controllerType = t
	p.routers = append(p.routers, route)
}

// 设置静态路由
var StaticDir map[string]string
func init(){
	StaticDir = make(map[string]string)
}
func (p *ControllerRegistor) SetStaticPath(url string, path string) {
	StaticDir[url] = path
}

// 转发路由
func (p *ControllerRegistor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func(){
		if err := recover(); err != nil {

			panic(err)  // 暴力输出 测试
		}
	}()

	var started bool
	for prefix, staticDir := range StaticDir {
		if strings.HasPrefix(r.URL.Path, prefix) {  // 匹配到静态路由
			file := staticDir + r.URL.Path[len(prefix):]  // 重写静态路由
			http.ServeFile(w, r, file)    // 调用 go 自带的静态文件处理函数 返回静态资源
			started = true
			return
		}
	}

	requestPath := r.URL.Path
	// 开始比对路由
	for _, route := range p.routers {
		if !route.regex.MatchString(requestPath) {
			continue
		}

		// 正则匹配路由 并获取 url 上带的参数
		matches := route.regex.FindStringSubmatch(requestPath)
		if len(matches[0]) != len(requestPath) {
			continue
		}
		params := make(map[string]string)
		fmt.Println("start r.URL.RawQuery : ", r.URL.RawQuery)
		if len(route.params) > 0 {  // 说明 url 上有参数
			values := r.URL.Query()   // 获取url get 传参
			fmt.Println("start r.URL.Query : ", values)
			for i, match := range matches[1:] {
				fmt.Println("matches[",i,"] : ", route.params[i])
				values.Add(route.params[i], match)
				fmt.Println("params[",i,"] added , r.URL.Query : ", values)
				params[route.params[i]] = match
			}
			r.URL.RawQuery = url.Values(values).Encode() + "&" + r.URL.RawQuery  // 添加参数到 原始url 上
			fmt.Println("changed r.URL.RawQuery : ", r.URL.RawQuery)
		}

		// 利用反射机制 获取 绑定的 控制器
		vc := reflect.New(route.controllerType)
		method := vc.MethodByName("Miss")
		//init := vc.MethodByName("init")
		if r.Method == "GET" {
			method = vc.MethodByName("Get")
		} else if r.Method == "POST" {
			method = vc.MethodByName("Post")
		} else {  // 其它 RESTFUL 风格就不写了 如果有需要请读者 自己 else if 下去

		}
		in := make([]reflect.Value, 0)
		method.Call(in)
		started = true
		break
	}

	if started == false {
		http.NotFound(w, r)
	}

}

func main(){
	p := &ControllerRegistor{}
	p.SetStaticPath("/img/", "/static/img/")
	p.Add("/user/:uid([0-9]+)", &UserController{})
	fmt.Println("regex : ",p.routers[0].regex)
	fmt.Println("params : ",p.routers[0].params)
	fmt.Println("start listen")
	err := http.ListenAndServe("127.0.0.1:9999", p)
	if err != nil {
		fmt.Println("listen failed : ", err)
	}

}
