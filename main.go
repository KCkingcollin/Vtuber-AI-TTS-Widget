package main

import (
	"fmt"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
    var inTE *walk.TextEdit

    MainWindow{
        Title:   "Text Input Example",
        MinSize: Size{Width: 100, Height: 75},
        Layout:  VBox{},
        Children: []Widget{
            TextEdit{
                AssignTo: &inTE,
                OnKeyPress: func(key walk.Key) {
                    if key == walk.KeyReturn {
                        inputText := strings.TrimSpace(inTE.Text())
                        if inputText != "" {
                            fmt.Println(inputText)
                            inTE.SetText("")
                        }
                    }
                },
            },
        },
    }.Run()
}
