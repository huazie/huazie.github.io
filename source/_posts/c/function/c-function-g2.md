---
title: C语言函数大全--g开头的函数（下）
date: 2023-04-07 23:43:58
updated: 2025-01-19 22:32:38
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
| `char * getmodename(int mode_name);`| 获取指定的图形模式名 |
| `void getmoderange(int graphdriver, int *lomode, int *himode);`|  获取给定图形驱动程序的模式范围 |
|`void getpalette(struct palettetype *palette);` | 获取有关当前调色板的信息  |
|`int getpixel(int x, int y);` | 获取得指定像素的颜色  |
|`char *gets(char *str);` |  从标准输入流中读取字符串，直至遇到到换行符或EOF时停止，并将读取的结果存放在 **buffer** 指针所指向的字符数组中。<br/> 换行符不作为读取串的内容，读取的换行符被转换为 `'\0'` 空字符，并由此来结束字符串。 |
|`void gettextsettings(struct textsettingstype *textinfo);` |  获取有关当前图形文本字体的信息 |
| `void getviewsettings(struct viewporttype *viewport);`| 获取有关当前视区的信息  |
|`int getw(FILE *strem);` |  从 stream 所指向文件读取下一个整数 |
| `int getx(void);`|  获取当前图形位置的 x 坐标 |
| `int gety(void);` |  获取当前图形位置的 y 坐标 |
|`struct tm *gmtime(long *clock);` | 把日期和时间转换为格林尼治标准时间(GMT)  |
|`void graphdefaults(void);` |  将所有图形设置复位为它们的缺省值 |
|`char * grapherrormsg(int errorcode);` | 返回一个错误信息串的指针  |
|`int graphresult(void);` |  返回最后一次不成功的图形操作的错误代码 |
| `int getmaxwidth(void);`| 获取屏幕的最大宽度  |
| `int getmaxheight(void);`| 获取屏幕的最大高度  |
|`int getdisplaycolor( int color );` |  根据 color ，返回要显示的颜色值  |
|`int getwindowwidth();` |  获取图形界面窗口宽度 |
|`int getwindowheight(void);` |  获取图形界面窗口高度 |
| `bool getrefreshingbgi(void);`|  获取刷新基础图形界面标识 |

# 1. getmodename
## 1.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `char * getmodename(int mode_name);`| 获取指定的图形模式名 |

**参数：**
- `mode_name` : 当前图形模式的模式代码。不同的图形模式对应不同的分辨率、颜色深度和显示模式。

## 1.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>

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

    mode = getgraphmode();
    sprintf(numname, "%d is the current mode number.", mode);
    sprintf(modename, "%s is the current graphics mode.", getmodename(mode));

    settextjustify(CENTER_TEXT, CENTER_TEXT);
    outtextxy(midx, midy, numname);
    outtextxy(midx, midy+2*textheight("W"), modename);

    getch();
    closegraph();
    return 0;
}
```
## 1.3 运行结果
![](getmodename.png)

# 2. getmoderange
## 2.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `void getmoderange(int graphdriver, int *lomode, int *himode);`|  获取给定图形驱动程序的模式范围 |

**参数：**

- `graphdriver` : 图形驱动程序的标识符。不同的图形驱动程序有不同的标识符，用于指定你希望使用的图形硬件或软件环境。例如，在某些图形库中，特定的数字或宏定义（如DETECT）可以用来自动检测可用的图形驱动程序。
- `lomode` : 一个指向整数的指针，用于接收指定图形驱动程序支持的最低显示模式编号。
- `himode` : 一个指向整数的指针，用于接收指定图形驱动程序支持的最高显示模式编号。


## 2.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>

int main()
{
    int gdriver = DETECT, gmode, errorcode;
    int midx, midy, low, high;
    char mrange[80];

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

    getmoderange(gdriver, &low, &high);

    sprintf(mrange, "This driver supports modes %d~%d", low, high);

    settextjustify(CENTER_TEXT, CENTER_TEXT);
    outtextxy(midx, midy, mrange);

    getch();
    closegraph();
    return 0;
}
```
## 2.3 运行结果

![](getmoderange.png)

# 3. getpalette
## 3.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void getpalette(struct palettetype *palette);` | 获取有关当前调色板的信息  |

