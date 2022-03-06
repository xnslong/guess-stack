package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

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

	thisLine := l.Line[0]
	anotherLine := anotherLoc.Line[0]
	return thisLine.Line == anotherLine.Line && thisLine.Function.ID == anotherLine.Function.ID
	//&&
	//	l.Function.StartLine == anotherLoc.Function.StartLine
}

type StackTraceElementSlice []*StackTraceElement

func (s StackTraceElementSlice) Len() int {
	return len(s)
}

func (s StackTraceElementSlice) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

func (s StackTraceElementSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type StackTrace struct {
	elements []*StackTraceElement
}

func (s *StackTrace) String() string {
	return fmt.Sprintf("%s", s.elements)
}

func (s *StackTrace) Path() []fix.PathNode {
	pn := make([]fix.PathNode, 0, len(s.elements))

	for _, loc := range s.elements {
		pn = append(pn, loc)
	}

	return pn
}

func (s *StackTrace) SetPath(path []fix.PathNode) {
	v := s.elements[:0]
	for _, node := range path {
		if loc, ok := node.(*StackTraceElement); ok {
			v = append(v, loc)
		} else {
			log.Panicf("invalid type, expect *Loc, but got %T", node)
		}
	}
	s.elements = v
}

func FixProfile(p *profile.Profile) {
	var path []fix.Path

	for _, sample := range p.Sample {
		st := SampleToStackTrace(sample)

		path = append(path, st)
	}

	fixer.Fix(path)

	for j, sample := range p.Sample {
		st := path[j].(*StackTrace)

		if len(st.elements) != len(sample.Location) {
			beforeStack := SampleToStackTrace(sample)

			log.Printf("#%d before: %s", j, beforeStack.String())
			log.Printf("#%d after : %s", j, st.String())
		}

		RecoverStackToSample(st, sample)
	}

}

func reverse(elements []*StackTraceElement) {
	for i, j := 0, len(elements)-1; i < j; {
		elements[i], elements[j] = elements[j], elements[i]
		i++
		j--
	}
}

func RecoverStackToSample(st *StackTrace, sample *profile.Sample) {
	reverse(st.elements)

	var loc []*profile.Location

	for _, element := range st.elements {
		loc = append(loc, element.Location)
	}

	sample.Location = loc
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

var fixer fix.Fixer

var overlapCountThreshold = flag.Int("overlap", 5, "guess to be the same stack if has overlapped stack node more than the threshold")
var outputFile = flag.String("o", "-", "guess to be the same stack if has overlapped stack node more than the threshold")

func main() {
	flag.Parse()
	log.Printf("output to: %s", *outputFile)
	fixer = &fix.CommonRootFixer{CommonCount: *overlapCountThreshold}

	p, err := profile.Parse(os.Stdin)
	if err != nil {
		log.Panicf("open profile error: %v", err)
	}

	marshal, err := json.MarshalIndent(p.Sample, "", "    ")
	ioutil.WriteFile("before.json", marshal, 0644)

	FixProfile(p)

	marshal, err = json.MarshalIndent(p.Sample, "", "    ")
	ioutil.WriteFile("after.json", marshal, 0644)

	if *outputFile == "-" {
		err = p.Write(os.Stdout)
		if err != nil {
			log.Panicf("write profile error: %v", err)
		}
	} else {
		buf := &bytes.Buffer{}
		err = p.Write(buf)
		if err != nil {
			log.Panicf("write profile error: %v", err)
		}
		ioutil.WriteFile(*outputFile, buf.Bytes(), 0644)
	}
}
