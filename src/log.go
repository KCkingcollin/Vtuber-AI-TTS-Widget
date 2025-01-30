package log

import (
    "errors"
    "fmt"
    "os"
    "time"
)

var logFile string
var fileLimit int

// location, name of the file, and a size in MB
func Init(dir, name string, limit int) error {
    fileLimit = limit*1e6
    p := dir + "/" + name
    info, err := os.Stat(p)
    if err != nil {
        if os.IsNotExist(err) {
            err := os.MkdirAll(dir, os.ModePerm)
            if err != nil {
                return errors.New("Unable to create directory: " + err.Error())
            }

            _, err = os.Create(p)
            if err != nil {
                return errors.New("Unable to create log file: " + err.Error())
            }

            logFile = p
            return nil
        }

        return errors.New("Unable to stat log file: " + err.Error())
    }

    if info.IsDir() {
        return errors.New("File name is directory")
    }

    logFile = p
    return nil
}

func Append(message string, verbose bool) {
    if verbose {fmt.Println(message)}
    file, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        fmt.Printf("Unable to open log file: %s\n%s\n", err, message)
        return
    }
    defer file.Close()

    if _, err = file.WriteString("[" + time.Now().Format(time.ANSIC) + "] " + message + "\n"); err != nil {
        fmt.Printf("Unable to write to log file: %s\n%s\n", err, message)
        return
    }

    fi, err := os.Stat(logFile)
    if err != nil {
        fmt.Println("Unable to stat log file: " + err.Error())
        return
    }

    fileSize := int(fi.Size())
    if fileSize > fileLimit {
        // // for testing
        // fmt.Println("testing: max file size reached")

        b, err := os.ReadFile(logFile)
        if err != nil {
            fmt.Println("Unable to read log file" + err.Error())
            return
        }

        err = file.Truncate(int64(fileLimit))
        if err != nil {
            fmt.Println("Unable to truncate log file" + err.Error())
            return
        }

        new := b[fileSize-fileLimit:]
        _, err = file.Write(new)
        if err != nil {
            fmt.Println("Unable to write to log file" + err.Error())
            return
        }
    }
}
