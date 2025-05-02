package codeUtil

import (
	"github.com/vincent78/butil/token"
	"strings"
)

func GeneMesId() string {
	return strings.ToLower(token.GenKsuid())
}
