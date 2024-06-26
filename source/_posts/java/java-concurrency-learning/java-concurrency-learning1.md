---
title: Java并发编程学习1-并发简介
date: 2021-01-29 18:12:24
updated: 2023-09-14 12:35:52
categories:
  - [开发语言-Java,Java并发编程]
tags:
  - Java
  - 并发编程
  - 并发简介
  - 线程的优势
  - 线程的风险
---



![](/images/java-concurrency-logo.png)

# 一、简介
在早期的计算机中不包含操作系统，它们从头到尾只执行一个程序，并且这个程序能访问计算机中的所有资源。在这种裸机环境中，不仅很难编写和运行程序，而且每次只能运行一个程序，这对于昂贵并稀有的计算机资源来说也是一种浪费。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

操作系统的出现使得计算机每次能运行多个程序，并且不同的程序都在单独的进程中运行：操作系统为各个独立执行的进程分配各种资源，包括内存，文件句柄以及安全证书等。如果需要的话，在不同的进程之间可以通过一些粗粒度的通信机制来交换数据，包括：套接字、信号处理器、共享内存、信号量以及文件等。

上世纪60年代，在操作系统中一直都是以**进程**作为能独立运行的基本单位。然而随着计算机技术的发展，进程慢慢出现了很多弊端，一是由于进程是资源拥有者，创建、撤消与切换存在较大的时空开销；二是由于对称多处理机（SMP）出现，可以满足多个运行单位，而多个进程并行开销过大。为了最大程度提高系统运行的性能，人们在80年代又提出了比进程更小的能独立运行的基本单位 — **线程**。

如果说，在操作系统中引入进程的目的，是为了使多个程序并发执行，以改善资源利用率及提高系统的吞吐量；那么，在操作系统中再引入线程则是为了减少程序并发执行时所付出的时空开销，使操作系统具有更好的并发性。

**线程**是Java语言中不可或缺的重要功能，它们能使复杂的异步代码变得更简单，从而极大地简化复杂系统的开发；线程会共享进程范围内的资源，例如内存句柄和文件句柄，但每个线程都有各自的程序计数器、栈以及局部变量等。线程还提供了一种直观的分解模式来充分利用多处理器系统中的硬件并行性，而在同一个程序中的多个线程也可以被同时调度到多个CPU上运行。


# 二、线程的优势
如果使用得当，线程可以有效地降低程序的开发和维护等成本，同时提升复杂应用程序的性能。线程能够将大部分的异步工作流转换成串行工作流，因此能更好地模拟人类的工作方式和交互方式。此外，线程还可以减低代码的复杂度，使代码更容易编写、阅读和维护。

## 2.1 发挥多处理器的强大能力
多线程程序可以同时在多个处理器上执行。如果设计正确，多线程程序可以通过提高处理器资源的利用率来提高系统吞吐率。在多线程程序中，如果一个线程在等待I/O操作完成，另一个线程可以继续运行，使程序能够在I/O阻塞期间继续运行。

## 2.2 建模的简单性
通过使用线程，可以将复杂并且异步的工作流进一步分解为一组简单并且同步的工作流，每个工作流在一个单独的线程中运行，并在特定的同步位置进行交互。
 
## 2.3 异步事件的简化处理
服务器应用程序在接受来自多个远程客户端的套接字连接请求时，如果为每个连接都分配其各自的线程并且使用同步I/O，那么就会降低这类程序的开发难度。如果每个请求都拥有自己的处理线程，那么在处理某个请求时发生的阻塞将不会影响其他请求的处理。

## 2.4 响应更灵敏的用户界面
传统的GUI应用程序通常都是单线程的，在代码的各个位置都需要调用poll方法来获得输入事件或者通过一个“主事件循环（Main Event Loop）”来间接地执行应用程序的所有代码。如果在主事件循环中调用的代码需要很长时间才能执行完成，那么用户界面就会“冻结”，只有当执行控制权返回到主事件循环后，才能处理后续的用户界面事件。如果将长时间运行的任务放在一个单独的线程中运行，事件线程就能及时地处理界面事件，从而使用户界面具有更高的灵敏度。在现代的GUI框架中，例如 AWT 和 Swing 等工具，都采用一个事件分发线程（Event Dispatch Thread，EDT）来替代主事件循环。

# 三、线程的风险
## 3.1 安全性问题
线程安全性可能是非常复杂的，在没有充足同步的情况下，多个线程中的操作执行顺序是不可预测的，甚至会产生奇怪的结果。由于多个线程要共享相同的内存地址空间，并且是并发执行，因此它们可能会访问或修改其他线程正在使用的变量。当多个线程同时访问和修改相同的变量时，将会在串行编程模型中引入非串行因素，而这种非串行性是很难分析的。要使多线程程序的行为可以预测，必须对共享变量的访问操作进行协同，这样才不会在线程之间发生彼此干扰。

## 3.2 活跃性问题
安全性的含义是“永远不发生糟糕的事情”，而活跃性则关注于“某件正确的事情最终会发生”。当某个操作无法继续执行下去时，就会发生活跃性问题。后续的笔记中会慢慢介绍各种形式的活跃性问题，以及如何避免这些问题，包括死锁，饥饿，以及活锁。

## 3.3 性能问题
与活跃性问题密切相关的是性能问题。性能问题包括多个方面，例如服务时间过长，响应不灵敏，吞吐率过低，资源消耗过高，或者可伸缩性较低等。

# 四、结语

Java并发编程的学习注定是个枯燥的过程，为了结合实战学习并发编程，笔者推荐目前正在学习的这本《Java并发编程实战》。笔者整理这一系列的初衷是打算能够通过写博的方式，巩固当前所学的并发编程知识，如果在这个过程中能够帮助到正在学习并发编程的小伙伴，那也算是一件值得开心的事情。相信有交流、有总结的学习过程，就不会那么的枯燥无聊了。那么下一篇我们开始了解线程安全性的相关基础知识。