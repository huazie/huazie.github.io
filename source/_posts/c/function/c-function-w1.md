---
title: C语言函数大全--w 开头的函数（1）
date: 2023-05-12 22:19:42
updated: 2025-07-22 21:57:02
categories:
  - 开发语言-C
tags:
  - C语言函数大全
  - w 开头的函数
---

![](/images/cplus-logo.png)

# 总览
| 函数声明 |  函数功能  |
|:--|:--|
|`wchar_t * wcscat(wchar_t *dest, const wchar_t *src);` | 用于将一个宽字符字符串追加到另一个宽字符字符串的末尾  |
|`wchar_t *wcschr(const wchar_t *str, wchar_t wc);` | 用于在宽字符串中查找指定字符的位置  |
|`int wcscmp(const wchar_t *str1, const wchar_t *str2);` |  用于比较两个宽字符串的大小；它将两个字符串逐个字符进行比较，直到遇到不同的字符或者其中一个字符串结束为止 |
|`int wcscoll(const wchar_t *str1, const wchar_t *str2);` | 用于比较两个宽字符串的大小  |
|`wchar_t *wcscpy(wchar_t *dest, const wchar_t *src);` | 用于将一个宽字符串复制到另一个字符串中  |
|`size_t wcsftime(wchar_t *str, size_t maxsize, const wchar_t *format, const struct tm *timeptr);` |  用于将日期和时间格式化为宽字符字符串 |
|`size_t wcslen(const wchar_t *str);` |  用于计算宽字符串的长度 |
|`wchar_t *wcsncat(wchar_t *dest, const wchar_t *src, size_t n);` | 用于将一个宽字符串的一部分追加到另一个宽字符串末尾  |
|`int wcsncmp(const wchar_t *str1, const wchar_t *str2, size_t n);` |  用于比较两个宽字符串的前若干个字符是否相同 |
|`wchar_t *wcsncpy(wchar_t *dest, const wchar_t *src, size_t n);` | 用于将一个宽字符串的一部分复制到另一个宽字符串中  |
|`size_t wcsrtombs(char *dest, const wchar_t **src, size_t n, mbstate_t *ps);` | 用于将宽字符串转换为多字节字符串  |
|`wchar_t *wcsstr(const wchar_t *haystack, const wchar_t *needle);` |  用于在一个宽字符串中查找另一个宽字符串 |
|`size_t wcsspn(const wchar_t *str, const wchar_t *accept);` | 用于查找宽字符串中连续包含某些字符集合中的字符的最长前缀  |
|`double wcstod(const wchar_t *nptr, wchar_t **endptr);` | 用于将宽字符串转换为双精度浮点数  |

# 1. wcscat
## 1.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`wchar_t * wcscat(wchar_t *dest, const wchar_t *src);` | 用于将一个宽字符字符串追加到另一个宽字符字符串的末尾  |

**参数：**
- **dest  ：**  目标字符串
- **src ：**  源字符串


## 1.2 演示示例
```c
#include <stdio.h>
#include <wchar.h>

int main() {
    wchar_t dest[30] = L"Hello";
    const wchar_t *src = L", Huazie!";

    wcscat(dest, src);

    wprintf(L"%ls\n", dest);

    return 0;
}
```

在上面的示例代码中，
- 首先，我们定义了一个大小为 `30` 的 `wchar_t` 数组 `dest`，并初始化为 `"Hello"`;
- 然后，定义了一个指向常量宽字符串的指针 `src`，指向字符串 `", Huazie!"`；
- 接着。调用 `wcscat()` 函数将 `src` 字符串中的所有字符追加到 `dest` 字符串的末尾，形成新的宽字符字符串 `"Hello, Huazie!"`；
- 最后，使用 `wprintf()` 函数将新的字符串输出到控制台。

**注意：** 在使用 `wcscat()` 函数时，需要确保目标字符串 `dest` 的空间足够大，以容纳源字符串 `src` 的所有字符和一个结束符（`\0`）。如果目标字符串的空间不足，可能会导致数据覆盖和未定义行为。

