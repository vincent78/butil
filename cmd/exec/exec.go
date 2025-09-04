package exec

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/vincent78/butil/model"
	"golang.org/x/text/encoding/simplifiedchinese"
)

//func run(name string, arg ...string) *Cmd {
//	return exec.Command("/bin/sh", "-c", name, arg...).Start()
//}

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

var EnvPath = ""

func ShellCommand(cmdStr string, arg ...string) error {
	if len(arg) > 0 {
		_, err := fmt.Fprintf(os.Stdout, "%s %v \n", cmdStr, arg)
		if err != nil {
			return err
		}
	} else {
		_, err := fmt.Fprintf(os.Stdout, "%s \n", cmdStr)
		if err != nil {
			return err
		}
	}

	var cmdTmp = cmdStr
	for str := range arg {
		cmdTmp = fmt.Sprintf("%v %v", cmdTmp, str)
	}
	cmd := run(cmdTmp)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func ShellCommandWithCtx(ctx context.Context, ch chan string, rt chan model.RespModel, cmd string) {
	c := runWithCtx(ctx, cmd)
	stdout, err := c.StdoutPipe()
	if err != nil {
		rt <- model.FailureRespWithErrModel(model.ErrorBaseNormal(err.Error()))
	} else {
		go func() {
			reader := bufio.NewReader(stdout)
			for {
				// 其实这段去掉程序也会正常运行，只是我们就不知道到底什么时候Command被停止了，而且如果我们需要实时给web端展示输出的话，这里可以作为依据 取消展示
				select {
				// 检测到ctx.Done()之后停止读取
				case <-ctx.Done():
					if ctx.Err() != nil {
						rt <- model.FailureRespWithErrModel(model.ErrorBaseNormal(ctx.Err().Error()))
					} else {
						rt <- model.FailureRespWithErrModel(model.ErrorBaseNormal("程序被终止"))
					}
					return
				default:
					readString, err := reader.ReadString('\n')
					if err != nil {
						if err != io.EOF {
							rt <- model.FailureRespWithErrModel(model.ErrorBaseNormal(err.Error()))
						}
						return
					}
					ch <- strings.TrimSuffix(readString, "\n")
				}
			}
		}()
		err = c.Start()
		if err != nil {
			rt <- model.FailureRespWithErrModel(model.ErrorBaseNormal(err.Error()))
		} else {
			err = c.Wait()
			rt <- model.FailureRespWithErrModel(model.ErrorBaseNormal(err.Error()))
		}
	}

}

func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}
