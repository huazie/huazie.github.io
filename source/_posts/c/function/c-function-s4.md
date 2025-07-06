---
title: C语言函数大全--s 开头的函数（4）
date: 2023-05-06 14:28:08
updated: 2025-07-06 19:22:36
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
|`char * strdup(const char *s);` | 用于将一个以 `NULL` 结尾的字符串复制到新分配的内存空间中  |
|`int stricmp(const char *s1, const char *s2);` | 用于比较两个字符串的字母序是否相等，忽略大小写  |
|`char *strerror(int errnum);` | 用于将指定的错误码转换为相应的错误信息  |
|`int strcmpi(const char *s1, const char *s2);` |  用于比较两个字符串的字母序是否相等，忽略大小写 |
|`int strncmp(const char *s1, const char *s2, size_t n);` | 用于比较两个字符串的前n个字符是否相等  |
|`int strncmpi(const char *s1, const char *s2, size_t n);` | 用于比较两个字符串的前n个字符是否相等，忽略大小写  |
|`char *strncpy(char *dest, const char *src, size_t n);` | 用于将一个字符串的一部分拷贝到另一个字符串中  |
|`int strnicmp(const char *s1, const char *s2, size_t n);` | 用于比较两个字符串的前n个字符是否相等，忽略大小写  |
|`char *strnset(char *str, int c, size_t n);` |  用于将一个字符串的前n个字符都设置为指定的字符 |
|`char *strpbrk(const char *str1, const char *str2);` | 用于在一个字符串中查找任意给定字符集合中的字符的第一次出现位置 |
|`char *strrchr(const char *str, int character);` | 在给定的字符串中查找指定字符的最后一个匹配项  |
|`char *strrev(char *str);` | 将给定字符串中的所有字符顺序颠倒，并返回颠倒后的字符串  |
|`char *strset(char *str, int character);` |  用于设置给定字符串中的所有字符为指定的值，并返回修改后的字符串 |
|`size_t strspn(const char *str1, const char *str2);` | 计算字符串 str1 中包含在字符串 str2 中的前缀子字符串长度，并返回该长度值  |
|`char *strstr(const char *str1, const char *str2);` | 在字符串 str1 中查找第一个出现的字符串 str2，如果找到了，则返回指向该位置的指针；否则，返回 NULL  |
|`double strtod(const char *str, char **endptr);` | 将字符串 str 转换为一个浮点数（double 类型），并返回该浮点数。如果发生了转换错误，则返回 0.0，并且可以通过 endptr 指针返回指向第一个无法转换的字符的指针。  |
|`char *strtok(char *str, const char *delim);` | 用于将一个字符串分割成多个子字符串  |
|`long int strtol(const char *str, char **endptr, int base);` |  将字符串 str 转换为一个长整型数（long int 类型） |
|`char *strupr(char *str);` | 将字符串 str 中的所有小写字母转换为大写字母，并返回指向该字符串的指针  |
| `void swab(const void *src, void *dest, ssize_t nbytes);`| 将源缓冲区中的每两个相邻字节进行互换，然后将结果存储到目标缓冲区中  |
|`int system(const char *command);` |  执行一个 `shell` 或 `cmd` 命令，并等待命令的完成 |


# 1. strdup
## 1.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`char * strdup(const char *s);` | 用于将一个以 `NULL` 结尾的字符串复制到新分配的内存空间中  |

**注意：** `strdup()` 函数返回指向新分配的内存空间的指针，如果空间不足则返回 `NULL`。调用者负责释放返回的指针所指向的内存空间。 `strdup()` 函数与`strcpy()` 函数类似，但是它会动态地分配内存空间，而 `strcpy()` 需要调用者提供足够大的目标缓冲区。

## 1.2 演示示例
```c
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

int main() 
{
    char *source = "Hello World!";
    char *destination = strdup(source);
    
    if(destination != NULL) 
    {
        printf("Original string: %s\n", source);
        printf("Duplicated string: %s\n", destination);
        free(destination); // 释放内存，避免内存泄漏
    }
    
    return 0;
}
```

## 1.3 运行结果
![](strdup.png)


# 2. stricmp
## 2.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int stricmp(const char *s1, const char *s2);` | 用于比较两个字符串的字母序是否相等，忽略大小写  |

**返回值：**
- 如果 `s1` 和 `s2` 代表的字符串相等（忽略大小写），则返回 `0`；
- 如果 `s1` 比 `s2` 小，则返回负数；
- 如果 `s1` 比 `s2` 大，则返回正数。

