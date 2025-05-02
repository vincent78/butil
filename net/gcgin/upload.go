package gcgin

import (
	"errors"
	"fmt"
	"gitee.com/vincent78/gcutil/utils/objUtil"
	"gitee.com/vincent78/gcutil/utils/timeUtil"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var extArray = []string{"jpg", "png"}

// Upload 上传
/*
curl --location --request POST 'http://127.0.0.1:8080/mdc/upload' \
--form 'file=@"/users/vincent/temp/install.sh"'
*/
func Upload(c *gin.Context, targetPath string, exts []string) (string, error) {

	file, err := c.FormFile("file")
	if err != nil {
		return "", err
	}
	ext := strings.Split(file.Filename, ".")
	if exts == nil {
		exts = extArray
	} else {
		exts = append(exts, extArray...)
	}
	if !objUtil.HasObj("*", exts) && !objUtil.HasObj(ext[1], exts) {
		//c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "文件格式不支持！"})
		return "", errors.New("文件格式不支持")
	}
	//filePath := fmt.Sprintf("./img/%s/", time.Now().Format("2006/01/02"))
	_, err = os.Stat(targetPath)
	if err != nil {
		if !os.IsExist(err) { // 目录不存在，创建目录
			err := os.MkdirAll(targetPath, os.ModePerm)
			if err != nil {
				fmt.Println(err)
				return "", err
			}
		}
	}
	fileName := fmt.Sprintf("%v.%v", timeUtil.NowFmtStr(timeUtil.CompactFormat), ext[1])
	targetFile := ""
	//if strings.HasSuffix(targetPath, "/") {
	//	targetFile = targetPath + fileName
	//} else {
	//	targetFile = fmt.Sprintf("%v/%v", targetPath, fileName)
	//}
	targetFile = filepath.Join(targetPath, fileName)

	err = c.SaveUploadedFile(file, targetFile)
	if err != nil {
		fmt.Printf("upload error:%v \n", err)
		return "", err
	}
	//c.JSON(http.StatusOK, gin.H{"msg": "ok", "path": targetPath})
	return fileName, nil
}

/*
*

	curl -X POST http://localhost:8080/upload \
	  -F "upload[]=@/users/huahua/gopath/src/gintest/html/bar.tmpl" \
	  -H "Content-Type: multipart/form-data"
*/
func uploads(c *gin.Context) {
	// Multipart form
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]

	for _, file := range files {
		log.Println(file.Filename)

		// Upload the file to specific dst.
		// c.SaveUploadedFile(file, dst)
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}
