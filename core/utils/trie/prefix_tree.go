package trie

type Trie interface {
	AddPath(path []interface{}, target interface{})
	PrefixFor([]interface{}) []Node
	VisitAllPath(f func(path []Node, target interface{}))
}

type Node interface {
	Element() interface{}    // to identify the node
	Attachment() interface{} // attachment to this node
	Attach(val interface{})  // set attachment to this node
}

type trie struct {
	root *node
}

func NewTrie() Trie {
	return &trie{
		root: newNode(nil),
	}
}

func (t *trie) AddPath(path []interface{}, target interface{}) {
	n := t.root

	for _, elem := range path {
		next := n.children[elem]
		if next == nil {
			next = newNode(elem)
			n.children[elem] = next
		}
		n = next
	}

	n.targets = append(n.targets, target)
}

func (t *trie) PrefixFor(i []interface{}) []Node {
	n := t.root

	path := make([]Node, 0)

	for _, elem := range i {
		next := n.children[elem]
		if next == nil {
			return path
		}

		path = append(path, next)
		n = next
	}

	return path
}

func (t *trie) VisitAllPath(f func(path []Node, target interface{})) {
	for _, n := range t.root.children {
		n.doVisit(&memoto{}, f)
	}
}

type memoto struct {
	path []Node
}

type node struct {
	children   map[interface{}]*node
	identity   interface{}
	targets    []interface{}
	attachment interface{}
}

func newNode(identity interface{}) *node {
	return &node{
		identity: identity,
		children: map[interface{}]*node{},
	}
}

func (n *node) Element() interface{} {
	return n.identity
}

func (n *node) Attachment() interface{} {
	return n.attachment
}

func (n *node) Attach(val interface{}) {
	n.attachment = val
}

func (n *node) doVisit(mem *memoto, f func(path []Node, target interface{})) {
	l := len(mem.path)
	mem.path = append(mem.path, n)

	// iterate over all targets for paths ends with current node
	for _, target := range n.targets {
		f(mem.path, target)
	}

	for _, n2 := range n.children {
		n2.doVisit(mem, f)
	}

	mem.path = mem.path[:l]
}
