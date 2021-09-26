package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const FORM = `
<html>
<body>
	<form action="#" method="post">
		<input type="text" name="name" />
		<input type="text" name="age" />
		<input type="submit" value="提交" />
	</form>
</body>
</html>
`

func SimpleServer(w http.ResponseWriter, request *http.Request) {
	gets := request.URL.Query()
	fmt.Println(gets)
	fmt.Printf("%T", gets.Get("age"))
	//_, _ = io.WriteString(w, "<h1>hello, world</h1>")
	fmt.Fprintf(w, `{"code":0}`)
}

func GetRequest(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	switch request.Method {
	case "GET" :
		_, _ = io.WriteString(w, FORM)
	case "POST" :
		_ = request.ParseForm()
		io.WriteString(w, request.Form["name"][0])
		io.WriteString(w, "<br>")
		io.WriteString(w, request.FormValue("age"))
	}
}

func JsonPostServer(w http.ResponseWriter, request *http.Request) {
	headerType := request.Header.Get("Content-Type")
	fmt.Println("header-type : ", headerType)
	decoder := json.NewDecoder(request.Body)
	var params map[string]interface{}
	decoder.Decode(&params)
	fmt.Println("获取到的数据 : ", params)
	fmt.Fprintf(w, `{"code":0}`)
}

func FormPostServer(w http.ResponseWriter, request *http.Request) {
	headerType := request.Header.Get("Content-Type")
	fmt.Println("header-type : ", headerType)
	username := request.Form.Get("name")
	userage := request.Form.Get("age")
	fmt.Printf("name=%s, age=%s", username, userage)
	fmt.Fprintf(w, `{"code":0}`)
}

type Person struct {
	Name string `json:"name"`
	Age int		`json:"age"`
}

type Response struct {
	Code int		`json:"code"`
	Msg string		`json:"msg"`
	Info []Person	`json:"info"`
}

func RetJsonServer(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	zyn := Person{
		Name:"zyn",
		Age:3,
	}
	sli := []Person{zyn}
	ret := Response{
		Code:0,
		Msg:"success",
		Info:sli,
	}
	json.NewEncoder(w).Encode(ret)
}

func main(){
	http.HandleFunc("/test1", SimpleServer)
	http.HandleFunc("/test2", GetRequest)
	http.HandleFunc("/test3", JsonPostServer)
	http.HandleFunc("/test4", FormPostServer)
	http.HandleFunc("/test5", RetJsonServer)
	fmt.Println("开始监听接口")
	if err := http.ListenAndServe("127.0.0.1:8099", nil); err != nil {
		fmt.Println(err)
		fmt.Println("链接失败")
	}
	fmt.Println("监听结束")
}