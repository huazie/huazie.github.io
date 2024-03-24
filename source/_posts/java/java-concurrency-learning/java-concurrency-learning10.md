---
title: Java并发编程学习10-任务执行与Executor框架
date: 2022-10-03 18:41:27
updated: 2024-03-19 17:21:35
categories:
  - [开发语言-Java,Java并发编程]
tags:
  - Java
  - 任务执行
  - Executor框架
  - 
---

[《开发语言-Java》](/categories/开发语言-Java/) [《Java并发编程》](/categories/开发语言-Java/Java并发编程/) 

![](/images/java-concurrency-logo.png)

# 任务执行
何为任务？ 任务通常是一些抽象且离散的工作单元。

大多数并发应用程序都是围绕着 **“任务执行”** 来构造的。而围绕着 **“任务执行”** 来设计应用程序结构时，首先要做的就是要找出清晰的任务边界。大多数服务器应用程序都提供了一种自然的任务边界选择方式：以独立的客户请求为边界。将独立的请求作为任务边界，既可以实现任务的独立性，又可以实现合理的任务规模。

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

## 1. 串行地执行任务

在应用程序中可以通过多种策略来调度任务，其中最简单的策略就是在单个线程中串行地执行各项任务。

下面我们来看如下的代码示例【**SingleThreadWebServer** 将串行地处理它的任务（即通过 **80** 端口接收到的 **HTTP** 请求）】：

```java
public class SingleThreadWebServer {
    public static void main(String[] args) throws IOException {
        ServerSocket socket = new ServerSocket(80);
        while (true) {
            Socket connection = socket.accept();
            handleRequest(connection);
        }
    }
}
```

在实际执行中，上述示例执行性能非常糟糕，因为它每次只能处理一个请求。在服务器应用程序中，串行处理机制通常都无法提供高吞吐率或快速响应性。 

> 在某些情况下，串行处理方式能带来简单性或安全性。大多数 **GUI** 框架都通过单一的线程来串行地处理任务。后面的博文我们会再次介绍串行模型。

## 2. 显式地为任务创建线程

下面我们来看如下的代码示例【通过为每个请求创建一个新的线程来提供服务，从而实现更高的响应性】：

```java
public class ThreadPerTaskWebServer {
    public static void main(String[] args) throws IOException {
        ServerSocket socket = new ServerSocket(80);
        while (true) {
            final Socket connection = socket.accept();
            Runnable task = new Runnable() {
                public void run() {
                    handleRequest(connection);
                }
            };
            new Thread(task).start();
        }
    }
}
```

上述 **ThreadPerTaskWebServer** 与前面的 单线程版本 **SingleThreadWebServer** 在结构上类似，主线程仍然不断地交替执行 “接受外部连接” 与 “处理相关请求” 等操作。区别在于，**ThreadPerTaskWebServer** 对于每个连接，主循环都将创建一个新线程来处理请求，而不是在主循环中进行处理。

在正常负载情况下，**“为每个任务分配一个线程”** 的方法能提升串行执行的性能【即请求的到达速率不超出服务器的请求处理能力】。

## 3. 无限制创建线程的不足

当需要创建大量的线程时，**“为每个任务分配一个线程”** 就存在如下的问题了：

- 线程生命周期的开销非常高。创建线程需要时间，又会延迟请求的处理。如果请求的到达率非常高且请求的处理过程是轻量级的，那么为每个请求创建一个新线程将消耗大量的计算资源。
- 资源消耗。活跃的线程会消耗系统资源，尤其是内存。如果可运行的线程数量多于可用处理器的数量，那么有些线程将闲置。大量空闲的线程会占用许多内存，给垃圾回收器带来压力，而且大量线程在竞争 **CPU** 资源时还将产生其他的性能开销。
- 稳定性。可创建线程的数量存在着一个限制，该限制随着平台不同而不同，并且受多个因素制约，包括JVM的启动参数、**Thread** 构造函数中请求的栈大小，以及底层操作系统对线程的限制等。如果破坏了这些限制，应用可能就会抛出 **OutOfMemoryError** 异常。

