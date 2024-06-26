---
title: flea-msg使用之JMS初识
date: 2023-02-11 09:00:00
updated: 2023-06-15 16:39:59
categories:
  - [开发框架-Flea,flea-msg]
tags:
  - flea-framework
  - flea-msg
  - JMS
  - 点对点模型
  - 发布/订阅模型
  - 请求/应答模式
---

# 1. JMS 基本概念
## 1.1 什么是 JMS ？
**Java** 消息服务【**Java Message Service**】，又简称 **JMS**，它是 **Java** 平台上有关面向消息中间件(**MOM**)的技术规范。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

## 1.2 JMS 规范
**JMS** 中定义了 **Java** 中访问消息中间件的接口，并没有给予实现，实现 **JMS** 接口的消息中间件称为 **JMS Provider**，例如 **ActiveMQ**、**RocketMQ** 等等。

下面我们来了解一下 **JMS** 中的一些专业名词：

|英文名称| 中文名称|描述  |
|--|--|:--|
| **JMS Provider** | **JMS** 提供者| 实现 **JMS** 接口和规范的消息中间件。提供者可以是 **Java** 平台的 **JMS** 实现，也可以是非 **Java** 平台的面向消息中间件的适配器。 |
|**JMS Message**| **JMS** 消息| 可以在 **JMS** 客户之间传递数据的对象。 它由如下三部分组成：<br/> （1） **消息头：** 每个消息头字段都有相应的 **getter** 和 **setter** 方法。<br/>（2） **消息属性：** 如果需要除消息头字段以外的值，那么可以使用消息属性。<br/> （3） **消息体：** 封装具体的消息数据。 |
| **JMS Client**| **JMS** 客户端|生产或消费基于消息的 **Java** 的应用程序或对象。 |
| **JMS Producer** | **JMS** 生产者|创建和发送 **JMS** 消息的客户端应用。 |
|**JMS Consumer** | **JMS** 消费者 |接收和处理 **JMS** 消息的客户端应用。可以有如下两种方法消费消息：<br/>（1）**同步消费：** 通过调用消费者的 **receive** 方法从目的地中显式提取消息，**receive** 方法可以一直阻塞到消息到达。<br/>（2）**异步消费：** 客户端可以为消费者注册一个消息监听器，当会话线程调用消息监听器对象的 onMessage() 方法时，客户端消费消息。|
| **JMS Domains**| **JMS** 消息传递域 | **JMS** 规范中定义了两种消息传递域：**点对点** （**point-topoint**，简称 **PTP**） 和 **发布/订阅** （**publish/subscribe**，简称 **pub/sub**）<br/><br/> **点对点消息传递域的特点**：<br/> （1）每个消息只能有一个消费者；<br>（2）消息的生产者和消费者之间没有时间上的关联性，无论消费者在生产者发送消息的时候是否处于运行状态，它都可以提取消息。生产者不需要在接收者消费该消息期间处于运行状态，接收者也同样不需要在消息发送时处于运行状态。<br/><br/> **发布/订阅消息传递域的特点**：<br/>（1）每个消息可以有多个消费者；<br/>（2）生产者和消费者之间有时间上的关联性。订阅一个主题的消费者只能消费自它订阅之后发布的消息。**JMS** 规范允许客户创建持久订阅，这在一定程度上放松了时间上的关联性要求。持久订阅允许消费者消费它在未处于激活状态时发送的消息。|
| **ConnectionFactory**| 连接工厂| 用来创建连接对象，以连接到 **JMS** 提供者。|
| **JMS Connection** | **JMS** 连接 | 封装了 **JMS** 客户端和 **JMS** 提供者【服务器端】 之间的一个活动的连接，是由客户端通过调用连接工厂的方法建立的。 |
| **JMS Session** |  **JMS** 会话| **JMS** 客户端 与 **JMS** 提供者【服务器端】 之间的会话状态。**JMS** 会话建立在 **JMS** 连接上，表示 客户端与 服务器端 之间的一个会话线程。它提供了一个事务性的上下文，在这个上下文中，一组发送和接收被组合到了一个原子操作中。|
| **JMS Destination**  | **JMS** 目的地 | 消息发送到的目的地，是实际的消息源。<br/> 在点对点消息传递域中，目的地被称为队列（**Queue**）；<br/> 在发布/订阅消息传递域中，目的地被称为主题（**Topic**）。<br/>![](to_2Domains.png) |
|**Administered Objects**| 管理对象| **JMS** 没有完全定义的两个消息传递元素是 **连接工厂** 和 **目的地**。尽管这些是 **JMS** 编程模型中的基本元素，但在提供者定义和管理这些对象的方式上存在许多现有的和预期的差异，因此创建一个通用的定义既不可能也不可取。因此，这两个对象通常是使用管理工具创建和配置的，而不是以编程方式创建的。然后将它们存储在（提供者）的对象存储区中，并由 **JMS** 客户端通过标准 **JNDI** 查找进行访问。<br/><br/> **连接工厂** 管理对象用于生成客户端到 **Broker** 的连接。它们封装了特定于提供者的信息，这些信息控制消息传递行为的某些方面：连接处理、客户端标识、消息头覆盖、可靠性和流控制等。从给定连接工厂派生的每个连接都显示为该工厂配置的行为。<br/><br/>**目的地** 管理对象用于引用 **Broker**上的物理目的地。它们封装了特定于提供者的命名（地址语法）约定，并指定了使用目的地的消息传递域：**队列（Queue）** 或 **主题（Topic）**。<br/><br/> 如下图显示了消息生产者和消息消费者如何使用目的地管理对象访问其对应的物理目的地。标记的步骤表示管理员和客户端应用程序使用此机制发送和接收消息所需采取的操作。<br/>![](to_JMSAppElements.png) <br/>**步骤1.** 管理员在 **Broker**上创建物理目的地。 <br/> **步骤2.**  管理员创建目的地管理对象，并通过指定其对应的物理目的地的名称及其类型（队列或主题）进行配置。<br/>**步骤3.** 消息生产者使用 **JNDI** 查找目的地管理对象。 <br/>**步骤4.** 消息生产者向目的地发送消息。 <br/>**步骤5.** 消息消费者查找其希望获取消息的目的地管理对象。<br/>**步骤6.** 消息消费者从目的地获取消息。|

