/**
 *我觉得这种方法适合http 接口
 */
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

func main(){
	resp, err := http.Get("http://t.djcapp.game.qq.com/daoju/igw/main/?_service=app.extpackage.list&acctype=guest&_biz_code=cjm")
	if err != nil {
		fmt.Println(err)
		fmt.Println("http.Get error")
		return
	}
	fmt.Println(resp)
	//data, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Println(err)
	//	fmt.Println("curl 请求失败")
	//	return
	//}
	//fmt.Println("Body is :",string(data))

	var json_map map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&json_map)
	for k, v := range json_map {
		fmt.Printf("json_map[%v] : %v, type - '%v'\n", k, v, reflect.TypeOf(v))
	}

	for k, v := range json_map["data"].([]interface{}) {
		a := v.(map[string]interface{})
		for kk, vv := range a {
			fmt.Printf("data[%d]['%s']:%v\n", k, kk, vv)
		}

	}
}
