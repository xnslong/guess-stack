package utils

import (
	"os"
	"runtime/pprof"
)

func DoWithPProf(outfile string, job func()) {
	file, _ := os.OpenFile(outfile, os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()

	job()
}
