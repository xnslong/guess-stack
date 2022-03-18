package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/google/pprof/profile"
	"github.com/xnslong/guess-stack/fix"
)

const DefaultStream = "-"

var fixer fix.StackFixer

func init() {
	InitFlags()
	fixer = &fix.CommonRootFixer{MinOverlaps: *overlapFlag}
}

func main() {
	p := OpenProfile()

	FixProfile(p)

	debugOutput(p)

	WriteProfile(p)
}

func debugOutput(p *profile.Profile) {
	if *verboseCounter < 2 {
		return
	}

	outFile := "outfix.json"
	if *outputFile != DefaultStream {
		outFile = path.Clean(*outputFile) + "_outfix.json"
	}
	marshal, _ := json.MarshalIndent(p.Sample, "", "    ")
	ioutil.WriteFile(outFile, marshal, 0644)
}

func FixProfile(p *profile.Profile) {
	var path []fix.Stack

	for _, sample := range p.Sample {
		st := SampleToStackTrace(sample, *baseCount)

		path = append(path, st)
	}

	toJoin := depthGreaterThanOrEqualTo(path, *minDepth)
	fixer.Fix(path, toJoin)

	for j, sample := range p.Sample {
		st := path[j].(*StackTrace)

		if len(st.Elements) != len(sample.Location) {
			if *verboseCounter >= 1 {
				beforeStack := SampleToStackTrace(sample, 0)
				log.Printf("#%d before: %s", j, beforeStack.String())
				log.Printf("#%d after : %s", j, st.String())
			}
		}

		StackTraceToSample(st, sample)
	}
}

func depthGreaterThanOrEqualTo(path []fix.Stack, threshold int) []bool {
	toJoin := make([]bool, 0, len(path))
	for _, stack := range path {
		if len(stack.Path()) >= threshold {
			toJoin = append(toJoin, true)
		} else {
			toJoin = append(toJoin, false)
		}
	}
	return toJoin
}

func WriteProfile(p *profile.Profile) {
	var out io.Writer
	if *outputFile == DefaultStream {
		out = os.Stdout
	} else {
		file, err := os.OpenFile(*outputFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Panic("write profile error", err)
		}
		defer file.Close()
		out = file
	}

	err := p.Write(out)
	if err != nil {
		log.Panic("write profile error", err)
	}
}

func OpenProfile() *profile.Profile {
	var in io.Reader
	if *inputFile == DefaultStream {
		in = os.Stdin
	} else {
		file, err := os.OpenFile(*inputFile, os.O_RDONLY, 0644)
		if err != nil {
			log.Panic("open profile error", err)
		}
		defer file.Close()
		in = file
	}

	p, err := profile.Parse(in)
	if err != nil {
		log.Panic("open profile error", err)
	}

	return p
}
