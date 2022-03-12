package fix

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type IntNode int

func (i IntNode) EqualsTo(another PathNode) bool {
	return i == another
}

type IntPath struct {
	arr []PathNode
}

func (i *IntPath) Path() []PathNode {
	return i.arr
}

func (i *IntPath) SetPath(path []PathNode) {
	i.arr = path
}

func Print(p []Path) {
	fmt.Println("---- print paths ----")
	for i, path := range p {
		fmt.Print(i, " {")
		nodes := path.Path()
		for _, node := range nodes {
			fmt.Print(node, ",")
		}
		fmt.Println("},")
	}
}

func makePath(ints [][]int) []Path {
	result := make([]Path, 0, len(ints))

	for _, arr := range ints {
		var path []PathNode
		for _, v := range arr {
			path = append(path, IntNode(v))
		}
		result = append(result, &IntPath{path})
	}

	return result
}

func pathToArray(path []Path) [][]int {
	result := make([][]int, 0, len(path))

	for _, p := range path {
		var stack []int
		for _, v := range p.Path() {
			iv := v.(IntNode)
			stack = append(stack, int(iv))
		}
		result = append(result, stack)
	}

	return result
}

func TestFix(t *testing.T) {
	paths := [][]int{
		{1, 2, 3, 4, 5, 6, 7},
		{6, 7, 8, 9, 10, 11, 12},
		{4, 5, 6, 7, 8, 9, 10},
		{8, 9, 10, 11, 12, 13},
	}

	c1 := makePath(paths)
	(&CommonRootFixer{1}).Fix(c1, []bool{true, true, true, true})
	assert.Equal(t, [][]int{
		{1, 2, 3, 4, 5, 6, 7},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
	}, pathToArray(c1))
	Print(c1)

	c2 := makePath(paths)
	(&CommonRootFixer{5}).Fix(c2, []bool{true, true, true, true})
	assert.Equal(t, [][]int{
		{1, 2, 3, 4, 5, 6, 7},
		{4, 5, 6, 7, 8, 9, 10, 11, 12},
		{4, 5, 6, 7, 8, 9, 10},
		{4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
	}, pathToArray(c2))
	Print(c2)
}

func TestCommonRootFixer_Fix(t *testing.T) {
	paths := [][]int{
		{2, 3, 4, 5, 6},
		{4, 5, 6, 7},
		{4, 5, 6, 8},
		{1, 2, 3, 4, 5},
	}

	c2 := makePath(paths)
	(&CommonRootFixer{2}).Fix(c2, []bool{true, true, true, true})

	expectOut := [][]int{
		{1, 2, 3, 4, 5, 6},
		{1, 2, 3, 4, 5, 6, 7},
		{1, 2, 3, 4, 5, 6, 8},
		{1, 2, 3, 4, 5},
	}

	outArray := pathToArray(c2)
	assert.Equal(t, expectOut, outArray)

	Print(c2)
}

func TestCommonRootFixer_SelectiveFix(t *testing.T) {
	paths := [][]int{
		{2, 3, 4, 5, 6},
		{4, 5, 6, 7},
		{4, 5, 6, 8},
		{1, 2, 3, 4, 5},
	}

	c2 := makePath(paths)
	(&CommonRootFixer{2}).Fix(c2, []bool{true, false, true, true})

	expectOut := [][]int{
		{1, 2, 3, 4, 5, 6},    // to fix
		{4, 5, 6, 7},          // not to fix
		{1, 2, 3, 4, 5, 6, 8}, // to fix
		{1, 2, 3, 4, 5},
	}

	outArray := pathToArray(c2)
	assert.Equal(t, expectOut, outArray)

	Print(c2)
}
