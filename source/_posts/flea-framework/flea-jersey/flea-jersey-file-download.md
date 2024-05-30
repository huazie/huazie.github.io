---
title: flea-jersey使用之文件下载接入
date: 2019-12-22 19:49:21
updated: 2023-06-28 22:40:36
categories:
  - [开发框架-Flea,flea-jersey]
tags:
  - flea-framework
  - flea-jersey
  - Flea RESTful接口
  - 文件下载接入
---

![](/images/flea-logo.png)


# 引言
本文将要介绍 **flea-jersey** 提供的文件下载功能。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

文件下载功能，需要引入Flea RESTful接口服务端和客户端依赖，详细如下所示：

# 1. 客户端依赖
```xml
  <!-- FLEA JERSEY CLIENT-->
    <dependency>
        <groupId>com.huazie.fleaframework</groupId>
        <artifactId>flea-jersey-client</artifactId>
        <version>2.0.0</version>
    </dependency>
```

# 2. 服务端依赖
```xml
  <!-- FLEA JERSEY SERVER-->
    <dependency>
        <groupId>com.huazie.fleaframework</groupId>
        <artifactId>flea-jersey-server</artifactId>
        <version>2.0.0</version>
    </dependency>
```

# 3. 文件下载接入讲解
Flea RESTful接口服务端和客户端接入，本篇博文不再赘述；可见笔者 [flea-jersey](/categories/开发框架-Flea/flea-jersey/) 下的文章。

## 3.1 服务端下载资源定义
下载资源 **DownloadResource** 继承文件GET资源 **FleaJerseyFGetResource** ，用于实现文件下载功能
```java
@Path("download")
public class DownloadResource extends FleaJerseyFGetResource implements JerseyPostResource {

    /**
     * @see JerseyPostResource#doPostResource(FleaJerseyRequest)
     */
    @Override
    public FleaJerseyResponse doPostResource(FleaJerseyRequest request) {
        return doResource(request);
    }
}
```

```java
public abstract class FleaJerseyFGetResource extends Resource implements JerseyFileDownloadResource {

    /**
     * @see JerseyFileDownloadResource#doFileDownloadResource(String)
     */
    @GET
    @Path("/fileDownload")
    @Override
    public FormDataMultiPart doFileDownloadResource(@QueryParam("REQUEST") String requestData) {
        return doCommonFileDownloadResource(requestData);
    }

}
```

```java
public interface JerseyFileDownloadResource {

    /**
     * <p> 处理文件下载资源数据 </p>
     *
     * @param requestData 请求数据字符串
     * @return Jersey响应对象
     * @since 1.0.0
     */
    @Consumes({MediaType.APPLICATION_JSON, MediaType.APPLICATION_XML})
    @Produces(MediaType.MULTIPART_FORM_DATA)
    FormDataMultiPart doFileDownloadResource(String requestData);
}
```
下载资源，配置参考如下：

![](flea-jersey-resource.png)

## 3.2 服务端文件下载服务定义
文件下载服务接口类，参考实现如下：

```java
public interface IFleaDownloadSV {

    /**
     * <p> 文件下载 </p>
     *
     * @param input 文件下载业务入参
     * @return 文件下载业务出参
     * @throws Exception
     * @since 1.0.0
     */
    OutputFileDownloadInfo fileDownload(InputFileDownloadInfo input) throws Exception;
}
```
文件下载服务实现类，参考实现如下：

```java
@Service
public class FleaDownloadSVImpl implements IFleaDownloadSV {

    private static final Logger LOGGER = LoggerFactory.getLogger(FleaDownloadSVImpl.class);

    @Override
    public OutputFileDownloadInfo fileDownload(InputFileDownloadInfo input) throws Exception {

        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug("FleaDownloadSVImpl#fileDownload(InputFileDownloadInfo) Start");
        }

        File file = new File("E:\\IMG.jpg");
        // 将文件添加到文件上下文中
        FleaJerseyManager.getManager().addFileDataBodyPart(file);

        OutputFileDownloadInfo output = new OutputFileDownloadInfo();
        output.setUploadAcctId("121212");
        output.setUploadSystemAcctId("1000");
        output.setUploadDate(DateUtils.date2String(null, DateFormatEnum.YYYYMMDDHHMMSS));

        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug("FleaDownloadSVImpl#fileDownload(InputFileDownloadInfo) End");
        }

        return output;
    }
}
```
文件下载服务，配置参考如下：

![](flea-jersey-res-service.png)

## 3.3 客户端文件下载配置
文件下载客户端，配置参考如下：

![](flea-jersey-res-client.png)

|   字段名                |    中文描述        |
|:------------------------|:----------------------|
|client_code          | 客户端编码          |
|resource_url        | 资源地址             |
|resource_code     | 资源编码            |
|service_code        | 服务编码             |
|request_mode      | 请求方式       |
|media_type          | 媒体类型            |
|client_input        |  客户端业务入参  |
|client_output      |  客户端业务出参  |

上述配置中 请求方式 为 fget，这里定义为文件GET请求，可参考枚举类 [RequestModeEnum](https://github.com/Huazie/flea-framework/blob/dev/flea-jersey/flea-jersey-client/src/main/java/com/huazie/fleaframework/jersey/client/request/RequestModeEnum.java)

```java
FGET("FGET", "com.huazie.fleaframework.jersey.client.request.impl.FGetFleaRequest", "文件GET请求")
```
文件GET请求具体实现，可至 GitHub查看 [FGetFleaRequest](https://github.com/Huazie/flea-framework/blob/dev/flea-jersey/flea-jersey-client/src/main/java/com/huazie/fleaframework/jersey/client/request/impl/FGetFleaRequest.java)

## 3.4 客户端文件下载调用
文件下载自测类，可至GitHub查看 [JerseyTest](https://github.com/Huazie/flea-framework/blob/dev/flea-jersey/flea-jersey-client/src/test/java/com/huazie/fleaframework/jersey/client/resource/JerseyTest.java)

```java
  @Test
    public void testDownloadFile() {
        try {
            String clientCode = "FLEA_CLIENT_FILE_DOWNLOAD";

            InputFileDownloadInfo input = new InputFileDownloadInfo();
            input.setToken(RandomCode.toUUID());

            FleaJerseyClient client = applicationContext.getBean(FleaJerseyClient.class);

            Response<OutputFileDownloadInfo> response = client.invoke(clientCode, input, OutputFileDownloadInfo.class);

            LOGGER.debug("result = {}", response);

            OutputFileDownloadInfo output = response.getOutput();

            // 获取文件信息
            FleaFileObject fileObject = FleaJerseyManager.getManager().getFileObject();
            String fileName = fileObject.getFileName();
            File downloadFile = fileObject.getFile();

            String uploadSystemAcctId = output.getUploadSystemAcctId();
            String uploadAcctId = output.getUploadAcctId();
            String uploadDate = output.getUploadDate();

            if (downloadFile.exists()) {
                IOUtils.toFile(new FileInputStream(downloadFile), "E:\\" + uploadDate + "_" + uploadSystemAcctId + "_" + uploadAcctId + "_" + fileName);
            }

        } catch (Exception e) {
            LOGGER.error("Exception = ", e);
        }
    }
```

# 总结

至此，文件下载接入告一段落； 

欢迎了解  [flea-jersey](/categories/开发框架-Flea/flea-jersey/)  其他内容 。
