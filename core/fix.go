package core

import (
	"log"
)

// FixerDecorator decorates a StackFixer, so that more feature will be introduced during the fix
type FixerDecorator interface {
	Decorate(underlying StackFixer) StackFixer
}

type fixerFunc func(stacks []Stack)

func (f fixerFunc) Fix(stacks []Stack) {
	f(stacks)
}

type FixOption struct {
	Overlap   int
	BaseCount int
	MinDepth  int
	Verbose   int
}

func Fix(p Profile, option FixOption) {
	fixer := &CommonRootFixer{MinOverlaps: option.Overlap}
	FixProfile(p, fixer, option)
}

func FixProfile(p Profile, sf StackFixer, option FixOption) {
	finalFixer := buildFixer(option, sf)

	stacks := p.Stacks()
	finalFixer.Fix(stacks)
}

func buildFixer(option FixOption, sf StackFixer) StackFixer {
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

	var finalFixer = sf

	for _, m := range middle {
		finalFixer = m.Decorate(finalFixer)
	}
	return finalFixer
}

type FixDeeperStacksDecorator struct {
	MinDepth int
}

func (o *FixDeeperStacksDecorator) Decorate(underlying StackFixer) StackFixer {
	return fixerFunc(func(stacks []Stack) {
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

func (v *VerboseDecorator) Decorate(underlying StackFixer) StackFixer {
	return fixerFunc(func(stacks []Stack) {
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

		log.Printf("fixed stacks: %d/%d", count, len(stacks))
	})
}

type WithBaseDecorator struct {
	BaseCount int
}

func (f *WithBaseDecorator) Decorate(underlying StackFixer) StackFixer {
	return fixerFunc(func(stacks []Stack) {
		base := make([][]StackNode, len(stacks))
		for i, stack := range stacks {
			path := stack.Path()
			base[i] = path[:f.BaseCount]
			stack.SetPath(path[f.BaseCount:])
		}

		underlying.Fix(stacks)

		for i, stack := range stacks {
			path := stack.Path()
			np := make([]StackNode, len(base[i])+len(path))
			copy(np, base[i])
			copy(np[f.BaseCount:], path)
			stack.SetPath(np)
		}
	})
}
