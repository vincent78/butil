package sys

import (
	"gitee.com/vincent78/gcutil/utils/fileUtil"
	"os"
	"strconv"
)

func GeneratePID(path, name string) (string, error) {
	if path == "" {
		path = fileUtil.CurrPath()
	}
	if name == "" {
		name = "pid"
	}
	//pid := os.Getegid()
	pid := os.Getpid()
	fileStr := fileUtil.Join(path, name)
	if fileUtil.Exist(fileStr) {
		fileUtil.DeleteFile(fileStr)
	}
	return fileStr, fileUtil.WriteFile(fileStr, strconv.Itoa(pid))

}
