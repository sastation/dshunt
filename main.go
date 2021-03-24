package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
)

func main() {
	parseConfig()
	fmt.Println(conf.cFile, conf.listen, conf.proxyHost, conf.dproxyHost, conf.IPFile, conf.PrintVer)

	// tcp连接，监听 conf.listen
	fmt.Printf("Listen %s\n", conf.listen)

	l, err := net.Listen("tcp", conf.listen)
	if err != nil {
		log.Panic(err)
	}

	// 无限循环，每当遇到连接时，调用handle
	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}

		go handle(client)
	}
}

// 处理客户端请求
func handle(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	var err error
	var urlConf URLConf
	//从客户端获取数据
	urlConf.n, err = client.Read(urlConf.b[:])
	if err != nil {
		log.Printf("client.Read error: %s\n", err)
		return
	}

	var URL string
	// 从客户端数据读入method，url
	fmt.Sscanf(string(urlConf.b[:bytes.IndexByte(urlConf.b[:], '\n')]), "%s%s", &urlConf.method, &URL)

	// 若方法是CONNECT，则为https协议
	if urlConf.method == "CONNECT" {
		URL = "https://" + URL
	}

	urlConf.hostURL, err = url.Parse(URL)
	if err != nil {
		log.Printf("url.Parse error: %s", err)
		return
	}

	urlConf.address = urlConf.hostURL.Host
	// 若host不带端口，则默认为80
	if !strings.Contains(urlConf.hostURL.Host, ":") {
		urlConf.address = urlConf.hostURL.Host + ":80"
	}

	urlConf.domain = strings.Split(urlConf.address, ":")[0]
	// 解析 domain获得IP
	ip := domainIP(urlConf.domain)

	// 判断 ip 所在位置: private | cn | oversea
	loc := IPLocation(ip)

	info := fmt.Sprintf("%v ==> %s %s [%s]", client.RemoteAddr(), urlConf.method, urlConf.address, loc)

	if loc == "CN" || loc == "Private" {
		log.Printf("%s [direct]\n", info)
		direct(client, urlConf)
	} else {
		log.Printf("%s [proxy]\n", info)
		proxy(client, urlConf)
	}

}
