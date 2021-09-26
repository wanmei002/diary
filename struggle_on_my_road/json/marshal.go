package main

import (
	"encoding/json"
	"fmt"
)
// JSON 转化前后的数据类型映射如下:
// 1. 布尔值转化为JSON后还是布尔类型
// 2. 浮点数和整型会被转化为 JSON 里边的常规数字
// 3. 字符串将以UTF-8编码转化输出为 Unicode 字符集的字符串, 特殊字符比如汉字将会被转义为 \u003c
// 4. 数组和切片会转化为 JSON 里边的数组, 但 `[]byte` 类型的值将会被转化为 Base64 编码后的字符串, slice 类型的零值会被转化为 `null`
// 5. 结构体会转化为 JSON 对象，并且只有结构体里面以大写字母开头的可被导出的字段才会被转化输出, 而这些可导出的字段会作为 JSON 对象的字符串索引
// 6. 转化一个 map 类型的数据结构时，该数据的类型必须是 `map[string]T` (T可以是 `encoding/json` 包支持的任意数据类型。)
type Person struct {
	Name string
	Age int
}
// 如果json 格式解析成结构体, 所对应的值会自动填充到结构体的目标字段上， json.Unmarshal() 将会遵循如下循序进行查找匹配:
// 1. 一个包含 对应标签的字段 `json:"name"`(不区分大小写)
// 2. 一个名为 对应字段名 或者除了首字母其它字母不区分大小写的 对应的字段(首字母必须大写)
// 3. 如果 json 里的字段 在结构体里找不到对应的字段，则会舍弃这个字段
// 4. 对于 JSON 中没有而 结构体 中定义的字段，会以对应数据类型的默认值填充
type Classes struct {
	Name string
	Info []Person

}

func main(){
	zyn := Person{
		Name : "zyn",
		Age:3,
	}
	slie := []Person{zyn}
	yey := Classes{
		Name:"lyl",
		Info:slie,
	}
	data, err := json.Marshal(yey)
	if err != nil {
		fmt.Println(err)
		fmt.Println("转换 json 失败")
		return
	}
	fmt.Println("转换的json : ", string(data))
	var cla Classes
	err = json.Unmarshal(data, &cla)
	if err != nil {
		fmt.Println(err)
		fmt.Println("json to obj error")
		return
	}
	fmt.Println("cla is : ",cla)
	person1 := cla.Info[0]
	fmt.Printf(
		"cla[name]:%s\n, cla[Info]:%v\n, cla[Info][0][Name]:%v\n",
		cla.Name, cla.Info, person1.Name)

}
