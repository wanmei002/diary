package controller

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

type UEditor struct {
	Controller
}

func (ue *UEditor) ControllerUE(){
	ue.IsTpl = false
	err := ue.Ct.Request.ParseForm()
	if err != nil {
		fmt.Fprintf(ue.Ct.ResponseWriter, "解析参数失败: %v", err)
		return
	}

	op := ue.Ct.Request.Form.Get("action")

	switch op {
	case "config" :
		file, err := os.Open("F:/goProgram/src/myblog/conf/config.json")
		if err != nil {
			fmt.Fprintf(ue.Ct.ResponseWriter, "打开文件错误 : %v", err)
			return
		}
		defer file.Close()
		fd, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Fprintf(ue.Ct.ResponseWriter, "读取文件失败 : %v", err)
			return
		}

		src := string(fd)
		re, _ := regexp.Compile(`\/\*[\S\s]+?\*\/`)  // 匹配里面的注释
		src = re.ReplaceAllString(src, "")
		tt := []byte(src)
		var r interface{}
		err = json.Unmarshal(tt, &r) //这个byte要解码
		if err != nil {
			fmt.Fprintf(ue.Ct.ResponseWriter, "json decode failed %v", err)
			return
		}

		tt , err = json.Marshal(r)
		if err != nil {
			fmt.Fprintf(ue.Ct.ResponseWriter, "json encode failed $v", err)
			return
		}
		fmt.Fprint(ue.Ct.ResponseWriter, string(tt))

	// 上传图片的功能
	case "uploadimage":
		err := ue.Ct.Request.ParseForm()
		if err != nil {
			fmt.Fprintf(ue.Ct.ResponseWriter, "uploadimage parseform fail : %v", err)
			return
		}
		fmt.Println("打印所有的请求 : ",ue.Ct.Request.PostForm)
		fmt.Println("打印 upfile :", ue.Ct.Request.PostForm.Get("upfile"))
		// 开始上传
		// 文件路径
		file, h, err := ue.Ct.Request.FormFile("upfile")
		defer file.Close()
		if err != nil {
			ue.Error(err)
		}
		fmt.Println("file header: ", h, "\n detail-filename: ", h.Filename, " ;\n fileHeader:",h.Header)
		// 文件路径
		filePath := "F:/goProgram/src/myblog/static/upload/"+time.Now().Format("20060102")
		err = os.MkdirAll(filePath, 0777)
		ue.Error(err)
		// 文件名
		fileName := filePath + "/" + makeFileName() + ".jpg"
		f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		ue.Error(err)
		defer f.Close()
		io.Copy(f, file)
		ret_json := map[string]interface{}{
			"state":"SUCCESS",
			"url":strings.Replace(fileName, "F:/goProgram/src/myblog", "",1),
			"title":h.Filename,
			"original":h.Filename,
			"type":h.Header.Get("Content-Type"),
		}
		json, _ := json.Marshal(ret_json)
		fmt.Fprintf(ue.Ct.ResponseWriter, string(json))


	default:
		fmt.Fprint(ue.Ct.ResponseWriter, `{"msg":"请求地址错误"}`)
	}
}

// 文件名的生成
func makeFileName() string {
	times := time.Now().Unix()
	rand.Seed(times)
	num := rand.Intn(9999)
	str := fmt.Sprintf("%v%v", times, num)
	//m := md5.New()
	//md5Str :=
	has := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", has)

}
