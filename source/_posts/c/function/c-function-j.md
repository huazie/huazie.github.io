---
title: C语言函数大全--j开头的函数
date: 2023-04-13 16:19:47
updated: 2025-02-21 20:15:12
categories:
  - 开发语言-C
tags:
  - C语言函数大全
  - j开头的函数
---

![](/images/cplus-logo.png)

# 总览
| 函数声明 |  函数功能  |
|:--|:--|
|`double j0 (double x);` | 计算 x 的 **第一类 0 阶贝塞尔函数**（double）  |
|`float j0f (float x);` |  计算 x 的 **第一类 0 阶贝塞尔函数**（float）【笔者本地windows环境，无此函数】 |
|`double j1 (double x);` | 计算 x 的 **第一类 1 阶贝塞尔函数**（double）|
|`float j1f (float x);` |  计算 x 的 **第一类 1 阶贝塞尔函数**（float）【笔者本地windows环境，无此函数】|
|`double jn (int n, double x);` | 计算 **x** 的 **第一类 n 阶贝塞尔函数**（double）|
|`float jnf (int n, float x);` | 计算 **x** 的 **第一类 n 阶贝塞尔函数**（float）【笔者本地windows环境，无此函数】  |
|`double jrand48();` | 生成伪随机数序列  |
|`int join(pthread_t thread, void **retval);` |  等待线程退出并回收资源 |
|`typedef _JBTYPE jmp_buf[_JBLEN];` | 它是一个数组类型，保存跳转目标地址的缓冲区。通常与 setjmp 和 longjmp 函数一起使用，用于实现非局部跳转 |
|`u32 jhash(const void *key, u32 length, u32 initval);` | 它是 Linux 内核头文件 `linux/jhash.h` 中的一个函数，用于实现一种高效的哈希算法。  |
|`unsigned long volatile jiffies;` | 它是 Linux 内核中的一个全局变量，表示内核启动后经过的节拍数。其中 `volatile` 关键字用于告知编译器在访问这个变量时不要使用缓存，以确保能够正确读取最新值。  |
|`u64 jiffies_64;`|它是 Linux 内核中的一个全局变量，类似于 `jiffies`，但是支持更大的取值范围。其中 `u64` 是 **64** 位无符号整型。|
|`clock_t jiffies_delta_to_clock_t(unsigned long delta);` | 它是 Linux 内核头文件 linux/jiffies.h 中的一个函数，用于计算两个 jiffies 值之间的时间差，并将结果转换为 clock_t 类型的值。  |
|`unsigned long jiffies_delta_to_msecs(unsigned long delta);` | 它是 Linux 内核头文件 `linux/jiffies.h` 中的一个函数，用于计算两个 **jiffies** 值之间的时间差，并将结果转换为毫秒数  |
| `clock_t jiffies_to_clock_t(unsigned long jiffies);`|  它是 Linux 内核头文件 linux/jiffies.h 中的一个函数，用于将 jiffies 值（内核节拍数）转换为 clock_t 类型的值。 |
| `unsigned long jiffies_to_msecs(const unsigned long j);`|  它是 Linux 内核头文件 linux/jiffies.h 中的一个函数，用于将 jiffies 值（内核节拍数）转换为毫秒数。 |
|`clock_t jiffies64_to_clock_t(u64 jiffies);` | 它是 Linux 内核头文件 linux/jiffies.h 中的一个函数，用于将 64 位 jiffies 值（内核节拍数）转换为 clock_t 类型的值。  |
| `u64 jiffies64_to_msecs(const u64 jiffies);`|  它是 Linux 内核头文件 linux/jiffies.h 中的一个函数，用于将 64 位 jiffies 值（内核节拍数）转换为毫秒数。 |
| `void jiffies_update_wallclock(void);`| 它是 Linux 内核头文件 `linux/time.h` 中的一个函数，用于更新系统时钟的时间戳。  |


# 1. j0，j0f
## 1.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double j0 (double x);` | 计算 x 的 **第一类 0 阶贝塞尔函数**（double）  |
|`float j0f (float x);` |  计算 x 的 **第一类 0 阶贝塞尔函数**（float）【笔者本地windows环境，无此函数】 |

