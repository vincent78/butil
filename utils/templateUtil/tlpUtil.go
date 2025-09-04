package templateUtil

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/vincent78/butil/utils/fileUtil"
)

func GenFile(tf, target string, d interface{}) error {
	if exist := fileUtil.Exist(tf); !exist {
		return fmt.Errorf("can't find template file [%v] ", tf)
	}
	if err := fileUtil.IsNotExistMkDir(target); err != nil {
		return err
	}
	if t, err := template.ParseFiles(tf); err != nil {
		return err
	} else {
		var b1 bytes.Buffer
		err = t.Execute(&b1, d)
		if err != nil {
			return err
		}
		err = fileUtil.WriteFile(target, b1.String())
		if err != nil {
			return err
		}
	}
	return nil
}