# 2. JMS 编程对象
所谓的 **JMS** 编程对象，即为实现 **JMS** 消息传递的对象，包括 连接工厂、连接、会话、生产者、消费者、消息 和 目的地。

下图就展示了上面这些 **JMS** 编程对象之间的联系：

![](to_JMSObjects.png)

我们在上图中可以看到，有两个对象（**连接工厂** 和 **目的地**）是在对象存储中的。它们通常是作为管理对象创建、配置和管理的。也就是说，**连接工厂** 和 **目的地** 是以管理方式（而不是以编程方式）创建的。

如下表格总结了 发送 和 接收 消息所需的步骤【从上图中也可看出一二】：
<table>
	<tr>
		<th>生产消息</th>
	    <th>消费消息</th>
	</tr >
	<tr>
		<td colspan="2" align="left" > 1. 管理员创建连接工厂管理的对象。 </td>
	</tr>
	<tr>
		<td colspan="2" align="left" > 2. 管理员创建物理目的地和引用它的管理对象。</td>
	</tr>
	<tr>
		<td colspan="2" align="left" > 3. 客户端通过 <b>JNDI</b> 查找获得连接工厂对象。 </td>
	</tr>
	<tr>
		<td colspan="2" align="left" > 4. 客户端通过 <b>JNDI</b> 查找获得目的地对象。</td>
	</tr>
	<tr>
		<td colspan="2" align="left" > 5. 客户端创建一个连接并设置针对此连接的属性。</td>
	</tr>
	<tr>
		<td colspan="2" align="left" > 6. 客户端创建一个会话并设置管理消息传递可靠性的属性。</td>
	</tr>
	<tr>
		<td align="left" > 7. 客户端创建消息生产者 </td>
		<td align="left" > 客户端创建消息消费者</td>
	</tr>
	<tr>
		<td align="left" > 8. 客户端创建消息 </td>
		<td align="left" > 客户端启动连接</td>
	</tr>
	<tr>
		<td align="left" > 9. 客户端发送消息 </td>
		<td align="left" > 客户端接收消息</td>
	</tr>
