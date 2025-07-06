---
title: C语言函数大全--s 开头的函数（2）
date: 2023-04-30 09:00:00
updated: 2025-07-06 11:05:06
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
|`void setlinestyle( int linestyle, unsigned upattern, int thickness );` |  设置当前绘图窗口的线条样式、线型模式和线条宽度 |
|`void *setmem(void *dest, size_t n, int c);` |  用于将指定内存区域的每个字节都设置为指定的值 |
|`int setmode(int fd, int mode);` | 它是 `Windows` 系统下的特定库函数，用于将指定文件的 `I/O` 模式设置为文本模式或二进制模式   |
|`void setpalette(int colornum, int color);` |  设置调色板的颜色 |
|`void setrgbpalette(int colornum, int red, int green, int blue);` | 用于设置当前绘图窗口的调色板中某个颜色的 `RGB` 值  |
|`void settextjustify(int horiz, int vert);` |  用于设置文本输出的对齐方式 |
|`void settextstyle(int font, int direction, int charsize);` |  用于设置当前文本输出的字体、方向和大小 |
|`void settime(struct time *timep);` |  设置当前系统时间 |
|`void setusercharsize(int multx, int dirx, int multy, int diry);` |  用于设置用户定义的字符大小 |
|`int setvbuf(FILE *stream, char *buf, int type, unsigned size);` | 用于设置文件流的缓冲方式  |
|`void setviewport(int left, int top, int right, int bottom, int clipflag);` | 用于设置绘图窗口的视口范围  |
|`void setvisualpage(int pagenum); ` | 用于设置图形窗口中用户可见的页面  |
|`void setwritemode(int mode); ` |  用于设置绘画操作的写入模式 |
|`void (*signal(int signum, void (*handler)(int)))(int);` | 用于设置指定信号的处理方式。当系统接收到某个信号时，会调用相应的信号处理函数来处理该信号。在调用 `signal` 函数时，需要指定要处理的信号以及相应的信号处理函数。 |
|`double sin(double x);` |  用于计算一个角度（以弧度为单位）的正弦值（double） |
|`float sinf(float x);` |  用于计算一个角度（以弧度为单位）的正弦值（float） |
|`long double sinl(long double x);` |  用于计算一个角度（以弧度为单位）的正弦值（long double） |
|`void sincos(double x, double* sinVal, double* cosVal);` | 用于同时计算一个角度（以弧度为单位）的正弦值和余弦值  |
|`void sincosf(float x, float* sinVal, float* cosVal);` | 用于同时计算一个角度（以弧度为单位）的正弦值和余弦值  |
|`void sincosl(long double x, long double* sinVal, long double* cosVal);` | 用于同时计算一个角度（以弧度为单位）的正弦值和余弦值  |
|`double sinh(double x);` | 用于计算一个数的双曲正弦值  |
|`float sinhf(float x);` | 用于计算一个数的双曲正弦值  |
|`long double sinhl(long double x);` | 用于计算一个数的双曲正弦值  |

# 1. setlinestyle
## 1.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void setlinestyle( int linestyle, unsigned upattern, int thickness );` |  设置当前绘图窗口的线条样式、线型模式和线条宽度 |

**参数：**
- **linestyle ：** 线条样式，取值范围为 0 到 4，不同的值对应着不同的线条样式，详见如下表格
- **upattern ：** 线型模式，它是一个 `16` 位的无符号整数，用二进制位表示线型模式，其中 `1` 表示绘制线条，`0` 表示空白。例如，如果 `upattern` 的值是 `0x00FF`，则表示绘制一段线条，然后空白一段，重复这个过程直到结束。如果 `upattern` 的值是 `0x5555`，则表示绘制一段虚线。
- **thickness ：** 线条宽度，取值范围为 1 到 10，表示线条的像素宽度

