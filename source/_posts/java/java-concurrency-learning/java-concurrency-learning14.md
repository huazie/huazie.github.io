---
title: Java并发编程学习14-任务关闭（上）
date: 2022-12-05 07:00:00
updated: 2024-03-19 17:21:35
categories:
  - [开发语言-Java,Java并发编程]
tags:
  - Java
  - 任务关闭
  - 关闭日志服务
  - ExecutorService
  - 毒丸对象
  - shutdownNow
---

[《开发语言-Java》](/categories/开发语言-Java/) [《Java并发编程》](/categories/开发语言-Java/Java并发编程/) 

![](/images/java-concurrency-logo.png)

# 引言
《任务关闭》由于篇幅较多，拆分了两篇来介绍各种任务和服务的关闭机制，以及如何编写任务和服务，使它们能够优雅地处理关闭。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

我们知道，应用程序通常会创建拥有多个线程的服务，例如线程池，并且这些服务的生命周期通常比创建它们的方法的生命周期更长。如果应用程序准备退出，那么这些服务所拥有的线程也需要结束。由于无法通过抢占式的方法来停止线程，因此它们需要自行结束。

线程应该有一个相应的所有者，即创建该线程的类。而线程池就是其工作者线程的所有者，如果要停止这些线程，那么应当通过线程池来操作。

应用程序可以拥有服务，服务也可以拥有工作者线程，但应用程序并不能拥有工作者线程，因此应用程序不能直接停止工作者线程。

工作者线程应当由它的拥有者来关闭，其拥有者需要提供生命周期方法来关闭它自己以及它拥有的线程。

例如，在 **ExecutorService** 中提供了 **shutdown** 和 **shutdownNow** 等方法。在其他拥有线程的服务中也应该提供类似的关闭机制。

> 对于拥有线程的服务，只要服务的存在时间大于创建线程的方法的存在时间，那么就应该提供生命周期方法。

# 主要内容
## 1. 关闭日志服务

下面我们来看一下如下的示例【不支持关闭的生产者-消费者日志服务】：

```java
public class LogWriter {

    private final BlockingQueue<String> queue;

    private final LoggerThread logger;

    public LogWriter(Writer writer) {
        this.queue = new LinkedBlockingQueue<>(100);
        this.logger = new LoggerThread(writer);
    }

    public void start() {
        logger.start();
    }

    public void log(String msg) throws InterruptedException {
        queue.put(msg);
    }

    private class LoggerThread extends Thread {

        private final PrintWriter writer;

        LoggerThread(Writer writer) {
            this.writer = (PrintWriter) writer;
        }

        public void run() {
            try {
                while (true)
                    writer.println(queue.take());
            } catch (InterruptedException e) {
                //
            } finally {
                writer.close();
            }
        }
    }
}
```

如上示例 **LogWriter** 中给出了一个简单的日志服务示例，其中日志操作在单独的日志线程中执行。产生日志消息的线程并不会将消息直接写入输出流，而是由 **LogWriter** 通过 **BlockingQueue** 将消息提交给日志线程，并由日志线程写入。

这是一种多生产者单消费者的设计方式：每个调用 **log** 的操作都相当于一个生产者，而后台的日志线程则相当于消费者。如果消费者的处理速度低于生产者的生成速度，那么 **BlockingQueue** 将阻塞生产者，直到日志线程有能力处理新的日志消息。

当然，上述 **LogWriter** 目前是无法关闭的。为了避免使 **JVM** 也无法正常关闭，**LogWriter** 还需要实现停止日志服务的逻辑。

**那 LogWriter 该如何实现停止日志服务呢？**

如上示例中的日志线程 **LoggerThread** 中会循环调用阻塞队列的 **take** 方法，而我们知道 **take** 方法可以响应中断。而且 **LoggerThread** 中已经包含了捕获 **InterruptedException** 时退出的逻辑，那么只需要中断日志线程 **LoggerThread** 就能停止日志服务。

不过，如果只是让日志线程退出，这还不是一种完备的关闭机制。

