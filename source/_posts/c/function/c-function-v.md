---
title: C语言函数大全--v 开头的函数
date: 2023-05-11 22:10:28
updated: 2025-07-21 23:57:09
categories:
  - 开发语言-C
tags:
  - C语言函数大全
  - v 开头的函数
---

![](/images/cplus-logo.png)

# 总览
| 函数声明 |  函数功能  |
|:--|:--|
|`void va_start(va_list ap, last_arg);` | 用于初始化一个 `va_list` 类型的变量，使其指向可变参数列表中的第一个参数  |
|`type va_arg(va_list ap, type);` | 用于从可变参数列表中获取下一个参数，并将其转换为指定的类型  |
|`void va_copy(va_list dest, va_list src);` | 用于将一个 `va_list` 类型的变量复制到另一个变量中  |
|`void va_end(va_list ap);` |  用于清理一个 `va_list` 类型的变量 |
|`int vfprintf(FILE *stream, const char *format, va_list arg);` | 用于将格式化输出写入到指定的文件流中  |
|`int vfscanf(FILE *stream, const char *format, va_list arg);` | 用于将指定文件流中的格式化输入读取到指定变量中  |
|`int vprintf(const char *format, va_list ap);` | 它使用格式化字符串 `format` 中的指令来指定输出的格式，并将后续的可变参数按照指令指定的格式输出到标准输出流 `stdout`  |
|`int vscanf(const char *format, va_list arg);` | 它使用格式化字符串 `format` 中的指令来指定输入的格式，并从标准输入流 `stdin` 中读取数据，并将数据按照指令指定的格式存储到相应的变量中  |
|`int vsprintf(char *str, const char *format, va_list ap);` | 它使用格式化字符串 `format` 中的指令来指定输出的格式，并将后续的可变参数按照指令指定的格式输出到字符数组 `str` 中。  |
|`int vsscanf(const char *str, const char *format, va_list ap);` |  它使用格式化字符串 `format` 中的指令来指定输入的格式，并从字符数组 `str` 中读取数据，并将数据按照指令指定的格式存储到相应的变量中 |


# 1. va_start
## 1.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void va_start(va_list ap, last_arg);` | 用于初始化一个 `va_list` 类型的变量，使其指向可变参数列表中的第一个参数  |

**参数：**
- **ap：** 一个指向 va_list 类型的变量的指针，表示要初始化的可变参数列表
- **last_arg：** 一个指向最后一个固定参数的指针，用于确定可变参数列表的起始位置

## 1.2 演示示例
```c
#include <stdio.h>
#include <stdarg.h>

void print_args(int count, ...) 
{
    va_list args1, args2;
    int i;

    va_start(args1, count); /* 初始化可变参数列表 */

    /* 复制可变参数列表 */
    va_copy(args2, args1);

    for (i = 0; i < count; i++) 
    {
        int arg1 = va_arg(args1, int); /* 获取下一个参数 */
        int arg2 = va_arg(args2, int);
        printf("arg[%d] = %d %d\n", i, arg1, arg2);
    }

    va_end(args2); /* 清理复制的可变参数列表 */
    va_end(args1); /* 清理原始可变参数列表 */
}

int main() 
{
    print_args(3, 10, 20, 30);
    return 0;
}
```

在上面的示例代码中，
- 我们首先调用自定义的 `print_args()` 函数，并传入了 4 个入参，第一个为 可变参数的个数，后面三个为具体的整数型可变参数。
- 然后在`print_args()`  函数内部，我们首先定义了两个 `va_list` 类型的变量 `args1` 和 `args2`，并使用 `va_start()` 函数初始化 `args1` 变量。
- 接着，我们使用 `va_copy()` 函数将 `args1` 复制到 `args2` 中，并使用 `for` 循环和 两个 `va_arg()` 函数来分别访问这两个可变参数列表，并依次输出每个参数的值。
- 最后，我们使用两次 `va_end()` 函数来清理这两个可变参数列表。


## 1.3 运行结果
![](va.png)

