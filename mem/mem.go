package mem

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
)

// GetVmRSS 获取指定进程的物理内存使用量（RSS），单位为KB
func GetVmRSS(pid int) (uint64, error) {
    file, err := os.Open(fmt.Sprintf("/proc/%d/status", pid))
    if err != nil {
        return 0, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "VmRSS:") {
            parts := strings.Fields(line)
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

// GetVmSize 获取指定进程的虚拟内存大小（VmSize），单位为KB
func GetVmSize(pid int) (uint64, error) {
    file, err := os.Open(fmt.Sprintf("/proc/%d/status", pid))
    if err != nil {
        return 0, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "VmSize:") {
            parts := strings.Fields(line)
            if len(parts) < 2 {
                return 0, fmt.Errorf("unexpected VmSize line format: %s", line)
            }
            sizeValue, err := strconv.ParseUint(parts[1], 10, 64)
            if err != nil {
                return 0, err
            }
            return sizeValue, nil
        }
    }
    return 0, fmt.Errorf("VmSize value not found for pid %d", pid)
}

// GetVmSwap 获取指定进程的交换内存使用量（VmSwap），单位为KB
func GetVmSwap(pid int) (uint64, error) {
    file, err := os.Open(fmt.Sprintf("/proc/%d/status", pid))
    if err != nil {
        return 0, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "VmSwap:") {
            parts := strings.Fields(line)
            if len(parts) < 2 {
                return 0, fmt.Errorf("unexpected VmSwap line format: %s", line)
            }
            swapValue, err := strconv.ParseUint(parts[1], 10, 64)
            if err != nil {
                return 0, err
            }
            return swapValue, nil
        }
    }
    return 0, fmt.Errorf("VmSwap value not found for pid %d", pid)
}

// GetTotalMemory 获取系统的总内存，单位为Bytes
func GetTotalMemory() (uint64, error) {
    file, err := os.Open("/proc/meminfo")
    if err != nil {
        return 0, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "MemTotal:") {
            parts := strings.Fields(line)
            if len(parts) < 2 {
                return 0, fmt.Errorf("unexpected MemTotal line format in /proc/meminfo: %s", line)
            }
            memTotalStr := strings.TrimSpace(parts[1])
            memTotal, err := strconv.ParseUint(memTotalStr, 10, 64)
            if err != nil {
                return 0, err
            }
            return memTotal * 1024, nil // 因为 MemTotal 的单位是 KB，所以乘以 1024 转换为 Bytes
        }
    }
    return 0, fmt.Errorf("MemTotal value not found in /proc/meminfo")
}

// HumanizeBytes 将字节转换为易读的格式
func HumanizeBytes(bytes uint64) string {
    units := []string{"B", "KB", "MB", "GB", "TB"}
    base := uint64(1024)
    unitIndex := 0
    for bytes >= base {
        bytes /= base
        unitIndex++
    }
    return fmt.Sprintf("%.2f %s", float64(bytes), units[unitIndex])
}