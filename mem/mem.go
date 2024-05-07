package mem

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
)

// getVmRSS 从 /proc/[pid]/status 获取 VmRSS 值
func getVmRSS(pid int) (uint64, error) {
    file, err := os.Open(fmt.Sprintf("/proc/%d/status", pid))
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
                return 0, fmt.Errorf("unexpected VmRSS line format: %s", line)
            }
            rssValue, err := strconv.ParseUint(parts[1], 10, 64)
            if err != nil {
                return 0, err
            }
            return rssValue, nil
        }
    }
    return 0, fmt.Errorf("VmRSS value not found for pid %d", pid)
}

// getTotalMemory 从 /proc/meminfo 获取总内存
func getTotalMemory() (uint64, error) {
    file, err := os.Open("/proc/meminfo")
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
                return 0, fmt.Errorf("unexpected MemTotal line format: %s", line)
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

// humanizeBytes 将字节转换为易读的格式
func humanizeBytes(bytes uint64) string {
    for _, unit := range []string{"B", "KB", "MB", "GB", "TB"} {
        if float64(bytes) < 1000 {
            return fmt.Sprintf("%.2f %s", float64(bytes), unit)
        }
        bytes /= 1000
    }
    return "n/a"
}

