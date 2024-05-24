---
title: C语言函数大全--b开头的函数
date: 2023-03-19 16:46:12
updated: 2023-06-25 23:24:36
categories:
  - 开发语言-C
tags:
  - C语言函数大全
  - b开头的函数
---

![](/images/cplus-logo.png)

# 总览

| 函数声明 |  函数功能  |
|:--|:--|
|`void bar(int left, int top, int right, int bottom);` | 画一个二维条形图  |
| `void bar3d(int left, int top, int right, int bottom, int depth, int topflag);`| 画一个三维条形图  |
| `int bdos(int dosfun, unsigned dosdx, unsigned dosal);`| DOS系统调用  |
| `int bdosptr(int dosfun, void *argument, unsigned dosal);`| DOS系统调用  |
| `int bioscom(int cmd, char abyte, int port);` |  串行I/O通信 |
|`int biosdisk(int cmd, int drive, int head, int track, int sector, int nsects, void *buffer);` |  软硬盘I/O |
|`int bioskey(int cmd);` |  直接使用BIOS服务的键盘接口 |
| `int biosmemory(void);` | 返回存储块大小，以K为单位  |
|`int biosprint(int cmd, int byte, int port);` |  直接使用BIOS服务的打印机I/O |
| `long biostime(int cmd, long newtime);`|  读取或设置BIOS时间 |
|`int brk(void *endds);` |  用来改变分配给调用程序的数据段的空间数量 |
|`void *bsearch(const void *key, const void *base, size_t *nelem,  size_t width, int(*fcmp)(const void *, const *));` |  二分法搜索 |


# 1. bar
## 1.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void bar(int left, int top, int right, int bottom);` | 画一个二维条形图  |

> **关注点：** 绘制二维条形图需要左上角和右下角的坐标。 **left** 指定左上角的 X 坐标，**top** 指定左上角的 Y 坐标，**right** 指定右下角的 X 坐标，**bottom** 指定右下角的 Y 坐标。 当前填充图案和填充颜色用于填充条形图。 要更改填充图案和填充颜色，请使用 **setfillstyle**。

## 1.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>
#include <conio.h>

int main(void)
{
    /* request auto detection */
    int gdriver = DETECT, gmode, errorcode;
    int midx, midy, i;

    /* initialize graphics and local variables */
    initgraph(&gdriver, &gmode, "");

    /* read result of initialization */
    errorcode = graphresult();
    if (errorcode != grOk)  /* an error occurred */
    {
        printf("Graphics error: %s\n", grapherrormsg(errorcode));
        printf("Press any key to halt:");
        getch();
        exit(1); /* terminate with an error code */
    }

    midx = getmaxx() / 2;
    midy = getmaxy() / 2;

    /* loop through the fill patterns */
    for (i=SOLID_FILL; i<USER_FILL; i++)
    {
        /* set the fill style */
        setfillstyle(i, getmaxcolor());

        /* draw the bar */
        bar(midx-50, midy-50, midx+50, midy+50);

        getch();
    }

    /* clean up */
    closegraph();
    return 0;
}
```
## 1.3 运行结果
![](bar.png)

# 2. bar3d
## 2.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `void bar3d(int left, int top, int right, int bottom, int depth, int topflag);`| 画一个三维条形图  |

> **关注点：**  绘制三维条形图需要条形左上角和右下角的坐标。 **left** 指定左上角的 X 坐标，**top** 指定左上角的 Y 坐标，**right** 指定右下角的 X 坐标，**bottom** 指定右下角的 Y 坐标，**depth** 指定条的深度 以像素为单位，**topflag** 确定是否将 3 维顶部放置在条形图上（如果它不为零，则放置否则不放置）。 当前填充图案和填充颜色用于填充条形图。 要更改填充图案和填充颜色，请使用 **setfillstyle**。

## 2.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>
#include <conio.h>

int main(void)
{
    /* request auto detection */
    int gdriver = DETECT, gmode, errorcode;
    int midx, midy, i;

    /* initialize graphics, local variables */
    initgraph(&gdriver, &gmode, "");

    /* read result of initialization */
    errorcode = graphresult();
    if (errorcode != grOk)  /* an error occurred */
    {
        printf("Graphics error: %s\n", grapherrormsg(errorcode));
        printf("Press any key to halt:");
        getch();
        exit(1); /* terminate with error code */
    }

    midx = getmaxx() / 2;
    midy = getmaxy() / 2;

    /* loop through the fill patterns */
    for (i=EMPTY_FILL; i<USER_FILL; i++)
    {
        /* set the fill style */
        setfillstyle(i, getmaxcolor());

        /* draw the 3-d bar */
        bar3d(midx-50, midy-50, midx+50, midy+50, 10, 1);

        getch();
    }

    /* clean up */
    closegraph();
    return 0;
}
```
## 2.3 运行结果
![](bar3d.png)

