---
title: Java并发编程学习11-任务执行演示
date: 2022-10-15 19:44:03
updated: 2024-03-19 17:21:35
categories:
  - [开发语言-Java,Java并发编程]
tags:
  - Java
  - 任务执行演示
  - Executor框架
---



![](/images/java-concurrency-logo.png)

# 引言

上一篇博文带大家了解了任务执行和 **Executor** 框架的基础知识，本篇将结合这些内容，演示一些不同版本的任务执行Demo，并且每个版本都实现了不同程度的并发性。

以下的示例是要实现浏览器程序中的页面渲染功能：将 HTML 页面绘制到图像缓存中【为了简便，假设 HTML 页面只包含标签文本、预定义大小的图片和URL】。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 1. 串行的页面渲染器

最简单实现页面渲染器功能就是对 **HTML** 文档进行串行处理。首先绘制文本元素，同时为图像预留出矩形的占位空间，在处理完第一遍文本后，程序再开始下载图像，并将它们绘制到相应的占位空间中。

```java
public class SingleThreadRenderer {  
    void renderPage (CharSequence source){
        renderText(source);
        List<ImageData> imageData = new ArrayList<>();
        for (ImageInfo imageInfo : sacanForImageInfo(source))
            imageData.add(imageInfo.downloadImage());
        for (ImageData data : imageData)
            renderImage(data);
    }
}
```

上述图像下载过程的大部分时间都是在等待 **I/O** 操作执行完成，在这期间 **CPU** 几乎不做任何工作。这种串行执行方法没有充分地利用 **CPU**，用户在看到最终页面需要等待过长的时间。

这个时候通过将上述串行执行的任务分解为多个独立的任务并发执行，就能够获得更高的 **CPU** 利用率和响应灵敏度。

# 2. 携带结果的任务

从[《任务执行和Executor框架》](/2022/10/03/java/java-concurrency-learning/java-concurrency-learning10/)的那篇博文中，我们知道 **Executor** 框架使用 **Runnable** 作为其基本的任务表示形式。但是 **Runnable** 也有自己的局限性，它不能 **返回一个值** 或 **抛出一个受检查的异常**。 

实际上，许多任务都是存在延迟的计算，比如：

- 执行数据库查询
- 从网络上获取资源
- 计算某个复杂的功能

对于这些延迟的任务，**Callable** 其实是个更好的任务表示形式，它的主入口点（即 **call**）将返回一个值，并可能抛出一个异常。在 `java.util.concurrent.Executors` 中包含了一些辅助方法【**callable**】能将其他类型的任务【**Runnable** 、**java.security.PrivilegedAction** 和 **java.security.PrivilegedExceptionAction**】封装为一个 **Callable**。

```java
public interface Callable<V> {
    V call() throws Exception;
}
```

> 可以使用 `Callable<Void>` 来表示无返回值的任务。 

从[《同步工具类》](/2022/09/17/java/java-concurrency-learning/java-concurrency-learning8/)的那篇博文中，我们知道 **Future** 表示一个任务的生命周期，它提供了相应的方法来判断是否已经完成或取消，以及获取任务的结果和取消任务等。在 **Future** 的规范中，任务的生命周期只能前进，不能后退，就像 **ExecutorService** 的生命周期一样。当某个任务完成后，它就永远停留在 “完成” 状态上。

```java
public interface Future<V> {
    boolean cancel(boolean mayInterruptIfRunning);
    boolean isCancelled();
    boolean isDone();
    V get() throws InterruptedException, ExecutionException, CancellationException;
    // 支持限时的获取结果
    V get(long timeout, TimeUnit unit) throws InterruptedException, ExecutionException, CancellationException, TimeoutException;
}
```

> 在 **Executor** 框架中，已提交但尚未开始的任务可以取消，但对于那些已经开始执行的任务，只有它们能响应中断时，才能取消。已经完成的任务可以随便取消，无任何影响。

那么如何创建一个 **Future** 来描述任务呢？

- **ExecutorService** 中的所有 **submit** 方法，可以将一个 **Runnable** 或 **Callable** 提交给 **Executor**，并得到一个 **Future** 用来获得任务的执行结果或者取消任务。
- 也可以显式为一个 **Runnable** 或 **Callable** 实例化一个 **FutureTask**，因为 **FutureTask** 实现了 **Runnable**，因此可以将它提交给 **Executor** 来执行【其实 **submit** 方法也是这么做的】。

