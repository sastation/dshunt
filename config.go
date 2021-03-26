package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// 版本号
const (
	version = "0.0.1 2021/3/21"
)

type DominType struct {
	Domain string
	Type   string
}

type Config struct {
	cFile         string
	listen        string
	proxyHost     string
	dproxyHost    string
	IPFile        string
	PrintVer      bool
	cnIPs         []*net.IPNet
	privateIPs    []*net.IPNet
	directDomains []DominType
	proxyDomains  []DominType
	rejectDomains []DominType
}

var conf Config

// 初始化配置参数
func initConfig() {
	conf.listen = "127.0.0.1:1234"
	conf.proxyHost = "s16.lan:8081"
	conf.dproxyHost = ""
	conf.IPFile = "china_ip_list.txt"
}

// 获取命令行参数
func parseCMDLine() {
	flag.StringVar(&conf.cFile, "conf", "rc.conf", "config file, default rc.conf")
	flag.BoolVar(&conf.PrintVer, "version", false, "print version")
	flag.Parse()

	if conf.PrintVer {
		fmt.Println("ZRouter version: ", version)
	}
}

// 处理<schema://>host:port字符串
func parseAddress(address string) string {

	addrs := strings.Split(address, "://")
	switch len(addrs) {
	case 1:
		_ = true // pass
	case 2:
		address = addrs[1]
	default:
		panic("Wrong format " + address)

	}

	return address
}

// 获取配置文件内参数
func parseConfig() {
	initConfig()
	parseCMDLine()

	fs, err := os.Open(conf.cFile)
	if err != nil {
		panic(err)
	}
	defer fs.Close()

	buf := bufio.NewScanner(fs)
	for buf.Scan() {
		line := strings.TrimSpace(buf.Text())
		if line == "" || line[0] == '#' {
			continue
		}

		v := strings.SplitN(line, "=", 2)
		if len(v) != 2 {
			continue
		}

		key, value := strings.TrimSpace(v[0]), strings.TrimSpace(v[1])
		switch strings.ToLower(key) {
		case "listen":
			conf.listen = parseAddress(value)
		case "proxy":
			conf.proxyHost = parseAddress(value)
		case "dproxy":
			conf.dproxyHost = parseAddress(value)
		case "ipfile":
			conf.IPFile = value
		default:
			log.Printf("Error: unknow parameter %s\n", key)
		}
	}
	if buf.Err() != nil {
		panic(buf.Err)
	}

	privateIP()
	cnIP()
}

// 读取conf.IPFile文件内的CN_IP数据
func cnIP() {
	fs, err := os.Open(conf.IPFile)
	if err != nil {
		panic(err)
	}
	defer fs.Close()

	buf := bufio.NewScanner(fs)
	for buf.Scan() {
		line := buf.Text()
		_, subnet, _ := net.ParseCIDR(line)
		conf.cnIPs = append(conf.cnIPs, subnet)
	}

	if buf.Err() != nil {
		panic(buf.Err())
	}
}

func privateIP() {
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	} {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(fmt.Errorf("parse error on %q: %v", cidr, err))
		}
		conf.privateIPs = append(conf.privateIPs, block)
	}
}
