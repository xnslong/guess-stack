package guess

import "github.com/xnslong/guess-stack/core/interfaces"

type stackNode struct {
	interfaces.StackNode
}

type stack struct {
	path    []*stackNode
	index   int
	needFix bool
}

var nodes = make(map[int][]*stackNode)

func intern(node interfaces.StackNode) *stackNode {
	hash := node.HashCode()

	hashNodes := nodes[hash]
	for _, hashNode := range hashNodes {
		if hashNode.EqualsTo(node) {
			return hashNode
		}
	}

	n := &stackNode{StackNode: node}
	nodes[hash] = append(nodes[hash], n)
	return n
}

func transform(stacks []interfaces.Stack) []*stack {
	result := make([]*stack, len(stacks))

	for i, s := range stacks {
		path := s.Path()

		stackNodes := make([]*stackNode, 0, len(path))
		for _, node := range path {
			n := intern(node)
			stackNodes = append(stackNodes, n)
		}

		result[i] = &stack{
			path:    stackNodes,
			index:   i,
			needFix: s.NeedFix(),
		}
	}

	return result
}

type node struct {
	Parent    *node
	StackNode *stackNode
	Children  map[*stackNode]*node

	MyIndex     int
	JoinPathIdx int
	JoinNodeIdx int
	Overlaps    int
}

func (n *node) visitLeaf(f func(leaf *node)) {
	if len(n.Children) == 0 {
		f(n)
		return
	}

	for _, child := range n.Children {
		child.visitLeaf(f)
	}
}

func buildTrie(stacks []*stack) *node {
	root := &node{
		Parent:    nil,
		StackNode: nil,
		Children:  make(map[*stackNode]*node),
	}

	for _, s := range stacks {
		if s.needFix {
			addStack(s, root)
		}
	}

	return root
}

func addStack(stack *stack, root *node) {
	current := root

	path := stack.path

	for len(path) > 0 {
		me := path[0]
		path = path[1:]
		next, ok := current.Children[me]
		if !ok {
			next = &node{
				Parent:    current,
				StackNode: me,
				Children:  make(map[*stackNode]*node),
			}
			current.Children[me] = next
		}
		current = next
	}

	current.MyIndex = stack.index
}

func computeStackJoints(stacks []*stack, joints []*joint) {
	trie := buildTrie(stacks)

	for i, s := range stacks {
		compare(trie, i, s)
	}

	collectJoints(trie, joints)
}

func collectJoints(trie *node, joints []*joint) {
	trie.visitLeaf(func(leaf *node) {
		idx := leaf.MyIndex
		n := leaf
		for n != nil && n.Overlaps == 0 {
			n = n.Parent
		}
		if n != nil {
			joints[idx].Overlaps = n.Overlaps
			joints[idx].JoinNodeIdx = n.JoinNodeIdx
			joints[idx].JoinPathIdx = n.JoinPathIdx
		}
	})
}

func compare(trie *node, iStack int, s *stack) {
	path := s.path

	for len(path) > 0 {
		updateOverlaps(trie, iStack, path)
		path = path[1:]
	}
}

func updateOverlaps(trie *node, iStack int, path []*stackNode) {
	if len(path) == 0 {
		return
	}

	path = path[1:]
	start := len(path)

	overlap := 0

	for len(path) > 0 {
		n := path[0]
		path = path[1:]

		next, ok := trie.Children[n]
		if !ok {
			break
		}

		overlap++
		if next.Overlaps < overlap {
			next.Overlaps = overlap
			next.JoinPathIdx = iStack
			next.JoinNodeIdx = start
		}

		trie = next
	}
}
