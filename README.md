# guess-stack

When taking samples of call stacks from a process in flight, some root stack nodes may be lost if the stack is too deep.
Then the stacks cannot align with each other anymore, it will very hard to analyze.

For example, following is a FlameGraph whose stacks have their root call nodes lost, it becomes hard to understand.

![mis aligned version](doc/before.png)

This tool will try to recover the stacks to the complete state by guessing the missing stack nodes, so that the stacks
can align with each other again. Following is the recovered version of the previous FlameGraph. Then we can analyze it
easily.

![aligned version](doc/after.png)

# Module introduction

Since there will be a lot of profile file formats, we will provide a variety of tools to fix each profile file type.

* [core](core): The core module providing an implementation to fix the missing nodes by guessing. This module will be
  used by each other modules.
* [guess-pprof](guess-pprof): a tool to fix the golang pprof profile. (taken
  by [pprof](https://github.com/google/pprof))
* [guess-fold](guess-fold): a tool to fix the folded stack type of call stacks defined by [FlameGraph].

[FlameGraph]: https://github.com/brendangregg/FlameGraph

# Installation

guess-pprof

```bash
go get github.com/xnslong/guess-stack/guess-pprof
```

guess-fold

```bash
go get github.com/xnslong/guess-stack/guess-fold
```

