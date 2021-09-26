package engine

import (
	"fmt"
	"log"
	"spider/fetcher"
)

// 引擎 控制整个程序的流程
func Run(seeds ...Request) {
	var requests []Request
	// 接收 main 函数传过来的值
	for _, r := range seeds {
		requests = append(requests, r)
	}

	// 利用传递过来的值进行 解析 及提取
	for len(requests) > 0 {
		// 获取第一个值
		r := requests[0]
		// 进行切片 把已经提取的内容筛选出去
		requests = requests[1:]
		// 接收 worker 函数的值
		ParseResult, err := worker(r)

		// 如果出现错误 结束本次循环
		if err != nil {
			continue
		}

		// requests被填满 requests 又得到新的 URL 和 运算函数 被抓取信息只要足够就可以一直运行下去
		requests = append(requests, ParseResult.Requests...)
		// 打印所有在PrintCityList 函数返回的Item值 Item值是任何类型可以使城市名也可以是用户信息
		for _, item := range ParseResult.Items {
			fmt.Printf("Got item %v\n", item)
		}
	}
}

func worker(r Request) (ParseResult, error) {
	// 第一次打印为 main 函数传入地址 然后每次打印是从 r.ParserFunc 函数中提取出的城市地址
	log.Printf("Fetching %s\n", r.Url)
	// 将不同URL传输进去 返回不同的页面源代码
	body, err := fetcher.Fetch(r.Url)
	// 判断 URL 是否正确 如果不正确 跳过此次循环
	if err != nil {
		log.Printf("Fetcher: err fetching url %s:%v", r.Url, err)
		return ParseResult{}, err
	}
	return r.ParseFunc(body), nil
}
