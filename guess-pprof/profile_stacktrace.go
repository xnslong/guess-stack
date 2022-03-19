package main

import (
	"fmt"
	"io"

	"github.com/google/pprof/profile"
	"github.com/xnslong/guess-stack/core/guess"
)

type StackTraceElement struct {
	*profile.Location
}

func (l *StackTraceElement) String() string {
	return fmt.Sprintf("%d", l.ID)
}

func (l *StackTraceElement) EqualsTo(another guess.StackNode) bool {
	anotherLoc, ok := another.(*StackTraceElement)
	if !ok {
		return false
	}

	return l.ID == anotherLoc.ID

}

type StackTrace struct {
	Elements []guess.StackNode
	*guess.StackExtraInfo
}

func (s *StackTrace) String() string {
	return fmt.Sprintf("%s", s.Elements)
}

func (s *StackTrace) Path() []guess.StackNode {
	return s.Elements
}

func (s *StackTrace) SetPath(path []guess.StackNode) {
	s.Elements = path
}

type Profile struct {
	PProf  *profile.Profile
	stacks []guess.Stack
}

func (p *Profile) Stacks() []guess.Stack {
	return p.stacks
}

func (p *Profile) WriteTo(writer io.Writer) error {
	for i, stack := range p.stacks {
		st := stack.(*StackTrace)
		StackTraceToSample(st, p.PProf.Sample[i])
	}

	return p.PProf.Write(writer)
}

func (p *Profile) ReadFrom(reader io.Reader) error {
	pprof, err := profile.Parse(reader)
	if err != nil {
		return fmt.Errorf("open profile error: %w", err)
	}

	stacks := make([]guess.Stack, 0, len(pprof.Sample))
	for _, sample := range pprof.Sample {
		stacks = append(stacks, SampleToStackTrace(sample))
	}

	p.PProf = pprof
	p.stacks = stacks

	return nil
}

func reverse(elements []guess.StackNode) {
	for i, j := 0, len(elements)-1; i < j; {
		elements[i], elements[j] = elements[j], elements[i]
		i++
		j--
	}
}

func StackTraceToSample(st *StackTrace, target *profile.Sample) {
	elem := make([]guess.StackNode, len(st.Elements))
	copy(elem, st.Elements)

	reverse(elem)

	var loc []*profile.Location

	for _, element := range elem {
		loc = append(loc, element.(*StackTraceElement).Location)
	}

	target.Location = loc
}

func SampleToStackTrace(sample *profile.Sample) *StackTrace {
	var v []guess.StackNode
	for _, location := range sample.Location {
		v = append(v, &StackTraceElement{location})
	}

	reverse(v)

	st := &StackTrace{
		Elements:       v,
		StackExtraInfo: guess.NewStackExtraInfo(),
	}
	return st
}