</table>

> 注意，**步骤1到6** 对于发送方和接收方是相同的。

下面我们来详细介绍下 **JMS** 编程对象：

## 2.1 连接工厂和连接
客户端使用连接工厂对象（**ConnectionFactory**）创建连接。连接对象（**Connection**）表示客户端到 **Broker** 的活跃连接。它使用基础连接服务，该服务在默认情况下启动，或者由该客户端的管理员显式启动。

通信资源的分配和客户端的身份验证都在创建连接时进行。它是一个相对重量级的对象，大多数客户端都使用一个连接来完成所有消息传递。连接支持并发使用：***任何数量的生产者和消费者都可以共享连接***。

创建连接工厂时，可以通过设置其属性来配置从该工厂派生的所有连接的行为。对于消息队列，它们可以指定如下信息：

- **Broker** 驻留的主机的名称、所需的连接服务以及客户端希望访问该服务的端口。
- 如果连接失败，应如何处理与  **Broker** 的自动重新连接。（如果连接丢失，此功能会将客户端重新连接到同一个（或不同的 **Broker**）。无法保证数据故障切换：当重新连接到其他代理时，持久消息和其他状态信息可能会丢失。）
- 需要 **Broker** 跟踪其持久订阅的客户端的ID。
- 尝试连接的用户的默认名称和密码。如果在连接时未指定密码，则此信息用于验证用户并授权操作。
- 对于那些不关心可靠性的客户端，是否应禁止 **Broker** 签收。
- 如何管理 **Broker** 和客户端运行时之间的控制流和有效负载消息。
- 应如何处理队列浏览（仅限Java客户端）。
- 是否应重写某些消息头字段。

可以从启动客户端应用程序的命令行来覆盖连接工厂属性。也可以通过设置那连接的属性来覆盖任何给定连接的属性。

您可以使用连接对象来创建会话对象、设置异常监听器 或 获取 **JMS** 版本和提供者信息。

