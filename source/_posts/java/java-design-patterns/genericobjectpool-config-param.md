---
title: 对象池 GenericObjectPool 配置参数详解
date: 2019-09-25 17:41:42
updated: 2024-03-19 17:21:35
categories:
  - [开发语言-Java,Java设计模式]
tags:
  - Java
  - 设计模式
  - 对象池
  - GenericObjectPool
---

![](/images/java-design-patterns-logo.png)

# 引言

使用 **GenericObjectPool** 之前，我们有必要了解一下 **GenericObjectPoolConfig**，下面将详细说明一下其相关的配置参数。



# 1. 对象池涉及依赖

```xml
    <dependency>
        <groupId>org.apache.commons</groupId>
        <artifactId>commons-pool2</artifactId>
        <version>2.4.3</version>
    </dependency>
```


# 2. 父类BaseObjectPoolConfig配置参数

`BaseObjectPoolConfig` 提供了子类共享的通用属性的实现。将使用公共常量定义的默认值创建此类的新实例。 特别注意，该类不是线程安全的。

```java
public abstract class BaseObjectPoolConfig extends BaseObject implements Cloneable {
    private boolean lifo = DEFAULT_LIFO;
    
    private boolean fairness = DEFAULT_FAIRNESS;
    
    private long maxWaitMillis = DEFAULT_MAX_WAIT_MILLIS;
    
    private long minEvictableIdleTimeMillis = DEFAULT_MIN_EVICTABLE_IDLE_TIME_MILLIS;
    
    private long evictorShutdownTimeoutMillis = DEFAULT_EVICTOR_SHUTDOWN_TIMEOUT_MILLIS;
    
    private long softMinEvictableIdleTimeMillis = DEFAULT_SOFT_MIN_EVICTABLE_IDLE_TIME_MILLIS;
    
    private int numTestsPerEvictionRun = DEFAULT_NUM_TESTS_PER_EVICTION_RUN;
    
    private String evictionPolicyClassName = DEFAULT_EVICTION_POLICY_CLASS_NAME;
    
    private boolean testOnCreate = DEFAULT_TEST_ON_CREATE;
    
    private boolean testOnBorrow = DEFAULT_TEST_ON_BORROW;
    
    private boolean testOnReturn = DEFAULT_TEST_ON_RETURN;
    
    private boolean testWhileIdle = DEFAULT_TEST_WHILE_IDLE;
    
    private long timeBetweenEvictionRunsMillis = DEFAULT_TIME_BETWEEN_EVICTION_RUNS_MILLIS;

    private boolean blockWhenExhausted = DEFAULT_BLOCK_WHEN_EXHAUSTED;

    private boolean jmxEnabled = DEFAULT_JMX_ENABLE;

    private String jmxNamePrefix = DEFAULT_JMX_NAME_PREFIX;

    private String jmxNameBase = DEFAULT_JMX_NAME_BASE;
}
```


## 2.1 lifo

提供了后进先出(`LIFO`)与先进先出(`FIFO`)两种行为模式的池；
默认 `DEFAULT_LIFO = true`， 当池中有空闲可用的对象时，调用`borrowObject` 方法会返回最近（后进）的实例。

`org.apache.commons.pool2.impl.GenericObjectPool`

```java
if (getLifo()) {
    idleObjects.addFirst(p);
} else {
    idleObjects.addLast(p);
}
```

## 2.2 fairness

当从池中获取资源或者将资源还回池中时,是否使用`java.util.concurrent.locks.ReentrantLock.ReentrantLock` 的公平锁机制。
默认 `DEFAULT_FAIRNESS = false`

`org.apache.commons.pool2.impl.GenericObjectPool`

```java
    idleObjects = new LinkedBlockingDeque<PooledObject<T>>(config.getFairness());
```

## 2.3 maxWaitMillis

当连接池资源用尽后，调用者获取连接时的最大等待时间（单位 ：毫秒）；
默认值 `DEFAULT_MAX_WAIT_MILLIS = -1L`， 永不超时。

`org.apache.commons.pool2.impl.GenericObjectPool`

```java
    @Override
    public T borrowObject() throws Exception {
        return borrowObject(getMaxWaitMillis());
    }
```

## 2.4 minEvictableIdleTimeMillis

连接的最小空闲时间，达到此值后该空闲连接可能会被移除（还需看是否已达最大空闲连接数）；
默认值 `DEFAULT_MIN_EVICTABLE_IDLE_TIME_MILLIS = 1000L * 60L * 30L`

`org.apache.commons.pool2.impl.GenericObjectPool`