**参数：**

- `palette` : 一个指向`palettetype`结构体的指针。`palettetype`结构体通常包含了一系列元素，每个元素代表调色板中的一个颜色条目。在标准的图形库中（如**Borland**的**BGI**图形库），`palettetype`结构体可能包含多个`unsigned char`类型的成员，每个成员对应调色板中的一个颜色通道（如红色、绿色、蓝色），以及可能的其他信息（如亮度或透明度）。

> **注意:** palettetype结构体的确切定义可能依赖于你使用的图形库。在某些实现中，它可能是一个简单的数组，每个元素代表一个颜色（可能是**RGB**值的一个组合），或者是一个更复杂的结构体，包含了关于每个颜色条目的更多信息。

## 3.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>

int main()
{
    int gdriver = DETECT, gmode, errorcode;
    struct palettetype pal;
    char psize[80], pval[20];
    int i, ht;
    int y = 10;

    initgraph(&gdriver, &gmode, "");

    errorcode = graphresult();
    if (errorcode != grOk)
    {
        printf("Graphics error: %s\n", grapherrormsg(errorcode));
        printf("Press any key to halt:");
        getch();
        exit(1);
    }

    getpalette(&pal);

    sprintf(psize, "The palette has %d modifiable entries.", pal.size);

    outtextxy(0, y, psize);
    if (pal.size != 0)
    {
        ht = textheight("W");
        y += 2*ht;
        outtextxy(0, y, "Here are the current values:");
        y += 2*ht;
        for (i=0; i < pal.size; i++, y+=ht)
        {
            sprintf(pval, "palette[%02d]: 0x%02X", i, pal.colors[i]);
            outtextxy(0, y, pval);
        }
    }

    getch();
    closegraph();
    return 0;
}
```
## 3.3 运行结果

![](getpalette.png)

# 4. getpixel
## 4.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int getpixel(int x, int y);` | 获取得指定像素的颜色  |

**参数：**

- `x` : 想要获取的像素颜色值的点的横坐标。坐标原点`(0, 0)`通常位于屏幕的左上角。
- `y` : 想要获取的像素颜色值的点的纵坐标。

**返回值：**

函数返回一个整数，该整数代表指定坐标 (x, y) 上像素的颜色编码。颜色编码的具体含义取决于你使用的图形库和当前的图形设置。在某些图形库中，这个整数可能直接代表一个RGB颜色值，其中不同的位或字节表示红色、绿色和蓝色通道的强度。在其他情况下，这个整数可能是一个索引值，指向当前调色板中的一个颜色条目。

## 4.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>
#include <time.h>

#define PIXEL_COUNT 1000
#define DELAY_TIME  100

int main(void)
{
    int gdriver = DETECT, gmode, errorcode;
    int midx, midy, x, y, color, maxx, maxy, maxcolor;
    char mPixel[50];

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

    maxx = getmaxx() + 1;
    maxy = getmaxy() + 1;
    maxcolor = getmaxcolor() + 1;

    while (!kbhit())
    {
        srand((unsigned)time(NULL));
        x = rand() % maxx;
        y = rand() % maxy;
        color = rand() % maxcolor;
        putpixel(x, y, color);

        sprintf(mPixel, "color of pixel at (%d,%d) = %d", x, y, getpixel(x, y));
        settextjustify(CENTER_TEXT, CENTER_TEXT);
        outtextxy(midx, midy, mPixel);

        delay(DELAY_TIME);

        cleardevice();
    }

    getch();
    closegraph();
    return 0;
}


```
## 4.3 运行结果

![](getpixel.gif)


# 5. gets
## 5.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`char *gets(char *str);` |  从标准输入流中读取字符串，直至遇到到换行符或EOF时停止，并将读取的结果存放在 **buffer** 指针所指向的字符数组中。<br/> 换行符不作为读取串的内容，读取的换行符被转换为 `'\0'` 空字符，并由此来结束字符串。 |

> **注意： gets** 函数可以无限读取，易发生溢出。如果溢出，多出来的字符将被写入到堆栈中，这就覆盖了堆栈原先的内容，破坏一个或多个不相关变量的值。
## 5.2 演示示例
```c
#include <stdio.h>

