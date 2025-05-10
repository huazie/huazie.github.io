---
title: C语言函数大全--m 开头的函数（下）
date: 2023-04-22 22:49:30
updated: 2025-05-10 17:23:23
categories:
  - 开发语言-C
tags:
  - C语言函数大全
  - m开头的函数
---

![](/images/cplus-logo.png)


# 总览
| 函数声明 |  函数功能  |
|:--|:--|
|`int mkdirat(int dirfd, const char *pathname, mode_t mode);` | 它是一个 `Linux` 系统下的系统调用函数，用于在指定目录下创建新的子目录  |
|`int mkfifo(const char *pathname, mode_t mode);` |  它是一个 `Linux` 系统下的系统调用函数，用于创建命名管道 |
|`int mkstemp(char *template);` | 用于在磁盘上创建一个唯一的临时文件并打开它以进行读写  |
|`int mkdir(const char *pathname, mode_t mode);` |它是一个 `Linux` 系统下的系统调用函数，用于创建新目录   |
|`int mkdir(const char *pathname);` | 它是在 `Windows` 系统下的系统调用函数，用于创建新目录  |
|`time_t mktime(struct tm *timeptr);` | 用于将表示时间的结构体（struct tm）转换为对应的 Unix 时间戳  |
|`int mlock(const void *addr, size_t len);` |  它是一个 `Linux` 系统下的系统调用函数，用于将指定内存区域锁定在物理内存中，防止其被交换到磁盘上  |
|`void *mmap(void *addr, size_t length, int prot, int flags, int fd, off_t offset);` |  它是一个 `Linux` 系统下的系统调用函数，可以将一个文件或者设备映射到内存中，并返回指向该内存区域的指针 |
|`double modf(double x, double *iptr);` | 用于将浮点数 `value` 拆分为其整数部分和小数部分（double）  |
|`float modff(float value, float *iptr);` | 用于将浮点数 `value` 拆分为其整数部分和小数部分（float）  |
|`long double modfl(long double value, long double *iptr);` | 用于将浮点数 `value` 拆分为其整数部分和小数部分（long double）  |
|`int mount(const char *source, const char *target, const char *filesystemtype, unsigned long mountflags, const void *data);` | 用于将文件系统挂载到指定的挂载点，并返回挂载点的文件描述符  |
|`int msync(void *addr, size_t length, int flags);` | 用于将指定内存区域的数据同步到文件中  |
|`int munmap(void *addr, size_t length);` | 用于取消内存映射区域，并释放与之相关的资源  |
|`int munlock(const void *addr, size_t len);` | 用于将之前使用mlock()函数锁定的内存区域解锁，使其可被操作系统交换出去或被回收  |


# 1. mkdirat
## 1.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int mkdirat(int dirfd, const char *pathname, mode_t mode);` | 它是一个 `Linux` 系统下的系统调用函数，用于在指定目录下创建新的子目录  |

**参数：**
- **dirfd ：** 要在其下创建新目录的父目录的文件描述符。如果值为 `AT_FDCWD`，则表示使用当前工作目录
- **pathname ：**  要创建的新目录的名称和路径
- **mode ：**  要创建的新目录的权限模式

**返回值：**
- 如果成功创建新目录时，则返回 `0`；
- 如果失败时，则返回 `-1`，并设置错误码（`errno`）。

## 1.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <fcntl.h>
#include <sys/stat.h>
#include <unistd.h>

int main() 
{
    int dirfd;
    if ((dirfd = open("/tmp", O_RDONLY)) == -1) 
    {
        printf("Error opening directory.\n");
        return 1;
    }

    if (mkdirat(dirfd, "testdir", S_IRWXU | S_IRGRP | S_IXGRP | S_IROTH | S_IXOTH) == -1) 
    {
        printf("Error creating directory.\n");
        return 1;
    }

    close(dirfd);
    return 0;
}
```

在上述的示例代码中，
- 首先，我们打开 `/tmp` 目录并获取其文件描述符 `dirfd`；
- 然后，我们调用 `mkdirat()` 函数，并将目录的文件描述符、要创建的新目录的名称和路径以及目录的权限模式作为参数传递给函数。如果函数调用成功，则新目录将在 `/tmp` 目录下创建。
- 最后，调用 `close()` 函数关闭文件。

**注意：** 
- 使用 `mkdirat()` 函数时，我们需要确保指定的父目录存在并具有适当的权限。
- 如果要使用相对路径创建新目录，需要确保当前工作目录正确设置。


# 2. mkfifo
## 2.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int mkfifo(const char *pathname, mode_t mode);` |  它是一个 `Linux` 系统下的系统调用函数，用于创建命名管道 |