从 **Java6** 开始，**ExecutorService** 实现可以改写 **AbstractExecutorService** 中的 **newTaskFor** 方法，从而根据已提交的 **Runnable** 或 **Callable** 来控制 **Future** 的实例化过程。

如下代码清单【**AbstractExecutorService** 中的 **newTaskFor** 方法的默认实现、**submit** 方法实现】：

```java
protected <T> RunnableFuture<T> newTaskFor(Runnable runnable, T value) {
    return new FutureTask<T>(runnable, value);
}

protected <T> RunnableFuture<T> newTaskFor(Callable<T> callable) {
    return new FutureTask<T>(callable);
}

public Future<?> submit(Runnable task) {
    if (task == null) throw new NullPointerException();
    RunnableFuture<Void> ftask = newTaskFor(task, null);
    execute(ftask);
    return ftask;
}

public <T> Future<T> submit(Runnable task, T result) {
    if (task == null) throw new NullPointerException();
    RunnableFuture<T> ftask = newTaskFor(task, result);
    execute(ftask);
    return ftask;
}

public <T> Future<T> submit(Callable<T> task) {
    if (task == null) throw new NullPointerException();
    RunnableFuture<T> ftask = newTaskFor(task);
    execute(ftask);
    return ftask;
}
```

# 3. 使用 Future 实现页面渲染器

为了使页面渲染器实现更高的并发性，首先将渲染过程分解为两个任务，一个是渲染所有的文本，另一个是下载所有的图像。

下面我们来看一下如下示例【使用 **Future** 等待图像下载】：

```java
public class FutureRenderer {

    private final ExecutorService executor = Executors.newCachedThreadPool();

    void renderPage (CharSequence source){
        final List<ImageInfo> imageInfos = scanForImageInfo(source);
        Callable<List<ImageData>> task = new Callable<List<ImageData>>() {
                public List<ImageData> call() {
                    List<ImageData> result = new ArrayList<ImageData>() ;
                    for (ImageInfo imageInfo : imageInfos)
                        result.add(imageInfo.downloadImage());
                    return result;
                }
            };
        Future<List<ImageData>> future = executor.submit(task);
        renderText(source);
        try {
            List<ImageData> imageData = future.get();
            for (ImageData data : imageData)
                renderImage(data);
        } catch (InterruptedException e){
            //重新设置线程的中断状态
            Thread.currentThread().interrupt() ;
            //由于不需要结果，因此取消任务
            future.cancel(true);
        } catch (ExecutionException e) {
            throw launderThrowable(e.getCause());
        }
    }
}
```
上述 **FutureRenderer** 中创建了一个 **Callable** 来下载所有的图像，并将其提交到一个 **ExecutorService**，这将返回一个描述任务执行情况的 **Future**。后面当主任务需要图像时，通过 **Future.get** 方法就可以获得所有下载的图像，即使还没下载好，至少下载任务已经开始了。

# 4. 使用 CompletionService 实现页面渲染器
在上面的 **FutureRenderer** 里 ，我们已经并行地执行了两个不同类型的任务--**下载图像** 与 **渲染文本**。如果渲染文本的速度远远高于下载图像的速度，那么程序的最终性能与串行执行时的性能差别不大，反而代码更加复杂了。其实用户不必等到所有的图像都下载好，而是希望每下载完一幅图像就立即显示出来。

下载图像的任务还可以继续细分，为每一幅图像的下载都创建一个独立任务，并在线程池中执行它们，从而将串行的下载过程转换为并行的过程，这样也就减少下载所有图像的总时间。

下面我们来看下如下的示例【使用 **CompletionService**，使页面元素在下载完成后立即显示出来】：

```java
public class CompletionServiceRenderer {

    private final ExecutorService executor;

    CompletionServiceRenderer(ExecutorService executor) {
        this.executor = executor;
    }

    void renderPage(CharSequence source) {

        List<ImageInfo> info = scanForImageInfo(source);

        CompletionService<ImageData> completionService = new ExecutorCompletionService<>(executor);

        for (final ImageInfo imageInfo : info)
            completionService.submit(new Callable<ImageData>() {
                public ImageData call() {
                    return imageInfo.downloadImage();
                }
            });

        renderText(source);

        try {
            for (int t = 0, n = info.size(); t < n; t++) {
                Future<ImageData> f = completionService.take();
                ImageData imageData = f.get();
                renderImage(imageData);
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        } catch (ExecutionException e) {
            throw launderThrowable(e.getCause());
        }
    }
}    
```

