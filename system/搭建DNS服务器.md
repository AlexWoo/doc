# 搭建 DNS 服务器

## 基本概念

### Bind 基本功能

一般 DNS 服务器使用 Bind 进行搭建，Bind 提供基本功能包括

- 域名解析服务
- 权威域名服务
- DNS 工具

### DNS 协议

域名：www.baidu.com -> www.baidu.com.

- 最后一个.: 根域
- .com: 一级域
- baidu.com: 二级域

解析过程

1. 名字解析服务器去根域服务器查询 www.baidu.com，根域返回 .com 所在的名字服务器
2. 名字解析服务器去 .com 所在的名字服务器查询 www.baidu.com，.com 所在的名字服务器返回 .baidu.com 所在的名字服务器
3. 名字解析服务器去 .baidu.com 所在的名字服务器查询 www.baidu.com，.baidu.com 返回 www.baidu.com 的 IP 地址。.baidu.com 即为 www.baidu.com 的权威域名服务器

DNS 记录分类

- A 记录：域名对应 IP 地址的记录为 A 记录
- CNAME 记录：将一个域名映射到另一个域名上，可以理解为别名，用于多个域名解析到同一 IP
- NS 记录：名字服务器不能提供权威解析时，返回另一个名字服务器的地址
- PTR 记录：反向解析，通过 IP 找域名，一般邮件服务会需要反向解析
- MX 记录：针对邮件服务器的域名解析

### 常用客户端工具

- nslookup

	普通查询
		
		nslookup www.baidu.com

- dig

	普通查询
		
		dig www.baidu.com

	指定 DNS 服务器查询
	
		dig @127.0.0.1 www.test.com
		
	反向解析
	
		dig -x mail.gmail.com @127.0.0.1

	查询制定类型记录
	
		dig -t soa www.baidu.com
		dig -t cname www.baidu.com
		dig -t ns baidu.com
		dig -t a www.baidu.com

- host

	普通查询
	
		host www.baidu.com

	查询指定类型记录
	
		host -t SOA www.baidu.com
		host -t CNAME www.baidu.com
		host -t NS baidu.com
		host -t A www.baidu.com

## Bind 安装

安装

	yum install -y bind bind-chroot bind-utils

查看安装内容

	rpm -ql bind | more

启动

- CentOS 6 下
	
		service named start

- CentOS 7 下
	
		systemctl start named

停止

- CentOS 6 下
	
		service named stop

- CentOS 7 下
	
		systemctl stop named


查看状态

- CentOS 6 下
	
		service named status

- CentOS 7 下
	
		systemctl status named

设置开机启动

- CentOS 7 下
	
		systemctl enable named

禁用开机启动

- CentOS 7 下
	
		systemctl disable named

## Bind 配置

### 主配置文件

	/etc/named.conf

- options

	全局配置
	
	- listen-on: 监听 IPv4 端口和地址，不配置默认监听所有地址的 TCP 和 UDP 53 端口
	- listen-on-v6: 监听 IPv6 端口和地址，不配置默认监听所有地址的 TCP 和 UDP 53 端口
	- directory: 域名等配置文件所在目录
	- dump-file: DNS 解析过的内容缓存
	- statistics-file: 静态解析文件，一般不使用
	- memstatistics-file: 内存内容统计
	- allow-query: 权限控制，不配置不开启权限控制
	- recursion: 是否允许递归查询，如果配置权威服务器不需要配置，如果需要递归查询则需要配置

- logging

	服务日志配置，本部分可以不配置

	- channel default_debug: 日志输出级别
	- file: 日志文件位置
	- severity: 日志输出级别

- zone

	DNS 域名解析

	- zone ".": 配置根域服务器, 配置文件默认为 /var/named/named.ca
	- type: DNS 类型，master 表示为主 DNS
	- file: 指定具体的域的配置文件



	zone 文件配置
	
	- @: @ 是 DNS 记录中的保留字，表示当前域名
	- SOA 记录: 每个Zone仅有一个SOA记录。SOA记录包括Zone的名字,一个技术联系人和各种不同的超时值

### 配置权威解析

- named.conf 配置(不包含 options 和 logging 部分)

		zone "test.com." {
			type master;
			file "test.com.zone";
		};

- zone 配置

		$TTL 3H
		@   IN  SOA @   wj19840501.gmail.com. (
		                     0   ; serial
		                     1D  ; refresh
		                     1H  ; retry
		                     1W  ; expire
		                     3H  ; minimum
		                     )
		
		 @       IN      NS          dns
		 dns     IN      A           192.168.0.22
		 www     IN      A           2.2.2.2

	1. 每个 zone 必须包含一个 SOA 类型的域名，对该域名解析从 SOA 开始
	2. 第二列 @ 表示 test.com.，为一个 NS 记录，NS 记录指向 dns.test.com.，表示可以通过 dns.test.com. 查询 test.com.。这里只写一个 dns 是 dns.test.com. 的缩写
	3. 对 dns.test.com. 配置了一个 A 记录，表示 dns.test.com. 的地址是 192.168.0.22
	4. 对 www.test.com. 配置了一个 A 记录，表示 www.test.com. 的地址是 2.2.2.2

### 配置 CNAME 解析

- named.conf 配置(不包含 options 和 logging 部分)

		zone "cname.com." {
			type master;
			file "cname.com.zone";
		};

- zone 配置

		$TTL 3H
		@   IN  SOA @   wj19840501.gmail.com. (
		                     0   ; serial
		                     1D  ; refresh
		                     1H  ; retry
		                     1W  ; expire
		                     3H  ; minimum
		                     )
		
		 @       IN      NS          dns
		 dns     IN      A           192.168.0.22
		 www     IN      CNAME       www.test.com

解析过程与上面相同，区别只是把 www.cname.com. 的 CNAME 记录设置为 www.test.com
