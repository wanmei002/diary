### go build 命令详解
#### 直接在目录运行 go build 
 - go build 会搜索当前目录里的 *.go 文件，编译这些文件, 生成可执行文件到当前目录下
 
#### go build filename_list
 - `go build tool.go main.go` 跟上面的差不多, 适合少量的文件编译
 
#### go build -o exec_file  filename_list
 - -o 后面紧跟 要生成的可执行文件名  
 
#### 其它可选参数
 - `-v` 编译时显示包名
 - `-`