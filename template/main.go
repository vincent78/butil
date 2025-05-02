package main

import (
	"fmt"
	"github.com/vincent78/butil/utils/fileUtil"
	"github.com/vincent78/butil/utils/templateUtil"
)

const basePath = "/Users/vincent/workspace/90_PROJECTS/zhejie/chain/butil/template"

func main() {
	logic1()
}

func logic1() {
	tf := fileUtil.Join(basePath, "example/simple/hello.go.template")
	tg := fileUtil.Join("/tmp/template", "hello.go")
	data := map[string]interface{}{"data": "\"hello world\""}
	if err := templateUtil.GenFile(tf, tg, data); err != nil {
		fmt.Printf("\n--- error: %v", err)
	} else {
		fmt.Println("--- success")
	}
}
