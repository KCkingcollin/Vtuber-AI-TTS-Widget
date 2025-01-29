package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("VATTS")

    myWindow.Resize(fyne.NewSize(300, 76))
	myWindow.SetFixedSize(true)
	myWindow.SetPadded(false)
	myWindow.CenterOnScreen()
	myWindow.SetMaster()

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

	myWindow.Show()
	myWindow.Canvas().Focus(entry)

	myApp.Run()
}

func sendText(text string) {
	println("Sending text:", text)
}
