---
title: C语言函数大全-- l 开头的 Linux 内核函数(链表管理函数)
date: 2023-04-21 11:44:23
updated: 2025-04-29 20:56:30
categories:
  - 开发语言-C
tags:
  - C语言函数大全
  - l开头的Linux 内核函数(链表管理函数)
---

![](/images/cplus-logo.png)


# 总览
| 函数声明 |  函数功能  |
|:--|:--|
|`void list_add(struct list_head *new, struct list_head *head);` | 它是 Linux 内核中双向链表操作的一个基本函数，用于将新节点添加到双向链表的头部  |
|`void list_add_tail(struct list_head *new, struct list_head *head);`|它是 Linux 内核中双向链表操作的一个基本函数，用于将新节点添加到链表尾部。|
|`void list_cut_before(struct list_head *new, struct list_head *head, struct list_head *entry);` |  它是 Linux 内核中双向链表操作的一个基本函数，用于将一段节点从原始链表中移动到另一个链表中，并将其添加到新链表的头部。 |
|`void list_cut_position(struct list_head *new, struct list_head *head, struct list_head *entry);`|它是 Linux 内核中双向链表操作的一个基本函数，用于将一段节点从原始链表中移动到另一个链表中，并将其添加到新链表的头部。与list_cut_before不同的是，该函数需要指定要移动的节点的具体位置，而不是直接指定一个节点。|
|`void list_del(struct list_head *entry);` | 用于从链表中删除一个节点，但不会修改该节点的指针信息。  |
| `void list_del_init(struct list_head *entry);`|  用于从链表中删除一个节点，但会将被删除的节点的指针信息初始化为NULL。 |
| `void list_del_init_careful(struct list_head *entry, struct list_head *prev, struct list_head *next);`| 用于从链表中删除一个节点，但需要指定该节点的前驱节点和后继节点，以确保链表结构正确。  |
| `int list_empty(const struct list_head *head);`|  用于判断链表是否为空，并返回非零值表示为空，返回0表示不为空 |
|`int list_empty_careful(const struct list_head *head);` |  用于判断链表是否为空，但会先检查链表头部的指针是否为空，以避免对空指针进行解引用。 |
|`void list_move(struct list_head *list, struct list_head *head);`|用于将一个节点移动到另外一个链表的头部。|
|`void list_move_tail(struct list_head *list, struct list_head *head);`|用于将一个节点移动到另外一个链表的尾部。|
|`void list_bulk_move_tail(struct list_head *list, int count, struct list_head *head);`|用于将多个节点从一个链表移动到另一个链表的尾部。|
|`void list_replace(struct list_head *old, struct list_head *new);`|用于用一个新节点替换指定节点。|
|`void list_replace_init(struct list_head *old, struct list_head *new);`|除了可以完成 list_replace 做的所有操作外，它还将原来的节点初始化为空。|
|`static inline void list_rotate_left(struct list_head *head)`|用于将链表向左旋转一个位置。|
|`void list_rotate_to_front(struct list_head *head, struct list_head *pivot);`|用于将指定节点移到链表头部，并旋转链表使得该节点成为新的头部。|
|`void list_splice(struct list_head *list, struct list_head *head);`|用于将一个链表中的所有节点插入到另一个链表的指定位置之前。|
|`void list_splice_tail(struct list_head *list, struct list_head *head);`|用于将一个链表中的所有节点插入到另一个链表的尾部。|
|`void list_splice_init(struct list_head *list, struct list_head *head);`|除了可以完成 `list_splice` 做的所有操作外，它还将原来的链表初始化为空。|
|`void list_splice_tail_init(struct list_head *list, struct list_head *head);`|除了可以完成 `list_splice_tail` 做的所有操作外，它还将原来的链表初始化为空。|
| `void list_swap(struct list_head *list1, struct list_head *list2);`| 交换两个链表头部的位置。  |


# 1. list_add，list_add_tail
## 1.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void list_add(struct list_head *new, struct list_head *head);` | 它是 Linux 内核中双向链表操作的一个基本函数，用于将新节点添加到双向链表的头部  |
|`void list_add_tail(struct list_head *new, struct list_head *head);`|它是 Linux 内核中双向链表操作的一个基本函数，用于将新节点添加到链表尾部。|
**参数：**
- **new ：** 要添加的新节点的指针
- **head ：** 链表头节点的指针。 
  - `list_add()` 函数会将 new 节点插入到链表头之前，使其成为新的链表头节点。
  - `list_add_tail()` 函数会根据 链表头节点找到链表尾节点，并将 new 节点添加到链表尾部。 