## 2.2 演示示例
```c
#include <stdio.h>
#include <string.h>

int main() 
{
    char *s1 = "Hello";
    char *s2 = "hElLo";
    
    int result = stricmp(s1, s2); // 忽略大小写比较
    
    if(result == 0) 
    {
        printf("The strings are equal.\n");
    } 
    else if(result < 0) 
    {
        printf("%s is less than %s.\n", s1, s2);
    }
    else 
    {
        printf("%s is greater than %s.\n", s1, s2);
    }
    
    return 0;
}
```

## 2.3 运行结果
![](stricmp.png)


# 3. strerror
## 3.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`char *strerror(int errnum);` | 用于将指定的错误码转换为相应的错误信息  |

**参数：**
- **errnum ：** 要转换为错误信息的错误码，通常是由系统调用或库函数返回的错误码。

**返回值：**
指向错误信息字符串的指针，该字符串描述了与错误码相关的错误

## 3.2 演示示例
```c
#include <stdio.h>
#include <string.h>
#include <errno.h>

int main() 
{
    FILE *file = fopen("nonexistent.txt", "r");
    
    if(file == NULL) 
    {
        int errcode = errno; // 获取发生的错误码
        printf("Error opening file: %s\n", strerror(errcode)); // 将错误码转换为错误信息并打印
    } 
    else 
    {
        printf("File opened successfully.\n");
        fclose(file);
    }
    
    return 0;
}
```

## 3.3 运行结果
![](strerror.png)


# 4. strcmpi
## 4.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int strcmpi(const char *s1, const char *s2);` |  用于比较两个字符串的字母序是否相等，忽略大小写 |

**返回值：**
- 如果 `s1` 和 `s2` 代表的字符串相等（忽略大小写），则返回 `0`；
- 如果 `s1` 比 `s2` 小，则返回负数；
- 如果 `s1` 比 `s2` 大，则返回正数。

## 4.2 演示示例
```c
#include <stdio.h>
#include <string.h>

int main() 
{
    char *s1 = "Hello";
    char *s2 = "hElLo";
    
    int result = strcmpi(s1, s2); // 忽略大小写比较
    
    if(result == 0) 
    {
        printf("The strings are equal.\n");
    } 
    else if(result < 0) 
    {
        printf("%s is less than %s.\n", s1, s2);
    } 
    else 
    {
        printf("%s is greater than %s.\n", s1, s2);
    }
    
    return 0;
}
```

看到这里，可能会疑惑 上面不是已经有了 忽略大小写的字符串比较了嘛？

那 `strcmpi` 和 `stricmp` 有什么区别呢？

虽然它们的实现功能相同，但是不同的编译器或操作系统可能会提供其中一个或两个函数。在具备这两个函数的系统中，`strcmpi` 常常作为 `VC（Visual C++）`和 `Borland C++` 的扩展库函数，而 `stricmp` 则是 `POSIX` 标准中定义的函数，在许多 类UNIX系统 上可用。因此，如果需要编写可移植的代码，则应该使用 `stricmp` 函数，而不是 `strcmpi` 函数。

## 4.3 运行结果
![](strcmpi.png)

# 5. strncmp
## 5.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int strncmp(const char *s1, const char *s2, size_t n);` | 用于比较两个字符串的前n个字符是否相等  |

**参数：**
- **s1 ：** 待比较的第一个字符串
- **s2 ：** 待比较的第二个字符串
- **n ：** 要比较的字符数

**返回值：**
- 如果 `s1` 和 `s2` 代表的字符串相等（忽略大小写），则返回 `0`；
- 如果 `s1` 比 `s2` 小，则返回负数；
- 如果 `s1` 比 `s2` 大，则返回正数。

## 5.2 演示示例
```c
#include <stdio.h>
#include <string.h>

int main() 
{
    char *s1 = "Hello World";
    char *s2 = "Hell!";
    
    int result = strncmp(s1, s2, 4); // 比较前4个字符
    
    if(result == 0) 
    {
        printf("The first 4 characters are equal.\n");
    } 
    else if(result < 0) 
    {
        printf("The first 4 characters of %s are less than %s.\n", s1, s2);
    } 
    else 
    {
        printf("The first 4 characters of %s are greater than %s.\n", s1, s2);
    }
    
    return 0;
}
```

## 5.3 运行结果
![](strncmp.png)