## 1.3 运行结果
![](wcscat.png)

# 2. wcschr
## 2.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`wchar_t *wcschr(const wchar_t *str, wchar_t wc);` | 用于在宽字符串中查找指定字符的位置  |

**参数：**
- **str ：**  要查找的宽字符串
- **wc ：**  要查找的宽字符

**返回值：**
- 如果在 `str` 字符串中查找到第一个与 `wc` 相等的宽字符，则返回该字符在字符串中的地址；
- 否则没有找到匹配字符，则返回空指针。

## 2.2 演示示例
```c
#include <stdio.h>
#include <wchar.h>

int main() 
{
    const wchar_t *str = L"hello, huazie";
    wchar_t c = L'u';

    wchar_t *p = wcschr(str, c);

    if (p != NULL) 
        printf("Found %lc at position %d.\n", c, p - str);
    else
        printf("%lc not found.\n", c);

    return 0;
}
```

在上面的示例代码中，
- 首先，我们定义了一个宽字符串 `str`，并初始化为 `"hello, huazie"`；
- 然后，定义了一个宽字符 `c`，值为 `'u'`；
- 接着调用 `wcschr()` 函数在 str 字符串中查找字符 `c`，并将返回结果保存在指针变量 `p` 中；
- 最后，根据返回值判断是否找到了匹配字符，并输出相应的信息。

## 2.3 运行结果
![](wcschr.png)


# 3. wcscmp
## 3.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int wcscmp(const wchar_t *str1, const wchar_t *str2);` |  用于比较两个宽字符串的大小；它将两个字符串逐个字符进行比较，直到遇到不同的字符或者其中一个字符串结束为止 |

**参数：**
- **str1：**  待比较的宽字符串1
- **str2：**  待比较的宽字符串2


**返回值：**
- 如果 `str1` 小于 `str2`，则返回一个负整数；
- 如果 `str1` 等于 `str2`，则返回 0；
- 如果 `str1` 大于 `str2`，则返回一个正整数。

## 3.2 演示示例
```c
#include <stdio.h>
#include <wchar.h>

int main() 
{
    const wchar_t *str1 = L"hello";
    const wchar_t *str2 = L"huazie";

    int result = wcscmp(str1, str2);

    if (result < 0) 
        printf("%ls is less than %ls.\n", str1, str2);
    else if (result == 0)
        printf("%ls is equal to %ls.\n", str1, str2);
    else
        printf("%ls is greater than %ls.\n", str1, str2);

    return 0;
}

```

在上述的示例代码中，
- 首先，我们定义了两个宽字符串 `str1` 和 `str2`，分别初始化为 `"hello"` 和 `"huazie"`；
- 然后，调用 `wcscmp()` 函数比较两个字符串的大小，并将返回结果保存在变量 `result` 中；
- 接着根据 `result` 的值，输出相应的比较结果；
- 最后结束程序。

## 3.3 运行结果
![](wcscmp.png)

# 4. wcscoll
## 4.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int wcscoll(const wchar_t *str1, const wchar_t *str2);` | 用于比较两个宽字符串的大小  |

**参数：**
- **str1：**  待比较的宽字符串1
- **str2：**  待比较的宽字符串2

## 4.2 演示示例
```c
#include <stdio.h>
#include <wchar.h>
#include <locale.h>

int main() 
{
    setlocale(LC_ALL, ""); // 设置本地化环境

    const wchar_t *str1 = L"hello";
    const wchar_t *str2 = L"huazie";

    int result = wcscoll(str1, str2);

    if (result < 0)
        printf("%ls is less than %ls.\n", str1, str2);
    else if (result == 0)
        printf("%ls is equal to %ls.\n", str1, str2);
    else
        printf("%ls is greater than %ls.\n", str1, str2);

    return 0;
}
```

