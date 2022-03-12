package fix

type StackFixer interface {
	Fix(paths []Path, toJoin []bool)
}

type PathNode interface {
	EqualsTo(another PathNode) bool
}

type Path interface {
	Path() []PathNode
	SetPath(path []PathNode)
}
