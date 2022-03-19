package guess

import "io"

type StackFixer interface {
	Fix(paths []Stack)
}

type StackNode interface {
	EqualsTo(another StackNode) bool
}

type Needed interface {
	NeedFix() bool
	SetNeedFix(need bool)
}

type Stack interface {
	Path() []StackNode
	SetPath(path []StackNode)
	Needed
}

type Profile interface {
	Stacks() []Stack
	WriteTo(writer io.Writer) error
	ReadFrom(reader io.Reader) error
}

// NewStackExtraInfo creates a StackExtraInfo with all field to the correct default value.
func NewStackExtraInfo() *StackExtraInfo {
	return &StackExtraInfo{
		need: true,
	}
}

// StackExtraInfo to add extra info conveniently. The implementations of the Stack do not need to care about
// the details, they just need to embed the type to themselves.
type StackExtraInfo struct {
	need bool
}

func (s *StackExtraInfo) NeedFix() bool {
	return s.need
}

func (s *StackExtraInfo) SetNeedFix(need bool) {
	s.need = need
}
