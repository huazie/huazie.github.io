---
title: Java并发编程学习5-对象的组合
date: 2022-09-08 10:00:25 
updated: 2024-03-19 17:21:35
categories:
  - [开发语言-Java,Java并发编程]
tags:
  - Java
  - 对象的组合
  - 实例封闭
  - 线程安全性的委托
  - 客户端加锁机制
---

[《开发语言-Java》](/categories/开发语言-Java/)

![](/images/java-concurrency-logo.png)

# 引言
前面的博文，我们已经了解了关于线程安全和同步的一些基础知识。本篇博文将介绍一些线程安全的组合模式，来帮助我们确保使用这些模式开发的程序是线程安全的。

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 1. 设计线程安全的类
我们考虑一下该如何设计一个线程安全的类？

首先，能想到的就是要确保对象中所有的状态变量都是在可控范围内的。

因此在设计线程安全类的过程中，需要包含如下三个基本要素：
- 找出构成对象状态的所有变量。
- 找出约束状态变量的不可变条件。
- 建立对象状态的并发访问管理策略。

那如何找构成对象状态的所有变量？ 

要分析对象的状态，首先得从对象的域开始。如果对象中所有的域都是基本类型的变量，那么这些域将构成对象的全部状态。如果在对象的域中引用了其他对象，那么该对象的状态将包含被引用对象的域。

下面我们来看一个示例：

```java
/**
 * <p> 使用 Java 监视器模式的线程安全计数器 </p>
 *
 * @author huazie
 */
public class Counter {

    @GuardedBy("this")
    private long value = 0;

    public synchronized long getValue() {
        return value;
    }

    public synchronized long increment() {
        if (value == Long.MAX_VALUE)
            throw new IllegalStateException("counter overflow");
        return ++value;
    }
}
```
如上，**Counter** 只有一个域 **value**，而且它是基本数据类型，因此这个域就是 **Counter** 的全部状态。

前两个基本要素都找到了，下一步我们就可以建立相应的并发访问管理策略，即同步策略（**Synchronization Policy**），它定义了如何在不违背对象不变性或后验条件的情况下对其状态的访问操作进行协同。同步策略规定了如何将不可变性、线程封闭与加锁机制等结合起来以维护线程的安全性，并且还规定了哪些变量由哪些锁来保护。

## 1.1 收集同步需求
要确保类的线程安全性，就需要确保它的**不变性条件**不会在并发访问的情况下被破坏，这就需要对其状态进行推断。

对象与变量都有一个状态空间，即所有可能的取值。状态空间越小，就越容易判断线程的状态。上述**Counter** 中的 **value** 域是 **long** 类型的变量，其状态空间为 从 `Long.MIN_VALUE` 到 `Long.MAX_VALUE`，即从 $-2^{63}$ 到 $2^{63}$-1，但 **Counter** 中 **value** 在取值范围上存在一个限制，即不可能是负值。

同样，在操作中还会包含一些**后验条件**来判断状态迁移是否是有效的。如果 **Counter** 的当前状态为 **1**，那么下一个有效的状态只能是 **2**。当下一个状态需要依赖当前状态时，这个操作就必须是一个复合操作。当然并非所有的操作都会在状态转换上施加限制。

由于不变性条件以及后验条件在状态及状态转换上施加了各种约束，因此就需要额外的同步与封装。

在类中也可以包含同时约束多个状态变量的不变性条件。比如一个表示数值范围的类中可以包含两个状态变量，分别表示范围的上界和下界，并且下界值应该小于等于上界值。上述情况下，这些相关变量必须在单个原子操作中进行读取或更新，不然可能会使对象处于无效的状态。如果在一个不变性条件中包含多个变量，那么在执行任何访问相关变量的操作时，都必须持有保护这些变量的锁。

