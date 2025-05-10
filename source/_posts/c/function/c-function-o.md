---
title: C语言函数大全--o 开头的函数
date: 2023-04-24 09:58:02
updated: 2025-05-10 21:36:59
categories:
  - 开发语言-C
tags:
  - C语言函数大全
  - o 开头的函数
---

![](/images/cplus-logo.png)

# 总览
| 函数声明 |  函数功能  |
|:--|:--|
|`void obstack_init(struct obstack *obstack_ptr);` |  它是 `POSIX` 标准库中的一个非标准函数，用于初始化对象堆栈。对象堆栈是一种可以动态增长以存储任意类型的对象的数据结构。其中，`obstack_ptr` 参数是一个指向 `struct obstack` 类型的指针，表示要初始化的对象堆栈。 |
|`void obstack_free(struct obstack *obstack_ptr, void *object_ptr);` | 用于释放通过对象堆栈分配的所有内存。其中，`obstack_ptr` 参数是一个指向 `struct obstack` 类型的指针，表示要释放的对象堆栈；`object_ptr` 参数是要释放的内存块。  |
|`void *obstack_alloc(struct obstack *obstack_ptr, int size);` | 用于从对象堆栈中分配指定大小的内存，并返回其地址。其中，`obstack_ptr` 参数是一个指向 `struct obstack` 类型的指针，表示要从中分配内存的对象堆栈；`size` 参数表示要分配的内存块的大小。  |
|`void *obstack_blank(struct obstack *obstack_ptr, int size);` | 用于向对象堆栈添加指定数量的空间，并返回指向添加的第一个字节的指针。其中，`obstack_ptr` 参数是一个指向 `struct obstack` 类型的指针，表示要添加空间的对象堆栈；`size` 参数表示要添加的空间大小。 |
|`void *obstack_grow(struct obstack *obstack_ptr, const void *data, int size);` |  用于将数据复制到对象堆栈，并返回指向添加的第一个字节的指针。其中，`obstack_ptr` 参数是一个指向 `struct obstack` 类型的指针，表示要添加数据的对象堆栈；`data` 参数是要复制的数据的指针；`size` 参数表示要复制的数据的大小。 |
|`#define offsetof(type, member) ((size_t)(&((type *)0)->member))` | 它是一个宏，用于获取结构体中某个成员的偏移量。  |
|`int open(const char *path, int oflag, ...);` |  用于打开文件 |
|`int openat(int dirfd, const char *pathname, int flags, mode_t mode);` | 它是 **Linux** 系统定义的一个函数，它可以打开一个相对于指定目录的文件。与 `open()` 函数相比，`openat()` 函数更加灵活，并支持更多的选项。  |
|`DIR *opendir(const char *name);` | 它是 **POSIX** 标准定义的一个函数，用于打开目录并返回一个指向 `DIR` 结构体类型的指针。  |
|`int openpty(int *amaster, int *aslave, char *name, const struct termios *termp, const struct winsize *winp);` | 它是 **POSIX** 标准定义的一个函数，用于打开一个伪终端（**PTY**）并返回与之关联的主从设备文件描述符。伪终端可以用于在进程之间建立通信，或者在程序中模拟终端行为。  |
|`int on_exit(void (*function)(int, void *), void *arg);` |  它是 **POSIX** 标准定义的一个函数，用于在进程退出时调用注册的回调函数。这个函数可以用于在程序异常退出或者正常退出时执行一些清理工作、记录日志等操作 |
|`void outtext(char *textstring);` |  在图形视区显示一个字符串|
|`void outtextxy(int x, int y, char *textstring);` | 在指定位置显示一字符串  |



