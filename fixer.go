package stack_fix

type Fixer interface {
	Fix(paths []Path)
}

type PathNode interface {
	EqualsTo(another PathNode) bool
}

type Path interface {
	Path() []PathNode
	SetPath(path []PathNode)
}
