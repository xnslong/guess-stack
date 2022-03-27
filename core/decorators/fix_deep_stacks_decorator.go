package decorators

import (
	"log"

	"github.com/xnslong/guess-stack/core/interfaces"
)

type fixDeepStacksDecorator struct {
	MinDepth int
	Verbose  int
}

func FixDeepStacks(minDepth, verbose int) interfaces.FixerDecorator {
	return &fixDeepStacksDecorator{
		MinDepth: minDepth,
		Verbose:  verbose,
	}
}

func (d *fixDeepStacksDecorator) Decorate(underlying interfaces.StackFixer) interfaces.StackFixer {
	return interfaces.FixerFunc(func(stacks []interfaces.Stack) {
		notNeed := 0
		for _, stack := range stacks {
			if stack.NeedFix() {
				if len(stack.Path()) < d.MinDepth {
					stack.SetNeedFix(false)
					notNeed++
				}
			}
		}
		if d.Verbose > 1 {
			log.Printf("stacks not deep enough are skipped: %d/%d", notNeed, len(stacks))
		}
		underlying.Fix(stacks)
	})
}