> **注意：** 如果操作成功，则返回 **x** 的 **第一类 0 阶贝塞尔函数**；如果 **x** 是 **NaN** 值，则返回 **NaN** 值；如果 **x** 太大或发生溢出范围错误，则返回 **0** 并将 **errno** 设置为 **ERANGE**。

## 1.2 演示示例
```c
#include <stdio.h>
#include <math.h>
int main()
{
    double x = 10.0, result;
    result = j0(x);

    printf("%lf 的 第一类 0 阶贝塞尔函数 : %lf", x, result);

    return 0;
} 
```

## 1.3 运行结果
![](j0.png)

# 2. j1，j1f
## 2.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double j1 (double x);` | 计算 x 的 **第一类 1 阶贝塞尔函数**（double）|
|`float j1f (float x);` |  计算 x 的 **第一类 1 阶贝塞尔函数**（float）【笔者本地windows环境，无此函数】|

**注意：** 如果操作成功，则返回 **x** 的 **第一类 1 阶贝塞尔函数**；如果 **x** 是 **NaN** 值，则返回 **NaN** 值；如果 **x** 太大或发生溢出范围错误，则返回 **0** 并将 **errno** 设置为 **ERANGE**。

## 2.2 演示示例
```c
#include <stdio.h>
#include <math.h>
int main()
{
    double x = 10.0, result;
    result = j1(x);

    printf("%lf 的 第一类 1 阶贝塞尔函数 : %lf", x, result);

    return 0;
} 
```

## 2.3 运行结果
![](j1.png)

# 3. jn，jnf
## 3.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double jn (int n, double x);` | 计算 **x** 的 **第一类 n 阶贝塞尔函数**（double）|
|`float jnf (int n, float x);` | 计算 **x** 的 **第一类 n 阶贝塞尔函数**（float）【笔者本地windows环境，无此函数】  |

> **注意：** 如果操作成功，则返回 **x** 的 **第一类 n 阶贝塞尔函数**；如果 **x** 是 **NaN** 值，则返回 **NaN** 值；如果 **x** 太大或发生溢出范围错误，则返回 **0** 并将 **errno** 设置为 **ERANGE**。

## 3.2 演示示例
```c
#include <stdio.h>
#include <math.h>

void jnPrint(int n, double x);

int main()
{
    double x = 10.0;
    jnPrint(2, x);
    jnPrint(3, x);
    jnPrint(4, x);
    return 0;
} 

void jnPrint(int n, double x)
{
    double result = jn(n, x);
    printf("%lf 的 第一类 %d 阶贝塞尔函数 : %lf\n", x, n, result);
}
```

## 3.3 运行结果
![](jn.png)

# 4. jrand48
## 4.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double jrand48();` | 生成伪随机数序列  |

> **jrand48** 函数是一个生成伪随机数序列的函数，并且它是可重入的，即可以在多个线程中同时调用而不会出现冲突。

## 4.2 演示示例
```c
#include <stdlib.h>
#include <stdio.h>
#include <time.h>

int main() {
    // 初始化种子
    srand48(time(NULL));
    
    // 生成10个随机数
    for (int i = 0; i < 5; ++i) {
        double r = jrand48();
        printf("%f\n", r);
    }
    
    return 0;
}
```
上述程序首先通过 **srand48** 函数初始化随机数生成器的种子，这里使用了当前系统时间作为种子。然后循环调用 **jrand48** 函数 **5** 次，每次输出一个伪随机数。注意，由于 **jrand48** 函数返回的是一个双精度浮点数（范围在 **[0, 1)** 内），因此输出时需要使用 **%f** 格式化符号。

# 5. join
## 5.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int join(pthread_t thread, void **retval);` |  等待线程退出并回收资源 |

> 在 **C** 语言中，**join** 函数不是标准库函数，也不是 **POSIX** 标准的函数。然而，一些操作系统（如 **UNIX/Linux**）提供了 **join** 函数用于等待线程退出并回收资源。在 **POSIX** 线程中，相应的函数是 **pthread_join**。

## 5.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>

void *thread_func(void *arg) {
    printf("Thread is running...\n");
    pthread_exit(NULL);
}