# 3. bdos
## 3.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `int bdos(int dosfun, unsigned dosdx, unsigned dosal);`| DOS系统调用  |

## 3.2 演示示例
```c
#include <stdio.h>
#include <dos.h>

/* Get current drive as 'A', 'B', ... */
char current_drive(void)
{
    char curdrive;

    /* Get current disk as 0, 1, ... */
    curdrive = bdos(0x19, 0, 0);
    return('A' + curdrive);
}

int main(void)
{
    printf("The current drive is %c:\n", current_drive());
    return 0;
}
```


# 4. bdosptr
## 4.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `int bdosptr(int dosfun, void *argument, unsigned dosal);`| DOS系统调用  |

## 4.2 演示示例
```c
#include <string.h>
#include <stdio.h>
#include <dir.h>
#include <dos.h>
#include <errno.h>
#include <stdlib.h>

#define  BUFLEN  80

int main(void)
{
    char  buffer[BUFLEN];
    int   test;

    printf("Enter full pathname of a directory\n");
    gets(buffer);

    test = bdosptr(0x3B,buffer,0);
    if(test)
    {
        printf("DOS error message: %d\n", errno);
        /* See errno.h for error listings */
        exit (1);
    }

    getcwd(buffer, BUFLEN);
    printf("The current directory is: %s\n", buffer);

    return 0;
}
```


# 5. bioscom
## 5.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `int bioscom(int cmd, char abyte, int port);` |  串行I/O通信 |

## 5.2 演示示例
```c
#include <bios.h>
#include <conio.h>

#define COM1       0
#define DATA_READY 0x100
#define TRUE       1
#define FALSE      0

#define SETTINGS ( 0x80 | 0x02 | 0x00 | 0x00)

int main(void)
{
    int in, out, status, DONE = FALSE;

    bioscom(0, SETTINGS, COM1);
    cprintf("... BIOSCOM [ESC] to exit ...\n");
    while (!DONE)
    {
        status = bioscom(3, 0, COM1);
        if (status & DATA_READY)
            if ((out = bioscom(2, 0, COM1) & 0x7F) != 0)
                putch(out);
        if (kbhit())
        {
            if ((in = getch()) == '\x1B')
                DONE = TRUE;
            bioscom(1, in, COM1);
        }
    }
    return 0;
}
```


# 6. biosdisk
## 6.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int biosdisk(int cmd, int drive, int head, int track, int sector, int nsects, void *buffer);` |  软硬盘I/O |

## 6.2 演示示例
```c
#include <bios.h>
#include <stdio.h>

int main(void)
{
    int result;
    char buffer[512];

    printf("Testing to see if drive a: is ready\n");
    result = biosdisk(4,0,0,0,0,1,buffer);
    result &= 0x02;
    (result) ? (printf("Drive A: Ready\n")) : (printf("Drive A: Not Ready\n"));

    return 0;
}
```


# 7. bioskey
## 7.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int bioskey(int cmd);` |  直接使用BIOS服务的键盘接口 |

## 7.2 演示示例
```c
#include <stdio.h>
#include <bios.h>
#include <ctype.h>

#define RIGHT  0x01
#define LEFT   0x02
#define CTRL   0x04
#define ALT    0x08

int main(void)
{
    int key, modifiers;

    /* function 1 returns 0 until a key is pressed */
    while (bioskey(1) == 0);

    /* function 0 returns the key that is waiting */
    key = bioskey(0);

    /* use function 2 to determine if shift keys were used */
    modifiers = bioskey(2);
    if (modifiers)
    {
        printf("[");
        if (modifiers & RIGHT) printf("RIGHT");
        if (modifiers & LEFT)  printf("LEFT");
        if (modifiers & CTRL)  printf("CTRL");
        if (modifiers & ALT)   printf("ALT");
        printf("]");
    }
    /* print out the character read */
    if (isalnum(key & 0xFF))
        printf("'%c'\n", key);
    else
        printf("%#02x\n", key);
    return 0;
}
```


