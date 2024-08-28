---
title: C语言函数大全--g开头的函数（上）
date: 2023-04-06 10:36:37
updated: 2024-07-22 20:19:23
categories:
  - 开发语言-C
tags:
  - C语言函数大全
  - g开头的函数
---

![](/images/cplus-logo.png)

# 总览
| 函数声明 |  函数功能  |
|:--|:--|
|`char *gcvt(double value, int ndigit, char *buf);` |  把浮点数转换成字符串，同时返回一个指向字符串的存储位置的指针的函数。 |
| `void getarccoords(struct arccoordstype *arccoords);`| 取最后一次调用arc的坐标  |
| `int getbkcolor(void);`|  获取当前背景颜色 |
|`int getc(FILE *stream);` |  从流中取字符 |
| `int getchar(void);`| 从 **stdin** 流中读字符  |
| `int getcolor(void);`| 当前画线的颜色  |
| `char *getcwd(char *buffer, int maxlen);`| 获取当前工作目录  |
|`struct palettetype* getdefaultpalette(void);` |  获取调色板定义结构 |
|`char *getdrivename(void);` |  获取当前图形驱动程序名字|
| `void getfillpattern(char *upattern);`| 将用户定义的填充模式拷贝到内存中  |
| `void getfillsettings(struct fillsettingstype *fillinfo);`|  获取有关当前填充模式和填充颜色的信息 |
|`int getgraphmode(void);` | 获取当前图形模式  |
|`void getimage(int left, int top, int right, int bottom, void *bitmap);` |  保存指定区域的屏幕上的像素图形到指定的内存区域 |
|`void getlinesettings(struct linesettingstype *lininfo);` | 取当前线型、模式和宽度  |
| `int getmaxcolor(void);`|  可以传给函数 setcolor 的最大颜色值 |
|`int getmaxx(void);` |  屏幕的最大x坐标 |
| `int getmaxy(void);`|  屏幕的最大y坐标 |

# 1. gcvt
## 1.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`char *gcvt(double value, int ndigit, char *buf);` |  把浮点数转换成字符串，同时返回一个指向字符串的存储位置的指针的函数。 |

> **参数：**
> **value：** 被转换的值。
> **ndigit：** 存储的有效数字位数。
> **buf：** 结果的存储位置。

> **注意：** **gcvt** 函数把一个浮点值转换成一个字符串 (包括一个小数点和可能的符号字节) 并存储该字符串在 **buffer** 中。该 **buffer** 应足够大以便容纳转换的值加上结尾的 结束符 `'\0'`，它是自动添加的。
> 如果一个缓冲区的大小为 **ndigit + 1**，则 **gcvt** 函数将覆盖该缓冲区的末尾。这是因为转换的字符串包括一个小数点以及可能包含符号和指数信息。
## 1.2 演示示例
```c
#include <stdlib.h>
#include <stdio.h>

int main(void)
{
   char str[25];
   double num;
   int sig = 5; 

   num = 1.23;
   gcvt(num, sig, str);
   printf("string = %s\n", str);

   num = -456.78912;
   gcvt(num, sig, str);
   printf("string = %s\n", str);

   num = 0.345e5;
   gcvt(num, sig, str);
   printf("string = %s\n", str);

   return(0);
}
```
## 1.3 运行结果
![](gcvt.png)

# 2. getarccoords
## 2.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `void getarccoords(struct arccoordstype *arccoords);`| 取最后一次调用arc的坐标  |

## 2.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>
#include <conio.h>