```java
    final EvictionConfig evictionConfig = new EvictionConfig(
                        getMinEvictableIdleTimeMillis(),
                        getSoftMinEvictableIdleTimeMillis(),
                        getMinIdle());
```

## 2.5 evictorShutdownTimeoutMillis

关闭驱逐线程的超时时间；
默认值 `DEFAULT_EVICTOR_SHUTDOWN_TIMEOUT_MILLIS = 10L * 1000L`

`org.apache.commons.pool2.impl.BaseGenericObjectPool`

```java
    final void startEvictor(final long delay) {
        synchronized (evictionLock) {
            if (null != evictor) {
                EvictionTimer.cancel(evictor, evictorShutdownTimeoutMillis, TimeUnit.MILLISECONDS);
                evictor = null;
                evictionIterator = null;
            }
            if (delay > 0) {
                evictor = new Evictor();
                EvictionTimer.schedule(evictor, delay, delay);
            }
        }
    }
```

## 2.6 softMinEvictableIdleTimeMillis

连接空闲的最小时间，达到此值后空闲链接将会被移除，且保留 `minIdle` 个空闲连接数；
默认值 `DEFAULT_SOFT_MIN_EVICTABLE_IDLE_TIME_MILLIS = -1`

## 2.7 numTestsPerEvictionRun

检测空闲对象线程每次运行时检测的空闲对象的数量；

*   如果 `numTestsPerEvictionRun>=0`, 则取 `numTestsPerEvictionRun` 和池内的连接数 的较小值 作为每次检测的连接数；
*   如果 `numTestsPerEvictionRun<0`，则每次检查的连接数是检查时池内连接的总数除以这个值的绝对值再向上取整的结果；

默认值 `DEFAULT_NUM_TESTS_PER_EVICTION_RUN = 3`

`org.apache.commons.pool2.impl.GenericObjectPool`

```java
    private int getNumTests() {
        final int numTestsPerEvictionRun = getNumTestsPerEvictionRun();
        if (numTestsPerEvictionRun >= 0) {
            return Math.min(numTestsPerEvictionRun, idleObjects.size());
        }
        return (int) (Math.ceil(idleObjects.size() / Math.abs((double) numTestsPerEvictionRun)));
    }
```

## 2.8 evictionPolicyClassName

驱逐策略的类名；
默认值 `DEFAULT_EVICTION_POLICY_CLASS_NAME = "org.apache.commons.pool2.impl.DefaultEvictionPolicy"`

`org.apache.commons.pool2.impl.GenericObjectPool`

```java
    public final void setEvictionPolicyClassName(final String evictionPolicyClassName) {
        try {
            Class<?> clazz;
            try {
                clazz = Class.forName(evictionPolicyClassName, true, Thread.currentThread().getContextClassLoader());
            } catch (final ClassNotFoundException e) {
                clazz = Class.forName(evictionPolicyClassName);
            }
            final Object policy = clazz.newInstance();
            if (policy instanceof EvictionPolicy<?>) {
                @SuppressWarnings("unchecked") // safe, because we just checked the class
                final EvictionPolicy<T> evicPolicy = (EvictionPolicy<T>) policy;
                this.evictionPolicy = evicPolicy;
            } else {
                throw new IllegalArgumentException("[" + evictionPolicyClassName + "] does not implement EvictionPolicy");
            }
        } catch (final ClassNotFoundException e) {
            throw new IllegalArgumentException("Unable to create EvictionPolicy instance of type " + evictionPolicyClassName, e);
        } catch (final InstantiationException e) {
            throw new IllegalArgumentException("Unable to create EvictionPolicy instance of type " + evictionPolicyClassName, e);
        } catch (final IllegalAccessException e) {
            throw new IllegalArgumentException("Unable to create EvictionPolicy instance of type " + evictionPolicyClassName, e);
        }
    }
```

## 2.9 testOnCreate

在创建对象时检测对象是否有效(`true : 是`) , 配置 `true` 会降低性能；
默认值 `DEFAULT_TEST_ON_CREATE = false`。

`org.apache.commons.pool2.impl.GenericObjectPool##borrowObject(final long borrowMaxWaitMillis)`

```java
    PooledObject<T> p = null;
    // ...省略
    // 配置true，新创建对象都会检测对象是否有效 【 create && getTestOnCreate() 】
    if (p != null && (getTestOnBorrow() || create && getTestOnCreate())) {
        boolean validate = false;
        Throwable validationThrowable = null;
        try {
            validate = factory.validateObject(p);
        } catch (final Throwable t) {
            PoolUtils.checkRethrow(t);
            validationThrowable = t;
        }
        if (!validate) {
            try {
                destroy(p);
                destroyedByBorrowValidationCount.incrementAndGet();
            } catch (final Exception e) {
                // Ignore - validation failure is more important
            }
            p = null;
            if (create) {
                final NoSuchElementException nsee = new NoSuchElementException("Unable to validate object");
                nsee.initCause(validationThrowable);
                throw nsee;
            }
        }
    }
```

