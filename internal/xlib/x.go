package xlib

type KeyCode uint

const (
	KeyPressMask   = (1 << 0)
	KeyReleaseMask = (1 << 1)

	KeyPress   = 2
	KeyRelease = 3

	GrabModeSync  = 0
	GrabModeAsync = 1

	ShiftMask   = (1 << 0)
	LockMask    = (1 << 1)
	ControlMask = (1 << 2)
	Mod1Mask    = (1 << 3)
	Mod2Mask    = (1 << 4)
	Mod3Mask    = (1 << 5)
	Mod4Mask    = (1 << 6)
	Mod5Mask    = (1 << 7)
)