下面我们来简单了解下 **CompletionService**【完成服务】：

- **CompletionService** 将 **Executor** 和 **BlockingQueue** 的功能融合在一起。可以将 **Callable** 任务提交给它来执行，然后使用类似于队列操作的 **take** 和 **poll** 等方法来获得已完成的结果，而这些结果将会在完成时被封装为 **Future**。

- **CompletionService** 有个子类实现为 **ExecutorCompletionService**。它的构造函数中会创建一个 **BlockingQueue** 来保存计算结果。当计算完成时，调用 **FutureTask** 中的 **done** 方法。当提交某个任务时，该任务将首先包装为一个 **QueueingFuture**，这是 **FutureTask** 的一个子类，它覆写了父类的 **done** 方法，并将结果放入 **BlockingQueue** 中。**take** 和 **poll** 方法委托给了 **BlockingQueue**，这些方法会在出结果之前阻塞。

如下为 JDK 1.8 中 **ExecutorCompletionService** 里的 **QueueingFuture** 实现【其他版本可能有差异，以实际为准】

```java
public class ExecutorCompletionService<V> implements CompletionService<V> {
    private final Executor executor;
    private final AbstractExecutorService aes;
    private final BlockingQueue<Future<V>> completionQueue;

    private class QueueingFuture extends FutureTask<Void> {
        QueueingFuture(RunnableFuture<V> task) {
            super(task, null);
            this.task = task;
        }
        protected void done() { completionQueue.add(task); }
        private final Future<V> task;
    }

    public ExecutorCompletionService(Executor executor) {
        if (executor == null)
            throw new NullPointerException();
        this.executor = executor;
        this.aes = (executor instanceof AbstractExecutorService) ?
            (AbstractExecutorService) executor : null;
        this.completionQueue = new LinkedBlockingQueue<Future<V>>();
    }
    
    public ExecutorCompletionService(Executor executor, BlockingQueue<Future<V>> completionQueue) {
        if (executor == null || completionQueue == null)
            throw new NullPointerException();
        this.executor = executor;
        this.aes = (executor instanceof AbstractExecutorService) ?
            (AbstractExecutorService) executor : null;
        this.completionQueue = completionQueue;
    }
    
    // 其他方法省略
}
```

> 从 **ExecutorCompletionService** 的构造函数可知，多个 **ExecutorCompletionService** 可以共享一个 **Executor**，因此可以创建一个对于特定计算私有，又能共享一个公共 **Executor** 的 **ExecutorCompletionService**。

# 5. 为任务设置时限

下面我们来看下如下的案例：

- 某个 **Web** 应用程序从外部的广告服务器上获取广告信息，但如果该应用程序在两秒内得不到响应，将直接显示默认的广告，这样即使无法获得广告信息，也不会降低站点的响应性能。
- 一个门户网站从多个数据源并行地获取数据，但可能只会在指定的时间内等待数据，如果超出了等待时间，那么将只显示已经获得的数据。

上述案例都规定了任务要在指定的时间内完成，如果某个任务无法在指定时间内完成，那么将不再需要它的结果，此时就应当放弃这个任务。

那么如何 **给任务设置时限** 呢？

前面提到的支持时间限制的 **Future.get** 支持给任务设置时限：当结果可用时，它将立即返回，如果在指定时限内没有计算出结果，那么将抛出 **TimeoutException**。

如果任务超时了该如何取消呢？

上述通过支持时间限制的 **Future.get** 获取任务结果。如果任务超时了，它会抛出 **TimeoutException**，这时可以通过 **Future.cancel** 来取消任务。

## 5.1 限时获取广告信息示例

下面我们来看下如下的示例【在指定时间内获取广告信息】：

```java
public class PageAdRenderer {

    private static final Long TIME_BUDGET = 2000000000L;

    private final ExecutorService executor = Executors.newCachedThreadPool();

    private final Ad DEFAULT_AD = new Ad();

    public Page renderPageWithAd() throws InterruptedException {
        long endNanos = System.nanoTime() + TIME_BUDGET;
        Future<Ad> f = executor.submit(new FetchAdTask());
        // 等待广告的同时显示页面
        Page page = renderPageBody();
        Ad ad;
        try {
            // 只等待指定的时间长度
            long timeLeft = endNanos - System.nanoTime();
            ad = f.get(timeLeft, NANOSECONDS);
        } catch (ExecutionException e) {
            ad = DEFAULT_AD;
        } catch (TimeoutException e) {
            ad = DEFAULT_AD;
            f.cancel(true);
        }
        page.setAd(ad);
        return page;
    }
}
```