## 1.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <stddef.h>

struct list_head {
    struct list_head *next, *prev;
};

struct node {
    int data;
    struct list_head link;
};

void init_list_head(struct list_head *head)
{
    head->next = head->prev = head;
}

void print_list(struct list_head *head)
{
    struct node *p;
    for (p = (struct node *)head->next; p != (struct node *)head; p = (struct node *)p->link.next) {
        printf("%d ", p->data);
    }
    printf("\n");
}

int main()
{
    struct list_head head = { NULL, NULL };
    init_list_head(&head); // 用于初始化双向链表头部节点。

    struct node *n1 = (struct node *)malloc(sizeof(struct node));
    n1->data = 1;
    list_add(&n1->link, &head); // 将节点添加到链表的头部

    struct node *n2 = (struct node *)malloc(sizeof(struct node));
    n2->data = 2;
    list_add(&n2->link, &head);

    struct node *n3 = (struct node *)malloc(sizeof(struct node));
    n3->data = 3;
    list_add_tail(&n3->link, &head); // 将节点添加到链表的尾部

    printf("The original list is: ");
    print_list(&head);

    return 0;
}
```

在上述示例代码中，我们首先定义了 `list_head` 和 `node` 两个结构体，并通过调用 `init_list_head()` 函数初始化链表头部。然后，我们创建了三个 `node` 类型的节点，前两个节点分别通过 `list_add()` 函数将它们添加到链表的头部，最后一个节点通过 `list_add_tail()` 函数添加到链表的尾部 。最后，我们调用 `print_list()` 函数打印链表中的元素。

> **注意：** 在使用 `list_add()` 和 `list_add_tail()` 函数之前，我们要为每个新节点分配内存空间。

# 2. list_cut_before，list_cut_position
## 2.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void list_cut_before(struct list_head *new, struct list_head *head, struct list_head *entry);` |  它是 Linux 内核中双向链表操作的一个基本函数，用于将一段节点从原始链表中移动到另一个链表中，并将其添加到新链表的头部。 |
|`void list_cut_position(struct list_head *new, struct list_head *head, struct list_head *entry);`|它是 Linux 内核中双向链表操作的一个基本函数，用于将一段节点从原始链表中移动到另一个链表中，并将其添加到新链表的头部。与list_cut_before不同的是，该函数需要指定要移动的节点的具体位置，而不是直接指定一个节点。|

**参数：**
- **new ：** 要添加的新链表头部；
- **head ：** 原始链表的头部
- **entry ：** 要移动的节点

`list_cut_before()` 函数会将 **entry** 节点及其前面的所有节点从原始链表中移动到 **new** 所指示的链表中，并将 **entry** 所在位置的前一个节点作为新链表的头节点。
`list_cut_position()` 函数会将 **entry** 节点及其后面的所有节点从原始链表中移动到 **new** 所指示的链表中，并将 **entry** 所在位置作为新链表的头节点。

## 2.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <stddef.h>

struct list_head {
    struct list_head *next, *prev;
};

struct node {
    int data;
    struct list_head link;
};

void init_list_head(struct list_head *head)
{
    head->next = head->prev = head;
}

void print_list(struct list_head *head)
{
    struct node *p;
    for (p = (struct node *)head->next; p != (struct node *)head; p = (struct node *)p->link.next) {
        printf("%d ", p->data);
    }
    printf("\n");
}

int main()
{
    struct list_head head1 = { NULL, NULL };
    init_list_head(&head1);

    struct node *n1 = (struct node *)malloc(sizeof(struct node));
    n1->data = 1;
    list_add_tail(&n1->link, &head1);

    struct node *n2 = (struct node *)malloc(sizeof(struct node));
    n2->data = 2;
    list_add_tail(&n2->link, &head1);

    struct node *n3 = (struct node *)malloc(sizeof(struct node));
    n3->data = 3;
    list_add_tail(&n3->link, &head1);

    printf("The original list is: ");
    print_list(&head1);

    struct list_head head2 = { NULL, NULL };
    init_list_head(&head2);

    // 移动节点n1和n2到另一个链表中
    list_cut_before(&head2, &head1, &n2->link);
    printf("The first list after move is: ");
    print_list(&head1);
    printf("The second list after move is: ");
    print_list(&head2);

    // 再次移动节点n3到另一个链表中
    list_cut_position(&head2, &head1, &n3->link);
    printf("The first list after second move is: ");
    print_list(&head1);
    printf("The second list after second move is: ");
    print_list(&head2);

    return 0;
}
```
在上述示例代码中，我们首先定义了 `list_head` 和 `node` 两个结构体，并通过调用 `init_list_head()` 函数初始化两个链表的头部。然后，我们创建了三个 `node` 类型的节点并分别将它们添加到第一个链表的尾部。接着，我们利用`list_cut_before()`函数和 `list_cut_position()` 函数将链表中的一段节点移动到第二个链表中。最后，我们调用 `print_list` 函数分别打印两个链表中的元素。

# 3. list_del，list_del_init，list_del_init_careful
## 3.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void list_del(struct list_head *entry);` | 用于从链表中删除一个节点，但不会修改该节点的指针信息。  |
| `void list_del_init(struct list_head *entry);`|  用于从链表中删除一个节点，但会将被删除的节点的指针信息初始化为NULL。 |
| `void list_del_init_careful(struct list_head *entry, struct list_head *prev, struct list_head *next);`| 用于从链表中删除一个节点，但需要指定该节点的前驱节点和后继节点，以确保链表结构正确。  |

