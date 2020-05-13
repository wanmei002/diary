package main

import (
	"fmt"
	"mvctest/route"
	"net/http"
)
func main(){
	// 开始实现我的 V 模式
	p := route.RouterRegister
	fmt.Println("start listen")
	err := http.ListenAndServe("127.0.0.1:9999", p)
	if err != nil {
		fmt.Println("listen failed : ", err)
	}
}