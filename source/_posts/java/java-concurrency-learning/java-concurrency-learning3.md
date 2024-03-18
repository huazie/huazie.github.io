---
title: Java并发编程学习3-可见性和对象发布
date: 2021-03-29 17:48:15
updated: 2023-09-14 12:36:42
categories:
  - [开发语言-Java,Java并发编程]
tags:
  - Java
  - 可见性
  - volatile
  - 对象发布与逸出
---

[《开发语言-Java》](/categories/开发语言-Java/) [《Java并发编程》](/categories/开发语言-Go/Go语言学习/) 

![](/images/java-concurrency-logo.png)

# 引言
书接上篇，我们了解了如何通过同步来避免多个线程在同一时刻访问相同的数据，而本篇将介绍如何共享和发布对象，从而使它们能够安全地由多个线程同时访问。

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 1. 可见性
线程安全性的内容，让我们知道了同步代码块和同步方法可以确保以原子的方式执行操作。但如果你认为关键字 **synchronized** 只能用于实现原子性或者确定“临界区（Critical Section）”，那就大错特错了。同步还有一个重要的方面：内存可见性（Memory Visibility）。我们不仅希望防止某个线程正在使用对象状态而另一个线程在同时修改该状态，而且希望确保当一个线程修改了对象状态后，其他线程能够看到发生的状态变化。

可见性是一种复杂的属性，在一般的单线程环境中，如果向某个变量先写入值，然后在没有其他写入操作的情况下读取这个变量，总是能够得到相同的值。然而，当读操作和写操作在不同的线程中执行时，因为无法确保执行读操作的线程能适时地看到其他线程写入的值，所以也就不能总是得到相同的值。为了确保多个线程之间对内存写入操作的可见性，必须使用同步机制。

介绍了这么多，还不如来看下代码示例：

```java
/**
 * <p> 在没有同步的情况下共享变量（不推荐使用） </p>
 */
public class NoVisibility {
    private static boolean ready;
    private static int number;

    private static class ReaderThread extends Thread {
        @Override
        public void run() {
            while (!ready) {
                Thread.yield();
            }
            System.out.println(number);
        }
    }

    public static void main(String[] args) {
        new ReaderThread().start();
        number = 42;
        ready = true;
    }
}

```
在上述代码中，主线程和读线程都将访问共享变量 **ready** 和 **number**。主线程启动读线程，然后将 **number** 设为 42， 并将 **ready** 设为 true。读线程一直循环直到发现 **ready** 的值变为 true，然后输出 **number** 的值。虽然 **NoVisibility** 看起来会输出 42，但事实上很可能输出0，或者根本无法终止。因为在代码中没有使用足够的同步机制，所以无法保证主线程写入的 **ready** 值 和 **number** 值对于读线程来说是可见的。 

如果你尝试运行该程序，大概率控制台还是会输出42，但这并不说明这块代码就总是能输出想要的结果。**NoVisibility** 可能会输出0，这是因为读线程可能看到了写入 **ready** 的值，但却没有看到之后写入 **number** 的值，这种现象被称为 “重排序”；**NoVisibility** 也可能会一直循环下去，因为读线程可能永远都看不到 **ready** 的值。

>在没有同步的情况下，编译器、处理器以及运行时等都可能对操作的执行顺序进行一些意想不到的调整。在缺乏足够同步的多线程程序中，要想对内存操作的执行顺序进行判断，几乎无法得出正确的结论。

## 1.1 失效数据
**NoVisibility** 展示了在缺乏同步的程序中可能产生错误结果的一种情况：失效数据。当读线程查看 ready 变量时，可能会得到一个已经失效的值。更糟糕的是，失效值可能不会同时出现：一个线程可能获得某个变量得最新值，而获得另一个变量得失效值。

下面再看一个代码示例：