int main()
{
    int gdriver = DETECT, gmode, errorcode;
    struct arccoordstype arcinfo;
    int midx, midy;
    int stangle = 45, endangle = 270;
    char sstr[80], estr[80];

    initgraph(&gdriver, &gmode, "");

    errorcode = graphresult();
    if (errorcode != grOk)
    {
      printf("Graphics error: %s\n", grapherrormsg(errorcode));
      printf("Press any key to halt:");
      getch();
      exit(1);
    }

    midx = getmaxx() / 2;
    midy = getmaxy() / 2;

    setcolor(getmaxcolor());
    arc(midx, midy, stangle, endangle, 100);
    // 取最后一次调用arc的坐标
    getarccoords(&arcinfo);

    sprintf(sstr, "*- (%d, %d)", arcinfo.xstart, arcinfo.ystart);
    sprintf(estr, "*- (%d, %d)", arcinfo.xend, arcinfo.yend);

    outtextxy(arcinfo.xstart, arcinfo.ystart, sstr);
    outtextxy(arcinfo.xend, arcinfo.yend, estr);

    getch();
    closegraph();
    return 0;
}

```
上述代码是一个简单的图形程序，使用了图形库函数 `arc` 来绘制一个弧线并显示其起始和结束点的坐标。

大致逻辑如下：
1. 初始化图形驱动和模式，创建一个空的图形窗口。
2. 检查图形操作是否成功，如果失败则输出错误信息并退出程序。
3. 计算屏幕的中心点坐标 `(midx, midy)`。
4. 设置绘图颜色为最大颜色值。
5. 在屏幕中心绘制一个弧线，起始角度为 `45` 度，结束角度为 `270` 度，半径为 `100` 像素。
6. 获取最后一次调用 `arc` 函数时的坐标信息，并将其存储在 `arcinfo` 结构体中。
7. 使用 `sprintf` 函数将起始点和结束点的坐标格式化为字符串。
8. 在屏幕上显示起始点和结束点的坐标信息。
9. 等待用户按键输入，然后关闭图形窗口并退出程序。
## 2.3 运行结果
![](getarccoords.png)

# 3. getbkcolor
## 3.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `int getbkcolor(void);`|  获取当前背景颜色 |


| 颜色值 | 英文枚举  | 中文描述 |
|--|--|--|
|0  | BLACK       |黑     |
|1  | BLUE      |蓝     |
|2  | GREEN       |绿     |
|3  | CYAN      |青     |
|4  | RED       |红     |
|5  | MAGENTA     |洋红   |
|6  | BROWN       |棕     |
|7  | LIGHTGRAY     |淡灰   |
|8  | DARKGRAY    |深灰   |
|9  | LIGHTBLUE     |淡兰   |
|10 |   LIGHTGREEN    |淡绿   |
|11 |   LIGHTCYAN     |淡青   |
|12 |   LIGHTRED    |淡红   |
|13 |   LIGHTMAGENTA  |淡洋红 |
|14 |   YELLOW      |黄     |
|15 |   WHITE       |白     |

## 3.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <conio.h>

int main(void)
{
    int gdriver = DETECT, gmode, errorcode;
    int bkcolor, midx, midy;
    char bkname[35];

    initgraph(&gdriver, &gmode, "");

    errorcode = graphresult();
    if (errorcode != grOk)
    {
        printf("Graphics error: %s\n", grapherrormsg(errorcode));
        printf("Press any key to halt:");
        getch();
        exit(1);
    }

    midx = getmaxx() / 2;
    midy = getmaxy() / 2;
    setcolor(getmaxcolor());

    settextjustify(CENTER_TEXT, CENTER_TEXT);
    cleardevice();

    for (int i = WHITE; i >= 0; i--)
    {
        setbkcolor(i);
        bkcolor = getbkcolor(); // 获取当前背景颜色
        if (i == WHITE) setcolor(BLACK);
        else setcolor(WHITE);
        itoa(bkcolor, bkname, 10);
        strcat(bkname," is the current background color.");

        outtextxy(midx, midy, bkname);
        getch();
        cleardevice();
    }

    getch();
    closegraph();
    return 0;
}

```
上述也是一个简单的图形程序，通过使用图形库来绘制背景颜色变化。

下面来简单总结下：

