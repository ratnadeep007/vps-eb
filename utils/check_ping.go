package utils

import (
	"fmt"
	"net"
	"time"
)

// CheckPing -> Check system is up
func CheckPing(addr string) bool {
	gotit := false
	for {
		conn, err := net.Dial("tcp", addr+":22")
		if err == nil {
			// log.Fatal(err.Error())
			fmt.Print("got conncetion")
			gotit = true
			conn.Close()
			break
		}
		time.Sleep(10 * time.Second)
	}
	return gotit
}
