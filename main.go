package main

import (
	"ops-monitor/mem"
	"time"
	"fmt"
	"os"
)

func main() {
    pid := os.Getpid() // 当前进程的 PID
    memUsageTicker := time.NewTicker(2 * time.Second) // 每2秒更新一次

    defer func() {
        // 程序结束时停止 ticker
        <-memUsageTicker.C
    }()

    fmt.Println("Memory usage of PID", pid, "will update every 2 seconds.")

    for {
        select {
        case <-memUsageTicker.C:
            rss, err := mem.GetVmRSS(pid)
            if err != nil {
                fmt.Println("Error reading VmRSS:", err)
                continue
            }

            totalMemory, err := mem.GetTotalMemory()
            if err != nil {
                fmt.Println("Error reading total memory:", err)
                continue
            }

            memUsagePercent := float64(rss) / float64(totalMemory) * 100
            fmt.Printf("\rMemory Usage: %s (%.2f%%)", mem.HumanizeBytes(rss*1024), memUsagePercent)
        }
    }
}