**那它会带来什么问题？又该如何理解呢？**

- 首先，这种直接关闭的做法会丢失那么正在等待被写入到日志的信息；

- 其次，其他线程将在调用 **log** 时被阻塞，因为日志消息队列是满的，因此这些线程将无法解除阻塞状态。

当取消一个 **生产者--消费者** 操作时，需要同时取消生产者和消费者。在上面示例中，在中断日志线程时会处理消费者，但由于生产者并不是专门的线程，因此要取消它们将非常困难。

那么我们还能想到什么方法可以用来关闭日志服务呢？

下面我们来看一下如下的示例【通过一种不可靠的方式为日志服务增加关闭支持】：

```java
    public void log(String msg) throws InterruptedException {
        if (!shutdownRequested)
            queue.put(msg);
        else
            throw new IllegalStateException("logger is shut down");
    }
```

上述示例，通过判断一个 “已请求关闭” 标志，以避免进一步提交日志消息。不过，这里的 **log** 方法存在着竞态条件问题。由于它是 “先判断再运行”，如果生产者调用 **log** 方法，在判断“已请求关闭” 标志时，发现该服务还没有关闭，此时在调用 **put** 之前，正好关闭服务，那生产者仍然会将日志消息放入队列，这同样会使得生产者可能在调用 **log** 时阻塞并且无法解除阻塞状态。

为了让 **LogWriter** 能够提供可靠的关闭操作，必须解决上面的竞态条件问题，这就需要使得上面 **log** 方法的日志消息的提交操作成为原子操作。然而，我们不希望在消息加入队列时去持有一个锁，因为 **put** 方法本身就可以阻塞。

下面我们来看一下如下的示例【向 **LogWriter** 添加可靠的取消操作】：

```java
public class LogService {

    private final BlockingQueue<String> queue;

    private final LoggerThread logger;

    private final PrintWriter writer;

    @GuardedBy("this")
    private boolean isShutdown;

    @GuardedBy("this")
    private int reservations;

    public LogService(Writer writer) {
        this.queue = new LinkedBlockingQueue<>(100);
        this.writer = (PrintWriter) writer;
        this.logger = new LoggerThread();
    }

    public void start() {
        logger.start();
    }

    public void stop() {
        synchronized (this) {
            isShutdown = true;
        }
        logger.interrupt();
    }

    public void log(String msg) throws InterruptedException {
        synchronized (this) {
            // 如果关闭标识为true，则抛出异常
            if (isShutdown)
                throw new IllegalStateException("logger is shut down");
            // 生产者生产消息，消息记录数 + 1
            ++reservations;
        }
        queue.put(msg);
    }

    private class LoggerThread extends Thread {
        public void run() {
            try {
                while (true) {
                    try {
                        synchronized (LogService.this) {
                            // 如果关闭标识为true，并且消息记录数为 0，则退出日志消费线程
                            // 如果关闭标识为true，并且消费记录数大于0，这时因为生产者继续调用log方法将会抛异常，所以消息记录数不会再增加；
                            // 那队列中剩余的数据就可以继续处理，直至消费记录数为0，然后再退出日志消费线程
                            if (isShutdown && reservations == 0)
                                break;
                        }
                        String msg = queue.take();
                        synchronized (LogService.this) {
                            // 消费者消费消息，消息记录数 - 1
                            --reservations;
                        }
                        writer.println(msg);
                    } catch (InterruptedException e) {
                        // retry
                    }
                }
            } finally {
                writer.close();
            }
        }
    }
}
```

如上示例 **LogService** 通过原子方法来检查关闭请求，并且有条件地递增一个计数器来 “保持” 提交消息的权利。

## 2. 关闭 ExecutorService

在前面的[《任务执行与Executor框架》](/2022/10/03/java/java-concurrency-learning/java-concurrency-learning10/)报文中，我们了解了 **ExecutorService** 提供了两种关闭方法：

