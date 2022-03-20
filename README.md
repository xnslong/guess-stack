# guess-stack

When taking samples of call stacks from a process in flight, some root stack nodes may be lost if the stack is too deep.
Then the stacks cannot align with each other anymore, it will be hard to analyze.

For example, following is a FlameGraph whose stacks have their root call nodes lost, it becomes hard to understand.

![mis-aligned version](doc/before.png)

This tool will try to recover the stacks by guessing the missing stack nodes, enabling the stacks to align with each other again. Following is the recovered version of the previous FlameGraph. Then we can analyze it easily.

![aligned version](doc/after.png)

## How to guess?

Suppose we have 2 stacks in a graph: A & B. When the root end nodes of a stack A overlaps nodes in stack B,
then we guess stack A should have trimmed some root nodes the same as those in stack B under the overlapping nodes
(shown in the following 'guess' graph).

![guess](doc/guess.png)

The longer the overlapping part is, the more trustable the guess is. 

If you are interested in the algorithm to realize it, please refer to [algorithm](core/README.md).


# Module introduction

Since there are a variety of profile file types, we will provide several tools to fix each profile file type.

* [guess-pprof](guess-pprof/README.md): a tool to fix the golang pprof profile. (taken
  by [pprof](https://github.com/google/pprof))
* [guess-fold](guess-fold/README.md): a tool to fix the folded stack type of call stacks defined by [FlameGraph]. it's will be more commonly used for programs written in a variety oflanguages.

[FlameGraph]: https://github.com/brendangregg/FlameGraph

The previous tools are working based on the following core module.

* [core](core): The core module providing an implementation to fix the missing nodes by guessing. This module will be
  used by each other modules.


# Installation

guess-pprof

```bash
go get github.com/xnslong/guess-stack/guess-pprof
```

guess-fold

```bash
go get github.com/xnslong/guess-stack/guess-fold
```

