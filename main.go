package main

import mem

func main() {
    pid := os.Getpid() // 获取当前进程的 PID
    ticker := time.NewTicker(2 * time.Second)

    defer func() {
        // 停止 ticker
        ticker.Stop()
    }()

    for range ticker.C {
        rss, err := getVmRSS(pid)
        if err != nil {
            fmt.Println("Error reading VmRSS:", err)
            continue
        }

        totalMemory, err := getTotalMemory()
        if err != nil {
            fmt.Println("Error reading total memory:", err)
            continue
        }

        printMemoryUsage(pid, rss, totalMemory)

        // 清屏，模仿 top 命令的效果
        cmd := exec.Command("clear")
        cmd.Stdout = os.Stdout
        err = cmd.Run()
        if err != nil {
            fmt.Println("Error clearing screen:", err)
        }
    }
}