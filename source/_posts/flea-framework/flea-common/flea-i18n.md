---
title: flea-common使用之本地国际化实现
date: 2021-12-13 08:45:00
updated: 2024-02-20 15:56:19
categories:
  - [开发框架-Flea,flea-common]
tags:
  - flea-framework
  - flea-common
  - Flea I18N
  - 本地国际化实现
---

![](/images/flea-logo.png)


# 引言
百度百科针对 **国际化** 的解释：

<!-- more -->

![](internationalization.png)

**本地国际化**，就是指应用程序根据所处语言环境的不同【如 **Java** 中可用 **国际化标识类** `java.util.Locale` 区分不同语言环境】，自动匹配应用内置的相应的语言环境下的资源配置【如 **Java** 中可用 **资源包类** `java.util.ResourceBundle` 来匹配】，从而获取并对外展示相应的语言环境下的资源信息。

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

话不多说，直接上干货：

# 1. 依赖

```xml
    <!-- FLEA COMMON-->
    <dependency>
        <groupId>com.huazie.fleaframework</groupId>
        <artifactId>flea-common</artifactId>
        <version>2.0.0</version>
    </dependency>
```

# 2. 实现
上面提到了 Java 中 的 [**国际化标识类** `java.util.Locale`](https://docs.oracle.com/javase/8/docs/api/java/util/Locale.html) 和  [**资源包类** `java.util.ResourceBundle`](https://docs.oracle.com/javase/8/docs/api/java/util/ResourceBundle.html)，这两者就是本地国际化实现的关键所在。

## 2.1 定义国际化资源相关配置 

[flea-config.xml](https://github.com/Huazie/flea-framework/blob/dev/flea-config/src/main/resources/flea/flea-config.xml) 用于特殊配置国际化资源的路径和文件前缀。

```xml
<flea-config>
    <!-- flea-common -->
    <config-items key="flea-i18n-config" desc="Flea国际化相关配置">
        <config-item key="error" desc="error国际化资源特殊配置，指定路径和文件前缀，逗号分隔">flea/i18n,flea_i18n</config-item>
    </config-items>
</flea-config>
```

## 2.2 定义Flea I18N 配置类 
在使用 [FleaI18nConfig](https://github.com/Huazie/flea-framework/blob/dev/flea-common/src/main/java/com/huazie/fleaframework/common/i18n/config/FleaI18nConfig.java) 之前，我们先了解下Flea国际化资源文件的组成，主要有如下 5 部分：

![](flea-i18n-config.png)

上述国际化资源也可以配置默认资源文件，即文件名中不需要包含国际化标识 。例如： `flea/i18n/flea_i18n_error.properties`

> **注意：** 国际化资源文件扩展名必须为 **properties**

好了，基础的认知有了，我们开始了解 **FleaI18nConfig**，如下贴出了实现：

```java
/**
 * Flea I18N 配置类，用于获取指定语言环境下的指定资源对应的国际化数据。
 *
 * <p> 它默认读取资源路径为 flea/i18n，资源文件前缀为 flea_i18n，当然
 * 也可以在 flea-config.xml 中为指定资源文件配置路径和前缀，从而可以
 * 实现读取任意位置的资源数据。
 *
 * @author huazie
 * @version 2.0.0
 * @since 1.0.0
 */
public class FleaI18nConfig {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(FleaI18nConfig.class);

    private static volatile FleaI18nConfig config;

    private ConcurrentMap<String, String> resFilePath = new ConcurrentHashMap<>(); // 资源文件路径集

    private ConcurrentMap<String, ResourceBundle> resources = new ConcurrentHashMap<>(); // 资源集

    /**
     * 只允许通过 getConfig() 获取 Flea I18N 配置类实例
     */
    private FleaI18nConfig() {
        init(); // 初始化资源文件相关配置
    }

    /**
     * 获取 Flea I18N 配置类实例
     *
     * @return Flea I18N 配置类实例
     * @since 1.0.0
     */
    public static FleaI18nConfig getConfig() {
        if (ObjectUtils.isEmpty(config)) {
            synchronized (FleaI18nConfig.class) {
                if (ObjectUtils.isEmpty(config)) {
                    config = new FleaI18nConfig();
                }
            }
        }
        return config;
    }

    /**
     * 初始化资源名和资源文件相关属性的映射关系
     *
     * @since 1.0.0
     */
    private void init() {
        ConfigItems fleaI18nItems = FleaConfigManager.getConfigItems(CommonConstants.FleaI18NConstants.FLEA_I18N_CONFIG_ITEMS_KEY);
        if (ObjectUtils.isNotEmpty(fleaI18nItems) && CollectionUtils.isNotEmpty(fleaI18nItems.getConfigItemList())) {
            for (ConfigItem configItem : fleaI18nItems.getConfigItemList()) {
                if (ObjectUtils.isNotEmpty(configItem) && StringUtils.isNotBlank(configItem.getKey()) && StringUtils.isNotBlank(configItem.getValue())) {
                    String[] valueArr = StringUtils.split(configItem.getValue(), CommonConstants.SymbolConstants.COMMA);
                    if (ArrayUtils.isNotEmpty(valueArr) && CommonConstants.NumeralConstants.INT_TWO == valueArr.length) {
                        // 获取资源文件路径
                        String filePath = StringUtils.trim(valueArr[0]);
                        // 获取资源文件前缀
                        String fileNamePrefix = StringUtils.trim(valueArr[1]);
                        if (StringUtils.isNotBlank(filePath) && StringUtils.isNotBlank(fileNamePrefix)) {
                            String configResFilePath;
                            // 如果资源文件路径最后没有 "/"，自动添加
                            if (CommonConstants.SymbolConstants.SLASH.equals(StringUtils.subStrLast(filePath, 1))) {
                                configResFilePath = filePath + fileNamePrefix;
                            } else {
                                configResFilePath = filePath + CommonConstants.SymbolConstants.SLASH + fileNamePrefix;
                            }
                            resFilePath.put(configItem.getKey(), configResFilePath);
                        }
                    }
                }
            }
        }
        // 添加默认资源文件路径
        String defaultResFilePath = CommonConstants.FleaI18NConstants.FLEA_I18N_FILE_PATH +
                CommonConstants.FleaI18NConstants.FLEA_I18N_FILE_NAME_PREFIX; // 默认资源文件路径（仅包含公共的部分）
        resFilePath.put(CommonConstants.SymbolConstants.ASTERISK, defaultResFilePath);
    }

    /**
     * 通过国际化数据的key，获取当前系统指定资源的国际化资源；
     * 其中国际化资源中使用 {} 标记的，需要values中的数据替换。
     *
     * @param key     国际化资源KEY
     * @param values  待替换字符串数组
     * @param resName 资源名
     * @param locale  国际化标识
     * @return 国际化资源数据
     * @since 2.0.0
     */
    public FleaI18nData getI18NData(String key, String[] values, String resName, Locale locale) {
        return new FleaI18nData(key, this.getI18NDataValue(key, values, resName, locale));
    }

    /**
     * 通过国际化数据的key，获取当前系统指定资源的国际化资源
     *
     * @param key     国际化资源KEY
     * @param resName 资源名
     * @param locale  国际化标识
     * @return 国际化资源数据
     * @since 1.0.0
     */
    public FleaI18nData getI18NData(String key, String resName, Locale locale) {
        return new FleaI18nData(key, this.getI18NDataValue(key, resName, locale));
    }

    /**
     * <p> 通过国际化数据的key，获取当前系统指定资源的国际化资源数据 </p>
     *
     * @param key     国际化资源KEY
     * @param values  国际化资源数据替换内容
     * @param resName 资源名
     * @param locale  国际化标识
     * @return 国际化资源数据
     * @since 1.0.0
     */
    public String getI18NDataValue(String key, String[] values, String resName, Locale locale) {
        String value = getI18NDataValue(key, resName, locale);
        if (ArrayUtils.isNotEmpty(values)) {
            StringBuilder builder = new StringBuilder(value);
            for (int i = 0; i < values.length; i++) {
                StringUtils.replace(builder, CommonConstants.SymbolConstants.LEFT_CURLY_BRACE + i + CommonConstants.SymbolConstants.RIGHT_CURLY_BRACE, values[i]);
            }
            value = builder.toString();
        }
        return value;
    }

    /**
     * <p> 通过国际化数据的key，获取当前系统指定资源的国际化资源数据 </p>
     *
     * @param key     国际化资源KEY
     * @param resName 资源名
     * @param locale  国际化标识
     * @return 国际化资源数据
     * @since 1.0.0
     */
    public String getI18NDataValue(String key, String resName, Locale locale) {
        Object obj = null;
        if (LOGGER.isDebugEnabled()) {
            obj = new Object() {
            };
            LOGGER.debug1(obj, "Find the key     : {}", key);
            LOGGER.debug1(obj, "Find the resName : {}", resName);
            LOGGER.debug1(obj, "Find the locale  : {} , {}", locale == null ? Locale.getDefault() : locale, locale == null ? Locale.getDefault().getDisplayLanguage() : locale.getDisplayLanguage());
        }
        ResourceBundle resource = getResourceBundle(resName, locale);

        String value = null;
        if (ObjectUtils.isNotEmpty(resource)) {
            value = resource.getString(key);
            if (StringUtils.isBlank(value)) { // 如果取不到数据，则使用key返回
                value = key;
            }
        }

        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug1(obj, "Find the value   : {} ", value);
        }
        return value;
    }

    /**
     * <p> 根据资源名和国际化标识获取指定国际化配置ResourceBundle对象 </p>
     *
     * @param resName 资源名
     * @param locale  国际化标识
     * @return 国际化配置ResourceBundle对象
     * @since 1.0.0
     */
    private ResourceBundle getResourceBundle(String resName, Locale locale) {

        String key = generateKey(resName, locale);

        Object obj = null;
        if (LOGGER.isDebugEnabled()) {
            obj = new Object() {
            };
            LOGGER.debug1(obj, "Find the resKey  : {}", key);
        }

        ResourceBundle resource = resources.get(key);

        // 获取资源文件名
        StringBuilder fileName = new StringBuilder(getResFilePath(resName));
        if (StringUtils.isNotBlank(resName)) {
            fileName.append(CommonConstants.SymbolConstants.UNDERLINE).append(resName);
        }

        if (LOGGER.isDebugEnabled()) {
            if (ObjectUtils.isEmpty(locale)) {
                LOGGER.debug1(obj, "Find the expected fileName: {}.properties", fileName);
            } else {
                LOGGER.debug1(obj, "Find the expected fileName: {}_{}.properties", fileName, locale);
            }
        }

        // 获取资源文件
        if (ObjectUtils.isEmpty(resource)) {
            if (ObjectUtils.isEmpty(locale)) {
                resource = ResourceBundle.getBundle(fileName.toString());
            } else {
                resource = ResourceBundle.getBundle(fileName.toString(), locale);
            }
            resources.put(key, resource);
        }

        if (LOGGER.isDebugEnabled()) {
            Locale realLocale = resource.getLocale();
            if (ObjectUtils.isEmpty(locale) || StringUtils.isBlank(realLocale.toString())) {
                LOGGER.debug1(obj, "Find the real fileName: {}.properties", fileName);
            } else {
                LOGGER.debug1(obj, "Find the real fileName: {}_{}.properties", fileName, realLocale);
            }
        }

        return resource;
    }

    /**
     * <p> 获取国际化资源文件KEY </p>
     * <p> 如果资源名不为空，则资源名作为key，同时如果国际化标识不为空，则取资源名+下划线+国际化语言作为key；
     *
     * @param resName 资源名
     * @param locale  国际化标识
     * @return 国际化资源文件KEY
     * @since 1.0.0
     */
    private String generateKey(String resName, Locale locale) {
        String key = "";
        if (StringUtils.isNotBlank(resName)) {
            key = resName;
            if (ObjectUtils.isNotEmpty(locale)) {
                key += CommonConstants.SymbolConstants.UNDERLINE + locale;
            }
        }
        return key;
    }

    /**
     * <p> 根据资源名，获取资源文件路径 </p>
     *
     * @param resName 资源名
     * @return 资源文件路径
     * @since 1.0.0
     */
    private String getResFilePath(String resName) {
        // 首先根据资源名，从 资源文件路径集中获取
        String resFilePathStr = resFilePath.get(resName);
        if (ObjectUtils.isEmpty(resFilePathStr)) {
            // 取默认资源文件路径
            resFilePathStr = resFilePath.get(CommonConstants.SymbolConstants.ASTERISK);
        }
        return resFilePathStr;
    }

}

```

## 2.3 定义Flea I18N 工具类 
[FleaI18nHelper](https://github.com/Huazie/flea-framework/blob/dev/flea-common/src/main/java/com/huazie/fleaframework/common/i18n/FleaI18nHelper.java) 封装了 I18N 资源数据获取的静态方法，主要包含如下4种：

```java
    public static String i18n(String key, String resName, Locale locale) {
        return FleaI18nConfig.getConfig().getI18NDataValue(key, resName, locale);
    }

    public static String i18n(String key, String[] values, String resName, Locale locale) {
        return FleaI18nConfig.getConfig().getI18NDataValue(key, values, resName, locale);
    }

    // 实际在调用该方法之前，可以通过 FleaFrameManager.getManager().setLocale(Locale) 设置当前线程的国际化标识。
    public static String i18n(String key, String resName) {
        return i18n(key, resName, FleaFrameManager.getManager().getLocale());
    }

    // 实际在调用该方法之前，可以通过 FleaFrameManager.getManager().setLocale(Locale) 设置当前线程的国际化标识。
    public static String i18n(String key, String[] values, String resName) {
        return i18n(key, values, resName, FleaFrameManager.getManager().getLocale());
    }

    // 其他是对具体资源的封装，如错误码资源error、授权资源auth 和 公共信息资源common
```

## 2.4 定义Flea I18N资源枚举 

[FleaI18nResEnum](https://github.com/Huazie/flea-framework/blob/dev/flea-common/src/main/java/com/huazie/fleaframework/common/i18n/FleaI18nResEnum.java) 定义了 `Flea I18N` 的资源文件类型

```java
/**
 * Flea I18N 资源枚举
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
public enum FleaI18nResEnum {

    ERROR("error", "异常信息国际码资源文件类型"),
    ERROR_CORE("error_core", "FLEA CORE异常信息国际码资源文件类型"),
    ERROR_DB("error_db", "FLEA DB异常信息国际码资源文件类型"),
    ERROR_JERSEY("error_jersey", "FLEA JERSEY异常信息国际码资源文件类型"),
    ERROR_AUTH("error_auth", "FLEA AUTH异常信息国际码资源文件类型"),
    AUTH("auth", "FLEA AUTH 国际码资源文件类型"),
    COMMON("common", "公共信息国际码资源文件类型");

    private String resName;
    private String resDesc;

    /**
     * <p> 资源文件类型枚举构造方法 </p>
     *
     * @param resName 资源名
     * @param resDesc 资源描述
     * @since 1.0.0
     */
    FleaI18nResEnum(String resName, String resDesc) {
        this.resName = resName;
        this.resDesc = resDesc;
    }

    public String getResName() {
        return resName;
    }

    public String getResDesc() {
        return resDesc;
    }

}
```
简单的介绍之后，初步了解了本地国际化的实现，下面就需要来实际测试一下了。

话不多说，开始操刀：
# 3. 自测
首先，我们先添加几个国际化配置文件，如下：

![](flea-i18n-config-file.png)

|                  资源文件                        			  |          国际化标识（语言环境）                |
|----------------------------------------------------------|------------------------------------------------------ |
| flea/i18n/flea_i18n_error.properties   		   |                          默认                                 |
| flea/i18n/flea_i18n_error_zh_CN.properties |                         中文（简体）                  |
| flea/i18n/flea_i18n_error_en_US.properties |                         英文（美式）                   |

> **注意：** 笔者电脑的本地语言环境为 **中文（简体）**。

## 3.1 匹配指定语言

```java
    @Test
    public void fleaI18nHelperTest1() {
        String value = FleaI18nHelper.i18n("ERROR0000000001", "error", Locale.US);
        LOGGER.debug("Value = {}", value);
    }
```

**测试结果：**

![](result1.png)

## 3.2 匹配本地语言

```java
    @Test
    public void fleaI18nHelperTest() {
        String value = FleaI18nHelper.i18n("ERROR0000000001", "error", Locale.FRANCE);
        LOGGER.debug("Value = {}", value);
    }
```
**测试结果：**

![](result2.png)

## 3.3 匹配默认资源

首先，我们将本地语言的资源文件删除，如下：

![](flea-i18n-config-file1.png)

```java
    @Test
    public void fleaI18nHelperTest() {
        String value = FleaI18nHelper.i18n("ERROR0000000001", "error", Locale.FRANCE);
        LOGGER.debug("Value = {}", value);
    }
```

**测试结果：**

![](result3.png)

## 3.4 无资源匹配
首先，我们将**本地语言** 和 **默认** 的 资源文件删除，如下：

![](flea-i18n-config-file2.png)

```java
    @Test
    public void fleaI18nHelperTest() {
        String value = FleaI18nHelper.i18n("ERROR0000000001", "error", Locale.FRANCE);
        LOGGER.debug("Value = {}", value);
    }
```
**测试结果：**

![](result4.png)

# 4. 接入
上面演示了 如何通过 **FleaI18nHelper** 获取本地国际化的资源数据，下面我们来看看在异常类中接入错误码国际化资源。

## 4.1 定义通用异常类

[CommonException](https://github.com/Huazie/flea-framework/blob/dev/flea-common/src/main/java/com/huazie/fleaframework/common/exception/CommonException.java) 定义了 Flea I18N 下的通用异常，由子类传入具体的国际化资源枚举类型

```java
/**
 * Flea I18N 通用异常，由子类传入具体的国际化资源枚举类型
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
public abstract class CommonException extends Exception {

    private static final long serialVersionUID = 1746312829236028651L;

    private String key;                     // 国际化资源数据关键字

    private Locale locale;                  // 国际化区域标识

    private FleaI18nResEnum i18nResEnum;    // 国际化资源类型

    public CommonException(String mKey, FleaI18nResEnum mI18nResEnum) {
        // 使用服务器当前默认的国际化区域设置
        this(mKey, mI18nResEnum, FleaFrameManager.getManager().getLocale());
    }

    public CommonException(String mKey, FleaI18nResEnum mI18nResEnum, String... mValues) {
        // 使用服务器当前默认的国际化区域设置
        this(mKey, mI18nResEnum, FleaFrameManager.getManager().getLocale(), mValues);
    }

    public CommonException(String mKey, FleaI18nResEnum mI18nResEnum, Locale mLocale) {
        // 使用指定的国际化区域设置
        this(mKey, mI18nResEnum, mLocale, new String[]{});
    }

    public CommonException(String mKey, FleaI18nResEnum mI18nResEnum, Locale mLocale, String... mValues) {
        // 使用指定的国际化区域设置
        super(convert(mKey, mValues, mI18nResEnum, mLocale));
        key = mKey;
        locale = mLocale;
        i18nResEnum = mI18nResEnum;
    }

    public CommonException(String mKey, FleaI18nResEnum mI18nResEnum, Throwable cause) {
        // 使用服务器当前默认的国际化区域设置
        this(mKey, mI18nResEnum, FleaFrameManager.getManager().getLocale(), cause);
    }

    public CommonException(String mKey, FleaI18nResEnum mI18nResEnum, Throwable cause, String... mValues) {
        // 使用服务器当前默认的国际化区域设置
        this(mKey, mI18nResEnum, FleaFrameManager.getManager().getLocale(), cause, mValues);
    }

    public CommonException(String mKey, FleaI18nResEnum mI18nResEnum, Locale mLocale, Throwable cause) {
        // 使用指定的国际化区域设置
        this(mKey, mI18nResEnum, mLocale, cause, new String[]{});
    }

    public CommonException(String mKey, FleaI18nResEnum mI18nResEnum, Locale mLocale, Throwable cause, String... mValues) {
        // 使用指定的国际化区域设置
        super(convert(mKey, mValues, mI18nResEnum, mLocale), cause);
        key = mKey;
        locale = mLocale;
        i18nResEnum = mI18nResEnum;
    }

    private static String convert(String key, String[] values, FleaI18nResEnum i18nResEnum, Locale locale) {
        if (ObjectUtils.isEmpty(locale)) {
            locale = FleaFrameManager.getManager().getLocale(); // 使用当前线程默认的国际化区域设置
        }
        if (ObjectUtils.isEmpty(i18nResEnum)) {
            i18nResEnum = FleaI18nResEnum.ERROR; // 默认使用 国际化资源名为 error
        }
        if (ArrayUtils.isNotEmpty(values)) {
            return FleaI18nHelper.i18n(key, values, i18nResEnum.getResName(), locale);
        } else {
            return FleaI18nHelper.i18n(key, i18nResEnum.getResName(), locale);
        }
    }

    public String getKey() {
        return key;
    }

    public Locale getLocale() {
        return locale;
    }

    public FleaI18nResEnum getI18nResEnum() {
        return i18nResEnum;
    }
}

```

## 4.2 定义业务逻辑层异常类

[ServiceException](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-common/src/main/java/com/huazie/fleaframework/db/common/exception/ServiceException.java) 定义了业务逻辑层抛出的异常，对应的国际化资源名为【error】

```java
/**
 * 业务逻辑层异常类，定义了业务逻辑层抛出的异常，
 * 其对应的国际化资源名为【error】
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
public class ServiceException extends CommonException {

    public ServiceException(String key) {
        super(key, FleaI18nResEnum.ERROR);
    }

    public ServiceException(String key, String... values) {
        super(key, FleaI18nResEnum.ERROR, values);
    }

    public ServiceException(String key, Throwable cause) {
        super(key, FleaI18nResEnum.ERROR, cause);
    }

    public ServiceException(String key, Throwable cause, String... values) {
        super(key, FleaI18nResEnum.ERROR, cause, values);
    }

}

```

# 总结

好了，Flea框架下的本地国际化实现已经介绍完毕，欢迎大家使用！
