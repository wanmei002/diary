package controller

import "fmt"

type UserController struct {
	Controller
}

func (this *UserController) Get(){
	_, _ = fmt.Fprintf(this.Ct.ResponseWriter, "this is UserController Get Method")
}

func (this *UserController) Post(){
	_, _ = fmt.Fprintf(this.Ct.ResponseWriter, "this is UserController Post Method")
}
