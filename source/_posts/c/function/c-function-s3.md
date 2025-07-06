---
title: C语言函数大全--s 开头的函数（3）
date: 2023-05-03 22:01:05
updated: 2025-07-06 12:51:11
categories:
  - 开发语言-C
tags:
  - C语言函数大全
  - s 开头的函数
---

![](/images/cplus-logo.png)

# 总览
| 函数声明 |  函数功能  |
|:--|:--|
|`unsigned int sleep(unsigned int seconds);` | 它是 C 语言标准库中的函数，用于使当前进程挂起一定的时间。在挂起期间，操作系统会将该进程从调度队列中移除，直到指定的时间过去为止。  |
|`void Sleep(DWORD milliseconds);` | 它是 `Windows API` 中的一部分，与 sleep 函数类似，它可以使当前线程挂起一段时间。   |
|`int sopen(const char* filename, int access, int sharemode, int shflag, ...);` | 它是 `Microsoft Visual C++` 中的一个函数，用于打开文件并返回文件句柄。与标准库中的 `fopen` 函数不同，`sopen` 函数支持以二进制方式打开文件，并且可以指定文件读写方式、共享模式和文件访问权限等参数。  |
|`void sound(int frequency);` | 用于发出声音，`sound` 函数会持续发出声音，直到调用 `nosound` 函数停止  |
|`int spawnl(int mode, const char *cmdname, const char *arg0, ..., NULL);` | 它是在 Windows 平台上使用的函数，用于启动另一个程序，并等待该程序运行结束后再继续执行本程序  |
|`int spawnle(int mode, const char *cmdname, const char *arg0, ..., const char *envp[]);` | 它是在 `Windows` 平台上使用的函数，可以启动另一个程序，并通过指定的环境变量传递参数  |
|`int sprintf(char *str, const char *format, ...);` | 用于将格式化的字符串输出到指定的缓冲区中  |
|`int snprintf(char *str, size_t size, const char *format, ...);` |  用于将格式化的字符串输出到指定的缓冲区中，类似于 `sprintf` 函数，但它可以限制输出字符串的长度，避免缓冲区溢出 |
|`double sqrt(double x);` |  计算 x 的平方根（double） |
|`float sqrtf(float x);` |  计算 x 的平方根 （float）|
|`long double sqrtl(long double x);` |  计算 x 的平方根 （long double）|
|`void srand(unsigned int seed);` | 用于初始化伪随机数生成器  |
|`int sscanf(const char *str, const char *format, ...);` | 用于从字符串中读取数据并进行格式化转换  |
|`int stat(const char *path, struct stat *buf);` | 用于获取文件或目录的属性信息，例如文件大小、创建时间、修改时间等。这些属性信息都被保存在一个名为 `struct stat` 的结构体中。  |
|`int stime(const time_t *t);` |  它是是 Unix/Linux 系统中的一个系统调用函数，用于设置系统时间 |
|`char *stpcpy(char *dest, const char *src);` | 用于复制一个字符串到另一个字符串缓冲区，并返回目标字符串的结尾指针  |
|`char* strcat(char* dest, const char* src);` |  用于将一个字符串拼接到另一个字符串的末尾 |
|`char* strchr(const char* str, int c);` | 用于查找字符串中第一次出现指定字符的位置，并返回该位置的指针  |
|`int strcmp(const char* str1, const char* str2);` |  用于比较两个字符串是否相等 |
|`char* strcpy(char* dest, const char* src);` | 用于将一个字符串复制到另一个字符串缓冲区中  |
|`size_t strcspn(const char* str, const char* charset);` | 用于查找字符串中第一次出现指定字符集合中任何字符的位置，并返回该位置的索引  |

# 1. sleep
## 1.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`unsigned int sleep(unsigned int seconds);` | 它是 C 语言标准库中的函数，用于使当前进程挂起一定的时间。在挂起期间，操作系统会将该进程从调度队列中移除，直到指定的时间过去为止。  |
|`void Sleep(DWORD milliseconds);` | 它是 `Windows API` 中的一部分，与 sleep 函数类似，它可以使当前线程挂起一段时间。   |

**sleep 函数参数：**
- **seconds ：** 要挂起的时间，单位为秒

**Sleep 函数参数：**
- **milliseconds：** 要挂起的时间，单位为毫秒

## 1.2 演示示例

