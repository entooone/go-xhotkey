package main

import (
	"log"

	"github.com/entooone/go-xhotkey"
	"github.com/entooone/go-xhotkey/internal/xlib"
)

func run() error {
	keymaps := []xhotkey.KeyMap{
		xhotkey.NewKeyMap(
			xhotkey.KeyInfo{
				KeyCode:   xlib.XK_e,
				Modifiers: xlib.ControlMask,
			},
			xhotkey.KeyInfo{
				KeyCode:   xlib.XK_End,
				Modifiers: 0,
			}),
		xhotkey.NewKeyMap(
			xhotkey.KeyInfo{
				KeyCode:   xlib.XK_a,
				Modifiers: xlib.ControlMask,
			},
			xhotkey.KeyInfo{
				KeyCode:   xlib.XK_Home,
				Modifiers: 0,
			}),
		xhotkey.NewKeyMap(
			xhotkey.KeyInfo{
				KeyCode:   xlib.XK_h,
				Modifiers: xlib.ControlMask,
			},
			xhotkey.KeyInfo{
				KeyCode:   xlib.XK_BackSpace,
				Modifiers: 0,
			}),
		xhotkey.NewKeyMap(
			xhotkey.KeyInfo{
				KeyCode:   xlib.XK_a,
				Modifiers: xlib.Mod4Mask,
			},
			xhotkey.KeyInfo{
				KeyCode:   xlib.XK_a,
				Modifiers: xlib.ControlMask,
			}),
	}

	for _, m := range keymaps {
		xhotkey.RegistHotKey(m)
	}

	if err := xhotkey.Run(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
