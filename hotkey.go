package xhotkey

// HotKey is an interface for hotkey.
type HotKey interface {
	GrabKey() KeyInfo
	Execute()
}

type hotkey struct {
	grabKey KeyInfo
	exec    func()
}

func (h *hotkey) GrabKey() KeyInfo {
	return h.grabKey
}

func (h *hotkey) Execute() {
	h.exec()
}

// NewHotKey returns a new HotKey. keyinfo is a key to grab, exec is a function to execute when the key is pressed.
func NewHotKey(keyinfo KeyInfo, exec func()) HotKey {
	return &hotkey{
		grabKey: keyinfo,
		exec:    exec,
	}
}

// NewRemap returns a new HotKey. grabkey is a key to grab, sendkey is a key to send when the grabkey is pressed.
func NewRemap(grabkey, sendkey KeyInfo) HotKey {
	return NewHotKey(grabkey, func() {
		sendKey(sendkey)
	})
}
