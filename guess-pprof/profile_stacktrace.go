package main

import (
	"fmt"
	"io"

	"github.com/google/pprof/profile"
	"github.com/xnslong/guess-stack/core/interfaces"
)

type StackTraceElement struct {
	*profile.Location
	name *nameStruct
}

type nameStruct struct {
	Value string
}

var names = make(map[string]*nameStruct)

func nameOf(name string) *nameStruct {
	val, ok := names[name]
	if ok {
		return val
	}

	val = &nameStruct{name}
	names[name] = val
	return val
}

func NewStackTraceElement(location *profile.Location) *StackTraceElement {
	var (
		name string
	)

	if len(location.Line) > 0 {
		f := location.Line[0].Function
		if f != nil {
			name = f.Name
		}
	}

	var n *nameStruct
	if len(name) != 0 {
		n = nameOf(name)
	}
	return &StackTraceElement{
		Location: location,
		name:     n,
	}
}

func (l *StackTraceElement) String() string {
	return fmt.Sprintf("%d", l.ID)
}

func (l *StackTraceElement) EqualsTo(another interfaces.StackNode) bool {
	anotherLoc, ok := another.(*StackTraceElement)
	if !ok {
		return false
	}

	if l.name != nil || anotherLoc.name != nil {
		return l.name == anotherLoc.name
	}

	return l.ID == anotherLoc.ID
}

type StackTrace struct {
	Elements []interfaces.StackNode
	*interfaces.StackExtraInfo
}

func (s *StackTrace) String() string {
	return fmt.Sprintf("%s", s.Elements)
}

func (s *StackTrace) Path() []interfaces.StackNode {
	return s.Elements
}

func (s *StackTrace) SetPath(path []interfaces.StackNode) {
	s.Elements = path
}

type Profile struct {
	PProf  *profile.Profile
	stacks []interfaces.Stack
}

func (p *Profile) Stacks() []interfaces.Stack {
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

	stacks := make([]interfaces.Stack, 0, len(pprof.Sample))
	for _, sample := range pprof.Sample {
		stacks = append(stacks, SampleToStackTrace(sample))
	}

	p.PProf = pprof
	p.stacks = stacks

	return nil
}

func reverse(elements []interfaces.StackNode) {
	for i, j := 0, len(elements)-1; i < j; {
		elements[i], elements[j] = elements[j], elements[i]
		i++
		j--
	}
}

func StackTraceToSample(st *StackTrace, target *profile.Sample) {
	elem := make([]interfaces.StackNode, len(st.Elements))
	copy(elem, st.Elements)

	reverse(elem)

	var loc []*profile.Location

	for _, element := range elem {
		loc = append(loc, element.(*StackTraceElement).Location)
	}

	target.Location = loc
}

func SampleToStackTrace(sample *profile.Sample) *StackTrace {
	var v []interfaces.StackNode
	for _, location := range sample.Location {
		v = append(v, NewStackTraceElement(location))
	}

	reverse(v)

	st := &StackTrace{
		Elements:       v,
		StackExtraInfo: interfaces.NewStackExtraInfo(),
	}
	return st
}
