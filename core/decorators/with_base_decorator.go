package decorators

import (
	"log"

	"github.com/xnslong/guess-stack/core/interfaces"
	"github.com/xnslong/guess-stack/utils"
)

type subStack []interfaces.StackNode

func (s subStack) EqualTo(another subStack) bool {
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
	Base   []subStack
	Groups [][]interfaces.Stack
}

func (g *groups) intern(base subStack, stack interfaces.Stack) {

	i := g.findGroupId(base)
	if i < 0 {
		g.Base = append(g.Base, base)
		g.Groups = append(g.Groups, nil)
		i = len(g.Base) - 1
	}

	g.Groups[i] = append(g.Groups[i], stack)
	return
}

func (g *groups) findGroupId(s subStack) int {
	for i, s2 := range g.Base {
		if s2.EqualTo(s) {
			return i
		}
	}
	return -1
}

type withBaseDecorator struct {
	BaseCount int
	Verbose   int
}

func WithBase(baseCount, verbose int) interfaces.FixerDecorator {
	return &withBaseDecorator{
		BaseCount: baseCount,
		Verbose:   verbose,
	}
}

func (f *withBaseDecorator) Decorate(underlying interfaces.StackFixer) interfaces.StackFixer {
	return interfaces.FixerFunc(func(stacks []interfaces.Stack) {
		g := &groups{}
		for _, stack := range stacks {
			path := stack.Path()
			b := utils.MinInt(f.BaseCount, len(path))
			stack.SetPath(path[b:])
			g.intern(path[:b], stack)
		}

		f.logBaseInfo(g)

		for i, group := range g.Groups {
			// fix each group respectively, so that stacks belonging to different group
			// won't refer to each other when guessing the root.
			underlying.Fix(group)
			addBaseForStacks(g.Base[i], group)
		}
	})
}

func addBaseForStacks(base subStack, group []interfaces.Stack) {
	for _, stack := range group {
		path := stack.Path()
		np := make([]interfaces.StackNode, len(base)+len(path))
		copy(np, base)
		copy(np[len(base):], path)
		stack.SetPath(np)
	}
}

func (f *withBaseDecorator) logBaseInfo(g *groups) {
	if f.Verbose > 0 {
		log.Printf("all stacks are grouped into %d groups by base nodes (depth=%d)", len(g.Base), f.BaseCount)
	}
}
