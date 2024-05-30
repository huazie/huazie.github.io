---
title: C语言函数大全--c开头的函数之复数篇
date: 2023-03-21 16:15:24 
updated: 2023-06-25 23:23:16
categories:
  - 开发语言-C
tags:
  - C语言函数大全
  - c开头的函数
  - 复数篇
---

![](/images/cplus-logo.png)

# 总览

| 函数声明 |  函数功能  |
|:--|:--|
| `double cabs (double complex z);`| 计算复数 z 的绝对值（double）  |
|`float cabsf (float complex z);`| 计算复数 z 的绝对值（float）|
|`long double cabsl (long double complex z);`|计算复数 z 的绝对值（long double）|
|`double creal (double complex z);` |  计算复数z的实部（double）|
|`float crealf (float complex z);`|计算复数z的实部（float）|
|`long double creall (long double complex z);`|计算复数z的实部（long double）|
|`double cimag (double complex z);` | 计算复数z的虚部（double）  |
|`float cimagf (float complex z);` |  计算复数z的虚部（float） |
|`long double cimagl (long double complex z);` |  计算复数z的虚部（long double） |
|`double carg (double complex z);` |  计算复数z的相位角 （double）|
|`float cargf (float complex z);`|计算复数z的相位角（float）|
|`long double cargl (long double complex z);`|计算复数z的相位角（long double）|
|`double complex cacos (double complex z);` |  计算复数z的反余弦 （double complex）|
|`float complex cacosf (float complex z);`|计算复数z的反余弦（float complex）|
|`long double complex cacosl (long double complex z);`|计算复数z的反余弦（long double complex）|
|`double complex cacosh (double complex z);`|计算复数z的反双曲余弦（double complex）|
|`float complex cacoshf (float complex z);`|计算复数z的反双曲余弦（float complex）|
|`long double complex cacoshl (long double complex z);`|计算复数z的反双曲余弦（long double complex）|
| `double complex casin (double complex z);` | 计算复数z的反正弦（double complex）  |
| `float complex casinf (float complex z);` | 计算复数z的反正弦（float complex）  |
| `long double complex casinl (long double complex z);` | 计算复数z的反正弦（long double complex）  |
| `double complex casinh (double complex z);` | 计算复数z的反双曲正弦（double complex）  |
| `float complex casinhf (float complex z);` | 计算复数z的反双曲正弦（float complex）  |
| `long double complex casinhl (long double complex z);` | 计算复数z的反双曲正弦（long double complex） |
| `double complex catan (double complex z);` | 计算复数z的反正切（double complex）  |
| `float complex catanf (float complex z);` | 计算复数z的反正切（float complex）  |
| `long double complex catanl (long double complex z);` | 计算复数z的反正切（long double complex） |
| `double complex catanh (double complex z);` | 计算复数z的反双曲正切（double complex）  |
| `float complex catanhf (float complex z);` | 计算复数z的反双曲正切（float complex）  |
| `long double complex catanhl (long double complex z);` | 计算复数z的反双曲正切（long double complex） |
| `double complex ccos (double complex z);` | 计算复数z的余弦（double complex）  |
| `float complex ccosf (float complex z);` |计算复数z的余弦（float complex）  |
| `long double complex ccosl (long double complex z);` | 计算复数z的余弦（long double complex） |
| `double complex ccosh (double complex z);` | 计算复数z的双曲余弦（double complex）  |
| `float complex ccoshf (float complex z);` |计算复数z的双曲余弦（float complex）  |
| `long double complex ccoshl (long double complex z);` | 计算复数z的双曲余弦（long double complex） |
| `double complex csin (double complex z);` | 计算复数z的正弦（double complex）  |
| `float complex csinf (float complex z);` |计算复数z的正弦（float complex）  |
| `long double complex csinl (long double complex z);` | 计算复数z的正弦（long double complex） |
| `double complex csinh (double complex z);` | 计算复数z的双曲正弦（double complex）  |
| `float complex csinhf (float complex z);` |计算复数z的双曲正弦（float complex）  |
| `long double complex csinhl (long double complex z);` | 计算复数z的双曲正弦（long double complex） |
| `double complex ctan (double complex z);` | 计算复数z的正切（double complex）  |
| `float complex ctanf (float complex z);` |计算复数z的正切（float complex）  |
| `long double complex ctanl (long double complex z);` | 计算复数z的正切（long double complex） |
| `double complex ctanh (double complex z);` | 计算复数z的双曲正切（double complex）  |
| `float complex ctanhf (float complex z);` |计算复数z的双曲正切（float complex）  |
| `long double complex ctanhl (long double complex z);` | 计算复数z的双曲正切（long double complex） |
| `double complex cexp (double complex z);` | 计算复数z的指数基数e（double complex）  |
| `float complex cexpf (float complex z);` |计算复数z的指数基数e（float complex）  |
| `long double complex cexpl (long double complex z);` | 计算复数z的指数基数e（long double complex） |
| `double complex clog (double complex z);` | 计算复数z的自然对数（以e为底）（double complex）  |
| `float complex clogf (float complex z);` |计算复数z的自然对数（以e为底）（float complex）  |
| `long double complex clogl (long double complex z);` | 计算复数z的自然对数（以e为底）（long double complex） |
| `double complex conj (double complex z);` | 计算复数z的[共轭](https://baike.baidu.com/item/%E5%85%B1%E8%BD%AD/31802?fr=aladdin)（double complex）  |
| `float complex conjf (float complex z);` |计算复数z的共轭（float complex）  |
| `long double complex conjl (long double complex z);` | 计算复数z的共轭（long double complex） |
|`double complex cpow (double complex x, double complex y);` |  计算x的y次方值 （double complex） |
|`float complex cpowf (float complex x, float complex y);` | 计算x的y次方值 （float complex）  |
|`long double complex cpowl (long double complex x, long double complex y);` |  计算x的y次方值 （double complex） |
| `double complex cproj (double complex z);` | 计算复数z在黎曼球面上的投影（double complex）  |
| `float complex cprojf (float complex z);` |计算复数z在黎曼球面上的投影（float complex）  |
| `long double complex cprojl (long double complex z);` | 计算复数z在黎曼球面上的投影（long double complex） |
| `double complex csqrt (double complex z);` | 计算复数z的平方根（double complex）  |
| `float complex csqrtf (float complex z);` |计算复数z的平方根（float complex）  |
| `long double complex csqrtl (long double complex z);` | 计算复数z的平方根（long double complex） |


# 1. cabs，cabsf，cabsl
## 1.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double cabs (double complex z);`| 计算复数 z 的绝对值（double）  |
|`float cabsf (float complex z);`| 计算复数 z 的绝对值（float）|
|`long double cabsl (long double complex z);`|计算复数 z 的绝对值（long double）|

## 1.2 演示示例
```c
// Huazie
#include <stdio.h>
#include <complex.h>

int main(void)
{
    double complex z;
    double x = 2.0, y = 2.0, val;
    z = x + y * I; // I 代指 虚数单位 i
    val = cabs(z); // 计算复数 z 的绝对值

    float complex zf;
    float xf = 2.0, yf = 2.0, valf;
    zf = xf + yf * I;
    valf = cabsf(zf);

    long double complex zL;
    long double xL = 2.0, yL = 2.0, valL;
    zL = xL + yL * I;
    valL = cabsl(zL);
    
    printf("The absolute value of (%.4lf + %.4lfi) is %.20lf\n", x, y, val);
    printf("The absolute value of (%.4f + %.4fi) is %.20f\n", xf, yf, valf);
    printf("The absolute value of (%.4Lf + %.4Lfi) is %.20Lf", xL, yL, valL);

    return 0;
}
```
## 1.3 运行结果
![](cabs.png)


# 2. creal，crealf，creall
## 2.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double creal (double complex z);` |  计算复数z的实部（double）|
|`float crealf (float complex z);`|计算复数z的实部（float）|
|`long double creall (long double complex z);`|计算复数z的实部（long double）|


## 2.2 演示示例
```c
// Huazie
#include <stdio.h>
#include <complex.h>

int main(void)
{
    double complex z;
    double x = 2.0, y = 1.0;
    z = x + y * I; // I 代指 虚数单位 i

    float complex zf;
    float xf = 3.0, yf = 1.0;
    zf = xf + yf * I;

    long double complex zL;
    long double xL = 4.0, yL = 1.0;
    zL = xL + yL * I;

    printf("The real part of (%.4lf + %.4lfi) is %.4lf\n", x, y, creal(z));
    printf("The real part of (%.4f + %.4fi) is %.4f\n", xf, yf, crealf(zf));
    printf("The real part of (%.4Lf + %.4Lfi) is %.4Lf", xL, yL, creall(zL));

    return 0;
}
```
## 2.3 运行结果
![](creal.png)


# 3. cimag，cimagf，cimagl
## 3.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double cimag (double complex z);` | 计算复数z的虚部（double）  |
|`float cimagf (float complex z);` |  计算复数z的虚部（float） |
|`long double cimagl (long double complex z);` |  计算复数z的虚部（long double） |

## 3.2 演示示例
```c
// Huazie
#include <stdio.h>
#include <complex.h>

int main(void)
{
    double complex z;
    double x = 1.0, y = 2.0;
    z = x + y * I; // I 代指 虚数单位 i

    float complex zf;
    float xf = 1.0, yf = 3.0;
    zf = xf + yf * I;

    long double complex zL;
    long double xL = 1.0, yL = 4.0;
    zL = xL + yL * I;

    printf("The imaginary part of (%.4lf + %.4lfi) is %.4lf\n", x, y, cimag(z));
    printf("The imaginary part of (%.4f + %.4fi) is %.4f\n", xf, yf, cimagf(zf));
    printf("The imaginary part of (%.4Lf + %.4Lfi) is %.4Lf", xL, yL, cimagl(zL));

    return 0;
}
```
## 3.3 运行结果
![](cimag.png)

# 4. carg，cargf，cargl
## 4.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double carg (double complex z);` |  计算复数z的相位角 （double）|
|`float cargf (float complex z);`|计算复数z的相位角（float）|
|`long double cargl (long double complex z);`|计算复数z的相位角（long double）|

相位角是描述波形在时间轴上的位置的一个重要参数，它决定了波形的起始位置和变化状态。在实际应用中，相位角的测量和控制对于电路设计和信号处理至关重要。通过对相位角的理解和应用，可以更好地分析和控制波动现象，从而实现对电力系统和通信系统的优化。

## 4.2 演示示例
```c
#include <stdio.h>
#include <complex.h>

int main(void)
{
    double complex z;
    z = 1.0 + 2.0 * I; // I 代指 虚数单位 i

    float complex zf;
    zf = 1.0f + 2.0f * I;

    long double complex zL;
    zL = (long double) 1.0 + (long double) 2.0 * I;

    printf("The phase angle of (%.4lf + %.4lfi) is %.60lf\n", creal(z), cimag(z), carg(z));
    printf("The phase angle of (%.4f + %.4fi) is %.60f\n", crealf(zf), cimagf(zf), cargf(zf));
    printf("The phase angle of (%.4Lf + %.4Lfi) is %.60Lf", creall(zL), cimagl(zL), cargl(zL));

    return 0;
}
```
## 4.3 运行结果
![](carg.png)


# 5. cacos，cacosf，cacosl
## 5.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double complex cacos (double complex z);` |  计算复数z的反余弦 （double complex）|
|`float complex cacosf (float complex z);`|计算复数z的反余弦（float complex）|
|`long double complex cacosl (long double complex z);`|计算复数z的反余弦（long double complex）|


## 5.2 演示示例
```c
#include <stdio.h>
#include <math.h>
#include <complex.h>

int main(void)
{
    double complex z, zcacos;
    z = 1.0 + 2.0 * I; // I 代指 虚数单位 i
    zcacos = cacos(z); // 计算复数z的反余弦

    float complex zf, zcacosf;
    zf = 1.0f + 2.0f * I;
    zcacosf = cacosf(zf);

    long double complex zL, zcacosl;
    zL = (long double) 1.0 + (long double) 2.0 * I;
    zcacosl = cacosl(zL);

    double zimag = cimag(zcacos);
    float zimagf = cimagf(zcacosf);
    long double zimagl = cimagl(zcacosl);
    if (zimag < 0) 
        printf("The arc cosine of (%.4lf + %.4lfi) is (%.20lf - %.20lfi)\n", creal(z), cimag(z), creal(zcacos), fabs(zimag));
    else 
        printf("The arc cosine of (%.4lf + %.4lfi) is (%.20lf + %.20lfi)\n", creal(z), cimag(z), creal(zcacos), zimag);       

    if (zimagf < 0) 
        printf("The arc cosine of (%.4f + %.4fi) is (%.20f - %.20fi)\n", crealf(zf), cimagf(zf), crealf(zcacosf), fabsf(zimagf));
    else 
        printf("The arc cosine of (%.4f + %.4fi) is (%.20f + %.20fi)\n", crealf(zf), cimagf(zf), crealf(zcacosf), zimagf);

    if (zimagl < 0) 
        printf("The arc cosine of (%.4Lf + %.4Lfi) is (%.20Lf - %.20Lfi)", creall(zL), cimagl(zL), creall(zcacosl), fabsl(zimagl));
    else 
        printf("The arc cosine of (%.4Lf + %.4Lfi) is (%.20Lf + %.20Lfi)", creall(zL), cimagl(zL), creall(zcacosl), zimagl);
    return 0;
}
```
## 5.3 运行结果
![](cacos.png)



# 6. cacosh，cacoshf，cacoshl
## 6.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double complex cacosh (double complex z);`|计算复数z的反双曲余弦（double complex）|
|`float complex cacoshf (float complex z);`|计算复数z的反双曲余弦（float complex）|
|`long double complex cacoshl (long double complex z);`|计算复数z的反双曲余弦（long double complex）|

## 6.2 演示示例
```c
#include <stdio.h>
#include <math.h>
#include <complex.h>

int main(void)
{
    double complex z, zcacosh;
    z = 1.0 + 2.0 * I; // I 代指 虚数单位 i
    zcacosh = cacosh(z); // 反双曲余弦

    float complex zf, zcacoshf;
    zf = 1.0f + 2.0f * I;
    zcacoshf = cacoshf(zf);

    long double complex zL, zcacoshl;
    zL = (long double) 1.0 + (long double) 2.0 * I;
    zcacoshl = cacoshl(zL);

    double zimag = cimag(zcacosh);
    float zimagf = cimagf(zcacoshf);
    long double zimagl = cimagl(zcacoshl);
    if (zimag < 0) 
        printf("The inverse hyperbolic cosine of (%.4lf + %.4lfi) is (%.20lf - %.20lfi)\n", creal(z), cimag(z), creal(zcacosh), fabs(zimag));
    else 
        printf("The inverse hyperbolic cosine of (%.4lf + %.4lfi) is (%.20lf + %.20lfi)\n", creal(z), cimag(z), creal(zcacosh), zimag);       

    if (zimagf < 0) 
        printf("The inverse hyperbolic cosine of (%.4f + %.4fi) is (%.20f - %.20fi)\n", crealf(zf), cimagf(zf), crealf(zcacoshf), fabsf(zimagf));
    else 
        printf("The inverse hyperbolic cosine of (%.4f + %.4fi) is (%.20f + %.20fi)\n", crealf(zf), cimagf(zf), crealf(zcacoshf), zimagf);

    if (zimagl < 0) 
        printf("The inverse hyperbolic cosine of (%.4Lf + %.4Lfi) is (%.20Lf - %.20Lfi)", creall(zL), cimagl(zL), creall(zcacoshl), fabsl(zimagl));
    else 
        printf("The inverse hyperbolic cosine of (%.4Lf + %.4Lfi) is (%.20Lf + %.20Lfi)", creall(zL), cimagl(zL), creall(zcacoshl), zimagl);
    return 0;
}
```
## 6.3 运行结果
![](cacosh.png)



# 7. casin，casinf，casinl
## 7.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double complex casin (double complex z);` | 计算复数z的反正弦（double complex）  |
| `float complex casinf (float complex z);` | 计算复数z的反正弦（float complex）  |
| `long double complex casinl (long double complex z);` | 计算复数z的反正弦（long double complex）  |
## 7.2 演示示例
```c
#include <stdio.h>
#include <math.h>
#include <complex.h>

int main(void)
{
    double complex z, zcasin;
    z = 1.0 + 2.0 * I; // I 代指 虚数单位 i
    zcasin = casin(z); // 反正弦

    float complex zf, zcasinf;
    zf = 1.0f + 2.0f * I;
    zcasinf = casinf(zf);

    long double complex zL, zcasinl;
    zL = (long double) 1.0 + (long double) 2.0 * I;
    zcasinl = casinl(zL);

    double zimag = cimag(zcasin);
    float zimagf = cimagf(zcasinf);
    long double zimagl = cimagl(zcasinl);
    if (zimag < 0) 
        printf("The arcsine of (%.4lf + %.4lfi) is (%.20lf - %.20lfi)\n", creal(z), cimag(z), creal(zcasin), fabs(zimag));
    else 
        printf("The arcsine of (%.4lf + %.4lfi) is (%.20lf + %.20lfi)\n", creal(z), cimag(z), creal(zcasin), zimag);       

    if (zimagf < 0) 
        printf("The arcsine of (%.4f + %.4fi) is (%.20f - %.20fi)\n", crealf(zf), cimagf(zf), crealf(zcasinf), fabsf(zimagf));
    else 
        printf("The arcsine of (%.4f + %.4fi) is (%.20f + %.20fi)\n", crealf(zf), cimagf(zf), crealf(zcasinf), zimagf);

    if (zimagl < 0) 
        printf("The arcsine of (%.4Lf + %.4Lfi) is (%.20Lf - %.20Lfi)", creall(zL), cimagl(zL), creall(zcasinl), fabsl(zimagl));
    else 
        printf("The arcsine of (%.4Lf + %.4Lfi) is (%.20Lf + %.20Lfi)", creall(zL), cimagl(zL), creall(zcasinl), zimagl);
    return 0;
}
```
## 7.3 运行结果
![](casin.png)



# 8. casinh，casinhf，casinhl
## 8.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double complex casinh (double complex z);` | 计算复数z的反双曲正弦（double complex）  |
| `float complex casinhf (float complex z);` | 计算复数z的反双曲正弦（float complex）  |
| `long double complex casinhl (long double complex z);` | 计算复数z的反双曲正弦（long double complex） |

## 8.2 演示示例
```c
#include <stdio.h>
#include <math.h>
#include <complex.h>

int main(void)
{
    double complex z, zcasinh;
    z = 1.0 + 2.0 * I; // I 代指 虚数单位 i
    zcasinh = casinh(z); // 反双曲正弦

    float complex zf, zcasinhf;
    zf = 1.0f + 2.0f * I;
    zcasinhf = casinhf(zf);

    long double complex zL, zcasinhl;
    zL = (long double) 1.0 + (long double) 2.0 * I;
    zcasinhl = casinhl(zL);

    double zimag = cimag(zcasinh);
    float zimagf = cimagf(zcasinhf);
    long double zimagl = cimagl(zcasinhl);
    if (zimag < 0) 
        printf("The inverse hyperbolic sine of (%.4lf + %.4lfi) is (%.20lf - %.20lfi)\n", creal(z), cimag(z), creal(zcasinh), fabs(zimag));
    else 
        printf("The inverse hyperbolic sine of (%.4lf + %.4lfi) is (%.20lf + %.20lfi)\n", creal(z), cimag(z), creal(zcasinh), zimag);       

    if (zimagf < 0) 
        printf("The inverse hyperbolic sine of (%.4f + %.4fi) is (%.20f - %.20fi)\n", crealf(zf), cimagf(zf), crealf(zcasinhf), fabsf(zimagf));
    else 
        printf("The inverse hyperbolic sine of (%.4f + %.4fi) is (%.20f + %.20fi)\n", crealf(zf), cimagf(zf), crealf(zcasinhf), zimagf);

    if (zimagl < 0) 
        printf("The inverse hyperbolic sine of (%.4Lf + %.4Lfi) is (%.20Lf - %.20Lfi)", creall(zL), cimagl(zL), creall(zcasinhl), fabsl(zimagl));
    else 
        printf("The inverse hyperbolic sine of (%.4Lf + %.4Lfi) is (%.20Lf + %.20Lfi)", creall(zL), cimagl(zL), creall(zcasinhl), zimagl);
    return 0;
}
```
## 8.3 运行结果
![](casinh.png)



# 9. catan，catanf，catanl
## 9.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double complex catan (double complex z);` | 计算复数z的反正切（double complex）  |
| `float complex catanf (float complex z);` | 计算复数z的反正切（float complex）  |
| `long double complex catanl (long double complex z);` | 计算复数z的反正切（long double complex） |

## 9.2 演示示例
```c
#include <stdio.h>
#include <math.h>
#include <complex.h>

int main(void)
{
    double complex z, zcatan;
    z = 1.0 + 2.0 * I; // I 代指 虚数单位 i
    zcatan = catan(z); // 反正切

    float complex zf, zcatanf;
    zf = 1.0f + 2.0f * I;
    zcatanf = catanf(zf);

    long double complex zL, zcatanl;
    zL = (long double) 1.0 + (long double) 2.0 * I;
    zcatanl = catanl(zL);

    double zimag = cimag(zcatan);
    float zimagf = cimagf(zcatanf);
    long double zimagl = cimagl(zcatanl);
    if (zimag < 0) 
        printf("The arc tangent of (%.4lf + %.4lfi) is (%.20lf - %.20lfi)\n", creal(z), cimag(z), creal(zcatan), fabs(zimag));
    else 
        printf("The arc tangent of (%.4lf + %.4lfi) is (%.20lf + %.20lfi)\n", creal(z), cimag(z), creal(zcatan), zimag);       

    if (zimagf < 0) 
        printf("The arc tangent of (%.4f + %.4fi) is (%.20f - %.20fi)\n", crealf(zf), cimagf(zf), crealf(zcatanf), fabsf(zimagf));
    else 
        printf("The arc tangent of (%.4f + %.4fi) is (%.20f + %.20fi)\n", crealf(zf), cimagf(zf), crealf(zcatanf), zimagf);

    if (zimagl < 0) 
        printf("The arc tangent of (%.4Lf + %.4Lfi) is (%.20Lf - %.20Lfi)", creall(zL), cimagl(zL), creall(zcatanl), fabsl(zimagl));
    else 
        printf("The arc tangent of (%.4Lf + %.4Lfi) is (%.20Lf + %.20Lfi)", creall(zL), cimagl(zL), creall(zcatanl), zimagl);
    return 0;
}
```
## 9.3 运行结果
![](catan.png)


# 10. catanh，catanhf，catanhl
## 10.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double complex catanh (double complex z);` | 计算复数z的反双曲正切（double complex）  |
| `float complex catanhf (float complex z);` | 计算复数z的反双曲正切（float complex）  |
| `long double complex catanhl (long double complex z);` | 计算复数z的反双曲正切（long double complex） |

## 10.2 演示示例
```c
#include <stdio.h>
#include <math.h>
#include <complex.h>

int main(void)
{
    double complex z, zcatanh;
    z = 1.0 + 2.0 * I; // I 代指 虚数单位 i
    zcatanh = catanh(z); // 反双曲正切

    float complex zf, zcatanhf;
    zf = 1.0f + 2.0f * I;
    zcatanhf = catanhf(zf);

    long double complex zL, zcatanhl;
    zL = (long double) 1.0 + (long double) 2.0 * I;
    zcatanhl = catanhl(zL);

    double zimag = cimag(zcatanh);
    float zimagf = cimagf(zcatanhf);
    long double zimagl = cimagl(zcatanhl);
    if (zimag < 0) 
        printf("The inverse hyperbolic tangent of (%.4lf + %.4lfi) is (%.20lf - %.20lfi)\n", creal(z), cimag(z), creal(zcatanh), fabs(zimag));
    else 
        printf("The inverse hyperbolic tangent of (%.4lf + %.4lfi) is (%.20lf + %.20lfi)\n", creal(z), cimag(z), creal(zcatanh), zimag);       

    if (zimagf < 0) 
        printf("The inverse hyperbolic tangent of (%.4f + %.4fi) is (%.20f - %.20fi)\n", crealf(zf), cimagf(zf), crealf(zcatanhf), fabsf(zimagf));
    else 
        printf("The inverse hyperbolic tangent of (%.4f + %.4fi) is (%.20f + %.20fi)\n", crealf(zf), cimagf(zf), crealf(zcatanhf), zimagf);

    if (zimagl < 0) 
        printf("The inverse hyperbolic tangent of (%.4Lf + %.4Lfi) is (%.20Lf - %.20Lfi)", creall(zL), cimagl(zL), creall(zcatanhl), fabsl(zimagl));
    else 
        printf("The inverse hyperbolic tangent of (%.4Lf + %.4Lfi) is (%.20Lf + %.20Lfi)", creall(zL), cimagl(zL), creall(zcatanhl), zimagl);
    return 0;
}
```
## 10.3 运行结果
![](catanh.png)


# 11. ccos，ccosf，ccosl
## 11.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double complex ccos (double complex z);` | 计算复数z的余弦（double complex）  |
| `float complex ccosf (float complex z);` |计算复数z的余弦（float complex）  |
| `long double complex ccosl (long double complex z);` | 计算复数z的余弦（long double complex） |

## 11.2 演示示例
```c
#include <stdio.h>
#include <math.h>
#include <complex.h>

int main(void)
{
    double complex z, zccos;
    z = 1.0 + 2.0 * I; // I 代指 虚数单位 i
    zccos = ccos(z); // 余弦

    float complex zf, zccosf;
    zf = 1.0f + 2.0f * I;
    zccosf = ccosf(zf);

    long double complex zL, zccosl;
    zL = (long double) 1.0 + (long double) 2.0 * I;
    zccosl = ccosl(zL);

    double zimag = cimag(zccos);
    float zimagf = cimagf(zccosf);
    long double zimagl = cimagl(zccosl);
    if (zimag < 0) 
        printf("The cosine of (%.4lf + %.4lfi) is (%.20lf - %.20lfi)\n", creal(z), cimag(z), creal(zccos), fabs(zimag));
    else 
        printf("The cosine of (%.4lf + %.4lfi) is (%.20lf + %.20lfi)\n", creal(z), cimag(z), creal(zccos), zimag);       

    if (zimagf < 0) 
        printf("The cosine of (%.4f + %.4fi) is (%.20f - %.20fi)\n", crealf(zf), cimagf(zf), crealf(zccosf), fabsf(zimagf));
    else 
        printf("The cosine of (%.4f + %.4fi) is (%.20f + %.20fi)\n", crealf(zf), cimagf(zf), crealf(zccosf), zimagf);

    if (zimagl < 0) 
        printf("The cosine of (%.4Lf + %.4Lfi) is (%.20Lf - %.20Lfi)", creall(zL), cimagl(zL), creall(zccosl), fabsl(zimagl));
    else 
        printf("The cosine of (%.4Lf + %.4Lfi) is (%.20Lf + %.20Lfi)", creall(zL), cimagl(zL), creall(zccosl), zimagl);
    return 0;
}
```
## 11.3 运行结果
![](ccos.png)

# 12. ccosh，ccoshf，ccoshl
## 12.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double complex ccosh (double complex z);` | 计算复数z的双曲余弦（double complex）  |
| `float complex ccoshf (float complex z);` |计算复数z的双曲余弦（float complex）  |
| `long double complex ccoshl (long double complex z);` | 计算复数z的双曲余弦（long double complex） |

## 12.2 演示示例
```c
#include <stdio.h>
#include <math.h>
#include <complex.h>

int main(void)
{
    double complex z, zccosh;
    z = 1.0 + 2.0 * I; // I 代指 虚数单位 i
    zccosh = ccosh(z); // 双曲余弦

    float complex zf, zccoshf;
    zf = 1.0f + 2.0f * I;
    zccoshf = ccoshf(zf);

    long double complex zL, zccoshl;
    zL = (long double) 1.0 + (long double) 2.0 * I;
    zccoshl = ccoshl(zL);

    double zimag = cimag(zccosh);
    float zimagf = cimagf(zccoshf);
    long double zimagl = cimagl(zccoshl);
    if (zimag < 0) 
        printf("The hyperbolic cosine of (%.4lf + %.4lfi) is (%.20lf - %.20lfi)\n", creal(z), cimag(z), creal(zccosh), fabs(zimag));
    else 
        printf("The hyperbolic cosine of (%.4lf + %.4lfi) is (%.20lf + %.20lfi)\n", creal(z), cimag(z), creal(zccosh), zimag);       

    if (zimagf < 0) 
        printf("The hyperbolic cosine of (%.4f + %.4fi) is (%.20f - %.20fi)\n", crealf(zf), cimagf(zf), crealf(zccoshf), fabsf(zimagf));
    else 
        printf("The hyperbolic cosine of (%.4f + %.4fi) is (%.20f + %.20fi)\n", crealf(zf), cimagf(zf), crealf(zccoshf), zimagf);

    if (zimagl < 0) 
        printf("The hyperbolic cosine of (%.4Lf + %.4Lfi) is (%.20Lf - %.20Lfi)", creall(zL), cimagl(zL), creall(zccoshl), fabsl(zimagl));
    else 
        printf("The hyperbolic cosine of (%.4Lf + %.4Lfi) is (%.20Lf + %.20Lfi)", creall(zL), cimagl(zL), creall(zccoshl), zimagl);
    return 0;
}
```
## 12.3 运行结果

![](ccosh.png)


# 13. csin，csinf，csinl
## 13.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double complex csin (double complex z);` | 计算复数z的正弦（double complex）  |
| `float complex csinf (float complex z);` |计算复数z的正弦（float complex）  |
| `long double complex csinl (long double complex z);` | 计算复数z的正弦（long double complex） |

## 13.2 演示示例
```c
#include <stdio.h>
#include <math.h>
#include <complex.h>

int main(void)
{
    double complex z, zcsin;
    z = 1.0 + 2.0 * I; // I 代指 虚数单位 i
    zcsin = csin(z); // 正弦

    float complex zf, zcsinf;
    zf = 1.0f + 2.0f * I;
    zcsinf = csinf(zf);

    long double complex zL, zcsinl;
    zL = (long double) 1.0 + (long double) 2.0 * I;
    zcsinl = csinl(zL);

    double zimag = cimag(zcsin);
    float zimagf = cimagf(zcsinf);
    long double zimagl = cimagl(zcsinl);
    if (zimag < 0) 
        printf("The sine of (%.4lf + %.4lfi) is (%.20lf - %.20lfi)\n", creal(z), cimag(z), creal(zcsin), fabs(zimag));
    else 
        printf("The sine of (%.4lf + %.4lfi) is (%.20lf + %.20lfi)\n", creal(z), cimag(z), creal(zcsin), zimag);       

    if (zimagf < 0) 
        printf("The sine of (%.4f + %.4fi) is (%.20f - %.20fi)\n", crealf(zf), cimagf(zf), crealf(zcsinf), fabsf(zimagf));
    else 
        printf("The sine of (%.4f + %.4fi) is (%.20f + %.20fi)\n", crealf(zf), cimagf(zf), crealf(zcsinf), zimagf);

    if (zimagl < 0) 
        printf("The sine of (%.4Lf + %.4Lfi) is (%.20Lf - %.20Lfi)", creall(zL), cimagl(zL), creall(zcsinl), fabsl(zimagl));
    else 
        printf("The sine of (%.4Lf + %.4Lfi) is (%.20Lf + %.20Lfi)", creall(zL), cimagl(zL), creall(zcsinl), zimagl);
    return 0;
}
```
## 13.3 运行结果
![](csin.png)


# 14. csinh，csinhf，csinhl
## 14.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double complex csinh (double complex z);` | 计算复数z的双曲正弦（double complex）  |
| `float complex csinhf (float complex z);` |计算复数z的双曲正弦（float complex）  |
| `long double complex csinhl (long double complex z);` | 计算复数z的双曲正弦（long double complex） |

## 14.2 演示示例
```c
#include <stdio.h>
#include <math.h>
#include <complex.h>

int main(void)
{
    double complex z, zcsinh;
    z = 1.0 + 2.0 * I; // I 代指 虚数单位 i
    zcsinh = csinh(z); // 双曲正弦

    float complex zf, zcsinhf;
    zf = 1.0f + 2.0f * I;
    zcsinhf = csinhf(zf);

    long double complex zL, zcsinhl;
    zL = (long double) 1.0 + (long double) 2.0 * I;
    zcsinhl = csinhl(zL);

    double zimag = cimag(zcsinh);
    float zimagf = cimagf(zcsinhf);
    long double zimagl = cimagl(zcsinhl);
    if (zimag < 0) 
        printf("The hyperbolic sine of (%.4lf + %.4lfi) is (%.20lf - %.20lfi)\n", creal(z), cimag(z), creal(zcsinh), fabs(zimag));
    else 
        printf("The hyperbolic sine of (%.4lf + %.4lfi) is (%.20lf + %.20lfi)\n", creal(z), cimag(z), creal(zcsinh), zimag);       

    if (zimagf < 0) 
        printf("The hyperbolic sine of (%.4f + %.4fi) is (%.20f - %.20fi)\n", crealf(zf), cimagf(zf), crealf(zcsinhf), fabsf(zimagf));
    else 
        printf("The hyperbolic sine of (%.4f + %.4fi) is (%.20f + %.20fi)\n", crealf(zf), cimagf(zf), crealf(zcsinhf), zimagf);

    if (zimagl < 0) 
        printf("The hyperbolic sine of (%.4Lf + %.4Lfi) is (%.20Lf - %.20Lfi)", creall(zL), cimagl(zL), creall(zcsinhl), fabsl(zimagl));
    else 
        printf("The hyperbolic sine of (%.4Lf + %.4Lfi) is (%.20Lf + %.20Lfi)", creall(zL), cimagl(zL), creall(zcsinhl), zimagl);
    return 0;
}
```
## 14.3 运行结果
![](csinh.png)



# 15. ctan，ctanf，ctanl
## 15.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double complex ctan (double complex z);` | 计算复数z的正切（double complex）  |
| `float complex ctanf (float complex z);` |计算复数z的正切（float complex）  |
| `long double complex ctanl (long double complex z);` | 计算复数z的正切（long double complex） |

## 15.2 演示示例
```c
#include <stdio.h>
#include <math.h>
#include <complex.h>

int main(void)
{
    double complex z, zctan;
    z = 1.0 + 2.0 * I; // I 代指 虚数单位 i
    zctan = ctan(z); // 正切

    float complex zf, zctanf;
    zf = 1.0f + 2.0f * I;
    zctanf = ctanf(zf);

    long double complex zL, zctanl;
    zL = (long double) 1.0 + (long double) 2.0 * I;
    zctanl = ctanl(zL);

    double zimag = cimag(zctan);
    float zimagf = cimagf(zctanf);
    long double zimagl = cimagl(zctanl);
    if (zimag < 0) 
        printf("The tangent of (%.4lf + %.4lfi) is (%.20lf - %.20lfi)\n", creal(z), cimag(z), creal(zctan), fabs(zimag));
    else 
        printf("The tangent of (%.4lf + %.4lfi) is (%.20lf + %.20lfi)\n", creal(z), cimag(z), creal(zctan), zimag);       

    if (zimagf < 0) 
        printf("The tangent of (%.4f + %.4fi) is (%.20f - %.20fi)\n", crealf(zf), cimagf(zf), crealf(zctanf), fabsf(zimagf));
    else 
        printf("The tangent of (%.4f + %.4fi) is (%.20f + %.20fi)\n", crealf(zf), cimagf(zf), crealf(zctanf), zimagf);

    if (zimagl < 0) 
        printf("The tangent of (%.4Lf + %.4Lfi) is (%.20Lf - %.20Lfi)", creall(zL), cimagl(zL), creall(zctanl), fabsl(zimagl));
    else 
        printf("The tangent of (%.4Lf + %.4Lfi) is (%.20Lf + %.20Lfi)", creall(zL), cimagl(zL), creall(zctanl), zimagl);
    return 0;
}
```
## 15.3 运行结果
![](ctan.png)



# 16. ctanh，ctanhf，ctanhl
## 16.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double complex ctanh (double complex z);` | 计算复数z的双曲正切（double complex）  |
| `float complex ctanhf (float complex z);` |计算复数z的双曲正切（float complex）  |
| `long double complex ctanhl (long double complex z);` | 计算复数z的双曲正切（long double complex） |

## 16.2 演示示例
```c
#include <stdio.h>
#include <math.h>
#include <complex.h>

int main(void)
{
    double complex z, zctanh;
    z = 1.0 + 2.0 * I; // I 代指 虚数单位 i
    zctanh = ctanh(z); // 双曲正切

    float complex zf, zctanhf;
    zf = 1.0f + 2.0f * I;
    zctanhf = ctanhf(zf);

    long double complex zL, zctanhl;
    zL = (long double) 1.0 + (long double) 2.0 * I;
    zctanhl = ctanhl(zL);

    double zimag = cimag(zctanh);
    float zimagf = cimagf(zctanhf);
    long double zimagl = cimagl(zctanhl);
    if (zimag < 0) 
        printf("The inverse hyperbolic tangent of (%.4lf + %.4lfi) is (%.20lf - %.20lfi)\n", creal(z), cimag(z), creal(zctanh), fabs(zimag));
    else 
        printf("The inverse hyperbolic tangent of (%.4lf + %.4lfi) is (%.20lf + %.20lfi)\n", creal(z), cimag(z), creal(zctanh), zimag);       

    if (zimagf < 0) 
        printf("The inverse hyperbolic tangent of (%.4f + %.4fi) is (%.20f - %.20fi)\n", crealf(zf), cimagf(zf), crealf(zctanhf), fabsf(zimagf));
    else 
        printf("The inverse hyperbolic tangent of (%.4f + %.4fi) is (%.20f + %.20fi)\n", crealf(zf), cimagf(zf), crealf(zctanhf), zimagf);

    if (zimagl < 0) 
        printf("The inverse hyperbolic tangent of (%.4Lf + %.4Lfi) is (%.20Lf - %.20Lfi)", creall(zL), cimagl(zL), creall(zctanhl), fabsl(zimagl));
    else 
        printf("The inverse hyperbolic tangent of (%.4Lf + %.4Lfi) is (%.20Lf + %.20Lfi)", creall(zL), cimagl(zL), creall(zctanhl), zimagl);
    return 0;
}
```
## 16.3 运行结果
![](ctanh.png)



# 17. cexp，cexpf，cexpl
## 17.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double complex cexp (double complex z);` | 计算复数z的指数基数e（double complex）  |
| `float complex cexpf (float complex z);` |计算复数z的指数基数e（float complex）  |
| `long double complex cexpl (long double complex z);` | 计算复数z的指数基数e（long double complex） |

## 17.2 演示示例
```c
#include <stdio.h>
#include <math.h>
#include <complex.h>

int main(void)
{
    double complex z, zcexp;
    z = 1.0 + 2.0 * I; // I 代指 虚数单位 i
    zcexp = cexp(z); // 指数基数e

    float complex zf, zcexpf;
    zf = 1.0f + 2.0f * I;
    zcexpf = cexpf(zf);

    long double complex zL, zcexpl;
    zL = (long double) 1.0 + (long double) 2.0 * I;
    zcexpl = cexpl(zL);

    double zimag = cimag(zcexp);
    float zimagf = cimagf(zcexpf);
    long double zimagl = cimagl(zcexpl);
    if (zimag < 0) 
        printf("The base-e exponential of (%.4lf + %.4lfi) is (%.20lf - %.20lfi)\n", creal(z), cimag(z), creal(zcexp), fabs(zimag));
    else 
        printf("The base-e exponential of (%.4lf + %.4lfi) is (%.20lf + %.20lfi)\n", creal(z), cimag(z), creal(zcexp), zimag);       

    if (zimagf < 0) 
        printf("The base-e exponential of (%.4f + %.4fi) is (%.20f - %.20fi)\n", crealf(zf), cimagf(zf), crealf(zcexpf), fabsf(zimagf));
    else 
        printf("The base-e exponential of (%.4f + %.4fi) is (%.20f + %.20fi)\n", crealf(zf), cimagf(zf), crealf(zcexpf), zimagf);

    if (zimagl < 0) 
        printf("The base-e exponential of (%.4Lf + %.4Lfi) is (%.20Lf - %.20Lfi)", creall(zL), cimagl(zL), creall(zcexpl), fabsl(zimagl));
    else 
        printf("The base-e exponential of (%.4Lf + %.4Lfi) is (%.20Lf + %.20Lfi)", creall(zL), cimagl(zL), creall(zcexpl), zimagl);
    return 0;
}
```
## 17.3 运行结果
![](cexp.png)

# 18. clog，clogf，clogl
## 18.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double complex clog (double complex z);` | 计算复数z的自然对数（以e为底）（double complex）  |
| `float complex clogf (float complex z);` |计算复数z的自然对数（以e为底）（float complex）  |
| `long double complex clogl (long double complex z);` | 计算复数z的自然对数（以e为底）（long double complex） |

## 18.2 演示示例
```c
#include <stdio.h>
#include <math.h>
#include <complex.h>

int main(void)
{
    double complex z, zclog;
    z = 1.0 + 2.0 * I; // I 代指 虚数单位 i
    zclog = clog(z); // 自然对数（以e为底）

    float complex zf, zclogf;
    zf = 1.0f + 2.0f * I;
    zclogf = clogf(zf);

    long double complex zL, zclogl;
    zL = (long double) 1.0 + (long double) 2.0 * I;
    zclogl = clogl(zL);

    double zimag = cimag(zclog);
    float zimagf = cimagf(zclogf);
    long double zimagl = cimagl(zclogl);
    if (zimag < 0) 
        printf("The natural (base-e) logarithm of (%.4lf + %.4lfi) is (%.20lf - %.20lfi)\n", creal(z), cimag(z), creal(zclog), fabs(zimag));
    else 
        printf("The natural (base-e) logarithm of (%.4lf + %.4lfi) is (%.20lf + %.20lfi)\n", creal(z), cimag(z), creal(zclog), zimag);       

    if (zimagf < 0) 
        printf("The natural (base-e) logarithm of (%.4f + %.4fi) is (%.20f - %.20fi)\n", crealf(zf), cimagf(zf), crealf(zclogf), fabsf(zimagf));
    else 
        printf("The natural (base-e) logarithm of (%.4f + %.4fi) is (%.20f + %.20fi)\n", crealf(zf), cimagf(zf), crealf(zclogf), zimagf);

    if (zimagl < 0) 
        printf("The natural (base-e) logarithm of (%.4Lf + %.4Lfi) is (%.20Lf - %.20Lfi)", creall(zL), cimagl(zL), creall(zclogl), fabsl(zimagl));
    else 
        printf("The natural (base-e) logarithm of (%.4Lf + %.4Lfi) is (%.20Lf + %.20Lfi)", creall(zL), cimagl(zL), creall(zclogl), zimagl);
    return 0;
}
```
## 18.3 运行结果
![](clog.png)


# 19. conj，conjf，conjl
## 19.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double complex conj (double complex z);` | 计算复数z的[共轭](https://baike.baidu.com/item/%E5%85%B1%E8%BD%AD/31802?fr=aladdin)（double complex）  |
| `float complex conjf (float complex z);` |计算复数z的共轭（float complex）  |
| `long double complex conjl (long double complex z);` | 计算复数z的共轭（long double complex） |


## 19.2 演示示例
```c
#include <stdio.h>
#include <math.h>
#include <complex.h>

int main(void)
{
    double complex z, zconj;
    z = 1.0 + 2.0 * I; // I 代指 虚数单位 i
    zconj = conj(z); // 共轭

    float complex zf, zconjf;
    zf = 1.0f + 2.0f * I;
    zconjf = conjf(zf);

    long double complex zL, zconjl;
    zL = (long double) 1.0 + (long double) 2.0 * I;
    zconjl = conjl(zL);

    double zimag = cimag(zconj);
    float zimagf = cimagf(zconjf);
    long double zimagl = cimagl(zconjl);
    if (zimag < 0) 
        printf("The conjugate of (%.4lf + %.4lfi) is (%.4lf - %.4lfi)\n", creal(z), cimag(z), creal(zconj), fabs(zimag));
    else 
        printf("The conjugate of (%.4lf + %.4lfi) is (%.4lf + %.4lfi)\n", creal(z), cimag(z), creal(zconj), zimag);       

    if (zimagf < 0) 
        printf("The conjugate of (%.4f + %.4fi) is (%.4f - %.4fi)\n", crealf(zf), cimagf(zf), crealf(zconjf), fabsf(zimagf));
    else 
        printf("The conjugate of (%.4f + %.4fi) is (%.4f + %.4fi)\n", crealf(zf), cimagf(zf), crealf(zconjf), zimagf);

    if (zimagl < 0) 
        printf("The conjugate of (%.4Lf + %.4Lfi) is (%.4Lf - %.4Lfi)", creall(zL), cimagl(zL), creall(zconjl), fabsl(zimagl));
    else 
        printf("The conjugate of (%.4Lf + %.4Lfi) is (%.4Lf + %.4Lfi)", creall(zL), cimagl(zL), creall(zconjl), zimagl);
    return 0;
}
```
## 19.3 运行结果
![](conj.png)



# 20. cpow，cpowf，cpowl
## 20.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`double complex cpow (double complex x, double complex y);` |  计算x的y次方值 （double complex） |
|`float complex cpowf (float complex x, float complex y);` | 计算x的y次方值 （float complex）  |
|`long double complex cpowl (long double complex x, long double complex y);` |  计算x的y次方值 （double complex） |

## 20.2 演示示例
```c
#include <stdio.h>
#include <math.h>
#include <complex.h>

int main(void)
{
    double complex x, y, z;
    x = 1.0 + 2.0 * I; // I 代指 虚数单位 i
    y = 2.0 + 1.0 * I; 
    z = cpow(x, y); // x的y次方值

    float complex xf, yf, zf;
    xf = 1.0f + 2.0f * I;
    yf = 2.0f + 1.0f * I; 
    zf = cpowf(xf, yf);

    long double complex xL, yL, zL;
    xL = (long double) 1.0 + (long double) 2.0 * I;
    yL = (long double) 2.0 + (long double) 1.0 * I;
    zL = cpowl(xL, yL);

    double zimag = cimag(z);
    float zimagf = cimagf(zf);
    long double zimagl = cimagl(zL);
    if (zimag < 0) 
        printf("the value of (%.4lf + %.4lfi) raised to the (%.4lf + %.4lfi) power is (%.20lf - %.20lfi)\n", 
            creal(x), cimag(x), creal(y), cimag(y), creal(z), fabs(zimag));
    else 
        printf("the value of (%.4lf + %.4lfi) raised to the (%.4lf + %.4lfi) power is (%.20lf + %.20lfi)\n", 
            creal(x), cimag(x), creal(y), cimag(y), creal(z), zimag);       

    if (zimagf < 0) 
        printf("the value of (%.4f + %.4fi) raised to the (%.4f + %.4fi) power is (%.20f - %.20fi)\n", 
            crealf(xf), cimagf(xf), crealf(yf), cimagf(yf), crealf(zf), fabs(zimagf));
    else 
        printf("the value of (%.4f + %.4fi) raised to the (%.4f + %.4fi) power is (%.20f + %.20fi)\n", 
            crealf(xf), cimagf(xf), crealf(yf), cimagf(yf), crealf(zf), zimagf);
    if (zimagl < 0) 
        printf("the value of (%.4Lf + %.4Lfi) raised to the (%.4Lf + %.4Lfi) power is (%.20Lf - %.20Lfi)\n", 
            creall(xL), cimagl(xL), creall(yL), cimagl(yL), creall(zL), fabs(zimagl));
    else 
        printf("the value of (%.4Lf + %.4Lfi) raised to the (%.4Lf + %.4Lfi) power is (%.20Lf + %.20Lfi)\n", 
            creall(xL), cimagl(xL), creall(yL), cimagl(yL), creall(zL), zimagl);    
    return 0;
}
```
## 20.3 运行结果
![](cpow.png)



# 21. cproj，cprojf，cprojl
## 21.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double complex cproj (double complex z);` | 计算复数z在黎曼球面上的投影（double complex）  |
| `float complex cprojf (float complex z);` |计算复数z在黎曼球面上的投影（float complex）  |
| `long double complex cprojl (long double complex z);` | 计算复数z在黎曼球面上的投影（long double complex） |

黎曼球面上的投影是一种将三维空间中的黎曼球面与二维复平面通过立体投影方式建立一一对应关系的映射。

## 21.2 演示示例
```c
#include <stdio.h>
#include <math.h>
#include <complex.h>

int main(void)
{
    double complex z, zcproj;
    z = 1.0 + 2.0 * I; // I 代指 虚数单位 i
    zcproj = cproj(z); // 计算复数z在黎曼球面上的投影

    float complex zf, zcprojf;
    zf = NAN + INFINITY * I;
    zcprojf = cprojf(zf);

    long double complex zL, zcprojl;
    zL = INFINITY + (long double) 3.0 * I; 
    zcprojl = cprojl(zL); // 结果相当于  INFINITY + i*copysign(0.0, cimag(z)).

    double zimag = cimag(zcproj);
    float zimagf = cimagf(zcprojf);
    long double zimagl = cimagl(zcprojl);
    if (zimag < 0) 
        printf("The projection of the (%.4lf + %.4lf i) onto the Riemann sphere is (%.4lf - %.4lf i)\n", creal(z), cimag(z), creal(zcproj), fabs(zimag));
    else 
        printf("The projection of the (%.4lf + %.4lf i) onto the Riemann sphere is (%.4lf + %.4lf i)\n", creal(z), cimag(z), creal(zcproj), zimag);       

    if (zimagf < 0) 
        printf("The projection of the (%.4f + %.4f i) onto the Riemann sphere is (%.4f - %.4f i)\n", crealf(zf), cimagf(zf), crealf(zcprojf), fabsf(zimagf));
    else 
        printf("The projection of the (%.4f + %.4f i) onto the Riemann sphere is (%.4f + %.4f i)\n", crealf(zf), cimagf(zf), crealf(zcprojf), zimagf);

    if (zimagl < 0) 
        printf("The projection of the (%.4Lf + %.4Lf i) onto the Riemann sphere is (%.4Lf - %.4Lf i)", creall(zL), cimagl(zL), creall(zcprojl), fabsl(zimagl));
    else 
        printf("The projection of the (%.4Lf + %.4Lf i) onto the Riemann sphere is (%.4Lf + %.4Lf i)", creall(zL), cimagl(zL), creall(zcprojl), zimagl);
    return 0;
}
```
## 21.3 运行结果
![](cproj.png)

# 22. csqrt，csqrtf，csqrtl
## 22.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `double complex csqrt (double complex z);` | 计算复数z的平方根（double complex）  |
| `float complex csqrtf (float complex z);` |计算复数z的平方根（float complex）  |
| `long double complex csqrtl (long double complex z);` | 计算复数z的平方根（long double complex） |

## 22.2 演示示例
```c
#include <stdio.h>
#include <math.h>
#include <complex.h>

int main(void)
{
    double complex z, zcsqrt;
    z = 9.0 + 9.0 * I; // I 代指 虚数单位 i
    zcsqrt = csqrt(z); // 平方根

    float complex zf, zcsqrtf;
    zf = 4.0f + 4.0f * I;
    zcsqrtf = csqrtf(zf);

    long double complex zL, zcsqrtl;
    zL = (long double) 16.0 + (long double) 16.0 * I;
    zcsqrtl = csqrtl(zL);

    double zimag = cimag(zcsqrt);
    float zimagf = cimagf(zcsqrtf);
    long double zimagl = cimagl(zcsqrtl);
    if (zimag < 0) 
        printf("The square root of (%.4lf + %.4lfi) is (%.20lf - %.20lfi)\n", creal(z), cimag(z), creal(zcsqrt), fabs(zimag));
    else 
        printf("The square root of (%.4lf + %.4lfi) is (%.20lf + %.20lfi)\n", creal(z), cimag(z), creal(zcsqrt), zimag);       

    if (zimagf < 0) 
        printf("The square root of (%.4f + %.4fi) is (%.20f - %.20fi)\n", crealf(zf), cimagf(zf), crealf(zcsqrtf), fabsf(zimagf));
    else 
        printf("The square root of (%.4f + %.4fi) is (%.20f + %.20fi)\n", crealf(zf), cimagf(zf), crealf(zcsqrtf), zimagf);

    if (zimagl < 0) 
        printf("The square root of (%.4Lf + %.4Lfi) is (%.20Lf - %.20Lfi)", creall(zL), cimagl(zL), creall(zcsqrtl), fabsl(zimagl));
    else 
        printf("The square root of (%.4Lf + %.4Lfi) is (%.20Lf + %.20Lfi)", creall(zL), cimagl(zL), creall(zcsqrtl), zimagl);
    return 0;
}
```
## 22.3 运行结果

![](csqrt.png)


# 参考
1. [【MATH-标准C库】](https://device.harmonyos.com/cn/docs/develop/apiref/math-0000001055228010#ZH-CN_TOPIC_0000001055228010__section733248854112508)



