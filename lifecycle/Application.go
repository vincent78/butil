package lifecycle

import (
	"fmt"

	"github.com/vincent78/butil/sys"
)

func AppBaseInit() {
	p, err := sys.GeneratePID("", "")
	if err != nil {
		fmt.Printf("generate PID error: %v\n", err)
	} else {
		fmt.Printf("pid file - %v \n", p)
	}
}

func AppBasePrepared() bool {
	return true
}

func AppBaseDestory() {

}

func AppBasePause() {

}

func AppBaseActive() {

}
