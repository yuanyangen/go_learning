package server

//实现功能：
//1： 能够对于每一个不同的请求， 能够单独开协程完成这个请求
//2： 解析出请求中的所有的参数， 就是说支持常见的http报文
//3： 能够直接将返回的结果拼装成http响应报文， 发送给客户端
//4： 长连接

import (
	"fmt"
	"net"
	"time"
)

var readChan = make(map[*net.Conn]chan string, 60000)
var writeChan = make(map[*net.Conn]chan string, 60000)

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
			time.Sleep(5 * time.Millisecond)
			continue
		}

		//这里实现的功能包括：
		//保留将连接保存起来，对于连接,
		//这里先不考虑长连接的问题，有一个连接过来的时候， 就直接为其复用同一个连接，这样说起里
		//实现长连接完全是客户端的事情了？因为连接时由客户端发起的， 如果客户端执意要建立新的连接
		//服务端也没办法吧
		//需要有两个channel， 主协程读数据， 主协程些数据
		go slaveProcess(conn)

	}
}

//下一步：解析http协议, 将请求解析，并生成到结构体中，然后就OK了
func slaveProcess(conn net.Conn) {
	fmt.Println("success create a conn and try to write test to it ")
	readChan[&conn] = make(chan string)
	writeChan[&conn] = make(chan string)
	go writeLoop(&conn)
	go readLoop(&conn)

	for {
		msg := <-readChan[&conn]
		time.Sleep(3 * time.Second)
		writeChan[&conn] <- msg
	}

}

func readLoop(conn *net.Conn) {
	co := *conn
	msg := make([]byte, 100)
	for {
		_, _ = co.Read(msg)
		readChan[conn] <- string(msg)
		continue
	}
}

func writeLoop(conn *net.Conn) {
	co := *conn
	for {
		msg := <-writeChan[conn]
		co.Write([]byte(msg))
		continue
	}
}