| 线条样式 | 值 |描述 |
|--|--|--|
| SOLID_LINE | 0 | 实线|
| DOTTED_LINE |1  |虚线|
| CENTER_LINE | 2 |点线|
| DASHED_LINE | 3 |长短线|
| USERBIT_LINE | 4 |双点线|

## 1.2 演示示例
```c
#include <graphics.h>
#include <string.h>

/* the names of the line styles supported */
char *lname[] = {
   "SOLID_LINE",
   "DOTTED_LINE",
   "CENTER_LINE",
   "DASHED_LINE",
   "USERBIT_LINE"
   };

int main(void)
{
    int gdriver = DETECT, gmode;

    int style, midx, midy, userpat;
    char stylestr[40];

    initgraph(&gdriver, &gmode, "");

    midx = getmaxx() / 2;
    midy = getmaxy() / 2;

    userpat = 1;

    for (style=SOLID_LINE; style<=USERBIT_LINE; style++)
    {
        setlinestyle(style, userpat, 1);

        strcpy(stylestr, lname[style]);

        line(0, 0, midx-10, midy);

        rectangle(100, 100, getmaxx() - 100, getmaxy() - 100);

        outtextxy(midx, midy, stylestr);

        getch();
        cleardevice();
    }

    closegraph();
    return 0;
}
```

## 1.3 运行结果
![](setlinestyle.gif)

# 2. setmem
## 2.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void *setmem(void *dest, size_t n, int c);` |  用于将指定内存区域的每个字节都设置为指定的值 |

**参数：**
- **dest ：** 要设置的内存区域的指针
- **n ：**  要设置的字节数
- **c ：**  要设置的值

·**注意：** setmem() 并不是标准 C 函数，而是 POSIX 标准库函数，因此可能并不被所有平台所支持。如果您的编译器或操作系统不支持 setmem() 函数，可以使用标准 C 库函数 memset() 来代替

## 2.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main()
{
    char *str = (char *)malloc(10);  // 分配 10 字节的内存空间

    // 将 str 中的每个字节都设置为 'A'
    setmem(str, 10, 'A');

    printf("%s\n", str);

    free(str);

    return 0;
}
```
在上面的示例程序中，
- 我们首先使用 `malloc()` 函数分配了 `10` 字节的内存空间，并将其赋值给指针变量 str。
- 然后，我们使用 `setmem()` 函数将 `str` 指向的内存区域的每个字节都设置为 `'A'`。
- 最后，我们输出 `str` 的内容并使用 `free()` 函数释放了分配的内存空间。

# 3. setmode
## 3.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int setmode(int fd, int mode);` | 它是 `Windows` 系统下的特定库函数，用于将指定文件的 `I/O` 模式设置为文本模式或二进制模式   |

**参数：**
- **fd ：** 要设置模式的文件描述符，通常使用 `fileno()` 函数将文件指针转换为文件描述符
- **mode ：** 要设置的模式，它可以取以下两个值中的一个：
    - `_O_BINARY`：二进制模式
    - `_O_TEXT`：文本模式

**注意：** 在 `Windows` 系统中，文本模式和二进制模式之间的区别在于换行符的处理方式。在文本模式下，`Windows` 将回车符（`\r`）和换行符（`\n`）组成的字符序列映射为单个换行符（`\n`），因此在读取文本文件时可以正确处理换行符。在二进制模式下，`Windows` 不对回车符和换行符做任何转换，因此在读取文本文件时可能会出现问题。

## 3.2 演示示例
```c
#include <stdio.h>
#include <fcntl.h>

int main()
{
    int result;
    result = setmode(fileno(stdin), O_TEXT);
    if (result == -1)
        perror("Mode not available\n");
    else
        printf("Mode successfully switched\n");
    return 0;
}
```
在上面的示例程序中，我们首先使用 `setmode()` 函数将标准输入流的模式从二进制模式切换到文本模式；如果调用成功，则返回 `0`，否则返回 `-1`，并将错误信息存储在全局变量 `errno` 中。在程序中，我们使用 `perror()` 函数来输出错误信息。如果调用成功，则输出一条提示信息。