int main() {
    pthread_t thread;
    if (pthread_create(&thread, NULL, thread_func, NULL)) {
        perror("pthread_create");
        exit(EXIT_FAILURE);
    }
    printf("Main thread is waiting for the child thread to exit...\n");
    join(thread, NULL);
    printf("Child thread has exited.\n");
    return EXIT_SUCCESS;
}
```

上述程序创建了一个新线程，并且主线程等待新线程退出后才继续执行。在新线程中，打印一条消息并调用 **pthread_exit** 函数退出线程。在主线程中，调用 **join** 函数等待新线程退出，并通过 **NULL** 参数指示不需要返回值。最终输出一条消息表示新线程已经退出。

# 6. jmp_buf 
## 6.1 类型说明
| 类型定义 |  描述  |
|:--|:--|
|`typedef _JBTYPE jmp_buf[_JBLEN];` | 它是一个数组类型，保存跳转目标地址的缓冲区。通常与 setjmp 和 longjmp 函数一起使用，用于实现非局部跳转 |

## 6.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <setjmp.h>

jmp_buf env;

void func() {
    printf("Entering func...\n");
    longjmp(env, 1);
    printf("This line will not be executed.\n");
}

int main() {
    int ret = setjmp(env);
    if (ret == 0) {
        printf("Jumping to func...\n");
        func();
    } else {
        printf("Returning from longjmp with value %d\n", ret);
    }
    return EXIT_SUCCESS;
}
```
上述程序定义了一个名为 **env** 的 **jmp_buf** 类型变量，用于保存当前执行状态。在主函数中，通过调用 **setjmp** 函数将当前状态保存到
 **env** 中，并返回 **0**。然后调用 **func** 函数，该函数打印一条消息并调用 **longjmp** 函数恢复之前保存的状态，这里传入参数值为 **1**。由于 **longjmp** 函数会导致程序跳转到 **setjmp** 函数继续执行，因此后面的 **printf** 语句会输出 `"Returning from longjmp with value 1"`。

> **需要注意的是**，在使用 **jmp_buf**、**setjmp** 和 **longjmp** 函数时需要遵循特定的使用规范，否则可能会导致未定义行为或错误。

## 6.3 运行结果
![](jmp_buf.png)

# 7. jhash
## 7.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`u32 jhash(const void *key, u32 length, u32 initval);` | 它是 Linux 内核头文件 `linux/jhash.h` 中的一个函数，用于实现一种高效的哈希算法。  |

> **参数：**
> - **key ：** 要进行哈希的数据
> - **length ：** 数据的长度（以字节为单位）
> - **initval ：** 哈希值的初始值。

## 7.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/jhash.h>

int my_init(void)
{
    char data[] = "Hello, world!";
    u32 hash;

    printk(KERN_INFO "Initializing module...\n");

    /* calculate hash value */
    hash = jhash(data, strlen(data), 0);
    printk(KERN_INFO "Hash value: %u\n", hash);

    return 0;
}

void my_exit(void)
{
    printk(KERN_INFO "Exiting module...\n");
}

MODULE_LICENSE("GPL");
module_init(my_init);
module_exit(my_exit);

