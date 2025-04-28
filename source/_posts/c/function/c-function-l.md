---
title: C语言函数大全--l开头的函数
date: 2023-04-17 23:31:45
updated: 2025-04-27 19:50:15
categories:
  - 开发语言-C
tags:
  - C语言函数大全
  - l开头的函数
---

![](/images/cplus-logo.png)


# 总览
| 函数声明 |  函数功能  |
|:--|:--|
| `long labs(long n);`| 计算长整型的绝对值  |
| `long long int llabs(long long int n);`| 计算long long int 类型整数的绝对值  |
| `double ldexp(double x, int exp);`|  计算 x 乘以 2 的指定次幂（double） |
| `float ldexpf(float x, int exp);`|  计算 x 乘以 2 的指定次幂（float） |
| `long double ldexpl(long double x, int exp);`|  计算 x 乘以 2 的指定次幂（long double） |
|`ldiv_t ldiv(long int numer, long int denom);` |  计算两个 long int 类型整数的商和余数 |
|`lldiv_t lldiv(long long int numer, long long int denom);`|计算两个 long long int 类型整数的商和余数|
|`void *lfind(const void *key, const void *base, size_t nmemb, size_t size, int (*compar)(const void *, const void *));` | 它是标准 C 库函数 `<search.h>` 中的一个函数，用于在一个数组中查找指定元素。  |
|`void line( int x1, int y1, int x2, int y2);` |  在指定两点间画一直线 |
| `void linerel(int dx, int dy);`| 从当前位置绘制一条指定长度和方向的线段。   |
| `void lineto(int x, int y);`|  从当前位置绘制一条直线到指定位置 |
|`struct tm *localtime(const time_t *timep);` | 将 UNIX 时间戳转换为本地时间  |
|`int lock(int fd, int cmd, off_t len);` | 它是标准 C 库函数 `<fcntl.h>` 中的一个函数，用于对文件进行加锁或解锁操作   |
| `double log(double x);`| 计算自然对数  |
|`double log10(double x);` |  计算以 10 为底的对数 |
|`void longjmp(jmp_buf env, int val);` |  跳转到指定的程序位置并恢复相应的上下文环境 |
|`void lowvideo(void);` | 用于将文本颜色设置为低对比度模式  |
|`unsigned long _lrotl(unsigned long value, int shift);` | 它是 Windows 系统特有的函数，用于将 32 位无符号整数按位循环左移。  |
| `void *lsearch(const void *key, void *base, size_t *nelp, size_t width, int (*compar)(const void *, const void *));`|  用于在指定的数组中查找指定元素，并返回该元素在数组中的地址 |
| ``long lseek(int handle, long offset, int fromwhere);``|  设置文件操作指针，即改变文件读取或写入的位置 |
| `char *ltoa(long value, char *str, int radix);`|  用于将长整型数值转换为字符串格式 |
| `char *lltoa(long long value, char *str, int radix);`| 用于将长长整型数值转换为字符串格式  |


# 1. labs，llabs
## 1.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `long labs(long n);`| 计算长整型的绝对值  |
| `long long int llabs(long long int n);`| 计算long long int 类型整数的绝对值  |

## 1.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>

int main(void)
{
    long result;
    long x = -12345678L;
    result= labs(x);
    printf("number: %ld , abs value: %ld\n", x, result);

    long long  resultL;
    long long int xL = -1234567890123456789;
    resultL = llabs(xL);
    printf("The absolute value of %lld is %lld\n", xL, resultL);

    return 0;
}
```
## 1.3 运行结果
![](labs.png)


# 2. ldexp，ldexpf，ldexpl
## 2.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double ldexp(double x, int exp);`|  计算 x 乘以 2 的指定次幂（double） |
| `float ldexpf(float x, int exp);`|  计算 x 乘以 2 的指定次幂（float） |
| `long double ldexpl(long double x, int exp);`|  计算 x 乘以 2 的指定次幂（long double） |

