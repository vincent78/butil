package logger3

import (
	"context"
	"fmt"
	"github.com/vincent78/butil/cmd/console"
	"github.com/vincent78/butil/config"
	"github.com/vincent78/butil/thread"
	"github.com/vincent78/butil/utils/fileUtil"
	"github.com/vincent78/butil/utils/timeUtil"
	"os"
	"runtime/debug"
)

type FileLogger struct {
	LogModel
	Name     string   `json:"name"`     // 文件名
	Path     string   `json:"path"`     // 日志路径
	MaxSize  int      `json:"maxSize"`  // 单文件的最大尺寸
	Compress bool     `json:"compress"` // 是否压缩文件
	currFile string   `json:"-"`
	file     *os.File `json:"-"`
}

var lineBytes = []byte{byte('\n')}

func NewFileLogger(ctx context.Context, config config.LoggerConfig) (*FileLogger, error) {
	log := &FileLogger{
		LogModel: LogModel{
			Level:     LoggerLeverl(config.Level + 1),
			MsgFormat: MsgFormatDefault,
			Queue:     make(chan *MsgModel, 1024),
			ctx:       ctx,
		},
		Name:     config.Prefix,
		Path:     config.Home,
		MaxSize:  512,
		Compress: true,
	}
	err := log.InitWithConfig(config)
	if err != nil {
		return nil, err
	}

	thread.Go(func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf(console.Yellow("context cancel"))
				_ = log.file.Sync()
				_ = log.file.Close()
				return
			case m := <-log.Queue:
				err = log.checkFile()
				if err != nil {
					fmt.Printf("fileLogger checkFile error: %v \n", err)
				} else {
					_, err = log.file.Write(lineBytes)
					_, err = log.file.Write(m.Bytes())
					if m.Level == LevelError {
						_, err = log.file.Write(debug.Stack())
					}
					if err != nil {
						fmt.Printf("fileLogger error: %v \n", err)
					}
				}
				m.SendBack()
			}
		}
	})
	return log, nil
}

func (log *FileLogger) checkFile() error {
	fileName := log.getName()
	filePath := fileUtil.GetAbs(fileUtil.Join(log.Path, fileName))
	err := fileUtil.IsNotExistMkDir(filePath)
	if err != nil {
		return err
	}
	if log.currFile != fileName {
		if log.file != nil {
			err := log.file.Close()
			fmt.Printf("close the old file[%v] point error: %v", log.file.Name(), err)
			go func(f string) {
				log.backup(f)
			}(log.currFile)
		}
		log.file, _ = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, fileUtil.GetFileMode(600))
		log.currFile = fileName
	}
	return nil
}

func (log *FileLogger) getName() string {
	fileName := timeUtil.NowFmtStr(timeUtil.DateFormat)
	if log.Name != "" {
		fileName = fmt.Sprintf("%v.%v", log.Name, fileName)
	}
	fileName = fmt.Sprintf("%v.log", fileName)
	return fileName
}

func (log *FileLogger) backup(filePath string) {
	fmt.Println(filePath)
}

func (log *FileLogger) InitWithConfig(config config.LoggerConfig) error {
	log.Level = LoggerLeverl(config.Level + 1)
	err := fileUtil.IsNotExistMkDir(config.Home)
	if err != nil {
		return err
	}
	return log.checkFile()
}
