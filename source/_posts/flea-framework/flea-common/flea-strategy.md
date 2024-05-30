---
title: flea-common使用之通用策略模式实现
date: 2021-08-24 16:22:00 
updated: 2024-02-20 18:06:30
categories:
  - [开发框架-Flea,flea-common]
tags:
  - flea-framework
  - flea-common
  - Flea Strategy
  - 通用策略模式
---

![](/images/flea-logo.png)


# 1. 概述
策略模式（**Strategy Pattern**）作为一种软件设计模式，用来实现对象的某个行为，该行为在不同的场景中拥有不同的实现逻辑。它定义了一组算法，同时将这些算法封装起来，并使它们之间可以互换。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

常用的策略模式有如下三个角色：

- **抽象策略角色 --- Strategy**
抽象策略类，通常为一个接口，其中定义了某个策略行为【即策略算法标识】。

- **具体策略角色 --- ConcreteStrategy**
具体策略类，实现抽象策略中的策略行为；每一个具体策略类即代表一种策略算法。

- **上下文角色 --- Context**
上下文类，起承上启下封装作用，屏蔽高层模块对策略、算法的直接访问，封装可能存在的变化。

本篇在上述常用的策略模式基础上，结合门面模式和调整后的策略上下文，构建了一套通用策略模式实现。

下面我们用这套通用策略模式来模拟一下各种动物的喊叫行为：

![](flea-strategy.png)

