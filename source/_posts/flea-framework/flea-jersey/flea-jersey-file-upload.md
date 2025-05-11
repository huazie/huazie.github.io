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
**Flea RESTful接口**服务端和客户端接入，本篇博文不再赘述；

可见笔者的如下的两篇文章：

- [《Flea RESTful接口服务端接入》](../../../../../../2019/11/29/flea-framework/flea-jersey/flea-jersey-server/)
- [《Flea RESTful接口客户端接入》](../../../../../../2019/12/15/flea-framework/flea-jersey/flea-jersey-client/)

## 3.1 服务端上传资源定义
上传资源继承 **FleaJerseyPostResource**，该类定义可至GitHub查看  [flea-jersey-server](https://github.com/Huazie/flea-framework/tree/dev/flea-jersey/flea-jersey-server)，具体如下所示：

```java
@Path("upload")
public class UploadResource extends FleaJerseyPostResource {

}
```

**Flea Jersey POST** 资源，包含 **通用 POST 资源API**，**文件上传 POST 资源API**。

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

这里的 `doPostResource` 方法实现的是**通用 POST 资源API**，可以看到该方法里面实际调用了 `FleaJerseyFPostResource` 中的 `doResource` 方法【实际上是资源父类 `Resource` 中的方法】。

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

**Flea Jersey 文件 POST 资源**，只包含**文件 POST 资源API**。

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

- 首先可以看到，**Flea Jersey 文件 POST 资源** 实现了 **Jersey** 文件上传资源接口，该接口就是提供处理文件上传**POST**资源数据的**API**。

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

- 其次也能看出，**Flea Jersey 文件 POST 资源** 继承了抽象资源父类 `Resource`，其中 `doFileUploadResource` 的方法中，就是调用该抽象父类中的 `doCommonFileUploadResource` 方法来实现处理文件上传资源数据的逻辑；

```java
/**
 * 处理文件上传资源数据 
 *
 * @param formDataMultiPart 表单数据
 * @return 响应对象
 * @since 1.0.0
 */
protected FleaJerseyResponse doCommonFileUploadResource(FormDataMultiPart formDataMultiPart) {

    FleaJerseyResponse fleaJerseyResponse = null;

    if (ObjectUtils.isNotEmpty(formDataMultiPart)) {

        // 生成文件上下文
        FleaJerseyFileContext fleaJerseyFileContext = new FleaJerseyFileContext();
        fleaJerseyFileContext.setFormDataMultiPart(formDataMultiPart);
        // 设置文件上下文
        FleaJerseyManager.getManager().getContext().setFleaJerseyFileContext(fleaJerseyFileContext);

        // 获取请求表单数据
        FormDataBodyPart formDataBodyPart = formDataMultiPart.getField(FleaJerseyConstants.FormDataConstants.FORM_DATA_KEY_REQUEST);
        // 获取请求参数
        String requestData = formDataBodyPart.getValueAs(String.class);
        // 处理文件上传资源数据
        fleaJerseyResponse = doResource(requestData);
    }

    return fleaJerseyResponse;
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

上传资源，配置参考如下：

![](flea-jersey-resource.png)

关键字段如下：

|   字段名                |    中文描述        |
|:------------------------|:----------------------|
|resource_code     | 资源编码            |
|resource_name        | 资源名称             |
|resource_packages      | 资源包名(如果存在多个，以逗号分隔) |

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

文件上传服务实现类，参考实现如下【演示使用】：

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

这里通过 `FleaJerseyManager` 获取 `FleaFileObject`，其中包含了客户端上传的文件信息。

文件上传服务，配置参考如下：

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

## 3.3 客户端文件上传配置

文件上传客户端，配置参考如下：

![](flea-jersey-res-client.png)

关键字段如下：

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

上述配置中 请求方式 为 `fpost`，这里定义为**文件POST请求**，可参考枚举类 [RequestModeEnum](https://github.com/Huazie/flea-frame/blob/master/flea-frame-jersey/flea-frame-jersey-client/src/main/java/com/huazie/frame/jersey/client/request/RequestModeEnum.java)

```java
FPOST("FPOST", "com.huazie.fleaframework.jersey.client.request.impl.FPostFleaRequest", "文件POST请求")
```

文件POST请求具体实现，可至 **GitHub** 查看 [FPostFleaRequest](https://github.com/Huazie/flea-framework/blob/dev/flea-jersey/flea-jersey-client/src/main/java/com/huazie/fleaframework/jersey/client/request/impl/FPostFleaRequest.java)

**文件 POST 请求**，对外提供了执行文件上传请求的能力。

注：服务端提供的资源入口方法需包含 `POST` 注解和 `Path` 注解【如：`@Path("/fileUpload")`】，这里从 `FleaJerseyFPostResource` 也可以看出来。

```java
public class FPostFleaRequest extends FleaRequest {

    public FPostFleaRequest() {
    }

    public FPostFleaRequest(RequestConfig config) {
        super(config);
    }

    @Override
    protected void init() {
        modeEnum = RequestModeEnum.FPOST;
    }

    @Override
    protected FleaJerseyResponse request(WebTarget target, FleaJerseyRequest request) throws CommonException {

        String requestData = JABXUtils.toXml(request, false);
        // 添加请求表单数据
        FleaJerseyManager.getManager().addFormDataBodyPart(requestData, FleaJerseyConstants.FormDataConstants.FORM_DATA_KEY_REQUEST);

        Entity<FormDataMultiPart> entity = Entity.entity(FleaJerseyManager.getManager().getFileContext().getFormDataMultiPart(), toMediaType());

        // 文件上传POST请求发送
        FleaJerseyResponse response = target
                .path(FleaJerseyConstants.FileResourceConstants.FILE_UPLOAD_PATH)
                .request()
                .post(entity, FleaJerseyResponse.class);

        return response;
    }
}
```

## 3.4 客户端文件上传调用
文件上传自测类，可至**GitHub**查看 [JerseyTest](https://github.com/Huazie/flea-framework/blob/dev/flea-jersey/flea-jersey-client/src/test/java/com/huazie/fleaframework/jersey/client/resource/JerseyTest.java)

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

有关 `FleaJerseyClient` 的使用，可以查看笔者的[《Flea RESTful接口客户端接入》](../../../../../../2019/12/15/flea-framework/flea-jersey/flea-jersey-client/)

# 总结

至此，文件上传已接入完毕，下篇博文将会讲解文件下载接入。