**参数：**
- **entry ：** 要删除的节点
- **prev ：** 该节点的前驱节点
- **next ：** 该节点的后继节点

## 3.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <stddef.h>

struct list_head {
    struct list_head *next, *prev;
};

struct node {
    int data;
    struct list_head link;
};

void init_list_head(struct list_head *head)
{
    head->next = head->prev = head;
}

void print_list(struct list_head *head)
{
    struct node *p;
    for (p = (struct node *)head->next; p != (struct node *)head; p = (struct node *)p->link.next) {
        printf("%d ", p->data);
    }
    printf("\n");
}

int main()
{
    struct list_head head = { NULL, NULL };
    init_list_head(&head);

    struct node *n1 = (struct node *)malloc(sizeof(struct node));
    n1->data = 1;
    list_add_tail(&n1->link, &head);

    struct node *n2 = (struct node *)malloc(sizeof(struct node));
    n2->data = 2;
    list_add_tail(&n2->link, &head);

    struct node *n3 = (struct node *)malloc(sizeof(struct node));
    n3->data = 3;
    list_add_tail(&n3->link, &head);

    printf("The original list is: ");
    print_list(&head);

    // 删除节点n2，但不改变其指针信息
    list_del(&n2->link);
    printf("The list after delete n2 is: ");
    print_list(&head);

    // 删除节点n3，并初始化其指针信息为NULL
    list_del_init(&n3->link);
    printf("The list after delete and init n3 is: ");
    print_list(&head);

    // 删除节点n1，并指定其前驱和后继节点
    list_del_init_careful(&n1->link, &head, head.next);
    printf("The list after careful delete n1 is: ");
    print_list(&head);

    return 0;
}
```

在上述示例代码中，我们首先定义了 `list_head` 和 `node` 两个结构体，并通过调用 `init_list_head()` 函数初始化链表头部。然后，我们创建了三个 `node` 类型的节点并分别将它们添加到链表的尾部。接下来，我们利用 `list_del()`、 `list_del_init()` 和 `list_del_init_careful()` 函数从链表中删除节点，并打印每次操作后的链表元素。

> **注意：** 在使用这些函数之前，我们要确保被删除的节点在链表中。

# 4. list_empty，list_empty_careful
## 4.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `int list_empty(const struct list_head *head);`|  用于判断链表是否为空，并返回非零值表示为空，返回0表示不为空 |
|`int list_empty_careful(const struct list_head *head);` |  用于判断链表是否为空，但会先检查链表头部的指针是否为空，以避免对空指针进行解引用。 |

**参数：**
- **head ：** 要判断的链表头部

## 4.2 演示示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <stddef.h>

struct list_head {
    struct list_head *next, *prev;
};

struct node {
    int data;
    struct list_head link;
};

void init_list_head(struct list_head *head)
{
    head->next = head->prev = head;
}

void print_list(struct list_head *head)
{
    struct node *p;
    for (p = (struct node *)head->next; p != (struct node *)head; p = (struct node *)p->link.next) {
        printf("%d ", p->data);
    }
    printf("\n");
}

int main()
{
    struct list_head head = { NULL, NULL };
    init_list_head(&head);

    printf("Is the list empty? %d\n", list_empty(&head));
    printf("Is the list empty carefully? %d\n", list_empty_careful(&head));

    struct node *n1 = (struct node *)malloc(sizeof(struct node));
    n1->data = 1;
    list_add_tail(&n1->link, &head);

    printf("The list after adding n1: ");
    print_list(&head);
    printf("Is the list empty? %d\n", list_empty(&head));
    printf("Is the list empty carefully? %d\n", list_empty_careful(&head));

    struct node *n2 = (struct node *)malloc(sizeof(struct node));
    n2->data = 2;
    list_add_tail(&n2->link, &head);

    printf("The list after adding n2: ");
    print_list(&head);
    printf("Is the list empty? %d\n", list_empty(&head));
    printf("Is the list empty carefully? %d\n", list_empty_careful(&head));

    return 0;
}
```
在上述示例代码中，我们首先定义了 `list_head` 和 `node` 两个结构体，并通过调用init_list_head()函数初始化链表头部。然后，我们利用list_empty()和list_empty_careful()函数分别判断链表是否为空，并打印其返回值。接下来，我们创建了两个node类型的节点并分别将它们添加到链表的尾部。每次添加节点后，我们再次使用list_empty和list_empty_careful函数判断链表是否为空，并打印其返回值。

