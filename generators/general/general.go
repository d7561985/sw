package general

import (
	"path/filepath"
)

type General struct {
}

func (g *General) RootPath(curpath string) string {
	return filepath.Join(curpath, "cmd", "root.go")
}
