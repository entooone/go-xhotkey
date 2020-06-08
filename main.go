package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/entooone/hotkey/internal/xlib"
)

func createKeyEvent(display *xlib.Display, win, winRoot *xlib.Window, press bool, keycode uint, modifiers uint) *xlib.XKeyEvent {
	v := &xlib.XKeyEventValues{
		Display:    display,
		Window:     win,
		Root:       winRoot,
		SubWindow:  xlib.NewWindow(),
		Time:       0,
		X:          1,
		Y:          1,
		XRoot:      1,
		YRoot:      1,
		SameScreen: true,
		Keycode:    uint(xlib.XKeysymToKeycode(display, xlib.KeySym(keycode))),
		State:      modifiers,
	}

	if press {
		v.Type = xlib.KeyPress
	} else {
		v.Type = xlib.KeyRelease
	}

	return xlib.NewXKeyEvent(v)
}

func sendKey(display *xlib.Display, keycode uint, modifiers uint) {
	winRoot := xlib.XDefaultRootWindow(display)

	winFocus := xlib.NewWindow()
	var revert int
	xlib.XGetInputFocus(display, winFocus, &revert)

	event := createKeyEvent(display, winFocus, winRoot, true, keycode, modifiers)
	xlib.XSendEvent(display, winFocus, true, xlib.KeyPressMask, event.ToXEvent())

	event = createKeyEvent(display, winFocus, winRoot, false, keycode, modifiers)
	xlib.XSendEvent(display, winFocus, true, xlib.KeyPressMask, event.ToXEvent())
}

func AddGrubKey(key KeyInfo) {
	xlib.XGrabKey(display, xlib.XKeysymToKeycode(display, xlib.KeySym(key.keycode)), key.modifiers, grabWindow, ownerEvents, pointerMode, keyboardMode)
}

func CheckPressedKey(key KeyInfo) bool {
	pk := uint(xlib.XkbKeycodeToKeysym(display, xlib.KeyCode(event.XKey().Keycode()), 0, 0))
	pm := event.XKey().State()
	return key.keycode == pk && key.modifiers == pm
}

func SendKey(key KeyInfo) {
	sendKey(display, key.keycode, key.modifiers)
}

type HotKey interface {
	Grab()
	Send()
	Check() bool
}

type KeyInfo struct {
	keycode   uint
	modifiers uint
}

type GrabKey struct {
	grabkey KeyInfo
}

func (g *GrabKey) Grab() {
	AddGrubKey(g.grabkey)
}

type KeyConfig struct {
	GrabKey
	sendkey KeyInfo
}

func (c *KeyConfig) Send() {
	SendKey(c.sendkey)
}

func (c *KeyConfig) Check() bool {
	return CheckPressedKey(c.grabkey)
}

type ShortcutKey struct {
	grabkey     KeyInfo
	sendcommand string
}

func (c *ShortcutKey) Grab() {
	AddGrubKey(c.grabkey)
}

func (c *ShortcutKey) Send() {
	cmd := exec.Command("qdbus",
		"org.kde.kglobalaccel",
		"/component/kwin",
		"org.kde.kglobalaccel.Component.invokeShortcut",
		c.sendcommand)
	cmd.Start()
}

func (c *ShortcutKey) Check() bool {
	return CheckPressedKey(c.grabkey)
}

