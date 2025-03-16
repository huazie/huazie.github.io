---
title: Windows部署多个Memcached和Redis服务
date: 2019-08-30 08:24:34
updated: 2023-07-11 11:52:39
categories:
  - [开发框架-Flea,flea-cache]
tags:
  - Windows
  - Memcached
  - Redis
---

![](/images/cache.png)

# 引言
考虑到各位读者朋友的本地开发环境大部分都是在 `windows` 系列系统上，本篇博文着眼于介绍如何在`Windows` 部署多个 `Memcached` 和 `Redis` 服务【这里不是 `Redis` 集群服务，以后有机会介绍在此基础上部署 `Redis` 集群服务】，以方便用于本地应用测试接入`Memcached` 和`Redis`。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

好了，废话不多说，我们开始下面的内容：
# 1. 部署多个Memcached服务
## 1.1 准备工作
相关安装已上传，大家可以直接从 百度网盘 下载 `Cache.rar`
链接：[http://pan.baidu.com/s/1pLSJ2Tt](http://pan.baidu.com/s/1pLSJ2Tt)  
密码：k8gj

下图是 `Memcached` 的相关文件，其中 `memcached.exe` 就是我们下面服务需要的可执行文件：

![](memcached.png)

## 1.2 创建服务
以管理员模式运行 `CMD`，执行如下命令，其中 `%MEMCACHED_PATH%` 为`memcached.exe` 的磁盘路径。如下这边创建了两个 `Memcached` 服务，分别是`memcached1` 和 `memcached2`，开放端口分别是 `31113` 和 `31114` ，这些服务名和端口可自行修改。

```cmd
sc create memcached1 start= auto binPath= "\"%MEMCACHED_PATH%\memcached.exe\" -d runservice -m 512 -c 2048 -p 31113"  DisplayName= "memcached1"

sc create memcached2 start= auto binPath= "\"%MEMCACHED_PATH%\memcached.exe\" -d runservice -m 512 -c 2048 -p 31114"  DisplayName= "memcached2"
```
这个时候，可以在服务面板（`ctrl+r` 输入 `services.msc`）查看，如下就是我创建的 `2` 个`memcached`服务，服务启动类型可以自行调整属性，赶快试试系统接入吧

![多个Memcached服务](memcached-service.png)

## 1.3 删除服务
删除服务除了服务面板直接删除，当然删除服务也可以通过命令执行，如下：
```cmd
sc delete memcached1
sc delete memcached2
```

# 2. 部署多个Redis服务
## 2.1 准备工作
相关安装已上传，大家可以直接从 百度网盘 下载 `Cache.rar`
链接：[http://pan.baidu.com/s/1pLSJ2Tt](http://pan.baidu.com/s/1pLSJ2Tt)  
密码：k8gj

下图是 `Redis` 的相关文件，其中 `redis-server.exe` 就是我们下面服务安装所需要的可执行文件：

![](redis.png)

## 2.2 创建服务
以管理员模式运行 `CMD`，切换到上述准备的 `Redis` 目录中，执行如下命令【服务名和端口可自行修改】：
```cmd
redis-server --service-install redis.windows.conf --loglevel verbose --service-name redis1 --port 10001
redis-server --service-install redis.windows.conf --loglevel verbose --service-name redis2 --port 10002
redis-server --service-install redis.windows.conf --loglevel verbose --service-name redis3 --port 10003
```
这个时候，可以在服务面板查看，如下就是我创建的3个Redis服务，服务启动类型可以自行调整属性，让我们来看看吧

![多个Redis服务](redis-service.png)

## 2.3 卸载服务
以管理员模式运行 `CMD`，切换到上述准备的 `Redis` 目录中，执行如下命令：
```cmd
redis-server --service-uninstall --service-name redis1 --port 10001
redis-server --service-uninstall --service-name redis2 --port 10002
redis-server --service-uninstall --service-name redis3 --port 10003
```

# 总结
好了，`Windows` 部署多个 `Memcached` 和 `Redis` 服务的相关内容已经介绍完毕，各位可以执行起来，用于本地测试验证了。
