package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"strings"
)

func main() {
	var inTE *walk.TextEdit
	var mainWin *walk.MainWindow

	MainWindow{
		AssignTo: &mainWin,
		Title:    "TTS Input",
		MinSize:  Size{Width: 250, Height: 120},
		Layout:   VBox{},
		Children: []Widget{
			TextEdit{
				AssignTo: &inTE,
				MinSize:  Size{Width: 200, Height: 40},
				MaxSize:  Size{Width: 200, Height: 40},
				Font:     Font{PointSize: 10},
				OnKeyPress: func(key walk.Key) {
					if key == walk.KeyReturn {
						sendText(inTE, mainWin)
					}
				},
			},
			PushButton{
				Text: "Send",
				MinSize: Size{Width: 80, Height: 30},
				OnClicked: func() {
					sendText(inTE, mainWin)
				},
			},
		},
	}.Run()
}

func sendText(inTE *walk.TextEdit, owner walk.Form) {
	inputText := strings.TrimSpace(inTE.Text())
	if inputText != "" {
		walk.MsgBox(owner, "Input Received", inputText, walk.MsgBoxIconInformation)
		inTE.SetText("")
	}
}
