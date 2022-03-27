package core

import (
	"github.com/xnslong/guess-stack/core/decorators"
	"github.com/xnslong/guess-stack/core/guess"
	"github.com/xnslong/guess-stack/core/interfaces"
)

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
	var middle []interfaces.FixerDecorator

	if option.Verbose > 0 {
		middle = append(middle, decorators.LogFixInfo(option.Verbose))
	}

	if option.BaseCount > 0 {
		middle = append(middle, decorators.WithBase(option.BaseCount, option.Verbose))
	}

	if option.MinDepth > 0 {
		middle = append(middle, decorators.FixDeepStacks(option.MinDepth, option.Verbose))
	}

	var fixer interfaces.StackFixer = &guess.CommonRootFixer{MinOverlaps: option.Overlap}
	for i := len(middle) - 1; i >= 0; i-- {
		m := middle[i]
		fixer = m.Decorate(fixer)
	}
	return fixer
}
