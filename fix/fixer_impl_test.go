package fix

import (
	"fmt"
	"testing"
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
		fmt.Println(i, path.Path())
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

func TestFix(t *testing.T) {
	paths := [][]int{
		{1, 2, 3, 5},
		{5, 1, 3, 4},
		{3, 5, 3, 4},
	}

	c1 := makePath(paths)
	Print(c1)
	(&CommonRootFixer{1}).Fix(c1)
	Print(c1)

	c2 := makePath(paths)
	(&CommonRootFixer{2}).Fix(c2)
	Print(c2)
}
