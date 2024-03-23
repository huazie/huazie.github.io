---
title: Java并发编程学习2-线程安全性
date: 2021-03-01 11:51:51
updated: 2024-03-18 14:22:22
categories:
  - [开发语言-Java,Java并发编程]
tags:
  - Java
  - 并发编程
  - 线程安全性
  - 竞态条件
  - 内置锁
  - 重入
---

[《开发语言-Java》](/categories/开发语言-Java/) [《Java并发编程》](/categories/开发语言-Java/Java并发编程/) 

![](/images/java-concurrency-logo.png)

# 引言
上篇我们初步了解了线程相关的知识，这篇我们深入了解下线程安全性的相关问题。

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 1. 什么是线程安全性？
线程安全性是一个在代码上使用的术语，它与对象或整个程序的状态相关的，只能应用于封装其状态的整个代码之中。在线程安全性的定义中，最核心的概念就是正确性。正确性的含义是，某个类的行为与其规范完全一致。当多个线程访问某个类时，不管运行时环境采用何种调度方式或者这些线程将如何交替执行，并且在主调代码中不需要任何额外的同步或协同，这个类始终都能表现出正确的行为，那么就称这个类是线程安全的。

## 1.1 如何编写线程安全的代码？
要想编写线程安全的代码，其核心在于要对状态访问操作进行管理，特别是对 **共享**（Shared）和 **可变**（Mutable）的状态的访问。“**共享**” 意味着变量可以有由多个线程同时访问，而 “**可变**” 则意味着变量的值在其生命周期内可以发生变化。要使得对象是线程安全的，需要采用同步机制来协同对对象可变状态的访问。Java中提供的主要同步机制是关键字 **synchronized** ，它提供了一种独占的加锁方式，但“同步”这个术语还包括 **volatile** 类型的变量，**显式锁（Explicit Lock）** 以及 **原子变量**。

## 1.2 线程安全类
从上面的定义中可以总结出：如果某个类满足线程安全性，那么就可以把它称作线程安全类。完全由线程安全类构成的程序并不一定就是线程安全的，而在线程安全类中也可以包含非线程安全的类。在后续的学习笔记中将会介绍如何组合使用线程安全类。在任何情况下，只有类中仅包含自己的状态时，线程安全类才是有意义的。在线程安全类中封装了必要的同步机制，因此客户端无须进一步采取同步措施。

## 1.3 无状态对象
无状态对象一定是线程安全的，这一点也很容易理解。

首先，我们来看一个简单的示例【**完整示例代码地址在文末提供**】：

```java
/**
 1. 一个基于Servlet的因数分解服务
 2. 这个Servlet从请求中提取数值，执行因数分解，然后将结果封装到该Servlet的响应中。
 */
public class StatelessFactorizer extends HttpServlet {

    protected void doGet(HttpServletRequest req, HttpServletResponse resp) throws IOException{
        BigInteger i = CommonUtils.extractFromRequest(req);
        BigInteger[] factors = Factor.factor(i);
        CommonUtils.encodeIntoResponse(resp, factors);
    }
}
```
然后我们来简单分析下：
（1）**StatelessFactorizer** 是无状态的：它既不包含任何域，也不包含任何对其他类中域的引用。

（2）上述示例的计算过程中的临时状态仅存在于线程栈上的局部变量中，并且只能由正在执行的线程访问，所以访问 **StatelessFactorizer** 的线程不会影响另一个访问同一个 **StatelessFactorizer** 的线程的计算结果。

最后根据分析，我们可以得出如下结论：

> 由于线程访问无状态对象的行为并不会影响其他线程中操作的正确性，因此无状态对象一定是线程安全的。


# 2. 原子性
下面我们在上述无状态对象中添加一个命中计数器的状态，用来统计所处理的请求数量。代码示例如下【**完整示例代码地址在文末提供**】：

