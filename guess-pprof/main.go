package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/akamensky/argparse"
	"github.com/google/pprof/profile"
	"github.com/xnslong/guess-stack/fix"
)

const DefaultStream = "-"

const Version = "v1.0.1-SNAPSHOT"

var fixer fix.StackFixer

var (
	overlapCountThreshold *int
	outputFile            *string
	inputFile             *string
	verbose               *bool
	version               *bool
	depth                 *int
	root                  *int
	verboseCounter        *int
)

func init() {
	parser := argparse.NewParser("guess-pprof", "to guess the missing root nodes for deep stacks, so that the stacks can align with each other again")
	inputFile = parser.String("i", "input", &argparse.Options{
		Required: false,
		Validate: nil,
		Help:     "input pprof file. \"-\" means read from the standard input stream",
		Default:  DefaultStream,
	})
	outputFile = parser.String("o", "output", &argparse.Options{
		Required: false,
		Validate: nil,
		Help:     "output pprof file, \"-\" means write to the standard output stream",
		Default:  DefaultStream,
	})
	overlapCountThreshold = parser.Int("", "overlap", &argparse.Options{
		Required: false,
		Validate: nil,
		Help:     "the minimal overlapping call node count",
		Default:  5,
	})
	depth = parser.Int("d", "depth", &argparse.Options{
		Required: false,
		Validate: nil,
		Help:     "the minimal depth of the stack who may be trimmed (the deep stacks still remains deep after trimmed)",
		Default:  1,
	})
	root = parser.Int("b", "base", &argparse.Options{
		Required: false,
		Validate: nil,
		Help:     "number of the base nodes always existing for all stacks (such as the process name), no matter whether the root call nodes are trimmed",
		Default:  0,
	})
	verboseCounter = parser.FlagCounter("v", "verbose", &argparse.Options{
		Required: false,
		Validate: nil,
		Help:     "",
		Default:  0,
	})
	version = parser.Flag("", "version", &argparse.Options{
		Required: false,
		Validate: nil,
		Help:     "",
		Default:  false,
	})
	parser.ExitOnHelp(true)
	err := parser.Parse(os.Args)
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}

	if *version {
		fmt.Println(os.Args[0], Version)
		os.Exit(0)
	}
	fixer = &fix.CommonRootFixer{MinOverlaps: *overlapCountThreshold}
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
		st := SampleToStackTrace(sample, *root)

		path = append(path, st)
	}

	toJoin := depthGreaterThanOrEqualTo(path, *depth)
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