## 2.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main()
{
    int n = 3;
    double x = 3.5, result;
    result = ldexp(x, n);

    float xf = 3.5f, resultf;
    resultf = ldexpf(xf, n);

    long double xL = 3.5L, resultL;
    resultL = ldexpl(xL, n); 

    printf("ldexp(%lf, %d) = %lf\n", x, n, result);
    printf("ldexpf(%f, %d) = %f\n", xf, n, resultf);
    printf("ldexpl(%Lf, %d) = %Lf\n", xL, n, resultL);

    return 0;
}
```

> **注意：ldexp，ldexpf，ldexpl** 函数会对参数进行溢出和下溢处理，因此可以处理很大或很小的数值。

## 2.3 运行结果
![](ldexp.png)


# 3. ldiv，lldiv
## 3.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`ldiv_t ldiv(long int numer, long int denom);` |  计算两个 long int 类型整数的商和余数 |
|`lldiv_t lldiv(long long int numer, long long int denom);`|计算两个 long long int 类型整数的商和余数|


**参数：**
- **numer ：** 被除数
- **denom ：** 除数


`ldiv` 函数的返回值类型 `ldiv_t` 是一个结构体类型，定义如下：

```c
typedef struct {
    long int quot;  // 商
    long int rem;   // 余数
} ldiv_t;
```

`lldiv` 函数的返回值类型 `lldiv_t` 是一个结构体类型，定义如下：

```c
typedef struct {
    long long int quot;  // 商
    long long int rem;   // 余数
} lldiv_t;
```

## 3.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>

int main()
{
    long int numer = 1234567890;
    long int denom = 987654321;
    ldiv_t result;

    result = ldiv(numer, denom);

    printf("%ld / %ld = %ld, %ld %% %ld = %ld\n", numer, denom, 
        result.quot, numer, denom, result.rem);

    long long int numerL = 1234567890123456789LL;
    long long int denomL = 987654321LL;
    lldiv_t resultL;

    resultL = lldiv(numerL, denomL);

    printf("%lld / %lld = %lld, %lld %% %lld = %lld\n", numerL, denomL,
        resultL.quot, numerL, denomL, resultL.rem);

    return 0;
}
```

> **注意：** 如果 **denom** 参数为零，则 `ldiv()` 函数会产生一个异常情况。此外，如果两个参数中有一个或两个都是负数，则商和余数的计算规则将根据 **C** 标准进行调整。

## 3.3 运行结果
![](ldiv.png)

# 4. lfind
## 4.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void *lfind(const void *key, const void *base, size_t nmemb, size_t size, int (*compar)(const void *, const void *));` | 它是标准 C 库函数 `<search.h>` 中的一个函数，用于在一个数组中查找指定元素。  |

**参数：**
- **key ：** 要查找的元素
- **base ：** 要查找的数组的首地址
- **nmemb ：** 数组元素个数
- **size ：** 每个数组元素的大小（以字节为单位）
- **compar ：** 比较函数，用于比较数组元素和要查找的元素。**compar** 函数需要返回一个整数值，表示两个元素之间的关系：
   - 如果第一个元素小于第二个元素，则返回一个负数。
   - 如果第一个元素等于第二个元素，则返回零。
   - 如果第一个元素大于第二个元素，则返回一个正数。


## 4.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <search.h>

int compare(const void *a, const void *b)
{
    return (*(int*)a - *(int*)b);
}

int main()
{
    int arr[] = {1, 2, 3, 4, 5};
    int n = sizeof(arr) / sizeof(arr[0]);
    for (int i = 0; i < n; i++)
        printf("%d ", arr[i]);

    unsigned int * number = (unsigned int *)&n;
    int key = 2;
    int *result;

    result = (int *)lfind(&key, arr, number, sizeof(int), compare);

    if (result != NULL) {
        printf("\n%d is found at index %lld\n", key, result - arr);
    } else {
        printf("\n%d is not found in the array\n", key);
    }

    return 0;
}

```

> **注意：** `lfind()` 函数使用线性搜索算法，因此对于大规模数据可能不太适用。除此之外，该函数还有一些变种函数，例如 `bsearch()` 和 `tfind()` 等，也可以用于在数组或树结构中查找元素。

## 4.3 运行结果
![](lfind.png)

# 5. line
## 5.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void line( int x1, int y1, int x2, int y2);` |  在指定两点间画一直线 |
**参数：**
> **(x1, y1)  ：** 第一个点的坐标
>  **(x2, y3)  ：** 第二个点的坐标

