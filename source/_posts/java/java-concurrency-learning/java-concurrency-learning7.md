---
title: Java并发编程学习7-阻塞队列
date: 2022-09-13 08:45:00
updated: 2024-03-19 17:21:35
categories:
  - [开发语言-Java,Java并发编程]
tags:
  - Java
  - Queue
  - BlockingQueue
  - Deque
  - BlockingDeque
  - 阻塞队列
  - 工作密取
---

[《开发语言-Java》](/categories/开发语言-Java/) [《Java并发编程》](/categories/开发语言-Java/Java并发编程/) 

![](/images/java-concurrency-logo.png)


# 引言

介绍阻塞队列之前，先来介绍下队列 **Queue**。**Queue** 用来临时保存一组等待处理的元素。它提供了几种非阻塞队列实现，如下：

- **ConcurrentLinkedQueue**，这是一个传统的先进先出队列。
- **PriorityQueue**，这是一个（非并发的）优先队列。

如上两个队列的操作不会阻塞，如果队列为空，那么获取元素的操作将返回空值。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

阻塞队列 **BlockingQueue** 扩展了 **Queue**，增加了可阻塞的 **put** 和 **take** 方法，以及支持定时的 **offer** 和 **poll** 方法。如果队列已经满了，那么 **put** 方法将阻塞直到有空间可用；如果队列为空，那么 **take** 方法将会阻塞直到有元素可用。队列可以是有界的也可以是无界的，无界队列永远都不会充满，因此无界队列上的 **put** 方法也永远不会阻塞。

阻塞队列支持 **生产者--消费者** 这种设计模式。当数据生成时，生产者把数据放入队列，而当消费者准备处理数据时，将从队列中获取数据。一种最常见的 **生产者--消费者** 设计模式就是线程池与工作队列的组合，在 **Executor** 任务执行框架中就体现了这种模式，这也是后面的博文中将要介绍的内容。

在 **Java** 类库中包含了 **BlockingQueue** 的多种实现，如下：

- **LinkedBlockingQueue** 和 **ArrayBlockingQueue** 是 **FIFO** 队列，二者分别与 **LinkedList** 和 **ArrayList** 类似，但比同步 **List** 拥有更好的并发性能。
- **PriorityBlockingQueue** 是一个按优先级排序的队列，它既可以根据元素的自然顺序来比较元素（前提是这些元素实现了**Comparable**方法），也可以使用 **Comparator** 来比较。
- **SynchronousQueue** ，实际上不能算一个队列，因为它不会为队列中元素维护存储空间。与其他队列不同的是，它维护一组线程，这些线程在等待着把元素加入或移出队列。因为 **SynchronousQueue** 没有存储功能，因此 **put** 和 **take** 会一直阻塞，直到有另一个线程已经准备好参与到交付过程中。仅当有足够多的消费者，并且总是有一个消费者准备好获取交付的工作时，才适合使用同步队列。 

# 主要内容
## 1. “桌面搜索” 示例

如下 **FileCrawler** 中给出了一个生产者任务，即在某个文件层次结构中搜索符合索引标准的文件，并将它们的名称放入工作队列。

```java
    public class FileCrawler implements Runnable {
        private final BlockingQueue<File> fileQueue;
        
        private final FileFilter fileFilter;
        
        private final File root;
        
        public FileCrawler(BlockingQueue<File> fileQueue, FileFilter fileFilter, File root) {
            this.fileQueue = fileQueue;
            this.fileFilter = fileFilter;
            this.root = root;
        }
        
        public void run() {
            try {
                crawl(root);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }    
        }
        
        private void crawl(File root) throws InterruptedException {
            File[] entries = root.listFiles(fileFilter);
            if (entries != null) {
                for (File entry : entries) 
                    if (entry.isDirectory()) 
                        crawl(entry);
                    else if (!FileRecord.alreadyIndexed(entry)) 
                        fileQueue.put(entry);
            }
        }
    }
```

如下 **Indexer** 中给出了一个消费者任务，即从队列中取出文件名称并对它们建立索引，它会一直运行下去。

```java
public class Indexer implements Runnable {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(Indexer.class);

    private final BlockingQueue<File> queue;

    public Indexer(BlockingQueue<File> queue) {
        this.queue = queue;
    }

    public void run() {
        try {
            while(true) {
                File file = queue.take();
                FileRecord.indexFile(file);
                if (LOGGER.isDebugEnabled()) {
                    LOGGER.debug(file.getAbsolutePath());
                }
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }
}
```

**生产者--消费者** 模式提供了一种适合线程的方法将桌面搜索问题分解为更简单的组件。将文件遍历与建立索引等功能分解为独立的操作，每个操作只需完成一个任务，并且阻塞队列将负责所有的控制流程，因此每个功能的代码都更加简单和清晰。

下面我们再看一个测试代码示例，用于启动桌面搜索。

```java
public class FileCrawlerTest {

    private static final int BOUND = 1000;

    private static final int N_CONSUMERS = 5;

    public static void main(String[] args) throws Exception {
        File file = new File("E:\\fleaworkspace");
        File file1 = new File("E:\\Software\\Maven\\Repository");
        File[] roots = {file, file1};
        startIndexing(roots);
    }

    private static void startIndexing(File[] roots) {
        BlockingQueue<File> queue = new LinkedBlockingQueue<>(BOUND);
        FileFilter fileFilter = new FileFilter() {
            public boolean accept(File file) {
                return true;
            }
        };

        for (File root : roots)
            new Thread(new FileCrawler(queue, fileFilter, root)).start();

        for (int i = 0; i < N_CONSUMERS; i++)
            new Thread(new Indexer(queue)).start();
    }
}
```

