package main

import (
	"github.com/xnslong/guess-stack/core"
	"github.com/xnslong/guess-stack/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	inputFile = kingpin.Flag("input", "input folded stack file. \"-\" means read from the standard input stream").
			Short('i').
			Default(utils.DefaultStream).
			String()
	outputFile = kingpin.Flag("output", "output folded stack file, \"-\" means write to the standard output stream").
			Short('o').
			Default(utils.DefaultStream).
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
	pprofOutfile = kingpin.Flag("self-pprof", "pprof for the program itself for improvement reasons").
			Short('p').
			Default("").
			String()
	verboseCounter = kingpin.Flag("verbose", "show verbose info on fixing the pprof").
			Short('v').
			Counter()
)

func InitFlags() {
	kingpin.Version(core.Version).
		Author(core.Author).
		Help = `A tool to fix the missing root call nodes of deep stacks in folded stacks, so that the stacks can align with each other. The output is also folded stacks.`
	kingpin.Parse()
}