# 1. obstack_init，obstack_free，obstack_alloc，obstack_blank，obstack_grow
## 1.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void obstack_init(struct obstack *obstack_ptr);` |  它是 `POSIX` 标准库中的一个非标准函数，用于初始化对象堆栈。对象堆栈是一种可以动态增长以存储任意类型的对象的数据结构。其中，`obstack_ptr` 参数是一个指向 `struct obstack` 类型的指针，表示要初始化的对象堆栈。 |
|`void obstack_free(struct obstack *obstack_ptr, void *object_ptr);` | 用于释放通过对象堆栈分配的所有内存。其中，`obstack_ptr` 参数是一个指向 `struct obstack` 类型的指针，表示要释放的对象堆栈；`object_ptr` 参数是要释放的内存块。  |
|`void *obstack_alloc(struct obstack *obstack_ptr, int size);` | 用于从对象堆栈中分配指定大小的内存，并返回其地址。其中，`obstack_ptr` 参数是一个指向 `struct obstack` 类型的指针，表示要从中分配内存的对象堆栈；`size` 参数表示要分配的内存块的大小。  |
|`void *obstack_blank(struct obstack *obstack_ptr, int size);` | 用于向对象堆栈添加指定数量的空间，并返回指向添加的第一个字节的指针。其中，`obstack_ptr` 参数是一个指向 `struct obstack` 类型的指针，表示要添加空间的对象堆栈；`size` 参数表示要添加的空间大小。 |
|`void *obstack_grow(struct obstack *obstack_ptr, const void *data, int size);` |  用于将数据复制到对象堆栈，并返回指向添加的第一个字节的指针。其中，`obstack_ptr` 参数是一个指向 `struct obstack` 类型的指针，表示要添加数据的对象堆栈；`data` 参数是要复制的数据的指针；`size` 参数表示要复制的数据的大小。 |

## 1.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <obstack.h>

int main(void)
{
    struct obstack my_obstack;
    const char *str1 = "Hello, ";
    const char *str2 = "World!";
    char *dst;

    obstack_init(&my_obstack);

    dst = (char *)obstack_alloc(&my_obstack, strlen(str1) + strlen(str2) + 1);
    strcpy(dst, str1);
    strcat(dst, str2);

    printf("%s\n", (char *)my_obstack.chunk);

    dst = (char *)obstack_blank(&my_obstack, sizeof(int)*2);
    int a = 100;
    int b = 200;
    memcpy(dst, &a, sizeof(int));
    memcpy(dst+sizeof(int), &b, sizeof(int));

    printf("%d %d\n", *(int *)(my_obstack.next_free-sizeof(int)*2), *(int *)(my_obstack.next_free-sizeof(int)));

    double d = 3.1415926;
    dst = (char *)obstack_grow(&my_obstack, &d, sizeof(double));

    printf("%f\n", *(double *)(my_obstack.next_free-sizeof(double)));

    obstack_free(&my_obstack, NULL);

    return 0;
}
```
在上述的程序中，
- 我们首先定义一个名为 `my_obstack` 的 `struct obstack` 类型变量，并将其传递给 `obstack_init()` 函数以初始化对象堆栈。
- 接着，我们使用 `obstack_alloc()` 函数从对象堆栈中分配一块内存，并将两个字符串连接起来。
- 然后，我们使用 `obstack_blank()` 函数向对象堆栈添加一块指定大小的空间，并使用 `memcpy()` 函数将两个整数复制到该空间中。
- 接下来，我们使用 `obstack_grow()` 函数向对象堆栈添加一个双精度浮点数，并返回指向该浮点数的指针。
- 最后，我们使用 `printf()` 函数将连接后的字符串、添加的整数和添加的双精度浮点数输出到终端，并使用 `obstack_free()` 函数释放通过对象堆栈分配的所有内存。

> 注意：在使用 `obstack_blank()` 函数向对象堆栈添加空间时，建议使用 `sizeof` 运算符来计算要添加的空间大小，并在使用 `memcpy()` 复制数据时也应该小心不要越界。同时，在使用 `obstack_grow()` 函数向对象堆栈添加数据时，需要小心指针是否正确，并且操作前需要先使用 `memcpy()` 将要添加的数据复制到一个临时变量中。

# 2. offsetof
## 2.1 宏说明
| 宏定义 |  宏描述  |
|:--|:--|
|`#define offsetof(type, member) ((size_t)(&((type *)0)->member))` | 它是一个宏，用于获取结构体中某个成员的偏移量。  |

**参数：**
- **type ：** 表示结构体类型
- **member ：** 表示结构体中的一个成员变量名

**返回值：** 一个 size_t 类型的值，表示该成员变量在结构体中的偏移量（单位是字节）。