**注意：** `setmode()` 函数只适用于 `Windows` 系统下的 `C/C++` 编程，并且不是标准库函数，因此在跨平台编程时应该避免使用它。在 `Unix/Linux` 系统下，也可以使用 `fcntl()` 函数来设置文件描述符的模式。

## 3.3 运行结果
![](setmode.png)

# 4. setpalette
## 4.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void setpalette(int colornum, int color);` |  设置调色板的颜色 |

**参数：**
- **colornum ：** 要设置的调色板中的颜色数量
- **color ：** 是一个整数类型的值，用于表示调色板中的颜色。这个参数可以是一个 `RGB` 值，也可以是一个预定义颜色名称，例如 `RED` 或 `BLUE`。

## 4.2 演示示例
```c
#include <graphics.h>

int main()
{
    int gd = DETECT, gm;
    initgraph(&gd, &gm, "");

    // 将第 5 种颜色设置为红色
    setpalette(5, RED);

    // 绘制一条红色的直线
    setcolor(5);
    line(100, 100, 200, 100);

    getch();
    closegraph();
    return 0;
}
```

## 4.3 运行结果
![](setpalette.png)

# 5. setrgbpalette
## 5.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void setrgbpalette(int colornum, int red, int green, int blue);` | 用于设置当前绘图窗口的调色板中某个颜色的 `RGB` 值  |

**参数：**
- **colornum ：** 要修改的颜色索引，取值范围为 `0~255`。
- **red、green 和 blue** 要设置的 `RGB` 值，取值范围为 `0~255`。

## 5.2 演示示例
```c
#include <graphics.h>

int main()
{
    int gd = DETECT, gm;
    initgraph(&gd, &gm, "");

    setbkcolor(WHITE);
    cleardevice();

    // 将第 5 种颜色设置为蓝绿色
    setrgbpalette(5, 0, 128, 128);

    // 绘制一条蓝绿色的斜线
    setcolor(5);
    line(100, 100, 200, 200);

    getch();
    closegraph();
    return 0;
}
```
在上述的这个示例程序中，
- 我们首先使用 `setbkcolor()` 函数设置背景颜色为白色，然后清除原有屏幕使前面设置生效。
- 接着我们使用 `setrgbpalette()` 函数将第 `5` 种颜色设置为蓝绿色，并使用 `setcolor()` 函数将绘图颜色设为索引值 `5`（即蓝绿色）；
- 最后使用 `line()` 函数绘制了一条斜线。

## 5.3 运行结果
![](setrgbpalette.png)

# 6. settextjustify
## 6.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void settextjustify(int horiz, int vert);` |  用于设置文本输出的对齐方式 |

**参数：**
- **horiz ：** 水平对齐方式，可以取以下值：
    - `LEFT_TEXT`：左对齐
    - `CENTER_TEXT`：居中对齐
    - `RIGHT_TEXT`：右对齐
- **vert ：** 垂直对齐方式，可以取以下值：
    - `TOP_TEXT`：顶部对齐
    - `CENTER_TEXT`：居中对齐
    - `BOTTOM_TEXT`：底部对齐

## 6.2 演示示例
```c
#include <graphics.h>

int main()
{
    int gd = DETECT, gm;
    initgraph(&gd, &gm, "");

    // 设置文本输出的对齐方式为居中对齐
    settextjustify(CENTER_TEXT, CENTER_TEXT);

    // 输出一行居中对齐的文本
    outtextxy(getmaxx() / 2, getmaxy() / 2, "Hello, world!");

    getch();
    closegraph();
    return 0;
}
```

在上述的示例程序中，我们使用 `settextjustify()` 函数将文本输出的对齐方式设置为居中对齐，然后使用 `outtextxy()` 函数在窗口的中心输出一行文本。

