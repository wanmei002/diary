#### Go modules
 1. 官网定义
    - 模块 是 go package 的集合 存储在文件树中的  go.mod 文件中
    - 从Go 1.11开始，当当前目录或任何父目录具有Go.mod时，Go命令将启用模块的使用，前提是该目录位于$GOPATH/src之外。（在$GOPATH/src中，为了兼容，go命令仍然在旧的GOPATH模式下运行，即使找到go.mod也是如此。有关详细信息，请参阅go命令文档。）从go 1.13开始，模块模式将是所有开发的默认模式。
 2. 人话 : 以后代码不用在 $GOPATH/src 目录下了, 随便一个目录里，
        只要当前目录或父目录里有 go.mod 文件，则就可以用 go build ... 来生成可执行文件
        
#### 步骤
 1. 在一个文件夹 中 运行 `go mod init modules_name`, 此时文件夹里会多出一个 go.mod 文件 里面保存着 modules 名称 和 golang 版本
 2. 如果要引入其它模块 可以 运行 `go get github.cm/xxx/xxx` 此时 文件夹下会生成一个go.sum 文件，里面保存着 包的各种依赖版本等等
 3. `go list -m all`  添加一个直接依赖项通常也会带来其他间接依赖项。命令go list-m all列出当前模块及其所有依赖项
 4.  `go mod tidy` 命令清除没有使用的依赖项
 
 
#### 迁移
 + 移动整个模块 在go run 或 go build 时 会自动检索 go.mod go.sum 中的依赖包并且下载更新(纯属个人理解 勿喷)