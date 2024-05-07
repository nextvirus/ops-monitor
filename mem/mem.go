package mem


import (
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "syscall"
    "time"
    "unsafe"
)

var (
    kernelPageSize = uint64(syscall.Getpagesize())
    pageSizes      = map[string]uint64{
        "4K": 4096,
        "2M": 2097152,
        "1G": 1073741824,
    }
)

// 读取并解析 /proc/meminfo 文件，返回总内存（单位：KB）
func getTotalMemory() (uint64, error) {
    const memInfoPath = "/proc/meminfo"
    file, err := os.Open(memInfoPath)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "MemTotal:") {
            parts := strings.Split(line, " ")
            if len(parts) < 2 {
                continue // 格式错误，跳过
            }
            memTotal, err := strconv.ParseUint(parts[1], 10, 64)
            if err != nil {
                return 0, err
            }
            return memTotal, nil
        }
    }
    return 0, fmt.Errorf("MemTotal value not found")
}


// 读取并解析 /proc/[pid]/status 文件，返回 VmRSS 值
func getVmRSS(pid int) (uint64, error) {
    statusFilePath := fmt.Sprintf("/proc/%d/status", pid)
    file, err := os.Open(statusFilePath)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "VmRSS:") {
            parts := strings.Split(line, " ")
            if len(parts) < 2 {
                continue // 格式错误，跳过
            }
            rssValue, err := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 64)
            if err != nil {
                return 0, err
            }
            return rssValue, nil
        }
    }
    return 0, fmt.Errorf("VmRSS value not found")
}

func printMemoryUsage(pid int, rss, totalMemory uint64) {
    fmt.Printf("%-6s %-8s %-10s\n", "PID", "RSS", "MEM USAGE")
    fmt.Printf("%-6d %-8s %-10.2f%%\n", pid, humanizeBytes(rss * pageSizes["4K"]), float64(rss)/float64(totalMemory)*100)
}

func humanizeBytes(bytes uint64, sizes []string) string {
    if bytes == 0 {
        return "0"
    }
    base := kernelPageSize
    for _, size := range sizes {
        newSize := base * pageSizes[size]
        if bytes >= newSize {
            return fmt.Sprintf("%d%s", bytes/newSize, size)
        }
        base = newSize
    }
    return fmt.Sprintf("%dB", bytes)
}