## 1.2 依赖状态的操作
前面提到，类的不变性条件与后验条件约束了在对象上有哪些状态和状态转换是有效的。当然，在某些对象的方法中还包含一些基于状态的先验条件。例如，不能从空队列中移除一个元素；在删除元素前，队列必须处于”非空“的状态。如果在某个操作中包含有基于状态的先验条件，那么这个操作就称为依赖状态的操作。

通过现有库中的类，例如阻塞队列[Blocking Queue]或信号量[Semaphore]，可以更简单地来实现依赖状态的行为，即某个等待先验条件为真时才执行的操作。
## 1.3 状态的所有权
所有权在 **Java** 中是属于类设计中的一个要素，不像 **C或C++**，需要认真考虑所有权的处理，**Java** 通过垃圾回收机制，减少了许多在引用共享方面常见的错误，降低了在所有权处理上的开销。

许多情况下，对象对它封装的状态拥有所有权。所有权意味着控制权。然而，如果发布了某个可变对象的引用，那么就不再拥有独占的控制权，最多是“共享控制权”。为了防止多个线程在并发访问同一个对象时产生的相互干扰，这些对象应该要么是线程安全的对象，要么是事实不可变的对象，或者由锁来保护的对象。

# 2. 实例封闭
如果某对象不是线程安全的，那么可以通过多种技术使其在多线程程序中安全地使用。封装简化了线程安全类的实现过程，它提供了一种实例封闭机制（**Instance Confinement**），简称”封闭“。

将数据封装在对象内部，可以将数据的访问限制在对象的方法上，从而更容易确保线程在访问数据时总能持有正确的锁。

如下示例 **PersonSet**，展示了如何通过封闭与加锁等机制使一个类成为线程安全的（即使这个类的状态变量并不是线程安全的）。

```java
@ThreadSafe
public class PersonSet {
    @GuardedBy("this")
    private final Set<Person> mySet = new HashSet<Person>();
    
    public synchronized void addPerson(Person p) {
        mySet.add(p);
    }
    
    public synchronized boolean containsPerson(Person p) {
        return mySet.contains(p);
    }
}
```
通过简单分析可知，**PersonSet** 的状态由 **HashSet** 管理的，而 **HashSet** 本身不是线程安全的。但由于 **mySet** 是私有的并且不会逸出，因而 **HashSet** 被封闭在 **PersonSet** 中。唯一能访问 **mySet** 的是 **addPerson** 与 **containsPerson**，在执行它们时都要获得 **PersonSet** 上的锁。**PersonSet** 的状态完全由它的内置锁保护，因而 **PersonSet** 是一个线程安全的类。

上述示例并未对 **Person** 的线程安全性做任何假设，但如果 **Person** 类是可变的，那么在访问从 **PersonSet** 中获得的 **Person** 对象时，还需要额外的同步。

在 **Java** 中，一些基本的容器类并非线程安全的，例如 **ArrayList** 和 **HashMap**，但类库中提供了包装器工厂方法（如 **Collections.synchronizedList** 及其类似的方法），使得这些非线程安全的类可以在多线程环境中安全地使用。这些工厂方法通过”装饰器“模式将容器类封装在一个同步的包装器对象中，而包装器能将接口中的每个方法都实现为同步方法，并将调用请求转发到底层的容器对象上。对底层容器对象的所有访问必须通过包装器来进行。

## 2.1 Java监视器模式
**Java** 监视器模式来自于 **Hoare** 对监视器机制的研究工作。在JVM中，进入和退出同步代码块的字节指令也称为 **monitorenter** 和 **monitorexit**，而 **Java** 的内置锁也称为 **监视器锁** 或 **监视器**。遵循 **Java** 监视器模式的对象会把对象的所有可变状态都封装起来，并由对象自己的内置锁来保护。

**Java** 监视器模式的简单使用示例可以参考上面的 **Counter** 类。**Java** 监视器模式模式仅仅是一种编写代码的约定，对于任何一种锁对象，只要自始至终都使用该锁对象，都可以用来保护对象的状态。

如下代码展示了如何使用私有锁来保护状态：

