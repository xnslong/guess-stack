package decorators

import (
	"log"
	"time"

	"github.com/xnslong/guess-stack/core/interfaces"
)

type showFixInfoDecorator struct {
	Verbose int
}

func LogFixInfo(verbose int) interfaces.FixerDecorator {
	return &showFixInfoDecorator{Verbose: verbose}
}

func (v *showFixInfoDecorator) Decorate(underlying interfaces.StackFixer) interfaces.StackFixer {
	return interfaces.FixerFunc(func(stacks []interfaces.Stack) {
		if v.Verbose == 0 {
			underlying.Fix(stacks)
			return
		}

		fixAndLogInfo(stacks, underlying)
	})
}

func fixAndLogInfo(stacks []interfaces.Stack, underlying interfaces.StackFixer) {
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
}