# 8. biosmemory
## 8.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `int biosmemory(void);` | 返回存储块大小，以K为单位  |

## 8.2 演示示例
```c
#include <stdio.h>
#include <bios.h>

int main(void)
{
    int memory_size;

    memory_size = biosmemory();  /* returns value up to 640K */
    printf("RAM size = %dK\n",memory_size);
    return 0;
}
```


# 9. biosprint
## 9.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int biosprint(int cmd, int byte, int port);` |  直接使用BIOS服务的打印机I/O |

## 9.2 演示示例
```c
#include <stdio.h>
#include <conio.h>
#include <bios.h>

int main(void)
{
    #define STATUS  2    /* printer status command */
    #define PORTNUM 0    /* port number for LPT1 */

    int status, abyte=0;

    printf("Please turn off your printer.  Press any key to continue\n");
    getch();
    status = biosprint(STATUS, abyte, PORTNUM);
    if (status & 0x01)
        printf("Device time out.\n");
    if (status & 0x08)
        printf("I/O error.\n");

    if (status & 0x10)
        printf("Selected.\n");
    if (status & 0x20)
        printf("Out of paper.\n");

    if (status & 0x40)
        printf("Acknowledge.\n");
    if (status & 0x80)
        printf("Not busy.\n");

    return 0;
}
```

# 10. biostime
## 10.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `long biostime(int cmd, long newtime);`|  读取或设置BIOS时间 |

## 10.2 演示示例
```c
#include <stdio.h>
#include <bios.h>
#include <time.h>
#include <conio.h>

int main(void)
{
    long bios_time;

    clrscr();
    cprintf("The number of clock ticks since midnight is:\r\n");
    cprintf("The number of seconds since midnight is:\r\n");
    cprintf("The number of minutes since midnight is:\r\n");
    cprintf("The number of hours since midnight is:\r\n");
    cprintf("\r\nPress any key to quit:");
    while(!kbhit())
    {
        bios_time = biostime(0, 0L);

        gotoxy(50, 1);
        cprintf("%lu", bios_time);

        gotoxy(50, 2);
        cprintf("%.4f", bios_time / CLK_TCK);

        gotoxy(50, 3);
        cprintf("%.4f", bios_time / CLK_TCK / 60);

        gotoxy(50, 4);
        cprintf("%.4f", bios_time / CLK_TCK / 3600);
    }
    return 0;
}
```


# 11. brk
## 11.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int brk(void *endds);` |  用来改变分配给调用程序的数据段的空间数量 |

## 11.2 演示示例
```c
#include <stdio.h>
#include <alloc.h>

int main(void)
{
    char *ptr;

    printf("Changing allocation with brk()\n");
    ptr = malloc(1);
    printf("Before brk() call: %lu bytes free\n", coreleft());
    brk(ptr+1000);
    printf(" After brk() call: %lu bytes free\n", coreleft());
    return 0;
}
```

# 12. bsearch
## 12.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void *bsearch(const void *key, const void *base, size_t *nelem,  size_t width, int(*fcmp)(const void *, const *));` |  二分法搜索 |

## 12.2 演示示例
```c
#include <stdlib.h>
#include <stdio.h>

#define NELEMS(arr) (sizeof(arr) / sizeof(arr[0]))

int numarray[] = {1, 2, 3, 5, 6, 8, 9, 10, 12, 14};

int numeric (const int *p1, const int *p2)
{
    return(*p1 - *p2);
}

int lookup(int key)
{
    int *itemptr;

    /* The cast of (int(*)(const void *,const void*)) is needed to avoid a type mismatch error at compile time */
    itemptr = (int(*))bsearch(&key, numarray, NELEMS(numarray), sizeof(int), (int(*)(const void *,const void *))numeric);
    return (itemptr != NULL);
}

int main(void)
{
    int a;
    printf("Please input key: ");
    scanf("%d", &a);
    if (lookup(a))
        printf("%d is in the table.\n", a);
    else
        printf("%d isn't in the table.\n", a);
    return 0;
}
```
## 12.3 运行结果
![](bsearch.png)
![](bsearch-1.png)

# 参考
1. [[API Reference Document]](https://www.apiref.com/c-zh/index.htm)
2. [\[c语言中的 bar 函数\]](https://21xrx.com/Articles/read_article/259)
3. [\[c语言中的 bar3d 函数\]](https://21xrx.com/Articles/read_article/260)