```java
    public class PrivateLock {
        private final Object myLock = new Object();
        
        @GuardedBy("myLock")
        Person person;
        
        void someMethod() {
            synchronized(myLock) {
                // 访问或修改Person的状态
            }
        }
    }
```

私有的锁对象可以将锁封装起来，使客户代码无法得到锁，但客户代码可以通过公有方法来访问锁，以便参与到它的同步策略中。

## 2.2 “车辆追踪” 示例

下面我们来看一个相比 **Counter** 类，更有用处的示例 -- **“车辆追踪”**：一个用于调度车辆的“车辆追踪器”。

首先使用监视器模式来构建车辆追踪器，代码清单如下所示，然后再尝试放宽某些封装性需求同时又保持线程安全性。

```java
@ThreadSafe
public class MonitorVehicleTracker {
    @GuardedBy("this")
    private final Map<String, MutablePoint> locations;
    
    public MonitorVehicleTracker(Map<String, MutablePoint> locations) {
        this.locations = deepCopy(locations);
    }
	
    /**
     * 当有大量车辆需要追踪的时候，这里执行的复制操作花费可能就会比较长，
     * 车辆追踪器的内置锁将一直被占用，这样会严重降低用户界面的响应灵敏度
     */
    public synchronized Map<String, MutablePoint> getLocations() {
        return deepCopy(locations);
    }
    
    public synchronized MutablePoint getLocations(String id) {
        MutablePoint loc = locations.getLocations(id);
        return loc == null ? null : new MutablePoint(loc);
    }
    
    public synchronized void setLocation(String id, int x, int y) {
        MutablePoint loc = locations.get(id);
        if (loc == null) 
            throw new IllegalArgumentException("No such ID" + id);
        loc.x = x;
        loc.y = y;
    }
    
    private static Map<String, MutablePoint> deepCopy(Map<String, MutablePoint> m) {
        Map<String, MutablePoint> result = new HashMap<String, MutablePoint>();
        for (String id : m.keySet())
            result.put(id, new MutablePoint(m.get(id)));
        return Collections.unmodifiableMap(result);
    }
}
```

每辆车都有一个 **ID** 标识，并且拥有一个相应的位置坐标（x,y）。

```java
@NotThreadSafe
public class MutablePoint {
    public int x, y;

    public MutablePoint() {
        x = 0;
        y = 0;
    }

    public MutablePoint(MutablePoint p) {
        this.x = p.x;
        this.y = p.y;
    }
}
```

虽然 **MutablePoint** 不是线程安全的，但追踪器类是线程安全的。它包含的 **Map** 对象和可变的 **MutablePoint** 对象都未曾发布。

在某种程度上，上述实现方式是通过在返回客户代码之前复制可变的数据来维持线程安全性的。通常情况下，这并不存在性能问题，但在车辆容量非常大的情况下将极大地降低性能【这里可以看下 **getLocations** 方法的备注】。

# 3. 线程安全性的委托

## 3.1 基于委托的车辆追踪器

下面我们介绍一个更实际的委托示例，构造一个委托给线程安全类的车辆追踪器。

首先，我们用一个**不可变的 Point** 来代替 **MutablePoint** ，用以保存位置。

```java
@Immutable
public class Point {
    public final int x, y;

    public Point(int x, int y) {
        this.x = x;
        this.y = y;
    }
}
```

由于 **Point** 类是不可变的，因而它是线程安全的。不可变的值可以被自由地共享与发布，因此在返回 **location** 时不需要复制。

下面我们来看看类 **DelegatingVehicleTracker**，它没有使用任何显式的同步，所有对状态的访问都由 **ConcurrentHashMap** 来管理，而且 **Map** 所有的键和值都是不可变的。

