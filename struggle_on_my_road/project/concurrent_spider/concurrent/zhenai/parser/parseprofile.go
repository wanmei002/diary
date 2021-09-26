package parser

import (
	"regexp"
	"spider/engine"
	"spider/model"
	"strconv"
)

// 获取正则表达式 并且在全局变量中定义

var Gender = regexp.MustCompile(`<td><span class="label">性别：</span><span field="">([^<]+)</span></td>`)
var ageRe = regexp.MustCompile(`<td><span class="label">年龄：</span>([\d]+)岁</td>`)
var Height = regexp.MustCompile(`<td><span class="label">身高：</span>([\d]+)CM</td>`)
var Weight = regexp.MustCompile(`<td><span class="label">体重：</span><span field="">([\d]+)KG</span></td>`)
var Income = regexp.MustCompile(`<td><span class="label">月收入：</span>([^<]+)</td>`)
var Marriage = regexp.MustCompile(`<td><span class="label">婚况：</span><span field=""> ([^<]+)</span></td>`)
var Education = regexp.MustCompile(`<td><span class="label">学历：</span>([^<]+)</td>`)
var Occupation = regexp.MustCompile(`<td><span class="label">职业： </span>([^<]+)</td>`)
var Hokou = regexp.MustCompile(`<td><span class="label">籍贯：</span>([^<]+)</td>`)
var Xingzuo = regexp.MustCompile(`<td><span class="label">星座：</span><span field="">([^<]+)</span></td>`)
var House = regexp.MustCompile(`<td><span class="label">住房条件：</span><span field="">([^<]+)</span></td>`)
var Car = regexp.MustCompile(`<td><span class="label">是否购车：</span><span field="">([^<]+)</span></td>`)

func ParseProfile(contents []byte, name string) engine.ParseResult {
	profile := model.Profile{}
	age, _ := strconv.Atoi(extractString(contents,ageRe))
	profile.Age = age
	height, _ := strconv.Atoi(extractString(contents,Height))
	profile.Height = height
	weight, _ := strconv.Atoi(extractString(contents,Weight))
	profile.Weight = weight
	profile.Name = name
	profile.Gender = extractString(contents,Gender)
	profile.Income = extractString(contents,Income)
	profile.Marriage = extractString(contents,Marriage)
	profile.Education = extractString(contents,Education)
	profile.Occupation = extractString(contents,Occupation)
	profile.Hokou = extractString(contents,Hokou)
	profile.Xingzuo = extractString(contents,Xingzuo)
	profile.House = extractString(contents,House)
	profile.Car = extractString(contents,Car)
	//只需要传入内容
	result := engine.ParseResult{
		Items: []interface{}{profile},
	}
	return result
}

func extractString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}