需要注意的是，在使用这些函数之前，我们要确保链表头部已经初始化。

# 5. Linux 内核中双向链表遍历相关宏

| 宏定义 |  宏描述  |
|:--|:--|
|`#define list_entry(ptr, type, member) ((type *)((char *)(ptr)-(unsigned long)(&((type *)0)->member)))` | 用于获取一个节点所在结构体的起始地址。  |
|`static inline int list_entry_is_head(const struct list_head *entry, const struct list_head *head) { return entry->prev == head; }`|   用于判断给定节点是否为链表头。 |
`#define list_first_entry(ptr, type, member) list_entry((ptr)->next, type, member)`|用于获取链表中第一个节点所在结构体的起始地址。|
|`#define list_first_entry_or_null(ptr, type, member) ({ struct list_head *__head = (ptr); struct list_head *__pos = __head->next; __pos != __head ? list_entry(__pos, type, member) : NULL; })`|用于获取链表中第一个节点所在结构体的起始地址，但会先检查链表是否为空，以避免对空指针进行解引用。|
|`#define list_next_entry(pos, member) list_entry((pos)->member.next, typeof(*(pos)), member)`|用于获取链表中紧随给定节点之后的节点所在结构体的起始地址。|
`#define list_last_entry(ptr, type, member) list_entry((ptr)->prev, type, member)`|用于获取链表中最后一个节点所在结构体的起始地址。|
|`#define list_prepare_entry(pos, ptr, member) ((pos) ? : list_entry(ptr, typeof(*pos), member))`|用于准备一个节点的数据结构指针。如果该指针为NULL，则将其初始化为链表的头部。|
|`#define list_prev_entry(pos, member) list_entry((pos)->member.prev, typeof(*(pos)), member)`|用于获取链表中紧靠给定节点之前的节点所在结构体的起始地址。|
`#define list_for_each(pos, head) for (pos = (head)->next; pos != (head); pos = pos->next)`|遍历链表中的所有节点|
|`#define list_for_each_continue(pos, head) for (pos = pos->next; pos != (head); pos = pos->next)`|从当前节点继续遍历链表中的剩余节点。|
`#define list_for_each_prev(pos, head) for (pos = (head)->prev; pos != (head); pos = pos->prev)`|从链表尾部开始遍历所有节点。|
|`#define list_for_each_safe(pos, n, head) for (pos = (head)->next, n = pos->next; pos != (head); pos = n, n = pos->next)`|与list_for_each函数类似，但允许在遍历过程中删除或添加节点。其中，n参数表示要处理的下一个节点。|
|`#define list_for_each_prev_safe(pos, n, head) for (pos = (head)->prev, n = pos->prev; pos != (head); pos = n, n = pos->prev)`|与list_for_each_safe函数类似，但遍历顺序是从链表尾部开始。|
|`#define list_for_each_entry(pos, head, member) for (pos = list_first_entry(head, typeof(*pos), member); &pos->member != (head); pos = list_next_entry(pos, member))`|用于在遍历链表时，获取每个节点所在结构体的起始地址。其中，pos参数表示当前节点所在结构体的指针；head参数表示要遍历的链表头部指针；member参数表示每个节点在结构体中的成员名称。|
|`#define list_for_each_entry_reverse(pos, head, member) for (pos = list_last_entry(head, typeof(*pos), member); &pos->member != (head); pos = list_prev_entry(pos, member))`|与list_for_each_entry函数类似，但遍历顺序是从链表尾部开始。|
|`#define list_for_each_entry_continue(pos, head, member) for (pos = list_next_entry(pos, member); &pos->member != (head); pos = list_next_entry(pos, member))`|用于从当前节点继续往后遍历链表，并获取每个节点所在结构体的起始地址。|
|`#define list_for_each_entry_continue_reverse(pos, head, member) for (pos = list_prev_entry(pos, member); &pos->member != (head); pos = list_prev_entry(pos, member))`|与list_for_each_entry_continue函数类似，但遍历顺序是从链表尾部开始。|
|`#define list_for_each_entry_from(pos, head, member) for (; &pos->member != (head); pos = list_next_entry(pos, member))`|用于从某个节点开始遍历链表，并获取每个节点所在结构体的起始地址。其中，pos参数表示当前要遍历的节点所在结构体的指针；head参数表示要遍历的链表头部指针；member参数表示每个节点在结构体中的成员名称。|
|`#define list_for_each_entry_from_reverse(pos, head, member) for (; &pos->member != (head); pos = list_prev_entry(pos, member))`|与list_for_each_entry_from函数类似，但遍历顺序是从链表尾部开始。|
|`#define list_for_each_entry_safe(pos, n, head, member) for (pos = list_first_entry(head, typeof(*pos), member), n = list_next_entry(pos, member); &pos->member != (head); pos = n, n = list_next_entry(n, member))`|与list_for_each_entry函数类似，但允许在遍历过程中删除或添加节点。其中，n参数表示要处理的下一个节点。|
|`#define list_for_each_entry_safe_continue(pos, n, head, member) for (pos = list_next_entry(pos, member), n = list_next_entry(pos, member); &pos->member != (head); pos = n, n = list_next_entry(n, member))`|用于从当前节点继续往后遍历链表，并允许在遍历过程中删除或添加节点。|
|`#define list_for_each_entry_safe_from(pos, n, head, member) for (n = list_next_entry(pos, member); &pos->member != (head); pos = n, n = list_next_entry(n, member))`|用于从某个节点开始遍历链表，并允许在遍历过程中删除或添加节点。|
|`#define list_for_each_entry_safe_reverse(pos, n, head, member) for (pos = list_last_entry(head, typeof(*pos), member), n = list_prev_entry(pos, member); &pos->member != (head); pos = n, n = list_prev_entry(n, member))`|与list_for_each_entry_reverse函数类似，但允许在遍历过程中删除或添加节点。|
|`#define list_is_first(pos, head) ((pos)->prev == (head))`|用于检查当前节点是否为链表中的第一个节点。其中，pos 参数表示要检查的节点指针；head 参数表示链表头部指针。|
|`#define list_is_last(pos, head) ((pos)->next == (head))`|用于检查当前节点是否为链表中的最后一个节点。其中，pos 参数表示要检查的节点指针；head 参数表示链表头部指针。|
|`#define list_is_head(pos, head) ((pos) == (head))`|用于检查当前节点是否为链表头部。其中，pos参数表示要检查的节点指针；head参数表示链表头部指针。|
|`#define list_is_singular(head) (!list_empty(head) && ((head)->next == (head)->prev))`|用于检查链表中是否只有一个节点。其中，head参数表示链表头部指针。|
|`#define list_safe_reset_next(curr, next, member) next = list_entry((curr)->member.next, typeof(*curr), member)`|用于安全地重置一个节点的后继节点指针，以便在遍历链表时删除当前节点。其中，curr 参数表示当前节点指针；next 参数表示当前节点的后继节点指针；member 参数表示节点结构体中 struct list_head 成员的名称。|


