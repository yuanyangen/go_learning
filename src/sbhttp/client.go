package sbhttp

import (
	"fmt"
	"io/ioutil"
	"net"
)

func Get(url string) (response string, err error) {
	//addr, err := net.ResolveTCPAddr("tcp4", "test.webid.360.cn:80")
	addr, err := net.ResolveTCPAddr("tcp4", "dev.yyg.youqian.360.cn:80")
	if err != nil {
		fmt.Println("look up ip failed")
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	//_, err = conn.Write([]byte("GET /test.html HTTP/1.1\r\nUser-Agent: Mozilla/5.0\r\nHost: test.webid.360.cn\r\nAccept: */*\r\n\r\n"))
	//request := strings.join(reqParam)
	_, err = conn.Write([]byte("GET /user/signQuery HTTP/1.1\r\nHost: dev.yyg.youqian.360.cn\r\nConnection: keep-alive\r\nCache-Control: max-age=0\r\nAccept: text/html\r\nUser-Agent: Mozilla/5.0\r\nAccept-Encoding: none\r\nAccept-Language: zh-CN;q=1\r\nCookie: Q=u%3Dlhnalnatra%26n%3D%26le%3DrKIuoayuozqyovH0ZQRlAv5wo20%3D%26m%3D%26qid%3D849734198%26im%3D1_t00df551a583a87f4e9%26src%3Dpcw_i360%26t%3D1; T=s%3D9f23781da9e4ef0b33960f8f56d14be0%26t%3D1448114324%26lm%3D%26lf%3D4%26sk%3Df9b79cfed72bf27a5fbe96332073f1e8%26mt%3D1448114324%26rc%3D%26v%3D2.0%26a%3D2\r\n\r\n"))
	result, err := ioutil.ReadAll(conn)
	response = string(result)
	if err != nil {
		fmt.Println("111111111111111111")

	}

	return
}

/*func CUrl(url string) (response string, err error) {*/

/*}*/