**注意：** 在使用 `settextjustify()` 函数设置对齐方式时，必须指定正确的参数值，并且同时考虑水平和垂直方向的对齐方式，否则可能会导致文本输出位置错误。

## 6.3 运行结果
![](settextjustify.png)

# 7. settextstyle
## 7.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void settextstyle(int font, int direction, int charsize);` |  用于设置当前文本输出的字体、方向和大小 |

**参数：**
- **font ：** 要使用的字体编号，可以取以下值：
    - `DEFAULT_FONT`：默认字体
    - `TRIPLEX_FONT`：粗体三线字体
    - `SMALL_FONT`：小号字体
    - `SANS_SERIF_FONT`：无衬线字体
    - `GOTHIC_FONT`：哥特式字体
    - `SCRIPT_FONT`：手写字体 
- **direction ：** 文本输出的方向，可以取以下值：
    - `HORIZ_DIR`：水平方向（从左到右）
    - `VERT_DIR`：垂直方向（从下到上）
- **horiz ：** 水平对齐方式，可以取以下值：
    - `DEFAULT_WIDTH` 和 `DEFAULT_HEIGHT`：默认大小
    - `TRIPLEX_WIDTH` 和 `TRIPLEX_HEIGHT`：大号尺寸
    - `SMALL_WIDTH` 和 `SMALL_HEIGHT`：小号尺寸

## 7.2 演示示例
```c
#include <graphics.h>

int main()
{
    int gd = DETECT, gm;
    initgraph(&gd, &gm, "");

    // 设置文本输出的字体、方向和大小
    settextstyle(TRIPLEX_FONT, HORIZ_DIR, 4);

    // 输出一行文本
    outtextxy(100, 100, "Hello, world!");

    getch();
    closegraph();
    return 0;
}
```

在上述的示例程序中，我们使用 `settextstyle()` 函数将文本输出的字体设置为**粗体三线字体**、方向设置为**水平方向**、大小设置为**大号尺寸**，然后使用 `outtextxy()` 函数在窗口的指定位置输出一行文本。

**注意：** 在使用 `settextstyle()` 函数设置文本样式时，必须指定正确的参数值，并且根据具体需求灵活选择字体、方向和大小，否则可能会导致文本输出不符合预期。

## 7.3 运行结果
![](settextstyle.png)

# 8. settime
## 8.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void settime(struct time *timep);` |  设置当前系统时间 |
**参数：**
- **timep ：** 用于存储当前系统时间的结构体变量

## 8.2 演示示例
```c
#include <stdio.h>
#include <dos.h>

int main(void)
{
    struct  time t;

    gettime(&t);
    printf("The current minute is: %d\n", t.ti_min);
    printf("The current hour is: %d\n", t.ti_hour);
    printf("The current hundredth of a second is: %d\n", t.ti_hund);
    printf("The current second is: %d\n", t.ti_sec);

    t.ti_min++;
    settime(&t);

    return 0;
}
```
在上述的程序中，
- 我们首先定义了一个 `struct time` 类型的变量 `t`，用于存储当前系统时间。然后使用 `gettime()` 函数获取当前时间，并将小时、分钟、秒和百分之一秒等信息存储到 `t` 变量的对应成员变量中。
- 接着，程序使用 `printf()` 函数输出了当前的分钟数、小时数、百分之一秒数和秒数。这里使用了 `%d` 占位符来指定输出整数类型的值。
- 最后，程序将 t 变量的分钟数加上了 1，然后使用 settime() 函数将修改后的时间写入系统时钟中。

**注意 ：** 这个程序依赖于 `DOS` 系统提供的日期和时间相关函数，可能无法在其他操作系统或环境下运行。此外，直接修改系统时间可能会对计算机的正常运行产生影响，因此应该谨慎使用。
在现代的 `Windows` 操作系统中，`DOS` 环境已经被废弃，因此这个程序可能无法正常工作。如果要获取和修改当前系统时间，可以使用操作系统提供的相关 `API` 或系统调用接口实现。例如，在 `Windows` 平台上，可以使用 `GetSystemTime()` 和 `SetSystemTime()` 等函数来获取和设置系统时间。