**参数：**
- **pathname ：**  要创建的命名管道的名称和路径
- **mode ：**  命名管道的权限模式

**返回值：**
- 如果成功创建命名管道时，则返回 `0`；
- 如果失败时，则返回 `-1`，并设置错误码（`errno`）。

## 2.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <sys/stat.h>
#include <fcntl.h>

int main() 
{
    // 如果函数调用成功，则在 /tmp 目录下创建一个名为 myfifo 的命名管道
    if (mkfifo("/tmp/myfifo", S_IRUSR | S_IWUSR | S_IRGRP | S_IWGRP) == -1) 
    {
        printf("Error creating named pipe.\n");
        return 1;
    }

    return 0;
}
```

**注意：** 
- 使用 `mkfifo()` 函数时，我们需要确保指定的路径可被访问并且不存在同名的文件或目录。
- 如果要使用相对路径创建命名管道，需要确保当前工作目录正确设置。

# 3. mkstemp
## 3.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int mkstemp(char *template);` | 用于在磁盘上创建一个唯一的临时文件并打开它以进行读写  |

**参数：**
- **template ：**  指向包含 `6` 个连续 `'X'` 的字符串的指针，这些 `'X'` 将被替换为随机字符以创建唯一的文件名。例如，`"/tmp/tempfile-XXXXXX"` 将会被替换为类似 `"/tmp/tempfile-5ZqYU2"` 的唯一文件名。

## 3.2 演示示例

```c
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main() 
{
    char temp_file_template[] = "tmp/tempfile-XXXXXX";
    int fd;

    if ((fd = mkstemp(temp_file_template)) == -1) 
    {
        printf("Error creating temporary file.\n");
        return 1;
    }

    printf("Temporary file created: %s\n", temp_file_template);

    // 读写临时文件

    // 关闭临时文件
    close(fd);
    return 0;
}
```

如上演示了，在 `windows` 下 创建临时文件：
- 首先，我们定义了一个字符串 `temp_file_template`，它包含连续的 `6` 个 `'X'`；
- 接着，我们调用 `mkstemp()` 函数，并将指向 `temp_file_template` 的指针作为参数传递给函数。如果函数调用成功，则返回新创建文件的文件描述符，并将 `temp_file_template` 中的 `'X'` 替换为随机字符以形成唯一的文件名；
- 然后，调用 `printf()` 函数，输出该临时文件的名称；
- 再接着，可以使用 fd 来操作临时文件【这里示例省略了】；
- 再然后，临时文件使用完毕，调用 `close()` 函数关闭临时文件。
- 最后，返回 `0` 表示程序执行成功。

> **注意：** 在使用 `mkstemp()` 函数时，我们需要确保提供的模板字符串至少包含 `6` 个 `'X'`，并且文件命名方式不能与现有文件冲突。

## 3.3 运行结果
![](mkstemp.png)

# 4. mkdir
## 4.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int mkdir(const char *pathname, mode_t mode);` |它是一个 `Linux` 系统下的系统调用函数，用于创建新目录   |
|`int mkdir(const char *pathname);` | 它是在 `Windows` 系统下的系统调用函数，用于创建新目录  |

**参数：**
- **pathname：**  要创建的新目录的名称和路径
- **mode ：**  要创建的新目录的权限模式

**返回值：**
- 如果成功创建新目录时，则返回 `0`；
- 如果失败时，则返回 `-1`，并设置错误码（`errno`）。

## 4.2 演示示例
### 4.2.1 Windows 下示例
```c
#include <stdio.h>
#include <stdlib.h> 
#include <sys/stat.h>

int main() 
{
    if (mkdir("tmp/newdir") == -1) 
    {
        printf("Error creating new directory.\n");
        return 1;
    }

    return 0;
}
```

### 4.2.1 Linux 下示例

```c
#include <stdio.h>
#include <stdlib.h> 
#include <sys/stat.h>
#include <sys/types.h>

int main() 
{
    if (mkdir("/tmp/newdir", S_IRWXU | S_IRGRP | S_IXGRP | S_IROTH | S_IXOTH) == -1) 
    {
        printf("Error creating new directory.\n");
        return 1;
    }

    return 0;
}
```

## 4.3 运行结果
**Windows 下示例运行结果**
![](mkdir.png)

# 5. mktime
## 5.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`time_t mktime(struct tm *timeptr);` | 用于将表示时间的结构体（struct tm）转换为对应的 Unix 时间戳  |

**参数：**
- **timeptr ：**  指向 `struct tm` 结构体的指针，其中包含要转换为 `Unix` 时间戳的日期和时间信息