# 2. va_arg
## 2.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`type va_arg(va_list ap, type);` | 用于从可变参数列表中获取下一个参数，并将其转换为指定的类型  |

**参数：**
- **ap ：** 一个指向 va_list 类型的变量的指针，表示要访问的可变参数列表
- **type ：** 一个类型说明符，表示下一个参数的类型

## 2.2 演示示例
可参考 1.2 中所示

# 3. va_copy
## 3.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void va_copy(va_list dest, va_list src);` | 用于将一个 `va_list` 类型的变量复制到另一个变量中  |

**参数：**
- **dest ：** 一个指向 va_list 类型的变量的指针，表示目标可变参数列表
- **src ：** 另一个指向 va_list 类型的变量的指针，表示要被复制的可变参数列表

## 3.2 演示示例
可参考 1.2 中所示

# 4. va_end
## 4.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void va_end(va_list ap);` |  用于清理一个 `va_list` 类型的变量 |

**参数：**
- **ap：** 一个指向 va_list 类型的变量的指针，表示要清理的可变参数列表

## 4.2 演示示例
可参考 1.2 中所示


# 5. vfprintf
## 5.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int vfprintf(FILE *stream, const char *format, va_list arg);` | 用于将格式化输出写入到指定的文件流中  |

**参数：**
- **stream ：** 要写入数据的文件流指针
- **format ：** 格式化字符串，用来指定要输出的内容和格式
- **arg ：** 一个 `va_list` 类型的变量，包含了可变参数列表

## 5.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>

FILE *fp;

int vfpf(const char *fmt, ...)
{
    va_list argptr;
    int cnt;

    va_start(argptr, fmt);
    cnt = vfprintf(fp, fmt, argptr);
    va_end(argptr);

    return cnt;
}

int main()
{
    int inumber;
    float fnumber;
    char string[4];

    fp = tmpfile();
    if (fp == NULL)
    {
        perror("tmpfile() call");
        exit(1);
    }

    vfpf("%d %f %s", 30, 90.0, "abc");

    rewind(fp);
    
    fscanf(fp,"%d %f %s", &inumber, &fnumber, string);
    
    printf("%d %.2f %s\n", inumber, fnumber, string);
    
    fclose(fp);

    return 0;
}
```

在上述的示例代码中，
- 我们首先声明了三个变量 `inumber`、`fnumber` 和 `string`；
- 然后，调用 `tmpfile()` 函数创建一个临时文件，并将返回的文件指针赋值给全局变量 `fp`。如果创建文件失败，则打印错误信息并退出程序；
- 接着，调用自定义的 `vfpf()` 函数来向临时文件中写入数据。它里面使用 `vfprintf()` 函数将格式化输出写入到一个文件流中；
- 再然后，调用 `rewind()` 函数将文件指针重新定位到文件开头；
- 再接着使用 `fscanf()` 函数从文件中读取数据，并使用 `printf()` 函数中输出从文件中读取的数据【其中浮点数部分保留两位小数】
- 最后调用 `fclose()` 函数关闭文件指针，并结束程序

## 5.3 运行结果
![](vfprintf.png)

# 6. vfscanf
## 6.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int vfscanf(FILE *stream, const char *format, va_list arg);` | 用于将指定文件流中的格式化输入读取到指定变量中  |

**参数：**
- **stream ：** 要读取数据的文件流指针
- **format ：** 格式化字符串，用来指定要读取的内容和格式
- **arg ：** 一个 `va_list` 类型的变量，包含了可变参数列表

**返回值：**
- 如果读取成功，返回成功读取并赋值给变量的项目数；
- 如果出现错误，则返回负数。

## 6.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>

FILE *fp;

int vfsf(const char *fmt, ...) {
    va_list argptr;
    int cnt;

    va_start(argptr, fmt);
    cnt = vfscanf(fp, fmt, argptr);
    va_end(argptr);

    return cnt;
}

