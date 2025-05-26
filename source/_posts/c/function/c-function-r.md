---
title: C语言函数大全--r 开头的函数
date: 2023-04-27 10:12:34
updated: 2025-05-26 19:57:38
categories:
  - 开发语言-C
tags:
  - C语言函数大全
  - r 开头的函数
---

![](/images/cplus-logo.png)

# 总览
| 函数声明 |  函数功能  |
|:--|:--|
|`int raise(int sig);` | 用于向当前进程发送指定的信号。  |
|`int rand(void);` | 用于生成伪随机数  |
|`ssize_t read(int fd, void *buf, size_t count);` | 用于从文件描述符读取数据的函数。  |
|`void *realloc(void *ptr, size_t size);` | 用于重新分配已经分配过内存的空间大小。  |
|`void rectangle( int left, int top, int right, int bottom);` |  画一个矩形 |
|`int registerbgidriver(void(*driver)(void));` | 用于将BGI（Borland Graphics Interface）驱动程序注册到系统中 |
|`int remove(char *filename);` |  用于删除指定的文件 |
|`int rename(char *oldname, char *newname); ` | 用于重命名或移动文件。  |
|``void restorecrtmode(void);`` |  将图形模式恢复到文本模式 |
|`void rewind(FILE *stream);` | 将文件指针 `stream` 指向的文件位置重置为文件开头，同时清除任何错误或文件结束标志。  |
|`int rmdir(const char *path);` |  用于删除一个空目录，即该目录必须为空。 |
|`double round(double x);` |  将传入的实数参数 `x` 四舍五入为最接近的整数（double） |
|`float roundf(float x);` | 将传入的实数参数 `x` 四舍五入为最接近的整数（float）  |
|`long double roundl(long double x);` | 将传入的实数参数 `x` 四舍五入为最接近的整数（long double）  |
|`double remainder(double x, double y);` | 用于计算两个浮点数的余数（即模运算结果）（double）|
|`float remainderf (float x, float y );` | 用于计算两个浮点数的余数（即模运算结果）（float）  |
|`long double remainderl (long double x, long double y);` | 用于计算两个浮点数的余数（即模运算结果）（long double）  |
|`double remquo(double x, double y, int *quo);` | 用于计算两个浮点数的余数，并返回商和余数。  |
|`float remquof(float x, float y, int *quo);` | 用于计算两个浮点数的余数，并返回商和余数。  |
|`long double remquol(long double x, long double y, int *quo);` |  用于计算两个浮点数的余数，并返回商和余数。 |
|`double rint(double x);` | 将 x 四舍五入到最接近的整数（double） |
|`float rintf(float x);` |  将 x 四舍五入到最接近的整数（float） |
|`long double rintl(long double x);` | 将 x 四舍五入到最接近的整数（long double）  |


# 1. raise
## 1.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int raise(int sig);` | 用于向当前进程发送指定的信号。  |

**参数：**
- **sig ：** 指定要发送的信号编号

**返回值：**
- 如果调用成功，`raise()` 函数将返回 `0`；
- 否则，它将返回一个非零值。

## 1.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <signal.h>

void sigint_handler(int sig) {
    printf("Caught signal %d\n", sig);
    exit(0);
}

int main() {
    signal(SIGINT, sigint_handler);

    while (1) {
        printf("Doing something...\n");
        getchar();
        // check for interrupt signal
        raise(SIGINT);
    }

    return 0;
}
```
在上面的示例中，
- 我们首先使用 `signal()` 函数设置了一个处理 `SIGINT` 信号的处理程序 `sigint_handler()`。
- 然后，在主循环中，我们随意输入一个字符后，就使用 raise() 函数向当前正在运行的进程发送 **SIGINT** 信号。当收到 **SIGINT** 信号时，程序将打印一条消息并退出。

>**注意：** `raise()` 函数只能向当前进程发送信号，不能向其他进程发送信号。如果要向其他进程发送信号，可以使用 `kill()` 函数。

## 1.3 运行结果
![](raise.png)

# 2. rand
## 2.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int rand(void);` | 用于生成伪随机数  |