```c
#include <stdio.h>
#include <unistd.h>

int main()
{
    printf("Start sleeping...\n");
    sleep(5);
    printf("Wake up!\n");

    return 0;
}
```
在使用 `sleep()` 函数时，将会使当前线程或者进程暂停指定的时间，以便给其他进程或线程执行机会，同时也可以用来控制程序的运行速度。

虽然 `sleep()` 函数很简单，但是需要注意以下几点：

1. `sleep()` 的精度并不高，它所挂起的时间可能会略微超过要求的时间。
2. `sleep()` 函数是阻塞式的，即在调用 `sleep()` 函数期间，该进程不能进行任何其他操作，包括响应信号等。
3. 在使用 `sleep()` 函数期间，如果发生信号，那么 `sleep()` 函数将被中断，该进程将继续运行。

```c
#include <stdio.h>
#include <Windows.h>

int main()
{
    printf("Start sleeping...\n");
    Sleep(5000); // 暂停 5 秒钟
    printf("Wake up!\n");

    return 0;
}
```
由于 `Sleep()` 函数是阻塞式的，因此该函数调用期间，当前线程将被阻塞。在函数调用结束后，该线程将恢复运行。

在 `Windows` 系统下使用 `Sleep()` 函数时，需要注意以下几点：
1. `Sleep()` 函数以毫秒为单位指定时间，精度比 `sleep()` 函数更高。
2. 在调用 `Sleep()` 函数期间，当前线程将被阻塞，不能进行任何其他操作，包括响应信号等。
3. 在使用 `Sleep()` 函数期间，如果发生信号，那么 `Sleep()` 函数将被中断，该线程将继续运行。

## 1.3 运行结果
![](sleep.png)

![](sleep1.png)

# 2. sopen
## 2.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int sopen(const char* filename, int access, int sharemode, int shflag, ...);` | 它是 `Microsoft Visual C++` 中的一个函数，用于打开文件并返回文件句柄。与标准库中的 `fopen` 函数不同，`sopen` 函数支持以二进制方式打开文件，并且可以指定文件读写方式、共享模式和文件访问权限等参数。  |

**参数：**
- **filename ：** 要打开的文件名
- **access ：** 指定文件的访问方式，可以是以下值之一：
    - `_O_RDONLY`：只读方式打开文件
    - `_O_WRONLY`：只写方式打开文件
    - `_O_RDWR`：读写方式打开文件
    - `_O_APPEND`：在文件末尾追加数据
    - `_O_CREAT`：如果文件不存在，则创建文件
    - `_O_TRUNC`：如果文件已存在，清空文件内容
    - `_O_EXCL`：与 `_O_CREAT` 配合使用，如果文件已经存在则打开失败
- **sharemode ：** 指定文件共享模式，可以是以下值之一：
    - `_SH_DENYRW`：独占方式打开文件，其他进程不能读取或写入该文件
    - `_SH_DENYWR`：共享读取方式打开文件，其他进程不能写入该文件
    - `_SH_DENYRD`：共享写入方式打开文件，其他进程不能读取该文件
    - `_SH_DENYNO`：共享方式打开文件，其他进程可以读取和写入该文件
- **shflag ：** 指定文件属性标志，可以是以下值之一：
    - `_S_IWRITE`：文件可写
    - `_S_IREAD`：文件可读
- **... ：** 可选参数。如果指定了 `_O_CREAT` 参数，则需要指定文件的访问权限

## 2.2 演示示例
```c
#include <stdio.h>
#include <fcntl.h>
#include <io.h>
#include <sys/stat.h>

#define _SH_DENYNO 0x40

int main()
{
    int handle;
    char buffer[1024];

    // 打开文件并读取数据
    handle = sopen("output.bin", _O_RDONLY, _SH_DENYNO, _S_IREAD);
    if (handle == -1) 
    {
        printf("Failed to open file!\n");
        return 1;
    }
    read(handle, buffer, sizeof(buffer));
    printf("File content: %s\n", buffer);

    // 关闭文件句柄
    close(handle);

    return 0;
}
```

指定文件共享模式，如果没有对应的宏常量，则可以定义如下：

```c
#define _SH_DENYRW 0x10
#define _SH_DENYWR 0x20
#define _SH_DENYRD 0x30
#define _SH_DENYNO 0x40
```

## 2.3 运行结果
![](sopen.png)
![](sopen1.png)


# 3. sound
## 3.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void sound(int frequency);` | 用于发出声音  |
|`void nosound(void);` | `sound` 函数会持续发出声音，直到调用 `nosound` 函数停止  |