# 6. list_move，list_move_tail，list_bulk_move_tail
## 6.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void list_move(struct list_head *list, struct list_head *head);`|用于将一个节点移动到另外一个链表的头部。|
|`void list_move_tail(struct list_head *list, struct list_head *head);`|用于将一个节点移动到另外一个链表的尾部。|
|`void list_bulk_move_tail(struct list_head *list, int count, struct list_head *head);`|用于将多个节点从一个链表移动到另一个链表的尾部。|


**参数：**
- **list ：** 要移动的节点指针
- **head ：** 目标链表头部指针
- **count ：** 要移动的节点数量

## 6.2 演示示例

```c
#include <stdio.h>
#include <stdlib.h>
#include "list.h"

struct my_struct {
    int data;
    struct list_head list;
};

int main() {
    struct list_head a, b;
    struct my_struct s1, s2, s3, *pos, *tmp;

    // 初始化两个链表
    INIT_LIST_HEAD(&a);
    INIT_LIST_HEAD(&b);

    // 添加三个结构体到链表 a 中
    s1.data = 10;
    list_add_tail(&s1.list, &a);

    s2.data = 20;
    list_add_tail(&s2.list, &a);

    s3.data = 30;
    list_add_tail(&s3.list, &a);

    // 将节点 s1 移动到链表 b 的头部
    printf("Before move:\n");
    printf("List a: ");
    list_for_each_entry(pos, &a, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    printf("List b: ");
    list_for_each_entry(pos, &b, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    list_move(&s1.list, &b);

    printf("After move:\n");
    printf("List a: ");
    list_for_each_entry(pos, &a, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    printf("List b: ");
    list_for_each_entry(pos, &b, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    // 将节点 s2 移动到链表 b 的尾部
    printf("Before move_tail:\n");
    printf("List a: ");
    list_for_each_entry(pos, &a, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    printf("List b: ");
    list_for_each_entry(pos, &b, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    list_move_tail(&s2.list, &b);

    printf("After move_tail:\n");
    printf("List a: ");
    list_for_each_entry(pos, &a, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    printf("List b: ");
    list_for_each_entry(pos, &b, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    // 将链表 a 中的所有节点移动到链表 b 的尾部
    printf("Before bulk_move_tail:\n");
    printf("List a: ");
    list_for_each_entry(pos, &a, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    printf("List b: ");
    list_for_each_entry(pos, &b, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    list_bulk_move_tail(&a, 3, &b);

    printf("After bulk_move_tail:\n");
    printf("List a: ");
    list_for_each_entry(pos, &a, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    printf("List b: ");
    list_for_each_entry(pos, &b, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    // 释放内存
    list_for_each_entry_safe(pos, tmp, &a, list) {
        list_del(&pos->list);
        free(pos);
    }

    list_for_each_entry_safe(pos, tmp, &b, list) {
        list_del(&pos->list);
        free(pos);
    }

    return 0;
}
```
上述示例代码中，我们首先创建了两个链表 `a` 和 `b`，然后向链表 `a` 中添加三个结构体。接着，我们使用 `list_move()` 函数将节点 `s1` 从链表 `a` 移动到链表 `b` 的头部，使用 `list_move_tail()` 函数将节点 `s2` 从链表 `a` 移动到链表 `b` 的尾部，最后使用 `list_bulk_move_tail()` 函数将链表 `a` 中的所有节点都移动到链表 `b` 的尾部。