```
上述示例程序中，在 `my_init()` 函数中定义了一个字符串 **data**，并使用 `jhash()` 函数计算出其哈希值，并打印出来。

> **注意：** 虽然 `jhash()` 函数可以用于快速查找和管理数据结构等，但在使用时必须充分理解其作用原理和使用方法，避免因为错误使用导致哈希冲突或其他问题。同时，应当根据具体情况选择合适的哈希算法，并考虑其效率和安全性等方面的因素。

# 8. jiffies，jiffies_64
## 8.1 变量说明
| 变量声明 |  变量描述  |
|:--|:--|
|`unsigned long volatile jiffies;` | 它是 Linux 内核中的一个全局变量，表示内核启动后经过的节拍数。其中 `volatile` 关键字用于告知编译器在访问这个变量时不要使用缓存，以确保能够正确读取最新值。  |
|`u64 jiffies_64;`|它是 Linux 内核中的一个全局变量，类似于 `jiffies`，但是支持更大的取值范围。其中 `u64` 是 **64** 位无符号整型。|

## 8.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/jiffies.h>

int my_init(void)
{
    u64 start = jiffies_64;
    u64 end;

    printk(KERN_INFO "Initializing module...\n");

    /* do some work */
    mdelay(1000);

    /* calculate elapsed time in jiffies_64 */
    end = jiffies_64 - start;
    printk(KERN_INFO "Elapsed time: %llu jiffies_64\n", end);

    return 0;
}

void my_exit(void)
{
    printk(KERN_INFO "Exiting module...\n");
}

MODULE_LICENSE("GPL");
module_init(my_init);
module_exit(my_exit);

```
上述示例程序中，在 `my_init()` 函数中获取了当前的 **jiffies_64** 值，并使用 `mdelay()` 函数让程序阻塞一段时间。在之后，通过计算当前的 `jiffies_64` 值与之前的 `jiffies_64` 值之差，计算出经过的时间并打印出来。

> **注意：** `jiffies` 和 `jiffies_64` 值每隔一段时间就会发生一次溢出，在处理 `jiffies` 和 `jiffies_64` 值时必须注意这个问题，避免计算结果错误。另外，`jiffies` 和 `jiffies_64` 变量只能在内核空间中使用，不能在用户空间中使用。

# 9. jiffies_delta_to_clock_t
## 9.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`clock_t jiffies_delta_to_clock_t(unsigned long delta);` | 它是 Linux 内核头文件 linux/jiffies.h 中的一个函数，用于计算两个 jiffies 值之间的时间差，并将结果转换为 clock_t 类型的值。  |

> **参数：**
> - **delta ：** 要计算的 jiffies 时间差值。

## 9.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/jiffies.h>

int my_init(void)
{
    unsigned long start = jiffies;
    unsigned long end;
    clock_t ticks;

    printk(KERN_INFO "Initializing module...\n");

    /* do some work */
    mdelay(1000);

    /* calculate elapsed time in jiffies and convert to ticks */
    end = jiffies;
    ticks = jiffies_delta_to_clock_t(end - start);
    printk(KERN_INFO "Elapsed time: %ld ticks\n", (long)ticks);

    return 0;
}

void my_exit(void)
{
    printk(KERN_INFO "Exiting module...\n");
}

MODULE_LICENSE("GPL");
module_init(my_init);
module_exit(my_exit);
```
上述示例程序中，在 `my_init()` 函数中获取了当前的 `jiffies` 值，并使用 `mdelay()` 函数让程序阻塞一段时间。在之后，通过计算当前的 `jiffies` 值与之前的 `jiffies` 值之差，并调用 `jiffies_delta_to_clock_t()` 函数，计算出经过的时间并打印出来。

> **注意：** 在使用 `jiffies_delta_to_clock_t()` 函数时，返回值类型是 `clock_t`，不同于 `jiffies_delta_to_msecs()` 函数的返回值类型是 `unsigned long`。另外，`clock_t` 的定义可能因系统而异，应当根据具体情况进行处理。

# 10. jiffies_delta_to_msecs
## 10.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`unsigned long jiffies_delta_to_msecs(unsigned long delta);` | 它是 Linux 内核头文件 `linux/jiffies.h` 中的一个函数，用于计算两个 **jiffies** 值之间的时间差，并将结果转换为毫秒数  |

> **参数：**
> - **delta ：** 要计算的 jiffies 时间差值。

## 10.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/jiffies.h>

int my_init(void)
{
    unsigned long start = jiffies;
    unsigned long end;
    unsigned long msecs;

    printk(KERN_INFO "Initializing module...\n");

    /* do some work */
    mdelay(1000);

    /* calculate elapsed time in jiffies and convert to milliseconds */
    end = jiffies;
    msecs = jiffies_delta_to_msecs(end - start);
    printk(KERN_INFO "Elapsed time: %lu ms\n", msecs);

    return 0;
}

void my_exit(void)
{
    printk(KERN_INFO "Exiting module...\n");
}

MODULE_LICENSE("GPL");
module_init(my_init);
module_exit(my_exit);
```
上述示例程序中，在 `my_init()` 函数中获取了当前的 `jiffies` 值，并使用 `mdelay()` 函数让程序阻塞一段时间。在之后，通过计算当前的 `jiffies` 值与之前的 `jiffies` 值之差，并调用 `jiffies_delta_to_msecs()` 函数，计算出经过的时间并打印出来。

