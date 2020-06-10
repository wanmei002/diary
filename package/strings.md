### strings

 - `strings.Count(s1 string, s2 string) int` s2 在 s1 中重复的次数
 
 - `strings.Contains(s1 string, s2 string) bool` s2(字串) 是否被s1包含, 如果 s2 为空, 则恒为 `true`
 
 - `strings.ContainsAny(s1 string, s2 string) bool` s1 字符串中是否包含 s2 中的任意一个字符， 如果 s2 为空则为 `false`
 
 - `strings.Index(s1 stirng, s2 string) int` s2第一次在s1中出现的位置[从0开始], 如果s2="" 返回0， 如果没有 返回 -1
 
 - `strings.LastIndex(s1 string, s2 string) int` 返回 s2 最后一次出现在 s1 中的位置, 如果找不到返回-1, 如果 s2 为空, 则返回 s1 的长度
 
 - `strings.IndexRune(s string, r rune) int` 返回 rune类型的 r 在 s 中第一次出现的位置, 如果找不到 返回 -1
 
 - `strings.IndexAny(s, chars string) int` 返回字符串 `chars` 中的任何一个字符在字符串 s 中第一次出现的位置，
 如果找不到 返回 -1, 如果 chars 为空 返回 -1
 
 - `strings.LastIndexAny(s, chars string) int` 返回字符串 chars 中的任何一个字符在字符串 s 中最后一次出现的位置, 如果找不到, 则返回-1，如果 chars 为空 返回 -1
 
 
 - `strings.SplitN(s, sep string, n int) []string ` splitN 以 sep 为分割符, 将 s 切分成多个子串, 结果中不包含 sep 本身，
 如果 sep 为空, 则将 s 切分成 Unicode 字符列表, 如果 s 中没有 sep 子串, 则将整个 s 作为 []string 的第一个元素返回
 参数 n 表示最多切分出几个子串, 超出的部分将不再切分. 如果 n 为0 ，则返回 nil, 如果 n 小于 0，则不限制切分个数, 全部切分
 
 - `strings.SplitAfterN(s, sep string, n int) []string` 以 sep 为分隔符, 将 s 切分成多个子串, 结果中包含 sep 本身, 其它跟 SplitN 一样
 
 - `strings.Split(s, sep string) []string` 跟 SplitN 类似，只是不在指定切分的次数
 
 - `strings.SplitAfter(s, sep string) []string` 跟 SplitAfterN 类似, 只是不再指定 N 切分的次数
 
 - `strings.Fields(s string) []string` 以连续的空白字符为分割符, 将 s 切分成多个子串, 结果中不包含空白字符本身,
    空白字符有 `\t`, `\n`, `\v`, `\f`, `" "`,  `U+0085 (NEL)`, `U+00A0 (NBSP)`, 如果 s 中只包含空白字符, 则返回一个空列表
    
 - `strings.FieldsFunc(s string, f func(rune) bool) []string` FieldsFunc 以一个或多个满足 f(rune) 的字符为分割符, 将 s 切分成多个子串, 结果中不包含分隔符本身.
  如果 s 中没有满足 f(rune) 的字符,则返回一个空列表
  ```go
    func isSlash(r rune) bool {
    	return r == '\\' || r =='/'
    }
    func main(){
    	s := "C:\\Windows\\System32\\FileName"
    	ss := strings.FieldsFunc(s, isSlash)
    	fmt.Printf("%q\n", ss) // ["C:" "Windows" "System32" "FileName"]
    }
  ```
  
 - `strings.Join(a []string, sep string) string` 将 a 中的子串连接成一个单独的字符串, 子串之间用 sep 分隔
 
 - `strings.HasPrefix(s, prefix string) bool` 判断字符串 s 是否以 prefix 开头
 
 - `strings.HasSuffix(s, suffix string) bool` 判断字符串 s 是否以 suffix 结尾
 
 - `strings.Map(mapping func(rune) rune, s string) string` Map 将 s 中满足 mapping(rune) 的字符替换为 mapping(rune) 的返回值， 如果 mapping(rune) 返回负数, 则相应的字符将被删除
 ```go
  func Slash(r rune) rune {
  	if r == '\\' {
  		return '/'
  	}
  	return r
  }

  func main(){
  	s := "C:\\Windows\\System32\\FileName"
    ms := strings.Map(Slash, s)
    fmt.Printf("%q\n", ms) // "C:/Windows/System32/FileName"
  }
 ```
 
 - `strings.Repeat(s string, count int) string` Repeat 将 count 个字符串 s 连接成一个新的字符串
 
 - `strings.ToUpper(string) string  strings.ToLower(string)string` 将 s 中的所有字符修改为其大写格式，对于非 ASCII字符, 它的大写(小写) 格式需要查表转换
 
 - `strings.Title(s string) string` 将 单词的首字母变成 大写
 
 
 
 
 
 
 
 