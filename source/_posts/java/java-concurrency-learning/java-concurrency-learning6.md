---
title: Java并发编程学习6-同步容器类和并发容器
date: 2022-09-11 08:45:00
updated: 2024-03-19 17:21:35
categories:
  - [开发语言-Java,Java并发编程]
tags:
  - Java
  - 同步容器类
  - 并发容器
  - ConcurrentHashMap
  - CopyOnWriteArrayList
---

[《开发语言-Java》](/categories/开发语言-Java/)

![](/images/java-concurrency-logo.png)

# 引言

本篇开始将要介绍 **Java** 平台类库下的一些最常用的 **并发基础构建模块**，以及使用这些模块来构造并发应用程序时的一些常用模式。

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 同步容器类

同步容器类包括 **Vector** 和 **Hashtable**，还有由 **Collections.synchronizedXxx** 等工厂方法创建的同步的封装器类。

这些类实现线程安全性的方法是：将它们的状态封装起来，并对每个公有方法都进行同步，使得每次只有一个线程能访问容器的状态。

## 1. 同步容器类的问题

同步容器类都是线程安全的，但在某些情况下可能需要额外的客户端加锁来保护复合操作。

容器里常见的复合操作包括：
- 迭代（反复访问元素，直到遍历完容器中的所有元素）
- 跳转（根据指定顺序找到当前元素的下一个元素）
- 条件运算（例如“若没有则添加”）

下面假设我们在 **Vector** 中定义两个复合操作的方法：**getLast** 和 **deleteLast**，它们都会执行 **“先检查后执行”** 操作。

```java
    public static Object getLast(Vector list) {
        int lastIndex = list.size() - 1;
        return list.get(lastIndex);
    }

    public static void deleteLast(Vector list) {
        int lastIndex = list.size() - 1;
        list.remove(lastIndex);
    }
```

上面定义的两个方法，看似没有任何问题，从某种程度上来看也的确如此，无论多少个线程同时调用它们，也不会破坏 **Vector**。

但如果从调用者的角度去看，如果线程 **A** 在包含 **10** 个元素的 **Vector** 上调用 **getLast**，同时线程 **B** 在同一个 **Vector** 上调用 **deleteLast**，这些操作交替执行如下图所示，**getLast** 将抛出 **ArrayIndexOutOfBoundsException** 异常。这里虽然很好地遵循了 **Vector** 的规范（如果请求一个不存在的元素，那么将抛出一个异常），但这并不是调用者所希望看到的结果，除非 **Vector** 一开始就是空的。

![](vector-problem.png)

同步容器类通过自身的锁来保护它的每个方法，因此只要获得容器类的锁，上面的 **getLast** 和 **deleteLast** 方法就可以成为原子操作。

下面看一下代码示例：

```java
    public static Object getLast(Vector list) {
        synchronized (list) {
            int lastIndex = list.size() - 1;
            return list.get(lastIndex);
        }
    }
    
    public static void deleteLast(Vector list) {
        synchronized (list) {
            int lastIndex = list.size() - 1;
            list.remove(lastIndex);
        }
    }
```

与 **getLast** 一样，如果在对 **Vector** 进行迭代时，另一个线程删除了一个元素，并且这两个操作交替执行，那么这种迭代方法也将抛出 **ArrayIndexOutOfBoundsException** 异常。

```java
    for (int i = 0; i < vector.size(); i++)
        doSomething(vector.get(i));
```

虽然上述迭代操作可能抛出异常，但并不意味着 **Vector** 就不是线程安全的。**Vector** 的状态仍然是有效的，而抛出的异常也与其规范保持一致。

像读取最后一个元素或者迭代这样的简单操作中抛出异常，显然是调用者不愿意看到的。我们可以通过在迭代期间持有 **Vector** 的锁，可以防止其他线程在迭代期间修改 **Vector**。当然这会导致其他线程在迭代期间无法访问它，从而降低了并发性。

```java
    synchronized (vector) {
        for (int i = 0; i < vector.size(); i++)
            doSomething(vector.get(i));
    }
```

## 2. 迭代器与 ConcurrentModificationException

在设计同步容器类的迭代器时并没有考虑到并发修改的问题，并且它们表现出的行为是 **“及时失败”** 的。这意味着，当它们发现容器在迭代过程中被修改时，就会抛出一个 **ConcurrentModificationException** 异常。