> **注意：** 在使用 jiffies_delta_to_msecs() 函数时，返回值类型是 unsigned long，不同于 jiffies_delta_to_clock_t() 函数的返回值类型是 clock_t。另外，由于 jiffies 的精度限制，计算结果可能存在一定的误差。

# 11. jiffies_to_clock_t
## 11.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `clock_t jiffies_to_clock_t(unsigned long jiffies);`|  它是 Linux 内核头文件 linux/jiffies.h 中的一个函数，用于将 jiffies 值（内核节拍数）转换为 clock_t 类型的值。 |

> **参数：**
> - **jiffies：** 要转换的 **jiffies** 值，它是 **Linux** 内核中的一个全局变量，表示内核启动后经过的节拍数。

## 11.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/jiffies.h>

int my_init(void)
{
    unsigned long j = jiffies;
    clock_t ticks;

    printk(KERN_INFO "Initializing module...\n");

    /* do some work */
    mdelay(1000);

    /* calculate elapsed time in ticks */
    ticks = jiffies_to_clock_t(jiffies - j);
    printk(KERN_INFO "Elapsed time: %ld ticks\n", (long)ticks);

    return 0;
}

void my_exit(void)
{
    printk(KERN_INFO "Exiting module...\n");
}

MODULE_LICENSE("GPL");
module_init(my_init);
module_exit(my_exit);
```
上述示例程序中，在 `my_init()` 函数中获取了当前的 **jiffies** 值，并使用 `mdelay()` 函数让程序阻塞一段时间。在之后，通过计算当前的 **jiffies** 值与之前的 **jiffies** 值之差并调用 `jiffies_to_clock_t()` 函数，计算出经过的时间，并打印出来。

> **注意：** 在使用 `jiffies_to_clock_t()` 函数时，返回值类型是 `clock_t`，不同于 `jiffies_to_msecs()` 函数的返回值类型是 `unsigned long`。另外，`clock_t` 的定义可能因系统而异，应当根据具体情况进行处理。

# 12. jiffies_to_msecs
## 12.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `unsigned long jiffies_to_msecs(const unsigned long j);`|  它是 Linux 内核头文件 linux/jiffies.h 中的一个函数，用于将 jiffies 值（内核节拍数）转换为毫秒数。 |

> **参数：**
> - **j：** 要转换的 **jiffies** 值，它是 **Linux** 内核中的一个全局变量，表示内核启动后经过的节拍数。

## 12.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/jiffies.h>

int my_init(void)
{
    unsigned long j = jiffies;
    unsigned long ms;

    printk(KERN_INFO "Initializing module...\n");

    /* do some work */
    mdelay(1000);

    /* calculate elapsed time in milliseconds */
    ms = jiffies_to_msecs(jiffies - j);
    printk(KERN_INFO "Elapsed time: %lu ms\n", ms);

    return 0;
}

void my_exit(void)
{
    printk(KERN_INFO "Exiting module...\n");
}

MODULE_LICENSE("GPL");
module_init(my_init);
module_exit(my_exit);
```
上述示例程序中，在 `my_init()` 函数中获取了当前的 **jiffies** 值，并使用 `mdelay()` 函数让程序阻塞一段时间。在之后，通过计算当前的 **jiffies** 值与之前的 **jiffies** 值之差并调用 `jiffies_to_msecs()` 函数，计算出经过的时间，并打印出来。

> **注意：** 在使用 `jiffies_to_msecs()` 函数时，必须十分小心地处理 **jiffies** 值的溢出等问题，以免计算结果错误。另外，`jiffies_to_msecs()` 函数只能用于内核空间中，不能在用户空间中使用。


# 13. jiffies64_to_clock_t
## 13.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`clock_t jiffies64_to_clock_t(u64 jiffies);` | 它是 Linux 内核头文件 linux/jiffies.h 中的一个函数，用于将 64 位 jiffies 值（内核节拍数）转换为 clock_t 类型的值。  |

