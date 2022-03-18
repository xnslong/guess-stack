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
}

func OpenProfile() *Profile {
	p := &Profile{}

	err := utils.ReadFromFile(*inputFile, func(reader io.Reader) error {
		return p.ReadFrom(reader)
	})
	if err != nil {
		log.Panic("open profile error", err)
	}

	return p
}
