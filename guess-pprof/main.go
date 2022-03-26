package main

import (
	"io"
	"log"
	"os"
	"path"

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
	printVersion()

	p := OpenProfile()

	core.Fix(p, core.FixOption{
		Overlap:   *overlapFlag,
		BaseCount: *baseCount,
		MinDepth:  *minDepth,
		Verbose:   *verboseCounter,
	})

	WriteProfile(p)
}

func printVersion() {
	if *verboseCounter > 0 {
		log.Printf("program: %s @%s\n", path.Base(os.Args[0]), core.Version)
	}
}

func WriteProfile(p *Profile) {
	err := utils.WriteToFile(*outputFile, func(writer io.Writer) error {
		return p.WriteTo(writer)
	})
	if err != nil {
		log.Fatalf("write profile error: %v", err)
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
		log.Fatalf("read profile error: %v", err)
	}

	if *verboseCounter > 0 {
		log.Printf("read profile success from: %s", *inputFile)
	}

	return p
}
