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
**Flea RESTful接口**服务端和客户端接入，本篇博文不再赘述；

可见笔者的如下的两篇文章：

- [《Flea RESTful接口服务端接入》](../../../../../../2019/11/29/flea-framework/flea-jersey/flea-jersey-server/)
- [《Flea RESTful接口客户端接入》](../../../../../../2019/12/15/flea-framework/flea-jersey/flea-jersey-client/)

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

这里的 `doPostResource` 方法实现的是**通用 POST 资源API**，可以看到该方法里面实际调用了 `FleaJerseyFGetResource` 中的 `doResource` 方法【实际上是资源父类 `Resource` 中的方法】。

```java
    /**
     * 处理资源数据 
     *
     * @param fleaJerseyRequest 请求对象
     * @return 响应对象
     * @since 1.0.0
     */
    protected FleaJerseyResponse doResource(FleaJerseyRequest fleaJerseyRequest) {
        initContext();
        return FleaJerseyFilterChainManager.getManager().doFilter(fleaJerseyRequest);
    }
```


**Flea Jersey 文件 GET 资源**，只包含**文件 GET 资源API**。

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

- 首先可以看到，**Flea Jersey 文件 GET 资源** 实现了 **Jersey 文件下载资源接口**，该接口就是提供处理文件下载资源数据的**API**。

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

- 其次也能看出，**Flea Jersey 文件 GET 资源** 继承了抽象资源父类 `Resource`，其中 `doFileDownloadResource` 的方法中，就是调用该抽象父类中的 `doCommonFileDownloadResource` 方法来实现处理文件下载资源数据的逻辑；

```java
/**
 * 处理文件下载资源数据 
 *
 * @param requestData 请求数据字符串
 * @return Jersey响应对象
 * @since 1.0.0
 */
protected FormDataMultiPart doCommonFileDownloadResource(String requestData) {

    FleaJerseyResponse fleaJerseyResponse = doResource(requestData);

    String responseData = JABXUtils.toXml(fleaJerseyResponse, false);

    FormDataMultiPart formDataMultiPart = null;

    try {
        // 将响应数据添加到表单中
        FleaJerseyManager.getManager().addFormDataBodyPart(responseData, FleaJerseyConstants.FormDataConstants.FORM_DATA_KEY_RESPONSE);
        formDataMultiPart = FleaJerseyManager.getManager().getFileContext().getFormDataMultiPart();
    } catch (CommonException e) {
        if (LOGGER.isErrorEnabled()) {
            LOGGER.error1(new Object() {}, "Exception occurs : \n", e);
        }
    }
    return formDataMultiPart;
}

/**
 * 处理资源数据 
 *
 * @param requestData 请求数据字符串
 * @return 响应对象
 * @since 1.0.0
 */
protected FleaJerseyResponse doResource(String requestData) {
    initContext();
    return FleaJerseyFilterChainManager.getManager().doFilter(requestData);
}
```

下载资源，配置参考如下：

![](flea-jersey-resource.png)

关键字段如下：

|   字段名                |    中文描述        |
|:------------------------|:----------------------|
|resource_code     | 资源编码            |
|resource_name        | 资源名称             |
|resource_packages      | 资源包名(如果存在多个，以逗号分隔) |

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

文件下载服务实现类，参考实现如下【演示使用】：

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

关键字段如下：

|   字段名                |    中文描述        |
|:------------------------|:----------------------|
|service_code        | 服务编码             |
|resource_code     | 资源编码            |
|service_name          | 服务名称          |
|service_interfaces        | 服务接口类             |
|service_method      | 服务方法       |
|service_input          | 服务入参            |
|service_output        |  服务出参  |

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

上述配置中 请求方式 为 `fget`，这里定义为**文件GET请求**，可参考枚举类 [RequestModeEnum](https://github.com/Huazie/flea-framework/blob/dev/flea-jersey/flea-jersey-client/src/main/java/com/huazie/fleaframework/jersey/client/request/RequestModeEnum.java)

```java
FGET("FGET", "com.huazie.fleaframework.jersey.client.request.impl.FGetFleaRequest", "文件GET请求")
```

文件GET请求具体实现，可至 **GitHub** 查看 [FGetFleaRequest](https://github.com/Huazie/flea-framework/blob/dev/flea-jersey/flea-jersey-client/src/main/java/com/huazie/fleaframework/jersey/client/request/impl/FGetFleaRequest.java)

**文件 GET 请求**，对外提供了执行文件下载请求的能力。

注：服务端提供的资源入口方法需包含 `GET` 注解和 `Path` 注解【如：`@Path("/fileDownload")`】，这里从 `FleaJerseyFGetResource` 也可以看出来。

```java
public class FGetFleaRequest extends FleaRequest {

    public FGetFleaRequest() {
    }

    public FGetFleaRequest(RequestConfig config) {
        super(config);
    }

    @Override
    protected void init() {
        modeEnum = RequestModeEnum.FGET;
    }

    @Override
    protected FleaJerseyResponse request(WebTarget target, FleaJerseyRequest request) throws CommonException {

        // 将请求报文转换成请求数据字符串
        String requestData = toRequestData(request);

        FleaJerseyResponse response = null;

        // 文件下载GET请求发送
        FormDataMultiPart formDataMultiPart = target
                .path(FleaJerseyConstants.FileResourceConstants.FILE_DOWNLOAD_PATH)
                .queryParam(FleaJerseyConstants.FormDataConstants.FORM_DATA_KEY_REQUEST, requestData)
                .request(toMediaType())
                .get(FormDataMultiPart.class);

        // 将表单添加到文件上下文中
        FleaJerseyManager.getManager().getFileContext().setFormDataMultiPart(formDataMultiPart);

        // 获取响应表单数据
        FormDataBodyPart responseFormData = FleaJerseyManager.getManager().getFormDataBodyPart(FleaJerseyConstants.FormDataConstants.FORM_DATA_KEY_RESPONSE);
        String responseData = responseFormData.getValue();

        if (StringUtils.isNotBlank(requestData)) {
            response = JABXUtils.fromXml(responseData, FleaJerseyResponse.class);
        }

        return response;
    }
}
```

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

有关 FleaJerseyClient 的使用，可以查看笔者的[《Flea RESTful接口客户端接入》](../../../../../../2019/12/15/flea-framework/flea-jersey/flea-jersey-client/)

# 总结

至此，文件下载接入告一段落； 

欢迎了解 [flea-jersey](/categories/开发框架-Flea/flea-jersey/) 的其他内容 。
