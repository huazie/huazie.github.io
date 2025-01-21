---
title: C语言函数大全--d开头的函数
date: 2023-03-28 00:24:20
updated: 2024-06-08 23:27:10
categories:
  - 开发语言-C
tags:
  - C语言函数大全
  - d开头的函数
---

![](/images/cplus-logo.png)

# 总览

| 函数声明 |  函数功能  |
|:--|:--|
|`void detectgraph(int *graphdriver, int *graphmode);` |  通过检测硬件确定图形驱动程序和模式 |
|`double difftime(time_t time2, time_t time1);`| 计算两个时刻之间的时间差  |
|`void disable(void);` | 屏蔽中断  |
|`div_t div(int number, int denom);` | 将两个整数相除, 返回商和余数  |
|`void drawpoly(int numpoints, int *polypoints);`| 画多边形  |
|`int dup(int handle);`| 复制文件描述符；若成功为新的文件描述，若出错为-1 |
|`int dup2(int oldhandle, int newhandle);` | 复制文件描述符；若成功为新的文件描述，若出错为-1。  |

# 1. detectgraph
## 1.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void detectgraph(int *graphdriver, int *graphmode);` |  通过检测硬件确定图形驱动程序和模式 |

## 1.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>
#include <conio.h>

/* names of the various cards supported */
char *dname[] = {
    "requests detection",
    "a CGA",
    "an MCGA",
    "an EGA",
    "a 64K EGA",
    "a monochrome EGA",
    "an IBM 8514",
    "a Hercules monochrome",
    "an AT&T 6300 PC",
    "a VGA",
    "an IBM 3270 PC"
};

int main(void)
{
    int gdriver, gmode, errorcode;

    detectgraph(&gdriver, &gmode);

    initgraph(&gdriver, &gmode, "");

    errorcode = graphresult();

    if (errorcode != grOk)
    {
        printf("Graphics error: %s\n", grapherrormsg(errorcode));
        printf("Press any key to halt:");
        getch();
        exit(1);
    }

    printf("You have [%s] video display card.\n", dname[gdriver]);
    printf("Press any key to halt:");
    getch();
    return 0;
}

```
## 1.3 运行结果
![](detectgraph.gif)

# 2. difftime
## 2.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double difftime(time_t time2, time_t time1);`| 计算两个时刻之间的时间差  |

## 2.2 演示示例
```c
#include <stdio.h>
#include <time.h>

int main(void)
{
    time_t first, second; // time_t 相当于 long
    first = time(NULL);  // Gets system time 
    getchar();
    second = time(NULL); // Gets system time again 
    printf("The difference is: %lf seconds\n", difftime(second, first));
    return 0;
}
```

在上述的示例代码中，
- 首先，我们定义了两个 `time_t` 类型的变量 `first` 和 `second`；
- 然后，调用 `time(NULL)` 函数获取当前的系统时间，并赋值给 `first`；
- 接着，调用 `getchar()` 函数等待用户输入，模拟延时的功能；
- 再然后，继续调用 `time(NULL)` 函数获取当前的系统时间，并赋值给 `second`；
- 再接着，调用 `difftime()` 函数计算 `first` 和 `second` 之间的时间差【单位：秒】
- 最终，输出时间差，并结束程序。

## 2.3 运行结果
![](difftime.gif)

# 3. disable
## 3.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|` void disable(void);` | 屏蔽中断  |

## 3.2 演示示例

```c
// 中断服务示例
#include <stdio.h>
#include <dos.h>
#include <conio.h>
// 定义一个宏INTR，代表时钟中断的十六进制向量值
#define INTR 0X1C  
// 声明一个函数指针oldhandler，该指针指向一个没有参数且返回void的函数。
// 这种函数通常用于中断处理程序
void interrupt (*oldhandler)(void);

int count=0;

void interrupt handler(void)
{
    // 禁用其他中断，确保此中断处理程序执行完毕前不再响应其他中断
    disable();
    count++;
    // 启用其他中断
    enable(); 
    // 调用原先的中断处理程序
    (*oldhandler)();
}

int main(void)
{
    // 获取时钟中断的原始处理程序，并将其存储在oldhandler指向的函数中
    oldhandler = getvect(INTR);
    // 将时钟中断的处理程序设置为handler函数
    setvect(INTR, handler);

    while (count < 20)
        printf("count is %d\n",count);
    // 恢复时钟中断的原始处理程序
    setvect(INTR, oldhandler);

    return 0;
}
```

**注意：** 这个程序可能无法在现代操作系统上直接运行，因为其中的一些函数（如`disable()`、`enable()`、`getvect()` 和 `setvect()`）是特定于 **DOS** 的。如果你想在现代操作系统（如 **Linux** 或 **Windows**）上运行这个程序，你可能需要使用更现代的方法来处理中断或使用 **DOS** 模拟器。

# 4. div
## 4.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`div_t div(int number, int denom);` | 将两个整数相除, 返回商和余数  |

## 4.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main(void)
{
    div_t x = div(10,3);
    // 商 和 余数
    printf("10 div 3 = %d remainder %d\n", x.quot, x.rem);
    return 0;
}
```
## 4.3 运行结果
![](div.png)

# 5. drawpoly
## 5.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `void drawpoly(int numpoints, int *polypoints);`| 画多边形  |

参数说明：
- `numpoints`：多边形顶点的数量

- `polypoints`：一个整数数组，包含多边形的各个顶点的坐标。

## 5.2 演示示例
```c
#include <graphics.h>
#include <stdio.h>

