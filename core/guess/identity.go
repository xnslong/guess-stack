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

// transformer transforms interfaces.StackNode to *idStackNode.
//
// For example, any 2 input node n1 and n2, where
//    n1.HashCode() == n2.HashCode()
//    n1.EqualsTo(n2)
// Then they will get the same *idStackNode pointer, so that they are comparable with the "==" and "!=" method.
//    t.intern(n1) == t.intern(n2)
// The Transform method will transfer all elements in the stacks into the *idStackNode type.
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
