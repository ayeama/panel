package stat

import (
	"os"
	"strconv"
	"strings"
)

func Cpu() (uint64, uint64) {
	c, err := os.ReadFile("/proc/stat")
	if err != nil {
		panic(err)
	}
	cpuLine := strings.SplitN(string(c), "\n", 2)[0]
	cpuValues := strings.FieldsFunc(cpuLine, func(c rune) bool { return c == ' ' })

	user, _ := strconv.ParseUint(cpuValues[1], 10, 64)
	nice, _ := strconv.ParseUint(cpuValues[2], 10, 64)
	system, _ := strconv.ParseUint(cpuValues[3], 10, 64)
	idle, _ := strconv.ParseUint(cpuValues[4], 10, 64)
	iowait, _ := strconv.ParseUint(cpuValues[5], 10, 64)
	irq, _ := strconv.ParseUint(cpuValues[6], 10, 64)
	softirq, _ := strconv.ParseUint(cpuValues[7], 10, 64)
	steal, _ := strconv.ParseUint(cpuValues[8], 10, 64)

	cpuTotal := user + nice + system + idle + iowait + irq + softirq + steal
	cpuIdle := idle + iowait

	return cpuTotal, cpuIdle
}

func CpuPercent(total *uint64, idle *uint64) float64 {
	totalNow, idleNow := Cpu()

	totalDelta := float64(totalNow - *total)
	idleDelta := float64(idleNow - *idle)

	*total = totalNow
	*idle = idleNow

	return (totalDelta - idleDelta) / totalDelta
}
