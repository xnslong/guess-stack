# README

When pprof takes sample of a running program, it will trim the stack elements from the root if the stack is too deep.
When one or more stack nodes are trimmed, the stack will not be able to align with each other, then the FlameGraph will be hard to read (the following 'before' graph). 

This tool will try to fix the pprof file by guessing the trimmed root nodes (the following 'After' graph).

Before:

![original](doc/before.png)

After:

![img.png](doc/after.png)

# Install

```bash
go get github.com/xnslong/guess-stack
```

# Usage

```bash
./guess-stack -h
Usage of ./guess-stack:
  -i string
        input file (default "-")
  -o string
        output file (default "-")
  -overlap int
        trustable overlap count (default 5)
  -v    show verbose info for debug
```

```bash
./guess-stack -i before.pprof -o after.pprof
```