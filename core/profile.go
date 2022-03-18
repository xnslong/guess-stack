package core

import "io"

type StackFixer interface {
	Fix(paths []Stack)
}

type StackNode interface {
	EqualsTo(another StackNode) bool
}

type Stack interface {
	Path() []StackNode
	SetPath(path []StackNode)
	NeedFix() bool
	SetNeedFix(need bool)
}

type Profile interface {
	Stacks() []Stack
	WriteTo(writer io.Writer) error
	ReadFrom(reader io.Reader) error
}