```java
@ThreadSafe
public class DelegatingVehicleTracker {
    private final ConcurrentHashMap<String, Point> locations;

    private final Map<String, Point> unmodifiableMap;

    public DelegatingVehicleTracker(Map<String, Point> points) {
        locations = new ConcurrentHashMap<String, Point>(points);
        unmodifiableMap = Collections.unmodifiableMap(locations);
    }

    public Map<String, Point> getLocations() {
        return unmodifiableMap;
    }

    public Point getLocation(String id) {
        return locations.get(id);
    }

    public void setLocation(String id, int x, int y) {
        if (locations.replace(id, new Point(x, y)) == null)
            throw new IllegalArgumentException("invalid vehicle name: " + id);
    }
}
```

在2.2中使用监视器模式的车辆追踪器中返回的是**车辆位置的快照**，而在使用委托的车辆追踪器中返回的是一个**不可修改但却实时的车辆位置视图**。

如果需要一个不发生变化的车辆视图，那么 **getLocations** 可以返回对 **locations** 这个 **Map** 对象的一个浅拷贝（**Shallow Copy**）。由于 **Map** 的内容是不可变的，因此只需复制 **Map** 的结构，而不用复制它的内容。

```java
/**
 * 返回 locations 的静态拷贝而非实时拷贝
 */
public Map<String, Point> getLocations() {
    return Collections.unmodifiableMap(new HashMap<String, Point>(locations));
}
```

## 3.2 独立的状态变量

上述示例，线程安全性仅仅委托给单个线程安全的状态变量。我们还可以将线程安全性委托给多个线程安全的状态变量，只要这些变量是彼此独立的，即组合而成的类并不会在其包含的多个状态变量上增加任何不变性条件。

下面我们来看一个代码示例 ：

```java
public class VisualComponent {
    private final List<KeyListener> keyListeners = new CopyOnWriteArrayList<KeyListener>();

    private final List<MouseListener> mouseListeners = new CopyOnWriteArrayList<MouseListener>();

    public void addKeyListener(KeyListener listener) {
        keyListeners.add(listener);
    }

    public void addMouseListener(MouseListener listener) {
        mouseListeners.add(listener);
    }

    public void removeKeyListener(KeyListener listener) {
        keyListeners.remove(listener);
    }

    public void removeMouseListener(MouseListener listener) {
        mouseListeners.remove(listener);
    }
}
```

**VisualComponent** 是一个图形组件，允许客户程序注册监控鼠标和键盘等事件的监听器。它为每种类型的事件都备有一个已注册监听器列表，因此但某个事件发生时，就会调用响应的监听器。在鼠标事件监听器与键盘事件监听器之间不存在任何关联，二者彼此独立，因此 **VisualComponent** 可以将线程安全性委托给这两个线程安全的监听器列表。

**VisualComponent** 使用 **CopyOnWriteArrayList** 来保存各个监听器列表。每个链表都是线程安全的，且各个状态之间不存在耦合关系。

## 3.3 委托失效

大多数组合对象不像 **VisualComponent** 这样简单，它们的状态变量之间存在着一些不变性条件。

下面我们来看一个代码示例：

```java
public class NumberRange {
    // 不变性条件 : lower <= upper
    private final AtomicInteger lower = new AtomicInteger(0);

    private final AtomicInteger upper = new AtomicInteger(0);

    public void setLower(int i) {
        // 注意：不安全的“先检查后执行”
        if (i > upper.get())
            throw new IllegalArgumentException("can't set lower to " + i + " > upper");
        lower.set(i);
    }

    public void setUpper(int i) {
        // 注意：不安全的“先检查后执行”
        if (i > lower.get())
            throw new IllegalArgumentException("can't set upper to " + i + " < lower");
        upper.set(i);
    }

    public boolean isInRange(int i) {
        return (i >= lower.get() && i <= upper.get());
    }
}
```

上面类 **NumberRange** 使用了两个 **AtomicInteger** 来管理状态，并且含有一个约束条件，即下界值 **lower** 要小于或等于上界值 **upper**。很显然，**NumberRange** 不是线程安全的，没有维持对下界和上界进行约束的不变性条件。**setLower** 和 **setUpper** 都是“先检查后执行”的操作，但他们没有使用足够的加锁机制来保证这些操作的原子性。由于状态变量 **lower** 和 **upper** 不是彼此独立的，因此 **NumberRange** 不能将线程安全性委托给它们。 

