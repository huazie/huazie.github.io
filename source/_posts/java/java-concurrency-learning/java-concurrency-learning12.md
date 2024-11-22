---
title: Java并发编程学习12-任务取消和线程中断
date: 2022-11-11 16:57:44
updated: 2024-03-19 17:21:35
categories:
  - [开发语言-Java,Java并发编程]
tags:
  - Java
  - 取消策略
  - 线程中断
  - 中断策略
  - 响应中断
---



![](/images/java-concurrency-logo.png)

# 引言

《任务取消》由于篇幅较多，拆分了两篇来介绍各种实现取消和中断的机制，以及如何编写任务和服务，使它们能对取消请求做出响应。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

**如何理解任务是可取消的 ？**

如果外部代码能在某个任务正常完成之前将其置入 “完成” 状态，那么这个任务就被认为是可取消的。

大多数任务，我们都希望让它们运行直到结束，或者让它们自行停止。但是也有**很多原因**，导致我们需要取消这些任务，如下所示：

- **用户请求取消**。用户点击图形界面程序中的 “取消” 按钮，或者通过管理接口来发出取消请求，例如 **JMX**（**Java Management Extensions**，即 **Java** 管理扩展）。
- **有时间限制的操作**。某个应用程序需要在有限时间内搜索问题空间，并在这个时间内选择最佳的解决方案。当计时器超时时，需要取消所有正在搜索的任务。
- **应用程序事件**。应用程序对某个问题空间进行分解并搜索，从而使不同的任务可以搜索问题空间中的不同区域。当其中一个任务找到了解决方案时，所有其他仍在搜索的任务都将被取消。
- **错误**。网页爬虫程序搜索相关的页面，并将页面或摘要数据保存到硬盘。当一个爬虫任务发生错误时（例如，磁盘空间已满），那么所有搜索任务都会取消，此时可能会记录它们的当前状态，以便稍后重新启动。
- **关闭**。当一个程序或服务关闭时，必须对正在处理和等待处理的工作执行某种操作。在平缓的关闭过程中，当前正在执行的任务将继续执行直到完成，而在立即关闭过程中，当前的任务则可能取消。

# 主要内容
## 1. 取消策略
**当我们需要取消任务时，该怎么操作呢？**

在 **Java** 中没有一种安全的抢占式方式来停止线程，因此也就没有安全的抢占式方法来停止任务。只有一种 **协作机制**，使请求取消的任务和代码都遵循一种协商好的协议。

下面我们来看一下如下示例【使用 **volatile** 类型的域来保存取消状态】：

```java
@ThreadSafe
public class PrimeGenerator implements Runnable {
    @GuardedBy("this")
    private final List<BigInteger> primes = new ArrayList<BigInteger>();
    // 为了使这个过程能可靠地工作，标志 cancelled 必须为 volatile 类型
    private volatile boolean cancelled;
    
    public void run() {
        BigInteger p = BigInteger.ONE;
        while(!cancelled) {
            LOGGER.debug("before = {}", p);
            p = p.nextProbablePrime();
            LOGGER.debug("prime = {}", p);
            synchronized (this) {
                primes.add(p);
            }
        }
    }
    
    public void cancel() {
        cancelled = true;
    }
    
    public synchronized List<BigInteger> get() {
        return new ArrayList<BigInteger>(primes);
    }
}
```

如上示例 **PrimeGenerator** 将会持续地枚举素数，直到它被取消。它使用了一种协作机制，**cancel** 方法将设置 **cancelled** 为 **true**，任务在 **搜索下一个素数之前** 会检查这个标志，如果标志为 **true**，则任务将会自行结束。

下面我们再来看一下如下示例【让上面的素数生成器运行 **1** 秒钟后取消】：

```java
public class PrimeTest {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(PrimeTest.class);

    @Test
    public void primeTest() throws InterruptedException {
        LOGGER.debug("Primes = {}", aSecondOfPrimes());
    }

    /**
     * 一个仅运行一秒钟的素数生成器
     */
    public List<BigInteger> aSecondOfPrimes() throws InterruptedException {
        PrimeGenerator generator = new PrimeGenerator();
        new Thread(generator).start();
        try {
            LOGGER.debug("sleep start");
            SECONDS.sleep(1);
        } finally {
            generator.cancel();
        }
        return generator.get();
    }
}
```

实际运行中，上述代码并不会刚好在运行 **1** 秒后停止，因为在请求取消的时刻和 **run** 方法中循环执行下一次检查之间可能存在延迟。**cancel** 方法由 **finally** 块调用，从而确保即使在调用 **sleep** 时被中断也能取消素数生成器的执行。如果 **cancel** 没有被调用，那么搜索素数的线程将永远运行下去。


