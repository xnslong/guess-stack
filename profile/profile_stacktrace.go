package profile

import (
	"fmt"
	"log"

	"github.com/google/pprof/profile"
	"github.com/xnslong/guess-stack/fix"
)

type StackTraceElement struct {
	*profile.Location
}

func (l *StackTraceElement) String() string {
	return fmt.Sprintf("%d", l.ID)
}

func (l *StackTraceElement) EqualsTo(another fix.PathNode) bool {
	anotherLoc, ok := another.(*StackTraceElement)
	if !ok {
		return false
	}

	if len(l.Line) == 0 {
		return false
	}

	if len(anotherLoc.Line) == 0 {
		return false
	}

	return l.ID == anotherLoc.ID

	//thisLine := l.Line[0]
	//anotherLine := anotherLoc.Line[0]
	//return thisLine.Line == anotherLine.Line // && thisLine.Function.ID == anotherLine.Function.ID
}

type StackTrace struct {
	Elements []*StackTraceElement
}

func (s *StackTrace) String() string {
	return fmt.Sprintf("%s", s.Elements)
}

func (s *StackTrace) Path() []fix.PathNode {
	pn := make([]fix.PathNode, 0, len(s.Elements))

	for _, loc := range s.Elements {
		pn = append(pn, loc)
	}

	return pn
}

func (s *StackTrace) SetPath(path []fix.PathNode) {
	v := s.Elements[:0]
	for _, node := range path {
		if loc, ok := node.(*StackTraceElement); ok {
			v = append(v, loc)
		} else {
			log.Panicf("invalid type, expect *Loc, but got %T", node)
		}
	}
	s.Elements = v
}

func reverse(elements []*StackTraceElement) {
	for i, j := 0, len(elements)-1; i < j; {
		elements[i], elements[j] = elements[j], elements[i]
		i++
		j--
	}
}

func StackTraceToSample(st *StackTrace, target *profile.Sample) {
	elem := make([]*StackTraceElement, len(st.Elements))
	copy(elem, st.Elements)

	reverse(elem)

	var loc []*profile.Location

	for _, element := range elem {
		loc = append(loc, element.Location)
	}

	target.Location = loc
}

func SampleToStackTrace(sample *profile.Sample) *StackTrace {
	var v []*StackTraceElement
	for _, location := range sample.Location {
		v = append(v, &StackTraceElement{location})
	}

	reverse(v)
	st := &StackTrace{v}
	return st
}
