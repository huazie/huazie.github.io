---
title: C语言函数大全--a开头的函数
date: 2023-03-15 23:22:23
updated: 2023-06-25 23:20:36
categories:
  - 开发语言-C
tags:
  - C语言函数大全
  - a开头的函数
---

![](/images/cplus-logo.png)

# 总览

| 函数声明 |  函数功能  |
|:--|:--|
|`void abort(void);` |  异常终止一个进程 |
|`int abs(int i);` |  求整数的绝对值 |
| `int absread(int drive, int nsects, int sectno, void *buffer);`|  从drive指定的驱动器磁盘上，sectno指定的逻辑扇区号开始读取nsects个(最多64K个)扇区的内容，储存于buffer所指的缓冲区中。 |
|`int abswrite(int drive, int nsects, int sectno, void *buffer);`| 将指定内容写入磁盘上的指定扇区|
|`int access(const char *filename, int amode);` | 确定文件的访问权限  |
|`double acos(double x);` |  反余弦函数 |
| `int allocmem(unsigned size, unsigned *seg);`|  分配DOS存储段 |
|`void arc(int x, int y, int stangle, int endangle, int radius);` | 画一弧线  |
|`char *asctime(const struct tm *tblock);` | 转换日期和时间为ASCII码  |
| `double asin(double x);`| 反正弦函数  |
| `void assert(int test);`|  测试一个条件并可能使程序终止 |
| `double atan(double x);`|  反正切函数 |
| `double atan2(double y, double x);`|  计算Y/X的反正切值 |
|`int atexit(atexit_t func);` |  注册终止函数 |
|`double atof(const char *nptr);` |  把字符串转换成浮点数 |
|`int atoi(const char *nptr);` |  把字符串转换成整型数 |
|`long atol(const char *nptr);` | 把字符串转换成长整型数  |


# 1. abort
## 1.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void abort(void);` |  异常终止一个进程 |

> **注意：** `abort()` 函数用于终止当前程序的执行。当程序调用 `abort()` 函数时，它会立即退出，并生成一个错误信号，通知操作系统程序非正常终止。如果程序已经打开了一些文件或句柄，但尚未关闭它们，则这些资源可能无法被正确释放。
## 1.2 演示示例

```c
#include <stdio.h>
#include <stdlib.h>

int main(void)
{
    printf("Calling abort()\n");
    abort();
    printf("already abort()"); // 这里永远也到不了
    return 0; 
}
```

## 1.3 运行结果

![](abort.png)

# 2. abs
## 2.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int abs(int i);` |  求整数的绝对值 |

## 2.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main(void)
{
    int number = -666;
  	printf("number: %d  absolute value: %d\n", number, 	abs(number));
  	return 0;
}
```
## 2.3 运行结果
![](abs.png)

# 3. absread 
## 3.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `int absread(int drive, int nsects, int sectno, void *buffer);`|  从drive指定的驱动器磁盘上，sectno指定的逻辑扇区号开始读取nsects个(最多64K个)扇区的内容，储存于buffer所指的缓冲区中。 |

## 3.2 演示示例
```c
#include <stdio.h>
#include <conio.h>
#include <process.h>
#include <dos.h>

int main(void)
{
  	int i, strt, ch_out, sector;
  	char buf[512];

  	printf("Insert a diskette into drive A and press any key\n");
  	getch();
  	sector = 0;
  	if (absread(0, 1, sector, &buf) != 0)
  	{
     	perror("Disk problem");
     	exit(1);
  	}
  	printf("Read OK\n");
  	strt = 3;
  	for (i=0; i<80; i++)
  	{
     	ch_out = buf[strt+i];
     	putchar(ch_out);
  	}
  	printf("\n");
  	return(0);
}
```
上述的代码实现了从 A 驱动器读取一个扇区的数据，并将其中一些字符输出到屏幕上。

- 首先提示用户插入一个软盘到 A 驱动器中。
- 然后读取 A 驱动器上第 0 个扇区的数据到缓冲区 buf 中。
- 接着检查读取是否成功。如果不成功，输出错误信息并退出程序。
- 最后将 buf 缓冲区中偏移量为 3 到偏移量为 82 的字符依次输出到屏幕上。

> **注意：** 程序中使用了一些 `DOS` 特定的函数，比如 `absread()` 和 `getch()`，可能不适用于其他操作系统或编译器环境。

# 4. abswrite
## 4.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int abswrite(int drive, int nsects, int sectno, void *buffer);`| 将指定内容写入磁盘上的指定扇区|

