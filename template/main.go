package main

import (
	"fmt"
	"github.com/vincent78/butil/utils/fileUtil"
	"github.com/vincent78/butil/utils/templateUtil"
)

func main() {
	logic1()
}

func logic1() {
	basePath := fileUtil.GetCurrentProjectPath()
	tf := fileUtil.Join(basePath, "template", "example/simple/hello.go.template")
	tg := fileUtil.Join("/tmp", "hello.go")
	data := map[string]interface{}{"msg": "\"hello world\""}
	if err := templateUtil.GenFile(tf, tg, data); err != nil {
		fmt.Printf("\n--- error: %v", err)
	} else {
		fmt.Println("--- success")
	}
}