**参数：**
- **frequency ：** 要发出的声音的频率，单位为赫兹（Hz）

## 3.2 演示示例

```c
#include <stdio.h>
#include <windows.h>

int main() 
{
    printf("Playing sound...\n");
    sound(1000); // 发出 1000 Hz 音调
    Sleep(5000); // 等待 5 秒钟
    nosound();   // 停止发声
    printf("Sound stopped.\n");
    return 0;
}
```
Windows 下如果上述出现 `error: 'sound' was not declared in this scope`，可以使用如下：
```c
#include <stdio.h>
#include <windows.h>

int main() 
{
    printf("Playing sound...\n");
    Beep(1000, 5000); // 发出 1000 Hz 音调，持续 5 秒钟
    printf("Sound stopped.\n");
    return 0;
}
```
**注意：** 在 `Windows` 平台上建议使用 `Beep()` 函数代替 `sound()` 函数，因为 `Beep()` 函数不需要特殊的硬件支持，并且可移植性更好。

# 4. spawnl
## 4.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int spawnl(int mode, const char *cmdname, const char *arg0, ..., NULL);` | 它是在 Windows 平台上使用的函数，用于启动另一个程序，并等待该程序运行结束后再继续执行本程序  |

**参数：**
- **mode ：** 执行模式，可以为 `P_WAIT` 或 `P_NOWAIT`
- **cmdname ：** 要执行的程序名称
- **arg0：** 要传递给程序的命令行参数，以 `NULL` 结尾


## 4.2 演示示例
### 4.2.1 test.c
```c
#include <stdio.h>

int main()
{
    printf("Hello World\n");
    return 0;
}
```

### 4.2.2 spawnl 演示

```c
#include <process.h>
#include <stdio.h>

int main(void)
{
    int result;

    // 要执行的程序名和参数列表
    const char* cmdname = "test.exe";
    const char* arg0 = NULL;

    // 执行程序
    result = spawnl(P_WAIT, cmdname, cmdname, arg0, NULL);
    if (result == -1)
    {
        perror("Error from spawnl");
        exit(1);
    }
    return 0;
}

```

如果在使用 `spawnl()` 函数时遇到了 `"Error from spawnl: Invalid argument"` 错误，有可能是由于参数传递不正确或要执行的程序不存在等原因导致的。

以下是一些可能导致该错误的情况：
1. 要执行的程序不存在或路径不正确。
2. `cmdname` 参数包含非法字符或格式不正确。
3. 参数列表没有以 `NULL` 结尾。
4. 要执行的程序需要管理员权限，但当前用户没有足够的权限。

## 4.3 运行结果
![](spawnl.png)


# 5. spawnle
## 5.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int spawnle(int mode, const char *cmdname, const char *arg0, ..., const char *envp[]);` | 它是在 `Windows` 平台上使用的函数，可以启动另一个程序，并通过指定的环境变量传递参数  |

**参数：**
- **mode ：** 执行模式，可以为 `P_WAIT` 或 `P_NOWAIT`
- **cmdname ：** 要执行的程序名称
- **arg0 ：** 要传递给程序的命令行参数，以 `NULL` 结尾
- **envp ：** 要传递给程序的环境变量

## 5.2 演示示例
### 5.2.1 SubEnvTest.c
```c
#include <process.h>
#include <processenv.h>
#include <stdio.h>

int main(int argc, char *argv[], char **envp)
{
    printf("SubEnvTest Command line arguments:\n");
    for (int i=0; i<argc; i++)
        printf("[%d] : %s\n", i, argv[i]);
    printf("exec %s, Hello, %s\n", argv[0], argv[1]);

    for (int i = 0; envp[i] != NULL; i++)
    {
        printf("%s\n", envp[i]);
    }
    return 0;
} 
```

### 5.2.2 spawnle 演示

```c
#include <stdio.h>
#include <stdlib.h>
#include <process.h>

int main(int argc, char *argv[], char *envp[]) 
{
    int result;

    // 启动 SubEnvTest.exe，并传递当前环境变量
    result = spawnle(P_WAIT, "SubEnvTest.exe", "SubEnvTest.exe", NULL, envp);
    if (result == -1) 
    {
        printf("Error: %d\n", errno);
        return 1;
    }

    return 0;
}
```

## 5.3 运行结果
![](spawnle.png)

