package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixDeeperStacksDecorator_Decorate(t *testing.T) {
	paths := [][]int{
		{1, 2, 3, 4, 5, 6, 7},
		{6, 7, 8, 9, 10, 11, 12},
		{4, 5, 6, 7, 8, 9, 10},
		{8, 9, 10, 11, 12, 13},
	}

	c1 := makePath(paths)
	fixer := &CommonRootFixer{1}

	d7Fixer := (&FixDeeperStacksDecorator{MinDepth: 7}).Decorate(fixer)
	d7Fixer.Fix(c1)

	assert.Equal(t, [][]int{
		{1, 2, 3, 4, 5, 6, 7},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		{8, 9, 10, 11, 12, 13}, // not fixed, depth not big enough
	}, pathToArray(c1))
	Print(c1)

	d2Fixer := (&FixDeeperStacksDecorator{MinDepth: 1}).Decorate(fixer)
	d2Fixer.Fix(c1)

	c2 := makePath(paths)
	d2Fixer.Fix(c2)
	// all fixed, depth are all big enough
	assert.Equal(t, [][]int{
		{1, 2, 3, 4, 5, 6, 7},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
	}, pathToArray(c2))
	Print(c2)
}

func TestWithBaseDecorator_Decorate(t *testing.T) {
	paths := [][]int{
		{0, 1, 2, 3, 4, 5, 6, 7},
		{0, 6, 7, 8, 9, 10, 11, 12},
		{0, 4, 5, 6, 7, 8, 9, 10},
		{0, 8, 9, 10, 11, 12, 13},
	}

	c1 := makePath(paths)
	fixer := &CommonRootFixer{1}

	(&WithBaseDecorator{BaseCount: 1}).Decorate(fixer).Fix(c1)

	// with 1 base, compute OK
	assert.Equal(t, [][]int{
		{0, 1, 2, 3, 4, 5, 6, 7},
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
	}, pathToArray(c1))
	Print(c1)

	c2 := makePath(paths)
	(&WithBaseDecorator{BaseCount: 0}).Decorate(fixer).Fix(c2)

	// with base, can not fix
	assert.Equal(t, [][]int{
		{0, 1, 2, 3, 4, 5, 6, 7},
		{0, 6, 7, 8, 9, 10, 11, 12},
		{0, 4, 5, 6, 7, 8, 9, 10},
		{0, 8, 9, 10, 11, 12, 13},
	}, pathToArray(c2))
	Print(c2)

}