## 5.2 演示示例
```c
#include <graphics.h>

int main(void)
{
    int gdriver = DETECT, gmode;
    int xmax, ymax;
    
    initgraph(&gdriver, &gmode, "");
    
    setcolor(getmaxcolor());
    xmax = getmaxx();
    ymax = getmaxy();

    // 在(0,0) 和(xmax, ymax)之间画一直线
    line(0, 0, xmax, ymax);

    // 在(0,ymax) 和(xmax, 0)之间画一直线
    line(0, ymax, xmax, 0);

    /* clean up */
    getch();
    closegraph();
    return 0;
}
```

## 5.3 运行结果
![](line.png)

# 6. linerel
## 6.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `void linerel(int dx, int dy);`| 从当前位置绘制一条指定长度和方向的线段。   |
**参数：**
- **dx ：** 线段在 X 轴上的位移量
- **dy ：** 线段在 Y 轴上的位移量

## 6.2 演示示例
```c
#include <graphics.h>

int main()
{
    int gd = DETECT, gm;
    initgraph(&gd, &gm, "");

    moveto(100, 100);  // 将当前点移动到 (100, 100)
    linerel(50, 0);    // 绘制长度为 50，方向为水平（X 轴正方向）的线段
    linerel(0, 50);    // 绘制长度为 50，方向为垂直（Y 轴正方向）的线段
    linerel(-50, 0);    // 绘制长度为 50，方向为水平（X 轴反方向）的线段
    linerel(0, -50);    // 绘制长度为 50，方向为垂直（Y 轴反方向）的线段

    getch();
    closegraph();

    return 0;
}
```

## 6.3 运行结果
![](linerel.png)

# 7. lineto
## 7.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `void lineto(int x, int y);`|  从当前位置绘制一条直线到指定位置 |
**参数：**
> **x ：** 线段终点的 X 坐标
> **y ：** 线段终点的 Y 坐标

## 7.2 演示示例
```c
#include <graphics.h>

int main()
{
    int gd = DETECT, gm;
    initgraph(&gd, &gm, "");

    moveto(100, 100);  // 将当前点移动到 (100, 100)
    lineto(150, 150);  // 绘制一条从当前点到 (150, 150) 的线段
    lineto(200, 100);  // 绘制一条从当前点到 (200, 100) 的线段
    lineto(150, 50);   // 绘制一条从当前点到 (150, 50) 的线段
    lineto(100, 100);   // 绘制一条从当前点到 (150, 50) 的线段

    getch();
    closegraph();

    return 0;
}
```

## 7.3 运行结果
![](lineto.png)


# 8. localtime
## 8.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`struct tm *localtime(const time_t *timep);` | 将 UNIX 时间戳转换为本地时间  |
**参数：**
- **timep ：** 指向 time_t 类型的指针，表示要转换的 UNIX 时间戳

**返回值：**
- **struct tm * ：** 一个指向 struct tm 类型的指针，该结构体包含了表示本地时间的各个字段，例如年、月、日、时、分、秒等。

## 8.2 演示示例
```c
#include <stdio.h>
#include <time.h>

int main()
{
    time_t now;
    struct tm *local;

    now = time(NULL);       // 获取当前时间戳
    local = localtime(&now);  // 将当前时间戳转换为本地时间

    printf("Current date and time: %s\n", asctime(local));

    return 0;
}
```
> **注意：** 在使用 `localtime()` 函数时需要注意结构体中的字段值是否正确，例如月份、星期等的表示方式可能因不同系统而异。

## 8.3 运行结果
![](localtime.png)


# 9. lock
## 9.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int lock(int fd, int cmd, off_t len);` | 它是标准 C 库函数 `<fcntl.h>` 中的一个函数，用于对文件进行加锁或解锁操作   |
**参数：**
- **fd ：** 文件描述符
- **cmd ：** 要执行的加锁或解锁操作（例如 F_LOCK、F_ULOCK 等）
- **len ：** 要锁定的字节数。

**返回值：**
> 返回一个整数值表示操作是否成功，若成功则返回 0，否则返回 -1

## 9.2 演示示例
```c
#include <stdio.h>
#include <fcntl.h>

int main()
{
    int fd;
    char buf[128];
    int nbytes;

    fd = open("test.txt", O_RDWR | O_CREAT, 0666);
    if (fd == -1) {
        perror("open");
        return 1;
    }

    // 加锁
    if (lock(fd, F_LOCK, 0) == -1) {
        perror("lock");
        return 1;
    }

    // 写入数据
    sprintf(buf, "Hello, world!\n");
    nbytes = write(fd, buf, sizeof(buf));
    if (nbytes == -1) {
        perror("write");
        return 1;
    }

    // 解锁
    if (lock(fd, F_ULOCK, 0) == -1) {
        perror("unlock");
        return 1;
    }

    close(fd);

    return 0;
}

```