# 6. strncmpi
## 6.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int strncmpi(const char *s1, const char *s2, size_t n);` | 用于比较两个字符串的前n个字符是否相等，忽略大小写  |

**参数：**
- **s1 ：** 待比较的第一个字符串
- **s2 ：** 待比较的第二个字符串
- **n ：** 要比较的字符数

**返回值：**
- 如果 `s1` 和 `s2` 代表的字符串相等（忽略大小写），则返回 `0`；
- 如果 `s1` 比 `s2` 小，则返回负数；
- 如果 `s1` 比 `s2` 大，则返回正数。

## 6.2 演示示例
```c
#include <stdio.h>
#include <string.h>

int main() 
{
    char *s1 = "Hello World";
    char *s2 = "heLL!";
    
    int result = strncmpi(s1, s2, 4); // 比较前4个字符，忽略大小写
    
    if(result == 0) 
    {
        printf("The first 4 characters are equal (case insensitive).\n");
    } 
    else if(result < 0) 
    {
        printf("The first 4 characters of %s are less than %s (case insensitive).\n", s1, s2);
    } 
    else 
    {
        printf("The first 4 characters of %s are greater than %s (case insensitive).\n", s1, s2);
    }
    
    return 0;
}
```
**注意：** `strncmpi` 函数不是 `C` 语言标准库中的函数，但在某些编译器或操作系统中可能会提供。

# 7. strncpy
## 7.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`char *strncpy(char *dest, const char *src, size_t n);` | 用于将一个字符串的一部分拷贝到另一个字符串中  |

**参数：**
- **dest ：** 目标字符串
- **src ：** 源字符串
- **n ：** 要拷贝的字符数

## 7.2 演示示例
```c
#include <stdio.h>
#include <string.h>

int main() 
{
    char dest[20];
    const char *src = "Hello";
    
    strncpy(dest, src, 8); // 拷贝前8个字符
    
    printf("The copied string is: %s\n", dest);
    
    return 0;
}
```
当源字符串长度小于 `n` 时，`strncpy()` 函数将在目标字符串的末尾填充 `\0` 字符以达到指定的拷贝长度 `n`

```c
#include <stdio.h>
#include <string.h>

int main() 
{
    char dest[6];
    const char *src = "Hello World";
    
    strncpy(dest, src, 5); // 拷贝前5个字符
    
    dest[5] = '\0'; // 手动添加末尾的\0字符
    
    printf("The copied string is: %s\n", dest);
    
    return 0;
}
```

如果源字符串长度大于等于 `n` 个字符，`strncpy()` 函数将会拷贝源字符串的前 `n` 个字符到目标字符串中，并且不会自动添加末尾的 `\0` 字符。这种情况下，目标字符串可能不是以 `\0` 字符结尾，因此需要手动在拷贝后的目标字符串中添加 `\0` 字符。

## 7.3 运行结果
![](strncpy.png)

![](strncpy1.png)

# 8. strnicmp
## 8.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int strnicmp(const char *s1, const char *s2, size_t n);` | 用于比较两个字符串的前n个字符是否相等，忽略大小写  |

**参数：**
- **s1 ：** 待比较的第一个字符串
- **s2 ：** 待比较的第二个字符串
- **n ：** 要比较的字符数

**返回值：**
- 如果 `s1` 和 `s2` 代表的字符串相等（忽略大小写），则返回 `0`；
- 如果 `s1` 比 `s2` 小，则返回负数；
- 如果 `s1` 比 `s2` 大，则返回正数。

## 8.2 演示示例
```c
#include <stdio.h>
#include <string.h>

int main() 
{
    const char *s1 = "Hello World";
    const char *s2 = "heLL!";
    
    int result = strnicmp(s1, s2, 4); // 比较前4个字符，忽略大小写
    
    if(result == 0) 
    {
        printf("The first 4 characters are equal (case insensitive).\n");
    } 
    else if(result < 0) 
    {
        printf("The first 4 characters of %s are less than %s (case insensitive).\n", s1, s2);
    } 
    else 
    {
        printf("The first 4 characters of %s are greater than %s (case insensitive).\n", s1, s2);
    }
    
    return 0;
}
```

## 8.3 运行结果
![](strnicmp.png)

# 9. strnset
## 9.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`char *strnset(char *str, int c, size_t n);` |  用于将一个字符串的前n个字符都设置为指定的字符 |

