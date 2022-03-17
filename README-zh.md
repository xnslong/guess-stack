# README

当使用 [pprof](https://github.com/google/pprof) 对一个运行的进程采样的时候，栈可能会因为太深而被截掉部分根部的栈节点。
从而导致火焰图中，栈没有办法对齐，很难分析。

当前工具就是通过猜测和补全这些被截断丢失的栈节点的方式，修复火焰图，使栈能重新对齐，便于pprof的分析。

```bash
guess-pprof -i before.pprof -o after.pprof
```

修复前：before.pprof

![before.pprof](doc/before.png "before")

修复后：after.pprof

![after.pprof](doc/after.png "after")

## 怎样补全栈节点？

设想有一个图中有两个栈：A 和 B，假设 A 栈的根部节点能够覆盖 B 栈的一部分节点，
那么我们就猜测，B 栈中重叠部分之下的根部栈节点，可能是 A 栈中被截断的部分栈节点，
于是可以将这些栈节点补全到 A 栈中。（如下图所示）

这里重叠的栈节点部分越长，这个猜测就越准确。

![guess.png](doc/guess.png "guess")

如果你对实现算法感兴趣，可以参考[算法](fix/README-zh.md).

# 安装

```bash
go get github.com/xnslong/guess-stack/guess-pprof
```

# 使用

```
usage: guess-pprof [-h|--help] [-i|--input "<value>"] [-o|--output "<value>"]
                   [--overlap <integer>] [-d|--depth <integer>] [-b|--base
                   <integer>] [-v|--verbose <integer>] [--version]

                   to guess the missing root nodes for deep stacks, so that the
                   stacks can align with each other again

Arguments:

  -h  --help     Print help information
  -i  --input    input pprof file. "-" means read from the standard input
                 stream. Default: -
  -o  --output   output pprof file, "-" means write to the standard output
                 stream. Default: -
      --overlap  the minimal overlapping call node count. Default: 5
  -d  --depth    the minimal depth of the stack who may be trimmed (the deep
                 stacks still remains deep after trimmed). Default: 1
  -b  --base     number of the base nodes always existing for all stacks (such
                 as the process name), no matter whether the root call nodes
                 are trimmed. Default: 0
  -v  --verbose 
      --version 
```

```bash
guess-pprof -i before.pprof -o after.pprof
```
