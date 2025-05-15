package stat

import (
	"os"
	"strconv"
	"strings"
)

func Uptime() float64 {
	c, err := os.ReadFile("/proc/uptime")
	if err != nil {
		panic(err)
	}
	uptime, err := strconv.ParseFloat(strings.Split(string(c), " ")[0], 64)
	if err != nil {
		panic(err)
	}
	return uptime
}
