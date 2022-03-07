package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/google/pprof/profile"
	"github.com/xnslong/guess-stack/fix"
)

const DefaultStream = "-"

var fixer fix.Fixer

var (
	overlapCountThreshold = flag.Int("overlap", 5, "trustable overlap length. when the number of overlapping elements is less than the length, it's not considered trustable for guessing")
	outputFile            = flag.String("o", DefaultStream, "output file")
	inputFile             = flag.String("i", DefaultStream, "input file")
	verbose               = flag.Bool("v", false, "show verbose info for debug")
)

func init() {
	flag.Parse()
	fixer = &fix.CommonRootFixer{CommonCount: *overlapCountThreshold}
}

func main() {
	p, err := OpenProfile()
	if err != nil {
		log.Panicf("open profile error: %v", err)
	}

	//marshal, err := json.MarshalIndent(p.Sample, "", "    ")
	//ioutil.WriteFile("before.json", marshal, 0644)

	FixProfile(p)

	debugOutput(p)

	err = WriteProfile(p)
	if err != nil {
		log.Panicf("write profile error: %v", err)
	}

}

func debugOutput(p *profile.Profile) {
	if !*verbose {
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
	var path []fix.Path

	for _, sample := range p.Sample {
		st := SampleToStackTrace(sample)

		path = append(path, st)
	}

	fixer.Fix(path)

	for j, sample := range p.Sample {
		st := path[j].(*StackTrace)

		if len(st.Elements) != len(sample.Location) {
			if *verbose {
				beforeStack := SampleToStackTrace(sample)
				log.Printf("#%d before: %s", j, beforeStack.String())
				log.Printf("#%d after : %s", j, st.String())
			}
		}

		StackTraceToSample(st, sample)
	}
}

func WriteProfile(p *profile.Profile) error {

	var out io.Writer
	if *outputFile == DefaultStream {
		out = os.Stdout
	} else {
		outFile, err := os.OpenFile(*outputFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("open file error: (file=%s) %w", *outputFile, err)
		}

		defer outFile.Close()
		out = outFile
	}

	err := p.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func OpenProfile() (*profile.Profile, error) {
	var in io.Reader
	if *inputFile == DefaultStream {
		in = os.Stdin
	} else {
		inFile, err := os.Open(*inputFile)
		if err != nil {
			return nil, fmt.Errorf("open file error: (file=%s) %w", *inputFile, err)
		}

		defer inFile.Close()
		in = inFile
	}

	p, err := profile.Parse(in)
	if err != nil {
		return nil, fmt.Errorf("parse profile error: %w", err)
	}
	return p, nil
}