```java
/**
 * <p> 在没有同步的情况下统计已经处理请求数量的Servlet </p>
 */
public class UnsafeCountingFactorizer extends HttpServlet {

    private long count = 0;

    public long getCount() {
        return count;
    }

    protected void doGet(HttpServletRequest req, HttpServletResponse resp) throws IOException{
        BigInteger i = CommonUtils.extractFromRequest(req);
        BigInteger[] factors = Factor.factor(i);
        ++count;
        CommonUtils.encodeIntoResponse(resp, factors);
    }
}
```
上述 **UnsafeCountingFactorizer** 是 **非线程安全** 的，这个类在并发环境下很可能会丢失一些更新操作。原因就在于++count，虽然递增操作++count这种语法看上去像是一个操作，但这个操作并非原子的，它实际上包含了三个独立的操作：**读取count的值**，**将值加1**，然后将**计算结果写入count**。

下图给出了两个线程在没有同步的情况下同时对一个计数器执行递增操作时发生的情况：

![](count.png)


如果计数器的初始值为0，在上图场景中Thread1和Thread2读到的count值都为0，接着执行递增操作，并且都将计数器的值设为1，显然这个结果并不是我们所希望看到的。

在并发编程中，这种由于不恰当的执行时序而出现不正确的结果的情况，有个专业的名词，我们称之为 **竞态条件**（Race Condition）。
## 2.1 竞态条件
当某个计算的正确性取决于多个线程的交替执行时序时，那么就会发生竞态条件，就比如上面的 **UnsafeCountingFactorizer**。

最常见的 **竞态条件** 类型就是 “**先检查后执行（Check-Then-Act）**”操作，即通过一个可能失效的观测结果来决定下一步的动作。怎么理解一个可能失效的观测结果呢？比如我们首先观测到某个条件为真（例如文件A不存在），然后根据这个观测结果采取相应的动作（例如创建文件A），但事实上，在我们观测到这个结果以及开始创建文件之间，观测结果可能变得无效（另一个线程在期间创建了文件A）。

## 2.2 延迟初始化
下面我们介绍使用 “**先检查后执行**” 的一种常见情况 ：**延迟初始化**。延迟初始化的目的是将对象的初始化操作推迟到实际被使用时才进行，同时要确保只被初始化一次。

首先我们来看一个 **延迟初始化** 的示例 【**完整示例代码地址在文末提供**】：

```java
/**
 * <p> 延迟初始化中的竞态条件（非线程安全，不推荐使用） </p>
 */
@NotThreadSafe
public class LazyInitRace {

    private ExpensiveObject instance = null;

    public ExpensiveObject getInstance() {
        if (instance == null)
            instance = new ExpensiveObject();
        return instance;
    }
}
```
上述 **LazyInitRace** 中包含了一个竞态条件，它可能会破坏该类的正确性，从而使得该类是非线程安全的。那么这里如何理解呢？假设线程 A 和线程 B 同时执行 getInstance。A 观测到 instance 为空，因而创建一个新的 ExpensiveObject 实例。B 同样需要判断 instance 是否为空。此时的 instance 是否为空，要取决于不可预测的时序，包括线程的调度方式，以及 A 需要花多长时间来初始化 ExpensiveObject 并设置 instance。如果当 B 检查时，instance 为空， 那么在两次调用 getInstance 时可能会得到不同的结果（即不同的 ExpensiveObject  实例对象）。

## 2.3 复合操作
上文的 UnsafeCountingFactorizer 和 LazyInitRace 都包含一组需要以原子方式执行（或者说不可分割）的操作。要避免竞态条件问题，就必须在某个线程修改该变量时，通过某种方式防止其他线程使用这个变量，从而确保其他线程只能在修改操作完成之前或之后读取和修改状态，而不是在修改状态的过程中。

> 假定有两个操作 A 和 B， 如果从执行 A 的线程来看， 当另一个线程执行 B 时， 要么将 B 全部执行完，要么完全不执行 B, 那么 A 和 B 对彼此来说是原子的。原子操作是指，对于访问同一个状态的所有操作（包括该操作本身）来说，这个操作是一个以原子方式执行的操作。

