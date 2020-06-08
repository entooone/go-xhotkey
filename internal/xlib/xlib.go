package xlib

/*
#cgo LDFLAGS: -lX11 -lXtst
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/extensions/XTest.h>
#include <X11/keysym.h>
#include <X11/XKBlib.h>

Window defaultRootWindow(Display* dpy) {
    return DefaultRootWindow(dpy);
}
*/
import "C"
import (
	"unsafe"
)

type Window struct {
	win C.Window
}

func NewWindow() *Window {
	return &Window{
		win: C.Window(0),
	}
}

// type XWindowAttributes struct {
// 	attr *C.XWindowAttributes
// }

// func (a *XWindowAttributes) Class()  {

// }

// func XGetWindowAttributes(display *Display, w *Window) *XWindowAttributes {
// 	a := &XWindowAttributes{
// 		attr: &C.XWindowAttributes{},
// 	}
// 	C.XGetWindowAttributes(display.dpy, w.win, a.attr)
// 	return a
// }

type XClassHint struct {
	hint *C.XClassHint
}

func (h *XClassHint) ResName() string {
	return C.GoString(h.hint.res_name)
}

func (h *XClassHint) ResClass() string {
	return C.GoString(h.hint.res_class)
}

func XGetClassHint(display *Display, w *Window) *XClassHint {
	h := &XClassHint{
		hint: &C.XClassHint{},
	}
	C.XGetClassHint(display.dpy, w.win, h.hint)
	return h
}

type Display struct {
	dpy *C.Display
}

type XKeyEvent struct {
	ev *C.XKeyEvent
}

func (e *XKeyEvent) State() uint {
	return uint(e.ev.state)
}

func (e *XKeyEvent) Keycode() uint {
	return uint(e.ev.keycode)
}

func (e *XKeyEvent) ToXEvent() *XEvent {
	return &XEvent{
		ev: (*C.XEvent)(unsafe.Pointer(e.ev)),
	}
}

type XKeyEventValues struct {
	Type         int
	Display      *Display
	Window       *Window
	Root         *Window
	SubWindow    *Window
	Time         uint64
	X, Y         int
	XRoot, YRoot int
	State        uint
	Keycode      uint
	SameScreen   bool
}

func NewXKeyEvent(values *XKeyEventValues) *XKeyEvent {
	var ss C.Bool
	if values.SameScreen {
		ss = C.True
	} else {
		ss = C.False
	}
	ev := &C.XKeyEvent{
		_type:       C.int(values.Type),
		display:     values.Display.dpy,
		window:      values.Window.win,
		root:        values.Root.win,
		subwindow:   values.SubWindow.win,
		time:        C.Time(values.Time),
		x:           C.int(values.X),
		y:           C.int(values.Y),
		x_root:      C.int(values.XRoot),
		y_root:      C.int(values.YRoot),
		state:       C.uint(values.State),
		keycode:     C.uint(values.Keycode),
		same_screen: ss,
	}
	return &XKeyEvent{
		ev: ev,
	}
}

type XAnyEvent struct {
	ev *C.XAnyEvent
}

func (e *XAnyEvent) Window() *Window {
	w := &Window{
		win: C.Window(0),
	}
	w.win = e.ev.window
	return w
}

type XEvent struct {
	ev *C.XEvent
}

func (e *XEvent) Type() int {
	t := (*C.int)(unsafe.Pointer(e.ev))
	return int(*t)
}

func (e *XEvent) XAnyEvent() *XAnyEvent {
	any := (*C.XAnyEvent)(unsafe.Pointer(e.ev))
	return &XAnyEvent{
		ev: any,
	}
}

func (e *XEvent) XKey() *XKeyEvent {
	xkey := (*C.XKeyEvent)(unsafe.Pointer(e.ev))
	return &XKeyEvent{
		ev: xkey,
	}
}

func NewXEvent() *XEvent {
	return &XEvent{
		ev: &C.XEvent{},
	}
}

func XOpenDisplay(x int) *Display {
	c := C.char(x)
	dpy := C.XOpenDisplay((*C.char)(unsafe.Pointer(&c)))
	return &Display{
		dpy: dpy,
	}
}

func DefaultRootWindow(display *Display) *Window {
	w := C.defaultRootWindow(display.dpy)
	return &Window{
		win: w,
	}
}

func XDefaultRootWindow(display *Display) *Window {
	w := C.XDefaultRootWindow(display.dpy)
	return &Window{
		win: w,
	}
}

func XKeysymToKeycode(display *Display, key KeySym) KeyCode {
	return KeyCode(C.XKeysymToKeycode(display.dpy, C.ulong(key)))
}

func XkbKeycodeToKeysym(display *Display, kc KeyCode, group, level int) KeySym {
	return KeySym(C.XkbKeycodeToKeysym(display.dpy, C.uchar(kc), C.int(group), C.int(level)))
}

func XGrabKey(display *Display, keycode KeyCode, modifiers uint, grabWindow *Window, ownerEvents bool, pointerMode int, keyboardMode int) {
	var cOwnerEvents C.int
	if ownerEvents {
		cOwnerEvents = C.int(1)
	}
	C.XGrabKey(display.dpy, C.int(keycode), C.uint(modifiers), grabWindow.win, cOwnerEvents, C.int(pointerMode), C.int(keyboardMode))
}

func XSelectInput(display *Display, window *Window, mask int64) {
	C.XSelectInput(display.dpy, window.win, C.long(mask))
}

func XPending(display *Display) bool {
	return C.XPending(display.dpy) != C.int(0)
}

func XNextEvent(display *Display, event *XEvent) {
	C.XNextEvent(display.dpy, event.ev)
}

func XGetInputFocus(display *Display, window *Window, revert *int) {
	r := (*C.int)(unsafe.Pointer(revert))
	w := (*C.Window)(unsafe.Pointer(&window.win))
	C.XGetInputFocus(display.dpy, w, r)
}

func XSendEvent(display *Display, window *Window, propagate bool, eventMask int64, eventSend *XEvent) {
	var cp C.Bool
	if propagate {
		cp = C.True
	} else {
		cp = C.False
	}
	C.XSendEvent(display.dpy, window.win, cp, C.long(eventMask), eventSend.ev)
}
