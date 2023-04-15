# go-xhotkey

[![Go Reference](https://pkg.go.dev/badge/github.com/entooone/go-xhotkey.svg)](https://pkg.go.dev/github.com/entooone/go-xhotkey)

This Go library introduces a HotKey system for the X Window System. By using this library, you can capture key events and call arbitrary functions defined in Go, or remap them to other keys. This library was developed to make it easy to implement customizable keyboard shortcuts.

## Usage

```go
import "github.com/entooone/go-xhotkey"

func main() {
    // Create a new option.
    option := &xhotkey.Option{
        HotKeys: []xhotkey.HotKey{
            // Remap Ctrl + e to End.
            xhotkey.NewRemap( 
                xhotkey.KeyInfo{
                KeyCode:   xhotkey.KeyCode_e,
                Modifiers: xhotkey.ModifierControl,
            },
            xhotkey.KeyInfo{
                KeyCode:   xhotkey.KeyCode_End,
                Modifiers: 0,
            }),

            // Ctrl + Shift + a prints "Ctrl + Shift + a" to stdout.
            xhotkey.NewHotKey(
                xhotkey.KeyInfo{
                KeyCode:   xhotkey.KeyCode_a,
                Modifiers: xhotkey.ModifierControl|xhotkey.ModifierShift,
                },
            func() {
                fmt.Println("Ctrl + Shift + a")
            }),
        },

        // If true, the library will print debug messages to stdout.
        IsDebug: true,
    }

    if err := xhotkey.Run(option); err != nil {
        panic(err)
    }
}
```

## References

1. <https://stackoverflow.com/questions/30156202/xlib-c-get-window-handle-sendevent>
1. <https://stackoverflow.com/questions/48962411/when-sending-an-xkeyevent-with-xsendevent-the-system-doesnt-respond-the-first>
1. <https://lists.x.org/archives/xorg/2010-October/051373.html>
