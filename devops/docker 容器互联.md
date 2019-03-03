# docker 容器互联
---

docker 容器互联总结的是在同一台宿主机上，多个 docker 容器文件共享和网络互联的方式

## docker 数据卷

docker 容器在创建时，使用了容器内部的文件系统，容器内部的进程可以对其进行访问。但是在很多场景下，需要从容器外部访问容器的文件(从宿主机或从其它容器)，如：

- 采集容器中 nginx 的日志进行分析
- mysql 的 data 数据需要持久化到宿主机上，防止容器挂掉后数据丢失

docker 可以在创建或启动时指定容器中相应的目录为数据卷：

- 在 Dockerfile 中指定 VOLUME
- 在运行容器是使用 -v 参数，如：

		docker run -v /var/log/nginx nginx

挂载数据卷类似于 Linux 下挂载文件系统，实际上是将宿主机上的目录或文件映射到容器内部，这样容器内部对数据的修改，都可以在数据机上反应出来，如上面完成以后，在宿主机上多出了一个文件夹，里面对应了 nginx 的日志目录：

	[root@iZ25dcta17gZ ~]# ls /var/lib/docker/volumes/f5994c19f011663445ce343bd3019216b553de7fba86ff1e337e1d7a2a898e77/_data/
	access.log  error.log

即时容器结束，该目录仍会存在，从而保证数据不会丢失，如果需要删除可以使用：

	docker volume rm f5994c19f011663445ce343bd3019216b553de7fba86ff1e337e1d7a2a898e77

一般情况下，使用数据卷的方式由于少了一层虚拟层，可以比容器文件系统的速度更快，对需要大量读写的文件目录，建议使用挂载数据卷

但是像上面的这种方式，产生的数据卷是随机的，每次启动容器都会变化，不便于管理，一般情况下，我们会指定挂载宿主机的路径，如：

	docker run -v /root/log:/var/log/nginx -p 8081:80 -d nginx

-v 参数指定了将本地目录或文件映射到容器中的目录或文件，冒号前面是宿主机的路径，后面是容器中的路径

在生产环境下，往往会有很多宿主机，同样的容器会启动到不同的宿主机上，也有可能启动到同一个宿主机上，指定宿主机的固定目录，可能会产生冲突。这种情况，则需要引入容器组，使用一个独立的容器作为数据卷容器，和其一组的容器将数据写入到数据容器中。如 kubernetes 中，Pod 即为一个容器组，他有一个 Pause 容器，如果 pod 需要持久化数据，pause 容器会挂载相应的数据卷，其他容器和该容器进行连接。

容器直接互相访问，使用 --volumes-from 参数，如：

1. 使用 centos 启动一个数据卷容器：

		docker run -it --rm --name log -v /var/log/nginx centos:7 /bin/bash

2. 使用 nginx 将目录挂载到该数据卷容器上：

		docker run --volumes-from log --name nginx -p 8081:80 -d nginx

	这样访问 nginx，即可在 log 容器中看到相应的日志

如果我们需要采集 nginx 日志，可以使用同样的方法，将 logstash 或 filebeat 绑定到 log 容器下。这里 log 作为一个基础的数据容器，将工作容器 nginx，日志采集容器 filebeat 整合为一个容器组。这种方式和 kubernetes 的  pod 就比较相像了

同样的使用该方法也可以对 mysql 数据进行备份。

## docker 简单网络

和数据卷相同，docker 也提供了相应的方式将容器内的端口映射到宿主机上：

- 在 Dockerfile 中指定 EXPOSE
- 在运行容器是使用 -p 参数，如：

		docker run -p 80 nginx

和 -v 参数一样，这种映射方式，映射出的是随机参数，使用 docker ps 可以看到映射到宿主机上的端口，如下例中的 32768：

	[root@iZ25dcta17gZ ~]# docker ps
	CONTAINER ID        IMAGE                 COMMAND                  CREATED             STATUS              PORTS                                            NAMES
	78d4cf3b3048        nginx                 "nginx -g 'daemon ..."   36 seconds ago      Up 36 seconds       0.0.0.0:32768->80/tcp                            flamboyant_clarke

同样，如果需要指定宿主机上的指定端口，可以使用如下命令：

	docker run -p 8081:80 nginx

-p 参数，冒号前面的宿主机端口，冒号后面为容器中的端口，这样我们在宿主机上访问 8081 端口，相当于访问容器中的 80 端口

同样的，可以使用 --link 参数实现容器间的网络互访，下面我们以一个简单的 elasticsearch 和 kibana 的互联为例

1. 启动一个 elasticsearch 容器(这里使用的是单机模式，不能用于生产环境)

		docker run -d --name elasticsearch -e "discovery.type=single-node" elasticsearch:6.4.0

2. 使用 --link 将 kibana 连接到 es 上：

		docker run -d --name kibana --link elasticsearch:es -p 5601:5601 kibana:6.4.0

这里，--link 参数，冒号前为要连接的容器名，冒号后为该容器在本容器的别名

接入容器后，可以看到，在 kibana 的 /etc/hosts 中增加了对 es 的域名：

	[root@iZ25dcta17gZ ~]# docker exec -it a75ae7d9b4f7 /bin/bash
	bash-4.2$
	bash-4.2$
	bash-4.2$
	bash-4.2$ cat /etc/host
	cat: /etc/host: No such file or directory
	bash-4.2$ cat /etc/hosts
	127.0.0.1	localhost
	::1	localhost ip6-localhost ip6-loopback
	fe00::0	ip6-localnet
	ff00::0	ip6-mcastprefix
	ff02::1	ip6-allnodes
	ff02::2	ip6-allrouters
	192.168.0.2	es a02fa05bbc23 elasticsearch
	192.168.0.3	a75ae7d9b4f7

这里的网络只是在同一宿主机上的容器间互联，在生产环境下面，容器间的相互访问，以及外部访问容器中的服务，后续再进行讨论