int main()
{
   char string[80];

   printf("Input a string:");
   gets(string);
   printf("The string input was: %s\n", string);
   return 0;
}
```
## 5.3 运行结果
![](gets.png)


# 6. gettextsettings
## 6.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void gettextsettings(struct textsettingstype *textinfo);` |  获取有关当前图形文本字体的信息 |

**参数：**

- `textinfo` : 一个指向 `textsettingstype` 结构体的指针。该结构体用于存储当前的文本设置。`textsettingstype` 结构体的具体定义取决于你使用的图形库。在不同的图形库中，这个结构体可能包含不同的成员，以反映该库支持的文本设置选项。

## 6.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>
#include <conio.h>

// 文本字体
char *font[] = { "DEFAULT_FONT",
                 "TRIPLEX_FONT",
                 "SMALL_FONT",
                 "SANS_SERIF_FONT",
                 "GOTHIC_FONT"
               };

// 文本方向
char *dir[] = { "HORIZ_DIR", "VERT_DIR" };

// 文本水平对齐方式
char *hjust[] = { "LEFT_TEXT", "CENTER_TEXT", "RIGHT_TEXT" };

// 文本垂直对齐方式
char *vjust[] = { "BOTTOM_TEXT", "CENTER_TEXT", "TOP_TEXT" };

int main()
{
    int gdriver = DETECT, gmode, errorcode;
    struct textsettingstype textinfo;
    int midx, midy, ht;
    char fontstr[80], dirstr[80], sizestr[80];
    char hjuststr[80], vjuststr[80];

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

    // 获取有关当前图形文本字体的信息
    gettextsettings(&textinfo);

    sprintf(fontstr, "%s is the text style.", font[textinfo.font]);
    sprintf(dirstr, "%s is the text direction.", dir[textinfo.direction]);
    sprintf(sizestr, "%d is the text size.", textinfo.charsize);
    sprintf(hjuststr, "%s is the horizontal justification.", hjust[textinfo.horiz]);
    sprintf(vjuststr, "%s is the vertical justification.", vjust[textinfo.vert]);

    ht = textheight("W");
    settextjustify(CENTER_TEXT, CENTER_TEXT);
    outtextxy(midx, midy, fontstr);
    outtextxy(midx, midy+2*ht, dirstr);
    outtextxy(midx, midy+4*ht, sizestr);
    outtextxy(midx, midy+6*ht, hjuststr);
    outtextxy(midx, midy+8*ht, vjuststr);

    getch();
    closegraph();
    return 0;
}

```
## 6.3 运行结果
![](gettextsettings.png)


# 7. getviewsettings
## 7.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `void getviewsettings(struct viewporttype *viewport);`| 获取有关当前视区的信息  |

**参数：**

- `viewport` : 一个指向 `viewporttype` 结构体的指针。该结构体用于存储当前的视口设置。调用 `getviewsettings` 函数后，这个结构体将被填充为当前的视口参数。 `viewporttype` 结构体的具体定义可能依赖于你使用的图形库，但通常它会包含以下成员：
    - `left, top`: 这两个成员定义了视口的左上角坐标。坐标原点通常位于屏幕的左上角。
    - `right, bottom`: 这两个成员定义了视口的右下角坐标。
    - `clip`: 一个用于指示视口是否启用裁剪的标志。如果启用了裁剪，那么任何在视口之外的图形输出都将被忽略。

## 7.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>

char *clip[] = { "OFF", "ON" };

int main()
{
    int gdriver = DETECT, gmode, errorcode;
    struct viewporttype viewinfo;
    int midx, midy, ht;
    char topstr[80], botstr[80], clipstr[80];

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

    // 获取有关当前视区的信息
    getviewsettings(&viewinfo);

    sprintf(topstr, "(%d, %d) is the upper left viewport corner.", viewinfo.left, viewinfo.top);
    sprintf(botstr, "(%d, %d) is the lower right viewport corner.", viewinfo.right, viewinfo.bottom);
    sprintf(clipstr, "Clipping is turned %s.", clip[viewinfo.clip]);

    settextjustify(CENTER_TEXT, CENTER_TEXT);
    ht = textheight("W");
    outtextxy(midx, midy, topstr);
    outtextxy(midx, midy+2*ht, botstr);
    outtextxy(midx, midy+4*ht, clipstr);

    getch();
    closegraph();
    return 0;
}

```
## 7.3 运行结果
![](getviewsettings.png)