> 注意：如果一个类是由多个独立且线程安全的状态变量组成，并且在所有的操作中都不包含无效状态转换，那么可以将线程安全性委托给底层的状态变量。

## 3.4 发布底层的状态变量

如果一个状态变量是线程安全的，并且没有任何不变性条件来约束它的值，在变量的操作上也不存在任何不允许的状态转换，那么就可以安全地发布这个变量。例如上文提到的，发布 **VisualComponent** 中的 **mouseListeners** 或 **keyListeners** 等变量就是安全的。

## 3.5 发布状态的车辆追踪器

下面我们还是以车辆追踪器为例，构造另一个版本，并在这个版本中发布底层的可变状态。

为了适应新的版本，位置信息就需要使用可变且线程安全的类 **SafePoint**，如下所示：

```java
@ThreadSafe
public class SafePoint {
    @GuardedBy("this") 
    private int x, y;

    private SafePoint(int[] a) {
        this(a[0], a[1]);
    }

    /**
     * 如果将拷贝构造函数实现为 this(p.x, p.y), 那么会产生竞态条件，而私有构造函数则可以避免这种竞态条件。
     * 这是私有构造函数捕获模式的一个实例。
     */
    public SafePoint(SafePoint p) {
        this(p.get());
    }

    public SafePoint(int x, int y) {
        this.x = x;
        this.y = y;
    }

    public synchronized int[] get() {
        return new int[] {x, y};
    }

    public synchronized void set(int x, int y) {
        this.x = x;
        this.y = y;
    }
}
```

**SafePoint** 提供的 **get** 方法同时获得 **x** 和 **y** 的值，并将二者放在一个数组中返回。如果为 **x** 和 **y** 分别提供 **get** 方法，那么在获得这两个不同坐标的操作之间，**x** 和 **y** 的值发生变化，导致出现车辆从来没有到达过的位置，其线程安全性被破坏了。

下面来看一下安全发布底层状态的车辆追踪器的代码示例：

```java
@ThreadSafe
public class PublishingVehicleTracker {
    private final Map<String, SafePoint> locations;
    private final Map<String, SafePoint> unmodifiableMap;

    public PublishingVehicleTracker(Map<String, SafePoint> locations) {
        this.locations = new ConcurrentHashMap<String, SafePoint>(locations);
        this.unmodifiableMap = Collections.unmodifiableMap(this.locations);
    }

    public Map<String, SafePoint> getLocations() {
        return unmodifiableMap;
    }

    public SafePoint getLocation(String id) {
        return locations.get(id);
    }

    public void setLocation(String id, int x, int y) {
        if (!locations.containsKey(id))
            throw new IllegalArgumentException("invalid vehicle name :" + id);
        locations.get(id).set(x, y);
    }
}
```


**PublishingVehicleTracker** 将其线程安全性委托给底层的 **ConcurrentHashMap**，只是 **Map** 中的元素是线程安全的且可变的 **SafePoint** 类。**getLocation** 方法返回底层 **Map** 对象的一个不可变副本。调用者可以通过修改返回 **Map** 中的 **SafePoint** 值来改变车辆的位置。

# 4. 在现有的线程安全类中添加功能

假设一个线程安全的链表，它需要提供一个“若没有则添加”的操作，而这个操作必须是原子操作，才能保证是线程安全的。

最安全的方法是修改原始的类，在理解原始代码的同步策略的基础上，将新方法添加到类中，更利于理解与维护。

当然如果修改原始的类很难做到，另一种方式是扩展这个类，如下面的代码类 BetterVector 所示：

```java
@ThreadSafe
public class BetterVector<E> extends Vector<E> {
    public synchronized boolean putIfAbsent(E e) {
        boolean isAbsent = !list.contains(e);
        if (isAbsent) 
            add(e);
        return isAbsent;
    }
}
```