我们把“ **先检查后执行** ” （例如延迟初始化）和 “ **读取 -- 修改 -- 写入**”（例如递增运算）等操作统称为 复合操作。为了确保线程安全性，这些操作必须保证是以原子方式执行的操作。

在下面的章节将介绍加锁机制，这是Java中用于确保原子性的内置机制。但目前，我们使用一个现有的线程安全类来解决这个问题，代码示例如下【**完整示例代码地址在文末提供**】：

```java
/**
 * <p> 使用 AtomicLong 类型的变量来统计已处理请求的数量 </p>
 */
@ThreadSafe
public class CountingFactorizer extends HttpServlet {

    private final AtomicLong count = new AtomicLong(0);

    public long getCount() {
        return count.get();
    }

    protected void doGet(HttpServletRequest req, HttpServletResponse resp) throws IOException{
        BigInteger i = CommonUtils.extractFromRequest(req);
        BigInteger[] factors = Factor.factor(i);
        count.incrementAndGet();
        CommonUtils.encodeIntoResponse(resp, factors);
    }
}
```
上述的 **CountingFactorizer** 中，通过用 **AtomicLong** 来代替 long 类型的计数器，能够确保所有对计数器状态的访问操作都是原子的。由于该 **CountingFactorizer** 的状态就是计数器的状态，并且计数器是线程安全的，所以 **CountingFactorizer** 是线程安全的。

至此，我们可以得到如下结论：
>当在无状态的类中添加一个状态时，如果该状态完全由线程安全的对象来管理，那么这个类仍然是线程安全的。

# 3. 加锁机制
当在无状态的类中添加一个状态变量时，可以通过线程安全的对象来管理它的状态以维护它的线程安全性。但如果想要添加更多的状态，那么是否只需添加更多的线程安全状态变量就足够了？

下面我们来看一个代码示例【**完整示例代码地址在文末提供**】：

```java
/**
 * <p> 该Servlet在没有足够原子性保证的情况下对其最近计算结果进行缓存（非线程安全，不推荐使用） </p>
 */
@NotThreadSafe
public class UnsafeCachingFactorizer extends HttpServlet {
    // 最近执行因数分解的数值
    private final AtomicReference<BigInteger> lastNumber = new AtomicReference<>();
    // 分解结果
    private final AtomicReference<BigInteger[]> lastFactors = new AtomicReference<>();

    protected void doGet(HttpServletRequest req, HttpServletResponse resp) throws IOException {
        BigInteger i = CommonUtils.extractFromRequest(req);
        if (i.equals(lastNumber.get())) {
            CommonUtils.encodeIntoResponse(resp, lastFactors.get());
        } else {
            BigInteger[] factors = Factor.factor(i);
            lastNumber.set(i);
            lastFactors.set(factors);
            CommonUtils.encodeIntoResponse(resp, factors);
        }
    }
}
```
> **知识点** ：
> AtomicLong 是一种替代 long 类型整数的线程安全类，类似地，AtomicReference 是一种替代对象引用的线程安全类。

上述代码要实现的功能：将最近的计算结果缓存起来，当两个连续的请求对相同的数值进行因数分解时，可以直接使用上一次的计算结果，而无需重新计算（当然这里并不是一种有效的缓存策略，后续的笔记里面将会讲解更好的策略，敬请期待）。这里包含两个状态：最近执行因数分解的数值 **lastNumber** ，以及分解结果 **lastFactors** 。

在线程安全性的定义中要求，多个线程之间的操作无论采用何种执行时序或交替方式，都要保证不变性条件不被破坏。上述 **UnsafeCachingFactorizer** 的不变性条件之一是：在 **lastFactors** 中缓存的因数之积应该等于在 **lastNumber** 中缓存的数值。只有确保了这个不变性条件不被破坏，这个缓存策略实现的功能才是正确的。

