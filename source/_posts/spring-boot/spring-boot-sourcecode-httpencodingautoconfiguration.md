---
title: 【Spring Boot 源码学习】HttpEncodingAutoConfiguration 详解
date: 2023-10-22 22:33:08 
updated: 2023-11-05 22:03:24
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - HttpEncodingAutoConfiguration
  - 自定义字符编码映射
---

[《Spring Boot 源码学习系列》](/categories/开发框架-Spring-Boot/)

![](/images/spring-boot-logo.png)

# 引言
前面的博文，我们从源码角度介绍了自动装配流程。虽然带大家从整体上有了清晰的认识，但是我们还不能熟练地运用。

本篇就以 **Spring Boot** 内置的 `http` 编码功能为例，来带大家分析一下 `HttpEncodingAutoConfiguration` 的整个自动配置的过程。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)
# 往期内容
在开始本篇的内容介绍之前，我们先来看看往期的系列文章【有需要的朋友，欢迎关注系列专栏】：

<table>
    <tr>
        <td rowspan="13" align="left"> 
            <a href="/categories/开发框架-Spring-Boot/">Spring Boot 源码学习</a> 
        </td>
    </tr>
    <tr>
        <td align="left"> 
            <a href="/2023/02/19/spring-boot/spring-boot-project-introduction/">Spring Boot 项目介绍</a> 
        </td>
    </tr>
    <tr>
        <td align="left"> 
            <a href="/2023/07/13/spring-boot/spring-boot-core-operating-principle/">Spring Boot 核心运行原理介绍</a> 
        </td>
    </tr>
    <tr>
        <td align="left"> 
            <a href="/2023/07/16/spring-boot/spring-boot-sourcecode-springbootapplication/">【Spring Boot 源码学习】@SpringBootApplication 注解</a> 
        </td>
    </tr>
    <tr>
        <td align="left"> 
            <a href="/2023/07/22/spring-boot/spring-boot-sourcecode-enableautoconfiguration/">【Spring Boot 源码学习】@EnableAutoConfiguration 注解</a> 
        </td>
    </tr>
    <tr>
        <td align="left"> 
            <a href="/2023/07/30/spring-boot/spring-boot-sourcecode-autoconfigurationimportselector/">【Spring Boot 源码学习】走近 AutoConfigurationImportSelector</a> 
        </td>
    </tr>
    <tr>
        <td align="left" > 
            <a href="/2023/08/06/spring-boot/spring-boot-sourcecode-autoconfigurationdetail-1/">【Spring Boot 源码学习】自动装配流程源码解析（上）</a> 
        </td>
    </tr>
    <tr>
        <td align="left" > 
            <a href="/2023/08/21/spring-boot/spring-boot-sourcecode-autoconfigurationdetail-2/">【Spring Boot 源码学习】自动装配流程源码解析（下）</a> 
        </td>
    </tr>
    <tr>
        <td align="left" > 
            <a href="/2023/09/08/spring-boot/spring-boot-sourcecode-filteringspringbootcondition/">【Spring Boot 源码学习】深入 FilteringSpringBootCondition</a> 
        </td>
    </tr>
    <tr>
        <td align="left" > 
            <a href="/2023/09/11/spring-boot/spring-boot-sourcecode-onclasscondition/">【Spring Boot 源码学习】OnClassCondition 详解</a> 
        </td>
    </tr>
    <tr>
        <td align="left" > 
            <a href="/2023/09/21/spring-boot/spring-boot-sourcecode-onbeancondition/">【Spring Boot 源码学习】OnBeanCondition 详解</a> 
        </td>
    </tr>
    <tr>
        <td align="left" > 
            <a href="/2023/10/06/spring-boot/spring-boot-sourcecode-onwebapplicationcondition/">【Spring Boot 源码学习】OnWebApplicationCondition 详解</a> 
        </td>
    </tr>
    <tr>
        <td align="left" > 
            <a href="/2023/10/15/spring-boot/spring-boot-sourcecode-conditional/">【Spring Boot 源码学习】@Conditional 条件注解</a> 
        </td>
    </tr>
</table>



# 主要内容
## 1. CharacterEncodingFilter 
在传统的 **web** 项目中，**Spring** 为 **web** 开发提供的一个过滤器【即 `org.springframework.web.filter.CharacterEncodingFilter` 】，用来防止 **web** 开发中出现的乱码问题，它是 **Spring** 通过在 **web** 请求中定义 **request** 和 **response** 的编码来实现。

