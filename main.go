package main

import (
    "fmt"
    "ops-monitor/mem"
    "time"
)

func main() {
    pid := 3889 // 获取当前进程的 PID
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

            vmSwap, err := mem.GetVmSwap(pid)
            if err != nil {
                fmt.Println("Error reading VmSwap:", err)
                continue
            }

            totalMemory, err := mem.GetTotalMemory()
            if err != nil {
                fmt.Println("Error reading total memory:", err)
                continue
            }

            // 计算总内存使用率，包括物理内存和交换内存
            totalMemUsage := uint64(rss) + uint64(vmSwap)
            memUsagePercent := float64(totalMemUsage) / float64(totalMemory) * 100

            fmt.Printf("\rMemory Usage of PID %d: %s (%.2f%%)", pid, mem.HumanizeBytes(uint64(totalMemUsage)*1024), memUsagePercent)
        }
    }
}