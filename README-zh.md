# README

当使用 [pprof](https://github.com/google/pprof) 对一个运行的进程采样的时候，栈可能会因为太深而被截断导致根部的栈节点丢失。
当采集的 pprof 文件中，每个栈都或多或少丢失了一些根部的栈节点，那么这些栈在火焰图中就没有办法对齐了，会很难分析。

当前工具就是通过补全这些被截断丢失的栈节点的方式，修复火焰图，使栈能重新对齐，辅助pprof的分析。

```bash
guess-pprof -i before.pprof -o after.pprof
```

修复前：before.pprof:

![before.pprof](doc/before.png)

修复后：after.pprof:

![after.pprof](doc/after.png)

## 怎样补全栈节点？

设想有一个图中有两个栈：A 和 B，假设 A 栈的根部节点能够覆盖 B 栈的一部分节点，
那么我们就猜测，B 栈中重叠部分之下的根部栈节点，可能是 A 栈中被截断的部分栈节点，
于是可以将这些栈节点补全到 A 栈中。（如下图所示）

这里重叠的栈节点部分越长，这个猜测就越准确。

![guess.png](doc/guess.png)

如果你对实现算法感兴趣，可以参考[算法](fix/README-zh.md).

# 安装

```bash
go get github.com/xnslong/guess-stack/guess-pprof
```

# 使用

```
Usage of guess-pprof:
  -d int
        only fix stack with depth greater than (or equal to) the threshold, because only deep stack may be trimmed (default 1)
  -i string
        input file (default "-")
  -o string
        output file (default "-")
  -overlap int
        trustable overlap length. when the number of overlapping elements is less than the length, it's not considered trustable for guessing (default 5)
  -v    show verbose info for debug
```

```bash
guess-pprof -i before.pprof -o after.pprof
```
