# README

This document will specify the algorithm of guessing stacks.

# Problem

We present a call stack in a line in this document, with the `root` element on the left side, and the `leaf` element on
the right side. The right stack element is always called by the left element. So, the left elements are more stable
while the right elements are more transient in the stacks.

```
Sn: f1, f2, f3, ..., fn

  (root)            (leaf)
```

Suppose the stacks have various continuous elements trimmed from the `root` side, then the stacks cannot align to each
other. We need to recover the original call stacks by filling the missing root elements, so that they can align to each
other again.

original (& the target of recovery)

```
S1: f1, f2, f3, f4, f5, f6
S2: f1, f2, f3, f4, f7
```

after trimmed

```
S1: f3, f4, f5, f6      (trimmed f1, f2)
S2: f1, f2, f3, f4, f7
```

# Analysis

If 2 stacks share the same root stack elements, and some root elements of one stack is trimed, there is possibility
where some of the shared root elements is preserved (the remainder). So the remainder root elements of the trimmed stack
may overlap elements in the other stack.

So, if we guess 2 stacks share the same root elements originally, then we can guess the trimmed root elements.

e.g.

If the root elements (`f3, f4`) of stack `S1` overlaps stack `S2` from the 3rd element. Then we guess the 2 stacks may
share the same root elements where element `f1, f2` has been trimmed from the root of stack `S1`.

```
S1:   ,   , f3, f4, f5, f6
S2: f1, f2, f3, f4, f7
```

The more elements the overlapping range contains, the more possible the guess is true.

**So, we can guess the trimmed root elements from other stacks overlapped by the current stack's root**

# Solution

If stack `S1`'s root elements overlaps `S2` with the longest range (with `l` elements) from the `s`-elements (from leaf)
, we call it the `max overlapping range` (`MOR`), take it down as:

```
MOR(S1, S2) = (s, l) // S1's root overlaps S2, from the s'th-element and ranges with l elements.
```

if the max overlapping range of `S1`, `S2` contains `0` elements, the 2 stack may not be sharing the same root. We take
it down as:

```
MOR(S1, S2) = (s, 0)
```

When there are many stacks (suppose `{S1, S2, S3, ... Sn}`), we can get the `max overlapping range` for a stack (`Si`)
to each other stacks (`Sj, where j != i`) `MOR(i,j)`. The stack `j` with the longest MOR range length is the most
possible stack that stack `i` shares root elements with.

| MOR | 1        | 2        | 3        | ... | n        | Possible MOR  |
|-----|----------|----------|----------|-----|----------|---------------|
| 1   | -        | MOR(1,2) | MOR(1,3) | ... | MOR(1,n) | max{MOR(1,k)} |
| 2   | MOR(2,1) | -        | MOR(2,3) | ... | MOR(2,n) | max{MOR(2,k)} |
| 3   | MOR(3,1) | MOR(3,2) | -        | ... | MOR(3,n) | max{MOR(3,k)} |
| ... | ...      | ...      | ...      | ... | ...      | ...           |
| n   | MOR(n,1) | MOR(n,2) | MOR(n,3) | ... | -        | max{MOR(n,k)} |

Then we can guess the missing root for stacks with the max `MOR` to those with min `MOR`.

Suppose `MOR(i, j) = (s, l)` is considered to be the best guess for stack `i`, i.e. stack `i` shares the root stack
element of stack `j`, then the elements of stack `i` will be updated. then `MOR` of other stacks to stack `i` may need
to be recalculated. We can assert the stacks need to be updated are those stack `k`s where `MOR(k, j)` is an upper sub
range of `MOR(i,j)` (sub range nearing the leaf side).

