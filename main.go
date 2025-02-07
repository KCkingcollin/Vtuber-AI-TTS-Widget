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
	"runtime"
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
var osEnv string
var mainApp fyne.App
var mainWindow fyne.Window

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to send GET request: %e", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %e", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %e", err)
	}

	log.Append(fmt.Sprintf("File downloaded successfully: %s", filepath), verbose)
	return nil
}

func checkForDepends() {
    _, err := os.Stat("./kokoro-tts/voices.json")
    if installLock("voices.json.loc") || os.IsNotExist(err) {
        installLock("voices.json.loc", true)
        mainWindow = textPopupWindow(mainWindow, "Downloading voices.json\nPlease do not close this window")
        mainWindow.Content().Refresh()
        err := downloadFile("https://github.com/thewh1teagle/kokoro-onnx/releases/download/model-files/voices.json", "kokoro-tts/voices.json")
        if err != nil {
            log.Append(fmt.Sprintf("failed to download file: %e", err), true)
            os.Exit(1)
        }
        installLock("voices.json.loc", false)
    }

    _, err = os.Stat("./kokoro-tts/voices.bin")
    if installLock("voices.bin.loc") || os.IsNotExist(err) {
        installLock("voices.bin.loc", true)
        mainWindow = textPopupWindow(mainWindow, "Downloading voices.bin\nPlease do not close this window")
        mainWindow.Content().Refresh()
        err := downloadFile("https://github.com/thewh1teagle/kokoro-onnx/releases/download/model-files/voices.bin", "kokoro-tts/voices.bin")
        if err != nil {
            log.Append(fmt.Sprintf("failed to download file: %e", err), true)
            os.Exit(1)
        }
        installLock("voices.bin.loc", false)
    }

    _, err = os.Stat("./kokoro-tts/voice.onnx")
    if installLock("voice.onnx.loc") || os.IsNotExist(err) {
        installLock("voice.onnx.loc", true)
        mainWindow = textPopupWindow(mainWindow, "Downloading voice.onnx\nPlease do not close this window")
        mainWindow.Content().Refresh()
        err := downloadFile("https://github.com/thewh1teagle/kokoro-onnx/releases/download/model-files/kokoro-v0_19.onnx", "kokoro-tts/voice.onnx")
        if err != nil {
            log.Append(fmt.Sprintf("failed to download file: %e", err), true)
            os.Exit(1)
        }
        installLock("voice.onnx.loc", false)
    }

    mainWindow = textPopupWindow(mainWindow, "Setting up python environment\nPlease do not close this window")
    mainWindow.Content().Refresh()

    if osEnv == "linux" {
        _, err = os.Stat("kokoro-tts/venv/bin/tts-server-py-bin")
    } else if osEnv == "windows" {
        _, err = os.Stat("kokoro-tts/venv/Scripts/tts-server-py-bin.exe")
    } else { 
        log.Append("os not supported", verbose)
        os.Exit(1)
    }

    venvInstall := os.IsNotExist(err) 

    if installLock("pipvenv.loc") || venvInstall {
        installLock("pipvenv.loc", true)
        cmd := exec.Command("python", "-m", "venv", "kokoro-tts/venv")

        var out bytes.Buffer
        cmd.Stdout = &out
        cmd.Stderr = &out

        err := cmd.Run()
        if err != nil {
            log.Append(fmt.Sprintf("failed to run command: %e", err), true)
            os.Exit(1)
        }

        log.Append(fmt.Sprintf("Command output: %s", out.String()), verbose)

        var sourceFile string
        var destFile string
        if osEnv == "linux" {
            sourceFile = "kokoro-tts/venv/bin/python"
            destFile = "kokoro-tts/venv/bin/tts-server-py-bin"
        } else if osEnv == "windows" {
            sourceFile = "kokoro-tts/venv/Scripts/pythonw.exe"
            destFile = "kokoro-tts/venv/Scripts/tts-server-py-bin.exe"
        } else { 
            log.Append("os not supported", true)
            os.Exit(1)
        }

        from, err := os.Open(sourceFile)
        if err != nil {
            error := fmt.Sprint(err)
            log.Append(error, true)
        }
        defer from.Close()

        to, err := os.OpenFile(destFile, os.O_RDWR|os.O_CREATE, 0755)
        if err != nil {
            error := fmt.Sprint(err)
            log.Append(error, true)
        }
        defer to.Close()

        _, err = io.Copy(to, from)
        if err != nil {
            error := fmt.Sprint(err)
            log.Append(error, true)
        }
        installLock("pipvenv.loc", false)
        time.Sleep(time.Second)
    }

    if installLock("pipdepend.loc") || venvInstall {
        installLock("pipdepend.loc", true)
        mainWindow = textPopupWindow(mainWindow, "Installing python dependencies\nThis could take a while, do not close this or the terminal window")
        mainWindow.Content().Refresh()

        var cmd *exec.Cmd
        if osEnv == "linux" {
            cmd = exec.Command("kokoro-tts/venv/bin/pip", "install", "onnxruntime", "kokoro_onnx", "sounddevice", "numpy", "psutil")
        } else if osEnv == "windows" {
            cmd = exec.Command("kokoro-tts/venv/Scripts/pip.exe", "install", "onnxruntime", "kokoro_onnx", "sounddevice", "numpy", "psutil")
        } else { 
            err := fmt.Errorf("os not supported")
            log.Append(fmt.Sprint(err), verbose)
            panic(err)
        }

        var out bytes.Buffer
        cmd.Stdout = &out
        cmd.Stderr = &out

        err = cmd.Run()
        if err != nil {
            log.Append(fmt.Sprintf("failed to run command: %e", err), verbose)
            os.Exit(1)
        }

        log.Append(fmt.Sprintf("Command output:%s", out.String()), verbose)
        installLock("pipdepend.loc", false)
        time.Sleep(time.Second*5)
    }
}