> **注意：** 在上述演示代码的最后，我们需要手动释放所有节点的内存空间，以免造成内存泄漏。

# 7. list_replace，list_replace_init
## 7.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void list_replace(struct list_head *old, struct list_head *new);`|用于用一个新节点替换指定节点。|
|`void list_replace_init(struct list_head *old, struct list_head *new);`|除了可以完成 list_replace 做的所有操作外，它还将原来的节点初始化为空。|

**参数：**
- **old ：** 要被替换的节点指针；
- **new ：** 新节点的指针。

## 7.2 演示示例
```c
#include <stdio.h>
#include "list.h"

struct my_struct {
    int data;
    struct list_head list;
};

int main() {
    struct list_head a;
    struct my_struct s1, s2, s3, s4;

    // 初始化链表
    INIT_LIST_HEAD(&a);

    // 添加三个结构体到链表中
    s1.data = 10;
    list_add_tail(&s1.list, &a);

    s2.data = 20;
    list_add_tail(&s2.list, &a);

    s3.data = 30;
    list_add_tail(&s3.list, &a);

    printf("Before replace:\n");
    list_for_each_entry(struct my_struct, pos, &a, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    // 替换第二个节点
    s4.data = 40;
    list_replace(&s2.list, &s4.list);

    printf("After replace:\n");
    list_for_each_entry(struct my_struct, pos, &a, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    // 替换第一个节点，并且清空原来的节点
    s4.data = 50;
    list_replace_init(&s1.list, &s4.list);

    printf("After replace_init:\n");
    list_for_each_entry(struct my_struct, pos, &a, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    // 释放内存
    list_del(&s3.list);
    list_del(&s4.list);

    return 0;
}
```

在上面的示例代码中，我们首先创建了一个链表 `a`，然后向其中添加三个结构体。接着，我们使用 `list_replace()` 函数将第二个节点 `s2` 替换成新节点 `s4`，并打印出替换后的链表元素；然后，我们使用 `list_replace_init()` 函数将第一个节点 `s1` 替换成新节点 `s4`，并清空原来的节点，同样打印出替换后的链表元素。

# 8. list_rotate_left，list_rotate_to_front
## 8.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`static inline void list_rotate_left(struct list_head *head)`|用于将链表向左旋转一个位置。|
|`void list_rotate_to_front(struct list_head *head, struct list_head *pivot);`|用于将指定节点移到链表头部，并旋转链表使得该节点成为新的头部。|

