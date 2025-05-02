package fileUtil

import (
	log "github.com/vincent78/butil/logger2/logger"
	"github.com/vincent78/butil/utils/strUtil"
	"os"
	"path"
	"path/filepath"

	"runtime"
	"strings"
)

// 最终方案-全兼容
func GetCurrentAbPath() string {
	dir := GetCurrentAbPathByExecutable()
	if strings.Contains(dir, GetTmpDir()) {
		return GetCurrentAbPathByCaller()
	}
	return dir
}

// 获取当前执行程序所在的绝对路径
func GetCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func GetCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

// 获取系统临时目录，兼容go run
func GetTmpDir() string {
	dir := os.Getenv("TEMP")
	if dir == "" {
		dir = os.Getenv("TMP")
	}
	res, _ := filepath.EvalSymlinks(dir)
	return res
}

func GetCurrentProjectPath() string {
	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// 假定工程根目录下有特定的文件或目录作为标识，例如 "go.mod"
	rootMarker := "go.mod"
	// 向上查找直到找到工程根目录标识
	for !CheckExist(Join(cwd, rootMarker)) {
		newDir := GetPrePath(cwd)
		if newDir == cwd { // 防止无限循环
			break
		}
		cwd = newDir
	}
	return cwd
}

// 判断文件是否存在（filepath 文件的绝对路径）
func CheckFileExist(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// IsNotExistMkDir 检查文件夹是否存在
// 如果不存在则新建文件夹
func IsNotExistMkDir(src string) error {
	if exist := !CheckExist(src); exist == false {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

// MkDir 新建文件夹
func MkDir(src string) error {
	err := os.MkdirAll(filepath.Dir(src), os.ModePerm)
	//err := os.MkdirAll(src, 0777)
	if err != nil {
		return err
	}

	return nil
}

// ps must is the full path name
func GetPrePath(ps string) string {
	pa := strings.Split(ps, string(os.PathSeparator))
	pl := len(pa)
	if pl >= 3 {
		l := pa[pl-1]
		if !strUtil.Constains(l, ".") {
			return strings.Join(pa[:pl-1], string(os.PathSeparator))
		} else {
			return strings.Join(pa[:pl-2], string(os.PathSeparator))
		}
	} else {
		if strings.Trim(pa[1], "") == "" {
			return string(os.PathSeparator)
		} else {
			return ps
		}
	}
}