int main()
{
    int inumber;
    float fnumber;
    char string[4];

    fp = tmpfile();
    if (fp == NULL)
    {
        perror("tmpfile() call");
        exit(1);
    }

    fprintf(fp, "%d %f %s", 30, 90.0, "abc");

    rewind(fp);
    
    vfsf("%d %f %s", &inumber, &fnumber, string);

    printf("%d %.2f %s\n", inumber, fnumber, string);

    fclose(fp);

    return 0;
}
```

在上面的示例代码中，
- 我们首先声明了三个变量 `inumber`、`fnumber` 和 `string`；
- 然后，程序调用 `tmpfile()` 函数创建一个临时文件，并将返回的文件指针赋值给全局变量 `fp`。如果创建文件失败，则打印错误信息并退出程序；
- 接着，我们使用 `fprintf()` 函数将三个数据（一个整型数字、一个浮点数和一个字符串）写入该文件中；
- 再然后，调用 `rewind()` 函数将文件指针重新定位到文件开头；
- 再接着，我们调用自定义的 `vfsf()` 函数，里面使用 `vfscanf()` 函数从文件中读取数据；
- 最后，打印出从文件中读取的数据，并关闭临时文件，退出程序。

## 6.3 运行结果
![](vfscanf.png)

# 7. vprintf
## 7.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int vprintf(const char *format, va_list ap);` | 它使用格式化字符串 `format` 中的指令来指定输出的格式，并将后续的可变参数按照指令指定的格式输出到标准输出流 `stdout`  |

**参数：**
- **format ：** 格式化字符串，用来指定要打印的内容和格式
- **va_list ：** 一个 `va_list` 类型的变量，包含了可变参数列表

## 7.2 演示示例
```c
#include <stdio.h>
#include <stdarg.h>

void myprint(const char *format, ...) 
{
    va_list args;
    va_start(args, format);
    vprintf(format, args);
    va_end(args);
}

int main() 
{
    int a = 10;
    float b = 3.14;
    char s[] = "hello";
    
    myprint("a=%d, b=%.2f, s=%s\n", a, b, s);
    return 0;
}
```
在上面的示例代码中，
- 我们首先定义了三个变量 整形 `a` 、浮点型 `b` 和 字符数组 `s`；
- 然后，调用自定义的 `myprint()` 函数将这些变量的值输出到标准输出流 stdout 中；
 在 `myprint()` 函数中， 
     - 我们首先使用 `va_start()` 宏初始化一个 `va_list` 变量 `args`；
     - 然后调用 `vprintf()` 函数将格式化字符串和参数列表传递给该函数进行输出；
     - 最后使用 `va_end()` 宏清理 `args` 变量。
- 最后结束程序。

## 7.3 运行结果
![](vprintf.png)


# 8. vscanf
## 8.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int vscanf(const char *format, va_list arg);` | 它使用格式化字符串 `format` 中的指令来指定输入的格式，并从标准输入流 `stdin` 中读取数据，并将数据按照指令指定的格式存储到相应的变量中  |

**参数：**
- **format ：** 格式化字符串，用来指定要读取的内容和格式
- **va_list ：** 一个 `va_list` 类型的变量，包含了可变参数列表

## 8.2 演示示例
```c
#include <stdio.h>
#include <stdarg.h>

void myscan(const char *format, ...) 
{
    va_list args;
    va_start(args, format);
    vscanf(format, args);
    va_end(args);
}

int main() 
{
    int a;
    float b;
    char s[10];
    
    myscan("%d%f%s", &a, &b, s);
    printf("a=%d, b=%.2f, s=%s\n", a, b, s);
    return 0;
}
```

在上面的示例代码中，
- 我们首先定义了三个变量 整形 `a` 、浮点型 `b` 和 字符数组 `s`；
- 然后，调用自定义的 `myscan()` 函数从标准输入流 `stdin` 中读取数据，并将数据存储到这些变量中
 在 `myscan()` 函数中， 
     - 我们首先使用 `va_start()` 宏初始化一个 `va_list` 变量 `args`；
     - 然后调用 `vscanf()` 函数将格式化字符串和参数列表传递给该函数进行输入；
     - 最后使用 `va_end()` 宏清理 `args` 变量。
