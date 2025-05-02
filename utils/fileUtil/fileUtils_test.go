package fileUtil

import (
	"fmt"
	"gitee.com/vincent78/gcutil/utils/strUtil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestReadFile(t *testing.T) {
	println(CurrPath())
}

func TestWriteFile(t *testing.T) {
	path := "/tmp/butil"
	if result := Exist(path); !result {
		um, _ := strconv.ParseInt(strconv.Itoa(700), 8, 0)
		CreatePath(path, os.FileMode(um))
	}

	file := filepath.Join(path, "write.txt")
	//println(timeUtil.NormalTimeStampMillis())
	for i := 0; i < 20; i++ {
		WriteFile(file, strUtil.GetRandomString(1024))
		WriteFile(file, "\n")
	}
}

func TestGetFiles(t *testing.T) {
	path := "/users/vincent/workspace/12_GO/src/github.com/mattn"
	files, _ := GetFiles(path, nil, NewGetFilesOption())
	for _, str := range files {
		println(str)
	}
}

func TestReplaceStr(t *testing.T) {
	src := "conf/container/config/params.json"
	dst := "/tmp/data/sim/CARRIER_192168001002/config/params.json"
	e := ReplaceStr(src, dst, map[string]string{
		"${agvCode}": "testAgvCode",
	})

	if e != nil {
		fmt.Errorf("replaceError: %v", e)
	} else {
		fmt.Println("replace success.")
	}
}
