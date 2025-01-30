package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	hook "github.com/robotn/gohook"
)

func main() {
    f, err := os.OpenFile("main.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    log.SetOutput(f)

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

	log.Println("Press Ctrl+M to toggle window minimize state")

    go func() {
        hook.Register(hook.KeyDown, []string{"ctrl", "shift", "m"}, func(e hook.Event) {
            if isMinimized {
                myWindow.Show()
                myWindow.RequestFocus()
                isMinimized = false
                log.Println("Window restored")
            } else {
                myWindow.Hide()
                isMinimized = true
                log.Println("Window minimized")
            }
        })

        s := hook.Start()
        <-hook.Process(s)
    }()

    myWindow.Canvas().Focus(entry)

    log.Println("Window show was called")
    myWindow.ShowAndRun()
}

func sendText(text string) {
    println("Sending text:", text)
}
