# ZRouter：选择性分流器

## 需求与目的
  在通常环境中需要分流到直连与代理，而在在某些环境中又需要分流到代理1、代理2。比如在某个办公环境CN网站分流到S01.lan，境外网站分流到S02.lan。目前世面上的绝大多数分流器默认情况下都只能分流直连与代理两种路径，若是要支持多重路径分流则需要编写复杂的规则而且调试困难。在此情况下结合本人魔改meow的经验，编写了ZRouter，用最简单的方式实现了3选2选择性分流器。
  
## 与魔改meow的区别
  1. CN_IP列表由内置于程序变为外置，便于更新
  2. 目前只支持 http 协议，包括监听与代理
  3. 程序大幅度简化（练手目的:)）
  4. 程序能自己下载 china_ip_list.txt 
  
## 使用方法
```bash
  Usage of ./zrouter:
  -conf string
          config file (default "rc.conf")
  -down
        download cn_ip.txt file
  -url string
        url for cn_ip.txt file (default "https://github.com/17mon/china_ip_list")
  -version
        print version
```
## 配置文件
  - 可以通过 "-conf" 选项在运行时指定，若未指定则默认为 rc.conf
  - 可用的配置为：
```bash
# 监听地址，默认为127.0.0.1:1234，设为0.0.0.0可以监听所有，只支持http协议
listen = http://0.0.0.0:8080

# 后端代理转发服务器，无默认值为必选项，只支持http协议
proxy = http://10.0.0.1:3128

# 默认为空，只支持http协议，若空_IP直连，非空CN_IP走dproxy指定代理
#dproxy = http://10.0.0.2:3128


# CN段IP地址列表文件名，默认为 cn_ip.txt
# 可用来源：https://github.com/17mon/china_ip_list
#cnip = cn_ip.txt

# 另有3个默认文件用于域名管理，为不可更改参数
# reject.txt - 拒绝域名列表
# direct.txt - 直连域名列表
# proxy.txt  - 代理域名列表
```
