package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

const Version = "1.0.1"

var (
	inputFile = kingpin.Flag("input", "input pprof file. \"-\" means read from the standard input stream").
			Short('i').
			Default(DefaultStream).
			String()
	outputFile = kingpin.Flag("output", "output pprof file, \"-\" means write to the standard output stream").
			Short('o').
			Default(DefaultStream).
			String()
	overlapFlag = kingpin.Flag("overlap", "the minimal overlapping call node count").
			Short('O').
			Default("5").
			Int()
	baseCount = kingpin.Flag("base", "number of the base nodes who always exist for all stacks, no matter whether the root call nodes are trimmed. (such as the process name for multi-process pprof)").
			Short('b').
			Default("0").
			Int()
	minDepth = kingpin.Flag("depth", "the minimal depth of the stack who may be trimmed (the deep stacks still remains deep after trimmed, base nodes not counted)").
			Short('d').
			Default("0").
			Int()
	verboseCounter = kingpin.Flag("verbose", "show verbose info on fixing the pprof").
			Short('v').
			Counter()
)

func InitFlags() {

	kingpin.Version(Version).
		Author("xnslong").
		Help = `A tool to fix the missing root call nodes of deep stacks in pprof result, so that the stacks can align with each other`
	kingpin.Parse()
}