上述示例程序中，首先通过 `open()` 函数打开一个名为 `test.txt` 的文件，并设置文件访问模式为可读写。接着，调用 `lock()` 函数对该文件进行加锁操作，保护写入数据的过程。然后，通过 `write()` 函数将数据写入到文件中。最后，调用 `lock()` 函数对该文件进行解锁操作，释放锁定的资源。

> **注意：** 在使用 `lock()` 函数时需要注意加锁和解锁的顺序、范围等问题，否则可能会造成死锁或其他问题。此外，该函数只适用于文件系统，不能用于套接字等其他类型的文件描述符。

# 10. log
## 10.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double log(double x);`| 计算自然对数  |
**参数：**
- **x ：** 要计算自然对数的数字。

**返回值：** 
- **x** 的自然对数，即 `ln(x)`。

## 10.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main()
{
    double x = 2.0;
    double result = log(x);

    printf("The natural logarithm of %lf is %lf.\n", x, result);

    return 0;
}

```

> **注意：** 由于 `log()` 函数接受的参数和返回值都是 `double` 类型，因此在使用时需要保证传入的参数类型正确，避免发生精度损失等问题。同时， `log()` 函数的参数不能为负数或零，否则会产生不可预知的行为。

## 10.3 运行结果
![](log.png)


# 11. log10
## 11.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double log10(double x);` |  计算以 10 为底的对数 |
**参数：**
- **x ：** 要计算以 10 为底的对数的数字

**返回值：** 
- **x** 的以 **10** 为底的对数，即 `log10(x)`

## 11.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main()
{
    double x = 100.0;
    double result = log10(x);

    printf("The logarithm base 10 of %lf is %lf.\n", x, result);

    return 0;
}
```

## 11.3 运行结果
![](log10.png)

# 12. longjmp
## 12.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void longjmp(jmp_buf env, int val);` |  跳转到指定的程序位置并恢复相应的上下文环境 |
**参数：**
- **env ：** 保存上下文环境的缓冲区
- **val ：** 跳转时返回的值

> 注意：在使用 `longjmp()` 函数之前，需要先调用 `setjmp()` 函数来设置上下文环境，并将其保存在 `jmp_buf` 数据类型中。然后，在程序执行过程中，如果需要跳转到之前设定的位置，就可以使用 `longjmp()` 函数进行跳转和上下文恢复。

## 12.2 演示示例
```c
#include <stdio.h>
#include <setjmp.h>

jmp_buf buf;

void do_something()
{
    printf("do_something() start.\n");

    // 跳转到 setjmp() 处
    longjmp(buf, 1);

    printf("do_something() end.\n");
}

int main()
{
    int val = 0;

    // 设置上下文环境
    if (setjmp(buf) == 0) {
        printf("setjmp() called.\n");
        do_something();
    } else {
        printf("longjmp() called.\n");
        val = 1;
    }

    printf("Program ends with value %d.\n", val);

    return 0;
}
```

上述示例程序中，首先在主函数中调用 `setjmp()` 函数设置上下文环境，并将其保存在 **buf** 变量中。然后，程序调用 `do_something()` 函数，在其中调用 `longjmp()` 函数跳转到之前设定的位置，并返回值为 **1**。

由于 `longjmp()` 调用后不会返回到调用它的位置，因此 `do_something()` 函数在被调用后并未执行完毕，而是直接跳转到了 `setjmp()` 所在的位置。当程序回到 `setjmp()` 处时，检测到了从 `longjmp()` 跳转过来的信号，并返回值为 **1**，表示跳转成功。

最后，程序输出 `"Program ends with value 1."`，结束运行。

> **注意：** 使用 `longjmp()` 和 `setjmp()` 函数进行跳转时，必须保证跳转的目标位置和之前设置的上下文环境是兼容的，否则可能会导致程序崩溃或其他严重问题。同时，尽管 `longjmp()` 可以快速跳出当前函数或代码块，但在实际应用中应该尽量避免使用它，以免造成代码逻辑混乱和难以调试的问题。

