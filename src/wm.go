package main

import (
	"context"
	"log"
	"os/exec"

	"github.com/ealvar3z/gwm/src/client"
	"github.com/jezek/xgb/xproto"
	"github.com/jezek/xgbutil"
	"github.com/jezek/xgbutil/xwindow"
)

type WindowManager struct {
	clients *client.Client
	config  *Config
	conn    *xgbutil.XUtil
	cursor  *xwindow.Window
}

func New() *WindowManager {
	// connect
	conn, err := xgbutil.NewConn()
	if err != nil {
		log.Fatalf("INFO: Unable to connect to X server: %v", err)
	}

	// load config.toml
	conf, err := GetConfig()
	if err != nil {
		log.Fatalf("INFO: Unable to load config.toml: %v", err)
	}

	// spin the client
	clients := client.New(conn, conf)

	// bring the cursor
	cursor, err := xwindow.Create(conn, conn.RootWin())
	if err != nil {
		log.Fatalf("INFO: Unable to load create cursor: %v", err)
	}

	return &WindowManager{
		clients: clients,
		config:  conf,
		conn:    conn,
		cursor:  cursor,
	}
}

func (w *WindowManager) Run(ctx context.Context, sub *FmtSubscriber) {
	screen := xwindow.New(w.conn, w.conn.Screen().Root)

	xgbutil.EwmhSetSupported(w.conn, []string{
		"_NET_SUPPORTED",
		"_NET_SUPPORTING_WM_CHECK",
		"_NET_ACTIVE_WINDOW",
		"_NET_CLIENT_LIST",
		"_NET_CURRENT_DESKTOP",
		"_NET_DESKTOP_NAMES",
		"_NET_NUMBER_OF_DESKTOPS",
		"_NET_WM_STATE",
		"_NET_WM_STATE_FULLSCREEN",
		"_NET_WM_WINDOW_TYPE",
		"_NET_WM_WINDOW_TYPE_DIALOG",
	})

	window := xwindow.New(w.conn, w.conn.NewId())

	window.CreateChecked(
		0,
		xproto.Window(w.conn.Screen().Root),
		0, 0, 1, 1,
		0,
		xproto.WindowClassInputOutput,
		w.conn.Screen().RootVisual,
		[]uint32{})

	xgbutil.EwmhSetSupportingWmCheck(w.conn, w.conn.Screen().Root, window.Id)
	xgbutil.EwmhSetWmName(w.conn, window.Id, "gwm")
	xgbutil.EwmhSetSupportingWmCheck(w.conn, w.conn.Screen().Root, w.conn.Screen().Root)

	for _, command := range w.config.Commands {
		grabKey(w.conn, command.Modifier, command.Keysym, w.conn.Screen().Root)
	}

	for _, action := range w.config.Actions {
		grabKey(w.conn, action.Modifier, action.Keysym, w.conn.Screen().Root)
	}

	for workspace := 1; workspace <= 9; workspace++ {
		grabKey(w.conn, w.config.WorkspaceModifier, xproto.Keycode(uint32(xproto.KeycodeXK_0)+uint32(workspace)), w.conn.Screen().Root)
		grabKey(w.conn, w.config.WorkspaceMoveWindowModifier, xproto.Keycode(uint32(xproto.KeycodeXK_0)+uint32(workspace)), w.conn.Screen().Root)
	}

	xgbutil.EwmhSetNumberOfDesktops(w.conn, 0, 9)
	xgbutil.EwmhSetCurrentDesktop(w.conn, 0, 1)

	attrs := map[xproto.AttrCW]uint32{
		xproto.CWEventMask: uint32(xproto.EventMaskSubstructureRedirect) | uint32(xproto.EventMaskSubstructureNotify),
	}

	if err := xproto.ChangeWindowAttributesChecked(w.conn, w.conn.Screen().Root, attrs).Check(); err != nil {
		panic("Unable to change window attributes. Is another window manager running?")
	}

	for _, program := range w.config.Autostart {
		if err := exec.Command(program).Start(); err != nil {
			log.Printf("Error starting program %s: %s", program, err)
		}
	}

	attrs = map[xproto.AttrCW]uint32{
		xproto.CWCursor: uint32(w.cursor),
	}

	if err := xproto.ChangeWindowAttributesChecked(w.conn, w.conn.Screen().Root, attrs).Check(); err != nil {
		panic("Unable to set cursor icon.")
	}

	log.Println("Started window manager.")

	for {
		if event, err := w.conn.WaitForEvent(); err == nil {
			go w.handle(event)
		}
	}
}
