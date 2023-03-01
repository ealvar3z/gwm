package client

import (
	"fmt"

	"github.com/jezek/xgbutil/xwindow"
)

type WorkSpaceError struct {
	workspace uint8
}

// Implements the error interface
func (e WorkSpaceError) Error() string {
	msg := "invalid workspace: %d"
	return fmt.Sprintf(msg, e.workspace)
}

type Client struct {
	Window     xwindow.Window
	Workspace  WorkSpaceError
	Visible    bool
	Controlled bool
	FullScreen bool
	PaddingTop bool
}
