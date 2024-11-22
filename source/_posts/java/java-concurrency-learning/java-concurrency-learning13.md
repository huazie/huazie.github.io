---
title: Java并发编程学习13-任务取消的进阶使用
date: 2022-11-15 15:16:18
updated: 2024-03-19 17:21:35
categories:
  - [开发语言-Java,Java并发编程]
tags:
  - Java
  - 任务取消
  - 任务限时运行
  - Future
  - newTaskFor
---



![](/images/java-concurrency-logo.png)

# 引言
《任务取消》由于篇幅较多，拆分了两篇来介绍各种实现取消和中断的机制，以及如何编写任务和服务，使它们能对取消请求做出响应。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 1. 任务限时运行

我们知道许多任务可能永远也无法处理完成（例如，枚举所有的素数），而某些任务，可能很快被处理掉，也可能很长时间才能处理完。这个时候如果能够对任务处理加个时间限制，比如指定 “最多花1分钟搜索答案” 或者 “枚举出1秒钟内能找到的素数”，那将会是非常有用的。

我们来思考一下，本篇开头提到的素数生成器 **PrimeGenerator**，如果它在指定时限内抛出一个未检查的异常，会怎么样呢？

可以肯定的是这个异常会被忽略掉，因为素数生成器在另一个独立的线程中运行，而这个线程并不会显式地处理异常。

大多数时候，我们还是希望能够知道在任务执行过程中是否会抛出异常的。

下面我们来看一下如下示例【在外部线程中安排中断，**不推荐使用，仅用于理解**】：

```java
public class TaskUtils {
    private static final ScheduledExecutorService cancelExec = Executors.newScheduledThreadPool(10);

    public static void timeRun(Runnable r, long timeout, TimeUnit unit) {
        final Thread taskThread = Thread.currentThread();
        cancelExec.schedule(new Runnable() {
            public void run() {
                taskThread.interrupt();
            }
        }, timeout, unit);
        r.run();
    }
}
```

上述示例给出了在指定时间内运行一个任意的 **Runnable** 的场景。**timeRun** 在调用线程中运行任务，并安排了一个取消任务，用于在运行指定的时间间隔后中断 **timeRun** 所在线程。从任务中抛出未检查异常的问题，也会被 **timeRun** 的调用者捕获。

下面我们来看一下如下测试场景【演示下 **1s** 后结束素数生成器的任务】：

```java
public class TimeRunTest {
    @Test
    public void timeRun() {
        LOGGER.debug("timeRun start");
        BlockingQueue<BigInteger> primes = new LinkedBlockingQueue<>(100);
        PrimeProducer producer = new PrimeProducer(primes);
        TaskUtils.timeRun(producer, 1, SECONDS);
        LOGGER.debug("timeRun end");
    }
}
```

读者可以自行调试下，虽然 **timeRun** 能实现限时执行的功能，但它是通过外部线程安排中断实现。

在前面的 [《任务取消和线程中断》](../../../../../../2022/11/11/java/java-concurrency-learning/java-concurrency-learning12/)中我们了解到，每个线程都有自己的中断策略，在中断线程之前，应该了解它的中断策略，否则就不应该中断该线程。

由于 **timeRun** 可以从任意一个线程中调用，因此它无法知道这个调用线程的中断策略。

**如果任务在超时之前完成，会怎么样呢？**

下面我们再来看一下如下测试场景【任务在超时之前完成】：

```java
public class TimeRunTest {
    @Test
    public void timeRun1() {
        LOGGER.debug("timeRun start");
        TaskUtils.timeRun(new Runnable() {
            @Override
            public void run() {
                LOGGER.debug("task");
            }
        }, 400, TimeUnit.MILLISECONDS);
        try {
            LOGGER.debug("sleep start");
            SECONDS.sleep(1);
            LOGGER.debug("sleep end");
        } catch (InterruptedException e) {
            LOGGER.debug("InterruptedException");
        }
        LOGGER.debug("timeRun end");
    }
}
```

读者可以自行调式下，运行如下：

![](result.png)