# 8. getw
## 8.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int getw(FILE *strem);` |  从 stream 所指向文件读取下一个整数 |

## 8.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>

#define FNAME "test.txt"

int main(void)
{
    FILE *fp;
    int word;

    fp = fopen(FNAME, "wb");
    if (fp == NULL)
    {
        printf("Error opening file %s\n", FNAME);
        exit(1);
    }

    word = 94;
    putw(word,fp);
    if (ferror(fp))
        printf("Error writing to file\n");
    else
        printf("Successful write\n");
    fclose(fp);

    fp = fopen(FNAME, "rb");
    if (fp == NULL)
    {
        printf("Error opening file %s\n", FNAME);
        exit(1);
    }

    word = getw(fp);
    if (ferror(fp))
        printf("Error reading file\n");
    else
        printf("Successful read: word = %d\n", word);

    fclose(fp);
    unlink(FNAME);

    return 0;
}
```
## 8.3 运行结果
![](getw.png)


# 9. getx，gety
## 9.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `int getx(void);`|  获取当前图形位置的 x 坐标 |
| `int gety(void);` |  获取当前图形位置的 y 坐标 |

## 9.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>

int main(void)
{
    int gdriver = DETECT, gmode, errorcode;
    char msg[80];

    initgraph(&gdriver, &gmode, "");

    errorcode = graphresult();
    if (errorcode != grOk)
    {
        printf("Graphics error: %s\n", grapherrormsg(errorcode));
        printf("Press any key to halt:");
        getch();
        exit(1);
    }

    moveto(getmaxx() / 2, getmaxy() / 2);

    sprintf(msg, "<-(%d, %d) is the here.", getx(), gety());

    outtext(msg);

    getch();
    closegraph();
    return 0;
}

```
## 9.3 运行结果
![](getxy.png)


# 10. gmtime
## 10.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`struct tm *gmtime(long *clock);` | 把日期和时间转换为格林尼治标准时间(GMT)  |

## 10.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <dos.h>

// 太平洋标准时间和夏令时
char *tzstr = "TZ=PST8PDT";

int main(void)
{
    time_t t;
    struct tm *gmt, *area;
    putenv(tzstr); // 用来改变或增加环境变量的内容
    tzset(); // UNIX时间兼容函数
    // 获取当前的系统时间，其值表示从协调世界时（Coordinated Universal Time）
    // 1970年1月1日00:00:00（称为UNIX系统的Epoch时间）到当前时刻的秒数。
    t = time(NULL); 
    area = localtime(&t); // 把从1970-1-1零点零分到当前时间系统所偏移的秒数时间转换为本地时间
    // asctime 把timeptr指向的tm结构体中储存的时间转换为字符串
    printf("Local time is: %s", asctime(area));
    // 把日期和时间转换为格林尼治标准时间(GMT)
    gmt = gmtime(&t);
    printf("GMT is:        %s", asctime(gmt));
    return 0;
}
```
## 10.3 运行结果

![](gmtime.png)

# 11. graphdefaults
## 11.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void graphdefaults(void);` |  将所有图形设置复位为它们的缺省值 |

## 11.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>
#include <conio.h>

int main(void)
{
    int gdriver = DETECT, gmode, errorcode;
    int maxx, maxy;

    initgraph(&gdriver, &gmode, "c:\\bor\\Borland\\bgi");
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

    setlinestyle(DOTTED_LINE, 0, 3);
    line(0, 0, maxx, maxy);
    outtextxy(maxx/2, maxy/3, "Before default values are restored.");
    getch();

    // 将所有图形设置复位为它们的缺省值
    graphdefaults();

    cleardevice();

    line(0, 0, maxx, maxy);
    outtextxy(maxx/2, maxy/3, "After restoring default values.");

    getch();
    closegraph();
    return 0;
}