**返回值：**
 每次调用它时会返回一个介于 `0` 和 `RAND_MAX` 之间的伪随机整数。其中，`RAND_MAX` 是一个常量，表示返回值的最大值，通常为 `32767`。

## 2.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <time.h>

int main() {
    // 使用系统时间作为随机数种子初始化
    srand(time(NULL));

    // 生成并输出 10 个随机数
    for (int i = 0; i < 10; i++) {
        printf("%d\n", rand());
    }

    return 0;
}
```
**注意：**
- 每次程序运行时，`rand()` 函数返回的随机数序列都是相同的。如果要生成不同的随机数序列，可以使用 `srand()` 函数提供的种子来初始化随机数发生器。
- 由于 `rand()` 函数只能生成伪随机数，不能保证其真正的随机性。因此，在需要高度安全性的应用程序中，建议使用更加安全的随机数生成器，如 `/dev/random` 和 `/dev/urandom` 等系统提供的硬件随机数生成器。

## 2.3 运行结果
![](rand.png)

# 3. read
## 3.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`ssize_t read(int fd, void *buf, size_t count);` | 用于从文件描述符读取数据的函数。  |

**参数：**
- **fd ：** 要读取的文件描述符
- **buf ：** 存储读取数据的缓冲区
- **count ：** 要读取的字节数。

## 3.2 演示示例
```c
#include <stdio.h>
#include <unistd.h>

#define BUFFER_SIZE 1024

int main() {
    char buffer[BUFFER_SIZE];
    ssize_t num_read;

    num_read = read(STDIN_FILENO, buffer, BUFFER_SIZE);
    if (num_read == -1) {
        perror("read");
        return 1;
    }

    printf("Read %ld bytes from standard input: %.*s", num_read, (int)num_read, buffer);

    return 0;
}
```

在上述的示例中，
- 我们首先定义了一个大小为 `BUFFER_SIZE` 的字符数组 `buffer`；
- 然后使用 `read()` 函数读取标准输入【`STDIN_FILENO` 表示标准输入的文件描述符】中的数据到 `buffer` 中；
- 最后将读取的结果输出到控制台上。

## 3.3 运行结果
![](read.png)

# 4. realloc
## 4.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void *realloc(void *ptr, size_t size);` | 用于重新分配已经分配过内存的空间大小。  |

**参数：**
- **ptr ：** 指向已分配内存区域的指针
- **size ：** 需要分配的新内存大小

**返回值：**
- 它会尝试重新调整已经分配给 `ptr` 指向的内存区域大小，并返回一个指向新内存起始地址的指针。
- 如果无法分配新的内存，则返回空指针 `NULL`。

注意： `realloc()` 函数并不保证原有的内存区域内容会被完全保留。当内存大小增加时，可能会分配新的内存并将原有的数据复制到新内存区域中；当内存大小减小时，则可能截断部分原有的数据。因此，在使用 `realloc()` 函数时，需要注意备份原有的数据，以免出现数据丢失的情况。

## 4.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>

int main() {
    // 分配 10 个整型变量的内存空间
    int *p = (int*) malloc(10 * sizeof(int));
    if (p == NULL) {
        printf("内存分配失败\n");
        exit(1);
    }

    // 更改为 20 个整型变量的内存空间
    int *q = (int*) realloc(p, 20 * sizeof(int));
    if (q == NULL) {
        printf("内存分配失败\n");
        exit(1);
    }

    // 输出新的内存地址
    printf("原内存地址：%p，新内存地址：%p\n", p, q);

    // 释放新分配的内存空间
    free(q);

    return 0;
}
```
在上面的示例中，
- 我们首先使用 `malloc()` 函数分配了 `10` 个整型变量的内存空间，并判断是否分配成功。
- 接着，我们使用 `realloc()` 函数将分配的内存空间大小扩充为 `20` 个整型变量，并判断是否分配成功。
- 最后，我们输出了原有的内存地址和新的内存地址，并释放了新分配的内存空间。

**注意：** 在使用 `realloc()` 函数时，应该始终检查返回值，以确保分配内存空间成功，避免因内存不足或其他原因导致程序崩溃。另外，也应该避免过度使用 `realloc()` 函数，因为频繁地重新分配内存会影响程序性能。

# 5. rectangle
## 5.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void rectangle( int left, int top, int right, int bottom);` |  画一个矩形 |

