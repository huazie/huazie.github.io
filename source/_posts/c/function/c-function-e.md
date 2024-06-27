---
title: C语言函数大全--e开头的函数
date: 2023-03-30 21:15:08
updated: 2024-06-08 23:27:10
categories:
  - 开发语言-C
tags:
  - C语言函数大全
  - e开头的函数
---

![](/images/cplus-logo.png)


# 总览
| 函数声明 |  函数功能  |
|:--|:--|
|`char ecvt(double value, int ndigit, int *decpt, int *sign);` | 把一个双精度浮点型数转换为字符串，转换结果中不包括十进制小数点 |
| `void ellipse(int x, int y, int stangle, int endangle, int xradius, int yradius);`| 画一段椭圆线 |
| `int eof(int *handle);`| 检测文件结束  |
|`int execl(const char *pathname, const char *arg0, ... const char *argn, NULL);` | 载入并运行其它程序  |
|`int execlp(char *pathname, char *arg0, ... const char *argn, NULL);` |  载入并运行其它程序 |
| `int execlpe(const char *pathname, const char *arg0, ... const char *argn, NULL, const char *const *envp);`|  载入并运行其它程序   |
|`int execv(const char *pathname, char *const *argv);` | 载入并运行其它程序  |
|`int execve(const char *pathname, char *const argv[], char *const envp[]);` | 载入并运行其它程序  |
|`int execvp(const char *pathname, char *const argv[]);` | 载入并运行其它程序  |
|`int execvpe(const char *pathname, char *const argv[], char *const envp[]);` | 载入并运行其它程序  |
| `void exit(int status);`|  终止程序 |
| `double exp(double x);`|  计算 x 的基数e指数（double） |
| `float expf(float x);`|  计算 x 的基数e指数（float）  |
| `long double expl(long double x);`|  计算 x 的基数e指数（long double） |
| `double exp2(double x);`|  计算 x 的基数为2的指数（double） |
| `float exp2f(float x);`|  计算 x 的基数为2的指数（float）  |
| `long double exp2l(long double x);`|  计算 x 的基数为2的指数（long double） |
| `double expm1 (double x);`| 计算 e 的 x 次方 减 1，即 (e^x) - 1  （double）|
| `float expm1f (float x);`| 计算 e 的 x 次方 减 1，即 (e^x) - 1 （float） |
| `long double expm1l (long double x);`| 计算 e 的 x 次方 减 1，即 (e^x) - 1  （long double）|
|`double erf (double x);` |  计算 x 的 [误差函数](https://baike.baidu.com/item/%E8%AF%AF%E5%B7%AE%E5%87%BD%E6%95%B0/5890875?fr=aladdin)（double） |
|`float erff (float x);` |  计算 x 的 [误差函数](https://baike.baidu.com/item/%E8%AF%AF%E5%B7%AE%E5%87%BD%E6%95%B0/5890875?fr=aladdin)（float） |
|`long double erfl (long double x);` |  计算 x 的 [误差函数](https://baike.baidu.com/item/%E8%AF%AF%E5%B7%AE%E5%87%BD%E6%95%B0/5890875?fr=aladdin)（long double） |
|`double erfc (double x);` |  计算 x 的互补误差函数（double） |
|`float erfcf (float x);` |  计算 x 的互补误差函数（float） |
|`long double erfcl (long double x);` |  计算 x 的互补误差函数（long double） |

# 1. ecvt
## 1.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`char ecvt(double value, int ndigit, int *decpt, int *sign);` | 把一个双精度浮点型数转换为字符串，转换结果中不包括十进制小数点 |

> **value** ： 待转换的双精度浮点数。
> **ndigit** ：存储的有效数字位数。这个函数存储最多 **ndigit** 个数字值作为一个字符串，并添加一个结束符(**'\0'**)，如果 **value** 中的数字个数超过 **ndigit**，低位数字被舍入。如果少于 **ndigit** 个数字，该字符串用 **0**填充。
> **decpt** ：指出给出小数点位置的整数值, 它是从该字符串的开头位置计算的。**0** 或负数指出小数点在第一个数字的左边。
> **sign**   ：指出一个指出转换的数的符号的整数。如果该整数为 **0**,这个数为正数,否则为负数。

## 1.2 演示示例
```c
#include <stdlib.h>
#include <stdio.h>
#include <conio.h>

int main()
{
    char *string;
    double value;
    int decpt, sign;

    int ndigit = 10;
    value = 9.876;
    string = ecvt(value, ndigit, &decpt, &sign);
    printf("string = %-16s decpt = %d sign = %d value = %lf\n", string, decpt, sign, value);

    value = -123.45;
    ndigit= 15;
    string = ecvt(value, ndigit, &decpt, &sign);
    printf("string = %-16s decpt = %d sign = %d value = %lf\n", string, decpt, sign, value);

    value = 0.6789e5; /* 科学记数法 scientific notation */
    ndigit = 5;
    string = ecvt(value, ndigit, &decpt, &sign);
    printf("string = %-16s decpt = %d sign = %d value = %lf\n", string, decpt, sign, value);

    return 0;
}
```
## 1.3 运行结果
![](ecvt.png)

# 2. ellipse
## 2.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `void ellipse(int x, int y, int stangle, int endangle, int xradius, int yradius);`| 画一段椭圆线 |

> 以 **(x, y)** 为中心，**xradius**、**yradius** 为 **x 轴** 和 **y 轴** 半径，从角 **stangle** 开始，**endangle** 结束，画一段椭圆线。当**stangle=0，endangle=360** 时，画出一个完整的椭圆

## 2.2 演示示例
```c
#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>

int main(void)
{
    // request auto detection
    int gdriver = DETECT, gmode, errorcode;
    int midx, midy;
    int stangle = 0, endangle = 360, midangle = 180;
    int xradius = 100, yradius = 50;

    // initialize graphics, local variables
    initgraph(&gdriver, &gmode, "");

    // read result of initialization
    errorcode = graphresult();
    if (errorcode != grOk) // an error occurred
    {
        printf("Graphics error: %s\n", grapherrormsg(errorcode));
        printf("Press any key to halt:");
        getch();
        exit(1);
    }

    midx = getmaxx() / 2;
    midy = getmaxy() / 2;
    setcolor(getmaxcolor());

    // draw ellipse
    ellipse(midx, 50, midangle, endangle, xradius, yradius);

    ellipse(midx, midy, stangle, endangle, xradius, yradius);

    ellipse(midx, getmaxy() - 50, stangle, midangle, xradius, yradius);

    getch();
    closegraph();
    return 0;
}

```
## 2.3 运行结果
![](ellipse.png)


# 3. eof
## 3.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `int eof(int *handle);`| 检测文件结束  |

## 3.2 演示示例
```c
#include <sys\stat.h>
#include <string.h>
#include <stdio.h>
#include <fcntl.h>
#include <io.h>

int main(void)
{
    int handle;
    char msg[] = "This is a test";
    char ch;

    // create a file 
    handle = open("STU.FIL", O_CREAT | O_RDWR, S_IREAD | S_IWRITE);

    // write some data to the file 
    write(handle, msg, strlen(msg));

    // seek to the beginning of the file 
    lseek(handle, 0L, SEEK_SET);

    // reads chars from the file until hit EOF
    do
    {
        read(handle, &ch, 1);
        printf("%c", ch);
    } while (!eof(handle));

    close(handle);
    return 0;
}
```
## 3.3 运行结果
![](eof.png)

# 4. execl
## 4.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int execl(const char *pathname, const char *arg0, ... const char *argn, NULL);` | 载入并运行其它程序  |

>**注意：** **execl** 函数，其后缀 `l` 代表 **list**，也就是参数列表的意思。第一个参数 **path** 字符指针指向要执行的文件路径， 接下来的参数代表执行该文件时传递的参数列表：`argv[0],argv[1]...` ，最后一个参数须用空指针 **NULL** 作结束。

## 4.2 演示示例
### 4.2.1 SubTest.c
```c
#include <stdio.h>
int main(int argc, char *argv[])
{
  printf("exec %s, Hello, %s", argv[0], argv[1]);
  return 0;
} 
```
### 4.2.2 Test.c
```c
#include <process.h>
#include <stdio.h>
#include <errno.h>

void main(int argc, char *argv[])
{
    int i;
    printf("Command line arguments:\n");
    for (i=0; i<argc; i++)
        printf("[%d] : %s\n", i, argv[i]);
    printf("exec SubTest with subargv ...\n");
    // 成功则不返回值， 失败返回-1， 失败原因存于errno中，可通过perror()打印
    // 第一个参数需要执行文件的全路径，这里写直接文件名，是因为和当前源码在同一目录中
    int result = execl("SubTest.exe", argv[0], "Huazie" , NULL);
    // 执行成功，这里不会执行到
    printf("result = %d\n", result);
    perror("exec error");
    exit(1);
}
```
## 4.3 运行结果
执行失败：
![](execl.png)
执行成功：
![](execl-1.png)

# 5. execle
## 5.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int execle(const char *pathname, const char *arg0, ... const char *argn, NULL, const char *const *envp);` | 载入并运行其它程序  |

> **注意：** **execl** 函数是用来执行参数 **path** 字符串所代表的文件路径。接下来的参数代表执行该文件时传递过去的 `argv[0], argv[1]…`，并且倒数第二个参数必须用空指针 **NULL** 作结束，最后一个参数为 **环境变量**。
## 5.2 演示示例
### 5.2.1 SubEnvTest.c
```c
#include <process.h>
#include <processenv.h>
#include <stdio.h>

int main(int argc, char *argv[], char **envp)
{
    printf("SubEnvTest Command line arguments:\n");
    for (int i=0; i<argc; i++)
        printf("[%d] : %s\n", i, argv[i]);
    printf("exec %s, Hello, %s\n", argv[0], argv[1]);

    for (int i = 0; envp[i] != NULL; i++)
    {
        printf("%s\n", envp[i]);
    }
    // while(*envp != NULL)
    //     printf("  %s\n", *(envp++));

    return 0;
} 
```

### 5.2.2 Test.c

```c
#include <process.h>
#include <stdio.h>
#include <errno.h>

void main(int argc, char *argv[])
{
    int i;
    printf("Command line arguments:\n");
    for (i=0; i<argc; i++)
        printf("[%d] : %s\n", i, argv[i]);
    printf("exec SubEnvTest with subargv ...\n");
    const char *envp[] = {"AUTHOR=Huazie", "DATE=2023-03-28", NULL}; // 环境变量
    
    // 成功则不返回值， 失败返回-1， 失败原因存于errno中，可通过perror()打印
    // 第一个参数需要执行文件的全路径，这里直接写文件名，是因为和当前源码在同一目录中
    int result = execle("SubEnvTest.exe", argv[0], "Huazie" , NULL, envp);
    // 执行成功，这里不会执行到
    printf("result = %d\n", result);
    perror("exec error");
    exit(1);
}
```

## 5.3 运行结果
![](execle.png)

# 6. execlp
## 6.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int execlp(char *pathname, char *arg0, ... const char *argn, NULL);` |  载入并运行其它程序 |

>**注意：** **execlp** 函数会从 **PATH** 环境变量所指的目录中查找符合参数 **pathname** 的文件名，找到后便执行该文件，然后将第二个以后的参数当做该文件的`arg0, arg1, …`，最后一个参数必须用 空指针 **NULL** 作结束。
## 6.2 演示示例
```c
#include <process.h>
#include <stdio.h>
#include <errno.h>

void main(int argc, char *argv[])
{
    int i;
    printf("Command line arguments:\n");
    for (i=0; i<argc; i++)
        printf("[%d] : %s\n", i, argv[i]);
    printf("exec java with subargv ...\n");
    // 成功则不返回值， 失败返回-1， 失败原因存于errno中，可通过perror()打印
    int result = execlp("java.exe", "java", "-version", NULL);
    // 执行成功，这里不会执行到
    printf("result = %d\n", result);
    perror("exec error");
    exit(1);
}
```
## 6.3 运行结果
![](execlp.png)

# 7. execlpe
## 7.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `int execlpe(const char *pathname, const char *arg0, ... const char *argn, NULL, const char *const *envp);`|  载入并运行其它程序   |

>**注意：** **execlp** 函数会从 **PATH** 环境变量所指的目录中查找符合参数 **pathname** 的文件名，找到后便执行该文件，然后将第二个以后的参数当做该文件的`arg0, arg1, …`，其中倒数第二个参数必须用 空指针 **NULL** 作结束，最后一个参数为 **环境变量**。
## 7.2 演示示例
### 7.2.1 SubEnvTest.c
参考 **5.2.1** 的 **SubEnvTest.c**
### 7.2.2 Test.c
```c
#include <process.h>
#include <stdio.h>
#include <errno.h>

void main(int argc, char *argv[])
{
    int i;
    printf("Command line arguments:\n");
    for (i=0; i<argc; i++)
        printf("[%d] : %s\n", i, argv[i]);
    printf("exec SubEnvTest with subargv ...\n");
    const char *envp[] = {"AUTHOR=Huazie", "DATE=2023-03-28", NULL}; // 环境变量
    // 成功则不返回值， 失败返回-1， 失败原因存于errno中，可通过perror()打印
    int result = execlpe("SubEnvTest", argv[0], "Huazie", NULL, envp);
    // 执行成功，这里不会执行到
    printf("result = %d\n", result);
    perror("exec error");
    exit(1);
}
```
## 7.3 运行结果
![](execlpe.png)

# 8. execv
## 8.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int execv(const char *pathname, char *const *argv);` | 载入并运行其它程序  |

> **注意：execv** 函数用来运行參数 **pathname** 字符串所指向的程序，第二个参数 **argv** 为參数列表【该数组的最后一个元素必须是空指针 **NULL**】。
## 8.2 演示示例
### 8.2.1 SubTest.c
参考 **4.2.1** 的 **SubTest.c**
### 8.2.2 Test.c
```c
#include <process.h>
#include <stdio.h>
#include <errno.h>

void main(int argc, char *argv[])
{
    int i;
    printf("Command line arguments:\n");
    for (i=0; i<argc; i++)
        printf("[%d] : %s\n", i, argv[i]);
    printf("exec SubTest with subargv ...\n");
    char *const subargv[]  = {argv[0], "Huazie" , NULL}; // 参数列表
    // 成功则不返回值， 失败返回-1， 失败原因存于errno中，可通过perror()打印
    // 第一个参数需要执行文件的全路径，这里写直接文件名，是因为和当前源码在同一目录中
    int result = execv("SubTest.exe", subargv);
    // 执行成功，这里不会执行到
    printf("result = %d\n", result);
    perror("exec error");
    exit(1);
}
```
## 8.3 运行结果
![](execv.png)

# 9. execve
## 9.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int execve(const char *pathname, char *const argv[], char *const envp[]);` | 载入并运行其它程序  |

> **注意：execve** 函数用来运行參数 **pathname** 字符串所指向的程序，第二个参数 **argv** 为參数列表【该数组的最后一个元素必须是空指针 **NULL**】，最后一个参数为 **环境变量**。

## 9.2 演示示例
### 9.2.1 SubEnvTest.c
参考 **5.2.1** 的 **SubEnvTest.c**

### 9.2.2 Test.c
```c
#include <process.h>
#include <stdio.h>
#include <errno.h>

void main(int argc, char *argv[])
{
    int i;
    printf("Command line arguments:\n");
    for (i=0; i<argc; i++)
        printf("[%d] : %s\n", i, argv[i]);
    printf("exec SubEnvTest with subargv ...\n");
    char *const subargv[]  = {argv[0], "Huazie" , NULL}; // 参数列表
    char *const envp[] = {"AUTHOR=Huazie", "DATE=2023-03-28", NULL}; // 环境变量
    // 成功则不返回值， 失败返回-1， 失败原因存于errno中，可通过perror()打印
    // 第一个参数需要执行文件的全路径，这里写直接文件名，是因为和当前源码在同一目录中
    int result = execve("SubEnvTest.exe", subargv, envp);
    // 执行成功，这里不会执行到
    printf("result = %d\n", result);
    perror("exec error");
    exit(1);
}
```
## 9.3 运行结果
![](execve.png)

# 10. execvp
## 10.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int execvp(const char *pathname, char *const argv[]);` | 载入并运行其它程序  |

> **注意：execvp** 函会从 **PATH** 环境变量所指的目录中查找符合参数 **pathname** 的文件名，找到后便执行该文件，第二个参数 **argv** 为參数列表【该数组的最后一个元素必须是空指针 **NULL**】。

## 10.2 演示示例
```c
#include <process.h>
#include <stdio.h>
#include <errno.h>

void main(int argc, char *argv[])
{
    int i;
    printf("Command line arguments:\n");
    for (i=0; i<argc; i++)
        printf("[%d] : %s\n", i, argv[i]);
    printf("exec java with subargv ...\n");
    char *const subargv[]  = {"java", "-version" , NULL}; // 参数列表
    // 成功则不返回值， 失败返回-1， 失败原因存于errno中，可通过perror()打印
    int result = execvp("java.exe", subargv);
    // 执行成功，这里不会执行到
    printf("result = %d\n", result);
    perror("exec error");
    exit(1);
}
```
## 10.3 运行结果
![](execvp.png)

# 11. execvpe
## 11.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int execvpe(const char *pathname, char *const argv[], char *const envp[]);` | 载入并运行其它程序  |

> **注意：execvpe** 函会从 **PATH** 环境变量所指的目录中查找符合参数 **pathname** 的文件名，找到后便执行该文件，第二个参数 **argv** 为參数列表【该数组的最后一个元素必须是空指针 **NULL**】，最后一个参数为 **环境变量**。

## 11.2 演示示例
### 11.2.1 SubEnvTest.c
参考 **5.2.1** 的 **SubEnvTest.c**

### 11.2.2 Test.c
```c
#include <process.h>
#include <stdio.h>
#include <errno.h>

void main(int argc, char *argv[])
{
    int i;
    printf("Command line arguments:\n");
    for (i=0; i<argc; i++)
        printf("[%d] : %s\n", i, argv[i]);
    printf("exec SubEnvTest with subargv ...\n");
    char *const subargv[]  = {argv[0], "Huazie" , NULL}; // 参数列表
    char *const envp[] = {"AUTHOR=Huazie", "DATE=2023-03-28", NULL}; // 环境变量
    // 成功则不返回值， 失败返回-1， 失败原因存于errno中，可通过perror()打印
    int result = execvpe("SubEnvTest.exe", subargv, envp);
    // 执行成功，这里不会执行到
    printf("result = %d\n", result);
    perror("exec error");
    exit(1);
}
```
## 11.3 运行结果
![](execvpe.png)

# 12. exit
## 12.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `void exit(int status);`|  终止程序 |

> **注意： exit** 函数通常是用来终结程序用的，使用后程序自动结束，跳回操作系统。`exit(0)` 表示程序正常退出【相当于主函数 `return 0;`】，`exit⑴` 表示程序异常退出【相当于主函数 `return 1;`】。

## 12.2 演示示例
```c
#include <stdlib.h>
#include <stdio.h>

int main(void)
{
    int status;

    printf("Enter either 1 or 2\n");
    status = getchar();
    /* Sets DOS errorlevel  */
    exit(status - '0');

    /* Note: this line is never reached */
    printf("this line is never reached!");
    return 0;
}
```
## 12.3 运行结果
![](exit.png)

# 13. exp，expf，expl
## 13.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double exp(double x);`|  计算 x 的基数e指数（double） |
| `float expf(float x);`|  计算 x 的基数e指数（float）  |
| `long double expl(long double x);`|  计算 x 的基数e指数（long double） |

## 13.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main(void)
{
    double result, x = 4.0;
    result = exp(x);

    float resultf, xf = 4.0;
    resultf = expf(xf);

    long double resultL, xL = 4.0;
    resultL = expl(xL);

    printf("\n'e' raised to the power of %lf (e ^ %lf) = %.20lf\n", x, x, result);
    printf("\n'e' raised to the power of %f (e ^ %f) = %.20f\n", xf, xf, resultf);
    printf("\n'e' raised to the power of %Lf (e ^ %Lf) = %.20Lf\n", xL, xL, resultL);

    return 0;
}
```
## 13.3 运行结果
![](exp.png)


# 14. exp2，exp2f，exp2l
## 14.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double exp2(double x);`|  计算 x 的基数为2的指数（double） |
| `float exp2f(float x);`|  计算 x 的基数为2的指数（float）  |
| `long double exp2l(long double x);`|  计算 x 的基数为2的指数（long double） |

## 14.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main(void)
{
    double result, x = 4.5;
    result = exp2(x);

    float resultf, xf = 4.5;
    resultf = exp2f(xf);

    long double resultL, xL = 4.5;
    resultL = exp2l(xL);

    printf("\n'2' raised to the power of %lf (2 ^ %lf) = %.20lf\n", x, x, result);
    printf("\n'2' raised to the power of %f (2 ^ %f) = %.20f\n", xf, xf, resultf);
    printf("\n'2' raised to the power of %Lf (2 ^ %Lf) = %.20Lf\n", xL, xL, resultL);

    return 0;
}
```
## 14.3 运行结果

![](exp2.png)

# 15. expm1，expm1f，expm1l
## 15.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double expm1 (double x);`| 计算 e 的 x 次方 减 1，即 (e^x) - 1  （double）|
| `float expm1f (float x);`| 计算 e 的 x 次方 减 1，即 (e^x) - 1 （float） |
| `long double expm1l (long double x);`| 计算 e 的 x 次方 减 1，即 (e^x) - 1  （long double）|

## 15.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main(void)
{
    double result, result1, x = 4.0;
    result = exp(x);
    result1 = expm1(x);

    float resultf, resultf1, xf = 4.0;
    resultf = expf(xf);
    resultf1 = expm1f(xf);

    long double resultL, resultL1, xL = 4.0;
    resultL = expl(xL);
    resultL1 = expm1l(xL);

    printf("\n'e' raised to the power of %lf (e ^ %lf) = %.20lf\n", x, x, result);
    printf("\n'e' raised to the power of %lf minus one (e ^ %lf - 1) = %.20lf\n", x, x, result1);

    printf("\n'e' raised to the power of %f (e ^ %f) = %.20f\n", xf, xf, resultf);
    printf("\n'e' raised to the power of %f minus one (e ^ %f - 1) = %.20f\n", xf, xf, resultf1);

    printf("\n'e' raised to the power of %Lf (e ^ %Lf) = %.20Lf\n", xL, xL, resultL);
    printf("\n'e' raised to the power of %Lf minus one (e ^ %Lf - 1) = %.20Lf\n", xL, xL, resultL1);

    return 0;
}
```
## 15.3 运行结果
![](expm1.png)

# 16. erf，erff，erfl
## 16.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double erf (double x);` |  计算 x 的 [误差函数](https://baike.baidu.com/item/%E8%AF%AF%E5%B7%AE%E5%87%BD%E6%95%B0/5890875?fr=aladdin)（double） |
|`float erff (float x);` |  计算 x 的 [误差函数](https://baike.baidu.com/item/%E8%AF%AF%E5%B7%AE%E5%87%BD%E6%95%B0/5890875?fr=aladdin)（float） |
|`long double erfl (long double x);` |  计算 x 的 [误差函数](https://baike.baidu.com/item/%E8%AF%AF%E5%B7%AE%E5%87%BD%E6%95%B0/5890875?fr=aladdin)（long double） |

## 16.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main(void)
{
    double result, x = 2.0;
    result = erf(x); // 高斯误差函数

    float resultf, xf = 2.0;
    resultf = erff(xf);

    long double resultL, xL = 2.0;
    resultL = erfl(xL);

    printf("the error function of %lf = %.20lf\n", x, result);
    printf("the error function of %f = %.20f\n", xf, resultf);
    printf("the error function of %Lf = %.20Lf\n", xL, resultL);    
    
    return 0;
}
```
## 16.3 运行结果
![](erf.png)


# 17. erfc，erfcf，erfcl
## 17.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double erfc (double x);` |  计算 x 的互补误差函数（double） |
|`float erfcf (float x);` |  计算 x 的互补误差函数（float） |
|`long double erfcl (long double x);` |  计算 x 的互补误差函数（long double） |


## 17.2 演示示例
```c
#include <stdio.h>
#include <math.h>

int main(void)
{
    double result, x = 2.0;
    result = erfc(x); // 互补误差函数

    float resultf, xf = 2.0;
    resultf = erfcf(xf);

    long double resultL, xL = 2.0;
    resultL = erfcl(xL);

    printf("the complementary error function of %lf = %.20lf\n", x, result);
    printf("the complementary error function of %f = %.20f\n", xf, resultf);
    printf("the complementary error function of %Lf = %.20Lf\n", xL, resultL);    
    
    return 0;
}
```
## 17.3 运行结果
![](erfc.png)

# 参考
1. [\[API Reference Document\]](https://www.apiref.com/c-zh/t_alph_c.htm)
2. [\[ecvt 函数\]](https://baike.baidu.com/item/ecvt/10942099?fr=aladdin)
3. [\[exec函数\]](https://learn.microsoft.com/zh-hk/cpp/c-runtime-library/reference/execl?view=msvc-170)