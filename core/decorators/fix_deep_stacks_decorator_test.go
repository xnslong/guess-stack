package decorators

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/xnslong/guess-stack/core/guess"
	"github.com/xnslong/guess-stack/core/mock"
)

func TestFixDeeperStacksDecorator_Decorate(t *testing.T) {
	paths := [][]int{
		{1, 2, 3, 4, 5, 6, 7},
		{6, 7, 8, 9, 10, 11, 12},
		{4, 5, 6, 7, 8, 9, 10},
		{8, 9, 10, 11, 12, 13},
	}

	c1 := mock.MakePath(paths)
	fixer := &guess.CommonRootFixer{MinOverlaps: 1}

	d7Fixer := FixDeepStacks(7, 0).Decorate(fixer)
	d7Fixer.Fix(c1)

	assert.Equal(t, [][]int{
		{1, 2, 3, 4, 5, 6, 7},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		{8, 9, 10, 11, 12, 13}, // not fixed, depth not big enough
	}, mock.PathToArray(c1))
	mock.Print(c1)

	d2Fixer := FixDeepStacks(1, 0).Decorate(fixer)
	d2Fixer.Fix(c1)

	c2 := mock.MakePath(paths)
	d2Fixer.Fix(c2)
	// all fixed, depth are all big enough
	assert.Equal(t, [][]int{
		{1, 2, 3, 4, 5, 6, 7},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
	}, mock.PathToArray(c2))
	mock.Print(c2)
}
