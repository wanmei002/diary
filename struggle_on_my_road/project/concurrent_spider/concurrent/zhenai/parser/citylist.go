package parser

import (
	"regexp"
	"spider/engine"
)

// ParseCityList afa
// @param contents
// return [][]string
func PrintCityList(contents []byte) engine.ParseResult {
	reg, err := regexp.Compile(`<a href="(http://www.zhenai.com/zhenghun/[a-z0-9]*)"[^>]*>([^<]+)</a>`)
	if err != nil {
		panic(err)
	}
	// regMatch := reg.FindAllString(string(contents), -1)
	regMath := reg.FindAllSubmatch(contents, -1)

	result := engine.ParseResult{}

	for _, m := range regMath {
		result.Items = append(result.Items, string(m[2]))
		result.Requests = append(result.Requests, engine.Request{
			string(m[1]),
			engine.Nilparser,
		})
	}
	return result
}