**参数：**
- **str：** 要进行操作的字符串
- **c：** 要设置的字符
- **n ：** 要设置的字符数

## 9.2 演示示例
```c
#include <stdio.h>
#include <string.h>

int main() {
    char str[20] = "Hello World";
    printf("Before: %s\n", str);

    strnset(str, '*', 5); // 将前5个字符都设置为*

    printf("After: %s\n", str);

    return 0;
}
```

**注意：** `strnset()` 函数是非标准函数，并不是所有的编译器和操作系统都支持该函数。如果需要跨平台兼容，请使用标准库函数 `memset()` 来实现类似的功能

```c
#include <stdio.h>
#include <string.h>

int main() 
{
    char str[20] = "Hello World";
    printf("Before: %s\n", str);

    memset(str, '*', 5); // 将前5个字符都设置为*

    printf("After: %s\n", str);

    return 0;
}
```
上述示例使用了标准库函数 `memset` 来将 `str` 的前 `5` 个字符都设置为 `*`

## 9.3 运行结果
![](strnset.png)

# 10. strpbrk
## 10.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`char *strpbrk(const char *str1, const char *str2);` | 用于在一个字符串中查找任意给定字符集合中的字符的第一次出现位置 |

**参数：**
- **str1 ：** 要进行搜索的字符串
- **str2 ：** 要查找的字符集合

**注意：** 如果在`str1`中没有找到`str2`中的任何字符，则 `strpbrk` 函数返回`NULL` 指针 

## 10.2 演示示例
```c
#include <stdio.h>
#include <string.h>

int main() 
{
    char str[] = "Hello World";
    char *ptr;

    ptr = strpbrk(str, "aeiou"); // 查找元音字母

    if (ptr != NULL) 
    {
        printf("Found vowel '%c' at position: %lld\n", *ptr, ptr - str + 1);
    } 
    else 
    {
        printf("No vowel found.\n");
    }

    return 0;
}
```

## 10.3 运行结果
![](strpbrk.png)

# 11. strrchr
## 11.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`char *strrchr(const char *str, int character);` | 在给定的字符串中查找指定字符的最后一个匹配项  |

**参数：**
- **str ：** 要进行搜索的字符串
- **character ：** 要查找的字符，其 ASCII 值由整数表示

**返回值：**
- 如果匹配到，返回匹配字符的地址；
- 如果没有找到匹配项，则返回 NULL。

## 11.2 演示示例
```c
#include <stdio.h>
#include <string.h>

int main () 
{
    const char str[] = "Hello World!";

    printf("在 '%s' 中查找字符 'o' 的最后一次出现：\n", str);
    char *last_o = strrchr(str, 'o');

    if (last_o) 
        printf("最后一个 'o' 的位置：%lld\n", last_o - str);
    else 
        printf("未找到匹配的字符。\n");

    printf("在 '%s' 中查找字符 'w' 的最后一次出现：\n", str);
    char *last_w = strrchr(str, 'w');

    if (last_w) 
        printf("最后一个 'w' 的位置：%lld\n", last_w - str);
    else 
        printf("未找到匹配的字符。\n");

    return 0;
}
```

## 11.3 运行结果
![](strrchr.png)


# 12. strrev
## 12.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`char *strrev(char *str);` | 将给定字符串中的所有字符顺序颠倒，并返回颠倒后的字符串  |

**参数：**
- **str ：** 要反转的的字符串

注意：因为 `strrev()` 函数是一个非标准的库函数，许多编译器可能并不支持该函数，所以在使用该函数之前，请确保你的编译器支持它。

## 12.2 演示示例
```c
#include <stdio.h>
#include <string.h>

int main () 
{
    char str[50];

    printf("输入一个字符串：");
    scanf("%s", str);

    printf("原始字符串：%s\n", str);
    printf("反转后的字符串：%s\n", strrev(str));

    return 0;
}
```

## 12.3 运行结果
![](strrev.png)


# 13. strset
## 13.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`char *strset(char *str, int character);` |  用于设置给定字符串中的所有字符为指定的值，并返回修改后的字符串 |

**参数：**
- **str ：** 要修改的字符串
- **character ：** 要设置的字符，其 ASCII 值由整数表示

## 13.2 演示示例
```c
#include <stdio.h>
#include <string.h>

int main () 
{
    char str[50];

    printf("请输入一个字符串：");
    scanf("%s", str);

    printf("原始字符串：%s\n", str);
    printf("将所有字符设置为 'X' 后的字符串：%s\n", strset(str, 'X'));

    return 0;
}
```

