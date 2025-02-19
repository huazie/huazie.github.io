---
title: C语言函数大全--h开头的函数
date: 2023-04-11 22:33:56
updated: 2025-02-19 21:22:32
categories:
  - 开发语言-C
tags:
  - C语言函数大全
  - h开头的函数
---

![](/images/cplus-logo.png)

# 总览
| 函数声明 |  函数功能  |
|:--|:--|
| `double hypot(double x, double y);`| 计算直角三角形的斜边长（double）  |
| `float hypotf (float x, float y);`| 计算直角三角形的斜边长（float）  |
| `long double hypot(long double x, long double y);`| 计算直角三角形的斜边长（long double）  |
|`#define HUGE_VAL _HUGE` |  正浮点常量表达式（double），这些表达式与浮点函数和运算符在溢出时返回的值相比较 |
|`#define HUGE_VALF __INFF` |  正浮点常量表达式（float），这些表达式与浮点函数和运算符在溢出时返回的值相比较 |
|`#define HUGE_VALL __INFL` | 正浮点常量表达式（long double），这些表达式与浮点函数和运算符在溢出时返回的值相比较  |
|`void harderr(int (*fptr)());` |  建立一个硬件错误处理程序 |
|`void hardresume(int rescode);`|  硬件错误处理函数 |
|`void hardretn(int rescode);` |  硬件错误处理函数 |
|`void highvideo(void);` | 选择高亮度文本字符  |
|`int hcreate(size_t nel);` | 根据条目数创建哈希表。  |
|`int hcreate_r(size_t nel, struct hsearch_data *htab);` | 根据条目数及其描述创建哈希表。|
|`ENTRY *hsearch(ENTRY item, ACTION action);` | 添加或搜索哈希条目。  |
|`int hsearch_r (ENTRY item, ACTION action, ENTRY ** retval, struct hsearch_data * htab )` | 搜索哈希表。  |
|`void hdestroy(void);` | 销毁哈希表，释放用于创建哈希表的内存。  |
|`void hdestroy_r(struct hsearch_data *htab);` | 销毁哈希表，释放指定哈希表所占用的内存。   |

# 1. hypot，hypotf，hypotl
## 1.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double hypot(double x, double y);`| 计算直角三角形的斜边长（double）  |
| `float hypotf (float x, float y);`| 计算直角三角形的斜边长（float）  |
| `long double hypot(long double x, long double y);`| 计算直角三角形的斜边长（long double）  |

## 1.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main(void)
{
   double result, x = 3.0, y = 4.0;
   result = hypot(x, y);

   float resultf, xf = 3.0, yf = 4.0;
   resultf = hypotf(xf, yf);

   long double resultL, xL = 3.0, yL = 4.0;
   resultL = hypotl(xL, yL);

   printf("The hypotenuse of a right triangle whose legs are %lf and %lf is %lf\n", x, y, result);
   printf("The hypotenuse of a right triangle whose legs are %f and %f is %f\n", xf, yf, resultf);
   printf("The hypotenuse of a right triangle whose legs are %Lf and %Lf is %Lf\n", xL, yL, resultL);

   return 0;
}
```
## 1.3 运行结果
![](hypot.png)

# 2. HUGE_VAL，HUGE_VALF，HUGE_VALL
## 2.1 函数说明
| 宏定义 |  宏描述  |
|:--|:--|
|`#define HUGE_VAL _HUGE` |  正浮点常量表达式（double），这些表达式与浮点函数和运算符在溢出时返回的值相比较 |
|`#define HUGE_VALF __INFF` |  正浮点常量表达式（float），这些表达式与浮点函数和运算符在溢出时返回的值相比较 |
|`#define HUGE_VALL __INFL` | 正浮点常量表达式（long double），这些表达式与浮点函数和运算符在溢出时返回的值相比较  |

## 2.2 演示示例
```c
#include<stdio.h>
#include<math.h>
int main()
{
    double result = 1.0/0.0;
    printf("1.0/0.0 = %lf\n", result);
    if (result == HUGE_VAL)
        puts("1.0/0.0 == HUGE_VAL\n");

    float resultf = 1.0f/0.0f;
    printf("1.0f/0.0f = %f\n", resultf);
    if (resultf == HUGE_VALF)
        puts("1.0f/0.0f == HUGE_VALF\n");

    long double resultL = 1.0L/0.0L;
    printf("1.0L/0.0L = %Lf\n", resultL);
    if (resultL == HUGE_VALL)
        puts("1.0L/0.0L == HUGE_VALL\n");

    return 0;  
}

```
## 2.3 运行结果
![](HUGE_VAL.png)