在 `web.xml` 中的配置示例如下：

```java
    <filter>  
        <filter-name>encodingFilter</filter-name>  
        <filter-class>org.springframework.web.filter.CharacterEncodingFilter
        </filter-class>  
        <init-param>  
             <param-name>encoding</param-name>  
             <param-value>UTF-8</param-value>  
        </init-param>  
        <init-param>  
             <param-name>forceEncoding</param-name>  
             <param-value>true</param-value>  
        </init-param>  
    </filter>
```



## 2. HttpEncodingAutoConfiguration

**那么在 **Spring Boot** 是如何实现的呢？**

**Spring Boot**  是通过内置的 `HttpEncodingAutoConfiguration` 配置类来完成这一功能。下面我们具体分析一下：

> **注意：** 以下涉及 **Spring Boot** 源码 均来自版本 `2.7.9`，其他版本有所出入，可自行查看源码。

### 2.1 加载自动配置组件
从之前的[《【Spring Boot 源码学习】自动装配流程源码解析（上）》](/2023/08/06/spring-boot/spring-boot-sourcecode-autoconfigurationdetail-1/)中，我们知道 Spring Boot 内部针对自动配置类，会读取如下两个配置文件：

- `META-INF/spring.factories`
- `META-INF/spring/org.springframework.boot.autoconfigure.AutoConfiguration.imports`

![](autoconfiguration.png)

实际上 在 `Spring Boot 2.7.9` 版本中， **Spring Boot** 自己内部的 `META-INF/spring.factories` 中有关自动配置的注册类的配置信息已经被去除掉了，不过其他外围的 **jar** 中可能有自己的 `META-INF/spring.factories` 文件，它里面也有关于自动配置注册类的配置信息；

而 Spring Boot 内置的 `HttpEncodingAutoConfiguration` 配置类，则是配置在上述的第二个配置文件 `META-INF/spring/org.springframework.boot.autoconfigure.AutoConfiguration.imports` 中。

![](httpencodingautoconfiguration.png)

### 2.2 过滤自动配置组件

上述自动配置加载完之后，就来到了 [《【Spring Boot 源码学习】自动装配流程源码解析（下）》](/2023/08/21/spring-boot/spring-boot-sourcecode-autoconfigurationdetail-2/) 介绍的 **过滤自动配置组件** 逻辑。

这部分数据对应的配置内容在 `META-INF/spring-autoconfigure-metadata.properties` 文件中：

```java
org.springframework.boot.autoconfigure.web.servlet.HttpEncodingAutoConfiguration=
org.springframework.boot.autoconfigure.web.servlet.HttpEncodingAutoConfiguration.ConditionalOnClass=org.springframework.web.filter.CharacterEncodingFilter
org.springframework.boot.autoconfigure.web.servlet.HttpEncodingAutoConfiguration.ConditionalOnWebApplication=SERVLET
```

显然这里涉及到了 `ConditionalOnClass` 和 `ConditionalOnWebApplication` 注解，我们翻看 `HttpEncodingAutoConfiguration` 配置类的源码，如下：

```java
@AutoConfiguration
@EnableConfigurationProperties(ServerProperties.class)
@ConditionalOnWebApplication(type = ConditionalOnWebApplication.Type.SERVLET)
@ConditionalOnClass(CharacterEncodingFilter.class)
@ConditionalOnProperty(prefix = "server.servlet.encoding", value = "enabled", matchIfMissing = true)
public class HttpEncodingAutoConfiguration {

    private final Encoding properties;

    public HttpEncodingAutoConfiguration(ServerProperties properties) {
        this.properties = properties.getServlet().getEncoding();
    }

    @Bean
    @ConditionalOnMissingBean
    public CharacterEncodingFilter characterEncodingFilter() {
        //  。。。
    }

    @Bean
    public LocaleCharsetMappingsCustomizer localeCharsetMappingsCustomizer() {
        // 。。。
    }
    
    // ...
}
```

#### 2.2.1 涉及注解 

