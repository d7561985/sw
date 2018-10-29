package buffalo

import (
	"github.com/d7561985/sw/generators/general"
	"path/filepath"
)

type Buffalo struct {
	general.General
}

func (*Buffalo) RootPath(curpath string) string {
	return filepath.Join(curpath, "buffalo", "root.go")
}
