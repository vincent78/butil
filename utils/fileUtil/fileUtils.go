package fileUtil

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	pathUtil "path/filepath"
	"strconv"
	"strings"

	imgext "github.com/shamsher31/goimgext"
)

func Join(str ...string) string {
	return pathUtil.Join(str...)
}

// perm ：文件权限，一个八进制数。r（读）04，w（写）02，x（执行）01。
func GetFileMode(num int) os.FileMode {
	var um, _ = strconv.ParseInt(strconv.Itoa(num), 8, 0)
	return os.FileMode(um)
}

func Exist(path string) bool {
	p, e := pathUtil.Abs(path)
	if e != nil {
		return false
	}
	_, err := os.Stat(p) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func ReadFile(filePath string) (string, error) {
	exist := Exist(filePath)
	if !exist {
		return "", fmt.Errorf("[%v] is not exist", filePath)
	}
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func WriteFileByMode(nameWithPath string, data []byte, perm os.FileMode) error {
	return os.WriteFile(nameWithPath, data, perm)
}

func WriteFile(filePath string, content string) (re error) {
	defer func() {
		if err := recover(); err != nil {
			re = err.(error)
		}
	}()
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, GetFileMode(600))
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	writer := bufio.NewWriter(f)
	_, err = AppendStrWithWriter(writer, content)
	return err
}

func WriteFileByBytes(content bytes.Buffer, name string) {
	file, err := os.Create(name)
	if err != nil {
		log.Println(err)
	}
	_, err = file.WriteString(content.String())
	if err != nil {
		log.Println(err)
	}
	err = file.Close()
	if err != nil {
		return
	}
}

func AppendStrWithWriter(writer *bufio.Writer, content string) (w *bufio.Writer, re error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("\nappendStrWithWriter error: %v \n", err)
			re = err.(error)
		}
	}()
	if writer != nil {
		_, err := writer.WriteString(content)
		if err != nil {
			fmt.Printf("\nwrite to file error: %v \n", err)
			re = err.(error)
			return
		}
		err = writer.Flush()
		if err != nil {
			fmt.Printf("\nwrite to file flush error: %v \n", err)
			return nil, err
		}
		return writer, nil
	} else {
		return writer, errors.New("writer is nil.")
	}

}

type GetFilesOption struct {
	subDir   string //子目录
	blankDir bool   //是否包含空目录
	hideObj  bool   //是否包含隐藏内容
	abs      bool   //包括绝对路径
	goNext   bool   //条件符合后是否查子目录
}

func NewGetFilesOption() GetFilesOption {
	return GetFilesOption{
		subDir:   "",
		blankDir: false,
		hideObj:  true,
		abs:      false,
		goNext:   false,
	}
}

// WalkDir  获取符合条件的子目录
// TODO:  当碰到已符合条件的目录后，是否继续查这个目录下的子目录？ 这个还未实现
func WalkDir(basePath string, filer func(name string) bool, option GetFilesOption) ([]string, error) {
	exist := Exist(basePath)
	if !exist {
		return nil, fmt.Errorf("[%v] is not exist", basePath)
	}
	result := make([]string, 0)
	err := pathUtil.WalkDir(basePath, func(path string, info fs.DirEntry, err error) error {
		curr, _ := pathUtil.Rel(basePath, path)
		if "." == curr {
			return nil
		}

		if !info.IsDir() {
			return nil
		}
		// 过滤隐藏文件
		if option.hideObj && strings.HasPrefix(info.Name(), ".") {
			return nil
		} else if filer(curr) {
			if !option.abs && pathUtil.IsAbs(path) {
				path, _ = pathUtil.Rel(basePath, path)
			}
			result = append(result, path)
		}
		return nil
	})
	return result, err
}

func Walk(basePath string, filer func(name string, isDir bool) bool, option GetFilesOption) ([]string, error) {
	exist := Exist(basePath)
	if !exist {
		return nil, fmt.Errorf("[%v] is not exist", basePath)
	}
	result := make([]string, 0)
	err := pathUtil.Walk(basePath, func(path string, info fs.FileInfo, err error) error {
		// 过滤隐藏文件
		if option.hideObj && strings.HasPrefix(info.Name(), ".") {
			return nil
		} else if filer(info.Name(), info.IsDir()) {
			if !option.abs && pathUtil.IsAbs(path) {
				path, _ = pathUtil.Rel(basePath, path)
			}
			result = append(result, path)
		}
		return nil
	})
	return result, err
}

func GetFiles(path string, filter func(name string, isDir bool) bool, option GetFilesOption) ([]string, error) {
	exist := Exist(path)
	if !exist {
		return nil, fmt.Errorf("[%v] is not exist", path)
	}

	entries, err := os.ReadDir(pathUtil.Clean(pathUtil.ToSlash(path)))
	if err != nil {
		return nil, err
	}
	objs := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		objs = append(objs, info)
	}

	var result []string
	for _, obj := range objs {
		// 过滤隐藏文件
		if option.hideObj && strings.HasPrefix(obj.Name(), ".") {
			continue
		}
		isDir := obj.IsDir()
		name := Join(option.subDir, obj.Name())
		if filter != nil && !filter(name, isDir) {
			continue
		}
		if !isDir {
			result = append(result, name)
		} else {
			subOptions := option
			subOptions.subDir = Join(option.subDir, obj.Name())
			subResult, re := GetFiles(Join(path, obj.Name()), filter, subOptions)
			if re == nil {
				result = append(result, subResult...)
			}
		}
	}
	return result, nil
}