# 3. harderr，hardresume，hardretn
## 3.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void harderr(int (*fptr)());` |  建立一个硬件错误处理程序 |
|`void hardresume(int rescode);`|  硬件错误处理函数 |
|`void hardretn(int rescode);` |  硬件错误处理函数 |

## 3.2 演示示例
```c
#include <stdio.h>
#include <conio.h>
#include <dos.h>
#define IGNORE  0
#define RETRY   1
#define ABORT   2
int buf[500];

//定义捕获磁盘问题的错误消息
static char *err_msg[] = {
    "write protect",
    "unknown unit",
    "drive not ready",
    "unknown command",
    "data error (CRC)",
    "bad request",
    "seek error",
    "unknown media type",
    "sector not found",
    "printer out of paper",
    "write fault",
    "read fault",
    "general failure",
    "reserved",
    "reserved",
    "invalid disk change"
};

int error_win(char *msg)
{
    int retval;
    cputs(msg);

    // 提示用户按键中止、重试、忽略
    while(1)
    {
        retval= getch();
        if (retval == 'a' || retval == 'A')
        {
            retval = ABORT;
            break;
        }
        if (retval == 'r' || retval == 'R')
        {
            retval = RETRY;
            break;
        }
        if (retval == 'i' || retval == 'I')
        {
            retval = IGNORE;
            break;
        }
    }

    return(retval);
}

/*
    pragma warn-par 减少了由于处理程序未使用参数errval、bp 和 si而产生的警告。
*/
#pragma warn -par

int handler(int errval,int ax,int bp,int si)
{
    static char msg[80];
    unsigned di;
    int drive;
    int errorno;
    di= _DI;

    // 如果这不是磁盘错误，那么是另一个设备出现故障
    if (ax < 0)
    {
        error_win("Device error");
        // 返回到直接请求中止的程序
        hardretn(ABORT);
    }
    // 否则就是磁盘错误
    drive = ax & 0x00FF;
    errorno = di & 0x00FF;
    sprintf(msg, "Error: %s on drive %c\r\nA)bort, R)etry, I)gnore: ", err_msg[errorno], 'A' + drive);
    // 通过dos中断0x23返回程序，并由用户输入中止、重试或忽略。
    hardresume(error_win(msg));
    return ABORT;
}

#pragma warn +par

int main(void)
{
    // 在硬件问题中断上安装我们的处理程序
    harderr(handler);
    clrscr();
    printf("Make sure there is no disk in drive A:\n");
    printf("Press any key ....\n");
    getch();
    printf("Trying to access drive A:\n");
    printf("fopen returned %p\n", fopen("A:temp.dat", "w"));
    return 0;
}
```

上述程序是一个基于**DOS**环境的磁盘错误处理示例。在磁盘操作出现错误时，向用户显示具体的错误消息，并提供 **“中止”**、**“重试”** 和 **“忽略”** 三种处理选项，根据用户的选择进行相应的处理。


# 4. highvideo
## 4.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void highvideo(void);` | 选择高亮度文本字符  |

## 4.2 演示示例
```c
#include <stdio.h>
#include <conio.h>

int main(void)
{
    clrscr();

    lowvideo();
    cprintf("Low Intensity text\r\n");
    highvideo();
    gotoxy(1,2);
    cprintf("High Intensity Text\r\n");

    return 0;
}
```

上述利用 `<conio.h>` 头文件中的函数实现特定的控制台文本显示效果。
- 首先通过 `clrscr()` 函数清空控制台窗口。
- 然后调用 `lowvideo()` 将文本显示设置为低亮度，使用 `cprintf()` 输出低亮度的 `"Low Intensity text"` 文本并换行。
- 接着调用 `highvideo()` 把文本显示切换为高亮度，利用 `gotoxy(1, 2)` 将光标定位到第 1 列第 2 行，再用 `cprintf()` 输出高亮度的 `"High Intensity Text"` 文本并换行，以此直观展示控制台中低亮度和高亮度两种不同的文本显示效果。

> 不过需要留意，`conio.h` 并非标准 **C** 库的一部分，它主要在像 **Turbo C** 这类旧的编译器中使用，而在现代开发环境里可能不被支持。

# 5. hcreate，hcreate_r
## 5.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int hcreate(size_t nel);` | 根据条目数创建哈希表。  |
|`int hcreate_r(size_t nel, struct hsearch_data *htab);` | 根据条目数及其描述创建哈希表。|

> **入参：**
> - **net ：** 哈希表中允许的最大项数。
> - **htab ：** 哈希表的结构体数据。

> **返回值：**
> - 如果操作成功，则返回一个非零值；
> - 如果操作失败，则返回 **0** 并将 **errno** 设置为一个值。

## 5.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <search.h>

