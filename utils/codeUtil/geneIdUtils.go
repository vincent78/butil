package codeUtil

import (
	"gitee.com/vincent78/gcutil/token"
	"strings"
)

func GeneMesId() string {
	return strings.ToLower(token.GenKsuid())
}