然而，在某些执行时序中，**UnsafeCachingFactorizer** 可能会破坏这个不变性条件。在使用原子引用的情况下，尽管 **lastNumber** 和 **lastFactors** 它们对各自 set 方法的每次调用都是原子的，但在 **UnsafeCachingFactorizer** 中仍然无法做到同时更新 **lastNumber** 和 **lastFactors**。如果只修改了其中一个变量，那么在这两次修改操作之间，其他线程将会发现不变性条件被破坏了。同样，我们也不能保证会同时获取两个值：在线程 A 获取这两个值的过程中，线程 B 可能修改了它们，这样线程 A 也会发现不变性条件被破坏了。

当在不变性条件中涉及多个变量时，各个变量之间并不是彼此独立的，而是某个变量的值会对其他变量的值产生约束，这时就需要在单个原子操作中更新所有相关的状态变量，才能保持状态状态的一致性。

## 3.1 内置锁
Java 提供了一种内置的锁机制来支持原子性：**同步代码块**（Synchronized Block）。同步代码块包括两部分：一个作为锁的对象引用，一个作为由这个锁保护的代码块。以关键字 synchronized 来修饰的方法就是一种横跨整个方法体的同步代码块，其中该同步代码块的锁就是方法调用所在的对象。静态的 synchronized 方法以 Class 对象作为锁。

```java
	synchronized(lock) {
	    // 访问或修改由锁保护的共享状态
	}
```

每个 Java 对象都可以用做一个实现同步的锁，这些锁被称为**内置锁**（Intrinsic Lock）或监视器锁（Monitor Lock）。线程在进入同步代码块之前会自动获得锁，并且在退出同步代码块时自动释放锁（无论是通过正常的控制路径退出，还是通过从代码块中抛出异常退出）。获得内置锁的唯一途径就是进入由这个锁保护的同步代码块或方法。

Java 的内置锁相当于一种互斥体（或互斥锁），这也就意味着最多只有一个线程能持有这种锁。当线程 A 尝试获取一个由线程 B 持有的锁时，线程 A 必须等待或着阻塞，直到线程 B 释放这个锁。如果 B 永远不释放锁，那么 A 也将永远地等下去。由于每次只能有一个线程执行内置锁保护的代码块，因此，由这个锁保护的同步代码块会以原子方式执行，多个线程在执行该代码块时也不会相互干扰。

下面我们先看一个使用 **synchronized** 的代码示例【**完整示例代码地址在文末提供**】：

```java
/**
 * <p> 该 Servlet 可以正确地缓存最新的计算结果，但并发性却非常糟糕（不推荐使用） </p>
 */
@ThreadSafe
public class SynchronizedFactorizer extends HttpServlet {
    @GuardedBy("this") private BigInteger lastNumber;
    @GuardedBy("this") private BigInteger[] lastFactors;

    protected synchronized void doGet(HttpServletRequest req, HttpServletResponse resp) throws IOException {
        BigInteger i = CommonUtils.extractFromRequest(req);
        if (i.equals(lastNumber)) {
            CommonUtils.encodeIntoResponse(resp, lastFactors);
        } else {
            BigInteger[] factors = Factor.factor(i);
            lastNumber = i;
            lastFactors = factors;
            CommonUtils.encodeIntoResponse(resp, factors);
        }
    }
}
```

如上代码演示了使用 **synchronized** 来修饰 **doGet** 方法，虽然现在是线程安全的，但这种方式过于极端，多个客户端无法同时使用因数分解 Servlet，服务的响应性非常低，并发性非常糟糕，不推荐使用。这里是一个性能问题，而不是线程安全问题，下面我们会介绍另外的方法来解决它。

## 3.2 重入
当某个线程请求一个由其他线程持有的锁时，发出请求的线程就会阻塞。然而，由于内置锁是可重入的，因此如果某个线程试图获得一个已经由它自己持有的锁，那么这个请求就会成功。

下面给出相关的代码示例：

