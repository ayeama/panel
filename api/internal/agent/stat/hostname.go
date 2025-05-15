package stat

import (
	"os"
	"strings"
)

func Hostname() string {
	c, err := os.ReadFile("/etc/hostname")
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(c))
}