由于采用扩展的方法，现在同步策略的实现被分布到多个单独维护的源代码文件中。如果底层的类改变了同步策略并选择了不同的锁来保护它的状态变量，那么子类因为在同步策略改变后无法再使用正确的锁来控制对基类状态的并发访问，从而破坏了子类的线程安全性。（在 **Vector** 的规范中定义了它的同步策略，因此 **BetterVector** 不存在这个问题）

## 4.1 客户端加锁机制

除了上面两种方式：**修改原类** 和 **扩展原类**，第三种方法是扩展类的功能，但并不是扩展类本身，而是将扩展代码放入一个 “辅助类” 中。

如下实现了一个包含 “若没有则添加” 操作的辅助类，但它是 **非线程安全的，慎用！！！**

```java
@NotThreadSafe
public class ListHelper<E> {
    public List<E> list = Collections.synchronizedList(new ArrayList<E>());

    public synchronized boolean putIfAbsent(E e) {
        boolean isAbsent = !list.contains(e);
        if (isAbsent) 
            list.add(e);
        return isAbsent;
    }
}
```	

那么为什么说 **ListHelper** 是非线程安全的？**putIfAbsent** 方法也声明了 **synchronized** 关键字，是不是？

这里的关键是 **putIfAbsent** 在错误的锁上进行了同步。可以明确的是，无论 **List** 使用哪一个锁来保护它的状态，这个锁一定不是 **ListHelper** 上的锁。

要想使这个方法可以正确地执行，必须使 **List** 在实现客户端加锁 或 外部加锁时使用同一个锁。客户端加锁是指，对于使用某个对象 **A** 的客户端代码，使用 **A** 本身用于保护其状态的锁来保护这段客户代码。当然要使用客户端锁，也就必须要知道对象 **A** 使用的是哪一个锁。

在 **Vector** 和 同步封装器类的文档中指出，它们通过使用 **Vector** 或 封装器容器的内置锁来支持客户端加锁。

下面我们来看一个代码示例，通过客户端加锁来实现 “若没有则添加” ：

```java
@ThreadSafe
public class ListHelper<E> {
    public List<E> list = Collections.synchronizedList(new ArrayList<E>());

    public boolean putIfAbsent(E e) {
        synchronized (list) {
            boolean isAbsent = !list.contains(e);
            if (isAbsent) 
                list.add(e);
            return isAbsent;
        }
    }
}
```

客户端加锁机制与扩展类机制有许多共同点，二者都是将派生类的行为与基类的实现耦合在一起。正如扩展会破坏实现的封装性，客户端加锁同样会破坏同步策略的封装性。

## 4.2 组合

当为现有的类添加一个原子操作时，有一种更好的方法：**组合（Composition）**。

下面我们来看一个代码示例，通过组合实现 “若没有则添加” ：

```java
@ThreadSafe
public class ImprovedList<T> implements List<T> {
    public final List<T> list;

    public ImprovedList(List<T> list) {
        this.list = list;
    }

    public synchronized boolean putIfAbsent(E e) {
        boolean isAbsent = !contains(e);
        if (isAbsent) 
            list.add(e);
        return isAbsent;
    }

    public synchronized void clear() {
        list.clear();
    }

    // ... 按照类似的方法委托 List 的其他方法
}
```

**ImprovedList** 通过自身的内置锁增加了一层额外的加锁。它并不关心底层的 **List** 是否是线程安全，即使 **List** 不是线程安全或者修改了它的加锁实现，**ImprovedList** 也会提供一致的加锁机制来实现线程安全性。事实上，我们使用了 **Java 监视器模式** 来封装现有的 **List**，并且只要在类中拥有指向底层 **List** 的唯一外部引用，就能确保线程安全性。


# 总结
本篇介绍了一些组合模式，可以很方便地保护类的线程安全性。下一篇我们将要学习Java类库中的并发基础构建模块，敬请期待！