1. 初始化图形驱动和模式，创建一个空的图形窗口。
2. 检查图形操作是否成功，如果失败则输出错误信息并退出程序。
3. 计算屏幕的中心点坐标 `(midx, midy)`。
4. 设置绘图颜色为最大颜色值。
5. 设置文本对齐方式为中心对齐。
6. 清空设备上的图形内容。
7. 循环遍历从白色到黑色的背景颜色，每次循环执行以下操作：
  - 设置当前背景颜色为循环变量 `i` 所代表的颜色。
  - 获取当前背景颜色并将其转换为字符串形式存储在 `bkcolor` 数组中。
  - 如果当前颜色是白色，则设置文本颜色为黑色；否则设置为白色。
  - 将背景颜色信息添加到 `bkname` 字符串中。
  - 在屏幕中心位置显示包含背景颜色信息的文本。
  - 等待用户按键输入，然后清空设备上的图形内容。
8. 等待用户按键输入，然后关闭图形窗口并退出程序。

## 3.3 运行结果
![](getbkcolor.gif)

# 4. getc
## 4.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int getc(FILE *stream);` |  从流中取字符 |

## 4.2 演示示例
```c
#include <stdio.h>

int main()
{
    char ch;
    printf("Input a character:");
    ch = getc(stdin);
    printf("The character input was: '%c'\n", ch);
    return 0;
}
```
## 4.3 运行结果
![](getc.png)

# 5. getchar
## 5.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `int getchar(void);`| 从 **stdin** 流中读字符  |

## 5.2 演示示例
```c
#include <stdio.h>

int main(void)
{
   int c;
   while ((c = getchar()) != '\n')
      printf("%c ", c);

   return 0;
}
```
## 5.3 运行结果
![](getchar.png)

# 6. getcolor
## 6.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `int getcolor(void);`| 当前画线的颜色  |

| 颜色值 | 英文枚举  | 中文描述 |
|--|--|--|
|0  | BLACK       |黑     |
|1  | BLUE      |蓝     |
|2  | GREEN       |绿     |
|3  | CYAN      |青     |
|4  | RED       |红     |
|5  | MAGENTA     |洋红   |
|6  | BROWN       |棕     |
|7  | LIGHTGRAY     |淡灰   |
|8  | DARKGRAY    |深灰   |
|9  | LIGHTBLUE     |淡兰   |
|10 |   LIGHTGREEN    |淡绿   |
|11 |   LIGHTCYAN     |淡青   |
|12 |   LIGHTRED    |淡红   |
|13 |   LIGHTMAGENTA  |淡洋红 |
|14 |   YELLOW      |黄     |
|15 |   WHITE       |白     |
## 6.2 演示示例
```c
#include <graphics.h>

int main(void)
{
    int gdriver = DETECT, gmode, errorcode;
    int color, midx, midy;
    char colname[35];

    initgraph(&gdriver, &gmode, "");

    errorcode = graphresult();
    /* an error occurred */
    if (errorcode != grOk)
    {
        printf("Graphics error: %s\n", grapherrormsg(errorcode));
        printf("Press any key to halt:");
        getch();
        exit(1);
    }

    midx = getmaxx() / 2;
    midy = getmaxy() / 2;
    setcolor(getmaxcolor());

    settextjustify(CENTER_TEXT, CENTER_TEXT);

    for (int i = WHITE; i > 0; i--)
    {
        color = getcolor();
        itoa(color, colname, 10);
        strcat(colname, " is the current drawing color.");
        outtextxy(midx, midy, colname);
        getch();
        cleardevice();
        setcolor(i - 1);
    }

    getch();
    closegraph();
    return 0;
}

```

上述程序，通过使用图形库，在一个循环中遍历所有颜色，每次在屏幕中心显示当前颜色的名称和值，等待用户按键后更改颜色并清除屏幕，直到所有颜色展示完毕。

## 6.3 运行结果
![](getcolor.gif)

# 7. getcwd
## 7.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `char *getcwd(char *buffer, int maxlen);`| 获取当前工作目录  |