**参数：**
left, top, right, bottom，它们分别表示矩形左上角和右下角的坐标。

## 5.2 演示示例
```c
#include <graphics.h>

int main(void)
{
    int gdriver = DETECT, gmode;
    int left, top, right, bottom;

    initgraph(&gdriver, &gmode, "");

    left = getmaxx() / 2 - 50;
    top = getmaxy() / 2 - 50;
    right = getmaxx() / 2 + 50;
    bottom = getmaxy() / 2 + 50;

    rectangle(left,top,right,bottom);

    getch();
    closegraph();
    return 0;
}
```

## 5.3 运行结果
![](rectangle.png)

# 6. registerbgidriver
## 6.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int registerbgidriver(void(*driver)(void));` | 用于将BGI（Borland Graphics Interface）驱动程序注册到系统中 |

**注意：**
它必须在使用任何 **BGI** 图形函数之前调用。该函数接受一个指向驱动程序结构的指针作为参数，并返回一个整数值以指示是否成功注册了驱动程序。**BGI** 驱动程序主要用于支持 **Borland C++** 等 **IDE** 环境下的图形绘制和显示操作，它们通常存储在一个单独的库文件中，例如 `graphics.h` 头文件需要使用的 **BGI driver** 位于 `libbgi.a` 文件中。

## 6.2 演示示例
```c
#include <graphics.h>

int main()
{
    int gd = DETECT, gm;
    initgraph(&gd, &gm, "");
    registerbgidriver(EGAVGA_driver); // 注册EGAVGA驱动程序
    circle(100, 100, 50); // 绘制圆形
    getch(); // 等待用户输入
    closegraph(); // 关闭图形界面
    return 0;
}
```


# 7. remove
## 7.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int remove(char *filename);` |  用于删除指定的文件 |

**参数：**
- **filename ：** 一个指向要删除文件的文件名的字符串指针

**返回值：** 
返回一个整数表示操作是否成功
- 如果文件成功删除，则返回 `0`；
- 否则，返回非零值以指示发生错误。

## 7.2 演示示例
```c
#include <stdio.h>

int main()
{
    int result;
    char filename[] = "example.txt";
    result = remove(filename);
    if (result == 0)
    {
        printf("The file %s has been removed.\n", filename);
    }
    else
    {
        printf("Error deleting the file %s.\n", filename);
    }
    return 0;
}
```
## 7.3 运行结果
![](remove.png)
![](remove-1.png)

# 8. rename
## 8.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int rename(char *oldname, char *newname); ` | 用于重命名或移动文件。  |

**参数：**
- `oldname ：` 指定要更改的文件名
- `newname ：` 指定新文件名或包含新路径的新文件名

**返回值：**
- 在执行成功时，该函数返回 `0`；
- 否则返回非零值以指示发生错误。

## 8.2 演示示例
```c
#include <stdio.h>

int main()
{
    char oldname[] = "temp.txt";
    char newname[] = "tempnew.txt";
    int result;
    
    result = rename(oldname, newname);
    
    if (result == 0)
    {
        printf("The file has been renamed.\n");
    }
    else
    {
        printf("Error renaming the file.\n");
    }
    
    return 0;
}
```

## 8.3 运行结果
![](rename.png)

# 9. restorecrtmode
## 9.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|``void restorecrtmode(void);`` |  将图形模式恢复到文本模式 |

## 9.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>