```java
/**
 * <p> 非线程安全的可变整数类 </p>
 */
@NotThreadSafe
public class MutableInteger {
    private int value;

    public int getValue() { return value; }
    public void setValue(int value) { this.value = value; }
}
```
上述代码中 **get** 和 **set** 方法都是在没有同步的情况下访问 **value** 的。如果某个线程调用了 **set** 方法，那么另一个正在调用 **get** 方法的线程可能会看到更新后的 **value** 值，也可能看不到。

下面我们通过对 **get** 和 **set** 方法进行同步，可以使 **MutableInteger** 成为一个线程安全的类。代码示例如下：

```java
/**
 * <p> 非线程安全的可变整数类 </p>
 */
public class SynchronizedInteger {
    @GuardedBy("this") private int value;

    public synchronized int getValue() { return value; }
    public synchronized void setValue(int value) { this.value = value; }
}
```
当然如果这里仅仅对 **set** 方法进行同步是不够的，调用 **get** 方法的线程仍然会看见失效值。

## 1.2 非原子的64位操作
上面我们了解到，当线程在没有同步的情况下读取变量时，可能会得到一个失效值，但至少这个值是由之前某个线程设置的值，而不是一个随机值。这种安全性保证也被称为最低安全性（out-of-thin-air-safety）。

最低安全性适用于绝大多数变量，但是非 **volatile** 类型的64位数值变量例外。Java内存模型要求，变量的读取操作和写入操作都必须是原子操作，但对于非 **volatile** 类型的 **long** 和 **double** 变量，JVM允许将64位的读操作或写操作分解为两个32位操作。当读取一个非 **volatile** 类型 的 **long** 变量时，如果对该变量的读操作和写操作在不同的线程中执行，那么很可能会读取到某个值的高32位和另一个值的低32位。

## 1.3 加锁与可见性
内置锁可以用于确保某个线程以一种可预测得方式来查看另一个线程的执行结果，如下图所示。当线程 A 执行某个同步代码块时，线程 B 随后进入由同一个锁保护的同步代码块，在这种情况下可以保证，在锁被释放之前，A 看到的变量值在 B 获得锁后同样可以由 B 看到。换句话说，当线程 B 执行由锁保护的同步代码块时，可以看到线程 A 之前在同一个同步代码块中的所有操作结果。

![](lock.png)

> 加锁的含义不仅仅局限于互斥行为，还包括内存可见性。为了确保所有线程都能看到共享变量的最新值，所有执行读操作或者写操作的线程都必须在同一个锁上同步。

## 1.4 volatile 变量
Java语言提供了一种稍弱的同步机制，即 **volatile** 变量，用来确保将变量的更新操作通知到其他线程。当变量声明为 **volatile** 类型后，编译器与运行时都会注意到这个变量是共享的，因此不会将该变量上的操作与其他内存操作一起重排序。**volatile** 变量不会被缓存在寄存器或者对其他处理器不可见的地方，因此在读取 **volatile** 类型的变量时总会返回最新写入的值。

当然，这里不建议过度依赖 **volatile** 变量提供的可见性。仅当 **volatile** 变量能简化代码的实现以及对同步策略的验证时，才应该使用它们。如果在验证正确性时需要对可见性进行复杂的判断，那么就建议使用 **volatile** 变量。

**volatile** 变量的正确使用方式包括：
1. 确保它们自身状态的可见性；
2. 确保它们所引用对象的状态的可见性；
3. 标识一些重要的程序生命周期事件的发生（初始化或关闭）

下面看一个利用 **volatile** 变量来数绵羊的代码示例：

```java
volatile boolean asleep;
// ...
    while (!asleep) 
        countSomeSheep();
```
在如上示例中，线程试图通过类似数绵羊的传统方式进入休眠状态。相比用锁来确保 **asleep** 更新操作的可见性，这里采用 **volatile** 变量，不仅满足了更新操作的可见性，而且代码逻辑也变得更加简单，更利于理解。

虽然 **volatile** 变量使用很方便，但它只能确保可见性，而加锁机制既可以确保可见性又可以确保原子性。