> **注意：getcwd** 函数是将当前工作目录的绝对路径复制到参数 **buffer** 所指的内存空间中，参数 **maxlen** 为 **buffer** 的空间大小。
## 7.2 演示示例
```c
#include <stdio.h>
#include <dir.h>

#define MAXPATH 1000

int main()
{
   char buffer[MAXPATH];
   getcwd(buffer, MAXPATH);
   printf("The current directory is: %s\n", buffer);
   return 0;
}
```
## 7.3 运行结果
![](getcwd.png)

# 8. getdefaultpalette
## 8.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`struct palettetype* getdefaultpalette(void);` |  获取调色板定义结构 |

## 8.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>

int main(void)
{
    int gdriver = DETECT, gmode, errorcode;
    int i, midx, midy;;

    struct palettetype far *pal=NULL;

    initgraph(&gdriver, &gmode, "");

    errorcode = graphresult();
    if (errorcode != grOk)
    {
        printf("Graphics error: %s\n", grapherrormsg(errorcode));
        printf("Press any key to halt:");
        getch();
        exit(1);
    }

    midx = getmaxx() / 3;
    midy = getmaxy() / 2;
    setcolor(getmaxcolor());
    // 获取调色板定义结构
    pal = getdefaultpalette();

    char buffer[100];
    for (i=BLACK; i<WHITE + 1; i++)
    {
        sprintf(buffer, "colors[%d] = %d", i, pal->colors[i]);
        outtextxy(midx, midy, buffer);
        getch();
        cleardevice();
    }

    getch();
    closegraph();
    return 0;
}

```
## 8.3 运行结果
![](getdefaultpalette.gif)

# 9. getdrivername
## 9.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`char *getdrivename(void);` |  获取当前图形驱动程序名字|

## 9.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>
#include <conio.h>

int main(void)
{
    int gdriver = DETECT, gmode, errorcode;
    char *drivername;
    initgraph(&gdriver, &gmode, "");
    errorcode = graphresult();
    if (errorcode != grOk)
    {
        printf("Graphics error: %s\n", grapherrormsg(errorcode));
        printf("Press any key to halt:");
        getch();
        exit(1);
    }

    setcolor(getmaxcolor());

    // 当前图形驱动程序名字
    drivername = getdrivername();

    settextjustify(CENTER_TEXT, CENTER_TEXT);

    outtextxy(getmaxx() / 2, getmaxy() / 2, drivername);

    getch();
    closegraph();
    return 0;
}

```
## 9.3 运行结果
![](getdrivename.png)

# 10. getfillpattern
## 10.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `void getfillpattern(char *upattern);`| 将用户定义的填充模式拷贝到内存中  |

## 10.2 演示示例
```c
#include <graphics.h>
#include <stdio.h>

int main(void)
{
    int gdriver = DETECT, gmode, errorcode;
    int maxx, maxy;
    char pattern[8] = {0x00, 0x70, 0x20, 0x27, 0x25, 0x27, 0x04, 0x04};

    initgraph(&gdriver, &gmode, "");

    errorcode = graphresult();
    if (errorcode != grOk)
    {
        printf("Graphics error: %s\n", grapherrormsg(errorcode));
        printf("Press any key to halt:");
        getch();
        exit(1);
    }

    maxx = getmaxx();
    maxy = getmaxy();
    setcolor(getmaxcolor());
    // 选择用户定义的填充模式
    setfillpattern(pattern, getmaxcolor());

    bar(0, 0, maxx, maxy);

    getch();
    // 将用户定义的填充模式拷贝到内存中
    getfillpattern(pattern);

    pattern[0] += 1;
    pattern[1] -= 2;
    pattern[2] += 3;
    pattern[3] -= 4;
    pattern[4] += 5;
    pattern[5] -= 6;
    pattern[6] += 7;
    pattern[7] -= 8;
    // 选择用户定义的填充模式
    setfillpattern(pattern, getmaxcolor());

    bar(0, 0, maxx, maxy);

    getch();
    closegraph();
    return 0;
}

```
## 10.3 运行结果
![](getfillpattern.gif)

