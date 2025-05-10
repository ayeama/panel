package domain

import (
	"fmt"
	"math/rand/v2"
	"net"
)

func freeHostPort() (uint16, error) {
	retries := 3
	var port uint16

	min := 10000
	max := 49152

	for i := 0; i < retries; i++ {
		port := uint16(rand.IntN(max-min) + min)
		listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
		if err != nil {
			continue
		}
		defer listen.Close()
		return port, err
	}

	return port, fmt.Errorf("panel: could not get free port in %d retries", retries)
}