# 6. sprintf 
## 6.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int sprintf(char *str, const char *format, ...);` | 用于将格式化的字符串输出到指定的缓冲区中  |

**参数：**
- **str ：** 指向字符数组（缓冲区）的指针，用于存储生成的格式化字符串
- **format ：** 用于指定要生成的格式化文本
- **... ：** 用于填充格式化字符串中的占位符

## 6.2 演示示例
```c
#include <stdio.h>

int main() 
{
    char buffer[128];
    int value = 123;

    // 将整数转换为字符串并保存到 buffer 中
    sprintf(buffer, "The value is %d.", value);
    printf("%s\n", buffer);

    return 0;
}
```
## 6.3 运行结果
![](sprintf.png)

# 7. snprintf
## 7.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int snprintf(char *str, size_t size, const char *format, ...);` |  用于将格式化的字符串输出到指定的缓冲区中，类似于 `sprintf` 函数，但它可以限制输出字符串的长度，避免缓冲区溢出 |

**参数：**
- **str ：** 指向字符数组（缓冲区）的指针，用于存储生成的格式化字符串
- **size ：** 指定缓冲区可写入的最大字节数
- **format ：** 用于指定要生成的格式化文本
- **... ：** 用于填充格式化字符串中的占位符

## 7.2 演示示例
```c
#include <stdio.h>

int main() 
{
    char buffer[5];
    int value = 123456;

    // 将整数转换为字符串并保存到 buffer 中，最大可写入长度为 sizeof(buffer) - 1
    snprintf(buffer, sizeof(buffer), "%d", value);
    printf("%s\n", buffer); // 输出 "1234"

    return 0;
}
```
## 7.3 运行结果
![](snprintf.png)

# 8. sqrt，sqrtf，sqrtl
## 8.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double sqrt(double x);` |  计算 x 的平方根（double） |
|`float sqrtf(float x);` |  计算 x 的平方根 （float）|
|`long double sqrtl(long double x);` |  计算 x 的平方根 （long double）|

**参数：**
- **x ：** 要计算平方根的数

## 8.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main() 
{
    double x = 2.0;
    float y = 3.0f;
    long double z = 4.0L;

    printf("sqrt(%.1f) = %.2f\n", x, sqrt(x));
    printf("sqrtf(%.1f) = %.2f\n", y, sqrtf(y));
    printf("sqrtl(%.1Lf) = %.2Lf\n", z, sqrtl(z));

    return 0;
}
```
## 8.3 运行结果
![](sqrt.png)

# 9. srand
## 9.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void srand(unsigned int seed);` | 用于初始化伪随机数生成器  |

**参数：**
- **seed  ：** 用于设置伪随机数生成器的种子值。不同的种子值会产生不同的随机数序列。

## 9.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <time.h>

int main() 
{
    int i;

    // 使用当前时间作为种子值
    srand(time(NULL));

    // 生成 10 个随机数并输出到控制台
    for (i = 0; i < 10; i++) 
    {
        printf("%d ", rand());
    }
    printf("\n");

    return 0;
}
```

**注意：** 如果不设置种子值，则每次程序运行时都会得到相同的随机数序列。因此，我们在实际开发中，常常使用时间戳或其他随机值来作为种子值，以确保生成的随机数具有更好的随机性。

## 9.3 运行结果
![](srand.png)

# 10. sscanf 
## 10.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int sscanf(const char *str, const char *format, ...);` | 用于从字符串中读取数据并进行格式化转换  |

**参数：**
- **str ：** 要读取数据的字符串
- **format ：** 格式字符串，用于指定要读取的数据的类型、数量和顺序
- **... ：** 指向要写入的变量的指针

## 10.2 演示示例
```c
#include <stdio.h>

int main() 
{
    char str[] = "hello world 123";
    char buf[16];
    int num;

    // 从字符串中读取一个字符串和一个整数
    if (sscanf(str, "%s %*s %d", buf, &num) == 2) 
    {
        printf("String: %s\n", buf);
        printf("Number: %d\n", num);
    }

    return 0;
}
```
在上述的示例代码中，我们使用 `sscanf()` 函数从字符串 `"hello world 123"` 中读取一个字符串和一个整数，并输出到控制台。

**注意：** 在格式字符串中，`%s` 表示读取一个字符串，`%d` 表示读取一个整数。另外，`%*s` 表示读取并忽略一个字符串。

## 10.3 运行结果
![](sscanf.png)