## 2.10 testOnBorrow

在从对象池获取对象时是否检测对象有效(`true : 是)` , 配置 `true` 会降低性能；
默认值 `DEFAULT_TEST_ON_BORROW = false`

`org.apache.commons.pool2.impl.GenericObjectPool##borrowObject(final long borrowMaxWaitMillis)`

```java
    // 配置true，从对象池获取对象，总是要检测对象是否有效 【 getTestOnBorrow() 】
    if (p != null && (getTestOnBorrow() || create && getTestOnCreate())) {
        // ...省略
    }
```

## 2.11 testOnReturn

在向对象池中归还对象时是否检测对象有效(`true : 是`) , 配置 `true` 会降低性能；
默认值 `DEFAULT_TEST_ON_RETURN = false`

`org.apache.commons.pool2.impl.GenericObjectPool##returnObject(final T obj)`

```java
    if (getTestOnReturn()) {
        if (!factory.validateObject(p)) {
            try {
                destroy(p);
            } catch (final Exception e) {
                swallowException(e);
            }
            try {
                ensureIdle(1, false);
            } catch (final Exception e) {
                swallowException(e);
            }
            updateStatsReturn(activeTime);
            return;
        }
    }
```

## 2.12 testWhileIdle

在检测空闲对象线程检测到对象不需要移除时，是否检测对象的有效性。建议配置为 `true`，不影响性能，并且保证安全性；
默认值 `DEFAULT_TEST_WHILE_IDLE = false`

`org.apache.commons.pool2.impl.GenericObjectPool##evict()`

```java
    final boolean testWhileIdle = getTestWhileIdle();
    // .... 代码省略
    // 配置为true， 检测空闲对象线程检测到对象不需要移除时，检测对象的有效性
    if (testWhileIdle) {
        boolean active = false;
        try {
            factory.activateObject(underTest);
            active = true;
        } catch (final Exception e) {
            destroy(underTest);
            destroyedByEvictorCount.incrementAndGet();
        }
        if (active) {
            if (!factory.validateObject(underTest)) {
                destroy(underTest);
                destroyedByEvictorCount.incrementAndGet();
            } else {
                try {
                    factory.passivateObject(underTest);
                } catch (final Exception e) {
                    destroy(underTest);
                    destroyedByEvictorCount.incrementAndGet();
                }
            }
        }
    }
```

## 2.13 timeBetweenEvictionRunsMillis

空闲连接检测的周期（单位毫秒）；如果为负值，表示不运行检测线程；
默认值 `DEFAULT_TIME_BETWEEN_EVICTION_RUNS_MILLIS = -1L`

`org.apache.commons.pool2.impl.GenericObjectPool`

```java
    public GenericObjectPool(final PooledObjectFactory<T> factory, final GenericObjectPoolConfig config) {
           super(config, ONAME_BASE, config.getJmxNamePrefix());
        if (factory == null) {
            jmxUnregister(); // tidy up
            throw new IllegalArgumentException("factory may not be null");
        }
        this.factory = factory;
        
        idleObjects = new LinkedBlockingDeque<PooledObject<T>>(config.getFairness());

        setConfig(config);
        // 启动空闲连接检测线程
        startEvictor(getTimeBetweenEvictionRunsMillis());
    }
```

## 2.14 blockWhenExhausted

当对象池没有空闲对象时，新的获取对象的请求是否阻塞（`true` 阻塞，`maxWaitMillis` 才生效； `false` 连接池没有资源立马抛异常）
默认值 `DEFAULT_BLOCK_WHEN_EXHAUSTED = true`

`org.apache.commons.pool2.impl.GenericObjectPool##borrowObject(final long borrowMaxWaitMillis)`

```java
    final boolean blockWhenExhausted = getBlockWhenExhausted();
    // ... 省略
    if (blockWhenExhausted) {
        if (p == null) {
            if (borrowMaxWaitMillis < 0) {
                p = idleObjects.takeFirst();
            } else {
                p = idleObjects.pollFirst(borrowMaxWaitMillis,
                        TimeUnit.MILLISECONDS);
            }
        }
        if (p == null) {
            throw new NoSuchElementException(
                    "Timeout waiting for idle object");
        }
    } 
```

## 2.15 jmxEnabled

是否注册 `JMX`
默认值 `DEFAULT_JMX_ENABLE = true`

`org.apache.commons.pool2.impl.BaseGenericObjectPool`

```java
    public BaseGenericObjectPool(final BaseObjectPoolConfig config, final String jmxNameBase, final String jmxNamePrefix) {
        if (config.getJmxEnabled()) {
            this.oname = jmxRegister(config, jmxNameBase, jmxNamePrefix);
        } else {
            this.oname = null;
        }

        // Populate the creation stack trace
        this.creationStackTrace = getStackTrace(new Exception());

        // save the current TCCL (if any) to be used later by the evictor Thread
        final ClassLoader cl = Thread.currentThread().getContextClassLoader();
        if (cl == null) {
            factoryClassLoader = null;
        } else {
            factoryClassLoader = new WeakReference<ClassLoader>(cl);
        }

        fairness = config.getFairness();
    }
```

## 2.16 jmxNamePrefix

`JMX` 前缀
默认值 `DEFAULT_JMX_NAME_PREFIX = "pool"`

`org.apache.commons.pool2.impl.GenericObjectPool`

```java
    // JMX specific attributes
    private static final String ONAME_BASE = "org.apache.commons.pool2:type=GenericObjectPool,name=";

    public GenericObjectPool(final PooledObjectFactory<T> factory, final GenericObjectPoolConfig config) {
        // 参见上述 jmxEnabled 部分
        super(config, ONAME_BASE, config.getJmxNamePrefix());
        // .....
    }
```

## 2.17 jmxNameBase

使用 `base + jmxNamePrefix + i` 来生成 `ObjectName`
默认值 `DEFAULT_JMX_NAME_BASE = null`，`GenericObjectPool` 构造方法使用 `ONAME_BASE` 初始化。

```java
    private ObjectName jmxRegister(final BaseObjectPoolConfig config,final String jmxNameBase, String jmxNamePrefix) {
        ObjectName objectName = null;
        final MBeanServer mbs = ManagementFactory.getPlatformMBeanServer();
        int i = 1;
        boolean registered = false;
        String base = config.getJmxNameBase();
        if (base == null) {
            base = jmxNameBase;
        }
        while (!registered) {
            try {
                ObjectName objName;
                if (i == 1) {
                    objName = new ObjectName(base + jmxNamePrefix);
                } else {
                    objName = new ObjectName(base + jmxNamePrefix + i);
                }
                mbs.registerMBean(this, objName);
                objectName = objName;
                registered = true;
            } catch (final MalformedObjectNameException e) {
                if (BaseObjectPoolConfig.DEFAULT_JMX_NAME_PREFIX.equals(jmxNamePrefix) && jmxNameBase.equals(base)) {
                    // 如果走到这步，就跳过jmx注册，应该不会发生
                    registered = true;
                } else {
                    // 前者使用的名称不是默认配置，则使用默认配置替代
                    jmxNamePrefix =
                            BaseObjectPoolConfig.DEFAULT_JMX_NAME_PREFIX;
                    base = jmxNameBase;
                }
            } catch (final InstanceAlreadyExistsException e) {
                // 增加一个索引，再试一次
                i++;
            } catch (final MBeanRegistrationException e) {
                // 如果走到这步，就跳过jmx注册，应该不会发生
                registered = true;
            } catch (final NotCompliantMBeanException e) {
                // 如果走到这步，就跳过jmx注册，应该不会发生
                registered = true;
            }
        }
        return objectName;
    }        
```

# 3. 子类GenericObjectPoolConfig配置参数

`GenericObjectPoolConfig` 是一个简单的配置类，封装了用于 GenericObjectPool 的配置。注意，该类也不是线程安全的，它仅用于在创建池时提供属性。

```java

public class GenericObjectPoolConfig extends BaseObjectPoolConfig {
    
    private int maxTotal = DEFAULT_MAX_TOTAL;
    
    private int maxIdle = DEFAULT_MAX_IDLE;
    
    private int minIdle = DEFAULT_MIN_IDLE;
}
```


## 3.1 maxTotal

**最大连接数**，默认值 `DEFAULT_MAX_TOTAL = 8`

## 3.2 maxIdle

**最大空闲连接数**， 默认值 `DEFAULT_MAX_IDLE = 8`

## 3.3 minIdle

**最小空闲连接数**， 默认值 `DEFAULT_MIN_IDLE = 0`


# 总结

了解这些配置参数对于正确设置和管理 `GenericObjectPool` 至关重要，本篇的介绍基于**commons-pool2-2.4.3**，其他版本可能有出入，请自行查看。