## 2.2 演示示例

```c
#include <stdio.h>
#include <stddef.h>

struct example {
    int a;
    char b;
    double c;
};

int main(void)
{
    size_t offset_b = offsetof(struct example, b);
    printf("Offset of 'b' in struct example: %zu\n", offset_b);

    return 0;
}
```

在这个程序中，
- 我们定义了一个名为 `example` 的结构体类型，并使用 `offsetof` 宏获取结构体中成员变量 `b` 的偏移量。
- 最后，我们使用 `printf()` 函数将结果输出到终端。

> **注意：** 在使用 `offsetof` 宏时，结构体类型名称必须使用括号括起来，否则代码会产生语法错误。此外，`offsetof` 宏的参数必须是已定义的结构体类型名称和该结构体类型中的成员变量名称，否则也会导致编译错误。

# 3. open
## 3.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int open(const char *path, int oflag, ...);` |  用于打开文件 |

**参数：**
- **path ：** 表示要打开的文件路径
- **oflag ：** 表示打开文件时的选项标志，可以为以下常量之一或多个按位或组合而成：
    - `O_RDONLY`：只读模式打开文件。
    - `O_WRONLY`：只写模式打开文件。
    - `O_RDWR`：读写模式打开文件。
    - `O_CREAT`：如果文件不存在，则创建它。
    - `O_TRUNC`：如果文件已存在，则将其长度截断为 0。
    - `O_APPEND`：在文件末尾追加数据。
- **可选参数 ：** 表示文件所有者、组和其他用户的访问权限。如果使用了 `O_CREAT` 选项，则必须提供这个参数

## 3.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>

int main(void)
{
    int fd = open("temp.txt", O_RDONLY);
    if (fd == -1) {
        perror("open");
        exit(1);
    }

    char buf[1024];
    ssize_t nread;
    while ((nread = read(fd, buf, sizeof(buf))) > 0) {
        if (write(STDOUT_FILENO, buf, nread) != nread) {
            perror("write");
            exit(1);
        }
    }

    if (nread == -1) {
        perror("read");
        exit(1);
    }

    if (close(fd) == -1) {
        perror("close");
        exit(1);
    }

    return 0;
}
```
在上述的程序中，
- 我们使用 `open()` 函数打开文件 `temp.txt`，并通过 `read()` 函数读取其中的数据。
- 然后，我们使用 `write()` 函数将数据写入标准输出，直到读取完整个文件。
- 最后，我们使用 `close()` 函数关闭文件。

> **注意：** 在使用 `open()` 函数打开文件时，返回值为负数则表示出现了错误。此时可以使用 `perror()` 函数输出错误信息，并使用 `exit()` 函数退出程序。同时，在使用 `read()` 函数和 `write()` 函数读写文件时也需要小心处理返回值，以避免出现不可预期的错误。

## 3.3 运行结果
![ ](open-file.png)

![ ](open.png)

# 4. openat
## 4.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int openat(int dirfd, const char *pathname, int flags, mode_t mode);` | 它是 **Linux** 系统定义的一个函数，它可以打开一个相对于指定目录的文件。与 `open()` 函数相比，`openat()` 函数更加灵活，并支持更多的选项。  |

**参数：**
- **dirfd ：** 表示要打开文件所在目录的文件描述符。如果传递的是 `AT_FDCWD`，则表示使用当前工作目录。
- **pathname ：** 表示要打开的文件路径
- **flags ：** 表示打开文件时的选项标志，可以为以下常量之一或多个按位或组合而成：
    - `O_RDONLY`：只读模式打开文件。
    - `O_WRONLY`：只写模式打开文件。
    - `O_RDWR`：读写模式打开文件。
    - `O_CREAT`：如果文件不存在，则创建它。
    - `O_TRUNC`：如果文件已存在，则将其长度截断为 0。
    - `O_APPEND`：在文件末尾追加数据。
    - `O_DIRECTORY`：要求打开的文件必须是一个目录。
    - `O_NOFOLLOW`：不跟随符号链接打开文件。
- **mode ：** 表示文件所有者、组和其他用户的访问权限。如果使用了 `O_CREAT` 选项，则必须提供这个参数

