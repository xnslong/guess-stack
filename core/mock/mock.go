package mock

import (
	"fmt"

	"github.com/xnslong/guess-stack/core/interfaces"
)

type IntNode int

func (i IntNode) EqualsTo(another interfaces.StackNode) bool {
	return i == another
}

type IntPath struct {
	arr []interfaces.StackNode
	*interfaces.StackExtraInfo
}

func (i *IntPath) Path() []interfaces.StackNode {
	return i.arr
}

func (i *IntPath) SetPath(path []interfaces.StackNode) {
	i.arr = path
}

func Print(p []interfaces.Stack) {
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

func MakePath(ints [][]int) []interfaces.Stack {
	needed := make([]bool, len(ints))
	for i := 0; i < len(needed); i++ {
		needed[i] = true
	}

	return MakePathNeed(ints, needed)
}

func MakePathNeed(ints [][]int, needFix []bool) []interfaces.Stack {
	result := make([]interfaces.Stack, 0, len(ints))

	for i, arr := range ints {
		var path []interfaces.StackNode
		for _, v := range arr {
			path = append(path, IntNode(v))
		}
		intPath := &IntPath{arr: path, StackExtraInfo: interfaces.NewStackExtraInfo()}
		intPath.SetNeedFix(needFix[i])
		result = append(result, intPath)
	}

	return result
}

func PathToArray(path []interfaces.Stack) [][]int {
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
