package fetcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Fetch 请求 URL 返回页面源码
var rateLimiter = time.Tick(10 *time.Millisecond)
func Fetch(url string) ([]byte, error) {
	// 增加定时器 每隔 10毫秒 请求一次 url 防止被封 也可以用 ip池
	<- rateLimiter
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code is error, ret code is %d", resp.StatusCode)
	}
	// 读取获得的内容
	return ioutil.ReadAll(resp.Body)
}