# 11. stat 
## 11.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int stat(const char *path, struct stat *buf);` | 用于获取文件或目录的属性信息，例如文件大小、创建时间、修改时间等。这些属性信息都被保存在一个名为 `struct stat` 的结构体中。  |

**参数：**
- **path ：** 要获取属性信息的文件或目录路径
- **buf ：** 指向 `struct stat` 结构体的指针，用于存储获取到的属性信息

**返回值：**
- 如果执行成功，则返回 `0`；
- 否则返回 `-1` 并设置相应的错误码（存储在 `errno` 变量中）

## 11.2 演示示例
```c
#include <stdio.h>
#include <sys/stat.h>

int main() 
{
    struct stat file_stat;

    // 获取文件属性信息
    if (stat("test.txt", &file_stat) == 0) {
        printf("File size: %ld bytes\n", file_stat.st_size);
        printf("Creation time: %ld\n", file_stat.st_ctime);
        printf("Modification time: %ld\n", file_stat.st_mtime);
    }

    return 0;
}
```

## 11.3 运行结果
![](stat.png)

![](stat1.png)

# 12. stime 
## 12.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int stime(const time_t *t);` |  它是是 Unix/Linux 系统中的一个系统调用函数，用于设置系统时间 |

**参数：**
- **t：** 指向一个 time_t 类型变量的指针，表示要设置的系统时间

**返回值：**
- 如果执行成功，则返回 `0`；
- 否则返回 `-1` 并设置相应的错误码（存储在 `errno` 变量中）

## 12.2 演示示例
```c
#include <stdio.h>
#include <time.h>
#include <unistd.h>

int main() 
{
    time_t t;

    // 获取当前时间
    time(&t);

    // 输出当前时间
    printf("Current time: %s", ctime(&t));

    // 设置系统时间为 2022 年 1 月 1 日
    struct tm new_time;
    new_time.tm_sec = 0;
    new_time.tm_min = 0;
    new_time.tm_hour = 0;
    new_time.tm_mday = 1;
    new_time.tm_mon = 0;
    new_time.tm_year = 122;
    new_time.tm_wday = 6;
    new_time.tm_yday = 0;
    new_time.tm_isdst = -1;

    t = mktime(&new_time);
    stime(&t);

    // 短暂等待，确保时间设置完成
    sleep(1);

    // 再次输出当前时间
    time(&t);
    printf("New time: %s", ctime(&t));

    return 0;
}
```
在如上的示例代码中，
- 我们首先使用 `time()` 函数获取当前时间，并输出到控制台；
- 然后，我们设置系统时间为 `2022 年 1 月 1 日`，并等待一秒钟以确保时间设置完成。
- 最后，我们再次输出当前时间，以验证时间设置是否成功。

**注意：** `stime()` 函数只能在 `Linux/Unix` 系统上使用，并且需要 `root` 权限才能调用。另外，在修改系统时间时应谨慎行事，以避免对系统和应用程序造成不可预料的影响。

# 13. stpcpy
## 13.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`char *stpcpy(char *dest, const char *src);` | 用于复制一个字符串到另一个字符串缓冲区，并返回目标字符串的结尾指针  |

**参数：**
- **dest ：** 目标字符串的缓冲区，必须具有足够的空间来存储源字符串
- **src ：** 要复制的源字符串。

**返回值：** 一个指向目标字符串结尾的指针

## 13.2 演示示例
```c
#include <stdio.h>
#include <string.h>

int main() 
{
    char src[] = "Hello world";
    char dest[32];

    char *end = stpcpy(dest, src);

    printf("Source string: %s\n", src);
    printf("Destination string: %s\n", dest);
    printf("End pointer: %p\n", end);

    return 0;
}
```
在上述的示例代码中，我们使用 `stpcpy()` 函数将源字符串 `"Hello world"` 复制到目标字符串 `dest` 中，并输出两个字符串以及目标字符串的结尾指针。

**注意：** `stpcpy()` 函数只能在支持 `ISO C99` 或 `POSIX.1-2001` 标准的系统上使用，对于其他系统，可能需要使用 `strcpy()` 函数代替。此外，应始终确保目标字符串缓冲区具有足够的空间来存储源字符串，以避免发生缓冲区溢出。


# 14. strcat
## 14.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`char* strcat(char* dest, const char* src);` |  用于将一个字符串拼接到另一个字符串的末尾 |

**参数：**
- **dest ：** 目标字符串的缓冲区，必须具有足够的空间来存储源字符串
- **src ：** 要拼接的源字符串

