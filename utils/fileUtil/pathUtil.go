package fileUtil

import (
	"os"
	"path"

	pathUtil "path/filepath"

	log "github.com/vincent78/butil/logger/logger2/logger"
	"github.com/vincent78/butil/utils/strUtil"

	"runtime"
	"strings"
)

// 最终方案-全兼容
func GetCurrentAbPath() string {
	dir := GetCurrentAbPathByExecutable()
	if strings.Contains(dir, GetTmpDir()) {
		//return GetCurrentAbPathByCaller()
		return GetCurrentProjectPath()
	}
	return dir
}

// 获取当前执行程序所在的绝对路径
func GetCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := pathUtil.EvalSymlinks(pathUtil.Dir(exePath))
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
	res, _ := pathUtil.EvalSymlinks(dir)
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
	err := os.MkdirAll(pathUtil.Dir(src), os.ModePerm)
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

func CreatePathWithDefaultMode(path string) string {
	if !strings.HasPrefix(path, string(os.PathSeparator)) {
		path = pathUtil.Join(CurrPath(), path)
	}
	var r = false
	if r = Exist(path); !r {
		r = CreatePath(path, GetFileMode(700))
	}
	if r {
		return path
	} else {
		return ""
	}
}

func CreatePath(path string, perm os.FileMode) bool {
	exist := Exist(path)
	if !exist {
		if te := os.MkdirAll(pathUtil.Clean(path), perm); te != nil {
			return false
		} else {
			return true
		}
	} else {
		return exist
	}
}

func CurrPath() string {
	dir, err := pathUtil.Abs(pathUtil.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func CurrPath2() string {
	_, file, _, _ := runtime.Caller(1)
	return path.Dir(file)
}

func GetAbs(path string) string {
	var err error
	if !pathUtil.IsAbs(path) {
		path = Join(CurrPath(), path)
		path, err = pathUtil.Abs(path)
		if err != nil {
			return ""
		}
	}
	return path
}

func LastPathName(fp string) string {
	return pathUtil.Base(fp)
}

// 获取可执行文件的名称
func executableName() string {
	path, err := os.Executable()
	if err != nil {
		return pathUtil.Base(os.Args[0])
	}
	return pathUtil.Base(path)
}
