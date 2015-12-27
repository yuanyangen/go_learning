package curl

//实现的功能是：
//1： post
//2： cookie
//3： 自定义header
//4： https
//5:  上传文件
//6:  允许执行的时间
//7:  允许执行多个请求
//8:  允许设置编码的类型
//9:  允许设置referer
//10: 允许设置的url
import (
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

//curl设置的option中，值为string的内容, 这个部分的内容在拼接http头的时候，直接
//将值进行拼接即可
type requestStringOption struct {
}

//curl设置的属性中， 值为bool的内容, 该部分的内容需要进过简单的判断
type requestBoolOption struct {
}

type Request struct {
	host         string
	port         int
	protocal     string
	path         string
	get_param    []string
	stringOption requestStringOption
	boolOption   requestBoolOption
}

type connection struct {
	host string
	port int32
}

var req = &Request{}

//实现curl 功能需要有3个部分：
// 1: request, 这个是一个结构体，其中包含很多的属性
// 2: connection， 需要根据request中的host属性及端口号建立连接
// 3: response， 解析返回的响应头，将常用的属性都写入到其中的结构体
func Get(url string) (response string) {
	parserRequest(url)
	conn := getConn(req)
	httpHeader := getHttpRequest()
	_, err := conn.Write(httpHeader)
	if err != nil {

	}
	msg, _ := ioutil.ReadAll(conn)
	//fmt.Println(string(msg))
	response = string(msg)
	return response
}

//获取http请求头信息
func getHttpRequest() (httpHeader []byte) {
	header := ""
	splitTag := "\r\n"
	endTag := "\r\n\r\n"
	header = header + "GET " + req.path + " " + req.protocal + "/1.1" + splitTag
	header = header + "host: " + req.host + splitTag
	/* header = header + "Connection: keep-alive" + splitTag*/
	////设置了referer
	//if req.referer != "" {
	//header = header + "Referer: " + req.referer
	/*}*/

	header += endTag

	fmt.Println(header)
	httpHeader = []byte(header)
	return
}
func getConn(req *Request) (conn *net.TCPConn) {
	host := req.host + ":" + strconv.Itoa(req.port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		fmt.Println("failed")
	}

	conn, err = net.DialTCP("tcp", nil, tcpAddr)
	return
}

func parserRequest(url string) {
	//将协议头取出来
	t := strings.Split(url, "://")
	//将协议头转换为大写
	req.protocal = strings.ToUpper(t[0])

	//将host及端口取出来，
	tmp := strings.Split(t[1], "/")
	//判断url中是否有写明端口
	if strings.Contains(tmp[0], ":") {
		tmp2 := strings.Split(tmp[0], ":")
		req.host = tmp2[0]
		port, err := strconv.Atoi(tmp2[1])
		if err != nil {
		}
		req.port = port
	} else {
		req.host = tmp[0]
		req.port = 80
	}

	//取出路径
	req.path = strings.TrimPrefix(t[1], tmp[0])
}