这种 **“及时失败”** 的迭代器只能作为并发问题的预警指示器。如果在迭代期间计数器被修改，那么 **hasNext** 或 **next** 将抛出 **ConcurrentModificationException**。然而，这种检查是在没有同步的情况下进行的，因此可能会看到失效的计数值，而迭代器可能并没有意识到已经发生了修改。这是一种设计上的权衡，从而降低并发修改操作的检测代码对程序性能带来的影响。

下面我们看一个代码示例，使用 **for-each** 循环语法对 **List** 容器进行迭代。

```java
    List<Person> personList = Collections.synchronizedList(new ArrayList<Person>());
    
    // 可能抛出 ConcurrentModificationException
    for (Person p : personList)
        doSomething(p);
```

从编译后的Class文件来看，上述 **for-each** 循环语法，**javac** 将生成使用 **Iterator** 的代码，反复调用 **hasNext** 和 **next** 来迭代 **List** 对象。 与迭代 **Vector** 一样，想要避免出现 **ConcurrentModificationException**，就必须在迭代过程中持有容器的锁。

如果不希望在迭代期间对容器加锁，那么可以“克隆”容器，并在副本上进行迭代。由于副本被封闭在线程内，因此其他线程不会在迭代期间对其进行修改，这样就避免了抛出 **ConcurrentModificationException**（在克隆过程中仍然需要对容器加锁）。

当然克隆容器存在显著的性能开销。这种方式的好坏，取决于容器的大小，在每个元素上执行的操作，迭代操作相对于容器其他操作的调用频率，以及在响应时间和吞吐量等方面的需求。

## 3. 隐藏迭代器

虽然加锁可以防止迭代器抛出 **ConcurrentModificationException**，但需要记住在所有对共享容器进行迭代的地方都需要加锁。

下面我们来看一个示例，在 **HiddenIterator** 中没有显式的容器迭代操作，但在 **System.out.pringln** 中将执行迭代操作。

```java
    @NotThreadSafe
    public class HiddenIterator {
        @GuardedBy("this")
        private final Set<Integer> set = new HashSet<Integer>();
        
        public synchronized void add(Integer i) {
            set.add(i);
        }
        
        public synchronized void remove(Integer i) {
            set.remove(i);
        }
        
        public void addTenThings() {
            Random r = new Random();
            for (int i = 0; i < 10; i++)
                add(r.nextInt());
            // 隐藏在字符串连接中的迭代操作
            System.out.println("DEBUG: added ten elements to " + set);
        }
    }
```

上述 **System.out.println** 代码中，编译器将字符串的连接操作转换为调用 **StringBuilder.append(Object)**，而这个方法又会调用容器的 **toString** 方法，标准容器的 **toString** 方法将迭代容器，并在每个元素上调用 **toString** 来生成容器内容的格式化表示。并发环境下，**addTenThings** 方法可能会抛出 **ConcurrentModificationException**。

如果状态与保护它的同步代码之间相隔越远，那么开发人员就越容易忘记在访问状态时使用正确的同步。如果 **HiddenIterator** 用 **synchronizedSet** 来包装 **HashSet**，并且对同步的代码进行封装，那么就不会发生这种错误。

> 正如封装对象的状态有助于维持不变性条件一样，封装对象的同步机制同样有助于确保实施同步策略。

除了 **toString** 对容器进行迭代，还有容器的 **hashCode**、**equals**、**containsAll**、**removeAll** 和 **retainAll** 等方法，以及把容器作为参数的构造函数，都会对容器进行迭代。所有这些间接的迭代操作都有可能抛出 **ConcurrentModificationException**。


# 并发容器

上面提到的同步容器，它是将所有对容器状态的访问都串行化，以实现它们的线程安全性。这种方式的代价就是严重降低并发性，当多个线程竞争容器的锁时，吞吐量将严重降低。

并发容器是针对多个线程并发访问而设计，如 **ConcurrentHashMap**，用于替代同步且基于散列的 Map；**CopyOnWriteArrayList**，用于在遍历操作为主要操作的情况下代替同步的 List。

> 通过并发容器来代替同步容器，可以极大地提高伸缩性并降低风险。

## 1. ConcurrentHashMap