```java
public class Widget {
	public synchronized void doSomething() {
		...
	}
}

public class LoggingWidget extends Widget {
	public synchronized void doSomething() {
		System.out.println(toString() + ": calling doSomething");
		super.doSomething();
	}
}
```
上述代码中，子类 **LogginWidget** 改写了父类 **Widget** 的 **synchronized** 方法，然后调用父类中的方法。由于 **Widget** 和 **LoggingWidget** 中 **doSomething** 方法都是 **synchronized** 方法，因此每个**doSomething** 方法在执行前都会获取 **Widget** 上的锁。此时如果内置锁不是可重入的，那么在调用**super.doSomething** 时将无法获得 **Widget** 上的锁，因为这个锁已经被持有，从而线程将永远停顿下去，等待一个永远也无法获得的锁。重入则避免了这种死锁情况的发生。

# 4. 用锁来保护状态
由于锁能使其保护的代码路径以串行形式来访问，因此可以通过锁来构造一些协议以实现对共享状态的独占访问。

如果在复合操作的执行过程中持有一个锁，那么会使复合操作成为原子操作。当然仅仅将复合操作封装到一个同步代码块中是不够的。如果用同步来协调对某个变量的访问，那么在访问这个变量的所有位置上都需要使用同步。而且，当使用锁来协调对某个变量的访问时，在访问变量的所有位置上都要使用同一个锁。一种常见的错误是认为，只有在写入共享变量时才需要使用同步，然后事实并非如此（下一篇笔记将进一步解释其中的原因）。

对象的内置锁与其状态之间没有内在的关联。当获取与对象关联的锁时，并不能阻止其他线程访问该对象，某个线程在获得对象的锁之后，只能阻止其他线程获得同一个锁。之所以每个对象都有一个内置锁，只是为了免去显式地创建锁对象。

每个共享的和可变的变量都应该只由一个锁来保护，从而使维护人员知道是哪一个锁。一种常见的加锁约定是，将所有的可变状态都封装在对象内部，并通过对象的内置锁对所有访问可变状态的代码路径进行同步，使得在该对象上不会发生并发访问。

当然并非所有数据都需要锁的保护，只有被多个线程同时访问的可变数据才需要通过锁来保护。当某个变量由锁来保护时，意味着在每次访问这个变量时都需要首先获得锁，这样就确保在同一时刻只有一个线程可以访问这个变量。当类的不变性条件涉及多个状态变量时，那么在不变性条件中的每个变量都必须由同一个锁来保护。

如果通过同步可以避免竞态条件问题，那么为什么不在每个方法声明时都使用关键字 **synchronized** 呢？ 因为如果不加区别地滥用 **synchronized**，可能导致程序中出现过多的同步。另外，如果只是将每个方法都作为同步方法，例如 如下代码示例，这里并不足以确保 **Vector** 上复合操作都是原子的：
```java
	if (!vector.contains(element))
		vector.add(element);
```
虽然 **contains** 和 **add** 方法都是原子方法，但在上述示例代码的操作中仍然存在竞态条件。如果需要把多个操作合并为一个复合操作，仅仅使用 **synchronized** 是不够的，它只能确保单个操作的原子性，还是需要额外的加锁机制（后续笔记将会了解如何在线程安全对象中添加原子操作的方法）。

# 5. 活跃性与性能
上面 **SynchronizedFactorizer** 虽然解决了线程安全性问题，但还遗留了性能问题。当多个请求同时到达 **SynchronizedFactorizer** 时，这些请求将排队等待。

那么我们有没有办法可以既确保 Servlet 的并发性，同时又可以维护线程安全性呢？ 当然是有办法的，我们可以通过缩小同步代码块的作用范围来实现。不过需要注意以下三点：

 - 要保证同步代码块的合理大小（通常需要在安全性、简单性和性能等各种需求之间进行权衡）； 
 - 不要将本应是原子的操作拆分到多个同步代码块中；
 - 尽量将不影响共享状态且执行时间较长的操作从同步代码块中分离出去，从而在这些操作的执行过程中，其他线程可以访问共享状态。

