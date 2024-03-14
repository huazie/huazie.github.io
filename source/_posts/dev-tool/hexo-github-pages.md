---
title: 【实操】基于 GitHub Pages + Hexo 搭建个人博客
date: 2024-01-17 15:14:55
updated: 2024-01-17 15:14:55
categories:
  - [开发工具]
tags:
  - Hexo
  - GitHub Pages
  - GitHub Actions
  - 个人博客搭建
---

[《开发工具系列》](/categories/开发工具/) 

![](/images/hexo-githubpages.png)

# 一、引言

相信很多学习技术的读者朋友们，都梦想能创建一个属于自己的个人博客。现在，这将不是梦想，下面跟着 **Huazie** 一起利用 **GitHub Pages + Hexo** 搭建一个属于自己的个人博客吧。

# 二、接入 Node.js 

## 2.1 下载并安装 Node.js 

[Node.js 官方下载地址](https://nodejs.org/en/download/)

> 注意：**Hexo** 官方建议使用 **Node.js 12.0** 及以上版本

笔者本地下载的是 **20.11.0 LTS**，这对大多数用户来说已经足够了

![](/images/dev-tool/nodejs-download.png)
 
笔者的 **Windows** 系统，下载完了是如下的 msi 安装包【其他系统自行去官网下载即可】：
![](/images/dev-tool/node-install-package.png)

这里直接双击安装即可，安装完了就可以去配置相关的环境变量了。

## 2.2 环境变量配置

现在，**Huazie** 以 **windows 11** 系统为例，介绍下配置环境变量，如下：

右击 **Window** 图标，打开下图并选择 **系统**：

![](/images/dev-tool/windows-system.png)


点击 **高级系统设置**，打开系统属性页面，点击 **环境变量** ：

![](/images/dev-tool/windows-env-config.png)

找到 **Path** 系统环境变量，配置上面你的 **Node.js** 的安装目录进去：

![](/images/dev-tool/nodejs-env-config.png)

环境变量配置好之后，我们就可以通过 CMD 命令行，检查：

- `npm -v` ：查看当前安装的 **npm** 的版本号
![](/images/dev-tool/npm-v.png)
- `node -v` : 查看当前安装的 **Node.js** 的版本号
![](/images/dev-tool/node-v.png)

# 三、接入 Git

## 3.1 下载并安装 Git

[Windows 下载地址](https://git-scm.com/download/win)，其他可参考 [【Hexo 官方文档里的安装 Git】](https://hexo.io/zh-cn/docs/)

笔者本地下载的版本如下【大家从上述地址下载的版本比我本地的高些】：
![](git-install-package.png)

这里也是一样直接双击安装即可，安装完了就可以去配置相关的环境变量了。

## 3.2 环境变量配置
我们先来看看 **Git** 的安装目录：

![](git-directory.png)
在上述的 **bin** 和 **cmd** 目录，我们都可以看到 **git.exe**，按需配置，我本地环境配置的是 **cmd** 目录。

参考上面 **Node.js** 环境变量配置，配置好之后，我们就可以在命令行输入如下命令查看：

![](git-version.png)

# 四、接入 Hexo

## 4.1 安装 Hexo

接入 **Node.js** 和 **Git** 之后，我们就可以使用 `npm` 安装 **Hexo**。

```shell
npm install -g hexo-cli
```
![](hexo-cli-install.png)

上述安装成功后，提示我 **npm** 有新的小版本更新，于是我进行了更新：

![](npm-install-new.png)
- `npm install -g npm` ：更新到最新版本
- `npm install -g npm@<version>` ：更新到特定的版本


这时我再查看当前安装的 **npm** 的版本号：

![](npm-v-10.3.0.png)
> 注意：上述更新不强制，大家按需更新即可

当然，对于熟悉 **npm** 的进阶用户，可以仅局部安装 **hexo** 包。

```shell
npm install hexo
```

安装 **Hexo** 以后，可以使用以下两种方式执行 **Hexo**：

- `npx hexo <command>`
- **Linux** 用户可以将 **Hexo** 所在的目录下的 **node_modules** 添加到环境变量之中即可直接使用 `hexo <command>`：	

	```shell
	echo 'PATH="$PATH:./node_modules/.bin"' >> ~/.profile
	```
## 4.2 建站

```shell
# 没有设置 folder 参数，Hexo 默认在当前文件夹下创建网站
hexo init <folder>
```

我们需要选个本地文件夹，然后输入上述命令，用于在指定文件夹下初始化一个本地网站。

下图即为 **Huazie** 本地在 `E:\fleaworkspace\blog` 目录开始初始化一个博客网站：
![](hexo-init.png)

因为要从 **GitHub** 克隆项目，这一步可能需要花点事件，请慢慢等待，不要关闭窗口

等待一会，如果如下图显示，就表示 **hexo** 初始化网站成功了。

![](hexo-init-1.png)

接着我们切换到上述初始化的网站目录，当然如果按笔者上述操作，当前目录就是我们的网站根目录。

接着我们输入 `npm install` 命令，用来下载我们网站必要的依赖包。

![](npm-install-package.png)

`npm install` 命令的作用包括：
1. **从 npm 注册表下载包**：`npm install` 会从 `npm` 注册表（一个在线仓库）中查找并下载指定的包。你可以指定包的名称和版本号，以获取正确的包版本。
2. **解析依赖**：`npm install` 会解析项目中的 `package.json` 文件，读取其中的 **dependencies** 和 **devDependencies** 字段，确定需要安装的依赖项及其版本。它会下载并安装所有必要的依赖项，以确保项目的正常运行。
3. **安装本地缓存**：`npm install` 会将下载的包和依赖项安装到项目的本地缓存中，这样其他开发者也可以共享相同的依赖项版本，确保项目的可移植性和一致性。
4. **生成 node_modules 目录**：在安装完成后，`npm install` 会生成一个 **node_modules** 目录，其中包含所有安装的包和依赖项

上述操作完成之后，可以查看我们初始化的网站目录，如下所示：

![](hexo-website-directory.png)

有关上述文件，我们这里简单介绍下：
- `_config.yml` ：网站的配置信息。
- `package.json` ：应用程序的信息。
- `scaffolds` ：模版文件夹。当您新建文章时，**Hexo** 会根据 **scaffold** 来创建文件。**Hexo** 的模板是指在新建的文章文件中默认填充的内容。例如，如果您修改 `scaffold/post.md` 中的 **Front-matter** 内容，那么每次新建一篇文章时都会包含这个修改。
- `source` ：资源文件夹是存放用户资源的地方。除 `_posts` 文件夹之外，开头命名为 _ (下划线)的文件 / 文件夹和隐藏的文件将会被忽略。`Markdown` 和 `HTML` 文件会被解析并放到 public 文件夹，而其他文件会被拷贝过去。
- `themes` ：主题文件夹。Hexo 会根据主题来生成静态页面。


## 4.3 本地启动服务器

我们可以在本地启动服务器。如下所示：
![](hexo-server.png)

执行完之后，不要关闭命令窗口，直接在浏览器打开 [http://localhost:4000/](http://localhost:4000/)，如下图所示：

![](hexo-server-localhost.png)


当然还有很多其他的命令，感兴趣的小伙伴，请查看 [官方指令文档](https://hexo.io/zh-cn/docs/commands)。


# 五、接入 GitHub Pages

## 5.1 初识 GitHub Pages

**GitHub Pages** 是一项静态站点托管服务，它直接从 **GitHub** 上的仓库获取 **HTML**、**CSS** 和 **JavaScript** 文件，（可选）通过构建过程运行文件，然后发布网站。 可以在 [GitHub Pages 示例集合](https://github.com/collections/github-pages-examples) 中看到 GitHub Pages 站点的示例。

你可以在 **GitHub** 的 `github.io` 域或自己的自定义域上托管站点。 有关详细信息，请参阅“[配置 GitHub Pages 站点的自定义域](https://docs.github.com/zh/pages/configuring-a-custom-domain-for-your-github-pages-site)”。

GitHub Pages 站点的类型，有三种：

- **项目** ：项目站点连接到 **GitHub** 上托管的特定项目，例如 **JavaScript** 库或配方集合
- **用户** ：用户站点连接到 **github.com** 上的特定帐户。若要发布用户站点，必须创建名为 `<username>.github.io` 的个人帐户拥有的存储库。
- **组织** ：组织站点连接到 **github.com** 上的特定帐户。若要发布组织站点，必须创建名为 `<organization>.github.io` 的组织帐户拥有的存储库。

> 除非使用的是自定义域，否则用户和组织站点在 `http(s)://<username>.github.io` 或 `http(s)://<organization>.github.io` 中可用。

> **GitHub Pages** 使用限制：
> **2016 年 6 月 15** 日后创建并使用 `github.io` 域的 **GitHub Pages** 站点通过 **HTTPS** 提供服务。 如果您在 **2016 年 6 月 15** 日之前创建站点，您可以为站点的流量启用 **HTTPS** 支持。 有关详细信息，请参阅“[使用 HTTPS 保护 GitHub Pages 站点](https://docs.github.com/zh/pages/getting-started-with-github-pages/securing-your-github-pages-site-with-https)”。

可以在将更改推送到特定分支时发布站点，也可以编写 **GitHub Actions** 工作流来发布站点。对于在 GitHub Pages 上部署 Hexo，请查看 [《官方文档》](https://hexo.io/zh-cn/docs/github-pages)，它就是使用 [GitHub Actions](https://docs.github.com/zh/actions) 部署至 **GitHub Pages**。

## 5.2 在 GitHub Pages 上部署 Hexo

下面 **Huazie** 来简单总结下：

1. 在你的 **GitHub** 上建立名为 <你的 GitHub 用户名>.github.io 的仓库。这里参考  [《GitHub Pages 快速入门》](https://docs.github.com/zh/pages/quickstart) 即可。
2. 使用 **GitHub 客户端** 克隆上述新建的仓库，并将 **4.2 中初始化的目录内容** 全部复制到新克隆的仓库中，或者 像官方那样自己推送到远端【参考[《在 GitHub Pages 上部署 Hexo》](https://hexo.io/zh-cn/docs/github-pages)】。
3. 在上面新克隆的仓库目录下，新建立 `.github/workflows/pages.yml` 【目录如果没有自己新建即可】
![](pages-yml.png)
`pages.yml` 中填入以下内容 (注意下面的 **Node.js** 的版本，我这里是 **20**，大家以自己本地安装的版本为准)：

	```yml
	name: Pages
	
	on:
	  push:
	    branches:
	      - main # default branch
	
	jobs:
	  build:
	    runs-on: ubuntu-latest
	    steps:
	      - uses: actions/checkout@v3
	        with:
	          token: ${{ secrets.GITHUB_TOKEN }}
	          # If your repository depends on submodule, please see: https://github.com/actions/checkout
	          submodules: recursive
	      - name: Use Node.js 20.x
	        uses: actions/setup-node@v2
	        with:
	          node-version: '20'
	      - name: Cache NPM dependencies
	        uses: actions/cache@v2
	        with:
	          path: node_modules
	          key: ${{ runner.OS }}-npm-cache
	          restore-keys: |
	            ${{ runner.OS }}-npm-cache
	      - name: Install Dependencies
	        run: npm install
	      - name: Build
	        run: npm run build
	      - name: Upload Pages artifact
	        uses: actions/upload-pages-artifact@v2
	        with:
	          path: ./public
	  deploy:
	    needs: build
	    permissions:
	      pages: write
	      id-token: write
	    environment:
	      name: github-pages
	      url: ${{ steps.deployment.outputs.page_url }}
	    runs-on: ubuntu-latest
	    steps:
	      - name: Deploy to GitHub Pages
	        id: deployment
	        uses: actions/deploy-pages@v2
	```
4. 使用 **GitHub** 客户端将上述仓库新增的文件推送到远端。
  ![](github-desktop.png)

5. 前往 **GitHub** 仓库，按下图顺序 **Settings > Pages > Source** ，并将 **Source** 改为 **GitHub Actions**。
  ![](github-actions.png)

6. 接着等待 **GitHub** 自动部署，然后就可以通过 `https://你的GitHub用户名.github.io/` 访问了
  ![](github-io.png)

# 六、总结

本篇 **Huazie** 带大家利用 **GitHub Pages + Hexo** 搭建了能访问的个人博客。一步步实操下来，相信大家都能见到实际的效果。当然要想做好个人博客，可不止这么一点点，**Huazie** 这里也只是抛砖引玉，后续的深入使用，需要发挥各位的主观能动性了。
# 七、参考

1. [《Hexo 官方文档》](https://hexo.io/zh-cn/docs/)
2. [《GitHub Actions 文档》](https://docs.github.com/zh/actions)
3. [《GitHub Pages 快速入门》](https://docs.github.com/zh/pages/quickstart)