在上述的示例代码中，
- 首先，我们调用 `setlocale()` 函数设置本地化环境为当前系统默认设置；
- 然后，定义了两个宽字符串 `str1` 和 `str2`，分别初始化为 "hello" 和 "huazie"；
- 接着，调用 `wcscoll()` 函数比较两个字符串的大小，并将返回结果保存在变量 `result` 中。
- 最后根据 `result` 的值，输出相应的比较结果。

**注意：** 在使用 `wcscoll()` 函数比较宽字符串大小时，需要确保本地化环境正确设置，以便该函数能够正常工作。如果没有设置本地化环境或者设置错误，可能会导致比较结果不准确。

## 4.3 运行结果
![](wcscoll.png)

# 5. wcscpy
## 5.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`wchar_t *wcscpy(wchar_t *dest, const wchar_t *src);` | 用于将一个宽字符串复制到另一个字符串中  |

**参数：**
- **dest  ：**  目标字符串
- **src ：**  源字符串

## 5.2 演示示例
```c
#include <stdio.h>
#include <wchar.h>

int main() 
{
    wchar_t dest[20];
    const wchar_t *src = L"Hello, huazie!";

    wcscpy(dest, src);

    wprintf(L"%ls\n", dest);

    return 0;
}
```

在上述的示例代码中，
- 首先，我们定义了一个大小为 `20` 的 `wchar_t` 数组 `dest`；
- 然后，定义了一个指向常量宽字符串的指针 `src`，指向字符串 "Hello, huazie!"；
- 接着，调用 `wcscpy()` 函数将 `src` 字符串中的所有字符复制到 `dest` 字符串中，形成新的宽字符字符串 `dest`；
- 最后，使用 `wprintf()` 函数将新的字符串输出到控制台。

## 5.3 运行结果
![](wcscpy.png)

# 6. wcsftime
## 6.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`size_t wcsftime(wchar_t *str, size_t maxsize, const wchar_t *format, const struct tm *timeptr);` |  用于将日期和时间格式化为宽字符字符串 |

**参数：**
- **str ：**  输出结果的缓冲区
- **maxsize ：**  缓冲区的大小
- **format ：**  格式化字符串
- **timeptr ：**  包含日期和时间信息的结构体指针

## 6.2 演示示例
```c
#include <stdio.h>
#include <wchar.h>
#include <time.h>

int main() 
{
    time_t current_time;
    struct tm *time_info;
    wchar_t buffer[80];

    time(&current_time);
    time_info = localtime(&current_time);

    wcsftime(buffer, 80, L"%c", time_info);
    wprintf(L"The current date and time is: %ls\n", buffer);

    return 0;
}
```

在上面的示例代码中，
- 首先，我们定义了一个变量 `current_time` 来存储当前时间，以及一个指向 `struct tm` 类型的指针 `time_info`；
- 然后，调用 `time()` 函数获取当前时间，并使用 `localtime()` 函数将时间转换为本地时间，存储在 `time_info` 指针变量中；
- 接着，调用 `wcsftime()` 函数将日期和时间格式化为宽字符字符串，并存储到缓冲区 `buffer` 中；
- 最后，使用 `wprintf()` 函数输出格式化后的字符串。

## 6.3 运行结果
![](wcsftime.png)

# 7. wcslen
## 7.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`size_t wcslen(const wchar_t *str);` |  用于计算宽字符串的长度 |

**参数：**
- **str ：**  要计算长度的宽字符串

## 7.2 演示示例
```c
#include <stdio.h>
#include <wchar.h>

int main() 
{
    const wchar_t *str = L"Hello, huazie!";
    size_t len = wcslen(str);

    wprintf(L"The length of '%ls' is %zu.\n", str, len);

    return 0;
}
```

在上面的示例代码中，
- 首先，我们定义了一个指向常量宽字符串的指针 `str`，指向字符串 `"Hello, huazie!"`；
- 然后，调用 `wcslen()` 函数计算 `str` 字符串的长度，并将结果保存在变量 `len` 中；
- 最后，使用 `wprintf()` 函数输出字符串的长度。

