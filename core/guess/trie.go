package guess

import (
	"github.com/xnslong/guess-stack/core/interfaces"
)

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

type trieNode struct {
	Parent    *trieNode
	StackNode *stackNode
	Children  map[*stackNode]*trieNode

	MyIndex     []int // leaf node will record the stack index, other nodes will not
	JoinPathIdx int
	JoinNodeIdx int
	Overlaps    int
}

func (n *trieNode) visitLeaf(f func(leaf *trieNode)) {
	if len(n.MyIndex) > 0 {
		f(n)
	}

	for _, child := range n.Children {
		child.visitLeaf(f)
	}
}

func (n *trieNode) addStack(path []*stackNode, id int) {
	if len(path) == 0 {
		n.MyIndex = append(n.MyIndex, id)
		return
	}

	node := path[0]
	next, ok := n.Children[node]
	if !ok {
		next = &trieNode{
			Parent:      n,
			StackNode:   node,
			Children:    make(map[*stackNode]*trieNode),
			JoinPathIdx: nonExistIndex,
			JoinNodeIdx: 0,
		}
		n.Children[node] = next
	}

	next.addStack(path[1:], id)
}

func buildTrie(stacks []*stack) *trieNode {
	root := &trieNode{
		Parent:      nil,
		StackNode:   nil,
		Children:    make(map[*stackNode]*trieNode),
		JoinPathIdx: nonExistIndex,
	}

	c := 0

	for _, s := range stacks {
		if s.needFix {
			root.addStack(s.path, s.index)
			c++
		}
	}

	return root
}

func computeStackJoints(stacks []*stack, joints []*joint) {
	trie := buildTrie(stacks)

	for i, s := range stacks {
		compare(trie, i, s)
	}

	collectJoints(trie, joints)
}

func collectJoints(trie *trieNode, joints []*joint) {
	trie.visitLeaf(func(leaf *trieNode) {
		indexes := leaf.MyIndex
		n := leaf
		for n != nil && n.Overlaps == 0 {
			n = n.Parent
		}
		if n != nil {
			for _, idx := range indexes {
				joints[idx].Overlaps = n.Overlaps
				joints[idx].JoinNodeIdx = n.JoinNodeIdx
				joints[idx].JoinPathIdx = n.JoinPathIdx
			}
		}
	})
}

func compare(trie *trieNode, iStack int, s *stack) {
	path := s.path

	for len(path) > 0 {
		updateOverlaps(trie, iStack, path)
		path = path[1:]
	}
}

func updateOverlaps(trie *trieNode, iStack int, path []*stackNode) {
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
