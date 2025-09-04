[toc]

# 测试

使用netcat发送测试数据

> $ echo -n -e "Hello Go-Netty\nhttps://go-netty.com\n" | nc 127.0.0.1 3333

# protobuf

## 安装相关包

```
go get -u github.com/golang/protobuf/proto
go get -u google.golang.org/protobuf
go get -u google.golang.org/grpc
```

## 转换命令

> protoc --go_out=plugins=grpc:./gcnet/proto.pb gcnet/proto/*.proto