## 7.3 运行结果
![](wcslen.png)

# 8. wcsncat
## 8.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`wchar_t *wcsncat(wchar_t *dest, const wchar_t *src, size_t n);` | 用于将一个宽字符串的一部分追加到另一个宽字符串末尾  |

**参数：**
- **dest  ：**  目标字符串
- **src ：**  源字符串
- **n ：**  要追加的字符数

## 8.2 演示示例
```c
#include <stdio.h>
#include <wchar.h>

int main() 
{
    wchar_t dest[20] = L"Hello, ";
    const wchar_t *src = L"huazie!";
    size_t n = 3;

    wcsncat(dest, src, n);

    wprintf(L"%ls\n", dest);

    return 0;
}
```

在上述的示例代码中，
- 首先，我们定义了一个大小为 `20` 的 `wchar_t` 数组 `dest`，并初始化为 `"Hello, "`；
- 然后，定义了一个指向常量宽字符串的指针 `src`，指向字符串 `"huazie!"`。
- 接着，调用 `wcsncat()` 函数将 `src` 字符串中的前 `3` 个字符追加到 `dest` 字符串的末尾，形成新的宽字符字符串 `dest`；
- 最后，使用 `wprintf()` 函数将新的字符串输出到控制台。

## 8.3 运行结果
![](wcsncat.png)

# 9. wcsncmp
## 9.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int wcsncmp(const wchar_t *str1, const wchar_t *str2, size_t n);` |  用于比较两个宽字符串的前若干个字符是否相同 |

**参数：**
- **str1 ：**  待比较的宽字符串1
- **str2 ：**  待比较的宽字符串2
- **n ：**  要比较的字符数


**返回值：**
将字符串 `str1` 和字符串 `str2` 中的前 `n` 个字符进行比较
- 如果两个字符串相同，返回值为 `0`；
- 如果字符串 `str1` 在前 `n` 个字符中的第一个不同于字符串 `str2` 对应字符的字符大于字符串 `str2` 对应字符的字符，返回值为正数；
- 如果字符串 `str1` 在前 `n` 个字符中的第一个不同于字符串 `str2` 对应字符的字符小于字符串 `str2` 对应字符的字符，返回值为负数。

## 9.2 演示示例
```c
#include <stdio.h>
#include <wchar.h>

int main() 
{
    const wchar_t *str1 = L"hello";
    const wchar_t *str2 = L"huazie";
    size_t n = 3;

    int result = wcsncmp(str1, str2, n);

    if (result < 0)
        printf("%ls is less than %ls in the first %zu characters.\n", str1, str2, n);
    else if (result == 0)
        printf("%ls is equal to %ls in the first %zu characters.\n", str1, str2, n);
    else
        printf("%ls is greater than %ls in the first %zu characters.\n", str1, str2, n);

    return 0;
}
```

## 9.3 运行结果
![](wcsncmp.png)

# 10. wcsncpy
## 10.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`wchar_t *wcsncpy(wchar_t *dest, const wchar_t *src, size_t n);` | 用于将一个宽字符串的一部分复制到另一个宽字符串中  |

**参数：**
- **dest  ：**  目标字符串
- **src ：**  源字符串
- **n ：**  要复制的字符数

## 10.2 演示示例
```c
#include <stdio.h>
#include <wchar.h>

int main() 
{
    wchar_t dest[20];
    const wchar_t *src = L"Hello, huazie!";
    size_t n = 5;

    wcsncpy(dest, src, n);
    dest[n] = L'\0';

    wprintf(L"%ls\n", dest);

    return 0;
}
```

