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
	if option.MinDepth > 0 {
		middle = append(middle, &FixDeeperStacksDecorator{&option})
	}

	if option.BaseCount > 0 {
		middle = append(middle, &WithBaseDecorator{&option})
	}

	if option.Verbose > 0 {
		middle = append(middle, &ShowFixInfoDecorator{&option})
	}

	var fixer interfaces.StackFixer = &guess.CommonRootFixer{MinOverlaps: option.Overlap}
	for _, m := range middle {
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
	G     []substack
	Count []int
}

func (g *groups) intern(s substack) (group int) {
	defer func() {
		if len(g.Count) < len(s)+1 {
			newArr := make([]int, len(s)+1)
			copy(newArr, g.Count)
			g.Count = newArr
		}
		g.Count[len(s)]++
	}()

	for i, s2 := range g.G {
		if s2.EqualTo(s) {
			return i
		}
	}

	g.G = append(g.G, s)
	return len(g.G) - 1
}

func (f *WithBaseDecorator) Decorate(underlying interfaces.StackFixer) interfaces.StackFixer {
	return fixerFunc(func(stacks []interfaces.Stack) {
		g := &groups{}
		base := make([][]interfaces.StackNode, len(stacks))
		for i, stack := range stacks {
			path := stack.Path()
			b := utils.MinInt(f.BaseCount, len(path))

			base[i] = path[:b]
			stack.SetPath(path[b:])

			gId := g.intern(base[i])
			stack.SetGroup(gId)
		}

		f.logBaseInfo(g)

		underlying.Fix(stacks)

		for i, stack := range stacks {
			path := stack.Path()
			np := make([]interfaces.StackNode, len(base[i])+len(path))
			copy(np, base[i])
			copy(np[len(base[i]):], path)
			stack.SetPath(np)
		}
	})
}

func (f *WithBaseDecorator) logBaseInfo(g *groups) {
	if f.Verbose > 1 {
		log.Printf("all stacks are grouped into %d groups by base nodes (depth=%d)", len(g.G), f.BaseCount)
	}
	if f.Verbose > 2 {
		for i, count := range g.Count {
			if count > 0 {
				log.Printf("\t%d samples counted %d base nodes", count, i)
			}
		}
	}
}