# 11. getfillsettings
## 11.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `void getfillsettings(struct fillsettingstype *fillinfo);`|  获取有关当前填充模式和填充颜色的信息 |

## 11.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>

// the names of the fill styles supported
char *fname[] = { "EMPTY_FILL",
                  "SOLID_FILL",
                  "LINE_FILL",
                  "LTSLASH_FILL",
                  "SLASH_FILL",
                  "BKSLASH_FILL",
                  "LTBKSLASH_FILL",
                  "HATCH_FILL",
                  "XHATCH_FILL",
                  "INTERLEAVE_FILL",
                  "WIDE_DOT_FILL",
                  "CLOSE_DOT_FILL",
                  "USER_FILL"
        };

int main()
{
   int gdriver = DETECT, gmode, errorcode;
   struct fillsettingstype fillinfo;
   int midx, midy;
   char patstr[40], colstr[40];

   initgraph(&gdriver, &gmode, "");

   errorcode = graphresult();
   if (errorcode != grOk)
   {
      printf("Graphics error: %s\n", grapherrormsg(errorcode));
      printf("Press any key to halt:");
      getch();
      exit(1);
   }

   midx = getmaxx() / 2;
   midy = getmaxy() / 2;

   // 获取有关当前填充模式和填充颜色的信息
   getfillsettings(&fillinfo);

   sprintf(patstr, "%s is the fill style.", fname[fillinfo.pattern]);
   sprintf(colstr, "%d is the fill color.", fillinfo.color);

   settextjustify(CENTER_TEXT, CENTER_TEXT);
   outtextxy(midx, midy, patstr);
   outtextxy(midx, midy+2*textheight("W"), colstr);

   getch();
   closegraph();
   return 0;
}

```
## 11.3 运行结果
![](getfillsettings.png)

# 12. getgraphmode
## 12.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int getgraphmode(void);` | 获取当前图形模式  |

## 12.2 演示示例
```c
#include <graphics.h>

int main(void)
{
    int gdriver = DETECT, gmode, errorcode;
    int midx, midy, mode;
    char numname[80], modename[80];

    initgraph(&gdriver, &gmode, "");

    errorcode = graphresult();
    if (errorcode != grOk)
    {
        printf("Graphics error: %s\n", grapherrormsg(errorcode));
        printf("Press any key to halt:");
        getch();
        exit(1);
    }

    midx = getmaxx() / 2;
    midy = getmaxy() / 2;

    // 获取当前图形模式
    mode = getgraphmode();
    sprintf(numname, "%d is the current mode number.", mode);
    sprintf(modename, "%s is the current graphics mode", getmodename(mode));

    settextjustify(CENTER_TEXT, CENTER_TEXT);
    outtextxy(midx, midy, numname);
    outtextxy(midx, midy+2*textheight("W"), modename);

    getch();
    closegraph();
    return 0;
}

```
## 12.3 运行结果
![](getgraphmode.png)

# 13. getimage
## 13.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void getimage(int left, int top, int right, int bottom, void *bitmap);` |  保存指定区域的屏幕上的像素图形到指定的内存区域 |

## 13.2 演示示例
```c
#include<graphics.h>
int main()
{
    int driver,mode;
    unsigned size;
    void *buf;

    driver=DETECT;
    mode=0; initgraph(&driver,&mode,"");

    setcolor(15);
    rectangle(20,20,200,200);

    setcolor(RED);
    line(20,20,200,200);

    setcolor(GREEN);
    line(20,200,200,20);

    getch();

    size=imagesize(20,20,200,200);
    if(size!=-1)
    {
        buf=malloc(size);
        if(buf)
        {
            getimage(20,20,200,200,buf);
            putimage(100,100, buf,COPY_PUT);
            putimage(300,50, buf,COPY_PUT);
        }
    }
    outtext("press a key");
    getch();
    return 0;
}

```
## 13.3 运行结果
![](getimage.gif)


