package stat

import (
	"os"
	"strconv"
	"strings"
)

func MemoryPercent() float64 {
	c, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		panic(c)
	}
	memoryLines := strings.Split(string(c), "\n")

	memoryTotalLine := memoryLines[0]
	memoryAvailableLine := memoryLines[2]

	memoryTotal := strings.FieldsFunc(memoryTotalLine, func(c rune) bool { return c == ' ' })
	memoryAvailable := strings.FieldsFunc(memoryAvailableLine, func(c rune) bool { return c == ' ' })

	total, _ := strconv.ParseUint(memoryTotal[1], 10, 64)
	available, _ := strconv.ParseUint(memoryAvailable[1], 10, 64)

	return float64(total-available) / float64(total)
}