一个可取消的任务必须拥有取消策略，在该策略中需要详细定义取消操作的三步骤：

- **How**。应用程序的其他代码如何（**How**）请求取消该任务。
- **When**。任务在何时（**When**）检查是否已经请求了取消。
- **What**。在响应取消请求时应该执行哪些（**What**）操作。

上述素数生成器 **PrimeGenerator** 就使用了一种简单的取消策略：

> 客户代码通过 **调用 cancel** 来请求取消，PrimeGenerator 在 **每次搜索素数前** 首先检查是否存在取消请求，如果存在则 **退出**。

## 2. 线程中断

介绍中断之前，我们首先来分析一下，上述素数生成器使用的取消机制目前存在的问题：

- 任务的退出过程仍然需要花费一定的时间。
- 任务中如果调用了一个阻塞的方法（如 **BlockingQueue.put**），它可能永远不会检查取消方法，从而永远不会结束。

下面我们再来看一下如下示例：


```java
public class BrokenPrimeProducer extends Thread {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(BrokenPrimeProducer.class);

    private final BlockingQueue<BigInteger> queue;

    private volatile boolean cancelled = false;

    public BrokenPrimeProducer(BlockingQueue<BigInteger> queue) {
        this.queue = queue;
    }

    @Override
    public void run() {
        try {
            BigInteger p = BigInteger.ONE;
            while (!cancelled) {
                LOGGER.debug("before = {}", p);
                queue.put(p = p.nextProbablePrime());
                LOGGER.debug("prime = {}", p);
            }
        } catch (InterruptedException e) {
            LOGGER.error("InterruptedException");
        }
    }

    public void cancel() {
        cancelled = true;
        LOGGER.debug("cancel");
    }
}

public class PrimeConsumer {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(PrimeConsumer.class);

    private static final int BOUND = 100;

    private int times = 0; // 消费次数

    public void consumePrimes() throws InterruptedException {
        BlockingQueue<BigInteger> primes = new LinkedBlockingQueue<>(BOUND);
        BrokenPrimeProducer producer = new BrokenPrimeProducer(primes);
        producer.start();
        try {
            while (needMorePrimes())
                consume(primes.take());
        } finally {
            producer.cancel();
        }
    }

    private void consume(BigInteger value) {
        times++;
        LOGGER.debug1(new Object() {}, "value = {}", value);
    }

    private boolean needMorePrimes() throws InterruptedException {
        return true;
    }
}
```

如上示例中，生产者线程 **BrokenPrimeProducer** 生成素数，并存放到阻塞队列 **queue** 中【其中 **queue** 由生产者的构造方法传入】。消费者 **PrimeConsumer** 的 **consumePrimes** 方法从生产者中获取素数并处理。

如果生产者的生产速度超过了消费者的处理速度，队列将被填满，其 **put** 方法也会一直阻塞下去。

**我们来思考下遇到上述情况，如果消费者想取消生产者任务，又该怎么办呢？**

消费者可以调用生产者的 **cancel** 方法来设置 **cancelled** 标志，但是因为消费者已经停止从队列中取素数，而阻塞队列 **queue** 的 put 方法将一直保持阻塞状态，导致生产者任务无法从阻塞的 **put** 方法中恢复过来，也就永远不会再检查 **cancelled** 标志，从而无法取消生产者任务。

下面我们来修改下消费者的 **needMorePrimes** 方法， 再用测试类 **PrimeConsumerTest** 来直观地演示下上述情况：

```java
private boolean needMorePrimes() throws InterruptedException {
    Thread.sleep(1000); // 等待1s 再判断
    return times < 5;
}

public class PrimeConsumerTest {

    public static void main(String[] args) throws InterruptedException {
        PrimeConsumer consumer = new PrimeConsumer();
        consumer.consumePrimes();
    }
}
```

运行结果如下：

![](result.png)

从上图可以看到，当消费次数达到 **5** 次后，消费者不再从队列中取素数并打印出来，从代码看它后面直接进入 **finally** 方法，并且调用生产者的 **cancel** 方法准备去取消生产者任务，但是生产者线程的打印日志除了 **cancel** 外，就一直保持 **before** 那，并且看程序也没有结束掉，说明此时生产者在 **put** 方法上一直阻塞着。

**那么我们如何改造生产者，能够在 put 方法阻塞的情况下，支持消费者取消生产者任务呢？**

想要实现这种功能，就不得不提到接下来要讲到的 **线程中断**。

其实在笔者前面的[《阻塞队列》](../../../../../../2022/09/13/java/java-concurrency-learning/java-concurrency-learning7/)博文中，曾简单介绍了阻塞方法与中断方法，大家可以快速去回顾一下。

说到线程中断，就不得不提到 **Thread** 类，下面简单介绍下：

