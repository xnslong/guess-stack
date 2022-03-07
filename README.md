# README

When [pprof](https://github.com/google/pprof) takes samples from a running program, 
it will trim some stack elements from the root side if the stack is too deep.
When one or more stack elements are trimmed from each stack, 
the stacks will not be able to align with each other, 
then the FlameGraph will be hard to read (the following 'before' graph). 

This tool will try to fix the pprof file by guessing the trimmed root nodes (the following 'After' graph).


Before:

![before.png](doc/before.png)

After:

![after.png](doc/after.png)

## How to guess?

When the root side elements of a stack (A) overlaps elements in another stack (B),
then we guess stack A should have trimmed root elements the same as those in stack B under the overlapping elements
(shown in the following 'guess' graph).

![guess.png](doc/guess.png)

The longer the overlapping part is, the more trustable the guess is.

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