## 2.2 会话
如果 **Connection** 表示客户端和 **Broker** 之间的通信通道，那 **Session** 就将代表它们之间的单个会话。后面我们主要使用会话对象来创建消息、消息生产者和消息消费者。创建会话时，您可以通过多个确认选项 或者 事务 来配置可靠的传递。有关详细信息，请参阅 [可靠性消息传递](https://docs.oracle.com/cd/E19717-01/819-7759/aerbz/index.html)。

根据 **JMS** 规范，会话是用于生产和消费消息的单线程上下文。您可以为一个会话创建多个消息生产者和消费者，但您只能连续使用它们。**Java** 和 **C** 客户端的线程实现略有不同，

还可以使用会话对象执行以下操作：

- 为那些不使用管理对象定义目的地的客户端创建和配置目的地。
- 创建和配置临时主题和队列；这些被用作请求-应答模式的一部分。请参阅 [请求-应答模式](https://docs.oracle.com/cd/E19717-01/819-7759/aerby/index.html)。
- 支持事务处理。
- 定义生产或消费消息的序列顺序。
- 为异步消费者序列化消息监听器的执行。
- 创建队列浏览器（仅限Java客户端）。

## 2.3 消息
上面我们了解到，消息由三部分组成，分别是 **消息头**、**消息属性** 和 **消息体**。

### 2.3.1 消息头
每个 **JMS** 消息都需要一个消息头。消息头包含十个预定义字段，这些字段参考如下表格：
|消息头字段| 描述  |
|:--|:--|
|JMSDestination  | 指定将消息发送到的目的地对象的名称（由提供者设置），也就是 **Queue** 和 **Topic**，自动分配。 |
| JMSDeliveryMode | 传送模式，指定消息是否持久（默认情况下，由 提供者 或 客户端 为生产者或单个消息显式设置）。有两种 ：持久模式和非持久模式。 |
| JMSExpiration | 指定消息过期的时间（默认情况下，由提供者 或 客户端为生产者或单个消息设置），它等于 **Destination** 的 **send** 方法中的 **timeToLive** 值加上发送时刻的 **GMT** 时间值。如果 **timeToLive** 值等于零，则**JMSExpiration** 被设为零，表示该消息永不过期。如果发送后，在消息过期时间之后消息还没有被发送到目的地，则该消息被清除。 |
| JMSPriority  | 指定0（低）到9（高）范围内的消息优先级（默认情况下，由提供者设置 或 客户端为生产者或单个消息显式设置），其中 **0-4** 是 **普通消息**，**5-9** 是 **加急消息**。**JMS** 不要求 **JMS Provider** 严格按照这十个优先级发送消息，但 ***必须保证加急消息要先于普通消息到达***。 |
| JMSMessageID |  为提供者上下文中的消息指定唯一ID（由提供者设置）|
|  JMSTimestamp| 指定提供者接收消息的时间（由提供者设置） |
|  JMSCorrelationID|  允许客户端定义两个消息之间的对应关系的值（如果需要，由客户端设置），典型的应用是在应答消息中连接到原消息。|
| JMSReplyTo | 指定消费者应发送回复的目的地（如果需要，由客户端设置） |
| JMSType |  消息类型的识别符，可以由消息选择器计算的值（如果需要，由客户端设置）|
| JMSRedelivered  | 指定消息是否已传递但未确认（由提供者设置）。如果一个客户端收到一个设置了 **JMSRedelivered** 属性的消息，则表示可能客户端曾经在早些时候收到过该消息，但并没有签收(acknowledged)。 |

通过查看上述表格，我们可以看出，消息头字段有多种用途：**标识消息**、**配置消息路由**、**提供有关消息处理的信息**等等。消息生产者可能需要配置消息头以获得某些消息传递行为；消息消费者可能需要读取消息头，以了解消息是如何路由的，以及它可能需要的进一步的处理。

**JMSDeliveryMode** 是最重要的字段之一，它决定了消息传递的可靠性。此字段指示消息是否持久。

- **持久消息**。保证消息传递并成功消费一次。如果消息服务失败，持久消息不会丢失。

- **非持久性消息**。保证消息最多传递一次。如果消息服务失败，非持久性消息可能会丢失。

### 2.3.2 消息属性
**JMS** 规范中包含如下三种类型的属性：

- 应用程序设置或添加的属性
- **JMS** 定义的属性。
- **JMS** 供应商特定的属性。

**JMS** 规范定义了九个标准属性，详见如下表格。其中一些由客户端设置，一些由提供者设置。它们的名称以保留字符 **“JMSX”** 开头。客户端或提供者可以使用这些属性来确定谁发送了消息、消息的状态、发送频率和时间。这些属性对于提供者路由消息和提供诊断信息很有用。

|属性名| 描述  |
|:--|:--|
| JMSXUserID | 发送消息的用户标识，发送时由提供者设置  |
| JMSXAppID| 发送消息的应用标识，发送时由提供者设置 |
| JMSXDeliveryCount|  转发消息重试次数,第一次是1，第二次是2，… ，发送时由提供者设置|
| JMSXGroupID| 消息所在消息组的标识，由客户端设置 |
| JMSXGroupSeq| 组内消息的序号第一个消息是1，第二个是2，…，由客户端设置 |
| JMSXProducerTXID| 产生消息的事务的事务标识，发送时由提供者设置 |
| JMSXConsumerTXID| 消费消息的事务的事务标识，接收时由提供者设置 |
| JMSXRcvTimestamp| JMS 转发消息到消费者的时间，接收时由提供者设置 |
| JMSXState| 假定存在一个消息仓库，它存储了每个消息的单独拷贝，且这些消息从原始消息被发送时开始。每个拷贝的状态有：1（等待），2（准备），3（到期）或4（保留）。由于状态与生产者和消费者无关，所以它不是由生产者和消费者来提供。它只和在仓库中查找消息相关，因此JMS没有提供这种API。由提供者设置  |

消息队列也定义了消息属性，这些属性用于标识压缩消息以及在无法传递消息时应如何处理消息。

### 2.3.1 消息体
消息体包含客户端要交换的数据。

**JMS** 消息的类型决定了消息体可能包含的内容以及消费者应该如何处理它，详见如下表格。另外，**Session** 对象中包含了每种类型的消息体的创建方法。

|消息类型|   描述|
|:--|:--|
| StreamMessage | 消息体包含 **Java** 原始值流的消息。它是按顺序填充和读取的。 |
| MapMessage | 消息体包含一组键值对的消息。未定义条目的顺序。 |
| TextMessage | 消息体包含Java字符串的消息，例如XML字符串消息。 |
| ObjectMessage |  消息体包含序列化Java对象的消息。|
| BytesMessage | 消息体包含未解释字节流的消息 |
| Message| 包含消息头和消息属性，但不包含消息体的消息 |

**Java** 客户端可以设置一个属性，让客户端运行时压缩正在发送的消息的消息体。消费者端的消息队列运行时在传递消息之前对消息进行解压缩。

## 2.4 生产者
上文中，我们知道生产者是创建和发送 **JMS** 消息的客户端应用，消息就是由消息生产者在连接和会话的上下文中发送或发布。生成消息其实非常简单：客户端使用消息生成器对象（**MessageProducer**）将消息发送到物理目的地（在 **JMS API** 中由目的地对象表示）。

创建生产者时，可以指定所有生产者发送消息的默认目的地。还可以为消息头字段指定默认值，这些字段控制持久性、优先级和生存时间。然后，从该生产者发出的所有消息都会使用这些默认值，除非在发送消息时通过指定备用目的地 或 为给定消息的消息头字段设置备用值 来覆盖这些默认值。

消息生产者还可以通过设置 **JMSReplyTo** 消息头字段来实现请求-应答模式。有关更多信息，请参阅 [请求-应答模式](https://docs.oracle.com/cd/E19717-01/819-7759/aerby/index.html)。

## 2.5 消费者
消费者是接收和处理 **JMS** 消息的客户端应用，消息就是由消息消费者在连接和会话的上下文中接收和处理的。客户端使用消息消费者对象（**MessageConsumer**）从指定的物理目的地（在 **JMS API** 中表示为目的地对象）接收消息。

需要注意，有如下三个因素影响 **Broker** 向消费者传递消息的方式：

- 消费是同步还是异步
- 是否使用选择器筛选传入消息
- 如果消息是从主题目标消费的，则订阅者是否持久

影响消息传递和客户端设计的另一个主要因素是消费者所需的可靠性程度。请参阅 [可靠性消息](https://docs.oracle.com/cd/E19717-01/819-7759/aerbz/index.html)。

### 2.5.1 同步和异步消费者
消息消费者可以支持消息的同步或异步消费。

- **同步消费**。它意味着消费者需明确请求传递消息，然后消费它。根据请求消息的方式，同步消费者可以选择（无限期地）等待消息到达，等待指定的消息时间，或者在没有消息可供使用时立即返回。（**"Consumed"** 表示客户端可以立即使用该对象。已成功发送但 **Broker** 尚未完成处理的消息【即尚未准备好消费】。）

- **异步消费**。它意味着消息将自动传递到已为消费者注册的消息监听器对象（**MessageListener**）上。当会话线程调用消息监听器对象的 `onMessage()` 方法时，客户端消费消息。

### 2.5.2 消息选择器
消息消费者可以使用消息选择器让消息服务仅传递其属性与特定选择条件匹配的消息。我们在创建消费者时可以指定此条件。

选择器使用类似 **SQL** 的语法来匹配消息属性。例如：

```sql
name = "Huazie"
age >= 18
```

**Java** 客户端还可以在浏览队列时指定选择器；这允许您查看 **有哪些选定的消息正在等待使用**。

### 2.5.3 持久订阅者
我们可以使用会话对象创建主题的持久订阅者。即使订阅者处于非活跃状态，**Broker** 也会保留这些订阅者的消息。

因为 **Broker** 必须维护订阅者的状态，并在订阅者被重新激活时恢复消息的传递，所以 **Broker** 必须能够在其来来往往的过程中识别给定订阅者。订阅者的标识是根据创建它的连接的 **ClientID** 属性和创建订阅者时指定的订阅者名称构造的。

# 3. JMS 点对点 模型
在 **点对点** 模型中，消息生产者称为发送者，消息消费者称为接收者。它们通过一个称为 **队列（Queue）** 的目的地交换消息：***发送方向队列生产消息，接收者消费队列中的消息***。

下图展示了 **点对点** 中一个最简单的消息传递操作。**MyQueueSender** 将 **Msg1** 发送到队列目的地 **MyQueue1**。然后，**MyQueueReceiver** 从**MyQueue1** 中获取消息。

![](to_simpleQ.png)
至于更为复杂的场景，我们可以看下图。两个发送方 **MyQSender1** 和 **MyQSender2** 使用 **相同的连接** 向 **MyQueue1** 发送消息。**MyQSender3** 使用额外的连接向**MyQueue1** 发送消息。在接收端，**MyQReceiver1** 使用来自 **MyQueue1** 的消息，**MyQRreceiver2** 和 **MyQRreceive3** 共享一个连接以使用来自 **MyQueue1** 的信息。

![](to_complexQ.png)
下面我们来总结一下，上图的场景中展示的 **点对点** 消息传递的一些附加要点：

- 多个生产者可以向队列发送消息。生产者可以共享一个连接或使用不同的连接，但他们都可以访问同一个队列。
- 多个接收方可以使用队列中的消息，但每个消息只能由一个接收方消费。因此，**Msg1**、**Msg2** 和 **Msg3** 由不同的接收器使用。
- 接收方可以共享一个连接或使用不同的连接，但它们都可以访问同一个队列。
- 发送方和接收方没有时间依赖性：无论客户端发送消息时消息是否正在运行，接收方都可以获取消息。
- 发送方和接收方可以在运行时动态添加和删除，从而允许消息传递系统根据需要进行扩展或收缩。
- 消息按照发送的顺序放置在队列中，但它们的消费顺序取决于消息过期日期、消息优先级以及是否使用选择器来使用消息等因素。

综合来说，**点对点** 模型具有如下的一些优势：

- 如果消息的接收顺序不重要，那么多个接收者可以消费同一队列中的消息，这一事实允许您平衡消息消耗。
- 即使没有接收方，也始终保留发往队列的消息。
- **Java** 客户端可以使用队列浏览器对象来检查队列的内容。然后，他们可以根据从检查中获得的信息消费消息。也就是说，尽管消费模型通常是FIFO（先进先出），但如果消费者通过使用消息选择器知道他们想要什么消息，他们可以消费不在队列头部的消息。管理客户端还可以使用队列浏览器监视队列的内容。

# 4. JMS 发布/订阅 模型
在 **发布/订阅** 模型中，消息生产者称为发布者，消息消费者称为订阅者。他们通过一个称为 **主题（Topic）** 的目的地交换消息：***发布者向主题发布消息；订阅者订阅主题并消费来自主题的消息***。

下图展示了发布/订阅域中的一个最简单的消息传递操作。**MyTopicPublisher** 将 **Msg1** 发布到 **MyTopic**。然后，**MyTopicSubscriber1** 和 **MyTopicSubscriber2**分别从 **MyTopic** 接收 **Msg1** 的副本。
![](to_simpleTopic.png)
虽然 **发布/订阅** 模型不需要有多个订阅者，但图中列出了两个订阅者，这就告诉我们该模型允许广播消息。主题的所有订阅者都会获得发布到该主题的任何消息的副本。

订阅服务器可以是持久的或者非持久的。**Broker** 将保留所有活跃订阅者的消息，但仅当这些订阅者是持久的，**Broker** 才会保留非活跃订阅者的信息。

下面我们来看下更为复杂的场景，如下图所示。三个生产者向 **Topic1** 发布消息，三个消费者消费来自 **Topic1** 的消息；除非订阅者使用选择器来筛选消息，否则每个订阅者都会获得发布到所选主题的所有消息【其中，**MyTSubscriber2** 过滤掉了 **Msg2**】。

![](to_complexTopic.png)
通过上图的场景，我们来总结一下其展示的 **发布/订阅** 消息传递的一些附加要点：

- 多个生产者可以向主题发布消息。生产者可以共享一个连接或使用不同的连接，但他们都可以访问同一主题。
- 多个订阅者可以消费来自主题的消息。订阅服务器检索发布到主题的所有消息，除非它们使用选择器筛选出消息，或者消息在使用之前过期。
- 订阅服务器可以共享一个连接或使用不同的连接，但它们都可以访问同一主题。
- 持久订阅者可以是活跃的或非活跃的。**Broker** 在它们处于非活跃状态时将为它们保留消息。
- 发布者和订阅者可以在运行时动态添加和删除，从而允许消息传递系统根据需要进行扩展或收缩。
- 消息按照发送的顺序发布到主题，但使用它们的顺序取决于消息过期日期、消息优先级以及是否使用选择器来使用消息等因素。
- 发布者和订阅者具有时间依赖性：***主题订阅者只能使用在创建订阅后发布的消息***。

**发布/订阅** 模型的主要优点是它允许 **向订阅者广播消息**。

# 5. JMS 请求-应答 模式
我们可以在同一个 **连接**（甚至使用统一API的 **会话**）中组合生产者和消费者。此外，**JMS API** 允许我们通过使用 **临时目的地** 来为 **消息传递操作** 实现 **请求-应答** 模式。

如果想要设置 **请求-应答** 模式，我们需要执行以下操作：

1. 创建一个消费者可以发送应答的临时目的地。
2. 在要发送的消息中，将消息头的 **JMSReplyTo** 字段设置为该临时目的地。

当消息消费者处理消息时，它检查消息的 **JMSReplyTo** 字段以确定是否需要应答，并将应答发送到指定的目的地。

**请求-应答** 机制为生产者省去了为应答目的地设置管理对象的麻烦，并使消费者更容易响应请求。当生产者在继续之前必须确保已经处理了请求时，该模式将非常有用。

下图就展示了 向主题发送消息并在临时队列中接收应答的 **请求-应答** 模式
![](to_ReplyTo.png)
如上图所示，**MyTopicPublisher** 向目标 **MyTopic** 生产了 **Msg1**。**MyTopicSubscriber1** 和 **MyTopicSubscriber2** 接收消息并向 **MyTempQueue** 发送应答，**MyTQReceiver** 从中检索消息。此模式可能适用于向大量客户端发布定价信息并将其订单排队进行顺序处理的应用程序。

临时目的地存在的时间仅与创建它们的连接一样长。任何生产者都可以发送到临时目的地，但唯一可以访问临时目的地的消费者是由创建目的地的同一连接创建的消费者。

由于 **请求-应答** 模式依赖于创建的临时目的地，所以在以下的情况下不应该使用此模式：

- 如果你预计创建临时目的地的连接可能会在发送应答之前终止。
- 如果需要将持久消息发送到临时目的地。

# 参考资料
1. [【百度百科--JMS】](https://baike.baidu.com/item/JMS/2836691)
2. [【JMS as a MOM Standard】](https://docs.oracle.com/cd/E19717-01/819-7759/aerar/index.html)
3. [【Client Programming Model】](https://docs.oracle.com/cd/E19717-01/819-7759/6n9mco7g2/index.html)