## 12.3 运行结果
![](longjmp.png)


# 13. lowvideo
## 13.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void lowvideo(void);` | 用于将文本颜色设置为低对比度模式  |

## 13.2 演示示例
```c
#include <conio.h>

int main(void)
{
    clrscr(); // 清空屏幕
    highvideo(); // 将文本颜色设置为高对比度模式
    cprintf("High Intesity Text\r\n");
    lowvideo(); // 将文本颜色设置为低对比度模式
    gotoxy(1,2); // 将光标移动到指定的坐标 (x, y)，其中 x 和 y 分别为列和行数
    cprintf("Low Intensity Text\r\n");

    return 0;
}
```
当该程序运行时，首先清空了控制台屏幕，然后将文本颜色设置为高对比度模式并输出一段文本。接着，将文本颜色设置为低对比度模式，并将光标移动到第二行第一个字符位置，输出另外一段文本。最后，程序执行结束，并返回 0。

> **注意：** `<conio.h>` 头文件中的函数在不同的操作系统和编译器下可能会有所不同，并且并非所有的平台都支持低对比度文本模式。在实际应用中，应该避免过度使用低对比度模式，以免影响用户体验和可读性。


# 14. _lrotl
## 14.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`unsigned long _lrotl(unsigned long value, int shift);` | 它是 Windows 系统特有的函数，用于将 32 位无符号整数按位循环左移。  |
**参数：**
- **value ：** 要进行循环左移的 32 位无符号整数
- **shift ：** 左移的位数

**返回值：**
- 左移后的结果。

## 14.2 演示示例
```c
#include <stdlib.h>
#include <stdio.h>

int main()
{
   unsigned long result;
   unsigned long value = 2;

   result = _lrotl(value, 2);
   printf("The value %lu rotated left one bit is: %lu\n", value, result);

   return 0;
}
```

> **注意：** `_lrotl()` 函数是 **Windows** 系统特有的函数，在其他操作系统或编译器下可能不可用或使用方式有所不同。此外，由于该函数只适用于 **32** 位无符号整数，如果需要对 **64** 位整数进行位移操作，则需要使用其他函数。

## 14.3 运行结果
![](lrotl.png)

# 15. lsearch
## 15.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `void *lsearch(const void *key, void *base, size_t *nelp, size_t width, int (*compar)(const void *, const void *));`|  用于在指定的数组中查找指定元素，并返回该元素在数组中的地址 |
**参数：**
- **key ：** 要查找的元素指针
- **base ：** 要进行查找的数组首地址
- **nelp ：** 当前数组中元素的个数
- **width ：** 数组中每个元素所占用的字节数
- **compar ：** 比较函数指针，用于比较两个元素的大小关系。

## 15.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <search.h>

int compare(const void *a, const void *b)
{
    return *(int *)a - *(int *)b;
}

int main()
{
    int arr[] = { 3, 1, 4, 1, 5, 9, 2, 6, 5 };
    int len = sizeof(arr) / sizeof(int);

    int key = 5;
    int *result = lsearch(&key, arr, &len, sizeof(int), compare);

    if (result != NULL) {
        printf("Found %d at index %ld.\n", *result, result - arr);
    } else {
        printf("%d not found in the array.\n", key);
    }

    return 0;
}

```

上述示例程序中，首先定义了一个整型数组 **arr** 并初始化为 `{ 3, 1, 4, 1, 5, 9, 2, 6, 5 }`。然后，将要查找的元素值 **key** 设置为 **5**，并调用 `lsearch()` 函数在数组中查找该元素。如果找到了，则输出该元素在数组中的下标；否则输出未找到的提示。

> **注意：** `lsearch()` 函数在查找数组元素时，只能够找到第一个匹配的元素，并返回其地址。如果数组中存在多个相同的元素，则无法区分它们的位置。此外，使用 `lsearch()` 函数进行查找时，必须保证数组已经按照指定的比较函数从小到大排好序，否则可能会导致查找失败或找到错误的元素。

## 15.3 运行结果
![](lsearch.png)


# 16. lseek
## 16.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| ``long lseek(int handle, long offset, int fromwhere);``|  设置文件操作指针，即改变文件读取或写入的位置 |
**参数：**
- **handle ：** 文件描述符
- **offset ：** 偏移量
- **whence ：** 偏移量的参考位置
>    - **SEEK_SET：** 从文件开头开始计算偏移量（即绝对位置）
>    - **SEEK_CUR：** 从当前位置开始计算偏移量（即相对位置）
>    - **SEEK_END：** 从文件结尾开始计算偏移量（即反向偏移）

**返回值：**
- 如果成功，则返回新的文件指针位置（即距离文件开头的字节数）；
>  - 如果发生错误，则返回 -1。

## 16.2 演示示例
```c
#include <stdio.h>
#include <unistd.h>
#include <fcntl.h>

