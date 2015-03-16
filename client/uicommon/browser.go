package uicommon

import (
	"strings"

	"github.com/gopherjs/gopherjs/js"
)

//SetCurrentPage causes the browser to switch to the page given by s.  This
//can be an absolute path, like /foo/bar.html, or a relative one, like
//bar.html.
func SetCurrentPage(s string) {
	if strings.HasPrefix(s, "/") {
		js.Global.Get("window").Get("location").Set("pathname", s)
		return
	}
	curr := js.Global.Get("window").Get("location").Get("pathname").String()
	parts := strings.Split(curr, "/")
	if len(parts) < 2 {
		parts = append(parts, s)
	} else {
		parts = append(parts[:len(parts)-1], s)
	}
	print("parts", parts)
	result := strings.Join(parts, "/")
	js.Global.Get("window").Get("location").Set("pathname", result)
}
