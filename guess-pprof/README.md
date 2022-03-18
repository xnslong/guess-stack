# guess-pprof

[中文README](README-zh.md)

When [pprof](https://github.com/google/pprof) takes samples from a running program, 
it will trim some stack elements from the root end if the stack is too deep.
When one or more stack elements are trimmed from each stack, 
the stacks will not be able to align with each other, 
then the FlameGraph will be hard to read (the following 'before.pprof' graph). 

This tool will try to fix the pprof file by guessing the trimmed root elements (the following 'after.pprof' graph).

```bash
guess-pprof -i before.pprof -o after.pprof
```

before.pprof:

![before.pprof](../doc/before.png)

after.pprof:

![after.pprof](../doc/after.png)

## How to guess?

Suppose we have 2 stacks in a graph: A & B. When the root end elements of a stack A overlaps elements in stack B,
then we guess stack A should have trimmed some root elements the same as those in stack B under the overlapping elements
(shown in the following 'guess' graph).

![guess.png](../doc/guess.png)

The longer the overlapping part is, the more trustable the guess is. 

If you are interested in the algorithm to realize it, please refer to [algorithm](../fix/README.md).

# Installation

```bash
go get github.com/xnslong/guess-stack/guess-pprof
```

# Usage

```
usage: guess-pprof [<flags>]

A tool to fix the missing root call nodes of deep stacks in pprof result, so that the stacks can align with each other

Flags:
      --help         Show context-sensitive help (also try --help-long and --help-man).
  -i, --input="-"    input pprof file. "-" means read from the standard input stream
  -o, --output="-"   output pprof file, "-" means write to the standard output stream
  -O, --overlap=5    the minimal overlapping call node count
  -b, --base=0       number of the base nodes who always exist for all stacks, no matter whether the root call nodes
                     are trimmed. (such as the process name for multi-process pprof)
  -d, --depth=0      the minimal depth of the stack who may be trimmed (the deep stacks still remains deep after
                     trimmed, base nodes not counted)
  -v, --verbose ...  show verbose info on fixing the pprof
      --version      Show application version.
```

```bash
guess-pprof -i before.pprof -o after.pprof
```
