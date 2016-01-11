package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"
)

func demo1(url string) (resp string) {
	tmp := strings.Split(url, "://")
	protocal := tmp[0]
	tmp2 := strings.Split(tmp[1], "/")
	host := tmp2[0]
	tmp3 := strings.Split(url, host)
	path := tmp3[1]

	httpHeader := ""
	splitTag := "\r\n"
	httpHeader = httpHeader + "GET " + path + " HTTP/1.1" + splitTag
	httpHeader += "Host: " + host + splitTag + splitTag
	ipAddr, _ := net.ResolveTCPAddr("tcp4", host+":80")
	conn, _ := net.DialTCP("tcp", nil, ipAddr)
	_, _ = conn.Write([]byte(httpHeader))
	info, _ := ioutil.ReadAll(conn)

	fmt.Println(string(info), protocal)
	return
}

func main() {
	demo1("http://youqian.360.cn/user/signQuery")
}
