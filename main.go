package main

import (
	"gitee.com/vincent78/gcutil/sys"
)

func main() {
	println("test")
	defer sys.RunMemProfile()()
	sys.BlockBySignal()
}