在一定的范围内，增加线程可以提高系统的吞吐率，但如果超出了这个范围，再创建再多的线程也只会降低程序的执行速度，甚至最后整个应用程序都将崩溃。


# Executor框架

任务是一组逻辑工作单元，而线程则是使任务异步执行的机制。

前面我们已经分析了两种通过线程来执行任务的策略：

- 把所有任务放在单个线程中串行执行【它的问题在于有着糟糕的响应性和吞吐量】
- 将每个任务放在各自的线程中执行【它的问题在于资源管理的复杂性】

在 **Java** 类库中，任务执行的接口是 **Executor**，如下：

```java
public interface Executor {
    void execute(Runnable command);
}
```

上述 **Executor** 虽然简单，但它却为 **java.util.concurrent** 下的异步任务执行框架提供了基础，该框架有着如下的特点：

- 支持多种不同类型的任务执行策略。
- 提供了一种标准的方法将任务的提交过程与执行过程解耦开来，并用 **Runnable** 来表示任务。

**Executor** 的实现还提供了对生命周期的支持，以及统计信息收集、应用程序管理机制和性能监视等机制。

**Executor** 基于生产者-消费者模式，提交任务的操作相当于生产者（生成待完成的工作单元），执行任务的线程则相当于消费者（执行完这些工作单元）。

## 1. 基于 Executor 的 Web 服务器

下面我们来看下如下的示例【用 **Executor** 代替了硬编码的线程，采用了一个固定长度的线程池，可以容纳 **100** 个线程】：

```java
public class TaskExecutionWebServer {
    private static final int NTHREADS = 100;
    
    private static final Executor exec = Executors.newFixedThreadPool(NTHREADS);
    
    public static void main(String[] args) throws IOException {
        ServerSocket socket = new ServerSocket(80);
        while (true) {
            final Socket connection = socket.accept();
            Runnable task = new Runnable() {
                public void run() {
                    handleRequest(connection);
                }
            };
            exec.execute(task);
        }
    }
}
```

上述 **TaskExecutionWebServer** 采用 **Executor**，将请求处理任务的提交与任务的实际执行解耦开来，并且只需采用另一种不同的 **Executor** 实现，就可以改变服务器的行为。

下面我们来将 **TaskExecutionWebServer** 修改为类似前面 **ThreadPerTaskWebServer** 的行为，只需要使用一个为每个请求都创建新线程的 **Executor**。

```java
public class ThreadPerTaskExecutor implements Executor {
    public void execute(Runnable r) {
    new Thread(r).start();
    }
}

public class TaskExecutionWebServer {
    
    private static final Executor exec = new ThreadPerTaskExecutor();
    
    public static void main(String[] args) throws IOException {
        ServerSocket socket = new ServerSocket(80);
        while (true) {
            final Socket connection = socket.accept();
            Runnable task = new Runnable() {
                public void run() {
                    handleRequest(connection);
                }
            };
            exec.execute(task);
        }
    }
}
```

同样，还可以编写一个 **Executor**，使得 **TaskExecutionWebServer** 的行为类似于单线程的行为，即以同步的方式执行每个任务，然后再返回。

```java
public class WithInThreadExecutor implements Executor {
    public void execute(Runnable r) {
    r.run();
    }
}

public class TaskExecutionWebServer {
    
    private static final Executor exec = new WithInThreadExecutor();
    
    public static void main(String[] args) throws IOException {
        ServerSocket socket = new ServerSocket(80);
        while (true) {
            final Socket connection = socket.accept();
            Runnable task = new Runnable() {
                public void run() {
                    handleRequest(connection);
                }
            };
            exec.execute(task);
        }
    }
}
```
从上面可以了解到，将任务的提交与执行解耦开来，通过简单的改动就可以为某种类型的任务指定和修改执行策略。
## 2. 执行策略
各种执行策略本质上都是一种资源管理工具，最佳策略还得取决于可用的计算资源以及对服务质量的需求。

