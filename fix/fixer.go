package fix

type StackFixer interface {
	Fix(paths []Stack, needFix []bool)
}

type StackNode interface {
	EqualsTo(another StackNode) bool
}

type Stack interface {
	Path() []StackNode
	SetPath(path []StackNode)
}
