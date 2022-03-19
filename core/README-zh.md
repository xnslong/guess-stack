# README

当前文档主要描述猜测栈节点的算法。

# 问题

我们把一个调用栈从根到叶子表示成一个列表。 那么调用顺序就是左侧调用右侧。因此左侧会相对更稳定，而右侧更多变。

```
Sn = [f1, f2, f3, ..., fn]
    (root)            (leaf)
```

假设每个列表左侧会有不定个数的连续节点被截断，那么列表之间就没有办法互相对齐了。 我们需要补齐这里面被截掉的节点，好让各个列表能重新互相对齐。

例如：原始栈

```
S1 = [f1, f2, f3, f4, f5, f6, f8, f9]
S2 = [f1, f2, f3, f4, f7]
```

假设`S1`左侧截断了`f1, f2`，则会变成

```
S1 = [f3, f4, f5, f6, f8, f9]      (截掉了 f1, f2)
S2 = [f1, f2, f3, f4, f7]
```

恢复出来后需要能重新变成

```
S1 = [f1, f2, f3, f4, f5, f6, f8, f9]  (恢复被截掉的 f1, f2)
S2 = [f1, f2, f3, f4, f7]
```

## 记号约定

为了下面分析方便，我们对列表本身的一些属性的表示，做如下约定

假设两个列表

```
A = [e1, e2, e3]
B = [e4, e5]
```

用`A+B`记先`A`后`B`的方式将元素连接起来形成一个新的列表，示例如下：

```
A+B = [e1, e2, e3, e4, e5]
B+A = [e4, e5, e1, e2, e3]
```

用`|A|`记列表`A`的长度，示例如下：

```
|A| = 3
|B| = 2
|A+B| = 5
|B+A| = 5
```

# 分析

如果两个栈根部原本相同， 即

```
S1 = A + A2
S2 = A + A3
```

其中`A2, A3`没有共同的前缀，

则当`S1`被截掉部分节点的时候，如果被截掉的节点是`A`的子集`A0`，假设`A = A0 + A1`，则截掉之后：

```
S1 = A1 + A2
S2 = A0 + A1 + A3
```

那么`S1`根部节点被截掉后，他里面剩余的根部节点`A1`应该能与`S2`的中间部分节点重叠。 因此，我们可以通过重叠的位置，推测出`S1`被截掉的是`A0`，从而补全它。

这里重叠部分`A1`越长，那么这种猜测成立的可能性就越大。

# 解决思路

如果栈`S1`的根部与`S2`最长重叠时可以按照如下分割：

```
S1 = A1 + A2
S2 = A0 + A1 + A3
```

其中 `|A0| > 0` 且 `A2`与 `A3` 没有共同的前缀节点，则记最长重叠（`MOR`, `max overlapping range`）如下：

```
起始：MORS(S1, S2) = |A1 + A3| 
长度：MORL(S1, S2) = |A1| 
```

其中

* `|A1 + A3|`为重叠部分的起始位置到叶子的距离，用于记录重叠部分的位置；
* `|A1|`为重叠部分本身的长度。

如果`MORL(S1, S2) = 0`，那么说明这两个栈节点没有重叠的部分，他们可能没有相同的根。

当存在多个栈时（如：`{S1, S2, S3, ..., Sn}`），对于`Si`，我们能够找到它与任意一个栈`Sj`的最长重叠（`MORL(Si, Sj)`）。
我们再记 `MORL(Si) = max{MORL(Si, Sj), j != i}`

| MORL | 1         | 2         | 3         | ... | n         | max MORL                 |
|------|-----------|-----------|-----------|-----|-----------|--------------------------|
| 1    | -         | MORL(1,2) | MORL(1,3) | ... | MORL(1,n) | MORL(1) = max{MORL(1,k)} |
| 2    | MORL(2,1) | -         | MORL(2,3) | ... | MORL(2,n) | MORL(2) = max{MORL(2,k)} |
| 3    | MORL(3,1) | MORL(3,2) | -         | ... | MORL(3,n) | MORL(3) = max{MORL(3,k)} |
| ...  | ...       | ...       | ...       | ... | ...       | ...                      |
| n    | MORL(n,1) | MORL(n,2) | MORL(n,3) | ... | -         | MORL(n) = max{MORL(n,k)} |

我们可以从`MORL(Si)`由长到短的顺序修复

## 不能循环

这里可以将上述的共享根节点过程理解为构建一颗树。共享根节点的过程即为将一个栈添加到另一个栈的某个节点的子节点列表中。

那么需要保证最终是一棵树，而不应该出现环

![loop.png](../doc/loop.png)

这样的限制总结起来就是：

1. 不应该嫁接在当前栈自己身上。
2. 嫁接的目标栈，不应该存在一条传递的嫁接链，最后重新嫁接回了当前栈。

## 不需动态计算`MOR`

假设栈`Si`的最佳猜测是补充`Sj`的部分根部节点，那么当`Si`补充相关节点后，其他的栈是否需要更新对`Si`的`MOR`？

结论是：不需要

**证明如下：**

假设 `Si` 为当前修复的栈，且最优解是参考`Sj`做修复，则可以将`Si`与`Sj`表示如下：

```
Si = A1 + A2
Sj = A0 + A1 + A3
```

其中`A2`与`A3`两个列表没有相同的前缀节点。此时 `MORL(i,j) = |A1|`

修复 `Si` 后，我们可以知道`Si`会获得`Sj`的`A0`这部分节点作为新的根节点。即：

```
Si = A0 + A1 + A2
```

假设`Sk`的`MORL(k)`会因为`Si`的修复而有所变化，那么新的更优的`MORL(k)`一定是起始于`A0`的某个节点的。

由于我们是按`MORL(i)`长度从大往小修复的，因此可以知道

```
MORL(k,j) <= MORL(k) < MORL(i) = | A1 |
```

可以知道`Sk`与`Sj`在起始`A0`这个范围内的重叠部分`A'`，一定是`A0 + A1`这个列表的子列表，而且`|A'| <= MORL(k,j) <= |A1|`。

因此可以肯定在新的`Si = A0 + A1 + A2`这个列表中，不会新产生出比当前`MORL(k,j) <= MORL(k)`更长的重叠列表。