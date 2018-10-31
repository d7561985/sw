package grifts

import (
	"github.com/d7561985/sw/mock/actions"

	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