- 通过限制并发任务的数量，可以确保应用程序不会由于资源耗尽而失败，或者由于在稀缺资源上发生竞争而严重影响性能。
- 通过将任务的提交与任务的执行策略分离开来，有助于在部署阶段选择与可用硬件资源最匹配的执行策略。

## 3. 线程池

**线程池**，是指管理一组同构工作线程的资源池。线程池与工作队列密切相关，其在工作队列中保存了所有等待执行的任务。

工作者线程是如何工作的呢？

它从工作队列中获取一个任务，执行任务，然后返回线程池并等待下一个任务。

采用线程池的好处有哪些呢？

- 通过重用现有的线程而不是创建新线程，可以在处理多个请求时分摊在线程创建和销毁过程中产生的巨大开销。
- 当请求到达时，工作线程通常已经存在了，因此不会由于等待创建线程而延迟任务的执行，从而提高了响应性。
- 通过适当调整线程池的大小，可以创建足够多的线程以便使处理器保持忙碌状态，同时还可以防止过多线程相互竞争资源而使应用程序耗尽内存而失败。

**Java** 类库中，可以通过调用 **Executors** 中的静态工厂方法来创建一个线程池：

- **newFixedThreadPool**。它将创建一个固定长度的线程池，每提交一个任务时就创建一个线程，直到达到线程池的最大数量，这时线程池的规模将不再变化（如果某个线程由于发生了未预期的 **Exception** 而结束，那么线程池会补充一个新的线程）。
- **newCachedThreadPool**。它将创建一个可缓存的线程池，如果线程池的当前规模超过了处理需求时，那么将回收空闲的线程，而当需求增加时，则可以添加新的线程，线程池的规模不存在任何限制。
- **newSingleThreadExecutor**。它是一个单线程的 **Executor**，它创建单个工作者线程来执行任务，如果这个线程异常结束，会创建另一个线程来替代。它能确保依照任务在队列中的顺序来串行执行（例如 **FIFO**、**LIFO**、**优先级**）。
- **newScheduledThreadPool**。它创建了一个固定长度的线程池，而且以延迟或定时的方式来执行任务，类似于 **Timer**。

## 4. Executor 的生命周期

由于 **Executor** 以异步方式来执行任务，因此在任何时刻，之前提交的任务的状态并不是立即可见的。这些任务可能的状态如下：

- **已经完成**
- **正在运行**
- **在队列中等待执行**

为了管理任务的不同状态，**Executor** 扩展了 **ExecutorService** 接口，添加了一些用于生命周期管理的方法，如下所示：

```java
public interface ExecutorService extends Executor {
    void shutdown();
    
    List<Runnable> shutdownNow();
    
    boolean isShutdown();
    
    boolean isTerminated();
    
    boolean awaitTermination(long timeout, TimeUnit unit) throws InterruptedException;
    // ...其他用于任务提交的方法
}
```

**ExecutorService** 的生命周期有 3 种状态：

- **运行**
- **关闭**
- **已终止**

**ExecutorService** 在初始创建时处于运行状态。**shutdown** 方法将执行平缓的关闭过程：不再接受新的任务，同时等待已经提交的任务执行完成--包括那些还未开始执行的任务。**shutdownNow** 方法将执行粗暴的关闭过程：它将尝试取消所有运行中的任务，并且不再启动队列中尚未开始执行的任务。

在 **ExecutorService** 关闭后提交的任务将由 **“拒绝执行处理器（Rejected Execution Handler）”** 来处理，它会抛弃任务，或者让 **execute** 方法抛出一个未检查的 **RejectedExecutionException**。等所有任务都完成后，**ExecutorService** 将转入终止状态。可以调用 **awaitTermination** 来等待 **ExecutorService** 到达终止状态，或者通过调用 **isTerminated** 来轮询 **ExecutorService** 是否已经终止。通常在调用 **awaitTermination** 之后会立即调用 **shutdown**，从而产生同步地关闭 **ExecutorService** 的效果。