每个线程都有一个 **boolean** 类型的中断状态。

当中断线程时，该线程的中断状态将被设置为 **true**。

线程中还包含了中断线程、查询线程中断状态的方法，如下所示：

- **interrupt** 方法能中断目标线程。
- **isInterrupted** 方法能返回目标线程的中断状态。
- 静态的 **interrupted** 方法将清除当前线程的中断状态，并返回它之前的值，这也是清除中断状态的唯一方法。

**注意：**
- 调用 **interrupt** 方法并不意味着立即停止目标线程正在进行的工作，而只是传递了请求中断的消息。
- 因为静态的 **interrupted** 方法会清除当前线程的中断状态。如果调用它时返回了 **true**，那么除非你想屏蔽这个中断，否则必须对它进行处理----可以抛出 **InterruptedException**，或者通过再次调用 **interrupt** 来恢复中断状态【可以结合[《阻塞队列》“阻塞方法与中断方法”](../../../../../../2022/09/13/java/java-concurrency-learning/java-concurrency-learning7/) 那块的内容进行思考】。

当线程在非阻塞状态下中断时，它的中断状态将被设置，然后根据将被取消的操作来检查中断状态以判断发生了中断。如果不触发 **InterruptedException**，那么中断状态将一直保持，直到明确地清除中断状态。

我们知道 **Java** 类库中的一些阻塞库方法支持中断。例如 **Thread.sleep** 和 **Object.wait** 等，它们都会检查中断状态，并且在发现中断时提前返回。

它们在响应中断时执行的操作包括：

- 清除中断状态
- 抛出 **InterruptedException**，表示阻塞操作由于中断而提前结束。

> 如果任务代码能够响应中断，那么可以使用中断作为取消机制，并且利用  **Java** 类库中提供的中断支持。

通过上面的了解，我们现在可以使用中断来请求取消生产者任务，以解决 **BrokenPrimeProducer** 存在的问题。

下面来看一下改造后的生产者的代码示例：

```java
public class PrimeProducer extends Thread {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(PrimeProducer.class);

    private final BlockingQueue<BigInteger> queue;

    public PrimeProducer(BlockingQueue<BigInteger> queue) {
        this.queue = queue;
    }

    @Override
    public void run() {
        try {
            BigInteger p = BigInteger.ONE;
            while (!Thread.currentThread().isInterrupted()) {
                LOGGER.debug("before = {}", p);
                queue.put(p = p.nextProbablePrime());
                LOGGER.debug("prime = {}", p);
            }
        } catch (InterruptedException e) {
            // 允许线程退出
            LOGGER.error("InterruptedException");
        }
    }

    public void cancel() {
        interrupt(); // 中断线程
        LOGGER.debug("interrupt");
    }
}

public class PrimeConsumer {

    public void consumePrimes() throws InterruptedException {
        // 。。。
        PrimeProducer producer = new PrimeProducer(primes);
        producer.start();
        // 。。。
    }
    
    // 。。。
}
```
上述生产者 **PrimeProducer** 在每次循环中，都有两个位置可以检测出中断：

 1. 阻塞的 **put** 方法调用中。
 2. **while** 循环的判断条件中。
 
当然因为这里调用阻塞的 **put** 方法，**while** 循环条件中显式的检测也可以去掉。但如果加上这段的话，可以使 **PrimeProducer** 对中断具有更高的响应性。


然后我们用自测类 **PrimeConsumerNewTest** 来演示一下使用中断来请求取消：

```java
public class PrimeConsumerNewTest {
    public static void main(String[] args) throws InterruptedException {
        PrimeConsumer consumer = new PrimeConsumer();
        consumer.consumePrimes();
    }
}
```

运行结果如下：

![](result-1.png)

从上图可以看到，当消费次数达到 **5** 次后，消费者不再从队列中取素数并打印出来，从代码看它后面直接进入 **finally** 方法，并且调用生产者的 **cancel** 方法，这里可以看到生产者日志打印了 **interrupt**，发出了中断请求；在发出中断请求之前，我们也从日志中看到生产最后停在了 **before** 处，说明此时生产者在 **put** 方法上阻塞着；而发出中断请求之后，put 方法响应了中断，并抛出了 **InterruptedException**，从日志看就是打印了 **InterruptedException**。

## 3. 中断策略

正如任务中应该包含取消策略一样，线程同样应该包含中断策略。

**那该如何理解中断策略呢？**

中断策略规定了当线程发现中断请求时，应该做哪些工作，哪些工作单元对于中断来说是原子操作，以及以多快的速度来响应中断。

**其中最合理的中断策略是什么呢？**

当线程发现中断请求后，就尽快退出，在必要时进行清理，并通知某个所有者该线程已经退出。

