package fetcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

//@param url
//@return []byte error
func Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code is error , ret code is %d", http.StatusOK)
	}
	// 读取所有的内容

	return ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(err)
	// }
	// cityAll := printCityAll(res)
	// for _, val := range cityAll {
	// 	fmt.Println(val)
	// }
}

func printCityAll(contents []byte) [][]string {
	reg, err := regexp.Compile(`<a href="(http://www.zhenai.com/zhenghun/[a-z0-9]*)"[^>]*>([^<]+)</a>`)
	if err != nil {
		panic(err)
	}
	// regMatch := reg.FindAllString(string(contents), -1)
	regMath := reg.FindAllSubmatch(contents, -1)
	var regMathString [][]string
	regMathString = make([][]string, len(regMath))
	for i := range regMathString {
		regMathString[i] = make([]string, 2)
	}
	for n, val := range regMath {
		for i, val1 := range val {
			if i > 0 {
				regMathString[n][i-1] = string(val1)
			}
		}
	}
	return regMathString
}