下面我们来看下如下的示例【通过增加生命周期支持来扩展 **Web** 服务器的功能】；

```java
public class LifecycleWebServer {
    private final Executor exec = Executors.newCachedThreadPool();
    
    public void start() throws IOException {
        ServerSocket socket = new ServerSocket(80);
        while (!exec.isShutdown()) {
            try {
                final Socket connection = socket.accept();
                exec.execute(new Runnable() {
                    public void run() {
                        handleRequest(connection);
                    }
                });
            } catch (RejectedExecutionException e) {
            if (!exec.isShutdown()) 
            
            }
        }
    }
    
    public void stop() {
        exec.shutdown();
    }
    
    void handleRequest(Socket connection) {
        Request request = readRequest(connection);
        if (isShutdownRequest(request))
            stop();
        else
            dispatchRequest(request);
    }
}
```

上述 **LifecycleWebServer** 可以通过两种方式来关闭 Web 服务器：

- 在程序中直接调用 **stop** 方法。
- 以客户端请求形式向 **Web** 服务器发送一个特定格式的HTTP请求。


## 5. 延迟任务与周期任务

在 **Java** 类库中，**Timer** 类负责管理延迟任务【指定时间后执行该任务】以及周期任务【指定周期执行一次该任务】。

不过，**Timer** 使用上存在着如下的缺陷 ：

- **Timer** 在执行所有的定时任务时只会创建一个线程。如果某个任务的执行时间过长，那么将破坏其他 **TimerTask** 的定时准确性。
- **Timer** 在执行任务时并不捕获异常。因此当 **TimerTask** 抛出了一个未检查的异常时，将会终止定时线程。这时，**Timer** 也不会恢复线程的执行，而是会错误地认为整个 **Timer** 都被取消了。已经被调度但尚未执行的 **TimerTask** 将不会再执行，新的任务也不能被调度。【这个问题也称之为 **“线程泄漏”**，后续的博文将会介绍】 

> **Timer** 支持基于绝对时间而不是相对时间的调度机制，因此任务的执行对系统时钟变化很敏感，而 **ScheduledThreadPoolExecutor** 只支持基于相对时间的调度。

使用 **线程池** 就可以弥补上述缺陷，它可以提供多个线程来执行延时任务和周期任务。

下面我们来看一个演示 **Timer** 问题的示例：

```java
public class OutOfTime {
    public static void main(String[] args) throws Exception {
        Timer timer = new Timer();
        timer.schedule(new ThrowTask(), 1);
        SECONDS.sleep(1);
        timer.schedule(new ThrowTask(), 1);
        SECONDS.sleep(5);
    }
    
    static class ThrowTask extends TimerTask {
        public void run() {
            throw new RuntimeException();
        }
    }
}
```

上述程序，你也许会认为运行 **6** 秒后退出，但实际情况是运行 **1** 秒就结束了，并抛出了一个异常消息 **“Timer already cancelled”**。 如下图所示：

![](result.png)


> 注意：在 **Java 5.0** 或 更高的 **JDK** 中，将很少再使用 **Timer**。如果需要，**ScheduledThreadPoolExecutor** 是个不错的选择。

# 总结

本篇我们重点讲了任务执行与 **Executor** 框架的基础知识。**Executor** 帮助指定执行策略，但如果要使用 **Executor**，必须将任务表述为一个 **Runnable**【即必须定义一个清晰的任务边界】。大多数服务器应用程序中都存在一个明显的任务边界：单个客户请求。但有时候，任务边界并非是显而易见的，需要进一步的揭示其粒度更细的并发性。说到这里，那么下一遍博文就将通过 **Demo** 演示如何去找出可利用的并行性需求，敬请期待！！！



