package interfaces

type Extra interface {
	Needed
	Grouped
}

type Needed interface {
	NeedFix() bool
	SetNeedFix(need bool)
}

type Grouped interface {
	Group() int     // stack belonging to different group won't join each other.
	SetGroup(g int) // set stack group. stack belonging to different group won't join each other.
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
	need  bool
	group int
}

func (s *StackExtraInfo) Group() int {
	return s.group
}

func (s *StackExtraInfo) SetGroup(g int) {
	s.group = g
}

func (s *StackExtraInfo) NeedFix() bool {
	return s.need
}

func (s *StackExtraInfo) SetNeedFix(need bool) {
	s.need = need
}