当然，除上外还可以建立其他的中断策略，如 **暂停服务** 或 **重新开始服务**。

我们知道任务不会在其自己拥有的线程中执行，而是在某个服务（例如线程池）拥有的线程中执行。对于非线程所有者的代码来说（例如，对于线程池而言，任何在线程池实现以外的代码），应该小心地保存中断状态，这样拥有线程的代码才能对中断做出响应，即使 “非所有者” 代码也可以做出响应。

这也就是为什么大多数可阻塞的库函数都只是抛出 **InterruptedException** 作为中断响应。它们永远不会在某个由自己拥有的线程中运行，因此它们为任务或库代码实现了最合理的取消策略：**`尽快退出执行流程，并把中断信息传递给调用者，从而使调用栈中的上层代码可以采取进一步的操作`**。

当检查到中断请求时，任务并不需要放弃所有的操作，它可以推迟处理中断请求，并直到某个更合适的时刻。因此就需要记住中断请求，并在完成当前任务后抛出 **InterruptedException** 或者 表示已收到中断请求。

无论任务把中断视为取消，还是其他某个中断响应操作，都应该小心地保存执行线程的中断状态。如果除了将 **InterruptedException** 传递给调用者外还需要下执行其他操作，那么应该在捕获 **InterruptedException** 之后恢复中断状态：
```java
    Thread.currentThread().interrupt();
```

线程应该只能由其所有者中断，所有者可以将线程的中断策略信息封装到某个合适的取消机制中，例如 **关闭方法**。

> 由于每个线程拥有各自的中断策略，因此除非你知道中断对该线程的含义，否则就不应该中断这个线程。

## 4. 响应中断

在笔者前面的[《阻塞队列》](../../../../../../2022/09/13/java/java-concurrency-learning/java-concurrency-learning7/)博文中，当调用可中断的阻塞函数时，例如 **Thread.sleep** 或 **BlockingQueue.put** 等，有两种常见的方法可用于处理 **InterruptedException** ：传递 **InterruptedException** 和 恢复中断。

不过需要注意的是，你不能在 **catch** 块中捕获到 **InterruptedException** 异常却不做任何处理，除非在你的代码中实现了线程的中断策略。

当然细心的小伙伴就会说了，上文中的生产者线程 **PrimeProducer** 不就捕获了 **InterruptedException** 异常却不做任何处理。这里需要解释下，虽然 **PrimeProducer** 屏蔽了中断，但因为它已经知道线程将要结束，并且在调用栈中已经没有上层代码需要知道中断信息。上述只是一类特殊的情况，由于大多数代码并不知道它们将在哪个线程中运行，因此应该保存中断状态。

> 只有实现了线程中断策略的代码才可以屏蔽中断请求。在常规的任务和库代码中都不应该屏蔽中断请求。

对于一些不支持取消但仍可以调用可中断阻塞方法的操作，它们必须在循环中调用这些方法，并在发现中断后重新尝试。在这种情况下，它们应该在本地保存中断状态，并在返回前恢复状态而不是在捕获 **InterruptedException** 时恢复状态，可参考如下示例：

```java
    public Task getNextTask(BlockingQueue<Task> queue) {
        boolean interrupted = false;
        try {
            while (true) {
                try {
                    return queue.take();
                } catch (InterruptedException e) {
                    interrupted = true;
                    // 重新尝试
                }
            }
        } finally {
            if (interrupted) 
                Thread.currentThread().interrupt(); // 恢复中断
        }
    }
```

如果过早地设置中断状态，就可能引起无限循环，因为大多数可中断的阻塞方法都会在入口处检查中断状态，并且当发现该状态已被设置时会立即抛出 **InterruptedException**。（通常，可中断的方法会在阻塞或进行重要的工作前首先检查中断，从而尽快地响应中断）。

如果代码不会调用可中断的阻塞方法，那么仍然可以通过在任务代码中轮询当前线程的中断状态来响应中断。要选择合适的轮询频率，就需要在效率和响应性之间进行权衡。如果响应性要求较高，那么就不应该调用那些执行时间较长并且不响应中断的方法。

在取消过程中可能涉及除了中断状态之外的其他状态，中断可以用来获得线程的注意，并且由中断线程保存的信息，可以为中断的线程提供进一步的指示。（当访问这些信息时，要确保使用同步。）

> **例如**，当一个由 **ThreadPoolExecutor** 拥有的工作者线程检测到中断时，它会检查线程池是否正在关闭。如果是，它会在结束之前执行一些线程池清理工作，否则它可能创建一个新线程将线程池恢复到合理的规模。

# 总结

本篇介绍了取消策略、线程中断、中断策略 和 响应中断的内容，下篇将要介绍如何编写任务和服务，使它们能对取消请求做出响应。
