---
title: Java并发编程学习8-同步工具类
date: 2022-09-17 15:16:21
updated: 2024-03-19 17:21:35
categories:
    - [开发语言-Java,Java并发编程]
tags:
    - Java
    - 同步工具类
    - 闭锁
    - FutureTask
    - 信号量
    - 栅栏
---

[《开发语言-Java》](/categories/开发语言-Java/) [《Java并发编程》](/categories/开发语言-Java/Java并发编程/) 

![](/images/java-concurrency-logo.png)

# 引言

同步工具类可以是任意一个对象，只要它根据其自身的状态来协调线程的控制流。阻塞队列可以作为同步工具类，类似地还有**信号量（Semaphore）**、**栅栏（Barrier）**以及**闭锁（Latch）**。当然 **Java** 平台类库中还有其他的一些同步工具类，如果这些都不能满足要求，那我们还可以创建自己的同步工具类【这块内容将在后续的博文中会介绍】。

同步工具类封装了一些状态，这些状态将决定执行同步工具类的线程是继续执行还是等待，此外还提供了一些方法对状态进行操作，以及另一些方法用于高效地等待同步工具类进入到预期状态。

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 主要内容
## 1. 闭锁

**闭锁** 是一种同步工具类，它可以延迟线程的进度直到其到达终止状态。

闭锁的作用相当于一扇门：在闭锁到达结束状态之前，这扇门一直是关闭的，并且没有任何线程都能通过，当到达结束状态时，这扇门会打开并允许所有的线程通过。当闭锁到达结束状态后，将不会再改变状态，因此这扇门将永远保持打开状态。

闭锁可以用来确保某些活动直到其他活动都完成后才继续执行，例如：

- 确保某个计算在其需要的所有资源都被初始化之后才继续执行。二元闭锁（包括两个状态）可以用来表示“资源 **R** 已经被初始化”，而所有需要R的操作都必须先在这个闭锁上等待。
- 确保某个服务在其依赖的所有其他服务都已经启动之后才启动。每个服务都有一个相关的二元闭锁。当启动服务 **S** 时，将首先在 **S** 依赖的其他服务的闭锁上等待，在所有依赖的服务都启动后会释放闭锁 **S**，这样其他依赖 **S** 的服务才能继续执行。
- 等待直到某个操作的所有参与者（如，在多玩家游戏中的所有玩家）都就绪，再继续执行。


**CountDownLatch** 是一种灵活的闭锁实现，可以在上诉情况下使用。它包含一个计数器，该计数器被初始化为一个正数，表示需要等待的事件数量。**countDown** 方法递减计数器，表示已经有一个事件发生了，而 **await** 方法等待计数器达到零，这表示所有需要等待的事件都已经发生。如果结束门计数器的值为非零，那么它的 **await** 方法会一直阻塞直到计数器的值为零，或者等待中的线程中断，或者等待超时。

下面我们来看一下如下的示例，**TestHarness** 中给出了闭锁的两种常见用法：

```java
public class TestHarness {
    
    public long timeTasks(int nThreads, final Runnable task) throws InterruptedException {
    
        final CountDownLatch startGate = new CountDownLatch(1);
    
        final CountDownLatch endGate = new CountDownLatch(nThreads);
    
        for (int i = 0; i < nThreads; i++) {
            Thread t = new Thread() {
                public void run() {
                    try {
                        startGate.await();
                        try {
                            task.run();
                        } finally {
                            endGate.countDown();
                        }
                    }catch (InterruptedException ignored) {
                    }
                }
            };
            t.start();
        }
    
        long start = System.nanoTime();
        startGate.countDown();
        endGate.await();
        long end = System.nanoTime();
        return end - start;
    }
}
```

**TestHarness**使用了两个闭锁，分别表示“**起始门（Starting Gate）**”和“**结束门（Ending Gate）**”。
起始门计数器的初始值为 **1**，而结束门计数器的初始值为**工作线程的数量**。每个工作线程首先要做的就是在起始门上等待，从而确保所有的线程都就绪后才开始执行。而每个线程要做的最后一件事情是将调用结束门的 **countDown** 方法减 **1**，这能使主线程高效地等待直到所有工作线程都执行完成，因此可以统计所消耗的时间。

## 2. FutureTask

**FutureTask** 也可以用作闭锁。它实现了 **Future** 语义，其表示的计算是通过 **Callable** 来实现的，相当于一种可生成结果的 **Runnable**，并且可以处于以下 **3** 种状态：

- **等待运行（Waiting to run）**
- **正在运行（Running）**
- **运行完成（Completed）**

“**运行完成**” 表示计算的所有可能的结束方式，包括正常结束、由于取消而结束和由于异常而结束等。当 **FutureTask** 进入完成状态后，它会永远停止在这个状态上。