1. 使用 **shutdown** 正常关闭
2. 使用 **shutdownNow** 强行关闭

这两种关闭方式的区别在于各自的安全性和响应性，如下：

- 强行关闭的速度更快，但风险更大，因为任务很可能在执行到一半时被结束。
- 正常关闭虽然速度慢，但却更安全，因为 **ExecutorService** 会一直等到队列中的所有正在执行的任务都执行完成后才关闭。

下面我们来看一下如下的示例【**LogService** 的一种变化形式，使用 **ExecutorService** 的日志服务】：

```java
public class LogService {

    private static final long TIMEOUT = 1;

    private static final TimeUnit UNIT = TimeUnit.SECONDS;

    private final ExecutorService exec = Executors.newSingleThreadExecutor();

    private final PrintWriter writer;

    public LogService(Writer writer) {
        this.writer = (PrintWriter) writer;
    }

    public void start() {

    }

    public void stop() throws InterruptedException {
        try {
            exec.shutdown();
            exec.awaitTermination(TIMEOUT, UNIT);
        } finally {
            writer.close();
        }
    }

    public void log(String msg) {
        try {
            exec.execute(new WriteTask(msg));
        } catch (RejectedExecutionException ignored) {
            //
        }
    }

    class WriteTask implements Runnable {
        
        String msg;

        WriteTask(String msg) {
            this.msg = msg;
        }

        @Override
        public void run() {
            writer.println(msg);
        }
    }
}
```

上述示例 **LogService** 将管理线程的工作委托给一个 **ExecutorService**，而不是由其自行管理。通过封装 **ExecutorService**，可以将所有权链从应用程序扩展到服务以及线程，所有权链上的各个成员都将管理它所拥有的服务或线程的生命周期。

## 3. "毒丸" 对象

另一种关闭生产者--消费者服务的方式就是使用 “**毒丸（Poison Pill）**” 对象：它是指一个放在队列上的对象，当从队列中取到该对象时，服务立即停止。

在 **FIFO**（先进先出）队列中，“毒丸” 对象将确保消费者在关闭之前首先完成队列中的所有工作，在提交 “毒丸” 对象之前提交的所有工作都会被处理，而生产者在提交了 “毒丸” 对象后，将不会再提交任何工作。

下面我们来看一下如下的示例【通过 “毒丸” 对象来关闭服务】

```java
public class IndexingService {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(IndexingService.class);

    private static final File POISON = new File("");

    private final CrawlerThread producer = new CrawlerThread();

    private final IndexerThread consumer = new IndexerThread();

    private final BlockingQueue<File> queue;

    private final FileFilter fileFilter;

    private final File root;

    public IndexingService(BlockingQueue<File> queue, FileFilter fileFilter, File root) {
        this.queue = queue;
        this.fileFilter = fileFilter;
        this.root = root;
    }

    public void start() {
        producer.start();
        consumer.start();
    }

    public void stop() {
        producer.interrupt();
    }

    public void awaitTermination() throws InterruptedException {
        consumer.join();
    }

    class CrawlerThread extends Thread {

        @Override
        public void run() {
            try {
                crawl(root);
            } catch (InterruptedException e) {
                // 发生异常
                LOGGER.debug("Occur InterruptedException");
            } finally {
                while (true) {
                    try {
                        // 生产者无论成功失败，最后都添加一个"毒丸"对象
                        queue.put(POISON);
                        break;
                    } catch (InterruptedException el) {
                        // 重新尝试
                    }
                }
            }
        }

        private void crawl(File root) throws InterruptedException {
            File[] entries = root.listFiles(fileFilter);
            if (entries != null) {
                for (File entry : entries)
                    if (entry.isDirectory())
                        crawl(entry);
                    else if (!FileRecord.alreadyIndexed(entry))
                        queue.put(entry);
            }
        }
    }

    class IndexerThread extends Thread {

        @Override
        public void run() {
            try {
                while (true) {
                    File file = queue.take();
                    // 如果遇到 "毒丸" 对象，则跳出循环，退出消费者线程
                    if (file == POISON)
                        break;
                    else
                        FileRecord.indexFile(file);
                    if (LOGGER.isDebugEnabled()) {
                        LOGGER.debug(file.getAbsolutePath());
                    }
                }
            } catch (InterruptedException e) {
                // 
            }
        }
    }
} 
```

