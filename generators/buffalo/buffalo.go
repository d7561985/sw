package buffalo

import (
	"dima/sw/generators/general"
	"path/filepath"
)

type Buffalo struct {
	general.General
}

func (*Buffalo) RootPath(curpath string) string {
	return filepath.Join(curpath, "buffalo", "root.go")
}