int main(void)
{
    // request auto detection
    int gdriver = DETECT, gmode, errorcode;
    int maxx, maxy;

    // our polygon array
    int poly[10];

    // initialize graphics and local variables
    initgraph(&gdriver, &gmode, "");

    // read result of initialization
    errorcode = graphresult();
    if (errorcode != grOk) // an error occurred
    {
        printf("Graphics error: %s\n", grapherrormsg(errorcode));
        printf("Press any key to halt:");
        getch();
        // terminate with an error code
        exit(1);
    }

    maxx = getmaxx();
    maxy = getmaxy();

    poly[0] = 20;
    poly[1] = maxy / 2;

    poly[2] = maxx - 20;
    poly[3] = 20;

    poly[4] = maxx - 50;
    poly[5] = maxy - 20;

    poly[6] = maxx / 2;
    poly[7] = maxy / 2;

    // drawpoly doesn't automatically close the polygon, so we close it.

    poly[8] = poly[0];
    poly[9] = poly[1];

    // draw the polygon
    drawpoly(5, poly);

    // clean up
    getch();
    closegraph();
    return 0;
}

```
## 5.3 运行结果
![](drawpoly.png)

# 6. dup
## 6.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `int dup(int handle);`| 复制文件描述符；若成功为新的文件描述，若出错为-1 |

**dup** 返回的新文件描述符一定是当前可用文件描述中的最小数值。

## 6.2 演示示例
```c
#include <sys/types.h>
#include <sys/stat.h>
#include <string.h>
#include <stdio.h>
#include <conio.h>
#include <fcntl.h>
#include <io.h>

void flush(FILE *stream);

int main(void)
{
    FILE *fp;
    char msg[] = "This is a test";

    fp = fopen("STU.FIL", "w");

    fwrite(msg, strlen(msg), 1, fp);

    int handle;

    handle = open("temp.txt", _O_RDWR | _O_CREAT, _S_IREAD | _S_IWRITE);
    printf("file hanlde : %d\n", handle);

    printf("Press any key to flush STU.FIL:");
    getchar();

    flush(fp);

    printf("\nFile was flushed, Press any key to quit:");
    getchar();
    return 0;
}

void flush(FILE *stream)
{
    int duphandle;

    fflush(stream);

    duphandle = dup(fileno(stream));

    printf("duplicate file hanlde : %d", duphandle);

    close(duphandle);
}
```

上述代码可以简单总结如下：
1. 首先，它打开一个名为`"STU.FIL"`的文件，并以写入模式打开。
2. 然后，将字符串`"This is a test"`写入该文件。
3. 接下来，它打开一个名为 `"temp.txt"` 的文件，并获取其文件句柄。
4. 然后，提示用户按下任意键以刷新 `"STU.FIL"` 文件。
5. 接着，调用自定义的`flush`函数来刷新文件缓冲区。
   - 首先调用`fflush`函数来刷新传入的文件流的缓冲区;
   - 然后，使用`dup`函数复制文件描述符，并将其存储在`duphandle`变量中;
   - 接着，打印出复制的文件句柄；
   - 最后，关闭复制的文件句柄。
6. 最后，再次提示用户按下任意键以退出程序。

## 6.3 运行结果
![](dup.png)

# 7. dup2
## 7.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int dup2(int oldhandle, int newhandle);` | 复制文件描述符；若成功为新的文件描述，若出错为-1。  |

**dup2** 可以用 **newhandle** 参数指定新的描述符数值。如果 **newhandle** 已经打开，则先关闭。若 **oldhandle = newhandle**，则 **dup2** 返回 **newhandle**，而不关闭它。

## 7.2 演示示例
```c
#include <sys\stat.h>
#include <stdio.h>
#include <string.h>
#include <fcntl.h>
#include <io.h>

int main(void)
{
    #define STDOUT 1

    int handle, oldstdout;
    char msg[] = "This is a test1";

    handle = open("STU.FIL", O_CREAT | O_RDWR, S_IREAD | S_IWRITE);
    printf("open file handle : %d\n", handle);

    oldstdout = dup(STDOUT);
    printf("dup file handle : %d", oldstdout);

    dup2(handle, STDOUT);

    close(handle);

    write(STDOUT, msg, strlen(msg));

    dup2(oldstdout, STDOUT);

    close(oldstdout);

    return 0;
}
```
上述代码简单分析如下：

- 定义常量 `STDOUT`，其值为 1，表示标准输出的文件描述符；
- 定义整型变量 `handle` 和 `oldstdout`，以及字符数组 `msg`，用于存储要写入文件的字符串；
- 使用 `open` 函数打开名为 `"STU.FIL"` 的文件，以创建和读写模式（`O_CREAT | O_RDWR`）打开，并设置文件权限为可读可写（`S_IREAD | S_IWRITE`）；将返回的文件描述符赋值给 `handle`，并打印出来；
- 使用 `dup` 函数备份当前的标准输出（`STDOUT`），将备份的文件描述符赋值给 `oldstdout`，并打印出来；
- 使用 `dup2` 函数将标准输出重定向到 `handle` 指向的文件，即将后续的输出内容写入到 `"STU.FIL"` 文件中；
- 关闭 `handle` 指向的文件句柄；
- 使用 `write` 函数将 `msg` 字符串写入到标准输出（此时已经重定向到文件），长度为字符串的长度；
- 使用 `dup2` 函数将标准输出恢复到备份的文件描述符 `oldstdout`，即将后续的输出内容输出到屏幕上。
- 关闭 `oldstdout` 指向的文件句柄。

## 7.3 运行结果
![](dup2.png)

# 参考
1. [\[API Reference Document\]](https://www.apiref.com/c-zh/t_alph_d.htm)


