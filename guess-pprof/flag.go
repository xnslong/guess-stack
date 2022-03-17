package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akamensky/argparse"
)

const Version = "v1.0.1-SNAPSHOT"

var (
	overlapCountThreshold *int
	outputFile            *string
	inputFile             *string
	version               *bool
	depth                 *int
	root                  *int
	verboseCounter        *int
)

func InitFlags() {
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
	version = parser.Flag("V", "version", &argparse.Options{
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
}