func textPopupWindow(textWindow fyne.Window, text string) fyne.Window {
    label := widget.NewLabel(text)
    label.Alignment = fyne.TextAlignCenter
    content := container.NewVBox(label)
    textWindow.SetContent(container.NewVBox(content))
    return textWindow
}

func installLock(name string, enable ...bool) bool {
    if len(enable) < 1 {
        _, err := os.Stat(name)
        if os.IsNotExist(err) {return false} else {return true}
    } else {
        if enable[0] {
            file, err := os.Create(name)
            if err != nil {
                log.Append(fmt.Sprint(err), verbose)
            }
            defer file.Close()
        } else {
            os.Remove(name)
        }
    }
    return false
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

    mainApp = app.New()
    mainWindow = mainApp.NewWindow("VATTS")

    var server *exec.Cmd
    go func() {
        osEnv = runtime.GOOS
        log.Append(fmt.Sprintf("Operating system: %s", osEnv), verbose)

        log.Append("Checking for dependencies", verbose)
        checkForDepends()
        log.Append("Dependencies installed", verbose)

        mainWindow = textPopupWindow(mainWindow, "Starting TTS engine")
        mainWindow.Content().Refresh()
        log.Append("Starting server", verbose)
        if osEnv == "linux" {
            server = exec.Command("kokoro-tts/venv/bin/tts-server-py-bin", "kokoro-tts/tts-server.py")
            err := server.Start()
            if err != nil {
                log.Append(fmt.Sprintf("failed to run server: %e", err), true)
                os.Exit(1)
            }
        } else if osEnv == "windows" {
            server = exec.Command("kokoro-tts/venv/Scripts/tts-server-py-bin.exe", "kokoro-tts/tts-server.py")
            err := server.Start()
            if err != nil {
                log.Append(fmt.Sprintf("failed to run server: %e", err), true)
                os.Exit(1)
            }
        } else {
            err := fmt.Errorf("os not supported")
            log.Append(fmt.Sprint(err), true)
            os.Exit(1)
        }
        log.Append("Server started", verbose)

        log.Append("Connecting to TTS server", verbose)
        var err error
        conn, err = net.Dial("tcp", "localhost:65432")
        for i := 0; err != nil; i++ {
            time.Sleep(time.Second)
            conn, err = net.Dial("tcp", "localhost:65432")
            if i >= 30 {
                log.Append(fmt.Sprintf("Error connecting: %e", err), true)
                os.Exit(1)
            }
        }
        log.Append("Connected to TTS server", verbose)

        mainWindow = textPopupWindow(mainWindow, "Looking for other running VATTS apps")
        mainWindow.Content().Refresh()
        log.Append("Looking for running VATTS apps", verbose)
        howMany, _ := robotgo.FindIds("VATTS")
        serverIDs, _ := robotgo.FindIds("tts-server-py-bin")
        if len(howMany) > 1 || len(serverIDs) > 1 {
            log.Append("To many processes", verbose)
            error := widget.NewLabel("Error: To many processes")
            text := widget.NewLabel("Just throw in a whole banana (peal included)?")
            howMany, _ = robotgo.FindIds("VATTS")
            serverIDs, _ = robotgo.FindIds("tts-server-py-bin")
            yes := widget.NewButton("YES", func() {
                for {
                    howMany, _ = robotgo.FindIds("VATTS")
                    if len(howMany) <= 1 {
                        message := widget.NewLabel("VATTS apps air subscription revoked")
                        popup := widget.NewModalPopUp(
                            container.NewVBox(message),
                            mainWindow.Canvas(),
                        )
                        popup.Show()
                        time.Sleep(time.Second*2)
                        server.Process.Kill()
                        os.Exit(0)
                    } else {
                        serverIDs, _ = robotgo.FindIds("tts-server-py-bin")
                        for i := 0; i < len(serverIDs); i++ {
                            if os.Getpid() != howMany[i] {
                                TTSprocess, err := os.FindProcess(serverIDs[i])
                                if err != nil {
                                    log.Append(fmt.Sprintf("%e", err), true)
                                }
                                TTSprocess.Kill()
                            }
                        }
                        howMany, _ = robotgo.FindIds("VATTS")
                        for i := 0; i < len(howMany); i++ {
                            if os.Getpid() != howMany[i] {
                                GUIprocess, err := os.FindProcess(howMany[i])
                                if err != nil {
                                    log.Append(fmt.Sprintf("%e", err), true)
                                }
                                GUIprocess.Kill()
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
            mainWindow.SetContent(container.NewVBox(content))
            for {time.Sleep(time.Second)}
        }
        log.Append("Found no other running VATTS apps", verbose)

        log.Append("Setting up main window", verbose)
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
        mainWindow.SetContent(content)
        mainWindow.CenterOnScreen()
        log.Append("Set up of main window done", verbose)
        mainWindow.Canvas().Focus(entry)

        log.Append("Key hook started", verbose)
        isHidden := false
        go func() {
            hook.Register(hook.KeyDown, []string{"ctrl", "shift", "m"}, func(e hook.Event) {
                if isHidden {
                    mainWindow.Show()
                    mainWindow.RequestFocus()
                    isHidden = false
                    log.Append("Window shown", verbose)
                } else {
                    mainWindow.Hide()
                    isHidden = true
                    log.Append("Window hidden", verbose)
                }
            })

            s := hook.Start()
            <-hook.Process(s)
        }()
    }()

    log.Append("Window show and run called", verbose)
    log.Append("Press Ctrl+Shift+M to toggle window minimize state", true)
    mainWindow.ShowAndRun()
    defer conn.Close()
    defer server.Process.Kill()
}

func sendText(text string) {
    log.Append(fmt.Sprintf("Sending text: %s", text), verbose)
		// Send text to server
        _, err := fmt.Fprintf(conn, "%s\n", text)
		if err != nil {
			log.Append(fmt.Sprintf("Error sending text: %e", err), true)
			return
		}

		// Wait for acknowledgment
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Append(fmt.Sprintf("Error reading response: %e", err), true)
			return
		}
        log.Append(fmt.Sprintf("Server response: %s", response), verbose)
}
