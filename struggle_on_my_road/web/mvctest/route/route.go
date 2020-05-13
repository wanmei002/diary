package route

import . "mvctest/controller"

func init(){
	// 全局变量 保存静态资源路径
	StaticDir = make(map[string]string)
	// 注册路由
	p := &ControllerRegistor{}
	register(p)
	RouterRegister = p
	AutoRender = false

}

func register(p *ControllerRegistor){
	p.SetStaticPath("/img/", "E:/static/img/")
	p.Add("/user/:uid", &UserController{}, "Get")
}