## 4.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>
#include <dirent.h>

int main(void)
{
    int dirfd = open(".", O_RDONLY | O_DIRECTORY);
    if (dirfd == -1) {
        perror("open");
        exit(1);
    }

    DIR *dirp = fdopendir(dirfd);
    if (dirp == NULL) {
        perror("fdopendir");
        exit(1);
    }

    struct dirent *entry;
    while ((entry = readdir(dirp)) != NULL) {
        printf("%s\n", entry->d_name);
    }

    if (closedir(dirp) == -1) {
        perror("closedir");
        exit(1);
    }

    return 0;
}
```
在这个程序中，
- 我们使用 `openat()` 函数打开当前目录，并通过 `fdopendir()` 函数将文件描述符转换为目录流。
- 然后，我们使用 `readdir()` 函数读取目录中的文件，并将文件名输出到终端。
- 最后，我们使用 `closedir()` 函数关闭目录。

> **注意：** 在使用 `openat()` 函数打开文件时，可以通过传递不同的文件描述符指定要打开的目录，从而实现更加灵活的操作。此外，在使用 `readdir()` 函数读取目录时也需要注意判断返回值是否为 `NULL`，以避免出现不可预期的错误。

# 5. opendir
## 5.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`DIR *opendir(const char *name);` | 它是 **POSIX** 标准定义的一个函数，用于打开目录并返回一个指向 `DIR` 结构体类型的指针。  |

**参数：**
- **name ：** 表示要打开的目录路径

**返回值：**
- 如果该函数执行成功，则返回一个指向 `DIR` 类型的指针；
- 否则返回 `NULL`。

## 5.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <dirent.h>

int main(void)
{
    DIR *dirp = opendir(".");
    if (dirp == NULL) {
        perror("opendir");
        exit(1);
    }

    struct dirent *entry;
    while ((entry = readdir(dirp)) != NULL) {
        printf("%s\n", entry->d_name);
    }

    if (closedir(dirp) == -1) {
        perror("closedir");
        exit(1);
    }

    return 0;
}
```
在上述的程序中，我们使用 `opendir()` 函数打开当前目录，并通过 `readdir()` 函数读取目录中的文件名，最后使用 `closedir()` 函数关闭目录。

> **注意：** 在使用 `opendir()` 函数打开目录时，返回值为 `NULL` 则表示出现了错误。此时可以使用 `perror()` 函数输出错误信息，并使用 `exit()` 函数退出程序。同时，在使用 `readdir()` 函数读取目录时也需要小心处理返回值，以避免出现不可预期的错误。

# 6. openpty
## 6.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int openpty(int *amaster, int *aslave, char *name, const struct termios *termp, const struct winsize *winp);` | 它是 **POSIX** 标准定义的一个函数，用于打开一个伪终端（**PTY**）并返回与之关联的主从设备文件描述符。伪终端可以用于在进程之间建立通信，或者在程序中模拟终端行为。  |

**参数：**
- **amaster ：** 表示要返回的主设备文件描述符
- **aslave ：** 表示要返回的从设备文件描述符
- **name ：** 表示从设备名称的缓冲区，如果不需要则可以传递 `NULL`
- **termp ：** 表示要使用的终端属性，如果不需要则可以传递 `NULL`
- **winp ：** 表示要使用的窗口大小，如果不需要则可以传递 `NULL`

**返回值：**
- 如果该函数执行成功，则返回值为 `0`；
- 否则返回 `-1`。

## 6.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>
#include <termios.h>