char *data[] = { "alpha", "bravo", "charlie", "delta",
     "echo", "foxtrot", "golf", "hotel", "india", "juliet",
     "kilo", "lima", "mike", "november", "oscar", "papa",
     "quebec", "romeo", "sierra", "tango", "uniform",
     "victor", "whisky", "x-ray", "yankee", "zulu"
};

int main(void)
{
    ENTRY e, *ep;
    int i;

    hcreate(30);

    for (i = 0; i < 24; i++) {
        e.key = data[i];
        // 数据只是一个整数，而不是指向某个东西的指针
        e.data = (void *) i;
        ep = hsearch(e, ENTER);
        // 这里不应该有失败场景
        if (ep == NULL) {
            fprintf(stderr, "entry failed\n");
            exit(EXIT_FAILURE);
        }
    }

    for (i = 22; i < 26; i++) {
        // 从表中打印两个条目，并显示其中两个不在表中
        e.key = data[i];
        ep = hsearch(e, FIND);
        printf("%9.9s -> %9.9s:%d\n", e.key, ep ? ep->key : "NULL", ep ? (int)(ep->data) : 0);
    }
    hdestroy();
    exit(EXIT_SUCCESS);
}
```
# 6. hsearch，hsearch_r
## 6.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`ENTRY *hsearch(ENTRY item, ACTION action);` | 添加或搜索哈希条目。  |
|`int hsearch_r (ENTRY item, ACTION action, ENTRY ** retval, struct hsearch_data * htab )` | 搜索哈希表。  |

> **注意：** 
> ``hsearch`` 和 ``hsearch_r``  函数根据指定的操作在哈希表中搜索条目。如果操作为 FIND，则仅执行搜索操作。如果操作为 ENTER，则未找到的条目将添加到哈希表中。``hsearch_r`` 函数与 ``hsearch`` 函数的不同之处在于，指向找到的项的指针以 ``*retval`` 形式返回，而不是作为函数结果。


> **入参：**
> - **item：** 要搜索的哈希表条目。
> - **action：** 功能操作。``ENTER`` 表示已添加条目，``FIND`` 表示已搜索条目。
> - **retval：** 指向找到的项的指针。
> - **htab ：** 哈希表的结构体数据。

> **hsearch 函数返回值：**
> - 如果操作成功，则返回指向哈希表的指针。

> **hsearch_r 函数返回值：**
> - 如果操作成功，则返回一个非零值；
> - 如果操作失败，则返回 0。


## 6.2 演示示例
参考 5.2

# 7. hdestroy，hdestroy_r
## 7.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void hdestroy(void);` | 销毁哈希表，释放用于创建哈希表的内存。  |
|`void hdestroy_r(struct hsearch_data *htab);` | 销毁哈希表，释放指定哈希表所占用的内存。   |

## 7.2 演示示例
参考 5.2

# 8. htonl, htons
## 8.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`uint32_t htonl(uint32_t hostlong); ` | 将 uint32_t（32位整数，如IPv4地址）从主机字节序转为网络字节序。  |
|`uint16_t htons(uint16_t hostshort); ` |  将 uint16_t（16位整数，如端口号）从主机字节序转为网络字节序（大端）。  |


## 8.2 演示示例

```c
#include <stdio.h>
#include <stdint.h>     // 提供 uint16_t, uint32_t 类型
#include <arpa/inet.h>  // Linux/macOS 头文件（网络字节序转换）

int main() {
    // 示例1：转换16位端口号（host to network short）
    uint16_t host_port = 0x1234;  // 主机字节序的端口（假设是小端模式）
    uint16_t net_port = htons(host_port);
    printf("主机端口: 0x%04x -> 网络端口: 0x%04x\n", host_port, net_port);

    // 示例2：转换32位IP地址（host to network long）
    uint32_t host_ip = 0x12345678;  // 主机字节序的IP地址
    uint32_t net_ip = htonl(host_ip);
    printf("主机IP: 0x%08x -> 网络IP: 0x%08x\n", host_ip, net_ip);

    return 0;
}
```

# 参考
1. [\[API Reference Document\]](https://www.apiref.com/c-zh/t_alph_h.htm)
2. [\[highvideo\]](https://baike.baidu.com/item/highvideo/6436214?fr=aladdin)
3. [\[hcreate,hsearch,hdestroy,hcreate_r,hsearch_r,hdestroy_r\]](https://linux.die.net/man/3/hcreate)
4. [\[UTILS-标准C库\]](https://device.harmonyos.com/cn/docs/documentation/apiref/utils-0000001055308029#ZH-CN_TOPIC_0000001055308029__gafb18cb23be808765135c3aa903df62fd)