# 9. setusercharsize
## 9.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void setusercharsize(int multx, int dirx, int multy, int diry);` |  用于设置用户定义的字符大小 |

**参数：**
- **multx ：** 水平放大倍数，取值为正整数
- **dirx ：** 水平方向，取值为 `1` 或 `-1`。当 `dirx` 的值为 `1` 时，字符不进行左右翻转；当 `dirx` 的值为 `-1` 时，字符进行左右翻转
- **multy ：** 垂直放大倍数，取值为正整数
- **diry ：** 垂直方向，取值为 `1` 或 `-1`。当 `diry` 的值为 `1` 时，字符不进行上下翻转；当 `diry` 的值为 `-1` 时，字符进行上下翻转。

## 9.2 演示示例
```c
#include <graphics.h>

int main()
{
    int gd = DETECT, gm;
    initgraph(&gd, &gm, "");

    // 设置字符的大小为水平方向放大 2 倍、垂直方向放大 3 倍
    setusercharsize(2, 1, 3, 1);

    // 输出一行文本
    outtextxy(100, 100, "Hello, world!");

    getch();
    closegraph();
    return 0;
}
```

在上述这个示例程序中，我们使用 `setusercharsize()` 函数将当前字符的大小设置为水平方向放大 `2` 倍、垂直方向放大 `3` 倍，然后使用 `outtextxy()` 函数在窗口的指定位置输出一行文本。

**注意：** 在使用 `setusercharsize()` 函数设置字符大小时，必须指定正确的参数值，并且考虑到水平和垂直方向的缩放倍数，否则可能会导致字符输出不符合预期。

## 9.3 运行结果
![](setusercharsize.png)

# 10. setvbuf
## 10.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int setvbuf(FILE *stream, char *buf, int type, unsigned size);` | 用于设置文件流的缓冲方式  |

**参数：**
- **stream ：** 要设置缓冲方式的文件指针，可以是标准输入流（`stdin`）、标准输出流（`stdout`）或标准错误流（`stderr`），也可以是通过 `fopen()` 函数打开的文件指针
- **buf ：** 缓冲区指针，可以是一个已经分配好的缓冲区，也可以是 `NULL`。如果 `buf` 参数为 `NULL`，则 `setvbuf()` 函数将自动为文件流分配一块缓冲区
- **type ：** 缓冲类型，可以取以下值：
    - `_IONBF`：不使用缓冲区，直接从或向设备读写数据
    - `_IOLBF`：行缓冲，每行数据结束后刷新缓冲区
    - `_IOFBF`：全缓冲，填满缓冲区后才进行读写操作
- **size：** 缓冲区大小，如果 `buf` 参数不为 `NULL`，则 `size` 参数指定了缓冲区大小；如果 `buf` 参数为 `NULL`，则 `size` 参数指定了系统为缓冲区分配的大小。`size` 参数只对全缓冲方式有效，行缓冲和无缓存方式忽略该参数

## 10.2 演示示例
```c
#include <stdio.h>

int main(void)
{
    FILE *input, *output;
    char bufr[512];

    input = fopen("file.in", "r+b");
    output = fopen("file.out", "w");

    if (setvbuf(input, bufr, _IOFBF, 512) != 0)
        printf("failed to set up buffer for input file\n");
    else
        printf("buffer set up for input file\n");

    if (setvbuf(output, NULL, _IOLBF, 132) != 0)
        printf("failed to set up buffer for output file\n");
    else
        printf("buffer set up for output file\n");

    fclose(input);
    fclose(output);
    return 0;
}
```

