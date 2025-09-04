# 概述

最基本的类库，可供其它项目引用，中间包括了所有独立的方法与函数

# 第三方框架

- [控制台框架-cli](https://github.com/urfave/cli)
- [配置文件-viper](https://github.com/spf13/viper)
- [缓存数据1-go-cache](https://github.com/patrickmn/go-cache)
- [缓存数据2-cache2go](github.com/muesli/cache2go)
- [日志框架-zap](https://github.com/uber-go/zap)
- [HTTP框架-gin](https://github.com/gin-gonic/gin)
- [Excel操作-excelize](https://github.com/qax-os/excelize)

# 编译部署相关

## linux下运行

> CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go

## window下运行

> CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