上述示例中，任务在超时之前完成，而中断 **timeRun** 所在线程的取消任务将在 **timeRun** 返回到调用者之后启动。其中 `SECONDS.sleep(1);` 响应了中断，并抛出了 **InterruptedException** 异常，示例代码捕获该异常后打印了 **InterruptedException**。

虽然我们的任务在超时之前已经运行完了，但是取消任务在指定时间后还是对 **timeRun** 所在线程发出了中断请求。我们不知道在这种情况下 **timeRun** 返回之后调用者将运行什么代码【`SECONDS.sleep(1);` 这段只是为了演示】，但结果一定是不好的。【当然这里可以使用 **schedule** 返回的 **ScheduledFuture** 来取消这个取消任务以避免这种风险，这种做法虽然可行，但却非常复杂。】

**如果任务不响应中断，会怎么样呢？**

下面我们再来看一下如下测试场景【任务不响应中断请求】：

```java
public class TimeRunTest {
    @Test
    public void timeRun2() {
        LOGGER.debug("timeRun start");
        TaskUtils.timeRun(new PrimeGenerator(), 400, TimeUnit.MILLISECONDS);
        try {
            LOGGER.debug("sleep start");
            SECONDS.sleep(1);
            LOGGER.debug("sleep end");
        } catch (InterruptedException e) {
            LOGGER.debug("InterruptedException");
        }
        LOGGER.debug("timeRun end");
    }
}
```

上述示例中，素数生成器任务采用了自定义的取消策略，并没有响应中断，结果就是 **timeRun** 一直等待素数生成器任务结束，而它却永远不会结束。

如果任务不响应中断，那么 **timeRun** 会在任务结束时才返回，此时可能已经超过了指定的时限（或者还没有超过时限）。如果某个限时运行的服务没有在指定的时间内返回，那么将对调用者带来负面的影响。

下面我们来看一下如下示例【在专门的线程中中断任务】：

```java
public class TaskUtils {
    private static final ScheduledExecutorService cancelExec = Executors.newScheduledThreadPool(10);

    public static void timeRunNew(Runnable r, long timeout, TimeUnit unit) throws InterruptedException {
        class RethrowableTask implements Runnable {
            private volatile Throwable t;

            public void run() {
                try {
                    r.run();
                } catch (Throwable t) {
                    this.t = t;
                }
            }

            void rethrow() {
                if (null != t)
                    throw launderThrowable(t);
            }
        }

        RethrowableTask task = new RethrowableTask();
        final Thread taskThread = new Thread(task);
        taskThread.start();
//        cancelExec.schedule(new Runnable() {
//            public void run() {
//                taskThread.interrupt();
//            }
//        }, timeout, unit);
        LOGGER.debug("join start");
        // 线程 taskThread 至多等待指定毫秒后结束
        taskThread.join(unit.toMillis(timeout));
        LOGGER.debug("join end");
        task.rethrow();
    }
}    
```

上述示例中，执行任务的线程拥有自己的执行策略，即使任务不响应中断，限时运行的方法仍能返回到它的调用者。

在启动任务线程之后，**timeRun** 将执行一个限时的 **join** 方法。在 **join** 返回后，它将检查任务中是否有异常抛出，如果有的话，则会在调用 **timeRun** 的线程中再次抛出该异常。由于 **Throwable** 将在两个线程之间共享，因此该变量被声明为 **volatile** 类型，从而确保安全地将其从任务线程发布到 **timeRun** 线程。

虽然上述示例代码解决了前面的问题，但是由于它依赖一个限时的 **join**，因此存在着 **join** 的不足： **`无法知道执行控制是因为线程正常退出而返回，还是因为 join 超时而返回`** 。

> 这是 **Thread API** 的一个缺陷，因为无论 **join** 是否成功地完成，在 **Java** 内存模型中都会有内存可见性结果，但 **join** 本身不会返回某个状态来表明它是否成功。

# 2. 通过 Future 来实现取消

在前面的[《同步工具类（闭锁、信号量和栅栏）》](../../../../../../2022/09/17/java/java-concurrency-learning/java-concurrency-learning8/)博文中，咱们已经初步了解 **Future**，它可以管理任务的生命周期、处理异常以及实现取消。

