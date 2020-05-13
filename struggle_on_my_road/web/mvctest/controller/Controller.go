package controller

import (
	"errors"
	"html/template"
	"net/http"
)

type ControllerInterface interface {
	Init(ct *Context, cn string)				// 初始化
	Prepare()										// 开始执行前的一些处理
	Get()											// method = GET 的处理 下面的一样
	Post()
	Delete()
	Put()
	Head()
	Patch()
	Options()
	Finish()										// 执行完之后的处理
	Render() error									// 执行完 method 对应的方法之后渲染页面
}

type Context struct {
	ResponseWriter http.ResponseWriter
	Request *http.Request
	Params map[string]string
}

// 基础类 实现接口
type Controller struct {
	Ct 			*Context
	Tpl			*template.Template
	Data 		map[interface{}]interface{}
	ChildName 	string
	TplNames 	string
	Layout []	string
	TplExt 		string
}

func (c *Controller) Init(ct *Context, cn string) {
	c.Data = make(map[interface{}]interface{})
	c.Layout = make([]string, 0)
	c.TplNames = ""
	c.ChildName = cn
	c.Ct = ct
	c.TplExt = "tpl"
}

func (c *Controller) Prepare(){

}

func (c *Controller) Finish(){

}

func (c *Controller) Get(){
	http.Error(c.Ct.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Post(){
	http.Error(c.Ct.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Delete(){
	http.Error(c.Ct.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Put(){
	http.Error(c.Ct.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Head(){
	http.Error(c.Ct.ResponseWriter, "Method Not Allowed", 405)
}
func (c *Controller) Patch(){
	http.Error(c.Ct.ResponseWriter, "Method Not Allowed", 405)
}
func (c *Controller) Options(){
	http.Error(c.Ct.ResponseWriter, "Method Not Allowed", 405)
}
func (c *Controller) Render() error {
	//http.Error(c.Ct.ResponseWriter, "Method Not Allowed", 405)
	return errors.New("Method Not Allowed")
}
