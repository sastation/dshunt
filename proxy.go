package main

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
)

/* 代理相关 */
type URLConf struct {
	b       [1024]byte // 用来存放客户端数据的缓冲区
	n       int        // 接受字节数
	method  string     // http请求的模式
	hostURL *url.URL   // type url.URL
	address string     // host:port string
	domain  string     // domain
}

// 进行直连
func direct(client net.Conn, urlConf URLConf) {
	//获得了请求的host和port，向服务端发起tcp连接
	server, err := net.Dial("tcp", urlConf.address)
	if err != nil {
		//fmt.Println(err)
		return
	}
	defer server.Close()

	//如果使用https协议，需先向客户端表示连接建立完毕
	if urlConf.method == "CONNECT" {
		fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else { //如果使用http协议，需将从客户端得到的http请求转发给服务端
		server.Write(urlConf.b[:urlConf.n])
	}

	//将客户端的请求转发至服务端，将服务端的响应转发给客户端。io.Copy为阻塞函数，文件描述符不关闭就不停止
	go io.Copy(server, client)
	io.Copy(client, server)
}

// 转发至后端代理服务器
func proxy(client net.Conn, urlConf URLConf, proxy string) {
	server, err := net.Dial("tcp", proxy) //连接代理服务器
	if err != nil {
		fmt.Println(err)
		return
	}
	defer server.Close()

	// 访问端与代理服务器直接交换数据
	server.Write(urlConf.b[:urlConf.n])
	go io.Copy(server, client)
	io.Copy(client, server)
}

/*域名判断相关*/

// 获得域名IP
func domainIP(domain string) (ip string) {
	addrs, err := net.LookupIP(domain)
	if err != nil {
		return ""
	}
	return addrs[0].String()
}

// 判断IP的地理位置: Private, CN, Oversea
func IPLocation(ip string) (loc string) {
	if ip == "" {
		loc = "None"
	} else {
		cidr := net.ParseIP(ip)
		if isPrivateIP(cidr) {
			loc = "Private"
		} else if isCNIP(cidr) {
			loc = "CN"
		} else {
			loc = "Oversea"
		}
	}

	return loc
}

// 判断IP是否为私有地址
func isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	for _, block := range conf.privateIPs {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// 判断IP是否为CN IP
func isCNIP(ip net.IP) bool {
	for _, block := range conf.cnIPs {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// 判断domain是否在conf.Domains列表中，若在则返回其类型，若不在则返回空
func DomainType(domain string) string {
	Type := ""
	for i := 0; i < len(conf.Domains); i++ {
		item := conf.Domains[i].Domain
		b1 := strings.HasSuffix(domain, item)
		b2 := strings.HasSuffix(domain, "."+item)
		if b1 || b2 {
			Type = conf.Domains[i].Type
			break
		}

	}
	return Type
}
