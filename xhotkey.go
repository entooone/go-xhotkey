package xhotkey

import (
	"os"
	"time"

	"github.com/entooone/go-xhotkey/internal/xlib"
	"golang.org/x/exp/slog"
)

var (
	display      *xlib.Display
	grabWindow   *xlib.Window
	ownerEvents  bool
	pointerMode  int
	keyboardMode int
	event        *xlib.XEvent
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

func sendKeyToDisplay(display *xlib.Display, keycode uint, modifiers uint) {
	winRoot := xlib.XDefaultRootWindow(display)

	winFocus := xlib.NewWindow()
	var revert int
	xlib.XGetInputFocus(display, winFocus, &revert)

	event := createKeyEvent(display, winFocus, winRoot, true, keycode, modifiers)
	xlib.XSendEvent(display, winFocus, true, xlib.KeyPressMask, event.ToXEvent())

	event = createKeyEvent(display, winFocus, winRoot, false, keycode, modifiers)
	xlib.XSendEvent(display, winFocus, true, xlib.KeyPressMask, event.ToXEvent())
}

type KeyInfo struct {
	KeyCode   uint
	Modifiers uint
}

func sendKey(key KeyInfo) {
	slog.Debug("sendkey:", "keycode", key.KeyCode, "modifiers", key.Modifiers)
	sendKeyToDisplay(display, key.KeyCode, key.Modifiers)
}

func addGrabKey(key KeyInfo) {
	xlib.XGrabKey(display, xlib.XKeysymToKeycode(display, xlib.KeySym(key.KeyCode)), key.Modifiers, grabWindow, ownerEvents, pointerMode, keyboardMode)
}

func isPressKey(key KeyInfo) bool {
	pk := uint(xlib.XkbKeycodeToKeysym(display, xlib.KeyCode(event.XKey().Keycode()), 0, 0))
	pm := event.XKey().State()
	return key.KeyCode == pk && key.Modifiers == pm
}

type Option struct {
	HotKeys []HotKey
	IsDebug bool
}

func Run(option *Option) error {
	level := slog.LevelInfo
	if option.IsDebug {
		level = slog.LevelDebug
	}

	handler := slog.HandlerOptions{
		Level: level,
	}.NewJSONHandler(os.Stdout)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	var isPressing bool

	for _, r := range option.HotKeys {
		addGrabKey(r.GrabKey())
	}

	for {
		for xlib.XPending(display) {
			xlib.XNextEvent(display, event)
			switch event.Type() {
			case xlib.KeyPress:
				if !isPressing {
					slog.Debug("hotkey pressed")
					for _, r := range option.HotKeys {
						if isPressKey(r.GrabKey()) {
							kinfo := r.GrabKey()
							slog.Debug("grubkey:", "keycode", kinfo.KeyCode, "modifiers", kinfo.Modifiers)
							r.Execute()
						}
					}
					isPressing = true
				}
			case xlib.KeyRelease:
				if isPressing {
					slog.Debug("hotkey released")
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
