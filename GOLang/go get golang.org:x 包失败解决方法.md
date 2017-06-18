# go get golang.org/x 包失败解决方法
---
由于限制问题，国内使用 go get 安装 golang 官方包可能会失败，如我自己在安装 collidermain 时，出现了以下报错：

	$ go get collidermain
	package golang.org/x/net/websocket: unrecognized import path "golang.org/x/net/websocket" (https fetch: Get https://golang.org/x/net/websocket?go-get=1: dial tcp 216.239.37.1:443: i/o timeout)

不翻墙的情况下怎么解决这个问题？其实 golang 在 github 上建立了一个[镜像库](https://github.com/golang)，如 https://github.com/golang/net 即是 https://golang.org/x/net 的镜像库

获取 golang.org/x/net 包，其实只需要以下步骤：

	mkdir -p $GOPATH/src/golang.org/x
	cd $GOPATH/src/golang.org/x
	git clone https://github.com/golang/net.git

其它 golang.org/x 下的包获取皆可使用该方法
