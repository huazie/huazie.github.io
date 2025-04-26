---
title: C语言函数大全--k开头的函数
date: 2023-04-14 14:43:51
updated: 2025-04-17 19:50:40
categories:
  - 开发语言-C
tags:
  - C语言函数大全
  - k开头的函数
---

![](/images/cplus-logo.png)

# 总览
| 函数声明 |  函数功能  |
|:--|:--|
|`void *kcalloc(size_t n, size_t size, gfp_t flags);` |  它是 Linux 内核中的一个函数，用于在内核空间分配一块连续的指定大小的内存，它与标准库函数 calloc() 的功能类似。 |
| `int kbhit(void);`| 在控制台中检测是否有按键被按下  |
|`void keep(void *ptr);` |  它是 Linux 内核中的一个函数，用于防止编译器将指定的符号优化掉。 |
|`asmlinkage int kernel_thread(int (*fn)(void *), void *arg, unsigned long flags);` | 它是 Linux 内核中的一个函数，用于在内核空间中创建一个新进程。  |
|`void kfree(void *ptr);` |  它是 Linux 内核中的一个函数，用于释放使用 kmalloc() 或者 kzalloc() 函数分配的内存空间。 |
| `int kill(pid_t pid, int sig);`|  向指定进程或进程组发送一个信号 |
|`int kill_proc(pid_t pid, int sig, int priv);` |   它是 Linux 内核中的一个函数，用于向指定进程发送信号。|
|`void *kmalloc(size_t size, gfp_t flags);` |  它是 Linux 内核中的一个函数，用于在内核中分配指定大小的内存空间。 |
|`void *kmap(struct page *page);`|它是 Linux 内核中的一个函数，用于将一个页映射到内核虚拟地址空间。|
|`void *kmap_high(struct page *page);` |  它是 Linux 内核中的一个函数，用于将高端内存映射到内核虚拟地址空间中。 |
| `void *kmem_cache_alloc(struct kmem_cache *cachep, int flags);`| 它是 Linux 内核中的一个函数，用于从指定的内存缓存中分配一个对象  |
| `struct kmem_cache *kmem_cache_create(const char *name, size_t size, size_t align, unsigned long flags, void (*ctor)(void *));`|  它是 Linux 内核中的一个函数，用于创建一个内存缓存区，可以用于高效地分配和释放指定大小的对象。 |
|`void kmem_cache_free(struct kmem_cache *cachep, void *objp);`|它是 Linux 内核中的一个函数，用于将之前使用 kmem_cache_alloc() 函数分配的对象释放回内存缓存池，以便下次再次分配使用。|
|`void kmem_cache_destroy(struct kmem_cache *cachep);`|它是 Linux 内核中的一个函数，用于销毁之前使用 kmem_cache_create() 函数创建的内存缓存区。|
|`void *kmem_cache_zalloc(struct kmem_cache *cache, gfp_t flags);`|它是 Linux 内核中的一个函数，用于从指定内存缓存区中分配一块指定大小的内存，并将其清零。|
|`void *kmemdup(const void *src, size_t len, gfp_t flags);` | 它是 Linux 内核中的一个函数，用于在内核空间中将一段指定大小的内存复制到另一段新的内存中，并返回这段新内存的指针。  |
|`void kprintf(const char *format, ...);` |  用于嵌入式系统中输出调试信息 |
| `void *krealloc(const void *ptr, size_t new_size, gfp_t flags);`|  它是 Linux 内核中的一个函数，用于动态调整已分配内存块的大小。 |
| `size_t ksize(const void *ptr);`| 它是 Linux 内核中的一个函数，用于获取已分配内存块的大小。  |
| `char *kstrdup(const char *s, gfp_t flags);`|  它是 Linux 内核中的一个函数，用于在内核空间中复制一个以 NULL 结尾的字符串，并返回这个新的字符串指针。 |
|`char *kstrndup(const char *s, size_t len, gfp_t flags);` |  它是 Linux 内核中的一个函数，用于在内核空间中复制一个以 NULL 结尾的字符串的一部分，并返回这个新的字符串指针。 |
|`void kstat_irqs_cpu(int cpu, int *irqs, unsigned long *stime);` |  它是 Linux 内核中的一个函数，用于查询指定 CPU 的中断统计信息。 |
|`struct task_struct *kthread_create(int (*threadfn)(void *data), void *data, const char *namefmt, ...);` | 它是 Linux 内核中的一个函数，用于创建一个内核线程。  |
|`int kthread_stop(struct task_struct *k);` |  它是 Linux 内核中的一个函数，用于停止由 kthread_create() 函数创建的内核线程。 |
| `void kunmap_high(struct page *page);`|  它是 Linux 内核中的一个函数，用于取消一个高端内存映射。 |
|`void *kzalloc(size_t size, gfp_t flags);` | 它是 Linux 内核中的一个函数，用于分配指定大小的内存空间，并将其初始化为零。  |