# 2. 参考
[flea-common使用之通用策略模式实现  源代码](https://github.com/Huazie/flea-framework/tree/dev/flea-common)

# 3. 实现
## 3.1 定义Flea策略接口类
[IFleaStrategy](https://github.com/Huazie/flea-framework/blob/dev/flea-common/src/main/java/com/huazie/fleaframework/common/strategy/IFleaStrategy.java) 定义了通用的策略行为，**T** 类型表示 **Flea** 策略执行结果对应的类型，**P** 类型表示 **Flea** 策略上下文参数。 具体代码如下：

```java
/**
 * Flea策略接口，定义统一的策略执行方法。
 *
 * @param <T> Flea策略执行结果对应的类型
 * @param <P> Flea策略上下文参数
 * @author huazie
 * @version 1.1.0
 * @since 1.1.0
 */
public interface IFleaStrategy<T, P> {

    /**
     * 策略执行
     *
     * @param contextParam Flea策略上下文参数
     * @return 策略执行结果对应的类型
     * @throws FleaStrategyException Flea策略异常
     * @since 1.1.0
     */
    T execute(final P contextParam) throws FleaStrategyException;
}
```
## 3.2 定义狗喊叫声策略类
 [DogVoiceStrategy](https://github.com/Huazie/flea-framework/blob/dev/flea-common/src/test/java/com/huazie/fleaframework/common/strategy/DogVoiceStrategy.java) 用于定义狗喊叫声策略，返回 `"阿狗【" + name + "】正在喊叫着；汪汪汪"`

```java
/**
 * 狗喊叫声策略
 *
 * @author huazie
 * @version 1.1.0
 * @since 1.1.0
 */
public class DogVoiceStrategy implements IFleaStrategy<String, String> {

    @Override
    public String execute(String name) throws FleaStrategyException {
        return "阿狗【" + name + "】正在喊叫着；汪汪汪";
    }

}
```
## 3.3 定义猫喊叫声策略类
[CatVoiceStrategy](https://github.com/Huazie/flea-framework/blob/dev/flea-common/src/test/java/com/huazie/fleaframework/common/strategy/CatVoiceStrategy.java) 用于定义猫喊叫声策略，返回 `"阿猫【" + name + "】正在喊叫着；喵喵喵"`

```java
/**
 * 猫喊叫声策略
 *
 * @author huazie
 * @version 1.1.0
 * @since 1.1.0
 */
public class CatVoiceStrategy implements IFleaStrategy<String, String> {

    @Override
    public String execute(String name) throws FleaStrategyException {
        return "阿猫【" + name + "】正在喊叫着；喵喵喵";
    }

}
```
## 3.4 定义鸭喊叫声策略类
[DuckVoiceStrategy](https://github.com/Huazie/flea-framework/blob/dev/flea-common/src/test/java/com/huazie/fleaframework/common/strategy/DuckVoiceStrategy.java) 定义鸭喊叫声策略，返回 `"阿鸭【" + name + "】正在喊叫着；嘎嘎嘎"`

```java
/**
 * 鸭喊叫声策略
 *
 * @author huazie
 * @version 1.1.0
 * @since 1.1.0
 */
public class DuckVoiceStrategy implements IFleaStrategy<String, String> {

    @Override
    public String execute(String name) throws FleaStrategyException {
        return "阿鸭【" + name + "】正在喊叫着；嘎嘎嘎";
    }

}
```

## 3.5 定义策略上下文接口类 
[IFleaStrategyContext](https://github.com/Huazie/flea-framework/blob/dev/flea-common/src/main/java/com/huazie/fleaframework/common/strategy/IFleaStrategyContext.java) 策略上下文接口，定义统一的策略上下文调用方法，同时可通过 **getContext** 获取上下文参数，**setContext** 设置上下文参数。
```java
/**
 * Flea策略上下文接口，定义统一的策略上下文调用方法。
 *
 * @param <T> Flea策略执行结果对应的类型
 * @param <P> Flea策略上下文参数
 * @author huazie
 * @version 1.1.0
 * @since 1.1.0
 */
public interface IFleaStrategyContext<T, P> {

    /**
     * 策略上下文调用
     *
     * @param strategy 策略名称
     * @return 策略执行结果对应的类型
     * @throws FleaStrategyException Flea策略异常
     * @since 1.1.0
     */
    T invoke(String strategy) throws FleaStrategyException;

    /**
     * 设置策略上下文参数
     *
     * @param contextParam 上下文参数对象
     * @since 1.1.0
     */
    void setContext(P contextParam);

    /**
     * 获取策略上下文参数
     *
     * @return 策略上下文参数
     * @since 1.1.0
     */
    P getContext();
}
```

## 3.6 定义Flea抽象策略上下文类
[FleaStrategyContext](https://github.com/Huazie/flea-framework/blob/dev/flea-common/src/main/java/com/huazie/fleaframework/common/strategy/FleaStrategyContext.java) **Flea** 抽象策略上下文类，封装了策略执行的逻辑，对外屏蔽高层模块对策略的直接访问。抽象方法 **init** 用于初始化 **Flea** 策略实现 **Map**，该方法由**Flea** 策略抽象上下文的子类实现，并在策略上下文子类实例化时，调用该方法完成具体初始化的工作。

```java
/**
 * Flea抽象策略上下文，封装了公共的策略执行逻辑，
 * 其中Flea策略Map由其子类进行初始化，键为策略名，
 * 值为具体的Flea策略实例。
 *
 * @param <T> Flea策略执行结果对应的类型
 * @param <P> Flea策略上下文参数
 * @author huazie
 * @version 1.1.0
 * @since 1.1.0
 */
public abstract class FleaStrategyContext<T, P> implements IFleaStrategyContext<T, P> {

    private Map<String, IFleaStrategy<T, P>> fleaStrategyMap; // Flea策略Map

    private P contextParam; // Flea策略上下文参数

    /**
     * 初始化策略上下文
     *
     * @since 1.1.0
     */
    public FleaStrategyContext() {
        fleaStrategyMap = init();
    }

    /**
     * 初始化策略上下文
     *
     * @param contextParam Flea策略上下文参数
     * @since 1.1.0
     */
    public FleaStrategyContext(P contextParam) {
        this();
        this.contextParam = contextParam;
    }

    /**
     * 初始化Flea策略Map
     *
     * @return Flea策略Map
     * @since 1.1.0
     */
    protected abstract Map<String, IFleaStrategy<T, P>> init();

    @Override
    public T invoke(String strategy) throws FleaStrategyException {
        if (ObjectUtils.isEmpty(fleaStrategyMap)) {
            throw new FleaStrategyException("The Strategy Map is not initialized!");
        }
        IFleaStrategy<T, P> fleaStrategy = fleaStrategyMap.get(strategy);
        if (ObjectUtils.isEmpty(fleaStrategy)) {
            throw new FleaStrategyNotFoundException("The Strategy [name =\"" + strategy + "\"] is not found!");
        }
        return fleaStrategy.execute(contextParam);
    }

    @Override
    public void setContext(P contextParam) {
        this.contextParam = contextParam;
    }

    @Override
    public P getContext() {
        return contextParam;
    }
}
```

## 3.7 定义动物喊叫声策略上下文类
[AnimalVoiceContext](https://github.com/Huazie/flea-framework/blob/dev/flea-common/src/test/java/com/huazie/fleaframework/common/strategy/AnimalVoiceContext.java) 动物喊叫声策略上下文，继承 **Flea** 抽象策略上下文，实现 **init** 方法，用于初始化 **Flea** 策略实现 **Map**，其中 **key** 为 策略名，**value** 为 具体的动物喊叫声策略实现类；**Collections.unmodifiableMap** 用于返回一个 **只读** 的 **Map**。

```java
/**
 * 动物喊叫声策略上下文
 *
 * @author huazie
 * @version 1.1.0
 * @since 1.1.0
 */
public class AnimalVoiceContext extends FleaStrategyContext<String, String> {

    private static Map<String, IFleaStrategy<String, String>> fleaStrategyMap;

    static {
        fleaStrategyMap = new HashMap<>();
        fleaStrategyMap.put("dog", new DogVoiceStrategy());
        fleaStrategyMap.put("cat", new CatVoiceStrategy());
        fleaStrategyMap.put("duck", new DuckVoiceStrategy());
        fleaStrategyMap = Collections.unmodifiableMap(fleaStrategyMap);
    }

    public AnimalVoiceContext() {
    }

    public AnimalVoiceContext(String contextParam) {
        super(contextParam);
    }

    @Override
    protected Map<String, IFleaStrategy<String, String>> init() {
        return fleaStrategyMap;
    }
}
```
## 3.8 定义Flea策略门面
[FleaStrategyFacade](https://github.com/Huazie/flea-framework/blob/dev/flea-common/src/main/java/com/huazie/fleaframework/common/strategy/FleaStrategyFacade.java) 定义 Flea 策略调用的统一入口

```java
/**
 * Flea策略门面，定义策略调用的统一入口。
 *
 * @author huazie
 * @version 1.1.0
 * @since 1.1.0
 */
public class FleaStrategyFacade {

    private FleaStrategyFacade() {
    }

    /**
     * 策略门面调用方法
     *
     * @param strategy            策略名
     * @param fleaStrategyContext 策略上下文
     * @param <T>                 Flea策略执行结果对应的类型
     * @param <P>                 Flea策略上下文参数
     * @return Flea策略执行结果对应的类型
     * @throws FleaStrategyException Flea策略异常
     * @since 1.1.0
     */
    public static <T, P> T invoke(String strategy, IFleaStrategyContext<T, P> fleaStrategyContext) throws FleaStrategyException {
        return fleaStrategyContext.invoke(strategy);
    }
}
```

# 4. 测试
单元自测类可查看 [FleaStrategyTest](https://github.com/Huazie/flea-framework/blob/dev/flea-common/src/test/java/com/huazie/fleaframework/common/strategy/FleaStrategyTest.java)。
```java
/**
 * @author huazie
 * @version 1.1.0
 * @since 1.1.0
 */
public class FleaStrategyTest {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(FleaStrategyTest.class);

    @Test
    public void testStrategy() {

        AnimalVoiceContext context = new AnimalVoiceContext("旺财");
        LOGGER.debug(FleaStrategyFacade.invoke("dog", context));

        context.setContext("Tom");
        LOGGER.debug(FleaStrategyFacade.invoke("cat", context));

        AnimalVoiceContext context1 = new AnimalVoiceContext();
        context1.setContext("Donald");
        LOGGER.debug(FleaStrategyFacade.invoke("duck", context1));

    }
}

```
单元测试类运行结果如下：

![](result.png)

# 总结
好了，通用策略模式模式实现--`Flea Strategy` 已讲解完毕，欢迎大家使用 ！
