package main

import (
	"io"
	"log"
	"os"
	"path"

	"github.com/xnslong/guess-stack/core"
	"github.com/xnslong/guess-stack/utils"
)

func main() {
	InitFlags()

	if *pprofOutfile != "" {
		utils.DoWithPProf(*pprofOutfile, FixProfile)
	} else {
		FixProfile()
	}
}

func FixProfile() {
	printVersion()

	p := &Profile{}
	InputProfile(p)

	core.Fix(p, core.FixOption{
		Overlap:   *overlapFlag,
		BaseCount: *baseCount,
		MinDepth:  *minDepth,
		Verbose:   *verboseCounter,
	})

	OutputProfile(p)
}

func printVersion() {
	if *verboseCounter > 0 {
		log.Printf("program: %s @%s", path.Base(os.Args[0]), core.Version)
	}
}

func InputProfile(p *Profile) {
	err := utils.ReadFromFile(*inputFile, func(reader io.Reader) error {
		return p.ReadFrom(reader)
	})
	if err != nil {
		log.Printf("read folded stack error: %v", err)
		os.Exit(1)
	}

	if *verboseCounter > 0 {
		log.Printf("read profile success from: %s", *inputFile)
	}
}

func OutputProfile(p *Profile) {
	err := utils.WriteToFile(*outputFile, func(writer io.Writer) error {
		return p.WriteTo(writer)
	})
	if err != nil {
		log.Printf("write folded stack error: %v", err)
		os.Exit(1)
	}
	if *verboseCounter > 0 {
		log.Printf("write profile success to: %s", *outputFile)
	}
}
