package netUtil

import (
	"fmt"
	"net/url"
)

func GetUrlPrefix(urlstr string) string {
	u, e := url.Parse(urlstr)
	if e != nil {
		return ""
	}
	r := fmt.Sprintf("%s://", u.Scheme)
	r = fmt.Sprintf("%s/%s", r, u.Host)
	return r
}