这里启动了多个文件搜索程序和索引简历程序，每个程序都在各自的线程中运行。前面讲到，消费者线程永远不会退出，因而程序无法终止，在后续的博文将介绍多种技术来解决这个问题。

## 2. 串行线程封闭

在 **java.util.concurrent** 中实现的各种阻塞队列都包含了足够的内部同步机制，从而安全地将对象从生产者线程发布到消费者线程。

对于可变对象，**生产者--消费者** 这种设计与阻塞队列一起，促进了串行线程封闭，从而将对象所有权从生产者交付给消费者。线程封闭对象只能由单个线程拥有，但可以通过安全地发布该对象来 “转移” 所有权。在转移所有权后，也只有另一个线程能获得这个对象的的访问权限，并且发布对象的线程不会再访问它。这种安全的发布确保了对象状态对于新的所有者来说是可见的，并且由于最初的所有者不会再访问它，因此对象被封闭在新的线程中。新的所有者线程可以对该对象做任意修改，因为它具有独占的访问权。

> 对象池利用了串行线程封闭，将对象“借给”一个请求线程。只要对象池包含足够的内部同步来安全地发布池中的对象，并且只要客户代码本身不会发布池中的对象，或者在将对象返回给对象池后就不再使用它，那么就可以安全地在线程之间传递所有权。

## 3. 双端队列与工作密取

**Java 6** 增加两种容器类型，**Deque** 和 **BlockingDeque**，他们分别对 **Queue** 和 **BlockingQueue** 进行了扩展。

**Deque** 是一个双端队列，实现了在队列头和队列尾的高效插入和移除，具体实现包括：

- **ArrayDeque**
- **LinkedBlockingDeque**

正如阻塞队列适用于 **生产者--消费者** 模式，双端队列同样适用另一种相关模式，即 **工作密取（Work Stealing**）。

在生产者--消费者模式中，所有消费者有一个共享的工作队列，而在工作密取设计中，每个消费者都有各自的双端队列。如果一个消费者完成了自己双端队列中的全部工作，那么它可以从其他消费者双端队列末尾秘密地获取工作。工作密取模式比传统的生产者--消费这模式具有更高的可伸缩性。在大多数时候，它们都只是访问自己的双端队列，从而极大地减少了竞争。当工作线程需要访问另一个队列时，它会从队列的尾部而不是从头部获取工作，因此进一步降低了队列上的竞争程度。

工作密取非常适用于既是消费者也是生产者问题---当执行某个工作时可能导致出现更多的工作。例如网页爬虫处理页面、搜索图的算法、在垃圾回收阶段对堆进行标记等。当一个工作线程找到新的任务单元时，它会将其放到自己队列的末尾（或者放入其他工作线程的队列中）。当双端队列为空时，它会在另一个线程的队列队尾查找新的任务，从而确保每个线程都保持忙碌状态。


## 4. 阻塞方法与中断方法

线程可能会阻塞或暂停执行，原因有多种：等待I/O操作结束，等待获得一个锁，等待从 **Thread.sleep** 方法中醒来，或是等待另一个线程的计算结果。当线程阻塞时，它通常被挂起，并处于某种阻塞状态（**BLOCKED**、**WAITING** 或 **TIME_WAITING**）。

**BlockingQueue** 的 **put** 和 **take** 等方法会抛出受检查异常 **InterruptedException**，这与类库中其他一些方法的做法相同，例如 **Thread.sleep**，当某方法抛出 **InterruptedException** 时，表示该方法是一个阻塞方法，如果这个方法被中断，那么它将努力提前结束阻塞状态。**Thread** 提供了 **interrupt** 方法，用于中断线程或者查询线程是否已经被中断。每个线程都有一个布尔类型的属性，表示线程的中断状态，当中断线程时将设置这个状态。

中断是一种协作机制。当线程 **A** 中断 **B** 时，**A** 仅仅是要求 **B** 在执行到某个可以暂停的地方停止正在执行的操作（当然前提是如果线程 **B** 愿意停止下来）。最常使用中断的情况就是取消某个操作，如果程序对中断请求的响应度越高，就越容易及时取消那些执行时间很长的操作。

当在代码中调用了一个将抛出 **InterruptedException** 的方法时，自身方法也就变成了一个阻塞方法，并且必须要处理对中断的响应。

这里有两种常见的方法：

- 传递 **InterruptedException**，只需要把 **InterruptedException** 传递给方法的调用者，要么根本不捕获异常，或者捕获该异常，然后在执行某种简单的清理工作后再次抛出这个异常。
- 恢复中断，当代码是 **Runnable** 的一部分时，在这种情况下必须捕获 **InterruptedException**，并通过调用当前线程上的 **interrupt** 方法恢复中断状态，这样在调用栈中更高层的代码将看到引发了一个中断。

下面看下恢复中断状态的示例：

```java
    public class TaskRunnable implements Runnable {
        BlockingQueue<Task> queue;
        
        // ...
        
        public void run() {
            try {
                processTask(queue.take());
            } catch (InterruptedException e) {
                // 恢复被中断的状态
                Thread.currentThread().interrupt();
            }
        }
    }
```

# 总结

当然还可以采用一些更复杂的中断处理方法，但上述两种方法已经可以应对大多数情况了。关于取消和中断等操作，这里只是简单提及，笔者将会在后续的博文中进一步介绍，敬请期待！！！