在上述的示例程序中，
- 我们首先定义了两个文件指针变量 `input` 和 `output`，分别表示输入文件和输出文件。
- 然后调用 `fopen()` 函数打开输入文件和输出文件，并将返回的文件指针保存到对应的变量中。
- 接着，程序使用 `setvbuf()` 函数分别为输入文件和输出文件设置缓冲方式。对于输入文件，使用 `_IOFBF` 缓冲类型和大小为 `512` 字节的缓冲区；对于输出文件，使用 `_IOLBF` 缓冲类型和大小为 `132` 字节的缓冲区（此处 `bufr` 缓冲区为空指针）。在设置完缓冲方式后，程序根据 `setvbuf()` 函数的返回值判断是否设置成功，如果返回值不为 0，则说明设置失败，否则说明设置成功，并输出相应的提示信息。
- 最后，程序使用 `fclose()` 函数关闭输入文件和输出文件。

**注意：** 在使用文件流进行读写操作时，必须在适当的时候进行缓冲区清理操作，以避免数据丢失或者读取到过期数据等问题。另外，需要根据具体需求选择合适的缓冲方式和缓冲区大小，以实现最优的性能和稳定性。

## 10.3 运行结果
![](setvbuf.png)

# 11. setviewport
## 11.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void setviewport(int left, int top, int right, int bottom, int clipflag);` | 用于设置绘图窗口的视口范围  |

**参数：**
- **left ：** 视口矩形的左上角横坐标，取值为正整数或 0
- **top：** 视口矩形的左上角纵坐标，取值为正整数或 0
- **right：** 视口矩形的右下角横坐标，取值为正整数
- **bottom：** 视口矩形的右下角纵坐标，取值为正整数
- **clipflag：** 裁剪标志，可以取以下值：
    - `CLIP_ON`：开启裁剪模式，只显示视口内的图形；
    - `CLIP_OFF`：关闭裁剪模式，显示整个画面。

## 11.2 演示示例
```c
#include <graphics.h>

#define CLIP_ON 1

int main()
{
    int gd = DETECT, gm;
    initgraph(&gd, &gm, "");

    setcolor(RED);
    rectangle(100, 100, 300, 200);

    // 将视口范围设置为矩形 (150, 150) - (250, 250)
    setviewport(150, 150, 250, 250, CLIP_ON);

    setcolor(GREEN);
    rectangle(0, 0, 400, 300);

    getch();
    closegraph();
    return 0;
}
```

在上面的示例程序中，
- 我们首先使用 `rectangle()` 函数绘制了一个红色的矩形；
- 然后使用 `setviewport()` 函数将视口范围设置为矩形 (150, 150) - (250, 250)；
- 最后使用 `rectangle()` 函数绘制了一个绿色的矩形，但这里只有在视口范围的矩形才显示出来。

**注意：** 在使用 `setviewport()` 函数设置视口范围时，必须指定正确的参数值，并考虑到裁剪模式的影响，否则可能会导致图形显示不符合预期。

## 11.3 运行结果
![](setviewport.png)

# 12. setvisualpage
## 12.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void setvisualpage(int pagenum); ` | 用于设置图形窗口中用户可见的页面  |

**参数：**
- **pagenum ：** 要设置的可视化页面页码，整数类型。

在双缓冲绘图中，我们通常会使用两个页面来绘制图像，一个是前台页面，另一个是后台页面。当我们绘制完一个完整的画面时，可以通过调用 `setactivepage()` 函数将后台页面切换到前台页面以显示出来。而 `setvisualpage()` 函数则用于设置用户当前看到的页面，它实际上是将指定的页面置于前台，从而更新屏幕上显示的内容。