而在另一篇[《任务执行演示》](../../../../../../2022/10/15/java/java-concurrency-learning/java-concurrency-learning11/)博文中，我们知道 `ExecutorService.submit` 将返回一个 **Future** 来描述任务。**Future** 拥有一个 **cancel** 方法，该方法带有一个 **boolean** 类型的参数 **mayInterruptIfRunning**，一个 **boolean** 类型的返回值。如果 **mayInterruptIfRunning** 为 **true** 并且任务当前正在某个线程中运行，那么这个线程能被中断。如果 **mayInterruptIfRunning** 为 **false**，则允许完成正在进行的任务，同时还未启动的任务也不再运行，这种方式适用于那些不处理中断的任务中。如果任务无法取消，则 **cancel** 方法返回 **false**，通常是因为任务已经正常完成；否则返回 **true**。

前文中我们一直强调，除非知道线程的中断策略，否则就不要中断线程。

那么使用 **Future** ，在什么情况下调用 **cancel** 可以将 **mayInterruptIfRunning** 参数指定为 **true** ？

执行任务的线程是由标准的 **Executor** 创建的，其实现了一种中断策略使得任务可以通过中断被取消。

当尝试取消某个任务时，不宜直接中断线程池，因为你并不知道当中断请求到达时正在运行什么任务--只能通过任务的 **Future** 来实现取消。

下面我们来看一下如下的示例【通过 **Future** 来取消任务】：

```java
public class TaskUtils {

    private static final ExecutorService taskExec = Executors.newCachedThreadPool();

    public static void timeRunByFuture(Runnable r, long timeout, TimeUnit unit) throws InterruptedException {
        Future<?> task = taskExec.submit(r);
        try {
            task.get(timeout, unit);
        } catch (ExecutionException e) {
            // 如果任务中抛出了异常，那么将重新抛出该异常，以便调用者处理异常
            throw launderThrowable(e.getCause());
        } catch (TimeoutException e) {
            // 任务超时，最终 finally 也会将任务取消
        } finally {
            // 如果任务已经结束，那么执行取消操作也不会带来任何影响
            // 如果任务正在运行，那么将被中断
            task.cancel(true);
        }
    }
}

```

上述示例应该很好理解，读者可以尝试跑下面的自测类来验证下。

```java
    /**
     * 任务运行中会响应中断请求
     */
    @Test
    public void timeRunByFuture() {
        LOGGER.debug("timeRun start");
        try {
            BlockingQueue<BigInteger> primes = new LinkedBlockingQueue<>(100);
            PrimeProducer producer = new PrimeProducer(primes);
            TaskUtils.timeRunByFuture(producer, 1, SECONDS);
        } catch (InterruptedException e) {
            LOGGER.debug("InterruptedException");
        }
        LOGGER.debug("timeRun end");
    }

    /**
     * 任务运行中不会响应中断请求
     */
    @Test
    public void timeRunByFuture1() {
        LOGGER.debug("timeRun start");
        try {
            TaskUtils.timeRunByFuture(new PrimeGenerator(), 500, TimeUnit.MILLISECONDS);
        } catch (InterruptedException e) {
            LOGGER.debug("InterruptedException");
        }
        LOGGER.debug("timeRun end");
    }

    /**
     * 任务超时之前完成
     */
    @Test
    public void timeRunByFuture2() {
        LOGGER.debug("timeRun start");
        try {
            TaskUtils.timeRunByFuture(new Runnable() {
                @Override
                public void run() {
                    LOGGER.debug("task");
                }
            }, 400, TimeUnit.MILLISECONDS);
        } catch (InterruptedException e) {
            LOGGER.debug("InterruptedException");
        }
        LOGGER.debug("timeRun end");
    }
```
# 3. 处理不可中断的阻塞

我们知道，为了方便开发人员构建出能响应取消请求的任务，在 **Java** 类库中的大多数可阻塞的方法都是通过提前返回或者抛出 **InterruptedException** 来响应中断请求的。

对于那些由于执行不可中断操作而被阻塞的线程，在知晓线程阻塞原因的前提下，我们也是可以使用类似中断的手段来停止这些线程。