在上面的示例代码中，
- 首先，我们定义了一个大小为 `20` 的 `wchar_t` 数组 `dest`，以及一个指向常量宽字符串的指针 `src`，指向字符串 `"Hello, huazie!"`；
- 然后，调用 `wcsncpy()` 函数将 `src` 字符串中的前 `5` 个字符复制到 `dest` 字符串中，形成新的宽字符字符串 `dest`；
- 最后，手动添加空字符 `\0` 以确保字符串结束，并使用 `wprintf()` 函数将新的字符串输出到控制台。

## 10.3 运行结果
![](wcsncpy.png)

# 11. wcsrtombs
## 11.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`size_t wcsrtombs(char *dest, const wchar_t **src, size_t n, mbstate_t *ps);` | 用于将宽字符串转换为多字节字符串  |

**参数：**
- **dest  ：**  输出结果的缓冲区
- **src ：**  指向源字符串的指针
- **n ：**  要转换的最大字符数
- **ps ：**  一个指向转换状态的指针

## 11.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <wchar.h>
#include <locale.h>

int main() 
{
    setlocale(LC_ALL, "");

    const wchar_t *src = L"Hello, huazie!";
    size_t n = wcslen(src) + 1;

    char *dest = (char *) malloc(n * sizeof(char));

    if (dest == NULL) 
    {
        fprintf(stderr, "Memory allocation failed.\n");
        return EXIT_FAILURE;
    }

    mbstate_t state = {0};

    wcsrtombs(dest, &src, n, &state);

    printf("%s\n", dest);

    free(dest);

    return 0;
}
```

在上面的示例代码中，
- 首先，我们调用 `setlocale()` 函数设置程序的本地化环境，以便可以正确地进行宽字符和多字节字符之间的转换；
- 然后，定义了一个指向常量宽字符串的指针 `src`，指向字符串 `"Hello, huazie!"`；
- 接着，根据源字符串的长度分配了足够大小的缓冲区 `dest`，并初始化为 `0`；
- 再然后，调用 `wcsrtombs()` 函数将宽字符串 `src` 转换为多字节字符串，并存储到缓冲区 `dest` 中；
- 最后，使用 `printf()` 函数输出多字节字符串，并释放目标缓存区的内存。

**注意：** 在使用 `wcsrtombs()` 函数进行宽字符和多字节字符转换时，需要确保程序的本地化环境已经正确设置，否则可能会导致转换失败或者输出结果不正确。此外，在分配缓冲区 `dest` 的大小时，可以考虑将源字符串的长度加 `1`，以容纳字符串的结尾空字符（`\0`）。最后在使用完毕后要记得释放缓冲区的内存。

## 11.3 运行结果
![](wcsrtombs.png)

# 12. wcsstr
## 12.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`wchar_t *wcsstr(const wchar_t *haystack, const wchar_t *needle);` |  用于在一个宽字符串中查找另一个宽字符串 |

**参数：**
- **haystack ：**  要查找的宽字符串
- **needle ：**  要查找的子串

**返回值：**
- 如果从 `haystack` 字符串中查找到第一个匹配 `needle` 子串的位置，则返回指向匹配位置的指针；
- 否则未找到匹配的子串，返回空指针（`NULL`）。

## 12.2 演示示例
```c
#include <stdio.h>
#include <wchar.h>