## 12.2 演示示例
```c
#include <graphics.h>

int main(void)
{
   int gdriver = EGA, gmode = EGAHI;
   int x, y, ht;

   initgraph(&gdriver, &gmode, "");

   x = getmaxx() / 2;
   y = getmaxy() / 2;
   ht = textheight("W");

   setactivepage(1);

   line(0, 0, getmaxx(), getmaxy());

   settextjustify(CENTER_TEXT, CENTER_TEXT);
   outtextxy(x, y, "This is page #1:");
   outtextxy(x, y+ht, "Press any key to halt:");

   setactivepage(0);

   outtextxy(x, y, "This is page #0.");
   outtextxy(x, y+ht, "Press any key to view page #1:");
   getch();

   setvisualpage(1);

   getch();
   closegraph();
   return 0;
}
```

上述示例将在屏幕上绘制两个页面，并允许用户通过按任意键查看第二个页面。

**注意：** `setvisualpage()` 函数只能用于已经创建的图形窗口，且参数 `pagenum` 必须在 `0` 到 `getmaxpages()` 函数返回值之间。

## 12.3 运行结果
![](setvisualpage.gif)

# 13. setwritemode
## 13.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void setwritemode(int mode); ` |  用于设置绘画操作的写入模式 |

**参数：**
- **mode：** 要设置的写入模式，整数类型。常用的写入模式有以下三种：
    - `COPY_PUT`：0，复制模式（默认），新绘制的图形将完全覆盖旧图形
    - `XOR_PUT`：1，异或模式，新绘制的图形将与旧图形进行异或运算后显示。在这种模式下，如果一个像素既在新图形中出现过，也在旧图形中出现过，则它最终不会被显示出来；如果一个像素只在新图形中出现过，或者只在旧图形中出现过，则它将会被显示出来。
    - `OR_PUT`：2，或模式，新绘制的图形将与旧图形进行或运算后显示。在这种模式下，如果一个像素既在新图形中出现过，也在旧图形中出现过，则它最终会被显示出来；如果一个像素只在新图形中出现过，或者只在旧图形中出现过，则它将会被显示出来。

## 13.2 演示示例
```c
#include <graphics.h>

int main()
{
   int gdriver = DETECT, gmode;
   int xmax, ymax;

   initgraph(&gdriver, &gmode, "");

   xmax = getmaxx();
   ymax = getmaxy();

   setwritemode(XOR_PUT);

   line(0, 0, xmax, ymax);
   getch();

   line(0, 0, xmax, ymax);
   getch();

   setwritemode(COPY_PUT);

   line(0, 0, xmax, ymax);

   getch();
   closegraph();
   return 0;
}
```

**注意：** `setwritemode()` 函数只对紧随其后的绘画操作生效，它不会影响之前已经绘制的图形。因此，在更改写入模式之前，必须先完成之前的绘画操作。

## 13.3 运行结果
![](setwritemode.gif)


# 14. signal
## 14.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void (*signal(int signum, void (*handler)(int)))(int);` | 用于设置指定信号的处理方式。当系统接收到某个信号时，会调用相应的信号处理函数来处理该信号。在调用 `signal` 函数时，需要指定要处理的信号以及相应的信号处理函数。 |

**参数：**
-  **signum ：** 要设置的信号编号，整数类型。常见的信号有很多种，如：
      - `SIGINT`（中断信号）：通常由用户按下 "`Ctrl + C`" 产生，用于中断正在运行的程序。
      - `SIGALRM`（闹钟信号）：用于在指定的时间后向进程发送信号，可以用于实现定时器等功能。
      - `SIGFPE`（浮点异常信号）：在发生浮点运算异常时发送该信号。
      - `SIGSEGV`（段错误信号）：在访问非法的内存地址时发送该信号。
- **handler ：** 要设置的信号处理函数，是一个指向函数的指针，其形式为 `void handler(int)`

**返回值：**
- 如果调用成功，返回之前对信号的处理方式（通常是一个函数指针）。如果之前没有设置过该信号的处理方式，返回默认的处理方式（通常是 `SIG_DFL` 或 `SIG_IGN`）。
- 如果调用失败，返回 `SIG_ERR`。

## 14.2 演示示例
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

## 14.3 运行结果
![](signal.png)