**Future.get** 的行为取决于任务的状态。如果任务已经完成，那么 **get** 方法会立即返回结果，否则 **get** 将阻塞直到任务进入完成状态，然后返回结果或者抛出异常。如果任务抛出了异常，那么 **get** 将该异常封装为 **ExecutionException** 并重新抛出，这时可以通过 **getCause** 来获得被封装的初始异常。如果任务被取消，那么 get 将抛出 **CancellationException**。 

**FutureTask** 在 **Executor** 框架中表示异步执行，此外还可以用来表示一些时间较长的计算，这些计算可以在使用计算结果之前启动。**FutureTask** 将计算结果从执行计算的线程传递到获取这个结果的线程，而 **FutureTask** 的规范确保了这种传递过程能实现结果的安全发布。

**下面来看一个示例：**

```java
public class Preloader {
    private final FutureTask<ProductInfo> future =
            new FutureTask<>(new Callable<ProductInfo>() {
                public ProductInfo call() throws DataLoadException {
                    return ProductInfoUtils.loadProductInfo();
                }
            });

    private final Thread thread = new Thread(future);

    public void start() {
        thread.start();
    }

    public ProductInfo get() throws DataLoadException, InterruptedException {
        try {
            return future.get();
        } catch (ExecutionException e) {
            Throwable cause = e.getCause();
            if (cause instanceof DataLoadException)
                throw (DataLoadException) cause;
            else
                throw ExceptionUtils.launderThrowable(cause);
        }
    }
}
```

上述 **Preloader** 创建了一个 **FutureTask**，其中包含从数据库加载产品信息的任务，以及一个执行运算的线程，同时提供了一个 **start** 方法来启动线程。当程序随后需要 **ProductInfo** 时，可以调用 **get** 方法，如果数据已经加载，那么将返回这些数据，否则将等待加载完成后再返回。

**Callable** 表示的任务可以抛出受检查的或未受检查的异常，并且任何代码都可能抛出一个 **Error**。无论任何代码抛出什么异常，都会被封装到一个 **ExecutionException** 中，并在 **Future.get** 中被重新抛出。

在 **Preloader** 中，当 **get** 方法抛出 **ExecutionException** 时，可能是以下三种情况之一：
- **Callable** 抛出的受检查异常
- **RuntimeException**
- **Error**

我们可以看到，**Preloader** 会首先检查已知的受检查异常，并重新抛出它们，这里是 **DataLoadException**。如果不是已知的受检查异常，将调用 **launderThrowable** 并抛出结果。

下面我们来看一下 **launderThrowable** 的代码：

```java
public static RuntimeException launderThrowable(Throwable t) {
    if (t instanceof RuntimeException)
        return (RuntimeException) t;
    else if (t instanceof Error) 
        throw (Error) t;
    else 
        throw new IllegalStateException("Not unchecked", t);
}
```

## 3. 信号量

信号量中管理着一组虚拟的许可，许可的初始数量可通过构造函数来指定。在执行操作时可以首先获取许可（只要还有剩余的许可），并在使用以后释放许可。如果没有许可，那么信号量的 **acquire** 方法将阻塞直到有许可（或者直到被中断或者操作超时）。**release** 方法将返回一个许可给信号量。

**计数信号量（Counting Semaphore）** 可以用来控制同时访问某个特定资源的操作数量，或者同时执行某个指定操作的数量。它还可以用来实现某种资源池，或者对容器施加边界。计数信号量的一种简化形式是 **二值信号量**，即初始值为 **1** 的信号量。二值型号量可以用作互斥体，并具备不可重入的加锁语义：谁拥有这个唯一的许可，谁就拥有了互斥锁。

> 在上述二值信号量的实现中，不包含真正的许可对象，并且信号量也不会将许可与线程关联起来，因此在一个线程中获得的许可可以在另一个线程中释放。可以将 **acquire** 操作视为是消费一个许可，而 **release** 操作是创建一个许可，信号量并不受限于它在创建时的初始许可数量。

信号量可以将任何一种容器变成有界阻塞容器，下面我们来看一下示例：

```java
public class BoundedHashSet<T> {
    private final Set<T> set;
    
    private final Semaphore sem;
    
    public BoundedHashSet(int bound) {
        this.set = Collections.synchronizedSet(new HashSet<T>());
        sem = new Semaphore(bound);
    }
    
    public boolean add(T t) throws InterruptedException {
        sem.acquire();
        boolean wasAdded = false;
        try {
            wasAdded = set.add(t);
            return wasAdded;
        } finally {
            if (!wasAdded)
                sem.release();
        }
    }
    
    public boolean remove(Object t) {
        boolean wasRemoved = set.remove(t);
        if (wasRemoved) 
            sem.release();
        return wasRemoved;
    }
}
```

上述 **BoundedHashSet** 中，信号量的计数值会初始化为容器容量的最大值。**add** 操作在向底层容器中添加一个元素之前，首先要获取一个许可。如果 **add** 操作没有添加任何元素，那么会立即释放许可。**remove** 操作释放一个许可，使更多的元素能够添加到容器中。当然底层的 **Set** 实现并不知道关于边界的任何信息，这是由 **BoundedHashSet** 来控制的。

