package curl

//实现的功能是：
//1： post done
//2： cookie done
//3： 自定义header done
//4： https
//5:  上传文件
//6:  允许执行的时间 done
//7:  允许执行多个请求
//8:  允许设置编码的类型
//9:  允许设置referer done
//10: 允许设置的url done

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

//实现curl 功能需要有3个部分：
// 1: Request, 这个是一个结构体，其中包含很多的属性
// 2: connection， 需要根据Request中的host属性及端口号建立连接
// 3: response， 解析返回的响应头，将常用的属性都写入到其中的结构体
//这里含有一个conn，说明对于不同的请求，使用的是短链接
type Ch struct {
	req  request
	resp response
	conn connection
}

//请求类型，描述一个http请求所需要或者含有的所有的资源
type request struct {
	method       string
	host         string
	port         int
	protocal     string
	path         string
	postField    string
	cookieJar    string
	customHeader string
	referer      string
	timeout      time.Duration
}

type response struct {
	//status int32
	header string
	body   string
	/* setCookie        string*/
	//contentType      string
	//date             string
	//transferEncoding string
	//connection       string
	//server           string
	/*cacheControl     string*/
}

type connection struct {
	host string
	port int32
}

//得到一个请求的指针，其指向的对象就是一次http请求的所有的信息
var defaultRequest = request{
	"GET",
	"",
	80,
	"HTTP 1.1",
	"",
	"",
	"",
	"",
	"",
	100,
}

var defaultResp = response{
	"",
	"",
}
var defaultConn = connection{
	"",
	80,
}

//返回的是一个Ch的对象，这个是因为
func Init() Ch {
	return Ch{defaultRequest, defaultResp, defaultConn}
}

//设置url的对象属性
func (this *Ch) SetUrl(url string) {
	//将协议头取出来
	t := strings.Split(url, "://")
	//将协议头转换为大写
	this.req.protocal = strings.ToUpper(t[0])

	//将host及端口取出来，
	tmp := strings.Split(t[1], "/")
	//判断url中是否有写明端口
	if strings.Contains(tmp[0], ":") {
		tmp2 := strings.Split(tmp[0], ":")
		this.req.host = tmp2[0]
		port, err := strconv.Atoi(tmp2[1])
		if err != nil {
		}
		this.req.port = port
	} else {
		this.req.host = tmp[0]
		this.req.port = 80
	}

	//取出路径
	this.req.path = strings.TrimPrefix(t[1], tmp[0])
}

//设置允许post方式进行参数的传递
func (this *Ch) SetPost() {
	this.req.method = "POST"
}

//设置post的参数， 由于是设置参数， 因此这里直接将得到的参数直接
//拼接到post的的域
func (this *Ch) SetPostField(key string, value string) {
	this.req.postField += key + "=" + url.QueryEscape(value) + "&"
}

//将cookie写入http头中
func (this *Ch) SetCookieJar(cookieJar string) {
	this.req.cookieJar = cookieJar
}

//设置referer
func (this *Ch) SetReferer(url string) {
	this.req.referer = url
}

//设置自定义的http头部信息
func (this *Ch) SetHeader(header string) {
	this.req.customHeader = header
}

func (this *Ch) SetTimeout(timeout time.Duration) {
	this.req.timeout = timeout
}

//运行一个curl的对象，并将返回信息返回
func (this *Ch) Execute() {
	//设置一个channel作为超时使用
	channel := make(chan bool, 1)
	go func() {
		conn := this.getConn()
		httpRequest := this.getHttpHeader()
		if this.req.method == "POST" {
			httpRequest += this.req.postField
		}
		_, err := conn.Write([]byte(httpRequest))
		if err != nil {
		}
		tmpInfo, _ := ioutil.ReadAll(conn)
		//这里应该将输出的字符串整理成response对象
		this.resp.processResponse(string(tmpInfo))
		channel <- true
	}()
	//如果超时, 返回的响应头中的header为timeout
	select {
	case <-channel:
		{
			return
		}
	case <-time.After(time.Second * this.req.timeout):
		{
			this.resp.header = "timeout"
		}
	}
	return
}

func (this *Ch) GetBody() (body string) {
	body = this.resp.body
	return
}
func (this *Ch) GetHeader() string {
	return this.resp.header
}

//解析http的响应，这里只是将http的头和body分开
func (this *response) processResponse(res string) {
	tmpResponse := strings.Split(res, "\r\n\r\n")
	//解析http response头
	this.header = tmpResponse[0]

	//如果返回的transfer-encoding是chunked，对响应的body进行解析
	if strings.Contains(this.header, "Transfer-Encoding: chunked") {
		//这里临时只取一个
		this.parserChunkedBody(tmpResponse[1])
	} else {
		this.body = tmpResponse[1]
	}
}

//这里的临时解决办法是只考虑只有一个分片的时候， 如果拥有多个分片， 目前的做法会有问题
func (this *response) parserChunkedBody(body string) {
	tmpBody := strings.Split(body, "\r\n")
	this.body = tmpBody[1]
}

//根据reqest对象的所有信息，拼装得到http请求头信息
func (this *Ch) getHttpHeader() (header string) {
	header = ""
	splitTag := "\r\n"

	//拼接http的方法， 这个一定要在最开始
	header = header + this.req.method + " " + this.req.path + " " + this.req.protocal + "/1.1" + splitTag

	//拼接host
	header = header + "host: " + this.req.host + splitTag

	//设置post的头和正文的长度
	if this.req.method == "POST" {
		this.req.postField = strings.TrimRight(this.req.postField, "&")
		header = header + "Content-Length: " + strconv.Itoa(strings.Count(this.req.postField, "")-1) + splitTag
		header = header + "Content-Type: application/x-www-form-urlencoded" + splitTag
	}

	//如果有referer就取出
	if this.req.referer != "" {
		header += "Referer: " + this.req.referer + splitTag
	}

	//拼接cookie头
	if this.req.cookieJar != "" {
		header += "Cookie: " + this.req.cookieJar + splitTag
	}

	//如果有自定义的头，就去出
	if this.req.customHeader != "" {
		header += this.req.customHeader + splitTag
	}

	header += splitTag
	return
}

//根据reqest的中的信息，得到对应的tcp连接
func (ch *Ch) getConn() (conn *net.TCPConn) {
	host := ch.req.host + ":" + strconv.Itoa(ch.req.port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		fmt.Println("failed")
	}

	conn, err = net.DialTCP("tcp", nil, tcpAddr)
	return
}