下面我们来看一个改造后的代码示例【**完整示例代码地址在文末提供**】：

```java
/**
 * <p> 缓存最近执行因数分解的数值及其计算结果 </p>
 */
@ThreadSafe
public class CachedFactorizer extends HttpServlet {
    @GuardedBy("this") private BigInteger lastNumber;
    @GuardedBy("this") private BigInteger[] lastFactors;
    @GuardedBy("this") private long hits;
    @GuardedBy("this") private long cacheHits;

    public synchronized long getHits() { return hits; }
    public synchronized double getCacheHitRatio() {
        return (double) cacheHits / (double) hits;
    }

    protected void doGet(HttpServletRequest req, HttpServletResponse resp) throws IOException {
        BigInteger i = CommonUtils.extractFromRequest(req);
        BigInteger[] factors = null;
        synchronized (this) {
            ++hits;
            if (i.equals(lastNumber)) {
                ++cacheHits;
                factors = lastFactors.clone();
            }
        }
        if (factors == null) {
            factors = Factor.factor(i);
            synchronized (this) {
                lastNumber = i;
                lastFactors = factors.clone();
            }
        }
        CommonUtils.encodeIntoResponse(resp, factors);
    }
}
```
上述 **CachedFactorizer** 将代码修改为使用两个独立的同步代码块，每个同步代码块都只包含一小段代码。其中一个同步代码块负责保护判断是否只需要返回缓存结果的 “先检查后执行” 操作序列，另一个同步代码块则负责确保对缓存的数值和因数分解结果进行同步更新。此外，我们还重新引入了 “请求命中”计数器 ，添加了一个 “缓存命中” 计数器，并在第一个同步代码块中更新这两个变量。由于这两个计数器也是共享可变状态的一部分，因此必须在所有访问它们的位置上都使用同步。位于同步代码块之外的代码将以独占方式来访问局部（位于栈上的）变量，这些变量不会在多个线程间共享，因此不需要同步。

重新构造后的 **CachedFactorizer** 实现了在简单性（对整个方法进行同步）与并发性（对尽可能短的代码路径进行同步）之间的平衡。通常这两者存在着相互制约因素。当实现某个同步策略时，一定不要盲目地为了性能而牺牲简单性（这可能会破坏安全性）。

无论是执行计算密集的操作，还是在执行某个可能阻塞的操作，如果持有锁的时间过长，那么都会带来活跃性或性能问题。因此，当执行时间较长的计算或者可能无法快速完成的操作时（例如，网络I/O 或控制台 I/O），**一定不要持有锁**。

# 结语
线程安全性的相关学习就总结到这，这节的内容相对较多，虽然看完可能只需几个小时，但前前后后总结和归纳，花了几天的时间；即便如此，我仍然乐此不疲地做着这件事情，因为我也在并发学习的过程中，慢慢体会到了并发编程的魅力，通过巧妙的构思，简单的代码，可以解决很多复杂的问题。希望笔者的并发编程学习笔记系列可以帮助到正在学习并发编程的读者们。

本篇所有示例代码地址 [请点击这里](https://github.com/Huazie/FleaJavaConcurrency/tree/main/thread-safety-demo)，其中的 Servlet 可以通过 [JettyStarter](https://github.com/Huazie/FleaJavaConcurrency/blob/main/thread-safety-demo/src/test/java/com/huazie/flea/concurrency/threadsafety/JettyStarter.java) 启动服务端，然后浏览器访问 [http://localhost:8080/demo1?factor=1231231234](http://localhost:8080/demo1?factor=1231231234) 或者 [使用JMeter 模拟多用户高并发请求](https://www.cnblogs.com/Marydon20170307/p/14110839.html)

下篇的博文将介绍并发编程中对象的共享相关的问题，内容可能较多（其中包括了加锁机制和其他同步机制中的一种重要方面：**可见性**），我将拆分几篇来总结，敬请期待。