## 4. 栅栏

**栅栏（Barrier）** 类似于闭锁，它能阻塞一组线程直到某个事件发生。栅栏与闭锁的关键区别在于，所有线程必须同时到达栅栏位置，才能继续执行。闭锁用于等待事件，而栅栏用于等待其他线程。

**CyclicBarrier** 可以使一定数量的参与方反复地在栅栏位置汇集，它在并行迭代算法中非常有用：这种算法通常将一个问题拆分成一系列相互独立的子问题。当线程到达栅栏位置时将调用 **await** 方法，这个方法将阻塞直到所有线程都到达栅栏位置。如果所有线程都到达了栅栏位置，那么栅栏将打开，此时所有线程都被释放，而栅栏将被重置以便下次使用。如果对 **await** 的调用超时，或者 **await** 阻塞的线程被中断，那么栅栏就被认为是打破了，所有阻塞的 **await** 调用都将终止被抛出 **BrokenBarrierException**。如果成功地通过栅栏，那么 **await** 将为每个线程返回一个唯一的到达索引号，我们可以利用这些索引来 “选举” 产生一个领导线程，并在下一次迭代中由该领导线程执行一些特殊的工作。

**CyclicBarrier** 还可以将一个栅栏操作传递给构造函数，该栅栏操作是一个 **Runnable**，当成功通过栅栏时会（在一个子任务线程中）执行它，但在阻塞线程被释放之前是不能执行的。

在一些模拟程序中通常需要使用栅栏，例如某个步骤中的计算可以并行执行，但必须等到该步骤中的所有计算都执行完毕才能进入下一个步骤。

下面我们来看一个示例【如下的 **CellularAutomata** 演示了如何通过栅栏来计算细胞的自动化模拟】：

```java
public class CellularAutomata {
    private final Board mainBroad;

    private final CyclicBarrier barrier;

    private final Worker[] workers;

    public CellularAutomata(Board board) {
        this.mainBroad = board;
        int count = Runtime.getRuntime().availableProcessors();
        // 栅栏的构造参数可以传入一个 Runnable 的匿名内部类
        this.barrier = new CyclicBarrier(count, new Runnable(){
            public void run() {
                mainBroad.commitNewValues();
            }
        });
        this.workers = new Worker[count];
        for (int i = 0; i < count; i++)
            workers[i] = new Worker(mainBroad.getSubBoard(count, i));
    }

    private class Worker implements Runnable {
        private final Board board;

        public Worker(Board board) {
            this.board = board;
        }

        public void run() {
            while (!board.hasConverged()) {
                for (int x = 0; x < board.getMaxX(); x++)
                    for (int y = 0; y < board.getMaxY(); y++)
                        board.setNewValue(x, y, computeValue(x, y));
                try {
                    barrier.await();
                } catch (InterruptedException | BrokenBarrierException ex) {
                    return;
                }
            }
        }

        private Board computeValue(int x, int y) {
            x = 2 * x + y;
            y = 2 * y + x;
            return new Board(x, y);
        }
    }

    public void start() {
        for (int i = 0; i < workers.length; i++)
            new Thread(workers[i]).start();
        mainBroad.waitForConvergence();
    }
}
```

上述 **CellularAutomata** 将问题分解为 $N_{cpu}$ 个子问题，其中 $N_{cpu}$ 等于可用 **CPU** 的数量，并将每个子问题分配给一个线程。在每个步骤中，工作线程都为各自子问题中的所有细胞计算新值。当所有工作线程都到达栅栏时，栅栏会把这些新值提交给数据模型。在栅栏的操作执行完以后，工作线程将开始下一步的计算，包括调用 **isDone** 方法来判断是否需要进行下一次迭代。

> 在这种不涉及 **I/O** 操作或共享数据访问的计算问题中，当线程数量为 $N_{cpu}$ 或 $N_{cpu}+1$ 时将获得最优的吞吐量。更多的线程并不会带来任何帮助，甚至在某种程度上会降低性能，因为多个线程将会在 **CPU** 和 **内存** 等资源上发生竞争。


另一种形式的栅栏是 **ExChanger**，它是一种 **两方（Two-Party）栅栏**，各方在栅栏位置上交换数据。当两方执行不对称的操作时，**ExChanger** 会非常有用。

> 例如，当一个线程向缓冲区写入数据，而另一个线程从缓冲区读取数据。这些线程可以使用 **Exchanger** 来汇合，并将满的缓冲区与空的缓冲区交换。当两个线程通过 **Exchanger** 交换对象时，这种交换就把这两个对象安全地发布给另一方。


# 总结
本篇介绍了 **Java** 平台类库中的一些常用的同步工具类，到目前为止，我们已经学到了很多的基础知识。下一篇博文将介绍如何利用前面学到的基础知识，构建一个高效且可伸缩的缓存，用于改进一个高计算开销的函数。