那么说了这么多，什么场景下我们才应该使用 **volatile** 变量呢？

当且仅当满足以下条件：
- 对变量的写入操作不依赖变量的当前值，或者你能确保只有单个线程更新变量的值。
- 该变量不会与其他状态变量一起纳入不变性条件中。
- 在访问变量时不需要加锁。

# 2. 发布与逸出

## 2.1 发布对象
本篇开头提到了 **发布对象**，它是指使对象能够在当前作用域之外的代码中使用。例如，将一个指向该对象的引用保存到其他代码可以访问的地方，或者在某一个非私有的方法中返回该引用，或者将引用传递到其他类的方法中。当某个不应该发布的对象被发布了，这种情况就被称为 **逸出**（Escape）。

发布对象的最简单方法是将对象的引用保存到一个公有的静态变量中，以便任何类和线程都能够看见该对象。

下面展示发布一个对象的代码示例：

```java
	public static Set<Secret> knownSecrets;
	
	public void initialize() {
	    knownSecrets = new HashSet<Secret>();
	}
```
上述代码中，在 **initialize** 方法中示例化一个新的 **HashSet** 对象，并将对象的引用保存到 **knownSecrets** 中以发布该对象。如果将一个 **Secret** 对象添加到集合 **knownSecrets** 中，那么同样会发布这个 **Secret** 对象，因为任何代码都可以遍历这个集合，并获得对这个新 **Secret** 对象的引用。

我们再来看一个代码示例：

```java
/**
 * <p> 使内部的可变状态逸出（不推荐使用） </p>
 */
public class UnsafeStates {
    private String[] states = new String[] {"HELLO", "HUAZIE"};

    public String[] getStates() { return states; }
}
```
上诉代码从非私有方法 **getStates** 中返回一个引用，这里同样会发布返回的引用的对象 **states** 。按上述方式来发布 **states**，就可能存在很大风险，因为任何调用者都能修改这个数组的内容。

如果一个已经发布的对象能够通过非私有的变量引用和方法调用到达其他的对象，那么这些对象也都会被发布。

最后一种发布对象或其内部状态的机制就是发布一个内部的类实例，如下代码示例：

```java
/**
 * <p> 隐式地使this引用逸出（不推荐使用） </p>
 */
public class ThisEscape {
    public ThisEscape(EventSource source) {
        source.registerListener(new EventListener(){
            public void onEvent(Event e) {
                doSomething(e);
            }
        });
    }

    private void doSomething(Event e) {
        // 事件处理
    }
}
```
当 **ThisEscape** 发布 **EventListener** 时，也隐含地发布了 **ThisEscape** 实例本身，因为在这个内部类的实例中包含了对 **ThisEscape** 实例的隐含引用。

## 2.2 安全的对象构造过程
在 **ThisEscape** 中给出了逸出的一个特殊示例，即 **this** 引用在构造函数中逸出。如果 **this** 引用在构造过程中逸出，那么这种对象就被认为是不正确构造。
>**注意**： 不要在构造过程中使 **this** 引用逸出

如果想在构造函数中注册一个事件监听器或启动进程，那么可以使用一个私有的构造函数和一个公共的工厂方法，从而避免不正确的构造过程。下面请看如下代码示例：

```java
/**
 * <p> 使用工厂方法来防止this引用在构造过程中逸出 </p>
 */
public class SafeListener {
    private final EventListener listener;

    private SafeListener() {
        listener = new EventListener(){
            public void onEvent(Event e) {
                doSomething(e);
            }
        };
    }

    public static SafeListener newInstance(EventSource source) {
        SafeListener safe = new SafeListener();
        source.registerListener(safe.listener);
        return safe;
    }

    private void doSomething(Event e) {
        // 事件处理
    }
}
```

# 结语
本篇我们一起了解了 可见性 和  对象的发布、逸出等相关内容；关于对象的共享的其他内容【线程封闭，不变性，安全发布】，还需要一篇博文才能介绍完，敬请期待！