## 4.2 演示示例
```c
#include <dos.h>
#include <stdio.h>

unsigned char buff[512];

int main()
{
    int i;
    char c;
    printf("\nQuick Format 1.44MB\n");
    printf("Program by ChenQingyang.\n");
    printf("ALL DATA IN THE FLOPPY DISK WILL BE LOST!!\n");
    printf("\nInsert a diskette for drive A:\n");
    printf("and press ENTER when ready. . .");
    c=getchar();
    printf("\n\nCleaning FAT area. . .");
    buff[0]=0xf0;
    buff[1]=buff[2]=0xff;
    for (i=3;i<512;i++) 
        buff[i]=0;
    abswrite(0,1,1,buff);
    abswrite(0,1,10,buff);
    for (i=0;i<512;i++) 
        buff[i]=0;
    for (i=2;i<10;i++) 
        abswrite (0,1,i,buff);
    for (i=11;i<19;i++) 
        abswrite (0,1,i,buff);
    printf("\nCleaning ROOT area. . .");
    for (i=19;i<33;i++) 
        abswrite (0,1,i,buff);
    printf("\n\nQuickFormat Completed!\n");
}
```
上述代码是一个使用 `DOS` 命令格式化软盘的程序。它会提示用户输入软盘，然后清空软盘的FAT和根目录区域，并在完成后打印 `“QuickFormat Completed!”` 的信息。程序使用了 `<dos.h>` 和 `<stdio.h>` 头文件，其中包含了一些 `DOS` 和标准输入输出函数。

# 5. access
## 5.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int access(const char *filename, int amode);` | 确定文件的访问权限  |

## 5.2 演示示例
```c
#include <stdio.h>
#include <io.h>

int file_exists(char *filename);

int main(void)
{
    printf("Does students1.txt exist: %s\n",
    file_exists("students1.txt") ? "YES" : "NO");
    return 0;
}

int file_exists(char *filename)
{
    return (access(filename, 0) == 0);
}
```
## 5.3 运行结果
![](access.png)
![](access-1.png)

# 6. acos
## 6.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double acos(double x);` |  反余弦函数 |

## 6.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main(void)
{
    double result;
    double x = 0.5;

    result = acos(x); 
    printf("The arc cosine of %lf is %lf\n", x, result);
    return 0;
}
```
## 6.3 运行结果
![](acos.png)

# 7. allocmem
## 7.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `int allocmem(unsigned size, unsigned *seg);`|  分配DOS存储段 |

## 7.2 演示示例
```c
#include <dos.h>
#include <alloc.h>
#include <stdio.h>

int main(void)
{
    unsigned int size, segp;
    int stat;
    size = 64; /* (64 x 16) = 1024 bytes */
    stat = allocmem(size, &segp);
    if (stat == -1)
        printf("Allocated memory at segment: %x\n", segp);
    else
        printf("Failed: maximum number of paragraphs available is %u\n",stat);
    return 0;
}
```
在上述的示例代码，
- 首先调用了 `allocmem()` 函数来分配内存，其中传递了两个参数：`size` 表示请求的内存大小（以段为单位），这里设置为 `64` 段；`&segp` 表示返回的内存段地址将存储在此变量中。
- 如果成功分配内存，`allocmem()`函数将返回 `-1`，并打印出已分配内存的段地址；
- 否则，它将返回最大可用段数，并打印出失败的消息。

# 8. arc
## 8.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void arc(int x, int y, int stangle, int endangle, int radius);` | 画一弧线  |

## 8.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>

int main(void)
{
    /* request auto detection */
    int gdriver = DETECT, gmode, errorcode;

    int midx, midy;
    int stangle = 45, endangle = 135;
    int radius = 100;

    /* initialize graphics and local variables */
    char ch[] = "";
    initgraph(&gdriver, &gmode, ch);

    /* read result of initialization */
    errorcode = graphresult();    /* an error occurred */
    if (errorcode != grOk)
    {
        printf("Graphics error: %s\n", grapherrormsg(errorcode));
        printf("Press any key to halt:");
        getch();
        exit(1);    /* terminate with an error code */
    }

    midx = getmaxx() / 2;
    midy = getmaxy() / 2;
    setcolor(getmaxcolor());

    /* draw arc */
    arc(midx, midy, stangle, endangle, radius);

    /* clean up */
    getch();
    closegraph();
    return 0;
}

