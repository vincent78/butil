package strUtil

import (
	"fmt"
	"strings"
)

func FmtSlice(f, space string, args ...any) string {
	if len(args) > 0 {
		sb := strings.Builder{}
		for _, v := range args {
			if sb.Len() > 0 {
				sb.WriteString(space)
			}
			sb.WriteString(fmt.Sprintf(f, v))
		}
		return sb.String()
	} else {
		return ""
	}
}

func FmtSliceByDefault(args ...any) string {
	f := "%v"
	space := " "
	if len(args) > 0 {
		sb := strings.Builder{}
		for _, v := range args {
			if sb.Len() > 0 {
				sb.WriteString(space)
			}
			sb.WriteString(fmt.Sprintf(f, v))
		}
		return sb.String()
	} else {
		return ""
	}
}
