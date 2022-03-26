package interfaces

import "io"

type StackFixer interface {
	Fix(paths []Stack)
}

type StackNode interface {
	EqualsTo(another StackNode) bool
	HashCode() int
}

type Stack interface {
	Path() []StackNode
	SetPath(path []StackNode)
	Extra
}

type Profile interface {
	Stacks() []Stack
	WriteTo(writer io.Writer) error
	ReadFrom(reader io.Reader) error
}