> **参数：**
> - **jiffies ：** 要转换的 **64** 位 **jiffies** 值。

## 13.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/jiffies.h>

int my_init(void)
{
    u64 start = jiffies_64;
    u64 end;
    clock_t ticks;

    printk(KERN_INFO "Initializing module...\n");

    /* do some work */
    mdelay(1000);

    /* calculate elapsed time in jiffies_64 and convert to ticks */
    end = jiffies_64;
    ticks = jiffies64_to_clock_t(end - start);
    printk(KERN_INFO "Elapsed time: %ld ticks\n", (long)ticks);

    return 0;
}

void my_exit(void)
{
    printk(KERN_INFO "Exiting module...\n");
}

MODULE_LICENSE("GPL");
module_init(my_init);
module_exit(my_exit);
```
上述示例程序中，在 `my_init()` 函数中获取了当前的 **64** 位 `jiffies_64` 值，并使用 `mdelay()` 函数让程序阻塞一段时间。在之后，通过计算当前的 **64** 位 `jiffies_64` 值与之前的 `jiffies_64` 值之差，并调用 `jiffies64_to_clock_t()` 函数，计算出经过的时间并打印出来。

# 14. jiffies64_to_msecs
## 14.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `u64 jiffies64_to_msecs(const u64 jiffies);`|  它是 Linux 内核头文件 linux/jiffies.h 中的一个函数，用于将 64 位 jiffies 值（内核节拍数）转换为毫秒数。 |
> **参数：**
> - **jiffies ：** 要转换的 **64** 位 **jiffies** 值。

## 14.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/jiffies.h>

int my_init(void)
{
    u64 start = jiffies_64;
    u64 end;
    u64 msecs;

    printk(KERN_INFO "Initializing module...\n");

    /* do some work */
    mdelay(1000);

    /* calculate elapsed time in jiffies_64 and convert to milliseconds */
    end = jiffies_64;
    msecs = jiffies64_to_msecs(end - start);
    printk(KERN_INFO "Elapsed time: %llu ms\n", msecs);

    return 0;
}

void my_exit(void)
{
    printk(KERN_INFO "Exiting module...\n");
}

MODULE_LICENSE("GPL");
module_init(my_init);
module_exit(my_exit);
```
上述示例程序中，在 `my_init()` 函数中获取了当前的 **64** 位 `jiffies_64` 值，并使用 `mdelay()` 函数让程序阻塞一段时间。在之后，通过计算当前的 **64** 位 `jiffies_64` 值与之前的 `jiffies_64` 值之差，并调用 `jiffies64_to_msecs()` 函数，计算出经过的时间并打印出来。

# 15. jiffies_update_wallclock
## 15.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `void jiffies_update_wallclock(void);`| 它是 Linux 内核头文件 `linux/time.h` 中的一个函数，用于更新系统时钟的时间戳。  |

## 15.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/jiffies.h>

int my_init(void)
{
    printk(KERN_INFO "Initializing module...\n");

    /* do some work */
    mdelay(1000);

    /* update wall clock */
    jiffies_update_wallclock();

    return 0;
}

void my_exit(void)
{
    printk(KERN_INFO "Exiting module...\n");
}

MODULE_LICENSE("GPL");
module_init(my_init);
module_exit(my_exit);
```
上述示例程序中，在 `my_init()` 函数中使用 `mdelay()` 函数让程序阻塞一段时间，在之后调用 `jiffies_update_wallclock()` 函数更新系统时钟的时间戳。

> **注意：** `jiffies_update_wallclock()` 函数只能在内核空间中使用，不能在用户空间中使用。另外，如果系统使用了 **NTP** 等网络时间同步服务，可能无法通过 `jiffies_update_wallclock()` 函数来准确更新系统时钟。

# 参考
1. [\[MATH-标准C库\]](https://device.harmonyos.com/cn/docs/documentation/apiref/math-0000001055228010#ZH-CN_TOPIC_0000001055228010__gaffb00730a1127dee798137075951ae21)
2. 《Linux内核API完全参考手册》
