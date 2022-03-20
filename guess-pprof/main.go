package main

import (
	"io"
	"log"

	"github.com/xnslong/guess-stack/core"
	"github.com/xnslong/guess-stack/utils"
)

func init() {
	InitFlags()
}

func main() {
	if *pprofOutfile != "" {
		utils.DoWithPProf(*pprofOutfile, FixProfile)
	} else {
		FixProfile()
	}
}

func FixProfile() {
	p := OpenProfile()

	core.Fix(p, core.FixOption{
		Overlap:   *overlapFlag,
		BaseCount: *baseCount,
		MinDepth:  *minDepth,
		Verbose:   *verboseCounter,
	})

	WriteProfile(p)
}

func WriteProfile(p *Profile) {
	err := utils.WriteToFile(*outputFile, func(writer io.Writer) error {
		return p.WriteTo(writer)
	})
	if err != nil {
		log.Panic("write profile error", err)
	}

	if *verboseCounter > 0 {
		log.Printf("write profile success to: %s", *outputFile)
	}
}

func OpenProfile() *Profile {
	p := &Profile{}

	err := utils.ReadFromFile(*inputFile, func(reader io.Reader) error {
		return p.ReadFrom(reader)
	})
	if err != nil {
		log.Panic("open profile error", err)
	}

	if *verboseCounter > 0 {
		log.Printf("read profile success from: %s", *inputFile)
	}

	return p
}