我们先来看看上述 `HttpEncodingAutoConfiguration` 配置类涉及到的注解，如下：
- `@AutoConfiguration` : 该类是一个自动配置类，**Spring Boot** 会根据项目中的依赖自动配置这个类的实例。
- `@EnableConfigurationProperties(ServerProperties.class)` ：启用 `ServerProperties` 类的配置属性，这样在配置文件中就可以使用 `server.servlet.encoding` 属性来配置字符编码。
- `@ConditionalOnWebApplication(type = ConditionalOnWebApplication.Type.SERVLET)` ：该配置类只有在基于 **servlet** 的 **web** 应用程序中才会被实例化。
- `@ConditionalOnClass(CharacterEncodingFilter.class)` ：只有在项目中存在 `CharacterEncodingFilter` 类时才会生效。
- `@ConditionalOnProperty(prefix = "server.servlet.encoding", value = "enabled", matchIfMissing = true)` ：只有在配置文件中 `server.servlet.encoding` 属性的值为 `"enabled"` 时才会生效；当然如果配置文件中没有这个属性，也默认会生效。
- `@Bean` ：用于声明一个方法创建的对象是一个 **Spring** 管理的 **Bean**。**Spring** 容器会自动管理这个 **Bean** 的生命周期，包括依赖注入、初始化和销毁等。
- `@ConditionalOnMissingBean` ：只有在当前 **Spring** 容器中不存在指定类型的 **Bean** 时，才会执行被注解的方法。这样可以用于确保在需要的时候才创建某个 **Bean**，避免重复创建。


#### 2.2.2 ServerProperties

其中 `ServerProperties` 类的属性值对应着 `application.yml` 或 `application.properties` 中的配置，通过注解`@ConfigurationProperties(prefix = "server", ignoreUnknownFields = true)` 实现的属性注入。

有关属性注入的内容后续笔者会另外介绍，我们先来看看`ServerProperties` 类相关的部分源码 和 对应的配置参数：
```java
@ConfigurationProperties(prefix = "server", ignoreUnknownFields = true)
public class ServerProperties {
    // 。。。
    private final Servlet servlet = new Servlet();
    // 。。。
    public static class Servlet {
        // 。。。
        @NestedConfigurationProperty
        private final Encoding encoding = new Encoding();
        // 。。。
    }
    // 。。。
}

public class Encoding {
    // 默认的HTTP编码，用于Servlet应用程序。
    public static final Charset DEFAULT_CHARSET = StandardCharsets.UTF_8;

    // HTTP请求和响应的字符集。如果未显式设置，将添加到"Content-Type"头中
    private Charset charset = DEFAULT_CHARSET;

    // 是否强制在HTTP请求和响应上使用配置的字符集的标志
    private Boolean force;

    // 是否强制在HTTP请求上使用配置的字符集的标志。当"force"未指定时，默认为true。
    private Boolean forceRequest;

    // 是否强制在HTTP响应上使用配置的字符集的标志。
    private Boolean forceResponse;

    // 将区域设置映射到字符集以进行响应编码的映射。
    private Map<Locale, Charset> mapping;
    // 。。。
}
```

当然在 `application.properties` 中，我们就可以添加如下的配置：

```java
server.servlet.encoding.force=true
server.servlet.encoding.charset=UTF-8
# server.servlet.encoding.force-request=true 
# ...其他配置省略
```

> **注意：** `server.servlet.encoding.force=true` 和 `server.servlet.encoding.force-request=true` 这两个配置项实际上具有相同的功能，它们都决定是否强制对客户端请求进行字符编码。当这些配置项设置为 `true`时，服务器将要求客户端发送的请求内容使用指定的字符集进行编码。
> 另外，从 **Spring Boot 2.3.5** 版本开始，`server.servlet.encoding.enabled` 配置项已被弃用。因此，推荐的做法是直接设置 `server.servlet.encoding.charset` 来指定字符集，然后通过设置`server.servlet.encoding.force=true` 来开启对请求/响应的编码集强制控制。

### 2.3 characterEncodingFilter 方法
先来看看 `characterEncodingFilter` 方法的源码【Spring Boot 2.7.9】：