int main() 
{
    const wchar_t *haystack = L"Hello, huazie!";
    const wchar_t *needle = L"huazie";

    wchar_t *result = wcsstr(haystack, needle);

    if (result != NULL)
        wprintf(L"The substring '%ls' was found at position %d.\n", needle, result - haystack);
    else
        wprintf(L"The substring '%ls' was not found in '%ls'.\n", needle, haystack);

    return 0;
}
```

在上述的示例代码中，
- 首先，我们定义了一个指向常量宽字符串的指针 `haystack`，指向字符串 `"Hello, huazie!"`。
- 然后，定义了一个指向常量宽字符串的指针 `needle`，指向字符串 `"huazie"`；
- 接着，调用 `wcsstr()` 函数在 `haystack` 字符串中查找子串 `needle`，并将结果指针保存在变量 `result` 中。
- 最后，根据 `result` 的值，输出相应的查找结果。

**注意：** 在使用 `wcsstr()` 函数查找子串时，该函数会自动遍历整个字符串，直到找到匹配的子串或者结束字符串。如果要查找的子串在字符串中多次出现，该函数将返回第一次出现的位置，并不会考虑后续的匹配。

## 12.3 运行结果
![](wcsstr.png)

# 13. wcsspn
## 13.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`size_t wcsspn(const wchar_t *str, const wchar_t *accept);` | 用于查找宽字符串中连续包含某些字符集合中的字符的最长前缀  |

**参数：**
- **str ：**  要查找的宽字符串
- **accept ：**  一个包含要接受字符集合的宽字符串

## 13.2 演示示例
```c
#include <stdio.h>
#include <wchar.h>

int main() 
{
    const wchar_t *str = L"123456789a0";
    const wchar_t *accept = L"0123456789";

    size_t length = wcsspn(str, accept);

    wprintf(L"The length of the prefix that contains digits is %zu.\n", length);

    return 0;
}
```
在上述的示例代码中，
- 首先，我们定义了一个指向常量宽字符串的指针 `str`，指向字符串 `"123456789a0"`；
- 然后，定义了一个指向常量宽字符串的指针 `accept`，指向字符串 `"0123456789"`，表示数字字符的集合；
- 接着，调用 `wcsspn()` 函数查找 `str` 字符串中连续包含数字字符集合中的字符的最长前缀，并将返回结果保存在变量 `length` 中。
- 最后，根据 `length` 的值，调用 `wprintf()`  函数 输出 最长前缀的长度。

**注意：** 在使用 `wcsspn()` 函数查找宽字符串中的字符集合时，该函数会自动遍历整个字符串，直到找到第一个不在字符集合中的字符或者结束字符串。如果要查找的字符集合为空串，则返回 `0`。

## 13.3 运行结果
![](wcsspn.png)


# 14. wcstod
## 14.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double wcstod(const wchar_t *nptr, wchar_t **endptr);` | 用于将宽字符串转换为双精度浮点数  |

**参数：**
- **nptr ：**  要转换的宽字符串
- **endptr ：**  一个指向指针的指针，用于存储第一个无法解析的字符的位置

**返回值：**
- 如果从 `nptr` 字符串中解析出一个双精度浮点数，则返回该数值；
- 如果 `nptr` 字符串不包含有效的浮点数，则返回 `0`；
- 如果存在无法解析的字符，存储其位置到 `endptr` 指针中。

## 14.2 演示示例
```c
#include <stdio.h>
#include <wchar.h>

int main() 
{
    const wchar_t *str = L"3.14159265358979323846";
    //const wchar_t *str = L"3.141592653589793a23846";
    //const wchar_t *str = L"a3.14159265358979323846";
    wchar_t *end;

    double number = wcstod(str, &end);

    if (end == str)
        wprintf(L"No digits were found.\n");
    else if (*end != L'\0')
        wprintf(L"Invalid character at position %ld: '%lc'.\n", end - str, *end);
    else
        wprintf(L"The parsed number is %.20lf.\n", number);

    return 0;
}
```

在上面的示例代码中，
- 首先，我们定义了一个指向常量宽字符串的指针 `str`，指向字符串 `"3.14159265358979323846"`；
- 然后，调用 `wcstod()` 函数将字符串转换为双精度浮点数，并将结果保存在变量 `number` 中。
- 最后，根据函数的返回结果 `number` 和 `endptr` 指针所指向的值，输出相应的转换结果。

**注意：** 在使用 `wcstod()` 函数转换宽字符串为双精度浮点数时，要确保字符串中只包含有效的浮点数表示，否则可能会导致转换错误或者未定义行为。

## 14.3 运行结果
![](wcstod.png)
![](wcstod1.png)
![](wcstod2.png)






