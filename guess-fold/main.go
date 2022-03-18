package main

import (
	"io"
	"log"
	"os"

	"github.com/xnslong/guess-stack/core"
	"github.com/xnslong/guess-stack/utils"
)

func main() {
	InitFlags()

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

func InputProfile(p *Profile) {
	err := utils.ReadFromFile(*inputFile, func(reader io.Reader) error {
		return p.ReadFrom(reader)
	})
	if err != nil {
		log.Printf("read folded stack error: %w", err)
		os.Exit(1)
	}
}

func OutputProfile(p *Profile) {
	err := utils.WriteToFile(*outputFile, func(writer io.Writer) error {
		return p.WriteTo(writer)
	})
	if err != nil {
		log.Printf("write folded stack error: %w", err)
		os.Exit(1)
	}
}
