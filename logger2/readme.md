# source: go-admin


# 配置文件
```yaml
settings:
  logger:
    # 日志存放路径
    path: /tmp/model/logs
    # 日志输出，file：文件，default：命令行，其他：命令行
    stdout: 'file' #控制台日志，启用后，不输出到文件
    # 日志等级, trace, debug, info, warn, error, fatal
    level: trace
    # 数据库日志开关
    enableddb: false
```

# 使用
## 入口初始化
```go 
	log.SetupLogger(
		logger.WithType(e.Type),
		logger.WithPath(e.Path),
		logger.WithLevel(e.Level),
		logger.WithStdout(e.Stdout),
		logger.WithCap(e.Cap),
	)
```