int main(void)
{
    int gdriver = DETECT, gmode;
    int x, y;

    initgraph(&gdriver, &gmode, "");

    x = getmaxx() / 2;
    y = getmaxy() / 2;

    settextjustify(CENTER_TEXT, CENTER_TEXT);
    outtextxy(x, y, "Press any key to exit graphics:");
    getch();

    // 将图形模式恢复到文本模式
    restorecrtmode();
    printf("We're now in text mode.\n");
    printf("Press any key to return to graphics mode:");
    getch();

    // 返回图形模式
    setgraphmode(getgraphmode());

    settextjustify(CENTER_TEXT, CENTER_TEXT);
    outtextxy(x, y, "We're back in graphics mode.");
    outtextxy(x, y+textheight("W"), "Press any key to halt:");

    getch();
    closegraph();
    return 0;
}
```

## 9.3 运行结果
![](restorecrtmode.gif)

# 10. rewind
## 10.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void rewind(FILE *stream);` | 将文件指针 `stream` 指向的文件位置重置为文件开头，同时清除任何错误或文件结束标志。  |

## 10.2 演示示例
```c
#include <stdio.h>

int main() {
    FILE *fp = fopen("tempnew.txt", "r");
    if (fp == NULL) {
        perror("Failed to open file");
        return 1;
    }
    
    // 读取文件内容
    char buffer[1024];
    while (fgets(buffer, sizeof(buffer), fp)) {
        printf("%s", buffer);
    }

    printf("\n");
    
    // 将文件指针重置到开头
    rewind(fp);
    
    // 再次读取文件内容
    while (fgets(buffer, sizeof(buffer), fp)) {
        printf("%s", buffer);
    }
    
    fclose(fp);
    return 0;
}
```

在上述的示例中，
- 我们首先打开一个名为 `tempnew.txt` 的文件；
- 然后使用 `fgets()` 函数从文件中读取文本行，并输出内容；
- 接着使用 `rewind()` 函数将文件指针重置到文件开头，并再次读取文件内容并输出；
- 最后关闭文件，退出程序。

## 10.3 运行结果
![](rewind.png)

![](rewind-1.png)

# 11. rmdir 
## 11.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int rmdir(const char *path);` |  用于删除一个空目录，即该目录必须为空。 |

**参数：**
- **path ：** 待删除的空目录路径

**返回值：**
 - 如果成功，则返回 `0`；
 - 否则返回 `-1`，并设置 `errno` 错误码。

## 11.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
int main() {
    const char *dir_name = "test";
    if (mkdir(dir_name) I= 0) {
        perror("Failed to create directory");
        return 1;
    }

    //在目录中创建一些文件
    FILE *fp1 = fopen("test/file1.txt","w");
    FILE *fp2 = fopen("test/file2.txt","w");
    fclose(fp1);
    fclose(fp2);
    
    /／尝试刪除非空目录
    if (rmdir(dir name) |= 0) {
        perror("Failed to remove directory");
        return 1;
    }
    printf("Directory removed successfullyIn");
    return 0;
}
```

在上述的示例中，
- 我们首先创建一个名为 `test` 的目录；
- 然后在其中创建两个文件 `file1.txt` 和 `file2.txt`；
- 接着使用 `rmdir()` 函数尝试删除该目录，但会失败，因为该目录不是空的。

**注意：** 如果要删除非空目录，可以使用上面的 `remove()` 函数来删除目录及其所有内容。

## 11.3 运行结果
![](rmdir.png)

# 12. round，roundf，roundl
## 12.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double round(double x);` |  将传入的实数参数 `x` 四舍五入为最接近的整数（double） |
|`float roundf(float x);` | 将传入的实数参数 `x` 四舍五入为最接近的整数（float）  |
|`long double roundl(long double x);` | 将传入的实数参数 `x` 四舍五入为最接近的整数（long double）  |