**参数：**
- **head ：** 链表头部指针
- **pivot ：** 要移到链表头部的节点指针。

## 8.2 演示示例
```c
#include <stdio.h>
#include "list.h"

struct my_struct {
    int data;
    struct list_head list;
};

int main() {
    struct list_head a;
    struct my_struct s1, s2, s3;

    // 初始化链表
    INIT_LIST_HEAD(&a);

    // 添加三个结构体到链表中
    s1.data = 10;
    list_add_tail(&s1.list, &a);

    s2.data = 20;
    list_add_tail(&s2.list, &a);

    s3.data = 30;
    list_add_tail(&s3.list, &a);

    printf("Before rotate:\n");
    list_for_each_entry(struct my_struct, pos, &a, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    // 将第二个节点移到链表头部，并旋转链表
    list_rotate_to_front(&a, &s2.list);

    printf("After rotate:\n");
    list_for_each_entry(struct my_struct, pos, &a, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    // 向左旋转一个位置
    list_rotate_left(&a);

    printf("After rotate_left:\n");
    list_for_each_entry(struct my_struct, pos, &a, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    // 释放内存
    list_del(&s1.list);
    list_del(&s2.list);
    list_del(&s3.list);

    return 0;
}
```
在上面的示例代码中，我们首先创建了一个链表 `a`，然后向其中添加三个结构体。接着，我们使用 `list_rotate_to_front()` 函数将第二个节点 `s2` 移到链表头部并旋转链表，打印出操作后的链表元素；然后，我们使用 `list_rotate_left()` 函数将链表向左旋转一个位置，同样打印出操作后的链表元素。

# 9. list_splice，list_splice_tail，list_splice_init，list_splice_tail_init
## 9.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
|`void list_splice(struct list_head *list, struct list_head *head);`|用于将一个链表中的所有节点插入到另一个链表的指定位置之前。|
|`void list_splice_tail(struct list_head *list, struct list_head *head);`|用于将一个链表中的所有节点插入到另一个链表的尾部。|
|`void list_splice_init(struct list_head *list, struct list_head *head);`|除了可以完成 `list_splice` 做的所有操作外，它还将原来的链表初始化为空。|
|`void list_splice_tail_init(struct list_head *list, struct list_head *head);`|除了可以完成 `list_splice_tail` 做的所有操作外，它还将原来的链表初始化为空。|
**参数：**
- **list ：** 要插入的链表头部指针
- **head ：** 
   - `list_splice()` 和 `list_splice_init()` 中表示目标链表插入的位置
   - `list_splice_tail()` 和 `list_splice_tail_init()` 中表示目标链表尾部的前一个节点

## 9.2 演示示例
```c
#include <stdio.h>
#include "list.h"

struct my_struct {
    int data;
    struct list_head list;
};

int main() {
    struct list_head a, b;
    struct my_struct s1, s2, s3, s4, s5;

    // 初始化两个链表
    INIT_LIST_HEAD(&a);
    INIT_LIST_HEAD(&b);

    // 向链表 a 中添加三个结构体
    s1.data = 10;
    list_add_tail(&s1.list, &a);

    s2.data = 20;
    list_add_tail(&s2.list, &a);

    s3.data = 30;
    list_add_tail(&s3.list, &a);

    // 向链表 b 中添加两个结构体
    s4.data = 40;
    list_add_tail(&s4.list, &b);

    s5.data = 50;
    list_add_tail(&s5.list, &b);

    printf("Before splice:\n");
    printf("List a:\n");
    list_for_each_entry(struct my_struct, pos, &a, list) {
        printf("%d ", pos->data);
    }
    printf("\nList b:\n");
    list_for_each_entry(struct my_struct, pos, &b, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    // 将链表 b 中的所有节点插入到链表 a 的头部
    list_splice(&b, &a);

    printf("After splice:\n");
    printf("List a:\n");
    list_for_each_entry(struct my_struct, pos, &a, list) {
        printf("%d ", pos->data);
    }
    printf("\nList b:\n");
    list_for_each_entry(struct my_struct, pos, &b, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    // 将链表 b 中的所有节点插入到链表 a 的尾部
    list_splice_tail(&b, &a);

    printf("After splice tail:\n");
    printf("List a:\n");
    list_for_each_entry(struct my_struct, pos, &a, list) {
        printf("%d ", pos->data);
    }
    printf("\nList b:\n");
    list_for_each_entry(struct my_struct, pos, &b, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    // 重新初始化链表 a 并将链表 b 中的所有节点插入到链表 a 的头部
    INIT_LIST_HEAD(&a);
    list_splice_init(&b, &a);

    printf("After splice init:\n");
    printf("List a:\n");
    list_for_each_entry(struct my_struct, pos, &a, list) {
        printf("%d ", pos->data);
    }
    printf("\nList b:\n");
    list_for_each_entry(struct my_struct, pos, &b, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    // 将链表 b 中的所有节点插入到链表 a 的尾部，并初始化链表 b
    INIT_LIST_HEAD(&b);
    list_splice_tail_init(&a, &b);

    printf("After splice tail init:\n");
    printf("List a:\n");
    list_for_each_entry(struct my_struct, pos, &a, list) {
        printf("%d ", pos->data);
    }
    printf("\nList b:\n");
    list_for_each_entry(struct my_struct, pos, &b, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    return 0;
}
```
上述演示代码中，我们创建了两个链表 `a` 和 `b`，并初始化为空。然后，我们向链表 `a` 中添加三个结构体，向链表 b 中添加两个结构体，并使用 `list_for_each_entry` 宏分别遍历两个链表并输出节点数据。

