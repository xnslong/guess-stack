package decorators

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/xnslong/guess-stack/core/guess"
	"github.com/xnslong/guess-stack/core/mock"
)

func TestWithBaseDecorator_Decorate(t *testing.T) {
	paths := [][]int{
		{0, 1, 2, 3, 4, 5, 6, 7},
		{0, 6, 7, 8, 9, 10, 11, 12},
		{0, 4, 5, 6, 7, 8, 9, 10},
		{0, 8, 9, 10, 11, 12, 13},
	}

	c1 := mock.MakePath(paths)
	fixer := &guess.CommonRootFixer{MinOverlaps: 1}

	WithBase(1, 0).Decorate(fixer).Fix(c1)

	// with 1 base, compute OK
	assert.Equal(t, [][]int{
		{0, 1, 2, 3, 4, 5, 6, 7},
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
	}, mock.PathToArray(c1))
	mock.Print(c1)

	c2 := mock.MakePath(paths)
	WithBase(0, 0).Decorate(fixer).Fix(c2)

	// with base, can not fix
	assert.Equal(t, [][]int{
		{0, 1, 2, 3, 4, 5, 6, 7},
		{0, 6, 7, 8, 9, 10, 11, 12},
		{0, 4, 5, 6, 7, 8, 9, 10},
		{0, 8, 9, 10, 11, 12, 13},
	}, mock.PathToArray(c2))
	mock.Print(c2)

}
