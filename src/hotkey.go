package hotKeys

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/robotn/gohook"
)

func HotKeys() {
	// The title of the window you want to control
	targetTitle := "VATTS"
	isMinimized := false

	fmt.Printf("Looking for window with title: %s\n", targetTitle)
	fmt.Println("Press Ctrl+Alt+M to toggle window minimize state")
	fmt.Println("Press Ctrl+C to exit")

	// Register the hotkey
	hook.Register(hook.KeyDown, []string{"alt", "ctrl", "m"}, func(e hook.Event) {
		// Find the window by title
		pid, _ := robotgo.FindIds(targetTitle)
		if pid[0] == 0 {
			fmt.Printf("Window with title '%s' not found\n", targetTitle)
			return
		}

		if isMinimized {
			robotgo.MaxWindow(pid[0])
			isMinimized = false
			fmt.Println("Window restored")
		} else {
			robotgo.MinWindow(pid[0])
			isMinimized = true
			fmt.Println("Window minimized")
		}
	})

	// Keep the program running
	s := hook.Start()
	<-hook.Process(s)
}