**返回值：** 一个指向目标字符串结尾的指针

## 14.2 演示示例
```c
#include <stdio.h>
#include <string.h>

int main() 
{
    char str1[32] = "Hello";
    char str2[] = " world!";

    printf("Before: str1=%s str2=%s\n", str1, str2);

    strcat(str1, str2);

    printf("After: str1=%s\n", str1);

    return 0;
}
```

## 14.3 运行结果
![](strcat.png)


# 15. strchr
## 15.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`char* strchr(const char* str, int c);` | 用于查找字符串中第一次出现指定字符的位置，并返回该位置的指针  |

**参数：**
- **src ：** 要查找的字符串
- **c：** 要查找的字符，是一个整数值

## 15.2 演示示例
```c
#include <stdio.h>
#include <string.h>

int main() 
{
    char str[] = "Hello world";
    char* p;

    p = strchr(str, 'o');

    if (p != NULL) {
        printf("Found character '%c' at position %ld\n", *p, p - str);
    }
    else {
        printf("Character not found.\n");
    }

    return 0;
}
```

**注意：** `strchr()` 函数只能查找单个字符，如果要查找一个子字符串，应使用 `strstr()` 函数代替。另外，在查找字符时，需要将字符转换为整数值传递给 `strchr()` 函数，以避免发生类型错误。

## 15.3 运行结果
![](strchr.png)

# 16. strcmp
## 16.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int strcmp(const char* str1, const char* str2);` |  用于比较两个字符串是否相等 |

**参数：**
- **str1 ：** 要比较的第一个字符串
- **str2 ：** 要比较的第二个字符串

**返回值：** 一个整数，表示两个字符串之间的大小关系
- 如果 `str1` 小于 `str2`，则返回负整数；
- 如果 `str1` 大于 `str2`，则返回正整数；
- 如果 `str1` 等于 `str2`，则返回 0。

## 16.2 演示示例
```c
#include <stdio.h>
#include <string.h>

int main() 
{
    char str1[] = "hello World";
    char str2[] = "hello world";

    int result = strcmp(str1, str2);

    if (result < 0) {
        printf("'%s' is less than '%s'\n", str1, str2);
    }
    else if (result 0) {
        printf("'%s' is greater than '%s'\n", str1, str2);
    }
    else {
        printf("'%s' is equal to '%s'\n", str1, str2);
    }

    return 0;
}
```
## 16.3 运行结果
![](strcmp.png)


# 17. strcpy
## 17.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`char* strcpy(char* dest, const char* src);` | 用于将一个字符串复制到另一个字符串缓冲区中  |

**参数：**
- **dest ：** 目标字符串的缓冲区，必须具有足够的空间来存储源字符串
- **src ：** 要复制的源字符串


**返回值：** 一个指向目标字符串结尾的指针

## 17.2 演示示例
```c
#include <stdio.h>
#include <string.h>

int main() {
    char src[] = "Hello world";
    char dest[32];

    strcpy(dest, src);

    printf("Source string: %s\n", src);
    printf("Destination string: %s\n", dest);

    return 0;
}
```

**注意：** `strcpy()` 函数只能用于复制以 `\0` 结尾的字符串，否则可能导致未定义的行为或内存损坏。在调用 `strcpy()` 函数之前，应确保目标字符串缓冲区具有足够的空间来容纳源字符串，以避免发生缓冲区溢出。

## 17.3 运行结果
![](strcpy.png)


# 18. strcspn
## 18.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`size_t strcspn(const char* str, const char* charset);` | 用于查找字符串中第一次出现指定字符集合中任何字符的位置，并返回该位置的索引  |

**参数：**
- **src ：** 要查找的字符串
- **charset ：** 要查找的字符集合

## 18.2 演示示例
```c
#include <stdio.h>
#include <string.h>

int main() 
{
    char str[] = "Hello world";
    char charset[] = "owd";

    size_t index = strcspn(str, charset);

    if (index != strlen(str)) 
    {
        printf("Found character '%c' at position %ld\n", str[index], index);
    }
    else 
    {
        printf("No matching characters found.\n");
    }

    return 0;
}
```

在上述的示例代码中，我们使用 `strcspn()` 函数在字符串 `"Hello world"` 中查找字符集合 `"owd"` 中的任何字符，并输出找到的字符及其在字符串中的位置索引。

## 18.3 运行结果
![](strcspn.png)
