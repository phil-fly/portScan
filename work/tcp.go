package work

import (
	"net"
	"time"
)

func IsOpenTCP(IpAddr,Port string) bool {
	conn, err := net.DialTimeout("tcp", IpAddr+":"+Port, time.Second*1)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}