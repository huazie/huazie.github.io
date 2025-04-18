---
title: 在WSL2 Ubuntu中部署FastDFS服务的完整指南
date: 2025-04-13 11:04:15
updated: 2025-04-13 11:04:15
categories:
  - [轻量级分布式文件系统-FastDFS]
tags:
  - FastDFS
  - WSL2
  - Ubuntu
  - 部署指南
---

![](/images/fastdfs-logo.png)


# 📖 前言

最近，**Huazie** 在开发文件服务器项目【[FleaFS](https://github.com/huazie/FleaFS)】，使用到了 **FastDFS** 来作为底层文件管理服务。为了方便本地测试，便部署了 **FastDFS** 的服务。

由于笔者的系统是 **Windows 11 家庭版**，因此本篇教程指南将指导大家在 **WSL2 Ubuntu** 子系统中直接安装和配置 **FastDFS** 服务，无需 **Docker** 或 **虚拟机**。

**方案特点**：

- ✅ 原生 **Linux** 环境部署
- ⚡ 高性能低延迟
- 🔄 与 **Windows** 文件系统无缝互通

<!-- more -->


# 🛠️ 环境准备

## 1. 系统要求
| 组件              | 要求                |
|-------------------|------------------------------|
| Windows版本       | Windows 10 2004+ 或 Windows 11 |
| WSL版本           | WSL2                |
| 内存              | 至少2GB空闲         |
| 磁盘空间          | 至少5GB可用空间     |

如果本地的**WSL**版本较低，可以查看[【官方文档】](https://learn.microsoft.com/zh-cn/windows/wsl/install-manual)进行手动更新。

![](wsl-enable.png)


![](wsl2.png)


![](vmp-enable.png)

## 2. Ubuntu应用

进入 [Micorsoft 应用市场](https://apps.microsoft.com/search?query=ubuntu&hl=zh-CN&gl=CN&department=Apps)，搜索 **Ubuntu**，选择 **20.04 LTS** 以上的版本安装即可。

![](ubuntu-search.png)

安装好了，本地系统可搜索**Ubuntu**应用【笔者本地安装的是 **Ubuntu 24.04.2 LTS** 】

![](ubuntu1.png)


![](ubuntu2.png)

# 🚀 安装服务

## 1. 更新系统
```bash
sudo apt update && sudo apt upgrade -y
```

## 2. 安装编译依赖
```bash
sudo apt install -y git gcc make perl libperl-dev
```

## 3. 下载源码

上述系统更新和编译依赖安装好了之后，我们就需要先后输入下面命令来克隆**fastdfs**相关的项目：

```bash
git clone https://github.com/happyfish100/libfastcommon.git
git clone https://github.com/happyfish100/libserverframe.git
git clone https://github.com/happyfish100/fastdfs.git
```

![](fastdfs.png)

如果上述命令执行后，提示未找到命令，则通过 `sudo apt install git` 安装即可。多数情况下，**WSL** 的 **Ubuntu** 环境会默认配备 **Git**、**Vim** 等常用工具，以便于快速使用。


## 4. 编译安装
```bash
# 安装libfastcommon
cd libfastcommon
sudo ./make.sh && sudo ./make.sh install

# 安装libserverframe
cd libserverframe
sudo ./make.sh && sudo ./make.sh install

# 安装FastDFS
cd ../fastdfs
sudo ./make.sh && sudo ./make.sh install
```

# 🔧 配置服务

## 1. 设置配置文件

在上述的 **fastdfs** 目录中执行如下命令，用来设置配置文件（主要是将**fastdfs**相关的模板配置文件复制到`/etc/fdfs`中）。该脚本 **不会覆盖** 现有的配置文件，请放心执行这个脚本。

```bash
sudo ./setup.sh /etc/fdfs
```

执行前：

![](etc-before.png)

执行后：


![](etc-after.png)


## 2. 创建数据目录

```bash
sudo mkdir -p /data/fdfs/{client,tracker,storage,storage0}
sudo chmod -R 777 /data/fdfs
```

## 3. 配置Tracker服务

通过命令来编辑`tracker.conf`，如下：

```bash
sudo nano /etc/fdfs/tracker.conf
```

> `nano` 是 **Linux** 系统中一个简单易用的 **命令行文本编辑器**，适合快速编辑文件（尤其是配置文件或脚本）。有关它的详细使用请自行搜索，这里不再赘述。

需要修改的关键配置如下：

```ini
port=22122
base_path=/data/fdfs/tracker
```

现在我们需要获取 **Tracker** 服务器的 **IP** 地址

本篇我们的示例 **Tracker** 服务器在 **WSL** 的 **Ubuntu** 中
在 **WSL** 的 **Ubuntu** 终端里，使用 `ip addr` 命令来获取网络接口信息。

![](ipaddr.png)

其中 **inet** 后面的 **IP** 地址（如 `172.22.204.57`）就是 **WSL** 的 **Ubuntu** 的实际被分配的 **IP** 地址，这个地址可作为 **Tracker** 服务器的 **IP** 地址。


## 4. 配置Storage服务

通过命令来编辑`storage.conf`，如下：

```bash
sudo nano /etc/fdfs/storage.conf
```

需要修改的关键配置如下：

```ini
group_name=group1
port=23000
tracker_server=172.22.204.57:22122
base_path=/data/fdfs/storage
store_path0=/data/fdfs/storage0
```
其中，**tracker_server** 就配置 **3** 中的 **IP**地址，端口 `22122`

## 5. 配置客户端

通过命令来编辑`client.conf`，如下：

```bash
sudo nano /etc/fdfs/client.conf
```

需要修改的关键配置如下：
```ini
base_path=/data/fdfs/client
tracker_server=172.22.204.57:22122
```

**tracker_server** 同上修改即可。

# 🚦 启动服务

## 1. 启动Tracker
```bash
sudo fdfs_trackerd /etc/fdfs/tracker.conf restart
```

`fdfs_trackerd` ：启动 **FastDFS** 的 **Tracker Server（跟踪服务器）** ，负责管理文件存储的元数据（如 **Storage** 节点的状态、文件映射关系）和调度客户端请求。


## 2. 启动Storage
```bash
sudo fdfs_storaged /etc/fdfs/storage.conf restart
```
`fdfs_storaged` ：启动 **FastDFS** 的 **Storage Server（存储服务器）** ，负责实际文件的存储、同步和访问。


## 3. 验证服务状态

**检查进程**

```bash
ps -ef | grep fdfs
```

![](ps-ef.png)

**查看端口监听**

```bash
sudo netstat -tunlp | grep -E '22122|23000'
```

![](netstat-tunlp.png)

# ✅ 功能验证

## 1. 新建并上传测试文件

**新建测试文件**

```bash
echo "Hello WSL FastDFS，Huazie" | sudo tee /data/fdfs/test.txt
```

![](tee-test.png)

**FastDFS** 的文件上传命令需要使用 `fdfs_upload_file`，如果不清楚具体使用，我们可以先通过如下命令查看：

```bash
fdfs_upload_file --help
```

![](fdfs-upload-file-help.png)

**上传测试文件**

```bash
sudo fdfs_upload_file /etc/fdfs/client.conf /data/fdfs/test.txt
```

上传成功返回路径：`group1/M00/00/00/rBbMOWf5O7KAVMDWAAAAG3MLPyM018.txt`


![](fdfs-upload-file.png)

下面 **Huazie** 来将返回的路径详细分析下：


```plaintext
group1/M00/00/00/rBbMOWf5O7KAVMDWAAAAG3MLPyM018.txt
│      │   │  │  └── 文件名，基于文件内容哈希、时间戳等生成唯一标识
│      │   │  └── 二级子目录（每级 256 个目录，编号 00 到 FF）
│      │   └── 一级子目录（每级 256 个目录，编号 00 到 FF）
│      └── 存储路径标识符，对应配置文件中的 store_path0、store_path1 等
└── 文件所属的存储组（Group），支持多组扩展
```


## 2. 查看存储文件

从配置服务小节的配置Storage服务中，我们可以看到存储数据 `store_path0` 的配置路径 `/data/fdfs/storage0`，因此上述上传文件，可以通过如下命令查看存储路径：

```bash
ls -rtl /data/fdfs/storage0/data/00/00
```

![](ls-rtl.png)

如上图中标红的文件就是本地我们上传的文件，可以看到文件名也和上面返回的内容对应上了。

## 3. 下载文件

**FastDFS** 的文件下载命令需要使用 `fdfs_download_file`，如果不清楚具体使用，我们可以先通过如下命令查看：

```bash
fdfs_download_file --help
```

![](fdfs-download-file-help.png)

```bash
sudo fdfs_download_file /etc/fdfs/client.conf group1/M00/00/00/rBbMOWf5O7KAVMDWAAAAG3MLPyM018.txt download.txt
```

![](fdfs-download-file.png)

```bash
cat download.txt
```

![](cat-download.png)


# 🌐 外部调用

## 1. Tracker Server 的 IP 地址

这里我们在上面配置 **Tracker** 服务时已经查看过，也可通过如下命令精确定位：

```bash
ip addr show eth0 | grep inet
```

![](ipaddr1.png)


## 2. Java客户端配置

**Huazie** 这里就以 **FleaFS** 项目的配置【[fdfs_client.conf](https://github.com/huazie/FleaFS/blob/main/fleafs-config/src/main/resources/fdfs_client.conf)】为例：


```conf
connect_timeout = 2
network_timeout = 30
charset = UTF-8
http.tracker_http_port = 9090
http.anti_steal_token = no
http.secret_key = FastDFS1234567890

tracker_server = 172.22.204.57:22122
#tracker_server = 172.22.204.58:22122

## Whether to open the connection pool, if not, create a new connection every time
connection_pool.enabled = true
## max_count_per_entry: max connection count per host:port , 0 is not limit
connection_pool.max_count_per_entry = 500
## connections whose the idle time exceeds this time will be closed, unit: second, default value is 3600
connection_pool.max_idle_time = 3600
## Maximum waiting time when the maximum number of connections is reached, unit: millisecond, default value is 1000
connection_pool.max_wait_time_in_ms = 1000

```

## 3. Java客户端示例

这里可查看 **FleaFS** 项目的 **FastDFS** 客户端测试类 [FastDFSClientTest](https://github.com/huazie/FleaFS/blob/main/fleafs-web/src/test/java/com/huazie/ffs/FastDFSClientTest.java)，该类演示了上传，下载，删除文件等功能，有需要的可以参考一下。


```java
public class FastDFSClientTest {

    private static final Logger LOGGER = LoggerFactory.getLogger(FastDFSClientTest.class);

    @Test
    public void uploadFile() {
        InputStream inputStream = IOUtils.getInputStreamFromClassPath("file/绿色田园风光.jpg");
        String fileName = "绿色田园风光.jpg";
        String fileId = FastDFSClient.uploadFile(inputStream, fileName);
        LOGGER.debug("FILE_ID : {}", fileId);
    }

    @Test
    public void downloadFile() throws Exception {
        InputStream inputStream = FastDFSClient.downloadFile("group1/M00/00/00/rBbMOWfqU3CAH_cdAARILk7ifpI048.jpg");
        File file = new File("E:\\绿色.jpg");
        FileUtils.copyInputStreamToFile(inputStream, file);
    }

    @Test
    public void deleteFile() {
        FastDFSClient.deleteFile("group1/M00/00/00/rBbMOWfqU3CAH_cdAARILk7ifpI048.jpg");
    }

}
```


# 📝 结语

通过本篇的指导，相信大家已经可以成功在本地系统的 **WSL2 Ubuntu** 应用中部署原生的 **FastDFS** 文件存储服务。

如有其他问题，欢迎大家在评论区交流讨论！
 
对 **FastDFS** 感兴趣的朋友，也可以研究研究 [FastDFS官方文档](https://github.com/happyfish100/fastdfs/wiki) 。
