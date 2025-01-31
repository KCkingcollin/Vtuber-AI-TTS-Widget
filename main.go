package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/KCkingcollin/Vtuber-AI-TTS-Widget/src"
	"github.com/go-vgo/robotgo"
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
    
    howMany, _ := robotgo.FindIds("main")
    if len(howMany) > 1 {
        log.Append("To many processes", verbose)
        error := widget.NewLabel("Error: To many processes")
        text := widget.NewLabel("Just throw in a whole banana (peal included)?")
        yes := widget.NewButton("YES", func() {
            for i := 0; i < len(howMany); i++ {
                process, err := os.FindProcess(howMany[i])
                if err != nil {
                    log.Append(fmt.Sprintf("%e", err), verbose)
                }
                process.Kill()
            }
        })
        no := widget.NewButton("no", func() {
            os.Exit(0)
        })
        text.Alignment = fyne.TextAlignCenter
        error.Alignment = fyne.TextAlignCenter
        buttons := container.NewHBox(yes, no)
        centeredButtons := container.NewCenter(buttons)
        content := container.NewVBox(error, text, centeredButtons)
        myWindow.SetContent(container.NewVBox(content))
        myWindow.ShowAndRun()
    }

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
    myWindow.CenterOnScreen()

    isHidden := false

	log.Append("Press Ctrl+M to toggle window minimize state", verbose)


    go func() {
        hook.Register(hook.KeyDown, []string{"ctrl", "shift", "m"}, func(e hook.Event) {
            if isHidden {
                myWindow.Show()
                myWindow.RequestFocus()
                isHidden = false
                log.Append("Window shown", verbose)
            } else {
                myWindow.Hide()
                isHidden = true
                log.Append("Window hidden", verbose)
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
