package server

import (
	"fmt"
	"net"
	"time"
)

func StartServer() bool {
	fmt.Println("start listen and serve")
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		fmt.Println("resolve tcp addr failed")
		return false
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	for {
		time.Sleep(time.Second)
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		fmt.Println("success create a conn and try to write test to it ")
		conn.Write([]byte("test\n"))
		for {
			msg := make([]byte, 100)
			_, _ = conn.Read(msg)
			fmt.Println(string(msg))
			sms := make([]byte, 100)
			_, _ = fmt.Scan(&sms)
			conn.Write(sms)
		}
		//conn.Close()

	}
}