只有在生产者和消费者的数量都已知的情况下，才可以使用 “毒丸” 对象。

在上述示例 **IndexingService** 中采用的解决方案可以拓展到多个生产者：只需要每个生产者都向队列中放入一个 “毒丸” 对象，并且消费者仅当在接收到 $N_{producers}$ 个 “毒丸” 对象时才停止。当然也可以拓展到多个消费者的情况，只需生产者将 $N_{consumers}$ 个 “毒丸” 对象放入队列。

## 4. 只执行一次的服务

如果某个方法需要处理一批任务，并且当所有任务都处理完成后才返回，那么可通过一个私有的 **Executor** 来简化服务的生命周期管理，其中该 **Executor** 的生命周期是由这个方法来控制的。

下面我们来看一下如下的示例【使用私有的 **Executor**，并且该 **Executor** 的生命周期受限于方法调用】

```java
    boolean checkMail(Set<String> hosts, long timeout, TimeUnit unit) throws InterruptedException {
        ExecutorService exec = Executors.newCachedThreadPool();
        final AtomicBoolean hasNewMail = new AtomicBoolean(false);
        try {
            for (final String host : hosts)
                exec.execute(new Runnable() {
                    public void run() {
                        if (checkMail(host))
                            hasNewMail.set(true);
                    }
                });
        } finally {
            exec.shutdown();
            exec.awaitTermination(timeout, unit);
        }
        return hasNewMail.get();
    }
```

如上示例的 **checkMail** 方法能在多台主机上并行地检查新邮件。它创建了一个私有的 **Executor**，并向每台主机提交一个任务。然后，当所有邮件检查任务都执行完成后，关闭 **Executor** 并等待结束。

> 之所以采用 **AtomicBoolean** 来代替 **volatile** 类型的 **boolean**，是因为能从内部的 **Runnable** 中访问 **hasNewMail** 标志，因此它必须是 **final** 类型以免被修改。


## 5. shutdownNow 的局限性

当通过 **shutdownNow** 来强行关闭 **ExecutorService** 时，它会尝试取消正在执行的任务，并返回所有已提交但尚未开始的任务，从而将这些任务写入日志或者保存起来以便之后进行处理。

> 需要关注的是，**shutdownNow** 返回的 **Runnable** 对象可能与提交给 **ExecutorService** 的 **Runnable** 对象并不相同：它们可能是被封装过的已提交任务。

在关闭过程中只会返回尚未开始的任务，而不会返回正在执行的任务。那要想知道当 **Executor** 关闭时那些任务正在执行，我们该怎么办呢？

下面我们来看一下如下的示例【在 **ExecutorService** 中跟踪在关闭之后被取消的任务】：

```java
public class TrackingExecutor extends AbstractExecutorService {

    private final ExecutorService exec;

    private final Set<Runnable> tasksCancelledAtShutdown = Collections.synchronizedSet(new HashSet<>());

    public TrackingExecutor(ExecutorService exec) {
        this.exec = exec;
    }

    public List<Runnable> getCancelledTasks() {
        if (!isTerminated())
            throw new IllegalStateException("The Executor is not terminated !");
        return new ArrayList<>(tasksCancelledAtShutdown);
    }

    @Override
    public void execute(final Runnable command) {
        exec.execute(new Runnable() {
            @Override
            public void run() {
                try {
                    command.run();
                } finally {
                    if (isShutdown() && Thread.currentThread().isInterrupted())
                        tasksCancelledAtShutdown.add(command);
                }
            }
        });
    }

    @Override
    public void shutdown() {
        exec.shutdown();
    }

    @Override
    public List<Runnable> shutdownNow() {
        return exec.shutdownNow();
    }

    @Override
    public boolean isShutdown() {
        return exec.isShutdown();
    }

    @Override
    public boolean isTerminated() {
        return exec.isTerminated();
    }

    @Override
    public boolean awaitTermination(long timeout, TimeUnit unit) throws InterruptedException {
        return exec.awaitTermination(timeout, unit);
    }
}
```