int main(void)
{
    int master, slave;
    char buf[1024];
    ssize_t nread;

    if (openpty(&master, &slave, NULL, NULL, NULL) == -1) {
        perror("openpty");
        exit(1);
    }

    printf("Slave device: /dev/pts/%d\n", slave);

    while ((nread = read(STDIN_FILENO, buf, sizeof(buf))) > 0) {
        if (write(master, buf, nread) != nread) {
            perror("write");
            exit(1);
        }
    }

    if (nread == -1) {
        perror("read");
        exit(1);
    }

    if (close(master) == -1) {
        perror("close");
        exit(1);
    }

    return 0;
}
```
在上述的程序中，
- 我们使用 `openpty()` 函数打开一个伪终端，并通过 `read()` 函数读取标准输入中的数据。
-  然后，我们将数据写入主设备文件描述符；
- 最后关闭该设备。

> **注意：** 在使用 `openpty()` 函数打开伪终端时，返回值为 `-1` 则表示出现了错误。此时可以使用 `perror()` 函数输出错误信息，并使用 `exit()` 函数退出程序。同时，在使用 `read()` 函数和 `write()` 函数读写文件时也需要小心处理返回值，以避免出现不可预期的错误。

# 7. on_exit
## 7.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int on_exit(void (*function)(int, void *), void *arg);` |  它是 **POSIX** 标准定义的一个函数，用于在进程退出时调用注册的回调函数。这个函数可以用于在程序异常退出或者正常退出时执行一些清理工作、记录日志等操作 |

**参数：**
- **function ：** 表示要注册的回调函数
- **arg ：** 表示传递给回调函数的参数

**返回值：**
- 如果该函数执行成功，则返回值为 `0`；
- 否则返回 `-1`。

## 7.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

void cleanup(int status, void *arg)
{
    printf("Cleanup function called with status %d\n", status);
}

int main(void)
{
    if (on_exit(cleanup, NULL) != 0) {
        perror("on_exit");
        exit(EXIT_FAILURE);
    }

    printf("This is the main program\n");

    return 0;
}
```
在如上的程序中，
- 我们使用 `on_exit()` 函数注册了一个回调函数 `cleanup()`，并将其参数设置为 NULL。
- 然后，在主函数中输出了一条消息。当程序退出时，会自动调用回调函数来进行清理操作。

# 8. outtext
## 8.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void outtext(char *textstring);` |  在图形视区显示一个字符串|

**参数：**
- `char *textstring` ：指向以空字符（'\0'）结尾的字符串的指针。

## 8.2 演示示例
```c
#include <graphics.h>

int main(void)
{
   int gdriver = DETECT, gmode, errorcode;
   int midx, midy;

   initgraph(&gdriver, &gmode, "");

   midx = getmaxx() / 2;
   midy = getmaxy() / 2;

   moveto(midx, midy);

   outtext("This ");
   outtext("is ");
   outtext("a ");
   outtext("test.");

   getch();
   closegraph();
   return 0;
}
```

在上述的程序中，
- 我们首先调用 `initgraph()` 函数初始化图形系统；
- 然后获取窗口的中心坐标；
- 接着使用 `moveto()` 函数将当前绘图位置移动到屏幕中心。
- 最后，使用 `outtext()` 函数输出一段文字，然后等待用户按下任意键，并关闭图形窗口。

## 8.3 运行结果
![ ](outtext.png)

# 9. outtextxy
## 9.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void outtextxy(int x, int y, char *textstring);` | 在指定位置显示一字符串  |

**参数：**
- `int x` : 字符串输出的水平起始坐标（单位为像素）。取值范围：与当前图形模式的屏幕分辨率相关（例如，640x480 模式下，x 范围为 0 到 639）。

- `int y` : 字符串输出的垂直起始坐标（单位为像素）。取值范围：与当前图形模式的屏幕分辨率相关（例如，480p 模式下，y 范围为 0 到 479）。

- `char *textstring` : 指向以空字符（'\0'）结尾的字符串的指针。

## 9.2 演示示例
```c
#include <graphics.h>

int main(void)
{
   int gdriver = DETECT, gmode, errorcode;
   int x, y;

   initgraph(&gdriver, &gmode, "");

   x = getmaxx() / 2;
   y = getmaxy() / 2;

   outtextxy(x, y, "Hello, world!");

   getch();
   closegraph();
   return 0;
}

```

在上述这个程序中，
- 我们首先通过 `initgraph()` 函数初始化图形系统并创建一个窗口;
- 然后，定义了一个坐标位置 `(x, y)` 并使用 `outtextxy()` 函数在该位置输出一段文本。
- 最后，使用 `getch` 函数等待用户按下任意键，然后关闭图形窗口。

## 9.3 运行结果
![ ](outtextxy.png)