int main()
{
    int fd = open("temp.txt", O_RDWR);
    if (fd == -1) {
        printf("Failed to open the file.\n");
        return -1;
    }

    off_t pos = lseek(fd, 5, SEEK_SET);
    if (pos == -1) {
        printf("Failed to seek the file.\n");
        close(fd);
        return -1;
    }

    char buf[10];
    ssize_t nread = read(fd, buf, 5);
    if (nread == -1) {
        printf("Failed to read the file.\n");
        close(fd);
        return -1;
    }

    buf[nread] = '\0';
    printf("Read %ld bytes from position %ld: %s\n", nread, pos, buf);

    close(fd);
    return 0;
}
```

上述示例程序中，首先使用 `open()` 函数打开名为 `"temp.txt"` 的文件，并获取其文件描述符。然后，调用 `lseek()` 函数将文件指针移动到距离文件开头 **5** 个字节处。接着，调用 `read()` 函数从该位置开始读取 **5** 个字节的数据，并输出读取结果。

> **注意：** `lseek()` 函数只能够对可寻址的文件进行操作，如磁盘文件、终端设备等，而不能对无法随机访问的流式数据进行操作，如管道、套接字等。同时，在使用 `lseek()` 函数时应该注意文件操作模式和文件共享模式，以免影响其他进程或线程的文件访问。

## 16.3 运行结果
![](lseek-1.png)

![](lseek-2.png)


# 17. ltoa
## 17.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `char *ltoa(long value, char *str, int radix);`|  用于将长整型数值转换为字符串格式 |
**参数：**
- **value ：** 要转换的长整型数值
- **str ：** 保存转换结果的字符缓冲区指针
- **radix ：** 要转换的进制数（如2进制、 10 进制、16 进制等），取值范围为 2～36

**返回值：**
- 指向转换结果的指针（即 str 参数的值）

## 17.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>

int main()
{
    long value = 123456789L;
    char str[20];

    ltoa(value, str, 2);
    printf("The result of converting %ld to binary string is: %s\n", value, str);

    ltoa(value, str, 10);
    printf("The result of converting %ld to decimal string is: %s\n", value, str);

    ltoa(value, str, 16);
    printf("The result of converting %ld to hexadecimal string is: %s\n", value, str);

    return 0;
}

```

> **注意：** `ltoa()` 函数在将长整型数值转换为字符串时，会将负数转换为相应的带符号字符串。如果要对无符号长整型进行转换，则需要使用其他函数或技巧。此外，由于 `ltoa()` 函数没有对输出缓冲区溢出进行检查，因此在使用时应该确保缓冲区足够大，以免发生错误。

## 17.3 运行结果

![](ltoa.png)


# 18. lltoa
## 18.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `char *lltoa(long long value, char *str, int radix);`| 用于将长长整型数值转换为字符串格式  |
**参数：**
- **value ：** 要转换的长长整型数值
- **str ：** 保存转换结果的字符缓冲区指针
- **radix ：** 要转换的进制数（如 2进制、10 进制、16 进制等），取值范围为 2～36

**返回值：**
- 指向转换结果的指针（即 str 参数的值）

## 18.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>

int main()
{
    long long value = 123456789012345LL;
    char str[20];

    lltoa(value, str, 2);
    printf("The result of converting %lld to binary string is: %s\n", value, str);

    lltoa(value, str, 10);
    printf("The result of converting %lld to decimal string is: %s\n", value, str);

    lltoa(value, str, 16);
    printf("The result of converting %lld to hexadecimal string is: %s\n", value, str);

    return 0;
}
```

## 18.3 运行结果
![](lltoa.png)


# 参考
1. [\[API Reference Document\]](https://www.apiref.com/c-zh/t_alph_l.htm)
