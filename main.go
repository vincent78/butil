package main

import (
	"github.com/vincent78/butil/sys"
)

func main() {
	println("test")
	defer sys.RunMemProfile()()
	sys.BlockBySignal()
}
