package main

import (
	"github.com/entooone/go-xhotkey"
)

func run() error {
	option := &xhotkey.Option{
		HotKeys: []xhotkey.HotKey{
			xhotkey.NewRemap(
				xhotkey.KeyInfo{
					KeyCode:   xhotkey.KeyCode_e,
					Modifiers: xhotkey.ModifierControl,
				},
				xhotkey.KeyInfo{
					KeyCode:   xhotkey.KeyCode_End,
					Modifiers: 0,
				}),
			xhotkey.NewRemap(
				xhotkey.KeyInfo{
					KeyCode:   xhotkey.KeyCode_a,
					Modifiers: xhotkey.ModifierControl,
				},
				xhotkey.KeyInfo{
					KeyCode:   xhotkey.KeyCode_Home,
					Modifiers: 0,
				}),
			xhotkey.NewRemap(
				xhotkey.KeyInfo{
					KeyCode:   xhotkey.KeyCode_h,
					Modifiers: xhotkey.ModifierControl,
				},
				xhotkey.KeyInfo{
					KeyCode:   xhotkey.KeyCode_BackSpace,
					Modifiers: 0,
				}),
		},
		IsDebug: true,
	}

	if err := xhotkey.Run(option); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
