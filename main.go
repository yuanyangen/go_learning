package main

import (
	"curl"
	"fmt"
)

func main() {
	str := curl.Get("http://youqian.360.cn/user/signQuery")
	fmt.Println(str)
	fmt.Println("for test")
}
