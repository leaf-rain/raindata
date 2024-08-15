package conf

import (
	"path/filepath"
	"strings"
)

func (a *AutoCode) WebRoot() string {
	webs := strings.Split(a.Web, "/")
	if len(webs) == 0 {
		webs = strings.Split(a.Web, "\\")
	}
	return filepath.Join(webs...)
}