**返回值：**
- 如果转换成功，则返回对应于输入时间的 `Unix` 时间戳；
- 如果转换失败，则返回 `-1`。

## 5.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <time.h>

int main() 
{
    struct tm my_time = {0};  // 初始化为 0，避免随机值影响结果
    my_time.tm_year = 2023 - 1900;  // 年份应该减去 1900
    my_time.tm_mon = 4 - 1;  // 月份从 0 开始计数，应该减去 1
    my_time.tm_mday = 15;
    my_time.tm_hour = 10;
    my_time.tm_min = 30;
    my_time.tm_sec = 0;

    time_t timestamp = mktime(&my_time);
    if (timestamp == -1) 
    {
        printf("Error converting time to timestamp.\n");
        return 1;
    }

    printf("Unix timestamp: %lld\n", timestamp);

    return 0;
}
```

在上面的示例代码中，
- 首先，我们创建一个 `struct tm` 结构体 `my_time`，并将其初始化为 `0`；
- 然后，设置结构体的年、月、日、时、分、秒等信息。这里的年份应该减去 `1900`，月份应该从 `0` 开始计数减去 `1`；
- 接着，调用 `mktime()` 函数，并将指向 `my_time` 结构体的指针作为参数传递给函数。如果函数调用成功，则返回对应于输入时间的 `Unix` 时间戳。
- 最后，输出该 `Unix` 时间戳。

**注意：** 
- 在使用 `mktime()` 函数时，我们需要确保提供的 `struct tm` 结构体中的所有字段都已正确设置。
- 由于 `mktime()` 函数所使用的时区可能与系统默认的时区不同，所以在某些情况下，转换结果可能会有一定偏差。

## 5.3 运行结果
![](mktime.png)

# 6. mlock
## 6.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int mlock(const void *addr, size_t len);` |  它是一个 `Linux` 系统下的系统调用函数，用于将指定内存区域锁定在物理内存中，防止其被交换到磁盘上  |

**参数：**
- **addr ：**  要锁定的内存区域的起始地址
- **len ：**  要锁定的内存区域的长度（字节数）

**返回值：**
- 如果成功锁定内存区域时，则返回 `0`；
- 如果失败时，则返回 `-1`，并设置错误码（`errno`）。

## 6.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>   // memset 函数所需头文件
#include <sys/mman.h> // mlock 函数所需头文件

#define PAGE_SIZE 4096  // 页大小

int main() 
{
    char *buf;
    size_t size = PAGE_SIZE;

    // 分配一段内存并清零
    buf = (char *)malloc(size);
    memset(buf, 0, size);

    // 锁定分配的内存区域
    if (mlock(buf, size) == -1) 
    {
        printf("Error locking memory.\n");
        return 1;
    }

    // 使用分配的内存...
    strncpy(buf, "Hello world!", size);
    printf("%s\n", buf);

    // 解锁内存区域
    if (munlock(buf, size) == -1) 
    {
        printf("Error unlocking memory.\n");
        return 1;
    }

    free(buf);
    return 0;
}
```
在上述的示例代码中，
- 首先，我们使用 `malloc()` 函数分配了一个页大小的内存区域，并使用 `memset()` 函数将其清零；
- 然后，调用 `mlock()` 函数，并将指向分配内存区域起始地址的指针以及内存区域的长度作为参数传递给函数。如果函数调用成功，则锁定分配的内存区域，防止其被交换到磁盘上；
- 接着，就可以对该内存区域进行读写操作【示例代码简单演示了使用 `strncpy()` 函数向上述的分配内存区域中写入字符串 `"Hello world!"`，并通过 `printf()` 函数输出该字符串】；
- 再接着，调用 `munlock()` 函数解除内存区域的锁定；
- 最后，释放分配内容，并正常结束程序。

> **注意：** 在使用 `mlock()` 函数时，我们需要确保指定的内存区域已正确分配并且足够大，以避免锁定错误的内存区域。

# 7. mmap
## 7.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void *mmap(void *addr, size_t length, int prot, int flags, int fd, off_t offset);` |  它是一个 `Linux` 系统下的系统调用函数，可以将一个文件或者设备映射到内存中，并返回指向该内存区域的指针 |

**参数：**
- **addr ：**  映射区域开始地址，通常为 `NULL`，由内核选定
- **length ：**  ：映射区域的长度（字节数）
- **prot ：**  映射区域的保护方式。其可能的取值如下【按位或（|）组合】：
    - `PROT_NONE`: 区域不能被访问。
    - `PROT_READ`: 区域可被读取。
    - `PROT_WRITE`: 区域可被写入。
    - `PROT_EXEC`: 区域可被执行。
