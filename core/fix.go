package core

import (
	"log"
	"time"

	"github.com/xnslong/guess-stack/core/guess"
	"github.com/xnslong/guess-stack/core/interfaces"
	"github.com/xnslong/guess-stack/utils"
)

// FixerDecorator decorates a StackFixer, so that more feature will be introduced during the fix
type FixerDecorator interface {
	Decorate(underlying interfaces.StackFixer) interfaces.StackFixer
}

type fixerFunc func(stacks []interfaces.Stack)

func (f fixerFunc) Fix(stacks []interfaces.Stack) {
	f(stacks)
}

type FixOption struct {
	Overlap   int
	BaseCount int
	MinDepth  int
	Verbose   int
}

func Fix(p interfaces.Profile, option FixOption) {
	finalFixer := buildFixer(option)

	stacks := p.Stacks()

	finalFixer.Fix(stacks)
}

func buildFixer(option FixOption) interfaces.StackFixer {
	var middle []FixerDecorator

	if option.Verbose > 0 {
		middle = append(middle, &ShowFixInfoDecorator{&option})
	}

	if option.BaseCount > 0 {
		middle = append(middle, &WithBaseDecorator{&option})
	}

	if option.MinDepth > 0 {
		middle = append(middle, &FixDeeperStacksDecorator{&option})
	}

	var fixer interfaces.StackFixer = &guess.CommonRootFixer{MinOverlaps: option.Overlap}
	for i := len(middle) - 1; i >= 0; i-- {
		m := middle[i]
		fixer = m.Decorate(fixer)
	}
	return fixer
}

type FixDeeperStacksDecorator struct {
	*FixOption
}

func (o *FixDeeperStacksDecorator) Decorate(underlying interfaces.StackFixer) interfaces.StackFixer {
	return fixerFunc(func(stacks []interfaces.Stack) {
		notNeed := 0
		for _, stack := range stacks {
			if stack.NeedFix() {
				if len(stack.Path()) < o.MinDepth {
					stack.SetNeedFix(false)
					notNeed++
				}
			}
		}
		if o.Verbose > 1 {
			log.Printf("stacks not deep enough are skipped: %d/%d", notNeed, len(stacks))
		}
		underlying.Fix(stacks)
	})
}

type ShowFixInfoDecorator struct {
	*FixOption
}

func (v *ShowFixInfoDecorator) Decorate(underlying interfaces.StackFixer) interfaces.StackFixer {
	return fixerFunc(func(stacks []interfaces.Stack) {
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
	*FixOption
}

type substack []interfaces.StackNode

func (s substack) EqualTo(another substack) bool {
	if len(s) != len(another) {
		return false
	}

	for i, node := range s {
		if !node.EqualsTo(another[i]) {
			return false
		}
	}

	return true
}

type groups struct {
	Base   []substack
	Groups [][]interfaces.Stack
}

func (g *groups) intern(base substack, stack interfaces.Stack) {

	i := g.findGroupId(base)
	if i < 0 {
		g.Base = append(g.Base, base)
		g.Groups = append(g.Groups, nil)
		i = len(g.Base) - 1
	}

	g.Groups[i] = append(g.Groups[i], stack)
	return
}

func (g *groups) findGroupId(s substack) int {
	for i, s2 := range g.Base {
		if s2.EqualTo(s) {
			return i
		}
	}
	return -1
}

func (f *WithBaseDecorator) Decorate(underlying interfaces.StackFixer) interfaces.StackFixer {
	return fixerFunc(func(stacks []interfaces.Stack) {
		g := &groups{}
		for _, stack := range stacks {
			path := stack.Path()
			b := utils.MinInt(f.BaseCount, len(path))
			stack.SetPath(path[b:])
			g.intern(path[:b], stack)
		}

		f.logBaseInfo(g)

		for i, group := range g.Groups {
			underlying.Fix(group)
			base := g.Base[i]

			for _, stack := range group {
				path := stack.Path()
				np := make([]interfaces.StackNode, len(base)+len(path))
				copy(np, base)
				copy(np[len(base):], path)
				stack.SetPath(np)
			}
		}
	})
}

func (f *WithBaseDecorator) logBaseInfo(g *groups) {
	if f.Verbose > 0 {
		log.Printf("all stacks are grouped into %d groups by base nodes (depth=%d)", len(g.Base), f.BaseCount)
	}
}