```
## 8.3 运行结果
![](arc.png)



# 9. asctime
## 9.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`char *asctime(const struct tm *tblock);` | 转换日期和时间为ASCII码  |

## 9.2 演示示例
```c
#include <stdio.h>
#include <string.h>
#include <time.h>

int main(void)
{
    struct tm t;
    char str[80];

    /* sample loading of tm structure  */

    t.tm_sec    = 1;  /* Seconds */
    t.tm_min    = 30; /* Minutes */
    t.tm_hour   = 9;  /* Hour */
    t.tm_mday   = 22; /* Day of the Month  */
    t.tm_mon    = 11; /* Month */
    t.tm_year   = 56; /* Year - does not include century */
    t.tm_wday   = 4;  /* Day of the week  */
    t.tm_yday   = 0;  /* Does not show in asctime  */
    t.tm_isdst  = 0;  /* Is Daylight SavTime; does not show in asctime */

    /* converts structure to null terminated
    string */

    strcpy(str, asctime(&t));
    printf("%s\n", str);

    return 0;
}
```
## 9.3 运行结果
![](asctime.png)

# 10. asin
## 10.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double asin(double x);`| 反正弦函数  |

## 10.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main(void)
{
    double result;
    double x = 0.5;

    result = asin(x);
    printf("The arc sin of %lf is %lf\n", x, result);
    return(0);
}
```
## 10.3 运行结果
![](asin.png)

# 11. assert
## 11.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `void assert(int test);`|  测试一个条件并可能使程序终止 |

## 11.2 演示示例
```c
#include <assert.h>
#include <stdio.h>
#include <stdlib.h>

struct ITEM {
   int key;
   int value;
};

/* add item to list, make sure list is not null */
void additem(struct ITEM *itemptr) {
    assert(itemptr != NULL);
   /* add item to list */
}

int main(void)
{
    additem(NULL);
    return 0;
}
```
## 11.3 运行结果
![](assert.png)

# 12. atan
## 12.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double atan(double x);`|  反正切函数 |

## 12.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main(void)
{
    double result;
    double x = 0.5;

    result = atan(x);
    printf("The arc tangent of %lf is %lf\n", x, result);
    return(0);
}
```
## 12.3 运行结果
![](atan.png)


# 13. atan2
## 13.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double atan2(double y, double x);`|  计算Y/X的反正切值 |

## 13.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main(void)
{
	double result;
	double x = 60.0, y = 30.0;

	result = atan2(y, x);
	printf("The arc tangent ratio of %lf is %lf\n", (y / x), result);
	return 0;
}
```
## 13.3 运行结果
![](atan2.png)

# 14. atexit
## 14.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int atexit(atexit_t func);` |  注册终止函数 |

## 14.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>

void exit_fn1(void)
{
    printf("Exit function #1 called\n");
}

void exit_fn2(void)
{
    printf("Exit function #2 called\n");
}

int main(void)
{
    /* post exit function #1 */
    atexit(exit_fn1);
    /* post exit function #2 */
    atexit(exit_fn2);
    return 0;
}
```
## 14.3 运行结果
![](atexit.png)

# 15. atof
## 15.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double atof(const char *nptr);` |  把字符串转换成浮点数 |

## 15.2 演示示例
```c
#include <stdlib.h>
#include <stdio.h>

int main(void)
{
    float f;
    char *str = "1234.5678";

    f = atof(str);
    printf("string = %s float = %f\n", str, f);
    return 0;
}
```
## 15.3 运行结果
![](atof.png)

# 16. atoi
## 16.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int atoi(const char *nptr);` |  把字符串转换成整型数 |

## 16.2 演示示例
```c
#include <stdlib.h>
#include <stdio.h>

int main(void)
{
    int n;
    char *str = "1234.5678";

    n = atoi(str);
    printf("string = %s integer = %d\n", str, n);
    return 0;
}
```
## 16.3 运行结果
![](atoi.png)

# 17. atol
## 17.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`long atol(const char *nptr);` | 把字符串转换成长整型数  |

## 17.2 演示示例
```c
#include <stdlib.h>
#include <stdio.h>

int main(void)
{
    long l;
    char *lstr = "87654321";

    l = atol(lstr);
    printf("string = %s integer = %ld\n", lstr, l);
    return(0);
}
```
## 17.3 运行结果
![](atol.png)


# 参考
1. [[API Reference Document]](https://www.apiref.com/c-zh/index.htm)