- **flags ：**  制定映射区域的类型和其他标志。其可能的取值如下【按位或（|）组合】：
    - `MAP_SHARED`: 允许多个进程共享该映射区域，对映射区域所做的修改将反映到所有共享该区域的进程中。
    - `MAP_PRIVATE`: 该映射区域只允许当前进程进行访问，对映射区域所做的修改不会反映到其他进程中。
    - `MAP_FIXED`: 强制将映射区域放置在指定的地址处（如果该地址已经被占用，则会导致错误）。
    - `MAP_ANONYMOUS`: 创建一个匿名映射区域，不与任何文件关联。
    - `MAP_FILE`: 将映射区域与文件关联，需要指定文件描述符和偏移量。
    - `MAP_LOCKED`: 指示内核在物理存储器中锁定映射区域的页面，以确保在访问该区域时不会发生缺页中断。
- **fd ：**  要映射的文件描述符
- **offset ：**  文件映射的起始偏移量

## 7.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <sys/mman.h>
#include <fcntl.h>
#include <unistd.h>

int main() 
{
    int fd;
    char *ptr;

    // 打开文件
    if ((fd = open("example.txt", O_RDWR)) == -1) 
    {
        printf("Error opening file.\n");
        return 1;
    }

    // 将文件映射到内存中
    ptr = mmap(NULL, 4096, PROT_READ | PROT_WRITE, MAP_SHARED, fd, 0);
    if (ptr == MAP_FAILED) 
    {
        printf("Error mapping file.\n");
        return 1;
    }

    // 使用映射区域进行读写操作
    strncpy(ptr, "Hello world!", 13);
    printf("%s\n", ptr);

    // 取消映射并关闭文件
    if (munmap(ptr, 4096) == -1) 
    {
        printf("Error unmapping file.\n");
        return 1;
    }
    close(fd);

    return 0;
}
```
上述的示例代码，演示了如何使用 `mmap()` 函数将一个文件映射到内存中，并使用指针 `ptr` 访问这个映射区域 :
- 首先，我们调用 `open()` 函数打开文件 `"example.txt"`，并检查是否成功打开。
- 接着，调用 `mmap()` 函数将文件的前 `4096` 字节映射到内存中，同时指定保护方式为可读写（`PROT_READ | PROT_WRITE`）以及共享属性（`MAP_SHARED`）。
- 然后，如果映射成功，则使用 `strncpy()` 函数向映射区域中写入字符串 `"Hello world!"`，并通过 `printf()` 函数输出该字符串。
- 最后，调用 `munmap()` 函数取消映射，并关闭文件。

# 8. modf，modff，modfl
## 8.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double modf(double x, double *iptr);` | 用于将浮点数 `value` 拆分为其整数部分和小数部分（double）  |
|`float modff(float value, float *iptr);` | 用于将浮点数 `value` 拆分为其整数部分和小数部分（float）  |
|`long double modfl(long double value, long double *iptr);` | 用于将浮点数 `value` 拆分为其整数部分和小数部分（long double）  |


**参数：**
- **value ：**  待处理的浮点数
- **iptr ：**  用于返回 `value` 的整数部分

## 8.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main() 
{
    double x = 3.141592653589793;
    double ipart;

    double fpart = modf(x, &ipart);

    printf("x = %f\n", x);
    printf("整数部分 = %.0f\n", ipart);
    printf("小数部分 = %f\n", fpart);

    float y = 2.718281828459045;
    float ipart_f;

    float fpart_f = modff(y, &ipart_f);

    printf("y = %f\n", y);
    printf("整数部分 = %.0f\n", ipart_f);
    printf("小数部分 = %f\n", fpart_f);

    long double z = 1.414213562373095;
    long double ipart_l;

    long double fpart_l = modfl(z, &ipart_l);

    printf("z = %Lf\n", z);
    printf("整数部分 = %.0Lf\n", ipart_l);
    printf("小数部分 = %Lf\n", fpart_l);

    return 0;
}
```

## 8.3 运行结果
![](modf.png)

# 9. mount
## 9.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int mount(const char *source, const char *target, const char *filesystemtype, unsigned long mountflags, const void *data);` | 用于将文件系统挂载到指定的挂载点，并返回挂载点的文件描述符  |

**参数：**
- **source ：**  要挂载的文件系统。可以是一个设备文件名、一个目录路径或者 `NULL`（表示根据 `filesystemtype` 参数自动选择默认源）
- **target ：**  要挂载到的目标点，即挂载点
- **filesystemtype ：**  要挂载的文件系统类型
- **mountflags ：**  挂载选项
- **data ：**  任意与文件系统相关的数据

