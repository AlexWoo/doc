# keepalived 负载均衡安装配置
---
## keepalived 安装

keepalived 是 LVS 的扩展项目，一般大家的了解是将其用于负载均衡主备模块的切换使用(类似于 Heartbeat)，实际上 keepalived 功能不限于此，从功能上，其划分为三个部分：

- 负载均衡(使用 IPVS 核心，不需要单独安装 ipvsadm)
- 对服务器池的健康检查
- 负载均衡器之间的失败切换

![](http://www.keepalived.org/images/Software%20Design.gif)

如上图，其核心模块为：

- WatchDog：守护模块，负责监控 Checkers 和 VRRP 进程的状况
- Checkers：负责对 RealServer 的健康检查
- VRRP Stack：负责负载均衡器之间的失败切换
- IPVS wrapper：用于发送设定规则到内核的 IPVS 模块
- Netlink Reflector：用于设定 VRRP 的 VIP

官方网站：[http://www.keepalived.org/index.html](http://www.keepalived.org/index.html)

中文文档链接：[http://www.keepalived.org/pdf/sery-lvs-cluster.pdf](http://www.keepalived.org/pdf/sery-lvs-cluster.pdf)

下载 keepalived：

	# wget http://www.keepalived.org/software/keepalived-1.2.19.tar.gz

keepalived 安装依赖于 openssl 开发库

	# rpm -qa |grep openssl
	openssl-libs-1.0.1e-42.el7.9.x86_64
	openssl-devel-1.0.1e-42.el7.9.x86_64
	openssl-1.0.1e-42.el7.9.x86_64

确保有 openssl-devel，如果没有，使用 yum 进行安装：

	# yum install openssl-devel

keepalived 中是包含了 ipvs 库的，实际上不需要安装 ipvsadm，但是需要安装 libnl-devel，popt-devel 和 popt-static库：

	# yum install libnl-devel
	# yum install popt-devel
	# yum install popt-static

编译安装(注意，安装 keepalived 前需要先安装 openssl)：

	# tar xzf keepalived-1.2.19.tar.gz
	# cd keepalived-1.2.19/
	# ./configure --prefix=/usr/local/keepalived
	# make && make install
	# ln -s /usr/local/keepalived/etc/rc.d/init.d/keepalived /etc/rc.d/init.d/keepalived
	# ln -s /usr/local/keepalived/sbin/keepalived /usr/sbin/keepalived
	# ln -s /usr/local/keepalived/etc/keepalived/keepalived.conf /etc/keepalived/keepalived.conf
	# ln -s /usr/local/keepalived/etc/sysconfig/keepalived /etc/sysconfig/keepalived

完成 ./configue 后，配置信息显示

	Keepalived configuration
	------------------------
	Keepalived version       : 1.2.19
	Compiler                 : gcc
	Compiler flags           : -g -O2 -DFALLBACK_LIBNL1
	Extra Lib                : -lssl -lcrypto -lcrypt  -lnl  
	Use IPVS Framework       : Yes
	IPVS sync daemon support : Yes
	IPVS use libnl           : Yes
	fwmark socket support    : Yes
	Use VRRP Framework       : Yes
	Use VRRP VMAC            : Yes
	SNMP support             : No
	SHA1 support             : No
	Use Debug flags          : No

启动 keepalived：

	# systemctl start keepalived.service

查看 keepalived 运行状态：

	# systemctl status keepalived.service

停止 keepalived：

	# systemctl stop keepalived.service

查看 keepalived 帮助文档：

	# man keepalived

查看 keepalived 配置文档：

	# man keepalived.conf

注：这里的系统是 CentOS 7，使用 systemd 替代传统的 SysV init，对于还在使用 SysV init，这里使用以下命令进行启动，查看状态和停止：

	# service keepalived start
	# service keepalived status
	# service keepalived stop

编译中可能会出现以下错误

	/usr/src/linux/include/linux/types.h:113:17: error: conflicting types for ‘int64_t’
	 typedef  __s64  int64_t;
	                 ^
	In file included from /usr/lib/gcc/x86_64-redhat-linux/4.8.2/include/stdint.h:9:0,
	                 from ../include/layer4.h:29,
	                 from layer4.c:24:
	/usr/include/stdint.h:40:19: note: previous declaration of ‘int64_t’ was here
	 typedef long int  int64_t;
	                   ^
	In file included from ../include/../libipvs-2.6/ip_vs.h:12:0,
	                 from ../include/check_data.h:38,
	                 from ../include/check_api.h:27,
	                 from ../include/layer4.h:37,
	                 from layer4.c:24:
	/usr/src/linux/include/linux/types.h:134:23: error: conflicting types for ‘blkcnt_t’
	 typedef unsigned long blkcnt_t;
	                       ^
	In file included from /usr/include/sys/uio.h:24:0,
	                 from /usr/include/sys/socket.h:27,
	                 from /usr/include/netinet/in.h:24,
	                 from /usr/include/netdb.h:27,
	                 from ../include/layer4.h:30,
	                 from layer4.c:24:
	/usr/include/sys/types.h:235:20: note: previous declaration of ‘blkcnt_t’ was here
	 typedef __blkcnt_t blkcnt_t;  /* Type to count number of disk blocks.  */
	                    ^
	make[2]: *** [layer4.o] Error 1
	make[2]: Leaving directory `/root/keepalived-1.2.19/keepalived/core'
	make[1]: *** [all] Error 1
	make[1]: Leaving directory `/root/keepalived-1.2.19/keepalived'
	make: *** [all] Error 2

检查 /usr/src/ 目录下，如果有 linux 软链接，将其删除(注：这块一般是手工安装过 ipvsadm 时，会出现该错误)

	# rm -f /usr/src/linux

注意不是 rm -f /usr/src/linux/，没有最后一个 /


## keepalived 的配置

按照上面方法进行安装，在 /usr/local/keepalived/etc/keepalived/samples 会有 keepalived 的配置样例，配置详细介绍可以使用命令查看：

	# man keepalived.conf

keepalived 配置分为三个部分

- GLOBAL CONFIGURATION
- VRRPD CONFIGURATION
- LVS CONFIGURATION

三个部分的配置方法具体使用 man 命令查看，下面以一个实例进行说明

### 配置规划

	VIP：10.1.63.2
	Virtual Server1：10.1.63.101 (MASTER)
	Virtual Server2：10.1.63.102 (BACKUP)
	Real Server1：10.1.63.103
	Real Server2：10.1.63.104
	策略：最小连接数
	条件：目的端口为 80
	LVS 负载模式：DR

### Virtual Server 配置

**Virtual Server 打开 IP 包转发**

	echo 1 > /proc/sys/net/ipv4/ip_forward

**Virtual Server 关闭 IP 包转发**

	echo 0 > /proc/sys/net/ipv4/ip_forward

**配置文件**

	! Configuration File for keepalived
	
	global_defs {
	   router_id LVS_HTTP
	}
	
	vrrp_instance VI_1 {
	    state BACKUP
	    interface ens32
	    mcast_src_ip 10.1.63.102
	    lvs_sync_daemon_inteface ens32
	    virtual_router_id 1
	    priority 50
		nopreempt
	    advert_int 1
	    authentication {
	        auth_type PASS
	        auth_pass 1111
	    }
	    virtual_ipaddress {
	        10.1.63.2
	    }
	}
	
	virtual_server 10.1.63.2 80 {
	    delay_loop 6
	    lb_algo rr 
	    lb_kind DR
	    protocol TCP
	
	    real_server 10.1.63.103 80 {
	        weight 1
	        TCP_CHECK {
	            connect_timeout 3
	            nb_get_retry 3
	            delay_before_retry 3
	            connect_port 80
	        }
	    }
	
	    real_server 10.1.63.104 80 {
	        weight 1
	        TCP_CHECK {
	            connect_timeout 3
	            nb_get_retry 3
	            delay_before_retry 3
	            connect_port 80
	        }
	    }
	}

**关键配置项的说明**

具体配置参考官方资料

- GLOBAL CONFIGURATION

	可以只保留一个 route_id，对于主备方式，该 route_id 可以相同，也可以不同

- VRRPD CONFIGURATION

	- state，主机配置为 MASTER，备机配置为 BACKUP
	- virtual_router_id，主备必须相同，主备争抢 VIP 时，会使用该参数来标识主备需要占用同一个 VIP
	- priority，优先级，主备机争抢 VIP 时会以该参数作为协商基准，数字越大优先级越高，优先级高的一方会抢到 VIP。一般配置为 MASTER 大于 BACKUP，这样 MASTER 宕掉，VIP 会切换到 BACKUP 上，MASTER 重新启动后，VIP 又会回切到 MASTER。如果 BACKUP 配置得比 MASTER 大，BACKUP 启动后会一直占用 VIP。配置成相同的时候，测试结果为 MASTER 重启后，会切换到 MASTER 上，然后 MASTER 争抢 VIP 失败，又会回到 BACKUP 上，所以不能将 priority 配成相同。
	- nopreempt，表示如果 keepalived 重新启动后，虚拟 IP 不回切
	- advert_int，相当于主备机心跳时长，单位为秒，一般配置为 1，如果配置得过大会导致 VIP 切换时间变长。测试情况：设置为 1 时，发生切换时，ping 包丢失 1-2 个，设置为 5 时，ping 包丢失 5-6 个
	- virtual_ipaddress，注意必须配置为 VIP

- LVS CONFIGURATION

	- delay_loop：Checker 链路检测时间间隔
	- lb_algo：轮徇策略，一般常用 rr(轮徇)，wrr(按权重轮徇)，lc(最小连接数)，wlc(按权重最小连接数)
	- lb_kind：LVS 模式，DR，NAT 和 TUN，配置成 DR 时，不要使用 lc 或 wlc，因为 DR 使用的是虚拟 IP 的方式，在负载均衡上没有连接状态，会导致所有请求都送到一台 Real Server 上。使用 NAT 时，需要添加 persistence_timeout，会话保持时长，以保证响应能回到正确的请求上。
	- real_server：配置相应 Real Server 的权重和链路检测参数

以上，如果需要添加一台 Real Server，在 virtual_server 配置块下增加一个 real_server 块
如果需要增加一个 Virtual Server，增加一个 virtual_server 块

### Real Server 配置

Real Server 关闭 ARP 广播响应

	echo "1" >/proc/sys/net/ipv4/conf/lo/arp_ignore
	echo "2" >/proc/sys/net/ipv4/conf/lo/arp_announce
	echo "1" >/proc/sys/net/ipv4/conf/all/arp_ignore
	echo "2" >/proc/sys/net/ipv4/conf/all/arp_announce

Real Server 恢复 ARP 广播响应

	echo "0" >/proc/sys/net/ipv4/conf/lo/arp_ignore
	echo "0" >/proc/sys/net/ipv4/conf/lo/arp_announce
	echo "0" >/proc/sys/net/ipv4/conf/all/arp_ignore
	echo "0" >/proc/sys/net/ipv4/conf/all/arp_announce

以上设置用于关闭 loopback 接口和其它所有网络接口的 ARP 广播响应，避免 Real Server 抢占 Virtual Server 的 VIP

### 测试结果

使用 ab 进行测试

测试命令

	[root@localhost ~]# ab -n 10000000 -c 10 http://10.1.63.2/

负载情况

	[root@eb63101 ~]# ipvsadm
	IP Virtual Server version 1.2.1 (size=4096)
	Prot LocalAddress:Port Scheduler Flags
	  -> RemoteAddress:Port           Forward Weight ActiveConn InActConn
	TCP  10.1.63.2:http rr
	  -> 10.1.63.103:http             Route   1      2          14113     
	  -> 10.1.63.104:http             Route   1      4          14112

## keepalived 日志

keepalived 日志存放在系统日志 /var/log/messages 中

- 查看 keepalived 所有日志

		grep Keepalived /var/log/messages

- 查看 VRRP 模块日志

		grep Keepalived_vrrp /var/log/messages

- 查看 Checkers 模块日志

		grep Keepalived_healthcheckers /var/log/messages

- 配置成服务时，也可以使用 status 查看，但只能显示最近的日志

		systemctl status keepalived
		service keepalived status

如果需要将 keepalived 日志存放在单独文件，需要对 syslog 进行配置，配置方法

- 修改 /etc/sysconfig/keepalived 配置

		KEEPALIVED_OPTIONS="-D -S 0"

	-S 0 表示将 keepalived 的 syslog facility 设置为 local0

- 修改 /etc/rsyslog.conf，添加如下内容：

		# Keepalived log
	    Keepalived.*		/var/log/keepalived.log

- 这样 keepalived 的日志会记录到 /var/log/keepalived.log 中，同时也会记录到 /var/log/messages 中。因此还需要对 /etc/rsyslog.conf 进行修改，找到如下行：

		*.info;mail.none;authpriv.none;cron.none                /var/log/messages

	修改为：

		*.info;mail.none;authpriv.none;cron.none;local0.none   /var/log/messages

	local0.none 就是表示所有 local0 的日志都不记录到 /var/log/messages 中

- 重启 syslog 服务：

		service rsyslog restart

## 参考资料

- keepalived 官方资料：[http://www.keepalived.org/pdf/UserGuide.pdf](http://www.keepalived.org/pdf/UserGuide.pdf)
- keepalived 中文资料：[http://www.keepalived.org/pdf/sery-lvs-cluster.pdf](http://www.keepalived.org/pdf/sery-lvs-cluster.pdf)