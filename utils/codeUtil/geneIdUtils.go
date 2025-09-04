package codeUtil

import (
	"strings"

	"github.com/vincent78/butil/token"
)

func GeneMesId() string {
	return strings.ToLower(token.GenKsuid())
}