## 13.3 运行结果
![](strset.png)


# 14. strspn
## 14.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`size_t strspn(const char *str1, const char *str2);` | 计算字符串 str1 中包含在字符串 str2 中的前缀子字符串长度，并返回该长度值  |

**参数：**
- **str1 ：** 要搜索的字符串
- **str2 ：** 包含要搜索的字符集合的 字符串

## 14.2 演示示例
```c
#include <stdio.h>
#include <string.h>

int main () 
{
    const char str[] = "Hello World!";
    const char charset[] = "lHeoWrd";

    size_t length = strspn(str, charset);

    printf("在 '%s' 中，最长的前缀子串 '%s' 包含于 '%s' 中。长度为 %zu\n", 
          str, strndup(str, length), charset, length);

    return 0;
}
```
上述示例代码运行后，如果出现 `error: 'strndup' was not declared in this scope` ，说明当前的编译器不支持 `strndup()` 函数。

因为 `strndup()` 函数是 `C11` 新增的函数，因此可能不是所有编译器都支持。

那我们就用如下的方式来替换一下 :

```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main () 
{
   const char str[] = "Hello World!";
   const char charset[] = "lHeoWrd";
   
   size_t length = strspn(str, charset);

   char *substr = (char *)malloc(length + 1);
   memcpy(substr, str, length);
   substr[length] = '\0';

   printf("在 '%s' 中，最长的前缀子串 '%s' 包含于 '%s' 中。长度为 %zu\n", 
          str, substr, charset, length);

   free(substr);
   return 0;
}
```

## 14.3 运行结果
![](strspn.png)


# 15. strstr
## 15.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`char *strstr(const char *str1, const char *str2);` | 在字符串 str1 中查找第一个出现的字符串 str2，如果找到了，则返回指向该位置的指针；否则，返回 NULL  |

**参数：**
- **str1 ：** 源字符串
- **str2 ：** 要查找的子字符串

## 15.2 演示示例
```c
#include <stdio.h>
#include <string.h>

int main () 
{
    const char str[] = "Hello World!";
    const char substr[] = "World";

    char *result = strstr(str, substr);

    if (result) 
        printf("在 '%s' 中找到了子串 '%s'。子串起始位置是 '%s'\n", str, substr, result);
    else 
        printf("在 '%s' 中未找到子串 '%s'\n", str, substr);

    return 0;
}
```

## 15.3 运行结果
![](strstr.png)


# 16. strtod
## 16.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double strtod(const char *str, char **endptr);` | 将字符串 str 转换为一个浮点数（double 类型），并返回该浮点数。如果发生了转换错误，则返回 0.0，并且可以通过 endptr 指针返回指向第一个无法转换的字符的指针。  |
**参数：**

- **str ：** 要转换为浮点数的字符串
- **endptr ：** 可选参数，用于存储第一个无法转换的字符的指针地址

## 16.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>

int main () 
{
    const char str[] = "3.14";
    // const char str[] = "3.1a4";
    // const char str[] = "a3.14";
    char *endptr;

    double num = strtod(str, &endptr);

    printf("输入字符串为 '%s'\n", str);
    printf("转换后的浮点数为 %f\n", num);

    if (endptr == str) 
        printf("无法进行任何转换。\n");
    else if (*endptr != '\0') 
        printf("字符串 '%s' 的末尾非法字符是 '%c'\n", str, *endptr);

    return 0;
}
```

## 16.3 运行结果
![](strtod.png)
![](strtod1.png)
![](strtod2.png)


# 17. strtok
## 17.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`char *strtok(char *str, const char *delim);` | 用于将一个字符串分割成多个子字符串  |

**参数：**
- **str ：** 要被分割的 字符串。第一次调用时，应该将其设置为要被分割的字符串的地址；以后的调用中，该参数应该为 NULL
- **delim：** 分隔符字符集合，用于指定子字符串的分割标准

**返回值：** 分割后的第一个子字符串，并在每次调用时修改传入的原始字符串 str，使其指向下一个要被分割的子字符串

## 17.2 演示示例
```c
#include <stdio.h>
#include <string.h>

