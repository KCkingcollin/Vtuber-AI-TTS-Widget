package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/KCkingcollin/Vtuber-AI-TTS-Widget/src"
	hook "github.com/robotn/gohook"
)

var verbose bool

func main() {
    if len(os.Args) > 1 {
        if os.Args[1] == "-v" || os.Args[1] == "--verbose" {
            verbose = true
        }
    }
    err := log.Init("./", "main.log", 1)
    if err != nil {
        panic(err)
    }

    myApp := app.New()
    myWindow := myApp.NewWindow("VATTS")
    
    entry := widget.NewEntry()
    entry.SetPlaceHolder("TTS Input...")
    
    button := widget.NewButton(" Press to send (or just press enter) ", func() {
        sendText(entry.Text)
        entry.SetText("")
    })
    
    entry.OnSubmitted = func(text string) {
        sendText(text)
        entry.SetText("")
    }
    
    content := container.NewVBox(entry, button)
    myWindow.SetContent(content)
    myWindow.SetFixedSize(true)
    myWindow.CenterOnScreen()
    
	isMinimized := false

	log.Append("Press Ctrl+M to toggle window minimize state", verbose)

    go func() {
        hook.Register(hook.KeyDown, []string{"ctrl", "shift", "m"}, func(e hook.Event) {
            if isMinimized {
                myWindow.Show()
                myWindow.RequestFocus()
                isMinimized = false
                log.Append("Window restored", verbose)
            } else {
                myWindow.Hide()
                isMinimized = true
                log.Append("Window minimized", verbose)
            }
        })

        s := hook.Start()
        <-hook.Process(s)
    }()

    myWindow.Canvas().Focus(entry)

    log.Append("Window show was called", verbose)
    myWindow.ShowAndRun()
}

func sendText(text string) {
    log.Append(fmt.Sprintf("Sending text: %s", text), verbose)
}