与 **HashMap** 一样，**ConcurrentHashMap** 也是一个基于散列的 **Map**, 但它使用了一种粒度更细的加锁机制来实现更大程度的共享，提供更高的并发性和伸缩性，这种机制称为分段锁（**Lock Striping**，以后的博文会讲解到）。在这种机制中，任意数量的读取线程可以并发地访问 Map，执行读取操作的线程和执行写入操作的线程可以并发地访问 **Map**，并且一定数量的写入线程可以并发地修改 **Map**。

**ConcurrentHashMap** 与其他并发容器一起增强了同步容器类，有如下的特点：

- 它们提供的迭代器不会抛出 **ConcurrentModificationException**，因此不需要再迭代过程中对容器加锁。
- **ConcurrentHashMap** 返回的迭代器具有弱一致性（**Weakly Consistent**），而并非 ”及时失败“。弱一致性的迭代器可以容忍并发的修改，当创建迭代器时会遍历已有的元素，并可以（但是不保证）在迭代器被构造后将修改操作反映给容器。

对于一些需要在整个 **Map** 上进行计算的方法，例如 **size** 和 **isEmpty**，这些方法的语义被略微减弱了以反映容器的并发特性。由于 **size** 返回的结果在计算时可能已经过期了，它实际上只是一个估计值，因此允许 **size** 返回一个近似值而不是一个精确值。事实上 **size** 和 **isEmpty** 这样的方法在并发环境下的用处很小，因为它们的返回值总是不断变化。因此，这些操作的需求被弱化了，以及换取对其他更重要操作的性能优化，包括 **get**、**put**、**containsKey** 和 **remove** 等。

在 **ConcurrentHashMap** 中没有实现对 **Map** 加锁以提供独占访问，而在 **Hashtable** 和 **synchronizedMap** 中，获得 **Map** 的锁能防止其他线程访问这个 **Map**。大多数情况下，用 **ConcurrentHashMap** 来代替同步 **Map** 能进一步提高代码的可伸缩性。

## 2. 额外的原子Map操作

由于 **ConcurrentHashMap** 不能被加锁来执行独占访问，因此无法使用客户端加锁来创造新的原子操作。 不过像 “若没有则添加”、“若相等则移除” 和 “若相等则替换” 等，都已经实现为原子操作并且在 **ConcurrentMap** 的接口中声明，如下代码所示：

```java
    public intercace ConcurrentHashMap<K, V> extends Map<K, V> {
        // 仅当 K 没有相应的映射值时才插入
        V putIfAbsent(K key, V value);
        
        // 仅当 K 被映射到 V 才移除
        boolean remove(K key, V value);
        
        // 仅当 K 被映射到 oldValue 时才替换为 newValue
        boolean replace(K key, V oldValue, V newValue);
        
        // 仅当K 被映射到某个值时才替换为 newValue
        V replace(K key, V newValue);
    }
```

如果你需要在现有的同步 **Map** 中添加如上的操作，那么也就意味着应该考虑使用 **ConcurrentMap** 了。

## 3. CopyOnWriteArrayList

**CopyOnWriteArrayList** 用于替代同步 **List**，在某些情况下它提供了更好的并发性能，并且在迭代期间不需要对容器进行加锁或复制。

> 类似地，**CopyOnWriteArraySet** 用于替代同步Set。

“写入时复制（**Copy-On-Write**）”容器的线程安全性在于，只要正确地发布一个事实不可变的对象，那么在访问该对象时就不再需要进一步的同步。在每次修改时，都会创建并重新发布一个新的容器副本，从而实现可变性。“写入时复制” 容器的迭代器保留一个指向底层基础数组的引用，这个数组当前位于迭代器的起始位置，由于它不会被修改，因此在对其进行同步时只需确保数组内容的可见性。

显然，每当修改容器时都会复制底层数组，这需要一定的开销，特别是当容器的规模较大时。仅当迭代操作远远多于修改操作时，才应该使用 “写入时复制” 容器。

> 许多事件通知系统中，在分发通知时需要迭代已注册监听器链表，并调用每一个监听器，在大多数情况下，注册和注销事件监听器的操作远少于接收事件通知的操作。


## 4. 阻塞队列
这块的篇幅较多，下一篇博文将会详细介绍，尽情期待！