- **java.io 包中的同步 Socket I/O**。在服务器应用程序中，最常见的阻塞 **I/O** 形式 就是对套接字进行读取和写入。虽然 **InputStream** 和 **OutputStream** 中的 **read** 和 **write** 等方法都不会响应中断，但通过关闭底层的套接字，可以使得由于执行 **read** 或 **write** 等方法而被阻塞的线程抛出一个 **SocketException**。

- **java.io 包中的同步 I/O**。当中断一个正在 **InterruptibleChannel**【可中断通道】上等待的线程时，将抛出 **ClosedByInterruptedException** 并关闭链路（这还会使得其他在这条链路上阻塞的线程同样抛出 **ClosedByInterruptedException**）。当关闭一个 **InterruptibleChannel** 时，将导致所有在链路操作上阻塞的线程抛出 **AsynchronousCloseException**。大多数标准的 **Channel** 都实现了 **InterruptibleChannel**。

- **Selector 的异步 I/O**。如果一个线程在调用 `Selector.select` 方法（在 `java.nio.channels` 中）时阻塞了，那么调用 **close** 或 **wakeup** 方法会使线程抛出 **ClosedSelectorException** 并提前返回。

- **获取某个锁**。如果一个线程由于等待某个内置锁而阻塞，那么将无法响应中断，因为线程认为它肯定会获得锁，所以将不会理会中断请求。不过，在 **Lock** 类中提供了 **lockInterruptibly** 方法，它允许在等待一个锁的同时仍能响应中断。

下面我们来看一下如下示例【通过改写 **interrput** 方法将非标准的取消操作封装在 **Thread** 中】：

```java
public class ReaderThread extends Thread {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(ReaderThread.class);

    private final Socket socket;

    private final InputStream in;

    public ReaderThread(Socket socket) throws IOException {
        this.socket = socket;
        this.in = socket.getInputStream();
    }

    @Override
    public void interrupt() {
        LOGGER.debug("interrupt");
        try {
            socket.close();
            LOGGER.debug("socket close");
        } catch (IOException e) {
            //
        } finally {
            super.interrupt();
        }
    }

    @Override
    public void run() {
        try {
            InputStreamReader inputStreamReader = new InputStreamReader(in);
            BufferedReader bufferedReader = new BufferedReader(inputStreamReader);
            String data;
            while ((data = bufferedReader.readLine()) != null) {
                processData(data);
            }
        } catch (IOException e) {
            // 允许线程退出
        }
    }

    /**
     * 输出 0 ~ data 区间内的素数
     */
    private void processData(String data) {
        LOGGER.debug("0 < All Primes < {}", data);
        BigInteger prime = BigInteger.ONE;
        while (!Thread.currentThread().isInterrupted() && prime.compareTo(BigInteger.valueOf(Long.valueOf(data))) < 0) {
            LOGGER.debug("prime = {}", prime);
            prime = prime.nextProbablePrime();
        }
    }
}
```

上述 **ReaderThread** 管理了一个套接字连接，它采用同步方式从该套接字中读取数据，并将接收到的数据传递给 **processData**。同时由于 **ReaderThread** 改写了 **interrupt** 方法，使其既能处理标准的中断，也能关闭底层的套接字。

感兴趣的读者，可以自行测试如下【先启动 **SocketServer** ，再运行 **SocketClient** 】：

```java
/**
 * Socket服务端
 */
public class SocketServer {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(SocketServer.class);

    private static final ScheduledExecutorService cancelExec = Executors.newScheduledThreadPool(10);

    public static void main(String[] args) throws IOException {
        // 创建服务端socket
        ServerSocket serverSocket = new ServerSocket(8888);

        //循环监听等待客户端的连接
        while (true) {
            // 监听客户端
            LOGGER.debug("start serverSocket.accept()");
            // 创建客户端socket
            Socket socket = serverSocket.accept();
            LOGGER.debug("end serverSocket.accept()");

            ReaderThread readerThread = new ReaderThread(socket);
            readerThread.start();
            
            // 演示 2s后中断 ReaderThread
            cancelExec.schedule(new Runnable() {
                public void run() {
                    readerThread.interrupt();
                }
            }, 2, SECONDS);

        }
    }
}

/**
 * Socket客户端
 */
public class SocketClient {

    public static void main(String[] args) throws IOException {
        // 和服务器创建连接
        Socket socket = new Socket("localhost", 8888);

        // 要发送给服务器的信息
        OutputStream os = socket.getOutputStream();
        PrintWriter pw = new PrintWriter(os);
        pw.write("1000000\n" +
                "10000");
        pw.flush();
        socket.shutdownOutput();

        pw.close();
        os.close();
        socket.close();
    }
}
```

