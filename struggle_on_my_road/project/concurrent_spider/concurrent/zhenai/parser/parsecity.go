package parser

import (
	"regexp"
	"spider/engine"
)

// 获取用户信息 URL 格式
const cityRe  = `<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`

// contents 为城市页面地址 从每个城市第一页中筛选信息
func ParseCity(contents []byte) engine.ParseResult {
	// 确定要查找的格式
	re, err := regexp.Compile(cityRe)
	if err != nil {
		panic("ParseCity 解析user_url err : " + err.Error())
	}
	// 搜索全部的格式相同的信息
	matches := re.FindAllSubmatch(contents, -1)
	// 创建结构体进行存放
	result := engine.ParseResult{}
	// 把第一张页面中的用户名及地址取出
	for _, m := range matches {
		// m[2] 为用户名
		name := string(m[2])
		// 在结构体中存入所有昵称名字 并标识为 User
		result.Items = append(result.Items, "User " + name)
		// 这个函数的中最关键的点
		// 将函数 ParseProfile 作为返回值 即确定了昵称 又没有改动结构体
		result.Requests = append(result.Requests, engine.Request{
			string(m[1]),
			func(c []byte) engine.ParseResult {
				return ParseProfile(c, name)
			},
		})

	}
	return result
}