# 15. significand，significandf
## 15.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double significand(double x);` | 用于分离浮点数 x 的尾数部分（double）  |
|`float significandf(float x);` | 用于分离浮点数 x 的尾数部分（float）  |

**参数：**
- **x ：** 要求尾数的浮点数。

## 15.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main()
{
    double x = 123.456789;
    float y = -2.5;

    printf("The significand of %lf is %lf\n", x, significand(x));
    printf("The significand of %f is %f\n", y, significandf(y));

    return 0;
}
```
**注意：** 在某些旧版本的编译器中，可能没有实现 significand 函数。这个时候可以考虑使用其他类似的函数来替代，如 `frexp`、`modf` 等。


# 16. sin，sinf，sinl
## 16.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double sin(double x);` |  用于计算一个角度（以弧度为单位）的正弦值（double） |
|`float sinf(float x);` |  用于计算一个角度（以弧度为单位）的正弦值（float） |
|`long double sinl(long double x);` |  用于计算一个角度（以弧度为单位）的正弦值（long double） |

**参数：**
- **x ：** 要求正弦值的角度，以弧度为单位

## 16.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main()
{
    double x = M_PI / 6.0; // M_PI 圆周率
    float y = M_PI / 4.0L;
    long double z = M_PI / 3.0L;

    printf("sin(%lf) = %lf\n", x, sin(x));
    printf("sinf(%f) = %f\n", y, sinf(y));
    printf("sinl(%Lf) = %Lf\n", z, sinl(z));

    return 0;
}
```

## 16.3 运行结果
![](sin.png)


# 17. sincos，sincosf，sincosl
## 17.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void sincos(double x, double* sinVal, double* cosVal);` | 用于同时计算一个角度（以弧度为单位）的正弦值和余弦值  |
|`void sincosf(float x, float* sinVal, float* cosVal);` | 用于同时计算一个角度（以弧度为单位）的正弦值和余弦值  |
|`void sincosl(long double x, long double* sinVal, long double* cosVal);` | 用于同时计算一个角度（以弧度为单位）的正弦值和余弦值  |

**参数：**
- **x ：** 要求正弦值和余弦值的角度，以弧度为单位
- **sinVal ：** 存放计算得到的正弦值的指针
- **cosVal ：** 存放计算得到的余弦值的指针

## 17.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main()
{
    double x = M_PI / 6.0;
    float y = M_PI / 4.0L;
    long double z = M_PI / 3.0L;

    double sinVal, cosVal;
    float sinfVal, cosfVal;
    long double sinlVal, coslVal;

    sincos(x, &sinVal, &cosVal);
    sincosf(y, &sinfVal, &cosfVal);
    sincosl(z, &sinlVal, &coslVal);

    printf("sincos(%lf) = (%lf, %lf)\n", x, sinVal, cosVal);
    printf("sincosf(%f) = (%f, %f)\n", y, sinfVal, cosfVal);
    printf("sincosl(%Lf) = (%Lf, %Lf)\n", z, sinlVal, coslVal);

    return 0;
}
```

## 17.3 运行结果
![](sincos.png)


# 18. sinh，sinhf，sinhl
## 18.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double sinh(double x);` | 用于计算一个数的双曲正弦值  |
|`float sinhf(float x);` | 用于计算一个数的双曲正弦值  |
|`long double sinhl(long double x);` | 用于计算一个数的双曲正弦值  |

**参数：**
- **x ：** 要求双曲正弦值的数，以弧度为单位

## 18.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main()
{
    double x = 1.0;
    float y = 2.0;
    long double z = 3.0;

    printf("sinh(%lf) = %lf\n", x, sinh(x));
    printf("sinhf(%f) = %f\n", y, sinhf(y));
    printf("sinhl(%Lf) = %Lf\n", z, sinhl(z));

    return 0;
}
```

## 18.3 运行结果
![](sinh.png)



