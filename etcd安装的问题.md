## 问题一 
获取 etcd client golang代码的时候，出现了如下问题:
github.com/coreos/bbolt: github.com/coreos/bbolt@v1.3.6: parsing go.mod:
        module declares its path as: go.etcd.io/bbolt
                but was required as: github.com/coreos/bbolt
                
                
进入项目 github.com/coreos/bbolt 
发现 go.mod 文件内 module go.etcd.io/bbolt
> 可见 module 不是随便定义的，这样相当于重定向
 
解决问题的方法在本地 module里重定向 `go mod edit -replace github.com/coreos/bbolt@v1.3.6=go.etcd.io/bbolt@v1.3.6`，
这样go get github.com/coreos/bbolt@v1.3.6 的时候就会直接去go.etcd.io/bbolt这里获取包

## 问题二
google.golang.org/grpc/naming: module google.golang.org/grpc@latest found (v1.38.0), but does not contain package google.golang.org/grpc/naming

这个网上查了下 grpc 的版本太新了 得降低版本到 v1.26.0

解决方法: `go mod edit -replace google.golang.org/grpc@v1.38.0=google.golang.org/grpc@v1.26.0`

## 问题三
finding github.com/coreos/go-systemd latest
go get github.com/coreos/go-systemd: no matching versions for query "upgrade"

不知道什么原因
直接在 GOAPTH/src 目录下 新建 mkdir -p github.com/coreos, 进入这个目录里，git clone https://github.com/coreos/go-systemd.git

然后 在本地项目里添加 replace github.com/coreos/go-systemd => GOPATH/src/github.com/coreos/go-systemd

代码直接指向本地项目路径