**Socket** 服务端启动后，执行 **Socket** 客户端，笔者 **Socket** 服务端运行结果如下【以实际运行为准】：

![](result-1.png)

# 4. 采用 newTaskFor 来封装非标准的取消

我们可以通过 **Java 6** 在 **ThreadPoolExecutor** 中新增的 **newTaskFor** 方法来进一步优化 **ReaderThread** 中封装非标准取消的技术。

当把一个 **Callable** 提交给 **ExecutorService** 时，**submit** 方法会返回一个 **Future**，我们可以使用这个 **Future** 来取消任务。

**newTaskFor** 是一个工厂方法，它将创建 **Future** 来代表任务。 **newTaskFor** 还能返回一个 **RunnableFuture** 接口，该接口扩展了 **Future** 和 **Runnable**（并由 **FutureTask** 实现）。

通过定制表示任务的 **Future** 可以改变 `Future.cancel` 的行为。定制的取消代码可以实现，例如：

- 日志记录
- 收集取消操作的统计信息
- 取消一些不响应中断的操作

下面我们来看一下如下示例【通过 **newTaskFor** 将非标准的取消操作封装到一个任务中】：

我们首先定义了一个 **CancellableTask** 接口，该接口扩展了 **Callable**，其中增加了一个 **取消方法** 和一个 **newTask** 工厂方法来构造 **RunnableFuture**。

```java
public interface CancellableTask<T> extends Callable<T> {
    void cancel();

    RunnableFuture<T> newTask();
}
```

然后我们定义抽象类 **SocketUsingTask** ，它实现了 **CancellableTask**，并通过 `Future.cancel` 来关闭套接字和调用 `super.cancel`。如果 **SocketUsingTask** 通过其自己的 **Future** 来取消，那么底层的套接字将被关闭并且线程将被中断。

```java
public abstract class SocketUsingTask<T> implements CancellableTask<T> {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(SocketUsingTask.class);

    @GuardedBy("this")
    private Socket socket;

    public SocketUsingTask(Socket socket) {
        this.socket = socket;
    }

    protected synchronized Socket getSocket() {
        return socket;
    }

    public synchronized void cancel() {
        LOGGER.debug("start custom cancel");
        try {
            if (socket != null) {
                socket.close();
                LOGGER.debug("socket close");
            }
        } catch (IOException e) {
            //
        }
        LOGGER.debug("end custom cancel");
    }

    public RunnableFuture<T> newTask() {
        return new FutureTask<T>(this) {
            @Override
            public boolean cancel(boolean mayInterruptIfRunning) {
                LOGGER.debug("start cancel");
                SocketUsingTask.this.cancel();
                boolean result = super.cancel(mayInterruptIfRunning);
                LOGGER.debug("end cancel");
                LOGGER.debug("cancel result = {}", result);
                return result;
            }
        };
    }
}
```

紧接着，我们定义 **CancellingExecutor** ，它扩展了 **ThreadPoolExecutor**，并通过改写 **newTaskFor** 使得 **CancellableTask** 可以创建自己的 **Future**。

```java
@ThreadSafe
public class CancellingExecutor extends ThreadPoolExecutor {

    public CancellingExecutor() {
        super(0, Integer.MAX_VALUE, 60L, TimeUnit.SECONDS, new SynchronousQueue<>());
    }

    @Override
    protected <T> RunnableFuture<T> newTaskFor(Callable<T> callable) {
        if (callable instanceof CancellableTask)
            return ((CancellableTask<T>) callable).newTask();
        else
            return super.newTaskFor(callable);
    }
}
```