上述示例生成的页面中包括响应用户请求的内容以及从广告服务器上获得的广告。它将获取广告的任务提交给一个 **Executor**，然后计算剩余的文本页面内容，最后等待广告信息，直到超出指定的时间。如果 **get** 超时，那么将取消广告获取任务，并使用默认的广告信息。

**注意：**
- 传递给 **get** 方法的 **timeout** 参数的计算方法是，将 **指定时限** 减去 **当前时间**。这可能会得到负数，但 `java.util.concurrent` 中所有 **与时限相关的方法** 都将 **负数视为零**，因此不需要额外的代码来处理这种情况。
- **Future.cancel** 的参数为 **true** 表示任务线程可以在运行过程中中断【在后续博文会详细介绍】。

##  5.2 旅行预订门户网站示例
下面我们来考虑这样一个旅行预订门户网站：

用户输入旅行的日期和其他要求，门户网站获取并显示来自多条航线、旅店或汽车租赁公司的**报价**。在获取不同公司报价的过程中，可能会调用 **Web** 服务、访问数据库、执行一个 EDI 事务或其他机制。在这种情况下，页面应该只显示在指定时间内收到的信息。对于没有及时响应的服务提供者，页面可以忽略它们，或者显示一个提示信息。

从一个公司获取报价的过程与从其他公司获得报价的过程无关，因此可以将获取报价的过程当成一个任务，从而使获得报价的过程能并发执行。

通过上面了解的支持限时的 **Future.get** ，我们很容易想到如下的获取报价的逻辑：

> 创建 **n** 个获取报价的任务，并将其提交到一个线程池，同时保留 **n** 个 **Future**，并使用限时的 **get** 方法通过 **Future** 串行地获取每一个结果。

虽然上面也可行，但是现在我们有更好的方法，下面来看一下如下示例【使用线程池的 **invokeAll** 方法】：

```java
public class TravelWebSite {

    private final ExecutorService executor = Executors.newCachedThreadPool();

    public List<TravelQuote> getRankedTravelQuotes(TravelInfo travelInfo,
                                                   Set<TravelCompany> companies,
                                                   Comparator<TravelQuote> ranking,
                                                   long time, TimeUnit unit) throws InterruptedException {
        List<QuoteTask> tasks = new ArrayList<>();

        for (TravelCompany company : companies) 
            tasks.add(new QuoteTask(company, travelInfo));

        List<Future<TravelQuote>> futures = executor.invokeAll(tasks, time, unit);

        List<TravelQuote> quotes = new ArrayList<>(tasks.size());

        Iterator<QuoteTask> taskIterator = tasks.iterator();

        for (Future<TravelQuote> future : futures) {
            QuoteTask task = taskIterator.next();
            try {
                quotes.add(future.get());
            } catch (ExecutionException e) {
                quotes.add(task.getFailureQuote(e.getCause()));
            } catch (CancellationException e) {
                quotes.add(task.getTimeoutQuote(e));
            }
        }

        Collections.sort(quotes, ranking);
        return quotes;
    }
}
```

如上示例使用了支持限时的 **invokeAll** 方法，将多个任务提交给一个 **ExecutorService** 并获得结果。

关于 **invokeAll** 方法，有如下几点需要了解：

- **invokeAll** 方法的参数为一组任务，并返回一组 **Future**。这两个集合有着相同的结构。
- **invokeAll** 方法按照任务集合中迭代器的顺序将所有的 **Future** 添加到返回的集合中，从而使调用者能将各个 **Future** 与其表示的 **Callable** 关联起来。
- 当所有任务都执行完毕时，或者调用线程被中断时，又或者超过指定时限时，**invokeAll** 将返回。
- 当超过指定时限后，任何还未完成的任务都会取消。
- 当 **invokeAll** 方法返回后，每个任务要么正常地完成，要么被取消，而客户端代码可以调用 **get** 或 **isCancelled** 来判断究竟是何种情况。

# 总结
本文以Demo的形式演示了如何寻找任务中更细粒度的并发场景，对我们的并发应用开发有着一定的借鉴意义。了解了任务执行的基本知识，下篇博文开始我们将介绍如何优雅地取消和关闭任务，敬请期待！