## 9.2 演示示例
```c
#include <stdio.h>
#include <sys/mount.h>

int main() 
{
    int ret;
    const char* source = "/dev/sda1";
    const char* target = "/mnt/usbdrive";
    const char* type = "ext4";

    ret = mount(source, target, type, 0, NULL);
    if (ret == -1) 
    {
        perror("mount");
        return 1;
    }

    printf("File system mounted successfully!\n");

    return 0;
}
```
在上面的示例代码中，我们调用 `mount()` 函数将 `/dev/sda1` 设备文件上的 `ext4` 文件系统挂载到 `/mnt/usbdrive` 目录下，并检查挂载操作是否成功。如果 `mount()` 函数返回值为 `-1`，则表示挂载失败，通过 `perror()` 函数输出错误信息并返回 `1`，程序异常结束；否则，打印 `“File system mounted successfully!”` 表示挂载成功，并返回 `0`，程序正常结束。

# 10. msync
## 10.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int msync(void *addr, size_t length, int flags);` | 用于将指定内存区域的数据同步到文件中  |

**参数：**
- **addr ：**  要同步的内存起始地址
- **length ：**  要同步的内存区域长度
- **flags ：**  同步选项，可以是以下值之一：
    - **MS_ASYNC：**  进行异步写操作（默认选项）
    - **MS_SYNC：**  进行同步写操作
    - **MS_INVALIDATE：**  使 `cache` 内容无效

## 10.2 演示示例
```c
#include <stdio.h>
#include <sys/mman.h>
#include <fcntl.h>
#include <unistd.h>
#include <string.h>

int main() 
{
    int fd = open("example.txt", O_RDWR);
    if (fd == -1) 
    {
        perror("open");
        return 1;
    }

    // 将文件映射到内存中
    char* ptr = mmap(NULL, 4096, PROT_READ | PROT_WRITE, MAP_SHARED, fd, 0);
    if (ptr == MAP_FAILED) 
    {
        perror("mmap");
        close(fd);
        return 1;
    }

    // 写入新的数据
    strncpy(ptr, "Hello world!", 13);

    // 同步数据到文件中
    int ret = msync(ptr, 4096, MS_SYNC);
    if (ret == -1) 
    {
        perror("msync");
        munmap(ptr, 4096);
        close(fd);
        return 1;
    }

    // 解除映射并关闭文件
    ret = munmap(ptr, 4096);
    if (ret == -1) 
    {
        perror("munmap");
        close(fd);
        return 1;
    }
    close(fd);

    return 0;
}
```
在上面的示例代码中，
- 首先，我们通过 `open()` 函数打开文件 `"example.txt"`；
- 然后，使用 `mmap()` 函数将文件的前 `4096` 字节映射到内存中；
- 接着，调用 `strncpy()` 函数向映射区域中写入新的数据，并通过 `msync()` 函数将修改后的数据同步回磁盘文件；
- 最后，调用 `munmap()` 函数解除映射并关闭文件。


# 11. munmap
## 11.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int munmap(void *addr, size_t length);` | 用于取消内存映射区域，并释放与之相关的资源  |

**参数：**
- **addr ：**  要取消映射的内存起始地址
- **length ：**  要取消映射的内存区域长度

在调用 `munmap()` 函数后，操作系统将取消指定的内存映射，并回收相应的资源，包括虚拟地址空间和关联的物理内存页（如果存在）。此外，取消映射还可能导致未同步到磁盘文件中的修改数据丢失。

> **注意：** 必须在对映射区域进行任何修改或者访问之前，先使用 `mmap()` 函数将文件映射到内存中；并在完成所有操作之后，再使用 `munmap()` 函数解除映射。否则，可能会引发各种错误或者异常情况。

## 11.2 演示示例
**参见 7.2 的 演示示例，这里不再赘述**


# 12. munlock
## 12.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int munlock(const void *addr, size_t len);` | 用于将之前使用mlock()函数锁定的内存区域解锁，使其可被操作系统交换出去或被回收  |

**参数：**
- **addr ：**  待解锁的内存区域的起始地址
- **len ：**  待解锁的内存区域的长度（以字节为单位）

**返回值：**
- 如果成功解锁时，则返回 `0`；
- 如果失败时，则返回 `-1`，并设置错误码（`errno`）。

> **注意：** 只有拥有相应权限的进程才能解锁该内存区域。

## 12.2 演示示例
**参见 6.2 的 演示示例，这里不再赘述**