# 14. getlinesettings
## 14.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void getlinesettings(struct linesettingstype *lininfo);` | 取当前线型、模式和宽度  |

## 14.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>

// the names of the line styles supported
char *lname[] = { "SOLID_LINE",
                  "DOTTED_LINE",
                  "CENTER_LINE",
                  "DASHED_LINE",
                  "USERBIT_LINE"
                };

int main()
{
    int gdriver = DETECT, gmode, errorcode;
    struct linesettingstype lineinfo;
    int midx, midy;
    char lstyle[80], lpattern[80], lwidth[80];

    initgraph(&gdriver, &gmode, "");
    errorcode = graphresult();
    if (errorcode != grOk)
    {
        printf("Graphics error: %s\n", grapherrormsg(errorcode));
        printf("Press any key to halt:");
        getch();
        exit(1);
    }

    midx = getmaxx() / 2;
    midy = getmaxy() / 2;

    // 取当前线型、模式和宽度
    getlinesettings(&lineinfo);

    sprintf(lstyle, "%s is the line style.", lname[lineinfo.linestyle]);
    sprintf(lpattern, "0x%X is the user-defined line pattern.", lineinfo.upattern);
    sprintf(lwidth, "%d is the line thickness.", lineinfo.thickness);

    settextjustify(CENTER_TEXT, CENTER_TEXT);
    outtextxy(midx, midy, lstyle);
    outtextxy(midx, midy+2*textheight("W"), lpattern);
    outtextxy(midx, midy+4*textheight("W"), lwidth);

    getch();
    closegraph();
    return 0;
}

```
## 14.3 运行结果
![](getlinesettings.png)


# 15. getmaxcolor
## 15.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `int getmaxcolor(void);`|  可以传给函数 `setcolor` 的最大颜色值 |

## 15.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>

int main(void)
{
    int gdriver = DETECT, gmode, errorcode;
    int midx, midy;
    char colstr[80];

    initgraph(&gdriver, &gmode, "");

    errorcode = graphresult();
    if (errorcode != grOk)
    {
        printf("Graphics error: %s\n", grapherrormsg(errorcode));
        printf("Press any key to halt:");
        getch();
        exit(1);
    }

    midx = getmaxx() / 2;
    midy = getmaxy() / 2;

    sprintf(colstr, "This mode supports colors 0~%d", getmaxcolor());

    settextjustify(CENTER_TEXT, CENTER_TEXT);
    outtextxy(midx, midy, colstr);

    getch();
    closegraph();
    return 0;
}

```
## 15.3 运行结果
![](getmaxcolor.png)


# 16. getmaxx，getmaxy
## 16.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int getmaxx(void);` |  屏幕的最大x坐标 |
| `int getmaxy(void);`|  屏幕的最大y坐标 |
## 16.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>

int main(void)
{
    int gdriver = DETECT, gmode, errorcode;
    int midx, midy;
    char xrange[80], yrange[80];

    initgraph(&gdriver, &gmode, "");

    errorcode = graphresult();
    if (errorcode != grOk)
    {
        printf("Graphics error: %s\n", grapherrormsg(errorcode));
        printf("Press any key to halt:");
        getch();
        exit(1);
    }

    midx = getmaxx() / 2;
    midy = getmaxy() / 2;

    sprintf(xrange, "X values range from 0~%d", getmaxx());
    sprintf(yrange, "Y values range from 0~%d", getmaxy());

    settextjustify(CENTER_TEXT, CENTER_TEXT);
    outtextxy(midx, midy, xrange);
    outtextxy(midx, midy+2*textheight("W"), yrange);

    getch();
    closegraph();
    return 0;
}
```
## 16.3 运行结果
![](getmax.png)

# 参考
1. [\[API Reference Document\]](https://www.apiref.com/c-zh/t_alph_g.htm)

