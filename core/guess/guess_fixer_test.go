package guess

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xnslong/guess-stack/core/mock"
)

func TestCommonRootFixer_Fix(t *testing.T) {
	paths := [][]int{
		{2, 3, 4, 5, 6},
		{4, 5, 6, 7},
		{4, 5, 6, 8},
		{4, 5, 6},
		{1, 2, 3, 4, 5},
	}

	c2 := mock.MakePath(paths)
	(&CommonRootFixer{2}).Fix(c2)

	expectOut := [][]int{
		{1, 2, 3, 4, 5, 6},
		{1, 2, 3, 4, 5, 6, 7},
		{1, 2, 3, 4, 5, 6, 8},
		{1, 2, 3, 4, 5, 6},
		{1, 2, 3, 4, 5},
	}

	outArray := mock.PathToArray(c2)
	assert.Equal(t, expectOut, outArray)

	mock.Print(c2)
}

func TestCommonRootFixer_FixOnMinOverlaps(t *testing.T) {
	paths := [][]int{
		{1, 2, 3, 4, 5, 6, 7},
		{6, 7, 8, 9, 10, 11, 12},
		{4, 5, 6, 7, 8, 9, 10},
		{8, 9, 10, 11, 12, 13},
	}

	c1 := mock.MakePath(paths)
	(&CommonRootFixer{1}).Fix(c1)
	assert.Equal(t, [][]int{
		{1, 2, 3, 4, 5, 6, 7},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
	}, mock.PathToArray(c1))
	mock.Print(c1)

	c2 := mock.MakePath(paths)
	(&CommonRootFixer{5}).Fix(c2)
	assert.Equal(t, [][]int{
		{1, 2, 3, 4, 5, 6, 7},
		{4, 5, 6, 7, 8, 9, 10, 11, 12},
		{4, 5, 6, 7, 8, 9, 10},
		{4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
	}, mock.PathToArray(c2))
	mock.Print(c2)
}

func TestCommonRootFixer_SelectiveFix(t *testing.T) {
	paths := [][]int{
		{2, 3, 4, 5, 6},
		{4, 5, 6, 7},
		{4, 5, 6, 8},
		{1, 2, 3, 4, 5},
	}

	c2 := mock.MakePathNeed(paths, []bool{true, false, true, false})
	(&CommonRootFixer{2}).Fix(c2)

	expectOut := [][]int{
		{1, 2, 3, 4, 5, 6},    // to fix
		{4, 5, 6, 7},          // not to fix
		{1, 2, 3, 4, 5, 6, 8}, // to fix
		{1, 2, 3, 4, 5},
	}

	outArray := mock.PathToArray(c2)
	assert.Equal(t, expectOut, outArray)

	mock.Print(c2)
}

func TestCommonRootFixer_NoLoop(t *testing.T) {
	paths := [][]int{
		{2, 3, 4, 5, 6},
		{4, 5, 6, 2, 3},
		{1, 2, 4, 5},
	}

	c2 := mock.MakePath(paths)
	(&CommonRootFixer{1}).Fix(c2)

	expectOut := [][]int{
		{1, 2, 3, 4, 5, 6},       // step 2) will join stack 2
		{1, 2, 3, 4, 5, 6, 2, 3}, // step 1) join to stack 1
		{1, 2, 4, 5},             //
	}

	outArray := mock.PathToArray(c2)
	assert.Equal(t, expectOut, outArray)

	mock.Print(c2)
}