```
## 11.3 运行结果

![](graphdefaults.gif)

# 12. grapherrormsg
## 12.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`char * grapherrormsg(int errorcode);` | 返回一个错误信息串的指针  |

## 12.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>

#define NONSENSE -50

int main(void)
{
    // FORCE AN ERROR TO OCCUR
    int gdriver = NONSENSE, gmode, errorcode;

    initgraph(&gdriver, &gmode, "");

    errorcode = graphresult();

    if (errorcode != grOk)
    {
        printf("Graphics error: %s\n", grapherrormsg(errorcode));
        printf("Press any key to halt:");
        getch();
        exit(1);
    }

    line(0, 0, getmaxx(), getmaxy());

    getch();
    closegraph();
    return 0;
}

```
## 12.3 运行结果
![](grapherrormsg.png)


# 13. graphresult
## 13.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int graphresult(void);` |  返回最后一次不成功的图形操作的错误代码 |

## 13.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>
#include <conio.h>

int main(void)
{
    // request auto detection
    int gdriver = DETECT, gmode, errorcode;

    initgraph(&gdriver, &gmode, "");

    errorcode = graphresult();
    if (errorcode != grOk)
    {
        printf("Graphics error: %s\n", grapherrormsg(errorcode));
        printf("Press any key to halt:");
        getch();
        exit(1);
    }

    line(0, 0, getmaxx(), getmaxy());

    getch();
    closegraph();
    return 0;
}

```
## 13.3 运行结果
![](graphresult.png)


# 14. getmaxwidth，getmaxheight
## 14.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `int getmaxwidth(void);`| 获取屏幕的最大宽度  |
| `int getmaxheight(void);`| 获取屏幕的最大高度  |


## 14.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>

int main()
{
    int gdriver = DETECT, gmode, errorcode;
    int midx, midy;
    char ch[80];

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

    sprintf(ch, "maxwidth = %d, maxheight = %d", getmaxwidth(), getmaxheight());

    settextjustify(CENTER_TEXT, CENTER_TEXT);
    outtextxy(midx, midy, ch);

    getch();
    closegraph();
    return 0;
}
```
## 14.3 运行结果
![](getmaxwh.png)


# 15. getdisplaycolor
## 15.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int getdisplaycolor( int color );` |  根据 color ，返回要显示的颜色值  |

> **注意：**  color = -1 , 则返回 WHITE = 15 的颜色值；color < - 1 或 color > 15，则输出一个8位整数。
## 15.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>

int main()
{
    int gdriver = DETECT, gmode, errorcode;
    int midx, midy;
    char ch[80];

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

    sprintf(ch, "color = %d, displaycolor(-1) = %d, displaycolor(16) = %d", getcolor(), getdisplaycolor(-1), getdisplaycolor(16));

    settextjustify(CENTER_TEXT, CENTER_TEXT);
    outtextxy(midx, midy, ch);

    getch();
    closegraph();
    return 0;
}


```
## 15.3 运行结果
![](getdisplaycolor.png)



# 16. getwindowwidth，getwindowheight
## 16.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int getwindowwidth(void);` |  获取图形界面窗口宽度 |
|`int getwindowheight(void);` |  获取图形界面窗口高度 |

## 16.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>

int main()
{
    int gdriver = DETECT, gmode, errorcode;
    int midx, midy, low, high;
    char ch[80];

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

    sprintf(ch, "windowwidth = %d, windowheight = %d", getwindowwidth(), getwindowheight());

    settextjustify(CENTER_TEXT, CENTER_TEXT);
    outtextxy(midx, midy, ch);

    getch();
    closegraph();
    return 0;
}

```
## 16.3 运行结果
![](getwindowwh.png)


# 17. getrefreshingbgi
## 17.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `bool getrefreshingbgi(void);`|  获取刷新基础图形界面标识 |

## 17.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>

int main()
{
    int gdriver = DETECT, gmode, errorcode;
    int midx, midy, low, high;
    char ch[80];

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

    sprintf(ch, "refreshingbgi = %d", getrefreshingbgi());

    settextjustify(CENTER_TEXT, CENTER_TEXT);
    outtextxy(midx, midy, ch);

    getch();
    closegraph();
    return 0;
}

```
## 17.3 运行结果

![](getrefreshingbgi.png)


# 参考
1. [\[API Reference Document\]](https://www.apiref.com/c-zh/t_alph_g.htm)
2. [\[gets\]](https://baike.baidu.com/item/gets?fromModule=lemma_search-box)