# 1. kcalloc
## 1.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void *kcalloc(size_t n, size_t size, gfp_t flags);` |  它是 Linux 内核中的一个函数，用于在内核空间分配一块连续的指定大小的内存，它与标准库函数 calloc() 的功能类似。 |

**参数：**
- **n ：** 要分配的元素个数
- **size ：** 每个元素的大小
- **flags ：** 用于控制内存分配行为的标志

## 1.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>

int my_init(void)
{
    int *arr;

    printk(KERN_INFO "Initializing module...\n");

    /* allocate and initialize array */
    arr = kcalloc(10, sizeof(int), GFP_KERNEL);
    if (!arr) {
        printk(KERN_ERR "Failed to allocate memory\n");
        return -ENOMEM;
    }
    for (int i = 0; i < 10; i++) {
        arr[i] = i + 1;
    }

    /* print array */
    printk(KERN_INFO "Array contents:\n");
    for (int i = 0; i < 10; i++) {
        printk(KERN_INFO "%d ", arr[i]);
    }
    printk(KERN_INFO "\n");

    /* free memory */
    kfree(arr);

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
上述示例程序中，在 `my_init()` 函数中使用 `kcalloc()` 函数分配了一个大小为 40 字节的整型数组，并将其初始化为 **1 到 10** 的连续整数。在之后，打印了数组内容并释放了内存。

> **注意：** 在使用 `kcalloc()` 函数时，必须确保请求的内存大小不会超过系统可用的物理内存大小，并且可以正确地处理内存分配失败等异常情况。另外，分配的内存应在不再需要时及时释放，以免造成内存泄漏等问题。


# 2. kbhit
## 2.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `int kbhit(void);`| 在控制台中检测是否有按键被按下  |

> 如果有按键被按下，该函数返回非零值，否则返回 0。

## 2.2 演示示例
```c
#include <stdio.h>
#include <conio.h>

int main()
{
    int ch;

    printf("Press any key to continue...\n");
    while (!kbhit()) {
        // 等待用户按键
    }
    ch = getch(); // 获取用户按下的键值
    printf("You pressed the '%c' key\n", ch);

    return 0;
}
```

## 2.3 运行结果
![](kbhit.png)

# 3. keep
## 3.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void keep(void *ptr);` |  它是 Linux 内核中的一个函数，用于防止编译器将指定的符号优化掉。 |

**参数：**
- **ptr：** 是指向要保留的符号的指针。

## 3.2 演示示例

```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>

static int my_symbol __attribute__((__used__));

int my_init(void)
{
    printk(KERN_INFO "Initializing module...\n");

    my_symbol = 123;

    /* do something with my_symbol */

    keep(&my_symbol); // 保留符号

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

上述示例程序中，在 `my_init()` 函数中定义了一个整型变量 **my_symbol**，并且对其进行了初始化。然后，在处理完该变量之后，调用 `keep()` 函数保留该符号，以避免编译器将其优化掉。

> **注意：** 在使用 `keep()` 函数时，必须确保所保留的符号不会被优化掉，否则可能会导致程序出错或崩溃。另外，由于 `keep()` 函数只是防止编译器优化符号，并不会改变其可见性或访问权限，因此在使用该函数时，应该确保所保留的符号在需要的位置上是可见和可访问的。

# 4. kernel_thread
## 4.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`asmlinkage int kernel_thread(int (*fn)(void *), void *arg, unsigned long flags);` | 它是 Linux 内核中的一个函数，用于在内核空间中创建一个新进程。  |

**参数：**
- **fn ：** 指向线程处理函数的指针
- **arg ：** 传递给线程处理函数的参数
- **flags ：** 用于控制进程创建方式的标志。

## 4.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/sched.h>

static int my_thread_func(void *data)
{
    printk(KERN_INFO "my_thread_func started\n");

    /* do something in the thread */

    printk(KERN_INFO "my_thread_func finished\n");
    return 0;
}

int my_init(void)
{
    printk(KERN_INFO "Initializing module...\n");

    kernel_thread(my_thread_func, NULL, CLONE_KERNEL);

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

上述示例程序中，在 `my_init()` 函数中使用 `kernel_thread()` 函数创建了一个新进程，并将其入口函数设置为 `my_thread_func()`。该进程可以执行任何需要在内核空间中进行的操作。在 `my_exit()` 函数中，什么也没做。

需要注意的是，在使用 `kernel_thread()` 函数时，必须确保线程处理函数能够正确地完成自己的任务并且及时退出，否则可能会导致内核崩溃或其他问题。另外，使用内核进程时需要特别小心，因为它们与内核数据结构和操作高度相关，并且可能会影响系统的稳定性和安全性。

# 5. kfree
## 5.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void kfree(void *ptr);` |  它是 Linux 内核中的一个函数，用于释放使用 kmalloc() 或者 kzalloc() 函数分配的内存空间。 |

**参数：** 
- **ptr ：** 指向要释放的内存块的指针。

## 5.2 演示示例
**参考 7.2 所示**

> **注意：** 在使用 `kfree()` 函数释放内存时，必须确保所释放的内存是由 `kmalloc()` 或者 `kzalloc()` 函数分配的，否则可能会导致内核崩溃或其他问题。另外，使用 `kfree()` 函数释放一个指针之后，应该将其设置为 `NULL`，以避免出现悬挂指针（**dangling pointer**）等问题。

# 6. kill
## 6.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `int kill(pid_t pid, int sig);`|  向指定进程或进程组发送一个信号 |

**参数：**
- **pid：** 目标进程的 ID（进程ID或进程组ID）
- **sig：** 要发送的信号编号

## 6.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <signal.h>
#include <unistd.h>

void my_handler(int signum)
{
    printf("Received signal %d\n", signum);
}

int main()
{
    pid_t pid;
    int ret;

    pid = fork(); // 创建子进程
    if (pid == 0) {
        // 子进程
        printf("Child process started\n");
        sleep(10); // 等待父进程发送信号
        exit(0);
    } else if (pid > 0) {
        // 父进程
        printf("Parent process started\n");
        sleep(5); // 等待子进程创建完毕
        ret = kill(pid, SIGINT); // 向子进程发送 SIGINT 信号
        if (ret == -1) {
            perror("kill failed");
            exit(EXIT_FAILURE);
        }
        printf("Signal sent successfully\n");
        wait(NULL); // 等待子进程结束
        printf("Child process finished\n");
        exit(EXIT_SUCCESS);
    } else {
        perror("fork failed");
        exit(EXIT_FAILURE);
    }

    return 0;
}
```

上述示例程序中，我们首先创建一个子进程，并在子进程中等待 10 秒钟。然后，在父进程中发送 **SIGINT** 信号给子进程，并等待子进程结束。当子进程收到 **SIGINT** 信号时，会调用 `my_handler()` 函数来处理信号。


# 7. kill_proc
## 7.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int kill_proc(pid_t pid, int sig, int priv);` |   它是 Linux 内核中的一个函数，用于向指定进程发送信号。|

**参数：**
- **pid ：** 要接收信号的进程的 PID
- **sig ：** 要发送的信号
- **priv ：** 表示是否对目标进程进行权限检查的标志。

## 7.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/sched.h>

static struct task_struct *my_thread;

static int my_thread_func(void *data)
{
    printk(KERN_INFO "my_thread_func started\n");

    /* send signal to current process */
    kill_proc(current->pid, SIGTERM, 1);

    printk(KERN_INFO "my_thread_func finished\n");
    return 0;
}

int my_init(void)
{
    printk(KERN_INFO "Initializing module...\n");

    my_thread = kthread_create(my_thread_func, NULL, "my_thread");
    if (IS_ERR(my_thread)) {
        printk(KERN_ERR "kthread_create failed\n");
        return -1;
    }

    wake_up_process(my_thread);

    return 0;
}

void my_exit(void)
{
    printk(KERN_INFO "Exiting module...\n");

    if (my_thread) {
        kthread_stop(my_thread);
    }
}

MODULE_LICENSE("GPL");
module_init(my_init);
module_exit(my_exit);

```
上述示例程序中，在 `my_init()` 函数中使用 `kthread_create()` 函数创建了一个名为 **my_thread** 的新内核线程，并将其启动。该线程的入口函数是 `my_thread_func()`，在这个函数中可以执行任何需要在内核空间中进行的操作。在这个函数中，使用 `kill_proc()` 函数向当前进程发送了 **SIGTERM** 信号，以请求终止当前进程。

> **注意：** 在使用 `kill_proc()` 函数时，必须确保目标进程存在且具有对应的权限，否则可能会导致系统出现不可预期的行为。另外，在发送信号之前，还需要先获得目标进程的 **PID**，这通常可以通过 **/proc** 文件系统中的相关接口或者其他方式来实现。

# 8. kmalloc
## 8.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void *kmalloc(size_t size, gfp_t flags);` |  它是 Linux 内核中的一个函数，用于在内核中分配指定大小的内存空间。 |

**参数：**
- **size：** 表示要分配的内存大小
- **flags：** 表示一组标志位，用于控制内存分配方式。

## 8.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/slab.h>

int my_init(void)
{
    void *ptr;

    printk(KERN_INFO "Initializing module...\n");

    ptr = kmalloc(1024, GFP_KERNEL);
    if (!ptr) {
        printk(KERN_ERR "kmalloc failed\n");
        return -1;
    }

    /* do something with ptr */

    kfree(ptr);

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

上述示例程序中，在 `my_init()` 函数中我们使用 `kmalloc()` 函数分配了一个大小为 **1024** 字节的内存空间，并且对于成功与否进行了检查。然后，在处理完该内存区域之后，使用 `kfree()` 函数释放了所占用的内存。

> **注意：** 在使用 `kmalloc()` 函数分配内存时，必须确保所分配的内存不会造成内核堆栈溢出或其他安全问题。另外，由于 `kmalloc()` 函数返回的内存地址可能不是连续的，因此在使用该函数分配大块内存时，需要特别注意内存对齐和分配方式等问题。

# 9. kmap，kmap_high
## 9.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void *kmap(struct page *page);`|它是 Linux 内核中的一个函数，用于将一个页映射到内核虚拟地址空间。|
|`void *kmap_high(struct page *page);` |  它是 Linux 内核中的一个函数，用于将高端内存映射到内核虚拟地址空间中。 |

**参数：**
- **page ：** 要映射的物理页面的指针。

## 9.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/highmem.h>

int my_init(void)
{
    struct page *page;
    void *ptr;

    printk(KERN_INFO "Initializing module...\n");

    /* allocate high memory page */
    page = alloc_pages(GFP_HIGHUSER, 0);
    if (!page) {
        printk(KERN_ERR "Failed to allocate page\n");
        return -ENOMEM;
    }

    /* map page to kernel virtual address space */
    ptr = kmap_high(page);
    if (!ptr) {
        printk(KERN_ERR "Failed to map page\n");
        __free_pages(page, 0);
        return -EFAULT;
    }

    /* do something with mapped memory */
    *(char *)ptr = 'A';
    printk(KERN_INFO "Value at mapped address: %c\n", *(char *)ptr);

    /* unmap and free page */
    kunmap_high(ptr);
    __free_pages(page, 0);

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

上述示例程序中，在 `my_init()` 函数中使用 `alloc_pages()` 函数分配了一块高端内存物理页面，并使用 `kmap_high()` 函数将其映射到内核虚拟地址空间中。在之后，对映射的内存进行操作并打印它的值。最后，使用 `kunmap_high()` 函数将映射解除并使用 `__free_pages()` 函数释放页面。

> **注意：** 在使用 kmap_high() 函数时，必须确保请求的页面大小不会超过系统可用的物理内存大小，并且可以正确地处理异常情况。另外，使用高端内存应格外小心，因为它们与物理内存管理和 **DMA** 操作相关，可能会影响系统的稳定性和安全性。

# 10. kmem_cache_alloc
## 10.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `void *kmem_cache_alloc(struct kmem_cache *cachep, int flags);`| 它是 Linux 内核中的一个函数，用于从指定的内存缓存中分配一个对象  |

**参数：**
- **cachep：** 指向所需缓存区的指针
- **flags：** 用于控制内存分配方式的标志。

## 10.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/slab.h>

struct my_struct {
    int a;
    char b;
};

static struct kmem_cache *my_cachep;

int my_init(void)
{
    struct my_struct *obj;

    printk(KERN_INFO "Initializing module...\n");

    my_cachep = kmem_cache_create("my_cache", sizeof(struct my_struct), 0, SLAB_HWCACHE_ALIGN, NULL);
    if (!my_cachep) {
        printk(KERN_ERR "kmem_cache_create failed\n");
        return -1;
    }

    obj = kmem_cache_alloc(my_cachep, GFP_KERNEL);
    if (!obj) {
        printk(KERN_ERR "kmem_cache_alloc failed\n");
        return -1;
    }

    obj->a = 123;
    obj->b = 'A';

    printk(KERN_INFO "obj->a = %d, obj->b = '%c'\n", obj->a, obj->b);

    return 0;
}

void my_exit(void)
{
    printk(KERN_INFO "Exiting module...\n");

    kmem_cache_free(my_cachep, obj);

    kmem_cache_destroy(my_cachep);
}

MODULE_LICENSE("GPL");
module_init(my_init);
module_exit(my_exit);
```

上述示例程序中，首先使用 `kmem_cache_create()` 函数创建了一个名为 **my_cache** 的内存缓存区，该缓存区能够容纳 `struct my_struct` 类型的对象。然后，在 `my_init()` 函数中使用 `kmem_cache_alloc()` 函数从缓存区中分配了一个 `struct my_struct` 类型的对象，并进行了初始化操作。最后，在 `my_exit()` 函数中使用 `kmem_cache_free()` 函数释放了该对象所占用的内存，然后销毁了整个缓存区。

> **注意：** `kmem_cache_alloc()` 函数和其他 **Linux** 内核函数在用户空间下无法直接使用，通常需要编写内核模块来调用这些函数。

# 11. kmem_cache_create
## 11.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `struct kmem_cache *kmem_cache_create(const char *name, size_t size, size_t align, unsigned long flags, void (*ctor)(void *));`|  它是 Linux 内核中的一个函数，用于创建一个内存缓存区，可以用于高效地分配和释放指定大小的对象。 |

**参数：**
- **name ：** 缓存区的名称
- **size ：** 要分配的对象的大小
- **align ：** 对齐方式
- **flags ：**  标志位
- **ctor ：**  构造函数指针。

## 11.2 演示示例
**参考 10.2 所示**

> **注意：** 在使用 `kmem_cache_create()` 函数时，必须确保请求的内存大小不会超过系统可用的物理内存大小，并且可以正确地处理内存分配失败等异常情况。另外，使用内存缓存区应格外小心，因为它们与内核数据结构和操作高度相关，并且可能会影响系统的稳定性和安全性。

# 12. kmem_cache_free
## 12.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void kmem_cache_free(struct kmem_cache *cachep, void *objp);`|它是 Linux 内核中的一个函数，用于将之前使用 kmem_cache_alloc() 函数分配的对象释放回内存缓存池，以便下次再次分配使用。|

**参数：**
- **cachep：** 指向之前使用的缓存区的指针
- **objp：** 要释放的对象的指针。

## 12.2 演示示例
**参考 10.2 所示**

> 注意： `kmem_cache_free()` 函数和其他 **Linux** 内核函数在用户空间下无法直接使用，通常需要编写内核模块来调用这些函数。


# 13. kmem_cache_destroy
## 13.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void kmem_cache_destroy(struct kmem_cache *cachep);`|它是 Linux 内核中的一个函数，用于销毁之前使用 kmem_cache_create() 函数创建的内存缓存区。|

**参数：**
- **cachep：** 指向要销毁的缓存区的指针

## 13.2 演示示例
**参考 10.2 所示**

> **注意：** 在销毁缓存区之前，必须确保所有从缓存区中分配的内存都已经被释放，否则可能会导致内存泄漏或其他问题。


# 14. kmem_cache_zalloc
## 14.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void *kmem_cache_zalloc(struct kmem_cache *cache, gfp_t flags);`|它是 Linux 内核中的一个函数，用于从指定内存缓存区中分配一块指定大小的内存，并将其清零。|

**参数：**
- **cache ：** 要分配内存的缓存区
- **flags ：** 用于控制内存分配行为的标志

## 14.2 演示示例

```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/slab.h>

static struct kmem_cache *my_cache;

struct my_struct {
    int value;
};

int my_init(void)
{
    struct my_struct *obj;

    printk(KERN_INFO "Initializing module...\n");

    /* create cache */
    my_cache = kmem_cache_create("my_cache", sizeof(struct my_struct), 0, 0, NULL);
    if (!my_cache) {
        printk(KERN_ERR "kmem_cache_create failed\n");
        return -ENOMEM;
    }

    /* allocate and initialize object */
    obj = kmem_cache_zalloc(my_cache, GFP_KERNEL);
    if (!obj) {
        printk(KERN_ERR "kmem_cache_zalloc failed\n");
        kmem_cache_destroy(my_cache);
        return -ENOMEM;
    }
    obj->value = 666;

    /* print object value */
    printk(KERN_INFO "Object value: %d\n", obj->value);

    /* free object and destroy cache */
    kmem_cache_free(my_cache, obj);
    kmem_cache_destroy(my_cache);

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
上述示例程序中，在 `my_init()` 函数中使用 `kmem_cache_create()` 函数创建了一个名为 **my_cache** 的内存缓存区，并使用 `kmem_cache_zalloc()` 函数从中分配了一块大小为 `sizeof(struct my_struct)` 的内存。在之后，将这块内存的值设置为 **666** 并打印它。最后，使用 `kmem_cache_free()` 函数释放内存并销毁缓存区。

> **注意：** 在使用 `kmem_cache_zalloc()` 函数时，必须确保请求的内存大小不会超过系统可用的物理内存大小，并且可以正确地处理内存分配失败等异常情况。另外，分配的内存应在不再需要时及时释放，以免造成内存泄漏等问题。

# 15. kmemdup
## 15.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void *kmemdup(const void *src, size_t len, gfp_t flags);` | 它是 Linux 内核中的一个函数，用于在内核空间中将一段指定大小的内存复制到另一段新的内存中，并返回这段新内存的指针。  |

**参数：**
- **src ：** 要复制的源内存地址
- **len ：** 要复制的内存字节数
- **flags ：** 用于控制内存分配行为的标志

## 15.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/slab.h>

int my_init(void)
{
    char *src = "Hello, world!";
    char *dst;

    printk(KERN_INFO "Initializing module...\n");

    /* allocate and copy memory */
    dst = kmemdup(src, strlen(src) + 1, GFP_KERNEL);
    if (!dst) {
        printk(KERN_ERR "kmemdup failed\n");
        return -ENOMEM;
    }

    /* print copied memory */
    printk(KERN_INFO "Copied string: %s\n", dst);

    /* free memory */
    kfree(dst);

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
上述示例程序中，在 `my_init()` 函数中使用 `kmemdup()` 函数将一段字符串 `Hello, world!` 复制到另一段新的内存中。在之后，打印了这段复制的内存，并释放了它。

> **注意：** 在使用 `kmemdup()` 函数时，必须确保请求的内存大小不会超过系统可用的物理内存大小，并且可以正确地处理内存分配失败等异常情况。另外，分配的内存应在不再需要时及时释放，以免造成内存泄漏等问题。

# 16. kprintf
## 16.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void kprintf(const char *format, ...);` |  用于嵌入式系统中输出调试信息 |

> **注意：** 该函数原型和使用方法与标准库中的 `printf()` 函数类似。不同的是，`kprintf()` 函数通常需要根据具体的嵌入式系统进行修改，以适应不同的输出方式。

## 16.2 演示示例
```c
#include <stdio.h>
#include <stdarg.h>
#include <fcntl.h>
#include <termios.h>

int uart_fd = -1; // 串口文件描述符

void kprintf(const char *format, ...)
{
    va_list ap;
    char buf[256];

    va_start(ap, format);
    vsnprintf(buf, sizeof(buf), format, ap);
    va_end(ap);

    if (uart_fd != -1) {
        write(uart_fd, buf, strlen(buf));
    }
}

int init_uart(const char *devname)
{
    struct termios tio;

    uart_fd = open(devname, O_RDWR | O_NOCTTY);
    if (uart_fd == -1) {
        perror("open failed");
        return -1;
    }

    memset(&tio, 0, sizeof(tio));
    cfmakeraw(&tio);
    cfsetspeed(&tio, B115200);
    tcsetattr(uart_fd, TCSANOW, &tio);

    return 0;
}

int main()
{
    int ret;

    ret = init_uart("/dev/ttyS0"); // 打开 ttyS0 串口
    if (ret == -1) {
        exit(EXIT_FAILURE);
    }

    kprintf("Hello, world!\n"); // 输出调试信息

    close(uart_fd); // 关闭串口文件描述符
    return 0;
}
```
上述示例程序中，首先通过 **init_uart()** 函数打开了 **ttyS0** 串口，并将其设置为 **RAW** 模式和波特率 **115200**。然后，在 `main()` 函数中调用了 `kprintf()` 函数来输出一条调试信息。该函数会将调试信息写入 **ttyS0** 串口中，并发送到外部设备（如 **PC**）上。

> **注意：**`kprintf()` 函数通常需要进行一定的修改以适应具体的嵌入式系统和调试工具，上面仅提供一个简单的示例，不能直接在所有系统中使用。

# 17. krealloc
## 17.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `void *krealloc(const void *ptr, size_t new_size, gfp_t flags);`|  它是 Linux 内核中的一个函数，用于动态调整已分配内存块的大小。 |

**参数：**
- **ptr ：** 指向原内存块的指针
- **new_size ：** 新的内存块大小
- **flags ：** 用于控制内存分配行为的标志。

## 17.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/slab.h>

int my_init(void)
{
    char *buf = kmalloc(16, GFP_KERNEL);

    printk(KERN_INFO "Initializing module...\n");

    if (!buf) {
        printk(KERN_ERR "kmalloc failed\n");
        return -ENOMEM;
    }

    /* resize buffer */
    buf = krealloc(buf, 32, GFP_KERNEL);
    if (!buf) {
        printk(KERN_ERR "krealloc failed\n");
        kfree(buf);
        return -ENOMEM;
    }

    /* print buffer size */
    printk(KERN_INFO "Buffer size: %lu\n", ksize(buf));

    /* free buffer */
    kfree(buf);

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
上述示例程序中，在 `my_init()` 函数中使用 `kmalloc()` 函数分配了一块大小为 **16** 字节的内存，并使用 `krealloc()` 函数将其调整为大小为 **32** 字节。在之后，打印了这块内存的大小并释放了它。

> **注意：** 在使用 `krealloc()` 函数时，必须确保操作合法且不会导致内存泄漏或其他问题。另外，应当避免过度依赖动态内存分配而导致系统性能下降或出现其他问题。


# 18. ksize
## 18.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `size_t ksize(const void *ptr);`| 它是 Linux 内核中的一个函数，用于获取已分配内存块的大小。  |

**参数：**
- **ptr ：** 指向已分配内存块的指针。

## 18.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/slab.h>

int my_init(void)
{
    char *buf = kmalloc(16, GFP_KERNEL);

    printk(KERN_INFO "Initializing module...\n");

    if (!buf) {
        printk(KERN_ERR "kmalloc failed\n");
        return -ENOMEM;
    }

    /* print buffer size */
    printk(KERN_INFO "Buffer size: %lu\n", ksize(buf));

    /* free buffer */
    kfree(buf);

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
上述示例程序中，在 `my_init()` 函数中使用 `kmalloc()` 函数分配了一块大小为 **16** 字节的内存，并使用 `ksize()` 函数获取了它的大小。在之后，打印了这块内存的大小并释放了它。

# 19. kstrdup
## 19.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `char *kstrdup(const char *s, gfp_t flags);`|  它是 Linux 内核中的一个函数，用于在内核空间中复制一个以 NULL 结尾的字符串，并返回这个新的字符串指针。 |

**参数：**
- **s ：** 要复制的源字符串
- **flags ：** 用于控制内存分配行为的标志

## 19.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/slab.h>

int my_init(void)
{
    const char *src = "Hello, world!";
    char *dst;

    printk(KERN_INFO "Initializing module...\n");

    /* duplicate string */
    dst = kstrdup(src, GFP_KERNEL);
    if (!dst) {
        printk(KERN_ERR "kstrdup failed\n");
        return -ENOMEM;
    }

    /* print duplicated string */
    printk(KERN_INFO "Duplicated string: %s\n", dst);

    /* free memory */
    kfree(dst);

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
上述示例程序中，在 `my_init()` 函数中使用 `kstrdup()` 函数将一个字符串 `Hello, world!` 复制到另一段新的内存中。在之后，打印了这段复制的字符串，并释放了它。

> **注意：** 在使用 `kstrdup()` 函数时，必须确保源字符串以 `NULL` 结尾，并且请求的内存大小不会超过系统可用的物理内存大小，并且可以正确地处理内存分配失败等异常情况。另外，分配的内存应在不再需要时及时释放，以免造成内存泄漏等问题。


# 20. kstrndup
## 20.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`char *kstrndup(const char *s, size_t len, gfp_t flags);` |  它是 Linux 内核中的一个函数，用于在内核空间中复制一个以 NULL 结尾的字符串的一部分，并返回这个新的字符串指针。 |

**参数：**
- **s ：** 要复制的源字符串
- **len ：** 要复制的字符串长度
- **flags ：** 用于控制内存分配行为的标志

## 20.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/slab.h>

int my_init(void)
{
    const char *src = "Hello, world!";
    char *dst;

    printk(KERN_INFO "Initializing module...\n");

    /* duplicate string */
    dst = kstrndup(src, 5, GFP_KERNEL);
    if (!dst) {
        printk(KERN_ERR "kstrndup failed\n");
        return -ENOMEM;
    }

    /* print duplicated string */
    printk(KERN_INFO "Duplicated string: %s\n", dst);

    /* free memory */
    kfree(dst);

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
上述示例程序中，在 `my_init()` 函数中使用 `kstrndup()` 函数将一个字符串 `Hello, world!` 的 **前五个字符** 复制到另一段新的内存中。在之后，打印了这段复制的字符串，并释放了它。

> **注意：** 在使用 `kstrndup()` 函数时，必须确保源字符串以 `NULL` 结尾，并且请求的内存大小不会超过系统可用的物理内存大小，并且可以正确地处理内存分配失败等异常情况。另外，分配的内存应在不再需要时及时释放，以免造成内存泄漏等问题。

# 21. kstat_irqs_cpu
## 21.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void kstat_irqs_cpu(int cpu, int *irqs, unsigned long *stime);` |  它是 Linux 内核中的一个函数，用于查询指定 CPU 的中断统计信息。 |

**参数：**
> **cpu ：** 要查询的 CPU 编号
> **irqs ：** 用于保存中断计数值的数组
> **stime ：** 用于保存中断处理时间戳的变量

## 21.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>

int my_init(void)
{
    int irqs[NR_IRQS];
    unsigned long stime;

    printk(KERN_INFO "Initializing module...\n");

    /* get IRQ statistics for CPU 0 */
    kstat_irqs_cpu(0, irqs, &stime);

    /* print IRQ statistics */
    printk(KERN_INFO "IRQ statistics for CPU 0:\n");
    for (int i = 0; i < NR_IRQS; i++) {
        if (irqs[i] > 0) {
            printk(KERN_INFO "IRQ %d: count=%d\n", i, irqs[i]);
        }
    }

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
上述示例程序中，在 `my_init()` 函数中使用 `kstat_irqs_cpu()` 函数查询了 **CPU 0** 的中断统计信息，并将结果打印到内核日志中。

> **注意：** 在使用 `kstat_irqs_cpu()` 函数时，必须确保传递给该函数的参数是正确的，并且具有足够的权限来访问相关的数据结构。另外，中断统计信息可能会随着时间的推移而发生变化，因此需要在使用之前及时更新数据。



# 22. kthread_create
## 22.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`struct task_struct *kthread_create(int (*threadfn)(void *data), void *data, const char *namefmt, ...);` | 它是 Linux 内核中的一个函数，用于创建一个内核线程。  |

**参数：**
> **threadfn ：** 指向线程处理函数的指针
> **data ：** 传递给线程处理函数的参数
> **namefmt ：** 用于命名线程的格式化字符串（类似于 **printf()** 函数的格式化字符串）。

## 22.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/kthread.h>

static struct task_struct *my_thread;

static int my_thread_func(void *data)
{
    printk(KERN_INFO "my_thread_func started\n");

    while (!kthread_should_stop()) {
        /* do something */
    }

    printk(KERN_INFO "my_thread_func finished\n");
    return 0;
}

int my_init(void)
{
    printk(KERN_INFO "Initializing module...\n");

    my_thread = kthread_create(my_thread_func, NULL, "my_thread");
    if (IS_ERR(my_thread)) {
        printk(KERN_ERR "kthread_create failed\n");
        return -1;
    }

    wake_up_process(my_thread);

    return 0;
}

void my_exit(void)
{
    printk(KERN_INFO "Exiting module...\n");

    if (my_thread) {
        kthread_stop(my_thread);
    }
}

MODULE_LICENSE("GPL");
module_init(my_init);
module_exit(my_exit);
```
上述示例程序中，在 `my_init()` 函数中使用 `kthread_create()` 函数创建了一个名为 **my_thread** 的新内核线程，并将其启动。该线程的入口函数是 `my_thread_func()`，在这个函数中可以执行任何需要在内核空间中进行的操作。在 `my_exit()` 函数中，使用 `kthread_stop()` 函数停止该线程的运行。

> **注意：** 在使用 `kthread_create()` 函数时，必须确保线程处理函数能够正确地完成自己的任务并且及时退出，否则可能会导致内核崩溃或其他问题。另外，使用内核线程时需要特别小心，因为它们与内核数据结构和操作高度相关，并且可能会影响系统的稳定性和安全性。

# 23. kthread_stop
## 23.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int kthread_stop(struct task_struct *k);` |  它是 Linux 内核中的一个函数，用于停止由 kthread_create() 函数创建的内核线程。 |

**参数：**
> **k ：** 指向要停止的内核线程的指针。

## 23.2 演示示例
**参考 22.2 中所示**

> **注意：** 在使用 `kthread_stop()` 函数时，必须确保线程处理函数能够正确地响应该函数并及时退出，并且不会造成内核资源泄漏等问题，否则可能会导致内核崩溃或其他问题。另外，在使用 `kthread_stop()` 函数之前，需要先调用 `kthread_should_stop()` 函数来检查线程是否应该停止。


# 24. kunmap_high
## 24.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `void kunmap_high(struct page *page);`|  它是 Linux 内核中的一个函数，用于取消一个高端内存映射。 |

**参数：**
- **page ：** 要取消映射的高端内存页

## 24.2 演示示例
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/highmem.h>

int my_init(void)
{
    struct page *page;
    char *kaddr;

    printk(KERN_INFO "Initializing module...\n");

    /* allocate high memory */
    page = alloc_highmem_pages(1);
    if (!page) {
        printk(KERN_ERR "alloc_highmem_pages failed\n");
        return -ENOMEM;
    }

    /* map high memory to kernel address space */
    kaddr = kmap(page);
    if (!kaddr) {
        printk(KERN_ERR "kmap failed\n");
        free_highmem_page(page);
        return -ENOMEM;
    }

    /* print mapped address */
    printk(KERN_INFO "Mapped address: %pK\n", kaddr);

    /* unmap high memory */
    kunmap_high(page);

    /* free high memory */
    free_highmem_page(page);

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
上述示例程序中，在 `my_init()` 函数中使用 `alloc_highmem_pages()` 函数分配了一块高端内存页，并使用 `kmap()` 函数将其映射到内核虚拟地址空间。在之后，打印了这个映射的内核地址，并使用 `kunmap_high()` 函数取消了这个映射。最后，释放了这块高端内存页。

> **注意：** 在使用高端内存时，必须确保操作合法且不会导致内存泄漏或其他问题。另外，在使用 `kmap()` 和 `kunmap_high()` 函数时，必须确保映射和取消映射是成对出现的，否则可能会导致内存泄漏或其他问题。


# 25. kzalloc
## 25.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void *kzalloc(size_t size, gfp_t flags);` | 它是 Linux 内核中的一个函数，用于分配指定大小的内存空间，并将其初始化为零。  |

**参数：**
- **size：** 表示要分配的内存大小
- **flags：** 表示一组标志位，用于控制内存分配方式。

## 25.2 演示示例
```c
c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/slab.h>

int my_init(void)
{
    void *ptr;

    printk(KERN_INFO "Initializing module...\n");

    ptr = kzalloc(1024, GFP_KERNEL);
    if (!ptr) {
        printk(KERN_ERR "kzalloc failed\n");
        return -1;
    }

    /* do something with ptr */

    kfree(ptr);

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

上述示例程序中，在 `my_init()` 函数中使用 `kzalloc()` 函数分配了一个大小为 **1024** 字节的内存空间，并将其初始化为零。然后，在处理完该内存区域之后，使用 `kfree()` 函数释放了所占用的内存。

> **注意：** 虽然可以使用 `kmalloc()` 函数分配内存然后手动初始化为零，但是使用 `kzalloc()` 函数可以更加高效和简单地完成这个操作。另外，由于 `kzalloc()` 函数返回的内存地址是经过清零的，因此在使用该函数分配内存时，无需显式调用 `memset()` 等函数进行初始化操作。


# 参考
1. [\[API Reference Document\]](https://www.apiref.com/c-zh/keep.htm)
2. [《Linux内核API完全参考手册》](https://baike.baidu.com/item/Linux%E5%86%85%E6%A0%B8API%E5%AE%8C%E5%85%A8%E5%8F%82%E8%80%83%E6%89%8B%E5%86%8C/7829665?fr=aladdin)