func GetSubDirs(path string, filter func(name string) bool) ([]string, error) {
	exist := Exist(path)
	if !exist {
		return nil, fmt.Errorf("[%v] is not exist", path)
	}

	entries, err := os.ReadDir(pathUtil.Clean(path))
	if err != nil {
		return nil, err
	}
	objs := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		objs = append(objs, info)
	}
	//objs, err := ioutil.ReadDir(pathUtil.Clean(path))
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(objs))
	for _, obj := range objs {
		if obj.IsDir() {
			name := obj.Name()
			if filter == nil || filter(name) {
				result = append(result, obj.Name())
			}
		}
	}
	return result, nil
}

// CopyFile over 是否覆盖
func CopyFile(source, target string, over bool) error {
	if ex := Exist(source); !ex {
		return errors.New("the source is not exist.")
	}

	if ex := Exist(target); ex && !over {
		return errors.New("the target is exist and require not override.")
	} else {
		os.Remove(target)
	}

	tp, _ := pathUtil.Split(target)
	if rt := CreatePath(tp, GetFileMode(700)); rt {
		src, err := os.Open(source)
		if err != nil {
			return err
		}
		defer src.Close()
		dst, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE, GetFileMode(644))
		if err != nil {
			return err
		}
		defer dst.Close()
		_, err = io.Copy(dst, src)
		return err
	} else {
		return errors.New("createPath error")
	}
}

func DeleteFile(target string) {
	if ex := Exist(target); ex {
		os.RemoveAll(target)
	}
}

func MoveFile(source, target string) {
	if err := CopyFile(source, target, true); err != nil {
		DeleteFile(source)
	}
}

// GetSize 获取文件大小
func GetSize(f multipart.File) (int, error) {
	content, err := io.ReadAll(f)

	return len(content), err
}

// GetExt 获取文件后缀
func GetExt(fileName string) string {

	fileType := pathUtil.Ext(fileName)
	if fileType != "" {
		return fileType[1:]
	} else {
		return ""
	}
}

func GetName(filePath string) string {
	fileNameWithExt := path.Base(filePath)
	fileExt := GetExt(fileNameWithExt)
	tmp := strings.TrimSuffix(fileNameWithExt, fileExt)
	// 删除末位的点
	if len(tmp) > 1 && strings.HasSuffix(tmp, ".") {
		tmp = tmp[0 : len(tmp)-1]
	}
	return tmp
}

// CheckExist 检查文件是否存在
func CheckExist(src string) bool {
	_, err := os.Stat(src)
	return !os.IsNotExist(err)
}

// CheckPermission 检查文件权限
func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

/*
*

	if err := fileUtil.MkDirByMode(file, 0700); err != nil {
			return gm.FailureResultWithError[string](500, fmt.Errorf("Could not create directory %s", file))
	}
*/
func MkDirByMode(src string, model os.FileMode) error {
	err := os.MkdirAll(pathUtil.Dir(src), model)
	if err != nil {
		return err
	}
	return nil
}

// Open 打开文件
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// GetImgType 获取Img文件类型
func GetImgType(p string) (string, error) {
	file, err := os.Open(p)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	buff := make([]byte, 512)

	_, err = file.Read(buff)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	filetype := http.DetectContentType(buff)

	ext := imgext.Get()

	for i := 0; i < len(ext); i++ {
		if strings.Contains(ext[i], filetype[6:len(filetype)]) {
			return filetype, nil
		}
	}

	return "", errors.New("Invalid image type")
}

// GetType 获取文件类型
func GetType(p string) (string, error) {
	file, err := os.Open(p)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	buff := make([]byte, 512)

	_, err = file.Read(buff)

	if err != nil {
		log.Println(err)
	}

	filetype := http.DetectContentType(buff)

	//ext := GetExt(p)
	//var list = strings.Split(filetype, "/")
	//filetype = list[0] + "/" + ext
	return filetype, nil
}

func ReplaceStr(src, dst string, rmap map[string]string) error {
	exist := Exist(src)
	if !exist {
		return fmt.Errorf("source file [%v] is not exist", src)
	}
	exist = Exist(dst)
	if exist {
		DeleteFile(dst)
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		fmt.Println("Open write file fail:", err)
		os.Exit(-1)
	}
	defer out.Close()

	br := bufio.NewReader(in)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		//newLine := strings.Replace(string(line), os.Args[2], os.Args[3], -1)
		newLine := string(line)
		for k, v := range rmap {
			newLine = strings.ReplaceAll(newLine, k, v)
		}
		_, err = out.WriteString(newLine + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// 获取文件大小
func GetFileSize(filename string) (int64, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

// FormatSize 转换文件大小的单位
func FormatSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

func GetFileSizeByFormat(filename string) (string, error) {
	size, err := GetFileSize(filename)
	if err != nil {
		return "", err
	}
	ret := FormatSize(size)
	return ret, nil
}

func FileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	_, _ = io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil)), nil
}
