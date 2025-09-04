[toc]

# 安装protoc

[官网版本](https://github.com/protocolbuffers/protobuf/releases)  
[当前用的版本 3.17.3](https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protobuf-all-3.17.3.zip)

环境设置
> export PROTOBUF_HOME=$DEVELOPER/protobuf/protobuf-3.17.3

# 安装protoc-gen-go

```
$ go get -u github.com/golang/protobuf/protoc-gen-go
$ cd $GOPATH/pkg/mod/github.com/golang/protobuf@v1.5.2/protoc-gen-go
$ go install //编译二进制文件，这时会生成二进制文件，保存到$GOPATH/bin/
$ cd $GOPATH/bin/
$ sudo cp protoc-gen-go /usr/sbin/
```

# 安装相关包

```
go get -u github.com/golang/protobuf/proto
// go get -u github.com/golang/protobuf/protoc-gen-go
go get -u google.golang.org/protobuf
go get -u google.golang.org/grpc
```

# 执行命令

在项目根目录下执行
> protoc --go_out=plugins=grpc:./gcrpc/proto.pb gcrpc/proto/helloworld.proto