int main () 
{
    const char str[] = "Hello,World!How are you?";
    const char delim[] = ", !";
    char *token;

    token = strtok((char *)str, delim);

    while (token != NULL) 
    {
        printf("%s\n", token);
        token = strtok(NULL, delim);
    }

    return 0;
}
```

在上述的示例中，
- 我们首先定义了字符串 `str` 和 字符集合 `delim` ；
- 然后使用 `strtok()` 函数将字符串 `str` 按照字符集合 `delim` 中的分隔符进行分割。每次调用 `strtok()` 函数时，它会返回分割出的第一个子字符串，并且会修改 `str` 指向下一个将要被分割的子字符串。
- 最后我们不断循环调用 `strtok()` 函数，并输出返回的每个子字符串，直到没有更多的子字符串可以分割为止

## 17.3 运行结果
![](strtok.png)


# 18. strtol
## 18.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`long int strtol(const char *str, char **endptr, int base);` |  将字符串 str 转换为一个长整型数（long int 类型） |

**参数：**
- **str ：** 要被转换为长整型数的字符串
- **endptr ：** 可选参数，用于存储第一个无法转换的字符的指针地址
- **base：** 转换的进制数，必须是 `2` 到 `36` 之间的有效数字或者 `0`。为 `0` 表示采用默认的进制数，即可以解释成合法的整数的最大进制数（一般是 `10`）

**返回值：** 
- 如果成功转换，返回一个长整型数；
- 如果发生了转换错误，则返回 `0`，并且可以通过 `endptr` 指针返回指向第一个无法转换的字符的指针。

## 18.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>

int main () 
{
    const char str[] = "101011";
    // const char str[] = "10101a1";
    // const char str[] = "a101011";
    char *endptr;

    long int num = strtol(str, &endptr, 2);

    printf("输入字符串为 '%s'\n", str);
    printf("转换后的十进制数为 %ld\n", num);

    if (endptr == str) 
        printf("无法进行任何转换。\n");
    else if (*endptr != '\0') 
        printf("字符串 '%s' 的末尾非法字符是 '%c'\n", str, *endptr);

    return 0;
}
```

## 18.3 运行结果
![](strtol.png)
![](strtol_1.png)
![](strtol_2.png)


# 19. strupr
## 19.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`char *strupr(char *str);` | 将字符串 str 中的所有小写字母转换为大写字母，并返回指向该字符串的指针  |

**参数：**
- **str ：** 要进行大小写转换的字符串

## 19.2 演示示例
```c
#include <stdio.h>
#include <string.h>

int main () 
{
    char str[] = "Hello, World!";
    printf("原始字符串为 '%s'\n", str);
    strupr(str);
    printf("转换后的字符串为 '%s'\n", str);
    return 0;
}
```

## 19.3 运行结果
![](strupr.png)


# 20. swab
## 20.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `void swab(const void *src, void *dest, ssize_t nbytes);`| 将源缓冲区中的每两个相邻字节进行互换，然后将结果存储到目标缓冲区中  |

**参数：**
- **str ：** 源缓冲区的指针
- **dest ：** 目标缓冲区的指针
- **nbytes ：** 需要交换的字节数

## 20.2 演示示例
```c
#include <stdio.h>
#include <string.h>
#include <unistd.h>

int main () 
{
    char src[] = "Hello, World!";
    char dest[15];
    printf("原始字符串为 '%s'\n", src);
    swab(src, dest, strlen(src));
    printf("转换后的字符串为 '%s'\n", dest);
    return 0;
}
```

上面示例中，我在一开始演示时，因为没有加上 `#include <unistd.h>` 这个头文件，导致出现如下错误：
![](swab_error.png)

## 20.3 运行结果

![](swab.png)

# 21. system
## 21.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`int system(const char *command);` |  执行一个 `shell` 或 `cmd` 命令，并等待命令的完成 |

**参数：**
- **command ：** 要执行的 `shell` 或 `cmd` 命令。

## 21.2 演示示例
下面演示 执行 `ls -l` 的 `shell` 命令，笔者是在 `windows`环境下运行，故会出错，详见运行结果那里

```c
#include <stdio.h>
#include <stdlib.h>

int main () 
{
    const char command[] = "ls -l";

    printf("执行命令: '%s'\n", command);

    int status = system(command);

    if (status != 0) 
        printf("执行出错！\n");

    return 0;
}
```

再来看下，演示使用 `dir` 命令，在 `windows` 下可以输出当前目录下的所有文件和子目录
```c
#include <stdlib.h>

int main () 
{
    const char command[] = "dir";
    system(command);
    return 0;
}
```

## 21.3 运行结果
![](system.png)

![](system1.png)
