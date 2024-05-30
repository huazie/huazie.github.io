---
title: flea-jersey使用之文件上传接入
date: 2019-12-18 08:50:12
updated: 2023-06-28 22:27:51
categories:
  - [开发框架-Flea,flea-jersey]
tags:
  - flea-framework
  - flea-jersey
  - Flea RESTful接口
  - 文件上传接入
---

![](/images/flea-logo.png)

# 引言
本文将要介绍 **flea-jersey** 提供的文件上传功能。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

文件上传功能，需要引入 **Flea RESTful** 接口服务端和客户端依赖，详细如下所示：

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

# 3. 文件上传接入讲解
Flea RESTful接口服务端和客户端接入，本篇博文不再赘述；可见笔者 [flea-jersey](/categories/开发框架-Flea/flea-jersey/) 下的文章。

## 3.1 服务端上传资源定义
上传资源继承 **FleaJerseyPostResource**，该类定义可至GitHub查看  [flea-jersey-server](https://github.com/Huazie/flea-framework/tree/dev/flea-jersey/flea-jersey-server)，具体如下所示：

```java
@Path("upload")
public class UploadResource extends FleaJerseyPostResource {

}
```

```java
public abstract class FleaJerseyPostResource extends FleaJerseyFPostResource implements JerseyPostResource {

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
public abstract class FleaJerseyFPostResource extends Resource implements JerseyFileUploadResource {

    /**
     * @see JerseyFileUploadResource#doFileUploadResource(FormDataMultiPart)
     */
    @POST
    @Path("/fileUpload")
    @Override
    public FleaJerseyResponse doFileUploadResource(FormDataMultiPart formDataMultiPart) {
        return doCommonFileUploadResource(formDataMultiPart);
    }

}
```

```java
public interface JerseyFileUploadResource {

    /**
     * <p> 处理文件上传POST资源数据 </p>
     *
     * @param formDataMultiPart 表单数据
     * @return 响应对象
     * @since 1.0.0
     */
    @Consumes(MediaType.MULTIPART_FORM_DATA)
    @Produces({MediaType.APPLICATION_JSON, MediaType.APPLICATION_XML})
    FleaJerseyResponse doFileUploadResource(FormDataMultiPart formDataMultiPart);

}
```

上传资源，配置参考如下：

![](flea-jersey-resource.png)

## 3.2 服务端文件上传服务定义
文件上传服务接口类，参考实现如下：

```java
public interface IFleaUploadSV {

    /**
     * <p> 文件上传 </p>
     *
     * @param input 文件上传入参（包含上传鉴权令牌）
     * @return 文件上传出参（包含文件编号）
     * @throws Exception
     * @since 1.0.0
     */
    OutputFileUploadInfo fileUpload(InputFileUploadInfo input) throws Exception;
}
```
文件上传服务实现类，参考实现如下：

```java
@Service
public class FleaUploadSVImpl implements IFleaUploadSV {

    private static final Logger LOGGER = LoggerFactory.getLogger(FleaUploadSVImpl.class);
    
    @Override
    public OutputFileUploadInfo fileUpload(InputFileUploadInfo input) throws Exception {

        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug("FleaUploadSVImpl#fileUpload(InputFileUploadInfo) Start");
        }

        String token = input.getToken();
        if (StringUtils.isBlank(token)) {
            // 入参【{0}】不能为空
            throw new ServiceException("ERROR-SERVICE0000000001", "上传鉴权令牌【token】");
        }

        FleaFileObject fileObject = FleaJerseyManager.getManager().getFileObject();
        String fileName = fileObject.getFileName();
        File uploadFile = fileObject.getFile();

        String fileId = DateUtils.date2String(null, DateFormatEnum.YYYYMMDD) + RandomCode.toUUID();
        IOUtils.toFile(new FileInputStream(uploadFile), "E:\\" + fileId + "_" +fileName);
        OutputFileUploadInfo outputFileUploadInfo = new OutputFileUploadInfo();
        outputFileUploadInfo.setFileId(fileId);

        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug("FleaUploadSVImpl#fileUpload(InputFileUploadInfo) Start");
        }
        return outputFileUploadInfo;
    }
}
```
文件上传服务，配置参考如下：

![](flea-jersey-res-service.png)

## 3.3 客户端文件上传配置
文件上传客户端，配置参考如下：

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

上述配置中 请求方式 为 fpost，这里定义为文件POST请求，可参考枚举类 [RequestModeEnum](https://github.com/Huazie/flea-frame/blob/master/flea-frame-jersey/flea-frame-jersey-client/src/main/java/com/huazie/frame/jersey/client/request/RequestModeEnum.java)

```java
FPOST("FPOST", "com.huazie.fleaframework.jersey.client.request.impl.FPostFleaRequest", "文件POST请求")
```
文件POST请求具体实现，可至 GitHub查看 [FPostFleaRequest](https://github.com/Huazie/flea-framework/blob/dev/flea-jersey/flea-jersey-client/src/main/java/com/huazie/fleaframework/jersey/client/request/impl/FPostFleaRequest.java)

## 3.4 客户端文件上传调用
文件上传自测类，可至GitHub查看 [JerseyTest](https://github.com/Huazie/flea-framework/blob/dev/flea-jersey/flea-jersey-client/src/test/java/com/huazie/fleaframework/jersey/client/resource/JerseyTest.java)

```java
  @Test
    public void testUploadFile() {
        try {
            String clientCode = "FLEA_CLIENT_FILE_UPLOAD";

            InputFileUploadInfo input = new InputFileUploadInfo();
            input.setToken(RandomCode.toUUID());
      // 添加待上传文件至文件上下文对象中
            File file = new File("E:\\IMG.jpg");
            FleaJerseyManager.getManager().addFileDataBodyPart(file);

            FleaJerseyClient client = applicationContext.getBean(FleaJerseyClient.class);
      // 调用文件上传服务
            Response<OutputFileUploadInfo> response = client.invoke(clientCode, input, OutputFileUploadInfo.class);

            LOGGER.debug("result = {}", response);
        } catch (Exception e) {
            LOGGER.debug("Exception = ", e);
        }
    }
```

# 总结

至此，文件上传已接入完毕，下篇博文将会讲解文件下载接入。