- 最后我们打印输出上面输入的两个变量的数据，并结束程序。

## 8.3 运行结果
![](vscanf.png)

# 9. vsprintf
## 9.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int vsprintf(char *str, const char *format, va_list ap);` | 它使用格式化字符串 `format` 中的指令来指定输出的格式，并将后续的可变参数按照指令指定的格式输出到字符数组 `str` 中。  |

**参数：**
- **str：** 用来存储格式化数据的字符数组的指针
- **format ：** 格式化字符串，用来指定要输出的内容和格式
- **va_list ：** 一个 `va_list` 类型的变量，包含了可变参数列表

## 9.2 演示示例
```c
#include <stdio.h>
#include <stdarg.h>

void myprint(const char *format, ...) 
{
    char buffer[100];
    va_list args;
    va_start(args, format);
    vsprintf(buffer, format, args);
    va_end(args);
    printf("%s", buffer);
}

int main() 
{
    int a = 10;
    float b = 3.1415;
    char s[] = "huazie";
    
    myprint("a=%d, b=%.4f, s=%s\n", a, b, s);
    return 0;
}
```

在上面的示例代码中，
- 我们首先定义了三个变量 整形 `a` 、浮点型 `b` 和 字符数组 `s`；
- 然后，调用自定义的 `myprint()` 函数将格式化字符串和这些变量的值输出到字符数组 `buffer` 中，并打印输出字符数组 `buffer` ；
 在 `myprint()` 函数中， 
     - 我们首先使用 `va_start()` 宏初始化一个 `va_list` 变量 `args`；
     - 然后调用 `vsprintf()` 函数将格式化字符串和参数列表传递给该函数进行输出，并将输出结果存储到 `buffer` 数组中；
     - 最后使用 `va_end()` 宏清理 `args` 变量。
- 最后结束程序。

## 9.3 运行结果
![](vsprintf.png)

# 10. vsscanf
## 10.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int vsscanf(const char *str, const char *format, va_list ap);` |  它使用格式化字符串 `format` 中的指令来指定输入的格式，并从字符数组 `str` 中读取数据，并将数据按照指令指定的格式存储到相应的变量中 |

**参数：**
- **str：** 用来读取的格式化数据的字符数组的指针
- **format ：** 格式化字符串，用来指定要读取的内容和格式
- **va_list ：** 一个 `va_list` 类型的变量，包含了可变参数列表

## 10.2 演示示例
```c
#include <stdio.h>
#include <stdarg.h>

void myscan(const char *str, const char *format, ...) 
{
    va_list args;
    va_start(args, format);
    vsscanf(str, format, args);
    va_end(args);
}

int main() 
{
    int a;
    float b;
    char s[10];
    char buffer[] = "10 3.1415 huazie";
    
    myscan(buffer, "%d%f%s", &a, &b, s);
    printf("a=%d, b=%.4f, s=%s\n", a, b, s);
    return 0;
}
```

在上面的示例代码中，
- 我们首先定义了四个变量 整形 `a` 、浮点型 `b` 、 字符数组 `s` 和 字符数组 `buffer`；
- 然后，调用自定义的 `myscan()` 函数从字符数组 `buffer` 中读取数据，并将数据存储到另外 `3` 个变量中；
 在 `myscan()` 函数中， 
     - 我们首先使用 `va_start()` 宏初始化一个 `va_list` 变量 `args`；
     - 然后调用 `vsscanf()` 函数将字符数组 `buffer` 和格式化字符串以及参数列表传递给该函数进行输入，并将数据存储到相应的变量中；
     - 最后使用 `va_end()` 宏清理 `args` 变量。
- 最后我们打印输出上面从字符数组 `buffer` 中读取并输入的三个变量的数据，并结束程序。

## 10.3 运行结果
![](vsscanf.png)

# 参考
1. [\[API Reference Document\]](https://www.apiref.com/c-zh/t_alph_v.htm)