上述示例 **TrackingExecutor** 给出了如何在关闭过程中判断正在执行的任务。通过封装 **ExecutorService** 并使得 **execute** 等方法记录那些任务时在关闭后取消的，**TrackingExecutor** 可以找出那些任务已经开始但还没有正常完成。在 **Executor** 结束后，**getCancelledTasks** 返回被取消的任务清单。

> 要使上述示例给出的方案能够奏效，任务在返回时必须维持线程的中断状态，这也是所有设计良好的任务中都会实现的功能。

我们知道，网页爬虫程序的工作通常是无穷尽的，因此当爬虫程序必须关闭时，我们通常希望保存它的状态，以便稍后重新启动。

下面我们再来看一个示例【使用 **TrackingExecutor** 来保存未完成的任务以备后续执行】：

```java
public abstract class WebCrawler {

    private static final long TIMEOUT = 1;

    private static final TimeUnit UNIT = TimeUnit.SECONDS;

    private volatile TrackingExecutor exec;

    private final Set<URL> urlsToCrawl = new HashSet<>();

    public synchronized void start() {
        exec = new TrackingExecutor(Executors.newCachedThreadPool());
        for (URL url : urlsToCrawl) {
            submitCrawlTask(url);
        }
        urlsToCrawl .clear();
    }

    public synchronized void stop() throws InterruptedException {
        try {
            saveUncrawled(exec.shutdownNow());
            if (exec.awaitTermination(TIMEOUT, UNIT))
                saveUncrawled(exec.getCancelledTasks());
        } finally {
            exec = null;
        }
    }

    protected abstract List<URL> processPage(URL url);

    private void saveUncrawled(List<Runnable> uncrawled) {
        for (Runnable task : uncrawled)
            urlsToCrawl.add(((CrawlTask) task).getPage());
    }

    private void submitCrawlTask(URL url) {
        exec.execute(new CrawlTask(url));
    }

    private class CrawlTask implements Runnable {

        private final URL url;

        CrawlTask(URL url) {
            this.url = url;
        }

        public URL getPage() {
            return url;
        }

        @Override
        public void run() {
            for (URL link : processPage(url)) {
                if (Thread.currentThread().isInterrupted())
                    return;
                submitCrawlTask(link);
            }
        }
    }
}
```

上述示例 **WebCrawler** 给出了 **TrackingExecutor** 的用法，它的 **CrawlTask** 方法提供了一个 **getPage** 方法，该方法能找出正在处理的页面。当爬虫程序关闭时，无论是还没有开始的任务，还是那些被取消的任务，都将记录它们的 **URL**，因此当爬虫程序重新启动时，就可以将这些 **URL** 的页面抓取任务加入到任务队列中。

讲到这里，我们需要注意到实际上在 **TrackingExecutor** 中会存在一个不可避免的竞态条件，从而产生 ”误报“ 问题：一些被认为已取消的任务实际上已经执行完成。

**那这个问题是怎么产生的呢？**

我们回到 **TrackingExecutor** 的任务执行的逻辑中，在任务执行最后一条指令以及线程池将任务记录为 ”结束“ 的两个时刻之间，线程池可能被关闭，那么这个任务也会被添加到 **tasksCancelledAtShutdown**【即被取消的任务列表】里。

如果任务是幂等的（即将任务执行两次与执行一次会得到相同的结果），那么这不会存在问题，就比如我们上面的网页爬虫程序就是这种情况。否则，设计者应当考虑到这种风险，做好应对。

# 总结
本篇介绍了基于线程的服务关闭相关的内容，下一篇介绍《任务关闭》剩下的内容，敬请期待！！！
