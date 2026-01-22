package blockchain

import (
    "fmt"
    "os"
    "time"
)

const dataDir = "data/"

func Now() int64 {
    return time.Now().Unix()
}

func appendLog(message string) {
    logFile := "data/miner.log"
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    full := "[" + timestamp + "] " + message

    f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Printf("? Failed to write log: %v", err)
        return
    }
    defer f.Close()
    _, err = f.WriteString(full)
    if err != nil {
        fmt.Printf("? Write error: %v", err)
    }
}