最后，我们定义了任务类 **PrimeSumTask** ，它继承了上面的抽象类 **SocketUsingTask**，call 方法用于计算指定范围内的 **素数总和** ，如下：

```java
public class PrimeSumTask extends SocketUsingTask<BigInteger> {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(PrimeSumTask.class);

    public PrimeSumTask(Socket socket) {
        super(socket);
    }

    @Override
    public BigInteger call() {
        BigInteger result = null;
        try {
            InputStream in = getSocket().getInputStream();
            InputStreamReader inputStreamReader = new InputStreamReader(in);
            BufferedReader bufferedReader = new BufferedReader(inputStreamReader);
            String data;
            while ((data = bufferedReader.readLine()) != null) {
                result = processData(data);
            }
        } catch (IOException e) {
            // 允许线程退出
        }
        return result;
    }

    /**
     * 计算 0 ~ data 区间内的素数总和
     */
    private BigInteger processData(String data) {
        LOGGER.debug("0 < All Primes < {}", data);
        BigInteger prime = BigInteger.ONE;
        BigInteger sum = BigInteger.ZERO;
        while (!Thread.currentThread().isInterrupted() && prime.compareTo(BigInteger.valueOf(Long.valueOf(data))) < 0) {
            sum = sum.add(prime);
            prime = prime.nextProbablePrime();
        }
        return sum;
    }
}
```

感兴趣的读者，可以自行测试如下【先启动 **SocketServer** ，再运行 **SocketClient** 】：

```java
public class SocketServer {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(SocketServer.class);

    private static CancellingExecutor executor = new CancellingExecutor();

    public static void main(String[] args) throws IOException {
        // 创建服务端socket
        ServerSocket serverSocket = new ServerSocket(8888);

        //循环监听等待客户端的连接
        while (true) {
            // 监听客户端
            LOGGER.debug("start serverSocket.accept()");
            // 创建客户端socket
            Socket socket = serverSocket.accept();
            LOGGER.debug("end serverSocket.accept()");

            PrimeSumTask primeSumTask = new PrimeSumTask(socket);
            Future<BigInteger> future = executor.submit(primeSumTask);

            try {
                BigInteger result = future.get(2, TimeUnit.SECONDS);
                LOGGER.debug("result = {}", result);
            } catch (ExecutionException e) {
                // 如果任务中抛出了异常，那么重新抛出该异常
                throw launderThrowable(e.getCause());
            } catch (TimeoutException e) {
                // 任务超时，最终 finally 也会将任务取消
                LOGGER.error("TimeoutException");
            } catch (InterruptedException e) {
                // 中断异常
            } finally {
                // 如果任务已经结束，那么执行取消操作也不会带来任何影响
                // 如果任务正在运行，那么将被中断
                LOGGER.debug( "task is done : {}", future.isDone());
                LOGGER.debug( "future cancel start");
                future.cancel(true);
                LOGGER.debug( "future cancel end");
                LOGGER.debug( "task is cancelled : {}", future.isCancelled());
            }

        }
    }
}

public class SocketClient {

    public static void main(String[] args) throws IOException {
        // 和服务器创建连接
        Socket socket = new Socket("localhost", 8888);

        // 要发送给服务器的信息
        OutputStream os = socket.getOutputStream();
        PrintWriter pw = new PrintWriter(os);
        pw.write("1000000");
        pw.flush();
        socket.shutdownOutput();

        pw.close();
        os.close();
        socket.close();
    }
}
```
**Socket** 服务端启动后，执行 **Socket** 客户端，笔者 **Socket** 服务端运行结果如下：

![](result-2.png)
上面场景是任务超时运行，接下来我们调整 `future.get` 的超时时间为 **5s**， 如下所示：

```java
    BigInteger result = future.get(5, TimeUnit.SECONDS);
```

再重新执行 **Socket** 客户端，此时运行结果如下：

![](result-3.png)

# 5. 总结
《任务取消》的内容已告一段落，下篇开始介绍各种任务和服务的关闭机制，以及如何编写任务和服务，使它们能够优雅地处理关闭。