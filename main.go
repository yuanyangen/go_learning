package main

import (
	"curl"
	"fmt"
	"webServer"
)

func main() {
	/* resp, _ := learning.Get("http://dev.yyg.youqian.360.cn/test.php")*/
	//msg, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(msg))
	server.StartServer()

}

var chann = make(chan int, 2)

//func main() {
//go b()
////	go a()
//i := 0
//WAIT:
//for {
//select {
//case <-chann:
//{
//i++
//if i > 0 {
//break WAIT
//}
//}
//case <-time.After(time.Second * 2):
//{
//fmt.Println("2s time out")
//break WAIT
//}
//}
//}
//}

func b() {
	targetUrl := "http://youqian.360.cn/user/signQuery"
	//targetUrl = "http://dev.yyg.youqian.360.cn/test.php"
	ch := curl.Init()
	ch.SetUrl(targetUrl)
	ch.SetPost()
	ch.SetPostField("key", "value")
	ch.SetCookieJar("Q=u%3Dlhnalnatra%26n%3D%26le%3DrKIuoayuozqyovH0ZQRlAv5wo20%3D%26m%3D%26qid%3D849734198%26im%3D1_t00df551a583a87f4e9%26src%3Dpcw_union_youqian%26t%3D1; T=s%3D9482f4bd7ebea3d5bd3312ce8bb441e4%26t%3D1450881248%26lm%3D%26lf%3D4%26sk%3D2ea00a41e96dfdaf71212f3c55a734b2%26mt%3D1450881248%26rc%3D%26v%3D2.0%26a%3D1")
	ch.SetHeader("der: hehehe")
	ch.SetReferer("http://www.baidu.com")
	ch.SetTimeout(1)
	ch.Execute()
	msg := ch.GetBody()
	msg = ch.GetHeader()
	fmt.Println(msg)
	chann <- 1

}

func a() {
	ch := curl.Init()
	targetUrl := "http://youqian.360.cn/psp_jump.html"
	ch.SetUrl(targetUrl)
	ch.Execute()
	msg := ch.GetBody()
	fmt.Println(msg)
	chann <- 2
}
