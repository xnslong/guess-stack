package core

import (
	"log"
	"time"

	"github.com/xnslong/guess-stack/core/guess"
)

// FixerDecorator decorates a StackFixer, so that more feature will be introduced during the fix
type FixerDecorator interface {
	Decorate(underlying guess.StackFixer) guess.StackFixer
}

type fixerFunc func(stacks []guess.Stack)

func (f fixerFunc) Fix(stacks []guess.Stack) {
	f(stacks)
}

type FixOption struct {
	Overlap   int
	BaseCount int
	MinDepth  int
	Verbose   int
}

func Fix(p guess.Profile, option FixOption) {
	finalFixer := buildFixer(option)

	stacks := p.Stacks()

	finalFixer.Fix(stacks)
}

func buildFixer(option FixOption) guess.StackFixer {
	var middle []FixerDecorator
	if option.MinDepth > 0 {
		middle = append(middle, &FixDeeperStacksDecorator{MinDepth: option.MinDepth})
	}

	if option.BaseCount > 0 {
		middle = append(middle, &WithBaseDecorator{BaseCount: option.BaseCount})
	}

	if option.Verbose > 0 {
		middle = append(middle, &VerboseDecorator{Verbose: option.Verbose})
	}

	var fixer guess.StackFixer = &guess.CommonRootFixer{MinOverlaps: option.Overlap}
	for _, m := range middle {
		fixer = m.Decorate(fixer)
	}
	return fixer
}

type FixDeeperStacksDecorator struct {
	MinDepth int
}

func (o *FixDeeperStacksDecorator) Decorate(underlying guess.StackFixer) guess.StackFixer {
	return fixerFunc(func(stacks []guess.Stack) {
		for _, stack := range stacks {
			if stack.NeedFix() {
				if len(stack.Path()) < o.MinDepth {
					stack.SetNeedFix(false)
				}
			}
		}
		underlying.Fix(stacks)
	})
}

type VerboseDecorator struct {
	Verbose int
}

func (v *VerboseDecorator) Decorate(underlying guess.StackFixer) guess.StackFixer {
	return fixerFunc(func(stacks []guess.Stack) {
		begin := time.Now()
		var nodeCount = make([]int, len(stacks))
		for i, stack := range stacks {
			nodeCount[i] = len(stack.Path())
		}

		underlying.Fix(stacks)

		count := 0
		for j, stack := range stacks {
			if nodeCount[j] != len(stack.Path()) {
				count++
			}
		}

		log.Printf("fixed stacks: %d/%d (%s elapsed)", count, len(stacks), time.Since(begin).Round(time.Millisecond))
	})
}

type WithBaseDecorator struct {
	BaseCount int
}

func (f *WithBaseDecorator) Decorate(underlying guess.StackFixer) guess.StackFixer {
	return fixerFunc(func(stacks []guess.Stack) {
		base := make([][]guess.StackNode, len(stacks))
		for i, stack := range stacks {
			path := stack.Path()
			base[i] = path[:f.BaseCount]
			stack.SetPath(path[f.BaseCount:])
		}

		underlying.Fix(stacks)

		for i, stack := range stacks {
			path := stack.Path()
			np := make([]guess.StackNode, len(base[i])+len(path))
			copy(np, base[i])
			copy(np[f.BaseCount:], path)
			stack.SetPath(np)
		}
	})
}
