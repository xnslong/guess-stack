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
go get github.com/xnslong/guess-stack/guess-pprof
```

# Usage

```bash
Usage of ./guess-pprof:
  -i string
        input file (default "-")
  -o string
        output file (default "-")
  -overlap int
        trustable overlap length. when the number of overlapping elements is less than the length, it's not considered trustable for guessing (default 5)
  -v    show verbose info for debug
```

```bash
./guess-pprof -i before.pprof -o after.pprof
```