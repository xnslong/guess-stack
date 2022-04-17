package guess

import (
	"github.com/xnslong/guess-stack/core/interfaces"
)

type idStackNode struct {
	interfaces.StackNode
}

type idStack struct {
	path []interface{}
}

func (s *idStack) Stack() []interface{} {
	return s.path
}

var defaultTransformer = newTransformer()

func newTransformer() *transformer {
	return &transformer{
		nodes: make(map[int][]*idStackNode),
	}
}

type transformer struct {
	nodes map[int][]*idStackNode
}

func (t *transformer) intern(node interfaces.StackNode) *idStackNode {
	hash := node.HashCode()

	hashNodes := t.nodes[hash]
	for _, hashNode := range hashNodes {
		if hashNode.EqualsTo(node) {
			return hashNode
		}
	}

	n := &idStackNode{StackNode: node}
	t.nodes[hash] = append(t.nodes[hash], n)
	return n
}

func (t *transformer) Transform(stacks []interfaces.Stack) []*idStack {
	result := make([]*idStack, len(stacks))

	for i, s := range stacks {
		path := s.Path()

		stackNodes := make([]interface{}, 0, len(path))
		for _, node := range path {
			n := t.intern(node)
			stackNodes = append(stackNodes, n)
		}

		result[i] = &idStack{
			path: stackNodes,
		}
	}

	return result
}
