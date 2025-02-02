package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/KCkingcollin/Vtuber-AI-TTS-Widget/src"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

var verbose bool
var conn net.Conn

func downloadFile(url, filepath string) error {
	// Step 1: Send an HTTP GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Step 2: Check if the response status code is OK (200)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Step 3: Create the file to save the downloaded content
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	// Step 4: Copy the response body (file content) to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	log.Append(fmt.Sprintf("File downloaded successfully: %s\n", filepath), verbose)
	return nil
}

func checkForDepends() error {
    // https://github.com/thewh1teagle/kokoro-onnx/releases/download/model-files/voices.json -o kokoro-tts/voices.json     
    // https://github.com/thewh1teagle/kokoro-onnx/releases/download/model-files/kokoro-v0_19.onnx -o kokoro-tts/kokoro-v0_19.onnx
    // https://github.com/thewh1teagle/kokoro-onnx/releases/download/model-files/voices.bin -o kokoro-tts/voices.bin

	_, err := os.Stat("./kokoro-tts/voices.json")
	if os.IsNotExist(err) {
        err := downloadFile("https://github.com/thewh1teagle/kokoro-onnx/releases/download/model-files/voices.json", "kokoro-tts/voices.json")
        if err != nil {
            return fmt.Errorf("failed to download file: %v", err)
        }
	}

	_, err = os.Stat("./kokoro-tts/voices.bin")
	if os.IsNotExist(err) {
        err := downloadFile("https://github.com/thewh1teagle/kokoro-onnx/releases/download/model-files/voices.bin", "kokoro-tts/voices.bin")
        if err != nil {
            return fmt.Errorf("failed to download file: %v", err)
        }
	}

	_, err = os.Stat("./kokoro-tts/voice.onnx")
	if os.IsNotExist(err) {
        err := downloadFile("https://github.com/thewh1teagle/kokoro-onnx/releases/download/model-files/kokoro-v0_19.onnx", "kokoro-tts/voice.onnx")
        if err != nil {
            return fmt.Errorf("failed to download file: %v", err)
        }
	}

    // python3 -m venv kokoro-tts/venv
    // kokoro-tts/venv/bin/pip install onnxruntime kokoro_onnx sounddevice numpy psutil
	_, err = os.Stat("./kokoro-tts/venv")
	if os.IsNotExist(err) {
        cmd := exec.Command("python3", "-m", "venv", "kokoro-tts/venv")

        var out bytes.Buffer
        cmd.Stdout = &out
        cmd.Stderr = &out

        err := cmd.Run()
        if err != nil {
            return fmt.Errorf("failed to run command: %v", err)
        }

        log.Append(fmt.Sprintf("Command output:\n%s\n", out.String()), verbose)

        cmd = exec.Command("kokoro-tts/venv/bin/pip", "install", "onnxruntime", "kokoro_onnx", "sounddevice", "numpy", "psutil")

        cmd.Stdout = &out
        cmd.Stderr = &out

        err = cmd.Run()
        if err != nil {
            return fmt.Errorf("failed to run command: %v", err)
        }

        log.Append(fmt.Sprintf("Command output:\n%s\n", out.String()), verbose)
	}

    return nil
}

func main() {
    if len(os.Args) > 1 {
        if os.Args[1] == "-v" || os.Args[1] == "--verbose" {
            verbose = true
        }
    }
    if err := log.Init("./", "main.log", 1); err != nil {
        panic(err)
    }
    if err := checkForDepends(); err != nil {
        log.Append(fmt.Sprintf("%e", err), verbose)
    }

    server := exec.Command("kokoro-tts/venv/bin/python", "kokoro-tts/tts-server.py")
    server.Start()

    var err error
    conn, err = net.Dial("tcp", "localhost:65432")
    for i := 0; err != nil; i++ {
        time.Sleep(time.Second)
        conn, err = net.Dial("tcp", "localhost:65432")
        if i >= 4 {
            log.Append(fmt.Sprintf("Error connecting: %e", err), verbose)
            return
        }
    }
    defer conn.Close()

    myApp := app.New()
    myWindow := myApp.NewWindow("VATTS")
    
    entry := widget.NewEntry()
    entry.SetPlaceHolder("TTS Input...")
    
    howMany, _ := robotgo.FindIds("VATTS")
    if len(howMany) > 1 {
        log.Append("To many processes", verbose)
        error := widget.NewLabel("Error: To many processes")
        text := widget.NewLabel("Just throw in a whole banana (peal included)?")
        howMany, _ = robotgo.FindIds("VATTS")
        yes := widget.NewButton("YES", func() {
            for true {
                howMany, _ = robotgo.FindIds("VATTS")
                if len(howMany) <= 1 {
                    message := widget.NewLabel("VATTS apps air subscription revoked")
                    popup := widget.NewModalPopUp(
                        container.NewVBox(message),
                        myWindow.Canvas(),
                    )
                    popup.Show()
                    time.Sleep(time.Second*3)
                    os.Exit(0)
                } else {
                    howMany, _ = robotgo.FindIds("VATTS")
                    for i := 0; i < len(howMany); i++ {
                        if os.Getpid() != howMany[i] {
                            process, err := os.FindProcess(howMany[i])
                            if err != nil {
                                log.Append(fmt.Sprintf("%e", err), verbose)
                            }
                            process.Kill()
                        }
                    }
                }
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
		// Send text to server
        _, err := fmt.Fprintf(conn, "%s\n", text)
		if err != nil {
			log.Append(fmt.Sprintf("Error sending text: %e", err), verbose)
			return
		}

		// Wait for acknowledgment
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Append(fmt.Sprintf("Error reading response: %e", err), verbose)
			return
		}
        log.Append(fmt.Sprintf("Server response: %s", response), verbose)
}
