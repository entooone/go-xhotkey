package main

import (
	"fmt"
	"time"

	"github.com/entooone/hotkey/internal/xlib"
)

func createKeyEvent(display *xlib.Display, win, winRoot *xlib.Window, press bool, keycode int, modifiers uint) *xlib.XKeyEvent {
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

func sendKey(display *xlib.Display, keycode int, modifiers uint) {
	winRoot := xlib.XDefaultRootWindow(display)

	winFocus := xlib.NewWindow()
	var revert int
	xlib.XGetInputFocus(display, winFocus, &revert)

	fmt.Printf("%t\n", winFocus)

	event := createKeyEvent(display, winFocus, winRoot, true, keycode, modifiers)
	xlib.XSendEvent(display, winFocus, true, xlib.KeyPressMask, event.ToXEvent())

	event = createKeyEvent(display, winFocus, winRoot, false, keycode, modifiers)
	xlib.XSendEvent(display, winFocus, true, xlib.KeyPressMask, event.ToXEvent())
}

func main() {
	var (
		dpy  = xlib.XOpenDisplay(0)
		root = xlib.DefaultRootWindow(dpy)
		ev   = xlib.NewXEvent()

		modifiers    uint = xlib.ControlMask
		keycode           = xlib.XKeysymToKeycode(dpy, xlib.XK_F)
		grabWindow        = root
		ownerEvents       = false
		pointerMode       = xlib.GrabModeAsync
		keyboardMode      = xlib.GrabModeAsync

		isPressing = false
	)
	xlib.XGrabKey(dpy, keycode, modifiers, grabWindow, ownerEvents, pointerMode, keyboardMode)
	xlib.XGrabKey(dpy, keycode+2, xlib.ControlMask+xlib.ShiftMask, grabWindow, ownerEvents, pointerMode, keyboardMode)
	xlib.XSelectInput(dpy, root, xlib.KeyPressMask)

	for {
		for xlib.XPending(dpy) {
			fmt.Println("Loop!!")
			xlib.XNextEvent(dpy, ev)

			switch ev.Type() {
			case xlib.KeyPress:
				if !isPressing {
					fmt.Println("Hot key pressed!")
					sendKey(dpy, xlib.XK_A, xlib.ControlMask)
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
}
