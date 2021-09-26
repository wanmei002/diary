### go 实现 百度富文本框(ueditor)
> 没有什么难点 主要参考 ueditor PHP 版

> 抓取 ueditor php版 请求的后台链接 用 golang 替换

1. ueditor.config.js 请求的后台地址替换掉
    + 把 `serverUrl: URL + "php/controller.php"` 替换成 `serverUrl: /controller` ,ueditor 一般都是跟这个地址交互
    + url 传参 action 的值来跟后台交互 我主要实现了 返回配置文件  上传图片的功能
    + 我把 /controller 地址绑定上了 UEditor.ControllerUE 这个方法处理，主要实现了 返回配置信息 和 上传图片功能~~~原谅懒惰的我
    + 具体实现请看本目录下的 `UEditor.go` 文件
    
### 文件介绍
 + `Create.html` 富文本框的 html 页面
 + `config.json` 后端返回的 json 文件
 + `ueditor.config.js` ueditor 的js 配置文件
 + `UEditor.go` 后端文件, 主要实现了 返回 配置信息 和 上传图片的功能