```java
public CharacterEncodingFilter characterEncodingFilter() {
    CharacterEncodingFilter filter = new OrderedCharacterEncodingFilter();
    filter.setEncoding(this.properties.getCharset().name());
    filter.setForceRequestEncoding(this.properties.shouldForce(Encoding.Type.REQUEST));
    filter.setForceResponseEncoding(this.properties.shouldForce(Encoding.Type.RESPONSE));
    return filter;
}
```
上述逻辑很好理解：
- 首先，新建一个 `CharacterEncodingFilter` 的实例对象 `filter` 。
- 然后，设置 `filter` 的 `encoding` 属性，即编码属性。其中 `this.properties.getCharset().name()` 就是上述 `application.properties` 中的 `server.servlet.encoding.charset=UTF-8`；如果没有配置，则默认是 **UTF-8**【可查看上述 `Encoding` 类】。
- 接着，设置 `filter` 的 `forceRequestEncoding` 和 `forceResponseEncoding` 属性。我们来直接查 `Encoding` 类的 `shouldForce` 即可：

    ```java
    public boolean shouldForce(Type type) {
        // Http请求，则取 server.servlet.encoding.force-request 配置
        // Http响应，则取 server.servlet.encoding.force-response 配置
        Boolean force = (type != Type.REQUEST) ? this.forceResponse : this.forceRequest;
        // 如果上述配置都没有
        if (force == null) {
            // 取 server.servlet.encoding.force 配置
            force = this.force;
        }
        if (force == null) {
            // 当 server.servlet.encoding.force 配置也未指定时，
            // 默认 强制在HTTP请求上使用配置的字符集。
            force = (type == Type.REQUEST);
        }
        return force;
    }
    ```

### 2.4 localeCharsetMappingsCustomizer 方法

话不多说，直接来看相关的源码【**Spring Boot 2.7.9**】

```java
    @Bean
    public LocaleCharsetMappingsCustomizer localeCharsetMappingsCustomizer() {
        return new LocaleCharsetMappingsCustomizer(this.properties);
    }

    static class LocaleCharsetMappingsCustomizer
            implements WebServerFactoryCustomizer<ConfigurableServletWebServerFactory>, Ordered {

        private final Encoding properties;

        LocaleCharsetMappingsCustomizer(Encoding properties) {
            this.properties = properties;
        }

        @Override
        public void customize(ConfigurableServletWebServerFactory factory) {
            if (this.properties.getMapping() != null) {
                factory.setLocaleCharsetMappings(this.properties.getMapping());
            }
        }

        @Override
        public int getOrder() {
            return 0;
        }
    }
```

上述 `LocaleCharsetMappingsCustomizer` 静态内部类实现了 `WebServerFactoryCustomizer` 接口，该接口是用于自定义 **Web** 服务器工厂的策略接口。此接口类型的任何 **bean** 都将在服务器本身启动之前获得与服务器工厂的回调，从而我们可以设置端口、地址、错误页面等。
> **注意：** 对此接口的调用通常由 `WebServerFactoryCustomizerBeanPostProcessor` 执行，它是一个 `BeanPostProcessor`（处于 `ApplicationContext` 生命周期中的非常早期的时候）。比较安全的做法是在包含 `BeanFactory` 中延迟查找依赖项，而不是使用 `@Autowired` 注入它们。

而 `LocaleCharsetMappingsCustomizer` 类实现的 `customize` 方法，则用于设置自定义字符编码映射，这就不得不提 `server.servlet.encoding.mapping` 配置属性。

默认情况下，**Spring Boot** 会根据请求头的 `Accept-Charset` 来设置响应的字符编码。但是，有时候我们可能需要根据不同的请求路径或请求参数来进行不同的字符编码映射。这时，就可以使用 `server.servlet.encoding.mapping` 来实现自定义的字符编码映射。

```java
# 当请求路径以 /en/ 开头时，将字符编码设置为 UTF-8；当请求路径以 /zh/ 开头时，将字符编码设置为 GBK。
server.servlet.encoding.mapping=/en/**=UTF-8,/zh/**=GBK
```

> **注意：** `server.servlet.encoding.mapping` 的配置优先级高于 `server.servlet.encoding.charset` 和 `server.servlet.encoding.force`。因此，如果同时存在多个配置项，`server.servlet.encoding.mapping` 会覆盖其他配置项。

# 总结
本篇我们以 **Spring Boot** 内置的 `http` 编码功能为例来分析一下整个自动配置的过程，深入讲解了 `HttpEncodingAutoConfiguration` 配置类的相关内容。相信大家后续在看其他配置类，也能知其所以然了。