var registers = []HotKey{
	&KeyConfig{
		GrabKey: GrabKey{
			grabkey: KeyInfo{
				keycode:   xlib.XK_e,
				modifiers: xlib.ControlMask,
			},
		},
		sendkey: KeyInfo{
			keycode:   xlib.XK_End,
			modifiers: 0,
		},
	},
	&KeyConfig{
		GrabKey: GrabKey{
			grabkey: KeyInfo{
				keycode:   xlib.XK_a,
				modifiers: xlib.ControlMask,
			},
		},
		sendkey: KeyInfo{
			keycode:   xlib.XK_Home,
			modifiers: 0,
		},
	},
	&KeyConfig{
		GrabKey: GrabKey{
			grabkey: KeyInfo{
				keycode:   xlib.XK_h,
				modifiers: xlib.ControlMask,
			},
		},
		sendkey: KeyInfo{
			keycode:   xlib.XK_BackSpace,
			modifiers: 0,
		},
	},
	&ShortcutKey{
		grabkey: KeyInfo{
			keycode:   xlib.XK_1,
			modifiers: xlib.Mod4Mask,
		},
		sendcommand: "Switch to Desktop 1",
	},
	&ShortcutKey{
		grabkey: KeyInfo{
			keycode:   xlib.XK_2,
			modifiers: xlib.Mod4Mask,
		},
		sendcommand: "Switch to Desktop 2",
	},
	&ShortcutKey{
		grabkey: KeyInfo{
			keycode:   xlib.XK_3,
			modifiers: xlib.Mod4Mask,
		},
		sendcommand: "Switch to Desktop 3",
	},
	&ShortcutKey{
		grabkey: KeyInfo{
			keycode:   xlib.XK_4,
			modifiers: xlib.Mod4Mask,
		},
		sendcommand: "Switch to Desktop 4",
	},
	&ShortcutKey{
		grabkey: KeyInfo{
			keycode:   xlib.XK_5,
			modifiers: xlib.Mod4Mask,
		},
		sendcommand: "Switch to Desktop 5",
	},
	&ShortcutKey{
		grabkey: KeyInfo{
			keycode:   xlib.XK_6,
			modifiers: xlib.Mod4Mask,
		},
		sendcommand: "Switch to Desktop 6",
	},
	&ShortcutKey{
		grabkey: KeyInfo{
			keycode:   xlib.XK_Tab,
			modifiers: xlib.Mod4Mask,
		},
		sendcommand: "Walk Through Windows Alternative",
	},
	&ShortcutKey{
		grabkey: KeyInfo{
			keycode:   xlib.XK_Tab,
			modifiers: xlib.ShiftMask | xlib.Mod4Mask,
		},
		sendcommand: "Walk Through Windows Alternative (Reverse)",
	},
	&ShortcutKey{
		grabkey: KeyInfo{
			keycode:   xlib.XK_f,
			modifiers: xlib.Mod4Mask,
		},
		sendcommand: "Window Maximize",
	},
}

var (
	display      *xlib.Display
	grabWindow   *xlib.Window
	ownerEvents  bool
	pointerMode  int
	keyboardMode int
	event        *xlib.XEvent
)

func init() {
	display = xlib.XOpenDisplay(0)
	root := xlib.DefaultRootWindow(display)
	grabWindow = root
	ownerEvents = false
	pointerMode = xlib.GrabModeAsync
	keyboardMode = xlib.GrabModeAsync
	event = xlib.NewXEvent()
	xlib.XSelectInput(display, grabWindow, xlib.KeyPressMask)
}

func Run() error {
	var (
		isPressing = false
	)

	for _, c := range registers {
		c.Grab()
	}

	for {
		for xlib.XPending(display) {
			xlib.XNextEvent(display, event)
			switch event.Type() {
			case xlib.KeyPress:
				if !isPressing {
					fmt.Println(xlib.XGetClassHint(display, event.XAnyEvent().Window()).ResName())
					fmt.Println(xlib.XGetClassHint(display, event.XAnyEvent().Window()).ResClass())
					fmt.Println("Hot key pressed!")
					for _, c := range registers {
						if c.Check() {
							fmt.Println(c)
							c.Send()
						}
					}
					isPressing = true
				}
			case xlib.KeyRelease:
				if isPressing {
					fmt.Println("Hot key released!")
					isPressing = false
				}
			}
			time.Sleep(time.Millisecond * 5)
		}
		time.Sleep(time.Millisecond * 5)
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}
