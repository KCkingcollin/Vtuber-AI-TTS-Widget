package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
    "src/hotkey.go"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("VATTS")

	myWindow.SetFixedSize(true)
	myWindow.CenterOnScreen()

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

	myWindow.Canvas().Focus(entry)
    fmt.Println("Window show was called")
    myWindow.ShowAndRun()
}

func sendText(text string) {
	println("Sending text:", text)
}