接着，我们使用 `list_splice()` 函数将链表 `b` 中的所有节点插入到链表 `a` 的头部，使用 `list_splice_tail()` 函数将链表 `b` 中的所有节点插入到链表 `a` 的尾部，并使用 `list_for_each_entry` 宏再次遍历两个链表并输出节点数据，可以看到链表 `a` 中包含了链表 `b` 中的所有节点。

接下来，我们使用 `INIT_LIST_HEAD` 宏重新初始化链表 `a` 并使用 `list_splice_init()` 函数将链表 `b` 中的所有节点插入到链表 `a` 的头部，使用 `INIT_LIST_HEAD` 宏重新初始化链表 `b` 并使用 `list_splice_tail_init()` 函数将链表 `a` 中的所有节点插入到链表 `b` 的尾部，并使用 `list_for_each_entry` 宏再次遍历两个链表并输出节点数据，可以看到两个链表中的节点顺序已经被重新排列。

# 10. list_swap
## 10.1 函数说明
| 函数声明 |  函数功能  |
|:--|:--|
| `void list_swap(struct list_head *list1, struct list_head *list2);`| 交换两个链表头部的位置。  |
**参数：**
**list1** 和 **list2** 分别指向两个要交换头部的链表。

## 10.2 演示示例
```c
#include <stdio.h>
#include "list.h"

struct my_struct {
    int data;
    struct list_head list;
};

int main() {
    struct list_head a, b;
    struct my_struct s1, s2, s3, s4;

    // 初始化两个链表
    INIT_LIST_HEAD(&a);
    INIT_LIST_HEAD(&b);

    // 向链表 a 中添加三个结构体
    s1.data = 10;
    list_add_tail(&s1.list, &a);

    s2.data = 20;
    list_add_tail(&s2.list, &a);

    s3.data = 30;
    list_add_tail(&s3.list, &a);

    // 向链表 b 中添加一个结构体
    s4.data = 40;
    list_add_tail(&s4.list, &b);

    printf("Before swap:\n");
    printf("List a:\n");
    list_for_each_entry(struct my_struct, pos, &a, list) {
        printf("%d ", pos->data);
    }
    printf("\nList b:\n");
    list_for_each_entry(struct my_struct, pos, &b, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    // 交换链表 a 和链表 b 的头部
    list_swap(&a, &b);

    printf("After swap:\n");
    printf("List a:\n");
    list_for_each_entry(struct my_struct, pos, &a, list) {
        printf("%d ", pos->data);
    }
    printf("\nList b:\n");
    list_for_each_entry(struct my_struct, pos, &b, list) {
        printf("%d ", pos->data);
    }
    printf("\n");

    return 0;
}
```
在上述示例中，我们创建了两个链表 `a` 和 `b`，并向链表 `a` 中添加三个节点，向链表 `b` 中添加一个节点。然后我们使用 `list_for_each_entry` 宏遍历两个链表并输出节点数据。

接着，我们使用 `list_swap()` 函数交换链表 `a` 和链表 `b` 的头部，并使用 `list_for_each_entry` 宏再次遍历两个链表并输出节点数据，可以看到链表 `a` 的头部变成了原来的链表 `b` 的头部，链表 `b` 的头部变成了原来的链表 `a` 的头部。


# 参考
1. [\[The Linux Kernel API\]](https://www.kernel.org/doc/html/v5.18/core-api/kernel-api.html)
