package xhotkey

import (
	"fmt"
	"time"

	"github.com/entooone/xhotkey/internal/xlib"
)

var (
	display      *xlib.Display
	grabWindow   *xlib.Window
	ownerEvents  bool
	pointerMode  int
	keyboardMode int
	event        *xlib.XEvent
)

var registerd = make(map[HotKey]struct{})

type KeyInfo struct {
	KeyCode   uint
	Modifiers uint
}

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

func SendKey(key KeyInfo) {
	sendKey(display, key.KeyCode, key.Modifiers)
}

func addGrabKey(key KeyInfo) {
	xlib.XGrabKey(display, xlib.XKeysymToKeycode(display, xlib.KeySym(key.KeyCode)), key.Modifiers, grabWindow, ownerEvents, pointerMode, keyboardMode)
}

func isPressKey(key KeyInfo) bool {
	pk := uint(xlib.XkbKeycodeToKeysym(display, xlib.KeyCode(event.XKey().Keycode()), 0, 0))
	pm := event.XKey().State()
	return key.KeyCode == pk && key.Modifiers == pm
}

type HotKey interface {
	GrabKey() KeyInfo
	Execute()
}

func RegistHotKey(hotkey HotKey) {
	registerd[hotkey] = struct{}{}
}

type KeyMap struct {
	grabkey KeyInfo
	sendkey KeyInfo
}

func NewKeyMap(grabkey, sendkey KeyInfo) KeyMap {
	return KeyMap{
		grabkey: grabkey,
		sendkey: sendkey,
	}
}

func (m KeyMap) GrabKey() KeyInfo {
	return m.grabkey
}

func (m KeyMap) Execute() {
	SendKey(m.sendkey)
}

func Run() error {
	var isPressing bool

	for r := range registerd {
		addGrabKey(r.GrabKey())
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
					for r := range registerd {
						if isPressKey(r.GrabKey()) {
							fmt.Println(r)
							r.Execute()
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
