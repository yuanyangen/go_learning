package server

//实现功能：
//1： 能够对于每一个不同的请求， 能够单独开协程完成这个请求
//2： 解析出请求中的所有的参数， 就是说支持常见的http报文
//3： 能够直接将返回的结果拼装成http响应报文， 发送给客户端
//4： 长连接

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type httpReq struct {
	method     string
	protocal   string
	path       string
	cookie     map[string]string
	referer    string
	mimeType   string
	connetcion string
	body       string
}
type httpResp struct {
	statusCode    string
	statusString  string
	setCookie     map[string]string
	body          string
	contentType   string
	cache_control string
}

type app struct {
	req       httpReq
	resp      httpResp
	writeChan chan string
	readChan  chan string
	conn      *net.Conn
}

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

		go slaveProcess(conn)
	}
}

//这里实现的功能包括：
//保留将连接保存起来，对于连接,
//这里先不考虑长连接的问题，有一个连接过来的时候， 就直接为其复用同一个连接，这样说起里
//实现长连接完全是客户端的事情了？因为连接时由客户端发起的， 如果客户端执意要建立新的连接
//服务端也没办法吧
//需要有两个channel， 主协程读数据， 主协程些数据

//下一步：解析http协议, 将请求解析，并生成到结构体中，然后就OK了
//由于http协议的特性，所以在一次http请求过程中， 只需要从IO读入一次数据， 然后将
//处理值写入到IO即
//子协程在
func slaveProcess(conn net.Conn) {
	fmt.Println("success create a conn and try to write test to it ")
	app := &app{}
	app.conn = &conn

	app.readChan = make(chan string)
	app.writeChan = make(chan string)
	go app.writeLoop()
	go app.readLoop()

	msg := <-app.readChan
	resp := app.handler(msg)
	app.writeChan <- resp
	conn.Close()
}

func (this *app) handler(msg string) (resp string) {
	tmp1 := strings.Split(msg, "\r\n")
	tmp2 := strings.Split(tmp1[0], " ")
	this.req.method = tmp2[0]
	this.req.path = tmp2[1]
	this.req.protocal = tmp2[2]
	fmt.Println(this)

	this.resp.statusCode = "200"
	this.resp.statusString = "OK"
	this.resp.contentType = "text/html"
	this.resp.body = "hello world!!!"
	resp = this.getResp()
	return
}

func (this *app) getResp() (response string) {
	splitTag := "\r\n"
	response = this.req.protocal + " " + this.resp.statusCode + " " + this.resp.statusString + splitTag
	//	response += "Server: nginx/1.2.9" + splitTag
	//	response += "Date: Wed, 20 Jan 2016 16:01:29 GMT" + splitTag
	//response += "Content-Type: " + this.resp.contentType + splitTag
	response += "Content-Length: " + strconv.Itoa(len(this.resp.body)) + splitTag
	response += "Connection: close" + splitTag
	//	response += "Accept-Ranges: bytes" + splitTag
	response += splitTag
	response += this.resp.body
	//response += this.resp.body + splitTag
	fmt.Println(response)

	return
}

func (this *app) readLoop() {
	conn := *this.conn
	msg := make([]byte, 100)
	for {
		_, _ = conn.Read(msg)
		this.readChan <- string(msg)
		continue
	}
}

func (this *app) writeLoop() {
	co := *this.conn
	for {
		msg := <-this.writeChan
		co.Write([]byte(msg))
		continue
	}
}