## 12.2 演示示例
```c
#include <math.h>
#include <stdio.h>

int main() {
    double x = 3.14159265;
    
    double rounded1 = round(x);     // 将 x 四舍五入为整数，结果为 3
    double rounded2 = round(x * 100) / 100;  // 将 x 精确到小数点后两位，结果为 3.14
    
    printf("rounded1: %lf\n", rounded1);
    printf("rounded2: %lf\n", rounded2);

    float xf = 2.5;
    printf("rountf(%f) = %f\n", xf, roundf(xf)); // 将 xf 四舍五入为整数，结果为 3
    
    long double xL = 2.4;
    printf("rountl(%Lf) = %Lf", xL, roundl(xL)); // 将 xL 四舍五入为整数，结果为 2
    return 0;
}
```

## 12.3 运行结果
![](round.png)

# 13. remainder，remainderf，remainderl
## 13.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double remainder(double x, double y);` | 用于计算两个浮点数的余数（即模运算结果）（double）|
|`float remainderf (float x, float y );` | 用于计算两个浮点数的余数（即模运算结果）（float）  |
|`long double remainderl (long double x, long double y);` | 用于计算两个浮点数的余数（即模运算结果）（long double）  |

**参数：**
- **x ：** 被除数
- **y ：** 除数

**返回值：** `x/y` 的余数

## 13.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main() {
    double x = 10.5, y = 3.2;
    printf("remainder(%lf, %lf) = %.20lf\n", x, y, remainder(x, y));


    float xf = 10.5f, yf = 3.2f;
    printf("remainder(%f, %f) = %.20f\n", xf, yf, remainderf(xf, yf));


    long double xL = 10.5L, yL = 3.2L;
    printf("remainder(%Lf, %Lf) = %.20Lf\n", xL, yL, remainderl(xL, yL));
    return 0;
}
```

## 13.3 运行结果

![](remainder.png)

# 14. remquo，remquof，remquol
## 14.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double remquo(double x, double y, int *quo);` | 用于计算两个浮点数的余数，并返回商和余数。  |
|`float remquof(float x, float y, int *quo);` | 用于计算两个浮点数的余数，并返回商和余数。  |
|`long double remquol(long double x, long double y, int *quo);` |  用于计算两个浮点数的余数，并返回商和余数。 |

**参数：**
- **x ：** 被除数
- **y ：** 除数
- **quo ：** 返回商的指针。如果不需要返回商，则可以传递 `NULL`

**返回值：** 返回 `x/y` 的余数，并通过 `quo` 指针返回商。

## 14.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main() {
    int quo;
    double x = 10.5, y = 3.2;
    double rem = remquo(x, y, &quo);
    printf("remquo(%lf, %lf, quo) = %.20lf, quo = %d\n", x, y, rem, quo);

    float xf = 10.5f, yf = 3.2f;
    float remf = remquof(xf, yf, &quo);
    printf("remquof(%f, %f, quo) = %.20f, quo = %d\n", xf, yf, remf, quo);

    long double xL = 10.5L, yL = 3.2L;
    long double remL = remquol(xL, yL, &quo);
    printf("remquol(%Lf, %Lf, quo) = %.20Lf, quo = %d\n", xL, yL, remL, quo);

    return 0;
}
```

## 14.3 运行结果
![](remquo.png)

# 15. rint，rintf，rintl
## 15.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double rint(double x);` | 将 x 四舍五入到最接近的整数（double） |
|`float rintf(float x);` |  将 x 四舍五入到最接近的整数（float） |
|`long double rintl(long double x);` | 将 x 四舍五入到最接近的整数（long double）  |

## 15.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main() {
    double d = 3.7;
    float f = 2.4f;
    long double ld = 5.9L;

    printf("rint(%f) = %.1f\n", d, rint(d));   // Output: rint(3.700000) = 4.0
    printf("rintf(%f) = %.1f\n", f, rintf(f)); // Output: rintf(2.400000) = 2.0
    printf("rintl(%Lf) = %.1Lf\n", ld, rintl(ld)); // Output: rintl(5.900000) = 6.0

    return 0;
}
```

## 15.3 运行结果
![](rint.png)

# 参考
1. [\[API Reference Document\]](https://www.apiref.com/c-zh/t_alph_r.htm)
